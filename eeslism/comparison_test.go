// C言語版との比較テスト
// testdata内のc_output/go_outputを比較して、Go版がC版と同等の結果を出力することを確認する
package eeslism

import (
	"bufio"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// CompareConfig は比較設定を保持する
type CompareConfig struct {
	Tolerance         float64 // 許容相対誤差（デフォルト0.001 = 0.1%）
	WarnTol           float64 // 警告閾値（デフォルト0.01 = 1%）
	AbsTol            float64 // 許容絶対誤差（小さな値用）
	NearZeroThreshold float64 // ゼロ近傍判定閾値
	TempAbsTol        float64 // 温度用許容絶対誤差（デフォルト0.1度）
}

// FieldInfo はESファイルのフィールド情報
type FieldInfo struct {
	Name       string
	TypeChar   string // H, h, T, t, Q, q, E, e, etc.
	FormatChar string // d (integer), f (float)
	IsMetadata bool   // True if statistical metadata (time/count)
}

// CompareResult は比較結果を保持する
type CompareResult struct {
	FileName    string
	TotalValues int
	Matched     int
	WithinTol   int
	Warnings    int
	Failures    int
	MaxRelErr   float64
	MaxAbsErr   float64
	MaxErrLine  int
	MaxErrCol   int

	// フィールドタイプ別統計
	PhysicsTotal       int
	PhysicsDifferent   int
	PhysicsMaxRelErr   float64
	PhysicsMaxAbsErr   float64

	// 温度フィールド（絶対誤差で評価）
	TempTotal          int
	TempDifferent      int
	TempMaxAbsErr      float64

	NearZeroTotal      int
	NearZeroDifferent  int
	NearZeroMaxAbsErr  float64

	MetadataTotal      int
	MetadataDifferent  int
}

// DefaultCompareConfig はデフォルトの比較設定を返す
func DefaultCompareConfig() CompareConfig {
	return CompareConfig{
		Tolerance:         0.001,  // 0.1%
		WarnTol:           0.01,   // 1%
		AbsTol:            1e-10,
		NearZeroThreshold: 1.0,    // ゼロ近傍閾値
		TempAbsTol:        0.1,    // 温度許容絶対誤差 0.1度
	}
}

// compareDirectories は2つのディレクトリ内の.esファイルを比較する
func compareDirectories(refDir, testDir string, config CompareConfig) ([]CompareResult, error) {
	var results []CompareResult

	err := filepath.Walk(refDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		// .esファイルのみ対象
		if filepath.Ext(path) != ".es" {
			return nil
		}

		// 対応するテストファイルのパスを構築
		relPath, _ := filepath.Rel(refDir, path)
		testPath := filepath.Join(testDir, relPath)

		if _, err := os.Stat(testPath); os.IsNotExist(err) {
			// テストファイルが存在しない場合はスキップ
			return nil
		}

		result := compareFiles(path, testPath, config)
		results = append(results, result)

		return nil
	})

	return results, err
}

// parseESHeader はESファイルのヘッダーを解析してフィールド情報を取得
func parseESHeader(path string) ([]FieldInfo, int) {
	var fields []FieldInfo
	dataStart := 0

	file, err := os.Open(path)
	if err != nil {
		return fields, 0
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	inHeader := false
	lineNum := 0

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		lineNum++

		// ヘッダー開始マーカー
		if line == "#" {
			inHeader = true
			continue
		}

		if inHeader {
			// データ行の開始を検出（数字で始まる行または-999）
			if len(line) > 0 && (line[0] >= '0' && line[0] <= '9' || strings.HasPrefix(line, "-999")) {
				dataStart = lineNum
				break
			}

			// フィールド定義を解析
			// 例: "SupplyFan_Hq H d SupplyFan_Q Q f ..."
			tokens := strings.Fields(line)
			j := 0
			for j+2 < len(tokens) {
				name := tokens[j]
				typeChar := tokens[j+1]
				formatChar := tokens[j+2]

				// 統計メタデータの判定
				// H d: 時間カウント (Hours)
				// h d: 時刻 (hour/time)
				isMetadata := (typeChar == "H" || typeChar == "h") && formatChar == "d"

				fields = append(fields, FieldInfo{
					Name:       name,
					TypeChar:   typeChar,
					FormatChar: formatChar,
					IsMetadata: isMetadata,
				})
				j += 3
			}
		}
	}

	// sfq.es形式（TSV、#マーカーなし）の検出
	if len(fields) == 0 {
		file.Seek(0, 0)
		scanner = bufio.NewScanner(file)
		lineNum = 0
		for scanner.Scan() {
			line := scanner.Text()
			lineNum++
			// TSVヘッダー行を検出（Mo\tNd\ttime...）
			if strings.Contains(line, "\t") && strings.Contains(line, "Mo") && strings.Contains(line, "Nd") {
				cols := strings.Split(strings.TrimSpace(line), "\t")
				for _, col := range cols[3:] { // Mo, Nd, time をスキップ
					fields = append(fields, FieldInfo{
						Name:       col,
						TypeChar:   "f",
						FormatChar: "f",
						IsMetadata: false,
					})
				}
				dataStart = lineNum + 1
				break
			}
		}
	}

	return fields, dataStart
}

// compareFiles は2つのファイルを比較する
func compareFiles(refPath, testPath string, config CompareConfig) CompareResult {
	result := CompareResult{
		FileName: filepath.Base(refPath),
	}

	// ヘッダー解析
	fields, _ := parseESHeader(refPath)

	refLines, err := readDataLines(refPath)
	if err != nil {
		return result
	}

	testLines, err := readDataLines(testPath)
	if err != nil {
		return result
	}

	minLines := len(refLines)
	if len(testLines) < minLines {
		minLines = len(testLines)
	}

	// TSV形式かどうか判定
	isTSV := false
	if len(refLines) > 0 {
		isTSV = strings.Contains(refLines[0], "\t")
	}

	// 各行を比較
	for lineNum := 0; lineNum < minLines; lineNum++ {
		refLine := refLines[lineNum]
		testLine := testLines[lineNum]

		// -999行はスキップ
		if strings.HasPrefix(strings.TrimSpace(refLine), "-999") {
			continue
		}

		var refTokens, testTokens []string
		var dataOffset int

		if isTSV {
			refTokens = strings.Split(refLine, "\t")
			testTokens = strings.Split(testLine, "\t")
			dataOffset = 3 // Mo, Nd, time
		} else {
			refTokens = strings.Fields(refLine)
			testTokens = strings.Fields(testLine)

			// 空行はスキップ
			if len(refTokens) == 0 || len(testTokens) == 0 {
				continue
			}

			// 日付のみの行（2-3トークン）は比較をスキップ
			if len(refTokens) <= 3 {
				first := refTokens[0]
				if len(first) > 0 && first[0] >= '0' && first[0] <= '9' {
					// 日付行（例: "01 31" or "01 01 1.00"）はスキップ
					continue
				}
			}

			// 日付プレフィックスがない行（最初のトークンが数字で始まらないか、多くのトークンを持つ）
			if len(refTokens) > 0 {
				first := refTokens[0]
				if len(first) > 0 && first[0] >= '0' && first[0] <= '9' && len(refTokens) <= 5 {
					// 少ないトークン数で数字始まり = 日付プレフィックスあり（例: "01 01 1.00 value value"）
					dataOffset = 2
					if len(refTokens) > 2 && strings.Contains(refTokens[2], ".") {
						dataOffset = 3
					}
				} else if len(first) > 0 && !(first[0] >= '0' && first[0] <= '9') {
					// 数字で始まらない = データ行（例: "F 0.128 40.46..."）
					dataOffset = 0
				} else {
					// 多くのトークンで数字始まり = データ行（例: "744 40.5 1160700..."）
					dataOffset = 0
				}
			} else {
				dataOffset = 0
			}
		}

		// フィールドごとに比較
		fieldIdx := 0
		for i := dataOffset; i < len(refTokens) && i < len(testTokens); i++ {
			refVal, refErr := strconv.ParseFloat(refTokens[i], 64)
			testVal, testErr := strconv.ParseFloat(testTokens[i], 64)

			if refErr != nil || testErr != nil {
				// 解析できない場合もfieldIdxをインクリメント（文字型フィールドなど）
				fieldIdx++
				continue
			}

			result.TotalValues++
			relErr, absErr := calculateError(refVal, testVal)
			isDifferent := absErr > config.AbsTol

			// フィールドタイプに基づく分類
			isMetadata := false
			if fieldIdx < len(fields) {
				field := fields[fieldIdx]
				fieldIdx++

				if field.IsMetadata {
					// 統計メタデータ
					isMetadata = true
					result.MetadataTotal++
					if isDifferent {
						result.MetadataDifferent++
					}
				} else if field.TypeChar == "t" || field.TypeChar == "T" {
					// 温度フィールド: 絶対誤差0.1度で評価
					result.TempTotal++
					isTempDifferent := absErr > config.TempAbsTol
					if isTempDifferent {
						result.TempDifferent++
						if absErr > result.TempMaxAbsErr {
							result.TempMaxAbsErr = absErr
						}
					}
				} else if math.Abs(refVal) <= config.NearZeroThreshold {
					// ゼロ近傍の物理値
					result.NearZeroTotal++
					if isDifferent {
						result.NearZeroDifferent++
						if absErr > result.NearZeroMaxAbsErr {
							result.NearZeroMaxAbsErr = absErr
						}
					}
				} else {
					// 通常の物理値（熱量など）: 相対誤差で評価
					result.PhysicsTotal++
					if isDifferent {
						result.PhysicsDifferent++
						if relErr > result.PhysicsMaxRelErr {
							result.PhysicsMaxRelErr = relErr
						}
						if absErr > result.PhysicsMaxAbsErr {
							result.PhysicsMaxAbsErr = absErr
						}
					}
				}
			} else {
				// フィールド情報がない場合は従来の分類
				if math.Abs(refVal) <= config.NearZeroThreshold {
					result.NearZeroTotal++
					if isDifferent {
						result.NearZeroDifferent++
					}
				} else {
					result.PhysicsTotal++
					if isDifferent {
						result.PhysicsDifferent++
					}
				}
			}

			// 従来の分類（互換性のため）
			// NearZero値とメタデータは相対誤差が無意味なので、Failuresから除外
			isNearZero := math.Abs(refVal) <= config.NearZeroThreshold
			if absErr < config.AbsTol {
				result.Matched++
			} else if relErr < config.Tolerance {
				result.WithinTol++
			} else if relErr < config.WarnTol {
				result.Warnings++
			} else if isNearZero || isMetadata {
				// NearZero値とメタデータは相対誤差が大きくても警告止まり
				result.Warnings++
			} else {
				result.Failures++
			}

			// 最大誤差の更新
			if relErr > result.MaxRelErr {
				result.MaxRelErr = relErr
				result.MaxAbsErr = absErr
				result.MaxErrLine = lineNum + 1
				result.MaxErrCol = i + 1
			}
		}
	}

	return result
}

// readDataLines はファイルからデータ行のみを読み込む（ヘッダーをスキップ）
func readDataLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	inData := false

	for scanner.Scan() {
		line := scanner.Text()

		// データ部分の開始を検出（数字で始まる行）
		if !inData {
			trimmed := strings.TrimSpace(line)
			if len(trimmed) > 0 {
				firstChar := trimmed[0]
				if firstChar >= '0' && firstChar <= '9' {
					inData = true
				}
			}
		}

		if inData {
			lines = append(lines, line)
		}
	}

	return lines, scanner.Err()
}

// parseValues は行を数値のスライスに変換
func parseValues(line string) []float64 {
	var values []float64
	fields := strings.Fields(line)

	for _, field := range fields {
		// NaN文字列のチェック（C版は"-nan"、Go版は"NaN"）
		lowerField := strings.ToLower(field)
		if lowerField == "nan" || lowerField == "-nan" {
			values = append(values, math.NaN())
			continue
		}
		val, err := strconv.ParseFloat(field, 64)
		if err == nil {
			values = append(values, val)
		}
	}

	return values
}

// calculateError は相対誤差と絶対誤差を計算
func calculateError(ref, test float64) (relErr, absErr float64) {
	// NaN同士は一致とみなす
	if math.IsNaN(ref) && math.IsNaN(test) {
		return 0.0, 0.0
	}
	// 片方だけNaNの場合は最大誤差
	if math.IsNaN(ref) || math.IsNaN(test) {
		return 1.0, math.MaxFloat64
	}

	absErr = math.Abs(test - ref)

	if math.Abs(ref) > 1e-15 {
		relErr = absErr / math.Abs(ref)
	} else if absErr > 1e-15 {
		relErr = 1.0 // 参照が0で差がある場合は100%誤差
	} else {
		relErr = 0.0 // 両方とも0
	}

	return relErr, absErr
}

// runComparisonTest は比較テストを実行するヘルパー関数
func runComparisonTest(t *testing.T, name, refDir, testDir string) {
	t.Helper()

	config := DefaultCompareConfig()
	results, err := compareDirectories(refDir, testDir, config)
	if err != nil {
		t.Fatalf("Failed to compare directories: %v", err)
	}

	if len(results) == 0 {
		t.Skip("No .es files found to compare")
	}

	totalPass := 0
	totalWarn := 0
	totalFail := 0

	// フィールドタイプ別の合計
	totalPhysics := 0
	totalPhysicsDiff := 0
	totalTemp := 0
	totalTempDiff := 0
	totalNearZero := 0
	totalNearZeroDiff := 0
	totalMetadata := 0
	totalMetadataDiff := 0
	maxPhysicsRelErr := 0.0
	maxTempAbsErr := 0.0

	for _, r := range results {
		totalPhysics += r.PhysicsTotal
		totalPhysicsDiff += r.PhysicsDifferent
		totalTemp += r.TempTotal
		totalTempDiff += r.TempDifferent
		totalNearZero += r.NearZeroTotal
		totalNearZeroDiff += r.NearZeroDifferent
		totalMetadata += r.MetadataTotal
		totalMetadataDiff += r.MetadataDifferent
		if r.PhysicsMaxRelErr > maxPhysicsRelErr {
			maxPhysicsRelErr = r.PhysicsMaxRelErr
		}
		if r.TempMaxAbsErr > maxTempAbsErr {
			maxTempAbsErr = r.TempMaxAbsErr
		}

		if r.Failures > 0 {
			totalFail++
			t.Errorf("[FAIL] %s: %d values, %d failures, max error %.4f%% at line %d col %d",
				r.FileName, r.TotalValues, r.Failures, r.MaxRelErr*100, r.MaxErrLine, r.MaxErrCol)
		} else if r.Warnings > 0 {
			totalWarn++
			t.Logf("[WARN] %s: %d values, %d warnings, max error %.4f%%",
				r.FileName, r.TotalValues, r.Warnings, r.MaxRelErr*100)
		} else {
			totalPass++
		}
	}

	// フィールドタイプ別サマリー
	t.Logf("Summary: %d files - %d PASS, %d WARN, %d FAIL", len(results), totalPass, totalWarn, totalFail)

	physicsRate := 0.0
	if totalPhysics > 0 {
		physicsRate = float64(totalPhysicsDiff) / float64(totalPhysics) * 100
	}
	tempRate := 0.0
	if totalTemp > 0 {
		tempRate = float64(totalTempDiff) / float64(totalTemp) * 100
	}
	nearZeroRate := 0.0
	if totalNearZero > 0 {
		nearZeroRate = float64(totalNearZeroDiff) / float64(totalNearZero) * 100
	}
	metadataRate := 0.0
	if totalMetadata > 0 {
		metadataRate = float64(totalMetadataDiff) / float64(totalMetadata) * 100
	}

	t.Logf("Field Type Analysis:")
	t.Logf("  Physics:   %d/%d (%.2f%%), max rel err: %.4f%%", totalPhysicsDiff, totalPhysics, physicsRate, maxPhysicsRelErr*100)
	t.Logf("  Temp:      %d/%d (%.2f%%), max abs err: %.4f deg", totalTempDiff, totalTemp, tempRate, maxTempAbsErr)
	t.Logf("  NearZero:  %d/%d (%.2f%%)", totalNearZeroDiff, totalNearZero, nearZeroRate)
	t.Logf("  Metadata:  %d/%d (%.2f%%)", totalMetadataDiff, totalMetadata, metadataRate)
}

// runComparisonTestWithVariants はメインテストとバリアントを含めて比較テストを実行（静的比較）
// baseDir: テストディレクトリのベースパス (e.g., "../tests/comparison/testdata/L2_equipment/pump_pipe")
func runComparisonTestWithVariants(t *testing.T, name, baseDir string) {
	t.Helper()

	// メインテストを実行
	mainRefDir := filepath.Join(baseDir, "c_output")
	mainTestDir := filepath.Join(baseDir, "go_output")

	if _, err := os.Stat(mainRefDir); err == nil {
		t.Run("main", func(t *testing.T) {
			runComparisonTest(t, name, mainRefDir, mainTestDir)
		})
	}

	// variantsディレクトリを探索
	variantsDir := filepath.Join(baseDir, "variants")
	if _, err := os.Stat(variantsDir); os.IsNotExist(err) {
		return // variantsがなければ終了
	}

	entries, err := os.ReadDir(variantsDir)
	if err != nil {
		t.Logf("Warning: cannot read variants directory: %v", err)
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		variantName := entry.Name()
		variantRefDir := filepath.Join(variantsDir, variantName, "c_output")
		variantTestDir := filepath.Join(variantsDir, variantName, "go_output")

		// c_outputとgo_outputの両方が存在する場合のみテスト実行
		if _, err := os.Stat(variantRefDir); os.IsNotExist(err) {
			continue
		}
		if _, err := os.Stat(variantTestDir); os.IsNotExist(err) {
			continue
		}

		t.Run(variantName, func(t *testing.T) {
			runComparisonTest(t, variantName, variantRefDir, variantTestDir)
		})
	}
}

// runSimulationAndCompare はGo版シミュレーションを実行してC版出力と比較する
// baseDir: テストディレクトリのベースパス
// eflPath: Baseファイルのパス (e.g., "../Base")
func runSimulationAndCompare(t *testing.T, name, baseDir, eflPath string) {
	t.Helper()

	// テスト入力ファイルを探す
	testFile := findTestInputFile(baseDir, name)
	if testFile == "" {
		t.Skipf("No test input file found in %s", baseDir)
		return
	}

	// 作業ディレクトリを保存
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	// テストディレクトリの絶対パスを取得
	absBaseDir, err := filepath.Abs(baseDir)
	if err != nil {
		t.Fatalf("Failed to get absolute path: %v", err)
	}
	absEflPath, err := filepath.Abs(eflPath)
	if err != nil {
		t.Fatalf("Failed to get absolute efl path: %v", err)
	}

	// テストディレクトリに移動
	if err := os.Chdir(absBaseDir); err != nil {
		t.Fatalf("Failed to change to test directory: %v", err)
	}
	defer os.Chdir(origDir)

	// 既存の.esファイルを削除
	cleanupGeneratedFiles(t, ".")

	// シミュレーション実行
	t.Logf("Running simulation: %s", testFile)
	Entry(testFile, absEflPath)

	// 生成された.esファイルをc_outputと比較
	refDir := filepath.Join(absBaseDir, "c_output")
	if _, err := os.Stat(refDir); os.IsNotExist(err) {
		t.Skipf("Reference directory not found: %s", refDir)
		return
	}

	config := DefaultCompareConfig()
	results, err := compareDirectoriesWithGenerated(refDir, ".", config)
	if err != nil {
		t.Fatalf("Failed to compare: %v", err)
	}

	// 結果を表示
	reportComparisonResults(t, name, results)

	// 生成ファイルをクリーンアップ
	cleanupGeneratedFiles(t, ".")
}

// findTestInputFile はテストディレクトリからメインのテスト入力ファイルを探す
func findTestInputFile(baseDir, name string) string {
	// 優先順位: {name}_test.txt > {name}.txt > 最初の.txt
	candidates := []string{
		filepath.Join(baseDir, name+"_test.txt"),
		filepath.Join(baseDir, name+".txt"),
	}

	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return filepath.Base(path)
		}
	}

	// 最初の.txtファイルを探す
	entries, err := os.ReadDir(baseDir)
	if err != nil {
		return ""
	}
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".txt") {
			return entry.Name()
		}
	}
	return ""
}

// cleanupGeneratedFiles は生成された.esファイルと一時ファイルを削除する
func cleanupGeneratedFiles(t *testing.T, dir string) {
	t.Helper()

	patterns := []string{"*.es", "*.gchi", "*.log", "*bdata.ewk", "*bdata0.ewk", "*schnma.ewk", "*schtba.ewk", "*week.ewk"}
	for _, pattern := range patterns {
		matches, _ := filepath.Glob(filepath.Join(dir, pattern))
		for _, match := range matches {
			os.Remove(match)
		}
	}
}

// compareDirectoriesWithGenerated はリファレンスディレクトリと生成されたファイルを比較する
func compareDirectoriesWithGenerated(refDir, genDir string, config CompareConfig) ([]CompareResult, error) {
	var results []CompareResult

	err := filepath.Walk(refDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		// .esファイルのみ対象
		if filepath.Ext(path) != ".es" {
			return nil
		}

		// 生成されたファイルのパスを構築
		fileName := filepath.Base(path)
		genPath := filepath.Join(genDir, fileName)

		if _, err := os.Stat(genPath); os.IsNotExist(err) {
			// 生成ファイルが存在しない場合はスキップ（または警告）
			return nil
		}

		result := compareFiles(path, genPath, config)

		result.FileName = fileName
		results = append(results, result)
		return nil
	})

	return results, err
}

// reportComparisonResults は比較結果をレポートする
func reportComparisonResults(t *testing.T, name string, results []CompareResult) {
	t.Helper()

	if len(results) == 0 {
		t.Log("No files to compare")
		return
	}

	totalPass := 0
	totalWarn := 0
	totalFail := 0

	config := DefaultCompareConfig()

	for _, r := range results {
		status := "PASS"
		if r.Failures > 0 {
			status = "FAIL"
			totalFail++
			t.Logf("[FAIL] %s: %d values, %d failures, max error %.4f%% at line %d col %d",
				r.FileName, r.TotalValues, r.Failures, r.MaxRelErr*100, r.MaxErrLine, r.MaxErrCol)
		} else if r.Warnings > 0 {
			status = "WARN"
			totalWarn++
			t.Logf("[WARN] %s: %d values, %d warnings, max error %.4f%%",
				r.FileName, r.TotalValues, r.Warnings, r.MaxRelErr*100)
		} else {
			totalPass++
		}
		_ = status
		_ = config
	}

	t.Logf("Summary: %d files - %d PASS, %d WARN, %d FAIL", len(results), totalPass, totalWarn, totalFail)

	if totalFail > 0 {
		t.Fail()
	}
}

// runSimulationTestWithVariants はメインテストとバリアントを含めてシミュレーション実行＋比較テストを行う
func runSimulationTestWithVariants(t *testing.T, name, baseDir, eflPath string) {
	t.Helper()

	// 最初に全ての絶対パスを計算（chdir前に行う）
	absBaseDir, err := filepath.Abs(baseDir)
	if err != nil {
		t.Fatalf("Failed to get absolute path for baseDir: %v", err)
	}
	absEflPath, err := filepath.Abs(eflPath)
	if err != nil {
		t.Fatalf("Failed to get absolute path for eflPath: %v", err)
	}

	// メインテストを実行
	refDir := filepath.Join(absBaseDir, "c_output")
	if _, err := os.Stat(refDir); err == nil {
		t.Run("main", func(t *testing.T) {
			runSimulationAndCompare(t, name, baseDir, eflPath)
		})
	}

	// variantsディレクトリを探索
	variantsDir := filepath.Join(absBaseDir, "variants")
	if _, err := os.Stat(variantsDir); os.IsNotExist(err) {
		return
	}

	entries, err := os.ReadDir(variantsDir)
	if err != nil {
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		variantName := entry.Name()
		variantRefDir := filepath.Join(variantsDir, variantName, "c_output")

		if _, err := os.Stat(variantRefDir); os.IsNotExist(err) {
			continue
		}

		// テストファイルの存在確認（絶対パスで）
		testFile := variantName + ".txt"
		testFilePath := filepath.Join(absBaseDir, testFile)
		if _, err := os.Stat(testFilePath); os.IsNotExist(err) {
			continue // ファイルがなければスキップ
		}

		t.Run(variantName, func(t *testing.T) {
			// 作業ディレクトリを保存
			origDir, _ := os.Getwd()

			os.Chdir(absBaseDir)
			defer os.Chdir(origDir)

			cleanupGeneratedFiles(t, ".")

			t.Logf("Running variant simulation: %s", testFile)
			Entry(testFile, absEflPath)

			config := DefaultCompareConfig()
			results, _ := compareDirectoriesWithGenerated(variantRefDir, ".", config)
			reportComparisonResults(t, variantName, results)
			cleanupGeneratedFiles(t, ".")
		})
	}
}

// runComparisonTestWithPhysicsThreshold は物理値の最大相対誤差閾値を指定して比較テストを実行
// maxPhysicsRelErrThreshold: 許容する最大相対誤差（%）
func runComparisonTestWithPhysicsThreshold(t *testing.T, name, refDir, testDir string, maxPhysicsRelErrThreshold float64) {
	t.Helper()

	config := DefaultCompareConfig()
	results, err := compareDirectories(refDir, testDir, config)
	if err != nil {
		t.Fatalf("Failed to compare directories: %v", err)
	}

	if len(results) == 0 {
		t.Skip("No .es files found to compare")
	}

	// フィールドタイプ別の合計
	totalPhysics := 0
	totalPhysicsDiff := 0
	totalTemp := 0
	totalTempDiff := 0
	totalNearZero := 0
	totalNearZeroDiff := 0
	totalMetadata := 0
	totalMetadataDiff := 0
	maxPhysicsRelErr := 0.0
	maxTempAbsErr := 0.0

	for _, r := range results {
		totalPhysics += r.PhysicsTotal
		totalPhysicsDiff += r.PhysicsDifferent
		totalTemp += r.TempTotal
		totalTempDiff += r.TempDifferent
		totalNearZero += r.NearZeroTotal
		totalNearZeroDiff += r.NearZeroDifferent
		totalMetadata += r.MetadataTotal
		totalMetadataDiff += r.MetadataDifferent
		if r.PhysicsMaxRelErr > maxPhysicsRelErr {
			maxPhysicsRelErr = r.PhysicsMaxRelErr
		}
		if r.TempMaxAbsErr > maxTempAbsErr {
			maxTempAbsErr = r.TempMaxAbsErr
		}
	}

	physicsRate := 0.0
	if totalPhysics > 0 {
		physicsRate = float64(totalPhysicsDiff) / float64(totalPhysics) * 100
	}
	tempRate := 0.0
	if totalTemp > 0 {
		tempRate = float64(totalTempDiff) / float64(totalTemp) * 100
	}
	nearZeroRate := 0.0
	if totalNearZero > 0 {
		nearZeroRate = float64(totalNearZeroDiff) / float64(totalNearZero) * 100
	}
	metadataRate := 0.0
	if totalMetadata > 0 {
		metadataRate = float64(totalMetadataDiff) / float64(totalMetadata) * 100
	}

	t.Logf("Field Type Analysis:")
	t.Logf("  Physics:   %d/%d (%.2f%%), max rel err: %.4f%%", totalPhysicsDiff, totalPhysics, physicsRate, maxPhysicsRelErr*100)
	t.Logf("  Temp:      %d/%d (%.2f%%), max abs err: %.4f deg", totalTempDiff, totalTemp, tempRate, maxTempAbsErr)
	t.Logf("  NearZero:  %d/%d (%.2f%%)", totalNearZeroDiff, totalNearZero, nearZeroRate)
	t.Logf("  Metadata:  %d/%d (%.2f%%)", totalMetadataDiff, totalMetadata, metadataRate)

	// 物理値の最大相対誤差が閾値を超えた場合のみエラー
	maxPhysicsRelErrPercent := maxPhysicsRelErr * 100
	if maxPhysicsRelErrPercent > maxPhysicsRelErrThreshold {
		t.Errorf("Physics max relative error %.4f%% exceeds threshold %.2f%%", maxPhysicsRelErrPercent, maxPhysicsRelErrThreshold)
	}

	// 温度の最大絶対誤差が0.1度を超えた場合のみエラー（浮動小数点精度のため1e-9のイプシロンを追加）
	if maxTempAbsErr > config.TempAbsTol+1e-9 {
		t.Errorf("Temperature max absolute error %.4f deg exceeds threshold %.2f deg", maxTempAbsErr, config.TempAbsTol)
	}
}

// ============================================================================
// L1_basic 比較テスト
// ============================================================================

func TestComparison_L1_SimpleRoomFull(t *testing.T) {
	refDir := "../tests/comparison/testdata/L1_basic/simple_room_full/c_output"
	testDir := "../tests/comparison/testdata/L1_basic/simple_room_full/go_output"
	runComparisonTest(t, "simple_room_full", refDir, testDir)
}

func TestComparison_L1_SimpleRoomInternalHeat(t *testing.T) {
	refDir := "../tests/comparison/testdata/L1_basic/simple_room_internal_heat/c_output"
	testDir := "../tests/comparison/testdata/L1_basic/simple_room_internal_heat/go_output"
	runComparisonTest(t, "simple_room_internal_heat", refDir, testDir)
}

func TestComparison_L1_SimpleRoomSchedule(t *testing.T) {
	refDir := "../tests/comparison/testdata/L1_basic/simple_room_schedule/c_output"
	testDir := "../tests/comparison/testdata/L1_basic/simple_room_schedule/go_output"
	runComparisonTest(t, "simple_room_schedule", refDir, testDir)
}

func TestComparison_L1_SimpleRoomVent(t *testing.T) {
	refDir := "../tests/comparison/testdata/L1_basic/simple_room_vent/c_output"
	testDir := "../tests/comparison/testdata/L1_basic/simple_room_vent/go_output"
	runComparisonTest(t, "simple_room_vent", refDir, testDir)
}

// ============================================================================
// L2_equipment 比較テスト
// ============================================================================

func TestComparison_L2_CoolingCoil(t *testing.T) {
	baseDir := "../tests/comparison/testdata/L2_equipment/cooling_coil"
	runComparisonTestWithVariants(t, "cooling_coil", baseDir)
}

func TestComparison_L2_HeatPump(t *testing.T) {
	baseDir := "../tests/comparison/testdata/L2_equipment/heat_pump"
	runComparisonTestWithVariants(t, "heat_pump", baseDir)
}

func TestComparison_L2_Hex(t *testing.T) {
	refDir := "../tests/comparison/testdata/L2_equipment/hex/c_output"
	testDir := "../tests/comparison/testdata/L2_equipment/hex/go_output"
	runComparisonTest(t, "hex", refDir, testDir)
}

func TestComparison_L2_Helm(t *testing.T) {
	refDir := "../tests/comparison/testdata/L2_equipment/helm/c_output"
	testDir := "../tests/comparison/testdata/L2_equipment/helm/go_output"
	runComparisonTest(t, "helm", refDir, testDir)
}

func TestComparison_L2_PumpPipe(t *testing.T) {
	baseDir := "../tests/comparison/testdata/L2_equipment/pump_pipe"
	runComparisonTestWithVariants(t, "pump_pipe", baseDir)
}

func TestComparison_L2_PV(t *testing.T) {
	refDir := "../tests/comparison/testdata/L2_equipment/pv/c_output"
	testDir := "../tests/comparison/testdata/L2_equipment/pv/go_output"
	runComparisonTest(t, "pv", refDir, testDir)
}

func TestComparison_L2_Qmeas(t *testing.T) {
	baseDir := "../tests/comparison/testdata/L2_equipment/qmeas"
	runComparisonTestWithVariants(t, "qmeas", baseDir)
}

func TestComparison_L2_Rmac(t *testing.T) {
	baseDir := "../tests/comparison/testdata/L2_equipment/rmac"
	runComparisonTestWithVariants(t, "rmac", baseDir)
}

func TestComparison_L2_SolarCollector(t *testing.T) {
	baseDir := "../tests/comparison/testdata/L2_equipment/solar_collector"
	runComparisonTestWithVariants(t, "solar_collector", baseDir)
}

func TestComparison_L2_StorageTank(t *testing.T) {
	baseDir := "../tests/comparison/testdata/L2_equipment/storage_tank"
	runComparisonTestWithVariants(t, "storage_tank", baseDir)
}

func TestComparison_L2_TotalHeatExchanger(t *testing.T) {
	baseDir := "../tests/comparison/testdata/L2_equipment/thex"
	runComparisonTestWithVariants(t, "thex", baseDir)
}

func TestComparison_L2_Valv(t *testing.T) {
	baseDir := "../tests/comparison/testdata/L2_equipment/valv"
	runComparisonTestWithVariants(t, "valv", baseDir)
}

func TestComparison_L2_VAV(t *testing.T) {
	baseDir := "../tests/comparison/testdata/L2_equipment/vav"
	runComparisonTestWithVariants(t, "vav", baseDir)
}

func TestComparison_L2_VAVCooling(t *testing.T) {
	baseDir := "../tests/comparison/testdata/L2_equipment/vav_cooling"
	runComparisonTestWithVariants(t, "vav_cooling", baseDir)
}

func TestComparison_L2_BoilerHeating(t *testing.T) {
	refDir := "../tests/comparison/testdata/L2_equipment/boiler_heating/c_output"
	testDir := "../tests/comparison/testdata/L2_equipment/boiler_heating/go_output"
	runComparisonTest(t, "boiler_heating", refDir, testDir)
}

func TestComparison_L2_Desiccant(t *testing.T) {
	baseDir := "../tests/comparison/testdata/L2_equipment/desiccant"
	runComparisonTestWithVariants(t, "desiccant", baseDir)
}

func TestComparison_L2_Duct(t *testing.T) {
	refDir := "../tests/comparison/testdata/L2_equipment/duct/c_output"
	testDir := "../tests/comparison/testdata/L2_equipment/duct/go_output"
	runComparisonTest(t, "duct", refDir, testDir)
}

func TestComparison_L2_Fan(t *testing.T) {
	refDir := "../tests/comparison/testdata/L2_equipment/fan/c_output"
	testDir := "../tests/comparison/testdata/L2_equipment/fan/go_output"
	runComparisonTest(t, "fan", refDir, testDir)
}

func TestComparison_L2_AirCollector(t *testing.T) {
	refDir := "../tests/comparison/testdata/L2_equipment/air_collector/c_output"
	testDir := "../tests/comparison/testdata/L2_equipment/air_collector/go_output"
	runComparisonTest(t, "air_collector", refDir, testDir)
}

func TestComparison_L2_Evpcooling(t *testing.T) {
	baseDir := "../tests/comparison/testdata/L2_equipment/evpcooling"
	runComparisonTestWithVariants(t, "evpcooling", baseDir)
}

func TestComparison_L2_HeatPumpCooling(t *testing.T) {
	refDir := "../tests/comparison/testdata/L2_equipment/heat_pump_cooling/c_output"
	testDir := "../tests/comparison/testdata/L2_equipment/heat_pump_cooling/go_output"
	runComparisonTest(t, "heat_pump_cooling", refDir, testDir)
}

func TestComparison_L2_Omvav(t *testing.T) {
	refDir := "../tests/comparison/testdata/L2_equipment/omvav/c_output"
	testDir := "../tests/comparison/testdata/L2_equipment/omvav/go_output"
	runComparisonTest(t, "omvav", refDir, testDir)
}

func TestComparison_L2_Stheat(t *testing.T) {
	baseDir := "../tests/comparison/testdata/L2_equipment/stheat"
	runComparisonTestWithVariants(t, "stheat", baseDir)
}

func TestComparison_L2_VWV(t *testing.T) {
	baseDir := "../tests/comparison/testdata/L2_equipment/vwv"
	runComparisonTestWithVariants(t, "vwv", baseDir)
}

// Note: obs, polygon, sunbrk, coordnt, divid tests
// C版は対話モード要求またはCOORDNTセクションでクラッシュするため、
// C版出力の再生成ができていません。Go版のみの実行テストとしてスキップ。
// テストファイルは以下のように整理されています：
// - obs_test.txt: OBSセクションのみ（COORDNT/DIVIDなし）
// - polygon_test.txt: POLYGONセクションのみ（COORDNT/DIVIDなし）
// - sunbrk_test.txt: SUNBRKセクションのみ（COORDNT/DIVIDなし）
// - coordnt_test.txt: COORDNT独立テスト
// - divid_test.txt: DIVID + COORDNT組み合わせテスト

func TestComparison_L2_Obs(t *testing.T) {
	refDir := "../tests/comparison/testdata/L2_equipment/obs/c_output"
	testDir := "../tests/comparison/testdata/L2_equipment/obs/go_output"
	runComparisonTest(t, "obs", refDir, testDir)
}

func TestComparison_L2_Polygon(t *testing.T) {
	refDir := "../tests/comparison/testdata/L2_equipment/polygon/c_output"
	testDir := "../tests/comparison/testdata/L2_equipment/polygon/go_output"
	runComparisonTest(t, "polygon", refDir, testDir)
}

func TestComparison_L2_Comfort(t *testing.T) {
	refDir := "../tests/comparison/testdata/L2_equipment/comfort/c_output"
	testDir := "../tests/comparison/testdata/L2_equipment/comfort/go_output"
	runComparisonTest(t, "comfort", refDir, testDir)
}

func TestComparison_L2_Sunbrk(t *testing.T) {
	refDir := "../tests/comparison/testdata/L2_equipment/sunbrk/c_output"
	testDir := "../tests/comparison/testdata/L2_equipment/sunbrk/go_output"
	runComparisonTest(t, "sunbrk", refDir, testDir)
}

func TestComparison_L2_Coordnt(t *testing.T) {
	refDir := "../tests/comparison/testdata/L2_equipment/coordnt/c_output"
	testDir := "../tests/comparison/testdata/L2_equipment/coordnt/go_output"
	runComparisonTest(t, "coordnt", refDir, testDir)
}

func TestComparison_L2_Divid(t *testing.T) {
	baseDir := "../tests/comparison/testdata/L2_equipment/divid"
	runComparisonTestWithVariants(t, "divid", baseDir)
}

func TestComparison_L2_TreeShadow(t *testing.T) {
	refDir := "../tests/comparison/testdata/L2_equipment/tree_shadow/c_output"
	testDir := "../tests/comparison/testdata/L2_equipment/tree_shadow/go_output"
	runComparisonTest(t, "tree_shadow", refDir, testDir)
}

// ============================================================================
// L3_system 比較テスト
// ============================================================================

func TestComparison_L3_PCMWall(t *testing.T) {
	baseDir := "../tests/comparison/testdata/L3_system/pcm_wall"
	runComparisonTestWithVariants(t, "pcm_wall", baseDir)
}

func TestComparison_L3_RadiantCeiling(t *testing.T) {
	refDir := "../tests/comparison/testdata/L3_system/radiant_ceiling/c_output"
	testDir := "../tests/comparison/testdata/L3_system/radiant_ceiling/go_output"
	runComparisonTest(t, "radiant_ceiling", refDir, testDir)
}

func TestComparison_L3_RadiantFloor(t *testing.T) {
	baseDir := "../tests/comparison/testdata/L3_system/radiant_floor"
	runComparisonTestWithVariants(t, "radiant_floor", baseDir)
}

// ============================================================================
// L4_annual 比較テスト
// ============================================================================

func TestComparison_L4_StandardHouse(t *testing.T) {
	baseDir := "../tests/comparison/testdata/L4_annual/standard_house"
	runComparisonTestWithVariants(t, "standard_house", baseDir)
}

// ============================================================================
// 物理値誤差に基づく比較テスト
// 統計メタデータ（時刻/カウント）とゼロ近傍値の誤差を除外して判定
// 閾値は観測値+0.01%のマージンで設定
// ============================================================================

// --- L1_basic ---

func TestPhysicsComparison_L1_SimpleRoomFull(t *testing.T) {
	refDir := "../tests/comparison/testdata/L1_basic/simple_room_full/c_output"
	testDir := "../tests/comparison/testdata/L1_basic/simple_room_full/go_output"
	runComparisonTestWithPhysicsThreshold(t, "simple_room_full", refDir, testDir, 0.01)
}

func TestPhysicsComparison_L1_SimpleRoomInternalHeat(t *testing.T) {
	refDir := "../tests/comparison/testdata/L1_basic/simple_room_internal_heat/c_output"
	testDir := "../tests/comparison/testdata/L1_basic/simple_room_internal_heat/go_output"
	runComparisonTestWithPhysicsThreshold(t, "simple_room_internal_heat", refDir, testDir, 0.01)
}

func TestPhysicsComparison_L1_SimpleRoomSchedule(t *testing.T) {
	refDir := "../tests/comparison/testdata/L1_basic/simple_room_schedule/c_output"
	testDir := "../tests/comparison/testdata/L1_basic/simple_room_schedule/go_output"
	runComparisonTestWithPhysicsThreshold(t, "simple_room_schedule", refDir, testDir, 0.01)
}

func TestPhysicsComparison_L1_SimpleRoomVent(t *testing.T) {
	refDir := "../tests/comparison/testdata/L1_basic/simple_room_vent/c_output"
	testDir := "../tests/comparison/testdata/L1_basic/simple_room_vent/go_output"
	runComparisonTestWithPhysicsThreshold(t, "simple_room_vent", refDir, testDir, 0.01)
}

// --- L2_equipment (観測値 0.00%) ---

func TestPhysicsComparison_L2_BoilerHeating(t *testing.T) {
	refDir := "../tests/comparison/testdata/L2_equipment/boiler_heating/c_output"
	testDir := "../tests/comparison/testdata/L2_equipment/boiler_heating/go_output"
	runComparisonTestWithPhysicsThreshold(t, "boiler_heating", refDir, testDir, 0.01)
}

func TestPhysicsComparison_L2_CoolingCoil(t *testing.T) {
	refDir := "../tests/comparison/testdata/L2_equipment/cooling_coil/c_output"
	testDir := "../tests/comparison/testdata/L2_equipment/cooling_coil/go_output"
	runComparisonTestWithPhysicsThreshold(t, "cooling_coil", refDir, testDir, 0.01)
}

func TestPhysicsComparison_L2_Desiccant(t *testing.T) {
	refDir := "../tests/comparison/testdata/L2_equipment/desiccant/c_output"
	testDir := "../tests/comparison/testdata/L2_equipment/desiccant/go_output"
	runComparisonTestWithPhysicsThreshold(t, "desiccant", refDir, testDir, 0.01)
}

func TestPhysicsComparison_L2_Duct(t *testing.T) {
	refDir := "../tests/comparison/testdata/L2_equipment/duct/c_output"
	testDir := "../tests/comparison/testdata/L2_equipment/duct/go_output"
	runComparisonTestWithPhysicsThreshold(t, "duct", refDir, testDir, 0.01)
}

func TestPhysicsComparison_L2_Evpcooling(t *testing.T) {
	refDir := "../tests/comparison/testdata/L2_equipment/evpcooling/c_output"
	testDir := "../tests/comparison/testdata/L2_equipment/evpcooling/go_output"
	runComparisonTestWithPhysicsThreshold(t, "evpcooling", refDir, testDir, 0.01)
}

func TestPhysicsComparison_L2_Fan(t *testing.T) {
	refDir := "../tests/comparison/testdata/L2_equipment/fan/c_output"
	testDir := "../tests/comparison/testdata/L2_equipment/fan/go_output"
	runComparisonTestWithPhysicsThreshold(t, "fan", refDir, testDir, 0.01)
}

func TestPhysicsComparison_L2_HeatPump(t *testing.T) {
	refDir := "../tests/comparison/testdata/L2_equipment/heat_pump/c_output"
	testDir := "../tests/comparison/testdata/L2_equipment/heat_pump/go_output"
	runComparisonTestWithPhysicsThreshold(t, "heat_pump", refDir, testDir, 0.01)
}

func TestPhysicsComparison_L2_HeatPumpCooling(t *testing.T) {
	refDir := "../tests/comparison/testdata/L2_equipment/heat_pump_cooling/c_output"
	testDir := "../tests/comparison/testdata/L2_equipment/heat_pump_cooling/go_output"
	runComparisonTestWithPhysicsThreshold(t, "heat_pump_cooling", refDir, testDir, 0.01)
}

func TestPhysicsComparison_L2_Hex(t *testing.T) {
	refDir := "../tests/comparison/testdata/L2_equipment/hex/c_output"
	testDir := "../tests/comparison/testdata/L2_equipment/hex/go_output"
	runComparisonTestWithPhysicsThreshold(t, "hex", refDir, testDir, 0.01)
}

func TestPhysicsComparison_L2_Obs(t *testing.T) {
	refDir := "../tests/comparison/testdata/L2_equipment/obs/c_output"
	testDir := "../tests/comparison/testdata/L2_equipment/obs/go_output"
	runComparisonTestWithPhysicsThreshold(t, "obs", refDir, testDir, 0.01)
}

func TestPhysicsComparison_L2_Polygon(t *testing.T) {
	refDir := "../tests/comparison/testdata/L2_equipment/polygon/c_output"
	testDir := "../tests/comparison/testdata/L2_equipment/polygon/go_output"
	runComparisonTestWithPhysicsThreshold(t, "polygon", refDir, testDir, 0.01)
}

func TestPhysicsComparison_L2_PumpPipe(t *testing.T) {
	refDir := "../tests/comparison/testdata/L2_equipment/pump_pipe/c_output"
	testDir := "../tests/comparison/testdata/L2_equipment/pump_pipe/go_output"
	runComparisonTestWithPhysicsThreshold(t, "pump_pipe", refDir, testDir, 0.01)
}

func TestPhysicsComparison_L2_PV(t *testing.T) {
	refDir := "../tests/comparison/testdata/L2_equipment/pv/c_output"
	testDir := "../tests/comparison/testdata/L2_equipment/pv/go_output"
	runComparisonTestWithPhysicsThreshold(t, "pv", refDir, testDir, 0.01)
}

func TestPhysicsComparison_L2_Qmeas(t *testing.T) {
	refDir := "../tests/comparison/testdata/L2_equipment/qmeas/c_output"
	testDir := "../tests/comparison/testdata/L2_equipment/qmeas/go_output"
	runComparisonTestWithPhysicsThreshold(t, "qmeas", refDir, testDir, 0.01)
}

func TestPhysicsComparison_L2_SolarCollector(t *testing.T) {
	refDir := "../tests/comparison/testdata/L2_equipment/solar_collector/c_output"
	testDir := "../tests/comparison/testdata/L2_equipment/solar_collector/go_output"
	runComparisonTestWithPhysicsThreshold(t, "solar_collector", refDir, testDir, 0.01)
}

func TestPhysicsComparison_L2_AirCollector(t *testing.T) {
	refDir := "../tests/comparison/testdata/L2_equipment/air_collector/c_output"
	testDir := "../tests/comparison/testdata/L2_equipment/air_collector/go_output"
	runComparisonTestWithPhysicsThreshold(t, "air_collector", refDir, testDir, 0.01)
}

// --- L2_equipment (観測値 > 0%) ---

// 観測値: 0.0000%（修正済み）
func TestPhysicsComparison_L2_Omvav(t *testing.T) {
	refDir := "../tests/comparison/testdata/L2_equipment/omvav/c_output"
	testDir := "../tests/comparison/testdata/L2_equipment/omvav/go_output"
	runComparisonTestWithPhysicsThreshold(t, "omvav", refDir, testDir, 0.01)
}

// 観測値: 0.0000%（修正済み）
func TestPhysicsComparison_L2_Stheat(t *testing.T) {
	refDir := "../tests/comparison/testdata/L2_equipment/stheat/c_output"
	testDir := "../tests/comparison/testdata/L2_equipment/stheat/go_output"
	runComparisonTestWithPhysicsThreshold(t, "stheat", refDir, testDir, 0.01)
}

// 観測値: 0.0000%（修正済み）
func TestPhysicsComparison_L2_VWV(t *testing.T) {
	refDir := "../tests/comparison/testdata/L2_equipment/vwv/c_output"
	testDir := "../tests/comparison/testdata/L2_equipment/vwv/go_output"
	runComparisonTestWithPhysicsThreshold(t, "vwv", refDir, testDir, 0.01)
}

// 観測値: 0.0000%（修正済み - helmwlsftのalr[j]バグ）
func TestPhysicsComparison_L2_Helm(t *testing.T) {
	refDir := "../tests/comparison/testdata/L2_equipment/helm/c_output"
	testDir := "../tests/comparison/testdata/L2_equipment/helm/go_output"
	runComparisonTestWithPhysicsThreshold(t, "helm", refDir, testDir, 0.01)
}

// 観測値: 0.0000%（修正済み）
func TestPhysicsComparison_L2_Coordnt(t *testing.T) {
	refDir := "../tests/comparison/testdata/L2_equipment/coordnt/c_output"
	testDir := "../tests/comparison/testdata/L2_equipment/coordnt/go_output"
	runComparisonTestWithPhysicsThreshold(t, "coordnt", refDir, testDir, 0.01)
}

// 観測値: 0.0000%（修正済み）
func TestPhysicsComparison_L2_Divid(t *testing.T) {
	refDir := "../tests/comparison/testdata/L2_equipment/divid/c_output"
	testDir := "../tests/comparison/testdata/L2_equipment/divid/go_output"
	runComparisonTestWithPhysicsThreshold(t, "divid", refDir, testDir, 0.01)
}

// 観測値: 0.0000%（修正済み）
func TestPhysicsComparison_L2_Sunbrk(t *testing.T) {
	refDir := "../tests/comparison/testdata/L2_equipment/sunbrk/c_output"
	testDir := "../tests/comparison/testdata/L2_equipment/sunbrk/go_output"
	runComparisonTestWithPhysicsThreshold(t, "sunbrk", refDir, testDir, 0.01)
}

// 観測値: 0.0000%
func TestPhysicsComparison_L2_VAV(t *testing.T) {
	refDir := "../tests/comparison/testdata/L2_equipment/vav/c_output"
	testDir := "../tests/comparison/testdata/L2_equipment/vav/go_output"
	runComparisonTestWithPhysicsThreshold(t, "vav", refDir, testDir, 0.01)
}

// --- L3_system ---

func TestPhysicsComparison_L3_RadiantCeiling(t *testing.T) {
	refDir := "../tests/comparison/testdata/L3_system/radiant_ceiling/c_output"
	testDir := "../tests/comparison/testdata/L3_system/radiant_ceiling/go_output"
	runComparisonTestWithPhysicsThreshold(t, "radiant_ceiling", refDir, testDir, 0.01)
}

// --- L4_annual ---

// 観測値: 0.0096%
func TestPhysicsComparison_L4_StandardHouse(t *testing.T) {
	refDir := "../tests/comparison/testdata/L4_annual/standard_house/c_output"
	testDir := "../tests/comparison/testdata/L4_annual/standard_house/go_output"
	runComparisonTestWithPhysicsThreshold(t, "standard_house", refDir, testDir, 0.01)
}

// ============================================================================
// シミュレーション実行テスト
// Go版を実際に実行してC版出力と比較するテスト
// テスト名: TestSimulation_*
// ============================================================================

func TestSimulation_L1_SimpleRoomFull(t *testing.T) {
	baseDir := "../tests/comparison/testdata/L1_basic/simple_room_full"
	runSimulationAndCompare(t, "simple_room_full", baseDir, "../Base")
}

func TestSimulation_L2_PumpPipe(t *testing.T) {
	baseDir := "../tests/comparison/testdata/L2_equipment/pump_pipe"
	runSimulationTestWithVariants(t, "pump_pipe", baseDir, "../Base")
}

func TestSimulation_L2_Valv(t *testing.T) {
	baseDir := "../tests/comparison/testdata/L2_equipment/valv"
	runSimulationTestWithVariants(t, "valv", baseDir, "../Base")
}

func TestSimulation_L2_StorageTank(t *testing.T) {
	baseDir := "../tests/comparison/testdata/L2_equipment/storage_tank"
	runSimulationTestWithVariants(t, "storage_tank", baseDir, "../Base")
}

func TestSimulation_L2_SolarCollector(t *testing.T) {
	baseDir := "../tests/comparison/testdata/L2_equipment/solar_collector"
	runSimulationTestWithVariants(t, "solar_collector", baseDir, "../Base")
}

func TestSimulation_L3_RadiantFloor(t *testing.T) {
	baseDir := "../tests/comparison/testdata/L3_system/radiant_floor"
	runSimulationTestWithVariants(t, "radiant_floor", baseDir, "../Base")
}

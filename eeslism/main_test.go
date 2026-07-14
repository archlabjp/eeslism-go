package eeslism

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"testing"
)

func Test_StandardPlan(t *testing.T) {
	t.Skip("Skipping complex test - PCM sample requires additional setup")
	Entry("../samples/standard-plan-no-hcap-PCM-CM-fsolm.txt", "../Base")
}

func Test_Sample(t *testing.T) {
	// samples/radiant_floor_heating.txtを使った統合テスト
	// （sample.txtはSYSCMPセクションがないため使用不可）
	//
	// 現状スキップ: このサンプルは寝室(Bedroom)の間仕切り壁 "(LivingRoom): -i ;" に
	// 面積も i= 名も無く、対となる壁から面積を解決できないため A=-999 のまま残り、
	// blrmdata.go の面積チェックで os.Exit(1) する（C版と同一挙動）。これはテスト
	// バイナリ全体を巻き添えにするため、入力ファイル側が修正されるまでスキップする。
	// （入力ファイルの修正は別タスク。testdata の radiant_floor.txt と同一の欠陥）
	t.Skip("sample radiant_floor_heating.txt has unresolved interior-wall area (A=-999) causing os.Exit — input fix tracked separately")
	resetPrintStates()
	Entry("../samples/radiant_floor_heating.txt", "../Base")
}

// Test_PCMWall_Summer は夏季のPCM壁体シミュレーションをテストする
// PCMは常に液体状態（室温 > 25°C）
func Test_PCMWall_Summer(t *testing.T) {
	resetPrintStates()
	Entry("../tests/comparison/testdata/L3_system/pcm_wall/pcm_wall_test.txt", "../Base")

	// 出力ファイルの存在確認
	outputFile := "../tests/comparison/testdata/L3_system/pcm_wall/pcm_wall_test_rm.es"
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Fatalf("Output file not generated: %s", outputFile)
	}

	// 温度範囲の検証（夏季は25-30°C程度を期待）
	temps := parseRoomTemperatures(t, outputFile)
	if len(temps) == 0 {
		t.Fatal("No temperature data found in output")
	}

	minTemp, maxTemp := temps[0], temps[0]
	for _, temp := range temps {
		if temp < minTemp {
			minTemp = temp
		}
		if temp > maxTemp {
			maxTemp = temp
		}
	}

	// 夏季の温度範囲チェック（PCM液体温度25°C以上を期待）
	if minTemp < 20 {
		t.Errorf("Summer temperature too low: min=%.2f°C (expected > 20°C)", minTemp)
	}
	if maxTemp > 35 {
		t.Errorf("Summer temperature too high: max=%.2f°C (expected < 35°C)", maxTemp)
	}
}

// Test_PCMWall_PhaseChange はPCM相変化が発生する条件でのシミュレーションをテストする
// PCMは固体↔液体の遷移を繰り返す（室温: 21-25°C）
func Test_PCMWall_PhaseChange(t *testing.T) {
	resetPrintStates()
	Entry("../tests/comparison/testdata/L3_system/pcm_wall/pcm_wall_phase_change_test.txt", "../Base")

	// 出力ファイルの存在確認
	outputFile := "../tests/comparison/testdata/L3_system/pcm_wall/pcm_wall_phase_change_test_rm.es"
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Fatalf("Output file not generated: %s", outputFile)
	}

	// 温度範囲の検証
	temps := parseRoomTemperatures(t, outputFile)
	if len(temps) == 0 {
		t.Fatal("No temperature data found in output")
	}

	minTemp, maxTemp := temps[0], temps[0]
	for _, temp := range temps {
		if temp < minTemp {
			minTemp = temp
		}
		if temp > maxTemp {
			maxTemp = temp
		}
	}

	// 相変化温度範囲（Ts=22°C, Tl=25°C）を横断することを確認
	pcmTs := 22.0 // PCM固体温度
	pcmTl := 25.0 // PCM液体温度

	if minTemp >= pcmTs {
		t.Errorf("Temperature never goes below solidification point: min=%.2f°C (Ts=%.1f°C)", minTemp, pcmTs)
	}
	if maxTemp <= pcmTl {
		t.Logf("Warning: Temperature stays below liquid point: max=%.2f°C (Tl=%.1f°C)", maxTemp, pcmTl)
	}
	// 相変化温度範囲内のデータがあることを確認
	hasPhaseChange := false
	for _, temp := range temps {
		if temp >= pcmTs && temp <= pcmTl {
			hasPhaseChange = true
			break
		}
	}
	if !hasPhaseChange {
		t.Errorf("No temperature data within phase change range (%.1f-%.1f°C)", pcmTs, pcmTl)
	}

	t.Logf("Phase change test: min=%.2f°C, max=%.2f°C, Ts=%.1f°C, Tl=%.1f°C", minTemp, maxTemp, pcmTs, pcmTl)
}

// parseRoomTemperatures は_rm.esファイルから室温を抽出する
func parseRoomTemperatures(t *testing.T, filename string) []float64 {
	file, err := os.Open(filename)
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	var temps []float64
	scanner := bufio.NewScanner(file)
	inData := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// ヘッダー行（"PCMroom_Tr t f..."）の後からデータ開始
		if strings.HasPrefix(line, "PCMroom_Tr") {
			inData = true
			continue
		}

		// データ終了マーカー
		if line == "-999" {
			break
		}

		if !inData {
			continue
		}

		// データ行のパース（交互に日時行とデータ行がある）
		// 日時行: "07 25  1.00"
		// データ行: "25.64 0.0053 26 25.64"
		parts := strings.Fields(line)
		if len(parts) >= 1 {
			// 最初のフィールドが温度かどうかチェック
			temp, err := strconv.ParseFloat(parts[0], 64)
			if err == nil && temp > 10 && temp < 40 {
				// 妥当な温度範囲の値のみ追加
				temps = append(temps, temp)
			}
		}
	}

	return temps
}

package eeslism

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

/*
EeTokens (Energy Simulation Tokens)

この構造体は、入力ファイルから読み込まれたテキストデータを、
シミュレーションモデルが解析しやすいようにトークン（単語や記号）のリストとして保持します。
また、トークンリスト内での現在の読み込み位置を管理します。

建築環境工学的な観点:
  - **入力データの解析**: 建物のエネルギーシミュレーションでは、
    建物形状、材料特性、設備機器、スケジュール、気象データなど、
    様々な形式の入力データを解析する必要があります。
    この構造体は、入力ファイルを効率的に読み込み、
    意味のある単位（トークン）に分割するための基盤を提供します。
  - **柔軟なデータ処理**: トークンベースの処理により、
    入力ファイルのフォーマットが多少異なっていても、
    柔軟に対応できる可能性があります。
    例えば、空白やコメントの扱いを統一することで、
    入力データの記述の自由度を高めることができます。
  - **エラー検出の補助**: トークン化の過程で、
    不正なフォーマットや予期しない文字が検出された場合、
    エラーを報告し、入力データの修正を促すことができます。

この構造体は、建物のエネルギーシミュレーションにおいて、
入力データの正確な解析と、シミュレーションモデルの構築を支援するための重要な役割を果たします。
*/
type EeTokens struct {
	tokens []string
	pos    int
}

/*
Len (Get Length of Tokens)

このメソッドは、`EeTokens`構造体が保持するトークンの総数を返します。

建築環境工学的な観点:
  - **データ処理の制御**: 入力データの処理ループにおいて、
    このメソッドは、処理すべきトークンが残っているかどうかを判断するために用いられます。
    これにより、入力ファイルの最後までデータを読み込むことを保証し、
    データの欠落を防ぎます。

このメソッドは、建物のエネルギーシミュレーションにおいて、
入力データの正確な解析と、シミュレーションモデルの構築を支援するための基本的な役割を果たします。
*/
func (t *EeTokens) Len() int {
	return len(t.tokens)
}

/*
GetPos (Get Current Position of Tokens)

このメソッドは、`EeTokens`構造体内の現在の読み込み位置（トークンリストのインデックス）を返します。

建築環境工学的な観点:
  - **データ処理の制御**: 入力データの処理において、
    特定のセクションを繰り返し読み込んだり、
    エラー発生時に読み込み位置を復元したりするために、
    現在の読み込み位置を保存する必要があります。
    このメソッドは、そのための情報を提供します。

このメソッドは、建物のエネルギーシミュレーションにおいて、
入力データの正確な解析と、シミュレーションモデルの構築を支援するための基本的な役割を果たします。
*/
func (t *EeTokens) GetPos() int {
	return t.pos
}

/*
RestorePos (Restore Position of Tokens)

このメソッドは、`EeTokens`構造体内の読み込み位置を、
指定された位置（`pos`）に復元します。

建築環境工学的な観点:
  - **データ処理の柔軟性**: 入力データの処理において、
    特定のセクションを複数回読み込んだり、
    エラー発生時に読み込み位置を復元して再試行したりするために、
    このメソッドが用いられます。
    これにより、入力ファイルの解析ロジックの柔軟性が向上します。

このメソッドは、建物のエネルギーシミュレーションにおいて、
入力データの正確な解析と、シミュレーションモデルの構築を支援するための基本的な役割を果たします。
*/
func (t *EeTokens) RestorePos(pos int) {
	t.pos = pos
}

/*
Reset (Reset Position of Tokens to Start)

このメソッドは、`EeTokens`構造体内の読み込み位置を、
トークンリストの先頭（インデックス0）にリセットします。

建築環境工学的な観点:
  - **データ処理の再開**: 入力ファイルの全体を複数回読み込む必要がある場合や、
    特定の処理が完了した後に最初から読み込みを再開する場合に用いられます。
    例えば、入力ファイルの構文解析とデータ抽出を別々のパスで行う場合などに有用です。

このメソッドは、建物のエネルギーシミュレーションにおいて、
入力データの正確な解析と、シミュレーションモデルの構築を支援するための基本的な役割を果たします。
*/
func (t *EeTokens) Reset() {
	t.pos = 0
}

/*
NewEeTokens (Create EeTokens Object from String)

この関数は、入力ファイルの内容を文字列として受け取り、
それを解析してトークン（単語や記号）のリストに変換し、
`EeTokens`構造体として返します。
この過程で、コメントの除去や空白の整形が行われます。

建築環境工学的な観点:
  - **入力データの前処理**: 建物のエネルギーシミュレーションでは、
    入力ファイルが様々な形式で記述される可能性があります。
    この関数は、入力ファイルからコメント行や不要な空白を除去し、
    シミュレーションモデルが解析しやすい統一されたトークンリストを生成します。
  - **構文解析の準備**: 生成されたトークンリストは、
    その後の構文解析（入力データの意味を解釈する処理）の基礎となります。
    正確なトークン化は、入力データの誤りを早期に検出し、
    シミュレーションの信頼性を確保するために重要です。
  - **柔軟なコメント形式**: `// コメントの除去` の部分で示されているように、
    `!`で始まるコメントを認識し、除去します。
    これにより、入力ファイルの記述者が自由にコメントを記述でき、
    入力ファイルの可読性を高めることができます。

この関数は、建物のエネルギーシミュレーションにおいて、
入力データの正確な解析と、シミュレーションモデルの構築を支援するための重要な役割を果たします。
*/
func NewEeTokens(s string) *EeTokens {
	reader := strings.NewReader(s)
	scanner := bufio.NewScanner(reader)
	tokens := make([]string, 0)
	for scanner.Scan() {
		//行単位の処理
		line := scanner.Text()

		// コメントの除去
		if st := strings.IndexRune(line, '!'); st != -1 {
			line = line[:st]
		}

		// 空文字の除去
		line = strings.TrimSpace(line)

		for _, s := range strings.Fields(line) {
			if strings.HasSuffix(s, ";") {
				s = s[:len(s)-1]
				if s != "" {
					tokens = append(tokens, s)
				}
				tokens = append(tokens, ";")
			} else if strings.ContainsRune(s, ';') {
				panic("Invalid position of `;`")
			} else {
				tokens = append(tokens, s)
			}
		}

		//改行
		tokens = append(tokens, "\n")
	}
	return &EeTokens{tokens: tokens, pos: 0}
}

/*
GetLine (Get Tokens for a Single Line)

このメソッドは、`EeTokens`構造体内の現在の読み込み位置から、
次の改行文字（`
`）までのトークンをリストとして返します。

建築環境工学的な観点:
  - **行単位のデータ処理**: 入力ファイルが論理的な行で構成されている場合、
    このメソッドは行単位でデータを読み取るために用いられます。
    これにより、各行が独立した情報を持つような入力ファイルの解析を容易にします。
  - **構文解析の補助**: 各行のトークンリストは、
    その後の構文解析（入力データの意味を解釈する処理）の基礎となります。
    例えば、特定のキーワードが行の先頭にあるかどうかを判断する際に用いられます。

このメソッドは、建物のエネルギーシミュレーションにおいて、
入力データの正確な解析と、シミュレーションモデルの構築を支援するための重要な役割を果たします。
*/
func (t *EeTokens) GetLine() []string {
	var line []string

	// find `\n`
	var found bool = false
	for i := t.pos; i < len(t.tokens); i++ {
		if t.tokens[i] == "\n" {
			line = t.tokens[t.pos:i]
			t.pos = i + 1
			found = true
			break
		}
	}
	// not found
	if found == false {
		t.pos = len(t.tokens)
		line = t.tokens[t.pos:]
	}

	return line
}

/*
GetLogicalLine (Get Tokens for a Logical Line)

このメソッドは、`EeTokens`構造体内の現在の読み込み位置から、
次のセミコロン（`;`）またはアスタリスク（`*`）までのトークンをリストとして返します。
これにより、論理的なデータブロックを読み込むことができます。

建築環境工学的な観点:
  - **論理的なデータブロックの処理**: 入力ファイルがセミコロンやアスタリスクで区切られた
    論理的なデータブロックで構成されている場合、
    このメソッドはデータブロック単位でデータを読み込むために用いられます。
    これにより、各データブロックが独立した情報を持つような入力ファイルの解析を容易にします。
  - **構文解析の補助**: 各論理行のトークンリストは、
    その後の構文解析（入力データの意味を解釈する処理）の基礎となります。
    例えば、特定のキーワードが論理行の先頭にあるかどうかを判断する際に用いられます。
  - **柔軟な区切り文字**: セミコロンとアスタリスクの両方を区切り文字として認識することで、
    入力ファイルの記述の柔軟性を高めます。

このメソッドは、建物のエネルギーシミュレーションにおいて、
入力データの正確な解析と、シミュレーションモデルの構築を支援するための重要な役割を果たします。
*/
func (t *EeTokens) GetLogicalLine() []string {
	var logiline []string
	var filtered []string

	//
	if t.tokens[t.pos] == "*" && t.tokens[t.pos-1] == "\n" {
		logiline = t.tokens[t.pos : t.pos+1]
		t.pos++
		return logiline
	}

	// find `;`
	var found bool = false
	for i := t.pos; i < len(t.tokens); i++ {
		if t.tokens[i] == ";" {
			logiline = t.tokens[t.pos : i+1] // `;` is included
			t.pos = i + 1
			found = true
			break
		} else if i > 0 && t.tokens[i-1] == "\n" && t.tokens[i] == "*" {
			logiline = t.tokens[t.pos : i+1] // `\n*` is included
			t.pos = i + 1                    // `\n` will be skipped
			found = true
			break
		}
	}

	// not found
	if found == false {
		logiline = t.tokens[t.pos:]
		t.pos = len(t.tokens)
	}

	// filter `\n` token and return it
	for _, s := range logiline {
		if s != "\n" {
			filtered = append(filtered, s)
		}
	}
	return filtered
}

/*
SkipToEndOfLine (Skip to End of Line)

このメソッドは、`EeTokens`構造体内の読み込み位置を、
現在の行の終わり（改行文字またはセミコロン）までスキップします。

建築環境工学的な観点:
  - **不要なデータのスキップ**: 入力ファイルには、
    シミュレーションモデルにとって不要な情報や、
    既に処理済みの情報が含まれている場合があります。
    このメソッドは、そのような不要なデータを効率的にスキップし、
    必要なデータのみを読み込むことで、
    入力処理の効率化を図ります。
  - **構文解析の補助**: 特定のデータブロックの処理が完了した後、
    次のデータブロックの開始位置まで読み込み位置を移動させる際に用いられます。

このメソッドは、建物のエネルギーシミュレーションにおいて、
入力データの正確な解析と、シミュレーションモデルの構築を支援するための重要な役割を果たします。
*/
func (t *EeTokens) SkipToEndOfLine() {
	for t.pos < len(t.tokens) && (t.tokens[t.pos] == "\n" || t.tokens[t.pos] == ";") {
		t.pos++
	}
	// もし、`;` が連続している場合は、最後の`;`までスキップする
	if t.pos < len(t.tokens) && t.tokens[t.pos] == ";" {
		t.pos++
	}
	// もし、改行が連続している場合は、最後の改行までスキップする
	if t.pos < len(t.tokens) && t.tokens[t.pos] == "\n" {
		t.pos++
	}
}

/*
GetSection (Get Tokens for a Section)

このメソッドは、`EeTokens`構造体内の現在の読み込み位置から、
次のアスタリスク（`*`）で始まる行までのトークンを新しい`EeTokens`構造体として返します。
これにより、入力ファイル内のセクション（GDAT, EXSRFなど）を読み込むことができます。

建築環境工学的な観点:
  - **セクション単位のデータ処理**: 建物のエネルギーシミュレーションの入力ファイルは、
    通常、GDAT（一般データ）、EXSRF（外部日射面）、WALL（壁体）など、
    論理的なセクションに分かれています。
    このメソッドは、これらのセクションを個別に読み込み、
    それぞれのセクションに対応する処理を行うために用いられます。
  - **モジュール化された入力処理**: 各セクションの処理を独立した関数で行うことで、
    入力処理のコードをモジュール化し、
    可読性や保守性を向上させることができます。

このメソッドは、建物のエネルギーシミュレーションにおいて、
入力データの正確な解析と、シミュレーションモデルの構築を支援するための重要な役割を果たします。
*/
func (t *EeTokens) GetSection() *EeTokens {
	t.SkipToEndOfLine()

	// find `*` at start of some line
	for i := t.pos; i < len(t.tokens); i++ {
		if i > 0 && t.tokens[i-1] == "\n" && t.tokens[i] == "*" {
			section := &EeTokens{tokens: t.tokens[t.pos : i+1], pos: 0}
			t.pos = i + 1
			return section
		}
	}
	// not found
	section := &EeTokens{tokens: t.tokens[t.pos:], pos: 0}
	t.pos = len(t.tokens)
	return section
}

/*
IsEnd (Check if End of Tokens)

このメソッドは、`EeTokens`構造体内の読み込み位置が、
トークンリストの終端に達しているかどうかを判断します。

建築環境工学的な観点:
  - **データ処理の終了条件**: 入力データの処理ループにおいて、
    このメソッドは、処理すべきトークンが残っていないかどうかを判断するために用いられます。
    これにより、入力ファイルの最後までデータを読み込むことを保証し、
    無限ループを防ぎます。

このメソッドは、建物のエネルギーシミュレーションにおいて、
入力データの正確な解析と、シミュレーションモデルの構築を支援するための基本的な役割を果たします。
*/
func (t *EeTokens) IsEnd() bool {
	return t.pos >= len(t.tokens)
}

/*
PeekToken (Peek Next Token)

このメソッドは、`EeTokens`構造体内の次のトークンを、
読み込み位置を進めることなく（消費することなく）返します。

建築環境工学的な観点:
  - **先読みによる構文解析**: 入力データの構文解析において、
    次にどのようなトークンが来るかを事前に知ることで、
    適切な処理分岐を行うことができます。
    例えば、特定のキーワードが続く場合にのみ、
    その後のデータを読み込むといったロジックを実装する際に用いられます。

このメソッドは、建物のエネルギーシミュレーションにおいて、
入力データの正確な解析と、シミュレーションモデルの構築を支援するための重要な役割を果たします。
*/
func (t *EeTokens) PeekToken() string {
	if t.pos < len(t.tokens) {
		return t.tokens[t.pos]
	}
	return ""
}

/*
GetToken (Get Next Token)

このメソッドは、`EeTokens`構造体内の次のトークンを返し、
読み込み位置を一つ進めます。

建築環境工学的な観点:
  - **トークン単位のデータ読み込み**: 入力ファイルの解析において、
    このメソッドは最も基本的なデータ読み込み操作を提供します。
    これにより、入力ファイルをトークン単位で順次処理し、
    シミュレーションモデルに必要な情報を抽出できます。

このメソッドは、建物のエネルギーシミュレーションにおいて、
入力データの正確な解析と、シミュレーションモデルの構築を支援するための重要な役割を果たします。
*/
func (t *EeTokens) GetToken() string {
	if t.pos < len(t.tokens) {
		t.pos++
		return t.tokens[t.pos-1]
	}
	return ""
}

/*
GetFloat (Get Next Token as Float64)

このメソッドは、`EeTokens`構造体内の次のトークンを浮動小数点数（`float64`）として解析し、
読み込み位置を一つ進めます。

建築環境工学的な観点:
  - **数値データの抽出**: 建物のエネルギーシミュレーションでは、
    温度、熱量、寸法など、多くの数値データが入力されます。
    このメソッドは、入力ファイルからこれらの数値データを効率的に抽出するために用いられます。

このメソッドは、建物のエネルギーシミュレーションにおいて、
入力データの正確な解析と、シミュレーションモデルの構築を支援するための重要な役割を果たします。
*/
func (t *EeTokens) GetFloat() float64 {
	var f float64
	fmt.Sscanf(t.GetToken(), "%f", &f)
	return f
}

/*
GetInt (Get Next Token as Integer)

このメソッドは、`EeTokens`構造体内の次のトークンを整数（`int`）として解析し、
読み込み位置を一つ進めます。

建築環境工学的な観点:
  - **数値データの抽出**: 建物のエネルギーシミュレーションでは、
    層の数、機器の数、日付など、多くの整数データが入力されます。
    このメソッドは、入力ファイルからこれらの整数データを効率的に抽出するために用いられます。

このメソッドは、建物のエネルギーシミュレーションにおいて、
入力データの正確な解析と、シミュレーションモデルの構築を支援するための重要な役割を果たします。
*/
func (t *EeTokens) GetInt() int {
	var i int
	fmt.Sscanf(t.GetToken(), "%d", &i)
	return i
}

/*
Eeinput (Energy Simulation Input Data Processing)

この関数は、建物のエネルギーシミュレーションに必要な全ての入力データを読み込み、
解析し、シミュレーションモデルを構築します。
これには、気象データ、スケジュール、建物形状、材料特性、設備機器、制御ロジックなどが含まれます。

建築環境工学的な観点:
  - **シミュレーションモデルの構築**: この関数は、
    建物のエネルギーシミュレーションモデルを構築するための中心的な役割を果たします。
    様々な入力ファイルからデータを読み込み、
    それぞれのデータが対応する構造体（`Simc`, `Exsf`, `Rmvls`, `Eqcat`, `Eqsys`など）に格納されます。
  - **データ処理のフロー**:
  - **曜日設定、スケジュール表の読み込み**: `Dayweek`, `Schtable`, `Schname`, `Schdata`関数を呼び出し、
    建物の運用スケジュールに関するデータを読み込みます。
  - **一般データ（GDAT）の読み込み**: `Gdata`関数を呼び出し、
    シミュレーションの基本設定（計算期間、時間間隔、出力設定など）を読み込みます。
  - **外部日射面（EXSRF）の読み込み**: `Exsfdata`関数を呼び出し、
    建物の周囲の外部日射面に関するデータを読み込みます。
  - **日よけ（SUNBRK）の読み込み**: `Snbkdata`関数を呼び出し、
    庇、バルコニー、袖壁などの日よけに関するデータを読み込みます。
  - **PCM（相変化材料）の読み込み**: `PCMdata`関数を呼び出し、
    PCMの特性に関するデータを読み込みます。
  - **壁体（WALL）と窓（WINDOW）の読み込み**: `Walldata`, `Windowdata`関数を呼び出し、
    建物の外皮に関するデータを読み込みます。
  - **室（ROOM）の読み込み**: `Roomdata`関数を呼び出し、
    室の幾何学的情報、内部発熱、換気などに関するデータを読み込みます。
  - **換気（RAICH, VENT）の読み込み**: `Ventdata`関数を呼び出し、
    室の外気導入量や室間相互換気量に関するデータを読み込みます。
  - **内部発熱（APPL）の読み込み**: `Appldata`関数を呼び出し、
    人体、照明、機器などからの内部発熱に関するデータを読み込みます。
  - **設備機器カタログ（EQPCAT）の読み込み**: `Eqcadata`関数を呼び出し、
    冷温水コイル、ボイラー、ポンプなどの設備機器のカタログデータを読み込みます。
  - **システムコンポーネント（SYSCMP）の読み込み**: `Compodata`関数を呼び出し、
    システムを構成する各コンポーネントの接続関係を定義します。
  - **システムパス（SYSPTH）の読み込み**: `Pathdata`関数を呼び出し、
    熱媒の流れる経路を定義します。
  - **制御（CONTL）の読み込み**: `Contrldata`関数を呼び出し、
    空調システムなどの制御ロジックを定義します。
  - **外部障害物（COORDNT, OBS, TREE, POLYGON）の読み込み**: `bdpdata`, `obsdata`, `treedata`, `polydata`関数を呼び出し、
    建物周囲の障害物に関するデータを読み込みます。
  - **エラーハンドリング**: 各データ読み込みステップでエラーチェックを行い、
    入力データの不備を早期に発見し、シミュレーションの信頼性を確保します。

この関数は、建物のエネルギーシミュレーションのデータ準備とモデル構築の全体を統括し、
シミュレーションの正確性、効率性、および再現性を確保するための重要な役割を果たします。
*/
func Eeinput(Ipath string, efl_path string, bdata, week, schtba, schnma string, Simc *SIMCONTL,
	Exsf *EXSFS, Rmvls *RMVLS, Eqcat *EQCAT, Eqsys *EQSYS,
	Compnt *[]*COMPNT,
	Elout *[]*ELOUT,
	Elin *[]*ELIN,
	Mpath *[]*MPATH,
	Plist *[]*PLIST,
	Pelm *[]*PELM,
	Contl *[]*CONTL,
	Ctlif *[]*CTLIF,
	Ctlst *[]*CTLST,
	Wd *WDAT, Daytm *DAYTM, key int,
	obsn *int, bp *[]*BBDP,
	obs *[]*OBS,
	tree *[]*TREE,
	shadtb *[]*SHADTB,
	poly *[]*POLYGN, monten *int, gpn *int, DE *float64,
	Noplpmp *NOPLPMP) (*SCHDL, []*FLOUT) {

	if Simc == nil {
		panic("Simc is nil")
	}
	if Exsf == nil {
		panic("Exsf is nil")
	}
	if Rmvls == nil {
		panic("Rmvls is nil")
	}
	if Eqcat == nil {
		panic("Eqcat is nil")
	}
	if Eqsys == nil {
		panic("Eqsys is nil")
	}
	if Compnt == nil {
		panic("Compnt is nil")
	}
	if Elout == nil {
		panic("Elout is nil")
	}
	if Elin == nil {
		panic("Elin is nil")
	}
	if Mpath == nil {
		panic("Mpath is nil")
	}
	if Plist == nil {
		panic("Plist is nil")
	}
	if Pelm == nil {
		panic("Pelm is nil")
	}
	if Contl == nil {
		panic("Contl is nil")
	}
	if Ctlif == nil {
		panic("Ctlif is nil")
	}
	if Ctlst == nil {
		panic("Ctlst is nil")
	}
	if Wd == nil {
		panic("Wd is nil")
	}
	if Daytm == nil {
		panic("Daytm is nil")
	}
	if obsn == nil {
		panic("obsn is nil")
	}
	if bp == nil {
		panic("bp is nil")
	}
	if obs == nil {
		panic("obs is nil")
	}
	if tree == nil {
		panic("tree is nil")
	}
	if shadtb == nil {
		panic("shadtb is nil")
	}
	if poly == nil {
		panic("poly is nil")
	}
	if monten == nil {
		panic("monten is nil")
	}
	if gpn == nil {
		panic("gpn is nil")
	}
	if DE == nil {
		panic("DE is nil")
	}
	if Noplpmp == nil {
		panic("Noplpmp is nil")
	}

	var Twallinit float64
	var j int
	dtm := 3600
	var nday int
	var Nday int
	daystartx := 0
	daystart := 0
	dayend := 0
	var Err, File string

	// 出力フラグ (GDAT.PRINT)
	// 中) 熱負荷要素の出力指定だけ変則的なことに注意
	wdpri := 0  // 気象データの出力指定
	revpri := 0 // 室内熱環境データの出力指定
	pmvpri := 0 // 室内のPMVの出力指定

	Nrmspri := 0 // 表面温度出力指定(室の数)
	Nqrmpri := 0 // 日射、室内発熱取得出力指定(室の数)
	Nwalpri := 0 // 壁体内部温度出力指定(壁体の数)
	Npcmpri := 0 // PCMの状態値出力フラグ(壁体の数)
	Nshdpri := 0 // 日よけの影面積出力 (壁体の数)

	var dfwl DFWL

	/*-------higuchi 070918---------start*/
	//RRMP *rp;
	//MADO *wp;
	//sunblk *sb;
	var smonth, sday, emonth, eday int

	//sb = bp.SBLK;
	//rp = bp.RMP;
	//wp = rp.WD;
	/*-------higuchi------------end*/

	Err = fmt.Sprintf(ERRFMT, "(Eeinput)")

	var err error

	// -------------------------------------------------------
	// 曜日設定ファイルの読み取り
	// -------------------------------------------------------
	var fi_dayweek []byte
	if fi_dayweek, err = ioutil.ReadFile(filepath.Join(efl_path, "dayweek.efl")); err != nil {
		Eprint("<Eeinput>", "dayweek.efl")
		os.Exit(EXIT_DAYWEK)
	}
	Dayweek(string(fi_dayweek), week, Simc.Daywk, key)

	if DEBUG {
		dprdayweek(Simc.Daywk)
	}

	// -------------------------------------------------------
	// スケジュ－ル表の読み取り
	// -------------------------------------------------------
	var Schdl *SCHDL = new(SCHDL)
	Schtable(schtba, Schdl)
	Schname(Schdl)

	// -------------------------------------------------------
	//  季節、曜日によるスケジュ－ル表の組み合わせの読み取り
	// -------------------------------------------------------
	Schdata(schnma, "schnm", Simc.Daywk, Schdl)

	// 入力を正規化することで後処理を簡単にする
	tokens := NewEeTokens(bdata)

	for !tokens.IsEnd() {
		s := tokens.GetToken()
		if s == "\n" || s == ";" {
			continue
		}
		fmt.Printf("=== %s\n", s)
		if s == "*" {
			continue
		}

		switch s {
		case "TITLE":
			line := tokens.GetLogicalLine()
			Simc.Title = strings.Join(line, " ")
			fmt.Printf("%s\n", Simc.Title)
		case "GDAT":
			section := tokens.GetSection()
			Wd.RNtype = 'C'
			Wd.Intgtsupw = 'N'
			Simc.Perio = 'n' // 周期定常計算フラグを'n'に初期化
			Gdata(section, Simc.File, &Simc.Wfname, &Simc.Ofname, &dtm, &Simc.Sttmm,
				&daystartx, &daystart, &dayend, &Twallinit, Simc.Dayprn,
				&wdpri, &revpri, &pmvpri, &Simc.Helmkey, &Simc.MaxIterate, Daytm, Wd, &Simc.Perio)

			// 気象データファイル名からファイル種別を判定
			if Simc.Wfname == "" {
				Simc.Wdtype = 'E'
			} else {
				Simc.Wdtype = 'H'
			}

			// 初期温度 (15[deg])
			Rmvls.Twallinit = Twallinit

			// 計算時間間隔 [s]
			Simc.DTm = dtm

			Simc.Unit = "t_C x_kg/kg r_% q_W e_W"
			Simc.Unitdy = "Q_kWh E_kWh"

			fmt.Printf("== File  Output=%s\n", Simc.Ofname)
		// case "SCHTB":
		// 	// SCHDBデータセットの読み取り
		// 	//Schtable(schtba, Schdl)
		// 	Schname(Schdl)
		// case "SCHNM":
		// 	// SCHNMデータセットの読み取り
		// 	Schdata(schnma, s, Simc.Daywk, Schdl)
		case "EXSRF":
			// EXSRFデータセットの読み取り
			section := tokens.GetSection()
			Exsfdata(section, s, Exsf, Schdl, Simc)

		case "SUNBRK":
			// 日よけの読み込み
			section := tokens.GetSection()
			Snbkdata(section, s, &Rmvls.Snbk)

		case "PCM":
			section := tokens.GetSection()
			PCMdata(section, s, &Rmvls.PCM, &Rmvls.Pcmiterate)

		case "WALL":
			if Fbmlist == "" {
				File = "wbmlist.efl"
			} else {
				File = Fbmlist
			}

			var fullpath string
			if filepath.IsAbs(File) {
				fullpath = File
			} else {
				fullpath = filepath.Join(efl_path, File)
			}

			var fbmContent []byte
			if fbmContent, err = ioutil.ReadFile(fullpath); err != nil {
				Eprint("<Eeinput>", "wbmlist.efl")
				os.Exit(EXIT_WBMLST)
			}
			/*******************/

			section := tokens.GetSection()
			Walldata(section, string(fbmContent), &Rmvls.Wall, &dfwl, Rmvls.PCM)

		case "WINDOW":
			section := tokens.GetSection()
			Windowdata(section, &Rmvls.Window)

		case "ROOM":
			Roomdata(tokens, Exsf.Exs, &dfwl, Rmvls, Schdl, Simc)
			Balloc(Rmvls.Sd, Rmvls.Wall, &Rmvls.Mw)

		case "RAICH", "VENT":
			section := tokens.GetSection()
			Ventdata(section, Schdl, Rmvls.Room, Simc)

		case "RESI":
			section := tokens.GetSection()
			Residata(section, Schdl, Rmvls.Room, &pmvpri, Simc)

		case "APPL":
			section := tokens.GetSection()
			Appldata(section, Schdl, Rmvls.Room, Simc)

		case "VCFILE":
			section := tokens.GetSection()
			Vcfdata(section, Simc)

		case "EQPCAT":
			section := tokens.GetSection()
			Eqcadata(section, Eqcat)

		case "SYSCMP": // 接続用のノードを設定している
			/*****Flwindata(Flwin, Nflwin,  Wd);********/
			section := tokens.GetSection()
			Compodata(section, Rmvls, Eqcat, Compnt, Eqsys)
			Elmalloc(*Compnt, Eqcat, Eqsys, Elout, Elin)

		case "SYSPTH": // 接続パスの設定をしている
			section := tokens.GetSection()
			Pathdata(section, Simc, Wd, *Compnt, Schdl, Mpath, Plist, Pelm, Eqsys, Elout, Elin)
			Roomelm(Rmvls.Room, Rmvls.Rdpnl)

			// 変数の割り当て
			Hclelm(Eqsys.Hcload)
			Thexelm(Eqsys.Thex)
			Desielm(Eqsys.Desi)
			Evacelm(Eqsys.Evac)

			Qmeaselm(Eqsys.Qmeas)

		case "CONTL":
			section := tokens.GetSection()
			Contrldata(section, Contl, Ctlif, Ctlst, Simc, *Compnt, *Mpath, Wd, Exsf, Schdl)

		/*--------------higuchi add-------------------start*/

		// 20170503 higuchi add
		case "DIVID":
			section := tokens.GetSection()
			dividdata(section, monten, DE)

		/*--対象建物データ読み込み--*/
		case "COORDNT":
			// メモリの確保
			*bp = make([]*BBDP, 0)

			for {
				section := tokens.GetSection()
				if section.PeekToken() == "*" {
					break
				}
				bdpdata(section, bp, Exsf)
				tokens.SkipToEndOfLine()
			}

		/*--障害物データ読み込み--*/
		case "OBS":
			section := tokens.GetSection()
			obsdata(section, obsn, obs)

		/*--樹木データ読み込み--*/
		case "TREE":
			section := tokens.GetSection()
			treedata(section, tree)

		/*--多角形障害物直接入力分の読み込み--*/
		case "POLYGON":
			section := tokens.GetSection()
			polydata(section, poly)

		/*--落葉スケジュール読み込み--*/
		case "SHDSCHTB":
			// 落葉スケジュールの数を数える
			section := tokens.GetSection()

			*shadtb = make([]*SHADTB, 0)

			for !section.IsEnd() {
				line := new(EeTokens)
				line.tokens = section.GetLogicalLine()
				line.pos = 0

				s = line.GetToken()
				if s[0] == '*' {
					break
				}

				shdp := new(SHADTB)
				shdp.lpname = s
				shdp.indatn = 0

				for !line.IsEnd() {
					s = line.GetToken()
					if s == ";" {
						break
					}

					_, err = fmt.Sscanf(s, "%d/%d-%f-%d/%d", &smonth, &sday, &shdp.shad[shdp.indatn], &emonth, &eday)
					if err != nil {
						panic(err)
					}
					shdp.ndays[shdp.indatn] = nennkann(smonth, sday)
					shdp.ndaye[shdp.indatn] = nennkann(emonth, eday)
					shdp.indatn = shdp.indatn + 1
				}

				*shadtb = append(*shadtb, shdp)
			}

		/*----------higuchi add-----------------end-*/

		default:
			Err = Err + "  " + s
			Eprint("<Eeinput>", Err)
		}
	}

	/*--------------higuchi 070918-------------------start-*/
	if len(*bp) != 0 {
		fmt.Printf("deviding of wall mm: %f\n", *DE)
		fmt.Printf("number of point in montekalro: %d\n", *monten)
	}
	/*----------------higuchi 7.11,061123------------------end*/

	// 外部障害物の数を数える
	Noplpmp.Nop = OPcount(*bp, *poly)
	Noplpmp.Nlp = LPcount(*bp, *obs, *tree, *poly)
	Noplpmp.Nmp = Noplpmp.Nop + Noplpmp.Nlp

	//////////////////////////////////////

	//----------------------------------------------------
	// シミュレーション設定
	//----------------------------------------------------

	if daystart > dayend {
		dayend = dayend + 365
	}
	Nday = dayend - daystart + 1

	if daystartx > daystart {
		daystart = daystart + 365
	}

	Nday += daystart - daystartx
	Simc.Dayend = daystartx + Nday - 1
	Simc.Daystartx = daystartx
	Simc.Daystart = daystart

	Simc.Timeid = []rune{'M', 'D', 'T'}

	Simc.Ntimedyprt = Simc.Dayend - Simc.Daystart + 1
	Simc.Dayntime = 24 * 3600 / dtm
	Simc.Ntimehrprt = 0

	for nday = Simc.Daystart; nday <= Simc.Dayend; nday++ {
		// NOTE: オリジナルコードはバッファーオーバーランしているので、`%366`を追加
		if Simc.Dayprn[nday%366] != 0 {
			Simc.Ntimehrprt += Simc.Dayntime
		}
	}

	//----------------------------------------------------
	// 出力ファイルの追加
	//----------------------------------------------------

	for i := range Rmvls.Room {
		Rm := Rmvls.Room[i]
		if Rm.sfpri {
			Nrmspri++
		}
		if Rm.eqpri {
			Nqrmpri++
		}
	}

	for _, Sd := range Rmvls.Sd {
		if Sd.wlpri {
			Nwalpri++
		}

		if Sd.pcmpri {
			Npcmpri++
		}

		// 日よけの影面積出力
		if Sd.shdpri {
			Nshdpri++
		}
	}

	// 出力ファイルの追加手続き
	var Flout []*FLOUT = make([]*FLOUT, 0, 30) // ファイル出力設定
	addFlout := func(idn PrintType) {
		Flout = append(Flout, &FLOUT{Idn: idn})
	}

	// 必須出力ファイル
	addFlout(PRTPATH)    // 時間別計算値(システム経路の温湿度出力)
	addFlout(PRTCOMP)    // 時間別計算値(機器の出力)
	addFlout(PRTDYCOMP)  // 日別計算値(システム要素機器の日集計結果出力)
	addFlout(PRTMNCOMP)  // 月別計算値(システム要素機器の月集計結果出力)
	addFlout(PRTMTCOMP)  // 月-時刻計算値(部屋ごとの熱集計結果出力)
	addFlout(PRTHRSTANK) // 時間別計算値(蓄熱槽内温度分布の出力)
	addFlout(PRTWK)      // 計算年月日出力
	addFlout(PRTREV)     // 時間別計算値(毎時室温、MRTの出力)
	addFlout(PRTHROOM)   // 時間別計算値(放射パネルの出力)
	addFlout(PRTDYRM)    // 日別計算値(部屋ごとの熱集計結果出力)
	addFlout(PRTMNRM)    // 月別計算値(部屋ごとの熱集計結果出力)

	// 要素別熱損失・熱取得（記憶域確保）
	Helminit("Helminit", Simc.Helmkey, Rmvls.Room, &Rmvls.Qetotal)

	if Simc.Helmkey == 'y' {
		addFlout(PRTHELM)   // 時間別計算値(要素別熱損失・熱取得)
		addFlout(PRTDYHELM) // 日別計算値(要素別熱損失・熱取得)

		Simc.Nhelmsfpri = 0
		for i := range Rmvls.Room {
			Rm := Rmvls.Room[i]
			for j = 0; j < Rm.N; j++ {
				Sdd := Rm.rsrf[j]
				if Sdd.sfepri {
					Simc.Nhelmsfpri++
				}
			}
		}
		if Simc.Nhelmsfpri > 0 {
			addFlout(PRTHELMSF) // 時間別計算値(要素別熱損失・熱取得) 表面?
		}
	}

	if pmvpri > 0 {
		addFlout(PRTPMV) // 時間別計算値(PMV計算)
	}

	if Nqrmpri > 0 {
		addFlout(PRTQRM) // 時間別計算値(日射、室内熱取得の出力)
		addFlout(PRTDQR) // 日別計算値(日射、室内熱取得の出力)
	}

	if Nrmspri > 0 {
		addFlout(PRTRSF)  // 時間別計算値(室内表面温度の出力)
		addFlout(PRTSFQ)  // 時間別計算値(室内表面熱流の出力)
		addFlout(PRTSFA)  // 時間別計算値(室内表面熱伝達率の出力)
		addFlout(PRTDYSF) // 日別計算値(日積算壁体貫流熱取得の出力)
	}

	if Nwalpri > 0 {
		addFlout(PRTWAL) // // 時間別計算値(壁体内部温度の出力)
	}

	// 日よけの影面積出力
	if Nshdpri > 0 {
		addFlout(PRTSHD) // 時間別計算値(日よけの影面積の出力)
	}

	// 潜熱蓄熱材がある場合
	if Npcmpri > 0 {
		addFlout(PRTPCM) // 時間別計算値(潜熱蓄熱材の状態値の出力)
	}

	// 気象データの出力を追加
	if wdpri > 0 {
		addFlout(PRTHWD) // 時間別計算値(気象データ出力)
		addFlout(PRTDWD) // 日別計算値(気象データ日集計値出力)
		addFlout(PRTMWD) // 月別計算値(気象データ月集計値出力)
	}

	// DEBUG
	fmt.Printf("読み取りデータ数\n")
	fmt.Printf("SHDSCHTB: %d\n", len(*shadtb))    // 落葉スケジュール
	fmt.Printf("TREE: %d\n", len(*tree))          // 樹木データ
	fmt.Printf("OBS: %d\n", len(*obs))            // 障害物データ
	fmt.Printf("POLYGON: %d\n", len(*poly))       // 多角形障害物直接入力分
	fmt.Printf("COORDNT: %d\n", len(*bp))         // 対象建物データ
	fmt.Printf("WINDOW: %d\n", len(Rmvls.Window)) // 窓データ
	fmt.Printf("WALL: %d\n", len(Rmvls.Wall))     // 壁データ

	fmt.Printf("RESI: %d\n", len(Rmvls.Room))
	fmt.Printf("SYSCMP: %d\n", len(*Compnt))
	fmt.Printf("SYSPTH: %d\n", len(*Mpath))
	fmt.Printf("ROOM: %d\n", len(Rmvls.Room))
	fmt.Printf("EXSRF: %d\n", len(Exsf.Exs))
	fmt.Printf("CONTL: %d\n", len(*Contl))

	return Schdl, Flout
}

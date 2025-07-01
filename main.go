// EESLISM Go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/akamensky/argparse"
	eeslism "github.com/archlabjp/eeslism-go/eeslism"
)

/*
main (Main Entry Point of EESLISM Go)

この関数は、EESLISM Goプログラムのメインエントリポイントであり、
コマンドライン引数の解析、初期設定、
そして建物のエネルギーシミュレーションの実行を統括します。

建築環境工学的な観点:
- **コマンドライン引数の処理**: 
  - `filename`: 入力データファイル名。建物の形状、材料特性、設備機器、スケジュールなど、
    シミュレーションに必要な全ての情報が記述されたファイルを指定します。
  - `efl_path`: EFL（Energy Flow Language）ファイルのディレクトリ。気象データファイルや、
    機器のカタログデータなどが格納されたディレクトリを指定します。
  これらの引数は、シミュレーションの入力条件を定義し、
  様々な建物のエネルギー性能を評価するための柔軟性を提供します。
- **シミュレーションの実行**: `eeslism.Entry(*filename, *efl_path)` を呼び出すことで、
  実際のエネルギーシミュレーションが開始されます。
  `eeslism.Entry`関数は、入力データの読み込み、モデルの初期化、
  時間ステップごとの計算ループ、そして結果の出力といった一連のプロセスを統括します。
- **ログ出力**: `log.SetFlags(log.Lmicroseconds)` は、
  ログメッセージにマイクロ秒単位のタイムスタンプを含める設定です。
  これは、シミュレーションの実行時間や、
  特定のイベントの発生時刻を詳細に追跡し、
  デバッグや性能分析に役立ちます。

この関数は、建物のエネルギーシミュレーションプログラムの起動と実行を制御し、
ユーザーが指定した条件に基づいて、
建物のエネルギー性能を評価するための重要な役割を果たします。
*/
func main() {
	log.SetFlags(log.Lmicroseconds)

	// コマンドライン引数の処理
	parser := argparse.NewParser("EESLISIM Go", "a general-purpose simulation program for building thermal-environmental control systems consisting of both buildings and facilities")

	filename := parser.StringPositional(&argparse.Options{
		Required: true,
		Help:     "Input data file name"})

	efl_path := parser.String("", "efl", &argparse.Options{
		Default: "Base",
		Help:    "EFLファイルのディレクトリ"})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
	}

	// if len(*efl_path) > 0 {
	// 	os.Chdir(*efl_path)
	// }

	eeslism.Entry(*filename, *efl_path)
}

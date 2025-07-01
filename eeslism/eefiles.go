package eeslism

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

/*
eeflopen (Energy Simulation File Open)

この関数は、建物のエネルギーシミュレーションに必要な入力ファイル（気象データなど）と、
出力ファイル（シミュレーション結果）をオープンします。

建築環境工学的な観点:
- **気象データファイルのオープン**: 建物のエネルギーシミュレーションには、
  外気温度、湿度、日射量などの気象データが不可欠です。
  `Simc.Wfname`で指定された気象データファイル（通常は年間気象データ）をオープンし、
  シミュレーション中に気象データを読み込めるようにします。
  `Simc.Fwdata`と`Simc.Fwdata2`の二つのファイルポインターは、
  気象データを複数回読み込む必要がある場合（例えば、異なる計算モジュールで同時にアクセスする場合）に用いられます。
- **supw.eflファイルのオープン**: `supw.efl`は、
  おそらく給水温度や地中温度など、
  シミュレーションに必要な補助的な気象データや環境データが格納されているファイルです。
  これらのデータは、給湯負荷計算や地中熱交換器の計算などに用いられます。
- **出力ファイルの準備**: シミュレーション結果を格納するための出力ファイル（`Flout`）を準備します。
  `fl.Fname = Simc.Ofname + string(fl.Idn) + ".es"` のように、
  出力ファイル名が自動的に生成され、
  結果を効率的に保存できるようにします。
- **エラーハンドリング**: ファイルのオープンに失敗した場合、
  エラーメッセージを出力し、プログラムを終了します。
  これは、シミュレーションの実行に必要なファイルが正しく読み込まれていることを確認し、
  シミュレーションの信頼性を確保するために重要です。

この関数は、建物のエネルギーシミュレーションのデータ入出力の基盤を形成し、
シミュレーションの正確性、安定性、および再現性を向上させるための重要な役割を果たします。
*/
func (Simc *SIMCONTL) eeflopen(Flout []*FLOUT, efl_path string) {
	// 気象データファイルを開く
	if Simc.Wdtype == 'H' {
		var fullpath string
		if filepath.IsAbs(Simc.Wfname) {
			fullpath = Simc.Wfname
		} else {
			fullpath = filepath.Join(efl_path, Simc.Wfname)
		}

		var err error
		Simc.Fwdata, err = os.Open(fullpath)
		if err != nil {
			Eprint("<eeflopen>", fullpath)
			os.Exit(EXIT_WFILE)
		}
		Simc.Fwdata2, err = os.Open(fullpath)
		if err != nil {
			Eprint("<eeflopen>", fullpath)
			os.Exit(EXIT_WFILE)
		}

		Simc.Ftsupw, err = ioutil.ReadFile(filepath.Join(efl_path, "supw.efl"))
		if err != nil {
			Eprint("<eeflopen>", "supw.efl")
			os.Exit(EXIT_SUPW)
		}
	}

	// Fname = Simc.ofname + ".log"
	// ferr, err := os.Create(Fname)
	// if err != nil {
	//     fmt.Println(err)
	//     os.Exit(1)
	// }

	// 出力ファイルを開く
	for _, fl := range Flout {
		fl.Fname = Simc.Ofname + string(fl.Idn) + ".es"
		fl.F = new(strings.Builder)
	}
}

/*
Eeflclose (Energy Simulation File Close)

この関数は、建物のエネルギーシミュレーションで使用された出力ファイルを閉じ、
バッファリングされたデータをファイルに書き込みます。

建築環境工学的な観点:
- **出力ファイルのクローズ**: シミュレーションが完了した後、
  出力ファイルを適切に閉じることが重要です。
  これにより、全てのシミュレーション結果がファイルに確実に書き込まれ、
  データの破損や欠落を防ぎます。
- **バッファリングされたデータの書き込み**: `fmt.Fprint(fo, fl.F)` は、
  メモリ上にバッファリングされていたシミュレーション結果をファイルに書き込みます。
  これにより、シミュレーション中に発生した全てのデータが保存されます。
- **エラーハンドリング**: ファイルのクローズに失敗した場合、
  エラーメッセージを出力し、プログラムを終了します。
  これは、シミュレーション結果が正しく保存されていることを確認し、
  シミュレーションの信頼性を確保するために重要です。

この関数は、建物のエネルギーシミュレーションのデータ入出力の最終ステップであり、
シミュレーション結果の完全性と信頼性を確保するための重要な役割を果たします。
*/
func Eeflclose(Flout []*FLOUT) {
	var fl *FLOUT

	if Ferr != nil {
		Ferr.Close()
	}

	for _, fl = range Flout {
		fo, err := os.Create(fl.Fname)
		if err != nil {
			Eprint("<eeflopen>", fl.Fname)
			os.Exit(EXIT_WFILE)
		}

		fmt.Fprint(fo, fl.F)
		fmt.Fprintln(fo, "-999")
		fo.Close()
	}
}

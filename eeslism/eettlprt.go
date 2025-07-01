package eeslism

import (
	"fmt"
	"io"
	"strings"
)

/* 標題、注記の出力（時刻別計算値ファイル） */

/*
__replace_dir_sep (Replace Directory Separator)

この関数は、ファイルパス内のディレクトリ区切り文字を統一します。

建築環境工学的な観点:
- **クロスプラットフォーム対応**: オペレーティングシステムによってディレクトリ区切り文字が異なる（Windowsでは`\`、Unix系では`/`）ため、
  この関数はパスを統一的な形式に変換することで、
  シミュレーション結果の出力ファイルが異なる環境でも正しく読み込めるようにします。

この関数は、建物のエネルギーシミュレーションにおいて、
出力データの互換性と可搬性を向上させるための補助的な役割を果たします。
*/
func __replace_dir_sep(path *string) {
	// Replace directory separator to unify them
	*path = strings.Replace(*path, "/", "\\", -1)
}

/*
ttlprint (Title Print for Hourly Output)

この関数は、建物のエネルギーシミュレーションの時刻別計算結果ファイルに、
ヘッダー情報（タイトル、バージョン、入力ファイル名、気象ファイル名、時間ID、単位、時間間隔、データ数など）を出力します。

建築環境工学的な観点:
- **シミュレーション結果のメタデータ**: シミュレーション結果を正確に解釈し、
  他のシミュレーションと比較するためには、
  その結果がどのような条件で得られたものかを示すメタデータが不可欠です。
  この関数は、以下の情報を提供します。
  - `fileid`: ファイルの識別子。
  - `EEVERSION`: シミュレーションプログラムのバージョン。
  - `simc.Title`: シミュレーションのタイトル。
  - `simc.File`: 入力データファイル名。
  - `simc.Wfname`: 気象データファイル名。
  - `simc.Timeid`: 時間ID（時、日、月など）。
  - `simc.Unit`: 出力データの単位。
  - `simc.DTm`: 計算時間間隔。
  - `simc.Ntimehrprt`: 時間別データ数。
- **出力ファイルの可読性**: ヘッダー情報を適切に記述することで、
  出力ファイルが人間にとっても、
  他の解析ツールにとっても読みやすくなります。

この関数は、建物のエネルギーシミュレーションにおいて、
出力データの品質と再現性を向上させるための重要な役割を果たします。
*/
func ttlprint(fo io.Writer, fileid string, simc *SIMCONTL) {
	// Replace directory separator to unify them
	__replace_dir_sep(&simc.File)

	fmt.Fprintf(fo, "%s#\n", fileid)
	fmt.Fprintf(fo, "-ver %s\n", EEVERSION)
	fmt.Fprintf(fo, "-t %s ;\n", simc.Title)
	fmt.Fprintf(fo, "-dtf %s\n", simc.File)
	fmt.Fprintf(fo, "-w %s\n", simc.Wfname)
	fmt.Fprint(fo, "-tid h\n")

	fmt.Fprint(fo, "-tmid ")
	for i := 0; i < len(simc.Timeid); i++ {
		fmt.Fprint(fo, string(simc.Timeid[i]))
	}
	fmt.Fprint(fo, "\n")

	fmt.Fprintf(fo, "-u %s ;\n", simc.Unit)
	fmt.Fprintf(fo, "-dtm %d\n", simc.DTm)
	fmt.Fprintf(fo, "-Ntime %d\n", simc.Ntimehrprt)
}

/* ---------------------------------------------------- */

/*
ttldyprint (Title Print for Daily Output)

この関数は、建物のエネルギーシミュレーションの日集計計算結果ファイルに、
ヘッダー情報（タイトル、バージョン、入力ファイル名、気象ファイル名、時間ID、単位、時間間隔、データ数など）を出力します。

建築環境工学的な観点:
- **シミュレーション結果のメタデータ**: シミュレーション結果を正確に解釈し、
  他のシミュレーションと比較するためには、
  その結果がどのような条件で得られたものかを示すメタデータが不可欠です。
  この関数は、以下の情報を提供します。
  - `fileid`: ファイルの識別子。
  - `EEVERSION`: シミュレーションプログラムのバージョン。
  - `simc.Title`: シミュレーションのタイトル。
  - `simc.File`: 入力データファイル名。
  - `simc.Wfname`: 気象データファイル名。
  - `simc.Timeid`: 時間ID（時、日、月など）。
  - `simc.Unit`, `simc.Unitdy`: 出力データの単位（時間別と日集計で異なる場合がある）。
  - `simc.DTm`: 計算時間間隔。
  - `simc.Ntimedyprt`: 日集計データ数。
- **出力ファイルの可読性**: ヘッダー情報を適切に記述することで、
  出力ファイルが人間にとっても、
  他の解析ツールにとっても読みやすくなります。

この関数は、建物のエネルギーシミュレーションにおいて、
出力データの品質と再現性を向上させるための重要な役割を果たします。
*/
func ttldyprint(fo io.Writer, fileid string, simc *SIMCONTL) {
	fmt.Fprintf(fo, "%s#\n", fileid)
	fmt.Fprintf(fo, "-ver %s\n", EEVERSION)
	fmt.Fprintf(fo, "-t %s ;\n", simc.Title)
	fmt.Fprintf(fo, "-dtf %s\n", simc.File)
	fmt.Fprintf(fo, "-w %s\n", simc.Wfname)
	fmt.Fprint(fo, "-tid d\n")

	fmt.Fprint(fo, "-tmid ")
	for i := 0; i < len(simc.Timeid)-1; i++ {
		fmt.Fprint(fo, string(simc.Timeid[i]))
	}
	fmt.Fprint(fo, "\n")

	fmt.Fprintf(fo, "-u %s %s ;\n", simc.Unit, simc.Unitdy)
	fmt.Fprintf(fo, "-dtm %d\n", simc.DTm)
	fmt.Fprintf(fo, "-Ntime %d\n", simc.Ntimedyprt)
}

/* ---------------------------------------------------- */

/*
ttlmtprint (Title Print for Monthly-Time-of-Day Output)

この関数は、建物のエネルギーシミュレーションの月・時刻別計算結果ファイルに、
ヘッダー情報（タイトル、バージョン、入力ファイル名、気象ファイル名、時間ID、単位、時間間隔、データ数など）を出力します。

建築環境工学的な観点:
- **シミュレーション結果のメタデータ**: シミュレーション結果を正確に解釈し、
  他のシミュレーションと比較するためには、
  その結果がどのような条件で得られたものかを示すメタデータが不可欠です。
  この関数は、以下の情報を提供します。
  - `fileid`: ファイルの識別子。
  - `EEVERSION`: シミュレーションプログラムのバージョン。
  - `simc.Title`: シミュレーションのタイトル。
  - `simc.File`: 入力データファイル名。
  - `simc.Wfname`: 気象データファイル名。
  - `simc.Timeid`: 時間ID（時、日、月など）。
  - `simc.Unit`, `simc.Unitdy`: 出力データの単位（時間別と日集計で異なる場合がある）。
  - `simc.DTm`: 計算時間間隔。
  - `Ntime`: 月・時刻別データ数（通常は24時間×12ヶ月）。
- **出力ファイルの可読性**: ヘッダー情報を適切に記述することで、
  出力ファイルが人間にとっても、
  他の解析ツールにとっても読みやすくなります。

この関数は、建物のエネルギーシミュレーションにおいて、
出力データの品質と再現性を向上させるための重要な役割を果たします。
*/
func ttlmtprint(fo io.Writer, fileid string, simc *SIMCONTL) {
	fmt.Fprintf(fo, "%s#\n", fileid)
	fmt.Fprintf(fo, "-ver %s\n", EEVERSION)
	fmt.Fprintf(fo, "-t %s ;\n", simc.Title)
	fmt.Fprintf(fo, "-dtf %s\n", simc.File)

	fmt.Fprintf(fo, "-w %s\n", simc.Wfname)
	fmt.Fprint(fo, "-tid h\n")

	fmt.Fprint(fo, "-tmid MT\n")

	fmt.Fprintf(fo, "-u %s %s ;\n", simc.Unit, simc.Unitdy)
	fmt.Fprintf(fo, "-dtm %d\n", simc.DTm)
	fmt.Fprint(fo, "-Ntime 288\n") // 24 * 12
}

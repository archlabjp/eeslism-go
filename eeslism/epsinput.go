package eeslism

import (
	"bufio"

	"fmt"
	"io"
	"strconv"
	"strings"
)

/*
esondat (Energy Simulation Output Data Input)

この関数は、建物のエネルギーシミュレーション結果ファイルから、
ヘッダー情報（タイトル、気象データファイル名、時間ID、単位、データ数など）を読み込み、
対応する構造体（`ESTL`）に格納します。

建築環境工学的な観点:
  - **シミュレーション結果のメタデータ**: シミュレーション結果を正確に解釈し、
    他のシミュレーションと比較するためには、
    その結果がどのような条件で得られたものかを示すメタデータが不可欠です。
    この関数は、以下の情報を提供します。
  - `Estl.Title`: シミュレーションのタイトル。
  - `Estl.Wdatfile`: 使用した気象データファイル名。
  - `Estl.Tid`: 入力データ種別（時刻別、日別など）。
  - `Estl.Unit`: 出力データの単位。
  - `Estl.Ntime`, `Estl.dtm`: データ数、計算時間間隔。
  - `Estl.Timeid`: 時刻データの表示形式。
  - `Estl.Flid`: ファイルの識別子。
  - `Estl.Catnm`: 要素カタログ名データ（機器数、データ項目数など）。
  - `Estl.Wdloc`: 地域情報（地名、緯度、経度など）。
  - **出力ファイルの解析**: シミュレーション結果ファイルは、
    通常、特定のフォーマットで記述されています。
    この関数は、そのフォーマットを解析し、
    必要な情報を抽出します。
  - **データ構造の準備**: 読み込んだメタデータに基づいて、
    シミュレーション結果のデータ構造（`TLIST`など）を準備します。

この関数は、建物のエネルギーシミュレーション結果を分析し、
省エネルギー設計、快適性評価、
および運用改善のための意思決定を支援するための重要な役割を果たします。
*/
func esondat(fi io.Reader, Estl *ESTL) {
	var s string
	var i, j, Nparm, Ndat int
	var catnm, C *CATNM

	Estl.Catnm = nil

	for {
		_, err := fmt.Fscanf(fi, "%s", &s)
		if err != nil {
			break
		}

		switch s {
		case "-t":
			fmt.Fscanf(fi, " %[^;];", &s)
			Estl.Title = s
		case "-w":
			fmt.Fscanf(fi, "%s", &s)
			Estl.Wdatfile = s
		case "-tid":
			// 時刻別データであることの指定
			fmt.Fscanf(fi, " %c", &Estl.Tid)
		case "-u":
			i := 0
			for {
				_, err := fmt.Fscanf(fi, "%s", &s)
				if err != nil || s == ";" {
					break
				}
				Estl.Unit[i] = s
				i++
				Estl.Nunit = i
			}
		case "-Ntime":
			// 項目ごとの全データ数
			fmt.Fscanf(fi, " %d", &Estl.Ntime)
		case "-dtm":
			// 時間間隔 [s]
			fmt.Fscanf(fi, " %d", &Estl.dtm)
		case "-tmid":
			// 時刻別データであることの指定
			fmt.Fscanf(fi, "%s", &s)
			Estl.Timeid = s
			Estl.Ntimeid = len(Estl.Timeid)
		case "-cat":
			N := CATNMMAX
			if N > 0 {
				Estl.Catnm = make([]CATNM, N)
			}

			if Estl.Catnm != nil {
				for ss := 0; ss < N; ss++ {
					C = &Estl.Catnm[ss]
					C.Name = ""
					C.N = 0
					C.Ncdata = 0
				}
			}

			Estl.Ndata = 0
			catidx := 0
			for {
				catnm = &Estl.Catnm[catidx]

				_, err := fmt.Fscanf(fi, "%s", &s)
				if err != nil || s == "*" {
					break
				}
				catnm.Name = s
				catnm.Ncdata = 0
				fmt.Fscanf(fi, "%d", &catnm.N)
				for i = 0; i < catnm.N; i++ {
					fmt.Fscanf(fi, "%s %d %d", &s, &Nparm, &Ndat)
					for j = 0; j < Nparm-1; j++ {
						fmt.Fscanf(fi, "%s", &s)
					}
					Estl.Ndata += Ndat
					catnm.Ncdata += Ndat
				}
				catidx++
			}
		case "-wdloc":
			// 地名
			fmt.Fscanf(fi, "%[^;];", &s)
			s += " ;"
			Estl.Wdloc = s
		case "-Ndata":
			// 各時刻ごとのデータ個数
			fmt.Fscanf(fi, " %d", &Estl.Ndata)
		default:
			if s[len(s)-1] == '#' {
				Estl.Flid = s
			} else {
				Eprint("<esondat>", s)
			}
		}
	}

	if Estl.Title != "" {
		fmt.Printf("esondat  title=%s\n", Estl.Title)
	}

	if Estl.Title != "" {
		fmt.Printf("esondat  w=%s\n", Estl.Wdatfile)
	}

	fmt.Printf("esondat  tid=%c\n", Estl.Tid)
	fmt.Printf("esondat  Ntime=%d\n", Estl.Ntime)
	fmt.Printf("esondat  tmdt=%s\n", Estl.Timeid)
	fmt.Printf("esondat  flid=%s\n", Estl.Flid)
}

/* ----------------------------------------------------------- */

/*
esoint (Energy Simulation Output Data Initialization)

この関数は、シミュレーション結果の各項目（温度、熱量、エネルギーなど）のデータ構造（`TLIST`）を初期化し、
出力形式や集計方法を設定します。

建築環境工学的な観点:
  - **シミュレーション結果の項目定義**: シミュレーション結果は、
    様々な種類のデータ（温度、熱量、エネルギーなど）から構成されます。
    この関数は、各項目のデータ種別（`Vtype`）、データ処理種別（`Stype`）、
    データ型（`Ptype`）を設定し、
    それぞれの項目がどのように集計され、出力されるかを定義します。
  - **データ処理種別の自動判定**: `Tlist.Vtype`に基づいて、
    データ処理種別（`Tlist.Stype`）を自動的に判定します。
    例えば、`H`, `Q`, `E`（積算値）の場合は`'t'`（積算）、
    `t`, `x`, `r`（瞬時値）の場合は`'a'`（平均）と設定します。
  - **出力形式の設定 (fofmt)**:
    `fofmt`関数を呼び出すことで、
    各項目の出力形式（小数点以下の桁数など）を設定します。
    これにより、出力ファイルの可読性を向上させます。
  - **選択項目の設定**: `Estl.Nrqlist`や`Estl.Nvreq`に基づいて、
    出力する項目を選択します。
    これにより、ユーザーが関心のある特定のデータのみを抽出できます。

この関数は、建物のエネルギーシミュレーション結果を分析し、
省エネルギー設計、快適性評価、
および運用改善のための意思決定を支援するための重要な役割を果たします。
*/
func esoint(fi io.Reader, err string, Ntime int, Estl *ESTL, _Tlist []TLIST) {
	var nm, id string
	var V *rune
	var st int
	var cat *CATNM
	var R *RQLIST
	var n int
	var catIdx = 0
	var rqIdx = 0

	cat = nil
	R = nil
	V = nil

	if Estl.Catnm != nil {
		cat = &Estl.Catnm[0]
	}
	// Rq = &Estl.Rq[0]

	for i := 0; i < Estl.Ndata; i++ {
		Tlist := &_Tlist[i]
		fmt.Fscanf(fi, " %[^_]_%s %c %c", &nm, &id, &Tlist.Vtype, &Tlist.Ptype)

		switch Tlist.Vtype {
		case 'H', 'Q', 'E', 'q', 'e', 'm':
			Tlist.Stype = 't'
		case 'T', 'X', 'R', 't', 'x', 'r':
			Tlist.Stype = 'a'
		case 'c':
			Tlist.Stype = 'c'
		default:
			switch id[len(id)-1] {
			case 'n', 'c':
				Tlist.Stype = 'n'
			case 'm', 'h', 'e', 'p':
				Tlist.Stype = 'm'
			default:
				s := fmt.Sprintf("xxxx %s xxx  %s %s %c %c %c\n", err, nm, id, id[len(id)-1], Tlist.Vtype, Tlist.Ptype)
				Eprint("<esoint>", s)
			}
		}

		if Estl.Catnm != nil {
			if n >= cat.Ncdata {
				catIdx++
				cat = &Estl.Catnm[catIdx]
				n = 0
			}
			Tlist.Cname = cat.Name
			n++
		} else {
			Tlist.Cname = "*"
		}

		Tlist.Name = nm
		Tlist.Id = id
		Tlist.Req = 'n'

		if Estl.Nrqlist == 0 && Estl.Nvreq == 0 {
			Tlist.Req = 'y'
		} else {
			R = &Estl.Rq[rqIdx]
			for j := 0; j < Estl.Nrqlist; j++ {
				if (Tlist.Name == R.Name || R.Name == "*") &&
					(Tlist.Id == R.Id || R.Id == "*") {
					Tlist.Req = 'y'
					break
				} else if st = strings.IndexRune(Tlist.Name, ':'); st != -1 {

					if Tlist.Name[:st] == R.Name[:st] && R.Id == "*" {
						Tlist.Req = 'y'
						break
					}
				}
			}

			for j := 0; j < Estl.Nvreq; j++ {
				V = &Estl.Vreq[j]
				if *V == Tlist.Vtype {
					Tlist.Req = 'y'
					break
				}
			}
		}

		switch Tlist.Ptype {
		case 'f':
			Tlist.Fval = make([]float64, Ntime)
		case 'd':
			Tlist.Ival = make([]int, Ntime)
		case 'c':
			Tlist.Cval = make([]rune, Ntime)
		}

		fofmt(Estl, Tlist)
	}
}

/*
fofmt (Format Output for Simulation Results)

この関数は、シミュレーション結果の各項目（温度、熱量、エネルギーなど）の
出力形式（小数点以下の桁数など）を設定します。

建築環境工学的な観点:
  - **出力の可読性向上**: シミュレーション結果を人間が読みやすく、
    また他の解析ツールで利用しやすいように、
    適切な数値の表示形式を設定します。
    例えば、温度は小数点以下1桁、絶対湿度は小数点以下4桁など、
    データの特性に応じたフォーマットを適用します。
  - **単位の表示**: `Tlist.Unit = Estl.Unit[2]` のように、
    各項目の単位を設定します。
    これにより、出力データの意味を明確にし、
    誤解を防ぎます。

この関数は、建物のエネルギーシミュレーションにおいて、
出力データの品質を向上させるための重要な役割を果たします。
*/
func fofmt(Estl *ESTL, Tlist *TLIST) {
	var fmt string

	switch Tlist.Vtype {
	case 't', 'T':
		// 温度は小数点以下１桁
		fmt = "%8.1f"
	case 'r', 'R':
		// 相対湿度は整数
		fmt = "%8.0f"
	case 'x', 'X':
		// 絶対湿度は小数点以下４桁
		fmt = "%8.4f"
	case 'q', 'e', 'Q', 'E':
		// 熱量は小数点以下1桁
		fmt = "%8.1f"
	case 'H':
		// 積算熱量は整数
		fmt = "%8d"
	case 'h':
		// 発生時刻は整数
		fmt = "%8d"
	case 'm':
		// 最小値は整数
		fmt = "%8d"
	case 'n':
		// 最大値は整数
		fmt = "%04d"
	case 'c':
		fmt = "%c"
	default:
		panic("TLIST.Vtype is wrong. Only 't', 'T', 'r', 'R', 'x', 'X', 'q', 'e', 'Q', 'E")
	}

	Tlist.Fmt = fmt

	for j := 0; j < len(Estl.Unit); j++ {
		if Estl.Unit[j][0] == byte(Tlist.Vtype) {
			Tlist.Unit = Estl.Unit[2]
		}
	}
}

/*
tmdata (Time Data Input for Simulation Results)

この関数は、シミュレーション結果ファイルから、
年、月、日、曜日、時刻などの時間データを読み込み、
対応する構造体（`TMDT`）に格納します。

建築環境工学的な観点:
  - **時間データの同期**: シミュレーション結果は、
    特定の時間軸に沿って生成されます。
    この関数は、結果ファイルから時間データを正確に読み込み、
    シミュレーションの時刻と同期させます。
  - **周期定常計算の考慮**: `perio`パラメータは、
    周期定常計算が行われているかどうかを示します。
    周期定常計算の場合、時間データの読み込み方法が異なる場合があります。
  - **データ読み込みの制御**: `Vcfile.Ic`や`Vcfile.Ad`は、
    ファイル内の読み込み位置を管理し、
    データの欠落や重複を防ぎます。
  - **エラーハンドリング**: 時間データの解析に失敗した場合、
    エラーメッセージを出力し、プログラムを終了します。
    これは、シミュレーション結果の信頼性を確保するために重要です。

この関数は、建物のエネルギーシミュレーション結果を分析し、
省エネルギー設計、快適性評価、
および運用改善のための意思決定を支援するための重要な役割を果たします。
*/
func tmdata(Vcfile *VCFILE, Tmdt *TMDT, Daytm *DAYTM, perio rune) int {
	var err error
	fi := Vcfile.Fi
	Estl := &Vcfile.Estl
	r := 1

	scanner := bufio.NewScanner(fi)
	for scanner.Scan() {
		s := scanner.Text()

		if s == "INAN" || s == "end" {
			if Vcfile.Ic != 2 || perio == 'y' {
				_, err := fi.Seek(Vcfile.Ad, io.SeekStart)
				if err != nil {
					panic(err)
				}
				Vcfile.Ic++
			} else {
				return 0
			}
		} else {
			D := 0
			for i := 0; i < Estl.Ntimeid; i++ {
				if i > 0 {
					if scanner.Scan() {
						s = scanner.Text()
					} else {
						return r
					}
				}

				switch Estl.Timeid[i] {
				case 'Y': //年
					Tmdt.CYear = s
					Tmdt.Year, err = strconv.Atoi(s)
					if err != nil {
						panic(err)
					}
					Tmdt.Dat[i] = &Tmdt.CYear

					if Tmdt.Year == Daytm.Year {
						D++
					}
				case 'M': //月
					Tmdt.CMon = s
					Tmdt.Mon, err = strconv.Atoi(s)
					if err != nil {
						panic(err)
					}
					Tmdt.Dat[i] = &Tmdt.CMon

					if Tmdt.Mon == Daytm.Mon {
						D++
					}
				case 'D': //日
					Tmdt.CDay = s
					Tmdt.Day, err = strconv.Atoi(s)
					if err != nil {
						panic(err)
					}
					Tmdt.Dat[i] = &Tmdt.CDay

					if Tmdt.Day == Daytm.Day {
						D++
					}
				case 'W': //曜日
					Tmdt.CWkday = s
					Tmdt.Dat[i] = &Tmdt.CWkday
				case 'T': //時刻
					Tmdt.CTime = s
					if st := strings.IndexByte(s, ':'); st != -1 {
						s = s[:st] + "." + s[st+1:]
					}
					fval, err := strconv.ParseFloat(s, 64)
					if err != nil {
						panic(err)
					}
					Tmdt.Time = int(fval*100 + 0.5)
					Tmdt.Dat[i] = &Tmdt.CTime

					if Tmdt.Time-int(100*Daytm.Time+0.5) == 0 {
						D++
					}
				}
			}

			if D == Estl.Ntimeid {
				return 1
			} else {
				for i := 0; i < Estl.Ndata; i++ {
					if scanner.Scan() {
						scanner.Text()
					} else {
						return r
					}
				}
			}
		}
	}

	return r
}

/*
esdatgt (Energy Simulation Data Get)

この関数は、シミュレーション結果ファイルから、
各項目（温度、熱量、エネルギーなど）の実際のデータ値を読み込み、
対応する`TLIST`構造体に格納します。

建築環境工学的な観点:
  - **シミュレーション結果のデータ抽出**: シミュレーション結果ファイルは、
    通常、数値データが羅列された形式で記述されています。
    この関数は、各項目のデータ型（`Tlist[j].Ptype`）に応じて、
    文字列を浮動小数点数（`float64`）、整数（`int`）、
    または文字（`rune`）に変換し、
    対応する`Fval`, `Ival`, `Cval`配列に格納します。
  - **選択項目の考慮**: `Tlist[j].Req == 'y'` の条件は、
    ユーザーが選択した項目のみを読み込むことを意味します。
    これにより、必要なデータのみを効率的に抽出し、
    メモリ使用量を削減できます。
  - **データ履歴の構築**: `Tlist[j].Fval[i]` のように、
    各項目のデータが時系列で格納されることで、
    シミュレーション期間中の値の変化を追跡できます。
  - **エラーハンドリング**: 数値の解析に失敗した場合、
    エラーメッセージを出力します。
    これは、シミュレーション結果の信頼性を確保するために重要です。

この関数は、建物のエネルギーシミュレーション結果を分析し、
省エネルギー設計、快適性評価、
および運用改善のための意思決定を支援するための重要な役割を果たします。
*/
func esdatgt(fi io.Reader, i int, Ndata int, Tlist []TLIST) {
	scanner := bufio.NewScanner(fi)
	for j := 0; j < Ndata; j++ {
		if scanner.Scan() {
			s := scanner.Text()
			if Tlist[j].Req == 'y' || Tlist[j].Vtype == 'h' || Tlist[j].Vtype == 'H' {
				switch Tlist[j].Ptype {
				case 'f':
					fval, err := strconv.ParseFloat(s, 64)
					if err != nil {
						fmt.Println(err)
					} else {
						Tlist[j].Fval[i] = fval
						// fmt.Printf("<<esdatgt>> j=%d (data=%s)  %s %s [%d]=%f\n",
						// j, s, Tlist[j].name, Tlist[j].id, i, Tlist[j].fval[i])
					}
				case 'd':
					ival, err := strconv.Atoi(s)
					if err != nil {
						fmt.Println(err)
					} else {
						Tlist[j].Ival[i] = ival
					}
				case 'c':
					Tlist[j].Cval[i] = rune(s[0])
				}
				if j > 0 {
					Tml := &Tlist[j-1]
					if Tml.Vtype == 'h' || Tml.Vtype == 'H' {
						Tlist[j].Pair = Tml
					}
				}
			}
		}
	}
}

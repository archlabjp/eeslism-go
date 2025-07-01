package eeslism

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
)

/*
monthday (Month and Day Calculation)

この関数は、与えられた月（`mo`）と日（`dayo`）から、
翌日の月と日を計算します。

建築環境工学的な観点:
- **時間管理の補助**: 建物のエネルギーシミュレーションでは、
  日単位で計算を進める際に、
  現在の日付から翌日の日付を正確に計算する必要があります。
  この関数は、各月の日数（2月は28日固定）を考慮して、
  翌日の月と日を返します。
- **シミュレーションの進行**: この関数は、
  シミュレーションのメインループにおいて、
  日付を更新するために用いられます。

この関数は、建物のエネルギーシミュレーションにおいて、
時間管理を正確に行い、
シミュレーションの進行を支援するための基礎的な役割を果たします。
*/
func monthday(mo int, dayo int) (int, int) {
	var Day = dayo + 1
	var Mon = mo

	switch mo {
	case 1, 3, 5, 7, 8, 10, 12:
		if dayo == 31 {
			Day = 1
			if mo == 12 {
				Mon = 1
			} else {
				Mon = mo + 1
			}
		}
	case 4, 6, 9, 11:
		if dayo == 30 {
			Day = 1
			Mon = mo + 1
		}
	case 2:
		if dayo == 28 {
			Day = 1
			Mon = mo + 1
		}
	}

	return Mon, Day
}

/*  気象デ－タの入力     */
var (
	Lat, Slat, Clat, Tlat, Lon, Ls, Isc float64
)

var __Weatherdt_ptt int = 25
var __Weatherdt_nc int = 0
var __Weatherdt_decl, __Weatherdt_E, __Weatherdt_tas, __Weatherdt_timedg float64
var __Weatherdt_dt [7][25]float64
var __Weatherdt_dtL [7][25]float64

/*
Weatherdt (Weather Data Processing)

この関数は、気象データファイルから現在の時刻の気象データを読み込み、
太陽位置を計算し、`WDAT`構造体（`Wd`）に格納します。
また、地表面温度の初期化や、給水温度の補間も行います。

建築環境工学的な観点:
- **気象データの読み込みと処理**: 建物のエネルギーシミュレーションでは、
  外気温度、湿度、日射量などの気象データが不可欠です。
  この関数は、`Simc.Wdtype`（気象データファイル種別）に応じて、
  HASP標準形式（`hspwdread`）またはVCFILE形式（`Wdflinput`）から気象データを読み込みます。
- **太陽位置の計算**: `FNDecl`（赤緯）、`FNE`（均時差）、`FNTtas`（真太陽時）、
  `Solpos`（太陽高度角、方位角）などの関数を呼び出し、
  現在の時刻における太陽位置を正確に計算します。
  これは、日射熱取得量や日影の計算に不可欠です。
- **地表面温度の初期化**: `EarthSrfTempInit`関数を呼び出し、
  地表面温度を初期化します。
  これは、地盤からの熱伝達をモデル化する際に用いられます。
- **給水温度の補間**: `Intgtsup`関数を呼び出し、
  給水温度を補間します。
  これは、給湯負荷計算などに用いられます。
- **時間間隔の考慮**: `Simc.DTm < 3600` の条件は、
  計算時間間隔が1時間未満の場合に、
  気象データを直線補間することを示唆します。
  これにより、より細かい時間ステップでのシミュレーションに対応できます。

この関数は、建物のエネルギーシミュレーションにおいて、
気象データを正確にモデル化し、
熱負荷計算、エネルギー消費量予測、
および省エネルギー対策の検討を行うための重要な役割を果たします。
*/
func Weatherdt(Simc *SIMCONTL, Daytm *DAYTM, Loc *LOCAT, Wd *WDAT, Exs []*EXSF, EarthSrfFlg bool) {
	var tt, Mon, Day int

	tt = Daytm.Tt

	if tt < __Weatherdt_ptt {
		if Simc.Wdtype == 'H' {
			if Simc.DTm < 3600 {
				_, Mon, Day, _ = hspwdread(Simc.Fwdata, Daytm.DayOfYear-1, Loc, &__Weatherdt_dtL)
				fmt.Printf("Mon=%d  Day=%d\n", Mon, Day)
			}

			_, Mon, Day, _ = hspwdread(Simc.Fwdata, Daytm.DayOfYear, Loc, &__Weatherdt_dt)
			if Daytm.Mon != Mon || Daytm.Day != Day {
				s := fmt.Sprintf("loop Mon/Day=%d/%d - data Mon/Day=%d/%d", Daytm.Mon, Daytm.Day, Mon, Day)
				Eprint("<Weatherdt>", s)
				Preexit()
				os.Exit(EXIT_MOND)
			}
		}

		if __Weatherdt_nc == 0 {
			Lat = Loc.Lat
			Lon = Loc.Lon
			Ls = Loc.Ls
			Sunint()
			Psyint()
			if Simc.Wdtype == 'H' {
				gtsupw(Simc.Ftsupw, Loc.Name, &(Loc.Daymxert), &(Loc.Tgrav), &(Loc.DTgr), &Loc.Twsup)

				Intgtsup(1, Loc.Twsup[:])
			}

			if EarthSrfFlg {
				Wd.EarthSurface = make([]float64, 366*25)
				EarthSrfTempInit(Simc, Loc, Wd)
			}
			__Weatherdt_nc = 1
		}

		__Weatherdt_decl = FNDecl(Daytm.DayOfYear)
		__Weatherdt_E = FNE(Daytm.DayOfYear)

		if Wd.Intgtsupw == 'N' {
			Wd.Twsup = Loc.Twsup[Daytm.Mon-1]
		} else {
			Wd.Twsup = Intgtsup(Daytm.DayOfYear, Loc.Twsup[:])
		}
	}

	Exsfter(Daytm.DayOfYear, Loc.Daymxert, Loc.Tgrav, Loc.DTgr, Exs, Wd, tt)
	__Weatherdt_timedg = float64(Daytm.Tt) + math.Mod(float64(Daytm.Ttmm), 100.0)/60.0

	__Weatherdt_tas = FNTtas(__Weatherdt_timedg, __Weatherdt_E)
	Wd.Sh, Wd.Sw, Wd.Ss, Wd.Solh, Wd.SolA = Solpos(__Weatherdt_tas, __Weatherdt_decl)

	if Simc.Wdtype == 'H' {
		// 計算時間間隔が1時間未満の場合には直線補完する
		if Simc.DTm < 3600 {
			wdatadiv(Daytm, Wd, __Weatherdt_dt, __Weatherdt_dtL)
		} else {
			dt2wdata(Wd, tt, __Weatherdt_dt)
		}
	} else {
		// VCFILE形式の気象データの読み込み
		Wdflinput(&Simc.Wdpt, Wd)
	}

	if DEBUG {
		fmt.Println("\n\n<Weatherdt>  ***** Wdata *****\n\n=================")
		fmt.Printf("\tT=%.1f\n\tx=%.4f\n\tRH=%.0f\n", Wd.T, Wd.X, Wd.RH)
		fmt.Printf("\tIdn=%.0f\n\tIsky=%.0f\n\tIhor=%.0f\n", Wd.Idn, Wd.Isky, Wd.Ihor)
		fmt.Printf("\tRN=%.0f\n\tCC=%.0f\n\tWdre=%.0f\n\tWv=%.1f\n", Wd.RN, Wd.CC, Wd.Wdre, Wd.Wv)
		fmt.Printf("\th=%.0f\n==================\n", Wd.H)
	}

	if dayprn && Ferr != nil {
		fmt.Fprintln(Ferr, "\n\n<Weatherdt>  ***** Wdata *****\n\n=================")
		fmt.Fprintf(Ferr, "\tT=%.1f\n\tx=%.4f\n\tRH=%.0f\n", Wd.T, Wd.X, Wd.RH)
		fmt.Fprintf(Ferr, "\tIdn=%.0f\n\tIsky=%.0f\n\tIhor=%.0f\n", Wd.Idn, Wd.Isky, Wd.Ihor)
		fmt.Fprintf(Ferr, "\tRN=%.0f\n\tCC=%.0f\n\tWdre=%.0f\n\tWv=%.1f\n", Wd.RN, Wd.CC, Wd.Wdre, Wd.Wv)
		fmt.Fprintf(Ferr, "\th=%.0f\n==================\n", Wd.H)
	}

	__Weatherdt_ptt = tt
}

var __gtsupw_ic int
/*
gtsupw (Ground Temperature for Supply Water)

この関数は、給水温度データファイル（`supw.efl`）から、
特定の地域（`loc`）の給水温度データ、
および地盤温度に関するパラメータ（`nmx`, `Tgrav`, `DTgr`）を読み込みます。

建築環境工学的な観点:
- **給水温度の重要性**: 給水温度は、
  給湯負荷計算や、地中熱交換器の性能評価において重要なパラメータです。
  この関数は、地域ごとの給水温度データを読み込み、
  シミュレーションに利用します。
- **地盤温度の考慮**: `Tgrav`（地盤温度）や`DTgr`（地盤温度の時定数）は、
  地盤からの熱伝達をモデル化する際に用いられます。
  これは、地下室や基礎からの熱損失・熱取得を評価する上で重要です。
- **データ読み込みの効率化**: `__gtsupw_ic == 0` の条件は、
  ファイルからのデータ読み込みを一度だけ行うことを意味します。
  これにより、計算効率を向上させます。
- **エラーハンドリング**: 指定された地域名が見つからない場合、
  エラーメッセージを出力し、プログラムを終了します。
  これは、入力データの不備を早期に発見し、
  シミュレーションの信頼性を確保するために重要です。

この関数は、建物のエネルギーシミュレーションにおいて、
給水温度や地盤温度を正確にモデル化し、
熱負荷計算、エネルギー消費量予測、
および省エネルギー対策の検討を行うための重要な役割を果たします。
*/
func gtsupw(fp []byte, loc string, nmx *int, Tgrav, DTgr *float64, Tsupw *[12]float64) {
	var flg int
	var s string

	if __gtsupw_ic == 0 {
		reader := strings.NewReader(string(fp))
		scanner := bufio.NewScanner(reader)
		scanner.Split(bufio.ScanWords)

		// Find location
		for scanner.Scan() {
			s = scanner.Text()
			if s == loc {
				flg = 1
				break
			}
		}

		// Check error
		if scanner.Err() != nil {
			E := fmt.Sprintf("supw.eflに%sが登録されていません。\n", s)
			Eprint("<gtsupw>", E)
			os.Exit(EXIT_GTSUPW)
		}

		// Read data of 12 months and nmx, Tgrav, DTgr
		var err error
		var values [15]float64
		for i := 0; i < 15; i++ {
			scanner.Scan()
			values[i], err = strconv.ParseFloat(scanner.Text(), 64)
			if err != nil {
				panic(err)
			}
		}
		copy(Tsupw[:], values[:12])
		*nmx = int(values[12])
		*Tgrav = values[13]
		*DTgr = values[14]

		__gtsupw_ic = 1

		if flg == 0 {
			E := fmt.Sprintf("supw.eflに%sが登録されていません。\n", s)
			Eprint("<gtsupw>", E)
			os.Exit(1)
		}
	}
}

/*   HASP標準気象デ－タファイルからの入力   */
var __hspwdread_ic int
var __hspwdread_recl int

/*
hspwdread (HASP Weather Data Read)

この関数は、HASP標準気象データファイルから、
指定された通日（`nday`）の気象データ（温度、湿度、日射量など）を読み込みます。

建築環境工学的な観点:
- **HASP標準気象データ**: HASP（Heating, Air-conditioning and Sanitary Engineering Program）は、
  日本における建築設備設計で広く用いられているシミュレーションプログラムであり、
  その気象データ形式は標準的です。
  この関数は、その標準形式の気象データを読み込むことで、
  シミュレーションの入力データとして利用します。
- **気象データの読み込みと処理**: 気象データファイルから、
  年、月、日、曜日、および各時刻の気象要素（温度、絶対湿度、直達日射量、拡散日射量、雲量、風向、風速）を読み込みます。
  読み込んだデータは、`dt`配列に格納され、
  その後の熱負荷計算や機器の運転制御に利用されます。
- **ファイルポインターの管理**: `fp.Seek`を用いてファイルポインターを移動させることで、
  指定された通日のデータを効率的に読み込みます。
- **エラーハンドリング**: ファイルの読み込みに失敗した場合や、
  日付の整合性が取れない場合、
  エラーメッセージを出力し、プログラムを終了します。
  これは、シミュレーションの信頼性を確保するために重要です。

この関数は、建物のエネルギーシミュレーションにおいて、
気象データを正確にモデル化し、
熱負荷計算、エネルギー消費量予測、
および省エネルギー対策の検討を行うための重要な役割を果たします。
*/
func hspwdread(fp io.ReadSeeker, nday int, Loc *LOCAT, dt *[7][25]float64) (year int, mon int, day int, wkdy int) {
	var d, a, b, c float64
	var k, t int
	Data := [24][3]byte{}
	Yr, Mon, Day, Wk, Sq := [2]byte{}, [2]byte{}, [2]byte{}, [1]byte{}, [2]byte{}
	var s string

	if nday > 365 {
		nday = nday - 365
	}
	if nday <= 0 {
		nday = nday + 365
	}

	if __hspwdread_ic == 0 {
		// NOTE: 独自のHASPフォーマットになっている
		fmt.Fscanf(fp, "%s %f %f %f %f %f %f ", &s, &Loc.Lat, &Loc.Lon, &Loc.Ls, &a, &b, &c)
		Loc.Name = s

		if Ferr != nil {
			fmt.Fprintf(Ferr, "\n------> <hspwdread> \n")
			fmt.Fprintf(Ferr, "\nName=%s\tLat=%.4g\tLon=%.4g\tLs=%.4g\ta=%.4g\tb=%.4g\tc=%.4g\n",
				Loc.Name, Loc.Lat, Loc.Lon, Loc.Ls, a, b, c)
		}

		//改行コードの判定
		fp.Seek(0, io.SeekEnd)
		fsize, _ := fp.Seek(0, io.SeekCurrent)
		__hspwdread_recl = int(fsize / 2556)
		//recl=82 -> Windows (CRLF)
		//recl=81 -> Linux/Mac(CR or LF)
	}

	if __hspwdread_ic != nday {
		fp.Seek(int64(__hspwdread_recl*7*(nday-1)+__hspwdread_recl), 0)
	}

	if __hspwdread_ic > 0 {
		for k = 0; k < 7; k++ {
			dt[k][0] = dt[k][24]
		}
	}

	for k = 0; k < 7; k++ {
		for t = 0; t < 24; t++ {
			fp.Read(Data[t][:])
		}

		fp.Read(Yr[:])
		fp.Read(Mon[:])
		fp.Read(Day[:])
		fp.Read(Wk[:])
		fp.Read(Sq[:])

		for t = 1; t < 25; t++ {
			var err error
			d, err = strconv.ParseFloat(strings.TrimSpace(string(Data[t-1][:])), 64)
			if err != nil {
				panic(err)
			}
			switch k {
			case 0:
				dt[k][t] = (d - 500.0) * 0.1
			case 1:
				dt[k][t] = 0.0001 * d
			case 6:
				dt[k][t] = 0.1 * d
			default:
				dt[k][t] = d
			}
		}
	}

	fmt.Printf("C Mon=%s Day=%s\n", Mon, Day)

	year, _ = strconv.Atoi(strings.TrimSpace(string(Yr[:])))
	mon, _ = strconv.Atoi(strings.TrimSpace(string(Mon[:])))
	day, _ = strconv.Atoi(strings.TrimSpace(string(Day[:])))
	wkdy, _ = strconv.Atoi(strings.TrimSpace(string(Wk[:])))

	if __hspwdread_ic == 0 {
		for k = 0; k < 7; k++ {
			dt[k][0] = dt[k][1]
		}
	}
	__hspwdread_ic = nday + 1

	if DEBUG {
		for t = 0; t < 25; t++ {
			fmt.Printf("%2d %5.1f %6.4f %5.0f %5.0f %2.0f %2.0f %5.1f\n", t, dt[0][t], dt[1][t], dt[2][t], dt[3][t], dt[4][t], dt[5][t], dt[6][t])
		}
	}

	return year, mon, day, wkdy
}

/*
dt2wdata (Data to Weather Data Conversion)

この関数は、気象データ配列（`dt`）から、
指定された時刻（`tt`）の気象要素（温度、湿度、日射量など）を抽出し、
`WDAT`構造体（`Wd`）に格納します。
また、日射量や夜間放射量の計算、
および相対湿度やエンタルピーの補完も行います。

建築環境工学的な観点:
- **気象データの抽出と変換**: 気象データ配列から、
  現在の時刻の気象データを抽出し、
  熱負荷計算や機器の運転制御に利用しやすい形式に変換します。
- **日射量と夜間放射量の計算**: 
  - `Wd.Idn`, `Wd.Isky`, `Wd.Ihor`: 直達日射量、天空日射量、全天日射量を計算します。
  - `Wd.RN`, `Wd.Rsky`: 夜間放射量、天空放射量を計算します。
  これらの計算は、日射熱取得や放射熱損失を正確にモデル化するために重要です。
- **データ補完と整合性**: 
  - `Wd.RH = FNRhtx(Wd.T, Wd.X)`: 温度と絶対湿度から相対湿度を計算し、
    データの整合性を確保します。
  - `Wd.H = FNH(Wd.T, Wd.X)`: 温度と絶対湿度からエンタルピーを計算します。

この関数は、建物のエネルギーシミュレーションにおいて、
気象データを正確にモデル化し、
熱負荷計算、エネルギー消費量予測、
および省エネルギー対策の検討を行うための重要な役割を果たします。
*/
func dt2wdata(Wd *WDAT, tt int, dt [7][25]float64) {
	Wd.T = dt[0][tt]
	Wd.X = dt[1][tt]
	Wd.Idn = dt[2][tt] / 0.86
	Wd.Isky = dt[3][tt] / 0.86
	Wd.Ihor = Wd.Idn*Wd.Sh + Wd.Isky

	if Wd.RNtype == 'C' {
		Br := 0.51 + 0.209*math.Sqrt(FNPwx(Wd.X))
		Wd.CC = dt[4][tt]
		Wd.RN = (1.0 - 0.62*Wd.CC/10.0) * (1.0 - Br) * Sgm * math.Pow(Wd.T+273.15, 4.0)
		Wd.Rsky = ((1.0-0.62*Wd.CC/10.0)*Br + 0.62*Wd.CC/10.0) * Sgm * math.Pow(Wd.T+273.15, 4.0)
	} else {
		Wd.CC = -999.0
		Wd.RN = dt[4][tt] / 0.86
		Wd.Rsky = Sgm*math.Pow(Wd.T+273.15, 4.0) - Wd.RN
	}

	Wd.Wdre = dt[5][tt]
	Wd.Wv = dt[6][tt]
	Wd.RH = FNRhtx(Wd.T, Wd.X)
	Wd.H = FNH(Wd.T, Wd.X)
}

/*
wddatadiv (Weather Data Division and Interpolation)

この関数は、気象データ配列（`dt`）から、
指定された時刻（`Daytm.Ttmm`）の気象要素（温度、湿度、日射量など）を抽出し、
必要に応じて直線補間を行います。

建築環境工学的な観点:
- **気象データの補間**: シミュレーションの計算時間間隔が、
  気象データの時間間隔よりも細かい場合、
  気象データを補間して連続的な値を得る必要があります。
  この関数は、`Lineardiv`関数を呼び出して直線補間を行います。
- **時間管理の考慮**: `Daytm.Ttmm%100 == 0` の条件は、
  時刻がちょうど時間の区切りである場合を示し、
  それ以外の場合は補間を行います。
- **日をまたぐ補間**: `Daytm.Ttmm > 100` の条件は、
  日をまたいで補間を行う場合を示し、
  前日のデータ（`dtL`）と当日のデータ（`dt`）を考慮します。

この関数は、建物のエネルギーシミュレーションにおいて、
気象データを正確にモデル化し、
熱負荷計算、エネルギー消費量予測、
および省エネルギー対策の検討を行うための重要な役割を果たします。
*/
func wdatadiv(Daytm *DAYTM, Wd *WDAT, dt [7][25]float64, dtL [7][25]float64) {
	var WdF, WdL WDAT
	var r float64

	WdL.RNtype = Wd.RNtype
	WdF.RNtype = Wd.RNtype

	r = math.Mod(float64(Daytm.Ttmm), 100.0) / 60.0

	if Daytm.Ttmm%100 == 0 {
		dt2wdata(Wd, Daytm.Tt, dt)
	} else if Daytm.Ttmm > 100 {
		dt2wdata(&WdL, Daytm.Tt, dt)
		dt2wdata(&WdF, Daytm.Tt+1, dt)

		WdLineardiv(Wd, &WdL, &WdF, r)
	} else {
		dt2wdata(&WdL, 24, dtL)
		dt2wdata(&WdF, 1, dt)

		WdLineardiv(Wd, &WdL, &WdF, r)
	}
}

/*
WdLineardiv (Weather Data Linear Interpolation)

この関数は、気象データ（温度、湿度、日射量など）を直線補間します。

建築環境工学的な観点:
- **気象データの補間**: シミュレーションの計算時間間隔が、
  気象データの時間間隔よりも細かい場合、
  気象データを補間して連続的な値を得る必要があります。
  この関数は、`Lineardiv`関数を呼び出して直線補間を行います。
- **データ整合性の確保**: 補間後も、
  相対湿度やエンタルピーなどの湿り空気の状態値が、
  物理的に整合性が取れるように再計算します。
- **風向の調整**: `Wd.Wdre > 16.0` の条件は、
  風向が特定の範囲を超えた場合に調整を行うことを示唆します。

この関数は、建物のエネルギーシミュレーションにおいて、
気象データを正確にモデル化し、
熱負荷計算、エネルギー消費量予測、
および省エネルギー対策の検討を行うための重要な役割を果たします。
*/
func WdLineardiv(Wd *WDAT, WdL *WDAT, WdF *WDAT, dt float64) {
	Wd.T = Lineardiv(WdL.T, WdF.T, dt)
	Wd.X = Lineardiv(WdL.X, WdF.X, dt)
	// Wd.RH = Lineardiv(WdL.RH, WdF.RH, dt)
	Wd.RH = FNRhtx(Wd.T, Wd.X)
	// Wd.h = Lineardiv(WdL.h, WdF.h, dt)
	Wd.H = FNH(Wd.T, Wd.X)
	Wd.Idn = Lineardiv(WdL.Idn, WdF.Idn, dt)
	Wd.Isky = Lineardiv(WdL.Isky, WdF.Isky, dt)
	// Wd.Ihor = Lineardiv(WdL.Ihor, WdF.Ihor, dt)
	Wd.Ihor = Wd.Idn*Wd.Sh + Wd.Isky
	Wd.CC = Lineardiv(WdL.CC, WdF.CC, dt)
	Wd.RN = Lineardiv(WdL.RN, WdF.RN, dt)
	Wd.Wv = Lineardiv(WdL.Wv, WdF.Wv, dt)
	Wd.Wdre = Lineardiv(WdL.Wdre, WdF.Wdre, dt)

	if Wd.Wdre > 16.0 {
		Wd.Wdre -= 16.0
	}
}

/*
EarthSrfTempInit (Earth Surface Temperature Initialization)

この関数は、地表面温度を計算し、初期化します。
これは、地盤からの熱伝達をモデル化する際に用いられる重要な情報です。

建築環境工学的な観点:
- **地中熱交換のモデル化**: 地盤に接する壁面や床面は、
  地盤からの熱伝達を受けます。
  この関数は、地盤の熱的特性（熱容量`Soic`、熱伝導率`Soil`、熱拡散率`a`）と、
  外部気象条件（外気温度`Wd.T`、日射量`Ihol`、夜間放射量`Wd.RN`）を考慮して、
  地表面温度を計算します。
- **熱伝導方程式の解法**: 地中温度の計算には、
  熱伝導方程式を数値的に解く必要があります。
  この関数は、地中を複数の層に分割し、
  各層の熱収支を連立方程式として解くことで、
  地中温度分布を計算します。
  `U`は係数行列、`oldT`は前時刻の温度、`T`は現在の温度です。
- **助走期間の考慮**: `for year = 0; year < 2; year++` のループは、
  地中温度が定常状態に達するまでの助走期間をシミュレーションすることを意味します。
  これにより、シミュレーション開始時の初期条件が結果に与える影響を最小限に抑えます。

この関数は、建物のエネルギーシミュレーションにおいて、
地盤からの熱伝達を正確にモデル化し、
熱負荷計算、エネルギー消費量予測、
および省エネルギー対策の検討を行うための重要な役割を果たします。
*/
func EarthSrfTempInit(Simc *SIMCONTL, Loc *LOCAT, Wd *WDAT) {
	var decl, E, ac, Te, tas float64
	var year, nday, tt int
	//var wkdy, Year, Mon, Day int
	dt := [7][25]float64{}
	var U, T, oldT []float64
	var i, j int
	var Ihol float64
	var Soic, Soil float64

	// 地表面温度の計算
	fmt.Println("地表面温度の計算開始")

	// 土壌の容積比熱[J/m3K]
	Soic = 3.34e6
	// 土壌の熱伝導率[W/mK]
	Soil = 1.047
	// 土壌の熱拡散率
	a := Soil / Soic
	ac = 23.0
	u := float64(Simc.DTm) * a / (0.5 * 0.5)
	U = make([]float64, 20*20)
	T = make([]float64, 20*20)
	oldT = make([]float64, 20*20)

	// 地中温度の初期温度は平均外気温度
	matinitx(oldT, 20, Loc.Tgrav)

	// 地中温度計算用の行列の作成
	for i = 0; i < 20; i++ {
		U[i*20+i] = 1.0 + 2.0*u
		if i > 0 {
			U[i*20+i-1] = -u
		}
		if i < 20-1 {
			U[i*20+i+1] = -u
		}

		if i == 0 {
			U[i*20+i] = 1.0 + float64(Simc.DTm)/(0.5*(0.0+Soic*0.5)*1.0/ac) + float64(Simc.DTm)/(0.5*(0.0+Soic*0.5)*0.5/Soil)
			U[i*20+i+1] = -float64(Simc.DTm) / (0.5 * (0.0 + Soic*0.5) * 0.5 / Soil)
		}
		if i == 19 {
			U[i*20+i] = 1.0 + float64(Simc.DTm)/(0.5*(0.0+Soic*0.5)*0.5/Soil)
			U[i*20+i-1] = -float64(Simc.DTm) / (0.5 * (0.0 + Soic*0.5) * 0.5 / Soil)
		}
	}

	//matfprint("%.2lf ", 20,U )
	// Uの逆行列の計算
	Matinv(U, 20, 20, "<EarthSrfTempInit>")

	// １年目は助走期間
	for year = 0; year < 2; year++ {
		for nday = 1; nday <= 365; nday++ {
			decl = FNDecl(nday)
			E = FNE(nday)

			hspwdread(Simc.Fwdata2, nday, Loc, &dt)

			for tt = 1; tt <= 24; tt++ {
				matinit(T, 20)

				tas = FNTtas(float64(tt), E)
				var Sh float64
				Sh, _, _, Wd.Solh, Wd.Solh = Solpos(tas, decl)
				dt2wdata(Wd, tt, dt)

				Ihol = Wd.Idn*Sh + Wd.Isky

				Te = Wd.T + (0.7*Ihol-Wd.RN)/ac

				oldT[0] += (Te * float64(Simc.DTm) / (0.5 * Soic * 0.5 / ac))

				Matmalv(U, oldT, 20, 20, T)

				if year == 1 {
					Wd.EarthSurface[nday*24+tt] = T[0]
				}

				for j = 0; j < 20; j++ {
					oldT[j] = T[j]
				}
			}
		}
	}

	fmt.Println("地表面温度の計算終了")
}

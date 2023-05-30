package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
)

/* 月日の設定（翌日の日時を設定する） */

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

func Weatherdt(Simc *SIMCONTL, Daytm *DAYTM, Loc *LOCAT, Wd *WDAT, Exs []EXSF, EarthSrfFlg rune) {
	var Sh, Sw, Ss float64
	var tt, Mon, Day int

	tt = Daytm.Tt

	if tt < __Weatherdt_ptt {
		if Simc.Wdtype == 'H' {
			if Simc.DTm < 3600 {
				_, Mon, Day, _ = hspwdread(Simc.Fwdata, Daytm.DayOfYear-1, Loc, &__Weatherdt_dtL)
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

			if EarthSrfFlg == 'Y' {
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
	Solpos(__Weatherdt_tas, __Weatherdt_decl, &Wd.Sh, &Wd.Sw, &Wd.Ss, &Wd.Solh, &Wd.SolA)

	Wd.Sh = Sh
	Wd.Sw = Sw
	Wd.Ss = Ss

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

	if Dayprn && Ferr != nil {
		fmt.Fprintln(Ferr, "\n\n<Weatherdt>  ***** Wdata *****\n\n=================")
		fmt.Fprintf(Ferr, "\tT=%.1f\n\tx=%.4f\n\tRH=%.0f\n", Wd.T, Wd.X, Wd.RH)
		fmt.Fprintf(Ferr, "\tIdn=%.0f\n\tIsky=%.0f\n\tIhor=%.0f\n", Wd.Idn, Wd.Isky, Wd.Ihor)
		fmt.Fprintf(Ferr, "\tRN=%.0f\n\tCC=%.0f\n\tWdre=%.0f\n\tWv=%.1f\n", Wd.RN, Wd.CC, Wd.Wdre, Wd.Wv)
		fmt.Fprintf(Ferr, "\th=%.0f\n==================\n", Wd.H)
	}

	__Weatherdt_ptt = tt
}

var __gtsupw_ic int

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
		fmt.Fscanf(fp, "%s %f %f %f %f %f %f ", &s, &Loc.Lat, &Loc.Lon, &Loc.Ls, &a, &b, &c)
		Loc.Name = s

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

func EarthSrfTempInit(Simc *SIMCONTL, Loc *LOCAT, Wd *WDAT) {
	var decl, E, ac, Te, tas, Sh, Sw, Ss float64
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
				Solpos(tas, decl, &Sh, &Sw, &Ss, &Wd.Solh, &Wd.Solh)
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

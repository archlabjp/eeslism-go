package eeslism

import (
	"fmt"
	"io"
)

var __Wdtsum_oldday, __Wdtsum_hrs, __Wdtsum_oldMon, __Wdtsum_hrsm int
var __Wdtsum_cffWh float64

func Wdtsum(Mon int, Day int, Nday int, ttmm int, Wd *WDAT, Nexs int, Exs []EXSF,
	Wdd *WDAT, Wdm *WDAT, Soldy []float64, Solmon []float64, Simc *SIMCONTL) {
	var e EXSF

	// 日集計の初期化
	if Nday != __Wdtsum_oldday {
		__Wdtsum_cffWh = Cff_kWh * 1000.0
		__Wdtsum_hrs = 0
		Wdd.T = 0.0
		Wdd.X = 0.0
		Wdd.Wv = 0.0
		Wdd.Idn = 0.0
		Wdd.Isky = 0.0
		Wdd.RN = 0.0

		for i := 0; i < Nexs; i++ {
			Soldy[i] = 0.0
		}

		__Wdtsum_oldday = Nday
	}

	// 月集計の初期化
	if Mon != __Wdtsum_oldMon {
		__Wdtsum_cffWh = Cff_kWh * 1000.0
		__Wdtsum_hrsm = 0
		Wdm.T = 0.0
		Wdm.X = 0.0
		Wdm.Wv = 0.0
		Wdm.Idn = 0.0
		Wdm.Isky = 0.0
		Wdm.RN = 0.0

		for i := 0; i < Nexs; i++ {
			Solmon[i] = 0.0
		}

		__Wdtsum_oldMon = Mon
	}

	// 日集計
	__Wdtsum_hrs++
	Wdd.T += Wd.T
	Wdd.X += Wd.X
	Wdd.Wv += Wd.Wv
	Wdd.Idn += Wd.Idn
	Wdd.Isky += Wd.Isky
	Wdd.RN += Wd.RN

	for i := 0; i < Nexs; i++ {
		e = Exs[i]

		if e.Typ != 'E' && e.Typ != 'e' {
			Soldy[i] += e.Iw
		} else {
			Soldy[i] += e.Tearth
		}
	}

	// 月集計
	__Wdtsum_hrsm++
	Wdm.T += Wd.T
	Wdm.X += Wd.X
	Wdm.Wv += Wd.Wv
	Wdm.Idn += Wd.Idn
	Wdm.Isky += Wd.Isky
	Wdm.RN += Wd.RN

	for i := 0; i < Nexs; i++ {
		e = Exs[i]

		if e.Typ != 'E' && e.Typ != 'e' {
			Solmon[i] += e.Iw
		} else {
			Solmon[i] += e.Tearth
		}
	}

	if ttmm == 2400 {
		Wdd.T /= float64(__Wdtsum_hrs)
		Wdd.X /= float64(__Wdtsum_hrs)
		Wdd.Wv /= float64(__Wdtsum_hrs)
		Wdd.Idn *= __Wdtsum_cffWh
		Wdd.Isky *= __Wdtsum_cffWh
		Wdd.RN *= __Wdtsum_cffWh

		for i := 0; i < Nexs; i++ {
			e = Exs[i]
			if e.Typ != 'E' && e.Typ != 'e' {
				Soldy[i] *= __Wdtsum_cffWh
			} else {
				Soldy[i] /= float64(__Wdtsum_hrs)
			}
		}
	}
	if IsEndDay(Mon, Day, Nday, Simc.Dayend) && __Wdtsum_hrsm > 0 && ttmm == 2400 {
		Wdm.T /= float64(__Wdtsum_hrsm)
		Wdm.X /= float64(__Wdtsum_hrsm)
		Wdm.Wv /= float64(__Wdtsum_hrsm)
		Wdm.Idn *= __Wdtsum_cffWh
		Wdm.Isky *= __Wdtsum_cffWh
		Wdm.RN *= __Wdtsum_cffWh

		for i := 0; i < Nexs; i++ {
			e = Exs[i]
			if e.Typ != 'E' && e.Typ != 'e' {
				Soldy[i] *= __Wdtsum_cffWh
			} else {
				Soldy[i] /= float64(__Wdtsum_hrsm)
			}
		}
	}
}

var __Wdtdprint_ic int

/* 気象データ日集計値出力 */
func Wdtdprint(fo io.Writer, title string, Mon int, Day int, Wdd *WDAT, Nexs int, Exs []EXSF, Soldy []float64) {
	if __Wdtdprint_ic == 0 {
		__Wdtdprint_ic++
		fmt.Fprintf(fo, "%s;\n %d\n", title, Nexs)

		fmt.Fprintf(fo, "Mo\tNd\tWd_T\tWd_x\tWd_Wv\tWd_RN\tWd_Idn\tWd_Isky\t")
		for i := 0; i < Nexs; i++ {
			e := Exs[i]
			fmt.Fprintf(fo, "%s[%c]\t", e.Name, e.Typ)
		}
		fmt.Fprintf(fo, "\n")
	}

	fmt.Fprintf(fo, "%d\t%d\t", Mon, Day)
	fmt.Fprintf(fo, "%.1f\t%.4f\t%.1f\t%.2f\t%.2f\t%4.2f", Wdd.T, Wdd.X, Wdd.Wv, Wdd.RN/1000., Wdd.Idn/1000., Wdd.Isky/1000.)

	for i := 0; i < Nexs; i++ {
		e := Exs[i]
		if e.Typ != 'E' && e.Typ != 'e' {
			fmt.Fprintf(fo, "\t%.2f", Soldy[i]/1000.)
		} else {
			fmt.Fprintf(fo, "\t%.1f", Soldy[i])
		}
	}
	fmt.Fprintf(fo, "\n")
}

var __Wdtprint_ic int

// 気象データの出力
func Wdtprint(fo io.Writer, title string, Mon, Day int, time float64, Wd *WDAT, Exsfst *EXSFS) {
	var Nexs, i int
	Nexs = Exsfst.Nexs

	if DEBUG {
		fmt.Printf("N=%d\t%d/%d\t%.2f\n", Nexs, Mon, Day, time)
		fmt.Printf("%s;\n %d\n", title, Nexs)
	}

	// ヘッダー部の出力
	if __Wdtprint_ic == 0 {
		__Wdtprint_ic++
		fmt.Fprintf(fo, "%s;\n %d\n", title, Nexs)
		fmt.Fprintf(fo, "Mon\tDay\tTime\tWd_T\tWd_x\tWd_RH\tWd_Wv\t")
		fmt.Fprintf(fo, "Wd_Wdre\tWd_RN\tWd_Idn\tWd_Isky\tsolh\tsolA\t")

		for i = 0; i < Nexs; i++ {
			e := Exsfst.Exs[i]
			if DEBUG {
				fmt.Printf("%s[%c]\t", e.Name, e.Typ)
			}

			fmt.Fprintf(fo, "%s[%c]\t", e.Name, e.Typ)
		}
		fmt.Fprintf(fo, "\n")
	}

	// 月・日・時刻の出力
	fmt.Fprintf(fo, "%d\t%d\t%.2f\t", Mon, Day, time)

	// 気象データの出力
	fmt.Fprintf(fo, "%.2f\t%.4f\t%.0f\t%.1f\t%.0f\t%.0f\t%.0f\t%.0f\t%.1f\t%.1f\t",
		Wd.T, Wd.X, Wd.RH, Wd.Wv, Wd.Wdre, Wd.RN, Wd.Idn, Wd.Isky, Wd.Solh, Wd.SolA)

	// 外表面の全日射・地中温度の出力
	for i = 0; i < Nexs; i++ {
		e := Exsfst.Exs[i]
		if e.Typ != 'E' && e.Typ != 'e' {
			fmt.Fprintf(fo, "%.0f\t", e.Iw)
		} else {
			fmt.Fprintf(fo, "%.1f\t", e.Tearth)
		}
	}
	fmt.Fprintf(fo, "\n")
}

var __Wdtmprint_ic int

// 気象データの出力
func Wdtmprint(fo io.Writer, title string, Mon, Day int, Wdm *WDAT, Nexs int, Exs []EXSF, Solmon []float64) {
	if __Wdtmprint_ic == 0 {
		__Wdtmprint_ic++
		fmt.Fprintf(fo, "%s;\n%d\n", title, Nexs)

		fmt.Fprintf(fo, "Mo\tNd\tWd_T\tWd_x\tWd_Wv\tWd_RN\tWd_Idn\tWd_Isky\t")
		for _, e := range Exs {
			fmt.Fprintf(fo, "%s[%c]\t", e.Name, e.Typ)
		}
		fmt.Fprintln(fo)
	}

	fmt.Fprintf(fo, "%d\t%d\t", Mon, Day)
	fmt.Fprintf(fo, "%.1f\t%.4f\t%.1f\t%.2f\t%.2f\t%4.2f",
		Wdm.T, Wdm.X, Wdm.Wv, Wdm.RN/1000., Wdm.Idn/1000., Wdm.Isky/1000.)

	for i, e := range Exs {
		if e.Typ != 'E' && e.Typ != 'e' {
			fmt.Fprintf(fo, "\t%.2f", Solmon[i]/1000.)
		} else {
			fmt.Fprintf(fo, "\t%.1f", Solmon[i])
		}
	}
	fmt.Fprintln(fo)
}

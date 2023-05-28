package main

import (
	"fmt"
	"os"
)

func Wdtsum(Mon int, Day int, Nday int, ttmm int, Wd *WDAT, Nexs int, Exs []EXSF,
	Wdd *WDAT, Wdm *WDAT, Soldy []float64, Solmon []float64, Simc *SIMCONTL) {
	var oldday, hrs, oldMon, hrsm int
	var cffWh float64
	var e EXSF

	// 日集計の初期化
	if Nday != oldday {
		cffWh = Cff_kWh * 1000.0
		hrs = 0
		Wdd.T = 0.0
		Wdd.X = 0.0
		Wdd.Wv = 0.0
		Wdd.Idn = 0.0
		Wdd.Isky = 0.0
		Wdd.RN = 0.0

		for i := 0; i < Nexs; i++ {
			Soldy[i] = 0.0
		}

		oldday = Nday
	}

	// 月集計の初期化
	if Mon != oldMon {
		cffWh = Cff_kWh * 1000.0
		hrsm = 0
		Wdm.T = 0.0
		Wdm.X = 0.0
		Wdm.Wv = 0.0
		Wdm.Idn = 0.0
		Wdm.Isky = 0.0
		Wdm.RN = 0.0

		for i := 0; i < Nexs; i++ {
			Solmon[i] = 0.0
		}

		oldMon = Mon
	}

	// 日集計
	hrs++
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
	hrsm++
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
		Wdd.T /= float64(hrs)
		Wdd.X /= float64(hrs)
		Wdd.Wv /= float64(hrs)
		Wdd.Idn *= cffWh
		Wdd.Isky *= cffWh
		Wdd.RN *= cffWh

		for i := 0; i < Nexs; i++ {
			e = Exs[i]
			if e.Typ != 'E' && e.Typ != 'e' {
				Soldy[i] *= cffWh
			} else {
				Soldy[i] /= float64(hrs)
			}
		}
	}
	if IsEndDay(Mon, Day, Nday, Simc.Dayend) && hrsm > 0 && ttmm == 2400 {
		Wdm.T /= float64(hrsm)
		Wdm.X /= float64(hrsm)
		Wdm.Wv /= float64(hrsm)
		Wdm.Idn *= cffWh
		Wdm.Isky *= cffWh
		Wdm.RN *= cffWh

		for i := 0; i < Nexs; i++ {
			e = Exs[i]
			if e.Typ != 'E' && e.Typ != 'e' {
				Soldy[i] *= cffWh
			} else {
				Soldy[i] /= float64(hrsm)
			}
		}
	}
}

var __Wdtdprint_ic int

/* 気象データ日集計値出力 */
func Wdtdprint(fo *os.File, title string, Mon int, Day int, Wdd *WDAT, Nexs int, Exs []EXSF, Soldy []float64) {
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

/* 気象データ出力 */

var __Wdtprint_ic int

func Wdtprint(fo *os.File, title string, Mon, Day int, time float64, Wd *WDAT, Exsfst *EXSFS) {
	var Nexs, i int
	Nexs = Exsfst.Nexs

	if DEBUG {
		fmt.Printf("N=%d\t%d/%d\t%.2f\n", Nexs, Mon, Day, time)
		fmt.Printf("%s;\n %d\n", title, Nexs)
	}

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

	fmt.Fprintf(fo, "%d\t%d\t%.2f\t", Mon, Day, time)
	fmt.Fprintf(fo, "%.2f\t%.4f\t%.0f\t%.1f\t%.0f\t%.0f\t%.0f\t%.0f\t%.1f\t%.1f\t",
		Wd.T, Wd.X, Wd.RH, Wd.Wv, Wd.Wdre, Wd.RN, Wd.Idn, Wd.Isky, Wd.Solh, Wd.SolA)

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

var Wdtmprint_ic int

func Wdtmprint(fo *os.File, title string, Mon, Day int, Wdm *WDAT, Nexs int, Exs []EXSF, Solmon []float64) {
	if Wdtmprint_ic == 0 {
		Wdtmprint_ic++
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

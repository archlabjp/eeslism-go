package main

import (
	"fmt"
)

/* ---------------------------------------------------------------- */
/* 毎時計算値出力　*/

func Eeprinth(Daytm *DAYTM, Simc *SIMCONTL, Nflout int, flout []*FLOUT, Rmvls *RMVLS, Exsfst *EXSFS, Nmpath int, Mpath []MPATH, Eqsys *EQSYS, Wd *WDAT) {
	if Daytm.Ddpri != 0 && Simc.Dayprn[Daytm.Day] != 0 {
		title := Simc.Title
		Mon := Daytm.Mon
		Day := Daytm.Day
		time := Daytm.Time

		for i := 0; i < Nflout; i++ {
			flo := flout[i]

			if DEBUG {
				fmt.Printf("Eeprinth MAX=%d flo[%d]=%s\n", Nflout, i, flo.Idn)
			}

			// 気象データの出力
			switch flo.Idn {
			case PRTHWD:
				if DEBUG {
					fmt.Println("<Eeprinth> xprsolrd")
				}
				Wdtprint(flo.F, title, Mon, Day, time, Wd, Exsfst)
			case PRTCOMP: // 毎時機器の出力
				Hcmpprint(flo.F, string(PRTCOMP), Simc, Mon, Day, time, Eqsys, Rmvls.Nrdpnl, Rmvls.Rdpnl)
			case PRTPATH: // システム経路の温湿度出力
				Pathprint(flo.F, title, Mon, Day, time, Nmpath, Mpath)
			case PRTHRSTANK: // 蓄熱槽内温度分布の出力
				Hstkprint(flo.F, title, Mon, Day, time, Eqsys)
			default:
				if SIMUL_BUILDG { // these blocks are only compiled in debug builds
					switch flo.Idn {
					case PRTREV: // 毎時室温、MRTの出力
						Rmevprint(flo.F, title, Rmvls.Room, Mon, Day, time)
					case PRTHROOM:
						Rmpnlprint(flo.F, string(PRTHROOM), Simc, Mon, Day, time, Rmvls.Nroom, Rmvls.Room) // 放射パネルの出力
					case PRTHELM:
						Helmprint(flo.F, string(PRTHELM), Simc, Mon, Day, time, Rmvls.Nroom, Rmvls.Room, &Rmvls.Qetotal)
					case PRTHELMSF:
						Helmsurfprint(flo.F, string(PRTHELMSF), Simc, Mon, Day, time, Rmvls.Nroom, Rmvls.Room)
					case PRTPMV: // PMV計算
						Pmvprint(flo.F, title, Rmvls.Room, Mon, Day, time)
					case PRTQRM: // Output of heat gain elements in room
						Qrmprint(flo.F, title, Mon, Day, time, Rmvls.Room, Rmvls.Qrm)
					case PRTRSF: // Output of indoor surface temperature
						Rmsfprint(flo.F, title, Mon, Day, time, Rmvls.Room, Rmvls.Sd)
					case PRTSHD: // Output of sunshade area
						Shdprint(flo.F, title, Mon, Day, time, Rmvls.Nsrf, Rmvls.Sd)
					case PRTWAL: // Output of wall internal temperature
						Wallprint(flo.F, title, Mon, Day, time, Rmvls.Nsrf, Rmvls.Sd)
					case PRTPCM: // Output of PCM state value
						PCMprint(flo.F, title, Mon, Day, time, Rmvls.Nsrf, Rmvls.Sd)
					case PRTSFQ: // Output of indoor surface heat flow
						Rmsfqprint(flo.F, title, Mon, Day, time, Rmvls.Room, Rmvls.Sd)
					case PRTSFA: // Output of indoor surface heat transfer coefficient
						Rmsfaprint(flo.F, title, Mon, Day, time, Rmvls.Room, Rmvls.Sd)
					}
				}
			}
		}
	}
}

/* ----------------------------------------------------------- */
/*  日集計値出力  */

var __Eeprintd_ic int

func Eeprintd(Daytm *DAYTM, Simc *SIMCONTL, Nflout int, flout []*FLOUT, Rmvls *RMVLS, Nexs int, Exs []EXSF, Soldy []float64, Eqsys *EQSYS, Wdd *WDAT) {
	var title string
	var Mon, Day, i int

	if Daytm.Ddpri != 0 {
		title = Simc.Title
		Mon = int(Daytm.Mon)
		Day = int(Daytm.Day)

		for i = 0; i < Nflout; i++ {
			flo := flout[i]

			switch flo.Idn {
			case PRTDWD:
				Wdtdprint(flo.F, title, Mon, Day, Wdd, Nexs, Exs, Soldy)
			case PRTWK:
				if __Eeprintd_ic == 0 {
					fmt.Fprintf(flo.F, "Mo Nd Day Week\n")
					__Eeprintd_ic = 1
				}

				fmt.Fprintf(flo.F, "%2d %2d %3d %s\n", Mon, Day, Daytm.day, DAYweek[Simc.Daywk[Daytm.Day]])
			case PRTDYCOMP:
				Compodyprt(flo.F, string(PRTDYCOMP), Simc, Mon, Day, Eqsys,
					Rmvls.Nrdpnl, Rmvls.Rdpnl)
			case PRTDYRM:
				Rmdyprint(flo.F, string(PRTDYRM), Simc, Mon, Day,
					Rmvls.Nroom, Rmvls.Room)
			case PRTDYHELM:
				Helmdyprint(flo.F, PRTDYHELM, Simc, Mon, Day,
					Rmvls.Nroom, Rmvls.Room, &Rmvls.Qetotal)
			case PRTDQR:
				Dyqrmprint(flo.F, title, Mon, Day,
					Rmvls.Room, Rmvls.Trdav, Rmvls.Qrmd)
			case PRTDYSF:
				Dysfprint(flo.F, title, Mon, Day, Rmvls.Room)
			}
			// #endif
		}
	}
}

/* ----------------------------------------------------------- */
/*  月集計値出力  */

func Eeprintm(daytm *DAYTM, simc *SIMCONTL, nflout int, flout []*FLOUT, rmvls *RMVLS, nexs int, exs []EXSF, solmon []float64, eqsys *EQSYS, wdm *WDAT) {
	var title string
	var mon, day int
	title = simc.Title
	mon = daytm.Mon
	day = daytm.Day
	if daytm.Ddpri != 0 {
		for i := 0; i < nflout; i++ {
			flo := flout[i]
			if flo.Idn == PRTMWD {
				Wdtmprint(flo.F, title, mon, day, wdm, nexs, exs, solmon)
			} else if flo.Idn == PRTMNCOMP {
				Compomonprt(flo.F, string(PRTMNCOMP), simc, mon, day, eqsys, rmvls.Nrdpnl, rmvls.Rdpnl)
			} else if flo.Idn == PRTMNRM {
				Rmmonprint(flo.F, string(PRTMNRM), simc, mon, day, rmvls.Nroom, rmvls.Room)
			}
		}
	}
}

/* ----------------------------------------------------------- */
/*  月－時刻集計値出力  */

func Eeprintmt(simc *SIMCONTL, nflout int, flout []*FLOUT, eqsys *EQSYS, nrdpnl int, rdpnl []RDPNL) {
	for i := 0; i < nflout; i++ {
		flo := flout[i]
		if flo.Idn == PRTMTCOMP {
			Compomtprt(flo.F, string(PRTMNCOMP), simc, eqsys, nrdpnl, rdpnl)
		}
	}
}

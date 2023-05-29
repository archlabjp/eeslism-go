package main

import (
	"fmt"
)

/* ---------------------------------------------------------------- */
/* 毎時計算値出力　*/

func Eeprinth(Daytm *DAYTM, Simc *SIMCONTL, flout []*FLOUT, Rmvls *RMVLS, Exsfst *EXSFS, Nmpath int, Mpath []MPATH, Eqsys *EQSYS, Wd *WDAT) {
	if Daytm.Ddpri != 0 && Simc.Dayprn[Daytm.Day] != 0 {
		title := Simc.Title
		Mon := Daytm.Mon
		Day := Daytm.Day
		time := Daytm.Time

		for i, flo := range flout {

			if DEBUG {
				fmt.Printf("Eeprinth MAX=%d flo[%d]=%s\n", len(flout), i, flo.Idn)
			}

			switch flo.Idn {
			case PRTHWD:
				if DEBUG {
					fmt.Println("<Eeprinth> xprsolrd")
				}
				// 気象データの出力
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
					case PRTREV:
						// 毎時室温、MRTの出力
						Rmevprint(flo.F, title, Rmvls.Room, Mon, Day, time)
					case PRTHROOM:
						// 放射パネルの出力
						Rmpnlprint(flo.F, string(PRTHROOM), Simc, Mon, Day, time, Rmvls.Nroom, Rmvls.Room)
					case PRTHELM:
						// 要素別熱損失・熱取得
						Helmprint(flo.F, string(PRTHELM), Simc, Mon, Day, time, Rmvls.Nroom, Rmvls.Room, &Rmvls.Qetotal)
					case PRTHELMSF:
						// 要素別熱損失・熱取得
						Helmsurfprint(flo.F, string(PRTHELMSF), Simc, Mon, Day, time, Rmvls.Nroom, Rmvls.Room)
					case PRTPMV:
						// PMV計算
						Pmvprint(flo.F, title, Rmvls.Room, Mon, Day, time)
					case PRTQRM:
						// 日射、室内熱取得の出力
						Qrmprint(flo.F, title, Mon, Day, time, Rmvls.Room, Rmvls.Qrm)
					case PRTRSF:
						// 室内表面温度の出力
						Rmsfprint(flo.F, title, Mon, Day, time, Rmvls.Room, Rmvls.Sd)
					case PRTSHD:
						// 日よけの影面積の出力
						Shdprint(flo.F, title, Mon, Day, time, Rmvls.Nsrf, Rmvls.Sd)
					case PRTWAL:
						// 壁体内部温度の出力
						Wallprint(flo.F, title, Mon, Day, time, Rmvls.Nsrf, Rmvls.Sd)
					case PRTPCM:
						// 潜熱蓄熱材の状態値の出力
						PCMprint(flo.F, title, Mon, Day, time, Rmvls.Nsrf, Rmvls.Sd)
					case PRTSFQ:
						// 室内表面熱流の出力
						Rmsfqprint(flo.F, title, Mon, Day, time, Rmvls.Room, Rmvls.Sd)
					case PRTSFA:
						// 室内表面熱伝達率の出力
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

func Eeprintd(Daytm *DAYTM, Simc *SIMCONTL, flout []*FLOUT, Rmvls *RMVLS, Nexs int, Exs []EXSF, Soldy []float64, Eqsys *EQSYS, Wdd *WDAT) {
	if Daytm.Ddpri != 0 {
		title := Simc.Title
		Mon := int(Daytm.Mon)
		Day := int(Daytm.Day)

		for _, flo := range flout {
			switch flo.Idn {
			case PRTDWD:
				// 気象データ日集計値出力
				Wdtdprint(flo.F, title, Mon, Day, Wdd, Nexs, Exs, Soldy)
			case PRTWK:
				// 計算年月日出力
				if __Eeprintd_ic == 0 {
					fmt.Fprintf(flo.F, "Mo Nd Day Week\n")
					__Eeprintd_ic = 1
				}

				fmt.Fprintf(flo.F, "%2d %2d %3d %s\n", Mon, Day, Daytm.DayOfYear, DAYweek[Simc.Daywk[Daytm.Day]])
			case PRTDYCOMP:
				// システム要素機器の日集計結果出力
				Compodyprt(flo.F, string(PRTDYCOMP), Simc, Mon, Day, Eqsys, Rmvls.Nrdpnl, Rmvls.Rdpnl)
			case PRTDYRM:
				// 部屋ごとの熱集計結果出力
				Rmdyprint(flo.F, string(PRTDYRM), Simc, Mon, Day, Rmvls.Nroom, Rmvls.Room)
			case PRTDYHELM:
				// 要素別熱損失・熱取得（日積算値出力）
				Helmdyprint(flo.F, string(PRTDYHELM), Simc, Mon, Day, Rmvls.Nroom, Rmvls.Room, &Rmvls.Qetotal)
			case PRTDQR:
				// 日射、室内熱取得の出力
				Dyqrmprint(flo.F, title, Mon, Day, Rmvls.Room, Rmvls.Trdav, Rmvls.Qrmd)
			case PRTDYSF:
				// 日積算壁体貫流熱取得の出力
				Dysfprint(flo.F, title, Mon, Day, Rmvls.Room)
			}
		}
	}
}

/* ----------------------------------------------------------- */
/*  月集計値出力  */

func Eeprintm(daytm *DAYTM, simc *SIMCONTL, flout []*FLOUT, rmvls *RMVLS, nexs int, exs []EXSF, solmon []float64, eqsys *EQSYS, wdm *WDAT) {
	var title string
	var mon, day int
	title = simc.Title
	mon = daytm.Mon
	day = daytm.Day
	if daytm.Ddpri != 0 {
		for _, flo := range flout {
			switch flo.Idn {
			case PRTMWD:
				// 気象データ月集計値出力
				Wdtmprint(flo.F, title, mon, day, wdm, nexs, exs, solmon)
			case PRTMNCOMP:
				// システム要素機器の月集計結果出力
				Compomonprt(flo.F, string(PRTMNCOMP), simc, mon, day, eqsys, rmvls.Nrdpnl, rmvls.Rdpnl)
			case PRTMNRM:
				// 部屋ごとの熱集計結果出力
				Rmmonprint(flo.F, string(PRTMNRM), simc, mon, day, rmvls.Nroom, rmvls.Room)
			}
		}
	}
}

/* ----------------------------------------------------------- */
/*  月－時刻集計値出力  */

func Eeprintmt(simc *SIMCONTL, flout []*FLOUT, eqsys *EQSYS, nrdpnl int, rdpnl []RDPNL) {
	for _, flo := range flout {
		if flo.Idn == PRTMTCOMP {
			Compomtprt(flo.F, string(PRTMNCOMP), simc, eqsys, nrdpnl, rdpnl)
		}
	}
}

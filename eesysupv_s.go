package main

import (
	"fmt"
)

func Sysupv(Nmpath int, Mpath []MPATH, Rmvls *RMVLS) {
	dayprn := 0
	var Rdpnl *RDPNL
	var Nrdpnl int
	var up *ELOUT

	for m := 0; m < Nmpath; m++ {
		mpath := &Mpath[m]
		/* 停止要素のシステム方程式からの除外 */

		if DEBUG {
			fmt.Printf("\n\n<< Sysupv >> m=%d  MAX=%d\n", m, Nmpath)
		}

		for i := 0; i < mpath.Nlpath; i++ {
			plist := &mpath.Plist[i]

			if DEBUG {
				fmt.Printf("\n<<Sysupv>  i=%d  iMAX=%d\n", i, mpath.Nlpath)
				fmt.Printf("OFF_SW=%c  Plist->control=%c\n", OFF_SW, plist.Control)
			}
			if dayprn != 0 && Ferr != nil {
				fmt.Fprintf(Ferr, "\n<<Sysupv>  i=%d  iMAX=%d\n", i, mpath.Nlpath)
				fmt.Fprintf(Ferr, "OFF_SW=%c  Plist->control=%c\n", OFF_SW, plist.Control)
			}
			if plist.Control != OFF_SW {
				var pelmStartIdx int
				if plist.Type == DIVERG_LPTP || plist.Type == CONVRG_LPTP {
					pelmStartIdx = 1
				} else {
					pelmStartIdx = 0
				}
				pelm := plist.Pelm[pelmStartIdx]
				if pelm.Out != nil && pelm.Out.Control != FLWIN_SW && pelm.In != nil {
					up = pelm.In.Upo
				}
				plist.Plmvb = nil
				for j := pelmStartIdx; j < plist.Nelm; j++ {
					if DEBUG {
						fmt.Printf("\n<< sysupv >> pelm=%d %s  MAX=%d\n", j, pelm.Cmp.Name, plist.Nelm)
						if pelm.Out != nil {
							fmt.Printf("<< Sysupv >> Pelm->out->control=%c\n", pelm.Out.Control)
						}
					}
					if dayprn != 0 && Ferr != nil {
						fmt.Fprintf(Ferr, "\n<< sysupv >> pelm=%d %s  MAX=%d\n", j, pelm.Cmp.Name, plist.Nelm)
						if pelm.Out != nil {
							fmt.Fprintf(Ferr, "<< Sysupv >> Pelm->out->control=%c\n", pelm.Out.Control)
						}
					}
					if pelm.Out == nil {
						pelm.In.Upv = up
						if plist.Plmvb == nil {
							plist.Plmvb = pelm
						}
					} else if pelm.Out.Control != OFF_SW {
						if DEBUG {
							fmt.Printf("<<<<<< Pelm->out->control=%c FLWIN_SW=%c\n", pelm.Out.Control, FLWIN_SW)
						}
						if dayprn != 0 && Ferr != nil {
							fmt.Fprintf(Ferr, "<<<<<< Pelm->out->control=%c FLWIN_SW=%c\n", pelm.Out.Control, FLWIN_SW)
						}
						if pelm.Out.Control == FLWIN_SW {
							up = pelm.Out
						} else {
							if pelm.In != nil {
								pelm.In.Upv = up
							}
							if DEBUG {
								fmt.Printf("<< Sysupv >> pelm=%s up=%s\n", pelm.Cmp.Name, pelm.In.Upv.Cmp.Name)
							}
							up = pelm.Out
							if plist.Plmvb == nil {
								plist.Plmvb = pelm
							}
						}
					} else if plist.Batch == 'y' && j == 0 {
						up = pelm.Out
					} else {
						if DEBUG {
							fmt.Printf("<Sysupv> 1\n")
						}
						if pelm.In != nil {
							pelm.In.Upv = nil
						}
					}
				}
			}
		}

		/* 分岐要素のシステム方程式からの除外 */

		for i := 0; i < mpath.Nlpath; i++ {
			Plist := &mpath.Plist[i]
			if DEBUG {
				fmt.Printf("  Sysupv  BRC  i=%d\n", i)
			}
			if Plist.Type == DIVERG_LPTP {
				Pelm := Plist.Pelm[0]
				Pelm.Out.Control = OFF_SW
				up := Plist.Pelm[0].Cmp.Elins[0].Upv
				if Pelm = Plist.Plmvb; Pelm != nil {
					Pelm.In.Upv = up
				}
			}
		}

	}
	if SIMUL_BUILDG {
		/*********************************/
		/* 放射パネル設置室についてのパネル上流要素 */

		for i := 0; i < Rmvls.Nroom; i++ {
			room := &Rmvls.Room[i]
			for j := 0; j < room.Nrp; j++ {
				rmpnl := &room.rmpnl[j]
				elin := room.cmp.Elins[room.Nachr+room.Ntr+j]
				elp := rmpnl.pnl.cmp.Elins[0]
				elin.Upv = elp.Upv
			}
		}

		for i := 0; i < Nrdpnl; i++ {
			Rdpnl = &Rmvls.Rdpnl[i]

			for j := 0; j < Rdpnl.MC; j++ {
				rm := Rdpnl.rm[j]
				rmpnl := &rm.rmpnl[j]
				elin := Rdpnl.cmp.Elins[Rdpnl.elinpnl[j]]
				for jj := 0; jj < Rdpnl.Nrp[j]; jj++ {
					elin.Upv = rmpnl.pnl.cmp.Elins[0].Upv
				}
			}
		}

		if DEBUG {
			fmt.Printf("  Sysupv end  ========\n")
		}
	}
}

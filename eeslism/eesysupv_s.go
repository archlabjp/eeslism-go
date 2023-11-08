package eeslism

import (
	"fmt"
)

// Upo, Upv の書き換え
// NOTE: おそらく、 Upoは経路要素における上流の要素を指す。
//       Upvは計算時に参照すべき上流要素を指す。多くの場合は Upo == Upv だと考えらえる。
func Sysupv(Mpath []*MPATH, Rmvls *RMVLS) {
	var Rdpnl *RDPNL
	var Nrdpnl int
	var up *ELOUT

	for m, mpath := range Mpath {
		/* 停止要素のシステム方程式からの除外 */

		if DEBUG {
			fmt.Printf("\n\n<< Sysupv >> m=%d  MAX=%d\n", m, len(Mpath))
		}

		for i, plist := range mpath.Plist {

			if DEBUG {
				fmt.Printf("\n<<Sysupv>  i=%d  iMAX=%d\n", i, len(mpath.Plist))
				fmt.Printf("OFF_SW=%c  Plist->control=%c\n", OFF_SW, plist.Control)
			}
			if dayprn && Ferr != nil {
				fmt.Fprintf(Ferr, "\n<<Sysupv>  i=%d  iMAX=%d\n", i, len(mpath.Plist))
				fmt.Fprintf(Ferr, "OFF_SW=%c  Plist->control=%c\n", OFF_SW, plist.Control)
			}

			if plist.Control != OFF_SW {
				// 末端経路が停止してなければ:
				//

				// 末端経路が分岐・合流の場合は、最初の要素を無視する
				var pelmStartIdx int
				if plist.Type == DIVERG_LPTP || plist.Type == CONVRG_LPTP {
					pelmStartIdx = 1
				} else {
					pelmStartIdx = 0
				}

				// Testcode
				pelm := plist.Pelm[pelmStartIdx]
				if pelm.Out != nil && pelm.Out.Control != FLWIN_SW && pelm.In != nil {
					up = pelm.In.Upo
				}
				plist.Plmvb = nil

				// 末端経路内の要素のループ
				for j := pelmStartIdx; j < len(plist.Pelm); j++ {
					pelm = plist.Pelm[j]
					if DEBUG {
						fmt.Printf("\n<< sysupv >> pelm=%d %s  MAX=%d\n", j, pelm.Cmp.Name, len(plist.Pelm))
						if pelm.Out != nil {
							fmt.Printf("<< Sysupv >> Pelm->out->control=%c\n", pelm.Out.Control)
						}
					}
					if dayprn && Ferr != nil {
						fmt.Fprintf(Ferr, "\n<< sysupv >> pelm=%d %s  MAX=%d\n", j, pelm.Cmp.Name, len(plist.Pelm))
						if pelm.Out != nil {
							fmt.Fprintf(Ferr, "<< Sysupv >> Pelm->out->control=%c\n", pelm.Out.Control)
						}
					}
					if pelm.Out == nil {
						pelm.In.Upv = up
						if plist.Plmvb == nil {
							// 末端経路内で出口がない要素のうち経路内で最上流のもの
							plist.Plmvb = pelm
						}
					} else if pelm.Out.Control != OFF_SW {
						if DEBUG {
							fmt.Printf("<<<<<< Pelm->out->control=%c FLWIN_SW=%c\n", pelm.Out.Control, FLWIN_SW)
						}
						if dayprn && Ferr != nil {
							fmt.Fprintf(Ferr, "<<<<<< Pelm->out->control=%c FLWIN_SW=%c\n", pelm.Out.Control, FLWIN_SW)
						}
						if pelm.Out.Control == FLWIN_SW {
							up = pelm.Out
						} else {
							if DEBUG {
								fmt.Printf("up->cmp->name=%s\n", up.Cmp.Name)
							}

							if dayprn && Ferr != nil {
								fmt.Fprintf(Ferr, "up->cmp->name=%s\n", up.Cmp.Name)
							}

							// Testcode
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
					} else if plist.Batch && j == 0 {
						up = pelm.Out
					} else {
						if DEBUG {
							fmt.Printf("<Sysupv> 1\n")
						}
						if pelm.In != nil {
							pelm.In.Upv = nil
						}
						if DEBUG {
							fmt.Printf("<Sysupv> 2\n")
						}
					}
				}

				if mpath.Type == CIR_PTYP {
					// 合流経路の場合:
					// ex) `> G4 G5 >` の場合、G5が停止していれば、

					ptermel := plist.Pelm[len(plist.Pelm)-1] // 末端要素
					if ptermel.Out.Control == OFF_SW {
						ptermel = plist.Plmvb
						ptermel.In.Upv = up
					}
				}
			} else {
				// 末端経路が停止している場合: 入力無し
				//
				for _, Pelm := range plist.Pelm {
					if Pelm.In != nil {
						Pelm.In.Upv = nil
					}
				}
			}
		}

		/* 分岐要素のシステム方程式からの除外 */

		for i, Plist := range mpath.Plist {
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

		for i := range Rmvls.Room {
			room := Rmvls.Room[i]
			for j := 0; j < room.Nrp; j++ {
				rmpnl := room.rmpnl[j]
				elin := room.cmp.Elins[room.Nachr+room.Ntr+j]
				elp := rmpnl.pnl.cmp.Elins[0]
				elin.Upv = elp.Upv
			}
		}

		for i := 0; i < Nrdpnl; i++ {
			Rdpnl = Rmvls.Rdpnl[i]

			for j := 0; j < Rdpnl.MC; j++ {
				rm := Rdpnl.rm[j]
				rmpnl := rm.rmpnl[j]
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

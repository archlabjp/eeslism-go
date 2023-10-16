package eeslism

import (
	"fmt"
	"math"
	"os"
)

/* --------------------------------------------------- */

/* システム要素の流量設定 */

func Pflow(Nmpath int, _Mpath []MPATH, Wd *WDAT) {
	var m, i, j, n, NG int
	var mpi *MPATH
	var Plist, pl *PLIST
	var Pelm *PELM
	var eli *ELIN
	var elo *ELOUT
	var cmp *COMPNT
	var vc, vcmb *VALV
	var Go float64
	var G float64
	var Err, s string
	/*---- Satoh Debug VAV  2000/12/6 ----*/
	var vav *VAV
	var G0 float64

	if Nmpath > 0 {
		for m = 0; m < Nmpath; m++ {
			Mpath := _Mpath[m]

			if DEBUG {
				fmt.Printf("m=%d mMAX=%d name=%s\n", m, Nmpath, Mpath.Name)
			}

			// 流量が既知の末端流量の初期化
			for i = 0; i < Mpath.Nlpath; i++ {
				Plist := &Mpath.Plist[i]
				if Plist.Go != nil && Plist.Nvalv == 0 {
					Plist.G = *Plist.Go
				}
			}

			for i = 0; i < Mpath.Nlpath; i++ {
				Plist := &Mpath.Plist[i]
				Plist.G = 0.0

				if DEBUG {
					fmt.Printf("i=%d iMAX=%d name=%s\n", i, Mpath.Nlpath, Plist.Name)
				}

				// 流量が既知の末端経路
				if Plist.Go != nil && Plist.Nvalv == 0 {
					Plist.G = *Plist.Go
				} else if Plist.Go != nil && Plist.Nvalv > 0 ||
					Plist.NOMVAV > 0 ||
					(Plist.Go == nil && Plist.Nvalv > 0 && Plist.UnknownFlow == 1) {
					if Plist.Go != nil && Plist.Valv != nil &&
						Plist.Valv.Cmp.Eqptype == VALV_TYPE {
						// 二方弁の計算
						Plist.G = *Plist.Go
						vc = Plist.Valv
						if vc == nil || vc.Org == 'y' {
							if vc.X < 0.0 {
								s = fmt.Sprintf("%s のバルブ開度 %f が不正です。", vc.Name, vc.X)
								Eprint("<Pflow>", s)
							}
							Plist.G = vc.X * *Plist.Go
						} else {
							vcmb = vc.Cmb.Eqp.(*VALV)
							Plist.G = (1.0 - vcmb.X) * *Plist.Go
						}
					} else if Plist.Valv != nil && Plist.Valv.MGo != nil &&
						*Plist.Valv.MGo > 0.0 && Plist.Control != OFF_SW {
						// 三方弁の計算

						vc = Plist.Valv
						vcmb = vc.Cmb.Eqp.(*VALV)

						if vc.Org == 'y' {
							Plist.G = vc.X * *vc.MGo
						} else {
							Plist.G = (1.0 - vcmb.X) * *vc.MGo
						}

						if Plist.G > 0. {
							Plist.Control = ON_SW
						}
					} else if Plist.Valv != nil && Plist.Valv.MGo != nil && *Plist.Valv.MGo <= 0.0 {
						Plist.G = 0.0
					} else if Plist.Valv != nil && Plist.Valv.Count > 0 {
						Plist.G = Plist.Gcalc
					} else if Plist.NOMVAV > 0 {
						Plist.G = OMflowcalc(Plist.OMvav, Wd)
					}

					if Plist.G <= 0.0 {
						lpathscdd(OFF_SW, Plist)
					}

					if Plist.G > 0. {
						Plist.Control = ON_SW
					}
				} else if Plist.Nvav > 0 {
					/*---- Satoh Debug VAV  2000/12/6 ----*/

					/* VAVユニット時の流量 */

					G = -999.0
					for j = 0; j < Plist.Nelm; j++ {
						Pelm = Plist.Pelm[j]

						if Pelm.Cmp.Eqptype == VAV_TYPE ||
							Pelm.Cmp.Eqptype == VWV_TYPE {
							vav = Pelm.Cmp.Eqp.(*VAV)

							if vav.Count == 0 {
								G = math.Max(G, vav.Cat.Gmax)
							} else {
								G = math.Max(G, vav.G)
							}
						}
					}
					Plist.G = G
				} else if Plist.Rate != nil {
					Plist.G = *Mpath.G0 * *Plist.Rate
				} else if !Plist.Batch {
					if Plist.Go != nil {
						Go = *Plist.Go
					} else {
						Go = 0.0
					}

					if Plist.Pelm != nil {
						var l int
						for l = 0; l < len(mpi.Plist); l++ {
							if &mpi.Plist[l] == Plist {
								break
							}
						}
						Err = fmt.Sprintf("Mpath=%s  lpath=%d  elm=%s  Go=%f\n", Mpath.Name, l, Plist.Pelm[0].Cmp.Name, Go)
					}
				}
			}

			NG = Mpath.NGv

			X := make([]float64, NG)
			Y := make([]float64, NG)
			A := make([]float64, NG*NG)

			for i = 0; i < NG; i++ {
				if DEBUG {
					fmt.Printf("i=%d iMAX=%d\n", i, NG)
				}

				cmp = Mpath.Cbcmp[i]

				if DEBUG {
					fmt.Printf("<Pflow> Name=%s\n", cmp.Name)
				}

				for j = 0; j < cmp.Nin; j++ {
					eli = cmp.Elins[j]

					if DEBUG {
						fmt.Printf("j=%d jMAX=%d\n", j, cmp.Nin)
					}

					if eli.Lpath.Go != nil ||
						eli.Lpath.Nvav != 0 ||
						eli.Lpath.Nvalv != 0 ||
						eli.Lpath.Rate != nil ||
						eli.Lpath.NOMVAV != 0 {
						Y[i] -= eli.Lpath.G
					} else {
						n = eli.Lpath.N

						if n < 0 || n >= NG {
							Err = fmt.Sprintf("n=%d", n)
							Eprint("<Pflow>", Err)
							os.Exit(EXIT_PFLOW)
						}

						A[i*NG+n] = 1.0
					}
				}

				////////

				for j = 0; j < cmp.Nout; j++ {
					elo = cmp.Elouts[j]

					if elo.Lpath.Go != nil ||
						elo.Lpath.Nvav != 0 ||
						elo.Lpath.Nvalv != 0 ||
						elo.Lpath.Rate != nil {
						Y[i] += elo.Lpath.G
					} else {
						n = elo.Lpath.N

						if n < 0 || n >= NG {
							Err = fmt.Sprintf(Err, "n=%d", n)
							Eprint("<Pflow>", Err)
							os.Exit(EXIT_PFLOW)
						}

						A[i*NG+n] = -1.0
					}
				}
			}

			if NG > 0 {

				if DEBUG {
					for i = 0; i < NG; i++ {
						fmt.Printf("%s\t", Mpath.Cbcmp[i].Name)

						for j = 0; j < NG; j++ {
							fmt.Printf("%6.1f", A[i*NG+j])
						}

						fmt.Printf("\t%.5f\n", Y[i])
					}
				}

				if dayprn && Ferr != nil {
					for i = 0; i < NG; i++ {
						fmt.Fprintf(Ferr, "%s\t", Mpath.Cbcmp[i].Name)

						for j = 0; j < NG; j++ {
							fmt.Fprintf(Ferr, "\t%.1g", A[i*NG+j])
						}

						fmt.Fprintf(Ferr, "\t\t%.2g\n", Y[i])
					}
				}

				if NG > 1 {
					Matinv(A, NG, NG, "<Pflow>")
					Matmalv(A, Y, NG, NG, X)
				} else {
					X[0] = Y[0] / A[0]
				}

				if DEBUG {
					fmt.Printf("<Pflow>  Flow Rate\n")
					for i = 0; i < NG; i++ {
						fmt.Printf("\t%6.2f\n", X[i])
					}
				}

				if dayprn && Ferr != nil {
					for i = 0; i < NG; i++ {
						fmt.Fprintf(Ferr, "\t\t%.2g\n", X[i])
					}
				}
			}

			for i = 0; i < NG; i++ {
				pl = Mpath.Pl[i]
				pl.G = X[i]
			}

			for i = 0; i < Mpath.Nlpath; i++ {
				Plist = &Mpath.Plist[i]

				if DEBUG {
					fmt.Printf("<< Pflow >> e i=%d iMAX=%d control=%c G=%g\n",
						i, Mpath.Nlpath, Plist.Control, Plist.G)
				}

				if Plist.Control == OFF_SW {
					Plist.G = 0.0
				} else if Plist.G <= 0.0 {
					// 負であればエラーを表示する
					//if (Plist.G < 0. )
					//	fmt.Printf("<%s>  流量が負になっています %g\n", Mpath.Name, Plist.G ) ;

					Plist.G = 0.0
					Plist.Control = OFF_SW
					lpathscdd(Plist.Control, Plist)
				}

				for j = 0; j < Plist.Nelm; j++ {
					Pelm = Plist.Pelm[j]

					if Pelm.Out != nil {
						Pelm.Out.G = Plist.G
					}

					if DEBUG {
						if Pelm.Out != nil {
							G0 = Pelm.Out.G
						} else {
							G0 = 0.0
						}

						fmt.Printf("< Pflow > j=%d\tjMAX=%d\tPelm-G=%g\tPlist.G=%g\n",
							j, Plist.Nelm, G0, Plist.G)
					}
				}
			}
		}
	}
}

package eeslism

import (
	"fmt"
)

// システム方程式の作成およびシステム変数の計算
func Syseqv(_Elout []*ELOUT, Syseq *SYSEQ) {
	var eleq, elosv []*ELOUT
	var sysmcf, syscv, Y []float64
	var i, m, n, Nsv int
	var mrk []rune

	Syseq.A = ' '
	eleq = nil
	elosv = nil
	sysmcf = nil
	syscv = nil
	Y = nil

	Nelout := len(_Elout)

	if Nelout > 0 {
		eleq = make([]*ELOUT, Nelout)

		for i = 0; i < Nelout; i++ {
			eleq[i] = nil
		}
	}

	if Nelout > 0 {
		elosv = make([]*ELOUT, Nelout)

		for i = 0; i < Nelout; i++ {
			elosv[i] = nil
		}
	}

	mrk = make([]rune, Nelout)

	for i, Elout := range _Elout {

		if DEBUG {
			fmt.Printf("xxx syseqv  Eo name=%s control=%c sysld=%c i=%d MAX=%d\n",
				Elout.Cmp.Name, Elout.Control, Elout.Sysld, i, Nelout)
		}

		if dayprn && Ferr != nil {
			fmt.Fprintf(Ferr, "xxx syseqv  Eo name=%s control=%c sysld=%c i=%d MAX=%d\n",
				Elout.Cmp.Name, Elout.Control, Elout.Sysld, i, Nelout)
		}

		if Elout.Control != LOAD_SW &&
			Elout.Control != FLWIN_SW &&
			Elout.Control != BATCH_SW {
			Elout.Sv = -1
			Elout.Sysv = 0.0
		}

		if Elout.Control == ON_SW {
			if DEBUG {
				fmt.Printf("ON_SW = [i=%d m=%d n=%d] %s  G=%f\n", i, m, n, Elout.Cmp.Name, Elout.G)
			}

			eleq[m] = Elout
			elosv[n] = Elout
			mrk[n] = SYSV_EQV
			Elout.Sv = n
			Elout.Sld = -1
			n++

			if Elout.Sysld == 'y' {
				elosv[n] = Elout
				mrk[n] = LOAD_EQV
				Elout.Sld = n
				n++
			}
			m++
		} else if Elout.Control == LOAD_SW && Elout.Sysld != 'y' {
			eleq[m] = Elout
			Elout.Sv = -1
			Elout.Sld = -1
			m++
		}
	}
	Nsv = n

	sysmcf = make([]float64, Nsv*Nsv)
	syscv = make([]float64, Nsv)
	Y = make([]float64, Nsv)

	for i = 0; i < Nsv; i++ {
		elout := eleq[i]
		a := sysmcf[Nsv*i : Nsv*i+Nsv]
		b := syscv[i : i+1]

		if DEBUG {
			fmt.Printf("xxx syseqv Elout=%d %s Ni=%d cfo=%f\n",
				i, elout.Cmp.Name, elout.Ni, elout.Coeffo)
		}

		if dayprn && Ferr != nil {
			fmt.Fprintf(Ferr, "xxx syseqv Elout=%d %s Ni=%d cfo=%f\n",
				i, elout.Cmp.Name, elout.Ni, elout.Coeffo)
		}

		c := a
		matinit(c, Nsv)

		b[0] = elout.Co

		if n = elout.Sv; n >= 0 {
			a[n] = elout.Coeffo
			if nn := elout.Sld; nn >= 0 {
				a[nn] = -1.0
			}
		} else {
			b[0] -= elout.Coeffo * elout.Sysv
		}

		for j := 0; j < elout.Ni; j++ {
			elin := elout.Elins[j]
			cfin := elout.Coeffin[j]
			elov := elin.Upv

			// DEBUG: RDPNLの全入力要素の状態を出力
			if DEBUG_RDPNL_COEFF && elout.Cmp != nil && elout.Cmp.Eqptype == RDPANEL_TYPE {
				elovName := "nil"
				elovCtrl := "-"
				if elov != nil {
					elovName = elov.Cmp.Name
					elovCtrl = string(elov.Control)
				}
				fmt.Printf("DEBUG Go RDPNL elin[%d]: elov=%s control=%s cfin=%.6f\n", j, elovName, elovCtrl, cfin)
			}

			if elov != nil {
				if DEBUG {
					fmt.Printf("xxx syseqv Elout=%d %s  in=%d elov=%s  control=%c sys=%f\n",
						i, elout.Cmp.Name, j,
						elov.Cmp.Name, elov.Control, elov.Sysv)
				}

				if dayprn && Ferr != nil {
					fmt.Fprintf(Ferr, "xxx syseqv Elout=%d %s  in=%d elov=%s  control=%c sys=%f\n",
						i, elout.Cmp.Name, j,
						elov.Cmp.Name, elov.Control, elov.Sysv)
				}

				if elov.Control == ON_SW {
					n = elin.Upv.Sv
					// DEBUG: RDPNLの係数加算を詳細出力
					if DEBUG_RDPNL_COEFF && elout.Cmp != nil && elout.Cmp.Eqptype == RDPANEL_TYPE {
						fmt.Printf("DEBUG Go Syseqv RDPNL input j=%d: n=%d cfin=%.10f elov=%s\n", j, n, cfin, elov.Cmp.Name)
					}
					a[n] += cfin
				} else if elov.Control == LOAD_SW ||
					elov.Control == FLWIN_SW ||
					elov.Control == BATCH_SW {
					if DEBUG {
						fmt.Printf("xxx syseqv elov=%s  control=%c sys=%f\n",
							elov.Cmp.Name, elov.Control, elov.Sysv)
					}

					if dayprn && Ferr != nil {
						fmt.Fprintf(Ferr, "xxx syseqv elov=%s  control=%c sys=%f\n",
							elov.Cmp.Name, elov.Control, elov.Sysv)
					}

					b[0] -= cfin * elov.Sysv
				}
			}
		}
		if DEBUG {
			fmt.Printf("xx syseqv  i=%d  b=%f\n", i, b[0])
		}
	}

	/********* 連立方程式 ***********/

	if DEBUG {
		Seqprint("%.6g\t", Nsv, sysmcf, "%.6g", syscv)

		//for ( i = 0; i < Nsv; i++ )
		//	fmt.Printf ( "%g\n", sysmcf[i+Nsv*7] ) ;
	}

	// DEBUG: パネル行のマトリクス値を出力
	if DEBUG_RDPNL_COEFF && Nsv > 0 {
		for i := 0; i < m; i++ {
			if eleq[i] != nil && eleq[i].Cmp != nil {
				if eleq[i].Cmp.Eqptype == RDPANEL_TYPE {
					fmt.Printf("DEBUG Go Syseqv RDPNL row %d: ", i)
					for j := 0; j < Nsv && j < 5; j++ {
						fmt.Printf("a[%d]=%.10f ", j, sysmcf[Nsv*i+j])
					}
					fmt.Printf("b=%.10f\n", syscv[i])
				}
			}
		}
	}

	if Nsv > 0 {
		Matinv(sysmcf, Nsv, Nsv, "<Syseqv>")
		Matmalv(sysmcf, syscv, Nsv, Nsv, Y)
	}

	for i = 0; i < Nsv; i++ {
		if mrk[i] == SYSV_EQV {
			elosv[i].Sysv = Y[i]
		} else if mrk[i] == LOAD_EQV {
			elosv[i].Load = Y[i]
		}
		if DEBUG {
			fmt.Printf("%d: %s = %f\n", i, elosv[i].Cmp.Name, Y[i])
		}
	}

	if DEBUG {
		for i = 0; i < Nsv; i++ {
			fmt.Printf("Y[%d]=%.5f  mrk=%c  Elo=%s\n",
				i, Y[i], mrk[i], elosv[i].Cmp.Name)
		}
		fmt.Printf("\n")
	}

	if dayprn && Ferr != nil {
		for i = 0; i < Nsv; i++ {
			fmt.Fprintf(Ferr, "Y[%d]=%6.3f  mrk=%c  Elo=%s\n",
				i, Y[i], mrk[i], elosv[i].Cmp.Name)
		}
		fmt.Fprintf(Ferr, "\n")
	}
}

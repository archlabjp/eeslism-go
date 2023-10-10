package eeslism

import (
	"fmt"
	"strings"
)

/* システム要素の接続経路の入力 */

func Pathdata(
	f *EeTokens,
	errkey string,
	Simc *SIMCONTL,
	Wd *WDAT,
	Ncompnt int,
	Compnt []COMPNT,
	Schdl *SCHDL,
	M *[]MPATH,
	Nmpath *int,
	Plst *[]PLIST,
	Plm *[]PELM,
	Npelm *int,
	Nplist *int,
	ID int,
	Eqsys *EQSYS,
) {
	//var Mpath *MPATH
	var mpi *MPATH
	var C *COMPNT
	var stank *STANK
	var Plist *PLIST
	var Pelm *PELM
	var Qmeas *QMEAS
	var s, ss, sss string
	var etyp EqpType
	var ci ELIOType = ' '
	var co ELIOType = ' '
	var elm, stv string

	var i, j, m, ncv, idci, idco, iswc, Np int
	var N []int
	var Nplst *int
	var k int
	var Go float64
	var Npl, Nm, Nplm int
	id := 0
	iPlist := 0

	if DEBUG {
		fmt.Printf("\n")
		for i := 0; i < Ncompnt; i++ {
			C = &Compnt[i]
			fmt.Printf("name=%s Nin=%d  Nout=%d\n", C.Name, C.Nin, C.Nout)
		}
	}

	Nm = Mpathcount(f, &Npl)
	*Nplist = Npl * 2

	if Nm > 0 {
		*M = make([]MPATH, Nm*2)

		Mpathinit(Nm*2, *M)

		for i := 0; i < Nm*2; i++ {
			mpi = &(*M)[i]
		}
	}

	if Npl > 0 {
		*Plst = make([]PLIST, Npl*2)

		Plistinit(Npl*2, *Plst)

		// for i := 0; i < Npl*2; i++ {
		// 	Pl = &(*Plst)[i]
		// }
	}

	var Mpath *MPATH
	if len(*M) > 0 {
		Mpath = &(*M)[0]
	}
	plistIdx := 0

	N = make([]int, Nm)

	Plcount(f, N)

	Nplm = Pelmcount(f)

	if Nplm > 0 {
		*Plm = make([]PELM, Nplm)
		Pelminit(Nplm, *Plm)
	}

	mpi = Mpath
	NpIdx := 0
	if len(N) > 0 {
		Nplst = &N[NpIdx]
	}
	*Npelm = 0
	pelmIdx := 0
	if len(*Plm) > 0 {
		Pelm = &(*Plm)[pelmIdx]
	}

	Mpath_idx := 0
	if ID == 0 {
		for f.IsEnd() == false {
			ss = f.GetToken()
			if ss[0] == '*' {
				break
			}

			if DEBUG {
				fmt.Printf("eepathdat.c  ss=%s\n", ss)
			}

			if *Nplst > 0 {
				Mpath.Pl = make([]*PLIST, *Nplst)
				for i := 0; i < *Nplst; i++ {
					Mpath.Pl[i] = nil
				}
			}

			if *Nplst > 0 {
				Mpath.Cbcmp = make([]*COMPNT, *Nplst*2)
				for i := 0; i < *Nplst*2; i++ {
					Mpath.Cbcmp[i] = nil
				}
			}

			Mpath.Name = ss
			Mpath.NGv = 0
			Mpath.NGv2 = 0
			Mpath.Ncv = 0
			Mpath.Nlpath = 0
			Mpath.Plist = *Plst
			//Pelmpre = nil

			for f.IsEnd() == false {
				s = f.GetToken()
				if s[0] == ';' {
					break
				}

				if DEBUG {
					fmt.Printf("eepathdat.c  s=%s\n", s)
				}

				Plist = &(*Plst)[plistIdx]

				if s[0] == '-' {
					ss = f.GetToken()

					if s[1:] == "sys" {
						Mpath.Sys = ss[0]
					} else if s[1:] == "f" {
						Mpath.Fluid = rune(ss[0])
					} else {
						Errprint(1, errkey, s)
					}
				} else if s[0] == '>' {
					sss = fmt.Sprintf("Path%d", iPlist)
					Plist.Plistname = sss
					Plist.Pelm = []*PELM{Pelm}

					for f.IsEnd() == false {
						s = f.GetToken()
						if s[0] == '>' {
							break
						}

						if DEBUG {
							fmt.Printf("eepathdat.c  s=%s\n", s)
						}

						if s[0] == '(' {
							i, err := fmt.Sscanf(s[1:], "%f", &Go)
							if err != nil {
								panic(err)
							}
							if i == 1 {
								Plist.Go = new(float64)
								*Plist.Go = Go
								if DEBUG {
									fmt.Printf("Go=%f\n", *Plist.Go)
								}
							} else {
								_, err := fmt.Sscanf(s[1:], "%[^)])", &ss)
								if err != nil {
									panic(err)
								}

								if DEBUG {
									fmt.Printf("s=%s ss=%s\n", s, ss)
								}

								if j = idsch(ss, Schdl.Sch, ""); j >= 0 {
									Plist.Go = new(float64)
									*Plist.Go = Schdl.Val[j]
								} else {
									Plist.Go = envptr(ss, Simc, Ncompnt, Compnt, Wd, nil)
								}

								if DEBUG {
									fmt.Printf("Go=%f\n", *Plist.Go)
								}
							}
						} else if s[0] == '[' {
							// 流量比率設定フラグのセット
							Mpath.Rate = 'Y'

							i, err := fmt.Scanf(s[1:], "%f", &Go)
							if err != nil {
								panic(err)
							}
							if i == 1 {
								Plist.Rate = new(float64)
								*Plist.Rate = Go
								if DEBUG {
									fmt.Printf("rate=%f\n", *Plist.Rate)
								}
							} else {
								_, err := fmt.Sscanf(s[1:], "%[^)])", &ss)
								if err != nil {
									panic(err)
								}

								if DEBUG {
									fmt.Printf("s=%s ss=%s\n", s, ss)
								}

								if j := idsch(ss, Schdl.Sch, ""); j >= 0 {
									Plist.Rate = new(float64)
									*Plist.Rate = Schdl.Val[j]
								} else {
									Plist.Rate = envptr(ss, Simc, Ncompnt, Compnt, Wd, nil)
								}

								if DEBUG {
									fmt.Printf("rate=%f\n", *Plist.Rate)
								}
							}
						} else {
							// 末端経路名称
							if s[:5] == "name=" {
								_, err := fmt.Sscanf(s, "%*[^=]=%s", &ss)
								if err != nil {
									panic(err)
								}
								Plist.Plistname = ss
							} else {
								var idx int
								if idx = strings.IndexRune(s, '/'); idx >= 0 {
									s = s[idx+1:]
								}
								if idx = strings.IndexRune(s, ':'); idx >= 0 {
									Plist.Name = s
									s = s[:idx]
									elm = s
								} else {
									if idx = strings.IndexRune(s, '['); idx >= 0 {
										co, ci = ELIOType(s[idx+1]), ELIOType(s[idx+1])
										s = s[:idx]
										elm = s
									} else {
										elm = s
										co = ELIOType(0)
										ci = ELIOType(0)
									}
								}
								err := 1

								for i := 0; i < Ncompnt; i++ {
									cmp := &Compnt[i]
									C = cmp
									if cmp.Name == elm {
										err = 0
										if cmp.Eqptype == FLIN_TYPE && Plist.Pelm[0] == Pelm {
											Plist.Type = IN_LPTP
										} else if cmp.Eqptype == VALV_TYPE || cmp.Eqptype == TVALV_TYPE {
											Plist.Nvalv++
											Plist.Valv = cmp.Eqp.(*VALV)
											Plist.Valv.Plist = Plist
										} else if cmp.Eqptype == OMVAV_TYPE {
											// Satoh OMVAV 2010/12/16
											Plist.NOMVAV++
											Plist.OMvav = cmp.Eqp.(*OMVAV)
											Plist.OMvav.Plist = Plist
										} else if cmp.Eqptype == VAV_TYPE || cmp.Eqptype == VWV_TYPE {
											/*---- Satoh Debug VAV  2000/12/6 ----*/
											Plist.Nvav++
										} else if cmp.Eqptype == QMEAS_TYPE {
											/*---- Satoh Debug QMEAS  2003/6/2 ----*/
											Qmeas = cmp.Eqp.(*QMEAS)
											if co == 'G' {
												Qmeas.G = &Plist.G
												Qmeas.PlistG = Plist
												Qmeas.Fluid = Mpath.Fluid
											} else if co == 'H' {
												Qmeas.PlistTh = Plist
												Qmeas.Nelmh = id
											} else if co == 'C' {
												Qmeas.PlistTc = Plist
												Qmeas.Nelmc = id
											}
										} else if cmp.Eqptype == STANK_TYPE {
											if stv != "" {
												Plist.Batch = 'y'
												stank = cmp.Eqp.(*STANK)
												for i := 0; i < stank.Nin; i++ {
													if stank.Pthcon[i] == co {
														if iswc = idscw(stv, Schdl.Scw, ""); iswc >= 0 {
															stank.Batchcon[i] = Schdl.Isw[iswc]
														}
													}
												}
											}
										}

										if cmp.Eqptype != VALV_TYPE && cmp.Eqptype != TVALV_TYPE &&
											cmp.Eqptype != QMEAS_TYPE && cmp.Eqptype != OMVAV_TYPE {
											(*Npelm)++
											Pelm.Out = nil
											Pelm.Cmp = cmp
											Pelm.Ci = ci
											Pelm.Co = co
											//Pelmpre = Pelm

											pelmIdx++
											Pelm = &(*Plm)[pelmIdx]
										}
										if cmp.Eqptype != VALV_TYPE && cmp.Eqptype != TVALV_TYPE &&
											cmp.Eqptype != QMEAS_TYPE && cmp.Eqptype != DIVERG_TYPE &&

											cmp.Eqptype != CONVRG_TYPE && cmp.Eqptype != DIVGAIR_TYPE &&
											cmp.Eqptype != CVRGAIR_TYPE && cmp.Eqptype != OMVAV_TYPE {
											id++
										}
										break
									}
								}

								Errprint(err, errkey, elm)

								if DEBUG {
									fmt.Printf("<<Pathdata>> Mp=%s  elm=%s Npelm=%d\n",
										Mpath.Name, elm, *Npelm)
								}
							}
						}
					}

					var n int
					for n = 0; n < len(Plist.Pelm); n++ {
						if Plist.Pelm[n] == Pelm {
							break
						}
					}
					Plist.Nelm = n
					Plist.Mpath = Mpath
					//Pelmpre = nil
					id = 0
					plistIdx++
					iPlist++
				} else {
					Errprint(1, errkey, s)
				}
			}
			var mn int
			for mn = 0; mn < len(Mpath.Plist); mn++ {
				if &Mpath.Plist[mn] == Plist {
					break
				}
			}
			Mpath.Nlpath = mn

			if DEBUG {
				var i int
				for i = 0; i < len((*M)); i++ {
					if &(*M)[i] == Mpath {
						break
					}
				}
				fmt.Printf("<<Pathdata>>  Mpath=%d fliud=%c\n", i, Mpath.Fluid)
			}

			if Mpath.Fluid == AIR_FLD {
				if DEBUG {
					fmt.Printf("<<Pathdata  a>> Mp=%s  Npelm=%d\n", Mpath.Name, *Npelm)
				}

				Mpath_idx++

				if *Nplst > 0 {
					Mpath.Pl = make([]*PLIST, *Nplst)
					for k = 0; k < *Nplst; k++ {
						Mpath.Pl[k] = nil
					}
				}

				if *Nplst > 0 {
					Mpath.Cbcmp = make([]*COMPNT, *Nplst*2)
					for k = 0; k < *Nplst*2; k++ {
						Mpath.Cbcmp[k] = nil
					}
				}

				Np = *Npelm

				// 空気系統用の絶対湿度経路へのコピー
				plistcpy(&(*M)[Mpath_idx], &(*M)[Mpath_idx-1], Npelm, (*Plm)[pelmIdx:], (*Plst)[plistIdx:], Ncompnt, Compnt)

				pelmIdx += *Npelm - Np
				Pelm = &(*Plm)[pelmIdx]
				plistIdx += Mpath.Nlpath
			}
			Mpath_idx++

			NpIdx++
			Nplst = &N[NpIdx]
		}
	}
	*Nmpath = Mpath_idx

	if DEBUG {
		if *Nmpath > 0 {
			mpi := *M
			plistprint(*Nmpath, mpi, mpi[0].Plist[0].Pelm, Compnt[0].Elouts, Compnt[0].Elins)
		}

		fmt.Printf("SYSPTH  Data Read end\n")
		fmt.Printf("Nmpath=%d\n", *Nmpath)
	}

	/* ============================================================================ */

	Mpath_idx = 0
	for i := 0; i < *Nmpath; i++ {
		if DEBUG {
			fmt.Printf("1----- MAX=%d  i=%d\n", *Nmpath, i)
		}

		ncv = 0

		Mpath = &(*M)[Mpath_idx]

		for j = 0; j < Mpath.Nlpath; j++ {
			Plist = &Mpath.Plist[j]

			if DEBUG {
				var MpathPos, PlistPos int
				for MpathPos = 0; MpathPos < len((*M)); MpathPos++ {
					if &(*M)[MpathPos] == Mpath {
						break
					}
				}
				for PlistPos = 0; PlistPos < len(Mpath.Plist); PlistPos++ {
					if &Mpath.Plist[PlistPos] == Plist {
						break
					}
				}
				fmt.Printf("eepath.c  Mpath.Nlpath=%d\n", Mpath.Nlpath)
				fmt.Printf("<<Pathdata>>  i=%d Mpath=%d  j=%d Plist=%d\n", i, MpathPos, j, PlistPos)
			}

			for m = 0; m < Plist.Nelm; m++ {
				Pelm = Plist.Pelm[m]

				idci = 1
				idco = 1
				etyp = Pelm.Cmp.Eqptype

				if m == 0 && etyp == FLIN_TYPE {
					idci = 0
				}

				if m == 0 && (etyp == CONVRG_TYPE || etyp == DIVERG_TYPE) {
					idci = 0
				} else if m == 0 && (etyp == CVRGAIR_TYPE || etyp == DIVGAIR_TYPE) {
					idci = 0
				}

				if m == Plist.Nelm-1 && (etyp == CONVRG_TYPE || etyp == DIVERG_TYPE) {
					idco = 0
				} else if m == Plist.Nelm-1 && (etyp == CVRGAIR_TYPE || etyp == DIVGAIR_TYPE) {
					idco = 0
				}

				if idci == 1 {
					pelmci(Mpath.Fluid, Pelm, errkey)
					Pelm.In.Lpath = Plist
				}
				if idco == 1 {
					pelmco(Mpath.Fluid, Pelm, errkey)

					Pelm.Out.Lpath = Plist
					Pelm.Out.Fluid = Mpath.Fluid
				}
			}
		}

		if DEBUG {
			plistprint(1, (*M)[Mpath_idx:], Mpath.Plist[0].Pelm, Compnt[0].Elouts, Compnt[0].Elins)

			fmt.Printf("i=%d\n", i)
		}

		Plist = &Mpath.Plist[0]
		if Mpath.Nlpath == 1 {
			Pelm = Plist.Pelm[0]

			if DEBUG {
				fmt.Printf("<<Pathdata>>   Plist.type=%c\n", Plist.Type)
			}

			if Plist.Type == IN_LPTP {
				Mpath.Type = THR_PTYP

				if DEBUG {
					fmt.Printf("<<Pathdata>>   Mpath.type=%c\n", Mpath.Type)
				}
			} else {
				Mpath.Type = CIR_PTYP
				Plist.Type = CIR_PTYP
				Pelm.In.Upo = Plist.Pelm[Plist.Nelm-1].Out
			}

			if DEBUG {
				fmt.Printf("<<Pathdata>> test end\n")
			}

			for m = 1; m < Plist.Nelm; m++ {
				Pelm := Plist.Pelm[m]
				PelmPrev := Plist.Pelm[m-1]
				Pelm.In.Upo = PelmPrev.Out
			}
		} else {
			Mpath.Type = BRC_PTYP

			if DEBUG {
				fmt.Printf("<<Pathdata>> Mpath i=%d  type=%c\n", i, Mpath.Type)
			}

			for j = 0; j < Mpath.Nlpath; j++ {
				Plist = &Mpath.Plist[j]

				Pelm = Plist.Pelm[0]
				etyp = Pelm.Cmp.Eqptype

				if DEBUG {
					fmt.Printf("<<Pathdata>> Plist j=%d name=%s eqptype=%s\n", j,
						Pelm.Cmp.Name, etyp)
				}

				if etyp == DIVERG_TYPE || etyp == DIVGAIR_TYPE {
					Plist.Type = DIVERG_LPTP
				}

				if etyp == CONVRG_TYPE || etyp == CVRGAIR_TYPE {
					Plist.Type = CONVRG_LPTP
					ncv++
				}

				Pelm := Plist.Pelm[Plist.Nelm-1]
				etyp = Pelm.Cmp.Eqptype
				if etyp != DIVERG_TYPE && etyp != CONVRG_TYPE &&
					etyp != DIVGAIR_TYPE && etyp != CVRGAIR_TYPE {
					Plist.Type = OUT_LPTP
				}

				for m = 1; m < Plist.Nelm; m++ {
					Pelm = Plist.Pelm[m]
					PelmPrev := Plist.Pelm[m-1]
					Pelm.In.Upo = PelmPrev.Out
				}

				if DEBUG {
					fmt.Printf("<<Pathdata>> Plist MAX=%d  j=%d  type=%c\n", Mpath.Nlpath, j, Plist.Type)
				}
			}
		}
		Mpath.Ncv = ncv

		if DEBUG {
			fmt.Printf("2----- MAX=%d  i=%d\n", *Nmpath, i)
		}
	}
	Mpath = mpi

	if DEBUG {
		if *Nmpath > 0 {
			mpi := *M
			plistprint(*Nmpath, mpi, mpi[0].Plist[0].Pelm, Compnt[0].Elouts, Compnt[0].Elins)
		}
	}

	// バルブがモニターするPlistの検索
	Valvinit(Eqsys.Nvalv, Eqsys.Valv, *Nmpath, *M)

	// 未知流量等の構造解析
	pflowstrct(*Nmpath, (*M)[Mpath_idx:])

	if DEBUG {
		if *Nmpath > 0 {
			plistprint(*Nmpath, *M, mpi.Plist[0].Pelm, Compnt[0].Elouts, Compnt[0].Elins)
		}
	}

	if DEBUG {
		fmt.Printf("\n")
		for i = 0; i < Ncompnt; i++ {
			C := &Compnt[i]
			fmt.Printf("name=%s Nin=%d  Nout=%d\n", C.Name, C.Nin, C.Nout)
		}
	}
}

/***********************************************************************/

func Mpathcount(fi *EeTokens, Pl *int) int {
	var N int
	var ad int
	var s string

	ad = fi.GetPos()
	*Pl = 0

	for fi.IsEnd() == false {
		s = fi.GetToken()

		if s == "*" {
			break
		}

		if s == ";" {
			N++
		}

		if s == ">" {
			*Pl++
		}
	}

	*Pl /= 2

	fi.RestorePos(ad)

	return N
}

/***********************************************************************/

func Plcount(fi *EeTokens, N []int) {
	i := 0
	M := 0
	ad := fi.GetPos()

	for fi.IsEnd() == false {
		s := fi.GetToken()

		if s == "*" {
			break
		}

		if s == ";" {
			N[M] = i
			M++
			i = 0
		}

		if s == ">" {
			i++
			fi.GetToken() // skip next token
		}
	}

	// Print the contents of the N slice for debugging purposes
	// for i := 0; i < len(*N); i++ {
	// 	fmt.Printf("i=%d pl=%d\n", i, (*N)[i])
	// }

	fi.RestorePos(ad)
}

/***********************************************************************/

func Pelmcount(fi *EeTokens) int {
	ad := fi.GetPos()
	i := 1
	N := 0

	for fi.IsEnd() == false {
		s := fi.GetToken()
		i = 1

		if s == "*" {
			break
		}

		for fi.IsEnd() == false {
			s = fi.GetToken()

			if s == ";" {
				break
			}

			if s == "-f" {
				t := fi.GetToken()

				if t == "W" || t == "a" {
					i = 1
				} else {
					i = 2
				}
			}

			if s == "-sys" {
				fi.GetToken()
			}

			if s != ">" && s[:1] != "(" && s[:1] != "-" && s[:1] != ";" {
				N += i
			}
		}
	}

	fi.RestorePos(ad)
	return N
}

/***********************************************************************/

func Elcount(N int, C []COMPNT) (int, int) {
	var Nelout, Nelin int = 0, 0

	for i := 0; i < N; i++ {
		e := C[i].Eqptype
		Nelout += C[i].Nout
		Nelin += C[i].Nin

		if e == HCLOADW_TYPE {
			Nelin += 8
		} else if e == THEX_TYPE {
			Nelout += 4
			Nelin += 14
		}
	}

	Nelout *= 4
	Nelin *= 4

	return Nelout, Nelin
}

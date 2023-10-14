package eeslism

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

/*   システム要素の入力 */
func Compodata(f *EeTokens, errkey string, Rmvls *RMVLS, Eqcat *EQCAT,
	Cmp *[]COMPNT, Ncompnt *int, Eqsys *EQSYS, Ncmpalloc *int, ID int) {
	var (
		//cmp    *[]COMPNT
		Compnt []COMPNT
		Ni, No int
		cio    ELIOType
		idi    []ELIOType
		ido    []ELIOType
		N      int
	)
	D := 0

	Nroom := len(Rmvls.Room)

	var Room []ROOM
	if Nroom > 0 {
		Room = Rmvls.Room
	} else {
		Room = []ROOM{}
	}

	Nrdpnl := Rmvls.Nrdpnl

	var Rdpnl []RDPNL
	if Nrdpnl > 0 {
		Rdpnl = Rmvls.Rdpnl
	} else {
		Rdpnl = []RDPNL{}
	}

	// Nairflow := Rmvls.Nairflow
	// var AirFlow *AIRFLOW
	// if Nairflow > 0 {
	//     AirFlow = Rmvls.airflow
	// } else {
	//     AirFlow = nil
	// }

	//コンポーネント数
	Ncmp := Compntcount(f)
	Ncmp += 2 + Nroom + Nrdpnl // + Nairflow

	*Cmp = make([]COMPNT, Ncmp)
	Compinit(Ncmp, *Cmp)

	// fmt.Printf("<Compodata> Compnt Alloc=%d\n", Ncmp)
	*Ncmpalloc = Ncmp
	//cmp = Cmp
	Compnt = *Cmp

	// if fi, err := os.Open("bdata.ewk"); err != nil {
	if f == nil {
		Eprint("bdata.ewk", "<Compodata>")
		os.Exit(EXIT_BDATA)
	}

	// 給水温度設定
	Compnt[0].Name = CITYWATER_NAME
	Compnt[0].Eqptype = FLIN_TYPE
	Compnt[0].Tparm = CITYWATER_PARM
	Compnt[0].Nin = 1
	Compnt[0].Nout = 1
	Eqsys.Nflin++
	D++

	// 取り入れ外気設定
	Compnt[1].Name = OUTDRAIR_NAME
	Compnt[1].Eqptype = FLIN_TYPE
	Compnt[1].Tparm = OUTDRAIR_PARM
	Compnt[1].Nin = 2
	Compnt[1].Nout = 2
	Eqsys.Nflin++
	D++

	/* 室およびパネル用     */

	var Ncrm int
	if SIMUL_BUILDG {

		// 室およびパネル用
		for i := 0; i < Nroom; i++ {
			Compnt[i+2].Name = Room[i].Name
			Compnt[i+2].Eqptype = ROOM_TYPE
			Compnt[i+2].Neqp = i
			Compnt[i+2].Eqp = &Room[i]
			Compnt[i+2].Nout = 2
			Compnt[i+2].Nin = 2*Room[i].Nachr + Room[i].Ntr + Room[i].Nrp
			Compnt[i+2].Nivar = 0
			Compnt[i+2].Ivparm = nil
			Compnt[i+2].Airpathcpy = 'y'
		}

		Ncrm = 2 + Nroom
		for i := 0; i < Nrdpnl; i++ {
			Compnt[Ncrm+i].Name = Rdpnl[i].Name
			Compnt[Ncrm+i].Eqptype = RDPANEL_TYPE
			Compnt[Ncrm+i].Neqp = i
			Compnt[Ncrm+i].Eqp = &Rdpnl[i]
			Compnt[Ncrm+i].Nout = 2
			Compnt[Ncrm+i].Nin = 3 + Rdpnl[i].Ntrm[0] + Rdpnl[i].Ntrm[1] + Rdpnl[i].Nrp[0] + Rdpnl[i].Nrp[1] + 1
			Compnt[Ncrm+i].Nivar = 0
			Compnt[Ncrm+i].Airpathcpy = 'y'
		}

		// エアフローウィンドウの機器メモリ
		// for i := 0; i < Nairflow; i++ {
		//     Compnt[Ncrm+Nrdpnl+i].name = AirFlow[i].name
		//     Compnt[Ncrm+Nrdpnl+i].eqptype = AIRFLOW_TYPE
		//     Compnt[Ncrm+Nrdpnl+i].neqp = i
		//     Compnt[Ncrm+Nrdpnl+i].eqp = AirFlow[i]
		//     Compnt[Ncrm+Nrdpnl+i].Nout = 2
		//     Compnt[Ncrm+Nrdpnl+i].Nin = 3
		//     Compnt[Ncrm+Nrdpnl+i].nivar = 0
		//     Compnt[Ncrm+Nrdpnl+i].airpathcpy = 'y'
		// }

	}

	//cp := (*COMPNT)(nil)
	var comp_num int

	if ID == 0 {
		for f.IsEnd() == false {
			s := f.GetToken()
			if s == "*" {
				break
			}
			cio = ELIO_SPACE

			Crm := true
			comp_num := 0

			if strings.HasPrefix(s, "(") {
				if len(s) == 1 {
					s = f.GetToken()
				} else {
					fmt.Sscanf(s, "(%s", &s)
				}
				Crm = false
				Compnt[comp_num].Name = s
				comp_num++

				s = f.GetToken()

				if strings.IndexRune(s, ')') != -1 {
					idx := strings.IndexRune(s, ')')
					s = s[idx+1:]
				}
				Compnt[comp_num].Name = s
				Compnt[comp_num-1].Valvcmp = &Compnt[comp_num]

				Eqsys.Nvalv++
			}

			if Crm {
				//部屋を検索
				cp := (*COMPNT)(nil)
				for i := 0; i < Ncrm; i++ {
					if s == Compnt[i].Name {
						cp = &Compnt[i]
						break
					}
				}
				if cp == nil {
					Compnt[comp_num].Name = s
					Compnt[comp_num].Ivparm = nil
					Compnt[comp_num].Tparm = ""
					Compnt[comp_num].Envname = ""
					Compnt[comp_num].Roomname = ""
					Compnt[comp_num].Nivar = 0
				}
			}

			for f.IsEnd() == false {
				s = f.GetToken()
				if s[0] == ';' {
					break
				}
				if strings.HasPrefix(s, "-") {
					/********************************/
					if cio == 'i' {
						Compnt[comp_num].Nin = Ni
						idi[Ni] = ELIO_None
						Compnt[comp_num].Idi = idi
					} else if cio == 'o' {
						Compnt[comp_num].Nout = No
						ido[No] = ELIO_None
						Compnt[comp_num].Ido = ido
						if Compnt[comp_num].Eqptype == DIVERG_TYPE {
							Compnt[comp_num].Nivar = No
							Compnt[comp_num].Ivparm = new(float64)
						}
					}

					/********************************/

					// ハイフンの後ろに続く文字列を取得
					ps := s[1:]

					switch ps {
					case "c":
						cio = 'c'
					case "type":
						cio = 't'
					case "Nin":
						cio = 'I'
					case "Nout":
						cio = 'O'
					case "in":
						cio = 'i'
						Ni = 0
						idi = []ELIOType{}
					case "out":
						cio = 'o'
						No = 0
						ido = []ELIOType{}
					case "L":
						cio = 'L'
						Compnt[comp_num].Nivar = 1
						Compnt[comp_num].Ivparm = new(float64)
					case "env":
						cio = 'e'
					case "room":
						cio = 'r'
					case "roomheff":
						cio = 'R'
					case "exs":
						cio = 's'
					case "S":
						cio = 'S'
					case "Tinit":
						cio = 'T'
					case "V":
						cio = 'V'
					case "hcc":
						cio = 'h'
					case "pfloor":
						cio = 'f'
					case "wet":
						Compnt[comp_num].Wetparm = ps

						/*---- Roh Debug for a constant outlet humidity model of wet coil  2003/4/25 ----*/
						cio = 'w'
						Compnt[comp_num].Ivparm = new(float64)
						*Compnt[comp_num].Ivparm = 90.0
					case "control":
						cio = 'M'
					case "monitor":
						cio = 'm'
					case "PCMweight":
						cio = 'P'
					default:
						Eprint(errkey, s)
					}
				} else if cio != 'V' && cio != 'S' && strings.IndexRune(s, '-') != -1 {
					idx := strings.IndexRune(s, '-')
					st := s[idx+1:]
					var err error
					switch {
					case strings.HasPrefix(s, "Ac"):
						Compnt[comp_num].Ac, err = strconv.ParseFloat(st, 64)
						if err != nil {
							panic(err)
						}
					case strings.HasPrefix(s, "PVcap"):
						Compnt[comp_num].PVcap, err = strconv.ParseFloat(st, 64)
						if err != nil {
							panic(err)
						}
					case strings.HasPrefix(s, "Area"):
						Compnt[comp_num].Area, err = strconv.ParseFloat(st, 64)
						if err != nil {
							panic(err)
						}
					}
				} else {
					switch cio {
					case 'c':
						if eqpcat(s, &Compnt[comp_num], Eqcat, Eqsys) {
							Eprint(errkey, s)
						}
						break
					case 't':
						Compnt[comp_num].Eqptype = EqpType(s)

						if Compnt[comp_num].Valvcmp != nil {
							Compnt[comp_num].Valvcmp.Eqptype = EqpType(s)
						}

						Compnt[comp_num].Neqp = 0
						Compnt[comp_num].Ncat = 0

						switch s {
						case CONVRG_TYPE, CVRGAIR_TYPE:
							Eqsys.Ncnvrg++
							Compnt[comp_num].Nout = 1
						case DIVERG_TYPE, DIVGAIR_TYPE:
							Compnt[comp_num].Nin = 1
						case FLIN_TYPE:
							Eqsys.Nflin++
						case HCLOAD_TYPE, HCLOADW_TYPE, RMAC_TYPE, RMACD_TYPE:
							Eqsys.Nhcload++
						case VALV_TYPE, TVALV_TYPE:
							Eqsys.Nvalv++
						case QMEAS_TYPE:
							Eqsys.Nqmeas++
						default:
							if s != DIVERG_TYPE && s != DIVGAIR_TYPE {
								Eprint(errkey, s)
							}
						}

						break

					case 'i':
						idi = append(idi, ELIO_None)
						Ni++
						break
					case 'o':
						ido = append(ido, ELIO_None)
						No++
						break
					case 'I':
						// Satoh DEBUG 1998/5/15
						if Crm {
							var err error
							if Compnt[comp_num].Eqptype == ROOM_TYPE {
								if SIMUL_BUILDG {
									room := Compnt[comp_num].Eqp.(*ROOM)
									room.Nasup, err = strconv.Atoi(s)
									if err != nil {
										panic(err)
									}

									N := room.Nasup
									if N > 0 {
										if room.Arsp == nil {
											room.Arsp = make([]AIRSUP, N)
										}
									}

									Compnt[comp_num].Nin += 2 * room.Nasup
								}
							} else {
								Compnt[comp_num].Nin, err = strconv.Atoi(s)
								if err != nil {
									panic(err)
								}
							}

							for i := 0; i < Compnt[comp_num].Nin; i++ {
								idi[i] = ' '
							}
							Compnt[comp_num].Idi = idi[:Compnt[comp_num].Nin]
						}
						break
					case 'O':
						var err error
						Compnt[comp_num].Nout, err = strconv.Atoi(s)
						if err != nil {
							panic(err)
						}
						ido = make([]ELIOType, Compnt[comp_num].Nout)
						for i := 0; i < Compnt[comp_num].Nout; i++ {
							ido[i] = ELIO_SPACE
						}
						Compnt[comp_num].Ido = ido
						break

					case 'L':
						var err error
						*Compnt[comp_num].Ivparm, err = strconv.ParseFloat(s, 64)
						if err != nil {
							panic(err)
						}
						break
					case 'e':
						Compnt[comp_num].Envname = s
						break
					case 'r':
						Compnt[comp_num].Roomname = s
						break
					case 'h':
						Compnt[comp_num].Hccname = s
						break
					case 'f':
						Compnt[comp_num].Rdpnlname = s
						break
					case 'R':
						var err error
						Compnt[comp_num].Roomname = s
						s = f.GetToken()
						Compnt[comp_num].Eqpeff, err = strconv.ParseFloat(s, 64)
						if err != nil {
							panic(err)
						}
						break
					case 's':
						Compnt[comp_num].Exsname = s
						break
					case 'M':
						Compnt[comp_num].Omparm = s
						break
					case 'S', 'V':
						s += "  "
						s += strings.Repeat(" ", len(s))
						_s := f.GetToken()
						s += _s + " *"
						Compnt[comp_num].Tparm = s
						break

					case 'T':
						if strings.HasPrefix(s, "(") {
							s += " "
							s += strings.Repeat(" ", len(s))
							_s := f.GetToken()
							s += _s
							Compnt[comp_num].Tparm = s
						} else {
							Compnt[comp_num].Tparm = s
						}
						break
					case 'w':
						var err error
						*Compnt[comp_num].Ivparm, err = strconv.ParseFloat(s, 64)
						if err != nil {
							panic(err)
						}
						break
					case 'm':
						Compnt[comp_num].MonPlistName = s
						Compnt[comp_num].Valvcmp.MonPlistName = s
						break
					case 'P':
						var err error
						Compnt[comp_num].MPCM, err = strconv.ParseFloat(s, 64)
						if err != nil {
							panic(err)
						}
						break
					}
				}
			}

			if Crm == false {
				comp_num++
				D++
			}
		}
	}

	ncmp := comp_num
	for i := 0; i < ncmp; i++ {
		cm := &Compnt[i]
		if cm.Eqptype == DIVGAIR_TYPE {
			s := cm.Name + ".x"
			Compnt[comp_num].Name = s
			Compnt[comp_num].Eqptype = cm.Eqptype
			Compnt[comp_num].Nout = cm.Nout
			Compnt[comp_num].Ido = cm.Ido
			comp_num++
		} else if cm.Eqptype == CVRGAIR_TYPE {
			s := cm.Name + ".x"
			Compnt[comp_num].Name = s
			Compnt[comp_num].Eqptype = cm.Eqptype
			Compnt[comp_num].Nin = cm.Nin
			Compnt[comp_num].Idi = cm.Idi
			Eqsys.Ncnvrg++
			comp_num++
		}
	}
	*Ncompnt = Ncmp

	//fmt.Printf("<<Compodata>> Ncompnt = %d\n", *Ncompnt)

	N = Eqsys.Nvalv
	if N > 0 {
		Eqsys.Valv = make([]VALV, N)
		Valv := Eqsys.Valv
		for i := 0; i < N; i++ {
			Valv[i].Name = ""
			Valv[i].Cmp = nil
			Valv[i].Cmb = nil
			Valv[i].Org = 'n'
			Valv[i].X = -999.0
			Valv[i].Xinit = nil
			Valv[i].Tin = nil
			Valv[i].Tset = nil
			Valv[i].Mon = nil
			Valv[i].Plist = nil
			Valv[i].MGo = nil
			Valv[i].Tout = nil
			Valv[i].MonPlist = nil
			//Valv[i].OMfan = nil
			//Valv[i].OMfanName = nil
		}
	}

	N = Eqsys.Nqmeas
	if N > 0 {
		Eqsys.Qmeas = make([]QMEAS, N)
		for i := 0; i < N; i++ {
			q := &Eqsys.Qmeas[i]
			q.Name = ""
			q.Cmp = nil
			q.Th = nil
			q.Tc = nil
			q.G = nil
			q.PlistG = nil
			q.PlistTc = nil
			q.PlistTh = nil
			q.Plistxc = nil
			q.Plistxh = nil
			q.Xc = nil
			q.Xh = nil
			q.Id = 0
			q.Nelmc = -999
			q.Nelmh = -999
		}
	}

	N = Eqsys.Nhcc
	Eqsys.Hcc = nil
	if N > 0 {
		Eqsys.Hcc = make([]HCC, N)
		for i := 0; i < N; i++ {
			hcc := &Eqsys.Hcc[i]
			hcc.Name = ""
			hcc.Cmp = nil
			hcc.Cat = nil
			hcc.Twin = 5.0
			hcc.Xain = FNXtr(25.0, 50.0)
		}
	}

	N = Eqsys.Nboi
	Eqsys.Boi = nil
	if N > 0 {
		Eqsys.Boi = make([]BOI, N)

		for i := 0; i < N; i++ {
			Boi := &Eqsys.Boi[i]
			Boi.Name = ""
			Boi.Cmp = nil
			Boi.Cat = nil
			Boi.Load = nil
			Boi.Mode = 'M'
			MtEdayinit(&Boi.MtEdy)
			MtEdayinit(&Boi.MtPhdy)
		}
	}

	N = Eqsys.Ncoll
	Eqsys.Coll = nil
	if N > 0 {
		Eqsys.Coll = make([]COLL, N)
		for i := 0; i < N; i++ {
			Coll := &Eqsys.Coll[i]
			Coll.Name = ""
			Coll.Cmp = nil
			Coll.Cat = nil
			Coll.sol = nil
			Coll.Te = 0.0
			Coll.Tcb = 0.0
			//Coll.Fd = 0.9
		}
	}

	N = Eqsys.Npv
	Eqsys.PVcmp = nil
	if N > 0 {
		Eqsys.PVcmp = make([]PV, N)
		for i := 0; i < N; i++ {
			PV := &Eqsys.PVcmp[i]
			PV.PVcap = -999.
			PV.Area = -999.
			PV.Name = ""
			PV.Cmp = nil
			PV.Cat = nil
			PV.Sol = nil
			PV.Ta = nil
			PV.V = nil
			PV.I = nil
			MtEdayinit(&PV.mtEdy)
		}
	}

	// Satoh OMVAV  2010/12/16
	N = Eqsys.Nomvav
	Eqsys.OMvav = nil

	if N > 0 {
		Eqsys.OMvav = make([]OMVAV, N)
		for i := 0; i < N; i++ {
			OMvav := &Eqsys.OMvav[i]
			OMvav.Name = ""
			OMvav.Cmp = nil
			OMvav.Cat = nil
			OMvav.Omwall = nil
			OMvav.Plist = nil
			OMvav.Nrdpnl = 0

			for j := 0; j < 4; j++ {
				OMvav.Rdpnl[j] = nil
			}
		}
	}

	N = Eqsys.Nrefa
	Eqsys.Refa = nil
	if N > 0 {
		Eqsys.Refa = make([]REFA, N)
		for i := 0; i < N; i++ {
			Rf := &Eqsys.Refa[i]
			Rf.Name = ""
			Rf.Cmp = nil
			Rf.Cat = nil
			Rf.Load = nil
			Rf.Room = nil
			Rf.Do = 0.0
			Rf.D1 = 0.0
			Rf.Tin = 0.0
			Rf.Te = 0.0
			Rf.Tc = 0.0
			MtEdayinit(&Rf.mtEdy)
			MtEdayinit(&Rf.mtPhdy)
		}
	}

	N = Eqsys.Npipe
	Eqsys.Pipe = nil
	if N > 0 {
		Eqsys.Pipe = make([]PIPE, N)
		for i := 0; i < N; i++ {
			Pi := &Eqsys.Pipe[i]
			Pi.Name = ""
			Pi.Cmp = nil
			Pi.Cat = nil
			Pi.Loadt = nil
			Pi.Loadx = nil
			Pi.Room = nil
		}
	}

	N = Eqsys.Nstank
	Eqsys.Stank = nil
	if N > 0 {
		Eqsys.Stank = make([]STANK, N)
		for i := 0; i < N; i++ {
			St := &Eqsys.Stank[i]
			St.Name = ""
			St.Cmp = nil
			St.Cat = nil
			St.Jin = nil
			St.Jout = nil
			St.Batchcon = nil
			St.Ihex = nil
			St.Batchcon = nil
			St.B = nil
			St.R = nil
			St.Fg = nil
			St.D = nil
			St.Tss = nil
			St.DtankF = nil
			St.Tssold = nil
			St.Dvol = nil
			St.Mdt = nil
			St.KS = nil
			St.Ihxeff = nil
			St.KA = nil
			St.KAinput = nil
			St.CGwin = nil
			St.EGwin = nil
			St.Twin = nil
			St.Q = nil
			St.Tenv = nil
			St.Stkdy = nil
			St.Mstkdy = nil
			St.MQlossdy = 0.0
			St.MQstody = 0.0
			St.Ncalcihex = 0
		}
	}

	N = Eqsys.Nhex
	Eqsys.Hex = nil
	if N > 0 {
		Eqsys.Hex = make([]HEX, N)
		for i := 0; i < N; i++ {
			Hx := &Eqsys.Hex[i]
			Hx.Name = ""
			Hx.Cmp = nil
			Hx.Cat = nil
			Hx.Id = 0
		}
	}

	N = Eqsys.Npump
	Eqsys.Pump = nil
	if N > 0 {
		Eqsys.Pump = make([]PUMP, N)

		for i := 0; i < N; i++ {
			Pp := &Eqsys.Pump[i]
			Pp.Name = ""
			Pp.Cmp = nil
			Pp.Cat = nil
			Pp.Sol = nil
			MtEdayinit(&Pp.MtEdy)
		}
	}

	N = Eqsys.Ncnvrg
	if N > 0 {
		Eqsys.Cnvrg = make([]*COMPNT, N)
		for i := 0; i < N; i++ {
			Eqsys.Cnvrg[i] = nil
		}
	}

	N = Eqsys.Nflin
	Eqsys.Flin = nil
	if N > 0 {
		Eqsys.Flin = make([]FLIN, N)
		for i := 0; i < N; i++ {
			Fl := &Eqsys.Flin[i]
			Fl.Name = ""
			Fl.Namet = ""
			Fl.Namex = ""
			Fl.Vart = nil
			Fl.Varx = nil
			Fl.Cmp = nil
		}
	}

	N = Eqsys.Nhcload
	Eqsys.Hcload = nil
	if N > 0 {
		Eqsys.Hcload = make([]HCLOAD, N)
		for i := 0; i < N; i++ {
			Hl := &Eqsys.Hcload[i]
			Hl.Name = ""
			Hl.Loadt = nil
			Hl.Loadx = nil
			Hl.Cmp = nil
			Hl.RMACFlg = 'N'
			MtEdayinit(&Hl.mtEdy)
			Hl.Ga = 0.0
			Hl.Gw = 0.0
			Hl.RHout = 50.0
		}
	}

	/*---- Satoh Debug VAV  2000/10/30 ----*/
	N = Eqsys.Nvav
	Eqsys.Vav = nil
	if N > 0 {
		Eqsys.Vav = make([]VAV, N)
		for i := 0; i < N; i++ {
			V := &Eqsys.Vav[0]
			V.Name = ""
			V.Cat = nil
			V.Hcc = nil
			V.Hcld = nil
			V.Cmp = nil
		}
	}

	N = Eqsys.Nstheat
	Eqsys.Stheat = nil
	if N > 0 {
		Eqsys.Stheat = make([]STHEAT, N)
		for i := 0; i < N; i++ {
			Sh := &Eqsys.Stheat[i]
			Sh.Name = ""
			Sh.Cat = nil
			Sh.Cmp = nil
			Sh.Room = nil
			Sh.Pcm = nil
			MtEdayinit(&Sh.MtEdy)
		}
	}

	// Satoh追加　デシカント槽 2013/10/23
	N = Eqsys.Ndesi
	Eqsys.Desi = nil
	if N > 0 {
		Eqsys.Desi = make([]DESI, N)
		for i := 0; i < N; i++ {
			Desi := &Eqsys.Desi[i]
			Desi.Name = ""
			Desi.Cat = nil
			Desi.Cmp = nil
			Desi.Room = nil
			Desi.Tenv = nil
		}
	}

	// Satoh追加　気化冷却器 2013/10/26
	N = Eqsys.Nevac
	Eqsys.Evac = nil
	if N > 0 {
		Eqsys.Evac = make([]EVAC, N)
		for i := 0; i < N; i++ {
			Evac := &Eqsys.Evac[i]
			Evac.Name = ""
			Evac.Cat = nil
			Evac.Cmp = nil
			Evac.M = nil
			Evac.Kx = nil
			Evac.Tdry = nil
			Evac.Twet = nil
			Evac.Xdry = nil
			Evac.Xwet = nil
			Evac.Ts = nil
			Evac.Xs = nil
			Evac.RHdry = nil
			Evac.RHwet = nil
			Evac.UXC = nil
			Evac.UX = nil
			//Evac.UXdry = nil
			//Evac.UXwet = nil
		}
	}

	N = Eqsys.Nthex
	Eqsys.Thex = nil
	if N > 0 {
		Eqsys.Thex = make([]THEX, N)
		for i := 0; i < N; i++ {
			thex := &Eqsys.Thex[i]
			thex.Name = ""
			thex.Cmp = nil
			thex.Cat = nil
			thex.Type = ' '
			thex.CGe = 0.0
			thex.Ge = 0.0
			thex.Go = 0.0
			thex.ET = -999.0
			thex.EH = -999.0
		}
	}
	//printf("<<Compodata>> end\n");
}

func Compntcount(fi *EeTokens) int {
	N := 0
	ad := fi.GetPos()

	for fi.IsEnd() == false {
		s := fi.GetToken()

		if s == ";" {
			N++

			s = fi.GetToken()

			if s == "*" {
				break
			} else if s == "(" {
				N++
			}
		} else if s == DIVGAIR_TYPE || s == CVRGAIR_TYPE {
			N++
		} else if s == "(" {
			N++
		}
	}

	fi.RestorePos(ad)

	return N
}

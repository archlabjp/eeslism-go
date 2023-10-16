package eeslism

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

/*   システム要素の入力 */
func Compodata(f *EeTokens, errkey string, Rmvls *RMVLS, Eqcat *EQCAT,
	Cmp *[]COMPNT, Eqsys *EQSYS, Ncmpalloc *int, ID int) {
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

	Nrdpnl := len(Rmvls.Rdpnl)

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
							// 合流要素
							Eqsys.Ncnvrg++
							Compnt[comp_num].Nout = 1
						case DIVERG_TYPE, DIVGAIR_TYPE:
							// 分岐要素
							Compnt[comp_num].Nin = 1
						case FLIN_TYPE:
							// 流入境界条件
							Eqsys.Nflin++
						case HCLOAD_TYPE, HCLOADW_TYPE, RMAC_TYPE, RMACD_TYPE:
							// 空調負荷
							Eqsys.Nhcload++
						case VALV_TYPE, TVALV_TYPE:
							// 弁・ダンパー
							Eqsys.Nvalv++
						case QMEAS_TYPE:
							// カロリーメータ
							//Eqsys.Nqmeas++
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

	//fmt.Printf("<<Compodata>> Ncompnt = %d\n", *Ncompnt)

	N = Eqsys.Nvalv
	if N > 0 {
		Eqsys.Valv = make([]VALV, N)
		Valv := Eqsys.Valv
		for i := 0; i < N; i++ {
			Valv[i] = NewVALV()
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
			Eqsys.Flin[i] = NewFLIN()
		}
	}

	N = Eqsys.Nhcload
	Eqsys.Hcload = nil
	if N > 0 {
		Eqsys.Hcload = make([]HCLOAD, N)
		for i := 0; i < N; i++ {
			Eqsys.Hcload[i] = NewHCLOAD()
		}
	}

	//printf("<<Compodata>> end\n");
}

func NewVALV() VALV {
	return VALV{
		Name:     "",
		Cmp:      nil,
		Cmb:      nil,
		Org:      'n',
		X:        -999.0,
		Xinit:    nil,
		Tin:      nil,
		Tset:     nil,
		Mon:      nil,
		Plist:    nil,
		MGo:      nil,
		Tout:     nil,
		MonPlist: nil,
		//OMfan : nil,
		//OMfanName : nil,
	}
}

func NewHCC() HCC {
	return HCC{
		Name: "",
		Cmp:  nil,
		Cat:  nil,
		Twin: 5.0,
		Xain: FNXtr(25.0, 50.0),
	}
}

func NewBOI() BOI {
	Boi := BOI{
		Name: "",
		Cmp:  nil,
		Cat:  nil,
		Load: nil,
		Mode: 'M',
	}
	MtEdayinit(&Boi.MtEdy)
	MtEdayinit(&Boi.MtPhdy)
	return Boi
}

func NewCOLL() COLL {
	return COLL{
		Name: "",
		Cmp:  nil,
		Cat:  nil,
		sol:  nil,
		Te:   0.0,
		Tcb:  0.0,
	}
}

func NewPV() PV {
	PV := PV{
		PVcap: -999.,
		Area:  -999.,
		Name:  "",
		Cmp:   nil,
		Cat:   nil,
		Sol:   nil,
		Ta:    nil,
		V:     nil,
		I:     nil,
	}
	MtEdayinit(&PV.mtEdy)
	return PV
}

func NewOMVAV() OMVAV {
	OMvav := OMVAV{
		Name:   "",
		Cmp:    nil,
		Cat:    nil,
		Omwall: nil,
		Plist:  nil,
		Nrdpnl: 0,
	}
	for j := 0; j < 4; j++ {
		OMvav.Rdpnl[j] = nil
	}
	return OMvav
}

func NewREFA() REFA {
	Rf := REFA{
		Name: "",
		Cmp:  nil,
		Cat:  nil,
		Load: nil,
		Room: nil,
		Do:   0.0,
		D1:   0.0,
		Tin:  0.0,
		Te:   0.0,
		Tc:   0.0,
	}
	MtEdayinit(&Rf.mtEdy)
	MtEdayinit(&Rf.mtPhdy)
	return Rf
}

func NewPIPE() PIPE {
	return PIPE{
		Name:  "",
		Cmp:   nil,
		Cat:   nil,
		Loadt: nil,
		Loadx: nil,
		Room:  nil,
	}
}

func NewSTANK() STANK {
	return STANK{
		Name:      "",
		Cmp:       nil,
		Cat:       nil,
		Jin:       nil,
		Jout:      nil,
		Batchcon:  nil,
		Ihex:      nil,
		B:         nil,
		R:         nil,
		Fg:        nil,
		D:         nil,
		Tss:       nil,
		DtankF:    nil,
		Tssold:    nil,
		Dvol:      nil,
		Mdt:       nil,
		KS:        nil,
		Ihxeff:    nil,
		KA:        nil,
		KAinput:   nil,
		CGwin:     nil,
		EGwin:     nil,
		Twin:      nil,
		Q:         nil,
		Tenv:      nil,
		Stkdy:     nil,
		Mstkdy:    nil,
		MQlossdy:  0.0,
		MQstody:   0.0,
		Ncalcihex: 0,
	}
}

func NewHEX() HEX {
	return HEX{
		Name: "",
		Cmp:  nil,
		Cat:  nil,
		Id:   0,
	}
}

func NewPUMP() PUMP {
	Pp := PUMP{
		Name: "",
		Cmp:  nil,
		Cat:  nil,
		Sol:  nil,
	}
	MtEdayinit(&Pp.MtEdy)
	return Pp
}

func NewFLIN() FLIN {
	return FLIN{
		Name:  "",
		Namet: "",
		Namex: "",
		Vart:  nil,
		Varx:  nil,
		Cmp:   nil,
	}
}

func NewHCLOAD() HCLOAD {
	Hl := HCLOAD{
		Name:    "",
		Loadt:   nil,
		Loadx:   nil,
		Cmp:     nil,
		RMACFlg: 'N',
	}
	MtEdayinit(&Hl.mtEdy)
	Hl.Ga = 0.0
	Hl.Gw = 0.0
	Hl.RHout = 50.0
	return Hl
}

func NewVAV() VAV {
	return VAV{
		Name: "",
		Cat:  nil,
		Hcc:  nil,
		Hcld: nil,
		Cmp:  nil,
	}
}

func NewSTHEAT() STHEAT {
	st := STHEAT{
		Name: "",
		Cat:  nil,
		Cmp:  nil,
		Room: nil,
		Pcm:  nil,
	}
	MtEdayinit(&st.MtEdy)
	return st
}

func NewDESI() DESI {
	return DESI{
		Name: "",
		Cat:  nil,
		Cmp:  nil,
		Room: nil,
		Tenv: nil,
	}
}

func NewEVAC() EVAC {
	return EVAC{
		Name:  "",
		Cat:   nil,
		Cmp:   nil,
		M:     nil,
		Kx:    nil,
		Tdry:  nil,
		Twet:  nil,
		Xdry:  nil,
		Xwet:  nil,
		Ts:    nil,
		Xs:    nil,
		RHdry: nil,
		RHwet: nil,
		UXC:   nil,
		UX:    nil,
		//UXdry: nil,
		//UXwet: nil,
	}
}

func NewTHEX() THEX {
	return THEX{
		Name: "",
		Cmp:  nil,
		Cat:  nil,
		Type: ' ',
		CGe:  0.0,
		Ge:   0.0,
		Go:   0.0,
		ET:   -999.0,
		EH:   -999.0,
	}
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

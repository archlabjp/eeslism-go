package eeslism

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// システム要素の入力
//
// - SYSCMPを`f`から読み込み、`Cmp`に使用機器情報を登録する.同時に、 `Eqsys`にメモリを確保する。
// - `Cmp`に使用機器情報を登録する際に `Eqcat`のカタログデータを参照する。
// - `Rmvls`に含まれる室および放射パネルは `Cmp`に使用機器として自動登録される。
//
// TODO:
// - この関数内で Eqsysにメモリを確保すると関数の責務を超えるので、分離独立させたほうが良い。
//
func Compodata(f *EeTokens, Rmvls *RMVLS, Eqcat *EQCAT, Cmp *[]*COMPNT, Eqsys *EQSYS) {
	var (
		Ni, No int
		cio    ELIOType
		idi    []ELIOType
		ido    []ELIOType
	)
	D := 0
	errkey := "Compodata"

	Room := Rmvls.Room
	Rdpnl := Rmvls.Rdpnl

	// Nairflow := Rmvls.Nairflow
	// var AirFlow *AIRFLOW
	// if Nairflow > 0 {
	//     AirFlow = Rmvls.airflow
	// } else {
	//     AirFlow = nil
	// }

	//コンポーネント数
	*Cmp = make([]*COMPNT, 0, 100)

	// fmt.Printf("<Compodata> Compnt Alloc=%d\n", Ncmp)
	//cmp = Cmp

	// if fi, err := os.Open("bdata.ewk"); err != nil {
	if f == nil {
		Eprint("bdata.ewk", "<Compodata>")
		os.Exit(EXIT_BDATA)
	}

	// ----------------------------------------------------------
	// 組込みシステム要素
	// 給水温度 _CW: -type FLI -V t=Twsup * ;
	// 外気温度・絶対湿度 _OA: -type FLI -V t=_Ta x=_xa * ;
	// 室 <室名>: ROOMデータセットの宣言に応じて自動で組み込まれる
	// 放射パネル <パネル名>:  放射パネルの
	// ----------------------------------------------------------

	// 給水温度設定 `_CW`
	Cmp1 := NewCOMPNT()
	Cmp1.Name = CITYWATER_NAME
	Cmp1.Eqptype = FLIN_TYPE
	Cmp1.Tparm = CITYWATER_PARM
	Cmp1.Nin = 1
	Cmp1.Nout = 1
	Eqsys.Flin = append(Eqsys.Flin, NewFLIN())
	D++

	// 取り入れ外気設定 `_OA`
	Cmp2 := NewCOMPNT()
	Cmp2.Name = OUTDRAIR_NAME
	Cmp2.Eqptype = FLIN_TYPE
	Cmp2.Tparm = OUTDRAIR_PARM
	Cmp2.Nin = 2
	Cmp2.Nout = 2
	Eqsys.Flin = append(Eqsys.Flin, NewFLIN())
	D++

	*Cmp = append(*Cmp, Cmp1, Cmp2)

	/* 室およびパネル用     */

	var Ncrm int
	if SIMUL_BUILDG {

		// 室用 `<室名>`
		for i := range Rmvls.Room {
			c := NewCOMPNT()
			c.Name = Room[i].Name
			c.Eqptype = ROOM_TYPE
			c.Neqp = i
			c.Eqp = Room[i]
			c.Nout = 2
			c.Nin = 2*Room[i].Nachr + Room[i].Ntr + Room[i].Nrp
			c.Nivar = 0
			c.Ivparm = nil
			c.Airpathcpy = true
			*Cmp = append(*Cmp, c)
			D++
		}
		Ncrm = 2 + len(Rmvls.Room) // 給水温度設定+取り入れ外気設定+室の数

		// パネル用 `<パネル名>`
		for i := range Rdpnl {
			c := NewCOMPNT()
			c.Name = Rdpnl[i].Name
			c.Eqptype = RDPANEL_TYPE
			c.Neqp = i
			c.Eqp = Rdpnl[i]
			c.Nout = 2
			c.Nin = 3 + Rdpnl[i].Ntrm[0] + Rdpnl[i].Ntrm[1] + Rdpnl[i].Nrp[0] + Rdpnl[i].Nrp[1] + 1
			c.Nivar = 0
			c.Airpathcpy = true
			*Cmp = append(*Cmp, c)
			D++
		}

		// エアフローウィンドウの機器メモリ
		// for i := 0; i < Nairflow; i++ {
		//	   c := NewCOMPNT()
		//     c.name = AirFlow[i].name
		//     c.eqptype = AIRFLOW_TYPE
		//     c.neqp = i
		//     c.eqp = AirFlow[i]
		//     c.Nout = 2
		//     c.Nin = 3
		//     c.nivar = 0
		//     c.airpathcpy = true
		//	   *Cmp = append(*Cmp, *c)
		// }

	}

	//cp := (*COMPNT)(nil)

	for f.IsEnd() == false {
		s := f.GetToken()
		if s == "*" {
			break
		}
		if s == "\n" {
			continue
		}
		cio = ELIO_SPACE

		var comp *COMPNT = NewCOMPNT()
		var Crm *COMPNT = comp

		// 三方弁（二方弁で連動する弁）を指定するとき
		// `(<elmname1> <elmname2>)` のように入力する
		// elmname1とelmname2が三方弁のように、逆作動の連動弁として機能するときの記述方法。
		// このように書くとelmname1とelmname2はともに2方弁であるが、elmname1が開くとelmname2は閉じる操作が行われる。
		//
		//   [elm1] --> [elm2]
		//
		if strings.HasPrefix(s, "(") {
			if len(s) == 1 {
				s = f.GetToken()
			} else {
				fmt.Sscanf(s, "(%s", &s)
			}
			Crm = nil

			// 1つ目の2方弁
			vc1 := NewCOMPNT()
			vc1.Name = s // <compname1>

			s = f.GetToken()
			if strings.IndexRune(s, ')') != -1 {
				idx := strings.IndexRune(s, ')')
				s = s[idx+1:]
			}

			// 2つ目の2方弁
			vc2 := NewCOMPNT()
			vc2.Name = s // <compname2>

			// 1つ目の2方弁から2つ目の2方弁を参照させる
			vc1.Valvcmp = vc2

			*Cmp = append(*Cmp, vc1, vc2)

			// 2つ目の要素に対して各種設定を反映させる
			comp = vc2

			// NOTE: ここでVALVを追加する理由が不明
			Eqsys.Valv = append(Eqsys.Valv, NewVALV())
		}

		if Crm != nil {
			// 組み込みの部屋要素を探す → 該当がなければ初期化 (必要性不明)
			Crm = FindCOMPNTByName(s, (*Cmp)[:Ncrm])
			if Crm == nil {
				comp.Name = s
				comp.Ivparm = nil
				comp.Tparm = ""
				comp.Envname = ""
				comp.Roomname = ""
				comp.Nivar = 0
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
					comp.Nin = Ni
					comp.Idi = make([]ELIOType, Ni)
				} else if cio == 'o' {
					comp.Nout = No
					comp.Ido = make([]ELIOType, No)
					if comp.Eqptype == DIVERG_TYPE {
						comp.Nivar = No
						//comp.Ivparm = new(float64)
					}
				}

				/********************************/

				// ハイフンの後ろに続く文字列を取得
				ps := s[1:]

				switch ps {
				case "c":
					// カタログ名
					cio = 'c'
				case "type":
					// 要素の種類
					cio = 't'
				case "Nin":
					// 合流数
					cio = 'I'
				case "Nout":
					// 分岐数
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
					// 配管長 [m]
					cio = 'L'
					comp.Nivar = 1
				case "env":
					// 周囲温度
					cio = 'e'
				case "room":
					// 機器設置空間の室名
					cio = 'r'
				case "roomheff":
					cio = 'R'
				case "exs":
					// 太陽熱集熱器の方位・傾斜名（EXSRFデータで入力）
					cio = 's'
				case "S":
					// ???
					cio = 'S'
				case "Tinit":
					// 蓄熱槽の初期水温
					cio = 'T'
				case "V":
					// ???
					cio = 'V'
				case "hcc":
					// VWV制御するときの制御対象熱交換器名称
					cio = 'h'
				case "pfloor":
					// VWV制御するときの制御対象床暖房
					cio = 'f'
				case "wet":
					// 冷却コイルで出口相対湿度一定の仮定で除湿の計算を行う。
					comp.Wetparm = ps

					/*---- Roh Debug for a constant outlet humidity model of wet coil  2003/4/25 ----*/
					cio = 'w'
					comp.Ivparm = CreateConstantValuePointer(90.0)
				case "control":
					// 集熱器が直列接続の場合に流れ方向に記載する
					cio = 'M'
				case "monitor":
					// Satoh Debug
					cio = 'm'
				case "PCMweight":
					// 電気蓄熱暖房器の潜熱蓄熱材重量（kg）
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
					// 集熱器面積 [m2]
					comp.Ac, err = readFloat(st)
					if err != nil {
						panic(err)
					}
				case strings.HasPrefix(s, "PVcap"):
					// 設置容量 [W] (PV)
					comp.PVcap, err = readFloat(st)
					if err != nil {
						panic(err)
					}
				case strings.HasPrefix(s, "Area"):
					// アレイ面積 [m2](PV)
					comp.Area, err = readFloat(st)
					if err != nil {
						panic(err)
					}
				}
			} else {
				switch cio {
				case 'c':
					// `-c <カタログ名>`
					if eqpcat(s, comp, Eqcat, Eqsys) {
						Eprint(errkey, s)
					}
					break
				case 't':
					// `-type <種類>`
					comp.Eqptype = EqpType(s)

					//  三方弁（二方弁で連動する弁）
					if comp.Valvcmp != nil {
						// 連動する対となる要素にも種類を指定
						comp.Valvcmp.Eqptype = EqpType(s)
					}

					comp.Neqp = 0
					comp.Ncat = 0

					switch s {
					case CONVRG_TYPE, CVRGAIR_TYPE:
						// 合流要素
						Eqsys.Cnvrg = append(Eqsys.Cnvrg, nil)
						comp.Nout = 1
					case DIVERG_TYPE, DIVGAIR_TYPE:
						// 分岐要素
						comp.Nin = 1
					case FLIN_TYPE:
						// 流入境界条件
						Eqsys.Flin = append(Eqsys.Flin, NewFLIN())
					case HCLOAD_TYPE, HCLOADW_TYPE, RMAC_TYPE, RMACD_TYPE:
						// 空調負荷
						Eqsys.Hcload = append(Eqsys.Hcload, NewHCLOAD())

					case VALV_TYPE, TVALV_TYPE:
						// 弁・ダンパー
						Eqsys.Valv = append(Eqsys.Valv, NewVALV())
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
					// `-in`
					idi = append(idi, ELIO_None)
					Ni++
					break
				case 'o':
					// `-out`
					ido = append(ido, ELIO_None)
					No++
					break
				case 'I':
					// `-Nin <合流数>`
					// Satoh DEBUG 1998/5/15
					var err error
					if Crm != nil && SIMUL_BUILDG {
						if Crm.Eqptype == ROOM_TYPE {
							room := Crm.Eqp.(*ROOM)
							room.Nasup, err = strconv.Atoi(s)
							if err != nil {
								panic(err)
							}

							N := room.Nasup
							if N > 0 {
								if room.Arsp == nil {
									room.Arsp = make([]*AIRSUP, N)
									for i := 0; i < N; i++ {
										room.Arsp[i] = new(AIRSUP)
									}
								}
							}

							Crm.Nin += 2 * room.Nasup
						}
					} else {
						comp.Nin, err = strconv.Atoi(s)
						if err != nil {
							panic(err)
						}
					}

					comp.Idi = make([]ELIOType, comp.Nin)
					for i := 0; i < comp.Nin; i++ {
						comp.Idi[i] = ' '
					}
					break
				case 'O':
					// `-Nout <分岐数>`
					var err error
					comp.Nout, err = strconv.Atoi(s)
					if err != nil {
						panic(err)
					}
					ido = make([]ELIOType, comp.Nout)
					for i := 0; i < comp.Nout; i++ {
						ido[i] = ELIO_SPACE
					}
					comp.Ido = ido
					break

				case 'L':
					// `-L <配管長>`
					l, err := readFloat(s)
					if err != nil {
						panic(err)
					}
					comp.Ivparm = CreateConstantValuePointer(l)
					break
				case 'e':
					// `-env <周囲温度>`
					comp.Envname = s
					break
				case 'r':
					// `-room <機器設置空間の室名>`
					comp.Roomname = s
					break
				case 'h':
					// `-hcc <VWV制御するときの制御対象熱交換器名称>`
					comp.Hccname = s
					break
				case 'f':
					// `-pfloor <VWV制御するときの制御対象床暖房>`
					comp.Rdpnlname = s
					break
				case 'R':
					// `-roomheff <ボイラ室内置き時の室内供給熱量率>`
					var err error
					comp.Roomname = s
					s = f.GetToken()
					comp.Eqpeff, err = readFloat(s)
					if err != nil {
						panic(err)
					}
					break
				case 's':
					// `-exs <太陽熱集熱器の方位・傾斜名>`
					comp.Exsname = s
					break
				case 'M':
					// `-control <集熱器が直列接続の場合に集熱器の要素名を流れ方向に記載する>`
					comp.Omparm = s
					break
				case 'S', 'V':
					// `-S <????>` or `-V <????>`
					for {
						_s := f.GetToken()
						if _s == "*" {
							break
						}
						s += " " + _s
					}
					s += " *"
					comp.Tparm = s
					break

				case 'T':
					// `-Tinit <蓄熱槽の初期水温>`
					if strings.HasPrefix(s, "(") {
						s += " "
						s += strings.Repeat(" ", len(s))
						_s := f.GetToken()
						s += _s
						comp.Tparm = s
					} else {
						comp.Tparm = s
					}
					break
				case 'w':
					// `-wet` 冷却コイルで出口相対湿度一定の仮定で除湿の計算を行う。
					wet, err := readFloat(s)
					if err != nil {
						panic(err)
					}
					comp.Ivparm = CreateConstantValuePointer(wet)
					break
				case 'm':
					// `-monitor` デバッグ用
					comp.MonPlistName = s
					comp.Valvcmp.MonPlistName = s
					break
				case 'P':
					// `-PCMweight <電気蓄熱暖房器の潜熱蓄熱材重量（kg）>`
					var err error
					comp.MPCM, err = readFloat(s)
					if err != nil {
						panic(err)
					}
					break
				}
			}
		}

		if Crm == nil {
			*Cmp = append(*Cmp, comp)
			D++
		}
	}

	// 空気の分岐要素・合流要素があった場合、内容をコピーして湿度用に分岐要素・合流要素を作成
	n := len(*Cmp)
	for i := 0; i < n; i++ {
		cm := (*Cmp)[i]
		if cm.Eqptype == DIVGAIR_TYPE {
			c := NewCOMPNT()
			c.Name = cm.Name + ".x"
			c.Eqptype = cm.Eqptype
			c.Nout = cm.Nout
			c.Ido = cm.Ido
			*Cmp = append(*Cmp, c)
		} else if cm.Eqptype == CVRGAIR_TYPE {
			c := NewCOMPNT()
			c.Name = cm.Name + ".x"
			c.Eqptype = cm.Eqptype
			c.Nin = cm.Nin
			c.Idi = cm.Idi
			// ここで合流要素一覧に追加する理由が不明
			Eqsys.Cnvrg = append(Eqsys.Cnvrg, nil)
			*Cmp = append(*Cmp, c)
		}
	}

	// for _, c := range *Cmp {
	// 	if len(c.Idi) != c.Nin {
	// 		panic(c.Name)
	// 	}
	// 	if len(c.Ido) != c.Nout {
	// 		panic(c.Name)
	// 	}
	// }

	//fmt.Printf("<<Compodata>> Ncompnt = %d\n", *Ncompnt)

	//printf("<<Compodata>> end\n");
}

func FindCOMPNTByName(s string, Compnt []*COMPNT) *COMPNT {
	cp := (*COMPNT)(nil)
	for i := range Compnt {
		if s == Compnt[i].Name {
			cp = Compnt[i]
			break
		}
	}
	return cp
}

func NewVALV() *VALV {
	return &VALV{
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

func NewHCC() *HCC {
	return &HCC{
		Name: "",
		Cmp:  nil,
		Cat:  nil,
		Twin: 5.0,
		Xain: FNXtr(25.0, 50.0),
	}
}

func NewBOI() *BOI {
	Boi := &BOI{
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

func NewCOLL() *COLL {
	return &COLL{
		Name: "",
		Cmp:  nil,
		Cat:  nil,
		sol:  nil,
		Te:   0.0,
		Tcb:  0.0,
	}
}

func NewPV() *PV {
	PV := &PV{
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

func NewOMVAV() *OMVAV {
	OMvav := &OMVAV{
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

func NewREFA() *REFA {
	Rf := &REFA{
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

func NewPIPE() *PIPE {
	return &PIPE{
		Name:  "",
		Cmp:   nil,
		Cat:   nil,
		Loadt: nil,
		Loadx: nil,
		Room:  nil,
	}
}

func NewSTANK() *STANK {
	return &STANK{
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

func NewHEX() *HEX {
	return &HEX{
		Name: "",
		Cmp:  nil,
		Cat:  nil,
		Id:   0,
	}
}

func NewPUMP() *PUMP {
	Pp := &PUMP{
		Name: "",
		Cmp:  nil,
		Cat:  nil,
		Sol:  nil,
	}
	MtEdayinit(&Pp.MtEdy)
	return Pp
}

func NewFLIN() *FLIN {
	return &FLIN{
		Name:  "",
		Namet: "",
		Namex: "",
		Vart:  nil,
		Varx:  nil,
		Cmp:   nil,
	}
}

func NewHCLOAD() *HCLOAD {
	Hl := &HCLOAD{
		Name:    "",
		Loadt:   nil,
		Loadx:   nil,
		Cmp:     nil,
		RMACFlg: 'N',
		Ga:      0.0,
		Gw:      0.0,
		RHout:   50.0,
	}
	MtEdayinit(&Hl.mtEdy)
	return Hl
}

func NewVAV() *VAV {
	return &VAV{
		Name: "",
		Cat:  nil,
		Hcc:  nil,
		Hcld: nil,
		Cmp:  nil,
	}
}

func NewSTHEAT() *STHEAT {
	st := &STHEAT{
		Name: "",
		Cat:  nil,
		Cmp:  nil,
		Room: nil,
		Pcm:  nil,
	}
	MtEdayinit(&st.MtEdy)
	return st
}

func NewDESI() *DESI {
	return &DESI{
		Name: "",
		Cat:  nil,
		Cmp:  nil,
		Room: nil,
		Tenv: nil,
	}
}

func NewEVAC() *EVAC {
	return &EVAC{
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

func NewTHEX() *THEX {
	return &THEX{
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

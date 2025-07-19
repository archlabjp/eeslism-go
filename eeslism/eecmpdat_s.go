package eeslism

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

/*
Compodata (Component Data Input and Initialization)

この関数は、建物のエネルギーシミュレーションで使用される様々なコンポーネント（機器）のデータを読み込み、
初期化します。
これには、組み込みのシステム要素（給水温度、外気温度など）や、
室、放射パネル、そしてカタログデータから読み込まれる各種設備機器が含まれます。

建築環境工学的な観点:
  - **システム要素の統合**: 建物のエネルギーシステムは、
    熱源設備、熱搬送設備、空調設備、換気設備など、
    様々なコンポーネントから構成されます。
    この関数は、これらのコンポーネントを統一的な`COMPNT`構造体としてモデル化し、
    システム全体を構築します。
  - **組み込みシステム要素**: `CITYWATER_NAME`（給水温度）や`OUTDRAIR_NAME`（外気温度・絶対湿度）のように、
    シミュレーションに不可欠な外部環境条件をシステム要素として組み込みます。
    これにより、建物と外部環境との熱交換を正確にモデル化できます。
  - **室と放射パネルの自動登録**: `Rmvls.Room`（室）や`Rmvls.Rdpnl`（放射パネル）は、
    建物の熱負荷計算の基本となる要素であり、
    この関数によって自動的にコンポーネントとして登録されます。
    これにより、室内の熱的挙動や放射冷暖房の効果をモデル化できます。
  - **カタログデータとの連携**: `Eqcat`（機器カタログ）から読み込まれたデータに基づいて、
    各機器のタイプ（`Eqptype`）、カタログ番号（`Ncat`）、
    入出力ポートの数（`Nin`, `Nout`）などを設定します。
    これにより、多様な設備機器の性能をモデル化し、
    システム全体のエネルギー消費量を評価できます。
  - **三方弁（二方弁で連動する弁）のモデル化**: `strings.HasPrefix(s, "(")` の条件で、
    三方弁のように連動する二方弁のペアをモデル化します。
    これは、熱源システムや熱搬送システムにおける流量制御を正確にシミュレーションするために重要です。
  - **熱湿気同時交換の考慮**: `comp.Airpathcpy = true` と設定される機器は、
    空気の温度だけでなく、湿度も同時に変化させる熱湿気同時交換を行う機器であることを示唆します。
    これにより、室内空気質や潜熱負荷の計算を正確に行えます。

この関数は、建物のエネルギーシミュレーションにおいて、
多様なコンポーネントを統合的にモデル化し、
システム全体のエネルギー消費量予測、省エネルギー対策の検討、
および最適な設備システム設計を行うための重要な役割を果たします。
*/
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

/*
FindCOMPNTByName (Find Component by Name)

この関数は、与えられた名称（`s`）に基づいて、
コンポーネントのリスト（`Compnt`）から該当するコンポーネントを検索します。

建築環境工学的な観点:
  - **コンポーネントの参照**: 建物のエネルギーシミュレーションでは、
    様々なコンポーネントが相互に接続され、熱や空気、水などをやり取りします。
    この関数は、あるコンポーネントが別のコンポーネントを参照する際に、
    その名称に基づいて対象のコンポーネントを効率的に見つけ出すために用いられます。
  - **システム構成の動的な構築**: シミュレーションモデルを構築する際、
    コンポーネント間の接続関係は、入力ファイルから読み込まれる名称に基づいて動的に設定されます。
    この関数は、その動的な接続を可能にするための基本的な検索機能を提供します。

この関数は、建物のエネルギーシミュレーションにおいて、
コンポーネント間の接続関係を確立し、
システム全体の熱・空気・水の流れをモデル化するための重要な役割を果たします。
*/
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

/*
NewVALV (New Valve/Damper Object)

この関数は、新しい弁（バルブ）またはダンパーのデータ構造を初期化します。
弁やダンパーは、熱媒（水、空気）の流量を制御し、
熱搬送システムや空調システムにおける熱供給量を調整する重要な機器です。

建築環境工学的な観点:
  - **流量制御のモデル化**: 弁やダンパーは、
    熱媒の流量を調整することで、
    熱源設備から熱利用設備への熱供給量を制御します。
    この関数は、弁やダンパーの初期状態（開度、制御モードなど）を定義し、
    流量制御のモデル化を可能にします。
  - **熱負荷への対応**: 流量制御によって、
    室の熱負荷変動に応じて熱供給量を調整し、
    室内温度の安定化や、熱源設備の効率的な運転に貢献します。
  - **省エネルギー運転**: 適切な流量制御は、
    不要な熱供給を防ぎ、エネルギーの無駄を削減できます。
    例えば、室温が設定値に達した場合に弁を閉じることで、
    熱供給を停止し、エネルギー消費を削減できます。

この関数は、熱搬送システムや空調システムにおける流量制御をモデル化し、
熱負荷計算、エネルギー消費量予測、
および省エネルギー対策の検討を行うための重要な役割を果たします。
*/
func NewVALV() *VALV {
	return &VALV{
		Name:     "",
		Cmp:      nil,
		Cmb:      nil,
		Org:      'n',
		X:        FNAN,
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

/*
NewHCC (New Heating/Cooling Coil Object)

この関数は、新しい冷温水コイルのデータ構造を初期化します。
冷温水コイルは、空調システムにおいて空気と熱媒（冷水または温水）の間で熱を交換する主要な機器です。

建築環境工学的な観点:
  - **熱交換器の初期化**: 冷温水コイルのシミュレーションを行う前に、
    その熱的特性（温度効率、エンタルピー効率、熱通過率と伝熱面積の積など）をデフォルト値で初期化します。
    これにより、後続のデータ入力や計算が正しく行われるための準備が整います。
  - **熱負荷処理のモデル化**: 冷温水コイルは、
    室内の顕熱負荷と潜熱負荷を処理し、
    室内の温湿度環境を維持する役割を担います。
    この関数で初期化されるパラメータは、
    コイルの熱交換能力や、空気と水の温度変化を決定する上で重要です。
  - **省エネルギー運転**: 効率の良い冷温水コイルを選定したり、
    コイルの運転条件を最適化したりすることで、
    空調システムのエネルギー消費量を削減できます。

この関数は、空調システムにおける熱交換器の性能をモデル化し、
室内の温湿度環境の維持、熱負荷計算、
およびエネルギー消費量予測を行うための重要な役割を果たします。
*/
func NewHCC() *HCC {
	return &HCC{
		Name: "",
		Cmp:  nil,
		Cat:  nil,
		Twin: 5.0,
		Xain: FNXtr(25.0, 50.0),
	}
}

/*
NewBOI (New Boiler Object)

この関数は、新しいボイラーのデータ構造を初期化します。
ボイラーは、建物に熱を供給する主要な熱源設備の一つです。

建築環境工学的な観点:
  - **熱源設備の初期化**: ボイラーのシミュレーションを行う前に、
    その性能（定格出力、効率、最小出力など）をデフォルト値で初期化します。
    これにより、後続のデータ入力や計算が正しく行われるための準備が整います。
  - **熱負荷への対応**: ボイラーは、
    建物の暖房負荷や給湯負荷に熱を供給する役割を担います。
    この関数で初期化されるパラメータは、
    ボイラーの熱供給能力や、熱媒の温度変化を決定する上で重要です。
  - **省エネルギー運転**: 効率の良いボイラーを選定したり、
    ボイラーの運転条件を最適化したりすることで、
    熱源システムのエネルギー消費量を削減できます。
  - **月・時刻別集計の準備**: `MtEdayinit(&Boi.MtEdy)`や`MtEdayinit(&Boi.MtPhdy)`は、
    月・時刻別のエネルギー消費量や電力消費量を集計するためのデータ構造を初期化します。
    これにより、ボイラーの運転状況を詳細に分析し、
    デマンドサイドマネジメントや、エネルギー供給計画を最適化する上で非常に有用な情報となります。

この関数は、熱源設備としてのボイラーの性能をモデル化し、
熱負荷計算、エネルギー消費量予測、
および省エネルギー対策の検討を行うための重要な役割を果たします。
*/
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

/*
NewCOLL (New Collector Object)

この関数は、新しい太陽熱集熱器のデータ構造を初期化します。
太陽熱集熱器は、太陽エネルギーを熱として利用し、
給湯や暖房に利用する再生可能エネルギー設備です。

建築環境工学的な観点:
  - **太陽熱利用の初期化**: 太陽熱集熱器のシミュレーションを行う前に、
    その性能（集熱効率、熱損失係数など）をデフォルト値で初期化します。
    これにより、後続のデータ入力や計算が正しく行われるための準備が整います。
  - **再生可能エネルギーの利用**: 太陽熱集熱器は、
    再生可能エネルギーである太陽光を直接利用するため、
    環境負荷の低減やエネルギー自立性の向上に貢献します。
    この関数で初期化されるパラメータは、
    集熱器の熱取得能力や、熱媒の温度変化を決定する上で重要です。
  - **省エネルギー運転**: 太陽熱集熱器を導入することで、
    従来の熱源設備からのエネルギー消費量を削減できます。
    この関数は、そのようなシステムのモデル化の基礎となります。

この関数は、太陽熱利用システムの性能をモデル化し、
熱負荷計算、エネルギー消費量予測、
および省エネルギー対策の検討を行うための重要な役割を果たします。
*/
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

/*
NewPV (New Photovoltaic Object)

この関数は、新しい太陽光発電（PV）システムのデータ構造を初期化します。
太陽光発電システムは、太陽エネルギーを電力に変換する再生可能エネルギー設備です。

建築環境工学的な観点:
  - **太陽光発電の初期化**: 太陽光発電システムのシミュレーションを行う前に、
    その性能（定格出力、面積など）をデフォルト値で初期化します。
    これにより、後続のデータ入力や計算が正しく行われるための準備が整います。
  - **再生可能エネルギーの利用**: 太陽光発電システムは、
    再生可能エネルギーである太陽光を直接利用するため、
    環境負荷の低減やエネルギー自立性の向上に貢献します。
    この関数で初期化されるパラメータは、
    発電量や、パネルの温度変化を決定する上で重要です。
  - **省エネルギー運転**: 太陽光発電システムを導入することで、
    従来の電力消費量を削減できます。
    この関数は、そのようなシステムのモデル化の基礎となります。
  - **月・時刻別集計の準備**: `MtEdayinit(&PV.mtEdy)`は、
    月・時刻別の発電量（エネルギー消費量）を集計するためのデータ構造を初期化します。
    これにより、太陽光発電システムの運転状況を詳細に分析し、
    エネルギーマネジメントや、電力系統への影響を評価する上で非常に有用な情報となります。

この関数は、太陽光発電システムの性能をモデル化し、
エネルギー消費量予測、および省エネルギー対策の検討を行うための重要な役割を果たします。
*/
func NewPV() *PV {
	PV := &PV{
		PVcap: FNAN,
		Area:  FNAN,
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

/*
NewOMVAV (New Outdoor Air VAV Object)

この関数は、新しい外気処理VAV（Variable Air Volume）ユニットのデータ構造を初期化します。
外気処理VAVユニットは、外気導入量を可変制御することで、
換気による熱負荷を最適化し、省エネルギーに貢献する空調システムです。

建築環境工学的な観点:
  - **外気処理の最適化**: 換気は室内空気質を維持するために不可欠ですが、
    同時に熱負荷を発生させます。
    外気処理VAVユニットは、室内の熱負荷やCO2濃度などに応じて外気導入量を可変制御することで、
    換気による熱負荷を最小限に抑え、省エネルギーに貢献します。
  - **流量制御のモデル化**: この関数は、
    外気処理VAVユニットの流量制御に関するパラメータ（最大流量、最小流量など）を初期化します。
    これにより、外気導入量の可変制御をモデル化し、
    換気による熱負荷の変動を正確に予測できます。
  - **システム統合**: 外気処理VAVユニットは、
    空調システム全体の一部として機能します。
    この関数で初期化されるパラメータは、
    外気処理VAVユニットが空調システム全体のエネルギー消費量や室内環境に与える影響を評価する上で重要です。

この関数は、換気システムにおける外気処理の最適化をモデル化し、
熱負荷計算、エネルギー消費量予測、
および省エネルギー対策の検討を行うための重要な役割を果たします。
*/
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

/*
NewREFA (New Refrigeration/Air Conditioning Unit Object)

この関数は、新しい冷凍機または空調機のデータ構造を初期化します。
冷凍機や空調機は、建物に冷熱を供給する主要な熱源設備の一つです。

建築環境工学的な観点:
  - **冷熱源設備の初期化**: 冷凍機や空調機のシミュレーションを行う前に、
    その性能（定格能力、効率、部分負荷特性など）をデフォルト値で初期化します。
    これにより、後続のデータ入力や計算が正しく行われるための準備が整います。
  - **冷房負荷への対応**: 冷凍機や空調機は、
    建物の冷房負荷や除湿負荷に冷熱を供給する役割を担います。
    この関数で初期化されるパラメータは、
    機器の冷熱供給能力や、熱媒の温度変化を決定する上で重要です。
  - **省エネルギー運転**: 効率の良い冷凍機や空調機を選定したり、
    運転条件を最適化したりすることで、
    冷熱源システムのエネルギー消費量を削減できます。
  - **月・時刻別集計の準備**: `MtEdayinit(&Rf.mtEdy)`や`MtEdayinit(&Rf.mtPhdy)`は、
    月・時刻別のエネルギー消費量や電力消費量を集計するためのデータ構造を初期化します。
    これにより、冷凍機や空調機の運転状況を詳細に分析し、
    デマンドサイドマネジメントや、エネルギー供給計画を最適化する上で非常に有用な情報となります。

この関数は、冷熱源設備としての冷凍機や空調機の性能をモデル化し、
熱負荷計算、エネルギー消費量予測、
および省エネルギー対策の検討を行うための重要な役割を果たします。
*/
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

/*
NewPIPE (New Pipe/Duct Object)

この関数は、新しい配管またはダクトのデータ構造を初期化します。
配管やダクトは、熱媒（水、空気）を搬送し、
熱源設備から熱利用設備へ熱を運ぶための重要な要素です。

建築環境工学的な観点:
  - **熱搬送のモデル化**: 配管やダクトは、
    熱媒を搬送する際に、熱損失や熱取得が発生します。
    この関数は、配管やダクトの熱的特性（熱損失係数など）をデフォルト値で初期化し、
    熱搬送における熱損失・熱取得をモデル化します。
  - **熱負荷への影響**: 配管やダクトからの熱損失・熱取得は、
    熱源設備や熱利用設備の負荷に影響を与えます。
    特に、長距離の搬送や、断熱性能が低い配管・ダクトでは、
    熱損失・熱取得が無視できない場合があります。
  - **省エネルギー運転**: 配管やダクトの断熱性能を向上させることで、
    熱損失・熱取得を削減し、熱搬送システムのエネルギー効率を向上させることができます。

この関数は、熱搬送システムにおける熱損失・熱取得をモデル化し、
熱負荷計算、エネルギー消費量予測、
および省エネルギー対策の検討を行うための重要な役割を果たします。
*/
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

/*
NewSTANK (New Storage Tank Object)

この関数は、新しい蓄熱槽のデータ構造を初期化します。
蓄熱槽は、熱源設備で発生した熱（または冷熱）を一時的に貯蔵し、
熱需要に応じて供給することで、熱負荷の平準化や熱源設備の効率的な運転を可能にします。

建築環境工学的な観点:
  - **蓄熱槽の初期化**: 蓄熱槽のシミュレーションを行う前に、
    その性能（容量、熱損失係数など）をデフォルト値で初期化します。
    これにより、後続のデータ入力や計算が正しく行われるための準備が整います。
  - **熱負荷平準化のモデル化**: 蓄熱槽は、
    熱源設備の容量を小さくしたり、
    電力料金の安い夜間電力などを利用したりすることができ、
    省エネルギーやランニングコストの削減に貢献します。
    この関数で初期化されるパラメータは、
    蓄熱槽の蓄熱能力や、熱損失を決定する上で重要です。
  - **温度成層の考慮**: 蓄熱槽内部の温度成層をモデル化するために、
    複数の層に分割された温度データ（`Tss`, `Tssold`）や、
    各層の熱容量（`Mdt`）、熱損失係数（`KS`）などを初期化します。
    これにより、蓄熱槽の有効利用率を評価できます。

この関数は、蓄熱槽の性能をモデル化し、
熱負荷平準化、エネルギー消費量予測、
および省エネルギー対策の検討を行うための重要な役割を果たします。
*/
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

/*
NewHEX (New Heat Exchanger Object)

この関数は、新しい熱交換器のデータ構造を初期化します。
熱交換器は、異なる温度の流体間で熱を交換する機器であり、
空調システム、給湯システム、熱回収システムなど、様々な用途で用いられます。

建築環境工学的な観点:
  - **熱交換器の初期化**: 熱交換器のシミュレーションを行う前に、
    その性能（効率、熱通過率と伝熱面積の積など）をデフォルト値で初期化します。
    これにより、後続のデータ入力や計算が正しく行われるための準備が整います。
  - **熱回収のモデル化**: 熱交換器は、排気や排水などから熱を回収し、
    新鮮な空気や水に熱を供給することで、エネルギーの有効利用を促進します。
    この関数で初期化されるパラメータは、
    熱交換器の熱交換能力や、各流体の温度変化を決定する上で重要です。
  - **省エネルギー運転**: 熱交換器を導入することで、
    熱源設備の負荷を軽減し、システム全体のエネルギー消費量を削減できます。

この関数は、熱交換器の性能をモデル化し、
熱回収システムや熱源システムの設計、熱負荷計算、
およびエネルギー消費量予測を行うための重要な役割を果たします。
*/
func NewHEX() *HEX {
	return &HEX{
		Name: "",
		Cmp:  nil,
		Cat:  nil,
		Id:   0,
	}
}

/*
NewPUMP (New Pump/Fan Object)

この関数は、新しいポンプまたはファンのデータ構造を初期化します。
ポンプやファンは、熱媒（水、空気）を搬送し、
熱を必要な場所へ運ぶための動力源です。

建築環境工学的な観点:
  - **熱搬送の動力の初期化**: ポンプやファンのシミュレーションを行う前に、
    その性能（定格流量、定格消費電力、部分負荷特性など）をデフォルト値で初期化します。
    これにより、後続のデータ入力や計算が正しく行われるための準備が整います。
  - **熱負荷への影響**: ポンプやファンは、
    熱媒を搬送する際に、その熱媒に動力を供給します。
    この動力は、熱媒の温度上昇として現れ、熱負荷の一部となります。
    この関数で初期化されるパラメータは、
    ポンプやファンの熱搬送能力や、エネルギー消費量を決定する上で重要です。
  - **省エネルギー運転**: 効率の良いポンプやファンを選定したり、
    運転条件を最適化したりすることで、
    熱搬送システムのエネルギー消費量を削減できます。
  - **月・時刻別集計の準備**: `MtEdayinit(&Pp.MtEdy)`は、
    月・時刻別のエネルギー消費量（電力消費量）を集計するためのデータ構造を初期化します。
    これにより、ポンプやファンの運転状況を詳細に分析し、
    デマンドサイドマネジメントや、エネルギー供給計画を最適化する上で非常に有用な情報となります。

この関数は、熱搬送システムにおけるポンプやファンの性能をモデル化し、
熱負荷計算、エネルギー消費量予測、
および省エネルギー対策の検討を行うための重要な役割を果たします。
*/
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

/*
NewFLIN (New Fluid Inlet Object)

この関数は、新しい流入境界条件（Fluid Inlet）のデータ構造を初期化します。
流入境界条件は、シミュレーションモデルに外部から流入する熱媒（空気、水など）の
温度や湿度、流量などの状態を定義するために用いられます。

建築環境工学的な観点:
  - **外部環境との相互作用のモデル化**: 建物のエネルギーシミュレーションでは、
    外部環境（外気、給水など）からの熱媒の流入が、
    建物全体の熱負荷や室内環境に大きな影響を与えます。
    この関数は、流入する熱媒の温度（`Vart`）や湿度（`Varx`）を定義し、
    外部環境との相互作用をモデル化します。
  - **熱負荷計算の基礎**: 流入する熱媒の温度や湿度は、
    空調システムや熱源設備の熱負荷計算の基礎となります。
    例えば、外気導入量が多い場合、外気温度や湿度が高いと、
    冷房負荷が増加する可能性があります。

この関数は、建物のエネルギーシミュレーションにおいて、
外部環境との相互作用を正確にモデル化し、
熱負荷計算、エネルギー消費量予測、
および省エネルギー対策の検討を行うための重要な役割を果たします。
*/
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

/*
NewHCLOAD (New Heating/Cooling Load Object)

この関数は、新しい空調負荷（Heating/Cooling Load）のデータ構造を初期化します。
空調負荷は、室内の温湿度環境を目標値に維持するために、
空調システムが供給または除去する必要がある熱量を示します。

建築環境工学的な観点:
  - **熱負荷のモデル化**: 室内の熱負荷は、
    透過熱負荷、日射熱負荷、内部発熱、換気熱負荷など、
    様々な要因によって構成されます。
    この関数は、空調負荷の計算に必要なパラメータ（目標温度`Loadt`、目標湿度`Loadx`など）を初期化し、
    熱負荷のモデル化を可能にします。
  - **室内温湿度環境の維持**: 空調システムは、
    この空調負荷を処理することで、
    室内の温湿度環境を目標値に維持します。
    この関数で初期化されるパラメータは、
    空調システムの設計や運転制御において重要な情報となります。
  - **省エネルギー運転**: 適切な空調負荷の計算は、
    過剰な冷暖房や除湿を防ぎ、エネルギーの無駄を削減できます。
    また、`RMACFlg`（室温追従制御フラグ）や`RHout`（出口相対湿度）などのパラメータは、
    空調システムの制御戦略をモデル化し、
    省エネルギー運転を検討する上で重要です。
  - **月・時刻別集計の準備**: `MtEdayinit(&Hl.mtEdy)`は、
    月・時刻別のエネルギー消費量（電力消費量）を集計するためのデータ構造を初期化します。
    これにより、空調システムの運転状況を詳細に分析し、
    デマンドサイドマネジメントや、エネルギー供給計画を最適化する上で非常に有用な情報となります。

この関数は、室の熱負荷をモデル化し、
空調システムの設計、エネルギー消費量予測、
および省エネルギー対策の検討を行うための重要な役割を果たします。
*/
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

/*
NewVAV (New Variable Air Volume Unit Object)

この関数は、新しいVAV（Variable Air Volume）ユニットのデータ構造を初期化します。
VAVユニットは、室の熱負荷に応じて送風量を可変制御することで、
省エネルギーと快適性を両立させる空調システムです。

建築環境工学的な観点:
  - **送風量制御のモデル化**: VAVユニットは、
    室内の熱負荷変動に応じて送風量を調整し、
    室温を目標値に維持します。
    この関数は、VAVユニットの制御に関するパラメータ（制御対象の熱交換器`Hcc`、
    空調負荷`Hcld`など）を初期化し、
    送風量制御のモデル化を可能にします。
  - **省エネルギー運転**: VAVシステムは、
    定風量システムに比べて送風機の消費電力を大幅に削減できるため、
    省エネルギーに貢献します。
    この関数で初期化されるパラメータは、
    VAVシステムの省エネルギー効果を評価する上で重要です。
  - **快適性の維持**: 送風量をきめ細かく制御することで、
    室内の温度変動を抑制し、
    居住者の快適性を向上させることができます。

この関数は、VAVシステムにおける送風量制御をモデル化し、
熱負荷計算、エネルギー消費量予測、
および省エネルギー対策の検討を行うための重要な役割を果たします。
*/
func NewVAV() *VAV {
	return &VAV{
		Name: "",
		Cat:  nil,
		Hcc:  nil,
		Hcld: nil,
		Cmp:  nil,
	}
}

/*
NewSTHEAT (New Sensible Heat Storage Object)

この関数は、新しい顕熱蓄熱式暖房器のデータ構造を初期化します。
顕熱蓄熱式暖房器は、電気ヒーターなどで熱を発生させ、
その熱を蓄熱材に顕熱として蓄え、必要な時に放熱することで暖房を行うシステムです。

建築環境工学的な観点:
  - **顕熱蓄熱の初期化**: 顕熱蓄熱式暖房器のシミュレーションを行う前に、
    その性能（定格熱量、効率、熱容量、熱通過率など）をデフォルト値で初期化します。
    これにより、後続のデータ入力や計算が正しく行われるための準備が整います。
  - **熱負荷平準化のモデル化**: 蓄熱式暖房器は、
    電力料金の安い夜間電力などを利用して蓄熱し、
    昼間のピークカットや熱負荷平準化に貢献します。
    この関数で初期化されるパラメータは、
    蓄熱式暖房器の蓄熱能力や、熱損失を決定する上で重要です。
  - **PCMの利用**: `Pcm`は、
    相変化材料（PCM）が利用されている場合に、
    その特性を格納するためのポインターです。
    PCMは、潜熱を利用して顕熱蓄熱よりも大きな熱量を蓄えることができ、
    よりコンパクトな蓄熱システムを実現できます。
  - **月・時刻別集計の準備**: `MtEdayinit(&st.MtEdy)`は、
    月・時刻別のエネルギー消費量（電力消費量）を集計するためのデータ構造を初期化します。
    これにより、蓄熱式暖房器の運転状況を詳細に分析し、
    デマンドサイドマネジメントや、エネルギー供給計画を最適化する上で非常に有用な情報となります。

この関数は、顕熱蓄熱式暖房器の性能をモデル化し、
熱負荷平準化、エネルギー消費量予測、
および省エネルギー対策の検討を行うための重要な役割を果たします。
*/
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

/*
NewDESI (New Desiccant Air Conditioner Object)

この関数は、新しいデシカント空調機のデータ構造を初期化します。
デシカント空調機は、吸湿材を用いて空気中の水蒸気を除去することで除湿を行い、
その後、顕熱交換によって温度を調整するシステムです。

建築環境工学的な観点:
  - **デシカント空調の初期化**: デシカント空調機のシミュレーションを行う前に、
    その性能（吸湿材の種類、量、デシカント槽の熱的特性など）をデフォルト値で初期化します。
    これにより、後続のデータ入力や計算が正しく行われるための準備が整います。
  - **潜熱負荷処理の最適化**: デシカント空調は、特に潜熱負荷（湿度）の処理に優れています。
    この関数で初期化されるパラメータは、
    デシカント空調機の除湿性能や、熱損失を決定する上で重要です。
  - **省エネルギー運転**: デシカント空調機は、
    吸湿材の再生に熱エネルギーを必要としますが、
    従来の冷媒を用いた空調システムに比べて、
    特に潜熱負荷が大きい場合に省エネルギー効果が期待できます。

この関数は、デシカント空調機の性能をモデル化し、
室内温湿度環境の予測、潜熱負荷の処理、
およびエネルギー消費量予測を行うための重要な役割を果たします。
*/
func NewDESI() *DESI {
	return &DESI{
		Name: "",
		Cat:  nil,
		Cmp:  nil,
		Room: nil,
		Tenv: nil,
	}
}

/*
NewEVAC (New Evaporative Cooler Object)

この関数は、新しい気化冷却器のデータ構造を初期化します。
気化冷却器は、水の蒸発潜熱を利用して空気を冷却するシステムであり、
特に乾燥地域において省エネルギーな冷房手段として注目されています。

建築環境工学的な観点:
  - **気化冷却の原理**: 気化冷却器は、
    水が蒸発する際に周囲から熱を奪う現象（蒸発潜熱）を利用して空気を冷却します。
    これにより、冷媒を用いた空調システムに比べて、
    大幅なエネルギー削減が期待できます。
  - **冷却性能の初期化**: 気化冷却器のシミュレーションを行う前に、
    その性能（冷却効率、空気流量など）をデフォルト値で初期化します。
    これにより、後続のデータ入力や計算が正しく行われるための準備が整います。
  - **湿度の影響**: 気化冷却器は、
    空気を冷却すると同時に加湿するため、
    特に高湿度の地域では適用が難しい場合があります。
    この関数で初期化されるパラメータは、
    気化冷却器の冷却性能と加湿特性をモデル化する上で重要です。
  - **省エネルギー運転**: 気化冷却器は、
    従来の空調システムに比べて消費電力が少ないため、
    省エネルギーに貢献します。

この関数は、気化冷却器の性能をモデル化し、
室内温湿度環境の予測、冷房負荷の処理、
およびエネルギー消費量予測を行うための重要な役割を果たします。
*/
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

/*
NewTHEX (New Total Heat Exchanger Object)

この関数は、新しい全熱交換器のデータ構造を初期化します。
全熱交換器は、排気から顕熱（温度）と潜熱（湿度）の両方を回収し、
給気に熱と湿気を供給することで、換気による熱損失を最小限に抑える機器です。

建築環境工学的な観点:
  - **全熱交換器の初期化**: 全熱交換器のシミュレーションを行う前に、
    その性能（顕熱交換効率、潜熱交換効率など）をデフォルト値で初期化します。
    これにより、後続のデータ入力や計算が正しく行われるための準備が整います。
  - **熱回収の重要性**: 換気は、室内空気質を維持するために不可欠ですが、
    同時に熱損失（または熱取得）を伴います。
    全熱交換器は、この換気による熱損失を大幅に削減し、
    空調負荷を軽減することで、建物のエネルギー消費量を削減します。
  - **熱湿気同時交換のモデル化**: 全熱交換器は、
    顕熱と潜熱の両方を交換するため、
    空気の温度と湿度の両方に影響を与えます。
    この関数で初期化されるパラメータは、
    全熱交換器の熱湿気交換能力をモデル化する上で重要です。

この関数は、換気システムにおける熱回収の設計、熱負荷計算、
およびエネルギー消費量予測を行うための重要な役割を果たします。
*/
func NewTHEX() *THEX {
	return &THEX{
		Name: "",
		Cmp:  nil,
		Cat:  nil,
		Type: ' ',
		CGe:  0.0,
		Ge:   0.0,
		Go:   0.0,
		ET:   FNAN,
		EH:   FNAN,
	}
}

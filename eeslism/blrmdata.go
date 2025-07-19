//This file is part of EESLISM.
//
//Foobar is free software : you can redistribute itand /or modify
//it under the terms of the GNU General Public License as published by
//the Free Software Foundation, either version 3 of the License, or
//(at your option) any later version.
//
//Foobar is distributed in the hope that it will be useful,
//but WITHOUT ANY WARRANTY; without even the implied warranty of
//MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.See the
//GNU General Public License for more details.
//
//You should have received a copy of the GNU General Public License
//along with Foobar.If not, see < https://www.gnu.org/licenses/>.

/*  rmdata.c   */
package eeslism

import (
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

/* -------------------------------------------- */

/*  室構成部材の入力  */

func Roomdata(tokens *EeTokens, Exs []*EXSF, dfwl *DFWL, Rmvls *RMVLS, Schdl *SCHDL, Simc *SIMCONTL) {
	// var Wall, w *WALL
	// var Window, W *WINDOW
	// var Snbk, S *SNBK
	// var Room, room, Rm, Rc, Rmchk *ROOM
	// var rdpnl, Rd *RDPNL
	// var Sd, rsd, nxsd, Sdj *RMSRF
	// var Nroom, Nrdpnl, Nsrf int
	// var N, Nnxrm int
	// var Nwindow, Nwall, Nexs, Nsnbk int
	// var Nairflow int
	// var e *EXSF
	// var Scw, Sch *SCH
	// var Schdl.Val []float64
	// var Ac *ACHIR
	// var ca, roa float64
	// var NSTOP int
	// var i int

	errkey := "Roomdata"

	//i := -1
	var j, n, nr, brs, ij, N2, k, l int
	n = -1
	brs = 0
	//var s, ss string
	//var st, ce, stt string
	var dexsname string // 外壁の名前 => EXSRFで定義される
	var dnxrname string // 内壁の名前
	var Er string
	var sfemark bool
	var RmnameEr string

	//stt := ""
	//sprintf(s, "No. 1") ;
	//HeapCheck(s) ;

	// 部屋数
	N := Roomcount(tokens)

	//printf ( "Nroom=%d\n", N ) ;

	if N > 0 {
		Rmvls.Room = make([]*ROOM, N)

		Roominit(N, Rmvls.Room)
	}

	// 部屋を構成する壁、床、天井等の数
	Nnxrm := Rmsrfcount(tokens)

	//printf ( "Nsrf=%d\n", N ) ;

	if Nnxrm > 0 {
		Rmvls.Sd = make([]*RMSRF, 0, Nnxrm)
	}

	//Wall := Rmvls.Wall
	//Window := Rmvls.Window
	//W := &Window[0]
	//Snbk := Rmvls.Snbk
	RmIdx := 0
	//Room := Rmvls.Room

	SdIdx := 0
	//rdpnl := Rmvls.Rdpnl

	Er = fmt.Sprintf(ERRFMT, errkey)
	RmIdx--
	SdIdx--

	var i int = -1
	for _i := 0; _i < N; _i++ {
		section := tokens.GetSection()

		// 部屋についての一行目の情報を処理
		s := section.GetToken()
		if s == "*" {
			break
		}

		/*****************************/

		if DEBUG {
			fmt.Printf("%s\n", s)
		}

		//err = Er + s

		i++
		RmIdx++
		Rm := Rmvls.Room[RmIdx]

		Rm.Name = s

		// Check duplication of room names
		for l = 0; l < i-1; l++ {
			Rmchk := Rmvls.Room[l]
			if Rm.Name == Rmchk.Name {
				RmnameEr = fmt.Sprintf("Room=%s is already defined name", Rm.Name)
				Eprint("<Roomdata>", RmnameEr)
			}
		}

		nr = -1

		// 室のデータの読み取り
		compnameStart := false // 部位の設定を読み取り中かどうかを示すフラグ
		for section.IsEnd() == false {
			s := section.PeekToken()
			if s == "\n" {
				section.GetToken()
				continue
			}
			if s == "*" {
				break
			}

			line := section.GetLogicalLine()

			n++
			nr++

			SdIdx++

			// 壁体データの初期化
			Sd := Rmsrfinit()
			var err error

			for i := 0; i < len(line)-1; i++ {
				s := line[i]
				if DEBUG {
					fmt.Printf("Roomdata1  s=%s\n", s)
				}

				var st int
				st = strings.IndexRune(s, '=')
				if st == -1 {
					if DEBUG {
						fmt.Printf("Roomdata2  s=%s\n", s)
					}
					/*******************/

					if s == "*s" {
						// outfile_sf.es への室内表面温度、
						// outfile_sfq.esへの部位別表面熱流、
						// outfile_sfa.esへの部位別表面熱伝達率の出力指定
						Rm.sfpri = true
					} else if s == "*q" {
						// outfile_rq.es、outfile_dqr.es への日射熱取得、
						// 室内発熱、隙間風熱取得要素の出力指定
						Rm.eqpri = true
					} else if s == "*sfe" {
						// 要素別壁体表面温度出力指定
						if !compnameStart {
							// 部屋の内表面に全体に反映させる
							sfemark = true
						} else {
							// 個別内表面にのみ反映させる
							Sd.sfepri = true
						}
					} else if s == "*p" {
						// 壁体内部温度出力指定
						Sd.wlpri = true
					} else if s == "*shd" {
						// 日よけの影面積出力
						Sd.shdpri = true
					} else if st = strings.Index(s, "*"); st != -1 {
						// 面積の指定 (幅・高さ指定)
						var X, Y float64
						fmt.Sscanf(s, "%f*%f", &X, &Y)
						Sd.A = X * Y
					} else if unicode.IsDigit(rune(s[0])) {
						// 面積の指定 (直接指定)
						_, err = fmt.Sscanf(s, "%f", &Sd.A)
						if err != nil {
							panic(err)
						}
					} else if s == "if" {
						// 動的にカーテンを開閉するロジックを追加  2012/2/25 Satoh
						Sd.Ctlif = new(CTLIF)

						// Read until the end of the if statement
						// ex `if (LD_Tr > 27)`
						// line[i]を読み進めて `)` まで読み進める
						ii := i + 1
						for {
							ss := line[ii]
							//ssが`)`が終わるか
							if ss[len(ss)-1] == ')' {
								break
							}
							ii = ii + 1
						}
						Sd.DynamicCode = strings.Trim(strings.Join(line[i+1:ii+1], " "), "()")

						// IF文の展開
						//ctifdecode(s, Sd->ctlif, Simc, Compnt, Nmpath, Mpath, Wd, Exsf, Schdl);

						// Read the True case window
						target_window := line[ii+1]
						found := false
						for j := range Rmvls.Window {
							W := Rmvls.Window[j]
							if W.Name == target_window {
								Sd.ifwin = W
								Sd.Rwall = W.Rwall
								Sd.CAPwall = 0.0
								found = true
								break
							}
						}

						if found == false {
							err := fmt.Sprintf("Room=%s <window> %s", Rm.Name, s)
							Eprint("<Roomdata>", err)
							os.Exit(1)
						}

						// 読み進めの記録
						i = ii + 1
					} else if strings.HasSuffix(s, ":") {
						// ex: `west: -E	7.50 ;`
						// ex: `(roomA):	-i	Room	i=roomA ;``

						ss := strings.Split(s, ":")

						// 部位の種類を判定するために、括弧を確認する
						const regexPattern = `\(([^)]*)\)`

						// 正規表現をコンパイルします。
						r, err := regexp.Compile(regexPattern)
						if err != nil {
							fmt.Printf("Error compiling regex: %v\n", err)
							return
						}

						// 文字列からマッチする部分を検索し、結果を取得します。
						match := r.FindStringSubmatch(ss[0])

						// マッチした内容が存在し、サブマッチ（ここではカッコ内）が存在する場合、それを表示します。
						if len(match) > 1 {
							// 内壁の名前
							dnxrname = match[1]
						} else {
							// 外壁の名前
							dexsname = ss[0]
						}
						compnameStart = true // 部位の設定を読み取り開始
					} else if s == "Fij" {
						// 形態係数の入力 （予め計算済みの形態係数を使用する場合）
						Rm.fij = 'F'

						// 室内の表面数 N
						N2, err = strconv.Atoi(section.GetToken())
						if err != nil {
							panic(err)
						}
						Rm.F = make([]float64, N2*N2)

						ij = 0
						for {
							var ss string
							ss = section.GetToken()
							if err != nil {
								panic(err)
							}
							if ss == ";" {
								break
							}
							Rm.F[ij], err = readFloat(ss)
							if err != nil {
								panic(err)
							}
							ij++
						}
					} else if s == "rsrnx" {
						Rm.rsrnx = true
					} else if s[0] == '-' {
						// `-E` の様に記述することで、部位の種類を指定する
						c := s[1]

						// 不正なBLEの場合はエラーを出力
						if c != 'E' && c != 'R' && c != 'F' && c != 'i' && c != 'c' && c != 'f' && c != 'W' {
							panic(fmt.Sprintf("Invalid ble '%s' at \"%s\"", string(c), strings.Join(line, " ")))
						}

						// 部位の種類を設定
						Sd.ble = BLEType(c)
					} else if compnameStart {
						c := Sd.ble
						if DEBUG {
							fmt.Printf("Roomdata3  s=%s  c=%c\n", s, c)
						}

						// 読み取り中の部位が窓の場合
						if c == BLE_Window {
							nf := Sd.Nfn
							Sd.Nfn++

							var stt string
							if sttIndex := strings.IndexByte(s, ':'); sttIndex != -1 {
								stt = s[sttIndex+1:]
								Sd.fnmrk[nf] = rune(s[0])
							} else {
								stt = s
							}

							for j := range Rmvls.Window {
								W := Rmvls.Window[j]
								if W.Name == stt {
									Sd.window = W
									Sd.fnd[nf] = j
									Sd.Rwall = W.Rwall
									Sd.CAPwall = 0.0
									Sd.RStrans = W.RStrans
									break
								}
							}

							if j == len(Rmvls.Window) {
								err := fmt.Sprintf("Room=%s <window> %s", Rm.Name, stt)
								Eprint("<Roomdata>", err)
								os.Exit(1)
							}
						} else {
							// 読み取り中の部位が窓以外の場合
							for j := range Rmvls.Wall {
								w := Rmvls.Wall[j]

								if DEBUG {
									fmt.Printf("!!!!Wall.name=%s  s=%s!!!!\n", get_string_or_null(w.name), s)
								}

								if w.name == s && w.ble == Sd.ble {
									if DEBUG {
										fmt.Printf("---- j=%d Wallname=%s n=%d\n", j, get_string_or_null(w.name), n)
									}

									Sd.wd = j
									Sd.Rwall = w.Rwall
									Sd.CAPwall = w.CAPwall
									Sd.PCMflg = w.PCMflg
									break
								}
							}

							if j == len(Rmvls.Wall) {
								err := fmt.Sprintf("Room=%s <wall> ble=%c %s Undefined in <WALL>", Rm.Name, Sd.ble, s)
								Eprint("<Roomdata>", err)
								os.Exit(1)
							}

						}
					}
				} else {
					key, value := s[:st], s[st+1:]

					if !compnameStart {
						// -- 室への設定 --

						if key == "Vol" {
							// 室容積
							Rm.VRM, err = readRoomVol(value)
							if err != nil {
								panic(err)
							}
						} else if key == "flrsr" {
							// 床の日射吸収比率
							Rm.flrsr = nil
							k, err = idsch(value, Schdl.Sch, "")
							if err == nil {
								Rm.flrsr = &Schdl.Val[k]
							} else {
								Rm.flrsr = envptr(value, Simc, nil, nil, nil)
							}
						} else if key == "alc" {
							// alc 室内表面熱伝達率[W/m2K]。
							k, err = idsch(value, Schdl.Sch, "")
							if err == nil {
								Rm.alc = &Schdl.Val[k]
							} else {
								Rm.alc = envptr(s[st+1:], Simc, nil, nil, nil)
							}
						} else if key == "Hcap" {
							// 室内空気に付加する熱容量 [J/K]
							Rm.Hcap, err = readFloat(value)
							if err != nil {
								panic(err)
							}
						} else if key == "Mxcap" {
							// 室内空気に付加する湿気容量 [kg/(kg/Kg)]
							Rm.Mxcap, err = readFloat(value)
							if err != nil {
								panic(err)
							}
						} else if key == "MCAP" {
							// 室内に置かれた物体の熱容量 [J/K]
							k, err := idsch(value, Schdl.Sch, "")
							if err == nil {
								Rm.MCAP = &Schdl.Val[k]
							} else {
								Rm.MCAP = envptr(value, Simc, nil, nil, nil)
							}
						} else if key == "CM" {
							// 室内に置かれた物体と室内空気との間の熱コンダクタンス [W/K]
							k, err = idsch(value, Schdl.Sch, "")
							if err == nil {
								Rm.CM = &Schdl.Val[k]
							} else {
								Rm.CM = envptr(value, Simc, nil, nil, nil)
							}
						} else if key == "fsolm" { // 家具への日射吸収割合
							k, err = idsch(value, Schdl.Sch, "")
							if err == nil {
								Rm.fsolm = &Schdl.Val[k]
							} else {
								Rm.fsolm = envptr(value, Simc, nil, nil, nil)
							}
						} else if key == "PCMFurn" {
							// PCM内臓家具の場合　(PCMname,mPCM)
							var PCMname, stbuf string
							s1 := s[:st]
							s2 := s[st+2:]
							PCMname = s2
							stbuf = PCMname

							st := strings.IndexRune(PCMname, ',')
							s1 = PCMname[:st]
							s2 = PCMname[st+1:]
							Rm.PCMfurnname = s1

							var err error
							st1 := strings.IndexRune(stbuf, ')')
							st2 := strings.IndexRune(stbuf, ',')
							Rm.mPCM, err = readFloat(stbuf[st2+1 : st1])
							if err != nil {
								panic(err)
							}

							for _, PCM := range Rmvls.PCM {
								if Rm.PCMfurnname == PCM.Name {
									Rm.PCM = PCM
								}
							}
							if Rm.PCM == nil {
								Er = fmt.Sprintf("Roomname=%s %sが見つかりません", Rm.Name, Rm.PCMfurnname)
								Eprint(Er, "<Roomdata>")
								os.Exit(1)
							}
						} else if key == "OTc" {
							// 作用温度設定時の対流成分重み係数の設定
							if k, err = idsch(value, Schdl.Sch, ""); err == nil {
								Rm.OTsetCwgt = &Schdl.Val[k]
							} else {
								Rm.OTsetCwgt = envptr(value, Simc, nil, nil, nil)
							}
						} else {
							Err := fmt.Sprintf("Room=%s s=%s", Rm.Name, s)
							Eprint("<Roomdata>", Err)
						}
					} else {
						// -- 部位設定 --
						//printf ( "st=%s  Sd.name=%s\n", st, Sd.name ) ;

						if strings.HasPrefix(s, "A=") {
							Sd.A, err = readFloat(s[2:])
							if err != nil {
								panic(err)
							}
						} else if strings.HasPrefix(s, "e=") {
							// 外表面の検索
							for j, e := range Exs {
								if e.Name == s[st+1:] {
									Sd.exs = j
									break
								}
							}
							// 見つからない場合
							if j == len(Exs) {
								err := fmt.Sprintf("Room=%s <exsrf> %s\n", Rm.Name, s)
								Eprint("<Roomdata>", err)
								os.Exit(1)
							}
						} else if strings.HasPrefix(s, "sb=") {
							// 日よけの検索
							for j := range Rmvls.Snbk {
								S := Rmvls.Snbk[j]
								if S.Name == s[st+1:] {
									Sd.sb = j
									break
								}
							}
							// 見つからない場合
							if j == len(Rmvls.Snbk) {
								err := fmt.Sprintf("Room=%s <Snbrk> %s\n", Rm.Name, s)
								Eprint("<Roomdata>", err)
								os.Exit(1)
							}
						} else if strings.HasPrefix(s, "r=") {
							// 隣室名
							Sd.nxrmname = s[st+1:]
						} else if strings.HasPrefix(s, "c=") {
							// 隣室温度係数
							Sd.c, err = strconv.ParseFloat(s[st+1:], 64)
							if err != nil {
								panic(err)
							}
						} else if strings.HasPrefix(s, "sw=") {
							// 窓変更設定番号
							Sd.fnsw, err = idscw(s[st+1:], Schdl.Scw, "")
							if err != nil {
								panic(err)
							}
						} else if strings.HasPrefix(s, "i=") {
							// 壁体名
							// 放射暖冷房パネル、部位一体型集熱器のときにSYSPTHでの要素名で使用する。
							Sd.Name = s[st+1:]
						} else if strings.HasPrefix(s, "alc=") {
							if k, err := idsch(s[st+1:], Schdl.Sch, ""); err == nil {
								Sd.alicsch = &Schdl.Val[k]
							} else {
								Sd.alicsch = envptr(s[st+1:], Simc, nil, nil, nil)
							}
						} else if strings.HasPrefix(s, "alr=") {
							if k, err := idsch(s[st+1:], Schdl.Sch, ""); err == nil {
								Sd.alirsch = &Schdl.Val[k]
							} else {
								Sd.alirsch = envptr(s[st+1:], Simc, nil, nil, nil)
							}
						} else if strings.HasPrefix(s, "fsol=") {
							Rm.Nfsolfix++
							Sd.ffix_flg = '*'
							if k, err := idsch(s[st+1:], Schdl.Sch, ""); err == nil {
								Sd.fsol = &Schdl.Val[k]
							} else {
								Sd.fsol = envptr(s[st+1:], Simc, nil, nil, nil)
							}
						} else if strings.HasPrefix(s, "rmp=") {
							// RMP名
							Sd.Sname = s[4:]
						} else if strings.HasPrefix(s, "PVcap=") {
							// 空気式集熱器で太陽電池(PV)付のときのアレイの定格発電量[W]
							Sd.PVwall.PVcap, err = readFloat(s[st+1:])
							if err != nil {
								panic(err)
							}
							Sd.PVwallFlg = true
						} else if strings.HasPrefix(s, "Wsu=") {
							// 集熱屋根の通気層上側の幅 [m]
							Sd.dblWsu, err = readFloat(s[st+1:])
							if err != nil {
								panic(err)
							}
						} else if strings.HasPrefix(s, "Wsd=") {
							// 集熱屋根の通気層下側の幅 [m]
							Sd.dblWsd, err = readFloat(s[st+1:])
							if err != nil {
								panic(err)
							}
						} else if strings.HasPrefix(s, "Ndiv=") {
							// 空気式集熱器のときの流れ方向（入口から出口）の分割数
							Sd.Ndiv, err = strconv.Atoi(s[st+1:])
							if err != nil {
								panic(err)
							}
							Sd.Tc = make([]float64, Sd.Ndiv)
						} else if strings.HasPrefix(s, "tnxt=") {
							// 当該部位への入射日射の隣接空間への日射分配（連続空間の隣室への日射分
							Sd.tnxt, err = readFloat(s[st+1:])
							if err != nil {
								panic(err)
							}
						} else {
							err := fmt.Sprintf("Room=%s ble=%c s=%s\n", Rm.Name, Sd.ble, s)
							Eprint("<Roomdata>", err)
							os.Exit(1)
						}
					}
				}
			}

			Sd.sfepri = sfemark
			Sd.Sname = ""

			Sd.rm = i
			Sd.room = Rm
			Sd.n = nr

			switch Sd.ble {
			case 'E', 'R', 'F', 'W':
				// 外壁, 屋根, 外気に接する床 or 窓の場合
				if Sd.exs == -1 {
					// 外壁の名前から外壁の情報を取得
					for j, e := range Exs {
						if e.Name == dexsname {
							Sd.exs = j
							break
						}
					}

					// 外壁の名前が見つからない場合
					if Sd.exs == -1 {
						err := fmt.Sprintf("Room=%s  (%s)\n --- %s", Rm.Name, dexsname, strings.Join(line, " "))
						Eprint("<Roomdata>", err)
						os.Exit(1)
					}
				}
			case 'i', 'c', 'f':
				// 内壁, 天井(内部) or 床(内部)の場合
				if Sd.nxrm == -1 && Sd.c < 0.0 {
					Sd.nxrmname = dnxrname //隣室名
				}
				if Sd.c < 0.0 {
					// 隣室温度係数 1.0
					Sd.c = 1.0
				}
			}

			// 窓を除く面積0より大きい壁体で、固有の壁体定義がない場合：
			// 既定の壁体定義番号を割り当てる
			if Sd.ble != BLE_Window && Sd.wd == -1 && Sd.A > 0.0 {
				switch Sd.ble {
				case BLE_ExternalWall:
					Sd.wd = dfwl.E // 外壁(壁体定義番号既定値)
				case BLE_Roof:
					Sd.wd = dfwl.R // 屋根(壁体定義番号既定値)
				case BLE_Floor:
					Sd.wd = dfwl.F // 外部に接する床(壁体定義番号既定値)
				case BLE_InnerWall:
					Sd.wd = dfwl.i // 内壁(壁体定義番号既定値)
				case BLE_Ceil:
					Sd.wd = dfwl.c // 天井(内部)(壁体定義番号既定値)
				case BLE_InnerFloor:
					Sd.wd = dfwl.f // 床(内部)(壁体定義番号既定値)
				}
			}

			if Sd.ble == BLE_Window {
				// 窓の場合
				Sd.typ = RMSRFType_W
				Sd.wd = -1
				Sd.tnxt = 0.0
			} else {
				// 窓以外の場合
				j := Sd.wd
				var jj int
				if jj = Sd.exs; jj >= 0 && Exs[jj].Typ == 'E' {
					Sd.typ = RMSRFType_E // 地下
				} else if jj = Sd.exs; jj >= 0 && Exs[jj].Typ == 'e' {
					Sd.typ = RMSRFType_e // 地表面
				} else {
					Sd.typ = RMSRFType_H // 壁
				}

				if j >= 0 {
					w := Rmvls.Wall[j]
					Sd.Eo = w.Eo
					Sd.Ei = w.Ei
					Sd.as = w.as
					Sd.fn = -1
					Sd.Rwall = w.Rwall
					Sd.CAPwall = w.CAPwall
					Sd.PCMflg = w.PCMflg
					if Sd.tnxt < 0.0 {
						Sd.tnxt = w.tnxt
					}
					Sd.tnxt = math.Max(Sd.tnxt, 0.0)
				}
			}

			// 壁体情報を追加
			Rmvls.Sd = append(Rmvls.Sd, Sd)
		}

		nr++
		Rm.N = nr

		N2 = nr * nr
		if Rm.fij != 'F' {
			Rm.F = make([]float64, N2)
		}
		Rm.alr = make([]float64, N2)
		Rm.XA = make([]float64, N2)
		Rm.Wradx = make([]float64, N2)

		Rm.Brs = brs
		Rm.rsrf = Rmvls.Sd[brs:]
		brs += nr

		Rm.GRM = Roa*Rm.VRM + Rm.Mxcap
		Rm.MRM = Ca*Roa*Rm.VRM + Rm.Hcap
	}
	i++
	Nroom := i

	n++
	Nsrf := n

	//printf ( "Nsrf=%d\n", Nsrf ) ;
	//Room = Rmvls.Room

	// 全ての壁体のループ
	for _, Sd := range Rmvls.Sd {
		// 隣室名を解決しておく
		if Sd.nxrmname != "" {
			err := fmt.Sprintf("%s%s", Er, Sd.nxrmname)
			var err2 error
			Sd.nxrm, err2 = idroom(Sd.nxrmname, Rmvls.Room, err)
			if err2 != nil {
				panic(err2)
			}
			Sd.nextroom = Rmvls.Room[Sd.nxrm]
		}
	}

	/******* 個別内壁 *****/

	// 隣室が設定されているの壁体のループ
	for i, Sd := range Rmvls.Sd {
		if Sd.nxrm < 0 {
			continue
		}

		Room := Rmvls.Room[Sd.nxrm]
		brs := Room.Brs
		bre := brs + Room.N

		switch Sd.ble {
		case BLE_InnerWall:
			// 内壁 [i]
			for j := brs; j < bre; j++ {
				Sdj := Rmvls.Sd[j]
				if Sdj.nxrm == Sd.rm && Sdj.ble == BLE_InnerWall {
					Sd.nxn = j
					if i == j {
						panic("")
					}
				}
			}
		case BLE_Ceil:
			// 天井(内部) [f]
			for j := brs; j < bre; j++ {
				Sdj := Rmvls.Sd[j]
				if Sdj.nxrm == Sd.rm && Sdj.ble == BLE_InnerFloor {
					Sd.nxn = j
					if i == j {
						panic("")
					}
				}
			}
		case BLE_InnerFloor:
			// 床(内部) [c]
			for j := brs; j < bre; j++ {
				Sdj := Rmvls.Sd[j]
				if Sdj.nxrm == Sd.rm && Sdj.ble == BLE_Ceil {
					Sd.nxn = j
					if i == j {
						panic("")
					}
				}
			}
		}
	}

	/***** 共用内壁 ******/

	// 共用内壁のループ
	for j, rsd := range Rmvls.Sd {
		if !((rsd.ble == BLE_InnerWall || rsd.ble == BLE_Ceil || rsd.ble == BLE_InnerFloor) && rsd.mwtype != RMSRFMwType_C) {
			continue
		}

		if !(rsd.Name != "" && rsd.wd >= 0 && rsd.A > 0.0) {
			continue
		}

		// 相手を探す
		flag := false
		for i, nxsd := range Rmvls.Sd {
			if !(nxsd.Name != "" && nxsd.A < 0.0) {
				continue
			}

			// 同じ名前で別の壁体
			if rsd.Name == nxsd.Name && rsd != nxsd {
				rsd.room.Ntr++
				nxsd.room.Ntr++

				nxsd.nextroom = rsd.room
				nxsd.nxsd = rsd
				nxsd.A = rsd.A

				nxsd.Ei = rsd.Eo
				nxsd.Eo = rsd.Ei
				nxsd.as = rsd.as
				nxsd.Rwall = rsd.Rwall
				nxsd.CAPwall = rsd.CAPwall

				nxsd.wd = rsd.wd
				nxsd.mwside = RMSRFMwSideType_M
				rsd.mwtype = RMSRFMwType_C
				nxsd.mwtype = RMSRFMwType_C
				nxsd.pcmpri = rsd.pcmpri
				nxsd.PCMflg = rsd.PCMflg

				nxsd.tnxt = rsd.tnxt

				rsd.nextroom = nxsd.room
				rsd.nxsd = nxsd

				if rsd.ble == BLE_InnerWall {
					nxsd.ble = BLE_InnerWall
				} else if rsd.ble == BLE_InnerFloor {
					nxsd.ble = BLE_Ceil
				} else if rsd.ble == BLE_Ceil {
					nxsd.ble = BLE_InnerFloor
				}

				var err error
				rsd.nxrm, err = idroom(rsd.nextroom.Name, Rmvls.Room, "")
				if err != nil {
					panic(err)
				}
				rsd.nxn = i
				nxsd.nxrm, err = idroom(nxsd.nextroom.Name, Rmvls.Room, "")
				if err != nil {
					panic(err)
				}
				nxsd.nxn = j

				flag = true
				break
			}
		}

		if !flag {
			fmt.Printf("name=%s 共用内壁が片側しか定義されていません。\n", rsd.Name)
		}

		if rsd.nxn < 0 && rsd.mwtype == RMSRFMwType_C {
			err := fmt.Sprintf("%s    room=%s  xxx  (%s):  -%c\n", Er, Rmvls.Room[rsd.rm].Name, Rmvls.Room[rsd.nxrm].Name, rsd.ble)
			Eprint("<Roomdata>", err)
			os.Exit(1)
		}
	}

	// 面積入力のチェック
	for _, rsd := range Rmvls.Sd {
		if rsd.A <= 0.0 {
			fmt.Printf("Room=%s  ble=%c  A=%f\n", rsd.room.Name, rsd.ble, rsd.A)
			os.Exit(1)
		}
	}

	/***** 放射パネル総数、室ごとのパネル数 *****/

	var Nairflow, Nrdpnl int
	for _, rsd := range Rmvls.Sd {
		if rsd.ble != BLE_Window {
			w := Rmvls.Wall[rsd.wd]
			if w.Ip >= 0 {
				rsd.room.Nrp++

				if rsd.mwside == RMSRFMwSideType_i {
					Nrdpnl++
				}
			}
		} else {
			// エアフローウィンドウの総数を数える
			Nairflow++
		}
	}

	for _, room := range Rmvls.Room {
		N := room.Ntr
		if N > 0 {
			room.trnx = make([]*TRNX, N)
		}

		if room.trnx != nil {
			for sk := 0; sk < N; sk++ {
				Tn := new(TRNX)
				Tn.nextroom = nil
				Tn.sd = nil

				room.trnx[sk] = Tn
			}
		}

		room.ARN = make([]float64, room.Ntr)

		N = room.Nrp
		if N > 0 {
			room.rmpnl = make([]*RPANEL, N)
		}

		if room.rmpnl != nil {
			for sk := 0; sk < N; sk++ {
				Rp := new(RPANEL)
				Rp.pnl = nil
				Rp.sd = nil
				Rp.elinpnl = 0

				room.rmpnl[sk] = Rp
			}
		}

		room.RMP = make([]float64, room.Nrp)
	}

	if Nrdpnl > 0 {
		Rmvls.Rdpnl = make([]*RDPNL, Nrdpnl)
	}

	if Rmvls.Rdpnl != nil {

		for sk := 0; sk < Nrdpnl; sk++ {
			Rd := new(RDPNL)
			Rd.Name = ""
			Rd.cmp = nil
			Rd.MC = 0
			Rd.eprmnx = 0
			Rd.epwtw = 0
			Rd.Loadt = nil
			Rd.Toset = FNAN
			Rd.cG = 0.0
			Rd.Ec = 0.0
			Rd.OMvav = nil
			MtEdayinit(&Rd.mtPVdy)

			for si := 0; si < 2; si++ {
				Rd.rm[si] = nil
				Rd.sd[si] = nil
				Rd.Ntrm[si] = 0.0
				Rd.Nrp[si] = 0.0
				Rd.elinpnl[si] = 0
			}

			Rmvls.Rdpnl[sk] = Rd
		}
	}

	for i := 0; i < Nsrf; i++ {
		rsd := Rmvls.Sd[i]
		rsd.WSRN = make([]float64, rsd.room.Ntr)
		rsd.WSPL = make([]float64, rsd.room.Nrp)
	}

	rdpnlIdx := 0
	for i := 0; i < Nroom; i++ {
		room := Rmvls.Room[i]
		room.Nisidermpnl = 0

		trnxIdx := 0
		rmpnlIdx := 0
		for n := 0; n < room.N; n++ {
			rsd := room.rsrf[n]

			// 共用壁の場合
			if rsd.mwtype == RMSRFMwType_C {
				trnx := room.trnx[trnxIdx]
				trnx.nextroom = rsd.nextroom
				trnx.sd = rsd
				trnxIdx++
			}

			if rsd.ble != BLE_Window {
				w := Rmvls.Wall[rsd.wd]
				if w.Ip >= 0 {
					if rsd.mwside == 'i' {
						rdpnl := Rmvls.Rdpnl[rdpnlIdx]

						if w.tra > 0. {
							rdpnl.Type = 'C'
						} else {
							rdpnl.Type = 'P'
						}

						rdpnl.Name = rsd.Name
						rdpnl.effpnl = w.effpnl
						rdpnl.MC = 1

						rdpnl.rm[0] = rsd.room
						rdpnl.sd[0] = rsd
						rdpnl.Ntrm[0] = rsd.room.Ntr
						rdpnl.Nrp[0] = rsd.room.Nrp

						rmpnl := room.rmpnl[rmpnlIdx]
						rmpnl.pnl = rdpnl
						rmpnl.sd = rsd

						rdpnl.elinpnl[0] = 1 + 1 + rdpnl.Ntrm[0]
						rmpnl.elinpnl = rdpnl.elinpnl[0]
						rmpnlIdx++
						room.Nisidermpnl++

						// 共用壁の場合
						if rsd.mwtype == RMSRFMwType_C {
							rdpnl.MC = 2
							nxsd := rsd.nxsd

							rdpnl.rm[1] = nxsd.room
							rdpnl.sd[1] = nxsd
							rdpnl.Ntrm[1] = nxsd.room.Ntr
							rdpnl.Nrp[1] = nxsd.room.Nrp
							rdpnl.elinpnl[1] = 1 + 1 + rdpnl.Ntrm[0] + rdpnl.Nrp[0] + 1 + rdpnl.Ntrm[1]
						}

						for j := 0; j < rdpnl.MC; j++ {
							rdpnl.EPR[j] = make([]float64, rdpnl.Ntrm[j])
							rdpnl.EPW[j] = make([]float64, rdpnl.Nrp[j])
						}

						rdpnlIdx++
					}
				}
			}
		}
	}

	for i := 0; i < Nroom; i++ {
		room := Rmvls.Room[i]
		rmpnlIdx := room.Nisidermpnl
		for n := 0; n < room.N; n++ {
			rsd := room.rsrf[n]

			if rsd.ble != BLE_Window {
				w := Rmvls.Wall[rsd.wd]
				if w.Ip > 0 && rsd.mwside == 'M' {
					rsd.rpnl = rsd.nxsd.rpnl

					rmpnl := room.rmpnl[rmpnlIdx]
					rmpnl.pnl = rsd.rpnl
					rmpnl.sd = rsd
					rmpnl.elinpnl = rsd.rpnl.elinpnl[1]
					rmpnlIdx++
				}
			}
		}
	}

	for i := 0; i < Nroom; i++ {
		Rm := Rmvls.Room[i]

		// 室間相互換気
		Rm.achr = make([]*ACHIR, 0, Nroom)
		for sk := range Rm.achr {
			Ac := Rm.achr[sk]
			Ac.rm = 0
			Ac.sch = 0
			Ac.room = nil
		}
		Rm.Nachr = len(Rm.achr)

		Rm.Arsp = nil
		Rm.rmld = nil
		Area := 0.0
		Rm.FArea = 0.0

		for j := 0; j < Rm.N; j++ {
			rsd := Rm.rsrf[j]
			Area += rsd.A
			if rsd.ble == BLE_Floor || rsd.ble == BLE_InnerFloor {
				Rm.Nflr++
				Rm.FArea += rsd.A
			}
		}
		Rm.Area = Area
		if Rm.fij != 'F' {
			Rm.fij = 'A'
			// 形態係数の近似計算（面積割）
			formfaprx(Rm.N, Area, Rmvls.Sd[Rm.Brs:], Rm.F)
		}
	}

	Rmvls.Trdav = make([]float64, len(Rmvls.Room))

	if len(Rmvls.Room) > 0 {
		N := len(Rmvls.Room)
		Rmvls.Qrm = make([]*QRM, N)
		Rmvls.Qrmd = make([]*QRM, N)
		Rmvls.Emrk = make([]rune, N)

		for i := 0; i < N; i++ {
			Rmvls.Qrm[i] = new(QRM)
			Rmvls.Qrmd[i] = new(QRM)
		}
	}
}

func readFloat(value string) (float64, error) {
	return strconv.ParseFloat(strings.Trim(value, "\""), 64)
}

func readRoomVol(value string) (float64, error) {
	// 室容積 [m3]入力室が直方体の場合には間口、奥行き、高さを'*'でつなげると、
	ast := strings.Split(value, "*")
	if len(ast) == 1 {
		return readFloat(ast[0])
	} else {
		// EESLISM内部で室容積を計算する。
		// Read Wi
		Wi, err := readFloat(ast[0])
		if err != nil {
			return 0.0, err
		}

		// Read H
		H, err := readFloat(ast[1])
		if err != nil {
			return 0.0, err
		}

		// Read D
		D, err := readFloat(ast[2])
		if err != nil {
			return 0.0, err
		}

		return Wi * H * D, nil
	}
}

/* ------------------------------------------------------------- */

/*  重量壁体の計算準備      */

func Balloc(Sd []*RMSRF, Wall []*WALL, Mwall *[]*MWALL) {
	var mw int
	for _, ssd := range Sd {
		if id := ssd.wd; id >= 0 && ssd.mwside == 'i' {
			mw++
		}
	}

	if mw > 0 {
		*Mwall = make([]*MWALL, mw)

		for n := 0; n < mw; n++ {
			(*Mwall)[n] = &MWALL{
				sd:   nil,
				nxsd: nil,
				wall: nil,
				ns:   0,
				rm:   0,
				n:    0,
				nxrm: 0,
				nxn:  0,
				M:    0,
				mp:   0,
				UX:   nil,
				res:  nil,
				cap:  nil,
				Tw:   nil,
				Told: nil,
				uo:   0.0,
				um:   0.0,
				Pc:   0.0,
			}
		}
	}

	mw = 0
	for n, ssd := range Sd {

		if id := ssd.wd; id >= 0 && ssd.mwside == 'i' {
			ssd.rmw = mw
			mwl := (*Mwall)[mw]
			W := Wall[id]
			ssd.mw = mwl // 壁体構造体のポインタ

			mwl.wall = W

			// 太陽光発電付のチェック
			sn := 0
			if ssd.mw.wall.ColType != "" {
				sn = len(ssd.mw.wall.ColType)
			}
			if sn == 2 || sn == 3 && ssd.mw.wall.ColType[2] != 'P' {
				ssd.PVwallFlg = false

				// 太陽電池の容量が入力されているときにはエラーを表示する
				if ssd.PVwall.PVcap > 0.0 {
					fmt.Printf("<%s> name=%s PVcap=%g ですが、WALLで太陽電池付が指定されていません\n",
						ssd.room.Name, ssd.Name, ssd.PVwall.PVcap)
					ssd.PVwall.PVcap = FNAN
					os.Exit(1)
				}
			}

			mwl.sd = ssd
			mwl.nxsd = ssd.nxsd
			mwl.ns = n
			mwl.rm = ssd.rm
			mwl.n = ssd.n
			mwl.nxrm = ssd.nxrm
			mwl.nxn = ssd.nxn
			mwl.M = W.M
			mwl.mp = W.mp

			M := mwl.M

			if mwl.res == nil {
				mwl.res = make([]float64, M+2)
			}

			if mwl.cap == nil {
				mwl.cap = make([]float64, M+2)
			}

			wres := W.res
			wcap := W.cap
			res := mwl.res
			cap := mwl.cap
			for m := 0; m <= M; m++ {
				res[m] = wres[m]
				cap[m] = wcap[m]
			}

			if ssd.typ == 'H' {
				M++
				mwl.M = M
				mwl.res[M] = 0.0
				mwl.cap[M] = 0.0
			}

			mwl.UX = make([]float64, M*M)

			// PCM状態値を保持する構造体
			ssd.pcmstate = make([]*PCMSTATE, M+1)
			pcmstate := ssd.pcmstate
			for m := 0; m <= M; m++ {
				PCM := mwl.wall.PCMLyr[m]
				pcmstate[m] = &PCMSTATE{
					Name:         nil,
					CapmL:        0.0,
					CapmR:        0.0,
					LamdaL:       0.0,
					LamdaR:       0.0,
					TempPCMave:   0.0,
					OldCapmL:     0.0,
					OldCapmR:     0.0,
					OldLamdaL:    0.0,
					OldLamdaR:    0.0,
					TempPCMNodeL: 0.0,
					TempPCMNodeR: 0.0,
				}
				if PCM != nil {
					pcmstate[m].Name = &PCM.Name
					ssd.Npcm++
					if ssd.wlpri {
						ssd.pcmpri = true
					}
				}
			}

			// prevLayer, startLayer := INAN, INAN
			// k := 0
			mw++
		} else {
			ssd.mw = nil
		}
	}

	for _, ssd := range Sd {
		if ssd.mwside == 'M' {
			ssd.mw = ssd.nxsd.mw
			M := ssd.mw.M
			ssd.rmw = ssd.nxsd.rmw

			ssd.PCMflg = ssd.nxsd.PCMflg
			ssd.pcmpri = ssd.nxsd.pcmpri
			ssd.Npcm = ssd.nxsd.Npcm

			// PCM状態値を保持する構造体
			// pcmstate := ssd.pcmstate
			// nxpcm2 := ssd.nxsd.pcmstate
			for m := 0; m <= M; m++ {
				// TODO: ここおかしい？
				//pcmstate[m] = nxpcm2[M-m]
			}
		}
	}

}

/* ------------------------------------------ */

// 壁体内部温度の初期値設定
func (Rmvls *RMVLS) Tinit() {
	Tini := Rmvls.Twallinit

	for _, rm := range Rmvls.Room {
		rm.Tr = Tini
		rm.Trold = Tini
		rm.Tsav = Tini
		rm.Tot = Tini
		rm.xrold = FNXtr(rm.Tr, 50.0)
		rm.xr = rm.xrold
		rm.hr = FNH(rm.Tr, rm.xr)
		rm.alrbold = FNAN
		rm.mrk = '*'
		rm.oldTM = Tini
		rm.TM = rm.oldTM
	}

	for _, Sd := range Rmvls.Sd {
		Sd.Ts = Tini
		Sd.mrk = '*'
	}

	for _, mw := range Rmvls.Mw {
		mw.Tw = make([]float64, mw.M)
		mw.Told = make([]float64, mw.M)
		mw.Toldd = make([]float64, mw.M)
		mw.Twd = make([]float64, mw.M)

		for m := 0; m < mw.M; m++ {
			mw.Tw[m] = Tini
			mw.Told[m] = Tini
			mw.Toldd[m] = Tini
			mw.Twd[m] = Tini
		}
	}

	for _, Room := range Rmvls.Room {
		if Room.rmqe == nil {
			continue
		}
		for j := 0; j < Room.N; j++ {
			rmsb := Room.rmqe.rmsb[j]
			Sd := Room.rsrf[j]
			if mw := Sd.mw; mw != nil {
				for m := 0; m < mw.M; m++ {
					Told := rmsb.Told[m]
					Tw := rmsb.Tw[m]

					helmclear(Told)
					Told.trs = Tini
					helmcpy(Told, Tw)
				}
			}
		}
	}
}

/********************************************************************/

func Roomcount(tokens *EeTokens) int {
	N := 0
	pos := tokens.GetPos()

	// Find empty section
	for tokens.IsEnd() == false {
		section := tokens.GetSection()
		s := section.GetToken()
		if s != "*" {
			N++
		} else {
			break
		}
	}

	// restore position
	tokens.RestorePos(pos)

	return N
}

/********************************************************************/

func Roominit(N int, Room []*ROOM) {
	for i := 0; i < N; i++ {
		B := new(ROOM)

		B.Name = ""
		B.PCM = nil
		B.PCMfurnname = ""
		B.mPCM = FNAN
		B.FunHcap = FNAN
		B.PCMQl = FNAN
		B.N = 0
		B.Brs = 0
		B.Nachr = 0
		B.Ntr = 0
		B.Nrp = 0
		B.Nflr = 0
		B.Nfsolfix = 0
		B.Nisidermpnl = 0
		B.Nasup = 0
		B.Brs = 0
		B.N = 0
		//B.Nairflow = 0 ;
		B.rsrf = nil
		B.achr = nil
		B.trnx = nil
		B.rmpnl = nil
		B.Arsp = nil
		B.cmp = nil
		B.elinasup = nil
		B.elinasupx = nil
		B.rmld = nil
		B.rmqe = nil
		B.F = nil
		B.alr = nil
		B.XA = nil
		B.Wradx = nil
		B.rsrnx = false
		B.fij = ' '
		B.sfpri = false
		B.eqpri = false
		B.mrk = ' '
		B.VRM = 0.0
		B.GRM = 0.0
		B.MRM = 0.0
		B.Area = 0.0
		B.FArea = 0.0
		B.flrsr = CreateConstantValuePointer(0.3)
		B.tfsol = 0.0
		B.alrbold = 0.0
		B.Hcap = 0.0
		B.Mxcap = 0.0
		B.Ltyp = ' '
		B.Nhm = 0.0
		B.Light = 0.0
		B.Apsc = 0.0
		B.Apsr = 0.0
		B.Apl = 0.0
		B.Gve = 0.0
		B.Gvi = 0.0
		B.alc = nil
		B.Vesc = nil
		B.Visc = nil
		//B.vesc = B.visc = 0 ;
		// B.hmwksc = B.hmnsc = B.lgtsc = B.apssc = B.aplsc = 0 ;
		//B.metsc = B.closc = B.wvsc = -1 ;
		B.Hc = 0.0
		B.Hr = 0.0
		B.HL = 0.0
		B.Lc = 0.0
		B.Lr = 0.0
		B.Ac = 0.0
		B.Ar = 0.0
		B.AL = 0.0
		B.eqcv = 0.5
		B.Qeqp = 0.0
		B.Gvent = 0.0
		B.RMt = 0.0
		B.ARN = nil
		B.RMP = nil
		B.RMC = 0.0
		B.RMx = 0.0
		B.RMXC = 0.0
		B.Tr = 0.0
		B.Trold = 0.0
		B.xr = 0.0
		B.xrold = 0.0
		B.RH = 0.0
		B.Tsav = 0.0
		B.Tot = 0.0
		B.PMV = 0.0
		B.AEsch = nil
		B.AGsch = nil
		B.AE = 0.0
		B.AG = 0.0
		B.Assch = nil
		B.Alsch = nil
		B.Lightsch = nil
		B.Hmsch = nil
		B.Metsch = nil
		B.Closch = nil
		B.Wvsch = nil
		B.Hmwksch = nil
		B.VAVcontrl = nil
		B.OTsetCwgt = nil // 作用温度設定時の対流成分重み係数
		//B.rairflow = nil ;
		B.MCAP = nil
		B.CM = nil
		B.QM = 0.0
		B.HM = 0.0
		B.fsolm = nil
		B.Srgm2 = 0.0
		B.TM = 15.0
		B.oldTM = 15.0
		B.SET = FNAN
		B.setpri = false

		Room[i] = B
	}
}

/************************************************************************/

func Rmsrfcount(tokens *EeTokens) int {
	N := 0

	//save current position
	pos := tokens.GetPos()

	for tokens.IsEnd() == false {
		s := tokens.GetToken()
		if strings.HasPrefix(s, "-") {
			N++
		}
	}

	// restore postion
	tokens.RestorePos(pos)

	return N
}

/************************************************************************/

func Rmsrfinit() *RMSRF {
	S := new(RMSRF)
	S.Ctlif = nil
	S.ifwin = nil
	S.Name = ""
	S.room = nil
	S.nextroom = nil
	S.DynamicCode = ""
	S.nxsd = nil
	S.mw = nil
	S.rpnl = nil
	S.pcmstate = nil
	S.Npcm = 0
	S.Nfn = 0
	S.pcmpri = false
	S.Rwall = FNAN
	S.CAPwall = FNAN
	S.A = 0.0
	S.Eo = 0.0
	S.as = 0.0
	S.c = 0.0
	S.tgtn = 0.0
	S.Bn = 0.0
	S.srg = 0.0
	S.srh = 0.0
	S.srl = 0.0
	S.sra = 0.0
	S.alo = 0.0
	S.ali = 0.0
	S.alic = 0.0
	S.alir = 0.0
	S.K = 0.0
	S.FI = 0.0
	S.FO = 0.0
	S.FP = 0.0
	S.CF = 0.0
	S.WSR = 0.0
	S.WSC = 0.0
	S.RS = 0.0
	S.RSsol = 0.0
	S.RSin = 0.0
	S.RSli = 0.0
	S.Qi = 0.0
	S.Qga = 0
	S.Qgt = 0.
	S.TeEsol = 0.0
	S.TeErn = 0.0
	S.Te = 0.0
	S.Tmrt = 0.0
	S.Ei = 0.0
	S.Ts = 0.0
	S.eqrd = 0.0
	S.alicsch = nil
	S.WSRN = nil
	S.WSPL = nil

	S.exs = -1
	S.sb = -1
	S.nxrm = -1
	S.nxn = -1
	S.wd = -1
	S.fn = -1
	S.c = -1.0
	S.A = FNAN
	//		S.Rwall = 0.0 ;
	S.mwside = RMSRFMwSideType_i
	S.mwtype = RMSRFMwType_I
	S.fnmrk = [10]rune{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '}
	S.alirsch = nil
	S.ffix_flg = '!'
	S.fsol = nil

	S.ColCoeff = FNAN
	S.oldTx = 20.0
	S.Iw = 0.0
	//S.Scol = 0.0 ;
	S.PVwall.Eff = 0.0
	S.PVwallFlg = false
	S.PVwall.PVcap = FNAN
	S.Ndiv = 0
	S.Tc = nil
	S.dblWsd = FNAN
	S.dblWsu = FNAN
	S.dblTf = 20.0
	S.dblTsd = 20.0
	S.dblTsu = 20.0
	S.ras = FNAN
	S.Tg = 20.

	S.tnxt = FNAN
	S.RStrans = false

	S.wlpri = false
	S.shdpri = false
	S.Iwall = 0.0
	S.fnsw = 0

	for j := 0; j < 10; j++ {
		f := &S.direct_heat_gain[j]
		g := &S.fnd[j]
		*f = 0
		*g = 0
	}

	return S
}

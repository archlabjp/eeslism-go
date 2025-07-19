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

/*   binit.c   */
package eeslism

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

/*
Walldata (Wall Data Input)

この関数は、建物の壁体（外壁、屋根、床、内壁など）の仕様データを読み込み、
対応する構造体に格納します。
これには、壁体の材料構成、熱的特性、および集熱器や放射パネルなどの特殊な機能に関する情報が含まれます。

建築環境工学的な観点:
  - **壁体の熱的特性の定義**: 建物の外皮は、外部環境と室内環境を隔てる重要な要素であり、
    その熱的特性は建物のエネルギー消費量や室内快適性に大きく影響します。
    この関数は、壁体を構成する各層の材料コード（`Mcode`）、熱伝導率（`Cond`）、
    容積比熱（`Cro`）などを読み込み、壁体全体の熱貫流率や熱容量を計算するための基礎データを提供します。
  - **部位ごとの特性設定**: `Wa.ble`（部位コード）によって、
    外壁、屋根、床、内壁などの部位を識別し、
    それぞれの部位に応じた熱的特性やデフォルト値を設定します。
    これにより、建物の各部位の熱的挙動を正確にモデル化できます。
  - **特殊な壁体のモデル化**:
  - **集熱器一体型壁 (Wa.ColType)**: 太陽熱集熱器や太陽光発電パネルが一体となった壁体の場合、
    その熱的特性（透過率`tra`、熱貫流率`Ksu`, `Ksd`、放射率`dblEsu`, `dblEsd`など）を読み込みます。
    これにより、パッシブソーラーシステムや再生可能エネルギー利用の効果を評価できます。
  - **放射パネル内蔵壁 (Wa.Ip)**: 床暖房などの放射パネルが内蔵された壁体の場合、
    その位置（`Wa.Ip`）や効率（`Wa.effpnl`）を読み込みます。
    これにより、放射冷暖房システムの効果をモデル化できます。
  - **PCM内蔵壁 (Wa.PCM)**: 相変化材料（PCM）が内蔵された壁体の場合、
    その特性を読み込み、蓄熱効果をモデル化します。
  - **層構成の定義**: 壁体を構成する各層の厚さ（`L`）や内部の分割数（`ND`）を読み込みます。
    これにより、壁体内部の温度分布や熱流を詳細にモデル化できます。

この関数は、建物の外皮の熱的挙動を詳細にモデル化し、
熱負荷計算、エネルギー消費量予測、
省エネルギー対策の検討、および快適性評価を行うための重要なデータ入力機能を提供します。
*/
func Walldata(section *EeTokens, fbmlist string, Wall *[]*WALL, dfwl *DFWL, pcm []*PCM) {
	var s string
	var i = -1
	var j, jj, jw, Nlyr, k = 0, 0, 0, 0, -1
	var Nbm, ndiv int // ic
	var dt float64
	var W []BMLST
	var Wl *BMLST

	W = nil
	// E = fmt.Sprintf(ERRFMT, dsn)

	NwLines := wbmlistcount(fbmlist)

	s = "Walldata wbmlist.efl--"

	k = 0

	for _, line := range NwLines {
		s := strings.Fields(line)

		Wl = new(BMLST)
		Wl.Mcode = s[0]

		for j := 0; j < k-1; j++ {
			Wc := &W[j]
			if Wl.Mcode == Wc.Mcode {
				message := fmt.Sprintf("wbmlist.efl duplicate code=<%s>", Wl.Mcode)
				Eprint("<Walldata>", message)
			}
		}

		// Cond：熱伝導率［W/mK］、Cro：容積比熱［kJ/m3K］
		var err error
		Wl.Cond, err = strconv.ParseFloat(s[1], 64)
		if err != nil {
			Eprint("<Walldata>", "wbmlist.efl Cond error")
		}

		Wl.Cro, err = strconv.ParseFloat(s[2], 64)
		if err != nil {
			Eprint("<Walldata>", "wbmlist.efl Cro error")
		}

		W = append(W, *Wl)

		for j := 0; j < k+1; j++ {
			fmt.Printf("k=%d code=%s\n", k, W[j].Mcode)
		}

		fmt.Printf("Walldata>> name=%s Con=%f Cro=%f\n", W[k].Mcode, W[k].Cond, W[k].Cro)
		W[k] = *Wl
		k++
	}

	k++
	Nbm = k

	//N = Wallcount(fi)

	s = "Walldata --"

	i = 0
	for section.IsEnd() == false {
		line := section.GetLogicalLine()
		if line[0] == "*" {
			break
		}

		Wa := NewWall()

		s = line[0]

		// for jj = 0; jj < len(Wa.welm); jj++ {
		// 	w = NewWelm()
		// }

		// for jj = 0; jj < len(Wa.PCM); jj++ {
		// 	Wa.PCM[jj] = nil
		// 	Wa.PCMrate[jj] = NAN
		// }

		// (1) 部位・壁体名の読み取り
		if strings.HasPrefix(s, "-") {
			// For `-E` or `-R:ROOF` or ...

			// 部位コードの指定
			Wa.ble = BLEType(s[1])

			if len(s) > 2 && s[2] == ':' {
				// 壁体名の指定(最初の1文字を必ず英字とする)
				Wa.name = s[3:]
			} else {
				// 部位のみ指定（既定値の定義）
				switch s[1] {
				case 'E': // 外壁
					dfwl.E = i
				case 'R': // 屋根
					dfwl.R = i
				case 'F': // 床(外部)
					dfwl.F = i
				case 'i': // 内壁
					dfwl.i = i
				case 'c': // 天井(内部)
					dfwl.c = i
				case 'f': // 床(内部)
					dfwl.f = i
				}
			}
		} else {
			// 部位コードの指定なし
			Wa.name = s
		}

		j = -1

		// (2) 部位・壁体のパラメータを読み取り
		// 例) `Eo=0.9 Ei=0.9 as=0.7 type=1 APR-100 APR-100/20 <P:80.3> ;`
		var layer []string
		for _, s = range line[1 : len(line)-1] {
			var err error
			st := strings.IndexRune(s, '=')
			if st != -1 {
				dt, err = strconv.ParseFloat(s[st+1:], 64)
				if err != nil {
					panic(err)
				}

				switch strings.TrimSpace(s[:st]) {
				case "Ei":
					Wa.Ei = dt // 室内表面放射率
				case "Eo":
					Wa.Eo = dt // 外表面放射率
				case "as":
					Wa.as = dt // 外表面日射吸収率
				case "type":
					Wa.ColType = s[st+1:] // 集熱器のタイプ
				case "tra":
					Wa.tra = dt // τα
				case "Ksu":
					Wa.Ksu = dt // 通気層内上側から屋外までの熱貫流率 [W/m2K]
				case "Ksd":
					Wa.Ksd = dt // 通気層内下側から裏面までの熱貫流率 [W/m2K]
				case "Ru":
					Wa.Ru = dt // 通気層から上面までの熱抵抗 [m2K/W]
				case "Rd":
					Wa.Rd = dt // 通気層から裏面までの熱抵抗 [m2K/W]
				case "fcu":
					Wa.fcu = dt // Kcu / Ksu （太陽光発電設置時のアレイ温度計算のみに使用）
				case "fcd":
					Wa.fcd = dt // Kcd / Ksd （太陽光発電設置時のアレイ温度計算のみに使用）
				case "Kc":
					Wa.Kc = dt
				case "Esu":
					Wa.dblEsu = dt // 通気層内上側の放射率
				case "Esd":
					Wa.dblEsd = dt // 通気層内下側の放射率
				case "Eg":
					Wa.Eg = dt // 透過体の中空層側表面の放射率
				case "Eb":
					Wa.Eb = dt // 集熱板の中空層側表面の放射率
				case "ag":
					Wa.ag = dt // 透過体の日射吸収率
				case "ta":
					Wa.ta = dt / 1000.0 // 中空層の厚さ [mm] -> [m]
				case "tnxt":
					Wa.tnxt = dt
				case "t":
					Wa.air_layer_t = dt / 1000.0
				case "KHD":
					Wa.PVwallcat.KHD = dt // 日射量年変動補正係数 (集熱板が太陽電池一体型のとき)
				case "KPD":
					Wa.PVwallcat.KPD = dt // 経時変化補正係数 (集熱板が太陽電池一体型のとき)
				case "KPM":
					Wa.PVwallcat.KPM = dt // アレイ負荷整合補正係数 (集熱板が太陽電池一体型のとき)
				case "KPA":
					Wa.PVwallcat.KPA = dt // アレイ回路補正係数 (集熱板が太陽電池一体型のとき)
				case "EffInv":
					Wa.PVwallcat.EffINO = dt // インバータ効率 (集熱板が太陽電池一体型のとき)
				case "apamax":
					Wa.PVwallcat.Apmax = dt // 最大出力温度係数 (集熱板が太陽電池一体型のとき)
				case "ap":
					Wa.PVwallcat.Ap = dt // 太陽電池裏面の対流熱伝達率
				case "Rcoloff":
					Wa.PVwallcat.Rcoloff = dt // 太陽電池から集熱器裏面までの熱抵抗 (集熱板が太陽電池一体型のとき)
					Wa.PVwallcat.Kcoloff = 1. / Wa.PVwallcat.Rcoloff
				default:
					Eprint("<Walldata>", s)
				}
			} else {
				layer = append(layer, s)
				j++
			}
		}

		// 層の数
		Nlyr = j + 1

		Wa.welm = []WELM{
			{
				Code: "ali", //内表面熱伝達率
				L:    FNAN,
				ND:   0,
				Cond: FNAN,
				Cro:  FNAN,
			},
		}

		// (3) 層の読み取り
		// `APR-100 APR-100/20 <P:80.3>`
		for j = 1; j <= Nlyr; j++ {
			if Wa.ble == BLE_Roof || Wa.ble == BLE_Ceil {
				jj = Nlyr - j
			} else {
				jj = j - 1
			}

			// `APR-100`, `APR-100/20` or `<P>:80.3` or `<C>`
			var err error
			st := strings.IndexRune(layer[jj], '-')
			if st != -1 {
				// 一般材料のとき または 一般材料で分割層数を指定するとき
				// For `APR-100` or `APR-100/20`
				var code string // 材料コード
				var ND int      // 内部分割数
				var L float64   // 部材厚さ [mm]

				// 1) 材料コード
				code = layer[jj][:st]

				// 2) 同一層内の内部分割数と部材の厚さ
				ss := layer[jj][st+1:]
				st = strings.IndexRune(ss, '/')
				if st != -1 {
					// "20.0/2"のように分割数が指定されている場合
					sss := strings.SplitN(ss, "/", 2)
					L, err = strconv.ParseFloat(sss[0], 64)
					if err != nil {
						panic(err)
					}

					// 分割数
					ndiv, err = strconv.Atoi(sss[1])
					if err != nil {
						panic(err)
					}
					ND = ndiv - 1
				} else {
					// "20.0"のように分割数が指定されていない場合
					L, err = strconv.ParseFloat(ss, 64)
					if err != nil {
						panic(err)
					}

					if L >= 50. {
						ND = int((L - 50.) / 50.)
					} else {
						// 50mm未満の場合は分割しない
						ND = 0
					}
				}

				Wa.welm = append(Wa.welm, WELM{
					Code: code,
					ND:   ND,
					L:    L / 1000.0, // [mm] -> [m]
				})

			} else if strings.HasPrefix(layer[jj], "<P") || layer[jj] == "<C>" {
				// 放射暖冷房パネル発熱面位置の指定の場合 For `<P:80.3>` or `<P>` or `<C>`
				// ※ `<C>`はマニュアルに記述が見当たらないが、`<P>`と同じ扱いにする
				Wa.Ip = len(Wa.welm) - 1
				if layer[jj][2] == ':' {
					// パネル効率が指定されている場合 `<P:80.3>`
					Wa.effpnl, err = strconv.ParseFloat(layer[jj][3:], 64)
					if err != nil {
						panic(err)
					}
				} else {
					// パネル効率が指定されない場合 `<C>`
					Wa.effpnl = 0.7
				}
			} else if layer[jj] != "alo" && layer[jj] != "ali" {
				// 表面熱伝達率、中空層熱コンダクタンスのとき、材料コードのみを指定する
				Wa.welm = append(Wa.welm, WELM{
					Code: layer[jj],
					ND:   0,
					L:    0.0,
				})
			}
		}
		jw = len(Wa.welm) - 1

		// 建材一体型空気集熱器の総合熱貫流率の計算、データ入力状況のチェック
		if Wa.Ip >= 0 {
			if Wa.tra > 0. {
				// 壁種類 -> 建材一体型空気集熱器
				Wa.WallType = 'C'

				if (Wa.Ksu > 0. && Wa.Ksd > 0.) || (Wa.Rd > 0. && (Wa.Ru >= 0. || Wa.ta > 0.)) {

					if strings.HasPrefix(Wa.ColType, "A") {
						// 集熱器のタイプ = A1,A2,A2P or A3
						if Wa.Ksu > 0. {
							// 熱抵抗が入力されていない
							Wa.Kcu = Wa.fcu * Wa.Ksu
							Wa.Kcd = Wa.fcd * Wa.Ksd
							Wa.Kc = Wa.Kcu + Wa.Kcd
							Wa.ku = Wa.Kcu / Wa.Kc
							Wa.kd = Wa.Kcd / Wa.Kc
						} else {
							// 熱抵抗が入力されている
							Wa.chrRinput = true
						}
					} else if strings.HasPrefix(Wa.ColType, "W") {
						// 集熱器のタイプ = W3
						Wa.Ko = Wa.Ksu + Wa.Ksd
						Wa.ku = Wa.Ksu / Wa.Ko
						Wa.kd = Wa.Ksd / Wa.Ko
					}

					// PVがある場合は事前計算をする
					if Wa.PVwallcat.KHD > 0. {
						PVwallPreCalc(&Wa.PVwallcat)
					}
				} else {
					s = fmt.Sprintf("ble=%c name=%s 建築一体型空気集熱の熱貫流率Ku、Kdが未定義です", Wa.ble, Wa.name)
					Eprint("<Walldata>", s)
				}

				if Wa.chrRinput == false && (Wa.Kc < 0. || Wa.Ksu < 0. || Wa.Ksd < 0.) {
					s = fmt.Sprintf("ble=%c name=%s 建築一体型空気集熱の熱貫流率Kc、Kdd、Kudが未定義です",
						Wa.ble, Wa.name)
					Eprint("<Walldata>", s)
				}

				if Wa.Ip == -1 {
					s = fmt.Sprintf("ble=%c name=%s 建築一体型空気集熱の空気流通層<P>が未定義です",
						Wa.ble, Wa.name)
					Eprint("<Walldata>", s)
				}
			} else {
				// 壁種類 -> 床暖房等放射パネル
				Wa.WallType = WallType_P
			}
		}

		Wa.N = jw + 1
		Walli(Nbm, W, Wa, pcm)

		*Wall = append(*Wall, Wa)
		i++
	}
}

/* ------------------------------------------------ */

/*
Windowdata (Window Data Input)

この関数は、建物の窓の仕様データを読み込み、
対応する構造体に格納します。
これには、窓の熱的特性、光学特性、および日射制御に関する情報が含まれます。

建築環境工学的な観点:
  - **窓の熱的・光学的特性の定義**: 窓は、建物の熱負荷や昼光利用に大きな影響を与える要素です。
    この関数は、窓の熱的・光学的特性を定義する以下のパラメータを読み込みます。
  - `tgtn`: 日射透過率。窓を透過する日射の割合を示し、日射熱取得量に影響します。
  - `Bn`: 吸収日射取得率。窓が日射を吸収する割合を示し、窓の表面温度上昇に影響します。
  - `Rwall`: 窓部材熱抵抗 [m2K/W]。窓の断熱性能を示し、熱損失・熱取得に影響します。
  - `Ei`, `Eo`: 室内・外表面放射率。窓表面からの放射熱伝達に影響します。
  - `Cidtype`: 入射角特性の種類。日射の入射角によって透過率や吸収率が変化する特性をモデル化します。
  - **日射制御の考慮**: `W.RStrans`は、
    室内透過日射が窓室内側への入射日射を屋外に透過するかどうかを示します。
    これは、窓からの日射熱取得を制御するための重要なパラメータです。
  - **熱負荷と昼光利用のバランス**: 窓は、日射熱取得による冷房負荷の増加や、
    熱損失による暖房負荷の増加といった負の側面を持つ一方で、
    昼光利用による照明エネルギーの削減や、
    視覚的な快適性の向上といった正の側面も持ちます。
    これらのパラメータを適切に設定することで、
    熱負荷と昼光利用のバランスを考慮した窓の設計を検討できます。

この関数は、建物の窓の熱的・光学的挙動を詳細にモデル化し、
熱負荷計算、エネルギー消費量予測、
省エネルギー対策の検討、および快適性評価を行うための重要なデータ入力機能を提供します。
*/
func Windowdata(section *EeTokens, Window *[]*WINDOW) {
	E := fmt.Sprintf(ERRFMT, "WINDOW")

	var N int
	for section.IsEnd() == false {
		line := section.GetLogicalLine()
		if line[0] != "*" {
			N++
		}
	}
	section.Reset()

	if N > 0 {
		*Window = make([]*WINDOW, N)

		for j := 0; j < N; j++ {
			(*Window)[j] = NewWINDOW()
		}
	}

	i := 0
	for section.IsEnd() == false {
		line := section.GetLogicalLine()

		if line[0] == "*" {
			break
		}

		W := (*Window)[i]

		// 名称
		W.Name = line[0]

		// 名前の重複確認
		for k := 0; k < i; k++ {
			Wc := (*Window)[k]
			if W.Name == Wc.Name {
				ss := fmt.Sprintf("<WINDOW>  WindowName Already Defined  (%s)", W.Name)
				Eprint("<Windowdata>", ss)
			}
		}

		// プロパティの設定
		for _, s := range line[1 : len(line)-1] {
			// 室内透過日射が窓室内側への入射日射を屋外に透過する場合'y'
			if s == "-RStrans" {
				W.RStrans = true
			} else {
				//キー・バリューの分離
				st := strings.IndexRune(s, '=')
				if st == -1 {
					panic("Windowdata: invalid format")
				} else {
					st++
				}
				key := s[:st-1]
				value := s[st:]

				// 小数読み取り
				var realValue float64
				var err error
				switch key {
				case "t", "B", "R", "Ei", "Eo":
					realValue, err = strconv.ParseFloat(value, 64)
					if err != nil {
						panic(err)
					}
				}

				// 値の設定
				switch key {
				case "t":
					W.tgtn = realValue // 日射透過率
				case "B":
					W.Bn = realValue // 吸収日射取得率
				case "R":
					W.Rwall = realValue // 窓部材熱抵抗 [m2K/W]
				case "Ei":
					W.Ei = realValue // 室内表面放射率(0.9)
				case "Eo":
					W.Eo = realValue // 外表面放射率(0.9)
				case "Cidtype":
					W.Cidtype = strings.Trim(value, "'") // 入射角特性の種類
				default:
					Err := fmt.Sprintf("%s %s\n", E, s)
					Eprint("<Windowdata>", Err)
				}

				//NOTE: 以下の項目を入力する箇所が不明
				// 窓ガラス面積 Ag, 開口面積 Ao, 幅 W, 高さ H, ??? K
			}
		}

		i++
	}
}

/* --------------------------------------------------- */

/*
Snbkdata (Sunbreak Data Input)

この関数は、日よけ（庇、袖壁、ルーバーなど）の仕様データを読み込み、
対応する構造体に格納します。
これらのデータは、日射遮蔽による日射熱取得の抑制や、
日影の形成による室内環境への影響を評価するために不可欠です。

建築環境工学的な観点:
  - **日射遮蔽の重要性**: 夏季の過度な日射熱取得は、
    冷房負荷の増加や室内温度の上昇を引き起こし、
    エネルギー消費量や快適性に悪影響を与えます。
    日よけは、窓からの日射侵入を効果的に遮蔽し、
    冷房負荷を軽減するパッシブな手法として重要です。
  - **日よけの形状と寸法**: 日よけの種類（庇、袖壁、ルーバーなど）や、
    その寸法（奥行き`D`、高さ`H`、幅`W`、窓からの距離`H1`, `W1`, `W2`など）は、
    日射遮蔽効果に直接影響します。
    この関数は、これらの幾何学的パラメータを読み込み、
    日影計算の基礎データを提供します。
  - **日影計算の基礎**: 日よけによって形成される日影の範囲は、
    太陽の高度角や方位角、そして日よけの形状と寸法によって決まります。
    正確な日影計算は、
  - **日射熱取得の予測**: 窓面への日射入射量を正確に予測し、
    冷房負荷を評価します。
  - **昼光利用の検討**: 日影によって昼光が遮られすぎないか、
    あるいは適切な昼光が確保されるかを検討します。
  - **日影規制の遵守**: 都市部における日影規制を遵守した建物設計を行う上で不可欠です。
  - **入力チェック**: `typstr`や`code`を用いた入力チェックは、
    入力データの整合性を確保し、シミュレーションの信頼性を高めます。

この関数は、建物の日射遮蔽計画をモデル化し、
日射熱取得の抑制、冷房負荷の軽減、
昼光利用の最適化、および日影計算を行うための重要なデータ入力機能を提供します。
*/
func Snbkdata(section *EeTokens, dsn string, Snbk *[]*SNBK) {
	// 入力チェック用パターン文字列
	typstr := []string{
		"HWDTLR.", // 庇
		"HWDTLRB", // 袖壁(その1)　(左右)
		"HWDTL.B", // 袖壁(その2)　(左のみ)
		"HWDT.RB", // 袖壁(その3)　(右のみ)
		"HWDT...", // 長い庇
		"HWD.LR.", // 長い袖壁(その1) (左右)
		"HWD.L..", // 長い袖壁(その2) (左のみ)
		"HWD..R.", // 長い袖壁(その3) (右のみ)
		"HWDTLRB", // ルーバー
	}

	Er := fmt.Sprintf(ERRFMT, dsn)
	Type := 0

	*Snbk = make([]*SNBK, 0)

	for section.IsEnd() == false {

		fields := section.GetLogicalLine()
		if fields[0] == "*" {
			break
		}

		S := NewSNBK()

		// 名前
		S.Name = fields[0]

		// 入力チェック用
		code := [7]rune{'.', '.', '.', '.', '.', '.', '.'}

		for _, s := range fields[1 : len(fields)-1] {
			// キー・バリューの分離
			st := strings.IndexRune(s, '=')
			if st == -1 {
				panic("Snbkdata: invalid format")
			}
			key := s[:st]
			v := s[st+1:]

			var err error
			switch key {
			case "type":
				var vs string
				if v[0] == '-' {
					S.Ksi = 1
					vs = v[1:]
				} else {
					S.Ksi = 0
					vs = v[0:]
				}

				if vs == "H" {
					Type = 1 // 一般の庇
				} else if vs == "HL" {
					Type = 5 // 長い庇
				} else if vs == "S" {
					Type = 2 // 袖壁
				} else if vs == "SL" {
					Type = 6 // 長い袖壁
				} else if vs == "K" {
					Type = 9 // 格子ルーバー
				} else {
					E := fmt.Sprintf("`%s` is invalid", vs)
					Eprint("<Snbkdata>", E)
				}

			case "window":
				// For `window=HhhhxWwww`
				hw := strings.Split(v, "x")
				h, w := hw[0], hw[1]

				if h[0] != 'H' || w[0] != 'W' {
					panic(fmt.Sprintf("Invaid window format: %s", v))
				}

				// 開口部の高さ
				S.H, err = strconv.ParseFloat(h[1:], 64)
				if err != nil {
					panic(err)
				}
				code[0] = 'H'

				// 開口部の幅
				S.W, err = strconv.ParseFloat(w[1:], 64)
				if err != nil {
					panic(err)
				}
				code[1] = 'W'

			case "D":
				// 庇の付け根から先端までの長さ
				S.D, err = strconv.ParseFloat(v, 64)
				if err != nil {
					panic(err)
				}
				code[2] = 'D'

			case "T":
				// 開口部の上端から壁の上端までの距離
				S.H1, err = strconv.ParseFloat(v, 64)
				if err != nil {
					panic(err)
				}
				code[3] = 'T'

			case "L":
				// 開口部の左端から壁の左端までの距離
				S.W1, err = strconv.ParseFloat(v, 64)
				if err != nil {
					panic(err)
				}
				code[4] = 'L'

			case "R":
				// 開口部の右端から壁の右端までの距離
				S.W2, err = strconv.ParseFloat(v, 64)
				if err != nil {
					panic(err)
				}
				code[5] = 'R'

			case "B":
				// 地面から開口部の下端までの高さ
				S.H2, err = strconv.ParseFloat(v, 64)
				if err != nil {
					panic(err)
				}
				code[6] = 'B'

			default:
				panic(fmt.Sprintf("Invaid window format: %s", s))
			}
		}

		// 日除けの種類
		S.Type = Type

		// 日除けの種類ごとに入力チェック
		switch Type {
		case 1, 5, 9:
			// 庇 or ルーバー
			if string(code[:]) != typstr[Type-1] {
				E := fmt.Sprintf("%s %s  type=%d %s\n", Er, fields[0], Type, string(code[:]))
				Eprint("<Snbkdata>", E)
			}
		case 2, 6:
			// 袖壁
			if string(code[:]) != typstr[Type-1] {
				for j := 1; j < 3; j++ {
					if string(code[:]) == typstr[Type+j-1] {
						S.Type = Type + j
						break
					}
					if j == 3 {
						E := fmt.Sprintf("%s %s  type=%d %s\n", Er, fields[0], Type, string(code[:]))
						Eprint("<Snbkdata>", E)
					}
				}
			}
		}

		*Snbk = append(*Snbk, S)
	}
}

/*
NewSNBK (New Sunbreak Object)

この関数は、新しい日よけ（庇、袖壁、ルーバーなど）のデータ構造を初期化します。

建築環境工学的な観点:
  - **日よけの初期化**: 日よけのシミュレーションを行う前に、
    その幾何学的パラメータ（幅`W`、高さ`H`、奥行き`D`、窓からの距離`W1`, `W2`, `H1`, `H2`など）を
    デフォルト値（通常はゼロ）で初期化します。
    これにより、後続のデータ入力や計算が正しく行われるための準備が整います。
  - **日影計算の基礎**: これらのパラメータは、
    日よけによって形成される日影の形状や範囲を計算するための基礎となります。
    正確な初期化は、日射遮蔽効果の評価や、
    日影規制の遵守をシミュレーションする上で重要です。

この関数は、建物の日射遮蔽計画をモデル化し、
日射熱取得の抑制、冷房負荷の軽減、
昼光利用の最適化、および日影計算を行うための基礎的な役割を果たします。
*/
func NewSNBK() *SNBK {
	S := new(SNBK)
	S.Name = ""
	S.W = 0.0
	S.H = 0.0
	S.D = 0.0
	S.W1 = 0.0
	S.W2 = 0.0
	S.H1 = 0.0
	S.H2 = 0.0
	S.Type = 0
	S.Ksi = 0
	return S
}

/************************************************************/

/*
wbmlistcount (Wall Building Material List Count)

この関数は、壁体構成材料のリストファイル（`wbmlist.efl`）を読み込み、
コメント行や空行を除外した有効なデータ行を抽出します。
これは、壁体の熱的特性を定義するための材料データを準備する際に用いられます。

建築環境工学的な観点:
  - **材料データの準備**: 建物の壁体は、様々な材料（コンクリート、断熱材、石膏ボードなど）で構成されます。
    これらの材料の熱伝導率や比熱などの物性値は、
    壁体全体の熱貫流率や熱容量を計算する上で不可欠です。
    この関数は、これらの材料データを効率的に読み込むための前処理を行います。
  - **入力ファイルの解析**: `wbmlist.efl`のような外部ファイルからデータを読み込む際、
    コメント行（`*`で始まる行）や、
    データ以外の情報（`!`で始まるコメント）を除外することで、
    正確なデータのみを抽出します。
    `strings.TrimSpace(s)`は、行頭・行末の空白を除去し、
    データの整形を行います。

この関数は、建物の壁体の熱的挙動を詳細にモデル化し、
熱負荷計算、エネルギー消費量予測、
省エネルギー対策の検討、および快適性評価を行うための重要なデータ準備機能を提供します。
*/
func wbmlistcount(fi string) []string {
	N := make([]string, 0)

	reader := strings.NewReader(fi)
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		s := scanner.Text()

		if strings.HasPrefix(s, "*") {
			break
		}

		//コメントを読み飛ばす
		st := strings.IndexRune(s, '!')
		if st != -1 {
			s = s[:st]
		}

		//ゴミ除去
		s = strings.TrimSpace(s)

		if s != "" {
			N = append(N, s)
		}
	}

	return N
}

/************************************************************/

/*
Wallcount (Wall Count)

この関数は、壁体データが記述された入力ストリームから、
コメント行や空行を除外した有効なデータ行を抽出し、その数をカウントします。
これは、壁体の熱的特性を定義するためのデータを準備する際に用いられます。

建築環境工学的な観点:
  - **壁体データの準備**: 建物の壁体は、様々な材料（コンクリート、断熱材、石膏ボードなど）で構成されます。
    これらの材料の熱伝導率や比熱などの物性値は、
    壁体全体の熱貫流率や熱容量を計算する上で不可欠です。
    この関数は、これらの材料データを効率的に読み込むための前処理を行います。
  - **入力ファイルの解析**: `scanner`からデータを読み込む際、
    コメント行（`*`で始まる行）や、
    データ以外の情報（`!`で始まるコメント）を除外することで、
    正確なデータのみを抽出します。
    `strings.TrimSpace(s)`は、行頭・行末の空白を除去し、
    データの整形を行います。

この関数は、建物の壁体の熱的挙動を詳細にモデル化し、
熱負荷計算、エネルギー消費量予測、
省エネルギー対策の検討、および快適性評価を行うための重要なデータ準備機能を提供します。
*/
func Wallcount(scanner *bufio.Scanner) []string {
	N := make([]string, 0)

	for scanner.Scan() {
		s := scanner.Text()

		if strings.HasPrefix(s, "*") {
			break
		}

		//コメントを読み飛ばす
		st := strings.IndexRune(s, '!')
		if st != -1 {
			s = s[:st]
		}

		//ゴミ除去
		s = strings.TrimSpace(s)

		if s != "" {
			N = append(N, s)
		}
	}

	return N
}

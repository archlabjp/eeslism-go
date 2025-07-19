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

/*  mcstank.c */

/*  95/11/17 rev  */

package eeslism

import (
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
	"unicode"
)

/*
Stankdata (Storage Tank Data Input)

この関数は、蓄熱槽の各種仕様（容量、熱損失係数など）を読み込み、
対応する蓄熱槽の構造体に格納します。
これらのデータは、蓄熱システムにおける蓄熱槽の性能評価、
熱負荷平準化、およびエネルギー消費量予測に不可欠です。

建築環境工学的な観点:
  - **蓄熱槽の役割**: 蓄熱槽は、熱源設備で発生した熱（または冷熱）を一時的に貯蔵し、
    熱需要に応じて供給することで、熱負荷の平準化や熱源設備の効率的な運転を可能にします。
    これにより、熱源設備の容量を小さくしたり、
    電力料金の安い夜間電力などを利用したりすることができ、
    省エネルギーやランニングコストの削減に貢献します。
  - **容量 (Vol)**:
    蓄熱槽が貯蔵できる熱媒の体積を示します。
    建物の熱負荷パターンや、蓄熱による熱負荷平準化の目標に応じて、適切な容量を選定する必要があります。
  - **熱損失係数 (KAside, KAtop, KAbtm)**:
    蓄熱槽からの熱損失は、その断熱性能や周囲環境との温度差に依存します。
    `KAside`（側面）、`KAtop`（上面）、`KAbtm`（下面）は、
    それぞれの部位からの熱損失を表す係数であり、
    蓄熱槽のエネルギー効率を評価する上で重要です。
    熱損失が大きいと、蓄熱された熱が有効に利用されず、エネルギーの無駄につながります。
  - **温度成層 (gxr)**:
    `gxr`は、蓄熱槽内の温度成層（温度の異なる層が形成される現象）の度合いを示すパラメータです。
    温度成層が良好な蓄熱槽は、熱媒の温度差を大きく保つことができ、
    熱源設備や熱利用設備との効率的な熱交換を可能にします。
    これにより、蓄熱槽の有効利用率が向上し、システム全体の効率が高まります。
  - **運転モード (Type, tparm)**:
    `Type`や`tparm`は、蓄熱槽の運転モード（例: 連続運転、バッチ運転）や、
    初期温度分布の設定方法を示唆します。

この関数は、蓄熱槽の性能をモデル化し、
熱負荷平準化、エネルギー消費量予測、
および省エネルギー対策の検討を行うための重要なデータ入力機能を提供します。
*/
func Stankdata(f *EeTokens, s string, Stankca *STANKCA) int {
	id := 0
	st := ""
	Stankca.gxr = 0.0

	var err error

	if stIdx := strings.IndexByte(s, '='); stIdx != -1 {
		s = strings.TrimSpace(s)
		st = s[stIdx+1:]

		switch {
		case strings.HasPrefix(s, "Vol"):
			Stankca.Vol, err = strconv.ParseFloat(st, 64)
		case strings.HasPrefix(s, "KAside"):
			Stankca.KAside, err = strconv.ParseFloat(st, 64)
		case strings.HasPrefix(s, "KAtop"):
			Stankca.KAtop, err = strconv.ParseFloat(st, 64)
		case strings.HasPrefix(s, "KAbtm"):
			Stankca.KAbtm, err = strconv.ParseFloat(st, 64)
		case strings.HasPrefix(s, "gxr"):
			Stankca.gxr, err = strconv.ParseFloat(st, 64)
		default:
			id = 1
		}

		if err != nil {
			fmt.Println(err)
		}

	} else if s == "-S" {
		st = ""
		s = f.GetToken()
		s += " *"
		Stankca.tparm = s
	} else {
		Stankca.name = s
		Stankca.Type = 'C'
		Stankca.tparm = ""
		Stankca.Vol = FNAN
		Stankca.KAside = FNAN
		Stankca.KAtop = FNAN
		Stankca.KAbtm = FNAN
		Stankca.gxr = 0.0
	}

	return id
}

/*
Stankmemloc (Storage Tank Memory Allocation)

この関数は、蓄熱槽のシミュレーションに必要なメモリ領域を確保し、
特に蓄熱槽の分割数や、入出力ポート、内蔵熱交換器などの情報を設定します。
これは、蓄熱槽内部の温度分布や熱交換を詳細にモデル化するために不可欠です。

建築環境工学的な観点:
  - **蓄熱槽の分割モデル (Ndiv)**:
    蓄熱槽内部の温度分布をより正確にモデル化するために、
    蓄熱槽を複数の層（分割数`Ndiv`）に分割して扱います。
    各層の温度を個別に計算することで、
    蓄熱槽内の温度成層（温度の異なる層が形成される現象）を再現し、
    蓄熱槽の有効利用率を評価できます。
  - **入出力ポートの設定 (Pthcon, Jin, Jout)**:
    蓄熱槽への熱媒の流入・流出経路（`Pthcon`）や、
    各流入・流出が蓄熱槽のどの層（`Jin`, `Jout`）に接続されているかを設定します。
    これにより、蓄熱槽の運転モード（例: 上部から流入、下部から流出）や、
    熱媒の循環経路をモデル化できます。
  - **内蔵熱交換器のモデル化 (Ihex, Ihxeff, KA, KAinput)**:
    蓄熱槽に内蔵された熱交換器の有無（`Ihex`）、
    その効率（`Ihxeff`）や熱通過率と伝熱面積の積（`KA`）を設定します。
    `KAinput`は、`KA`値が入力値として与えられるか、
    あるいは熱交換器の寸法（`Dbld0`, `DblL`）から計算されるかを示唆します。
    これにより、蓄熱槽を介した熱源設備や熱利用設備との熱交換を詳細にモデル化できます。
  - **温度分布の記憶 (Tss, Tssold)**:
    各層の現在の温度（`Tss`）と前時刻の温度（`Tssold`）を記憶するためのメモリを確保します。
    これにより、蓄熱槽内部の温度変化を追跡し、
    蓄熱量や熱損失量を計算できます。

この関数は、蓄熱槽の熱的挙動を詳細にモデル化し、
熱負荷平準化、エネルギー消費量予測、
および省エネルギー対策の検討を行うための重要な初期設定機能を提供します。
*/
func Stankmemloc(errkey string, Stank *STANK) {
	var np, Ndiv, Nin int
	var st, stt, ss string
	var parm []string = make([]string, 0)

	st = Stank.Cat.tparm[:]

	// 読み飛ばし処理
	np = 0
	for {
		_, err := fmt.Sscanf(st, "%s", &ss)
		if err != nil || ss == "*" {
			break
		}

		parm = append(parm, st)
		np++
		st = st[len(ss):]
		for st[0] == ' ' || st[0] == '\t' {
			st = st[1:]
		}
	}

	Stank.Pthcon = make([]ELIOType, np)
	Stank.Batchcon = make([]ControlSWType, np)
	Stank.Ihex = make([]rune, np)
	Stank.Jin = make([]int, np)
	Stank.Jout = make([]int, np)
	Stank.Ihxeff = make([]float64, np)
	Stank.KA = make([]float64, np)
	Stank.KAinput = make([]rune, np)

	i := 0

	for j := 0; j < np; j++ {
		_, err := fmt.Sscanf(parm[j], "%s", &ss)
		if err != nil {
			panic(err)
		}

		if strings.HasPrefix(ss, "N=") {
			Stank.Ndiv, err = strconv.Atoi(ss[2:])
			if err != nil {
				panic(err)
			}
		} else if stIdx := strings.IndexRune(ss, ':'); stIdx != -1 {
			Stank.Pthcon[i] = ELIOType(ss[0])
			tmp, err := strconv.Atoi(ss[stIdx+1:])
			if err != nil {
				panic(err)
			} else {
				Stank.Jin[i] = tmp - 1
			}

			if sttIdx := strings.IndexRune(ss[stIdx+1:], '-'); sttIdx != -1 {
				stt = ss[stIdx+1:]
				Stank.Ihex[i] = 'n'
				Stank.Ihxeff[i] = 1.0
				tmp, err := strconv.Atoi(stt)
				if err != nil {
					panic(err)
				} else {
					Stank.Jout[i] = tmp - 1
				}
			} else if sttIdx := strings.IndexRune(ss[stIdx+1:], '_'); sttIdx != -1 {
				stt = ss[stIdx+1 : sttIdx]
				Stank.Ihex[i] = 'y'

				if stt[1] == 'e' { // 温度効率が入力されている場合
					Stank.Ihxeff[i], err = strconv.ParseFloat(stt[5:], 64)
					if err != nil {
						panic(err)
					}
				} else if stt[1] == 'K' { // 内蔵熱交のKAが入力されている場合
					Stank.KAinput[i] = 'Y'
					Stank.KA[i], err = strconv.ParseFloat(stt[4:], 64)
					if err != nil {
						panic(err)
					}
				} else if stt[1] == 'd' {
					Stank.KAinput[i] = 'C' // 内蔵熱交換器の内径と長さが入力されている場合
					stpIdx := strings.IndexRune(stt[4:], '_')
					Stank.Dbld0, err = strconv.ParseFloat(stt[4:], 64)
					if err != nil {
						panic(err)
					}
					Stank.DblL, err = strconv.ParseFloat(stt[stpIdx+1:], 64)
					if err != nil {
						panic(err)
					}
					Stank.Ncalcihex++
				}

				Stank.Jout[i] = Stank.Jin[i]

				i++
			}
		}
	}

	Stank.Nin = i
	Nin = i

	Ndiv = Stank.Ndiv
	Stank.DtankF = make([]rune, Ndiv)

	Stank.B = make([]float64, Ndiv*Ndiv)
	Stank.R = make([]float64, Ndiv)
	Stank.D = make([]float64, Ndiv)
	Stank.Fg = make([]float64, Ndiv*Nin)
	Stank.Tss = make([]float64, Ndiv)

	Stank.Tssold = make([]float64, Ndiv)
	Stank.Dvol = make([]float64, Ndiv)
	Stank.Mdt = make([]float64, Ndiv)
	Stank.KS = make([]float64, Ndiv)
	Stank.CGwin = make([]float64, Nin)
	Stank.EGwin = make([]float64, Nin)
	Stank.Twin = make([]float64, Nin)
	Stank.Q = make([]float64, Nin)
	if Nin > 0 {
		Stank.Stkdy = make([]STKDAY, Nin)
	}
	if Nin > 0 {
		Stank.Mstkdy = make([]STKDAY, Nin)
	}
}

/*
Stankint (Storage Tank Initialization)

この関数は、蓄熱槽の初期設定を行います。
特に、蓄熱槽内部の初期温度分布や、周囲環境への熱損失に関するパラメータを設定します。

建築環境工学的な観点:
  - **初期温度分布の設定**: 蓄熱槽のシミュレーションを開始する際の、
    蓄熱槽内部の初期温度分布を設定します。
    `stank.DtankF`は各層が満水か空かを示し、
    `stank.Tssold`は各層の初期温度を示します。
    これにより、シミュレーション開始時の蓄熱槽の状態を正確にモデル化できます。
  - **周囲環境への熱損失**: 蓄熱槽は、周囲環境との温度差によって熱損失が発生します。
    `stank.Tenv`は蓄熱槽周囲の環境温度へのポインターであり、
    熱損失計算に用いられます。
    `stank.Cat.KAside`, `stank.Cat.KAtop`, `stank.Cat.KAbtm`は、
    蓄熱槽の側面、上面、下面からの熱損失係数であり、
    蓄熱槽の断熱性能を評価する上で重要です。
  - **熱損失係数の計算 (stoint)**:
    `stoint`関数は、蓄熱槽の容量、熱損失係数、分割数などに基づいて、
    各層の熱容量（`stank.Mdt`）や熱損失係数（`stank.KS`）を計算します。
    これにより、蓄熱槽内部の温度変化や熱損失を正確にモデル化できます。
  - **パラメータの妥当性チェック**: `if stank.Cat.Vol < 0.0` のようなエラーチェックは、
    入力されたパラメータが物理的に妥当な範囲内にあるかを確認するために重要です。
    不適切なパラメータは、シミュレーション結果の信頼性を損なう可能性があります。

この関数は、蓄熱槽の熱的挙動を詳細にモデル化し、
熱負荷平準化、エネルギー消費量予測、
および省エネルギー対策の検討を行うための重要な初期設定機能を提供します。
*/
func Stankint(Stank []*STANK, Simc *SIMCONTL, Compnt []*COMPNT, Wd *WDAT) {
	var s, ss, Err, E string
	var mrk rune
	var Tso float64

	E = "Stankint"

	for _, stank := range Stank {

		// 内蔵熱交換器の熱伝達率計算用温度の初期化
		stank.DblTa = 20.0
		stank.DblTw = 20.0

		s = stank.Cmp.Tparm
		if s != "" {
			if s[0] == '(' {
				s = s[1:]
				for j := 0; j < stank.Ndiv; j++ {
					_, err := fmt.Sscanf(s, " %s ", &ss)
					if err != nil {
						panic(err)
					}

					if ss[0] == TANK_EMPTY {
						stank.DtankF[j] = TANK_EMPTY
						stank.Tssold[j] = TANK_EMPTMP
					} else {
						stank.DtankF[j] = TANK_FULL
						stank.Tssold[j], err = strconv.ParseFloat(ss, 64)
						if err != nil {
							panic(err)
						}
					}
					s = s[len(ss):]
					for s[0] == ' ' {
						s = s[1:]
					}
				}
			} else {
				if s[0] == TANK_EMPTY {
					mrk = TANK_EMPTY
					Tso = TANK_EMPTMP
				} else {
					var err error
					mrk = TANK_FULL
					Tso, err = strconv.ParseFloat(s, 64)
					if err != nil {
						panic(err)
					}
				}
				for j := 0; j < stank.Ndiv; j++ {
					stank.DtankF[j] = mrk
					stank.Tssold[j] = Tso
				}
			}
		}

		stank.Tenv = envptr(stank.Cmp.Envname, Simc, Compnt, Wd, nil)
		stoint(stank.Ndiv, stank.Cat.Vol, stank.Cat.KAside, stank.Cat.KAtop, stank.Cat.KAbtm,
			stank.Dvol, stank.Mdt, stank.KS, stank.Tss, stank.Tssold, &stank.Jva, &stank.Jvb)

		if stank.Cat.Vol < 0.0 {
			Err = fmt.Sprintf("Name=%s  Vol=%.4g", stank.Cmp.Name, stank.Cat.Vol)
			Eprint(E, Err)
		}
		if stank.Cat.KAside < 0.0 {
			Err = fmt.Sprintf("Name=%s  KAside=%.4g", stank.Cmp.Name, stank.Cat.KAside)
			Eprint(E, Err)
		}
		if stank.Cat.KAtop < 0.0 {
			Err = fmt.Sprintf("Name=%s  KAtop=%.4g", stank.Cmp.Name, stank.Cat.KAtop)
			Eprint(E, Err)
		}
		if stank.Cat.KAbtm < 0.0 {
			Err = fmt.Sprintf("Name=%s  KAbtm=%.4g", stank.Cmp.Name, stank.Cat.KAbtm)
			Eprint(E, Err)
		}
	}
}

/*
Stankcfv (Storage Tank Characteristic Function Value Calculation)

この関数は、蓄熱槽の運転特性を評価し、
熱媒の流量、熱容量流量、そして内蔵熱交換器の効率などを計算します。
これは、蓄熱槽が熱源設備や熱利用設備とどのように熱を交換するかをモデル化するために不可欠です。

建築環境工学的な観点:
  - **熱媒の熱容量流量 (cGwin)**:
    蓄熱槽への熱媒の流入量と比熱から計算される熱容量流量は、
    蓄熱槽が受け入れる熱量に直接影響します。
    `Spcheat('W') * elin.Lpath.G` のように、熱媒の比熱と質量流量から計算されます。
  - **内蔵熱交換器の効率 (ihxeff)**:
    蓄熱槽に内蔵された熱交換器の効率は、
    蓄熱槽を介した熱源設備や熱利用設備との熱交換能力に影響します。
    `stank.KAinput[j] == 'C'` の場合、熱交換器の寸法（内径`Dbld0`、長さ`DblL`）から、
    熱伝達率（`ho`, `hi`）を計算し、それに基づいて`KA`値を算出し、
    最終的に効率（`ihxeff`）を計算します。
    `1.0 - math.Exp(-NTU)` の式は、NTU（Number of Transfer Units）法に基づいた効率計算です。
  - **有効熱容量流量 (EGwin)**:
    `*EGwin = *cGwin * *ihxeff` のように、
    熱容量流量に熱交換器の効率を乗じることで、
    実際に熱交換に寄与する有効な熱容量流量を計算します。
  - **蓄熱槽内部の熱伝達モデル (stofc)**:
    `stofc`関数は、蓄熱槽の分割数、入出力ポート、内蔵熱交換器の特性などに基づいて、
    蓄熱槽内部の熱伝達を記述する係数行列（`B`, `R`, `D`, `Fg`）を計算します。
    これにより、蓄熱槽内部の温度分布や熱流を詳細にモデル化できます。

この関数は、蓄熱槽の熱的挙動を詳細にモデル化し、
熱負荷平準化、エネルギー消費量予測、
および省エネルギー対策の検討を行うための重要な役割を果たします。
*/
func Stankcfv(Stank []*STANK) {
	for _, stank := range Stank {
		for j := 0; j < stank.Nin; j++ {
			elin := stank.Cmp.Elins[j]
			cGwin := &stank.CGwin[j]
			EGwin := &stank.EGwin[j]
			ihxeff := &stank.Ihxeff[j]
			ihex := &stank.Ihex[j]

			if elin.Lpath.Batch {
				*cGwin = 0.0
			} else {
				*cGwin = Spcheat('W') * elin.Lpath.G
			}

			// 内蔵熱交のKAが入力されている場合
			if *ihex == 'y' && *cGwin > 0.0 {
				// 内蔵熱交換器の内径、管長が入力されている場合
				if stank.KAinput[j] == 'C' {
					dblT := (stank.DblTa + stank.DblTw) / 2.0
					// 内蔵熱交換器の表面温度は内外流体の平均温度で代用
					ho := FNhoutpipe(stank.Dbld0, dblT, stank.DblTw)
					// 流速の計算
					dblv := elin.Lpath.G / Row / (math.Pi * math.Pow(stank.Dbld0/2.0, 2.0))
					hi := FNhinpipe(stank.Dbld0, stank.DblL, dblv, dblT)
					stank.KA[j] = 1.0 / (1.0/ho + 1.0/hi) * math.Pi * stank.Dbld0 * stank.DblL
				}
				if stank.KAinput[j] == 'Y' || stank.KAinput[j] == 'C' {
					NTU := stank.KA[j] / *cGwin
					*ihxeff = 1.0 - math.Exp(-NTU)
				}
			}
			*EGwin = *cGwin * *ihxeff
		}

		stofc(stank.Ndiv, stank.Nin, stank.Jin,
			stank.Jout, stank.Ihex, stank.Ihxeff, stank.Jva, stank.Jvb,
			stank.Mdt, stank.KS, stank.Cat.gxr, stank.Tenv,
			stank.Tssold, stank.CGwin, stank.EGwin, stank.B, stank.R, stank.D, stank.Fg)

		fgIdx := 0
		cfinIdx := 0
		for j := 0; j < stank.Nin; j++ {
			Eo := stank.Cmp.Elouts[j]
			Eo.Coeffo = 1.0
			Eo.Co = stank.D[stank.Jout[j]]

			for k := 0; k < stank.Nin; k++ {
				Eo.Coeffin[cfinIdx] = -stank.Fg[fgIdx]
				cfinIdx++
				fgIdx++
			}
		}
	}
}

/*
stankvptr (Storage Tank Internal Water Temperature Pointer Setting)

この関数は、蓄熱槽内部の各層の水温へのポインターを設定します。
これにより、蓄熱槽の運転を特定の目標温度に追従させる制御をモデル化できます。

建築環境工学的な観点:
  - **蓄熱槽の温度制御**: 蓄熱槽は、熱源設備で発生した熱を貯蔵し、
    熱需要に応じて供給することで、熱負荷の平準化を図ります。
    この際、蓄熱槽内部の温度を適切に制御することが、
    蓄熱効率や熱供給の安定性に影響します。
  - **制御対象の指定**: `key[1]`が`"Ts"`の場合、
    蓄熱槽内部の水温を制御対象とすることを意味します。
    `key[2]`によって、蓄熱槽の最上層（`'t'`）、最下層（`'b'`）、
    または特定の層（数値）の温度を選択できます。
    `vptr.Ptr = &Stank.Tssold[i]` は、
    選択された層の温度変数へのポインターを設定し、
    `vptr.Type = VAL_CTYPE` は、そのポインターが制御値であることを示します。
  - **温度成層の利用**: 蓄熱槽内の温度成層を利用することで、
    熱源設備からの高温水は上層に、熱利用設備への低温水は下層から供給されるように制御し、
    蓄熱槽の有効利用率を高めることができます。
    このポインター設定は、そのような温度成層を考慮した制御のモデル化を可能にします。
  - **フィードバック制御の基礎**: このポインター設定は、
    蓄熱槽のフィードバック制御の基礎となります。
    シミュレーションの各時間ステップで、
    現在の蓄熱槽内部温度と目標温度を比較し、その差に基づいて熱源設備や熱利用設備の運転を調整します。

この関数は、蓄熱槽の制御ロジックをモデル化し、
熱負荷平準化、エネルギー消費量予測、
および省エネルギー対策の検討を行うための重要な役割を果たします。
*/
func stankvptr(key []string, Stank *STANK) (VPTR, error) {
	var err error
	var vptr VPTR
	var s string
	if key[1] == "Ts" {
		s = key[2]
		if unicode.IsLetter(rune(s[0])) {
			if s[0] == 't' {
				vptr.Ptr = &Stank.Tssold[0]
				vptr.Type = VAL_CTYPE
			} else if s[0] == 'b' {
				vptr.Ptr = &Stank.Tssold[Stank.Ndiv-1]
				vptr.Type = VAL_CTYPE
			} else {
				err = errors.New("'t' or 'b' is expected")
			}
		} else {
			i, _ := strconv.Atoi(s)
			if i >= 0 && i < Stank.Ndiv {
				vptr.Ptr = &Stank.Tssold[i]
				vptr.Type = VAL_CTYPE
			} else {
				err = errors.New("numeric value is expected")
			}
		}
	} else {
		err = errors.New("'Ts' is expected")
	}

	return vptr, err
}

/*
Stanktss (Storage Tank Internal Water Temperature Calculation and Stratification Check)

この関数は、蓄熱槽内部の各層の水温を計算し、
水温分布の逆転（温度成層の崩壊）が発生していないかをチェックします。
これは、蓄熱槽の効率的な運用と、シミュレーションの安定性に不可欠です。

建築環境工学的な観点:
  - **蓄熱槽の温度成層**: 蓄熱槽は、温度の異なる水が層状に分かれる「温度成層」を形成することで、
    熱源設備からの高温水と熱利用設備への低温水を効率的に供給できます。
    この関数は、熱媒の流入・流出や熱損失を考慮して、
    各層の温度（`stank.Tss`）を計算します。
  - **水温分布逆転のチェック (stotsexm)**:
    蓄熱槽の運転状況によっては、温度成層が崩れて水温分布が逆転する（下層の水温が上層より高くなる）ことがあります。
    これは、蓄熱槽の効率を低下させ、熱源設備や熱利用設備の運転に悪影響を与える可能性があります。
    `stotsexm`関数は、この水温分布の逆転を検出し、
    必要に応じてシミュレーションの再計算（`*TKreset = 1`）を促します。
  - **シミュレーションの安定性**: 水温分布の逆転は、
    シミュレーションモデルの不安定性を示す場合があり、
    正確な結果を得るためには、この問題を解決する必要があります。
    このチェックと再計算のメカニズムは、
    シミュレーションの安定性と信頼性を確保するために重要です。
  - **入水温度の考慮**: `stank.Twin[j] = eli.Sysvin` のように、
    各入水ポートからの熱媒の温度（`stank.Twin`）を考慮して、
    蓄熱槽内部の温度分布を計算します。

この関数は、蓄熱槽の熱的挙動を詳細にモデル化し、
熱負荷平準化、エネルギー消費量予測、
および省エネルギー対策の検討を行うための重要な役割を果たします。
*/
func Stanktss(Stank []*STANK, TKreset *int) {
	for _, stank := range Stank {

		for j := 0; j < stank.Nin; j++ {
			eli := stank.Cmp.Elins[j]
			stank.Twin[j] = eli.Sysvin
		}

		stotss(stank.Ndiv, stank.Nin, stank.Jin, stank.B, stank.R, stank.EGwin, stank.Twin,
			stank.Tss)

		stotsexm(stank.Ndiv, stank.Tss, &stank.Jva, &stank.Jvb,
			stank.DtankF, &stank.Cfcalc)

		if stank.Cfcalc == 'y' {
			*TKreset = 1
		}
	}
}

/*
Stankene (Storage Tank Energy Calculation)

この関数は、蓄熱槽の供給熱量、熱損失量、および蓄熱量を計算します。
これは、蓄熱槽のエネルギー収支を詳細に分析し、
熱負荷平準化の効果やエネルギー消費量を評価する上で不可欠です。

建築環境工学的な観点:
  - **供給熱量 (stank.Q)**:
    蓄熱槽から熱利用設備へ供給される熱量です。
    `EGwin * (stank.Tss[Jo] - Twin)` のように、
    有効熱容量流量と蓄熱槽の出口温度（`stank.Tss[Jo]`）および入口温度（`Twin`）の差から計算されます。
    これは、蓄熱槽が熱需要にどれだけ貢献できるかを示します。
  - **熱損失量 (stank.Qloss)**:
    蓄熱槽から周囲環境へ逃げる熱量です。
    `stank.KS[j] * (stank.Tss[j] - *stank.Tenv)` のように、
    各層の熱損失係数と層の温度および周囲環境温度の差から計算されます。
    熱損失は、蓄熱槽のエネルギー効率に影響を与え、
    特に長期間の蓄熱では無視できない要因となります。
  - **蓄熱量 (stank.Qsto)**:
    蓄熱槽内部に蓄えられた熱量の変化量です。
    `stank.Mdt[j] * (stank.Tss[j] - stank.Tssold[j])` のように、
    各層の熱容量と温度変化から計算されます。
    これは、蓄熱槽が熱をどれだけ貯蔵できたかを示します。
  - **バッチモードの考慮**: `stank.Batchop == BTFILL` の条件は、
    蓄熱槽がバッチモードで運転されている場合に、
    熱媒の供給方法をモデル化します。
    これにより、特定の運転戦略における蓄熱槽の挙動を再現できます。
  - **内蔵熱交換器の考慮**: `stank.KAinput[j] == 'C'` の条件は、
    内蔵熱交換器が設置されている場合に、
    その熱交換器の熱伝達率を計算するために必要な温度（`stank.DblTa`, `stank.DblTw`）を更新します。

この関数は、蓄熱槽のエネルギー収支を詳細に分析し、
熱負荷平準化、エネルギー消費量予測、
および省エネルギー対策の検討を行うための重要な役割を果たします。
*/
func Stankene(Stank []*STANK) {
	for _, stank := range Stank {
		// バッチモードチェック（各層が空かどうかをチェック）
		for k := 0; k < stank.Ndiv; k++ {
			if stank.DtankF[k] == TANK_EMPTY {
				stank.Tss[k] = TANK_EMPTMP
			}
		}

		// バッチモードの水供給
		if stank.Batchop == BTFILL {
			Tsm := 0.0
			for k := 0; k < stank.Ndiv; k++ {
				if stank.DtankF[k] == TANK_EMPTY {
					stank.DtankF[k] = TANK_FULL
					for j := 0; j < stank.Nin; j++ {
						if stank.Batchcon[j] == BTFILL {
							stank.Tss[k] = stank.Twin[j]
						}
					}
				}
				Tsm += stank.Tss[k]
			}
			Tsm /= float64(stank.Ndiv)
			for k := 0; k < stank.Ndiv; k++ {
				stank.Tss[k] = Tsm
			}
		}

		for j := 0; j < stank.Nin; j++ {
			Jo := stank.Jout[j]
			Q := &stank.Q[j]
			EGwin := stank.EGwin[j]
			Twin := stank.Twin[j]
			// ihex := stank.Ihex[j]

			*Q = EGwin * (stank.Tss[Jo] - Twin)

			// // 内蔵熱交換器の場合
			if stank.KAinput[j] == 'C' {
				stank.DblTa = stank.Tss[Jo]
				if EGwin > 0.0 {
					stank.DblTw = Twin
				}
			}
		}

		stank.Qloss = 0.0
		stank.Qsto = 0.0
		for j := 0; j < stank.Ndiv; j++ {
			if stank.DtankF[j] == TANK_FULL {
				stank.Qloss += stank.KS[j] * (stank.Tss[j] - *stank.Tenv)
				if stank.Tssold[j] > -273.0 {
					stank.Qsto += stank.Mdt[j] * (stank.Tss[j] - stank.Tssold[j])
				}
			}
			stank.Tssold[j] = stank.Tss[j]
		}
	}
}

/* ------------------------------------------------------- */

// 代表日の出力
func stankcmpprt(fo io.Writer, id int, Stank []*STANK) {
	switch id {
	case 0:
		if len(Stank) > 0 {
			fmt.Fprintf(fo, "%s %d\n", STANK_TYPE, len(Stank))
		}

		for _, stank := range Stank {
			fmt.Fprintf(fo, "%s:%d", stank.Name, stank.Nin)
			for i := 0; i < stank.Nin; i++ {
				fmt.Fprintf(fo, "%c", stank.Cmp.Idi[i])
			}

			fmt.Fprintf(fo, " 1 %d\n", stank.Nin*5+2+stank.Ndiv+stank.Ncalcihex)
		}
	case 1:
		for _, stank := range Stank {
			for i := 0; i < stank.Nin; i++ {
				c := stank.Cmp.Idi[i]
				fmt.Fprintf(fo, "%s:%c_c c c %s:%c_G m f %s:%c_Ti t f %s:%c_To t f %s:%c_Q q f  ",
					stank.Name, c, stank.Name, c, stank.Name, c, stank.Name, c, stank.Name, c)
				if stank.KAinput[i] == 'C' {
					fmt.Fprintf(fo, "%s:%c_KA q f  ", stank.Name, c)
				}
				fmt.Fprintln(fo)
			}
			fmt.Fprintf(fo, "%s_Qls q f %s_Qst q f\n ", stank.Name, stank.Name)
			for i := 0; i < stank.Ndiv; i++ {
				fmt.Fprintf(fo, "%s_Ts[%d] t f ", stank.Name, i+1)
			}
			fmt.Fprintln(fo)
		}
	default:
		for _, stank := range Stank {
			Tss := &stank.Tss[0]
			for i := 0; i < stank.Nin; i++ {
				Ei := stank.Cmp.Elins[i]
				Twin := &stank.Twin[i]
				Q := &stank.Q[i]
				Eo := stank.Cmp.Elouts[i]
				fmt.Fprintf(fo, "%c %.5g %4.1f %4.1f %3.0f  ", Ei.Lpath.Control,
					Eo.G, *Twin, Eo.Sysv, *Q)

				if stank.KAinput[i] == 'C' {
					if Eo.G > 0.0 {
						fmt.Fprintf(fo, "%.2f  ", stank.KA[i])
					} else {
						fmt.Fprintf(fo, "%.2f  ", 0.0)
					}
				}
			}
			fmt.Fprintf(fo, "%2.0f %3.0f\n", stank.Qloss, stank.Qsto)

			for i := 0; i < stank.Ndiv; i++ {
				fmt.Fprintf(fo, " %4.1f", *Tss)
				Tss = &stank.Tss[i+1]
			}
			fmt.Fprintln(fo)
		}
	}
}

/* ------------------------------------------------------- */
func stankivprt(fo io.Writer, id int, Stank []*STANK) {
	if id == 0 && len(Stank) > 0 {
		for m, stank := range Stank {
			fmt.Fprintf(fo, "m=%d  %s  %d\n", m, stank.Name, stank.Ndiv)
		}
	} else {
		for m, stank := range Stank {
			fmt.Fprintf(fo, "m=%d  ", m)

			for i := 0; i < stank.Ndiv; i++ {
				fmt.Fprintf(fo, " %5.1f", stank.Tss[i])
			}
			fmt.Fprintln(fo)
		}
	}
}

/* --------------------------- */

/* 日積算値に関する処理 */

func stankdyint(Stank []*STANK) {
	for _, stank := range Stank {
		stank.Qlossdy = 0.0
		stank.Qstody = 0.0

		for j := 0; j < stank.Nin; j++ {
			s := &stank.Stkdy[j]
			svdyint(&s.Tidy)
			svdyint(&s.Tsdy)
			qdyint(&s.Qdy)
		}
	}
}

func stankmonint(Stank []*STANK) {
	for _, stank := range Stank {
		stank.MQlossdy = 0.0
		stank.MQstody = 0.0

		for j := 0; j < stank.Nin; j++ {
			s := &stank.Mstkdy[j]
			svdyint(&s.Tidy)
			svdyint(&s.Tsdy)
			qdyint(&s.Qdy)
		}
	}
}

// 日集計、月集計
func stankday(Mon, Day, ttmm int, Stank []*STANK, Nday, SimDayend int) {
	for _, stank := range Stank {

		// 日集計
		Ts := 0.0

		S := &stank.Stkdy[0]
		for j := 0; j < stank.Ndiv; j++ {
			Ts += stank.Tss[j] / float64(stank.Ndiv)
		}
		svdaysum(int64(ttmm), ON_SW, Ts, &S.Tsdy)

		stank.Qlossdy += stank.Qloss
		stank.Qstody += stank.Qsto

		for j := 0; j < stank.Nin; j++ {
			Ei := stank.Cmp.Elins[j]
			S := &stank.Stkdy[j]
			svdaysum(int64(ttmm), Ei.Lpath.Control, stank.Twin[j], &S.Tidy)
			qdaysum(int64(ttmm), Ei.Lpath.Control, stank.Q[j], &S.Qdy)
		}

		// 月集計
		S = &stank.Mstkdy[0]
		svmonsum(Mon, Day, ttmm, ON_SW, Ts, &S.Tsdy, Nday, SimDayend)

		stank.MQlossdy += stank.Qloss
		stank.MQstody += stank.Qsto

		for j := 0; j < stank.Nin; j++ {
			Ei := stank.Cmp.Elins[j]
			S := &stank.Mstkdy[j]
			svmonsum(Mon, Day, ttmm, Ei.Lpath.Control, stank.Twin[j], &S.Tidy, Nday, SimDayend)
			qmonsum(Mon, Day, ttmm, Ei.Lpath.Control, stank.Q[j], &S.Qdy, Nday, SimDayend)
		}
	}
}

// 日集計の出力
func stankdyprt(fo io.Writer, id int, Stank []*STANK) {
	switch id {
	case 0:
		if len(Stank) > 0 {
			fmt.Fprintf(fo, "%s %d\n", STANK_TYPE, len(Stank))
		}

		for _, stank := range Stank {
			fmt.Fprintf(fo, "%s:%d", stank.Name, stank.Nin)

			for i := 0; i < stank.Nin; i++ {
				fmt.Fprintf(fo, "%c", stank.Cmp.Idi[i])
			}

			fmt.Fprintf(fo, " 1 %d\n", stank.Nin*14+2+1)
		}

	case 1:
		for _, stank := range Stank {
			fmt.Fprintf(fo, "%s_Ts t f \n", stank.Name)

			for i := 0; i < stank.Nin; i++ {
				c := stank.Cmp.Idi[i]
				fmt.Fprintf(fo, "%s:%c_Ht H d %s:%c_T T f ", stank.Name, c, stank.Name, c)
				fmt.Fprintf(fo, "%s:%c_ttn h d %s:%c_Tn t f %s:%c_ttm h d %s:%c_Tm t f\n",
					stank.Name, c, stank.Name, c, stank.Name, c, stank.Name, c)
				fmt.Fprintf(fo, "%s:%c_Hh H d %s:%c_Qh Q f %s:%c_Hc H d %s:%c_Qc Q f\n",
					stank.Name, c, stank.Name, c, stank.Name, c, stank.Name, c)
				fmt.Fprintf(fo, "%s:%c_th h d %s:%c_qh q f %s:%c_tc h d %s:%c_qc q f\n",
					stank.Name, c, stank.Name, c, stank.Name, c, stank.Name, c)
			}
			fmt.Fprintf(fo, "%s_Qls Q f %s_Qst Q f\n\n", stank.Name, stank.Name)
		}

	default:
		for _, stank := range Stank {
			S := &stank.Stkdy[0]

			fmt.Fprintf(fo, "%.1f\n", S.Tsdy.M)
			for j := 0; j < stank.Nin; j++ {
				S := &stank.Stkdy[j]

				fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ",
					S.Tidy.Hrs, S.Tidy.M,
					S.Tidy.Mntime, S.Tidy.Mn,
					S.Tidy.Mxtime, S.Tidy.Mx)

				fmt.Fprintf(fo, "%1d %3.1f ", S.Qdy.Hhr, S.Qdy.H)
				fmt.Fprintf(fo, "%1d %3.1f ", S.Qdy.Chr, S.Qdy.C)
				fmt.Fprintf(fo, "%1d %2.0f ", S.Qdy.Hmxtime, S.Qdy.Hmx)
				fmt.Fprintf(fo, "%1d %2.0f ", S.Qdy.Cmxtime, S.Qdy.Cmx)
			}
			fmt.Fprintf(fo, " %3.1f %3.1f\n",
				stank.Qlossdy*Cff_kWh, stank.Qstody*Cff_kWh)
		}
	}
}

// 月集計の出力
func stankmonprt(fo io.Writer, id int, Stank []*STANK) {
	switch id {
	case 0:
		if len(Stank) > 0 {
			fmt.Fprintf(fo, "%s %d\n", STANK_TYPE, len(Stank))
		}

		for _, stank := range Stank {
			fmt.Fprintf(fo, "%s:%d", stank.Name, stank.Nin)

			for i := 0; i < stank.Nin; i++ {
				fmt.Fprintf(fo, "%c", stank.Cmp.Idi[i])
			}

			fmt.Fprintf(fo, " 1 %d\n", stank.Nin*14+2+1)
		}

	case 1:
		for _, stank := range Stank {
			fmt.Fprintf(fo, "%s_Ts t f \n", stank.Name)

			for i := 0; i < stank.Nin; i++ {
				c := stank.Cmp.Idi[i]
				fmt.Fprintf(fo, "%s:%c_Ht H d %s:%c_T T f ", stank.Name, c, stank.Name, c)
				fmt.Fprintf(fo, "%s:%c_ttn h d %s:%c_Tn t f %s:%c_ttm h d %s:%c_Tm t f\n",
					stank.Name, c, stank.Name, c, stank.Name, c, stank.Name, c)
				fmt.Fprintf(fo, "%s:%c_Hh H d %s:%c_Qh Q f %s:%c_Hc H d %s:%c_Qc Q f\n",
					stank.Name, c, stank.Name, c, stank.Name, c, stank.Name, c)
				fmt.Fprintf(fo, "%s:%c_th h d %s:%c_qh q f %s:%c_tc h d %s:%c_qc q f\n",
					stank.Name, c, stank.Name, c, stank.Name, c, stank.Name, c)
			}
			fmt.Fprintf(fo, "%s_Qls Q f %s_Qst Q f\n\n", stank.Name, stank.Name)
		}

	default:
		for _, stank := range Stank {
			S := &stank.Mstkdy[0]

			fmt.Fprintf(fo, "%.1f\n", S.Tsdy.M)
			for j := 0; j < stank.Nin; j++ {
				S := &stank.Mstkdy[j]

				fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ",
					S.Tidy.Hrs, S.Tidy.M,
					S.Tidy.Mntime, S.Tidy.Mn,
					S.Tidy.Mxtime, S.Tidy.Mx)

				fmt.Fprintf(fo, "%1d %3.1f ", S.Qdy.Hhr, S.Qdy.H)
				fmt.Fprintf(fo, "%1d %3.1f ", S.Qdy.Chr, S.Qdy.C)
				fmt.Fprintf(fo, "%1d %2.0f ", S.Qdy.Hmxtime, S.Qdy.Hmx)
				fmt.Fprintf(fo, "%1d %2.0f ", S.Qdy.Cmxtime, S.Qdy.Cmx)
			}
			fmt.Fprintf(fo, " %3.1f %3.1f\n",
				stank.MQlossdy*Cff_kWh, stank.MQstody*Cff_kWh)
		}
	}
}

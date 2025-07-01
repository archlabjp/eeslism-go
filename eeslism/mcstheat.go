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

package eeslism

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

/*  mcstheat.c  */
/*
Stheatdata (Sensible Heat Storage Data Input)

この関数は、顕熱蓄熱式暖房器の各種仕様（定格熱量、効率、熱容量、熱通過率など）を読み込み、
対応する構造体に格納します。
これらのデータは、顕熱蓄熱システムにおける機器の性能評価、
熱負荷平準化、およびエネルギー消費量予測に不可欠です。

建築環境工学的な観点:
- **顕熱蓄熱の原理**: 顕熱蓄熱式暖房器は、
  電気ヒーターなどで熱を発生させ、その熱を蓄熱材（レンガ、コンクリートなど）に顕熱として蓄え、
  必要な時に放熱することで暖房を行うシステムです。
  電力料金の安い夜間電力などを利用して蓄熱し、
  昼間のピークカットや熱負荷平準化に貢献します。
- **定格熱量 (Q)**:
  蓄熱式暖房器が供給できる最大の熱量を示します。
  建物の暖房負荷に対して十分な能力があるか、
  あるいは複数台の機器を組み合わせる必要があるかを判断する際に用いられます。
- **効率 (Eff)**:
  投入されたエネルギーに対して、どれだけの熱エネルギーを蓄熱・供給できるかを示す指標です。
  効率が高いほど、エネルギーの無駄が少なく、省エネルギーに貢献します。
- **熱容量 (Hcap)**:
  蓄熱材が蓄えることができる熱量を示します。
  `Hcap`が大きいほど、より多くの熱を蓄えることができ、
  長時間の放熱や、大きな熱負荷変動に対応できます。
- **熱通過率 (KA)**:
  蓄熱材から周囲への熱損失を表す係数です。
  `KA`が大きいと、蓄熱された熱が有効に利用されず、エネルギーの無駄につながります。
- **PCMの利用 (PCMName)**:
  `PCMName`が設定されている場合、
  相変化材料（PCM）を蓄熱材として利用していることを示唆します。
  PCMは、潜熱を利用して顕熱蓄熱よりも大きな熱量を蓄えることができ、
  よりコンパクトな蓄熱システムを実現できます。

この関数は、顕熱蓄熱式暖房器の性能をモデル化し、
熱負荷平準化、エネルギー消費量予測、
および省エネルギー対策の検討を行うための重要なデータ入力機能を提供します。
*/
func Stheatdata(s string, stheatca *STHEATCA) int {
	var id int

	if st := strings.IndexRune(s, '='); st == -1 {
		stheatca.Name = s
		stheatca.Eff = -999.0
		stheatca.Q = -999.0
		stheatca.Hcap = -999.0
		stheatca.KA = -999.0
	} else {
		sval := s[st+1:]

		if s == "PCM" {
			stheatca.PCMName = sval
		} else {
			dt, err := strconv.ParseFloat(sval, 64)
			if err != nil {
				panic(err)
			}

			if s == "Q" {
				stheatca.Q = dt
			} else if s == "KA" {
				stheatca.KA = dt
			} else if s == "eff" {
				stheatca.Eff = dt
			} else if s == "Hcap" {
				stheatca.Hcap = dt
			} else {
				id = 1
			}
		}
	}
	return id
}

/* --------------------------- */

/*
Stheatint (Sensible Heat Storage Initialization)

この関数は、顕熱蓄熱式暖房器のシミュレーションに必要な初期設定と、
各種パラメータの妥当性チェックを行います。
特に、蓄熱材の初期温度や、PCM（相変化材料）が利用されている場合の関連付けを行います。

建築環境工学的な観点:
- **初期温度の設定**: シミュレーション開始時の蓄熱材の温度（`stheat.Tsold`）を初期化します。
  この初期値は、シミュレーションの収束性や、
  初期段階での蓄熱式暖房器の挙動に影響を与える可能性があります。
- **PCMの関連付け**: `stheat.Cat.PCMName`が設定されている場合、
  対応するPCMのデータ（`_PCM`）を検索し、`stheat.Pcm`にポインターを設定します。
  これにより、PCMの潜熱蓄熱効果をモデルに組み込むことができます。
  PCMは、顕熱蓄熱よりも大きな熱量を蓄えることができ、
  よりコンパクトな蓄熱システムを実現できます。
- **周囲環境への熱損失**: `stheat.Tenv`は蓄熱式暖房器周囲の環境温度へのポインターであり、
  熱損失計算に用いられます。
  `stheat.Room`は、設置室の温度を周囲環境温度として利用する場合に設定されます。
- **パラメータの妥当性チェック**: `if st.Q < 0.0` のようなエラーチェックは、
  入力されたパラメータが物理的に妥当な範囲内にあるかを確認するために重要です。
  不適切なパラメータは、シミュレーション結果の信頼性を損なう可能性があります。

この関数は、顕熱蓄熱式暖房器の熱的挙動を詳細にモデル化し、
熱負荷平準化、エネルギー消費量予測、
および省エネルギー対策の検討を行うための重要な初期設定機能を提供します。
*/
func Stheatint(_stheat []*STHEAT, Simc *SIMCONTL, Compnt []*COMPNT, Wd *WDAT, _PCM []*PCM) {
	for i := range _stheat {
		stheat := _stheat[i]
		if stheat.Cmp.Envname != "" {
			stheat.Tenv = envptr(stheat.Cmp.Envname, Simc, Compnt, Wd, nil)
		} else {
			stheat.Room = roomptr(stheat.Cmp.Roomname, Compnt)
		}

		if stheat.Cat.PCMName != "" {
			for j := range _PCM {
				if stheat.Cat.PCMName == _PCM[j].Name {
					stheat.Pcm = _PCM[j]
				}
			}
			if stheat.Pcm == nil {
				Err := fmt.Sprintf("STHEAT %s のPCM=%sが見つかりません", stheat.Name, stheat.Cat.PCMName)
				Eprint(Err, "<Stheatint>")
				os.Exit(1)
			}
		}

		st := stheat.Cat

		if st.Q < 0.0 {
			Err := fmt.Sprintf("Name=%s  Q=%.4g", stheat.Name, st.Q)
			Eprint("Stheatinit", Err)
		}
		if stheat.Pcm == nil && st.Hcap < 0.0 {
			Err := fmt.Sprintf("Name=%s  Hcap=%.4g", stheat.Name, st.Hcap)
			Eprint("Stheatinit", Err)
		}
		if st.KA < 0.0 {
			Err := fmt.Sprintf("Name=%s  KA=%.4g", stheat.Name, st.KA)
			Eprint("Stheatinit", Err)
		}
		if st.Eff < 0.0 {
			Err := fmt.Sprintf("Name=%s  eff=%.4g", stheat.Name, st.Eff)
			Eprint("Stheatinit", Err)
		}

		var err error
		stheat.Tsold, err = strconv.ParseFloat(stheat.Cmp.Tparm, 64)
		if err != nil {
			panic(err)
		}

		// 内臓PCMの質量
		stheat.MPCM = stheat.Cmp.MPCM
	}
}

/* --------------------------- */
/*  特性式の係数  */

//
//    +--------+ --> [OUT 1]
//    | STHEAT |
//    +--------+ --> [OUT 2]
//


/*
Stheatcfv (Sensible Heat Storage Characteristic Function Value Calculation)

この関数は、顕熱蓄熱式暖房器の運転特性を評価し、
熱媒の熱容量流量、熱容量、熱通過率、そして電力消費量などを考慮した係数を計算します。
これは、蓄熱式暖房器が熱を蓄え、放熱するプロセスをモデル化するために不可欠です。

建築環境工学的な観点:
- **熱容量 (stheat.Hcap)**:
  蓄熱材の熱容量は、蓄熱式暖房器が蓄えることができる熱量を示します。
  PCMが利用されている場合は、その潜熱蓄熱効果も熱容量に加算されます。
- **熱媒の熱容量流量 (stheat.CG)**:
  蓄熱式暖房器を通過する空気の熱容量流量は、
  放熱能力に影響します。
  `Spcheat(Eo1.Fluid) * Eo1.G` のように、空気の比熱と質量流量から計算されます。
- **熱通過率 (KA)**:
  蓄熱材から周囲への熱損失を表す係数です。
- **電力消費量 (stheat.E)**:
  蓄熱式暖房器が熱を蓄えるために消費する電力です。
  `stheat.Cat.Q`（定格熱量）が設定されている場合、
  その値が電力消費量として扱われます。
- **熱収支方程式の構築**: `d = stheat.Hcap/DTM + eff*cG + KA` のように、
  蓄熱材の熱容量、熱媒との熱交換、周囲への熱損失を考慮した熱収支方程式の係数を計算します。
  これにより、蓄熱材の温度変化を予測できます。
- **出口温度の係数設定**: 計算された蓄熱式暖房器の特性に基づいて、
  出口空気温度（`Eo1`）に関する係数（`Coeffo`, `Co`, `Coeffin`）を設定します。
  これらの係数は、システム全体の熱収支方程式に組み込まれ、
  室内の温度を予測するために用いられます。

この関数は、顕熱蓄熱式暖房器の熱的挙動を詳細にモデル化し、
熱負荷平準化、エネルギー消費量予測、
および省エネルギー対策の検討を行うための重要な役割を果たします。
*/
func Stheatcfv(_stheat []*STHEAT) {
	for i := range _stheat {
		stheat := _stheat[i]

		// 作用温度 ?
		var Te float64
		if stheat.Cmp.Envname != "" {
			Te = *(stheat.Tenv)
		} else {
			Te = stheat.Room.Tot
		}

		Eo1 := stheat.Cmp.Elouts[0]
		eff := stheat.Cat.Eff
		stheat.CG = Spcheat(Eo1.Fluid) * Eo1.G
		KA := stheat.Cat.KA
		Tsold := stheat.Tsold
		pcm := stheat.Pcm
		if pcm != nil {
			//NOTE: FNPCMState のシグネチャがヘッダと一致しない。。。
			// stheat.Hcap = stheat.MPCM *
			// 	FNPCMState(pcm.Cros, pcm.Crol, pcm.Ql, pcm.Ts, pcm.Tl, Tsold, nil)
			panic("Cannot call FNPCMState")
		} else {
			stheat.Hcap = stheat.Cat.Hcap
		}
		cG := stheat.CG

		d := stheat.Hcap/DTM + eff*cG + KA
		if stheat.Cmp.Control != OFF_SW {
			stheat.E = stheat.Cat.Q
		} else {
			stheat.E = 0.0
		}

		//  空気が流れていれば出入口温度の関係式係数を作成する
		if Eo1.Control != OFF_SW {
			Eo1.Coeffo = 1.0
			Eo1.Co = eff * (stheat.Hcap/DTM*Tsold + KA*Te + stheat.E) / d
			Eo1.Coeffin[0] = eff - 1.0 - eff*eff*cG/d

			Eo2 := stheat.Cmp.Elouts[1]
			Eo2.Coeffo = 1.0
			Eo2.Co = 0.0
			Eo2.Coeffin[0] = -1.0
		} else {
			Eo1.Coeffo = 1.0
			Eo1.Co = 0.0
			Eo1.Coeffin[0] = -1.0

			Eo2 := stheat.Cmp.Elouts[1]
			Eo2.Coeffo = 1.0
			Eo2.Co = 0.0
			Eo2.Coeffin[0] = -1.0
		}
	}
}

/* --------------------------- */

/*
Stheatene (Sensible Heat Storage Energy Calculation)

この関数は、顕熱蓄熱式暖房器が空気側に供給する熱量、
周囲への熱損失量、および蓄熱量の変化を計算します。
これは、蓄熱式暖房器のエネルギー収支を詳細に分析し、
熱負荷平準化の効果やエネルギー消費量を評価する上で不可欠です。

建築環境工学的な観点:
- **供給熱量 (stheat.Q)**:
  蓄熱式暖房器が空気側に供給する熱量です。
  `stheat.CG * (stheat.Tout - stheat.Tin)` のように、
  空気側の熱容量流量と入口・出口空気温度差から計算されます。
  これは、蓄熱式暖房器が暖房負荷にどれだけ貢献できるかを示します。
- **熱損失量 (stheat.Qls)**:
  蓄熱式暖房器から周囲環境へ逃げる熱量です。
  `stheat.Cat.KA * (Te - stheat.Ts)` のように、
  熱通過率と蓄熱材温度および周囲環境温度の差から計算されます。
  熱損失は、蓄熱式暖房器のエネルギー効率に影響を与え、
  特に長時間の蓄熱では無視できない要因となります。
- **蓄熱量 (stheat.Qsto)**:
  蓄熱式暖房器内部に蓄えられた熱量の変化量です。
  `stheat.Hcap / DTM * (stheat.Ts - stheat.Tsold)` のように、
  熱容量と温度変化から計算されます。
  これは、蓄熱式暖房器が熱をどれだけ貯蔵できたかを示します。
- **蓄熱材温度の更新 (stheat.Tsold = stheat.Ts)**:
  計算された蓄熱材の温度（`stheat.Ts`）を、
  次の時間ステップの計算のために`stheat.Tsold`に更新します。
  これにより、蓄熱材の熱的履歴が考慮され、
  より正確な動的熱応答のシミュレーションが可能となります。
- **設置室内部発熱への影響**: `if stheat.Room != nil { stheat.Room.Qeqp += (-stheat.Qls) }` のように、
  蓄熱式暖房器からの熱損失が設置室の内部発熱として計上されることで、
  建物全体の熱収支モデルに組み込まれます。

この関数は、顕熱蓄熱式暖房器のエネルギー収支を詳細に分析し、
熱負荷平準化、エネルギー消費量予測、
および省エネルギー対策の検討を行うための重要な役割を果たします。
*/
func Stheatene(_stheat []*STHEAT) {
	var elo *ELOUT
	var cat *STHEATCA
	var Te float64

	for i := range _stheat {
		stheat := _stheat[i]
		elo = stheat.Cmp.Elouts[0]
		stheat.Tin = elo.Elins[0].Sysvin

		cat = stheat.Cat

		if stheat.Cmp.Envname != "" {
			Te = *(stheat.Tenv)
		} else {
			Te = stheat.Room.Tot
		}

		stheat.Tout = elo.Sysv
		stheat.Ts = (stheat.Hcap/DTM*stheat.Tsold +
			cat.Eff*stheat.CG*stheat.Tin +
			cat.KA*Te + stheat.E) /
			(stheat.Hcap/DTM + cat.Eff*stheat.CG + cat.KA)

		stheat.Q = stheat.CG * (stheat.Tout - stheat.Tin)

		stheat.Qls = stheat.Cat.KA * (Te - stheat.Ts)

		stheat.Qsto = stheat.Hcap / DTM * (stheat.Ts - stheat.Tsold)

		stheat.Tsold = stheat.Ts

		if stheat.Room != nil {
			stheat.Room.Qeqp += (-stheat.Qls)
		}
	}
}

/*
stheatvptr (Sensible Heat Storage Internal Variable Pointer Setting)

この関数は、顕熱蓄熱式暖房器の制御で使用される内部変数（蓄熱材温度、制御状態など）へのポインターを設定します。
これにより、蓄熱式暖房器の運転を特定の目標値に追従させる制御をモデル化できます。

建築環境工学的な観点:
- **蓄熱材温度の制御**: 蓄熱式暖房器は、
  蓄熱材の温度を監視し、それに基づいて電力供給を制御することで、
  目標とする蓄熱量や放熱量を維持します。
  `key[1]`が`"Ts"`の場合、蓄熱材温度（`Stheat.Tsold`）を制御対象とすることを意味します。
  `vptr.Ptr`は、この変数へのポインターを設定し、
  `vptr.Type = VAL_CTYPE`は、そのポインターが制御値であることを示します。
- **運転制御のポインター**: `key[1]`が`"control"`の場合、
  蓄熱式暖房器の運転制御状態（`Stheat.Cmp.Control`）へのポインターを設定します。
  これにより、外部からの信号やスケジュールに基づいて、
  蓄熱式暖房器の運転をON/OFFしたり、運転モードを切り替えたりする制御をモデル化できます。
- **フィードバック制御の基礎**: このポインター設定は、
  蓄熱式暖房器のフィードバック制御の基礎となります。
  シミュレーションの各時間ステップで、
  現在の蓄熱材温度や運転状態と目標値を比較し、その差に基づいて運転を調整します。
  これにより、蓄熱式暖房器の効率的な運転や、
  熱負荷平準化の効果を最大化することができます。

この関数は、顕熱蓄熱式暖房器の制御ロジックをモデル化し、
熱負荷平準化、エネルギー消費量予測、
および省エネルギー対策の検討を行うための重要な役割を果たします。
*/
func stheatvptr(key []string, Stheat *STHEAT) (VPTR, VPTR, error) {
	var err error
	var vptr, vpath VPTR

	if key[1] == "Ts" {
		vptr = VPTR{
			Ptr:  &Stheat.Tsold,
			Type: VAL_CTYPE,
		}
	} else if key[1] == "control" {
		vpath = VPTR{
			Type: 's',
			Ptr:  Stheat,
		}
		vptr = VPTR{
			Ptr:  &Stheat.Cmp.Control,
			Type: SW_CTYPE,
		}
	} else {
		err = errors.New("'Ts' or 'control' is expected")
	}

	return vptr, vpath, err
}

/* ---------------------------*/

func stheatprint(fo io.Writer, id int, stheat []*STHEAT) {
	switch id {
	case 0:
		if len(stheat) > 0 {
			fmt.Fprintf(fo, "%s %d\n", STHEAT_TYPE, len(stheat))
		}
		for i := range stheat {
			fmt.Fprintf(fo, " %s 1 11\n", stheat[i].Name)
		}
	case 1:
		for i := range stheat {
			fmt.Fprintf(fo, "%s_c c c %s_G m f %s_Ts t f %s_Ti t f %s_To t f %s_Q q f %s_Qsto q f ",
				stheat[i].Name, stheat[i].Name, stheat[i].Name, stheat[i].Name, stheat[i].Name, stheat[i].Name, stheat[i].Name)
			fmt.Fprintf(fo, "%s_Qls q f %s_Ec c c %s_E e f ",
				stheat[i].Name, stheat[i].Name, stheat[i].Name)
			fmt.Fprintf(fo, "%s_Hcap q f\n", stheat[i].Name)
		}
	default:
		for i := range stheat {
			fmt.Fprintf(fo, "%c %6.4g %4.1f %4.1f %4.1f %2.0f %.4g ",
				stheat[i].Cmp.Elouts[0].Control, stheat[i].Cmp.Elouts[0].G,
				stheat[i].Ts,
				stheat[i].Tin, stheat[i].Tout, stheat[i].Q, stheat[i].Qsto)
			fmt.Fprintf(fo, "%.4g %c %2.0f ",
				stheat[i].Qls, stheat[i].Cmp.Control, stheat[i].E)
			fmt.Fprintf(fo, "%.0f\n", stheat[i].Hcap)
		}
	}
}

/* --------------------------- */
/* 日積算値に関する処理 */

/*
stheatdyint (Sensible Heat Storage Daily Integration Initialization)

この関数は、顕熱蓄熱式暖房器の日積算値（日ごとの入口・出口空気温度、蓄熱材温度、
供給熱量、エネルギー消費量、熱損失量、蓄熱量など）をリセットします。
これは、日単位での蓄熱式暖房器の運転状況や熱交換量を集計し、
空調システムの性能評価やエネルギー消費量の分析に用いられます。

建築環境工学的な観点:
- **日単位の性能評価**: 顕熱蓄熱式暖房器の運転状況は、日中の熱負荷変動に応じて大きく変化します。
  日積算値を集計することで、日ごとの熱負荷平準化の効果、
  蓄熱式暖房器の稼働時間、部分負荷運転の割合などを把握できます。
  これにより、特定の日の蓄熱特性を分析したり、
  蓄熱式暖房器の運転効率を日単位で評価したりすることが可能になります。
- **運用改善の指標**: 日積算データは、蓄熱システムの運用改善のための重要な指標となります。
  例えば、外気温度や熱負荷などの気象条件と蓄熱式暖房器の熱交換量の関係を分析したり、
  設定温度や換気量などの運用条件が蓄熱式暖房器の性能に与える影響を評価したりすることで、
  より効率的な運転方法を見つけることができます。
- **データ集計の準備**: この関数は、新しい日の集計を開始する前に、
  前日のデータをクリアする役割を担います。
  `svdyint`や`qdyint`、`edyint`といった関数は、
  それぞれ温度、熱量、エネルギーなどの日積算値をリセットするためのものです。

この関数は、顕熱蓄熱式暖房器の運転状況と熱交換量を日単位で詳細に分析し、
蓄熱システムの運用改善や省エネルギー対策の効果評価を行うための基礎的な役割を果たします。
*/
func stheatdyint(stheat []*STHEAT) {
	for i := range stheat {
		stheat[i].Qlossdy = 0.0
		stheat[i].Qstody = 0.0

		svdyint(&stheat[i].Tidy)
		svdyint(&stheat[i].Tsdy)
		svdyint(&stheat[i].Tody)
		qdyint(&stheat[i].Qdy)
		edyint(&stheat[i].Edy)
	}
}

/*
stheatmonint (Sensible Heat Storage Monthly Integration Initialization)

この関数は、顕熱蓄熱式暖房器の月積算値（月ごとの入口・出口空気温度、蓄熱材温度、
供給熱量、エネルギー消費量、熱損失量、蓄熱量など）をリセットします。
これは、月単位での蓄熱式暖房器の運転状況や熱交換量を集計し、
空調システムの性能評価やエネルギー消費量の分析に用いられます。

建築環境工学的な観点:
- **月単位の性能評価**: 顕熱蓄熱式暖房器の運転状況は、月単位で変動します。
  月積算値を集計することで、月ごとの熱負荷平準化の効果、
  蓄熱式暖房器の稼働時間、部分負荷運転の割合などを把握できます。
  これにより、特定の月の蓄熱特性を分析したり、
  蓄熱式暖房器の運転効率を月単位で評価したりすることが可能になります。
- **運用改善の指標**: 月積算データは、蓄熱システムの運用改善のための重要な指標となります。
  例えば、季節ごとの蓄熱特性の傾向を把握したり、
  月ごとの気象条件と蓄熱式暖房器の熱交換量の関係を分析したりすることで、
  より効率的な運転方法を見つけることができます。
- **データ集計の準備**: この関数は、新しい月の集計を開始する前に、
  前月のデータをクリアする役割を担います。
  `svdyint`や`qdyint`、`edyint`といった関数は、
  それぞれ温度、熱量、エネルギーなどの月積算値をリセットするためのものです。

この関数は、顕熱蓄熱式暖房器の運転状況と熱交換量を月単位で詳細に分析し、
蓄熱システムの運用改善や省エネルギー対策の効果評価を行うための基礎的な役割を果たします。
*/
func stheatmonint(stheat []*STHEAT) {
	for i := range stheat {
		stheat[i].MQlossdy = 0.0
		stheat[i].MQstody = 0.0

		svdyint(&stheat[i].MTidy)
		svdyint(&stheat[i].MTsdy)
		svdyint(&stheat[i].MTody)
		qdyint(&stheat[i].MQdy)
		edyint(&stheat[i].MEdy)
	}
}

/*
stheatday (Sensible Heat Storage Daily and Monthly Data Aggregation)

この関数は、顕熱蓄熱式暖房器の運転データ（入口・出口空気温度、蓄熱材温度、
供給熱量、エネルギー消費量、熱損失量、蓄熱量など）を、
日単位および月単位で集計します。
これにより、蓄熱式暖房器の性能評価やエネルギー消費量の分析が可能になります。

建築環境工学的な観点:
- **日次集計 (svdaysum, qdaysum, edaysum)**:
  日次集計は、蓄熱式暖房器の運転状況を日単位で詳細に把握するために重要です。
  例えば、特定の日の熱負荷変動に対する蓄熱式暖房器の応答、
  あるいは日中のピーク負荷時の熱交換量などを分析できます。
  これにより、日ごとの運用改善点を見つけ出すことが可能になります。
- **月次集計 (svmonsum, qmonsum, emonsum)**:
  月次集計は、季節ごとの熱負荷変動や、
  蓄熱式暖房器の年間を通じた熱交換量の傾向を把握するために重要です。
  これにより、年間を通じた省エネルギー対策の効果を評価したり、
  熱交換量の予測精度を向上させたりすることが可能になります。
- **月・時刻のクロス集計 (emtsum)**:
  `MtEdy`のようなクロス集計は、
  特定の月における時間帯ごとのエネルギー消費量を分析するために用いられます。
  これにより、例えば、冬季の朝の時間帯に暖房負荷が集中する傾向があるか、
  あるいは夏季の夜間に給湯負荷が高いかなどを詳細に把握できます。
  これは、デマンドサイドマネジメントや、
  エネルギー供給計画を最適化する上で非常に有用な情報となります。
- **データ分析の基礎**: この関数で集計されるデータは、
  蓄熱式暖房器の性能評価、熱交換量のベンチマーキング、
  省エネルギー対策の効果検証、そして運用改善のための意思決定の基礎となります。

この関数は、顕熱蓄熱式暖房器の運転状況と熱交換量を多角的に分析し、
蓄熱システムの運用改善や省エネルギー対策の効果評価を行うための重要なデータ集計機能を提供します。
*/
func stheatday(Mon, Day, ttmm int, stheat []*STHEAT, Nday, SimDayend int) {
	Mo := Mon - 1
	tt := ConvertHour(ttmm)

	for i := range stheat {
		// 日集計
		stheat[i].Qlossdy += stheat[i].Qls
		stheat[i].Qstody += stheat[i].Qsto
		svdaysum(int64(ttmm), stheat[i].Cmp.Control, stheat[i].Tin, &stheat[i].Tidy)
		svdaysum(int64(ttmm), stheat[i].Cmp.Control, stheat[i].Tout, &stheat[i].Tody)
		svdaysum(int64(ttmm), stheat[i].Cmp.Control, stheat[i].Ts, &stheat[i].Tsdy)
		qdaysum(int64(ttmm), stheat[i].Cmp.Control, stheat[i].Q, &stheat[i].Qdy)
		edaysum(ttmm, stheat[i].Cmp.Control, stheat[i].E, &stheat[i].Edy)

		// 月集計
		stheat[i].MQlossdy += stheat[i].Qls
		stheat[i].MQstody += stheat[i].Qsto
		svmonsum(Mon, Day, ttmm, stheat[i].Cmp.Control, stheat[i].Tin, &stheat[i].MTidy, Nday, SimDayend)
		svmonsum(Mon, Day, ttmm, stheat[i].Cmp.Control, stheat[i].Tout, &stheat[i].MTody, Nday, SimDayend)
		svmonsum(Mon, Day, ttmm, stheat[i].Cmp.Control, stheat[i].Ts, &stheat[i].MTsdy, Nday, SimDayend)
		qmonsum(Mon, Day, ttmm, stheat[i].Cmp.Control, stheat[i].Q, &stheat[i].MQdy, Nday, SimDayend)
		emonsum(Mon, Day, ttmm, stheat[i].Cmp.Control, stheat[i].E, &stheat[i].MEdy, Nday, SimDayend)

		// 月・時刻のクロス集計
		emtsum(Mon, Day, ttmm, stheat[i].Cmp.Control, stheat[i].E, &stheat[i].MtEdy[Mo][tt])
	}
}

func stheatdyprt(fo io.Writer, id int, stheat []*STHEAT) {
	switch id {
	case 0:
		if len(stheat) > 0 {
			fmt.Fprintf(fo, "%s %d\n", STHEAT_TYPE, len(stheat))
		}
		for i := range stheat {
			fmt.Fprintf(fo, " %s 1 32\n", stheat[i].Name)
		}
	case 1:
		for i := range stheat {
			fmt.Fprintf(fo, "%s_Ht H d %s_Ti T f ", stheat[i].Name, stheat[i].Name)
			fmt.Fprintf(fo, "%s_ttn h d %s_Tin t f %s_ttm h d %s_Tim t f\n",
				stheat[i].Name, stheat[i].Name, stheat[i].Name, stheat[i].Name)
			fmt.Fprintf(fo, "%s_Ht H d %s_To T f ", stheat[i].Name, stheat[i].Name)
			fmt.Fprintf(fo, "%s_ttn h d %s_Ton t f %s_ttm h d %s_Tom t f\n",
				stheat[i].Name, stheat[i].Name, stheat[i].Name, stheat[i].Name)
			fmt.Fprintf(fo, "%s_Ht H d %s_Ts T f ", stheat[i].Name, stheat[i].Name)
			fmt.Fprintf(fo, "%s_ttn h d %s_Tsn t f %s_ttm h d %s_Tsm t f\n",
				stheat[i].Name, stheat[i].Name, stheat[i].Name, stheat[i].Name)
			fmt.Fprintf(fo, "%s_Hh H d %s_Qh Q f %s_Hc H d %s_Qc Q f\n",
				stheat[i].Name, stheat[i].Name, stheat[i].Name, stheat[i].Name)
			fmt.Fprintf(fo, "%s_th h d %s_qh q f %s_tc h d %s_qc q f\n",
				stheat[i].Name, stheat[i].Name, stheat[i].Name, stheat[i].Name)
			fmt.Fprintf(fo, "%s_He H d %s_E E f %s_te h d %s_Em e f\n",
				stheat[i].Name, stheat[i].Name, stheat[i].Name, stheat[i].Name)
			fmt.Fprintf(fo, "%s_Qls Q f %s_Qst Q f\n\n",
				stheat[i].Name, stheat[i].Name)
		}
	default:
		for i := range stheat {
			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ",
				stheat[i].Tidy.Hrs, stheat[i].Tidy.M,
				stheat[i].Tidy.Mntime, stheat[i].Tidy.Mn,
				stheat[i].Tidy.Mxtime, stheat[i].Tidy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ",
				stheat[i].Tody.Hrs, stheat[i].Tody.M,
				stheat[i].Tody.Mntime, stheat[i].Tody.Mn,
				stheat[i].Tody.Mxtime, stheat[i].Tody.Mx)
			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ",
				stheat[i].Tsdy.Hrs, stheat[i].Tsdy.M,
				stheat[i].Tsdy.Mntime, stheat[i].Tsdy.Mn,
				stheat[i].Tsdy.Mxtime, stheat[i].Tsdy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f ", stheat[i].Qdy.Hhr, stheat[i].Qdy.H)
			fmt.Fprintf(fo, "%1d %3.1f ", stheat[i].Qdy.Chr, stheat[i].Qdy.C)
			fmt.Fprintf(fo, "%1d %2.0f ", stheat[i].Qdy.Hmxtime, stheat[i].Qdy.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f ", stheat[i].Qdy.Cmxtime, stheat[i].Qdy.Cmx)
			fmt.Fprintf(fo, "%1d %3.1f ", stheat[i].Edy.Hrs, stheat[i].Edy.D)
			fmt.Fprintf(fo, "%1d %2.0f ", stheat[i].Edy.Mxtime, stheat[i].Edy.Mx)
			fmt.Fprintf(fo, " %3.1f %3.1f\n",
				stheat[i].Qlossdy*Cff_kWh, stheat[i].Qstody*Cff_kWh)
		}
	}
}

func stheatmonprt(fo io.Writer, id int, stheat []*STHEAT) {
	switch id {
	case 0:
		if len(stheat) > 0 {
			fmt.Fprintf(fo, "%s %d\n", STHEAT_TYPE, len(stheat))
		}
		for i := range stheat {
			fmt.Fprintf(fo, " %s 1 32\n", stheat[i].Name)
		}
	case 1:
		for i := range stheat {
			fmt.Fprintf(fo, "%s_Ht H d %s_Ti T f ", stheat[i].Name, stheat[i].Name)
			fmt.Fprintf(fo, "%s_ttn h d %s_Tin t f %s_ttm h d %s_Tim t f\n",
				stheat[i].Name, stheat[i].Name, stheat[i].Name, stheat[i].Name)
			fmt.Fprintf(fo, "%s_Ht H d %s_To T f ", stheat[i].Name, stheat[i].Name)
			fmt.Fprintf(fo, "%s_ttn h d %s_Ton t f %s_ttm h d %s_Tom t f\n",
				stheat[i].Name, stheat[i].Name, stheat[i].Name, stheat[i].Name)
			fmt.Fprintf(fo, "%s_Ht H d %s_Ts T f ", stheat[i].Name, stheat[i].Name)
			fmt.Fprintf(fo, "%s_ttn h d %s_Tsn t f %s_ttm h d %s_Tsm t f\n",
				stheat[i].Name, stheat[i].Name, stheat[i].Name, stheat[i].Name)
			fmt.Fprintf(fo, "%s_Hh H d %s_Qh Q f %s_Hc H d %s_Qc Q f\n",
				stheat[i].Name, stheat[i].Name, stheat[i].Name, stheat[i].Name)
			fmt.Fprintf(fo, "%s_th h d %s_qh q f %s_tc h d %s_qc q f\n",
				stheat[i].Name, stheat[i].Name, stheat[i].Name, stheat[i].Name)
			fmt.Fprintf(fo, "%s_He H d %s_E E f %s_te h d %s_Em e f\n",
				stheat[i].Name, stheat[i].Name, stheat[i].Name, stheat[i].Name)
			fmt.Fprintf(fo, "%s_Qls Q f %s_Qst Q f\n\n",
				stheat[i].Name, stheat[i].Name)
		}
	default:
		for i := range stheat {
			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ",
				stheat[i].MTidy.Hrs, stheat[i].MTidy.M,
				stheat[i].MTidy.Mntime, stheat[i].MTidy.Mn,
				stheat[i].MTidy.Mxtime, stheat[i].MTidy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ",
				stheat[i].MTody.Hrs, stheat[i].MTody.M,
				stheat[i].MTody.Mntime, stheat[i].MTody.Mn,
				stheat[i].MTody.Mxtime, stheat[i].MTody.Mx)
			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ",
				stheat[i].MTsdy.Hrs, stheat[i].MTsdy.M,
				stheat[i].MTsdy.Mntime, stheat[i].MTsdy.Mn,
				stheat[i].MTsdy.Mxtime, stheat[i].MTsdy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f ", stheat[i].MQdy.Hhr, stheat[i].MQdy.H)
			fmt.Fprintf(fo, "%1d %3.1f ", stheat[i].MQdy.Chr, stheat[i].MQdy.C)
			fmt.Fprintf(fo, "%1d %2.0f ", stheat[i].MQdy.Hmxtime, stheat[i].MQdy.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f ", stheat[i].MQdy.Cmxtime, stheat[i].MQdy.Cmx)
			fmt.Fprintf(fo, "%1d %3.1f ", stheat[i].MEdy.Hrs, stheat[i].MEdy.D)
			fmt.Fprintf(fo, "%1d %2.0f ", stheat[i].MEdy.Mxtime, stheat[i].MEdy.Mx)
			fmt.Fprintf(fo, " %3.1f %3.1f\n",
				stheat[i].MQlossdy*Cff_kWh, stheat[i].MQstody*Cff_kWh)
		}
	}
}

func stheatmtprt(fo io.Writer, id int, stheat []*STHEAT, Mo, tt int) {
	switch id {
	case 0:
		if len(stheat) > 0 {
			fmt.Fprintf(fo, "%s %d\n", STHEAT_TYPE, len(stheat))
		}
		for i := range stheat {
			fmt.Fprintf(fo, " %s 1 1\n", stheat[i].Name)
		}
	case 1:
		for i := range stheat {
			fmt.Fprintf(fo, "%s_E E f \n", stheat[i].Name)
		}
	default:
		for i := range stheat {
			fmt.Fprintf(fo, " %.2f\n", stheat[i].MtEdy[Mo-1][tt-1].D*Cff_kWh)
		}
	}
}

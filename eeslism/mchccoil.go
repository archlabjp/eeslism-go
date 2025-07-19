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

/*  hccoil.c  */
package eeslism

import (
	"fmt"
	"io"
	"math"
	"strings"
)

/*
Hccdata (Heating/Cooling Coil Data Input)

この関数は、冷温水コイルの各種仕様（温度効率、エンタルピー効率、熱通過率と伝熱面積の積など）を読み込み、
対応するコイルの構造体に格納します。
これらのデータは、空調システムにおける熱交換器の性能評価、
熱負荷への対応、およびエネルギー消費量予測に不可欠です。

建築環境工学的な観点:
  - **熱交換器のモデル化**: 冷温水コイルは、空調システムにおいて空気と熱媒（冷水または温水）の間で熱を交換する主要な機器です。
    その性能を正確にモデル化することは、室内の温湿度環境を維持するために必要な熱量や、
    空調システムのエネルギー消費量を評価する上で非常に重要です。
  - **温度効率 (et)**:
    コイルを通過する空気の温度変化と、熱媒の温度変化の比率で定義される効率です。
    コイルの熱交換性能を示す指標であり、顕熱処理能力に影響します。
    `Hccca.et`が設定されている場合、コイルの温度効率が固定値として扱われることを示唆します。
  - **エンタルピー効率 (eh)**:
    コイルを通過する空気のエンタルピー変化と、熱媒のエンタルピー変化の比率で定義される効率です。
    コイルの全熱交換性能（顕熱と潜熱の両方）を示す指標であり、
    特に除湿運転時の潜熱処理能力に影響します。
    `Hccca.eh`が設定されている場合、コイルが湿りコイルとして機能し、
    除湿能力を持つことを示唆します。
  - **熱通過率と伝熱面積の積 (KA)**:
    コイルの熱交換能力を総合的に示すパラメータです。
    `KA`が大きいほど、コイルの熱交換能力が高いことを意味します。
    `Hccca.KA`が設定されている場合、コイルの温度効率が負荷に応じて変動するタイプとして扱われることを示唆します。
  - **乾きコイルと湿りコイル**: コイル表面が露点温度以下になると、空気中の水蒸気が凝縮し、
    潜熱交換（除湿）が行われます。
    `eh`が`0.0`より大きいか、`KA`が`0.0`より大きいかによって、
    コイルが乾きコイル（顕熱交換のみ）として機能するか、
    湿りコイル（顕熱と潜熱の両方）として機能するかが決定されます。

この関数は、空調システムにおける熱交換器の性能をモデル化し、
室内の温湿度環境の維持、熱負荷計算、およびエネルギー消費量予測を行うための重要なデータ入力機能を提供します。
*/
func Hccdata(s string, Hccca *HCCCA) int {
	var st string
	var dt float64
	id := 0

	if stIdx := strings.IndexRune(s, '='); stIdx == -1 {
		Hccca.name = s
		Hccca.eh = 0.0
		Hccca.et = FNAN
		Hccca.KA = FNAN
	} else {
		st = s[stIdx+1:]
		dt, _ = readFloat(st)

		if s == "et" {
			// コイル温度効率
			Hccca.et = dt
		} else if s == "eh" {
			// コイルエンタルピー効率
			Hccca.eh = dt
		} else if s == "KA" {
			// コイルの熱通過率と伝熱面積の積 [W/K]
			Hccca.KA = dt
		} else {
			id = 1
		}
	}

	return id
}

/*
Hccdwint (Heating/Cooling Coil Dry/Wet Initialization)

この関数は、冷温水コイルが乾きコイル（顕熱交換のみ）として機能するか、
湿りコイル（顕熱と潜熱の両方）として機能するかを初期設定します。
また、コイルの温度効率が固定タイプか変動タイプかを判定します。

建築環境工学的な観点:
  - **乾きコイルと湿りコイルの判定**: 空調システムにおいて、
    コイルが顕熱交換のみを行う「乾きコイル」として機能するか、
    顕熱と潜熱の両方を行う「湿りコイル」として機能するかは、
    コイル表面温度と空気の露点温度の関係によって決まります。
  - `hcc.Cat.eh > 1.0e-10`: エンタルピー効率がゼロより大きい場合、
    コイルが潜熱交換能力を持つと判断し、`hcc.Wet = 'w'`（湿りコイル）と設定します。
    これは、コイルが除湿運転を行う可能性があることを意味します。
  - それ以外の場合: `hcc.Wet = 'd'`（乾きコイル）と設定します。
    これは、コイルが顕熱交換のみを行うことを意味します。
    この判定は、コイルの熱負荷計算において、顕熱と潜熱の分離を適切に行うために不可欠です。
  - **温度効率タイプの判定**: コイルの温度効率の計算方法には、
    固定値を用いるタイプと、負荷に応じて変動するタイプがあります。
  - `hcc.Cat.et > 0.0`: 温度効率が固定値として与えられている場合、
    `hcc.Etype = 'e'`（定格温度効率タイプ）と設定します。
  - `hcc.Cat.KA > 0.0`: 熱通過率と伝熱面積の積（`KA`）が与えられている場合、
    コイルの温度効率が負荷に応じて変動するタイプとして、`hcc.Etype = 'k'`（変動タイプ）と設定します。
    この場合、`FNhccet`関数などを用いて、コイルの熱交換能力を計算します。
    この判定は、コイルの熱交換性能を正確にモデル化し、
    様々な運転条件下での熱負荷計算の精度を向上させるために重要です。

この関数は、冷温水コイルの熱交換特性を初期設定し、
空調システムの熱負荷計算やエネルギー消費量予測を行うための基礎的な役割を果たします。
*/
func Hccdwint(_hcc []*HCC) {
	for _, hcc := range _hcc {

		// 乾きコイルと湿りコイルの判定
		if hcc.Cat.eh > 1.0e-10 {
			hcc.Wet = 'w' // 湿りコイル
		} else {
			hcc.Wet = 'd' // 乾きコイル
		}

		// 温度効率固定タイプと変動タイプの判定
		if hcc.Cat.et > 0.0 {
			hcc.Etype = 'e' // 定格(温度効率固定タイプ)
		} else if hcc.Cat.KA > 0.0 {
			hcc.Etype = 'k' // 変動タイプ
		} else {
			fmt.Printf("Hcc %s  Undefined Character et or KA\n", hcc.Name)
			hcc.Etype = '-'
		}

		// 入口水温、入口空気絶対湿度を初期化
		//Hcc.Twin = 5.0
		//Hcc.xain = FNXtr(25.0, 50.0)
	}
}

/* ------------------------------------------ */
/*
Hcccfv (Heating/Cooling Coil Characteristic Function Value Calculation)

この関数は、冷温水コイルの運転特性を評価し、
空気側および水側の流量、熱容量流量、そして温度効率やエンタルピー効率を計算します。
これは、コイルが空気と熱媒の間で熱と湿気をどのように交換するかをモデル化するために不可欠です。

建築環境工学的な観点:
- **空気側・水側の熱容量流量 (hcc.cGa, hcc.cGw)**:
  熱交換器の性能は、空気側と水側の熱容量流量（質量流量と比熱の積）に大きく依存します。
  `hcc.cGa`は空気側の熱容量流量、`hcc.cGw`は水側の熱容量流量を表します。
  これらの値は、コイルの熱交換能力や、空気と水の温度変化を決定する上で重要です。
- **温度効率 (hcc.et)**:
  コイルの顕熱交換能力を示す指標です。
  `hcc.Etype`が`'e'`（定格温度効率タイプ）の場合は、設定された固定値を使用し、
  `'k'`（変動タイプ）の場合は、`FNhccet`関数を用いて熱容量流量と`KA`値から計算します。
  これにより、コイルの顕熱処理能力を正確にモデル化できます。
- **エンタルピー効率 (hcc.eh)**:
  コイルの全熱交換能力（顕熱と潜熱の両方）を示す指標です。
  特に、湿りコイルの場合、除湿能力を評価する上で重要です。
- **乾きコイルと湿りコイルの処理 (wcoil)**:
  `wcoil`関数は、コイルが乾きコイルか湿りコイルかに応じて、
  空気と水の熱交換を計算します。
  湿りコイルの場合、空気中の水蒸気の凝縮による潜熱交換が考慮され、
  コイルの出口空気の絶対湿度や、凝縮水の発生量が計算されます。
- **出口温度・湿度の係数設定**: 計算されたコイルの特性に基づいて、
  出口空気温度（`eo_ta`）、出口空気絶対湿度（`eo_xa`）、
  出口水温度（`eo_tw`）に関する係数（`Coeffo`, `Co`, `Coeffin`）を設定します。
  これらの係数は、空調システム全体の熱収支方程式に組み込まれ、
  室内の温湿度環境を予測するために用いられます。

この関数は、冷温水コイルの熱交換特性を詳細にモデル化し、
空調システムの熱負荷計算、室内温湿度環境の予測、
およびエネルギー消費量予測を行うための重要な役割を果たします。
*/
func Hcccfv(_hcc []*HCC) {
	for _, hcc := range _hcc {
		hcc.Ga = 0.0
		hcc.Gw = 0.0
		hcc.et = 0.0
		hcc.eh = 0.0

		// 経路が停止していなければ
		if hcc.Cmp.Control == OFF_SW {
			continue
		}

		// 機器出力は3つ
		if len(hcc.Cmp.Elouts) != 3 || len(hcc.Cmp.Elins) != 0 {
			panic("HCCの機器出力数は3、機器入力は0です。")
		}

		eo_ta := hcc.Cmp.Elouts[0] // 排気温度
		eo_xa := hcc.Cmp.Elouts[1] // 排気湿度
		eo_tw := hcc.Cmp.Elouts[2] // 排水温度

		var AirSW, WaterSW ControlSWType

		// 排気量・排気熱量
		hcc.Ga = eo_ta.G                        // 排気量
		hcc.cGa = Spcheat(eo_ta.Fluid) * hcc.Ga // 排気熱量
		if hcc.Ga > 0.0 {
			AirSW = ON_SW
		} else {
			AirSW = OFF_SW
		}

		// 排水量・排水熱量
		hcc.Gw = eo_tw.G                        // 排水量
		hcc.cGw = Spcheat(eo_tw.Fluid) * hcc.Gw // 排水熱量
		if hcc.Gw > 0.0 {
			WaterSW = ON_SW
		} else {
			WaterSW = OFF_SW
		}

		// 温度効率
		if hcc.Etype == 'e' {
			// 定格温度効率
			hcc.et = hcc.Cat.et
		} else if hcc.Etype == 'k' {
			// 温度効率を計算
			hcc.et = FNhccet(hcc.cGa, hcc.cGw, hcc.Cat.KA)
		} else {
			panic(hcc.Etype)
		}

		// エンタルピ効率 [-]
		hcc.eh = hcc.Cat.eh

		// 冷温水コイルの処理熱量
		hcc.Et, hcc.Ex, hcc.Ew = wcoil(AirSW, WaterSW, hcc.Wet, hcc.Ga*hcc.et, hcc.Ga*hcc.eh, hcc.Xain, hcc.Twin)

		// 排気温度に関する係数の設定
		eo_ta.Coeffo = hcc.cGa
		eo_ta.Co = -(hcc.Et.C)
		eo_ta.Coeffin[0] = hcc.Et.T - hcc.cGa
		eo_ta.Coeffin[1] = hcc.Et.X
		eo_ta.Coeffin[2] = -(hcc.Et.W)

		// 排気湿度に関する係数の設定
		eo_xa.Coeffo = hcc.Ga
		eo_xa.Co = -(hcc.Ex.C)
		eo_xa.Coeffin[0] = hcc.Ex.T
		eo_xa.Coeffin[1] = hcc.Ex.X - hcc.Ga
		eo_xa.Coeffin[2] = -(hcc.Ex.W)

		// 排水温度に関する係数の設定
		eo_tw.Coeffo = hcc.cGw
		eo_tw.Co = hcc.Ew.C
		eo_tw.Coeffin[0] = -(hcc.Ew.T)
		eo_tw.Coeffin[1] = -(hcc.Ew.X)
		eo_tw.Coeffin[2] = hcc.Ew.W - hcc.cGw
	}
}

/*
Hccdwreset (Heating/Cooling Coil Dry/Wet Reset)

この関数は、冷温水コイルの運転状態（乾きコイルか湿りコイルか）を再判定し、
必要に応じてコイルの特性係数を再計算します。
これは、コイル表面での結露の有無が、コイルの熱交換性能に大きく影響するためです。

建築環境工学的な観点:
  - **乾きコイルから湿りコイルへの遷移**: コイルを通過する空気の露点温度が、
    コイル表面温度（または熱媒温度）を下回ると、コイル表面で結露が発生し、
    コイルは乾きコイルから湿りコイルへと遷移します。
    この遷移は、コイルの潜熱交換能力が発揮されることを意味し、
    空調システムの除湿運転において非常に重要です。
  - **露点温度の計算**: `Tdp := FNDp(FNPwx(xain))` は、
    コイル入口空気の絶対湿度（`xain`）から飽和水蒸気圧を計算し、
    それに対応する露点温度（`Tdp`）を算出しています。
    この露点温度と熱媒温度（`Twin`）を比較することで、
    コイル表面での結露の有無を判定します。
  - **特性係数の再計算**: コイルが乾きコイルから湿りコイルへ、
    あるいは湿りコイルから乾きコイルへと状態が変化した場合、
    その熱交換特性も変化します。
    そのため、`reset = true` となった場合に `Hcccfv(Hcc[i : i+1])` を呼び出し、
    コイルの特性係数を再計算することで、
    コイルの熱交換性能を正確にモデル化し続けることができます。
  - **空調システムの制御**: この乾き/湿り状態の判定と特性係数の再計算は、
    空調システムの制御戦略にも影響を与えます。
    例えば、除湿運転が必要な場合にコイルを湿り状態に保つための制御や、
    過度な除湿を防ぐための制御などが考えられます。

この関数は、冷温水コイルの動的な熱交換特性をモデル化し、
空調システムの熱負荷計算、特に潜熱負荷の処理、
および室内温湿度環境の予測精度を向上させるために不可欠な役割を果たします。
*/
func Hccdwreset(Hcc []*HCC, DWreset *int) {
	for i, hcc := range Hcc {
		xain := hcc.Cmp.Elins[1].Sysvin // <給気>絶対湿度 [kg/kg]
		Twin := hcc.Cmp.Elins[2].Sysvin // <給水>温水の温度 [C]

		reset := false
		if hcc.Cat.eh > 1.0e-10 {
			Tdp := FNDp(FNPwx(xain)) // 露点温度
			if hcc.Wet == 'w' && Twin > Tdp {
				// 露点温度を上回った => 結露なし (乾きコイル)
				hcc.Wet = 'd'
				reset = true
			} else if hcc.Wet == 'd' && Twin < Tdp {
				// 露点温度を上回った => 結露あり (湿りコイル)
				hcc.Wet = 'w'
				reset = true
			}

			if reset {
				(*DWreset)++
				Hcccfv(Hcc[i : i+1])
			}
		}
	}
}

/*
Hccene (Heating/Cooling Coil Energy Calculation)

この関数は、冷温水コイルが空気側および水側に供給または除去する熱量（顕熱、潜熱、全熱）を計算します。
これは、空調システムにおける熱負荷の処理能力や、エネルギー消費量を評価する上で不可欠です。

建築環境工学的な観点:
  - **顕熱量 (hcc.Qs)**:
    空気の温度変化に伴う熱量であり、室内の温度制御に直接関係します。
    `hcc.cGa * (hcc.Taout - hcc.Tain)` のように、空気側の熱容量流量と入口・出口空気温度差から計算されます。
    暖房時には正の値、冷房時には負の値となります。
  - **潜熱量 (hcc.Ql)**:
    空気中の水蒸気量の変化に伴う熱量であり、室内の湿度制御に直接関係します。
    `Ro * hcc.Ga * (Xaout - hcc.Xain)` のように、空気の密度、空気流量、入口・出口空気絶対湿度差から計算されます。
    除湿時には負の値、加湿時には正の値となります。
  - **水側熱量 (hcc.Qt)**:
    熱媒（冷水または温水）の温度変化に伴う熱量であり、
    熱源設備（チラー、ボイラーなど）との熱交換量を表します。
    `hcc.cGw * (hcc.Twout - hcc.Twin)` のように、水側の熱容量流量と入口・出口水温度差から計算されます。
  - **熱負荷のバランス**: これらの熱量は、室内の熱負荷（顕熱負荷、潜熱負荷）とバランスが取れるように計算されます。
    空調システムは、これらの熱量を適切に供給または除去することで、
    室内の温湿度環境を目標値に維持します。
  - **エネルギー消費量への影響**: コイルが処理する熱量は、
    熱源設備や搬送動力（ファン、ポンプ）のエネルギー消費量に直接影響します。
    これらの熱量を正確に把握することで、
    空調システム全体のエネルギー効率を評価し、省エネルギー対策の効果を検証できます。

この関数は、冷温水コイルの熱交換性能を定量的に評価し、
空調システムの熱負荷処理能力、室内温湿度環境の予測、
およびエネルギー消費量予測を行うための重要な役割を果たします。
*/
func Hccene(Hcc []*HCC) {
	for _, hcc := range Hcc {
		hcc.Tain = hcc.Cmp.Elins[0].Sysvin // <給気>空気温度 [C]
		hcc.Xain = hcc.Cmp.Elins[1].Sysvin // <給気>絶対湿度 [kg/kg]
		hcc.Twin = hcc.Cmp.Elins[2].Sysvin // <給水>温水の温度 [C]

		if hcc.Cmp.Control != OFF_SW {
			// <排気>空気温度 [C]
			hcc.Taout = hcc.Cmp.Elouts[0].Sysv
			hcc.Qs = hcc.cGa * (hcc.Taout - hcc.Tain)

			// <排気>空気絶対湿度 [kg/kg]
			Xaout := hcc.Cmp.Elouts[1].Sysv
			hcc.Ql = Ro * hcc.Ga * (Xaout - hcc.Xain)

			// <排水>温水の温度 [C]
			hcc.Twout = hcc.Cmp.Elouts[2].Sysv
			hcc.Qt = hcc.cGw * (hcc.Twout - hcc.Twin)
		} else {
			// 経路が停止している場合は熱供給しない
			hcc.Qs = 0.0
			hcc.Ql = 0.0
			hcc.Qt = 0.0
		}
	}
}

/* ------------------------------------------ */

// 冷温水コイルHccの状態をfoに出力する。
func hccprint(fo io.Writer, id int, Hcc []*HCC) {
	switch id {
	case 0:
		if len(Hcc) > 0 {
			fmt.Fprintf(fo, "%s %d\n", HCCOIL_TYPE, len(Hcc))
		}
		for _, hcc := range Hcc {
			fmt.Fprintf(fo, " %s 1 16\n", hcc.Name)
		}
	case 1:
		for _, hcc := range Hcc {
			fmt.Fprintf(fo, "%s_ca c c %s_Ga m f %s_Ti t f %s_To t f %s_Qs q f\n", hcc.Name, hcc.Name, hcc.Name, hcc.Name, hcc.Name)
			fmt.Fprintf(fo, "%s_cx c c %s_xi x f %s_xo x f %s_Ql q f\n", hcc.Name, hcc.Name, hcc.Name, hcc.Name)
			fmt.Fprintf(fo, "%s_cw c c %s_Gw m f %s_Twi t f %s_Two t f %s_Qt q f\n", hcc.Name, hcc.Name, hcc.Name, hcc.Name, hcc.Name)
			fmt.Fprintf(fo, "%s_et m f %s_eh m f\n\n", hcc.Name, hcc.Name)
		}
	default:
		for _, hcc := range Hcc {
			// 給排気温度に関する事項
			eo_ta := hcc.Cmp.Elouts[0]
			fmt.Fprintf(fo, "%c %6.4g %4.1f %4.1f %2.0f ", eo_ta.Control, hcc.Ga, hcc.Tain, eo_ta.Sysv, hcc.Qs)

			// 給排気湿度に関する事項
			eo_xa := hcc.Cmp.Elouts[1]
			fmt.Fprintf(fo, "%c %5.3f %5.3f %2.0f ", eo_xa.Control, hcc.Xain, eo_xa.Sysv, hcc.Ql)

			// 給排水温度に関する事項
			eo_tw := hcc.Cmp.Elouts[2]
			fmt.Fprintf(fo, "%c %6.4g %4.1f %4.1f %2.0f ", eo_tw.Control, hcc.Gw, hcc.Twin, eo_tw.Sysv, hcc.Qt)

			// 温度効率、エンタルピー
			fmt.Fprintf(fo, "%6.4g %6.4g\n", hcc.et, hcc.eh)
		}
	}
}

/* ------------------------------ */

/*
hccdyint (Heating/Cooling Coil Daily Integration Initialization)

この関数は、冷温水コイルの日積算値（日ごとの入口空気温度、絶対湿度、水温度、
顕熱量、潜熱量、水側熱量など）をリセットします。
これは、日単位でのコイルの運転状況や熱交換量を集計し、
空調システムの性能評価やエネルギー消費量の分析に用いられます。

建築環境工学的な観点:
  - **日単位の性能評価**: 冷温水コイルの運転状況は、日中の熱負荷変動に応じて大きく変化します。
    日積算値を集計することで、日ごとの顕熱・潜熱負荷の割合、
    コイルの稼働時間、部分負荷運転の割合などを把握できます。
    これにより、特定の日の空調負荷特性を分析したり、
    コイルの運転効率を日単位で評価したりすることが可能になります。
  - **運用改善の指標**: 日積算データは、空調システムの運用改善のための重要な指標となります。
    例えば、外気温度や湿度などの気象条件とコイルの熱交換量の関係を分析したり、
    設定温度や換気量などの運用条件がコイルの性能に与える影響を評価したりすることで、
    より効率的な運転方法を見つけることができます。
  - **データ集計の準備**: この関数は、新しい日の集計を開始する前に、
    前日のデータをクリアする役割を担います。
    `svdyint`や`qdyint`といった関数は、
    それぞれ温度、湿度、熱量などの日積算値をリセットするためのものです。

この関数は、冷温水コイルの運転状況と熱交換量を日単位で詳細に分析し、
空調システムの運用改善や省エネルギー対策の効果評価を行うための基礎的な役割を果たします。
*/
func hccdyint(Hcc []*HCC) {
	for _, hcc := range Hcc {
		svdyint(&hcc.Taidy)
		svdyint(&hcc.xaidy)
		svdyint(&hcc.Twidy)
		qdyint(&hcc.Qdys)
		qdyint(&hcc.Qdyl)
		qdyint(&hcc.Qdyt)
	}
}

/*
hccmonint (Heating/Cooling Coil Monthly Integration Initialization)

この関数は、冷温水コイルの月積算値（月ごとの入口空気温度、絶対湿度、水温度、
顕熱量、潜熱量、水側熱量など）をリセットします。
これは、月単位でのコイルの運転状況や熱交換量を集計し、
空調システムの性能評価やエネルギー消費量の分析に用いられます。

建築環境工学的な観点:
  - **月単位の性能評価**: 冷温水コイルの運転状況は、月単位で変動します。
    月積算値を集計することで、月ごとの顕熱・潜熱負荷の割合、
    コイルの稼働時間、部分負荷運転の割合などを把握できます。
    これにより、特定の月の空調負荷特性を分析したり、
    コイルの運転効率を月単位で評価したりすることが可能になります。
  - **運用改善の指標**: 月積算データは、空調システムの運用改善のための重要な指標となります。
    例えば、季節ごとの空調負荷の傾向を把握したり、
    月ごとの気象条件とコイルの熱交換量の関係を分析したりすることで、
    より効率的な運転方法を見つけることができます。
  - **データ集計の準備**: この関数は、新しい月の集計を開始する前に、
    前月のデータをクリアする役割を担います。
    `svdyint`や`qdyint`といった関数は、
    それぞれ温度、湿度、熱量などの月積算値をリセットするためのものです。

この関数は、冷温水コイルの運転状況と熱交換量を月単位で詳細に分析し、
空調システムの運用改善や省エネルギー対策の効果評価を行うための基礎的な役割を果たします。
*/
func hccmonint(Hcc []*HCC) {
	for _, hcc := range Hcc {
		svdyint(&hcc.mTaidy)
		svdyint(&hcc.mxaidy)
		svdyint(&hcc.mTwidy)
		qdyint(&hcc.mQdys)
		qdyint(&hcc.mQdyl)
		qdyint(&hcc.mQdyt)
	}
}

/*
hccday (Heating/Cooling Coil Daily and Monthly Data Aggregation)

この関数は、冷温水コイルの運転データ（入口空気温度、絶対湿度、水温度、
顕熱量、潜熱量、水側熱量など）を、日単位および月単位で集計します。
これにより、コイルの性能評価やエネルギー消費量の分析が可能になります。

建築環境工学的な観点:
  - **日次集計 (svdaysum, qdaysum)**:
    日次集計は、コイルの運転状況を日単位で詳細に把握するために重要です。
    例えば、特定の日の顕熱・潜熱負荷の変動に対するコイルの応答、
    あるいは日中のピーク負荷時の熱交換量などを分析できます。
    これにより、日ごとの運用改善点を見つけ出すことが可能になります。
  - **月次集計 (svmonsum, qmonsum)**:
    月次集計は、季節ごとの熱負荷変動や、
    コイルの年間を通じた熱交換量の傾向を把握するために重要です。
    これにより、年間を通じた省エネルギー対策の効果を評価したり、
    熱交換量の予測精度を向上させたりすることが可能になります。
  - **データ分析の基礎**: この関数で集計されるデータは、
    冷温水コイルの性能評価、熱交換量のベンチマーキング、
    省エネルギー対策の効果検証、そして運用改善のための意思決定の基礎となります。

この関数は、冷温水コイルの運転状況と熱交換量を多角的に分析し、
空調システムの運用改善や省エネルギー対策の効果評価を行うための重要なデータ集計機能を提供します。
*/
func hccday(Mon, Day, ttmm int, Hcc []*HCC, Nday, SimDayend int) {
	for _, hcc := range Hcc {
		// 日集計
		svdaysum(int64(ttmm), hcc.Cmp.Control, hcc.Tain, &hcc.Taidy)
		svdaysum(int64(ttmm), hcc.Cmp.Control, hcc.Xain, &hcc.xaidy)
		svdaysum(int64(ttmm), hcc.Cmp.Control, hcc.Twin, &hcc.Twidy)
		qdaysum(int64(ttmm), hcc.Cmp.Control, hcc.Qs, &hcc.Qdys)
		qdaysum(int64(ttmm), hcc.Cmp.Control, hcc.Ql, &hcc.Qdyl)
		qdaysum(int64(ttmm), hcc.Cmp.Control, hcc.Qt, &hcc.Qdyt)

		// 月集計
		svmonsum(Mon, Day, ttmm, hcc.Cmp.Control, hcc.Tain, &hcc.mTaidy, Nday, SimDayend)
		svmonsum(Mon, Day, ttmm, hcc.Cmp.Control, hcc.Xain, &hcc.mxaidy, Nday, SimDayend)
		svmonsum(Mon, Day, ttmm, hcc.Cmp.Control, hcc.Twin, &hcc.mTwidy, Nday, SimDayend)
		qmonsum(Mon, Day, ttmm, hcc.Cmp.Control, hcc.Qs, &hcc.mQdys, Nday, SimDayend)
		qmonsum(Mon, Day, ttmm, hcc.Cmp.Control, hcc.Ql, &hcc.mQdyl, Nday, SimDayend)
		qmonsum(Mon, Day, ttmm, hcc.Cmp.Control, hcc.Qt, &hcc.mQdyt, Nday, SimDayend)
	}
}

func hccdyprt(fo io.Writer, id int, Hcc []*HCC) {
	switch id {
	case 0:
		if len(Hcc) > 0 {
			fmt.Fprintf(fo, "%s %d\n", HCCOIL_TYPE, len(Hcc))
		}
		for _, hcc := range Hcc {
			fmt.Fprintf(fo, "%s 4 42 14 14 14\n", hcc.Name)
		}
	case 1:
		for _, hcc := range Hcc {
			fmt.Fprintf(fo, "%s_Ht H d %s_T T f ", hcc.Name, hcc.Name)
			fmt.Fprintf(fo, "%s_ttn h d %s_Tn t f %s_ttm h d %s_Tm t f\n",
				hcc.Name, hcc.Name, hcc.Name, hcc.Name)
			fmt.Fprintf(fo, "%s_Hsh H d %s_Qsh Q f %s_Hsc H d %s_Qsc Q f\n",
				hcc.Name, hcc.Name, hcc.Name, hcc.Name)
			fmt.Fprintf(fo, "%s_tsh h d %s_qsh q f %s_tsc h d %s_qsc q f\n",
				hcc.Name, hcc.Name, hcc.Name, hcc.Name)

			fmt.Fprintf(fo, "%s_Hx H d %s_x X f ", hcc.Name, hcc.Name)
			fmt.Fprintf(fo, "%s_txn h d %s_xn x f %s_txm h d %s_xm c f\n",
				hcc.Name, hcc.Name, hcc.Name, hcc.Name)
			fmt.Fprintf(fo, "%s_Hlh H d %s_Qlh Q f %s_Hlc H d %s_Qlc Q f\n",
				hcc.Name, hcc.Name, hcc.Name, hcc.Name)
			fmt.Fprintf(fo, "%s_tlh h d %s_qlh q f %s_tlc h d %s_qlc q f\n",
				hcc.Name, hcc.Name, hcc.Name, hcc.Name)

			fmt.Fprintf(fo, "%s_Htw H d %s_Tw T f ", hcc.Name, hcc.Name)
			fmt.Fprintf(fo, "%s_ttwn h d %s_Twn t f %s_ttwm h d %s_Twm t f\n",
				hcc.Name, hcc.Name, hcc.Name, hcc.Name)
			fmt.Fprintf(fo, "%s_Hth H d %s_Qth Q f %s_Htc H d %s_Qtc Q f\n",
				hcc.Name, hcc.Name, hcc.Name, hcc.Name)
			fmt.Fprintf(fo, "%s_tth h d %s_qth q f %s_ttc h d %s_qtc q f\n\n",
				hcc.Name, hcc.Name, hcc.Name, hcc.Name)
		}
	default:
		for _, hcc := range Hcc {
			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ",
				hcc.Taidy.Hrs, hcc.Taidy.M,
				hcc.Taidy.Mntime, hcc.Taidy.Mn,
				hcc.Taidy.Mxtime, hcc.Taidy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f ", hcc.Qdys.Hhr, hcc.Qdys.H)
			fmt.Fprintf(fo, "%1d %3.1f ", hcc.Qdys.Chr, hcc.Qdys.C)
			fmt.Fprintf(fo, "%1d %2.0f ", hcc.Qdys.Hmxtime, hcc.Qdys.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f ", hcc.Qdys.Cmxtime, hcc.Qdys.Cmx)

			fmt.Fprintf(fo, "%1d %5.4f %1d %5.4f %1d %5.4f ",
				hcc.xaidy.Hrs, hcc.xaidy.M,
				hcc.xaidy.Mntime, hcc.xaidy.Mn,
				hcc.xaidy.Mxtime, hcc.xaidy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f ", hcc.Qdyl.Hhr, hcc.Qdyl.H)
			fmt.Fprintf(fo, "%1d %3.1f ", hcc.Qdyl.Chr, hcc.Qdyl.C)
			fmt.Fprintf(fo, "%1d %2.0f ", hcc.Qdyl.Hmxtime, hcc.Qdyl.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f ", hcc.Qdyl.Cmxtime, hcc.Qdyl.Cmx)

			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ",
				hcc.Twidy.Hrs, hcc.Twidy.M,
				hcc.Twidy.Mntime, hcc.Twidy.Mn,
				hcc.Twidy.Mxtime, hcc.Twidy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f ", hcc.Qdyt.Hhr, hcc.Qdyt.H)
			fmt.Fprintf(fo, "%1d %3.1f ", hcc.Qdyt.Chr, hcc.Qdyt.C)
			fmt.Fprintf(fo, "%1d %2.0f ", hcc.Qdyt.Hmxtime, hcc.Qdyt.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f\n", hcc.Qdyt.Cmxtime, hcc.Qdyt.Cmx)
		}
	}
}

func hccmonprt(fo io.Writer, id int, Hcc []*HCC) {
	switch id {
	case 0:
		if len(Hcc) > 0 {
			fmt.Fprintf(fo, "%s %d\n", HCCOIL_TYPE, len(Hcc))
		}
		for _, hcc := range Hcc {
			fmt.Fprintf(fo, "%s 4 42 14 14 14\n", hcc.Name)
		}
	case 1:
		for _, hcc := range Hcc {
			fmt.Fprintf(fo, "%s_Ht H d %s_T T f ", hcc.Name, hcc.Name)
			fmt.Fprintf(fo, "%s_ttn h d %s_Tn t f %s_ttm h d %s_Tm t f\n",
				hcc.Name, hcc.Name, hcc.Name, hcc.Name)
			fmt.Fprintf(fo, "%s_Hsh H d %s_Qsh Q f %s_Hsc H d %s_Qsc Q f\n",
				hcc.Name, hcc.Name, hcc.Name, hcc.Name)
			fmt.Fprintf(fo, "%s_tsh h d %s_qsh q f %s_tsc h d %s_qsc q f\n",
				hcc.Name, hcc.Name, hcc.Name, hcc.Name)

			fmt.Fprintf(fo, "%s_Hx H d %s_x X f ", hcc.Name, hcc.Name)
			fmt.Fprintf(fo, "%s_txn h d %s_xn x f %s_txm h d %s_xm c f\n",
				hcc.Name, hcc.Name, hcc.Name, hcc.Name)
			fmt.Fprintf(fo, "%s_Hlh H d %s_Qlh Q f %s_Hlc H d %s_Qlc Q f\n",
				hcc.Name, hcc.Name, hcc.Name, hcc.Name)
			fmt.Fprintf(fo, "%s_tlh h d %s_qlh q f %s_tlc h d %s_qlc q f\n",
				hcc.Name, hcc.Name, hcc.Name, hcc.Name)

			fmt.Fprintf(fo, "%s_Htw H d %s_Tw T f ", hcc.Name, hcc.Name)
			fmt.Fprintf(fo, "%s_ttwn h d %s_Twn t f %s_ttwm h d %s_Twm t f\n",
				hcc.Name, hcc.Name, hcc.Name, hcc.Name)
			fmt.Fprintf(fo, "%s_Hth H d %s_Qth Q f %s_Htc H d %s_Qtc Q f\n",
				hcc.Name, hcc.Name, hcc.Name, hcc.Name)
			fmt.Fprintf(fo, "%s_tth h d %s_qth q f %s_ttc h d %s_qtc q f\n\n",
				hcc.Name, hcc.Name, hcc.Name, hcc.Name)
		}
	default:
		for _, hcc := range Hcc {
			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ",
				hcc.mTaidy.Hrs, hcc.mTaidy.M,
				hcc.mTaidy.Mntime, hcc.mTaidy.Mn,
				hcc.mTaidy.Mxtime, hcc.mTaidy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f ", hcc.mQdys.Hhr, hcc.mQdys.H)
			fmt.Fprintf(fo, "%1d %3.1f ", hcc.mQdys.Chr, hcc.mQdys.C)
			fmt.Fprintf(fo, "%1d %2.0f ", hcc.mQdys.Hmxtime, hcc.mQdys.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f ", hcc.mQdys.Cmxtime, hcc.mQdys.Cmx)

			fmt.Fprintf(fo, "%1d %5.4f %1d %5.4f %1d %5.4f ",
				hcc.mxaidy.Hrs, hcc.mxaidy.M,
				hcc.mxaidy.Mntime, hcc.mxaidy.Mn,
				hcc.mxaidy.Mxtime, hcc.mxaidy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f ", hcc.mQdyl.Hhr, hcc.mQdyl.H)
			fmt.Fprintf(fo, "%1d %3.1f ", hcc.mQdyl.Chr, hcc.mQdyl.C)
			fmt.Fprintf(fo, "%1d %2.0f ", hcc.mQdyl.Hmxtime, hcc.mQdyl.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f ", hcc.mQdyl.Cmxtime, hcc.mQdyl.Cmx)

			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ",
				hcc.mTwidy.Hrs, hcc.mTwidy.M,
				hcc.mTwidy.Mntime, hcc.mTwidy.Mn,
				hcc.mTwidy.Mxtime, hcc.mTwidy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f ", hcc.mQdyt.Hhr, hcc.mQdyt.H)
			fmt.Fprintf(fo, "%1d %3.1f ", hcc.mQdyt.Chr, hcc.mQdyt.C)
			fmt.Fprintf(fo, "%1d %2.0f ", hcc.mQdyt.Hmxtime, hcc.mQdyt.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f\n", hcc.mQdyt.Cmxtime, hcc.mQdyt.Cmx)
		}
	}
}

/*
FNhccet (Function for Heating/Cooling Coil Temperature Effectiveness)

この関数は、冷温水コイルの温度効率（熱交換効率）を計算します。
特に、熱通過率と伝熱面積の積（KA値）が与えられている場合に、
熱容量流量比を考慮したコイルの熱交換性能を評価します。

建築環境工学的な観点:
  - **熱交換器の効率**: 熱交換器の効率は、
    熱交換器がどれだけ効率的に熱を伝達できるかを示す重要な指標です。
    温度効率は、顕熱交換の効率を表し、
    空調システムにおける顕熱負荷の処理能力に直接関係します。
  - **NTU法**: この関数は、熱交換器の設計や性能評価に用いられるNTU（Number of Transfer Units）法に基づいています。
  - `NTU = KA / Ws`: NTUは、熱交換器の熱交換能力（`KA`）と、
    最小熱容量流量（`Ws`、ここでは空気側熱容量流量`Wa`と水側熱容量流量`Ww`の小さい方）の比で定義されます。
    NTUが大きいほど、熱交換器の熱交換能力が高いことを意味します。
  - `C = Ws / Wl`: 熱容量流量比。最大熱容量流量（`Wl`）に対する最小熱容量流量の比です。
    この比は、熱交換器の効率に影響を与えます。
  - **向流コイルのモデル**: この関数は、向流（Counterflow）コイルの熱交換モデルを適用しています。
    向流は、熱媒と空気の流れが逆方向であるため、
    並流（Parallel flow）に比べて熱交換効率が高くなるという特徴があります。
    ` (1.0 - exB) / (1.0 - C*exB)` の式は、向流熱交換器の効率を計算する一般的な式です。
  - **省エネルギー設計**: 熱交換器の効率を正確に評価することは、
    空調システムの省エネルギー設計において非常に重要です。
    効率の高いコイルを選定したり、
    コイルの運転条件を最適化したりすることで、
    エネルギー消費量を削減できます。

この関数は、冷温水コイルの熱交換性能を定量的に評価し、
空調システムの設計最適化や、エネルギー消費量予測を行うための重要な役割を果たします。
*/
func FNhccet(Wa, Ww, KA float64) float64 {
	Ws := Wa
	Wl := Ww

	NTU := KA / Ws
	C := Ws / Wl
	B := (1.0 - C) * NTU

	if math.Abs(Ws-Wl) < 1.0e-5 {
		return NTU / (1.0 + NTU)
	} else {
		if exB := math.Exp(-B); math.IsInf(exB, 0) {
			return 1.0 / C
		} else {
			return (1.0 - exB) / (1.0 - C*exB)
		}
	}
}

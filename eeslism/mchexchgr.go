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

/* hexchgr.c */

package eeslism

import (
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
)

// -------------------------------------------------------------
// 熱交換器
//
// 冷風入力 [IN  1] ---> +-----+ <--- [IN  2] 温風入力
//                       | HEX |
// 冷風出力 [OUT 1] <--- +-----+ ---> [OUT 2] 温風出力
//
// -------------------------------------------------------------

/*
Hexdata (Heat Exchanger Data Input)

この関数は、熱交換器の各種仕様（効率、熱通過率と伝熱面積の積など）を読み込み、
対応する熱交換器の構造体に格納します。
これらのデータは、熱回収システムや熱源システムにおける熱交換器の性能評価、
熱負荷への対応、およびエネルギー消費量予測に不可欠です。

建築環境工学的な観点:
- **熱交換器のモデル化**: 熱交換器は、異なる温度の流体間で熱を交換する機器であり、
  空調システム、給湯システム、熱回収システムなど、様々な用途で用いられます。
  その性能を正確にモデル化することは、システム全体のエネルギー効率を評価する上で非常に重要です。
- **効率 (eff)**:
  熱交換器の熱交換効率を示す指標です。
  投入された熱量に対して、どれだけの熱量を回収できたかを表します。
  効率が高いほど、エネルギーの無駄が少なく、省エネルギーに貢献します。
  `Hexca.eff`が設定されている場合、熱交換器の効率が固定値として扱われることを示唆します。
- **熱通過率と伝熱面積の積 (KA)**:
  熱交換器の熱交換能力を総合的に示すパラメータです。
  `KA`が大きいほど、熱交換器の熱交換能力が高いことを意味します。
  `Hexca.KA`が設定されている場合、熱交換器の効率が負荷に応じて変動するタイプとして扱われることを示唆します。
- **熱回収の重要性**: 熱交換器は、排気や排水などから熱を回収し、
  新鮮な空気や水に熱を供給することで、エネルギーの有効利用を促進します。
  これにより、熱源設備の負荷を軽減し、システム全体のエネルギー消費量を削減できます。

この関数は、熱交換器の性能をモデル化し、
熱回収システムや熱源システムの設計、熱負荷計算、
およびエネルギー消費量予測を行うための重要なデータ入力機能を提供します。
*/
func Hexdata(s string, Hexca *HEXCA) int {
	st := strings.IndexByte(s, '=')
	if st == -1 {
		Hexca.Name = s
		Hexca.eff = -999.0
		Hexca.KA = -999.0
	} else {
		s1 := s[:st]
		s2 := s[st+1:]
		if s1 == "eff" {
			e, err := strconv.ParseFloat(s2, 64)
			if err != nil {
				return 1
			}
			Hexca.eff = e
		} else if s1 == "KA" {
			ka, err := strconv.ParseFloat(s2, 64)
			if err != nil {
				return 1
			}
			Hexca.KA = ka
		} else {
			return 1
		}
	}
	return 0
}

/* --------------------------- */

/*
Hexcfv (Heat Exchanger Characteristic Function Value Calculation)

この関数は、熱交換器の運転特性を評価し、
熱媒の流量、熱容量流量、そして熱交換効率を計算します。
これは、熱交換器が異なる温度の流体間で熱をどのように交換するかをモデル化するために不可欠です。

建築環境工学的な観点:
- **熱媒の熱容量流量 (hex.CGc, hex.CGh)**:
  熱交換器の性能は、冷媒側（`hex.CGc`）と温媒側（`hex.CGh`）の熱容量流量に大きく依存します。
  これらの値は、熱交換器の熱交換能力や、各流体の温度変化を決定する上で重要です。
- **熱交換効率 (hex.Eff)**:
  熱交換器の熱交換効率を示す指標です。
  `hex.Etype`が`'e'`（効率固定タイプ）の場合は、設定された固定値を使用し、
  `'k'`（KA値変動タイプ）の場合は、`FNhccet`関数を用いて熱容量流量と`KA`値から計算します。
  これにより、熱交換器の熱交換性能を正確にモデル化できます。
- **熱交換能力 (eCGmin)**:
  `eCGmin = hex.Eff * math.Min(hex.CGc, hex.CGh)` は、
  熱交換器が実際に交換できる熱量の上限を示します。
  これは、熱交換器の効率と、熱容量流量の小さい方の値によって決定されます。
- **出口温度の係数設定**: 計算された熱交換器の特性に基づいて、
  冷媒側出口温度（`eoc`）と温媒側出口温度（`eoh`）に関する係数（`Coeffin`, `Coeffo`, `Co`）を設定します。
  これらの係数は、システム全体の熱収支方程式に組み込まれ、
  各流体の出口温度を予測するために用いられます。

この関数は、熱交換器の熱交換特性を詳細にモデル化し、
熱回収システムや熱源システムの設計、熱負荷計算、
およびエネルギー消費量予測を行うための重要な役割を果たします。
*/
func Hexcfv(Hex []*HEX) {
	for _, hex := range Hex {

		// 計算準備
		if hex.Id == 0 {
			/* 温度効率固定タイプと変動タイプの判定 */
			if hex.Cat.eff > 0.0 {
				hex.Etype = 'e'
			} else if hex.Cat.KA > 0.0 {
				hex.Etype = 'k'
			} else {
				fmt.Printf("Hex %s  Undefined Character eff or KA\n", hex.Name)
				hex.Etype = '-'
			}

			hex.Id = 1
		}

		if hex.Cmp.Control != OFF_SW {
			hex.Eff = hex.Cat.eff

			if hex.Eff < 0.0 {
				errMsg := fmt.Sprintf("Name=%s  eff=%.4g", hex.Cmp.Name, hex.Eff)
				Eprint("Hexcfv", errMsg)
			}

			eoh := hex.Cmp.Elouts[1]
			eoc := hex.Cmp.Elouts[0]
			hex.CGc = Spcheat(eoc.Fluid) * eoc.G
			hex.CGh = Spcheat(eoh.Fluid) * eoh.G

			if hex.Etype == 'k' {
				hex.Eff = FNhccet(hex.CGc, hex.CGh, hex.Cat.KA)
			}

			eCGmin := hex.Eff * math.Min(hex.CGc, hex.CGh)
			hex.ECGmin = eCGmin
			eoc.Coeffin[0] = -hex.CGc + eCGmin
			eoc.Coeffin[1] = -eCGmin
			eoc.Coeffo = hex.CGc
			eoc.Co = 0.0

			eoh.Coeffin[0] = -eCGmin
			eoh.Coeffin[1] = -hex.CGh + eCGmin
			eoh.Coeffo = hex.CGh
			eoh.Co = 0.0
		}
	}
}

/* --------------------------- */

/*
Hexene (Heat Exchanger Energy Calculation)

この関数は、熱交換器が冷媒側および温媒側に供給または除去する熱量を計算します。
これは、熱交換器の熱交換性能や、システム全体のエネルギー収支を評価する上で不可欠です。

建築環境工学的な観点:
- **冷媒側熱量 (hex.Qci)**:
  冷媒側流体が熱交換器を通過する際に得られる熱量です。
  `hex.CGc * (hex.Cmp.Elouts[0].Sysv - hex.Tcin)` のように、
  冷媒側の熱容量流量と入口・出口温度差から計算されます。
  これは、冷媒側が熱を吸収する能力を示します。
- **温媒側熱量 (hex.Qhi)**:
  温媒側流体が熱交換器を通過する際に失う熱量です。
  `hex.CGh * (hex.Cmp.Elouts[1].Sysv - hex.Thin)` のように、
  温媒側の熱容量流量と入口・出口温度差から計算されます。
  これは、温媒側が熱を放出する能力を示します。
- **熱回収の評価**: 熱交換器は、温媒側から冷媒側へ熱を移動させることで、
  エネルギーを回収します。
  `hex.Qci`と`hex.Qhi`は、この熱回収量を定量的に評価するために用いられます。
  理想的な熱交換器では、`hex.Qci`と`hex.Qhi`は等しくなりますが、
  実際の熱交換器では熱損失などにより差が生じます。
- **システム全体のエネルギー効率**: 熱交換器による熱回収は、
  熱源設備の負荷を軽減し、システム全体のエネルギー消費量を削減します。
  これらの熱量を正確に把握することで、
  熱回収システムの省エネルギー効果を評価し、
  システム全体のエネルギー効率を向上させるための設計検討に役立ちます。

この関数は、熱交換器の熱交換性能を定量的に評価し、
熱回収システムや熱源システムの設計、熱負荷計算、
およびエネルギー消費量予測を行うための重要な役割を果たします。
*/
func Hexene(Hex []*HEX) {
	for _, hex := range Hex {

		// 流入
		hex.Tcin = hex.Cmp.Elins[0].Sysvin
		hex.Thin = hex.Cmp.Elins[1].Sysvin

		if hex.Cmp.Control != OFF_SW {
			// 流出
			hex.Qci = hex.CGc * (hex.Cmp.Elouts[0].Sysv - hex.Tcin)
			hex.Qhi = hex.CGh * (hex.Cmp.Elouts[1].Sysv - hex.Thin)
		} else {
			hex.Qci = 0.0
			hex.Qhi = 0.0
		}
	}
}

/* --------------------------- */

func hexprint(fo io.Writer, id int, Hex []*HEX) {
	switch id {
	case 0:
		if len(Hex) > 0 {
			fmt.Fprintf(fo, "%s %d\n", HEXCHANGR_TYPE, len(Hex))
		}
		for _, hex := range Hex {
			fmt.Fprintf(fo, " %s 1 9\n", hex.Name)
		}
	case 1:
		for _, hex := range Hex {
			fmt.Fprintf(fo, "%s_c c c %s:c_G m f %s:c_Ti t f %s:c_To t f %s:c_Q q f\n",
				hex.Name, hex.Name, hex.Name, hex.Name, hex.Name)
			fmt.Fprintf(fo, "%s:h_G m f %s:h_Ti t f %s:h_To t f %s:h_Q q f\n",
				hex.Name, hex.Name, hex.Name, hex.Name)
		}
	default:
		for _, hex := range Hex {
			eo_Tc := hex.Cmp.Elouts[0]
			eo_Th := hex.Cmp.Elouts[1]
			fmt.Fprintf(fo, "%c %6.4g %4.1f %4.1f %2.0f", hex.Cmp.Control, eo_Tc.G, hex.Tcin, eo_Tc.Sysv, hex.Qci)
			fmt.Fprintf(fo, " %6.4g %4.1f %4.1f %2.0f\n", eo_Th.G, hex.Thin, eo_Th.Sysv, hex.Qhi)
		}
	}
}

/*
hexdyint (Heat Exchanger Daily Integration Initialization)

この関数は、熱交換器の日積算値（日ごとの入口温度、熱量など）をリセットします。
これは、日単位での熱交換器の運転状況や熱交換量を集計し、
熱回収システムの性能評価やエネルギー消費量の分析に用いられます。

建築環境工学的な観点:
- **日単位の性能評価**: 熱交換器の運転状況は、日中の熱負荷変動に応じて大きく変化します。
  日積算値を集計することで、日ごとの熱回収量、
  熱交換器の稼働時間、部分負荷運転の割合などを把握できます。
  これにより、特定の日の熱回収特性を分析したり、
  熱交換器の運転効率を日単位で評価したりすることが可能になります。
- **運用改善の指標**: 日積算データは、熱回収システムの運用改善のための重要な指標となります。
  例えば、外気温度や熱源温度などの気象条件と熱交換器の熱交換量の関係を分析したり、
  設定温度や流量などの運用条件が熱交換器の性能に与える影響を評価したりすることで、
  より効率的な運転方法を見つけることができます。
- **データ集計の準備**: この関数は、新しい日の集計を開始する前に、
  前日のデータをクリアする役割を担います。
  `svdyint`や`qdyint`といった関数は、
  それぞれ温度、熱量などの日積算値をリセットするためのものです。

この関数は、熱交換器の運転状況と熱交換量を日単位で詳細に分析し、
熱回収システムの運用改善や省エネルギー対策の効果評価を行うための基礎的な役割を果たします。
*/
func hexdyint(Hex []*HEX) {
	for _, hex := range Hex {
		svdyint(&hex.Tcidy)
		svdyint(&hex.Thidy)
		qdyint(&hex.Qcidy)
		qdyint(&hex.Qhidy)
	}
}

/*
hexmonint (Heat Exchanger Monthly Integration Initialization)

この関数は、熱交換器の月積算値（月ごとの入口温度、熱量など）をリセットします。
これは、月単位での熱交換器の運転状況や熱交換量を集計し、
熱回収システムの性能評価やエネルギー消費量の分析に用いられます。

建築環境工学的な観点:
- **月単位の性能評価**: 熱交換器の運転状況は、月単位で変動します。
  月積算値を集計することで、月ごとの熱回収量、
  熱交換器の稼働時間、部分負荷運転の割合などを把握できます。
  これにより、特定の月の熱回収特性を分析したり、
  熱交換器の運転効率を月単位で評価したりすることが可能になります。
- **運用改善の指標**: 月積算データは、熱回収システムの運用改善のための重要な指標となります。
  例えば、季節ごとの熱回収量の傾向を把握したり、
  月ごとの気象条件と熱交換器の熱交換量の関係を分析したりすることで、
  より効率的な運転方法を見つけることができます。
- **データ集計の準備**: この関数は、新しい月の集計を開始する前に、
  前月のデータをクリアする役割を担います。
  `svdyint`や`qdyint`といった関数は、
  それぞれ温度、熱量などの月積算値をリセットするためのものです。

この関数は、熱交換器の運転状況と熱交換量を月単位で詳細に分析し、
熱回収システムの運用改善や省エネルギー対策の効果評価を行うための基礎的な役割を果たします。
*/
func hexmonint(Hex []*HEX) {
	for _, hex := range Hex {
		svdyint(&hex.MTcidy)
		svdyint(&hex.MThidy)
		qdyint(&hex.MQcidy)
		qdyint(&hex.MQhidy)
	}
}

/*
hexday (Heat Exchanger Daily and Monthly Data Aggregation)

この関数は、熱交換器の運転データ（入口温度、熱量など）を、
日単位および月単位で集計します。
これにより、熱交換器の性能評価やエネルギー消費量の分析が可能になります。

建築環境工学的な観点:
- **日次集計 (svdaysum, qdaysum)**:
  日次集計は、熱交換器の運転状況を日単位で詳細に把握するために重要です。
  例えば、特定の日の熱負荷変動に対する熱交換器の応答、
  あるいは日中のピーク負荷時の熱交換量などを分析できます。
  これにより、日ごとの運用改善点を見つけ出すことが可能になります。
- **月次集計 (svmonsum, qmonsum)**:
  月次集計は、季節ごとの熱負荷変動や、
  熱交換器の年間を通じた熱交換量の傾向を把握するために重要です。
  これにより、年間を通じた省エネルギー対策の効果を評価したり、
  熱交換量の予測精度を向上させたりすることが可能になります。
- **データ分析の基礎**: この関数で集計されるデータは、
  熱交換器の性能評価、熱交換量のベンチマーキング、
  省エネルギー対策の効果検証、そして運用改善のための意思決定の基礎となります。

この関数は、熱交換器の運転状況と熱交換量を多角的に分析し、
熱回収システムの運用改善や省エネルギー対策の効果評価を行うための重要なデータ集計機能を提供します。
*/
func hexday(Mon, Day, ttmm int, Hex []*HEX, Nday, SimDayend int) {
	for _, hex := range Hex {
		// 日集計
		svdaysum(int64(ttmm), hex.Cmp.Control, hex.Tcin, &hex.Tcidy)
		svdaysum(int64(ttmm), hex.Cmp.Control, hex.Thin, &hex.Thidy)
		qdaysum(int64(ttmm), hex.Cmp.Control, hex.Qci, &hex.Qcidy)
		qdaysum(int64(ttmm), hex.Cmp.Control, hex.Qhi, &hex.Qhidy)

		// 月集計
		svmonsum(Mon, Day, ttmm, hex.Cmp.Control, hex.Tcin, &hex.MTcidy, Nday, SimDayend)
		svmonsum(Mon, Day, ttmm, hex.Cmp.Control, hex.Thin, &hex.MThidy, Nday, SimDayend)
		qmonsum(Mon, Day, ttmm, hex.Cmp.Control, hex.Qci, &hex.MQcidy, Nday, SimDayend)
		qmonsum(Mon, Day, ttmm, hex.Cmp.Control, hex.Qhi, &hex.MQhidy, Nday, SimDayend)
	}
}

func hexdyprt(fo io.Writer, id int, Hex []*HEX) {
	var c byte

	switch id {
	case 0:
		if len(Hex) > 0 {
			fmt.Fprintf(fo, "%s %d\n", HEXCHANGR_TYPE, len(Hex))
		}
		for _, hex := range Hex {
			fmt.Fprintf(fo, " %s 1 28\n", hex.Name)
		}
	case 1:
		for _, hex := range Hex {
			for j := 0; j < 2; j++ {
				if j == 0 {
					c = 'c'
				} else {
					c = 'h'
				}
				fmt.Fprintf(fo, "%s:%c_Ht H d %s:%c_T T f ", hex.Name, c, hex.Name, c)
				fmt.Fprintf(fo, "%s:%c_ttn h d %s:%c_Tn t f %s:%c_ttm h d %s:%c_Tm t f\n",
					hex.Name, c, hex.Name, c, hex.Name, c, hex.Name, c)
				fmt.Fprintf(fo, "%s:%c_Hh H d %s:%c_Qh Q f %s:%c_Hc H d %s:%c_Qc Q f\n",
					hex.Name, c, hex.Name, c, hex.Name, c, hex.Name, c)
				fmt.Fprintf(fo, "%s:%c_th h d %s:%c_qh q f %s:%c_tc h d %s:%c_qc q f\n",
					hex.Name, c, hex.Name, c, hex.Name, c, hex.Name, c)
			}
		}
	default:
		for _, hex := range Hex {
			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ",
				hex.Tcidy.Hrs, hex.Tcidy.M,
				hex.Tcidy.Mntime, hex.Tcidy.Mn,
				hex.Tcidy.Mxtime, hex.Tcidy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f ", hex.Qcidy.Hhr, hex.Qcidy.H)
			fmt.Fprintf(fo, "%1d %3.1f ", hex.Qcidy.Chr, hex.Qcidy.C)
			fmt.Fprintf(fo, "%1d %2.0f ", hex.Qcidy.Hmxtime, hex.Qcidy.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f ", hex.Qcidy.Cmxtime, hex.Qcidy.Cmx)
			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ",
				hex.Thidy.Hrs, hex.Thidy.M,
				hex.Thidy.Mntime, hex.Thidy.Mn,
				hex.Thidy.Mxtime, hex.Thidy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f ", hex.Qhidy.Hhr, hex.Qhidy.H)
			fmt.Fprintf(fo, "%1d %3.1f ", hex.Qhidy.Chr, hex.Qhidy.C)
			fmt.Fprintf(fo, "%1d %2.0f ", hex.Qhidy.Hmxtime, hex.Qhidy.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f\n", hex.Qhidy.Cmxtime, hex.Qhidy.Cmx)
		}
	}
}

func hexmonprt(fo io.Writer, id int, Hex []*HEX) {
	var c byte

	switch id {
	case 0:
		if len(Hex) > 0 {
			fmt.Fprintf(fo, "%s %d\n", HEXCHANGR_TYPE, len(Hex))
		}
		for _, hex := range Hex {
			fmt.Fprintf(fo, " %s 1 28\n", hex.Name)
		}
	case 1:
		for _, hex := range Hex {
			for j := 0; j < 2; j++ {
				if j == 0 {
					c = 'c'
				} else {
					c = 'h'
				}
				fmt.Fprintf(fo, "%s:%c_Ht H d %s:%c_T T f ", hex.Name, c, hex.Name, c)
				fmt.Fprintf(fo, "%s:%c_ttn h d %s:%c_Tn t f %s:%c_ttm h d %s:%c_Tm t f\n",
					hex.Name, c, hex.Name, c, hex.Name, c, hex.Name, c)
				fmt.Fprintf(fo, "%s:%c_Hh H d %s:%c_Qh Q f %s:%c_Hc H d %s:%c_Qc Q f\n",
					hex.Name, c, hex.Name, c, hex.Name, c, hex.Name, c)
				fmt.Fprintf(fo, "%s:%c_th h d %s:%c_qh q f %s:%c_tc h d %s:%c_qc q f\n",
					hex.Name, c, hex.Name, c, hex.Name, c, hex.Name, c)
			}
		}
	default:
		for _, hex := range Hex {
			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ",
				hex.MTcidy.Hrs, hex.MTcidy.M,
				hex.MTcidy.Mntime, hex.MTcidy.Mn,
				hex.MTcidy.Mxtime, hex.MTcidy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f ", hex.MQcidy.Hhr, hex.MQcidy.H)
			fmt.Fprintf(fo, "%1d %3.1f ", hex.MQcidy.Chr, hex.MQcidy.C)
			fmt.Fprintf(fo, "%1d %2.0f ", hex.MQcidy.Hmxtime, hex.MQcidy.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f ", hex.MQcidy.Cmxtime, hex.MQcidy.Cmx)
			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ",
				hex.MThidy.Hrs, hex.MThidy.M,
				hex.MThidy.Mntime, hex.MThidy.Mn,
				hex.MThidy.Mxtime, hex.MThidy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f ", hex.MQhidy.Hhr, hex.MQhidy.H)
			fmt.Fprintf(fo, "%1d %3.1f ", hex.MQhidy.Chr, hex.MQhidy.C)
			fmt.Fprintf(fo, "%1d %2.0f ", hex.MQhidy.Hmxtime, hex.MQhidy.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f\n", hex.MQhidy.Cmxtime, hex.MQhidy.Cmx)
		}
	}
}

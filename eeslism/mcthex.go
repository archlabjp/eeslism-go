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

/* mcthex.c */

package eeslism

import (
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
)

/*
Thexdata (Total Heat Exchanger Data Input)

この関数は、全熱交換器の各種仕様（顕熱交換効率、潜熱交換効率など）を読み込み、
対応する全熱交換器の構造体に格納します。
これらのデータは、換気システムにおける熱回収の性能評価、
熱負荷への対応、およびエネルギー消費量予測に不可欠です。

建築環境工学的な観点:
- **全熱交換器の役割**: 全熱交換器は、排気から顕熱（温度）と潜熱（湿度）の両方を回収し、
  給気に熱と湿気を供給することで、換気による熱損失を最小限に抑える機器です。
  特に、外気導入量が多い建物や、室内外の温湿度差が大きい地域において、
  省エネルギー効果が期待できます。
- **顕熱交換効率 (et)**:
  全熱交換器が顕熱をどれだけ効率的に交換できるかを示す指標です。
  排気と給気の温度差に対して、給気の温度がどれだけ排気に近づいたかを表します。
  `Thexca.et`が設定されている場合、全熱交換器の顕熱交換効率が固定値として扱われることを示唆します。
- **潜熱交換効率 (eh)**:
  全熱交換器が潜熱（水蒸気）をどれだけ効率的に交換できるかを示す指標です。
  排気と給気の絶対湿度差に対して、給気の絶対湿度がどれだけ排気に近づいたかを表します。
  `Thexca.eh`が設定されている場合、全熱交換器の潜熱交換効率が固定値として扱われることを示唆します。
- **熱回収の重要性**: 換気は、室内空気質を維持するために不可欠ですが、
  同時に熱損失（または熱取得）を伴います。
  全熱交換器は、この換気による熱損失を大幅に削減し、
  空調負荷を軽減することで、建物のエネルギー消費量を削減します。

この関数は、全熱交換器の性能をモデル化し、
換気システムにおける熱回収の設計、熱負荷計算、
およびエネルギー消費量予測を行うための重要なデータ入力機能を提供します。
*/
func Thexdata(s string, Thexca *THEXCA) int {
	var st int
	var dt float64
	var id int

	if st = strings.IndexRune(s, '='); st == -1 {
		Thexca.Name = s
	} else {
		stval := strings.Replace(s[st:], "=", "", 1)
		dt, _ = strconv.ParseFloat(stval, 64)

		if s == "et" {
			Thexca.et = dt
		} else if s == "eh" {
			Thexca.eh = dt
		} else {
			id = 1
		}
	}

	return id
}

/* ------------------------------------------------------ */

/*
Thexint (Total Heat Exchanger Initialization)

この関数は、全熱交換器の初期設定を行います。
特に、全熱交換器が顕熱交換のみを行うタイプか、
顕熱と潜熱の両方を行うタイプかを判定し、
初期の絶対湿度を設定します。

建築環境工学的な観点:
- **全熱交換器のタイプ判定**: 全熱交換器には、
  顕熱交換のみを行うタイプ（温度交換効率のみを持つ）と、
  顕熱と潜熱の両方を行うタイプ（温度交換効率とエンタルピー交換効率を持つ）があります。
  `thex.Cat.eh < 0.0` の条件は、
  潜熱交換効率が設定されていない場合に、
  全熱交換器を顕熱交換のみを行うタイプ（`thex.Type = 't'`）として扱います。
  それ以外の場合は、顕熱と潜熱の両方を行うタイプ（`thex.Type = 'h'`）として扱います。
  この判定は、全熱交換器の熱湿気交換特性を正確にモデル化するために重要です。
- **効率の妥当性チェック**: `if thex.Cat.et < 0.0` のように、
  顕熱交換効率が設定されていない場合にエラーメッセージを出力することで、
  入力データの不備を早期に発見し、シミュレーションの信頼性を確保します。
- **初期絶対湿度の設定**: `thex.Xeinold`などの絶対湿度を初期化します。
  これは、シミュレーション開始時の全熱交換器の熱湿気交換プロセスを正確にモデル化するために重要です。

この関数は、全熱交換器の熱湿気交換特性を初期設定し、
換気システムにおける熱回収の設計、熱負荷計算、
およびエネルギー消費量予測を行うための基礎的な役割を果たします。
*/
func Thexint(Thex []*THEX) {
	for _, thex := range Thex {
		if thex.Cat.eh < 0.0 {
			thex.Type = 't'
			thex.Cat.eh = 0.0
		} else {
			thex.Type = 'h'
		}

		if thex.Cat.et < 0.0 {
			s := fmt.Sprintf("Name=%s catname=%s et=%f", thex.Name, thex.Cat.Name, thex.Cat.et)
			Eprint("<Thexint>", s)
		}

		thex.Xeinold = FNXtr(26.0, 50.0)
		thex.Xeoutold = thex.Xeinold
		thex.Xoinold = thex.Xeinold
		thex.Xooutold = thex.Xeinold
	}
}

/*  全熱交換器出口空気温湿度に関する変数割当  */

/*
Thexelm (Total Heat Exchanger Element Matrix Setup)

この関数は、全熱交換器における要素方程式の変数を設定し、
特に排気系統と給気系統の空気温度と絶対湿度の相互関係をモデル化します。
これは、全熱交換器の熱湿気同時交換プロセスをシミュレーションするために不可欠です。

建築環境工学的な観点:
- **熱湿気同時交換のモデル化**: 全熱交換器は、
  排気と給気の間で顕熱と潜熱の両方を交換します。
  この関数は、出口空気温度の計算が入口空気絶対湿度に依存し、
  出口空気絶対湿度の計算が入口空気温度に依存するという、
  これらの変数の相互関係を要素方程式に組み込むための設定を行います。
  `Upo`や`Upv`といった変数は、
  要素方程式における他の変数の影響を考慮するためのポインターを示唆します。
- **排気系統と給気系統の連携**: 全熱交換器は、
  排気系統と給気系統の熱湿気を同時に交換します。
  この関数は、それぞれの系統の入出力変数を適切に割り当てることで、
  両系統間の熱湿気交換を正確にモデル化します。
- **システムモデルの構築**: この設定は、
  全熱交換器を構成する各要素の熱湿気交換特性を、
  連立方程式として解くための基礎となります。
  これにより、全熱交換器が様々な運転条件下で、
  空気の温度と湿度をどのように変化させるかを正確に予測できます。

この関数は、全熱交換器の熱湿気同時交換プロセスをモデル化し、
換気システムにおける熱回収の設計、熱負荷計算、
およびエネルギー消費量予測を行うための重要な役割を果たします。
*/
func Thexelm(Thex []*THEX) {
	var E, E1, E2, E3 *ELOUT
	var elin, elin2 *ELIN

	for _, thex := range Thex {
		E = thex.Cmp.Elouts[0]
		E1 = thex.Cmp.Elouts[1]
		E2 = thex.Cmp.Elouts[2]
		E3 = thex.Cmp.Elouts[3]

		// Tein variable assignment
		// E: Teout calculation, elin2: Tein
		elin2 = E.Elins[0]

		// E+2: Toout calculation, elin: Tein
		elin = E2.Elins[1]
		elin.Upo = elin2.Upo
		elin.Upv = elin2.Upo

		if thex.Cat.eh > 0.0 {
			// E+1: xeout calculation, elin:
			elin = E1.Elins[1]
			elin.Upo = elin2.Upo
			elin.Upv = elin2.Upo

			elin = E3.Elins[3]
			elin.Upo = elin2.Upo
			elin.Upv = elin2.Upo
		}

		// Toin variable assignment
		elin2 = E.Elins[1]

		elin = E2.Elins[0]
		elin.Upo = elin2.Upo
		elin.Upv = elin2.Upo

		if thex.Cat.eh > 0.0 {
			elin = E1.Elins[3]
			elin.Upo = elin2.Upo
			elin.Upv = elin2.Upo

			elin = E3.Elins[1]
			elin.Upo = elin2.Upo
			elin.Upv = elin2.Upo

			// Teoutの変数割り当て
			elin = E1.Elins[2]
			elin.Upo = E
			elin.Upv = E2

			// Tooutの割り当て
			elin = E3.Elins[2]
			elin.Upo = E2
			elin.Upv = E2

			// xeinの割り当て
			elin = E1.Elins[0]
			elin2 = E3.Elins[4]
			elin2.Upo = elin.Upo
			elin2.Upv = elin.Upo

			// xoinの割り当て
			elin = E1.Elins[4]
			elin2 = E3.Elins[0]
			elin2.Upo = elin.Upo
			elin2.Upv = elin.Upo
		}
	}
}

/* ------------------------------------------------------ */

//
//  [IN 1] --(E)-->  +------+ --(E)--> [OUT 1] 排気系統（温度）
//  [IN 2] --(e)-->  |      | --(e)--> [OUT 2] 排気系統（エンタルピー）
//                   | THEX |
//  [IN 3] --(O)-->  |      | --(O)--> [OUT 3] 給気系統（温度）
//  [IN 4] --(o)-->  +------+ --(o)--> [OUT 4] 給気系統（エンタルピー）
//


/*
Thexcfv (Total Heat Exchanger Characteristic Function Value Calculation)

この関数は、全熱交換器の運転特性を評価し、
排気系統と給気系統の熱容量流量、そして顕熱交換効率と潜熱交換効率を計算します。
これは、全熱交換器が空気の熱湿気をどのように処理するかをモデル化するために不可欠です。

建築環境工学的な観点:
- **熱容量流量 (thex.CGe, thex.CGo)**:
  排気系統（`thex.CGe`）と給気系統（`thex.CGo`）の熱容量流量は、
  全熱交換器が処理できる熱量に影響します。
  `Spcheat(Fluid) * G` のように、空気の比熱と質量流量から計算されます。
- **顕熱交換効率 (thex.ET)**:
  全熱交換器の顕熱交換能力を示す指標です。
  `thex.Cat.et`から設定されます。
- **潜熱交換効率 (thex.EH)**:
  全熱交換器の潜熱交換能力を示す指標です。
  `thex.Cat.eh`から設定されます。
- **有効熱容量流量 (etCGmin, ehGmin)**:
  - `etCGmin = thex.ET * math.Min(thex.CGe, thex.CGo)` は、
    顕熱交換器が実際に交換できる顕熱量の上限を示します。
  - `ehGmin = thex.EH * math.Min(thex.Ge, thex.Go)` は、
    潜熱交換器が実際に交換できる潜熱量の上限を示します。
  これらの値は、全熱交換器の効率と、熱容量流量の小さい方の値によって決定されます。
- **出口温度・湿度の係数設定**: 計算された全熱交換器の特性に基づいて、
  排気系統の出口温度（`Eoet`）、出口エンタルピー（`Eoex`）、
  給気系統の出口温度（`Eoot`）、出口エンタルピー（`Eoox`）に関する係数（`Coeffo`, `Co`, `Coeffin`）を設定します。
  これらの係数は、換気システム全体の熱収支方程式に組み込まれ、
  各系統の出口温度や湿度を予測するために用いられます。

この関数は、全熱交換器の熱湿気同時交換特性を詳細にモデル化し、
換気システムにおける熱回収の設計、熱負荷計算、
およびエネルギー消費量予測を行うための重要な役割を果たします。
*/
func Thexcfv(Thex []*THEX) {
	var Eoet, Eoot, Eoex, Eoox *ELOUT
	var etCGmin, ehGmin, Aeout, Aein, Aoout, Aoin float64

	for _, thex := range Thex {
		if thex.Cmp.Control != OFF_SW {
			thex.ET = thex.Cat.et
			thex.EH = thex.Cat.eh

			Eoet = thex.Cmp.Elouts[0] // 排気系統（温度）
			Eoex = thex.Cmp.Elouts[1] // 排気系統（エンタルピー）
			Eoot = thex.Cmp.Elouts[2] // 給気系統（温度）
			Eoox = thex.Cmp.Elouts[3] // 給気系統（エンタルピー）

			thex.Ge = Eoet.G
			thex.Go = Eoot.G

			if DEBUG {
				fmt.Printf("<Thexcfv>  %s Ge=%f Go=%f\n", thex.Cmp.Name, thex.Ge, thex.Go)
			}

			thex.CGe = Spcheat(Eoet.Fluid) * thex.Ge
			thex.CGo = Spcheat(Eoot.Fluid) * thex.Go
			etCGmin = thex.ET * math.Min(thex.CGe, thex.CGo)
			ehGmin = thex.EH * math.Min(thex.Ge, thex.Go)

			Aein = Ca + Cv*thex.Xeinold
			Aeout = Ca + Cv*thex.Xeoutold
			Aoin = Ca + Cv*thex.Xoinold
			Aoout = Ca + Cv*thex.Xooutold

			// 排気系統（温度）の熱収支
			Eoet.Coeffo = thex.CGe
			Eoet.Co = 0.0
			cfin := Eoet.Coeffin
			cfin[0] = etCGmin - thex.CGe
			cfin[1] = -etCGmin

			// 給気系統（温度）の熱収支
			Eoot.Coeffo = thex.CGo
			Eoot.Co = 0.0
			cfin = Eoot.Coeffin
			cfin[0] = etCGmin - thex.CGo
			cfin[1] = -etCGmin

			if thex.Type == 'h' {
				// 排気系統（エンタルピー）の熱収支
				Eoex.Coeffo = thex.Ge * Ro
				Eoex.Co = 0.0
				cfin = Eoex.Coeffin
				cfin[0] = Ro * (ehGmin - thex.Ge)
				cfin[1] = Aein * (ehGmin - thex.Ge)
				cfin[2] = Aeout * thex.Ge
				cfin[3] = -ehGmin * Aoin
				cfin[4] = -ehGmin * Ro

				// 給気系統（エンタルピー）の熱収支
				Eoox.Coeffo = thex.Go * Ro
				Eoox.Co = 0.0
				cfin = Eoox.Coeffin
				cfin[0] = Ro * (ehGmin - thex.Go)
				cfin[1] = Aoin * (ehGmin - thex.Go)
				cfin[2] = thex.Go * Aoout
				cfin[3] = -ehGmin * Aein
				cfin[4] = -ehGmin * Ro
			} else {
				Eoex.Coeffo = 1.0
				Eoex.Coeffin[0] = -1.0

				Eoox.Coeffo = 1.0
				Eoox.Coeffin[0] = -1.0
			}
		}
	}
}

/*
Thexene (Total Heat Exchanger Energy Calculation)

この関数は、全熱交換器が排気系統および給気系統で交換する顕熱量、潜熱量、および全熱量を計算します。
これは、全熱交換器の熱交換性能や、システム全体のエネルギー収支を評価する上で不可欠です。

建築環境工学的な観点:
- **排気系統の熱量 (thex.Qes, thex.Qel, thex.Qet)**:
  - `thex.Qes`: 排気系統で交換される顕熱量。
    空気の温度変化に伴う熱量であり、`Ca * thex.Ge * (thex.Teout - thex.Tein)` のように計算されます。
  - `thex.Qel`: 排気系統で交換される潜熱量。
    空気中の水蒸気量の変化に伴う熱量であり、`Ro * thex.Ge * (thex.Xeout - thex.Xein)` のように計算されます。
  - `thex.Qet`: 排気系統で交換される全熱量（顕熱と潜熱の合計）。
- **給気系統の熱量 (thex.Qos, thex.Qol, thex.Qot)**:
  - `thex.Qos`: 給気系統で交換される顕熱量。
  - `thex.Qol`: 給気系統で交換される潜熱量。
  - `thex.Qot`: 給気系統で交換される全熱量（顕熱と潜熱の合計）。
  これらの熱量は、排気系統から給気系統へ熱が回収されることを示します。
- **熱回収の評価**: 全熱交換器は、排気から熱を回収し、給気に供給することで、
  換気による熱損失を削減します。
  これらの熱量を定量的に評価することで、
  全熱交換器の省エネルギー効果を把握し、
  換気システムの設計最適化に役立てることができます。
- **絶対湿度の更新**: `thex.Xeinold = thex.Xein` のように、
  前時刻の絶対湿度を更新することで、
  次の時間ステップの計算に現在の状態を反映させ、
  より正確な動的シミュレーションを可能にします。

この関数は、全熱交換器の熱湿気同時交換性能を定量的に評価し、
換気システムにおける熱回収の設計、熱負荷計算、
およびエネルギー消費量予測を行うための重要な役割を果たします。
*/
func Thexene(Thex []*THEX) {
	for _, thex := range Thex {
		Eoet := thex.Cmp.Elouts[0] // 排気系統（温度）
		Eoex := thex.Cmp.Elouts[1] // 排気系統（エンタルピー）
		Eoot := thex.Cmp.Elouts[2] // 給気系統（温度）
		Eoox := thex.Cmp.Elouts[3] // 給気系統（エンタルピー）

		thex.Tein = Eoet.Elins[0].Upo.Sysv
		thex.Teout = Eoet.Sysv
		thex.Xein = Eoex.Elins[0].Upo.Sysv
		thex.Xeout = Eoex.Sysv

		thex.Toin = Eoot.Elins[0].Upo.Sysv
		thex.Toout = Eoot.Sysv
		thex.Xoin = Eoox.Elins[0].Upo.Sysv
		thex.Xoout = Eoox.Sysv

		thex.Hein = FNH(thex.Tein, thex.Xein)
		thex.Heout = FNH(thex.Teout, thex.Xeout)
		thex.Hoin = FNH(thex.Toin, thex.Xoin)
		thex.Hoout = FNH(thex.Toout, thex.Xoout)

		if thex.Cmp.Control != OFF_SW {
			// 交換熱量の計算
			thex.Qes = Ca * thex.Ge * (thex.Teout - thex.Tein)
			thex.Qel = Ro * thex.Ge * (thex.Xeout - thex.Xein)
			thex.Qet = thex.Qes + thex.Qel

			thex.Qos = Ca * thex.Go * (thex.Toout - thex.Toin)
			thex.Qol = Ro * thex.Go * (thex.Xoout - thex.Xoin)
			thex.Qot = thex.Qos + thex.Qol

			// 前時刻の絶対湿度の入れ替え
			thex.Xeinold = thex.Xein
			thex.Xeoutold = thex.Xeout
			thex.Xoinold = thex.Xoin
			thex.Xooutold = thex.Xoout
		} else {
			thex.Qes = 0.0
			thex.Qel = 0.0
			thex.Qet = 0.0
			thex.Qos = 0.0
			thex.Qol = 0.0
			thex.Qot = 0.0
			thex.Ge = 0.0
			thex.Tein = 0.0
			thex.Teout = 0.0
			thex.Xein = 0.0
			thex.Xeout = 0.0
			thex.Hein = 0.0
			thex.Heout = 0.0
			thex.Go = 0.0
			thex.Toin = 0.0
			thex.Toout = 0.0
			thex.Xoin = 0.0
			thex.Xoout = 0.0
			thex.Hoin = 0.0
			thex.Hoout = 0.0
		}
	}
}

func Thexprint(fo io.Writer, id int, Thex []*THEX) {
	var el *ELOUT

	switch id {
	case 0:
		if len(Thex) > 0 {
			fmt.Fprintf(fo, "%s %d\n", THEX_TYPE, len(Thex))
		}
		for _, thex := range Thex {
			fmt.Fprintf(fo, " %s 1 22\n", thex.Name)
		}

	case 1:
		for _, thex := range Thex {
			fmt.Fprintf(fo, "%s_ce c c %s_Ge m f %s_Tei t f %s_Teo t f %s_xei t f %s_xeo t f\n",
				thex.Name, thex.Name, thex.Name, thex.Name, thex.Name, thex.Name)
			fmt.Fprintf(fo, "%s_hei h f %s_heo h f %s_Qes q f %s_Qel q f %s_Qet q f\n",
				thex.Name, thex.Name, thex.Name, thex.Name, thex.Name)

			fmt.Fprintf(fo, "%s_co c c %s_Go m f %s_Toi t f %s_Too t f %s_xoi t f %s_xoo t f\n",
				thex.Name, thex.Name, thex.Name, thex.Name, thex.Name, thex.Name)
			fmt.Fprintf(fo, "%s_hoi h f %s_hoo h f %s_Qos q f %s_Qol q f %s_Qot q f\n",
				thex.Name, thex.Name, thex.Name, thex.Name, thex.Name)
		}

	default:
		for _, thex := range Thex {
			el = thex.Cmp.Elouts[0]
			fmt.Fprintf(fo, "%c %6.4g %4.2f %4.2f %.4f %.4f ",
				el.Control, thex.Ge, thex.Tein, thex.Teout, thex.Xein, thex.Xeout)
			fmt.Fprintf(fo, "%.0f %.0f %.2f %.2f %.2f\n",
				thex.Hein, thex.Heout, thex.Qes, thex.Qel, thex.Qet)

			el = thex.Cmp.Elouts[2]
			fmt.Fprintf(fo, "%c %6.4g %4.2f %4.2f %.4f %.4f ",
				el.Control, thex.Go, thex.Toin, thex.Toout, thex.Xoin, thex.Xoout)
			fmt.Fprintf(fo, "%.0f %.0f %.2f %.2f %.2f\n",
				thex.Hoin, thex.Hoout, thex.Qos, thex.Qol, thex.Qot)
		}
	}
}

/*
Thexdyint (Total Heat Exchanger Daily Integration Initialization)

この関数は、全熱交換器の日積算値（日ごとの入口・出口空気温度、絶対湿度、
顕熱量、潜熱量、全熱量など）をリセットします。
これは、日単位での全熱交換器の運転状況や熱湿気交換量を集計し、
換気システムの性能評価やエネルギー消費量の分析に用いられます。

建築環境工学的な観点:
- **日単位の性能評価**: 全熱交換器の運転状況は、日中の熱湿気負荷変動に応じて大きく変化します。
  日積算値を集計することで、日ごとの熱回収量、
  全熱交換器の稼働時間、部分負荷運転の割合などを把握できます。
  これにより、特定の日の熱回収特性を分析したり、
  全熱交換器の運転効率を日単位で評価したりすることが可能になります。
- **運用改善の指標**: 日積算データは、換気システムの運用改善のための重要な指標となります。
  例えば、外気温度や湿度などの気象条件と全熱交換器の熱湿気交換量の関係を分析したり、
  設定温度や換気量などの運用条件が全熱交換器の性能に与える影響を評価したりすることで、
  より効率的な運転方法を見つけることができます。
- **データ集計の準備**: この関数は、新しい日の集計を開始する前に、
  前日のデータをクリアする役割を担います。
  `svdyint`や`qdyint`といった関数は、
  それぞれ温度、湿度、熱量などの日積算値をリセットするためのものです。

この関数は、全熱交換器の運転状況と熱湿気交換量を日単位で詳細に分析し、
換気システムの運用改善や省エネルギー対策の効果評価を行うための基礎的な役割を果たします。
*/
func Thexdyint(Thex []*THEX) {
	for _, thex := range Thex {
		svdyint(&thex.Teidy)
		svdyint(&thex.Teody)
		svdyint(&thex.Xeidy)
		svdyint(&thex.Xeody)

		svdyint(&thex.Toidy)
		svdyint(&thex.Toody)
		svdyint(&thex.Xoidy)
		svdyint(&thex.Xoody)

		qdyint(&thex.Qdyes)
		qdyint(&thex.Qdyel)
		qdyint(&thex.Qdyet)

		qdyint(&thex.Qdyos)
		qdyint(&thex.Qdyol)
		qdyint(&thex.Qdyot)
	}
}

/*
Thexmonint (Total Heat Exchanger Monthly Integration Initialization)

この関数は、全熱交換器の月積算値（月ごとの入口・出口空気温度、絶対湿度、
顕熱量、潜熱量、全熱量など）をリセットします。
これは、月単位での全熱交換器の運転状況や熱湿気交換量を集計し、
換気システムの性能評価やエネルギー消費量の分析に用いられます。

建築環境工学的な観点:
- **月単位の性能評価**: 全熱交換器の運転状況は、月単位で変動します。
  月積算値を集計することで、月ごとの熱回収量、
  全熱交換器の稼働時間、部分負荷運転の割合などを把握できます。
  これにより、特定の月の熱回収特性を分析したり、
  全熱交換器の運転効率を月単位で評価したりすることが可能になります。
- **運用改善の指標**: 月積算データは、換気システムの運用改善のための重要な指標となります。
  例えば、季節ごとの熱回収量の傾向を把握したり、
  月ごとの気象条件と全熱交換器の熱湿気交換量の関係を分析したりすることで、
  より効率的な運転方法を見つけることができます。
- **データ集計の準備**: この関数は、新しい月の集計を開始する前に、
  前月のデータをクリアする役割を担います。
  `svdyint`や`qdyint`といった関数は、
  それぞれ温度、湿度、熱量などの月積算値をリセットするためのものです。

この関数は、全熱交換器の運転状況と熱湿気交換量を月単位で詳細に分析し、
換気システムの運用改善や省エネルギー対策の効果評価を行うための基礎的な役割を果たします。
*/
func Thexmonint(Thex []*THEX) {
	for _, thex := range Thex {
		svdyint(&thex.MTeidy)
		svdyint(&thex.MTeody)
		svdyint(&thex.MXeidy)
		svdyint(&thex.MXeody)

		svdyint(&thex.MToidy)
		svdyint(&thex.MToody)
		svdyint(&thex.MXoidy)
		svdyint(&thex.MXoody)

		qdyint(&thex.MQdyes)
		qdyint(&thex.MQdyel)
		qdyint(&thex.MQdyet)

		qdyint(&thex.MQdyos)
		qdyint(&thex.MQdyol)
		qdyint(&thex.MQdyot)
	}
}

/*
Thexday (Total Heat Exchanger Daily and Monthly Data Aggregation)

この関数は、全熱交換器の運転データ（入口・出口空気温度、絶対湿度、
顕熱量、潜熱量、全熱量など）を、日単位および月単位で集計します。
これにより、全熱交換器の性能評価やエネルギー消費量の分析が可能になります。

建築環境工学的な観点:
- **日次集計 (svdaysum, qdaysum)**:
  日次集計は、全熱交換器の運転状況を日単位で詳細に把握するために重要です。
  例えば、特定の日の熱湿気負荷変動に対する全熱交換器の応答、
  あるいは日中のピーク負荷時の熱湿気交換量などを分析できます。
  これにより、日ごとの運用改善点を見つけ出すことが可能になります。
- **月次集計 (svmonsum, qmonsum)**:
  月次集計は、季節ごとの熱湿気負荷変動や、
  全熱交換器の年間を通じた熱湿気交換量の傾向を把握するために重要です。
  これにより、年間を通じた省エネルギー対策の効果を評価したり、
  熱湿気交換量の予測精度を向上させたりすることが可能になります。
- **データ分析の基礎**: この関数で集計されるデータは、
  全熱交換器の性能評価、熱湿気交換量のベンチマーキング、
  省エネルギー対策の効果検証、そして運用改善のための意思決定の基礎となります。

この関数は、全熱交換器の運転状況と熱湿気交換量を多角的に分析し、
換気システムの運用改善や省エネルギー対策の効果評価を行うための重要なデータ集計機能を提供します。
*/
func Thexday(Mon, Day, ttmm int, Thex []*THEX, Nday, SimDayend int) {
	for _, thex := range Thex {
		// 日集計
		svdaysum(int64(ttmm), thex.Cmp.Control, thex.Tein, &thex.Teidy)
		svdaysum(int64(ttmm), thex.Cmp.Control, thex.Teout, &thex.Teody)
		svdaysum(int64(ttmm), thex.Cmp.Control, thex.Xein, &thex.Xeidy)
		svdaysum(int64(ttmm), thex.Cmp.Control, thex.Xeout, &thex.Xeody)

		svdaysum(int64(ttmm), thex.Cmp.Control, thex.Toin, &thex.Toidy)
		svdaysum(int64(ttmm), thex.Cmp.Control, thex.Toout, &thex.Toody)
		svdaysum(int64(ttmm), thex.Cmp.Control, thex.Xoin, &thex.Xoidy)
		svdaysum(int64(ttmm), thex.Cmp.Control, thex.Xoout, &thex.Xoody)

		qdaysum(int64(ttmm), thex.Cmp.Control, thex.Qes, &thex.Qdyes)
		qdaysum(int64(ttmm), thex.Cmp.Control, thex.Qel, &thex.Qdyel)
		qdaysum(int64(ttmm), thex.Cmp.Control, thex.Qet, &thex.Qdyet)

		qdaysum(int64(ttmm), thex.Cmp.Control, thex.Qos, &thex.Qdyos)
		qdaysum(int64(ttmm), thex.Cmp.Control, thex.Qol, &thex.Qdyol)
		qdaysum(int64(ttmm), thex.Cmp.Control, thex.Qot, &thex.Qdyot)

		// 月集計
		svmonsum(Mon, Day, ttmm, thex.Cmp.Control, thex.Tein, &thex.MTeidy, Nday, SimDayend)
		svmonsum(Mon, Day, ttmm, thex.Cmp.Control, thex.Teout, &thex.MTeody, Nday, SimDayend)
		svmonsum(Mon, Day, ttmm, thex.Cmp.Control, thex.Xein, &thex.MXeidy, Nday, SimDayend)
		svmonsum(Mon, Day, ttmm, thex.Cmp.Control, thex.Xeout, &thex.MXeody, Nday, SimDayend)

		svmonsum(Mon, Day, ttmm, thex.Cmp.Control, thex.Toin, &thex.MToidy, Nday, SimDayend)
		svmonsum(Mon, Day, ttmm, thex.Cmp.Control, thex.Toout, &thex.MToody, Nday, SimDayend)
		svmonsum(Mon, Day, ttmm, thex.Cmp.Control, thex.Xoin, &thex.MXoidy, Nday, SimDayend)
		svmonsum(Mon, Day, ttmm, thex.Cmp.Control, thex.Xoout, &thex.MXoody, Nday, SimDayend)

		qmonsum(Mon, Day, ttmm, thex.Cmp.Control, thex.Qes, &thex.MQdyes, Nday, SimDayend)
		qmonsum(Mon, Day, ttmm, thex.Cmp.Control, thex.Qel, &thex.MQdyel, Nday, SimDayend)
		qmonsum(Mon, Day, ttmm, thex.Cmp.Control, thex.Qet, &thex.MQdyet, Nday, SimDayend)

		qmonsum(Mon, Day, ttmm, thex.Cmp.Control, thex.Qos, &thex.MQdyos, Nday, SimDayend)
		qmonsum(Mon, Day, ttmm, thex.Cmp.Control, thex.Qol, &thex.MQdyol, Nday, SimDayend)
		qmonsum(Mon, Day, ttmm, thex.Cmp.Control, thex.Qot, &thex.MQdyot, Nday, SimDayend)
	}
}

func Thexdyprt(fo io.Writer, id int, Thex []*THEX) {
	switch id {
	case 0:
		if len(Thex) > 0 {
			fmt.Fprintf(fo, "%s %d\n", THEX_TYPE, len(Thex))
		}
		for _, thex := range Thex {
			fmt.Fprintf(fo, " %s 1 48\n", thex.Name)
		}
	case 1:
		for _, thex := range Thex {
			fmt.Fprintf(fo, "%s_Hte H d %s_Te T f ", thex.Name, thex.Name)
			fmt.Fprintf(fo, "%s_ttne h d %s_Ten t f %s_ttme h d %s_Tem t f\n",
				thex.Name, thex.Name, thex.Name, thex.Name)
			fmt.Fprintf(fo, "%s_Hto H d %s_To T f ", thex.Name, thex.Name)
			fmt.Fprintf(fo, "%s_ttno h d %s_Ton t f %s_ttmo h d %s_Tom t f\n",
				thex.Name, thex.Name, thex.Name, thex.Name)

			fmt.Fprintf(fo, "%s_Hxe H d %s_xe T f ", thex.Name, thex.Name)
			fmt.Fprintf(fo, "%s_txne h d %s_xen t f %s_txme h d %s_xem t f\n",
				thex.Name, thex.Name, thex.Name, thex.Name)
			fmt.Fprintf(fo, "%s_Hxo H d %s_xo T f ", thex.Name, thex.Name)
			fmt.Fprintf(fo, "%s_txno h d %s_xon t f %s_txmo h d %s_xom t f\n",
				thex.Name, thex.Name, thex.Name, thex.Name)

			fmt.Fprintf(fo, "%s_Hhs H d %s_Qsh Q f %s_Hcs H d %s_Qsc Q f\n",
				thex.Name, thex.Name, thex.Name, thex.Name)
			fmt.Fprintf(fo, "%s_ths h d %s_qsh q f %s_tcs h d %s_qsc q f\n\n",
				thex.Name, thex.Name, thex.Name, thex.Name)
			fmt.Fprintf(fo, "%s_Hhl H d %s_Qlh Q f %s_Hcl H d %s_Qlc Q f\n",
				thex.Name, thex.Name, thex.Name, thex.Name)
			fmt.Fprintf(fo, "%s_thl h d %s_qlh q f %s_tcl h d %s_qlc q f\n\n",
				thex.Name, thex.Name, thex.Name, thex.Name)
			fmt.Fprintf(fo, "%s_Hht H d %s_Qth Q f %s_Hct H d %s_Qtc Q f\n",
				thex.Name, thex.Name, thex.Name, thex.Name)
			fmt.Fprintf(fo, "%s_tht h d %s_qth q f %s_tct h d %s_qtc q f\n\n",
				thex.Name, thex.Name, thex.Name, thex.Name)
		}
	default:
		for _, thex := range Thex {
			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ",
				thex.Teidy.Hrs, thex.Teidy.M,
				thex.Teidy.Mntime, thex.Teidy.Mn,
				thex.Teidy.Mxtime, thex.Teidy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f\n",
				thex.Toidy.Hrs, thex.Toidy.M,
				thex.Toidy.Mntime, thex.Toidy.Mn,
				thex.Toidy.Mxtime, thex.Toidy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ",
				thex.Xeidy.Hrs, thex.Xeidy.M,
				thex.Xeidy.Mntime, thex.Xeidy.Mn,
				thex.Xeidy.Mxtime, thex.Xeidy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f\n",
				thex.Xoidy.Hrs, thex.Xoidy.M,
				thex.Xoidy.Mntime, thex.Xoidy.Mn,
				thex.Xoidy.Mxtime, thex.Xoidy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f ", thex.Qdyes.Hhr, thex.Qdyes.H)
			fmt.Fprintf(fo, "%1d %3.1f ", thex.Qdyes.Chr, thex.Qdyes.C)
			fmt.Fprintf(fo, "%1d %2.0f ", thex.Qdyes.Hmxtime, thex.Qdyes.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f\n", thex.Qdyes.Cmxtime, thex.Qdyes.Cmx)

			fmt.Fprintf(fo, "%1d %3.1f ", thex.Qdyel.Hhr, thex.Qdyel.H)
			fmt.Fprintf(fo, "%1d %3.1f ", thex.Qdyel.Chr, thex.Qdyel.C)
			fmt.Fprintf(fo, "%1d %2.0f ", thex.Qdyel.Hmxtime, thex.Qdyel.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f\n", thex.Qdyel.Cmxtime, thex.Qdyel.Cmx)

			fmt.Fprintf(fo, "%1d %3.1f ", thex.Qdyet.Hhr, thex.Qdyet.H)
			fmt.Fprintf(fo, "%1d %3.1f ", thex.Qdyet.Chr, thex.Qdyet.C)
			fmt.Fprintf(fo, "%1d %2.0f ", thex.Qdyet.Hmxtime, thex.Qdyet.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f\n", thex.Qdyet.Cmxtime, thex.Qdyet.Cmx)
		}
	}
}
func Thexmonprt(fo io.Writer, id int, Thex []*THEX) {
	switch id {
	case 0:
		if len(Thex) > 0 {
			fmt.Fprintf(fo, "%s %d\n", THEX_TYPE, len(Thex))
		}
		for _, thex := range Thex {
			fmt.Fprintf(fo, " %s 1 48\n", thex.Name)
		}
	case 1:
		for _, thex := range Thex {
			fmt.Fprintf(fo, "%s_Hte H d %s_Te T f ", thex.Name, thex.Name)
			fmt.Fprintf(fo, "%s_ttne h d %s_Ten t f %s_ttme h d %s_Tem t f\n",
				thex.Name, thex.Name, thex.Name, thex.Name)
			fmt.Fprintf(fo, "%s_Hto H d %s_To T f ", thex.Name, thex.Name)
			fmt.Fprintf(fo, "%s_ttno h d %s_Ton t f %s_ttmo h d %s_Tom t f\n",
				thex.Name, thex.Name, thex.Name, thex.Name)

			fmt.Fprintf(fo, "%s_Hxe H d %s_xe T f ", thex.Name, thex.Name)
			fmt.Fprintf(fo, "%s_txne h d %s_xen t f %s_txme h d %s_xem t f\n",
				thex.Name, thex.Name, thex.Name, thex.Name)
			fmt.Fprintf(fo, "%s_Hxo H d %s_xo T f ", thex.Name, thex.Name)
			fmt.Fprintf(fo, "%s_txno h d %s_xon t f %s_txmo h d %s_xom t f\n",
				thex.Name, thex.Name, thex.Name, thex.Name)

			fmt.Fprintf(fo, "%s_Hhs H d %s_Qsh Q f %s_Hcs H d %s_Qsc Q f\n",
				thex.Name, thex.Name, thex.Name, thex.Name)
			fmt.Fprintf(fo, "%s_ths h d %s_qsh q f %s_tcs h d %s_qsc q f\n\n",
				thex.Name, thex.Name, thex.Name, thex.Name)
			fmt.Fprintf(fo, "%s_Hhl H d %s_Qlh Q f %s_Hcl H d %s_Qlc Q f\n",
				thex.Name, thex.Name, thex.Name, thex.Name)
			fmt.Fprintf(fo, "%s_thl h d %s_qlh q f %s_tcl h d %s_qlc q f\n\n",
				thex.Name, thex.Name, thex.Name, thex.Name)
			fmt.Fprintf(fo, "%s_Hht H d %s_Qth Q f %s_Hct H d %s_Qtc Q f\n",
				thex.Name, thex.Name, thex.Name, thex.Name)
			fmt.Fprintf(fo, "%s_tht h d %s_qth q f %s_tct h d %s_qtc q f\n\n",
				thex.Name, thex.Name, thex.Name, thex.Name)
		}
	default:
		for _, thex := range Thex {
			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ",
				thex.MTeidy.Hrs, thex.MTeidy.M,
				thex.MTeidy.Mntime, thex.MTeidy.Mn,
				thex.MTeidy.Mxtime, thex.MTeidy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f\n",
				thex.MToidy.Hrs, thex.MToidy.M,
				thex.MToidy.Mntime, thex.MToidy.Mn,
				thex.MToidy.Mxtime, thex.MToidy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ",
				thex.MXeidy.Hrs, thex.MXeidy.M,
				thex.MXeidy.Mntime, thex.MXeidy.Mn,
				thex.MXeidy.Mxtime, thex.MXeidy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f\n",
				thex.MXoidy.Hrs, thex.MXoidy.M,
				thex.MXoidy.Mntime, thex.MXoidy.Mn,
				thex.MXoidy.Mxtime, thex.MXoidy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f ", thex.MQdyes.Hhr, thex.MQdyes.H)
			fmt.Fprintf(fo, "%1d %3.1f ", thex.MQdyes.Chr, thex.MQdyes.C)
			fmt.Fprintf(fo, "%1d %2.0f ", thex.MQdyes.Hmxtime, thex.MQdyes.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f\n", thex.MQdyes.Cmxtime, thex.MQdyes.Cmx)

			fmt.Fprintf(fo, "%1d %3.1f ", thex.MQdyel.Hhr, thex.MQdyel.H)
			fmt.Fprintf(fo, "%1d %3.1f ", thex.MQdyel.Chr, thex.MQdyel.C)
			fmt.Fprintf(fo, "%1d %2.0f ", thex.MQdyel.Hmxtime, thex.MQdyel.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f\n", thex.MQdyel.Cmxtime, thex.MQdyel.Cmx)

			fmt.Fprintf(fo, "%1d %3.1f ", thex.MQdyet.Hhr, thex.MQdyet.H)
			fmt.Fprintf(fo, "%1d %3.1f ", thex.MQdyet.Chr, thex.MQdyet.C)
			fmt.Fprintf(fo, "%1d %2.0f ", thex.MQdyet.Hmxtime, thex.MQdyet.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f\n", thex.MQdyet.Cmxtime, thex.MQdyet.Cmx)
		}
	}
}

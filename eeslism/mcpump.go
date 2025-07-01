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

/*  pump.c  */

package eeslism

import (
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
)


/*
Pumpdata (Pump/Fan Data Input)

この関数は、ポンプまたはファンの各種仕様（定格流量、定格消費電力、部分負荷特性など）を読み込み、
対応する構造体に格納します。
これらのデータは、熱搬送システムや空調システムにおけるポンプ・ファンの性能評価、
熱負荷への対応、およびエネルギー消費量予測に不可欠です。

建築環境工学的な観点:
- **熱搬送の動力**: ポンプやファンは、熱媒（水、空気）を搬送し、
  熱を必要な場所へ運ぶための動力源です。
  その性能を正確にモデル化することは、システム全体のエネルギー効率を評価する上で非常に重要です。
- **定格流量 (Go)**:
  ポンプやファンが搬送できる最大の流量を示します。
  建物の最大熱負荷や換気量に対して十分な能力があるか、
  あるいは複数台の機器を組み合わせる必要があるかを判断する際に用いられます。
- **定格消費電力 (Wo)**:
  定格流量を搬送する際に消費する電力です。
  これは、機器のエネルギー消費量を評価する上で基本的なパラメータです。
- **部分負荷特性 (qef, val, pfcmp)**:
  実際の建物では、ポンプやファンが定格能力で運転される時間は限られています。
  部分負荷特性は、流量が定格流量よりも少ない場合に、
  消費電力がどのように変化するかを示すものです。
  - `qef`: 部分負荷時の効率係数。
  - `Type == "P"` の場合: 太陽電池ポンプの特性を示唆し、
    日射量に応じて流量が変化するモデル（`val`配列に係数を格納）が適用されます。
  - `pfcmp`: ポンプ・ファンの部分負荷特性曲線（`pumpfanlst.efl`から読み込まれる）へのポインターです。
    これにより、より詳細な部分負荷特性をモデル化し、
    実際の運用におけるエネルギー消費量を正確に予測できます。
- **ポンプとファンの区別**: `cattype`によってポンプ（`PUMP_TYPE`）とファン（`FAN_TYPE`）を区別し、
  それぞれの特性に応じたパラメータを設定します。

この関数は、熱搬送システムや空調システムにおけるポンプ・ファンの性能をモデル化し、
熱負荷計算、エネルギー消費量予測、および省エネルギー対策の検討を行うための重要なデータ入力機能を提供します。
*/
func Pumpdata(cattype EqpType, s string, Pumpca *PUMPCA, pfcmp []*PFCMP) int {
	st := strings.IndexByte(s, '=')
	var dt float64
	var id int

	if st == -1 {
		Pumpca.name = s
		Pumpca.Type = ""
		Pumpca.Wo = -999.0
		Pumpca.Go = -999.0
		Pumpca.qef = -999.0
		Pumpca.val = nil

		if cattype == PUMP_TYPE {
			Pumpca.pftype = PUMP_PF
		} else if cattype == FAN_TYPE {
			Pumpca.pftype = FAN_PF
		} else {
			Pumpca.pftype = rune(OFF_SW)
		}
	} else {
		s1, s2 := s[:st], s[st+1:]

		if s1 == "type" {
			Pumpca.Type = s2
			if Pumpca.Type == "P" {
				Pumpca.val = make([]float64, 4)
			}

			for _, pfc := range pfcmp {
				if Pumpca.pftype == pfc.pftype && Pumpca.Type == pfc.Type {
					Pumpca.pfcmp = pfc
					break
				}
			}
		} else {
			dt, _ = strconv.ParseFloat(s[st+1:], 64)
			if s == "qef" {
				Pumpca.qef = dt
			} else {
				if Pumpca.Type != "P" {
					switch s {
					case "Go":
						Pumpca.Go = dt
					case "Wo":
						Pumpca.Wo = dt
					default:
						id = 1
					}
				} else {
					switch s {
					case "a0":
						Pumpca.val[0] = dt
					case "a1":
						Pumpca.val[1] = dt
					case "a2":
						Pumpca.val[2] = dt
					case "Ic":
						Pumpca.val[3] = dt
					default:
						id = 1
					}
				}
			}
		}
	}
	return id
}

/* --------------------------- */

/*
Pumpint (Pump Initialization)

この関数は、ポンプの初期設定を行います。
特に、太陽電池ポンプの場合、その運転に影響を与える太陽電池パネルの方位設定（日射量データへのリンク）を行います。

建築環境工学的な観点:
- **太陽電池ポンプの特性**: 太陽電池ポンプは、太陽光エネルギーを直接利用して運転されるポンプです。
  その流量や消費電力は、太陽電池パネルが受ける日射量に依存します。
  この関数は、`p.Cat.Type == "P"`（太陽電池ポンプ）の場合に、
  対応する太陽電池パネルの日射量データ（`Exs`）へのポインターを設定します。
- **日射量とポンプ流量の関係**: 太陽電池ポンプの流量は、
  日射量が少ない時間帯や曇天時には低下し、
  日射量が多い時間帯には増加します。
  この特性を正確にモデル化することで、
  太陽光エネルギーの利用効率を評価し、
  システム全体のエネルギー消費量を予測できます。
- **再生可能エネルギーの利用**: 太陽電池ポンプは、
  再生可能エネルギーである太陽光を直接利用するため、
  環境負荷の低減やエネルギー自立性の向上に貢献します。
  この関数は、そのようなシステムのモデル化の基礎となります。
- **エラーハンドリング**: `if p.Sol == nil` のように、
  対応する太陽電池パネルが見つからない場合にエラーメッセージを出力することで、
  入力データの不備を早期に発見し、シミュレーションの信頼性を確保します。

この関数は、太陽電池ポンプのような再生可能エネルギーを利用した熱搬送システムのモデル化において、
日射量とポンプ運転の連動を正確に表現するための重要な役割を果たします。
*/
func Pumpint(Pump []*PUMP, Exs []*EXSF) {
	for _, p := range Pump {
		if p.Cat.Type == "P" {
			p.Sol = nil
			for j := 0; j < len(Exs); j++ {
				if p.Cmp.Exsname == Exs[j].Name {
					p.Sol = Exs[j]
					break
				}
			}
			if p.Sol == nil {
				Eprint("Pumpint", p.Cmp.Exsname)
			}
		}
	}
}

/* --------------------------- */

/*
Pumpflow (Pump Flow Rate Calculation)

この関数は、ポンプの流量を計算します。
特に、太陽電池ポンプの場合、太陽電池パネルが受ける日射量に基づいて流量を決定します。
それ以外のポンプについては、運転状態に応じて定格流量を設定します。

建築環境工学的な観点:
- **太陽電池ポンプの流量制御**: 太陽電池ポンプは、
  日射量（`S`）に応じて流量（`p.G`）が変化する特性を持ちます。
  `S > p.Cat.val[3]` の条件は、
  太陽電池が発電を開始するしきい値日射量を示唆し、
  その後の`p.G = p.Cat.val[0] + (p.Cat.val[1]+p.Cat.val[2]*S)*S` の式は、
  日射量と流量の関係をモデル化したものです。
  これにより、日射量に応じた熱媒の搬送能力を正確にシミュレーションできます。
- **定格流量運転**: 太陽電池ポンプ以外のポンプについては、
  `p.Cmp.Control != OFF_SW` の条件で運転中であれば、
  定格流量（`p.Cat.Go`）で運転されるとモデル化されます。
  これは、一般的な定流量ポンプの運転特性を反映しています。
- **エネルギー消費量への影響**: ポンプの流量は、
  熱搬送システム全体の熱供給能力に直接影響します。
  また、流量の変化はポンプの消費電力（`p.E`）にも影響を与え、
  システム全体のエネルギー消費量を左右します。
- **熱負荷への対応**: ポンプの流量が適切に制御されることで、
  建物の熱負荷変動に対して熱媒の供給量を調整し、
  室内温度の安定化や、熱源設備の効率的な運転に貢献します。

この関数は、ポンプの流量をモデル化し、
熱搬送システム全体の熱供給能力、エネルギー消費量、
および熱負荷への応答をシミュレーションするために不可欠な役割を果たします。
*/
func (eqsys *EQSYS) Pumpflow() {
	for i, p := range eqsys.Pump {
		if p.Cat.Type == "P" {
			S := p.Sol.Iw

			if DEBUG {
				fmt.Printf("<Pumpflow> i=%d S=%f Ic=%f a0=%f a1=%e a2=%e\n",
					i, S, p.Cat.val[3], p.Cat.val[0],
					p.Cat.val[1], p.Cat.val[2])
			}

			if S > p.Cat.val[3] {
				p.G = p.Cat.val[0] + (p.Cat.val[1]+p.Cat.val[2]*S)*S
			} else {
				p.G = -999.0
			}

			p.E = 0
		} else {
			if p.Cmp.Control != OFF_SW {
				p.G = p.Cat.Go
				p.E = p.Cat.Wo
			} else {
				p.G = -999.0
				p.E = 0.0
			}

			if DEBUG {
				fmt.Printf("<Pumpflow>  control=%c G=%f E=%f\n",
					p.Cmp.Control, p.G, p.E)
			}
		}
	}
}

/* --------------------------- */

/*
Pumpcfv (Pump/Fan Characteristic Function Value Calculation)

この関数は、ポンプまたはファンの運転特性を評価し、
熱媒の熱容量流量、そして部分負荷時の消費電力に関する係数を計算します。
これは、ポンプやファンが熱搬送システムにおいて熱と動力をどのように供給するかをモデル化するために不可欠です。

建築環境工学的な観点:
- **熱容量流量 (p.CG)**:
  ポンプやファンが搬送する熱媒（空気または水）の熱容量流量は、
  熱搬送システム全体の熱供給能力に直接影響します。
  `Spcheat(Eo1.Fluid) * Eo1.G` のように、熱媒の比熱と質量流量から計算されます。
- **部分負荷特性 (p.PLC)**:
  `PumpFanPLC`関数を用いて計算される`p.PLC`は、
  ポンプやファンが定格流量よりも少ない流量で運転される場合に、
  消費電力がどのように変化するかを示す部分負荷特性係数です。
  実際の建物では、機器が定格能力で運転される時間は限られているため、
  部分負荷特性の考慮はエネルギー消費量予測の精度向上に不可欠です。
- **消費電力の計算**: `Eo1.Co = p.Cat.qef * p.E * p.PLC` のように、
  定格消費電力（`p.E`）に部分負荷特性係数（`p.PLC`）と効率係数（`p.Cat.qef`）を乗じることで、
  実際の消費電力を計算します。
  これにより、ポンプやファンのエネルギー消費量を正確に予測できます。
- **出口温度・湿度の係数設定**: 計算されたポンプやファンの特性に基づいて、
  出口温度（`Eo1`）や出口湿度（`Eo2`、ファンのみ）に関する係数（`Coeffo`, `Co`, `Coeffin`）を設定します。
  これらの係数は、システム全体の熱収支方程式に組み込まれ、
  各流体の出口温度や湿度を予測するために用いられます。

この関数は、ポンプやファンの運転特性を詳細にモデル化し、
熱搬送システム全体の熱供給能力、エネルギー消費量、
および熱負荷への応答をシミュレーションするために不可欠な役割を果たします。
*/
func Pumpcfv(Pump []*PUMP) {
	for _, p := range Pump {
		if p.Cmp.Control != OFF_SW {
			Eo1 := p.Cmp.Elouts[0]
			cG := Spcheat(Eo1.Fluid) * Eo1.G
			p.CG = cG
			Eo1.Coeffo = cG
			p.PLC = PumpFanPLC(Eo1.G/p.G, p)
			Eo1.Co = p.Cat.qef * p.E * p.PLC
			Eo1.Coeffin[0] = -cG

			if p.Cat.pftype == FAN_PF {
				Eo2 := p.Cmp.Elouts[1]
				Eo2.Coeffo = Eo2.G
				Eo2.Co = 0.0
				Eo2.Coeffin[0] = -Eo2.G
			}
		} else {
			p.G = 0.0
			p.E = 0.0
		}
	}
}

/*
PumpFanPLC (Pump/Fan Partial Load Characteristic Curve)

この関数は、ポンプまたはファンの部分負荷特性曲線に基づいて、
部分負荷時の消費電力係数を計算します。
これは、機器が定格流量よりも少ない流量で運転される場合に、
消費電力がどのように変化するかをモデル化するために不可欠です。

建築環境工学的な観点:
- **部分負荷運転の重要性**: 実際の建物では、ポンプやファンが定格能力で運転される時間は限られており、
  多くの場合、部分負荷で運転されます。
  部分負荷時の消費電力を正確に評価することは、
  システム全体のエネルギー消費量予測の精度向上に不可欠です。
- **部分負荷特性曲線**: ポンプやファンの消費電力は、
  流量（または回転数）に対して非線形な特性を示します。
  この関数は、`cat.pfcmp.dblcoeff`に格納された多項式係数を用いて、
  流量比（`XQ`）に対する消費電力係数（`Buff`）を計算します。
  これにより、様々な機器の部分負荷特性を柔軟にモデル化できます。
- **省エネルギー運転の評価**: 部分負荷特性を考慮することで、
  - **変流量制御の効果**: 変流量制御（VAVシステムなど）を導入した場合の省エネルギー効果を定量的に評価できます。
    流量を絞ることで消費電力が大幅に削減される機器は、変流量制御に適しています。
  - **機器選定の最適化**: 部分負荷時の効率が良い機器を選定することで、
    実際の運用におけるエネルギー消費量を削減できます。
- **エラーハンドリング**: `if cat.pfcmp == nil` のように、
  部分負荷特性データが設定されていない場合にエラーメッセージを出力することで、
  入力データの不備を早期に発見し、シミュレーションの信頼性を確保します。

この関数は、ポンプやファンの部分負荷運転時のエネルギー消費量を正確にモデル化し、
熱搬送システムや空調システムの省エネルギー設計、
および運用改善のための意思決定に不可欠な役割を果たします。
*/
func PumpFanPLC(XQ float64, Pump *PUMP) float64 {
	var Buff, dQ float64
	var i int
	cat := Pump.Cat

	dQ = math.Min(1.0, math.Max(XQ, 0.25))

	if cat.pfcmp == nil {
		Err := fmt.Sprintf("<PumpFanPLC>  PFtype=%c  type=%s", cat.pftype, cat.Type)
		Eprint("PUMP oir FAN", string(Err[:]))
		Buff = 0.0
	} else {
		Buff = 0.0

		for i = 0; i < 5; i++ {
			Buff += cat.pfcmp.dblcoeff[i] * math.Pow(dQ, float64(i))
		}
	}
	return Buff
}

/* --------------------------- */

/*
Pumpene (Pump Energy Calculation)

この関数は、ポンプが熱媒に供給する熱量（動力）を計算します。
これは、ポンプのエネルギー消費量と、熱搬送システムにおける熱負荷の処理能力を評価する上で不可欠です。

建築環境工学的な観点:
- **熱媒への動力供給**: ポンプは、熱媒を搬送する際に、その熱媒に動力を供給します。
  この動力は、熱媒の温度上昇として現れ、熱負荷の一部となります。
  `p.Q = p.CG * (Eo.Sysv - p.Tin)` のように、
  熱媒の熱容量流量と入口・出口温度差から計算されます。
  これは、ポンプが熱媒に与えるエネルギー量を示します。
- **エネルギー消費量への影響**: ポンプの消費電力は、
  熱媒に供給される動力とポンプの効率によって決まります。
  この関数で計算される熱量`p.Q`は、
  ポンプの消費電力と密接に関連しており、
  システム全体のエネルギー消費量を評価する上で重要な要素となります。
- **熱負荷への影響**: ポンプが熱媒に供給する動力は、
  熱搬送システムにおける熱負荷の一部として計上されます。
  特に、大規模なシステムや長距離の搬送では、
  ポンプの動力による熱負荷が無視できない場合があります。

この関数は、ポンプのエネルギー消費量と、
熱搬送システムにおける熱負荷の処理能力を定量的に評価し、
システム全体のエネルギー効率を向上させるための設計検討に役立ちます。
*/
func Pumpene(Pump []*PUMP) {
	for _, p := range Pump {
		p.Tin = p.Cmp.Elins[0].Sysvin
		Eo := p.Cmp.Elouts[0]

		if Eo.Control != OFF_SW {
			p.Q = p.CG * (Eo.Sysv - p.Tin)
		} else {
			p.Q = 0.0
		}
	}
}

/* --------------------------- */

func pumpprint(fo io.Writer, id int, Pump []*PUMP) {
	var G float64

	switch id {
	case 0:
		if len(Pump) > 0 {
			fmt.Fprintf(fo, "%s  %d\n", PUMP_TYPE, len(Pump))
		}
		for _, p := range Pump {
			fmt.Fprintf(fo, " %s 1 6\n", p.Name)
		}
	case 1:
		for _, p := range Pump {
			fmt.Fprintf(fo, "%s_c c c %s_Ti t f %s_To t f ", p.Name, p.Name, p.Name)
			fmt.Fprintf(fo, "%s_Q q f  %s_E e f %s_G m f\n", p.Name, p.Name, p.Name)
		}
	default:
		for _, p := range Pump {
			if p.Cmp.Elouts[0].G > 0.0 && p.Cmp.Elouts[0].Control != OFF_SW {
				G = p.Cmp.Elouts[0].G
			} else {
				G = 0.0
			}
			fmt.Fprintf(fo, "%c %4.1f %4.1f %4.0f %4.0f %.5g\n", p.Cmp.Elouts[0].Control,
				p.Tin, p.Cmp.Elouts[0].Sysv, p.Q, p.E*p.PLC, G)
		}
	}
}

/* --------------------------- */

/*
pumpdyint (Pump/Fan Daily Integration Initialization)

この関数は、ポンプまたはファンの日積算値（日ごとの熱量、エネルギー消費量、流量など）をリセットします。
これは、日単位でのポンプ・ファンの運転状況やエネルギー消費量を集計し、
熱搬送システムの性能評価やエネルギー消費量の分析に用いられます。

建築環境工学的な観点:
- **日単位の性能評価**: ポンプやファンの運転状況は、日中の熱負荷変動に応じて大きく変化します。
  日積算値を集計することで、日ごとの熱搬送量、
  ポンプ・ファンの稼働時間、部分負荷運転の割合などを把握できます。
  これにより、特定の日の熱搬送特性を分析したり、
  ポンプ・ファンの運転効率を日単位で評価したりすることが可能になります。
- **運用改善の指標**: 日積算データは、熱搬送システムの運用改善のための重要な指標となります。
  例えば、外気温度や熱負荷などの気象条件とポンプ・ファンのエネルギー消費量の関係を分析したり、
  設定温度や流量などの運用条件がポンプ・ファンの性能に与える影響を評価したりすることで、
  より効率的な運転方法を見つけることができます。
- **データ集計の準備**: この関数は、新しい日の集計を開始する前に、
  前日のデータをクリアする役割を担います。
  `edyint`といった関数は、
  それぞれ熱量、エネルギー、流量などの日積算値をリセットするためのものです。

この関数は、ポンプやファンの運転状況とエネルギー消費量を日単位で詳細に分析し、
熱搬送システムの運用改善や省エネルギー対策の効果評価を行うための基礎的な役割を果たします。
*/
func pumpdyint(Pump []*PUMP) {
	for _, p := range Pump {
		edyint(&p.Qdy)
		edyint(&p.Edy)
		edyint(&p.Gdy)
	}
}

/*
pumpmonint (Pump/Fan Monthly Integration Initialization)

この関数は、ポンプまたはファンの月積算値（月ごとの熱量、エネルギー消費量、流量など）をリセットします。
これは、月単位でのポンプ・ファンの運転状況やエネルギー消費量を集計し、
熱搬送システムの性能評価やエネルギー消費量の分析に用いられます。

建築環境工学的な観点:
- **月単位の性能評価**: ポンプやファンの運転状況は、月単位で変動します。
  月積算値を集計することで、月ごとの熱搬送量、
  ポンプ・ファンの稼働時間、部分負荷運転の割合などを把握できます。
  これにより、特定の月の熱搬送特性を分析したり、
  ポンプ・ファンの運転効率を月単位で評価したりすることが可能になります。
- **運用改善の指標**: 月積算データは、熱搬送システムの運用改善のための重要な指標となります。
  例えば、季節ごとの熱搬送量の傾向を把握したり、
  月ごとの気象条件とポンプ・ファンのエネルギー消費量の関係を分析したりすることで、
  より効率的な運転方法を見つけることができます。
- **データ集計の準備**: この関数は、新しい月の集計を開始する前に、
  前月のデータをクリアする役割を担います。
  `edyint`といった関数は、
  それぞれ熱量、エネルギー、流量などの月積算値をリセットするためのものです。

この関数は、ポンプやファンの運転状況とエネルギー消費量を月単位で詳細に分析し、
熱搬送システムの運用改善や省エネルギー対策の効果評価を行うための基礎的な役割を果たします。
*/
func pumpmonint(Pump []*PUMP) {
	for _, p := range Pump {
		edyint(&p.MQdy)
		edyint(&p.MEdy)
		edyint(&p.MGdy)
	}
}

/*
pumpday (Pump/Fan Daily and Monthly Data Aggregation)

この関数は、ポンプまたはファンの運転データ（熱量、エネルギー消費量、流量など）を、
日単位および月単位で集計します。
これにより、ポンプ・ファンの性能評価やエネルギー消費量の分析が可能になります。

建築環境工学的な観点:
- **日次集計 (edaysum)**:
  日次集計は、ポンプやファンの運転状況を日単位で詳細に把握するために重要です。
  例えば、特定の日の熱負荷変動に対するポンプ・ファンの応答、
  あるいは日中のピーク負荷時の熱搬送量などを分析できます。
  これにより、日ごとの運用改善点を見つけ出すことが可能になります。
- **月次集計 (emonsum)**:
  月次集計は、季節ごとの熱負荷変動や、
  ポンプ・ファンの年間を通じた熱搬送量の傾向を把握するために重要です。
  これにより、年間を通じた省エネルギー対策の効果を評価したり、
  熱搬送量の予測精度を向上させたりすることが可能になります。
- **月・時刻のクロス集計 (emtsum)**:
  `MtEdy`のようなクロス集計は、
  特定の月における時間帯ごとのエネルギー消費量を分析するために用いられます。
  これにより、例えば、冬季の朝の時間帯に暖房負荷が集中する傾向があるか、
  あるいは夏季の夜間に給湯負荷が高いかなどを詳細に把握できます。
  これは、デマンドサイドマネジメントや、
  エネルギー供給計画を最適化する上で非常に有用な情報となります。
- **データ分析の基礎**: この関数で集計されるデータは、
  ポンプやファンの性能評価、エネルギー消費量のベンチマーキング、
  省エネルギー対策の効果検証、そして運用改善のための意思決定の基礎となります。

この関数は、ポンプやファンの運転状況とエネルギー消費量を多角的に分析し、
熱搬送システムの運用改善や省エネルギー対策の効果評価を行うための重要なデータ集計機能を提供します。
*/
func pumpday(Mon, Day, ttmm int, Pump []*PUMP, Nday, SimDayend int) {
	Mo := Mon - 1
	tt := ConvertHour(ttmm)

	for _, p := range Pump {
		// 日集計
		edaysum(ttmm, p.Cmp.Elouts[0].Control, p.Q, &p.Qdy)
		edaysum(ttmm, p.Cmp.Elouts[0].Control, p.E, &p.Edy)
		edaysum(ttmm, p.Cmp.Elouts[0].Control, p.G, &p.Gdy)

		// 月集計
		emonsum(Mon, Day, ttmm, p.Cmp.Elouts[0].Control, p.Q, &p.MQdy, Nday, SimDayend)
		emonsum(Mon, Day, ttmm, p.Cmp.Elouts[0].Control, p.E, &p.MEdy, Nday, SimDayend)
		emonsum(Mon, Day, ttmm, p.Cmp.Elouts[0].Control, p.G, &p.MGdy, Nday, SimDayend)

		// 月・時刻のクロス集計
		emtsum(Mon, Day, ttmm, p.Cmp.Elouts[0].Control, p.E, &p.MtEdy[Mo][tt])
	}
}

func pumpdyprt(fo io.Writer, id int, Pump []*PUMP) {
	switch id {
	case 0:
		if len(Pump) > 0 {
			fmt.Fprintf(fo, "%s  %d\n", PUMP_TYPE, len(Pump))
		}
		for _, p := range Pump {
			fmt.Fprintf(fo, " %s 1 12\n", p.Name)
		}
	case 1:
		for _, p := range Pump {
			fmt.Fprintf(fo, "%s_Hq H d %s_Q Q f %s_tq h d %s_Qm q f\n",
				p.Name, p.Name, p.Name, p.Name)
			fmt.Fprintf(fo, "%s_He H d %s_E E f %s_te h d %s_Em e f\n",
				p.Name, p.Name, p.Name, p.Name)
			fmt.Fprintf(fo, "%s_Hg H d %s_G M f %s_tg h d %s_Gm m f\n\n",
				p.Name, p.Name, p.Name, p.Name)
		}
	default:
		for _, p := range Pump {
			fmt.Fprintf(fo, "%1d %3.1f ", p.Qdy.Hrs, p.Qdy.D)
			fmt.Fprintf(fo, "%1d %2.0f ", p.Qdy.Mxtime, p.Qdy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f ", p.Edy.Hrs, p.Edy.D)
			fmt.Fprintf(fo, "%1d %2.0f ", p.Edy.Mxtime, p.Edy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f ", p.Gdy.Hrs, p.Gdy.D)
			fmt.Fprintf(fo, "%1d %2.0f\n", p.Gdy.Mxtime, p.Gdy.Mx)
		}
	}
}

func pumpmonprt(fo io.Writer, id int, Pump []*PUMP) {
	switch id {
	case 0:
		if len(Pump) > 0 {
			fmt.Fprintf(fo, "%s  %d\n", PUMP_TYPE, len(Pump))
		}
		for _, p := range Pump {
			fmt.Fprintf(fo, " %s 1 12\n", p.Name)
		}
	case 1:
		for _, p := range Pump {
			fmt.Fprintf(fo, "%s_Hq H d %s_Q Q f %s_tq h d %s_Qm q f\n",
				p.Name, p.Name, p.Name, p.Name)
			fmt.Fprintf(fo, "%s_He H d %s_E E f %s_te h d %s_Em e f\n",
				p.Name, p.Name, p.Name, p.Name)
			fmt.Fprintf(fo, "%s_Hg H d %s_G M f %s_tg h d %s_Gm m f\n\n",
				p.Name, p.Name, p.Name, p.Name)
		}
	default:
		for _, p := range Pump {
			fmt.Fprintf(fo, "%1d %3.1f ", p.MQdy.Hrs, p.MQdy.D)
			fmt.Fprintf(fo, "%1d %2.0f ", p.MQdy.Mxtime, p.MQdy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f ", p.MEdy.Hrs, p.MEdy.D)
			fmt.Fprintf(fo, "%1d %2.0f ", p.MEdy.Mxtime, p.MEdy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f ", p.MGdy.Hrs, p.MGdy.D)
			fmt.Fprintf(fo, "%1d %2.0f\n", p.MGdy.Mxtime, p.MGdy.Mx)
		}
	}
}
func pumpmtprt(fo io.Writer, id int, Pump []*PUMP, Mo, tt int) {
	switch id {
	case 0:
		if len(Pump) > 0 {
			fmt.Fprintf(fo, "%s %d\n", PUMP_TYPE, len(Pump))
		}
		for _, p := range Pump {
			fmt.Fprintf(fo, " %s 1 1\n", p.Name)
		}
	case 1:
		for _, p := range Pump {
			fmt.Fprintf(fo, "%s_E E f \n", p.Name)
		}
	default:
		for _, p := range Pump {
			fmt.Fprintf(fo, " %.2f \n", p.MtEdy[Mo-1][tt-1].D*Cff_kWh)
		}
	}
}

func NewPFCMP() *PFCMP {
	Pfcmp := new(PFCMP)
	Pfcmp.pftype = ' '
	Pfcmp.Type = ""
	matinit(Pfcmp.dblcoeff[:], 5)
	return Pfcmp
}

/*
PFcmpdata (Pump/Fan Characteristic Data)

この関数は、ポンプおよびファンの部分負荷特性曲線に関するデータを、
外部ファイル（`pumpfanlst.efl`）から読み込みます。
これらのデータは、機器が定格流量よりも少ない流量で運転される場合に、
消費電力がどのように変化するかをモデル化するために不可欠です。

建築環境工学的な観点:
- **部分負荷特性の重要性**: 実際の建物では、ポンプやファンが定格能力で運転される時間は限られており、
  多くの場合、部分負荷で運転されます。
  部分負荷時の消費電力を正確に評価することは、
  システム全体のエネルギー消費量予測の精度向上に不可欠です。
- **経験式または近似式**: `pumpfanlst.efl`ファイルには、
  ポンプやファンの部分負荷特性を近似するための多項式係数（`dblcoeff`）が格納されています。
  これらの係数は、機器メーカーのデータや実測データに基づいて作成された経験式や近似式を表します。
- **機器選定と省エネルギー**: 部分負荷特性を考慮することで、
  - **変流量制御の効果**: 変流量制御（VAVシステムなど）を導入した場合の省エネルギー効果を定量的に評価できます。
    流量を絞ることで消費電力が大幅に削減される機器は、変流量制御に適しています。
  - **機器選定の最適化**: 部分負荷時の効率が良い機器を選定することで、
    実際の運用におけるエネルギー消費量を削減できます。
- **データ管理の外部化**: 部分負荷特性データを外部ファイルから読み込むことで、
  プログラムの柔軟性が高まります。
  新しい機器の特性を追加したり、既存の機器の特性を変更したりする際に、
  プログラムコードを修正することなく対応できます。

この関数は、ポンプやファンの部分負荷運転時のエネルギー消費量を正確にモデル化し、
熱搬送システムや空調システムの省エネルギー設計、
および運用改善のための意思決定に不可欠な役割を果たします。
*/
func PFcmpdata() []*PFCMP {
	var s string
	var c byte
	var i int

	fl, err := os.Open("pumpfanlst.efl")
	if err != nil {
		Eprint(" file ", "pumpfanlst.efl")
	}
	Pfcmp := make([]*PFCMP, 0)

	for {
		_, err := fmt.Fscanf(fl, "%s", &s)
		if err != nil || s[0] == '*' {
			break
		}

		if s == "!" {
			for {
				_, err = fmt.Fscanf(fl, "%c", &c)
				if err != nil || c == '\n' {
					break
				}
			}
		} else {
			pfcmp := NewPFCMP()

			if s == string(PUMP_TYPE) {
				pfcmp.pftype = PUMP_PF
			} else if s == string(FAN_TYPE) {
				pfcmp.pftype = FAN_PF
			} else {
				Eprint("<pumpfanlst.efl>", s)
			}

			_, err = fmt.Fscanf(fl, "%s", &s)
			if err != nil {
				break
			}

			pfcmp.Type = s

			i = 0
			for {
				_, err = fmt.Fscanf(fl, "%s", &s)
				if err != nil || s[0] == ';' {
					break
				}

				var err error
				pfcmp.dblcoeff[i], err = strconv.ParseFloat(s, 64)
				if err != nil {
					panic(err)
				}
				i++
			}

			Pfcmp = append(Pfcmp, pfcmp)
		}
	}

	fl.Close()

	return Pfcmp
}

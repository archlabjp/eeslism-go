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

/*  boiler.c  */
package eeslism

import (
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
)

/*
Boidata (Boiler Data Input)

この関数は、ボイラーの各種仕様（定格出力、効率、最小出力など）を読み込み、
対応するボイラーの構造体に格納します。
これらのデータは、建物の熱源設備としてのボイラーの性能評価、
熱負荷への対応、およびエネルギー消費量予測に不可欠です。

建築環境工学的な観点:
  - **熱源設備のモデル化**: ボイラーは、建物に熱を供給する主要な熱源設備の一つです。
    その性能を正確にモデル化することは、建物の暖房負荷や給湯負荷への対応能力、
    そしてエネルギー消費量を評価する上で非常に重要です。
  - **定格出力 (Qo)**:
    ボイラーが供給できる最大の熱量を示します。
    建物の最大熱負荷に対して十分な能力があるか、
    あるいは複数台のボイラーを組み合わせる必要があるかを判断する際に用いられます。
  - **効率 (eff)**:
    投入されたエネルギーに対して、どれだけの熱エネルギーを供給できるかを示す指標です。
    効率が高いほど、燃料消費量が少なくなり、省エネルギーに貢献します。
    部分負荷時の効率特性も重要であり、実際の運用におけるエネルギー消費量を左右します。
  - **最小出力 (Qmin)**:
    ボイラーが安定して運転できる最小の熱量を示します。
    最小出力以下の負荷では、ボイラーが頻繁にON/OFFを繰り返す（短サイクル運転）ことになり、
    効率の低下や機器の寿命短縮につながる可能性があります。
    `belowmin`パラメータは、最小出力以下の負荷時にボイラーを停止させるか、
    最小出力で運転を継続するかの制御ロジックを示唆します。
  - **運転モード (p, en)**:
    `p`（部分負荷特性）や`en`（エネルギー源）などのパラメータは、
    ボイラーの運転特性や使用燃料の種類を示唆します。
    これにより、様々な種類のボイラーをモデル化し、
    それぞれのエネルギー消費量や環境負荷を評価できます。
  - **無制限容量 (unlimcap)**:
    `unlimcap`が`'y'`の場合、ボイラーの容量を無制限と見なす設定です。
    これは、初期の設計検討段階で、熱源設備の容量制約を考慮せずに、
    建物の熱負荷特性を把握したい場合などに用いられることがあります。

この関数は、建物の熱源設備としてのボイラーの性能をモデル化し、
熱負荷計算、エネルギー消費量予測、および省エネルギー対策の検討を行うための重要なデータ入力機能を提供します。
*/
func Boidata(s string, boica *BOICA) int {
	var id int

	st := strings.IndexRune(s, '=')
	if st == -1 && s[0] != '-' {
		boica.name = s
		boica.unlimcap = 'n'
		boica.ene = ' '
		boica.plf = ' '
		boica.Qo = nil
		boica.eff = 1.0
		boica.Ph = FNAN
		boica.Qmin = FNAN
		//boica.mode = 'n'
		boica.Qostr = ""
	} else if s == "-U" {
		boica.unlimcap = 'y'
	} else {
		if st >= 0 {
			s1, s2 := s[:st], s[st+1:]
			switch s1 {
			case "p":
				boica.plf = rune(s2[0])
			case "en":
				boica.ene = rune(s2[0])
			case "blwQmin":

				if s2 == "ON" {
					// 負荷が最小出力以下のときに最小出力でONとする
					boica.belowmin = ON_SW
				} else if s2 == "OFF" {
					// 負荷が最小出力以下のときにOFFとする
					boica.belowmin = OFF_SW
				} else {
					id = 1
				}
			case "Qo":
				boica.Qostr = s2
			case "Qmin", "eff", "Ph":
				dt, err := strconv.ParseFloat(s2, 64)
				if err != nil {
					id = 1
				} else {
					switch s1 {
					case "Qmin":
						boica.Qmin = dt
					case "eff":
						boica.eff = dt
					case "Ph":
						boica.Ph = dt
					}
				}
			default:
				id = 1
			}
		}
	}
	return id
}

func (eqcat *EQCAT) Boicaint(Simc *SIMCONTL, Compnt []*COMPNT, Wd *WDAT, Exsf *EXSFS, Schdl *SCHDL) {
	for _, Boica := range eqcat.Boica {
		if idx, err := idsch(Boica.Qostr, Schdl.Sch, ""); err == nil {
			Boica.Qo = &Schdl.Val[idx]
		} else {
			Boica.Qo = envptr(Boica.Qostr, Simc, nil, nil, nil)
		}
	}
}

/* --------------------------- */

/*
Boicfv (Boiler Characteristic Function Value Calculation)

この関数は、ボイラーの運転特性を評価し、
熱供給量や熱媒の流量、温度などの関係を定義する係数を計算します。
これは、ボイラーが建物の熱負荷にどのように応答するかをモデル化するために不可欠です。

建築環境工学的な観点:
  - **熱媒の流量と比熱 (cG)**:
    ボイラーが供給する熱量は、熱媒（水など）の流量と比熱、
    そして入口と出口の温度差によって決まります。
    `cG = Spcheat(Eo1.Fluid) * Eo1.G` は、熱媒の比熱と質量流量を乗じることで、
    熱容量流量（熱媒が単位時間あたりに運ぶことができる熱量）を計算しています。
    これは、熱源設備の能力を評価する上で基本的なパラメータです。
  - **熱供給量 (Do)**:
    `boi.Do = Qocat` は、ボイラーが供給すべき熱量（要求熱量）を設定しています。
    この要求熱量は、建物の熱負荷計算の結果に基づいて決定されます。
  - **運転モード (HCmode)**:
    `boi.HCmode`は、ボイラーが暖房（`HEATING_LOAD`）または冷房（`COOLING_LOAD`）のどちらのモードで運転しているかを示します。
    ボイラーは通常、暖房用途に用いられますが、
    システムによっては熱源として冷房サイクルにも関与する場合があります。
  - **制御ロジック**: `boi.Cmp.Control != OFF_SW` の条件は、
    ボイラーが運転中であることを示し、その後の制御ロジックが実行されます。
  - **出口温度制御**: `Eo1.Sysld == 'y'` の場合、
    ボイラーの出口温度を設定温度に制御するモードを示唆します。
    これは、熱需要に応じて熱媒の温度を調整することで、
    熱供給の安定化や省エネルギー化を図る制御です。
  - **最大/最小能力運転**: `boi.Mode == 'M'` の場合、
    ボイラーが最大能力で運転することを意味し、
    `boi.Cat.Qmin`は最小能力運転を示唆します。
    これらのモードは、熱負荷の変動に対応するためのボイラーの運転戦略をモデル化します。

この関数は、ボイラーの運転特性を詳細にモデル化し、
建物の熱負荷変動に対する熱源設備の応答をシミュレーションするために不可欠な役割を果たします。
*/
func Boicfv(Boi []*BOI) {
	var cG, Qocat, Temp float64

	if len(Boi) != len(Boi) {
		panic("len(Boi) != len(Boi)")
	}

	for _, boi := range Boi {

		Eo1 := boi.Cmp.Elouts[0]

		if boi.Cmp.Control != OFF_SW {
			Temp = math.Abs(*boi.Cat.Qo - (FNAN))
			if math.Abs(Temp) < 1e-3 {
				Qocat = 0.0
			} else {
				Qocat = *boi.Cat.Qo
			}

			if Qocat > 0.0 {
				boi.HCmode = HEATING_LOAD
			} else {
				boi.HCmode = COOLING_LOAD
			}

			boi.Do = Qocat

			if (boi.Do < 0.0 && boi.HCmode == HEATING_LOAD) || (boi.Do > 0.0 && boi.HCmode == COOLING_LOAD) || boi.HCmode == 'n' {
				fmt.Printf("<BOI> name=%s  Qo=%.4g\n", boi.Cmp.Name, boi.Do)
			}

			boi.D1 = 0.0

			cG = Spcheat(Eo1.Fluid) * Eo1.G
			boi.cG = cG
			Eo1.Coeffo = cG

			if Eo1.Control != OFF_SW {
				if Eo1.Sysld == 'y' {
					// 出口を設定温度に制御
					Eo1.Co = 0.0
					Eo1.Coeffin[0] = -cG
				} else {
					if boi.Mode == 'M' {
						// 最大能力
						Eo1.Co = boi.Do
					} else {
						// 最小能力
						Eo1.Co = boi.Cat.Qmin
					}
					Eo1.Coeffin[0] = boi.D1 - cG
				}
			}
		} else {
			// 機器が停止
			Eo1.Co = 0.0
			Eo1.Coeffo = 1.0
			Eo1.Coeffin[0] = -1.0
		}
	}
}

/*
Boiene (Boiler Energy Calculation)

この関数は、ボイラーの供給熱量とエネルギー消費量を計算します。
また、ボイラーの運転状態（ON/OFF、部分負荷運転など）を制御し、
熱負荷に対する応答をモデル化します。

建築環境工学的な観点:
  - **供給熱量 (boi.Q)**:
    ボイラーが実際に建物に供給した熱量です。
    これは、熱媒の熱容量流量（`boi.cG`）と、
    ボイラーの入口温度（`boi.Tin`）と出口温度（`Eo.Sysv`）の差から計算されます。
    この供給熱量は、建物の熱負荷を賄うために必要なエネルギー量を示します。
  - **エネルギー消費量 (boi.E)**:
    ボイラーが熱を供給するために消費したエネルギー量です。
    供給熱量（`boi.Q`）をボイラーの効率（`boi.Cat.eff`）で割ることで計算されます。
    この値は、建物のエネルギー消費量全体に占める熱源設備の割合を評価し、
    省エネルギー対策の効果を定量的に把握する上で重要です。
  - **運転制御ロジック**: ボイラーの運転は、熱負荷の変動や最小出力の制約などに応じて制御されます。
  - **加熱/冷却モードの不一致**: `(Qcheck < 0.0 && boi.HCmode == 'H') || (Qcheck > 0.0 && boi.HCmode == 'C')` の条件は、
    ボイラーが暖房モードなのに冷却負荷が発生している、
    あるいは冷却モードなのに加熱負荷が発生している場合に、ボイラーを停止させるロジックです。
    これは、システムの無駄な運転を防ぎ、エネルギー効率を向上させるために重要です。
  - **最小出力制御**: `Qmin > 0.0 && Qcheck < Qmin` の条件は、
    熱負荷がボイラーの最小出力を下回る場合の制御です。
    `boi.Cat.belowmin == OFF_SW` の場合はボイラーを停止させ、
    それ以外の場合は最小出力で運転を継続する（`boi.Mode = 'm'`）など、
    ボイラーの特性に応じた運転戦略をモデル化します。
  - **過負荷状態のチェック**: `boi.Cat.unlimcap == 'n'` の場合、
    ボイラーの定格出力（`Qocat`）を超えた負荷が発生していないかをチェックします。
    過負荷状態では、ボイラーが要求される熱量を供給できないため、
    室内温度の低下や快適性の悪化につながる可能性があります。

この関数は、ボイラーの動的な運転挙動をモデル化し、
建物の熱負荷変動に対する熱源設備の応答、
およびエネルギー消費量を正確にシミュレーションするために不可欠な役割を果たします。
*/
func Boiene(Boi []*BOI, BOIreset *int) {
	for i, boi := range Boi {
		boi.Tin = boi.Cmp.Elins[0].Sysvin
		Qmin := boi.Cat.Qmin
		if math.Abs(Qmin-(FNAN)) < 1.0e-5 {
			Qmin = 0.0
		}

		Eo := boi.Cmp.Elouts[0]
		reset := 0

		if Eo.Control != OFF_SW {
			boi.Q = boi.cG * (Eo.Sysv - boi.Tin)

			// 次回ループの機器制御判定用の熱量
			Qcheck := boi.Q

			// 加熱設定での冷却、冷却設定での加熱時はボイラを止める
			if (Qcheck < 0.0 && boi.HCmode == 'H') || (Qcheck > 0.0 && boi.HCmode == 'C') {
				boi.Cmp.Control = OFF_SW
				Eo.Control = ON_SW
				Eo.Emonitr.Control = ON_SW
				Eo.Sysld = 'n'

				reset = 1
			} else if Qmin > 0.0 && Qcheck < Qmin {
				// 最小出力以下はOFFにする場合
				if boi.Cat.belowmin == OFF_SW {
					boi.Cmp.Elouts[0].Control = OFF_SW
					boi.Cmp.Control = OFF_SW
					Eo.Control = ON_SW
					Eo.Emonitr.Control = ON_SW
					Eo.Sysld = 'n'
				} else {
					Eo.Control = ON_SW
					Eo.Emonitr.Control = ON_SW
					Eo.Sysld = 'n'
					boi.Mode = 'm'
				}

				reset = 1
			} else if boi.Cat.unlimcap == 'n' {
				// 過負荷状態のチェック
				Qocat := 0.0
				if math.Abs(*boi.Cat.Qo-(FNAN)) < 1.0e-3 {
					Qocat = 0.0
				} else {
					Qocat = *boi.Cat.Qo
				}

				reset0 := maxcapreset(Qcheck, Qocat, boi.HCmode, Eo)

				if reset == 0 {
					reset = reset0
				}
			}

			if reset == 1 {
				Boicfv(Boi[i : i+1])
				(*BOIreset)++
			}

			boi.E = boi.Q / boi.Cat.eff
			boi.Ph = boi.Cat.Ph
		} else {
			boi.Q = 0.0
			boi.E = 0.0
			boi.Ph = 0.0
		}
	}
}

/* --------------------------- */

/*
boildptr (Boiler Load Pointer Setting)

この関数は、ボイラーの負荷計算において、
制御対象となるパラメータ（例: 出口温度）へのポインターを設定します。
これにより、ボイラーの運転を特定の目標値に追従させる制御をモデル化できます。

建築環境工学的な観点:
  - **負荷計算と制御**: 建物の熱負荷は常に変動するため、
    熱源設備であるボイラーは、その負荷変動に応じて熱供給量を調整する必要があります。
    この調整は、熱媒の流量を制御したり、出口温度を制御したりすることで行われます。
  - **制御対象の指定**: `key[1]`が`"Tout"`の場合、
    ボイラーの出口温度（`Boi.Toset`）を制御対象とすることを意味します。
    `vptr.Ptr = &Boi.Toset` は、この出口温度の変数へのポインターを設定し、
    `vptr.Type = VAL_CTYPE` は、そのポインターが制御値であることを示します。
  - **フィードバック制御の基礎**: このポインター設定は、
    ボイラーのフィードバック制御の基礎となります。
    シミュレーションの各時間ステップで、
    現在の出口温度と目標出口温度を比較し、その差に基づいてボイラーの運転を調整します。
    これにより、室内温度の安定化や、熱供給の効率化を図ることができます。

この関数は、ボイラーの制御ロジックをモデル化し、
建物の熱負荷変動に対する熱源設備の応答をシミュレーションするために不可欠な役割を果たします。
*/
func boildptr(load *ControlSWType, key []string, Boi *BOI) (VPTR, error) {
	var err error
	var vptr VPTR

	if strings.Compare(key[1], "Tout") == 0 {
		vptr.Ptr = &Boi.Toset
		vptr.Type = VAL_CTYPE
		Boi.Load = load
	} else {
		err = errors.New("Tout expected")
	}
	return vptr, err
}

/* --------------------------- */

/*
boildschd (Boiler Load Schedule Setting)

この関数は、ボイラーの負荷計算において、
スケジュールに基づいて運転を制御するための設定を行います。
特に、目標出口温度が設定されている場合に、その温度に応じてボイラーの運転をON/OFFしたり、
目標温度を設定したりするロジックを実装します。

建築環境工学的な観点:
  - **スケジュール運転**: 建物の熱負荷は、時間帯や曜日、季節によって変動します。
    ボイラーをスケジュールに基づいて運転することで、
    不要な時間帯の運転を停止したり、熱需要に応じて運転モードを切り替えたりすることができ、
    エネルギー消費量の削減に貢献します。
  - **目標出口温度による制御**: `Boi.Toset`は、ボイラーの目標出口温度を示します。
    `Boi.Toset > TEMPLIMIT` の条件は、
    目標温度が有効な範囲内にある場合にボイラーを運転することを意味します。
    `Eo.Control = LOAD_SW` は、ボイラーが負荷追従運転モードであることを示し、
    `Eo.Sysv = Boi.Toset` は、ボイラーの出口温度を目標温度に設定します。
  - **省エネルギー運転**: 目標出口温度を適切に設定することで、
    過剰な熱供給を防ぎ、エネルギーの無駄を削減できます。
    例えば、外気温度が高い時期には目標出口温度を低く設定することで、
    ボイラーの運転時間を短縮したり、効率を向上させたりすることが可能です。
  - **システム連携**: この関数は、ボイラーの運転制御が、
    熱供給先のシステム（例: 空調機、給湯器）の要求温度と連携して行われることを示唆します。
    これにより、建物全体の熱供給システムを統合的にモデル化し、
    エネルギーマネジメント戦略を評価できます。

この関数は、ボイラーのスケジュール運転と目標温度制御をモデル化し、
建物の熱負荷変動に対する熱源設備の応答、
およびエネルギー消費量をシミュレーションするために不可欠な役割を果たします。
*/
func boildschd(Boi *BOI) {
	Eo := Boi.Cmp.Elouts[0]

	if Boi.Load != nil {
		if Eo.Control != OFF_SW {
			if Boi.Toset > TEMPLIMIT {
				Eo.Control = LOAD_SW
				Eo.Sysv = Boi.Toset
			} else {
				Eo.Control = OFF_SW
			}
		}
	}
}

/* --------------------------- */

func boiprint(fo io.Writer, id int, Boi []*BOI) {
	for _, boi := range Boi {
		switch id {
		case 0:
			if len(Boi) > 0 {
				fmt.Fprintf(fo, "%s  %d\n", BOILER_TYPE, len(Boi))
			}
			fmt.Fprintf(fo, " %s 1 7\n", boi.Name)
		case 1:
			fmt.Fprintf(fo, "%s_c c c %s_G m f %s_Ti t f %s_To t f ", boi.Name, boi.Name, boi.Name, boi.Name)
			fmt.Fprintf(fo, "%s_Q q f  %s_E e f %s_P e f\n", boi.Name, boi.Name, boi.Name)
		default:
			fmt.Fprintf(fo, "%c %.4g %4.2f %4.2f %4.0f %4.0f %2.0f\n",
				boi.Cmp.Elouts[0].Control, boi.Cmp.Elouts[0].G,
				boi.Tin, boi.Cmp.Elouts[0].Sysv, boi.Q, boi.E, boi.Ph)
		}
	}
}

/* --------------------------- */

/*
boidyint (Boiler Daily Integration Initialization)

この関数は、ボイラーの日積算値（日ごとの熱供給量、エネルギー消費量など）をリセットします。
これは、日単位でのエネルギー消費量を集計し、
建物の運用状況や省エネルギー対策の効果を評価するために用いられます。

建築環境工学的な観点:
  - **日単位のエネルギー評価**: 建物のエネルギー消費量は、日単位で変動します。
    日積算値を集計することで、日ごとの熱負荷変動や、
    ボイラーの運転状況（稼働時間、部分負荷運転の割合など）を把握できます。
    これにより、特定の日のエネルギー消費が多かった原因を分析したり、
    省エネルギー対策の効果を日単位で評価したりすることが可能になります。
  - **運用改善の指標**: 日積算データは、建物の運用改善のための重要な指標となります。
    例えば、休日や夜間のエネルギー消費量が過剰でないかを確認したり、
    外気温度や日射量などの気象条件とエネルギー消費量の関係を分析したりすることで、
    より効率的な運転方法を見つけることができます。
  - **データ集計の準備**: この関数は、新しい日の集計を開始する前に、
    前日のデータをクリアする役割を担います。
    `svdyint`, `qdyint`, `edyint`, `phdyint`といった関数は、
    それぞれ温度、熱量、エネルギー、電力などの日積算値をリセットするためのものです。

この関数は、建物のエネルギー消費量を日単位で詳細に分析し、
運用改善や省エネルギー対策の効果評価を行うための基礎的な役割を果たします。
*/
func boidyint(Boi []*BOI) {
	for _, boi := range Boi {
		// 日集計のリセット
		svdyint(&boi.Tidy)
		qdyint(&boi.Qdy)
		edyint(&boi.Edy)
		edyint(&boi.Phdy)
	}
}

/* --------------------------- */

/*
boimonint (Boiler Monthly Integration Initialization)

この関数は、ボイラーの月積算値（月ごとの熱供給量、エネルギー消費量など）をリセットします。
これは、月単位でのエネルギー消費量を集計し、
建物の運用状況や省エネルギー対策の効果を評価するために用いられます。

建築環境工学的な観点:
  - **月単位のエネルギー評価**: 建物のエネルギー消費量は、月単位で変動します。
    月積算値を集計することで、月ごとの熱負荷変動や、
    ボイラーの運転状況（稼働時間、部分負荷運転の割合など）を把握できます。
    これにより、特定の月のエネルギー消費が多かった原因を分析したり、
    省エネルギー対策の効果を月単位で評価したりすることが可能になります。
  - **運用改善の指標**: 月積算データは、建物の運用改善のための重要な指標となります。
    例えば、季節ごとのエネルギー消費量の傾向を把握したり、
    月ごとの気象条件とエネルギー消費量の関係を分析したりすることで、
    より効率的な運転方法を見つけることができます。
  - **データ集計の準備**: この関数は、新しい月の集計を開始する前に、
    前月のデータをクリアする役割を担います。
    `svdyint`, `qdyint`, `edyint`, `phdyint`といった関数は、
    それぞれ温度、熱量、エネルギー、電力などの月積算値をリセットするためのものです。

この関数は、建物のエネルギー消費量を月単位で詳細に分析し、
運用改善や省エネルギー対策の効果評価を行うための基礎的な役割を果たします。
*/
func boimonint(Boi []*BOI) {
	for _, boi := range Boi {
		// 日集計のリセット
		svdyint(&boi.mTidy)
		qdyint(&boi.mQdy)
		edyint(&boi.mEdy)
		edyint(&boi.mPhdy)
	}
}

/*
boiday (Boiler Daily and Monthly Data Aggregation)

この関数は、ボイラーの運転データ（入口温度、供給熱量、エネルギー消費量、電力消費量など）を、
日単位および月単位で集計します。
これにより、ボイラーの性能評価やエネルギー消費量の分析が可能になります。

建築環境工学的な観点:
  - **日次集計 (svdaysum, qdaysum, edaysum)**:
    日次集計は、ボイラーの運転状況を日単位で詳細に把握するために重要です。
    例えば、特定の日の熱負荷変動に対するボイラーの応答、
    あるいは日中のピーク負荷時のエネルギー消費量などを分析できます。
    これにより、日ごとの運用改善点を見つけ出すことが可能になります。
  - **月次集計 (svmonsum, qmonsum, emonsum)**:
    月次集計は、季節ごとの熱負荷変動や、
    ボイラーの年間を通じたエネルギー消費量の傾向を把握するために重要です。
    これにより、年間を通じた省エネルギー対策の効果を評価したり、
    エネルギー消費量の予測精度を向上させたりすることが可能になります。
  - **月・時刻のクロス集計 (emtsum)**:
    `MtEdy`や`MtPhdy`のようなクロス集計は、
    特定の月における時間帯ごとのエネルギー消費量や電力消費量を分析するために用いられます。
    これにより、例えば、冬季の朝の時間帯に暖房負荷が集中する傾向があるか、
    あるいは夏季の夜間に給湯負荷が高いかなどを詳細に把握できます。
    これは、デマンドサイドマネジメントや、
    エネルギー供給計画を最適化する上で非常に有用な情報となります。
  - **データ分析の基礎**: この関数で集計されるデータは、
    ボイラーの性能評価、エネルギー消費量のベンチマーキング、
    省エネルギー対策の効果検証、そして運用改善のための意思決定の基礎となります。

この関数は、建物のエネルギー消費量を多角的に分析し、
運用改善や省エネルギー対策の効果評価を行うための重要なデータ集計機能を提供します。
*/
func boiday(Mon, Day, ttmm int, Boi []*BOI, Nday, SimDayend int) {
	var Mo, tt int

	Mo = Mon - 1
	tt = ConvertHour(ttmm)
	for _, boi := range Boi {
		// 日集計
		svdaysum(int64(ttmm), boi.Cmp.Control, boi.Tin, &boi.Tidy)
		qdaysum(int64(ttmm), boi.Cmp.Control, boi.Q, &boi.Qdy)
		edaysum(ttmm, boi.Cmp.Control, boi.E, &boi.Edy)
		edaysum(ttmm, boi.Cmp.Control, boi.Ph, &boi.Phdy)

		// 月集計
		svmonsum(Mon, Day, ttmm, boi.Cmp.Control, boi.Tin, &boi.mTidy, Nday, SimDayend)
		qmonsum(Mon, Day, ttmm, boi.Cmp.Control, boi.Q, &boi.mQdy, Nday, SimDayend)
		emonsum(Mon, Day, ttmm, boi.Cmp.Control, boi.E, &boi.mEdy, Nday, SimDayend)
		emonsum(Mon, Day, ttmm, boi.Cmp.Control, boi.Ph, &boi.mPhdy, Nday, SimDayend)

		// 月・時刻のクロス集計
		emtsum(Mon, Day, ttmm, boi.Cmp.Control, boi.E, &boi.MtEdy[Mo][tt])
		emtsum(Mon, Day, ttmm, boi.Cmp.Control, boi.E, &boi.MtPhdy[Mo][tt])
	}
}

func boidyprt(fo io.Writer, id int, Boi []*BOI) {
	switch id {
	case 0:
		if len(Boi) > 0 {
			fmt.Fprintf(fo, "%s  %d\n", BOILER_TYPE, len(Boi))
		}
		for _, boi := range Boi {
			fmt.Fprintf(fo, " %s 1 22\n", boi.Name)
		}
	case 1:
		for _, boi := range Boi {
			fmt.Fprintf(fo, "%s_Ht H d %s_T T f ", boi.Name, boi.Name)
			fmt.Fprintf(fo, "%s_ttn h d %s_Tn t f %s_ttm h d %s_Tm t f\n",
				boi.Name, boi.Name, boi.Name, boi.Name)
			fmt.Fprintf(fo, "%s_Hh H d %s_Qh Q f %s_Hc H d %s_Qc Q f\n",
				boi.Name, boi.Name, boi.Name, boi.Name)
			fmt.Fprintf(fo, "%s_th h d %s_qh q f %s_tc h d %s_qc q f\n",
				boi.Name, boi.Name, boi.Name, boi.Name)
			fmt.Fprintf(fo, "%s_He H d %s_E E f %s_te h d %s_Em e f\n",
				boi.Name, boi.Name, boi.Name, boi.Name)
			fmt.Fprintf(fo, "%s_Hp H d %s_P E f %s_tp h d %s_Pm e f\n\n",
				boi.Name, boi.Name, boi.Name, boi.Name)
		}
	default:
		for _, boi := range Boi {
			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ",
				boi.Tidy.Hrs, boi.Tidy.M,
				boi.Tidy.Mntime, boi.Tidy.Mn,
				boi.Tidy.Mxtime, boi.Tidy.Mx)

			fmt.Fprintf(fo, "%1d %3.1f ", boi.Qdy.Hhr, boi.Qdy.H)
			fmt.Fprintf(fo, "%1d %3.1f ", boi.Qdy.Chr, boi.Qdy.C)
			fmt.Fprintf(fo, "%1d %2.0f ", boi.Qdy.Hmxtime, boi.Qdy.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f ", boi.Qdy.Cmxtime, boi.Qdy.Cmx)

			fmt.Fprintf(fo, "%1d %3.1f ", boi.Edy.Hrs, boi.Edy.D)
			fmt.Fprintf(fo, "%1d %2.0f ", boi.Edy.Mxtime, boi.Edy.Mx)

			fmt.Fprintf(fo, "%1d %3.1f ", boi.Phdy.Hrs, boi.Phdy.D)
			fmt.Fprintf(fo, "%1d %2.0f\n", boi.Phdy.Mxtime, boi.Phdy.Mx)
		}
	}
}

func boimonprt(fo io.Writer, id int, Boi []*BOI) {
	switch id {
	case 0:
		if len(Boi) > 0 {
			fmt.Fprintf(fo, "%s  %d\n", BOILER_TYPE, len(Boi))
		}
		for _, boi := range Boi {
			fmt.Fprintf(fo, " %s 1 22\n", boi.Name)
		}
	case 1:
		for _, boi := range Boi {
			fmt.Fprintf(fo, "%s_Ht H d %s_T T f ", boi.Name, boi.Name)
			fmt.Fprintf(fo, "%s_ttn h d %s_Tn t f %s_ttm h d %s_Tm t f\n",
				boi.Name, boi.Name, boi.Name, boi.Name)
			fmt.Fprintf(fo, "%s_Hh H d %s_Qh Q f %s_Hc H d %s_Qc Q f\n",
				boi.Name, boi.Name, boi.Name, boi.Name)
			fmt.Fprintf(fo, "%s_th h d %s_qh q f %s_tc h d %s_qc q f\n",
				boi.Name, boi.Name, boi.Name, boi.Name)
			fmt.Fprintf(fo, "%s_He H d %s_E E f %s_te h d %s_Em e f\n",
				boi.Name, boi.Name, boi.Name, boi.Name)
			fmt.Fprintf(fo, "%s_Hp H d %s_P E f %s_tp h d %s_Pm e f\n\n",
				boi.Name, boi.Name, boi.Name, boi.Name)
		}
	default:
		for _, boi := range Boi {
			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ",
				boi.mTidy.Hrs, boi.mTidy.M,
				boi.mTidy.Mntime, boi.mTidy.Mn,
				boi.mTidy.Mxtime, boi.mTidy.Mx)

			fmt.Fprintf(fo, "%1d %3.1f ", boi.mQdy.Hhr, boi.mQdy.H)
			fmt.Fprintf(fo, "%1d %3.1f ", boi.mQdy.Chr, boi.mQdy.C)
			fmt.Fprintf(fo, "%1d %2.0f ", boi.mQdy.Hmxtime, boi.mQdy.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f ", boi.mQdy.Cmxtime, boi.mQdy.Cmx)

			fmt.Fprintf(fo, "%1d %3.1f ", boi.mEdy.Hrs, boi.mEdy.D)
			fmt.Fprintf(fo, "%1d %2.0f ", boi.mEdy.Mxtime, boi.mEdy.Mx)

			fmt.Fprintf(fo, "%1d %3.1f ", boi.mPhdy.Hrs, boi.mPhdy.D)
			fmt.Fprintf(fo, "%1d %2.0f\n", boi.mPhdy.Mxtime, boi.mPhdy.Mx)
		}
	}
}

func boimtprt(fo io.Writer, id int, Boi []*BOI, Mo int, tt int) {
	switch id {
	case 0:
		if len(Boi) > 0 {
			fmt.Fprintf(fo, "%s %d\n", BOILER_TYPE, len(Boi))
		}
		for _, boi := range Boi {
			fmt.Fprintf(fo, " %s 1 2\n", boi.Name)
		}
	case 1:
		for _, boi := range Boi {
			fmt.Fprintf(fo, "%s_E E f %s_Ph E f \n", boi.Name, boi.Name)
		}
	default:
		for _, boi := range Boi {
			fmt.Fprintf(fo, " %.2f %.2f\n",
				boi.MtEdy[Mo-1][tt-1].D*Cff_kWh, boi.MtPhdy[Mo-1][tt-1].D*Cff_kWh)
		}
	}
}

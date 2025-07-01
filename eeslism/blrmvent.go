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

/* rmvent.c  */
package eeslism

import (
	"fmt"
	"regexp"
	"strings"
)

/* ------------------------------------------------------------------ */

/*
Ventdata (Ventilation Data Input)

この関数は、各室の外気導入量（換気量）とすきま風量を設定するためのデータを読み込み、
対応する室の構造体に格納します。
これらの換気量は、室内の熱負荷計算、空気質評価、およびエネルギー消費量予測において重要な要素です。

建築環境工学的な観点:
- **換気量の設定 (Vent)**:
  換気は、室内の汚染物質（CO2、VOCなど）を排出し、新鮮な外気を導入することで、
  室内空気質を維持するために不可欠です。
  また、換気によって室内の熱が排出されたり、外気の熱が導入されたりするため、
  建物の熱負荷にも大きく影響します。
  `Rm.Gve`は基準となる換気量（質量流量[kg/s]）を示し、
  `Rm.Vesc`は時間帯によって換気量を変化させるためのスケジュール設定値を示唆します。
  これにより、居住者の在室状況や活動レベルに応じた適切な換気計画をモデル化できます。
- **すきま風量の設定 (Inf)**:
  すきま風（infiltration）は、建物の隙間から非意図的に侵入する外気のことです。
  これは、建物の気密性能に依存し、特に冬季の暖房負荷や夏季の冷房負荷に大きな影響を与えます。
  `Rm.Gvi`は基準となるすきま風量、`Rm.Visc`はすきま風量のスケジュール設定値を示唆します。
  すきま風は、計画的な換気とは異なり、制御が難しいため、
  建物の設計段階での気密性の確保が重要となります。
- **室内空気質と熱負荷のバランス**: 換気量を適切に設定することは、
  室内空気質を確保しつつ、過剰な換気による熱損失（または熱取得）を抑え、
  省エネルギーと快適性を両立させる上で重要です。
  この関数で設定される換気データは、これらのバランスを評価するための基礎となります。
- **スケジュール制御**: 換気量やすきま風量をスケジュールで制御できることは、
  実際の建物の運用状況をより忠実に再現するために重要です。
  例えば、夜間や不在時には換気量を絞ることで、エネルギー消費を削減できます。

この関数は、建物の換気計画をモデル化し、
室内空気質、熱負荷、およびエネルギー消費量を評価するための重要なデータ入力機能を提供します。
*/
func Ventdata(fi *EeTokens, Schdl *SCHDL, Room []*ROOM, Simc *SIMCONTL) {
	var Rm *ROOM
	var name1, ss, E string
	var k int

	E = fmt.Sprintf(ERRFMT, "RAICH or VENT")
	for fi.IsEnd() == false {
		line := fi.GetLogicalLine()
		if line[0] == "*" {
			break
		}

		// 室名
		name1 = line[0]

		// 室検索
		i, err := idroom(name1, Room, E+name1)
		if err != nil {
			panic(err)
		}
		Rm = Room[i] //室の参照

		for _, s := range line[1 : len(line)-1] {
			_ss := strings.SplitN(s, "=", 2)
			key := _ss[0]
			valstr := _ss[1]

			switch key {
			case "Vent":
				// 換気量
				// Vent=(基準値[kg/s],換気量設定値名)
				regex := regexp.MustCompile(`\(([^,]+),([^,]+)\)`)
				match := regex.FindStringSubmatch(valstr)
				if len(match) == 3 {
					// 基準値[kg/s]
					Rm.Gve, err = readFloat(match[1])
					if err != nil {
						panic(err)
					}

					// 換気量設定値名
					ss = match[2]
					if k, err := idsch(ss, Schdl.Sch, ""); err == nil {
						Rm.Vesc = &Schdl.Val[k]
					} else {
						Rm.Vesc = envptr(ss, Simc, nil, nil, nil)
					}
				} else {
					fmt.Println("No match found.")
				}

			case "Inf":
				// すきま風
				// Inf=(基準値[kg/s],隙間風量設定値名)
				regex := regexp.MustCompile(`\(([^,]+),([^,]+)\)`)
				match := regex.FindStringSubmatch(valstr)
				if len(match) == 3 {
					// 基準値[kg/s]
					Rm.Gvi, err = readFloat(match[1])
					if err != nil {
						panic(err)
					}

					// 隙間風量設定値名
					ss = match[2]
					if k, err = idsch(ss, Schdl.Sch, ""); err == nil {
						Rm.Visc = &Schdl.Val[k]
					} else {
						Rm.Visc = envptr(ss, Simc, nil, nil, nil)
					}
				} else {
					fmt.Println("No match found.")
				}

			default:
				err := fmt.Sprintf("Room=%s  %s", Rm.Name, key)
				Eprint("<Ventedata>", err)
			}
		}
	}
}

/* ------------------------------------------------------------------ */

/*
Aichschdlr (Air Change Schedule for Room-to-Room Ventilation)

この関数は、室間相互換気量（隣接する室間での空気の移動量）をスケジュールに基づいて設定します。
これは、複数の室からなる建物において、各室の熱負荷や空気質を評価する上で重要な要素です。

建築環境工学的な観点:
- **室間相互換気**: 建物内では、ドアの開閉、内部の圧力差、あるいは意図的な開口部を通じて、
  室間で空気が移動します。
  この室間相互換気は、ある室の熱や汚染物質が隣の室へ移動する経路となり、
  各室の熱負荷や室内空気質に影響を与えます。
- **スケジュール制御**: `val`パラメータは、時間帯や季節に応じて室間相互換気量を変化させるためのスケジュール値を示唆します。
  例えば、昼間はドアを開放して換気を促進し、夜間は閉鎖して熱の移動を抑えるといった運用をモデル化できます。
  `achr.Gvr`は、実際に適用される室間相互換気量（質量流量）を表します。
- **熱負荷と空気質の相互作用**: 室間相互換気は、
  - **熱負荷**: 温度差のある室間で空気が移動することで、熱が輸送されます。
    これにより、ある室の暖房負荷が隣室に影響を与えたり、冷房負荷が軽減されたりする可能性があります。
  - **室内空気質**: 汚染物質（例: CO2、臭気）が排出される室から隣室へ移動することで、
    建物全体の空気質分布に影響を与えます。
    特に、汚染源のある室とそうでない室が隣接する場合に重要です。
- **ゾーン間の熱・物質移動**: 複数のゾーン（室）からなる建物のエネルギーシミュレーションにおいて、
  ゾーン間の熱・物質移動を正確にモデル化することは、
  建物全体のエネルギー消費量や室内環境の評価精度を向上させる上で不可欠です。

この関数は、建物内の空気の流れと熱・物質移動をモデル化し、
各室の熱負荷、室内空気質、および建物全体のエネルギー性能を評価するための重要な機能を提供します。
*/
func Aichschdlr(val []float64, rooms []*ROOM) {
	for i := range rooms {
		room := rooms[i]

		for j := 0; j < room.Nachr; j++ {
			achr := room.achr[j]
			v := val[achr.sch]
			if v > 0.0 {
				achr.Gvr = v
			} else {
				achr.Gvr = 0.0
			}
		}
	}
}

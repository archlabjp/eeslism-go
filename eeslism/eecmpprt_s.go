package eeslism

import (
	"fmt"
	"io"
)

var __Hcmpprint_id int = 0
var __Hstkprint_id int = 0

/*
Hcmpprint (Hourly Component Output)

この関数は、建物のエネルギーシミュレーションにおける各設備機器の
時刻ごとの運転データや熱的挙動を整形して出力します。
これにより、機器の性能、熱負荷への対応、およびエネルギー消費量の詳細な分析が可能になります。

建築環境工学的な観点:
- **機器の運転状況の把握**: シミュレーションの各時間ステップで、
  ボイラー、冷凍機、太陽熱集熱器、冷温水コイル、配管、熱交換器、蓄熱槽、ポンプ、空調負荷、VAV、顕熱蓄熱器、全熱交換器、カロリーメータ、太陽光発電、デシカント空調機、気化冷却器など、
  様々な設備機器の運転データを出力します。
  これにより、各機器が熱負荷にどのように応答しているか、
  また、その際のエネルギー消費量がどうなっているかを詳細に把握できます。
- **出力形式の制御**: `__Hcmpprint_id`によって出力形式を制御し、
  ヘッダー情報（`ttlprint`）やカテゴリ情報（`-cat`）を出力します。
  これにより、出力データを解析ツールなどで利用しやすくなります。
- **熱負荷への対応と省エネルギー**: 各機器の供給熱量や消費エネルギーを出力することで、
  - **熱負荷への対応能力**: 各機器が室の熱負荷にどれだけ貢献しているかを評価できます。
  - **省エネルギー効果**: 高効率機器の導入や、運転条件の最適化による省エネルギー効果を定量的に把握できます。
- **システム統合**: 各機器のデータが統合的に出力されることで、
  建物全体のエネルギーシステムを俯瞰し、
  システム間の相互作用や、エネルギーフローを分析できます。

この関数は、建物のエネルギーシミュレーションにおいて、
各設備機器の運転状況とエネルギー消費量を時刻ごとに詳細に分析し、
省エネルギー対策の効果評価や、最適な設備システム設計を行うための重要なデータ出力機能を提供します。
*/
func Hcmpprint(fo io.Writer, mrk string, Simc *SIMCONTL, mon, day int, time float64, Eqsys *EQSYS, Rdpnl []*RDPNL) {
	var j int

	if __Hcmpprint_id == 0 {
		ttlprint(fo, mrk, Simc)

		for j = 0; j < 2; j++ {
			if j == 0 {
				fmt.Fprintln(fo, "-cat")
			}

			boiprint(fo, __Hcmpprint_id, Eqsys.Boi)
			refaprint(fo, __Hcmpprint_id, Eqsys.Refa)
			collprint(fo, __Hcmpprint_id, Eqsys.Coll)
			hccprint(fo, __Hcmpprint_id, Eqsys.Hcc)
			pipeprint(fo, __Hcmpprint_id, Eqsys.Pipe)
			hexprint(fo, __Hcmpprint_id, Eqsys.Hex)
			stankcmpprt(fo, __Hcmpprint_id, Eqsys.Stank)
			pumpprint(fo, __Hcmpprint_id, Eqsys.Pump)
			hcldprint(fo, __Hcmpprint_id, Eqsys.Hcload)
			vavprint(fo, __Hcmpprint_id, Eqsys.Vav)
			stheatprint(fo, __Hcmpprint_id, Eqsys.Stheat)
			Thexprint(fo, __Hcmpprint_id, Eqsys.Thex)
			Qmeasprint(fo, __Hcmpprint_id, Eqsys.Qmeas)
			PVprint(fo, __Hcmpprint_id, Eqsys.PVcmp)
			Desiprint(fo, __Hcmpprint_id, Eqsys.Desi)
			Evacprint(fo, __Hcmpprint_id, Eqsys.Evac)

			if j == 0 {
				fmt.Fprintln(fo, "*")
				fmt.Fprintln(fo, "#")
			}

			__Hcmpprint_id++
		}
	}

	fmt.Fprintf(fo, "%02d %02d %5.2f\n", mon, day, time)
	boiprint(fo, __Hcmpprint_id, Eqsys.Boi)
	refaprint(fo, __Hcmpprint_id, Eqsys.Refa)
	collprint(fo, __Hcmpprint_id, Eqsys.Coll)
	hccprint(fo, __Hcmpprint_id, Eqsys.Hcc)
	pipeprint(fo, __Hcmpprint_id, Eqsys.Pipe)
	hexprint(fo, __Hcmpprint_id, Eqsys.Hex)
	stankcmpprt(fo, __Hcmpprint_id, Eqsys.Stank)
	pumpprint(fo, __Hcmpprint_id, Eqsys.Pump)
	hcldprint(fo, __Hcmpprint_id, Eqsys.Hcload)
	vavprint(fo, __Hcmpprint_id, Eqsys.Vav)
	stheatprint(fo, __Hcmpprint_id, Eqsys.Stheat)
	Thexprint(fo, __Hcmpprint_id, Eqsys.Thex)
	Qmeasprint(fo, __Hcmpprint_id, Eqsys.Qmeas)
	PVprint(fo, __Hcmpprint_id, Eqsys.PVcmp)
	Desiprint(fo, __Hcmpprint_id, Eqsys.Desi)
	Evacprint(fo, __Hcmpprint_id, Eqsys.Evac)

	if SIMUL_BUILDG {
		panelprint(fo, __Hcmpprint_id, Rdpnl)
	}

}

/*
Hstkprint (Hourly Storage Tank Output)

この関数は、蓄熱槽の時刻ごとの内部水温分布を整形して出力します。
これにより、蓄熱槽の熱的挙動や温度成層の状態を詳細に分析できます。

建築環境工学的な観点:
- **蓄熱槽の温度成層の把握**: 蓄熱槽は、
  温度の異なる水が層状に分かれる「温度成層」を形成することで、
  熱源設備からの高温水と熱利用設備への低温水を効率的に供給できます。
  この関数は、蓄熱槽内部の各層の温度（`Eqsys.Stank`）を出力することで、
  温度成層が適切に維持されているか、
  あるいは崩壊の兆候がないかを詳細に把握できます。
- **熱負荷平準化の評価**: 蓄熱槽の温度分布は、
  熱負荷平準化の効果を評価する上で重要です。
  例えば、夜間に蓄熱された熱が昼間にどのように放熱されているか、
  あるいはピーク負荷時に蓄熱槽がどれだけ熱需要を賄っているかなどを分析できます。
- **出力形式の制御**: `__Hstkprint_id`によって出力形式を制御し、
  ヘッダー情報（`title`）を出力します。
  これにより、出力データを解析ツールなどで利用しやすくなります。

この関数は、蓄熱槽の熱的挙動を詳細に分析し、
熱負荷平準化、エネルギー消費量予測、
および省エネルギー対策の検討を行うための重要なデータ出力機能を提供します。
*/
func Hstkprint(fo io.Writer, title string, mon int, day int, time float64, Eqsys *EQSYS) {
	if __Hstkprint_id == 0 {
		fmt.Fprintf(fo, "%s ;\n", title)
		stankivprt(fo, __Hstkprint_id, Eqsys.Stank)
		__Hstkprint_id++
	}
	if len(Eqsys.Stank) > 0 {
		fmt.Fprintf(fo, "%02d %02d %5.2f  ", mon, day, time)
		stankivprt(fo, __Hstkprint_id, Eqsys.Stank)
	}
	fmt.Fprintln(fo, " ;")
}

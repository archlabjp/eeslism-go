package eeslism

import (
	"fmt"
	"io"
)

var __Pathprint_id int = 0

/*
Pathprint (Path Print)

この関数は、建物のエネルギーシミュレーションにおける熱媒（空気、水など）の
流れる経路（パス）に関する計算結果を整形して出力します。
これにより、熱搬送システムや空調システムの運転状況、
およびエネルギーフローの詳細な分析が可能になります。

建築環境工学的な観点:
- **システム構成の可視化**: 出力されるヘッダー情報には、
  システム経路の名称（`Mpath.Name`）、種別（`Mpath.Type`）、
  流体種別（`Mpath.Fluid`）、末端経路数（`len(Mpath.Plist)`）などが含まれます。
  これにより、シミュレーションモデルのシステム構成を視覚的に把握できます。
- **流量と温度・湿度の追跡**: 各末端経路の流量（`Pli.G`）と制御情報（`Pli.Control`）、
  および経路内の各機器出口における熱媒の温度（`Pelm.Out.Sysv`）や湿度（空気の場合）が出力されます。
  これにより、熱媒がシステム内をどのように流れ、
  その温度や湿度がどのように変化しているかを詳細に追跡できます。
- **運転状態の把握**: 経路や機器の制御情報（`ControlSWType`）を出力することで、
  システムがON/OFF、負荷追従、バッチ運転など、
  どのようなモードで運転しているかを把握できます。
  これは、システムが熱負荷に適切に応答しているか、
  あるいは無駄な運転が発生していないかを評価する上で重要です。
- **エネルギーフローの分析**: 流量と温度・湿度のデータから、
  各経路や機器における熱量（顕熱、潜熱）の移動を計算できます。
  これにより、システム全体のエネルギーフローを分析し、
  エネルギーの無駄を特定し、
  省エネルギー対策の効果を定量的に把握できます。

この関数は、建物のエネルギーシミュレーションにおいて、
熱搬送システムや空調システムの運転状況とエネルギーフローを詳細に分析し、
省エネルギー対策の効果評価や、最適な設備システム設計を行うための重要なデータ出力機能を提供します。
*/
func Pathprint(fo io.Writer, title string, mon int, day int, time float64, _Mpath []*MPATH) {
	// ** ヘッダーの出力 **
	if __Pathprint_id == 0 {
		__Pathprint_id++
		fmt.Fprintf(fo, "%s ;\n", title)
		fmt.Fprintf(fo, "%d\n", len(_Mpath))

		// 経路のループ
		for _, Mpath := range _Mpath {
			// 経路名、種別、流体種別、末端経路数を出力
			fmt.Fprintf(fo, "%s %c %c %d\n", Mpath.Name, Mpath.Type, Mpath.Fluid, len(Mpath.Plist))

			if Mpath.Plist[0].Pelm[0].Ci == ELIO_IN {
				fmt.Fprint(fo, " >")
			}

			// 末端経路のループ
			for _, Pli := range Mpath.Plist {
				// 末端経路の名前を出力
				if Pli.Name != "" {
					fmt.Fprintf(fo, "%s", Pli.Name)
				} else {
					fmt.Fprint(fo, "?")
				}

				// 末端経路の種別と機器要素数を出力
				// T: ,C: , B:
				// b:合流, c:分岐, i:流入境界, o: 流出境界
				fmt.Fprintf(fo, " %c %d\n", rune(Pli.Type), len(Pli.Pelm))

				// 専用要素が '>' の場合はコンポーネント名を出力
				if Pli.Pelm[0].Ci == ELIO_IN {
					fmt.Fprintf(fo, " %s", Pli.Pelm[0].Cmp.Name)
				}

				// 各要素のコンポーネント名を出力
				for _, Pelm := range Pli.Pelm {
					fmt.Fprintf(fo, " %s", Pelm.Cmp.Name)
				}
				fmt.Fprint(fo, "\n")
			}
		}
	}

	// ** 状態表示 **

	fmt.Fprintf(fo, "%02d %02d %5.2f\n", mon, day, time)

	// 経路のループ
	for _, Mpath := range _Mpath {
		// 連立方程式の答えの出力形式
		fm := get_format_by_fluidtype(Mpath.Fluid)

		// 経路の制御情報Control
		fmt.Fprintf(fo, "[%c]", get_control_sw_rune(Mpath.Control))

		// 末端経路のループ
		for _, Pli := range Mpath.Plist {

			// 末端経路の流量Gと制御情報Control
			fmt.Fprintf(fo, " %5.3g %c: ", Pli.G, get_control_sw_rune(Pli.Control))

			if Pli.Pelm[0].Ci == ELIO_IN {
				fmt.Fprintf(fo, fm, Pli.Pelm[0].In.Sysvin)
			}

			// 末端経路内の機器のループ
			for _, Pelm := range Pli.Pelm {
				if Pelm.Out != nil {
					// 連立方程式の答え
					fmt.Fprintf(fo, fm, Pelm.Out.Sysv)

					// 機器出口の制御情報
					fmt.Fprintf(fo, " %c ", get_control_sw_rune(Pelm.Out.Control))
				}
			}
			fmt.Fprint(fo, "\n")
		}
	}
	fmt.Fprintf(fo, " ;\n")
}

/*
get_format_by_fluidtype (Get Format String by Fluid Type)

この関数は、熱媒の種類（`fluid_type`）に応じて、
出力する数値のフォーマット文字列を返します。

建築環境工学的な観点:
- **出力の可読性向上**: 熱媒の種類によって、
  適切な数値の表示形式が異なります。
  例えば、空気の絶対湿度は非常に小さな値になるため、
  小数点以下の桁数を多く表示する必要があります。
  この関数は、空気の絶対湿度（`AIRx_FLD`）の場合には`"%6.4f"`（小数点以下4桁）を、
  それ以外の場合には`"%4.1f"`（小数点以下1桁）を返すことで、
  出力の可読性を向上させます。

この関数は、建物のエネルギーシミュレーションにおいて、
出力データの品質を向上させるための補助的な役割を果たします。
*/
func get_format_by_fluidtype(fluid_type FliudType) string {
	if fluid_type == AIRx_FLD {
		return "%6.4f"
	}
	return "%4.1f"
}

/*
get_control_sw_rune (Get Control Switch Rune)

この関数は、制御情報（`control`）を印字用の文字に変換して返します。

建築環境工学的な観点:
- **運転状態の可視化**: シミュレーション結果の出力において、
  機器や経路の運転状態を簡潔な文字で表示することで、
  出力の可読性を向上させ、
  運転状況を視覚的に分かりやすくします。
  例えば、`OFF_SW`（停止中）の場合は`'x'`、`ON_SW`（動作中）の場合は`'-'`を返します。
- **デバッグと検証**: 機器の運転状態が意図通りに制御されているかを確認する際に、
  この出力は非常に役立ちます。

この関数は、建物のエネルギーシミュレーションにおいて、
出力データの品質を向上させるための補助的な役割を果たします。
*/
func get_control_sw_rune(control ControlSWType) rune {
	var c = rune(control)
	if c == 0 {
		c = '?'
	}
	return c
}

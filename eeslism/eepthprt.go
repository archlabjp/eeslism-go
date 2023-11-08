package eeslism

import (
	"fmt"
	"io"
)

var __Pathprint_id int

// 経路に沿ったシステム要素の出力
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

func get_format_by_fluidtype(fluid_type FliudType) string {
	if fluid_type == AIRx_FLD {
		return "%6.4f"
	}
	return "%4.1f"
}

// 制御情報を印字用文字列に変換する
// 不適切な状態であった場合は '?' を返す
func get_control_sw_rune(control ControlSWType) rune {
	var c = rune(control)
	if c == 0 {
		c = '?'
	}
	return c
}

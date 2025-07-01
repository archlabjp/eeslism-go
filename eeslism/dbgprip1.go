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
	"fmt"
	"io"
	"os"
)

/*
dprdayweek (Display Day of Week for Debugging)

この関数は、年間を通じた各日の曜日情報をデバッグ目的で出力します。
これは、スケジュール設定や気象データとの関連性を確認し、
シミュレーションの挙動を理解するために用いられます。

建築環境工学的な観点:
- **曜日情報の確認**: 建物のエネルギー消費量や室内環境は、
  曜日によって大きく変動します。
  例えば、平日はオフィスビルが稼働し、週末は住宅のエネルギー消費が増加する傾向があります。
  この関数は、`daywk`配列に格納された各日の曜日情報（0:日曜日, 1:月曜日など）を出力することで、
  スケジュール設定が意図通りに適用されているかを確認できます。
- **デバッグと検証**: シミュレーションが期待通りの結果を出さない場合、
  日付と曜日の対応関係の誤りが原因であることがあります。
  このデバッグ出力は、曜日情報が正しいかを確認し、
  問題の特定に役立ちます。
- **モデルの理解**: 曜日情報は、
  建物の運用パターンやエネルギー消費量を理解する上で重要です。
  この出力は、曜日ごとの挙動を視覚的に確認し、
  モデルの理解を深めるのに役立ちます。

この関数は、建物の熱的挙動を詳細に分析し、
シミュレーションの信頼性を確保するための重要なデバッグ機能を提供します。
*/
func dprdayweek(daywk []int) {
	const dmax = 366

	fmt.Print("---  Day of week -----\n   ")
	for d := 0; d < 8; d++ {
		fmt.Printf("  %s=%d  ", DAYweek[d], d)
	}
	fmt.Println()

	k := 1
	for d := 1; d < dmax; d++ {
		if FNNday(k, 1) == d {
			fmt.Printf("\n%2d - ", k)
			k++
		}
		fmt.Printf("%2d", daywk[d])
	}
	fmt.Println()
}

/* ----------------------------------------------------------------- */

/*
dprschtable (Display Schedule Table for Debugging)

この関数は、スケジュール設定（季節、曜日、1日の設定値、1日の切替スケジュール）を
デバッグ目的で出力します。
これは、スケジュール設定が正しいかを確認し、
シミュレーションの挙動を理解するために用いられます。

建築環境工学的な観点:
- **スケジュール設定の確認**: 建物の空調システムや換気システムは、
  時間帯や曜日、季節に応じて運転モードや設定値が変化するスケジュールに基づいて制御されます。
  この関数は、`Schdl.Seasn`（季節設定）、`Schdl.Wkdy`（曜日設定）、
  `Schdl.Dsch`（1日の設定値スケジュール）、`Schdl.Dscw`（1日の切替スケジュール）を出力することで、
  スケジュール設定が意図通りに設定されているかを確認できます。
- **デバッグと検証**: シミュレーションが期待通りの結果を出さない場合、
  スケジュール設定の誤りが原因であることがあります。
  このデバッグ出力は、スケジュール設定が正しいかを確認し、
  問題の特定に役立ちます。
- **モデルの理解**: スケジュール設定は、
  建物の運用パターンやエネルギー消費量を理解する上で重要です。
  この出力は、スケジュール設定の挙動を視覚的に確認し、
  モデルの理解を深めるのに役立ちます。

この関数は、建物の熱的挙動を詳細に分析し、
シミュレーションの信頼性を確保するための重要なデバッグ機能を提供します。
*/
func (Schdl *SCHDL) dprschtable() {

	Ssn, Wkd, Dh, Dw := Schdl.Seasn, Schdl.Wkdy, Schdl.Dsch, Schdl.Dscw

	Ns := len(Ssn)
	Nw := len(Wkd)
	Nsc := len(Dh)
	Nsw := len(Dw)

	if DEBUG {
		fmt.Printf("\n*** dprschtable  ***\n")
		fmt.Printf("\n=== Schtable end  is=%d  iw=%d  sc=%d  sw=%d\n", Ns, Nw, Nsc, Nsw)

		// 季節設定の出力
		for _, Seasn := range Ssn {
			fmt.Printf("\n- %s", Seasn.name)

			for js := range Seasn.sday {
				sday := Seasn.sday[js]
				eday := Seasn.eday[js]
				fmt.Printf("  %4d-%4d", sday, eday)
			}
		}

		// 曜日設定の出力
		for _, Wkdy := range Wkd {
			fmt.Printf("\n- %s", Wkdy.name)

			for _, wday := range Wkdy.wday {
				if wday {
					fmt.Printf("   1")
				} else {
					fmt.Printf("   0")
				}
			}
		}

		// 1日の設定値スケジュールの出力
		for sc, Dsch := range Dh {
			fmt.Printf("\n-VL   %10s (%2d) ", Dsch.name, sc)

			for jsc := range Dsch.stime {
				stime := Dsch.stime[jsc]
				val := Dsch.val[jsc]
				etime := Dsch.etime[jsc]
				fmt.Printf("  %4d-(%6.2f)-%4d", stime, val, etime)
			}
		}

		// 1日の切替スケジュールの出力
		for sw, Dscw := range Dw {
			fmt.Printf("\n-SW   %10s (%2d) ", Dscw.name, sw)

			for jsw := range Dscw.stime {
				stime := Dscw.stime[jsw]
				mode := Dscw.mode[jsw]
				etime := Dscw.etime[jsw]
				fmt.Printf("  %4d-( %c)-%4d", stime, mode, etime)
			}
		}
	}

	if Ferr != nil {
		fmt.Fprintf(Ferr, "\n*** dprschtable  ***\n")
		fmt.Fprintf(Ferr, "\n=== Schtable end  is=%d  iw=%d  sc=%d  sw=%d\n", Ns, Nw, Nsc, Nsw)

		// 季節設定の出力
		for _, Seasn := range Ssn {
			fmt.Fprintf(Ferr, "\n\t%s", Seasn.name)

			for js := range Seasn.sday {
				sday := Seasn.sday[js]
				eday := Seasn.eday[js]
				fmt.Fprintf(Ferr, "\t%d-%d", sday, eday)
			}
		}

		// 曜日の出力
		for j := range DAYweek {
			fmt.Fprintf(Ferr, "\t%s", DAYweek[j])
		}

		// 曜日設定の出力
		for _, Wkdy := range Wkd {
			fmt.Fprintf(Ferr, "\n%s", Wkdy.name)

			for _, wday := range Wkdy.wday {
				if wday {
					fmt.Fprintf(Ferr, "\t1")
				} else {
					fmt.Fprintf(Ferr, "\t0")
				}
			}
		}

		// 1日の設定値スケジュールの出力
		for sc, Dsch := range Dh {
			fmt.Fprintf(Ferr, "\nVL\t%s\t[%d]", Dsch.name, sc)

			for jsc := range Dsch.stime {
				stime := Dsch.stime[jsc]
				val := Dsch.val[jsc]
				etime := Dsch.etime[jsc]
				fmt.Fprintf(Ferr, "\t%d-(%.2g)-%d", stime, val, etime)
			}
		}

		// 1日の切替スケジュールの出力
		for sw, Dscw := range Dw {
			fmt.Fprintf(Ferr, "\nSW\t%s\t[%d]", Dscw.name, sw)

			for jsw := range Dscw.stime {
				stime := Dscw.stime[jsw]
				mode := Dscw.mode[jsw]
				etime := Dscw.etime[jsw]
				fmt.Fprintf(Ferr, "\t%d-(%c)-%d", stime, mode, etime)
			}
		}

		fmt.Fprintf(Ferr, "\n\n")
	}
}

/* ----------------------------------------------------------------- */

/*
dprschdata (Display Schedule Data for Debugging)

この関数は、年間を通じたスケジュールデータ（設定値スケジュール、切替スケジュール）を
デバッグ目的で出力します。
これは、スケジュール設定が正しいかを確認し、
シミュレーションの挙動を理解するために用いられます。

建築環境工学的な観点:
- **年間を通じたスケジュール設定の確認**: 建物の空調システムや換気システムは、
  年間を通じて運転モードや設定値が変化するスケジュールに基づいて制御されます。
  この関数は、`Sh`（設定値スケジュール）と`Sw`（切替スケジュール）を出力することで、
  年間を通じたスケジュール設定が意図通りに設定されているかを確認できます。
- **デバッグと検証**: シミュレーションが期待通りの結果を出さない場合、
  スケジュール設定の誤りが原因であることがあります。
  このデバッグ出力は、スケジュール設定が正しいかを確認し、
  問題の特定に役立ちます。
- **モデルの理解**: スケジュール設定は、
  建物の運用パターンやエネルギー消費量を理解する上で重要です。
  この出力は、年間を通じたスケジュール設定の挙動を視覚的に確認し、
  モデルの理解を深めるのに役立ちます。

この関数は、建物の熱的挙動を詳細に分析し、
シミュレーションの信頼性を確保するための重要なデバッグ機能を提供します。
*/
func dprschdata(Sh []SCH, Sw []SCH) {
	const dmax = 366

	Nsc := len(Sh)
	Nsw := len(Sw)

	if DEBUG {
		fmt.Printf("\n*** dprschdata  ***\n")
		fmt.Printf("\n== len(Sch)=%d   len(Scw)=%d\n", Nsc, Nsw)

		for i := 0; i < Nsc; i++ {
			Sch := &Sh[i]
			fmt.Printf("\nSCH= %s (%2d) ", Sch.name, i)

			k := 1
			for d := 1; d < dmax; d++ {
				day := Sch.day[d]
				if FNNday(k, 1) == d {
					fmt.Printf("\n%2d - ", k)
					k++
				}
				fmt.Printf("%2d", day)
			}
		}

		for i := 0; i < Nsw; i++ {
			Scw := &Sw[i]
			fmt.Printf("\nSCW= %s (%2d) ", Scw.name, i)
			k := 1
			for d := 1; d < dmax; d++ {
				day := Scw.day[d]
				if FNNday(k, 1) == d {
					fmt.Printf("\n%2d - ", k)
					k++
				}
				fmt.Printf("%2d", day)
			}
		}
		fmt.Printf("\n")
	}

	if Ferr != nil {
		fmt.Fprintf(Ferr, "\n*** dprschdata  ***\n")
		fmt.Fprintf(Ferr, "\n== len(Sch)=%d   len(Scw)=%d\n", Nsc, Nsw)

		for i := 0; i < Nsc; i++ {
			Sch := &Sh[i]
			fmt.Fprintf(Ferr, "\nSCH=%s\t[%d]\t", Sch.name, i)

			k := 1
			for d := 1; d < dmax; d++ {
				day := Sch.day[d]
				if FNNday(k, 1) == d {
					fmt.Fprintf(Ferr, "\n%2d - ", k)
					k++
				}
				fmt.Fprintf(Ferr, "%2d", day)
			}
		}

		for i := 0; i < Nsw; i++ {
			Scw := &Sw[i]
			fmt.Fprintf(Ferr, "\nSCW= %s (%2d) ", Scw.name, i)
			k := 1
			for d := 1; d < dmax; d++ {
				day := Scw.day[d]
				if FNNday(k, 1) == d {
					fmt.Fprintf(Ferr, "\n%2d - ", k)
					k++
				}
				fmt.Fprintf(Ferr, "%2d", day)
			}
		}
		fmt.Fprintf(Ferr, "\n")
	}
}

/* ----------------------------------------------------------------- */

/*
dprachv (Display Room-to-Room Air Change for Debugging)

この関数は、室間相互換気（隣接する室間での空気の移動）に関するデータを
デバッグ目的で出力します。
これは、室間相互換気モデルの入力が正しいかを確認し、
シミュレーションの挙動を理解するために用いられます。

建築環境工学的な観点:
- **室間相互換気モデルの確認**: 建物内では、ドアの開閉、内部の圧力差、
  あるいは意図的な開口部を通じて、室間で空気が移動します。
  この室間相互換気は、ある室の熱や汚染物質が隣の室へ移動する経路となり、
  各室の熱負荷や室内空気質に影響を与えます。
  この関数は、各室の室間相互換気に関する情報（接続先の室、スケジュールなど）を出力することで、
  室間相互換気モデルの入力が意図通りに設定されているかを確認できます。
- **デバッグと検証**: シミュレーションが期待通りの結果を出さない場合、
  室間相互換気モデルの誤りが原因であることがあります。
  このデバッグ出力は、室間相互換気モデルの入力が正しいかを確認し、
  問題の特定に役立ちます。
- **モデルの理解**: 室間相互換気は、
  室内の空気質や熱負荷、エネルギー消費量を理解する上で重要です。
  この出力は、室間相互換気モデルの挙動を視覚的に確認し、
  モデルの理解を深めるのに役立ちます。

この関数は、建物の熱的挙動を詳細に分析し、
シミュレーションの信頼性を確保するための重要なデバッグ機能を提供します。
*/
func dprachv(Room []ROOM) {

	f := func(s io.Writer) {
		fmt.Fprintln(Ferr, "\n*** dprachv***")

		for i := range Room {
			Rm := Room[i]
			fmt.Fprintf(Ferr, "to rm: %-10s   from rms(sch):", Rm.Name)

			for j := 0; j < Rm.Nachr; j++ {
				A := Rm.achr[j]
				fmt.Fprintf(Ferr, "  %-10s (%3d)", Room[A.rm].Name, A.sch)
			}
			fmt.Fprintln(Ferr)
		}
	}

	if DEBUG {
		f(os.Stdout)
	}

	if Ferr != nil {
		f(Ferr)
	}
}

/* ----------------------------------------------------------------- */

/*
dprexsf (Display External Surface Data for Debugging)

この関数は、外部日射面（`EXSF`）に関するデータ（名称、タイプ、方位角、傾斜角、
地盤反射率、標高、拡散日射補正係数など）をデバッグ目的で出力します。
これは、外部日射面モデルの入力が正しいかを確認し、
シミュレーションの挙動を理解するために用いられます。

建築環境工学的な観点:
- **外部日射面モデルの確認**: 建物の日射熱取得量や太陽光発電システムの発電量予測は、
  正確な外部日射面データに大きく依存します。
  この関数は、各外部日射面が受ける日射量や、
  その方位角、傾斜角などを出力することで、
  外部日射面モデルの入力が意図通りに設定されているかを確認できます。
- **デバッグと検証**: シミュレーションが期待通りの結果を出さない場合、
  外部日射面モデルの誤りが原因であることがあります。
  このデバッグ出力は、外部日射面モデルの入力が正しいかを確認し、
  問題の特定に役立ちます。
- **モデルの理解**: 外部日射面は、
  建物の熱負荷やエネルギー消費量を理解する上で重要です。
  この出力は、外部日射面モデルの挙動を視覚的に確認し、
  モデルの理解を深めるのに役立ちます。

この関数は、建物の熱的挙動を詳細に分析し、
シミュレーションの信頼性を確保するための重要なデバッグ機能を提供します。
*/
func (exsfs *EXSFS) dprexsf() {
	if exsfs.Exs == nil {
		return
	}

	if DEBUG {
		fmt.Println("\n*** dprexsf ***")
		for i, Exs := range exsfs.Exs {
			fmt.Printf("%2d  %-11s  typ=%c Wa=%6.2f Wb=%5.2f Rg=%4.2f  z=%5.2f edf=%6.2e\n",
				i, Exs.Name, Exs.Typ, Exs.Wa, Exs.Wb, Exs.Rg, Exs.Z, Exs.Erdff)
		}
	}

	if Ferr != nil {
		fmt.Fprintln(Ferr, "\n*** dprexsf ***")
		fmt.Fprintln(Ferr, "\tNo.\tName\ttyp\tWa\tWb\tRg\tz\tedf")

		for i, Exs := range exsfs.Exs {
			fmt.Fprintf(Ferr, "\t%d\t%s\t%c\t%.4g\t%.4g\t%.2g\t%.2g\t%.2g\n",
				i, Exs.Name, Exs.Typ, Exs.Wa, Exs.Wb, Exs.Rg, Exs.Z, Exs.Erdff)
		}
	}
}

/* ----------------------------------------------------------------- */

/*
dprwwdata (Display Wall and Window Data for Debugging)

この関数は、壁体と窓の仕様データ（熱抵抗、放射率、日射吸収率、層構成など）を
デバッグ目的で出力します。
これは、壁体と窓の熱的モデルの入力が正しいかを確認し、
シミュレーションの挙動を理解するために用いられます。

建築環境工学的な観点:
- **壁体と窓の熱的モデルの確認**: 建物の熱負荷は、
  壁体と窓の熱的特性に大きく依存します。
  この関数は、各壁体の熱抵抗（`Wall.Rwall`）、
  放射率（`Wall.Ei`, `Wall.Eo`）、日射吸収率（`Wall.as`）、
  および層構成（`Wall.welm`）を出力することで、
  壁体と窓の熱的モデルの入力が意図通りに設定されているかを確認できます。
- **デバッグと検証**: シミュレーションが期待通りの結果を出さない場合、
  壁体と窓の熱的モデルの誤りが原因であることがあります。
  このデバッグ出力は、壁体と窓の熱的モデルの入力が正しいかを確認し、
  問題の特定に役立ちます。
- **モデルの理解**: 壁体と窓は、
  建物の熱負荷やエネルギー消費量を理解する上で重要です。
  この出力は、壁体と窓の熱的モデルの挙動を視覚的に確認し、
  モデルの理解を深めるのに役立ちます。

この関数は、建物の熱的挙動を詳細に分析し、
シミュレーションの信頼性を確保するための重要なデバッグ機能を提供します。
*/
func (Rmvls *RMVLS) dprwwdata() {
	if DEBUG {
		fmt.Printf("\n*** dprwwdata ***\nWALLdata\n")

		for i, Wall := range Rmvls.Wall {
			fmt.Printf("\nWall i=%d %s R=%5.3f IP=%d Ei=%4.2f Eo=%4.2f as=%4.2f\n", i, get_string_or_null(Wall.name), Wall.Rwall, Wall.Ip, Wall.Ei, Wall.Eo, Wall.as)

			for j := 0; j < Wall.N; j++ {
				w := &Wall.welm[j]
				fmt.Printf("   %2d  %-10s %5.3f %2d\n", j, w.Code, w.L, w.ND)
			}
		}

		fmt.Printf("\nWINDOWdata\n")

		for _, Window := range Rmvls.Window {
			fmt.Printf("windows  %s\n", Window.Name)
			fmt.Printf(" R=%f t=%f B=%f  Ei=%f Eo=%f\n", Window.Rwall, Window.tgtn, Window.Bn, Window.Ei, Window.Eo)
		}
	}

	if Ferr != nil {
		fmt.Fprintf(Ferr, "\n*** dprwwdata ***\nWALLdata\n")

		for i, Wall := range Rmvls.Wall {
			fmt.Fprintf(Ferr, "\nWall[%d]\t%s\tR=%.3g\tIP=%d\tEi=%.2g\tEo=%.2g\tas=%.2g\n", i, Wall.name, Wall.Rwall, Wall.Ip, Wall.Ei, Wall.Eo, Wall.as)

			fmt.Fprintf(Ferr, "\tNo.\tcode\tL\tND\n")

			for j := 0; j < Wall.N; j++ {
				w := &Wall.welm[j]
				fmt.Fprintf(Ferr, "\t%d\t%s\t%.3g\t%d\n", j, w.Code, w.L, w.ND)
			}
		}

		fmt.Fprintf(Ferr, "\nWINDOWdata\n")

		for i, Window := range Rmvls.Window {
			fmt.Fprintf(Ferr, "windows[%d]\t%s\n", i, Window.Name)
			fmt.Fprintf(Ferr, "\tR=%.3g\tt=%.2g\tB=%.2g\tEi=%.2g\tEo=%.2g\n", Window.Rwall,
				Window.tgtn, Window.Bn, Window.Ei, Window.Eo)
		}
	}
}

/* ----------------------------------------------------------------- */

/*
dprroomdata (Display Room Data for Debugging)

この関数は、各室の仕様データ（名称、熱容量、床面積、表面積、換気量、内部発熱など）を
デバッグ目的で出力します。
これは、室の熱的モデルの入力が正しいかを確認し、
シミュレーションの挙動を理解するために用いられます。

建築環境工学的な観点:
- **室の熱的モデルの確認**: 室の熱負荷は、
  室の熱容量、換気量、内部発熱、および各表面からの熱伝達に大きく依存します。
  この関数は、各室の熱容量（`Room.MRM`）、床面積（`Room.FArea`）、
  表面積（`Room.Area`）、換気量（`Room.Gve`, `Room.Gvi`）、
  内部発熱（`Room.Light`, `Room.Nhm`）、
  および各表面の熱的特性（`Sdd.ble`, `Sdd.typ`, `Sdd.A`など）を出力することで、
  室の熱的モデルの入力が意図通りに設定されているかを確認できます。
- **デバッグと検証**: シミュレーションが期待通りの結果を出さない場合、
  室の熱的モデルの誤りが原因であることがあります。
  このデバッグ出力は、室の熱的モデルの入力が正しいかを確認し、
  問題の特定に役立ちます。
- **モデルの理解**: 室の熱的特性は、
  建物の熱負荷やエネルギー消費量を理解する上で重要です。
  この出力は、室の熱的モデルの挙動を視覚的に確認し、
  モデルの理解を深めるのに役立ちます。

この関数は、建物の熱的挙動を詳細に分析し、
シミュレーションの信頼性を確保するための重要なデバッグ機能を提供します。
*/
func (Rmvls *RMVLS) dprroomdata() {
	if DEBUG {
		fmt.Printf("\n*** dprroomdata ***\n")

		for i, Room := range Rmvls.Room {
			fmt.Printf("\n==room=(%d)    %s   N=%d  Ntr=%d Nrp=%d  V=%8.1f   MRM=%10.4e\n",
				i, Room.Name, Room.N, Room.Ntr, Room.Nrp, Room.VRM, Room.MRM)
			fmt.Printf("   Floor area=%6.2f   Total surface area=%6.2f\n", Room.FArea, Room.Area)

			fmt.Printf("   Gve=%f    Gvi=%f\n",
				Room.Gve, Room.Gvi)
			fmt.Printf("   Light=%f  Ltyp=%c  ", Room.Light, Room.Ltyp)
			fmt.Printf("  Nhm=%f\n",
				Room.Nhm)
			fmt.Printf("  Apsc=%f  Apsr=%f   ",
				Room.Apsc, Room.Apsr)
			fmt.Printf("  Apl=%f \n", Room.Apl)

			for j := 0; j < Room.N; j++ {
				Sdd := Rmvls.Sd[Room.Brs+j]
				fmt.Printf(" %2d  ble=%c typ=%c name=%8s exs=%2d nxrm=%2d nxn=%2d ",
					Room.Brs+j, Sdd.ble, Sdd.typ, get_string_or_null(Sdd.Name), Sdd.exs, Sdd.nxrm, Sdd.nxn)
				fmt.Printf("wd=%2d Nfn=%2d A=%5.1f mwside=%c mwtype=%c Ei=%.2f Eo=%.2f\n",
					Sdd.wd, Sdd.Nfn, Sdd.A, Sdd.mwside, Sdd.mwtype, Sdd.Ei, Sdd.Eo)
			}
		}
	}

	if Ferr != nil {
		fmt.Fprintf(Ferr, "\n*** dprroomdata ***\n")

		for i, Room := range Rmvls.Room {
			fmt.Fprintf(Ferr, "\n==room=(%d)\t%s\tN=%d\tNtr=%d\tNrp=%d\tV=%.3g\tMRM=%.2g\n",
				i, Room.Name, Room.N, Room.Ntr, Room.Nrp, Room.VRM, Room.MRM)
			fmt.Fprintf(Ferr, "\tFloor_area=%.3g\tTotal_surface_area=%.2g\n", Room.FArea, Room.Area)

			fmt.Fprintf(Ferr, "\tGve=%.2g\tGvi=%.2g\n", Room.Gve, Room.Gvi)
			fmt.Fprintf(Ferr, "\tLight=%.2g\tLtyp=%c", Room.Light, Room.Ltyp)
			fmt.Fprintf(Ferr, "\tNhm=%.2g\n", Room.Nhm)
			fmt.Fprintf(Ferr, "\tApsc=%.2g\tApsr=%.2g", Room.Apsc, Room.Apsr)
			fmt.Fprintf(Ferr, "\tApl=%.2g\n", Room.Apl)

			fmt.Fprintf(Ferr, "\tNo.\tble\ttyp\tname\texs\tnxrmd\tnxn\t")
			fmt.Fprintf(Ferr, "wd\tNfn\tA\tmwside\tmwtype\tEi\tEo\n")

			for j := 0; j < Room.N; j++ {
				Sdd := Rmvls.Sd[Room.Brs+j]
				fmt.Fprintf(Ferr, "\t%d\t%c\t%c\t%s\t%d\t%d\t%d\t", Room.Brs+j, Sdd.ble, Sdd.typ, Sdd.Name, Sdd.exs, Sdd.nxrm, Sdd.nxn)
				fmt.Fprintf(Ferr, "%d\t%d\t%.3g\t%c\t%c\t%.2f\t%.2f\n", Sdd.wd, Sdd.Nfn, Sdd.A, Sdd.mwside, Sdd.mwtype, Sdd.Ei, Sdd.Eo)
			}
		}
	}
}

/* ----------------------------------------------------------------- */

/*
dprballoc (Display Wall Allocation Data for Debugging)

この関数は、壁体（`MWALL`）と室表面（`RMSRF`）の関連付けに関するデータを
デバッグ目的で出力します。
これは、壁体と室表面のモデルの入力が正しいかを確認し、
シミュレーションの挙動を理解するために用いられます。

建築環境工学的な観点:
- **壁体と室表面の関連付けの確認**: 建物の熱負荷計算では、
  壁体と室表面の関連付けが正確に行われていることが重要です。
  この関数は、各壁体（`Mw`）がどの室表面（`Sd`）に割り当てられているか、
  およびその熱的特性（`Mw.M`：壁体層数、`Mw.sd.A`：表面積など）を出力することで、
  壁体と室表面のモデルの入力が意図通りに設定されているかを確認できます。
- **デバッグと検証**: シミュレーションが期待通りの結果を出さない場合、
  壁体と室表面の関連付けの誤りが原因であることがあります。
  このデバッグ出力は、壁体と室表面のモデルの入力が正しいかを確認し、
  問題の特定に役立ちます。
- **モデルの理解**: 壁体と室表面の関連付けは、
  建物の熱負荷やエネルギー消費量を理解する上で重要です。
  この出力は、壁体と室表面のモデルの挙動を視覚的に確認し、
  モデルの理解を深めるのに役立ちます。

この関数は、建物の熱的挙動を詳細に分析し、
シミュレーションの信頼性を確保するための重要なデバッグ機能を提供します。
*/
func (Rmvls *RMVLS) dprballoc() {
	if DEBUG {
		fmt.Println("\n*** dprballoc ***")

		for mw, Mw := range Rmvls.Mw {
			id := Rmvls.Sd[Mw.ns].wd
			fmt.Printf(" %2d n=%2d  rm=%2d  nxrm=%2d wd=%2d wall=%s M=%2d A=%.2f\n",
				mw, Mw.ns, Mw.rm, Mw.nxrm, id, get_string_or_null(Mw.wall.name), Mw.M, Mw.sd.A)
		}
	}
	if Ferr != nil {
		fmt.Fprintln(Ferr, "\n*** dprballoc ***")
		fmt.Fprintln(Ferr, "\tNo.\tn\trm\tnxrm\twd\twall\tM\tA")

		for mw, Mw := range Rmvls.Mw {
			id := Rmvls.Sd[Mw.ns].wd
			fmt.Fprintf(Ferr, "\t%d\t%d\t%d\t%d\t%d\t%s\t%d\t%.2g\n",
				mw, Mw.ns, Mw.rm, Mw.nxrm, id, Mw.wall.name, Mw.M, Mw.sd.A)
		}
	}
}

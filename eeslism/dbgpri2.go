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

/*  dbgpri2.c   */

package eeslism

import "fmt"

/*
xprroom (Export Room Coefficients for Debugging)

この関数は、室の熱収支計算に関する主要な係数（熱容量、熱伝達、熱負荷など）を
デバッグ目的で出力します。
これは、室の熱収支モデルの計算結果が正しいかを確認し、
シミュレーションの挙動を理解するために用いられます。

建築環境工学的な観点:
- **室の熱収支モデルの確認**: 室の熱収支は、
  壁体からの熱伝達、日射熱取得、内部発熱、換気など、
  様々な要因によって決まります。
  この関数は、以下の主要な係数を出力することで、
  室の熱収支モデルの計算結果が意図通りに設定されているかを確認できます。
  - `Room.MRM`, `Room.GRM`: 室の顕熱・潜熱熱容量。
  - `Room.RMt`, `Room.RMx`: 室温・絶対湿度に関する係数。
  - `Room.ARN`, `Room.RMP`: 隣室からの熱伝達、放射パネルからの熱供給に関する係数。
  - `Room.RMC`, `Room.RMXC`: 室の顕熱・潜熱熱負荷に関する定数項。
- **デバッグと検証**: シミュレーションが期待通りの結果を出さない場合、
  室の熱収支モデルの誤りが原因であることがあります。
  このデバッグ出力は、室の熱収支モデルの計算結果が正しいかを確認し、
  問題の特定に役立ちます。
- **モデルの理解**: これらの係数は、
  室の熱的挙動を理解する上で重要です。
  この出力は、室の熱収支モデルの挙動を視覚的に確認し、
  モデルの理解を深めるのに役立ちます。

この関数は、建物の熱的挙動を詳細に分析し、
シミュレーションの信頼性を確保するための重要なデバッグ機能を提供します。
*/
func xprroom(R []*ROOM) {
	var j int
	var ARN []float64
	var RMP []float64
	var Room *ROOM

	Room = R[0]
	if DEBUG {
		fmt.Println("--- xprroom")
		for i := range R {
			Room = R[i]
			fmt.Printf(" Room:  name=%s  MRM=%f  GRM=%f\n", Room.Name, Room.MRM, Room.GRM)
			fmt.Printf("     RMt=%f", Room.RMt)

			ARN = Room.ARN
			for j = 0; j < Room.Ntr; j++ {
				fmt.Printf(" ARN=%f", ARN[j])
			}

			RMP = Room.RMP
			for j = 0; j < Room.Nrp; j++ {
				fmt.Printf(" RMP=%f", RMP[j])
			}

			fmt.Printf(" RMC=%f\n", Room.RMC)
			fmt.Printf("     RMx=%f          RMXC=%f\n", Room.RMx, Room.RMXC)
		}
	}

	Room = R[0]
	if Ferr != nil {
		fmt.Fprintln(Ferr, "--- xprroom")
		for i := range R {
			Room = R[i]
			fmt.Fprintf(Ferr, "Room:\tname=%s\tMRM=%.4g\tGRM=%.4g\n", Room.Name, Room.MRM, Room.GRM)
			fmt.Fprintf(Ferr, "\tRMt=%.4g\n", Room.RMt)

			ARN = Room.ARN
			for j = 0; j < Room.Ntr; j++ {
				fmt.Fprintf(Ferr, "\tARN[%d]=%.4g", j, ARN[j])
			}
			fmt.Fprintln(Ferr)

			RMP = Room.RMP
			for j = 0; j < Room.Nrp; j++ {
				fmt.Fprintf(Ferr, "\tRMP[%d]=%.4g", j, RMP[j])
			}
			fmt.Fprintln(Ferr)

			fmt.Fprintf(Ferr, "\tRMC=%.4g\n", Room.RMC)
			fmt.Fprintf(Ferr, "\tRMx=%.2g\t\tRMXC=%.2g\n", Room.RMx, Room.RMXC)
		}
	}
}

/* ----------------------------------------- */

/*
xprschval (Export Schedule Values for Debugging)

この関数は、スケジュールデータ（`val`）とスイッチ状態（`isw`）を
デバッグ目的で出力します。
これは、スケジュール設定が正しいかを確認し、
シミュレーションの挙動を理解するために用いられます。

建築環境工学的な観点:
- **スケジュール設定の確認**: 建物の空調システムや換気システムは、
  時間帯や曜日、季節に応じて運転モードや設定値が変化するスケジュールに基づいて制御されます。
  この関数は、`val`（スケジュール値）と`isw`（スイッチ状態）を出力することで、
  スケジュール設定が意図通りに設定されているかを確認できます。
- **デバッグと検証**: シミュレーションが期待通りの結果を出さない場合、
  スケジュール設定の誤りが原因であることがあります。
  このデバッグ出力は、スケジュール設定が正しいかを確認し、
  問題の特定に役立ちます。
- **モデルの理解**: スケジュール設定は、
  建物のエネルギー消費量や室内環境を理解する上で重要です。
  この出力は、スケジュール設定の挙動を視覚的に確認し、
  モデルの理解を深めるのに役立ちます。

この関数は、建物の熱的挙動を詳細に分析し、
シミュレーションの信頼性を確保するための重要なデバッグ機能を提供します。
*/
func xprschval(val []float64, isw []ControlSWType) {
	fmt.Println("--- xprschval")

	for j := range val {
		fmt.Printf("--- val=(%d) %f\n", j, val[j])
	}

	for j := range isw {
		fmt.Printf("--- isw=(%d) %c\n", j, isw[j])
	}
}

/* --------------------------------------------- */

/*
xprqin (Export Room Internal Heat Gains for Debugging)

この関数は、室内の内部発熱（人体、照明、機器など）に関するデータを
デバッグ目的で出力します。
これは、内部発熱モデルの入力が正しいかを確認し、
シミュレーションの挙動を理解するために用いられます。

建築環境工学的な観点:
- **内部発熱の確認**: 室内の内部発熱は、
  建物の熱負荷に大きな影響を与えます。
  この関数は、各室の内部発熱（`r.Hc`, `r.Hr`, `r.HL`, `r.Lc`, `r.Lr`, `r.Ac`, `r.Ar`, `r.AL`）を出力することで、
  内部発熱モデルの入力が意図通りに設定されているかを確認できます。
  - `Hc`: 人体からの顕熱。
  - `Hr`: 人体からの放射熱。
  - `HL`: 人体からの潜熱。
  - `Lc`: 照明からの顕熱。
  - `Lr`: 照明からの放射熱。
  - `Ac`: 機器からの顕熱。
  - `Ar`: 機器からの放射熱。
  - `AL`: 機器からの潜熱。
- **デバッグと検証**: シミュレーションが期待通りの結果を出さない場合、
  内部発熱モデルの誤りが原因であることがあります。
  このデバッグ出力は、内部発熱モデルの入力が正しいかを確認し、
  問題の特定に役立ちます。
- **モデルの理解**: 内部発熱は、
  建物の熱負荷やエネルギー消費量を理解する上で重要です。
  この出力は、内部発熱モデルの挙動を視覚的に確認し、
  モデルの理解を深めるのに役立ちます。

この関数は、建物の熱的挙動を詳細に分析し、
シミュレーションの信頼性を確保するための重要なデバッグ機能を提供します。
*/
func xprqin(Room []*ROOM) {
	fmt.Printf("--- xprqin  Nroom=%d\n", len(Room))

	for i := range Room {
		r := Room[i]
		fmt.Printf("  [%d] Hc=%f Hr=%f HL=%f Lc=%f Lr=%f Ac=%f Ar=%f AL=%f\n",
			i, r.Hc, r.Hr, r.HL, r.Lc, r.Lr, r.Ac, r.Ar, r.AL)
	}
}

/* --------------------------------------------- */

/*
xprvent (Export Room Ventilation Data for Debugging)

この関数は、各室の換気量（外気導入量、室間相互換気量）に関するデータを
デバッグ目的で出力します。
これは、換気モデルの入力が正しいかを確認し、
シミュレーションの挙動を理解するために用いられます。

建築環境工学的な観点:
- **換気量データの確認**: 換気は、室内の空気質を維持し、
  熱負荷に影響を与える重要な要素です。
  この関数は、各室の外気導入量（`Room.Gvent`）や、
  室間相互換気量（`A.Gvr`）を出力することで、
  換気モデルの入力が意図通りに設定されているかを確認できます。
- **デバッグと検証**: シミュレーションが期待通りの結果を出さない場合、
  換気モデルの誤りが原因であることがあります。
  このデバッグ出力は、換気モデルの入力が正しいかを確認し、
  問題の特定に役立ちます。
- **モデルの理解**: 換気量は、
  室内の空気質や熱負荷、エネルギー消費量を理解する上で重要です。
  この出力は、換気モデルの挙動を視覚的に確認し、
  モデルの理解を深めるのに役立ちます。

この関数は、建物の熱的挙動を詳細に分析し、
シミュレーションの信頼性を確保するための重要なデバッグ機能を提供します。
*/
func xprvent(R []*ROOM) {
	var j int
	var A *ACHIR
	var Room *ROOM

	if DEBUG {
		fmt.Println("--- xprvent")

		for i := range R {
			Room = R[i]
			fmt.Printf("  [%d] %-10s  Gvent=%f  -- Gvr:", i, Room.Name, Room.Gvent)

			for j = 0; j < Room.Nachr; j++ {
				A = Room.achr[j]
				fmt.Printf(" <%d>=%f", A.rm, A.Gvr)
			}
			fmt.Println()
		}
	}

	if Ferr != nil {
		fmt.Fprintln(Ferr, "\n\n--- xprvent")

		for i := range R {
			Room = R[i]
			fmt.Fprintf(Ferr, "\t[%d]\t%s\tGvent=%.3g\n\t\t", i, Room.Name, Room.Gvent)

			for j = 0; j < Room.Nachr; j++ {
				A = Room.achr[j]
				fmt.Fprintf(Ferr, "\t<%d>=%.2g", A.rm, A.Gvr)
			}
			fmt.Fprintln(Ferr)
		}
	}
}

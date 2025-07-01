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
)

/*
xprtwallinit (Export Wall Temperature Initialization for Debugging)

この関数は、壁体内部温度の初期値や、シミュレーション開始時の壁体温度分布をデバッグ目的で出力します。
これは、壁体モデルの初期設定が正しいかを確認し、
シミュレーションの挙動を理解するために用いられます。

建築環境工学的な観点:
- **壁体温度分布の初期状態**: 壁体は熱容量を持つため、
  その初期温度分布はシミュレーションの初期段階の熱的挙動に影響を与えます。
  この関数は、各壁体の各層の初期温度（`M[j].Told[m]`）を出力することで、
  モデルの初期状態が意図通りに設定されているかを確認できます。
- **デバッグと検証**: シミュレーションが期待通りの結果を出さない場合、
  初期条件の誤りが原因であることがあります。
  このデバッグ出力は、壁体モデルの初期設定が正しいかを確認し、
  問題の特定に役立ちます。
- **モデルの理解**: 壁体内部の温度分布は、
  壁体の熱貫流特性や蓄熱効果を理解する上で重要です。
  この出力は、壁体モデルの挙動を視覚的に確認し、
  モデルの理解を深めるのに役立ちます。

この関数は、建物の熱的挙動を詳細に分析し、
シミュレーションの信頼性を確保するための重要なデバッグ機能を提供します。
*/
func xprtwallinit(Nmwall int, M []*MWALL) {
	Max := 0
	for j := 0; j < Nmwall; j++ {
		if M[j].M > Max {
			Max = M[j].M
		}
	}

	if DEBUG {
		fmt.Println("--- xprtwallinit")
		for j := 0; j < Nmwall; j++ {
			fmt.Printf("Told  j=%2d", j)
			for m := 0; m < M[j].M; m++ {
				fmt.Printf("  %2d%5.1f", m, M[j].Told[m])
			}
			fmt.Println()
		}
	}

	if Ferr != nil {
		fmt.Fprintln(Ferr, "--- xprtwallinit")
		fmt.Fprint(Ferr, "\tNo.")
		for j := 0; j < Max; j++ {
			fmt.Fprintf(Ferr, "\tT[%d]", j)
		}
		fmt.Fprintln(Ferr)
		for j := 0; j < Nmwall; j++ {
			fmt.Fprintf(Ferr, "\t%d", j)
			for m := 0; m < M[j].M; m++ {
				fmt.Fprintf(Ferr, "\t%.3g", M[j].Told[m])
			}
			fmt.Fprintln(Ferr)
		}
	}
}

/* -------------------------------------------- */

/*
xprsolrd (Export Solar Radiation Data for Debugging)

この関数は、外部日射面（`EXSF`）に関する日射量データ（直達日射、拡散日射、全天日射、夜間放射など）を
デバッグ目的で出力します。
これは、日射量モデルの入力が正しいかを確認し、
シミュレーションの挙動を理解するために用いられます。

建築環境工学的な観点:
- **日射量データの確認**: 建物の日射熱取得量や太陽光発電システムの発電量予測は、
  正確な日射量データに大きく依存します。
  この関数は、各外部日射面が受ける日射量（`Exs.Idre`, `Exs.Idf`, `Exs.Iw`）や、
  夜間放射量（`Exs.Rn`）、清澄度指数（`Exs.Cinc`）などを出力することで、
  日射量モデルの入力が意図通りに設定されているかを確認できます。
- **デバッグと検証**: シミュレーションが期待通りの結果を出さない場合、
  日射量データの誤りが原因であることがあります。
  このデバッグ出力は、日射量モデルの入力が正しいかを確認し、
  問題の特定に役立ちます。
- **モデルの理解**: 日射量データは、
  建物の熱負荷やエネルギー消費量を理解する上で重要です。
  この出力は、日射量モデルの挙動を視覚的に確認し、
  モデルの理解を深めるのに役立ちます。

この関数は、建物の日射環境を詳細に分析し、
シミュレーションの信頼性を確保するための重要なデバッグ機能を提供します。
*/
func xprsolrd(E []*EXSF) {
	if DEBUG {
		fmt.Println("--- xprsolrd")
		for i, Exs := range E {
			fmt.Printf("EXSF[%2d]=%s  Id=%5.0f  Idif=%5.0f  Iw=%5.0f RN=%5.0f cinc=%5.3f\n",
				i, Exs.Name, Exs.Idre, Exs.Idf, Exs.Iw, Exs.Rn, Exs.Cinc)
		}
	}

	if Ferr != nil {
		fmt.Fprintln(Ferr, "--- xprsolrd")
		fmt.Fprintln(Ferr, "\tNo.\tName\tId\tIdif\tIw\tRN\tcinc")
		for i, Exs := range E {
			fmt.Fprintf(Ferr, "\t%d\t%s\t%.0f\t%.0f\t%.0f\t%.0f\t%.3f\n",
				i, Exs.Name, Exs.Idre, Exs.Idf, Exs.Iw, Exs.Rn, Exs.Cinc)
		}
	}
}

/* ---------------------------------------------------------- */

/*
xpralph (Export Surface Heat Transfer Coefficients for Debugging)

この関数は、室内の各表面における熱伝達率（室内側表面熱伝達率、放射熱伝達率、対流熱伝達率）を
デバッグ目的で出力します。
これは、表面熱伝達モデルの計算結果が正しいかを確認し、
シミュレーションの挙動を理解するために用いられます。

建築環境工学的な観点:
- **表面熱伝達の確認**: 室内表面からの熱伝達は、
  室内の温熱環境や快適性に直接影響します。
  この関数は、各表面の熱伝達率（`Sd.alo`, `Sd.alir`, `Sd.alic`, `Sd.ali`）を出力することで、
  表面熱伝達モデルの計算結果が意図通りに設定されているかを確認できます。
  - `alo`: 室外側表面熱伝達率。
  - `alir`: 室内側放射熱伝達率。
  - `alic`: 室内側対流熱伝達率。
  - `ali`: 室内側総合熱伝達率（放射＋対流）。
- **形態係数の確認**: `Room.alr`は、
  室内表面間の放射形態係数を示す行列です。
  この出力は、形態係数の設定が正しいかを確認し、
  放射熱交換モデルの挙動を理解するのに役立ちます。
- **デバッグと検証**: シミュレーションが期待通りの結果を出さない場合、
  表面熱伝達モデルの誤りが原因であることがあります。
  このデバッグ出力は、表面熱伝達モデルの計算結果が正しいかを確認し、
  問題の特定に役立ちます。

この関数は、建物の熱的挙動を詳細に分析し、
シミュレーションの信頼性を確保するための重要なデバッグ機能を提供します。
*/
func xpralph(_Room []*ROOM, S []*RMSRF) {
	fmt.Println("--- xpralph")

	for i := range _Room {
		Room := _Room[i]
		N := Room.N
		brs := Room.Brs

		fmt.Println(" alr(i,j)")
		Matfprint("  %5.1f", N, Room.alr)
		fmt.Println(" alph")

		for n := brs; n < brs+N; n++ {
			Sd := S[n]
			fmt.Printf("  %3d  alo=%5.1f  alir=%5.1f alic=%5.1f  ali=%5.1f\n",
				n, Sd.alo, Sd.alir, Sd.alic, Sd.ali)
		}
	}
}

/* ---------------------------------------------------------- */

/*
xprxas (Export Room Surface Coefficients for Debugging)

この関数は、室内の各表面における熱収支計算に関する係数（熱応答係数、放射熱交換係数など）を
デバッグ目的で出力します。
これは、室の熱収支モデルの計算結果が正しいかを確認し、
シミュレーションの挙動を理解するために用いられます。

建築環境工学的な観点:
- **室の熱収支モデルの確認**: 室の熱収支は、
  壁体からの熱伝達、日射熱取得、内部発熱、換気など、
  様々な要因によって決まります。
  この関数は、各表面の熱応答係数（`Sd.FI`, `Sd.FO`, `Sd.FP`）、
  放射熱交換係数（`Sd.WSR`, `Sd.WSRN`, `Sd.WSPL`, `Sd.WSC`）、
  および熱伝達係数（`Sd.K`, `Sd.alo`, `Sd.CF`）などを出力することで、
  室の熱収支モデルの計算結果が意図通りに設定されているかを確認できます。
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
func xprxas(R []*ROOM, S []*RMSRF) {
	if DEBUG {
		fmt.Printf("--- xprxas\n")

		for _, Room := range R {
			N := Room.N
			brs := Room.Brs

			fmt.Printf(" XA(i,j)\n")
			Matprint("%7.4f", N, Room.XA)

			for n := brs; n < brs+N; n++ {
				Sd := S[n]
				fmt.Printf("%2d  K=%f  alo=%f  FI=%f FO=%f FP=%f  CF=%f\n",
					n, Sd.K, Sd.alo, Sd.FI, Sd.FO, Sd.FP, Sd.CF)
				fmt.Printf("            WSR=%f", Sd.WSR)

				for j := 0; j < Room.Ntr; j++ {
					fmt.Printf(" WSRN=%f", Sd.WSRN[j])
				}

				for j := 0; j < Room.Nrp; j++ {
					fmt.Printf(" WSPL=%f", Sd.WSPL[j])
				}

				fmt.Printf("   WSC=%f\n", Sd.WSC)
			}
		}
	}

	if Ferr != nil {
		fmt.Fprintf(Ferr, "--- xprxas\n")

		for _, Room := range R {
			N := Room.N
			brs := Room.Brs

			fmt.Fprintf(Ferr, "Room=%s\tXA(i,j)\n", Room.Name)
			Matfiprint(Ferr, "\t%.1g", N, Room.XA)

			for n := brs; n < brs+N; n++ {
				Sd := S[n]
				fmt.Fprintf(Ferr, "\n\n\t%d\tK=%.2g\talo=%.2g\tFI=%.2g\tFO=%.2g\tFP=%.2g\tCF=%.2g\t",
					n, Sd.K, Sd.alo, Sd.FI, Sd.FO, Sd.FP, Sd.CF)
				fmt.Fprintf(Ferr, "\t\tWSR=%.3g\n\t", Sd.WSR)

				for j := 0; j < Room.Ntr; j++ {
					fmt.Fprintf(Ferr, "\tWSRN[%d]=%.3g", j, Sd.WSRN[j])
				}
				fmt.Fprintf(Ferr, "\n\t")

				for j := 0; j < Room.Nrp; j++ {
					fmt.Fprintf(Ferr, "\tWSPL[%d]=%.3g", j, Sd.WSPL[j])
				}
				fmt.Fprintf(Ferr, "\n")

				fmt.Fprintf(Ferr, "\t\tWSC=%.3g\n", Sd.WSC)
			}
		}
	}
}

/*
xprtwsrf (Export Wall Surface Temperature for Debugging)

この関数は、室内の各表面（壁、窓など）の温度（`Sd.Ts`）、
平均放射温度（`Sd.Tmrt`）、相当外気温度（`Sd.Te`）、
および日射吸収量（`Sd.RS`）をデバッグ目的で出力します。
これは、表面温度モデルの計算結果が正しいかを確認し、
シミュレーションの挙動を理解するために用いられます。

建築環境工学的な観点:
- **表面温度の確認**: 室内表面温度は、
  居住者の快適性（特に放射快適性）に直接影響を与える重要な要素です。
  この関数は、各表面の温度を出力することで、
  表面温度モデルの計算結果が意図通りに設定されているかを確認できます。
- **平均放射温度の確認**: 平均放射温度（MRT）は、
  居住者が感じる放射熱環境を代表する温度です。
  この出力は、MRTの計算が正しいかを確認し、
  快適性評価の妥当性を検証するのに役立ちます。
- **デバッグと検証**: シミュレーションが期待通りの結果を出さない場合、
  表面温度モデルの誤りが原因であることがあります。
  このデバッグ出力は、表面温度モデルの計算結果が正しいかを確認し、
  問題の特定に役立ちます。

この関数は、建物の熱的挙動を詳細に分析し、
シミュレーションの信頼性を確保するための重要なデバッグ機能を提供します。
*/
func (Rmvls *RMVLS) xprtwsrf() {
	fmt.Println("--- xprtwsrf")

	for n, Sd := range Rmvls.Sd {
		fmt.Printf("  n=%2d  rm=%d nr=%d  Ts=%6.2f  Tmrt=%6.2f  Te=%6.2f  RS=%7.1f\n",
			n, Sd.rm, Sd.n, Sd.Ts, Sd.Tmrt, Sd.Te, Sd.RS)
	}
}

/* -------------------------------------------------------------------- */

/*
xprrmsrf (Export Room Surface Temperatures for Debugging)

この関数は、室内の各表面（壁、窓など）の温度をデバッグ目的で出力します。
これは、表面温度モデルの計算結果が正しいかを確認し、
シミュレーションの挙動を理解するために用いられます。

建築環境工学的な観点:
- **表面温度の確認**: 室内表面温度は、
  居住者の快適性（特に放射快適性）に直接影響を与える重要な要素です。
  この関数は、各表面の温度を出力することで、
  表面温度モデルの計算結果が意図通りに設定されているかを確認できます。
- **デバッグと検証**: シミュレーションが期待通りの結果を出さない場合、
  表面温度モデルの誤りが原因であることがあります。
  このデバッグ出力は、表面温度モデルの計算結果が正しいかを確認し、
  問題の特定に役立ちます。

この関数は、建物の熱的挙動を詳細に分析し、
シミュレーションの信頼性を確保するための重要なデバッグ機能を提供します。
*/
func (Rmvls *RMVLS) xprrmsrf() {
	fmt.Println("--- xprrmsf")

	for n, Sd := range Rmvls.Sd {
		fmt.Printf("  [%d]=%6.2f", n, Sd.Ts)
	}
	fmt.Println()
}

/* -------------------------------------------------------------------- */

/*
xprtwall (Export Wall Internal Temperatures for Debugging)

この関数は、壁体内部の各層の温度をデバッグ目的で出力します。
これは、壁体内部の熱伝導モデルの計算結果が正しいかを確認し、
シミュレーションの挙動を理解するために用いられます。

建築環境工学的な観点:
- **壁体内部温度の確認**: 壁体は熱容量を持つため、
  その内部温度は外気温度や室内温度の変化に対して時間遅れを伴って応答します。
  この関数は、各壁体の各層の温度（`Mw.Tw[m]`）を出力することで、
  壁体内部の熱伝導モデルの計算結果が意図通りに設定されているかを確認できます。
- **デバッグと検証**: シミュレーションが期待通りの結果を出さない場合、
  壁体内部の熱伝導モデルの誤りが原因であることがあります。
  このデバッグ出力は、壁体内部の熱伝導モデルの計算結果が正しいかを確認し、
  問題の特定に役立ちます。
- **モデルの理解**: 壁体内部の温度分布は、
  壁体の熱貫流特性や蓄熱効果を理解する上で重要です。
  この出力は、壁体内部の熱伝導モデルの挙動を視覚的に確認し、
  モデルの理解を深めるのに役立ちます。

この関数は、建物の熱的挙動を詳細に分析し、
シミュレーションの信頼性を確保するための重要なデバッグ機能を提供します。
*/
func (Rmvls *RMVLS) xprtwall() {
	fmt.Println("--- xprtwall")

	for j, Mw := range Rmvls.Mw {
		if Mw.Pc > 0 {
			fmt.Printf("Tw j=%2d", j)

			for m := 0; m < Mw.M; m++ {
				Tw := Mw.Tw[m]
				fmt.Printf("  [%d]=%6.2f", m, Tw)
			}

			fmt.Println()
		}
	}
}

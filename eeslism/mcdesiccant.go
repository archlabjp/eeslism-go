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

/*  mcdessicant.C  */
/*  バッチ式デシカント空調機 */

package eeslism

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

/*
Desielm (Desiccant Element Matrix Setup)

この関数は、デシカント空調機における要素方程式の変数を設定し、
特に空気温度と絶対湿度の相互関係をモデル化します。
これは、デシカント空調機の熱湿気同時交換プロセスをシミュレーションするために不可欠です。

建築環境工学的な観点:
- **デシカント空調の原理**: デシカント空調は、
  吸湿材（デシカント）を用いて空気中の水蒸気を除去することで除湿を行い、
  その後、顕熱交換によって温度を調整するシステムです。
  従来の冷媒を用いた空調システムとは異なり、
  潜熱と顕熱を独立して処理できる点が特徴です。
- **熱湿気同時交換**: デシカント空調機では、
  空気の温度と湿度が同時に変化する熱湿気同時交換が行われます。
  この関数は、出口空気温度の計算が入口空気絶対湿度に依存し、
  出口空気絶対湿度の計算が入口空気温度に依存するという、
  これらの変数の相互関係を要素方程式に組み込むための設定を行います。
  `Upo`や`Upv`といった変数は、
  要素方程式における他の変数の影響を考慮するためのポインターを示唆します。
- **システムモデルの構築**: この設定は、
  デシカント空調機を構成する各要素（吸湿ローター、熱交換器など）の熱湿気交換特性を、
  連立方程式として解くための基礎となります。
  これにより、デシカント空調機が様々な運転条件下で、
  空気の温度と湿度をどのように変化させるかを正確に予測できます。

この関数は、デシカント空調機の熱湿気同時交換プロセスをモデル化し、
室内温湿度環境の予測、潜熱負荷の処理、
およびエネルギー消費量予測を行うための重要な役割を果たします。
*/
func Desielm(Desi []*DESI) {
	for _, desi := range Desi {
		Eot := desi.Cmp.Elouts[0] // 空気温度出口
		Eox := desi.Cmp.Elouts[1] // 空気湿度出口

		elin := Eot.Elins[1]
		elin.Upo = Eox.Elins[0].Upo // 出口空気温度の要素方程式の2つめの変数は絶対湿度
		elin.Upv = Eox.Elins[0].Upo

		elin = Eox.Elins[1]
		elin.Upo = Eot.Elins[0].Upo // 出口絶対湿度の要素方程式の2つめの変数は空気温度
		elin.Upv = Eot.Elins[0].Upo // 空気温度の要素方程式の2つ目の変数（空気入口温度）のupo、upvに空気湿度をつなげる
	}
}

/*
Desiccantdata (Desiccant Data Input)

この関数は、デシカント空調機を構成する吸湿材（シリカゲルなど）や、
デシカント槽の各種仕様を読み込み、対応する構造体に格納します。
これらのデータは、デシカント空調機の除湿性能や熱損失を評価する上で不可欠です。

建築環境工学的な観点:
- **デシカント空調の性能パラメータ**: デシカント空調機の性能は、
  吸湿材の種類、量、形状、そしてデシカント槽の熱的特性に大きく依存します。
  この関数で設定されるパラメータは、以下のようなデシカント空調機の特性を定義します。
  - `Uad`: シリカゲル槽壁面の熱貫流率 [W/m2K]。デシカント槽からの熱損失に影響します。
  - `A`: シリカゲル槽表面積 [m2]。デシカント槽からの熱損失に影響します。
  - `ms`: シリカゲル質量 [g]。吸湿能力に影響します。
  - `r`: シリカゲル平均直径 [cm]。吸湿・脱湿速度に影響します。
  - `rows`: シリカゲル充填密度 [g/cm3]。吸湿能力に影響します。
  - `Vm`, `eps`, `P0`, `kp`, `cps`: シリカゲルの吸湿特性や熱的特性に関連するパラメータ。
- **潜熱負荷処理の最適化**: デシカント空調は、特に潜熱負荷（湿度）の処理に優れています。
  これらのパラメータを適切に設定することで、
  高湿度の外気条件や、室内での水蒸気発生が多い空間（例: 厨房、プール）において、
  効率的な除湿運転を行うためのシステム設計を検討できます。
- **エネルギー消費量への影響**: デシカント空調機は、
  吸湿材の再生に熱エネルギーを必要とします。
  `Uad`や`A`などのパラメータは、デシカント槽からの熱損失に影響し、
  再生に必要なエネルギー消費量に影響を与えます。

この関数は、デシカント空調機の性能をモデル化し、
室内温湿度環境の予測、潜熱負荷の処理、
およびエネルギー消費量予測を行うための重要なデータ入力機能を提供します。
*/
func Desiccantdata(s string, desica *DESICA) int {
	st := strings.Index(s, "=")
	var dt float64
	var id int

	if st == -1 {
		desica.name = s
		desica.r = -999.0
		desica.rows = -999.0
		desica.Uad = -999.0
		desica.A = -999.0
		desica.Vm = 18.0
		desica.eps = 0.4764
		desica.P0 = 0.4
		desica.kp = 0.0012
		desica.cps = 710.0
		desica.ms = -999.0
	} else {
		st++
		dt, _ = strconv.ParseFloat(s[st:], 64)

		switch {
		case s == "Uad": // シリカゲル槽壁面の熱貫流率[W/m2K]
			desica.Uad = dt
		case s == "A": // シリカゲル槽表面積[m2]
			desica.A = dt
		case s == "ms": // シリカゲル質量[g]
			desica.ms = dt
		case s == "r": // シリカゲル平均直径[cm]
			desica.r = dt
		case s == "rows": // シリカゲル充填密度[g/cm3]
			desica.rows = dt
		default:
			id = 1
		}
	}
	return id
}

/*
Desiint (Desiccant Initialization)

この関数は、デシカント空調機のシミュレーションに必要な初期設定と、
各種パラメータの妥当性チェックを行います。
特に、デシカント槽の熱損失係数や、吸湿材と槽内空気の熱伝達面積などを計算します。

建築環境工学的な観点:
- **初期温度・湿度の設定**: シミュレーション開始時のデシカント槽の温度（`desi.Tsold`）や、
  空気の絶対湿度（`desi.Xsold`）を初期化します。
  これらの初期値は、シミュレーションの収束性や、
  初期段階でのデシカント空調機の挙動に影響を与える可能性があります。
- **デシカント槽熱損失係数 (desi.UA)**:
  `desi.UA = Desica.Uad * Desica.A` のように、
  デシカント槽の熱通過率（`Desica.Uad`）と表面積（`Desica.A`）から計算されます。
  この係数は、デシカント槽から周囲への熱損失を評価する上で重要であり、
  熱損失が大きいと、吸湿材の再生に必要なエネルギーが増加し、
  システム全体の効率が低下する可能性があります。
- **吸湿材と槽内空気の熱伝達面積 (desi.Asa)**:
  `desi.Asa`は、吸湿材と槽内空気の間で熱と湿気が交換される有効な面積を表します。
  この面積が大きいほど、吸湿材と空気間の熱湿気交換が活発に行われ、
  デシカント空調機の除湿性能が向上します。
  シリカゲルの質量（`Desica.ms`）や直径（`Desica.r`）、充填密度（`Desica.rows`）などから計算されます。
- **パラメータの妥当性チェック**: `if Desica.Uad < 0.0` のようなエラーチェックは、
  入力されたパラメータが物理的に妥当な範囲内にあるかを確認するために重要です。
  不適切なパラメータは、シミュレーション結果の信頼性を損なう可能性があります。

この関数は、デシカント空調機の性能をモデル化し、
室内温湿度環境の予測、潜熱負荷の処理、
およびエネルギー消費量予測を行うための重要な初期設定と検証機能を提供します。
*/
func Desiint(Desi []*DESI, Simc *SIMCONTL, Compnt []*COMPNT, Wd *WDAT) {
	var Err string
	var Desica *DESICA

	for _, desi := range Desi {

		if desi.Cmp.Envname != "" {
			desi.Tenv = envptr(desi.Cmp.Envname, Simc, Compnt, Wd, nil)
		} else {
			desi.Room = roomptr(desi.Cmp.Roomname, Compnt)
		}

		Desica = desi.Cat

		if Desica.Uad < 0.0 {
			Err = fmt.Sprintf("Name=%s  Uad=%.4g", Desica.name, Desica.Uad)
			Eprint("Desiint", Err)
		}
		if Desica.A < 0.0 {
			Err = fmt.Sprintf("Name=%s  A=%.4g", Desica.name, Desica.A)
			Eprint("Desiint", Err)
		}
		if Desica.r < 0.0 {
			Err = fmt.Sprintf("Name=%s  r=%.4g", Desica.name, Desica.r)
			Eprint("Desiint", Err)
		}
		if Desica.rows < 0.0 {
			Err = fmt.Sprintf("Name=%s  rows=%.4g", Desica.name, Desica.rows)
			Eprint("Desiint", Err)
		}
		if Desica.ms < 0.0 {
			Err = fmt.Sprintf("Name=%s  ms=%.4g", Desica.name, Desica.ms)
			Eprint("Desiint", Err)
		}

		// 初期温度、出入口温度の初期化
		desi.Tsold = 20.0
		desi.Xsold = FNXtr(desi.Tsold, 50.0)
		desi.Tain = desi.Tsold
		desi.Taout = desi.Tsold
		desi.Xain = desi.Xsold
		desi.Xaout = desi.Xsold

		// デシカント槽熱損失係数の計算
		desi.UA = Desica.Uad * Desica.A

		// 吸湿量の初期化
		desi.Pold = Desica.P0

		// シリカゲルと槽内空気の熱伝達面積[m2]
		desi.Asa = 3.0 * Desica.ms * 1000.0 * (1.0 - Desica.eps) / (1.0e4 * (Desica.r / 10.0) * Desica.rows)

		// 逆行列
		desi.UX = make([]float64, 5*5)
		desi.UXC = make([]float64, 5)
	}
}

/* --------------------------- */

/*
Desicfv (Desiccant Characteristic Function Value Calculation)

この関数は、デシカント空調機の運転特性を評価し、
空気の温度と湿度の変化をモデル化するための係数行列を作成します。
これは、デシカント空調機が空気の熱湿気をどのように処理するかをシミュレーションするために不可欠です。

建築環境工学的な観点:
- **熱湿気同時交換のモデル化**: デシカント空調機では、
  吸湿材と空気の間で熱と湿気が同時に交換されます。
  この関数は、この複雑なプロセスを線形方程式系として表現するための係数行列（`U`）と定数行列（`C`）を構築します。
  これにより、入口空気の状態から出口空気の状態を予測できます。
- **熱容量流量 (desi.CG)**:
  空気側の熱容量流量は、デシカント空調機が処理できる熱量に影響します。
  `desi.CG = Spcheat(Eo1.Fluid) * Eo1.G` のように、空気の比熱と質量流量から計算されます。
- **対流熱伝達率 (hsa) と湿気伝達率 (hsad)**:
  - `hsa`: シリカゲルと槽内空気の間の対流熱伝達率。
    この値が大きいほど、顕熱交換が活発に行われます。
  - `hsad`: シリカゲルと槽内空気の間の湿気伝達率。
    この値が大きいほど、潜熱交換（除湿）が活発に行われます。
  これらの伝達率は、デシカント空調機の性能を決定する重要な要素です。
- **吸湿材の吸湿特性 (h, i, j)**:
  `desi.Pold`（吸湿材の含水率）に応じて`h`, `i`, `j`といった係数が変化することから、
  吸湿材の吸湿特性が含水率によって非線形に変化することをモデル化していることが伺えます。
  これは、吸湿材の再生サイクルや、除湿性能の変動を正確に予測するために重要です。
- **逆行列の計算 (Matinv)**:
  構築された係数行列の逆行列を計算することで、
  入口空気の状態から出口空気の状態を直接求めることができるようになります。

この関数は、デシカント空調機の熱湿気同時交換プロセスを詳細にモデル化し、
室内温湿度環境の予測、潜熱負荷の処理、
およびエネルギー消費量予測を行うための重要な役割を果たします。
*/
func Desicfv(Desi []*DESI) {
	var Eo1 *ELOUT
	var h, i, j float64
	var Te, hsa, hsad, hAsa, hdAsa float64
	var Desica *DESICA
	var U, C, Cmat []float64

	N := 5
	N2 := N * N
	for _, desi := range Desi {
		Desica = desi.Cat

		// 係数行列のメモリ確保
		U = make([]float64, N2)
		// 定数行列のメモリ確保
		C = make([]float64, N)

		if desi.Cmp.Envname != "" {
			Te = *desi.Tenv
		} else {
			Te = desi.Room.Tot
		}

		Eo1 = desi.Cmp.Elouts[0]
		// 熱容量流量の計算
		desi.CG = Spcheat(Eo1.Fluid) * Eo1.G

		// シリカゲルと槽内空気の対流熱伝達率の計算
		if Eo1.Cmp.Control == OFF_SW {
			hsa = 4.614
		} else {
			hsa = 40.0
		}

		// シリカゲルと槽内空気の湿気伝達率の計算
		hsad = hsa / Ca

		hAsa = hsa * desi.Asa
		hdAsa = hsad * desi.Asa

		if desi.Pold >= 0.25 {
			h = 0.001319
			i = 0.103335
			j = -0.05416
		} else {
			h = 0.001158
			i = 0.149479
			j = -0.05835
		}

		// 定数行列Cの作成
		Cmat = C
		Cmat[0] = Desica.ms * Desica.cps / DTM * desi.Tsold
		Cmat[1] = desi.UA * Te
		Cmat[3] = Desica.ms / DTM * desi.Pold
		Cmat[4] = -j

		// 係数行列の作成
		U[0*N+0] = Desica.ms*Desica.cps/DTM + hAsa
		U[0*N+1] = -hAsa
		U[0*N+2] = -hdAsa * Ro
		U[0*N+3] = hdAsa * Ro
		U[1*N+0] = -hAsa
		U[1*N+1] = Ca*Eo1.G + hAsa + desi.UA
		U[2*N+2] = Eo1.G + hdAsa
		U[2*N+3] = -hdAsa
		U[3*N+2] = -hdAsa
		U[3*N+3] = hdAsa
		U[3*N+4] = Desica.ms / DTM
		U[4*N+0] = h
		U[4*N+3] = -1.0
		U[4*N+4] = i

		// 逆行列の計算
		Matinv(U, N, N, "<Desicfv U>")

		// 行列のコピー
		matinit(desi.UX, N2)
		matcpy(U, desi.UX, N2)

		// {UXC}=[UX]*{C}の作成
		matinit(desi.UXC, N)
		for ii := 0; ii < N; ii++ {
			for jj := 0; jj < N; jj++ {
				desi.UXC[ii] += desi.UX[ii*N+jj] * C[jj]
			}
		}

		// 出口温度の要素方程式
		Eo1.Coeffo = -1.0
		Eo1.Co = -desi.UXC[1]
		Eo1.Coeffin[0] = desi.UX[1*N+1] * Eo1.G * Ca
		Eo1.Coeffin[1] = desi.UX[1*N+2] * Eo1.G

		// 出口湿度の要素方程式
		Eo2 := desi.Cmp.Elouts[1]
		Eo2.Coeffo = -1.0
		Eo2.Co = -desi.UXC[2]
		Eo2.Coeffin[0] = desi.UX[2*N+2] * Eo2.G
		Eo2.Coeffin[1] = desi.UX[2*N+1] * Eo2.G * Ca
	}
}

/*
Desiene (Desiccant Energy Calculation)

この関数は、デシカント空調機が空気から除去する顕熱量、潜熱量、および全熱量を計算します。
また、デシカント槽からの熱損失も評価し、
デシカント空調機のエネルギー収支を詳細に分析します。

建築環境工学的な観点:
- **顕熱量 (desi.Qs)**:
  デシカント空調機が空気から除去する顕熱量です。
  空気の温度変化に伴う熱量であり、室内の温度制御に直接関係します。
  `desi.CG * (desi.Taout - desi.Tain)` のように、空気側の熱容量流量と入口・出口空気温度差から計算されます。
- **潜熱量 (desi.Ql)**:
  デシカント空調機が空気から除去する潜熱量です。
  空気中の水蒸気量の変化に伴う熱量であり、室内の湿度制御に直接関係します。
  `elo.G * Ro * (desi.Xaout - desi.Xain)` のように、空気流量と入口・出口空気絶対湿度差から計算されます。
  デシカント空調の主要な機能であり、特に高湿度の環境下での除湿能力を評価する上で重要です。
- **全熱量 (desi.Qt)**:
  顕熱量と潜熱量の合計であり、デシカント空調機が空気から除去する総熱量です。
- **デシカント槽からの熱損失 (desi.Qloss)**:
  デシカント槽から周囲へ逃げる熱量です。
  `desi.UA * (Te - desi.Ta)` のように、デシカント槽の熱損失係数と槽内外の温度差から計算されます。
  この熱損失は、デシカント空調機のエネルギー効率に影響を与え、
  特に吸湿材の再生に必要なエネルギー消費量を増加させる要因となります。
- **設置室内部発熱への影響**: `if desi.Room != nil { desi.Room.Qeqp += (-desi.Qloss) }` のように、
  デシカント槽からの熱損失が設置室の内部発熱として計上されることで、
  建物全体の熱収支モデルに組み込まれます。

この関数は、デシカント空調機の熱湿気同時交換性能を定量的に評価し、
潜熱負荷の処理能力、室内温湿度環境の予測、
およびエネルギー消費量予測を行うための重要な役割を果たします。
*/
func Desiene(Desi []*DESI) {
	Sin := make([]float64, 5)
	S := make([]float64, 5)

	N := 5
	//N2 := N * N
	for _, desi := range Desi {
		matinit(Sin, N)
		matinit(S, N)
		elo := desi.Cmp.Elouts[0]
		elox := desi.Cmp.Elouts[1]
		elix := elo.Elins[1]
		desi.Tain = elo.Elins[0].Sysvin
		desi.Xain = elix.Sysvin

		var Te float64
		if desi.Cmp.Envname != "" {
			Te = *desi.Tenv
		} else {
			Te = desi.Room.Tot
		}

		desi.Taout = elo.Sysv
		desi.Xaout = elox.Sysv

		// 入口状態行列Sinの作成
		Sin[1] = Ca * elo.G * desi.Tain
		Sin[2] = elo.G * desi.Xain
		// 内部状態値の計算
		for ii := 0; ii < N; ii++ {
			for jj := 0; jj < N; jj++ {
				S[ii] += desi.UX[ii*N+jj] * Sin[jj]
			}
			S[ii] += desi.UXC[ii]
		}
		// 変数への格納
		desi.Tsold = S[0]
		desi.Ta = S[1]
		desi.Xa = S[2]
		desi.Xsold = S[3]
		desi.Pold = S[4]
		desi.RHold = FNRhtx(desi.Tsold, desi.Xsold)

		// 顕熱の計算
		desi.Qs = desi.CG * (desi.Taout - desi.Tain)
		desi.Ql = elo.G * Ro * (desi.Xaout - desi.Xain)
		desi.Qt = desi.Qs + desi.Ql

		// デシカント槽からの熱損失の計算
		desi.Qloss = desi.UA * (Te - desi.Ta)

		// 設置室内部発熱の計算
		if desi.Room != nil {
			desi.Room.Qeqp += (-desi.Qloss)
		}
	}
}

/*
Desivptr (Desiccant Internal Variable Pointer Setting)

この関数は、デシカント空調機の制御で使用される内部変数（吸湿材温度、含水率、相対湿度など）へのポインターを設定します。
これにより、デシカント空調機の運転を特定の目標値に追従させる制御をモデル化できます。

建築環境工学的な観点:
- **デシカント空調の制御**: デシカント空調機は、
  吸湿材の温度や含水率、処理空気の相対湿度などを監視し、
  それに基づいて再生熱量や空気流量を制御することで、
  目標とする温湿度環境を維持します。
- **制御対象の指定**: `key[1]`が`"Ts"`の場合、吸湿材温度（`Desi.Tsold`）を、
  `"xs"`の場合、吸湿材の含水率（`Desi.Xsold`）を、
  `"RH"`の場合、吸湿材の相対湿度（`Desi.RHold`）を制御対象とすることを意味します。
  `vptr.Ptr`は、これらの変数へのポインターを設定し、
  `vptr.Type = VAL_CTYPE`は、そのポインターが制御値であることを示します。
- **フィードバック制御の基礎**: このポインター設定は、
  デシカント空調機のフィードバック制御の基礎となります。
  シミュレーションの各時間ステップで、
  現在の内部状態と目標値を比較し、その差に基づいて運転を調整します。
  これにより、室内温湿度環境の安定化や、除湿効率の向上を図ることができます。

この関数は、デシカント空調機の制御ロジックをモデル化し、
室内温湿度環境の予測、潜熱負荷の処理、
およびエネルギー消費量予測を行うための重要な役割を果たします。
*/
func Desivptr(key []string, Desi *DESI) (VPTR, error) {
	var err error
	var vptr VPTR

	switch key[1] {
	case "Ts":
		vptr.Ptr = &Desi.Tsold
		vptr.Type = VAL_CTYPE
	case "xs":
		vptr.Ptr = &Desi.Xsold
		vptr.Type = VAL_CTYPE
	case "RH":
		vptr.Ptr = &Desi.RHold
		vptr.Type = VAL_CTYPE
	default:
		err = errors.New("'Ts', 'xs' or 'RH' is expected")
	}

	return vptr, err
}

///* ---------------------------*/
//
func Desiprint(fo io.Writer, id int, Desi []*DESI) {
	switch id {
	case 0:
		if len(Desi) > 0 {
			fmt.Fprintf(fo, "%s %d\n", DESI_TYPE, len(Desi))
		}
		for _, desi := range Desi {
			fmt.Fprintf(fo, " %s 1 14\n", desi.Name)
		}
	case 1:
		for _, desi := range Desi {
			fmt.Fprintf(fo, "%s_c c c %s_G m f %s_Ts t f %s_Ti t f %s_To t f %s_Qs q f ", desi.Name, desi.Name, desi.Name, desi.Name, desi.Name, desi.Name)
			fmt.Fprintf(fo, "%s_xs x f %s_RHs r f %s_xi x f %s_xo x f %s_Ql q f %s_Qt q f ", desi.Name, desi.Name, desi.Name, desi.Name, desi.Name, desi.Name)
			fmt.Fprintf(fo, "%s_Qls q f %s_P m f\n", desi.Name, desi.Name)
		}
	default:
		for _, desi := range Desi {
			fmt.Fprintf(fo, "%c %6.4g %4.1f %4.1f %4.1f %2.0f  ", desi.Cmp.Elouts[0].Control, desi.Cmp.Elouts[0].G, desi.Tsold, desi.Tain, desi.Taout, desi.Qs)
			fmt.Fprintf(fo, "%.3f %.0f %.3f %.3f %2.0f %2.0f  ", desi.Xsold, desi.RHold, desi.Xain, desi.Xaout, desi.Ql, desi.Qt)
			fmt.Fprintf(fo, "%.0f %.3f\n", desi.Qloss, desi.Pold)
		}
	}
}

/*
Desidyint (Desiccant Daily Integration Initialization)

この関数は、デシカント空調機の日積算値（日ごとの入口・出口空気温度、絶対湿度、
吸湿材温度、顕熱量、潜熱量、全熱量、熱損失など）をリセットします。
これは、日単位でのデシカント空調機の運転状況や熱湿気交換量を集計し、
空調システムの性能評価やエネルギー消費量の分析に用いられます。

建築環境工学的な観点:
- **日単位の性能評価**: デシカント空調機の運転状況は、日中の熱湿気負荷変動に応じて大きく変化します。
  日積算値を集計することで、日ごとの顕熱・潜熱負荷の割合、
  デシカント空調機の稼働時間、部分負荷運転の割合などを把握できます。
  これにより、特定の日の空調負荷特性を分析したり、
  デシカント空調機の運転効率を日単位で評価したりすることが可能になります。
- **運用改善の指標**: 日積算データは、空調システムの運用改善のための重要な指標となります。
  例えば、外気温度や湿度などの気象条件とデシカント空調機の熱湿気交換量の関係を分析したり、
  設定温度や換気量などの運用条件がデシカント空調機の性能に与える影響を評価したりすることで、
  より効率的な運転方法を見つけることができます。
- **データ集計の準備**: この関数は、新しい日の集計を開始する前に、
  前日のデータをクリアする役割を担います。
  `svdyint`や`qdyint`といった関数は、
  それぞれ温度、湿度、熱量などの日積算値をリセットするためのものです。

この関数は、デシカント空調機の運転状況と熱湿気交換量を日単位で詳細に分析し、
空調システムの運用改善や省エネルギー対策の効果評価を行うための基礎的な役割を果たします。
*/
func Desidyint(Desi []*DESI) {
	for _, desi := range Desi {
		svdyint(&desi.Tidy)
		svdyint(&desi.Tsdy)
		svdyint(&desi.Tody)
		svdyint(&desi.xidy)
		svdyint(&desi.xsdy)
		svdyint(&desi.xody)
		qdyint(&desi.Qsdy)
		qdyint(&desi.Qldy)
		qdyint(&desi.Qtdy)
		qdyint(&desi.Qlsdy)
	}
}

/*
Desiday (Desiccant Daily Data Aggregation)

この関数は、デシカント空調機の運転データ（入口・出口空気温度、絶対湿度、
吸湿材温度、顕熱量、潜熱量、全熱量、熱損失など）を、日単位で集計します。
これにより、デシカント空調機の性能評価やエネルギー消費量の分析が可能になります。

建築環境工学的な観点:
- **日次集計 (svdaysum, qdaysum)**:
  日次集計は、デシカント空調機の運転状況を日単位で詳細に把握するために重要です。
  例えば、特定の日の顕熱・潜熱負荷の変動に対するデシカント空調機の応答、
  あるいは日中のピーク負荷時の熱湿気交換量などを分析できます。
  これにより、日ごとの運用改善点を見つけ出すことが可能になります。
- **データ分析の基礎**: この関数で集計されるデータは、
  デシカント空調機の性能評価、熱湿気交換量のベンチマーキング、
  省エネルギー対策の効果検証、そして運用改善のための意思決定の基礎となります。

この関数は、デシカント空調機の運転状況と熱湿気交換量を日単位で詳細に分析し、
空調システムの運用改善や省エネルギー対策の効果評価を行うための重要なデータ集計機能を提供します。
*/
func Desiday(Mon, Day, ttmm int, Desi []*DESI, Nday, SimDayend int) {
	// Mo := Mon - 1
	// tt := ConvertHour(ttmm)

	for _, desi := range Desi {
		// 日集計
		svdaysum(int64(ttmm), desi.Cmp.Control, desi.Tain, &desi.Tidy)
		svdaysum(int64(ttmm), desi.Cmp.Control, desi.Taout, &desi.Tody)
		svdaysum(int64(ttmm), desi.Cmp.Control, desi.Tsold, &desi.Tsdy)
		svdaysum(int64(ttmm), desi.Cmp.Control, desi.Xain, &desi.xidy)
		svdaysum(int64(ttmm), desi.Cmp.Control, desi.Xaout, &desi.xody)
		svdaysum(int64(ttmm), desi.Cmp.Control, desi.Xsold, &desi.xsdy)
		qdaysum(int64(ttmm), desi.Cmp.Control, desi.Qs, &desi.Qsdy)
		qdaysum(int64(ttmm), desi.Cmp.Control, desi.Ql, &desi.Qldy)
		qdaysum(int64(ttmm), desi.Cmp.Control, desi.Qt, &desi.Qtdy)
		qdaysum(int64(ttmm), desi.Cmp.Control, desi.Qloss, &desi.Qlsdy)
	}
}

func Desidyprt(fo io.Writer, id int, Desi []*DESI) {
	switch id {
	case 0:
		if len(Desi) > 0 {
			fmt.Fprintf(fo, "%s %d\n", DESI_TYPE, len(Desi))
		}
		for _, desi := range Desi {
			fmt.Fprintf(fo, " %s 1 68\n", desi.Name)
		}
	case 1:
		for _, desi := range Desi {

			fmt.Fprintf(fo, "%s_Ht H d %s_Ti T f ", desi.Name, desi.Name)
			fmt.Fprintf(fo, "%s_ttn h d %s_Tin t f %s_ttm h d %s_Tim t f\n", desi.Name, desi.Name, desi.Name, desi.Name)

			fmt.Fprintf(fo, "%s_Ht H d %s_To T f ", desi.Name, desi.Name)
			fmt.Fprintf(fo, "%s_ttn h d %s_Ton t f %s_ttm h d %s_Tom t f\n", desi.Name, desi.Name, desi.Name, desi.Name)

			fmt.Fprintf(fo, "%s_Ht H d %s_Ts T f ", desi.Name, desi.Name)
			fmt.Fprintf(fo, "%s_ttn h d %s_Tsn t f %s_ttm h d %s_Tsm t f\n", desi.Name, desi.Name, desi.Name, desi.Name)

			fmt.Fprintf(fo, "%s_Ht H d %s_xi T f ", desi.Name, desi.Name)
			fmt.Fprintf(fo, "%s_ttn h d %s_xin t f %s_ttm h d %s_xim t f\n", desi.Name, desi.Name, desi.Name, desi.Name)

			fmt.Fprintf(fo, "%s_Ht H d %s_xo T f ", desi.Name, desi.Name)
			fmt.Fprintf(fo, "%s_ttn h d %s_xon t f %s_ttm h d %s_xom t f\n", desi.Name, desi.Name, desi.Name, desi.Name)

			fmt.Fprintf(fo, "%s_Ht H d %s_xs T f ", desi.Name, desi.Name)
			fmt.Fprintf(fo, "%s_ttn h d %s_xsn t f %s_ttm h d %s_xsm t f\n", desi.Name, desi.Name, desi.Name, desi.Name)

			fmt.Fprintf(fo, "%s_Hsh H d %s_Qsh Q f %s_Hsc H d %s_Qsc Q f\n", desi.Name, desi.Name, desi.Name, desi.Name)
			fmt.Fprintf(fo, "%s_tsh h d %s_qsh q f %s_tsc h d %s_qsc q f\n", desi.Name, desi.Name, desi.Name, desi.Name)

			fmt.Fprintf(fo, "%s_Hlh H d %s_Qlh Q f %s_Hlc H d %s_Qlc Q f\n", desi.Name, desi.Name, desi.Name, desi.Name)
			fmt.Fprintf(fo, "%s_tlh h d %s_qlh q f %s_tlc h d %s_qlc q f\n", desi.Name, desi.Name, desi.Name, desi.Name)

			fmt.Fprintf(fo, "%s_Hth H d %s_Qth Q f %s_Htc H d %s_Qtc Q f\n", desi.Name, desi.Name, desi.Name, desi.Name)
			fmt.Fprintf(fo, "%s_tth h d %s_qth q f %s_ttc h d %s_qtc q f\n", desi.Name, desi.Name, desi.Name, desi.Name)

			fmt.Fprintf(fo, "%s_Hlsh H d %s_Qlsh Q f %s_Hlsc H d %s_Qlsc Q f\n", desi.Name, desi.Name, desi.Name, desi.Name)
			fmt.Fprintf(fo, "%s_tlsh h d %s_qlsh q f %s_tlsc h d %s_qlsc q f\n", desi.Name, desi.Name, desi.Name, desi.Name)
		}
	default:
		for _, desi := range Desi {
			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ", desi.Tidy.Hrs, desi.Tidy.M, desi.Tidy.Mntime, desi.Tidy.Mn, desi.Tidy.Mxtime, desi.Tidy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ", desi.Tody.Hrs, desi.Tody.M, desi.Tody.Mntime, desi.Tody.Mn, desi.Tody.Mxtime, desi.Tody.Mx)
			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ", desi.Tsdy.Hrs, desi.Tsdy.M, desi.Tsdy.Mntime, desi.Tsdy.Mn, desi.Tsdy.Mxtime, desi.Tsdy.Mx)

			fmt.Fprintf(fo, "%1d %.4f %1d %.4f %1d %.4f ", desi.xidy.Hrs, desi.xidy.M, desi.xidy.Mntime, desi.xidy.Mn, desi.xidy.Mxtime, desi.xidy.Mx)
			fmt.Fprintf(fo, "%1d %.4f %1d %.4f %1d %.4f ", desi.xody.Hrs, desi.xody.M, desi.xody.Mntime, desi.xody.Mn, desi.xody.Mxtime, desi.xody.Mx)
			fmt.Fprintf(fo, "%1d %.4f %1d %.4f %1d %.4f ", desi.xsdy.Hrs, desi.xsdy.M, desi.xsdy.Mntime, desi.xsdy.Mn, desi.xsdy.Mxtime, desi.xsdy.Mx)

			fmt.Fprintf(fo, "%1d %3.1f ", desi.Qsdy.Hhr, desi.Qsdy.H)
			fmt.Fprintf(fo, "%1d %3.1f ", desi.Qsdy.Chr, desi.Qsdy.C)
			fmt.Fprintf(fo, "%1d %2.0f ", desi.Qsdy.Hmxtime, desi.Qsdy.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f ", desi.Qsdy.Cmxtime, desi.Qsdy.Cmx)

			fmt.Fprintf(fo, "%1d %3.1f ", desi.Qldy.Hhr, desi.Qldy.H)
			fmt.Fprintf(fo, "%1d %3.1f ", desi.Qldy.Chr, desi.Qldy.C)
			fmt.Fprintf(fo, "%1d %2.0f ", desi.Qldy.Hmxtime, desi.Qldy.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f ", desi.Qldy.Cmxtime, desi.Qldy.Cmx)

			fmt.Fprintf(fo, "%1d %3.1f ", desi.Qtdy.Hhr, desi.Qtdy.H)
			fmt.Fprintf(fo, "%1d %3.1f ", desi.Qtdy.Chr, desi.Qtdy.C)
			fmt.Fprintf(fo, "%1d %2.0f ", desi.Qtdy.Hmxtime, desi.Qtdy.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f ", desi.Qtdy.Cmxtime, desi.Qtdy.Cmx)

			fmt.Fprintf(fo, "%1d %3.1f ", desi.Qlsdy.Hhr, desi.Qlsdy.H)
			fmt.Fprintf(fo, "%1d %3.1f ", desi.Qlsdy.Chr, desi.Qlsdy.C)
			fmt.Fprintf(fo, "%1d %2.0f ", desi.Qlsdy.Hmxtime, desi.Qlsdy.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f\n", desi.Qlsdy.Cmxtime, desi.Qlsdy.Cmx)
		}
	}
}

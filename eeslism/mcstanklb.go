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

/* mcstanklb.c  */

package eeslism

import "math"

const TSTOLE = 0.04

/*
stoint (Storage Tank Initialization for Layers)

この関数は、蓄熱槽を複数の仮想的な層に分割し、
各層の熱容量、熱損失係数、および初期温度を設定します。
これは、蓄熱槽内部の温度分布や熱損失を詳細にモデル化するために不可欠です。

建築環境工学的な観点:
- **蓄熱槽の仮想分割 (N)**:
  蓄熱槽内部の温度分布をより正確にモデル化するために、
  蓄熱槽を`N`個の仮想的な層に分割して扱います。
  これにより、蓄熱槽内の温度成層（温度の異なる層が形成される現象）を再現し、
  蓄熱槽の有効利用率を評価できます。
- **各層の熱容量 (Mdt)**:
  `Mdt[i] = (Cw * Row * Vol / float64(N)) / DTM` のように、
  各層の体積（`Vol / float64(N)`）、熱媒の比熱（`Cw`）、密度（`Row`）、
  そして時間ステップ（`DTM`）から計算されます。
  これは、各層がどれだけの熱を蓄えることができるかを示します。
- **各層の熱損失係数 (KS)**:
  `KS[i] = KAside / float64(N)` のように、
  蓄熱槽の側面からの熱損失係数（`KAside`）を各層に均等に配分します。
  さらに、最上層には上面からの熱損失（`KAtop`）、
  最下層には下面からの熱損失（`KAbtm`）が加算されます。
  これにより、蓄熱槽の断熱性能や、各層からの熱損失を詳細にモデル化できます。
- **初期温度の設定 (Tss, Tssold)**:
  各層の初期温度（`Tssold`）を現在の温度（`Tss`）に設定します。
  これにより、シミュレーション開始時の蓄熱槽の状態を正確にモデル化できます。

この関数は、蓄熱槽の熱的挙動を詳細にモデル化し、
熱負荷平準化、エネルギー消費量予測、
および省エネルギー対策の検討を行うための重要な初期設定機能を提供します。
*/
func stoint(N int, Vol float64, KAside float64, KAtop float64, KAbtm float64,
	dvol, Mdt, KS, Tss, Tssold []float64, Jva, Jvb *int) {

	for i := 0; i < N; i++ {
		dvol[i] = Vol / float64(N)
		Mdt[i] = (Cw * Row * Vol / float64(N)) / DTM
		KS[i] = KAside / float64(N)

		Tss[i] = Tssold[i]
	}

	KS[0] += KAtop
	KS[N-1] += KAbtm

	*Jva = 0
	*Jvb = 0
}

/* ----------------------------------------------------------- */

/*
stofc (Storage Tank Coefficient Calculation)

この関数は、蓄熱槽内部の熱伝達を記述する係数行列（`B`）と定数行列（`R`, `d`, `fg`）を作成します。
これは、蓄熱槽内部の温度分布や熱交換を詳細にモデル化するために不可欠です。

建築環境工学的な観点:
- **熱伝達方程式の構築**: 蓄熱槽内部の各層における熱収支は、
  熱容量、熱損失、層間の熱伝達、そして熱媒の流入・流出によって決まります。
  この関数は、これらの要素を連立方程式として表現するための係数行列`B`を構築します。
  `Mdt[j]`は各層の熱容量、`KS[j]`は各層の熱損失係数、
  `gxr`は温度成層の度合いを示す係数であり、層間の熱伝達に影響します。
- **熱媒の流入・流出の考慮**: `Jcin`（流入層）、`Jcout`（流出層）、
  `cGwin`（流入熱容量流量）、`EGwin`（有効熱容量流量）などのパラメータを用いて、
  熱媒の流入・流出が蓄熱槽内部の温度分布に与える影響をモデル化します。
  特に、内蔵熱交換器（`ihex == 'y'`）がある場合は、
  その効率（`ihxeff`）も考慮されます。
- **温度分布の計算**: 係数行列`B`の逆行列を計算し（`Matinv`）、
  定数行列`R`と乗算する（`Matmalv`）ことで、
  蓄熱槽内部の各層の温度（`d`）を計算できます。
- **熱交換器の係数 (fg)**:
  `fg`は、内蔵熱交換器を介した熱交換に関する係数であり、
  熱源設備や熱利用設備との熱交換をモデル化するために用いられます。

この関数は、蓄熱槽の熱的挙動を詳細にモデル化し、
熱負荷平準化、エネルギー消費量予測、
および省エネルギー対策の検討を行うための重要な役割を果たします。
*/
func stofc(N, Nin int, Jcin, Jcout []int,
	ihex []rune, ihxeff []float64, Jva, Jvb int, Mdt, KS []float64,
	gxr float64, Tenv *float64, Tssold, cGwin, EGwin, B, R, d, fg []float64) {
	N2 := N * N
	for j := 0; j < N2; j++ {
		B[j] = 0.0
	}

	for j := 0; j < N; j++ {
		B[j*N+j] = Mdt[j] + KS[j]
		R[j] = Mdt[j]*Tssold[j] + KS[j]**Tenv
	}

	for j := 0; j < N-1; j++ {
		B[j*N+j+1] = -Mdt[j] * gxr
		B[(j+1)*N+j] = -Mdt[j+1] * gxr
	}

	if Jva >= 0 {
		for j := Jva; j <= Jvb; j++ {
			B[j*N+j+1] = -Mdt[j] * 1.0e6
			B[(j+1)*N+j] = -Mdt[j] * 1.0e6
		}
	}

	for i := 0; i < Nin; i++ {
		Jin := Jcin[i]
		if cGwin[i] > 0.0 {
			B[Jin*N+Jin] += EGwin[i]

			if Jin < Jcout[i] {
				for j := Jin + 1; j <= Jcout[i]; j++ {
					B[j*N+j-1] -= cGwin[i]
				}
			} else if Jin > Jcout[i] {
				for j := Jcout[i]; j < Jin; j++ {
					B[j*N+j+1] -= cGwin[i]
				}
			}
		}
	}

	for j := 1; j < N-1; j++ {
		B[j*N+j] += math.Abs(B[j*N+j-1]) + math.Abs(B[j*N+j+1])
	}

	B[0] += math.Abs(B[1])
	B[N*N-1] += math.Abs(B[N*N-2])

	Matinv(B, N, N, "<stofc>")
	Matmalv(B, R, N, N, d)

	fgIndex := 0
	for k := 0; k < Nin; k++ {
		Jo := Jcout[k]
		if ihex[k] == 'y' {
			d[Jo] *= ihxeff[k]
			for i := 0; i < Nin; i++ {
				Jin := Jcin[i]
				fg[fgIndex] = B[Jo*N+Jin] * EGwin[i] * ihxeff[k]
				if k == i {
					fg[fgIndex] += (1.0 - ihxeff[k])
				}
				fgIndex++
			}
		} else {
			for i := 0; i < Nin; i++ {
				Jin := Jcin[i]
				fg[fgIndex] = B[Jo*N+Jin] * EGwin[i]
				fgIndex++
			}
		}
	}
}

/* -------------------------------------------------------------- */

/*
stotss (Storage Tank Water Temperature Calculation)

この関数は、蓄熱槽内部の各層の水温を計算します。
これは、蓄熱槽の熱的挙動を詳細にモデル化し、
熱負荷平準化の効果を評価するために不可欠です。

建築環境工学的な観点:
- **熱収支方程式の解法**: この関数は、
  `stofc`関数で構築された係数行列`B`と定数行列`R`を用いて、
  蓄熱槽内部の各層の熱収支方程式を解き、
  各層の温度（`Tss`）を算出します。
  `Matmalv(B, R, N, N, Tss)` は、行列の乗算によって温度を計算しています。
- **熱媒の流入による温度変化**: `R[Jin] += EGwin[i] * Twin[i]` のように、
  各入水ポートからの熱媒の流入が、
  対応する層の熱収支に与える影響を考慮しています。
  `EGwin[i]`は有効熱容量流量、`Twin[i]`は流入熱媒の温度です。
- **温度成層の維持**: 蓄熱槽内の温度成層を維持することは、
  蓄熱効率を高める上で重要です。
  この関数で計算される温度分布は、
  温度成層が適切に維持されているか、
  あるいは崩壊の兆候がないかを評価するために用いられます。

この関数は、蓄熱槽の熱的挙動を詳細にモデル化し、
熱負荷平準化、エネルギー消費量予測、
および省エネルギー対策の検討を行うための重要な役割を果たします。
*/
func stotss(N, Nin int, Jcin []int, B, R, EGwin, Twin, Tss []float64) {
	for i := 0; i < Nin; i++ {
		Jin := Jcin[i]
		R[Jin] += EGwin[i] * Twin[i]
	}

	Matmalv(B, R, N, N, Tss)
}

/* -------------------------------------------------------------- */

/*
stotsexm (Storage Tank Temperature Stratification Examination)

この関数は、蓄熱槽内部の水温分布をチェックし、
水温分布の逆転（温度成層の崩壊）が発生していないかを判定します。
これは、蓄熱槽の効率的な運用と、シミュレーションの安定性に不可欠です。

建築環境工学的な観点:
- **温度成層の維持**: 蓄熱槽は、温度の異なる水が層状に分かれる「温度成層」を形成することで、
  熱源設備からの高温水と熱利用設備への低温水を効率的に供給できます。
  温度成層が崩れると、蓄熱槽の有効利用率が低下し、
  熱源設備や熱利用設備の運転効率に悪影響を与える可能性があります。
- **水温分布逆転の検出**: この関数は、
  下層の温度が上層の温度よりも高くなる「水温分布の逆転」を検出します。
  `Tss[j+1] > (Tss[j] + TSTOLE)` の条件は、
  隣接する層間で一定の温度差（`TSTOLE`）を超えて逆転が発生した場合に、
  それを異常と判断することを示唆します。
- **シミュレーションの安定性**: 水温分布の逆転は、
  シミュレーションモデルの不安定性を示す場合があり、
  正確な結果を得るためには、この問題を解決する必要があります。
  `*cfcalc = 'y'` と設定されることで、
  シミュレーションの再計算が必要であることを示します。
- **運用改善への示唆**: 水温分布の逆転が頻繁に発生する場合、
  それは蓄熱槽の設計（例: 入出力ポートの位置、槽の形状）や、
  運転方法（例: 流量制御、温度制御）に問題がある可能性を示唆します。
  このチェックは、蓄熱槽の運用改善のための重要な情報を提供します。

この関数は、蓄熱槽の熱的挙動を詳細にモデル化し、
熱負荷平準化、エネルギー消費量予測、
および省エネルギー対策の検討を行うための重要な役割を果たします。
*/
func stotsexm(N int, Tss []float64, Jva, Jvb *int, dtankF []rune, cfcalc *rune) {
	*Jvb = -1
	*Jva = -1

	for j := N - 2; j >= 0; j-- {
		if dtankF[j] == TANK_FULL {
			if Tss[j+1] > (Tss[j] + TSTOLE) {
				*Jvb = j
			}
			if *Jvb >= 0 {
				break
			}
		}
	}

	if *Jvb >= 0 {
		for j := *Jvb - 1; j >= 0; j-- {
			if dtankF[j] == TANK_FULL {
				if Tss[*Jvb+1] > (Tss[j] + TSTOLE) {
					*Jva = j
				}
			}
		}
		if *Jva == -1 {
			*Jva = *Jvb
		}
	}

	if *Jva < 0 {
		*cfcalc = 'n'
	} else {
		*cfcalc = 'y'
	}
}

/*-----------------------------------------------------------------*/

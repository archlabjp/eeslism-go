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

import "math"

// 水、空気の物性値の計算
// パソコンによる空気調和計算法より作成

/*
FNarow (Function for Air Density)

この関数は、与えられた温度（`dblT`）における空気の密度を計算します。

建築環境工学的な観点:
- **空気の密度と熱容量**: 空気の密度は、
  換気量から質量流量を計算したり、
  空気の熱容量を計算したりする際に不可欠な要素です。
  空気の密度は温度によって変化するため、
  正確な熱負荷計算や空調システムの設計には、
  温度に応じた密度を考慮する必要があります。
- **換気量計算の基礎**: 換気システムでは、
  通常、体積流量（m3/s）で換気量が指定されますが、
  熱負荷計算では質量流量（kg/s）を用いることが多いため、
  この関数で計算される密度が変換に用いられます。

この関数は、建物の熱負荷計算、換気システム設計、
およびエネルギー消費量予測を行うための基礎的な役割を果たします。
*/
func FNarow(dblT float64) float64 {
	return 1.293 / (1.0 + dblT/273.15)
}

/*
FNac (Function for Air Specific Heat)

この関数は、空気の比熱を返します。

建築環境工学的な観点:
- **空気の比熱と顕熱**: 空気の比熱は、
  空気の温度変化に伴う顕熱量を計算する際に不可欠な要素です。
  顕熱は、空調システムの冷暖房負荷の主要な構成要素の一つです。
- **熱負荷計算の基礎**: 空調システムが室内の温度を変化させるために必要な熱量を計算する際に、
  この比熱が用いられます。

この関数は、建物の熱負荷計算、空調システム設計、
およびエネルギー消費量予測を行うための基礎的な役割を果たします。
*/
func FNac() float64 {
	return 1005.0
}

/*
FNalam (Function for Air Thermal Conductivity)

この関数は、与えられた温度（`dblT`）における空気の熱伝導率を計算します。

建築環境工学的な観点:
- **空気の熱伝導率と熱伝達**: 空気は、
  断熱材として広く用いられる材料であり、
  その熱伝導率は熱伝達計算において重要なパラメータです。
  窓の空気層や、壁体内部の空気層からの熱伝達をモデル化する際に用いられます。
- **温度依存性**: 空気の熱伝導率は温度によって変化するため、
  正確な熱伝達計算には、
  温度に応じた熱伝導率を考慮する必要があります。

この関数は、建物の熱負荷計算、断熱設計、
およびエネルギー消費量予測を行うための基礎的な役割を果たします。
*/
func FNalam(dblT float64) float64 {
	var dblTemp float64

	if dblT > -50.0 && dblT < 100.0 {
		dblTemp = 0.0241 + 0.000077*dblT
	} else {
		dblTemp = -999.0
	}

	return dblTemp
}

/*
FNamew (Function for Air Viscosity)

この関数は、与えられた温度（`dblT`）における空気の粘性係数を計算します。

建築環境工学的な観点:
- **空気の粘性と熱伝達**: 空気の粘性係数は、
  空気の流れ（対流）による熱伝達をモデル化する際に重要なパラメータです。
  特に、自然対流や強制対流による熱伝達率の計算に用いられます。
- **温度依存性**: 空気の粘性係数は温度によって変化するため、
  正確な熱伝達計算には、
  温度に応じた粘性係数を考慮する必要があります。

この関数は、建物の熱負荷計算、換気システム設計、
およびエネルギー消費量予測を行うための基礎的な役割を果たします。
*/
func FNamew(dblT float64) float64 {
	return (0.0074237 / (dblT + 390.15) * math.Pow((dblT+273.15)/293.15, 1.5))
}

/*
FNanew (Function for Air Kinematic Viscosity)

この関数は、与えられた温度（`dblT`）における空気の動粘性係数を計算します。
動粘性係数は、空気の流れ（対流）による熱伝達をモデル化する際に重要なパラメータです。

建築環境工学的な観点:
- **空気の動粘性と熱伝達**: 空気の動粘性係数は、
  空気の流れ（対流）による熱伝達をモデル化する際に重要なパラメータです。
  特に、レイノルズ数やグラスホフ数などの無次元数の計算に用いられ、
  自然対流や強制対流による熱伝達率の計算に不可欠です。
- **温度依存性**: 空気の動粘性係数は温度によって変化するため、
  正確な熱伝達計算には、
  温度に応じた動粘性係数を考慮する必要があります。

この関数は、建物の熱負荷計算、換気システム設計、
およびエネルギー消費量予測を行うための基礎的な役割を果たします。
*/
func FNanew(dblT float64) float64 {
	return FNamew(dblT) / FNarow(dblT)
}

/*
FNabeta (Function for Air Thermal Expansion Coefficient)

この関数は、与えられた温度（`dblT`）における空気の膨張率を計算します。

建築環境工学的な観点:
- **空気の膨張率と自然対流**: 空気の膨張率は、
  温度変化による空気の密度変化の度合いを示すもので、
  自然対流による熱伝達をモデル化する際に重要なパラメータです。
  特に、グラスホフ数などの無次元数の計算に用いられ、
  自然対流による熱伝達率の計算に不可欠です。
- **温度依存性**: 空気の膨張率は温度によって変化するため、
  正確な熱伝達計算には、
  温度に応じた膨張率を考慮する必要があります。

この関数は、建物の熱負荷計算、換気システム設計、
およびエネルギー消費量予測を行うための基礎的な役割を果たします。
*/
func FNabeta(dblT float64) float64 {
	return 1.0 / (dblT + 273.15)
}

/*
FNwrow (Function for Water Density)

この関数は、与えられた温度（`dblT`）における水の密度を計算します。

建築環境工学的な観点:
- **水の密度と熱容量**: 水の密度は、
  流量から質量流量を計算したり、
  水の熱容量を計算したりする際に不可欠な要素です。
  水の密度は温度によって変化するため、
  正確な熱負荷計算や熱搬送システムの設計には、
  温度に応じた密度を考慮する必要があります。
- **熱搬送量計算の基礎**: 熱搬送システムでは、
  通常、体積流量（m3/s）で流量が指定されますが、
  熱負荷計算では質量流量（kg/s）を用いることが多いため、
  この関数で計算される密度が変換に用いられます。

この関数は、建物の熱負荷計算、熱搬送システム設計、
およびエネルギー消費量予測を行うための基礎的な役割を果たします。
*/
func FNwrow(dblT float64) float64 {
	var dblTemp float64

	if dblT > 0.0 && dblT < 100.0 {
		dblTemp = 1000.5 - 0.068737*dblT - 0.0035781*dblT*dblT
	} else if dblT < 200.0 {
		dblTemp = 1008.7 - 0.28735*dblT - 0.0021643*dblT*dblT
	} else {
		dblTemp = -999.0
	}

	return dblTemp
}

/*
FNwc (Function for Water Specific Heat)

この関数は、与えられた温度（`dblT`）における水の比熱を計算します。

建築環境工学的な観点:
- **水の比熱と顕熱**: 水の比熱は、
  水の温度変化に伴う顕熱量を計算する際に不可欠な要素です。
  水は、熱媒として広く用いられるため、
  熱搬送システムや蓄熱システムにおける熱量計算に重要です。
- **温度依存性**: 水の比熱は温度によって変化するため、
  正確な熱量計算には、
  温度に応じた比熱を考慮する必要があります。

この関数は、建物の熱負荷計算、熱搬送システム設計、
およびエネルギー消費量予測を行うための基礎的な役割を果たします。
*/
func FNwc(dblT float64) float64 {
	var dblTemp float64

	if dblT > 0.0 && dblT < 100.0 {
		dblTemp = 4210.4 - 1.356*dblT + 0.014588*dblT*dblT
	} else if dblT < 200.0 {
		dblTemp = 4306.8 - 2.7913*dblT + 0.018773*dblT*dblT
	} else {
		dblTemp = -999.0
	}

	return dblTemp
}

/*
FNwlam (Function for Water Thermal Conductivity)

この関数は、与えられた温度（`dblT`）における水の熱伝導率を計算します。

建築環境工学的な観点:
- **水の熱伝導率と熱伝達**: 水の熱伝導率は、
  水が関与する熱伝達計算において重要なパラメータです。
  特に、配管内の熱伝達や、蓄熱槽内部の熱伝達をモデル化する際に用いられます。
- **温度依存性**: 水の熱伝導率は温度によって変化するため、
  正確な熱伝達計算には、
  温度に応じた熱伝導率を考慮する必要があります。

この関数は、建物の熱負荷計算、熱搬送システム設計、
およびエネルギー消費量予測を行うための基礎的な役割を果たします。
*/
func FNwlam(dblT float64) float64 {
	var dblTemp float64

	if dblT > 0.0 && dblT < 100.0 {
		dblTemp = 0.56871 + 0.0018421*dblT - 7.0427e-6*dblT*dblT
	} else if dblT < 200.0 {
		dblTemp = 0.60791 + 0.0012032*dblT - 4.7025e-6*dblT*dblT
	} else {
		dblTemp = -999.0
	}

	return dblTemp
}

/*
FNwnew (Function for Water Kinematic Viscosity)

この関数は、与えられた温度（`dblT`）における水の動粘性係数を計算します。

建築環境工学的な観点:
- **水の動粘性と熱伝達**: 水の動粘性係数は、
  水が関与する熱伝達計算において重要なパラメータです。
  特に、配管内の流体の流れや、自然対流による熱伝達率の計算に用いられます。
- **温度依存性**: 水の動粘性係数は温度によって変化するため、
  正確な熱伝達計算には、
  温度に応じた動粘性係数を考慮する必要があります。

この関数は、建物の熱負荷計算、熱搬送システム設計、
およびエネルギー消費量予測を行うための基礎的な役割を果たします。
*/
func FNwnew(dblT float64) float64 {
	var dblTemp float64

	if dblT > 0.0 && dblT < 50.0 {
		dblTemp = math.Exp(-13.233 - 0.032516*dblT + 0.000068997*dblT*dblT +
			0.0000069513*dblT*dblT*dblT - 0.00000009386*dblT*dblT*dblT*dblT)
	} else if dblT < 100.0 {
		dblTemp = math.Exp(-13.618 - 0.015499*dblT - 0.000022461*dblT*dblT +
			0.00000036334*dblT*dblT*dblT)
	} else if dblT < 200.0 {
		dblTemp = math.Exp(-13.698 - 0.016782*dblT + 0.000034425*dblT*dblT)
	} else {
		dblTemp = -999.0
	}

	return dblTemp
}

/*
FNwbeta (Function for Water Thermal Expansion Coefficient)

この関数は、与えられた温度（`dblT`）における水の膨張率を計算します。

建築環境工学的な観点:
- **水の膨張率と自然対流**: 水の膨張率は、
  温度変化による水の密度変化の度合いを示すもので、
  自然対流による熱伝達をモデル化する際に重要なパラメータです。
  特に、グラスホフ数などの無次元数の計算に用いられ、
  自然対流による熱伝達率の計算に不可欠です。
- **温度依存性**: 水の膨張率は温度によって変化するため、
  正確な熱伝達計算には、
  温度に応じた膨張率を考慮する必要があります。

この関数は、建物の熱負荷計算、熱搬送システム設計、
およびエネルギー消費量予測を行うための基礎的な役割を果たします。
*/
func FNwbeta(dblT float64) float64 {
	var dblTemp float64

	if dblT > 0.0 && dblT < 50.0 {
		dblTemp = -0.060159 + 0.018725*dblT - 0.00045278*dblT*dblT +
			0.0000098148*dblT*dblT*dblT - 0.000000083333*dblT*dblT*dblT*dblT
	} else if dblT < 100.0 {
		dblTemp = -0.46048 + 0.03104*dblT - 0.000325*dblT*dblT +
			0.0000013889*dblT*dblT*dblT
	} else if dblT < 200.0 {
		dblTemp = 0.33381 + 0.002847*dblT + 0.000016154*dblT*dblT
	} else {
		dblTemp = -999.0
	}

	return dblTemp
}

/*
FNwmew (Function for Water Viscosity)

この関数は、与えられた温度（`dblT`）における水の粘性係数を計算します。

建築環境工学的な観点:
- **水の粘性と熱伝達**: 水の粘性係数は、
  水が関与する熱伝達計算において重要なパラメータです。
  特に、配管内の流体の流れや、自然対流による熱伝達率の計算に用いられます。
- **温度依存性**: 水の粘性係数は温度によって変化するため、
  正確な熱伝達計算には、
  温度に応じた粘性係数を考慮する必要があります。

この関数は、建物の熱負荷計算、熱搬送システム設計、
およびエネルギー消費量予測を行うための基礎的な役割を果たします。
*/
func FNwmew(dblT float64) float64 {
	return FNwnew(dblT) / FNwrow(dblT)
}

/*
FNaa (Function for Air Thermal Diffusivity)

この関数は、与えられた温度（`dblT`）における空気の熱拡散率を計算します。

建築環境工学的な観点:
- **空気の熱拡散率と熱伝達**: 空気の熱拡散率は、
  熱が空気中をどれだけ速く伝わるかを示すもので、
  熱伝達計算において重要なパラメータです。
  特に、自然対流や強制対流による熱伝達率の計算に用いられます。
- **温度依存性**: 空気の熱拡散率は温度によって変化するため、
  正確な熱伝達計算には、
  温度に応じた熱拡散率を考慮する必要があります。

この関数は、建物の熱負荷計算、換気システム設計、
およびエネルギー消費量予測を行うための基礎的な役割を果たします。
*/
func FNaa(dblT float64) float64 {
	return FNalam(dblT) / FNac() / FNarow(dblT)
}

/*
FNwa (Function for Water Thermal Diffusivity)

この関数は、与えられた温度（`dblT`）における水の熱拡散率を計算します。

建築環境工学的な観点:
- **水の熱拡散率と熱伝達**: 水の熱拡散率は、
  熱が水中をどれだけ速く伝わるかを示すもので、
  熱伝達計算において重要なパラメータです。
  特に、配管内の流体の流れや、自然対流による熱伝達率の計算に用いられます。
- **温度依存性**: 水の熱拡散率は温度によって変化するため、
  正確な熱伝達計算には、
  温度に応じた熱拡散率を考慮する必要があります。

この関数は、建物の熱負荷計算、熱搬送システム設計、
およびエネルギー消費量予測を行うための基礎的な役割を果たします。
*/
func FNwa(dblT float64) float64 {
	return FNwlam(dblT) / FNwc(dblT) / FNwrow(dblT)
}

/*
FNPr (Function for Prandtl Number)

この関数は、与えられた流体（`strF`）と温度（`dblT`）におけるプラントル数を計算します。
プラントル数は、流体の運動量拡散率と熱拡散率の比を示す無次元数であり、
熱伝達計算において重要なパラメータです。

建築環境工学的な観点:
- **プラントル数と熱伝達**: プラントル数は、
  流体中の熱と運動量の伝達メカニズムの相対的な重要性を示します。
  対流熱伝達率の計算に用いられ、
  特に強制対流や自然対流における熱伝達の特性を評価する上で不可欠です。
- **流体種別の考慮**: `strF`によって流体種別（空気`'A'`または水`'W'`）を区別し、
  それぞれの流体に応じた動粘性係数（`FNanew`, `FNwnew`）と熱拡散率（`FNaa`, `FNwa`）を用いて計算します。

この関数は、建物の熱負荷計算、熱搬送システム設計、
およびエネルギー消費量予測を行うための基礎的な役割を果たします。
*/
func FNPr(strF byte, dblT float64) float64 {
	var dblTemp float64

	if strF == 'A' || strF == 'a' {
		dblTemp = FNanew(dblT) / FNaa(dblT)
	} else if strF == 'W' || strF == 'w' {
		dblTemp = FNwnew(dblT) / FNwa(dblT)
	} else {
		dblTemp = -999.0
	}

	return dblTemp
}

/*
FNGr (Function for Grashof Number)

この関数は、与えられた流体（`strF`）、表面温度（`dblTs`）、主流温度（`dblTa`）、
および代表長さ（`dblx`）におけるグラスホフ数を計算します。
グラスホフ数は、自然対流の強さを示す無次元数であり、
熱伝達計算において重要なパラメータです。

建築環境工学的な観点:
- **グラスホフ数と自然対流**: グラスホフ数は、
  流体中の浮力と粘性力の比を示すもので、
  自然対流による熱伝達の強さを評価する際に用いられます。
  特に、壁面や窓からの自然対流による熱伝達率の計算に不可欠です。
- **流体種別の考慮**: `strF`によって流体種別（空気`'A'`または水`'W'`）を区別し、
  それぞれの流体に応じた膨張率（`FNabeta`, `FNwbeta`）と動粘性係数（`FNanew`, `FNwnew`）を用いて計算します。
- **温度差の考慮**: `dbldT`は、
  表面温度と主流温度の差であり、
  自然対流の駆動力となります。

この関数は、建物の熱負荷計算、換気システム設計、
およびエネルギー消費量予測を行うための基礎的な役割を果たします。
*/
func FNGr(strF byte, dblTs, dblTa, dblx float64) float64 {
	const dblg = 9.80665

	dblT := (dblTs + dblTa) / 2.

	// 温度差の計算
	dbldT := math.Max(math.Abs(dblTs-dblTa), 0.1)

	var dblBeta, n float64
	if strF == 'A' || strF == 'a' {
		dblBeta = FNabeta(dblT)
		n = FNanew(dblT)
	} else if strF == 'W' || strF == 'w' {
		dblBeta = FNwbeta(dblT)
		n = FNwnew(dblT)
	} else {
		dblBeta = -999.0
		n = -999.0
	}

	return dblg * dblBeta * dbldT * math.Pow(dblx, 3.0) / (n * n)
}

/*
FNCNu (Function for Nusselt Number Calculation)

この関数は、ヌセルト数（`Nu`）を計算します。
ヌセルト数は、対流熱伝達の効率を示す無次元数であり、
熱伝達率の計算に不可欠なパラメータです。

建築環境工学的な観点:
- **ヌセルト数と熱伝達率**: ヌセルト数は、
  対流熱伝達率と熱伝導率の比を示すもので、
  対流による熱伝達の強さを評価する際に用いられます。
  この関数は、`Nu = C * (Pr Gr)^m` のような経験式を用いてヌセルト数を計算します。
  `dblC`と`dblm`は経験的な定数であり、
  `dblPrGr`はプラントル数とグラスホフ数（またはレイノルズ数）の積です。
- **対流熱伝達のモデル化**: ヌセルト数を計算することで、
  様々な条件下での対流熱伝達率を推定できます。
  これは、壁面や窓からの熱伝達、
  あるいは空調システムにおける熱交換器の性能評価に不可欠です。

この関数は、建物の熱負荷計算、熱搬送システム設計、
およびエネルギー消費量予測を行うための基礎的な役割を果たします。
*/
func FNCNu(dblC, dblm, dblPrGr float64) float64 {
	return dblC * math.Pow(dblPrGr, dblm)
}

/*
FNhinpipe (Function for Internal Convective Heat Transfer Coefficient in Pipe)

この関数は、配管内部を流れる水と管壁間の強制対流熱伝達率を計算します。

建築環境工学的な観点:
- **配管内の熱伝達**: 配管は、
  熱媒（水）を搬送する際に、
  熱損失や熱取得が発生します。
  この関数は、配管内部の流体の流れ（流速`dblv`、内径`dbld`）と、
  流体および管壁の温度（`dblT`）を考慮して、
  熱伝達率を計算します。
- **レイノルズ数とプラントル数**: 
  - `dblRe`: レイノルズ数。流体の流れの様式（層流か乱流か）を示す無次元数であり、
    熱伝達率に大きく影響します。
  - `dblPr`: プラントル数。流体中の熱と運動量の伝達メカニズムの相対的な重要性を示します。
  これらの無次元数を用いて、
  配管内部の強制対流熱伝達率を推定します。
- **層流と乱流の考慮**: `dblRe < 2200.` の条件は、
  流体の流れが層流であるか乱流であるかを判断し、
  それぞれに応じた経験式を適用します。

この関数は、建物の熱搬送システムにおける熱損失・熱取得をモデル化し、
熱負荷計算、エネルギー消費量予測、
および省エネルギー対策の検討を行うための重要な役割を果たします。
*/
func FNhinpipe(dbld, dblL, dblv, dblT float64) float64 {
	dblnew := FNwnew(dblT)
	dblRe := dblv * dbld / dblnew
	dblPr := FNPr('W', dblT)
	dbldL := dbld / dblL
	dblld := FNwlam(dblT) / dbld

	var dblTemp float64
	if dblRe < 2200. {
		dblTemp = (3.66 + 0.0668*dbldL*dblRe*dblPr/(1.+0.04*math.Pow(dbldL*dblRe*dblPr, 2./3.))) * dblld
	} else {
		dblTemp = 0.023 * math.Pow(dblRe, 0.8) * math.Pow(dblPr, 0.4) * dblld
	}

	return dblTemp
}

/*
FNhoutpipe (Function for External Natural Convective Heat Transfer Coefficient for Pipe)

この関数は、円管外部の自然対流熱伝達率を計算します。

建築環境工学的な観点:
- **配管外部の熱伝達**: 配管は、
  周囲の空気と熱を交換します。
  この関数は、管外部の空気の自然対流による熱伝達率を計算します。
- **プラントル数とグラスホフ数**: 
  - `dblPr`: プラントル数。流体中の熱と運動量の伝達メカニズムの相対的な重要性を示します。
  - `dblGr`: グラスホフ数。自然対流の強さを示す無次元数であり、
    熱伝達率に大きく影響します。
  これらの無次元数を用いて、
  配管外部の自然対流熱伝達率を推定します。
- **温度差の考慮**: `dblTs`（表面温度）と`dblTa`（周囲温度）の差は、
  自然対流の駆動力となります。

この関数は、建物の熱搬送システムにおける熱損失・熱取得をモデル化し、
熱負荷計算、エネルギー消費量予測、
および省エネルギー対策の検討を行うための重要な役割を果たします。
*/
func FNhoutpipe(dbld, dblTs, dblTa float64) float64 {
	dblC := 0.5
	dbln := 0.25

	dblPr := FNPr('W', (dblTs+dblTa)/2.)
	dblGr := FNGr('W', dblTs, dblTa, dbld)

	dblNu := dblC * math.Pow(dblPr*dblGr, dbln)

	return dblNu / dbld
}

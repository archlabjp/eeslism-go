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

/*   snbklib.c  */

package eeslism

import (
	"fmt"
	"math"
)

/*
FNFsdw (Function for Shaded Area Ratio of Window)

この関数は、窓面に対する日よけ（庇、袖壁、ルーバーなど）による影面積率を計算します。
これは、日射遮蔽による日射熱取得の抑制効果を定量的に評価するために不可欠です。

建築環境工学的な観点:
- **日射遮蔽の重要性**: 夏季の過度な日射熱取得は、
  冷房負荷の増加や室内温度の上昇を引き起こし、
  エネルギー消費量や快適性に悪影響を与えます。
  日よけは、窓からの日射侵入を効果的に遮蔽し、
  冷房負荷を軽減するパッシブな手法として重要です。
- **影面積率 (Fsdw)**:
  窓の総面積に対する影の面積の割合を示します。
  この値が大きいほど、日射遮蔽効果が高いことを意味します。
- **太陽位置と日よけの形状**: 影の形状や大きさは、
  太陽の高度角や方位角（`Xazm`, `Xprf`）と、
  日よけの形状や寸法（`D`, `Wr`, `Hr`, `Wi1`, `Hi1`, `Wi2`, `Hi2`）によって決まります。
  - `Xazm`: 太陽方位角の正接。
  - `Xprf`: 太陽高度角の正接。
  - `D`: 日よけの奥行き。
  - `Wr`, `Hr`: 窓の幅と高さ。
  - `Wi1`, `Hi1`, `Wi2`, `Hi2`: 袖壁やルーバーの寸法。
- **日よけの種類に応じた計算**: `Ksdw`によって日よけの種類（庇、袖壁、ルーバーなど）を識別し、
  それぞれの種類に応じた影面積の計算関数（`FNAsdw1`, `FNAsdw2`, `FNAsdw3`）を呼び出します。
- **日射熱取得の予測**: この関数で計算される影面積率は、
  窓を透過して室内に侵入する日射熱量を正確に予測するために用いられます。
  これにより、冷房負荷を正確に評価し、
  省エネルギー対策の効果を定量的に把握できます。

この関数は、建物の日射遮蔽計画をモデル化し、
日射熱取得の抑制、冷房負荷の軽減、
昼光利用の最適化、および日影計算を行うための重要な役割を果たします。
*/
func FNFsdw(Ksdw, Ksi int, Xazm, Xprf, D, Wr, Hr, Wi1, Hi1, Wi2, Hi2 float64) float64 {

	if DEBUG {
		fmt.Printf("----- FNFsdw  Ksdw=%d Ksi=%d Xazm=%f Xprf=%f D=%f Wr=%f Hr=%f Wi1=%f Hi1=%f Wi2=%f Hi2=%f\n",
			Ksdw, Ksi, Xazm, Xprf, D, Wr, Hr, Wi1, Hi1, Wi2, Hi2)
	}

	if Ksdw == 0 {
		return 0.0
	}

	Da := D * Xazm
	Dp := D * Xprf
	if Ksdw == 2 || Ksdw == 6 {
		Da = math.Abs(Da)
	}
	if Ksdw == 4 || Ksdw == 8 {
		Da = -Da
	}

	Asdw := 0.0
	switch Ksdw {
	case 1:
		Asdw = FNAsdw1(Da, Dp, Wr, Hr, Wi1, Hi1, Wi2)
	case 2, 3, 4:
		Asdw = FNAsdw1(Dp, Da, Hr, Wr, Hi1, Wi1, Hi2)
	case 5:
		Asdw = FNAsdw2(Dp, Hr, Wr, Hi1)
	case 6, 7, 8:
		Asdw = FNAsdw2(Da, Wr, Hr, Wi1)
	case 9:
		Asdw = FNAsdw3(Da, Dp, Wr, Hr, Wi1, Hi1, Wi2, Hi2)
	}

	Fsdw := Asdw / (Wr * Hr)

	if Ksi == 1 {
		Fsdw = 1.0 - Fsdw
	}

	return Fsdw
}

/*  -----------------------------------------------------  */

/*
FNAsdw1 (Function for Shaded Area Calculation for Type 1 Sunbreak)

この関数は、庇や袖壁などの日よけ（タイプ1）による窓面の影面積を計算します。
これは、日射遮蔽による日射熱取得の抑制効果を定量的に評価するために不可欠です。

建築環境工学的な観点:
- **日よけの影面積**: 庇や袖壁は、窓面の一部に影を落とすことで、
  日射の侵入を抑制します。
  この関数は、日よけの幾何学的パラメータ（`D`, `Wr`, `Hr`, `Wi1`, `Hi`, `Wi2`）と、
  太陽位置（`Da`, `Dp`）に基づいて、影の形状を計算し、その面積を算出します。
  - `Da`: 太陽方位角の正接に日よけの奥行きを乗じた値。
  - `Dp`: 太陽高度角の正接に日よけの奥行きを乗じた値。
  - `Wr`, `Hr`: 窓の幅と高さ。
  - `Wi1`, `Hi`, `Wi2`: 袖壁の幅と高さ、および窓からの距離。
- **影の形状の複雑性**: 影の形状は、太陽の動きや日よけの形状によって複雑に変化します。
  この関数は、`math.Max`, `math.Min`, `math.Abs`などの関数を用いて、
  影の境界線を正確に計算し、その面積を算出します。
- **日射熱取得の予測**: この関数で計算される影面積は、
  窓を透過して室内に侵入する日射熱量を正確に予測するために用いられます。
  これにより、冷房負荷を正確に評価し、
  省エネルギー対策の効果を定量的に把握できます。

この関数は、建物の日射遮蔽計画をモデル化し、
日射熱取得の抑制、冷房負荷の軽減、
昼光利用の最適化、および日影計算を行うための重要な役割を果たします。
*/
func FNAsdw1(Da, Dp, Wr, Hr, Wi1, Hi, Wi2 float64) float64 {
	if Dp <= 0.0 {
		return 0.0
	}

	Wi := Wi1
	if Da < 0.0 {
		Wi = Wi2
	}

	Daa := math.Abs(Da)

	Dha := Wi*Dp/math.Max(Wi, Daa) - Hi
	Dha = math.Min(math.Max(0.0, Dha), Hr)

	Dhb := (Wi+Wr)*Dp/math.Max(Wi+Wr, Daa) - Hi
	Dhb = math.Min(math.Max(0.0, Dhb), Hr)

	var Dwa float64
	if Hi >= Dp {
		Dwa = 0.0
	} else {
		Dwa = (Wi + Wr) - Hi*Daa/Dp
		Dwa = math.Min(math.Max(0.0, Dwa), Wr)
	}

	Dwb := (Wi + Wr) - (Hi+Hr)*Daa/math.Max(Hi+Hr, Dp)
	Dwb = math.Min(math.Max(0.0, Dwb), Wr)

	Asdw := Dwa*Dha + 0.5*(Dwa+Dwb)*(Dhb-Dha)

	return Asdw
}

/*  -----------------------------------------------------  */

/*
FNAsdw2 (Function for Shaded Area Calculation for Type 2 Sunbreak)

この関数は、庇や袖壁などの日よけ（タイプ2）による窓面の影面積を計算します。
これは、日射遮蔽による日射熱取得の抑制効果を定量的に評価するために不可欠です。

建築環境工学的な観点:
- **日よけの影面積**: 庇や袖壁は、窓面の一部に影を落とすことで、
  日射の侵入を抑制します。
  この関数は、日よけの幾何学的パラメータ（`Dp`, `Hr`, `Wr`, `Hi`）と、
  太陽位置（`Dp`）に基づいて、影の形状を計算し、その面積を算出します。
  - `Dp`: 太陽高度角の正接に日よけの奥行きを乗じた値。
  - `Hr`, `Wr`: 窓の高さと幅。
  - `Hi`: 窓の上端から日よけまでの距離。
- **影の形状の簡略化**: この関数は、タイプ1よりも簡略化された影面積の計算を行います。
  これは、特定の日よけ形状や太陽位置の条件下で、
  影の形状が単純になる場合に適用されます。
- **日射熱取得の予測**: この関数で計算される影面積は、
  窓を透過して室内に侵入する日射熱量を正確に予測するために用いられます。
  これにより、冷房負荷を正確に評価し、
  省エネルギー対策の効果を定量的に把握できます。

この関数は、建物の日射遮蔽計画をモデル化し、
日射熱取得の抑制、冷房負荷の軽減、
昼光利用の最適化、および日影計算を行うための重要な役割を果たします。
*/
func FNAsdw2(Dp, Hr, Wr, Hi float64) float64 {
	if Dp <= 0.0 {
		return 0.0
	}

	Dh := math.Min(math.Max(0.0, Dp-Hi), Hr)
	Asdw := Wr * Dh
	return Asdw
}

/*  -----------------------------------------------------  */

/*
FNAsdw3 (Function for Shaded Area Calculation for Type 3 Sunbreak)

この関数は、ルーバーなどの日よけ（タイプ3）による窓面の影面積を計算します。
これは、日射遮蔽による日射熱取得の抑制効果を定量的に評価するために不可欠です。

建築環境工学的な観点:
- **ルーバーの影面積**: ルーバーは、複数の水平または垂直なブレードで構成され、
  日射の侵入を抑制します。
  この関数は、ルーバーの幾何学的パラメータ（`Wr`, `Hr`, `Wi1`, `Hi1`, `Wi2`, `Hi2`）と、
  太陽位置（`Da`, `Dp`）に基づいて、影の形状を計算し、その面積を算出します。
  - `Da`: 太陽方位角の正接に日よけの奥行きを乗じた値。
  - `Dp`: 太陽高度角の正接に日よけの奥行きを乗じた値。
  - `Wr`, `Hr`: 窓の幅と高さ。
  - `Wi1`, `Hi1`, `Wi2`, `Hi2`: ルーバーのブレードの寸法や間隔。
- **影の形状の複雑性**: ルーバーによる影は、
  太陽の動きやブレードの角度によって複雑に変化します。
  この関数は、`math.Min`, `math.Max`などの関数を用いて、
  影の境界線を正確に計算し、その面積を算出します。
- **日射熱取得の予測**: この関数で計算される影面積は、
  窓を透過して室内に侵入する日射熱量を正確に予測するために用いられます。
  これにより、冷房負荷を正確に評価し、
  省エネルギー対策の効果を定量的に把握できます。

この関数は、建物の日射遮蔽計画をモデル化し、
日射熱取得の抑制、冷房負荷の軽減、
昼光利用の最適化、および日影計算を行うための重要な役割を果たします。
*/
func FNAsdw3(Da, Dp, Wr, Hr, Wi1, Hi1, Wi2, Hi2 float64) float64 {
	Dw1 := math.Min(math.Max(0.0, Da-Wi1), Wr)
	Dw2 := math.Min(math.Max(0.0, -Da-Wi2), Wr)
	Dh1 := math.Min(math.Max(0.0, Dp-Hi1), Hr)
	Dh2 := math.Min(math.Max(0.0, -Dp-Hi2), Hr)
	Asdw := Wr*(Dh1+Dh2) + (Dw1+Dw2)*(Hr-Dh1-Dh2)

	return Asdw
}

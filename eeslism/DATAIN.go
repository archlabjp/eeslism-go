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
	"os"
	"strconv"
)

/*

  入力データの読み込み
  FILE=DATAIN.c
  Create Date=1999.6.7
*/

/*
HISASHI (Awning Data Input)

この関数は、庇（ひさし）の仕様データを入力ファイルから読み込み、
対応する構造体（`sunblk`）に格納します。
庇は、窓からの日射侵入を抑制し、冷房負荷を軽減する重要な日射遮蔽部材です。

建築環境工学的な観点:
- **日射遮蔽の幾何学的モデル化**: 庇の形状と位置は、
  日射遮蔽効果に直接影響します。
  この関数は、庇の幾何学的パラメータを読み込みます。
  - `snbname`: 庇の名称。
  - `x`, `y`: 庇の基準点（通常は窓の左下隅など）の相対座標。
  - `D`: 庇の奥行き。窓からの突出寸法。
  - `W`: 庇の幅。窓の幅方向の寸法。
  - `WA`: 庇の傾斜角（方位角）。
  これらのパラメータは、太陽位置と組み合わせて、
  窓面への影の形状と面積を計算するために用いられます。
- **日射熱取得の抑制**: 庇は、特に夏季の太陽高度が高い時間帯に、
  窓からの直達日射を遮蔽することで、
  室内の日射熱取得量を削減し、冷房負荷を軽減します。
  この関数で入力されるデータは、
  この日射遮蔽効果を定量的に評価するための基礎となります。
- **色彩の考慮**: `rgb`は、庇の表面色をRGB値で指定します。
  表面色は、日射吸収率や反射率に影響を与え、
  庇自体の温度上昇や、周囲への放射熱伝達に影響する可能性があります。

この関数は、建物の日射遮蔽計画をモデル化し、
日射熱取得の抑制、冷房負荷の軽減、
昼光利用の最適化、および日影計算を行うための重要なデータ入力機能を提供します。
*/
func HISASHI(fi *EeTokens, sb *sunblk) {
	// 付設障害物名
	sb.snbname = fi.GetToken()

	// 色の初期値
	sb.rgb[0] = 0.0
	sb.rgb[1] = 0.2
	sb.rgb[2] = 0.0

	for !fi.IsEnd() {
		NAME := fi.GetToken()
		if NAME[0] == ';' {
			break
		}

		if NAME == "-xy" {
			sb.x = fi.GetFloat()
			sb.y = fi.GetFloat()
		} else if NAME == "-DW" {
			sb.D = fi.GetFloat()
			sb.W = fi.GetFloat()
		} else if NAME == "-a" {
			sb.WA = fi.GetFloat()
		} else if NAME == "-rgb" {
			// 色
			sb.rgb[0] = fi.GetFloat()
			sb.rgb[1] = fi.GetFloat()
			sb.rgb[2] = fi.GetFloat()
		} else {
			fmt.Printf("ERROR parameter----HISASI: %s\n", NAME)

			os.Exit(1)
		}
	}
}

/*--------------------------------------------------------------*/

/*
BARUKO (Balcony Data Input)

この関数は、バルコニーの仕様データを入力ファイルから読み込み、
対応する構造体（`sunblk`）に格納します。
バルコニーは、庇と同様に窓からの日射侵入を抑制し、
冷房負荷を軽減する日射遮蔽部材として機能します。

建築環境工学的な観点:
- **日射遮蔽の幾何学的モデル化**: バルコニーの形状と位置は、
  日射遮蔽効果に直接影響します。
  この関数は、バルコニーの幾何学的パラメータを読み込みます。
  - `snbname`: バルコニーの名称。
  - `x`, `y`: バルコニーの基準点（通常は窓の左下隅など）の相対座標。
  - `D`: バルコニーの奥行き。窓からの突出寸法。
  - `H`: バルコニーの高さ。手すりなどの高さ。
  - `W`: バルコニーの幅。窓の幅方向の寸法。
  - `h`: バルコニーの床の厚さ。
  これらのパラメータは、太陽位置と組み合わせて、
  窓面への影の形状と面積を計算するために用いられます。
- **日射熱取得の抑制**: バルコニーは、特に夏季の太陽高度が高い時間帯に、
  窓からの直達日射を遮蔽することで、
  室内の日射熱取得量を削減し、冷房負荷を軽減します。
  この関数で入力されるデータは、
  この日射遮蔽効果を定量的に評価するための基礎となります。
- **反射率の考慮**: `ref`は、バルコニー表面の反射率を指定します。
  反射率が高いバルコニーは、日射を反射することで、
  窓面への反射日射を増加させ、日射熱取得を促進する可能性があります。
  そのため、反射率の適切な設定は、日射遮蔽効果を正確に評価する上で重要です。
- **色彩の考慮**: `rgb`は、バルコニーの表面色をRGB値で指定します。
  表面色は、日射吸収率や反射率に影響を与え、
  バルコニー自体の温度上昇や、周囲への放射熱伝達に影響する可能性があります。

この関数は、建物の日射遮蔽計画をモデル化し、
日射熱取得の抑制、冷房負荷の軽減、
昼光利用の最適化、および日影計算を行うための重要なデータ入力機能を提供します。
*/
func BARUKO(fi *EeTokens, sb *sunblk) {
	// 反射率の初期値
	sb.ref = 0.0

	// 色の初期値
	sb.rgb[0] = 0.0
	sb.rgb[1] = 0.2
	sb.rgb[2] = 0.0

	// 付設障害物名
	sb.snbname = fi.GetToken()

	for fi.IsEnd() == false {
		NAME := fi.GetToken()
		if NAME[0] == ';' {
			break
		}

		if NAME == "-xy" {
			sb.x = fi.GetFloat()
			sb.y = fi.GetFloat()
		} else if NAME == "-DHWh" {
			sb.D = fi.GetFloat()
			sb.H = fi.GetFloat()
			sb.W = fi.GetFloat()
			sb.h = fi.GetFloat()
		} else if NAME == "-ref" {
			// 反射率
			sb.ref = fi.GetFloat()
		} else if NAME == "-rgb" {
			// 色
			sb.rgb[0] = fi.GetFloat()
			sb.rgb[1] = fi.GetFloat()
			sb.rgb[2] = fi.GetFloat()
		} else {
			fmt.Printf("ERROR parameter----WBARUKONI: %s\n", NAME)

			os.Exit(1)
		}
	}
}

/*------------------------------------------------------------------*/

/*
SODEK (Side Wall Data Input)

この関数は、袖壁（そでかべ）の仕様データを入力ファイルから読み込み、
対応する構造体（`sunblk`）に格納します。
袖壁は、窓からの日射侵入を抑制し、冷房負荷を軽減する日射遮蔽部材として機能します。

建築環境工学的な観点:
- **日射遮蔽の幾何学的モデル化**: 袖壁の形状と位置は、
  日射遮蔽効果に直接影響します。
  この関数は、袖壁の幾何学的パラメータを読み込みます。
  - `snbname`: 袖壁の名称。
  - `x`, `y`: 袖壁の基準点（通常は窓の左下隅など）の相対座標。
  - `D`: 袖壁の奥行き。窓からの突出寸法。
  - `H`: 袖壁の高さ。窓の高さ方向の寸法。
  - `WA`: 袖壁の傾斜角（方位角）。
  これらのパラメータは、太陽位置と組み合わせて、
  窓面への影の形状と面積を計算するために用いられます。
- **日射熱取得の抑制**: 袖壁は、特に夏季の太陽高度が低い時間帯（朝夕）に、
  窓からの直達日射を遮蔽することで、
  室内の日射熱取得量を削減し、冷房負荷を軽減します。
  この関数で入力されるデータは、
  この日射遮蔽効果を定量的に評価するための基礎となります。
- **色彩の考慮**: `rgb`は、袖壁の表面色をRGB値で指定します。
  表面色は、日射吸収率や反射率に影響を与え、
  袖壁自体の温度上昇や、周囲への放射熱伝達に影響する可能性があります。

この関数は、建物の日射遮蔽計画をモデル化し、
日射熱取得の抑制、冷房負荷の軽減、
昼光利用の最適化、および日影計算を行うための重要なデータ入力機能を提供します。
*/
func SODEK(fi *EeTokens, sb *sunblk) {
	// 色の初期値
	sb.rgb[0] = 0.0
	sb.rgb[1] = 0.2
	sb.rgb[2] = 0.0

	// 付設障害物名
	sb.snbname = fi.GetToken()

	for !fi.IsEnd() {
		NAME := fi.GetToken()
		if NAME[0] == ';' {
			break
		}

		if NAME == "-xy" {
			sb.x = fi.GetFloat()
			sb.y = fi.GetFloat()
		} else if NAME == "-DH" {
			sb.D = fi.GetFloat()
			sb.H = fi.GetFloat()
		} else if NAME == "-a" {
			sb.WA = fi.GetFloat()
		} else if NAME == "-rgb" {
			// 色
			sb.rgb[0] = fi.GetFloat()
			sb.rgb[1] = fi.GetFloat()
			sb.rgb[2] = fi.GetFloat()
		} else {
			fmt.Printf("ERROR parameter----SODEKABE: %s\n", NAME)

			os.Exit(1)
		}
	}
}

/*-----------------------------------------------------------------------*/

/*
SCREEN (Window Shade Data Input)

この関数は、窓日よけ（ブラインド、ルーバーなど）の仕様データを入力ファイルから読み込み、
対応する構造体（`sunblk`）に格納します。
窓日よけは、窓からの日射侵入を抑制し、冷房負荷を軽減する日射遮蔽部材として機能します。

建築環境工学的な観点:
- **日射遮蔽の幾何学的モデル化**: 窓日よけの形状と位置は、
  日射遮蔽効果に直接影響します。
  この関数は、窓日よけの幾何学的パラメータを読み込みます。
  - `snbname`: 窓日よけの名称。
  - `x`, `y`: 窓日よけの基準点（通常は窓の左下隅など）の相対座標。
  - `D`: 窓日よけの奥行き。窓からの突出寸法。
  - `H`: 窓日よけの高さ。窓の高さ方向の寸法。
  - `W`: 窓日よけの幅。窓の幅方向の寸法。
  これらのパラメータは、太陽位置と組み合わせて、
  窓面への影の形状と面積を計算するために用いられます。
- **日射熱取得の抑制**: 窓日よけは、特に夏季の太陽高度が高い時間帯に、
  窓からの直達日射を遮蔽することで、
  室内の日射熱取得量を削減し、冷房負荷を軽減します。
  この関数で入力されるデータは、
  この日射遮蔽効果を定量的に評価するための基礎となります。
- **色彩の考慮**: `rgb`は、窓日よけの表面色をRGB値で指定します。
  表面色は、日射吸収率や反射率に影響を与え、
  窓日よけ自体の温度上昇や、周囲への放射熱伝達に影響する可能性があります。

この関数は、建物の日射遮蔽計画をモデル化し、
日射熱取得の抑制、冷房負荷の軽減、
昼光利用の最適化、および日影計算を行うための重要なデータ入力機能を提供します。
*/
func SCREEN(fi *EeTokens, sb *sunblk) {
	sb.rgb[0] = 0.0
	sb.rgb[1] = 0.2
	sb.rgb[2] = 0.0

	sb.snbname = fi.GetToken()

	for !fi.IsEnd() {
		NAME := fi.GetToken()
		if NAME[0] == ';' {
			break
		}

		if NAME == "-xy" {
			sb.x = fi.GetFloat()
			sb.y = fi.GetFloat()
		} else if NAME == "-DHW" {
			sb.D = fi.GetFloat()
			sb.H = fi.GetFloat()
			sb.W = fi.GetFloat()
		} else if NAME == "-rgb" {
			sb.rgb[0] = fi.GetFloat()
			sb.rgb[1] = fi.GetFloat()
			sb.rgb[2] = fi.GetFloat()
		} else {
			fmt.Printf("ERROR paramater---MADOHIYOKE: %s\n", NAME)

			os.Exit(1)
		}
	}
}

/*----------------------------------------------------------------*/
/*
rmpdata (Room Main Plane Data Input)

この関数は、室の主面（壁、床、天井など）の仕様データを入力ファイルから読み込み、
対応する構造体（`RRMP`）に格納します。
これには、主面の位置、寸法、反射率、および窓に関する情報が含まれます。

建築環境工学的な観点:
- **室の幾何学的モデル化**: 室の熱負荷計算や昼光利用計算では、
  室を構成する各面（壁、床、天井など）の幾何学的情報が不可欠です。
  この関数は、主面の位置（`xb0`, `yb0`）、寸法（`Rw`, `Rh`）、
  そして窓の位置や寸法を読み込み、室の形状を正確にモデル化します。
- **反射率の考慮**: `ref`は、主面の反射率を指定します。
  表面の反射率は、日射熱取得や昼光利用に影響を与えます。
  例えば、反射率の高い壁面は、日射を室内に反射することで、
  日射熱取得を促進したり、昼光利用を向上させたりする可能性があります。
- **窓のモデル化**: `WD`は、主面内に含まれる窓の情報を格納しており、
  窓の名称（`winname`）、相対位置（`xr`, `yr`）、寸法（`Ww`, `Wh`）、
  反射率（`ref`）、色（`rgb`）などを読み込みます。
  これにより、窓からの日射熱取得や昼光利用を詳細にモデル化できます。
- **前面地面の代表点までの距離 (grpx)**:
  `grpx`は、主面から前面地面の代表点までの距離を示します。
  これは、地面からの反射日射や、地盤からの熱伝達を考慮する際に用いられる可能性があります。

この関数は、室の幾何学的モデル化と熱的特性の定義を行い、
熱負荷計算、エネルギー消費量予測、
昼光利用の最適化、および快適性評価を行うための重要なデータ入力機能を提供します。
*/
func rmpdata(fi *EeTokens) *RRMP {
	rp := RRMPInit()

	// RMP名
	rp.rmpname = fi.GetToken()

	// 壁名称
	rp.wallname = fi.GetToken()

	for !fi.IsEnd() {
		NAME := fi.GetToken()
		if NAME[0] == ';' {
			fi.SkipToEndOfLine()
			break
		}

		if NAME == "-xyb" {
			// 左下頂点座標
			rp.xb0 = fi.GetFloat()
			rp.yb0 = fi.GetFloat()
		} else if NAME == "-WH" {
			// 巾、高さ
			rp.Rw = fi.GetFloat()
			rp.Rh = fi.GetFloat()
		} else if NAME == "-ref" {
			// 反射率
			rp.ref = fi.GetFloat()
		} else if NAME == "-grpx" {
			// 前面地面の代表点までの距離
			rp.grpx = fi.GetFloat()
		} else if NAME == "-rgb" {
			// 色
			rp.rgb[0] = fi.GetFloat()
			rp.rgb[1] = fi.GetFloat()
			rp.rgb[2] = fi.GetFloat()
		} else {
			fmt.Printf("ERROR parameter----RMP: %s\n", NAME)
			os.Exit(1)
		}
	}

	// ex: `WD  window  -xyr 1.325 1.05 -WH 3.3 1.05`
	rp.WD = make([]*MADO, 0)

	pos := fi.GetPos()

	for !fi.IsEnd() {
		line := new(EeTokens)
		line.tokens = fi.GetLogicalLine()
		line.pos = 0

		// 空行の場合は終了
		if len(line.tokens) == 1 && line.tokens[0] == ";" {
			break
		}

		NAME := line.GetToken()

		if NAME == "WD" {
			wp := MADOInit()

			wp.winname = line.GetToken()

			for !line.IsEnd() {
				NAME := line.GetToken()
				if NAME[0] == ';' {
					break
				}

				if NAME == "-xyr" {
					// 左下頂点座標
					wp.xr = line.GetFloat()
					wp.yr = line.GetFloat()
				} else if NAME == "-WH" {
					// 巾、高さ
					wp.Ww = line.GetFloat()
					wp.Wh = line.GetFloat()
				} else if NAME == "-ref" {
					// 反射率
					wp.ref = line.GetFloat()
				} else if NAME == "-grpx" {
					// 前面地面の代表点までの距離
					wp.grpx = line.GetFloat()
				} else if NAME == "-rgb" {
					// 色
					wp.rgb[0] = line.GetFloat()
					wp.rgb[1] = line.GetFloat()
					wp.rgb[2] = line.GetFloat()
				} else {
					fmt.Printf("ERROR parameter----WD: %s\n", NAME)
					os.Exit(1)
				}
			}

			rp.WD = append(rp.WD, wp)
		} else if NAME == "RMP" {
			fi.RestorePos(pos)
			break
		}
	}

	return rp
}

/*
MADOInit (Window Initialization)

この関数は、新しい窓のデータ構造を初期化します。

建築環境工学的な観点:
- **窓の初期化**: 窓のシミュレーションを行う前に、
  その幾何学的パラメータ（幅`Ww`、高さ`Wh`、相対位置`xr`, `yr`など）や、
  熱的・光学的特性（反射率`ref`、色`rgb`など）をデフォルト値で初期化します。
  これにより、後続のデータ入力や計算が正しく行われるための準備が整います。
- **日射熱取得と昼光利用の基礎**: これらのパラメータは、
  窓からの日射熱取得量や昼光の取り込み量を計算するための基礎となります。
  正確な初期化は、窓の設計が建物のエネルギー消費量や室内快適性に与える影響を
  シミュレーションする上で重要です。

この関数は、建物の窓の熱的・光学的挙動をモデル化し、
熱負荷計算、エネルギー消費量予測、
昼光利用の最適化、および快適性評価を行うための重要な役割を果たします。
*/
func MADOInit() *MADO {
	wp := new(MADO)
	wp.winname = ""
	matinit(wp.rgb[:], 3)
	wp.grpx = 1.0 // 前面地面の代表点までの距離 = 1
	wp.ref = 0.0
	wp.Wh = 0.0
	wp.xr = 0.0
	wp.yr = 0.0

	wp.ref = 0.0 // 反射率 = 0

	wp.rgb[0] = 0.0
	wp.rgb[1] = 0.3
	wp.rgb[2] = 0.8

	return wp
}

/*------------------------------------------------------------------*/

/*
rectdata (Rectangular Obstacle Data Input)

この関数は、長方形（平面）の外部障害物の仕様データを入力ファイルから読み込み、
対応する構造体（`OBS`）に格納します。
これらのデータは、日影計算や日射量計算において、
周囲の建物や地形による日影の影響を評価するために不可欠です。

建築環境工学的な観点:
- **外部障害物のモデル化**: 建物の周囲に存在する他の建物や構造物、
  あるいは地形などは、日影を形成し、
  対象建物への日射入射量に影響を与えます。
  この関数は、長方形の障害物の幾何学的パラメータを読み込みます。
  - `obsname`: 障害物の名称。
  - `x`, `y`, `z`: 障害物の基準点（通常は左下隅）の座標。
  - `W`, `H`: 障害物の幅と高さ。
  - `Wa`, `Wb`: 障害物の方位角と傾斜角。
  これらのパラメータは、太陽位置と組み合わせて、
  対象建物への影の形状と面積を計算するために用いられます。
- **日影計算の基礎**: この関数で入力されるデータは、
  対象建物への日影の影響を定量的に評価するための基礎となります。
  これにより、日射熱取得の予測精度を向上させ、
  冷房負荷を正確に評価できます。
- **反射率と色彩の考慮**: `ref`は、障害物表面の反射率を指定します。
  反射率の高い障害物は、日射を反射することで、
  対象建物への反射日射を増加させる可能性があります。
  `rgb`は、障害物の表面色をRGB値で指定します。

この関数は、建物の日射環境を正確にモデル化し、
日射熱取得の抑制、冷房負荷の軽減、
昼光利用の最適化、および日影計算を行うための重要なデータ入力機能を提供します。
*/
func rectdata(fi *EeTokens, obs *OBS) {
	obs.ref[0] = 0.0

	obs.rgb[0] = 0.7
	obs.rgb[1] = 0.7
	obs.rgb[2] = 0.7

	// 名前
	obs.obsname = fi.GetToken()

	for !fi.IsEnd() {
		NAME := fi.GetToken()
		if NAME[0] == ';' {
			break
		}

		if NAME == "-xyz" {
			// 左下頂点座標
			obs.x = fi.GetFloat()
			obs.y = fi.GetFloat()
			obs.z = fi.GetFloat()
		} else if NAME == "-WH" {
			// 巾、高さ
			obs.W = fi.GetFloat()
			obs.H = fi.GetFloat()
		} else if NAME == "-WaWb" {
			// 方位角、傾斜角
			obs.Wa = fi.GetFloat()
			obs.Wb = fi.GetFloat()
		} else if NAME == "-ref" {
			// 反射率
			obs.ref[0] = fi.GetFloat()
		} else if NAME == "-rgb" {
			// 色
			obs.rgb[0] = fi.GetFloat()
			obs.rgb[1] = fi.GetFloat()
			obs.rgb[2] = fi.GetFloat()
		} else {
			fmt.Printf("ERROR parameter----OBS.rect: %s\n", NAME)
			os.Exit(1)
		}
	}
}

/*------------------------------------------------------------------*/

/*
cubdata (Cubic Obstacle Data Input)

この関数は、直方体（立方体）の外部障害物の仕様データを入力ファイルから読み込み、
対応する構造体（`OBS`）に格納します。
これらのデータは、日影計算や日射量計算において、
周囲の建物や地形による日影の影響を評価するために不可欠です。

建築環境工学的な観点:
- **外部障害物のモデル化**: 建物の周囲に存在する他の建物や構造物、
  あるいは地形などは、日影を形成し、
  対象建物への日射入射量に影響を与えます。
  この関数は、直方体の障害物の幾何学的パラメータを読み込みます。
  - `obsname`: 障害物の名称。
  - `x`, `y`, `z`: 障害物の基準点（通常は左下隅）の座標。
  - `W`, `D`, `H`: 障害物の幅、奥行き、高さ。
  - `Wa`: 障害物の方位角。
  これらのパラメータは、太陽位置と組み合わせて、
  対象建物への影の形状と面積を計算するために用いられます。
- **日影計算の基礎**: この関数で入力されるデータは、
  対象建物への日影の影響を定量的に評価するための基礎となります。
  これにより、日射熱取得の予測精度を向上させ、
  冷房負荷を正確に評価できます。
- **反射率と色彩の考慮**: `ref0`, `ref1`, `ref2`, `ref3`は、
  障害物の各面（通常は東西南北の4面）の反射率を指定します。
  反射率の高い障害物は、日射を反射することで、
  対象建物への反射日射を増加させる可能性があります。
  `rgb`は、障害物の表面色をRGB値で指定します。

この関数は、建物の日射環境を正確にモデル化し、
日射熱取得の抑制、冷房負荷の軽減、
昼光利用の最適化、および日影計算を行うための重要なデータ入力機能を提供します。
*/
func cubdata(fi *EeTokens, obs *OBS) {
	for i := 0; i < 3; i++ {
		obs.ref[i] = 0.0
	}

	obs.rgb[0] = 0.7
	obs.rgb[1] = 0.7
	obs.rgb[2] = 0.7

	// 名前
	obs.obsname = fi.GetToken()

	for fi.IsEnd() == false {
		NAME := fi.GetToken()
		if NAME[0] == ';' {
			break
		}

		if NAME == "-xyz" {
			// 左下頂点座標
			obs.x = fi.GetFloat()
			obs.y = fi.GetFloat()
			obs.z = fi.GetFloat()
		} else if NAME == "-WDH" {
			// 巾、奥行き、高さ
			obs.W = fi.GetFloat()
			obs.D = fi.GetFloat()
			obs.H = fi.GetFloat()
		} else if NAME == "-Wa" {
			// 方位角
			obs.Wa = fi.GetFloat()
		} else if NAME == "-ref0" {
			// 反射率
			obs.ref[0] = fi.GetFloat()
		} else if NAME == "-ref1" {
			// 反射率
			obs.ref[1] = fi.GetFloat()
		} else if NAME == "-ref2" {
			// 反射率
			obs.ref[2] = fi.GetFloat()
		} else if NAME == "-ref3" {
			// 反射率
			obs.ref[3] = fi.GetFloat()
		} else if NAME == "-rgb" {
			// 色
			obs.rgb[0] = fi.GetFloat()
			obs.rgb[1] = fi.GetFloat()
			obs.rgb[2] = fi.GetFloat()
		} else {
			fmt.Printf("ERROR parameter----OBS.cube: %s\n", NAME)
			os.Exit(1)
		}
	}
}

/*-------------------------------------------------------------------*/

/*
tridata (Triangular Obstacle Data Input)

この関数は、三角形の外部障害物の仕様データを入力ファイルから読み込み、
対応する構造体（`OBS`）に格納します。
これらのデータは、日影計算や日射量計算において、
周囲の建物や地形による日影の影響を評価するために不可欠です。

建築環境工学的な観点:
- **外部障害物のモデル化**: 建物の周囲に存在する他の建物や構造物、
  あるいは地形などは、日影を形成し、
  対象建物への日射入射量に影響を与えます。
  この関数は、三角形の障害物の幾何学的パラメータを読み込みます。
  - `obsname`: 障害物の名称。
  - `x`, `y`, `z`: 障害物の基準点（通常は左下隅）の座標。
  - `W`, `H`: 障害物の幅と高さ。
  - `Wa`, `Wb`: 障害物の方位角と傾斜角。
  これらのパラメータは、太陽位置と組み合わせて、
  対象建物への影の形状と面積を計算するために用いられます。
- **日影計算の基礎**: この関数で入力されるデータは、
  対象建物への日影の影響を定量的に評価するための基礎となります。
  これにより、日射熱取得の予測精度を向上させ、
  冷房負荷を正確に評価できます。
- **反射率と色彩の考慮**: `ref`は、障害物表面の反射率を指定します。
  反射率の高い障害物は、日射を反射することで、
  対象建物への反射日射を増加させる可能性があります。
  `rgb`は、障害物の表面色をRGB値で指定します。

この関数は、建物の日射環境を正確にモデル化し、
日射熱取得の抑制、冷房負荷の軽減、
昼光利用の最適化、および日影計算を行うための重要なデータ入力機能を提供します。
*/
func tridata(fi *EeTokens, obs *OBS) {
	obs.ref[0] = 0.0

	obs.rgb[0] = 0.7
	obs.rgb[1] = 0.7
	obs.rgb[2] = 0.7

	// 名前
	obs.obsname = fi.GetToken()

	for !fi.IsEnd() {
		NAME := fi.GetToken()
		if NAME[0] == ';' {
			break
		}

		if NAME == "-xyz" {
			// 左下頂点座標
			obs.x = fi.GetFloat()
			obs.y = fi.GetFloat()
			obs.z = fi.GetFloat()
		} else if NAME == "-WH" {
			// 巾、高さ
			obs.W = fi.GetFloat()
			obs.H = fi.GetFloat()
		} else if NAME == "-WaWb" {
			// 方位角、傾斜角
			obs.Wa = fi.GetFloat()
			obs.Wb = fi.GetFloat()
		} else if NAME == "-ref" {
			// 反射率
			obs.ref[0] = fi.GetFloat()
		} else if NAME == "-rgb" {
			// 色
			obs.rgb[0] = fi.GetFloat()
			obs.rgb[1] = fi.GetFloat()
			obs.rgb[2] = fi.GetFloat()
		} else {
			fmt.Printf("ERROR parameter----OBS.triangle: %s\n", NAME)
			os.Exit(1)
		}
	}
}

/*-------------------------------------------------------------------*/
/*
dividdata (Division Data Input)

この関数は、壁の分割数やモンテカルロ法の計算回数など、
シミュレーションの精度や計算方法に関するデータを入力ファイルから読み込みます。

建築環境工学的な観点:
- **壁の分割数 (DE)**:
  壁体内部の熱伝導計算において、壁を仮想的に分割する細かさを示します。
  分割数を増やすことで、壁体内部の温度分布や熱流をより正確にモデル化できますが、
  計算時間も増加します。
  適切な分割数の設定は、シミュレーションの精度と計算効率のバランスを考慮する必要があります。
- **モンテカルロ法の計算回数 (monten)**:
  モンテカルロ法は、乱数を用いてシミュレーションを行う手法であり、
  特に複雑な日影計算や放射熱伝達計算に用いられます。
  計算回数を増やすことで、計算結果の精度が向上しますが、
  計算時間も増加します。
  適切な計算回数の設定は、シミュレーションの精度と計算効率のバランスを考慮する必要があります。
- **シミュレーションの精度と計算効率**: この関数で入力されるデータは、
  シミュレーションの精度と計算効率に直接影響します。
  適切な設定を行うことで、
  限られた計算資源の中で、より信頼性の高いシミュレーション結果を得ることができます。

この関数は、建物のエネルギーシミュレーションの精度と計算効率を制御し、
より信頼性の高いシミュレーション結果を得るための重要なデータ入力機能を提供します。
*/
func dividdata(fi *EeTokens, monten *int, DE *float64) {
	var NAME string

	for !fi.IsEnd() {
		NAME = fi.GetToken()
		if NAME[0] == ';' {
			break
		}

		if NAME == "DE" {
			// DE: 壁の分割[mm]
			var err error
			s := fi.GetToken()
			*DE, err = strconv.ParseFloat(s, 64)
			if err != nil {
				fmt.Printf("ERROR parameter----DIVID: %s\n", NAME)
			}
		} else if NAME == "MONT" {
			// monte: モンテカルロ法の計算回数
			var err error
			s := fi.GetToken()
			*monten, err = strconv.Atoi(s)
			if err != nil {
				fmt.Printf("ERROR parameter----DIVID: %s\n", NAME)
			}
		} else {
			fmt.Printf("ERROR parameter----DIVID: %s\n", NAME)

			os.Exit(1)
		}
	}
}

/*
treedata (Tree Data Input)

この関数は、樹木の仕様データを入力ファイルから読み込み、
対応する構造体（`TREE`）に格納します。
樹木は、日影を形成し、日射熱取得に影響を与える外部障害物としてモデル化されます。

建築環境工学的な観点:
- **樹木による日影のモデル化**: 樹木は、その形状や葉の密度によって、
  日影を形成し、対象建物への日射入射量に影響を与えます。
  この関数は、樹木の幾何学的パラメータを読み込みます。
  - `treetype`: 樹木の形状タイプ（例: `treeA`）。
  - `treename`: 樹木の名称。
  - `x`, `y`, `z`: 樹木の基準点（通常は幹の根元）の座標。
  - `W1`, `H1`: 幹の太さと高さ。
  - `W2`, `H2`: 葉部下面の幅と高さ。
  - `W3`, `H3`: 葉部中央の幅と高さ。
  - `W4`: 葉部上面の幅。
  これらのパラメータは、太陽位置と組み合わせて、
  対象建物への影の形状と面積を計算するために用いられます。
- **日影計算の基礎**: この関数で入力されるデータは、
  対象建物への日影の影響を定量的に評価するための基礎となります。
  これにより、日射熱取得の予測精度を向上させ、
  冷房負荷を正確に評価できます。
- **色彩の考慮**: `rgb`は、樹木の表面色をRGB値で指定します。
  表面色は、日射吸収率や反射率に影響を与え、
  樹木自体の温度上昇や、周囲への放射熱伝達に影響する可能性があります。

この関数は、建物の日射環境を正確にモデル化し、
日射熱取得の抑制、冷房負荷の軽減、
昼光利用の最適化、および日影計算を行うための重要なデータ入力機能を提供します。
*/
func treedata(fi *EeTokens, tree *[]*TREE) {
	var tred *TREE

	*tree = make([]*TREE, 0)

	for !fi.IsEnd() {
		tred = treeinit()

		line := new(EeTokens)
		line.tokens = fi.GetLogicalLine()
		line.pos = 0

		var NAME string
		NAME = line.GetToken()
		if NAME[0] == '*' {
			break
		}

		// 樹木の形
		tred.treetype = NAME

		// 名前
		tred.treename = line.GetToken()

		if tred.treetype == "treeA" {
			for !line.IsEnd() {
				NAME = line.GetToken()
				if NAME[0] == ';' {
					break
				}

				if NAME == "-xyz" {
					// 幹部下面の中心座標
					tred.x = line.GetFloat()
					tred.y = line.GetFloat()
					tred.z = line.GetFloat()
				} else if NAME == "-WH1" {
					// 幹太さ、幹高さ
					tred.W1 = line.GetFloat()
					tred.H1 = line.GetFloat()
				} else if NAME == "-WH2" {
					// 葉部下面巾、葉部下側高さ
					tred.W2 = line.GetFloat()
					tred.H2 = line.GetFloat()
				} else if NAME == "-WH3" {
					// 葉部中央巾、葉部上側高さ
					tred.W3 = line.GetFloat()
					tred.H3 = line.GetFloat()
				} else if NAME == "-WH4" || NAME == "-W4" {
					// 葉部上面巾
					tred.W4 = line.GetFloat()
				} else {
					fmt.Printf("ERROR parameter----TREE: %s %s\n", tred.treename, NAME)
					os.Exit(1)
				}
			}
		} else {
			fmt.Printf("ERROR parameter----TREE: %s\n", tred.treetype)
			os.Exit(1)
		}

		*tree = append(*tree, tred)
	}
}

/*
treeinit (Tree Initialization)

この関数は、新しい樹木のデータ構造を初期化します。

建築環境工学的な観点:
- **樹木の初期化**: 樹木のシミュレーションを行う前に、
  その幾何学的パラメータ（幹の太さ`W1`、高さ`H1`、葉部の幅`W2`, `W3`, `W4`、高さ`H2`, `H3`など）を
  デフォルト値（通常はゼロ）で初期化します。
  これにより、後続のデータ入力や計算が正しく行われるための準備が整います。
- **日影計算の基礎**: これらのパラメータは、
  樹木によって形成される日影の形状や範囲を計算するための基礎となります。
  正確な初期化は、日射遮蔽効果の評価や、
  日影規制の遵守をシミュレーションする上で重要です。

この関数は、建物の日射環境をモデル化し、
日射熱取得の抑制、冷房負荷の軽減、
昼光利用の最適化、および日影計算を行うための重要な役割を果たします。
*/
func treeinit() *TREE {
	tred := new(TREE)

	tred.treename = ""
	tred.treetype = ""
	tred.x = 0.0
	tred.y = 0.0
	tred.z = 0.0
	tred.W1 = 0.0
	tred.W2 = 0.0
	tred.W3 = 0.0
	tred.W4 = 0.0
	tred.H1 = 0.0
	tred.H2 = 0.0
	tred.H3 = 0.0

	return tred
}

/*-------------------------*/

// `POLYGON
//   <polyknd> <polyd> <polyname> <wallname> -xyz [<x> <y> <z>]+ -rgb <r> <g> <b> -ref <ref> -refg <refg> -grpx <grpx> ;`
//   <polyknd> <polyd> <polyname> <wallname> -xyz [<x> <y> <z>]+ -rgb <r> <g> <b> -ref <ref> -refg <refg> -grpx <grpx> ;`
//  *`
/*
polydata (Polygon Data Input)

この関数は、多角形データで定義された障害物や受光面の仕様を入力ファイルから読み込み、
対応する構造体（`POLYGN`）に格納します。
これにより、複雑な形状の建物や周囲の環境を柔軟にモデル化できます。

建築環境工学的な観点:
- **複雑な形状のモデル化**: 建物や周囲の環境は、
  単純な直方体や長方形では表現できない複雑な形状を持つことがあります。
  この関数は、多角形の頂点座標（`P`）を読み込むことで、
  これらの複雑な形状を正確にモデル化し、
  日影計算や日射量計算の精度を向上させます。
- **障害物と受光面の区別**: `polyknd`によって、
  多角形が障害物（`OBS`）として機能するか、
  受光面（`RMP`）として機能するかを区別します。
  これにより、それぞれの役割に応じた計算ロジックが適用されます。
- **反射率と色彩の考慮**: `ref`は、多角形表面の反射率を指定します。
  `refg`は、前面地面の反射率を指定します。
  反射率の高い表面は、日射を反射することで、
  周囲の環境や対象建物への日射熱取得に影響を与えます。
  `rgb`は、多角形の表面色をRGB値で指定します。

この関数は、建物の日射環境を正確にモデル化し、
日射熱取得の抑制、冷房負荷の軽減、
昼光利用の最適化、および日影計算を行うための重要なデータ入力機能を提供します。
*/
func polydata(fi *EeTokens, poly *[]*POLYGN) {
	var i int
	var Npoly int

	*poly = make([]*POLYGN, Npoly)

	for !fi.IsEnd() {
		polyp := polyinit()

		line := new(EeTokens)
		line.tokens = fi.GetLogicalLine()
		line.pos = 0

		var NAME string
		NAME = line.GetToken()
		if NAME[0] == '*' {
			break
		}

		// 前面地面の代表点までの距離の初期値
		polyp.grpx = 1.0

		// 色の初期値
		polyp.rgb[0] = 0.9
		polyp.rgb[1] = 0.9
		polyp.rgb[2] = 0.9

		// ポリゴン種類
		polyp.polyknd = NAME

		if polyp.polyknd != "RMP" && polyp.polyknd != "OBS" {
			fmt.Printf("ERROR parameter----POLYGON: %s  <RMP> or <OBS> \n", polyp.polyknd)
			os.Exit(1)
		}

		// 頂点数
		polyp.polyd = line.GetInt()
		polyp.P = make([]XYZ, polyp.polyd)

		// 名前
		polyp.polyname = line.GetToken()

		// 壁名
		polyp.wallname = line.GetToken()

		for !line.IsEnd() {
			NAME = line.GetToken()
			if NAME[0] == ';' {
				break
			}

			if NAME == "-xyz" {
				// 頂点座標
				for i = 0; i < polyp.polyd; i++ {
					polyp.P[i].X = line.GetFloat()
					polyp.P[i].Y = line.GetFloat()
					polyp.P[i].Z = line.GetFloat()
				}

			} else if NAME == "-rgb" {
				// 色
				polyp.rgb[0] = line.GetFloat()
				polyp.rgb[1] = line.GetFloat()
				polyp.rgb[2] = line.GetFloat()
			} else if NAME == "-ref" {
				// 反射率
				polyp.ref = line.GetFloat()
			} else if NAME == "-refg" {
				// 前面地面の反射率
				polyp.refg = line.GetFloat()
			} else if NAME == "-grpx" {
				// 前面地面の代表点までの距離
				polyp.grpx = line.GetFloat()
			} else {
				fmt.Printf("ERROR parameter----POLYGON: %s\n", NAME)
				os.Exit(1)
			}
		}

		*poly = append(*poly, polyp)
	}
}

/*
polyinit (Polygon Initialization)

この関数は、新しい多角形データ構造を初期化します。

建築環境工学的な観点:
- **多角形の初期化**: 多角形のシミュレーションを行う前に、
  その幾何学的パラメータ（頂点座標`P`、頂点数`polyd`など）や、
  熱的・光学的特性（反射率`ref`、色`rgb`など）をデフォルト値で初期化します。
  これにより、後続のデータ入力や計算が正しく行われるための準備が整います。
- **日影計算と日射量計算の基礎**: これらのパラメータは、
  多角形によって形成される日影の形状や範囲、
  および多角形表面への日射入射量を計算するための基礎となります。
  正確な初期化は、建物のエネルギー消費量や室内快適性に与える影響を
  シミュレーションする上で重要です。

この関数は、建物の日射環境をモデル化し、
日射熱取得の抑制、冷房負荷の軽減、
昼光利用の最適化、および日影計算を行うための重要な役割を果たします。
*/
func polyinit() *POLYGN {
	polyp := new(POLYGN)

	// ポリゴン種類
	polyp.polyknd = ""

	// 名前
	polyp.polyname = ""

	// 壁名
	polyp.wallname = ""

	// 反射率
	polyp.ref = 0.0

	// 前面地面の反射率
	polyp.refg = 0.0

	// 前面地面の代表点までの距離 = 1
	polyp.grpx = 1.0

	// 頂点
	polyp.polyd = 0
	polyp.P = nil

	// 色
	matinit(polyp.rgb[:], 3)

	return polyp
}

/*---------------------------------------------------------------------------*/
/*
bdpdata (Building Data Point Data Input)

この関数は、建物の基準点（Building Data Point, BDP）のデータと、
それに付随する日よけ（SBLK）や室の主面（RMP）のデータを入力ファイルから読み込み、
対応する構造体に格納します。
これは、建物全体の幾何学的モデルを構築し、
日影計算や日射量計算を行うための基礎となります。

建築環境工学的な観点:
- **建物全体の幾何学的モデル化**: BDPは、
  建物の基準となる位置（`x0`, `y0`, `z0`）と、
  建物の向き（方位角`Wa`、傾斜角`Wb`）を定義します。
  これにより、建物全体を3次元空間に配置し、
  太陽位置との相対関係を正確にモデル化できます。
- **日よけと室の主面の関連付け**: BDPに付随するSBLK（日よけ）やRMP（室の主面）は、
  そのBDPの座標系に基づいて定義されます。
  これにより、建物全体の日射遮蔽や日射熱取得を統合的に評価できます。
- **外部日射面との連携**: `exsfname`は、
  外部日射面（`EXSF`）の名称を指定し、
  その方位角と傾斜角をBDPに設定します。
  これにより、建物が受ける外部日射の影響を正確にモデル化できます。
- **データ入力の階層構造**: この関数は、
  BDP、SBLK、RMPという階層的なデータ構造を読み込みます。
  これにより、建物の複雑な幾何学的情報を効率的に管理し、
  シミュレーションモデルに組み込むことができます。

この関数は、建物の幾何学的モデルを構築し、
日射熱取得の抑制、冷房負荷の軽減、
昼光利用の最適化、および日影計算を行うための重要なデータ入力機能を提供します。
*/
func bdpdata(fi *EeTokens, bp *[]*BBDP, Exsf *EXSFS) {

	var sb *sunblk

	// BDPデータの読み込み
	for !fi.IsEnd() {
		bbdp := bdpinit()

		var NAME string
		NAME = fi.GetToken()
		if NAME[0] == '*' {
			break
		}

		if NAME != "BDP" {
			fmt.Printf("error BDP\n")
			os.Exit(1)
		}

		bbdp.bdpname = fi.GetToken()

		for !fi.IsEnd() {
			NAME = fi.GetToken()
			if NAME[0] == ';' {
				fi.SkipToEndOfLine()
				break
			}

			if NAME == "-xyz" {
				bbdp.x0 = fi.GetFloat()
				bbdp.y0 = fi.GetFloat()
				bbdp.z0 = fi.GetFloat()
			} else if NAME == "-WA" {
				bbdp.Wa = fi.GetFloat()
			} else if NAME == "-WB" {
				bbdp.Wb = fi.GetFloat()
			} else if NAME == "-WH" {
				bbdp.exw = fi.GetFloat()
				bbdp.exh = fi.GetFloat()
			} else if NAME == "-exs" {
				// Satoh修正（2018/1/23）
				bbdp.exsfname = fi.GetToken()

				//外表面の検索して、Wa,Wbを設定
				id := false
				for _, Exs := range Exsf.Exs {
					if bbdp.exsfname == Exs.Name {
						bbdp.Wa = Exs.Wa
						bbdp.Wb = Exs.Wb
						id = true
						break
					}
				}
				if id == false {
					fmt.Printf("BDP<%s> %s is not found in EXSRF\n", bbdp.bdpname, bbdp.exsfname)
				}
			} else {
				fmt.Printf("ERROR parameter----BDP %s\n", NAME)
				os.Exit(1)
			}
		}

		// SBLKの個数を数えてメモリを確保
		bbdp.SBLK = make([]*sunblk, 0)

		// RMPのメモリを確保
		bbdp.RMP = make([]*RRMP, 0)

		// if rp != nil {
		// 	wp = rp.WD
		// }

		// SBLK, RMPの読み込み
		for !fi.IsEnd() {
			//bbdp = (*bp)[i]

			NAME = fi.GetToken()
			if NAME[0] == '*' {
				break
			}

			if NAME == "SBLK" {
				// 日よけ
				// `SBLK <sbfname> <snbname> -xy <x> <y> -DW <D> <W> -a <WA> -rgb <r> <g> <b> ;`
				sb = SBLKInit()
				sb.ref = 0.0
				sb.sbfname = fi.GetToken()

				// 日よけの種類に応じた読み取り処理
				if sb.sbfname == "HISASI" {
					HISASHI(fi, sb)
				} else if sb.sbfname == "BARUKONI" {
					BARUKO(fi, sb)
				} else if sb.sbfname == "SODEKABE" {
					SODEK(fi, sb)
				} else if sb.sbfname == "MADOHIYOKE" {
					SCREEN(fi, sb)
				} else {
					fmt.Printf("ERROR----\nhiyoke no syurui <HISASI> or <BARUKONI> or <SODEKABE> or <MADOHIYOKE> : %s \n", sb.sbfname)
					os.Exit(1)
				}

				fi.SkipToEndOfLine()

				bbdp.SBLK = append(bbdp.SBLK, sb)
			} else if NAME == "RMP" {

				// RMP 読み取り処理
				rp := rmpdata(fi)

				fi.SkipToEndOfLine()

				bbdp.RMP = append(bbdp.RMP, rp)
			} else {
				fmt.Printf("ERROR----<SBLK> or <RMP> : %s \n", NAME)
				os.Exit(1)
			}
		}

		(*bp) = append(*bp, bbdp)
	}
}

/*
bdpinit (Building Data Point Initialization)

この関数は、新しい建物の基準点（BDP）のデータ構造を初期化します。

建築環境工学的な観点:
- **BDPの初期化**: 建物のシミュレーションを行う前に、
  その幾何学的パラメータ（基準点座標`x0`, `y0`, `z0`、方位角`Wa`、傾斜角`Wb`など）を
  デフォルト値（通常はゼロ）で初期化します。
  これにより、後続のデータ入力や計算が正しく行われるための準備が整います。
- **建物全体の幾何学的モデル化の基礎**: これらのパラメータは、
  建物全体を3次元空間に配置し、
  太陽位置との相対関係を正確にモデル化するための基礎となります。
  正確な初期化は、日射熱取得、日影、昼光利用、
  そして太陽光発電システムの発電量予測をシミュレーションする上で重要です。

この関数は、建物の幾何学的モデルを構築し、
日射熱取得の抑制、冷房負荷の軽減、
昼光利用の最適化、および日影計算を行うための重要な役割を果たします。
*/
func bdpinit() *BBDP {
	bbdp := new(BBDP)
	bbdp.bdpname = ""
	bbdp.exh = 0
	bbdp.exw = 0.
	bbdp.x0 = 0
	bbdp.y0 = 0
	bbdp.z0 = 0.
	bbdp.Wa = 0
	bbdp.Wb = 0.
	bbdp.SBLK = nil
	bbdp.RMP = nil
	bbdp.exsfname = ""
	return bbdp
}

/*
RRMPInit (Room Main Plane Initialization)

この関数は、室の主面（壁、床、天井など）の新しいデータ構造を初期化します。

建築環境工学的な観点:
- **室の幾何学的モデル化の準備**: 室の熱負荷計算や昼光利用計算では、
  室を構成する各面（壁、床、天井など）の幾何学的情報が不可欠です。
  この関数は、主面の位置（`xb0`, `yb0`）、寸法（`Rw`, `Rh`）、
  反射率（`ref`）、色（`rgb`）などをデフォルト値で初期化します。
  これにより、後続のデータ入力や計算が正しく行われるための準備が整います。
- **熱的・光学的特性の定義**: これらのパラメータは、
  主面からの熱伝達や、日射の反射、昼光の取り込みなどを計算するための基礎となります。
  正確な初期化は、室の熱的・光学的挙動をシミュレーションする上で重要です。

この関数は、室の幾何学的モデル化と熱的特性の定義を行い、
熱負荷計算、エネルギー消費量予測、
昼光利用の最適化、および快適性評価を行うための重要な役割を果たします。
*/
func RRMPInit() *RRMP {
	rp := new(RRMP)

	// RMP名
	rp.rmpname = ""

	// 壁名称
	rp.wallname = ""

	// 反射率
	rp.ref = 0.0

	// 左下頂点座標
	rp.xb0 = 0.0
	rp.yb0 = 0.0

	// 巾、高さ
	rp.Rw = 0.0
	rp.Rh = 0.0

	// 色
	matinit(rp.rgb[:], 3)
	rp.rgb[0] = 0.9
	rp.rgb[1] = 0.9
	rp.rgb[2] = 0.9

	// 窓
	rp.WD = nil

	// 前面地面の代表点までの距離
	rp.grpx = 1.0

	return rp
}

/*
SBLKInit (Sunbreak Initialization)

この関数は、新しい日よけ（庇、バルコニー、袖壁、窓日よけなど）のデータ構造を初期化します。

建築環境工学的な観点:
- **日よけの初期化**: 日よけのシミュレーションを行う前に、
  その幾何学的パラメータ（奥行き`D`、高さ`H`、幅`W`、窓からの距離`h`など）や、
  熱的・光学的特性（反射率`ref`、色`rgb`など）をデフォルト値で初期化します。
  これにより、後続のデータ入力や計算が正しく行われるための準備が整います。
- **日射遮蔽の基礎**: これらのパラメータは、
  日よけによって形成される日影の形状や範囲を計算するための基礎となります。
  正確な初期化は、日射遮蔽効果の評価や、
  日影規制の遵守をシミュレーションする上で重要です。

この関数は、建物の日射遮蔽計画をモデル化し、
日射熱取得の抑制、冷房負荷の軽減、
昼光利用の最適化、および日影計算を行うための重要な役割を果たします。
*/
func SBLKInit() *sunblk {
	sb := new(sunblk)

	sb.D = 0.0
	sb.H = 0.0
	sb.h = 0.0
	sb.ref = 0.0
	sb.W = 0.0
	sb.WA = 0.0
	sb.x = 0.0
	sb.y = 0.0
	sb.sbfname = ""
	sb.snbname = ""
	matinit(sb.rgb[:], 3)

	return sb
}

/*--------------------------------------------------------------------------*/
/*
obsdata (Obstacle Data Input)

この関数は、外部障害物（長方形、直方体、三角形など）の仕様データを入力ファイルから読み込み、
対応する構造体（`OBS`）に格納します。
これらのデータは、日影計算や日射量計算において、
周囲の建物や地形による日影の影響を評価するために不可欠です。

建築環境工学的な観点:
- **外部障害物のモデル化**: 建物の周囲に存在する他の建物や構造物、
  あるいは地形などは、日影を形成し、
  対象建物への日射入射量に影響を与えます。
  この関数は、障害物の種類（`fname`）に応じて、
  それぞれの幾何学的パラメータを読み込みます。
  - `rect`: 長方形（平面）の障害物。
  - `cube`: 直方体の障害物。
  - `r_tri`, `i_tri`: 三角形の障害物。
- **日影計算の基礎**: この関数で入力されるデータは、
  対象建物への日影の影響を定量的に評価するための基礎となります。
  これにより、日射熱取得の予測精度を向上させ、
  冷房負荷を正確に評価できます。
- **反射率と色彩の考慮**: `ref`は、障害物表面の反射率を指定します。
  反射率の高い障害物は、日射を反射することで、
  対象建物への反射日射を増加させる可能性があります。
  `rgb`は、障害物の表面色をRGB値で指定します。

この関数は、建物の日射環境を正確にモデル化し、
日射熱取得の抑制、冷房負荷の軽減、
昼光利用の最適化、および日影計算を行うための重要なデータ入力機能を提供します。
*/
func obsdata(fi *EeTokens, obsn *int, obs *[]*OBS) {
	var i, Nobs int
	var obsp *OBS

	*obs = make([]*OBS, Nobs)
	*obsn = 0

	for !fi.IsEnd() {
		obsp = obsinit()

		line := new(EeTokens)
		line.tokens = fi.GetLogicalLine()
		line.pos = 0

		NAME := line.GetToken()
		if NAME[0] == '*' {
			break
		}

		// 外部障害物の種類
		obsp.fname = NAME

		// 反射率の初期化
		for i = 0; i < 4; i++ {
			obsp.ref[i] = 0.0
		}

		// 外部障害物の種類に応じた読み取り処理
		if obsp.fname == "rect" {
			// 長方形（平面）
			rectdata(line, obsp)
		} else if obsp.fname == "cube" {
			// 直方体
			cubdata(line, obsp)
		} else if obsp.fname == "r_tri" || obsp.fname == "i_tri" {
			// 三角形
			tridata(line, obsp)
		} else {
			fmt.Printf("ERROR parameter----OBS : %s\n", obsp.fname)
			os.Exit(1)
		}

		*obs = append(*obs, obsp)
		(*obsn)++
	}
}

/*
obsinit (Obstacle Initialization)

この関数は、新しい外部障害物のデータ構造を初期化します。

建築環境工学的な観点:
- **外部障害物の初期化**: 外部障害物のシミュレーションを行う前に、
  その幾何学的パラメータ（座標`x`, `y`, `z`、幅`W`、奥行き`D`、高さ`H`、
  方位角`Wa`、傾斜角`Wb`など）や、
  熱的・光学的特性（反射率`ref`、色`rgb`など）をデフォルト値で初期化します。
  これにより、後続のデータ入力や計算が正しく行われるための準備が整います。
- **日影計算の基礎**: これらのパラメータは、
  外部障害物によって形成される日影の形状や範囲を計算するための基礎となります。
  正確な初期化は、日射遮蔽効果の評価や、
  日影規制の遵守をシミュレーションする上で重要です。

この関数は、建物の日射環境をモデル化し、
日射熱取得の抑制、冷房負荷の軽減、
昼光利用の最適化、および日影計算を行うための重要な役割を果たします。
*/
func obsinit() *OBS {
	obsp := new(OBS)

	// 外部障害物の種類
	obsp.fname = ""

	// 名前
	obsp.obsname = ""

	// 左下頂点座標
	obsp.x = 0.0
	obsp.y = 0.0
	obsp.z = 0.0

	// 巾、奥行き、高さ
	obsp.H = 0.0
	obsp.D = 0.0
	obsp.W = 0.0

	// 方位角、傾斜角
	obsp.Wa = 0.0
	obsp.Wb = 0.0

	// 反射率
	matinit(obsp.ref[:], 4)

	// 色
	matinit(obsp.rgb[:], 3)

	return obsp
}

/*
OPcount (Opening Plane Count)

この関数は、日射量計算や昼光利用計算で用いられる「受光面（Opening Plane, OP）」の総数をカウントします。

建築環境工学的な観点:
- **受光面の総数把握**: 建物のエネルギーシミュレーションでは、
  日射熱取得や昼光利用を評価するために、
  窓や壁面などの受光面をモデル化します。
  この関数は、BDP（建物の基準点）に付随するRMP（室の主面）やWD（窓）、
  および多角形データで直接定義されたRMPの総数をカウントします。
- **シミュレーションの準備**: 受光面の総数を事前に把握することで、
  シミュレーションに必要なメモリ領域を確保したり、
  計算ループの回数を決定したりすることができます。
- **日射熱取得と昼光利用の評価**: 受光面の総数は、
  建物全体の日射熱取得ポテンシャルや、
  昼光利用の可能性を評価する上で基本的な情報となります。

この関数は、建物の日射環境を正確にモデル化し、
日射熱取得の抑制、冷房負荷の軽減、
昼光利用の最適化、および日影計算を行うための重要な役割を果たします。
*/
func OPcount(_Bdp []*BBDP, _poly []*POLYGN) int {
	Nop := 0

	// BDPの受光面 = RMP + WD
	// (日よけは含まない)
	for _, Bdp := range _Bdp {
		Nop += len(Bdp.RMP)
		for _, RMP := range Bdp.RMP {
			Nop += len(RMP.WD)
		}
	}

	// ポリゴン指定の受照面
	for _, poly := range _poly {
		if poly.polyknd == "RMP" {
			// OBS(Obstacle)は受光面に含まない
			Nop++
		}
	}

	return Nop
}

/*
LPcount (Light-Receiving Plane Count)

この関数は、日影計算や日射量計算で用いられる「被受照面（Light-Receiving Plane, LP）」の総数をカウントします。

建築環境工学的な観点:
- **被受照面の総数把握**: 建物のエネルギーシミュレーションでは、
  日影の影響を評価するために、
  庇、バルコニー、袖壁、樹木、外部障害物などの被受照面をモデル化します。
  この関数は、BDP（建物の基準点）に付随するSBLK（日よけ）、
  OBS（外部障害物）、TREE（樹木）、
  および多角形データで直接定義されたRMPやOBSの総数をカウントします。
- **シミュレーションの準備**: 被受照面の総数を事前に把握することで、
  シミュレーションに必要なメモリ領域を確保したり、
  計算ループの回数を決定したりすることができます。
- **日影計算の複雑性**: 日よけの種類（バルコニーなど）や、
  外部障害物の種類（立方体など）によって、
  一つのオブジェクトが複数の被受照面を持つ場合があるため、
  その複雑性を考慮してカウントを行います。

この関数は、建物の日射環境を正確にモデル化し、
日射熱取得の抑制、冷房負荷の軽減、
昼光利用の最適化、および日影計算を行うための重要な役割を果たします。
*/
func LPcount(_Bdp []*BBDP, _Obs []*OBS, _tree []*TREE, _poly []*POLYGN) int {
	Nlp := 0

	// BDPの被光面 = 日よけ
	for _, Bdp := range _Bdp {
		for _, snbk := range Bdp.SBLK {
			if snbk.sbfname == "BARUKONI" {
				// 日よけの種類がバルコニーの場合
				Nlp += 5
			} else {
				Nlp++
			}
		}
	}

	// OBSの被光面
	for _, Obs := range _Obs {
		if Obs.fname == "cube" {
			// 外部依存障害物の種類が立方体の場合
			Nlp += 4
		} else {
			Nlp++
		}
	}

	// 樹木用
	Nlp += len(_tree) * 20

	// ポリゴン指定の被照面
	for _, poly := range _poly {
		if poly.polyknd == "RMP" || poly.polyknd == "OBS" {
			// RMPは受光面であり被受光面でもある
			Nlp++
		}
	}

	return Nlp
}

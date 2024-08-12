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

// 庇
// `SBLK HISASI <snbname> -xy <x> <y> -DW <D> <W> -a <WA> -rgb <r> <g> <b> ;`
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

// バルコニー
// `SBLK BARUKONI <snbname> -xy <x> <y> -DHWh <D> <H> <W> <h> -ref <ref> -rgb <r> <g> <b> ;`
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

// 袖壁
// `SBLK SODEKABE <snbname> -xy <x> <y> -DH <D> <H> -rgb <r> <g> <b> ;`
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

// 日よけ
// `SBLK MADOHIYOKE <snbname> -xy <x> <y> -DHW <D> <W> <H> -rgb <r> <g> <b> ;`
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
// 以下のようなRMPデータ(WD含む)を読み取る。
// WDデータ複数行ある場合があるので注意する。
// ---------------------------------------------------------
// RMP <rmpname> <wallname> -xyb <xb0> <yb0> -WH <Rw> <Rh> -ref <ref> -grpx <grpx> -rgb <r> <g> <b> ;
//   WD <winname> -xyr <xr> <yr> -WH <Ww> <Wh> -ref <ref> -grpx <grpx> -rgb <r> <g> <b> ;
//   WD <winname> -xyr <xr> <yr> -WH <Ww> <Wh> -ref <ref> -grpx <grpx> -rgb <r> <g> <b> ;
// ;
// ---------------------------------------------------------
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

// `rect <obsname> -xyz <x> <y> <z> -WH <W> <H> -WaWb <Wa> <Wb> -ref <ref> -rgb <r> <g> <b> ;`
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

// `cube <obsname> -xyz <x> <y> <z> -WDH <W> <D> <H> -Wa <Wa> -ref0 <ref0> -ref1 <ref1> -ref2 <ref2> -ref3 <ref3> -rgb <r> <g> <b> ;`
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

// `r_tri <obsname> -xyz <x> <y> <z> -WH <W> <H> -WaWb <Wa> <Wb> -ref <ref> -rgb <r> <g> <b> ;`
// `i_tri <obsname> -xyz <x> <y> <z> -WH <W> <H> -WaWb <Wa> <Wb> -ref <ref> -rgb <r> <g> <b> ;`
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
// fi から入力データを読み取り、monten と DE に値を設定する。
// 20170503 higuchi add
// ------------------------------
// `DIVID
//    DE <DE> MONT <monte> ;
//  *`
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

// `TREE <treetype> <treename> -xyz <x> <y> <z> -WH1 <W1> <H1> -WH2 <W2> <H2> -WH3 <W3> <H3> -WH4 <W4> -rgb <r> <g> <b> ;`
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
// BDP, SBLK データの読み込み
// 例: BDP Ssrfs -xyz 0 0 0.5 -exs south -WH 5.95 2.9 ;
// 例: SBLK HISASI Ssblk -xy 1.125 2.9 -DW 0.9 3.7 -a 90 ;
// 例: RMP Swall LD -xyb 0 0 -WH 5.95 2.9 -ref 0.1 ;
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
// OBS
//   rect obs0 -xyz 22.5 -7.5 0.0 -WH 9.000 8.150 -WaWb 180 90 -ref 0.3 ;
//   rect obs11 -xyz -13.5 15.0 0 -WH 9.000 8.150 -WaWb 0 90 -ref 0.3 ;
//   cube obs12 -xyz 25 20 0 -WDH 20 20 30 -Wa 0 ;
// *
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

// OP(受照面)のカウント？
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

// LP(被受照面)のカウント？
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

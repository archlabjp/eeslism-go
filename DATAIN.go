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

package main

import (
	"fmt"
	"io"
	"os"
)

/*

  入力データの読み込み
  FILE=DATAIN.c
  Create Date=1999.6.7
*/

func HISASHI(fi *os.File, sb *sunblk) {
	var NAME string

	fmt.Fscanf(fi, "%s", &NAME)
	sb.snbname = NAME

	// 色の初期値
	sb.rgb[0] = 0.0
	sb.rgb[1] = 0.2
	sb.rgb[2] = 0.0

	for {
		fmt.Fscanf(fi, "%s", &NAME)
		if NAME[0] == ';' {
			break
		}

		if NAME == "-xy" {
			fmt.Fscanf(fi, "%f", &sb.x)
			fmt.Fscanf(fi, "%f", &sb.y)
		} else if NAME == "-DW" {
			fmt.Fscanf(fi, "%f", &sb.D)
			fmt.Fscanf(fi, "%f", &sb.W)
		} else if NAME == "-a" {
			fmt.Fscanf(fi, "%f", &sb.WA)
		} else if NAME == "-rgb" {
			fmt.Fscanf(fi, "%f", &sb.rgb[0])
			fmt.Fscanf(fi, "%f", &sb.rgb[1])
			fmt.Fscanf(fi, "%f", &sb.rgb[2])
		} else {
			fmt.Printf("ERROR parameter----HISASI: %s\n", NAME)

			os.Exit(1)
		}
	}
}

/*--------------------------------------------------------------*/

func BARUKO(fi *os.File, sb *sunblk) {
	var NAME string

	sb.ref = 0.0

	sb.rgb[0] = 0.0
	sb.rgb[1] = 0.2
	sb.rgb[2] = 0.0

	fmt.Fscanf(fi, "%s", &NAME)
	sb.snbname = NAME

	for {
		fmt.Fscanf(fi, "%s", &NAME)
		if NAME[0] == ';' {
			break
		}

		if NAME == "-xy" {
			fmt.Fscanf(fi, "%f", &sb.x)
			fmt.Fscanf(fi, "%f", &sb.y)
		} else if NAME == "-DHWh" {
			fmt.Fscanf(fi, "%f", &sb.D)
			fmt.Fscanf(fi, "%f", &sb.H)
			fmt.Fscanf(fi, "%f", &sb.W)
			fmt.Fscanf(fi, "%f", &sb.h)
		} else if NAME == "-ref" {
			fmt.Fscanf(fi, "%f", &sb.ref)
		} else if NAME == "-rgb" {
			fmt.Fscanf(fi, "%f", &sb.rgb[0])
			fmt.Fscanf(fi, "%f", &sb.rgb[1])
			fmt.Fscanf(fi, "%f", &sb.rgb[2])
		} else {
			fmt.Printf("ERROR parameter----WBARUKONI: %s\n", NAME)

			os.Exit(1)
		}
	}
}

/*------------------------------------------------------------------*/

func SODEK(fi *os.File, sb *sunblk) {
	var NAME string

	sb.rgb[0] = 0.0
	sb.rgb[1] = 0.2
	sb.rgb[2] = 0.0

	fmt.Fscanf(fi, "%s", &NAME)
	sb.snbname = NAME

	for {
		fmt.Fscanf(fi, "%s", &NAME)
		if NAME[0] == ';' {
			break
		}

		if NAME == "-xy" {
			fmt.Fscanf(fi, "%f", &sb.x)
			fmt.Fscanf(fi, "%f", &sb.y)
		} else if NAME == "-DH" {
			fmt.Fscanf(fi, "%f", &sb.D)
			fmt.Fscanf(fi, "%f", &sb.H)
		} else if NAME == "-a" {
			fmt.Fscanf(fi, "%f", &sb.WA)
		} else if NAME == "-rgb" {
			fmt.Fscanf(fi, "%f", &sb.rgb[0])
			fmt.Fscanf(fi, "%f", &sb.rgb[1])
			fmt.Fscanf(fi, "%f", &sb.rgb[2])
		} else {
			fmt.Printf("ERROR parameter----SODEKABE: %s\n", NAME)

			os.Exit(1)
		}
	}
}

/*-----------------------------------------------------------------------*/

func SCREEN(fi *os.File, sb *sunblk) {
	var NAME string

	sb.rgb[0] = 0.0
	sb.rgb[1] = 0.2
	sb.rgb[2] = 0.0

	fmt.Fscanf(fi, "%s", &NAME)
	sb.snbname = NAME

	for {
		fmt.Fscanf(fi, "%s", &NAME)
		if NAME[0] == ';' {
			break
		}

		if NAME == "-xy" {
			fmt.Fscanf(fi, "%f", &sb.x)
			fmt.Fscanf(fi, "%f", &sb.y)
		} else if NAME == "-DHW" {
			fmt.Fscanf(fi, "%f", &sb.D)
			fmt.Fscanf(fi, "%f", &sb.H)
			fmt.Fscanf(fi, "%f", &sb.W)
		} else if NAME == "-rgb" {
			fmt.Fscanf(fi, "%f", &sb.rgb[0])
			fmt.Fscanf(fi, "%f", &sb.rgb[1])
			fmt.Fscanf(fi, "%f", &sb.rgb[2])
		} else {
			fmt.Printf("ERROR paramater---MADOHIYOKE: %s\n", NAME)

			os.Exit(1)
		}
	}
}

/*----------------------------------------------------------------*/

func rmpdata(fi *os.File, rp *RRMP, _wp []MADO) {
	var NAME string

	rp.ref = 0.0
	rp.grpx = 1.0

	rp.rgb[0] = 0.9
	rp.rgb[1] = 0.9
	rp.rgb[2] = 0.9

	fmt.Fscanf(fi, "%s", &NAME)
	rp.rmpname = NAME
	fmt.Fscanf(fi, "%s", &NAME)
	rp.wallname = NAME

	for {
		fmt.Fscanf(fi, "%s", &NAME)
		if NAME[0] == ';' {
			break
		}

		if NAME == "-xyb" {
			fmt.Fscanf(fi, "%f", &rp.xb0)
			fmt.Fscanf(fi, "%f", &rp.yb0)
		} else if NAME == "-WH" {
			fmt.Fscanf(fi, "%f", &rp.Rw)
			fmt.Fscanf(fi, "%f", &rp.Rh)
		} else if NAME == "-ref" {
			fmt.Fscanf(fi, "%f", &rp.ref)
		} else if NAME == "-grpx" {
			fmt.Fscanf(fi, "%f", &rp.grpx)
		} else if NAME == "-rgb" {
			fmt.Fscanf(fi, "%f", &rp.rgb[0])
			fmt.Fscanf(fi, "%f", &rp.rgb[1])
			fmt.Fscanf(fi, "%f", &rp.rgb[2])
		} else {
			fmt.Printf("ERROR parameter----RMP: %s\n", NAME)
			os.Exit(1)
		}
	}

	rp.sumWD = 0
	for _, wp := range _wp {
		fmt.Fscanf(fi, "%s", &NAME)
		if NAME[0] == ';' {
			break
		}

		wp.ref = 0.0
		wp.grpx = 1.0

		wp.rgb[0] = 0.0
		wp.rgb[1] = 0.3
		wp.rgb[2] = 0.8

		if NAME != "WD" {
			fmt.Printf("ERROR parameter----WD: %s\n", NAME)
			os.Exit(1)
		}

		rp.sumWD++

		fmt.Fscanf(fi, "%s", &NAME)
		wp.winname = NAME

		for {
			fmt.Fscanf(fi, "%s", &NAME)
			if NAME[0] == ';' {
				break
			}

			if NAME == "-xyr" {
				fmt.Fscanf(fi, "%f", &wp.xr)
				fmt.Fscanf(fi, "%f", &wp.yr)
			} else if NAME == "-WH" {
				fmt.Fscanf(fi, "%f", &wp.Ww)
				fmt.Fscanf(fi, "%f", &wp.Wh)
			} else if NAME == "-ref" {
				fmt.Fscanf(fi, "%f", &wp.ref)
			} else if NAME == "-grpx" {
				fmt.Fscanf(fi, "%f", &wp.grpx)
			} else if NAME == "-rgb" {
				fmt.Fscanf(fi, "%f", &wp.rgb[0])
				fmt.Fscanf(fi, "%f", &wp.rgb[1])
				fmt.Fscanf(fi, "%f", &wp.rgb[2])
			} else {
				fmt.Printf("ERROR parameter----WD: %s\n", NAME)
				os.Exit(1)
			}
		}
	}
}

/*------------------------------------------------------------------*/
func rectdata(fi *os.File, obs *OBS) {
	var NAME string

	obs.ref[0] = 0.0

	obs.rgb[0] = 0.7
	obs.rgb[1] = 0.7
	obs.rgb[2] = 0.7

	fmt.Fscanf(fi, "%s", &NAME)
	obs.obsname = NAME

	for {
		fmt.Fscanf(fi, "%s", &NAME)
		if NAME[0] == ';' {
			break
		}

		if NAME == "-xyz" {
			fmt.Fscanf(fi, "%f", &obs.x)
			fmt.Fscanf(fi, "%f", &obs.y)
			fmt.Fscanf(fi, "%f", &obs.z)
		} else if NAME == "-WH" {
			fmt.Fscanf(fi, "%f", &obs.W)
			fmt.Fscanf(fi, "%f", &obs.H)
		} else if NAME == "-WaWb" {
			fmt.Fscanf(fi, "%f", &obs.Wa)
			fmt.Fscanf(fi, "%f", &obs.Wb)
		} else if NAME == "-ref" {
			fmt.Fscanf(fi, "%f", &obs.ref[0])
		} else if NAME == "-rgb" {
			fmt.Fscanf(fi, "%f", &obs.rgb[0])
			fmt.Fscanf(fi, "%f", &obs.rgb[1])
			fmt.Fscanf(fi, "%f", &obs.rgb[2])
		} else {
			fmt.Printf("ERROR parameter----OBS.rect: %s\n", NAME)
			os.Exit(1)
		}
	}
}

/*------------------------------------------------------------------*/
func cubdata(fi *os.File, obs *OBS) {
	var NAME string

	for i := 0; i < 3; i++ {
		obs.ref[i] = 0.0
	}

	obs.rgb[0] = 0.7
	obs.rgb[1] = 0.7
	obs.rgb[2] = 0.7

	fmt.Fscanf(fi, "%s", &NAME)
	obs.obsname = NAME

	for {
		fmt.Fscanf(fi, "%s", &NAME)
		if NAME[0] == ';' {
			break
		}

		if NAME == "-xyz" {
			fmt.Fscanf(fi, "%f", &obs.x)
			fmt.Fscanf(fi, "%f", &obs.y)
			fmt.Fscanf(fi, "%f", &obs.z)
		} else if NAME == "-WDH" {
			fmt.Fscanf(fi, "%f", &obs.W)
			fmt.Fscanf(fi, "%f", &obs.D)
			fmt.Fscanf(fi, "%f", &obs.H)
		} else if NAME == "-Wa" {
			fmt.Fscanf(fi, "%f", &obs.Wa)
		} else if NAME == "-ref0" {
			fmt.Fscanf(fi, "%f", &obs.ref[0])
		} else if NAME == "-ref1" {
			fmt.Fscanf(fi, "%f", &obs.ref[1])
		} else if NAME == "-ref2" {
			fmt.Fscanf(fi, "%f", &obs.ref[2])
		} else if NAME == "-ref3" {
			fmt.Fscanf(fi, "%f", &obs.ref[3])
		} else if NAME == "-rgb" {
			fmt.Fscanf(fi, "%f", &obs.rgb[0])
			fmt.Fscanf(fi, "%f", &obs.rgb[1])
			fmt.Fscanf(fi, "%f", &obs.rgb[2])
		} else {
			fmt.Printf("ERROR parameter----OBS.cube: %s\n", NAME)
			os.Exit(1)
		}
	}
}

/*-------------------------------------------------------------------*/
func tridata(fi *os.File, obs *OBS) {
	var NAME string

	obs.ref[0] = 0.0

	obs.rgb[0] = 0.7
	obs.rgb[1] = 0.7
	obs.rgb[2] = 0.7

	fmt.Fscanf(fi, "%s", &obs.obsname)

	for {
		fmt.Fscanf(fi, "%s", &NAME)
		if NAME[0] == ';' {
			break
		}

		if NAME == "-xyz" {
			fmt.Fscanf(fi, "%f", &obs.x)
			fmt.Fscanf(fi, "%f", &obs.y)
			fmt.Fscanf(fi, "%f", &obs.z)
		} else if NAME == "-WH" {
			fmt.Fscanf(fi, "%f", &obs.W)
			fmt.Fscanf(fi, "%f", &obs.H)
		} else if NAME == "-WaWb" {
			fmt.Fscanf(fi, "%f", &obs.Wa)
			fmt.Fscanf(fi, "%f", &obs.Wb)
		} else if NAME == "-ref" {
			fmt.Fscanf(fi, "%f", &obs.ref[0])
		} else if NAME == "-rgb" {
			fmt.Fscanf(fi, "%f", &obs.rgb[0])
			fmt.Fscanf(fi, "%f", &obs.rgb[1])
			fmt.Fscanf(fi, "%f", &obs.rgb[2])
		} else {
			fmt.Printf("ERROR parameter----OBS.triangle: %s\n", NAME)
			os.Exit(1)
		}
	}
}

/*-------------------------------------------------------------------*/
// 20170503 higuchi add
func dividdata(fi *os.File, monten *int, DE *float64) {
	var NAME string

	for {
		fmt.Fscanf(fi, "%s", &NAME)
		if NAME[0] == ';' {
			break
		}

		if NAME == "DE" {
			fmt.Fscanf(fi, "%f", DE)
		} else if NAME == "MONT" {
			fmt.Fscanf(fi, "%d", monten)
		} else {
			fmt.Printf("ERROR parameter----DIVID: %s\n", NAME)

			os.Exit(1)
		}
	}

	fmt.Fscanf(fi, "%s", &NAME)
}

func treedata(fi *os.File, treen *int, tree *[]TREE) {
	var i int
	var Ntree int
	var tred *TREE

	// BDPの数を数える
	Ntree = InputCount(fi, ";")
	fmt.Printf("<treedata> Ntree=%d\n", Ntree)

	if Ntree > 0 {
		*tree = make([]TREE, Ntree)

		// 構造体の初期化
		for i = 0; i < Ntree; i++ {
			tred = &(*tree)[i]

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
		}
	}

	*treen = 0

	for i = 0; i < Ntree; i++ {
		tred = &(*tree)[i]

		var NAME string
		fmt.Fscanf(fi, "%s", &NAME)
		if NAME[0] == '*' {
			break
		}

		tred.treetype = NAME

		fmt.Fscanf(fi, "%s", &NAME)
		tred.treename = NAME

		if tred.treetype == "treeA" {
			for {
				fmt.Fscanf(fi, "%s", &NAME)
				if NAME[0] == ';' {
					break
				}

				if NAME == "-xyz" {
					fmt.Fscanf(fi, "%f", &tred.x)
					fmt.Fscanf(fi, "%f", &tred.y)
					fmt.Fscanf(fi, "%f", &tred.z)
				} else if NAME == "-WH1" {
					fmt.Fscanf(fi, "%f", &tred.W1)
					fmt.Fscanf(fi, "%f", &tred.H1)
				} else if NAME == "-WH2" {
					fmt.Fscanf(fi, "%f", &tred.W2)
					fmt.Fscanf(fi, "%f", &tred.H2)
				} else if NAME == "-WH3" {
					fmt.Fscanf(fi, "%f", &tred.W3)
					fmt.Fscanf(fi, "%f", &tred.H3)
				} else if NAME == "-WH4" {
					fmt.Fscanf(fi, "%f", &tred.W4)
				} else {
					fmt.Printf("ERROR parameter----TREE: %s %s\n", tred.treename, NAME)
					os.Exit(1)
				}
			}
		} else {
			fmt.Printf("ERROR parameter----TREE: %s\n", tred.treetype)
			os.Exit(1)
		}

		(*treen)++
	}
}

/*-------------------------*/
func polydata(fi *os.File, polyn *int, poly *[]POLYGN) {
	var i int
	var Npoly int
	var polyp *POLYGN

	// BDPの数を数える
	Npoly = InputCount(fi, ";")
	fmt.Printf("<polydata> Npoly=%d\n", Npoly)

	if Npoly > 0 {
		*poly = make([]POLYGN, Npoly)

		// 構造体の初期化
		for i = 0; i < Npoly; i++ {
			polyp = &(*poly)[i]
			polyp.polyknd = ""
			polyp.polyname = ""
			polyp.wallname = ""
			polyp.polyd = 0
			polyp.ref = 0.0
			polyp.refg = 0.0
			polyp.grpx = 0.0
			polyp.P = nil
			matinit(polyp.rgb[:], 3)
		}
	}

	*polyn = 0
	for i = 0; i < Npoly; i++ {
		polyp = &(*poly)[i]

		var NAME string
		fmt.Fscanf(fi, "%s", &NAME)
		if NAME[0] == '*' {
			break
		}

		polyp.grpx = 1.0

		polyp.rgb[0] = 0.9
		polyp.rgb[1] = 0.9
		polyp.rgb[2] = 0.9

		polyp.polyknd = NAME

		if polyp.polyknd != "RMP" && polyp.polyknd != "OBS" {
			fmt.Printf("ERROR parameter----POLYGON: %s  <RMP> or <OBS> \n", polyp.polyknd)
			os.Exit(1)
		}

		fmt.Fscanf(fi, "%d", &polyp.polyd)
		polyp.P = make([]XYZ, polyp.polyd)

		fmt.Fscanf(fi, "%s", &NAME)
		polyp.polyname = NAME
		fmt.Fscanf(fi, "%s", &NAME)
		polyp.wallname = NAME

		for {
			fmt.Fscanf(fi, "%s", &NAME)
			if NAME[0] == ';' {
				break
			}

			if NAME == "-xyz" {
				for i = 0; i < polyp.polyd; i++ {
					fmt.Fscanf(fi, "%f", &polyp.P[i].X)
					fmt.Fscanf(fi, "%f", &polyp.P[i].Y)
					fmt.Fscanf(fi, "%f", &polyp.P[i].Z)
				}

			} else if NAME == "-rgb" {
				fmt.Fscanf(fi, "%f", &polyp.rgb[0])
				fmt.Fscanf(fi, "%f", &polyp.rgb[1])
				fmt.Fscanf(fi, "%f", &polyp.rgb[2])
			} else if NAME == "-ref" {
				fmt.Fscanf(fi, "%f", &polyp.ref)
			} else if NAME == "-refg" {
				fmt.Fscanf(fi, "%f", &polyp.refg)
			} else if NAME == "-grpx" {
				fmt.Fscanf(fi, "%f", &polyp.grpx)
			} else {
				fmt.Printf("ERROR parameter----POLYGON: %s\n", NAME)
				os.Exit(1)
			}
		}
		(*polyn)++
	}
}

/*---------------------------------------------------------------------------*/
func bdpdata(fi *os.File, bdpn *int, bp *[]BBDP, Exsf *EXSFS) {

	var rp *RRMP
	var wp *MADO
	var sb *sunblk
	var Nbdp int
	var bbdp *BBDP

	// BDPの数を数える
	Nbdp = InputCount(fi, "*")
	//printf("<bdpdata> Nbdp=%d\n", Nbdp)

	if Nbdp > 0 {
		*bp = make([]BBDP, Nbdp)
		if *bp == nil {
			fmt.Printf("<bdpdata> bpのメモリが確保できません\n")
		}

		for i := 0; i < Nbdp; i++ {
			bbdp = &(*bp)[i]
			bbdp.bdpname = ""
			bbdp.exh = 0
			bbdp.exw = 0.
			bbdp.sumRMP = 0
			bbdp.sumsblk = 0
			bbdp.x0 = 0
			bbdp.y0 = 0
			bbdp.z0 = 0.
			bbdp.Wa = 0
			bbdp.Wb = 0.
			bbdp.SBLK = nil
			bbdp.RMP = nil
			bbdp.exsfname = ""
		}
	}

	*bdpn = 0

	for i := 0; i < Nbdp; i++ {
		bbdp = &(*bp)[i]

		var NAME string
		fmt.Fscanf(fi, "%s", &NAME)
		if NAME[0] == '*' {
			break
		}

		if NAME != "BDP" {
			fmt.Printf("error BDP\n")
			os.Exit(1)
		}

		fmt.Fscanf(fi, "%s", &NAME)
		bbdp.bdpname = NAME

		for {
			fmt.Fscanf(fi, "%s", &NAME)
			if NAME[0] == ';' {
				break
			}

			if NAME == "-xyz" {
				fmt.Fscanf(fi, "%f", &bbdp.x0)
				fmt.Fscanf(fi, "%f", &bbdp.y0)
				fmt.Fscanf(fi, "%f", &bbdp.z0)
			} else if NAME == "-WA" {
				fmt.Fscanf(fi, "%f", &bbdp.Wa)
			} else if NAME == "-WB" {
				fmt.Fscanf(fi, "%f", &bbdp.Wb)
			} else if NAME == "-WH" {
				fmt.Fscanf(fi, "%f", &bbdp.exw)
				fmt.Fscanf(fi, "%f", &bbdp.exh)
			} else if NAME == "-exs" {
				// Satoh修正（2018/1/23）
				fmt.Fscanf(fi, "%s", &NAME)
				bbdp.exsfname = NAME

				//外表面の検索
				id := false
				for i := 0; i < Exsf.Nexs; i++ {
					Exs := Exsf.Exs[i]
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
		Nsblk := SBLKCount(fi)
		if Nsblk > 0 {
			bbdp.SBLK = make([]sunblk, Nsblk)

			for i := 0; i < Nsblk; i++ {
				sb = &bbdp.SBLK[i]

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
			}
		}

		// RMPの個数を数えてメモリを確保
		Nrmp := RMPCount(fi)
		if Nrmp > 0 {
			bbdp.RMP = make([]RRMP, Nrmp)
			for i := 0; i < Nrmp; i++ {
				rp = &bbdp.RMP[i]
				rp.rmpname = ""
				rp.wallname = ""
				rp.sumWD = 0
				rp.ref = 0.0
				rp.xb0 = 0.0
				rp.yb0 = 0.0
				rp.Rw = 0.0
				rp.Rh = 0.0
				rp.grpx = 0.0
				matinit(rp.rgb[:], 3)
				rp.WD = nil
			}
		}

		// if rp != nil {
		// 	wp = rp.WD
		// }

		sb_idx := 0
		rp_idx := 0
		for i := 0; i < len(*bp); i++ {
			bbdp = &(*bp)[i]

			sb = &bbdp.SBLK[sb_idx]
			rp = &bbdp.RMP[rp_idx]

			fmt.Fscanf(fi, "%s", &NAME)
			if NAME[0] == '*' {
				break
			}

			if NAME == "SBLK" {
				sb.ref = 0.0
				fmt.Fscanf(fi, "%s", &NAME)
				sb.sbfname = NAME

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

				sb_idx++
				bbdp.sumsblk++
			} else if NAME == "RMP" {
				// WDの数を数えてメモリを確保
				Nwd := WDCount(fi)

				if Nwd > 0 {
					rp.WD = make([]MADO, Nwd)
					for i := 0; i < Nwd; i++ {
						wp = &rp.WD[i]
						wp.winname = ""
						matinit(wp.rgb[:], 3)
						wp.grpx = 0.0
						wp.ref = 0.0
						wp.Wh = 0.0
						wp.xr = 0.0
						wp.yr = 0.0
					}
				}
				rp.ref = 0.0
				bbdp.sumRMP++
				rmpdata(fi, rp, rp.WD)

				rp_idx++
			} else {
				fmt.Printf("ERROR----<SBLK> or <RMP> : %s \n", NAME)
				os.Exit(1)
			}
		}

		(*bdpn)++
	}
}

/*--------------------------------------------------------------------------*/
func obsdata(fi *os.File, obsn *int, obs *[]OBS) {
	var i, Nobs int
	var obsp *OBS

	// Count the number of OBS entries
	Nobs = InputCount(fi, ";")
	if Nobs > 0 {
		*obs = make([]OBS, Nobs)
		for i = 0; i < Nobs; i++ {
			obsp = &(*obs)[i]
			obsp.fname = ""
			obsp.obsname = ""
			obsp.x = 0.0
			obsp.y = 0.0
			obsp.z = 0.0
			obsp.H = 0.0
			obsp.D = 0.0
			obsp.W = 0.0
			obsp.Wa = 0.0
			obsp.Wb = 0.0
			matinit(obsp.ref[:], 4)
			matinit(obsp.rgb[:], 3)
		}
	}

	*obsn = 0
	var NAME string
	for i = 0; i < Nobs; i++ {
		obsp = &(*obs)[i]

		_, err := fmt.Fscanf(fi, "%s", &NAME)
		if err != nil || NAME[0] == '*' {
			break
		}

		obsp.fname = NAME

		for i = 0; i < 4; i++ {
			obsp.ref[i] = 0.0
		}

		if obsp.fname == "rect" {
			rectdata(fi, obsp)
		} else if obsp.fname == "cube" {
			cubdata(fi, obsp)
		} else if obsp.fname == "r_tri" || obsp.fname == "i_tri" {
			tridata(fi, obsp)
		} else {
			fmt.Printf("ERROR parameter----OBS : %s\n", obsp.fname)
			os.Exit(1)
		}

		(*obsn)++
	}
}

func InputCount(fi *os.File, key string) int {
	N := 0
	ad, _ := fi.Seek(0, io.SeekCurrent)

	var s string
	for {
		_, err := fmt.Fscanf(fi, "%s", &s)
		if err == io.EOF || s == "*" {
			break
		}

		N++

		for {
			_, err := fmt.Fscanf(fi, "%s", &s)
			if err == io.EOF {
				break
			}

			if s == key {
				break
			}
		}
	}

	_, _ = fi.Seek(ad, io.SeekStart)
	return N
}

func SBLKCount(fi *os.File) int {
	N := 0
	ad, _ := fi.Seek(0, io.SeekCurrent)

	var s string
	for {
		_, err := fmt.Fscanf(fi, "%s", &s)
		if err == io.EOF || s == "*" {
			break
		}

		if s == "SBLK" {
			N++
		}
	}

	_, _ = fi.Seek(ad, io.SeekStart)
	return N
}

func RMPCount(fi *os.File) int {
	N := 0
	ad, _ := fi.Seek(0, io.SeekCurrent)

	var s string
	for {
		_, err := fmt.Fscanf(fi, "%s", &s)
		if err == io.EOF || s == "*" {
			break
		}

		if s == "RMP" {
			N++
		}
	}

	_, _ = fi.Seek(ad, io.SeekStart)
	return N
}

func WDCount(fi *os.File) int {
	N := 0
	ad, _ := fi.Seek(0, io.SeekCurrent)

	Flg := 0
	var s string
	for {
		_, err := fmt.Fscanf(fi, "%s", &s)
		if err == io.EOF {
			break
		}

		if s == "WD" {
			N++
		}

		if s == ";" {
			if Flg == 1 {
				break
			} else {
				Flg = 1
			}
		} else {
			Flg = 0
		}
	}

	_, _ = fi.Seek(ad, io.SeekStart)
	return N
}

func OPcount(Nbdp int, _Bdp []BBDP, Npoly int, _poly []POLYGN) int {
	Nop := 0

	for i := 0; i < Nbdp; i++ {
		Bdp := &_Bdp[i]
		Nop += Bdp.sumRMP
		for j := 0; j < Bdp.sumRMP; j++ {
			RMP := &Bdp.RMP[i]
			Nop += RMP.sumWD
		}
	}

	for i := 0; i < Npoly; i++ {
		poly := &_poly[i]
		if poly.polyknd == "RMP" {
			Nop++
		}
	}

	return Nop
}

func LPcount(Nbdp int, _Bdp []BBDP, Nobs int, _Obs []OBS, Ntree int, Npoly int, _poly []POLYGN) int {
	Nlp := 0

	//初期化
	for i := 0; i < Nbdp; i++ {
		Bdp := &_Bdp[i]
		for j := 0; j < Bdp.sumsblk; j++ {
			snbk := &Bdp.SBLK[j]
			if snbk.sbfname == "BARUKONI" {
				Nlp += 5
			} else {
				Nlp++
			}
		}
	}

	for i := 0; i < Nobs; i++ {
		Obs := &_Obs[i]
		if Obs.fname == "cube" {
			Nlp += 4
		} else {
			Nlp++
		}
	}

	// 樹木用
	Nlp += Ntree * 20

	// ポリゴン
	for i := 0; i < Npoly; i++ {
		poly := &_poly[i]
		if poly.polyknd == "RMP" || poly.polyknd == "OBS" {
			Nlp++
		}
	}

	return Nlp
}

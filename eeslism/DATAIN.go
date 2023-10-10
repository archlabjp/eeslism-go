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

func HISASHI(fi *EeTokens, sb *sunblk) {
	sb.snbname = fi.GetToken()

	// 色の初期値
	sb.rgb[0] = 0.0
	sb.rgb[1] = 0.2
	sb.rgb[2] = 0.0

	for fi.IsEnd() == false {
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

func BARUKO(fi *EeTokens, sb *sunblk) {
	sb.ref = 0.0

	sb.rgb[0] = 0.0
	sb.rgb[1] = 0.2
	sb.rgb[2] = 0.0

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
			sb.ref = fi.GetFloat()
		} else if NAME == "-rgb" {
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

func SODEK(fi *EeTokens, sb *sunblk) {
	sb.rgb[0] = 0.0
	sb.rgb[1] = 0.2
	sb.rgb[2] = 0.0

	sb.snbname = fi.GetToken()

	for fi.IsEnd() == false {
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

func SCREEN(fi *EeTokens, sb *sunblk) {
	sb.rgb[0] = 0.0
	sb.rgb[1] = 0.2
	sb.rgb[2] = 0.0

	sb.snbname = fi.GetToken()

	for fi.IsEnd() == false {
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

func rmpdata(fi *EeTokens, rp *RRMP, _wp []MADO) {
	rp.ref = 0.0
	rp.grpx = 1.0

	rp.rgb[0] = 0.9
	rp.rgb[1] = 0.9
	rp.rgb[2] = 0.9

	rp.rmpname = fi.GetToken()
	rp.wallname = fi.GetToken()

	for fi.IsEnd() == false {
		NAME := fi.GetToken()
		if NAME[0] == ';' {
			break
		}

		if NAME == "-xyb" {
			rp.xb0 = fi.GetFloat()
			rp.yb0 = fi.GetFloat()
		} else if NAME == "-WH" {
			rp.Rw = fi.GetFloat()
			rp.Rh = fi.GetFloat()
		} else if NAME == "-ref" {
			rp.ref = fi.GetFloat()
		} else if NAME == "-grpx" {
			rp.grpx = fi.GetFloat()
		} else if NAME == "-rgb" {
			rp.rgb[0] = fi.GetFloat()
			rp.rgb[1] = fi.GetFloat()
			rp.rgb[2] = fi.GetFloat()
		} else {
			fmt.Printf("ERROR parameter----RMP: %s\n", NAME)
			os.Exit(1)
		}
	}

	rp.sumWD = 0
	for _, wp := range _wp {
		NAME := fi.GetToken()
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

		wp.winname = fi.GetToken()

		for fi.IsEnd() == false {
			NAME := fi.GetToken()
			if NAME[0] == ';' {
				break
			}

			if NAME == "-xyr" {
				wp.xr = fi.GetFloat()
				wp.yr = fi.GetFloat()
			} else if NAME == "-WH" {
				wp.Ww = fi.GetFloat()
				wp.Wh = fi.GetFloat()
			} else if NAME == "-ref" {
				wp.ref = fi.GetFloat()
			} else if NAME == "-grpx" {
				wp.grpx = fi.GetFloat()
			} else if NAME == "-rgb" {
				wp.rgb[0] = fi.GetFloat()
				wp.rgb[1] = fi.GetFloat()
				wp.rgb[2] = fi.GetFloat()
			} else {
				fmt.Printf("ERROR parameter----WD: %s\n", NAME)
				os.Exit(1)
			}
		}
	}
}

/*------------------------------------------------------------------*/
func rectdata(fi *EeTokens, obs *OBS) {
	obs.ref[0] = 0.0

	obs.rgb[0] = 0.7
	obs.rgb[1] = 0.7
	obs.rgb[2] = 0.7

	obs.obsname = fi.GetToken()

	for fi.IsEnd() == false {
		NAME := fi.GetToken()
		if NAME[0] == ';' {
			break
		}

		if NAME == "-xyz" {
			obs.x = fi.GetFloat()
			obs.y = fi.GetFloat()
			obs.z = fi.GetFloat()
		} else if NAME == "-WH" {
			obs.W = fi.GetFloat()
			obs.H = fi.GetFloat()
		} else if NAME == "-WaWb" {
			obs.Wa = fi.GetFloat()
			obs.Wb = fi.GetFloat()
		} else if NAME == "-ref" {
			obs.ref[0] = fi.GetFloat()
		} else if NAME == "-rgb" {
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
func cubdata(fi *EeTokens, obs *OBS) {
	for i := 0; i < 3; i++ {
		obs.ref[i] = 0.0
	}

	obs.rgb[0] = 0.7
	obs.rgb[1] = 0.7
	obs.rgb[2] = 0.7

	obs.obsname = fi.GetToken()

	for fi.IsEnd() == false {
		NAME := fi.GetToken()
		if NAME[0] == ';' {
			break
		}

		if NAME == "-xyz" {
			obs.x = fi.GetFloat()
			obs.y = fi.GetFloat()
			obs.z = fi.GetFloat()
		} else if NAME == "-WDH" {
			obs.W = fi.GetFloat()
			obs.D = fi.GetFloat()
			obs.H = fi.GetFloat()
		} else if NAME == "-Wa" {
			obs.Wa = fi.GetFloat()
		} else if NAME == "-ref0" {
			obs.ref[0] = fi.GetFloat()
		} else if NAME == "-ref1" {
			obs.ref[1] = fi.GetFloat()
		} else if NAME == "-ref2" {
			obs.ref[2] = fi.GetFloat()
		} else if NAME == "-ref3" {
			obs.ref[3] = fi.GetFloat()
		} else if NAME == "-rgb" {
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
func tridata(fi *EeTokens, obs *OBS) {
	obs.ref[0] = 0.0

	obs.rgb[0] = 0.7
	obs.rgb[1] = 0.7
	obs.rgb[2] = 0.7

	obs.obsname = fi.GetToken()

	for fi.IsEnd() == false {
		NAME := fi.GetToken()
		if NAME[0] == ';' {
			break
		}

		if NAME == "-xyz" {
			obs.x = fi.GetFloat()
			obs.y = fi.GetFloat()
			obs.z = fi.GetFloat()
		} else if NAME == "-WH" {
			obs.W = fi.GetFloat()
			obs.H = fi.GetFloat()
		} else if NAME == "-WaWb" {
			obs.Wa = fi.GetFloat()
			obs.Wb = fi.GetFloat()
		} else if NAME == "-ref" {
			obs.ref[0] = fi.GetFloat()
		} else if NAME == "-rgb" {
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
// 20170503 higuchi add
func dividdata(fi *EeTokens, monten *int, DE *float64) {
	var NAME string

	for fi.IsEnd() == false {
		NAME = fi.GetToken()
		if NAME[0] == ';' {
			break
		}

		if NAME == "DE" {
			var err error
			s := fi.GetToken()
			*DE, err = strconv.ParseFloat(s, 64)
			if err != nil {
				fmt.Printf("ERROR parameter----DIVID: %s\n", NAME)
			}
		} else if NAME == "MONT" {
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

	NAME = fi.GetToken()
}

func treedata(fi *EeTokens, treen *int, tree *[]TREE) {
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
		NAME = fi.GetToken()
		if NAME[0] == '*' {
			break
		}

		tred.treetype = NAME

		NAME = fi.GetToken()
		tred.treename = NAME

		if tred.treetype == "treeA" {
			for fi.IsEnd() == false {
				NAME = fi.GetToken()
				if NAME[0] == ';' {
					break
				}

				if NAME == "-xyz" {
					tred.x = fi.GetFloat()
					tred.y = fi.GetFloat()
					tred.z = fi.GetFloat()
				} else if NAME == "-WH1" {
					tred.W1 = fi.GetFloat()
					tred.H1 = fi.GetFloat()
				} else if NAME == "-WH2" {
					tred.W2 = fi.GetFloat()
					tred.H2 = fi.GetFloat()
				} else if NAME == "-WH3" {
					tred.W3 = fi.GetFloat()
					tred.H3 = fi.GetFloat()
				} else if NAME == "-WH4" {
					tred.W4 = fi.GetFloat()
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
func polydata(fi *EeTokens, polyn *int, poly *[]POLYGN) {
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
		NAME = fi.GetToken()
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

		polyp.polyd = fi.GetInt()
		polyp.P = make([]XYZ, polyp.polyd)

		polyp.polyname = fi.GetToken()
		polyp.wallname = fi.GetToken()

		for fi.IsEnd() == false {
			NAME = fi.GetToken()
			if NAME[0] == ';' {
				break
			}

			if NAME == "-xyz" {
				for i = 0; i < polyp.polyd; i++ {
					polyp.P[i].X = fi.GetFloat()
					polyp.P[i].Y = fi.GetFloat()
					polyp.P[i].Z = fi.GetFloat()
				}

			} else if NAME == "-rgb" {
				polyp.rgb[0] = fi.GetFloat()
				polyp.rgb[1] = fi.GetFloat()
				polyp.rgb[2] = fi.GetFloat()
			} else if NAME == "-ref" {
				polyp.ref = fi.GetFloat()
			} else if NAME == "-refg" {
				polyp.refg = fi.GetFloat()
			} else if NAME == "-grpx" {
				polyp.grpx = fi.GetFloat()
			} else {
				fmt.Printf("ERROR parameter----POLYGON: %s\n", NAME)
				os.Exit(1)
			}
		}
		(*polyn)++
	}
}

/*---------------------------------------------------------------------------*/
func bdpdata(fi *EeTokens, bdpn *int, bp *[]BBDP, Exsf *EXSFS) {

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
		NAME = fi.GetToken()
		if NAME[0] == '*' {
			break
		}

		if NAME != "BDP" {
			fmt.Printf("error BDP\n")
			os.Exit(1)
		}

		bbdp.bdpname = fi.GetToken()

		for fi.IsEnd() == false {
			NAME = fi.GetToken()
			if NAME[0] == ';' {
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

			NAME = fi.GetToken()
			if NAME[0] == '*' {
				break
			}

			if NAME == "SBLK" {
				sb.ref = 0.0
				sb.sbfname = fi.GetToken()

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
func obsdata(fi *EeTokens, obsn *int, obs *[]OBS) {
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
	for i = 0; i < Nobs; i++ {
		obsp = &(*obs)[i]

		NAME := fi.GetToken()
		if NAME[0] == '*' {
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

func InputCount(fi *EeTokens, key string) int {
	N := 0
	ad := fi.GetPos()

	for fi.IsEnd() == false {
		s := fi.GetToken()
		if s == "*" {
			break
		}

		N++

		for fi.IsEnd() == false {
			s = fi.GetToken()

			if s == key {
				break
			}
		}
	}

	fi.RestorePos(ad)
	return N
}

func SBLKCount(fi *EeTokens) int {
	N := 0
	ad := fi.GetPos()

	var s string
	for fi.IsEnd() == false {
		s = fi.GetToken()
		if s == "*" {
			break
		}

		if s == "SBLK" {
			N++
		}
	}

	fi.RestorePos(ad)
	return N
}

func RMPCount(fi *EeTokens) int {
	N := 0
	ad := fi.GetPos()

	var s string
	for fi.IsEnd() == false {
		s = fi.GetToken()
		if s == "*" {
			break
		}

		if s == "RMP" {
			N++
		}
	}

	fi.RestorePos(ad)

	return N
}

func WDCount(fi *EeTokens) int {
	N := 0
	ad := fi.GetPos()

	Flg := 0
	for fi.IsEnd() == false {
		s := fi.GetToken()

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

	fi.RestorePos(ad)

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

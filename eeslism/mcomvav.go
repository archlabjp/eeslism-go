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

/*  mcomvav.c  */

/*  OM用変風量コントローラ */

package eeslism

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

/* ------------------------------------------ */

/* 機器仕様入力　　　　　　*/

/*---- Satoh OMVAV  2010/12/16 ----*/
func OMVAVdata(s string, OMvavca *OMVAVCA) int {
	id := 0

	st := strings.Split(s, "=")
	if len(st) == 1 {
		OMvavca.Name = s
		OMvavca.Gmax = -999.0
		OMvavca.Gmin = -999.0
	} else {
		switch st[0] {
		case "Gmax":
			dt, err := strconv.ParseFloat(st[1], 64)
			if err != nil {
				id = 1
				break
			}
			OMvavca.Gmax = dt
		case "Gmin":
			dt, err := strconv.ParseFloat(st[1], 64)
			if err != nil {
				id = 1
				break
			}
			OMvavca.Gmin = dt
		default:
			id = 1
		}
	}

	return id
}

func CollTout(Tcin, G float64, Sd *RMSRF) float64 {
	var Kc, ca float64

	Wall := Sd.mw.wall
	if Wall.chrRinput == 'Y' {
		Kc = Sd.dblKc
	} else {
		Kc = Wall.Kc
	}

	return Sd.Tcole - (Sd.Tcole-Tcin)*math.Exp(-Kc*Sd.A/(ca*G))
}

func OMflowcalc(OMvav *OMVAV, Wd *WDAT) float64 {
	var Tcout float64
	var CollTout func(Tcin, G float64, Sd *RMSRF) float64
	var loop int

	G := 0.0
	//EPS := 0.00001
	dGp := 0.01
	LoopMax := 100

	if OMvav.Plist.Control != OFF_SW {
		omwall := OMvav.Omwall
		//Wall := omwall.mw.wall
		Tcin := Wd.T
		Tcoutset := omwall.rpnl.Toset
		G0 := OMvav.Cat.Gmin
		//dG := OMvav.Cat.Gmin * 0.001
		G2 := OMvav.Cat.Gmax

		/********************************************************/
		// 棟温度の計算（最小風量の場合）
		Tcin = Wd.T
		for i := 0; i < OMvav.Nrdpnl; i++ {
			Sd := OMvav.Rdpnl[i].sd[0]
			Tcout := CollTout(Tcin, G0, Sd)

			// 集熱器の入り口温度は上流集熱器の出口温度
			Tcin = Tcout
		}
		Tcoutmin := Tcout

		// 棟温度の計算（最大風量の場合）
		Tcin = Wd.T
		for i := 0; i < OMvav.Nrdpnl; i++ {
			Sd := OMvav.Rdpnl[i].sd[0]
			Tcout := CollTout(Tcin, G2, Sd)

			// 集熱器の入り口温度は上流集熱器の出口温度
			Tcin = Tcout
		}
		Tcoutmax := Tcout

		//fmt.Printf("Tcoutmin=%.2f Tcoutmax=%.2f\n", Tcoutmin, Tcoutmax)
		if Tcoutmin < Tcoutset {
			G = G0
		} else if Tcoutmax > Tcoutset {
			G = G2
		} else {
			// ニュートンラプソン法のループ
			for loop := 0; loop < LoopMax; loop++ {
				Tcin := Wd.T
				G := G0 + float64(loop)*dGp*(G2-G0)
				for i := 0; i < OMvav.Nrdpnl; i++ {
					Sd := OMvav.Rdpnl[i].sd[0]
					Tcout := CollTout(Tcin, G, Sd)

					// 集熱器の入り口温度は上流集熱器の出口温度
					Tcin = Tcout
				}

				FG := Tcoutset - Tcout

				if FG > 0.0 {
					break
				}
			}

			if loop == LoopMax {
				fmt.Printf("%s  風量が収束しませんでした。 G=%f\n", OMvav.Name, G)
			}
		}
		G = math.Min(math.Max(OMvav.Cat.Gmin, G), OMvav.Cat.Gmax)
	}

	OMvav.G = G

	return G
}

func OMvavControl(OMvav *OMVAV, Compnt []COMPNT, Ncompnt int) {
	colname := strings.Split(OMvav.Cmp.Omparm, "-")
	OMvav.Nrdpnl = len(colname)

	if len(Compnt) != Ncompnt {
		panic("OMvavControl: len(Compnt) != Ncompnt")
	}

	for j, name := range colname {
		for i := 0; i < len(Compnt); i++ {
			if name == Compnt[i].Name {
				OMvav.Rdpnl[j] = Compnt[i].Eqp.(*RDPNL)
				break
			}
		}
	}
}

func strCompcount(st string, key byte) int {
	count := 0

	for i := 0; i < len(st); i++ {
		if st[i] == key {
			count++
		}
	}

	return count
}

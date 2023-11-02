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

/*   mc_reflib.c                     */

package eeslism

import (
	"bufio"
	"io"
	"math"
	"strconv"
	"strings"
)

/*  圧縮式冷凍機定格特性入力    */

func Refcmpdat(frf io.Reader, Rfcmp *[]*RFCMP) {
	scanner := bufio.NewScanner(frf)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "*" {
			break
		}

		fields := strings.Fields(line)

		rfcmp := new(RFCMP)
		rfcmp.name = fields[0]
		rfcmp.cname = fields[1]

		for i := 0; i < 4; i++ {
			val, _ := strconv.ParseFloat(fields[i+2], 64)
			rfcmp.e[i] = val
		}
		for i := 0; i < 4; i++ {
			val, _ := strconv.ParseFloat(fields[i+6], 64)
			rfcmp.d[i] = val
		}
		for i := 0; i < 4; i++ {
			val, _ := strconv.ParseFloat(fields[i+10], 64)
			rfcmp.w[i] = val
		}
		val, _ := strconv.ParseFloat(fields[14], 64)
		rfcmp.Teo[0] = val
		val, _ = strconv.ParseFloat(fields[15], 64)
		rfcmp.Teo[1] = val
		val, _ = strconv.ParseFloat(fields[16], 64)
		rfcmp.Tco[0] = val
		val, _ = strconv.ParseFloat(fields[17], 64)
		rfcmp.Tco[1] = val
		val, _ = strconv.ParseFloat(fields[18], 64)
		rfcmp.Meff = val

		*Rfcmp = append(*Rfcmp, rfcmp)
	}
}

/* ----------------------------------- */

/*  冷凍機の蒸発温度と冷凍能力の一次式の係数  */

func Compca(e, d *[4]float64, EGex float64, Teo [2]float64, Ta float64, Ho, He *float64) {
	var Tc, Te float64
	var Qo [2]float64

	for i := 0; i < 2; i++ {
		Te = Teo[i]
		Tc = (d[0] + d[1]*Te + EGex*Ta) / (EGex - d[2] - d[3]*Te)
		Qo[i] = e[0] + e[1]*Te + (e[2]+e[3]*Te)*Tc
	}
	*He = (Qo[0] - Qo[1]) / (Teo[1] - Teo[0])
	*Ho = Qo[0] + *He*Teo[0]
}

/* ------------------------------------------------------------ */

/*  ヒ－トポンプの凝縮温度と冷凍能力の一次式の係数  */

func Compha(e, d *[4]float64, EGex float64, Tco [2]float64, Ta float64, Ho, He *float64) {
	var Tc, Te float64
	var Qo [2]float64

	for i := 0; i < 2; i++ {
		Tc = Tco[i]
		Te = (e[0] + e[2]*Tc + EGex*Ta) / (EGex - e[1] - e[3]*Tc)
		Qo[i] = d[0] + d[2]*Tc + (d[1]+d[3]*Tc)*Te
	}
	*He = (Qo[0] - Qo[1]) / (Tco[1] - Tco[0])
	*Ho = Qo[0] + *He*Tco[0]
}

/* --------------------------------------- */

/*  冷凍機／ヒ－トポンプの軸動力の計算　　 */

func Refpow(Rf *REFA, QP float64) float64 {
	var W, Te, Tc float64
	if math.Abs(QP) > 1.0 {
		if Rf.Chmode == COOLING_SW {
			Te = QP/(Rf.Cat.cool.eo*Rf.cG) + Rf.Tin
			Tc = (QP - Rf.c_e[0] - Rf.c_e[1]*Te) / (Rf.c_e[2] + Rf.c_e[3]*Te)
			W = Rf.c_w[0] + Rf.c_w[1]*Te + Rf.c_w[2]*Tc + Rf.c_w[3]*Te*Tc
		} else if Rf.Chmode == HEATING_SW {
			Tc = QP/(Rf.Cat.heat.eo*Rf.cG) + Rf.Tin
			Te = (QP - Rf.h_d[0] - Rf.h_d[2]*Tc) / (Rf.h_d[1] + Rf.h_d[3]*Tc)
			W = Rf.h_w[0] + Rf.h_w[1]*Te + Rf.h_w[2]*Tc + Rf.h_w[3]*Te*Tc
		}

		Rf.Te = Te
		Rf.Tc = Tc
	} else {
		W = 0.0
	}

	return W
}

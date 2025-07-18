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
	"math"
	"os"
	"strings"
)

// 圧縮式冷凍機定格特性入力
// reflist.efl ファイルから読み取ります。
func Refcmpdat() []*RFCMP {
	frf, err := os.Open("reflist.efl")
	if err != nil {
		Eprint(" file ", "reflist.efl")
	}
	defer frf.Close()

	return _Refcmpdat(frf)
}

func _Refcmpdat(frf *os.File) []*RFCMP {
	Rfcmp := make([]*RFCMP, 0)
	scanner := bufio.NewScanner(frf)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "*" {
			break
		}

		fields := strings.Fields(line)

		rfcmp := NewRFCMP()
		rfcmp.name = fields[0]
		rfcmp.cname = fields[1]
		for i := 0; i < 4; i++ {
			rfcmp.e[i], _ = readFloat(fields[i+2])
		}
		for i := 0; i < 4; i++ {
			rfcmp.d[i], _ = readFloat(fields[i+6])
		}
		for i := 0; i < 4; i++ {
			rfcmp.w[i], _ = readFloat(fields[i+10])
		}
		rfcmp.Teo[0], _ = readFloat(fields[14])
		rfcmp.Teo[1], _ = readFloat(fields[15])
		rfcmp.Tco[0], _ = readFloat(fields[16])
		rfcmp.Tco[1], _ = readFloat(fields[17])
		rfcmp.Meff, _ = readFloat(fields[18])

		Rfcmp = append(Rfcmp, rfcmp)
	}

	return Rfcmp
}

func NewRFCMP() *RFCMP {
	Rf := new(RFCMP)
	Rf.cname = ""
	for j := 0; j < 4; j++ {
		Rf.d[j] = 0.0
		Rf.e[j] = 0.0
		Rf.w[j] = 0.0
	}
	Rf.Meff = 0.0
	Rf.name = ""
	for j := 0; j < 2; j++ {
		Rf.Tco[j] = 0.0
		Rf.Teo[j] = 0.0
	}
	return Rf
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

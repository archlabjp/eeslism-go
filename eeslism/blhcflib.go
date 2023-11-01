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

/*    bhcflib.c            */

package eeslism

import (
	"fmt"
	"math"
)

const (
	Alidmy  = 9.3
	ALITOLE = 1.e-5
)

/* -------------------------------------------------------------- */

/*    熱伝達率に関する計算  */

func Htrcf(alc, alo *float64, alotype AloType, Exs []*EXSF, Tr float64, N int, alr []float64, _Sd []RMSRF, RMmrk *rune, Wd *WDAT) {
	var n int
	var alic float64
	var hc *float64
	//var dT float64

	if DEBUG {
		fmt.Printf("Htrcf Start\n")
	}

	for n = 0; n < N; n++ {
		Sd := &_Sd[n]

		if DEBUG {
			fmt.Printf("n=%d name=%s\n", n, Sd.Name)
		}

		// 室内側対流熱伝達率の計算
		alic = -1.0

		if alc != nil {
			if *alc >= 0.01 {
				alic = *alc
			}
		}

		if DEBUG {
			fmt.Printf("alic=%f\n", alic)
		}

		hc = nil

		if DEBUG {
			fmt.Printf("Sd->alicsch\n")
		}

		if Sd.alicsch != nil {
			hc = Sd.alicsch
			if *hc >= 0.00 {
				alic = *hc
			}
		}

		if DEBUG {
			fmt.Printf("hc set end\n")
		}
		if alic < 0.0 {
			dT := Sd.Ts - Tr
			switch Sd.ble {
			case 'F', 'f':
				if dT > 0 {
					alic = alcvup(dT)
				} else {
					alic = alcvdn(dT)
				}
			case 'R', 'c':
				if math.Abs(dT) <= 1.0e-3 {
					alic = 0.0
				} else if dT < 0 {
					alic = alcvup(dT)
				} else {
					alic = alcvdn(dT)
				}
			default:
				alic = alcvh(dT)
			}
		}

		if DEBUG {
			fmt.Printf("----- Htrcf n=%d mrk=%c alic=%f  Sd->alic=%f\n",
				n, Sd.mrk, alic, Sd.alic)
		}

		if math.Abs(alic-Sd.alic) >= ALITOLE || Sd.mrk == '*' || Sd.PCMflg {
			*RMmrk = '*'
			Sd.mrk = '*'
			Sd.alic = alic

			switch Sd.ble {
			case BLE_Window, BLE_ExternalWall, BLE_Floor, BLE_Roof:
				Sd.alo = *Exs[Sd.exs].Alo
			default:
				Sd.alo = Alidmy
			}
		}

		if Sd.mrk == '*' {
			if Sd.alirsch == nil {
				Sd.alir = alr[n*N+n]
			} else if *Sd.alirsch > 0.0 {
				Sd.alir = *Sd.alirsch
			} else {
				Sd.alir = 0.0
			}

			Sd.ali = alic + Sd.alir
		}

		if DEBUG {
			fmt.Printf("----- Htrcf n=%2d ble=%c Ts=%.1f Tr=%.1f alic=%.3f alir=%.3f rmname=%s\n",
				n, Sd.ble, Sd.Ts, Tr, Sd.alic, Sd.alir, Sd.room.Name)
		}

		//if dayprn && Ferr != nil {
		fmt.Fprintf(Ferr, "----- Htrcf n=%2d ble=%c Ts=%f Tr=%f alic=%f alir=%f rmname=%s\n",
			n, Sd.ble, Sd.Ts, Tr, Sd.alic, Sd.alir, Sd.room.Name)
		//}
	}
}

/*-----------------------------------------------------*/
// 屋外側熱伝達率の計算
func alov(Exs *EXSF, Wd *WDAT) float64 {
	var u float64

	Wv := Wd.Wv
	Wdre := -180.0 + 360.0/16.0*Wd.Wdre

	Wadiff := math.Abs(Exs.Wa - Wdre)
	Wadiff = math.Mod(Wadiff, 360.0)
	if Wadiff < 45.0 {
		if Wv <= 2.0 {
			u = 0.5
		} else {
			u = 0.25 * Wv
		}
	} else {
		u = 0.3 + 0.05*Wv
	}

	return 3.5 + 5.6*u
}

/* ---------------------------------------- */

/*  室内表面対流熱伝達率の計算     */

func alcvup(dT float64) float64 {
	return 2.18 * math.Pow(math.Abs(dT), 0.31)
}

func alcvdn(dT float64) float64 {
	return 0.138 * math.Pow(math.Abs(dT), 0.25)
}

func alcvh(dT float64) float64 {
	return 1.78 * math.Pow(math.Abs(dT), 0.32)
}

/* --------------------------------------- */

/*  室内表面間放射熱伝達率の計算  */

func Radcf0(Tsav float64, alrbold *float64, N int, Sd []RMSRF, W, alr []float64) {
	var n int
	var alir, TA float64

	TA = Tsav + 273.15
	alir = 4.0 * Sgm * math.Pow(TA, 3.0)

	/*****/
	if DEBUG {
		fmt.Printf("----- Radcf0   alir=%f  alrbold=%f ALITOLE\n", alir, *alrbold)
	}
	/*****/

	if math.Abs(alir-*alrbold) >= ALITOLE {
		*alrbold = alir

		for n = 0; n < N; n++ {
			Sd[n].mrk = '*'
		}

		for n = 0; n < N*N; n++ {
			alr[n] = alir * math.Abs(W[n])

			/*****fmt.Printf("----- Radcf0  n=%d alr=%f\n",n, alr[n])
			 *****/
		}
	}
}

/* ------------------------------------------- */

/*  放射伝達係数の計算  */

func radex(N int, Sd []RMSRF, F, W []float64) {
	wk := make([]float64, N*N)
	Ff := make([]float64, N*N)

	for l, n := 0, 0; n < N; n++ {
		for j := 0; j < N; j++ {
			wk[l] = -F[l] * (1.0 - Sd[j].Ei) / Sd[j].Ei
			Ff[l] = -F[l]
			l++
		}
		nn := n*N + n
		wk[nn] += 1.0 / Sd[n].Ei
		Ff[nn] += 1.0
	}

	Matinv(wk, N, N, "<radex>")

	for i := 0; i < N; i++ {
		for j := 0; j < N; j++ {
			c := 0.0
			for k := 0; k < N; k++ {
				c += wk[N*i+k] * Ff[N*k+j]
			}
			W[N*i+j] = c
		}
	}
}

/* ------------------------------------------- */

/* 形態係数の近似計算　　*/

func formfaprx(N int, Aroom float64, Sd []RMSRF, F []float64) {
	var i, n int

	for n = 0; n < N; n++ {
		F[n] = Sd[n].A / Aroom

		for i = 0; i < N; i++ {
			F[i*N+n] = F[n]
		}
	}
}

/* ------------------------------------------- */

/*  短波長放射吸収係数 */

func Radshfc(N int, FArea, Aroom float64, Sd0 []RMSRF, tfsol, eqcv float64, Rmname string, fsolm *float64) {
	var Sumsrg2, dblTemp, Srgchk float64
	Room := Sd0[0].room

	Sumsrg2 = 0.0

	// tfsol:定義済みの部位日射吸収係数の合計値
	for n := 0; n < N; n++ {
		Sd := &Sd0[n]
		Sd.eqrd = (1.0 - eqcv) * Sd.A / Aroom

		dblTemp = math.Max((1.0-tfsol), 0.0) * Sd.A / Aroom
		if Sd.fsol != nil {
			v := *(Sd.fsol) + dblTemp
			Sd.srg = v
			Sd.srh = v
			Sd.srl = v
			Sd.sra = v
		} else {
			Sd.srg = dblTemp
			Sd.srh = dblTemp
			Sd.srl = dblTemp
			Sd.sra = dblTemp
		}

		Sd.srg2 = Sd.srg
		if Sd.RStrans {
			Sd.srg2 = 0.0
		}
		if Sd.tnxt > 0.0 {
			Sd.srg2 = Sd.srg * (1.0 - Sd.tnxt)
		}
		Sumsrg2 += Sd.srg2
	}

	Room.Srgm2 = 0.0
	if fsolm != nil {
		Room.Srgm2 = *fsolm
	}

	//  各種部位の吸収係数のチェック
	Srgchk = 0.0

	// 家具の日射吸収割合を加算
	if fsolm != nil {
		Srgchk += *fsolm
		Sumsrg2 += *fsolm
	}

	for n := 0; n < N; n++ {
		Sd := &Sd0[n]
		Srgchk += Sd.srg
	}

	if math.Abs(Srgchk-1.0) > 1.0e-3 {
		fmt.Printf("xxxxx (%s)  室内部位への日射吸収比率の合計が不適 %.3f (本来、1となるべき)\n", Rmname, Srgchk)
	}

	if tfsol > 1.0 {
		fmt.Printf("xxxxx (%s)  室内部位への日射吸収比率を指定したものだけで合計が不適 %.3f (本来、1未満となるべき)\n", Rmname, tfsol)
	}

	// 最終日射吸収比率の計算（Sumsrg2で基準化）
	for n := 0; n < N; n++ {
		Sd := &Sd0[n]
		Sd.srg2 /= Sumsrg2
	}

	if fsolm != nil {
		Room.Srgm2 /= Sumsrg2
	}
}

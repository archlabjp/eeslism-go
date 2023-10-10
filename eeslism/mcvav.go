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

/*  mcvav.c  */

/*  VAVコントローラ */

package eeslism

import (
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
)

/* ------------------------------------------ */

/* 機器仕様入力　　　　　　*/

/*---- Satoh Debug VAV  2000/10/30 ----*/

func VAVdata(cattype EqpType, s string, vavca *VAVCA) int {
	var st string
	var dt float64
	var id int

	if cattype == VAV_TYPE {
		vavca.Type = VAV_PDT
	} else if cattype == VWV_TYPE {
		vavca.Type = VWV_PDT
	}

	if st = strings.Split(s, "=")[1]; st == "" {
		vavca.Name = s
		vavca.Gmax = -999.0
		vavca.Gmin = -999.0
		vavca.dTset = -999.0
	} else {
		dt, _ = strconv.ParseFloat(st, 64)

		switch s {
		case "Gmax":
			vavca.Gmax = dt
		case "Gmin":
			vavca.Gmin = dt
		case "dTset":
			vavca.dTset = dt
		default:
			id = 1
		}
	}

	return id
}

func VWVint(Nvav int, VAVs []VAV, Ncompnt int, Compn []COMPNT) {
	for i := 0; i < Nvav; i++ {
		vav := &VAVs[i]
		vav.Hcc = nil
		vav.Hcld = nil
		vav.Mon = '-'

		if vav.Cat.Type == VWV_PDT {
			if vav.Cmp.Hccname != "" {
				vav.Hcc = hccptr('c', vav.Cmp.Hccname, Ncompnt, Compn, &vav.Mon).(*HCC)
			} else if vav.Cmp.Rdpnlname != "" {
				vav.Rdpnl = rdpnlptr(vav.Cmp.Rdpnlname, Ncompnt, Compn)
				if vav.Rdpnl != nil {
					vav.Mon = 'f'
				}
			}

			if vav.Mon == '-' {
				vav.Hcld = hccptr('h', vav.Cmp.Hccname, Ncompnt, Compn, &vav.Mon).(*HCLOAD)
			}

			if vav.Mon == '-' {
				s := fmt.Sprintf("VWV(%s)=%s", vav.Name, vav.Cmp.Hccname)
				fmt.Printf("xxxxxxxxx %s xxxxxxxxxxx\n", s)
			}
		}
	}
}

/*  特性式の係数  */
/*---- Satoh Debug VAV  2000/11/8 ----*/
/*********************/
func VAVcfv(Nvav int, vav []VAV) {
	for i := 0; i < Nvav; i++ {
		v := &vav[i]
		Eo := v.Cmp.Elouts

		if v.Cmp.Control != OFF_SW && Eo[0].Control != OFF_SW {
			if v.Cat.Gmax < 0.0 {
				Err := fmt.Sprintf("Name=%s  Gmax=%.5g", v.Name, v.Cat.Gmax)
				Eprint("VAVcfv", Err)
			}
			if v.Cat.Gmin < 0.0 {
				Err := fmt.Sprintf("Name=%s  Gmin=%.5g", v.Name, v.Cat.Gmin)
				Eprint("VAVcfv", Err)
			}

			if v.Count == 0 {
				v.G = Eo[0].G
				v.CG = Spcheat(Eo[0].Fluid) * v.G

				Eo[0].Coeffo = v.CG
				Eo[0].Co = 0.0
				Eo[0].Coeffin[0] = -v.CG
			} else {
				Eo[0].Coeffo = 1.0
				Eo[0].Co = 0.0
				Eo[0].Coeffin[0] = -1.0
			}

			if v.Cat.Type == VAV_PDT {
				Eo[1].Coeffo = 1.0
				Eo[1].Co = 0.0
				Eo[1].Coeffin[0] = -1.0
			}
		}
	}
}

/************************:/
/* ------------------------------------------ */

/* VAVコントローラ再熱部分の計算 */
/*---- Satoh Debug VAV  2000/11/27 ----*/
/*******************/
func VAVene(Nvav int, vav []VAV, VAVrest *int) {
	var i, rest int
	var elo *ELOUT
	var Tr, Go, dTset float64

	for i = 0; i < Nvav; i++ {
		rest = 0

		elo = vav[i].Cmp.Elouts[0]
		vav[i].Tin = elo.Elins[0].Sysvin

		if vav[i].Cmp.Control != OFF_SW && elo.Control != OFF_SW {
			Go = vav[i].G
			vav[i].Tout = elo.Sysv

			if vav[i].Cat.Type == VAV_PDT {
				Tr = vav[i].Cmp.Elouts[0].Emonitr.Sysv

				vav[i].Q = Spcheat(elo.Fluid) * Go * (vav[i].Tout - Tr)

				if math.Abs(vav[i].Tin-Tr) > 1.0e-3 {
					vav[i].G = (vav[i].Tout - Tr) / (vav[i].Tin - Tr) * Go
				} else {
					vav[i].G = vav[i].Cat.Gmin
				}
			} else {
				if vav[i].Mon == 'c' && vav[i].Count < VAVCountMAX-1 {
					vav[i].Qrld = -vav[i].Hcc.Qt
				} else if vav[i].Mon == 'f' && vav[i].Count < VAVCountMAX-1 {
					vav[i].Qrld = -vav[i].Rdpnl.Q
				} else if vav[i].Mon == 'h' {
					vav[i].Qrld = vav[i].Hcld.Qt
				}

				vav[i].Q = vav[i].Qrld

				dTset = vav[i].Cat.dTset

				if vav[i].Mon == 'h' && dTset <= 0.0 {
					fmt.Printf("<VAVene> VWV SetDifferencialTemp=%.1f\n", vav[i].Cat.dTset)
				}

				if vav[i].Chmode == COOLING_SW {
					dTset = -dTset
				}

				if vav[i].Mon == 'h' || vav[i].Mon == 'f' {
					vav[i].G = vav[i].Q / (Spcheat(elo.Fluid) * dTset)
				} else if vav[i].Mon == 'c' {
					vav[i].G = FNVWVG(&vav[i])
				}
			}

			elo.Control = ON_SW
			elo.Sysld = 'n'

			if vav[i].Mon != 'h' {
				elo.Emonitr.Control = ON_SW
			}

			rest = chvavswreset(vav[i].Q, vav[i].Chmode, &vav[i])

			if rest == 1 || math.Abs(vav[i].G-Go) > 1.0e-5 {
				(*VAVrest)++
			}
		} else {
			vav[i].Q = 0.0
			vav[i].G = vav[i].Cat.Gmin

			if vav[i].Count == 0 {
				(*VAVrest)++
			}
		}
	}
}

func VAVcountreset(Nvav int, VAVs []VAV) {
	for i := 0; i < Nvav; i++ {
		v := &VAVs[i]
		v.Count = 0
	}
}

func VAVcountinc(Nvav int, VAVs []VAV) {
	for i := 0; i < Nvav; i++ {
		v := &VAVs[i]
		v.Count++
	}
}

func vavswptr(key []string, VAV *VAV, vptr *VPTR) int {
	err := 0

	if key[1] == "chmode" {
		vptr.Ptr = &VAV.Chmode
		vptr.Type = SW_CTYPE
	} else if key[1] == "control" {
		vptr.Ptr = &VAV.Cmp.Elouts[0].Control
		vptr.Type = SW_CTYPE
	} else {
		err = 1
	}

	return err
}

func chvavswreset(Qload float64, chmode rune, vav *VAV) int {
	if (chmode == HEATING_SW && Qload < 0.0) ||
		(chmode == COOLING_SW && Qload > 0.0) {
		vav.G = vav.Cat.Gmin
		return 1
	} else {
		if vav.G < vav.Cat.Gmin {
			vav.G = vav.Cat.Gmin
			return 1
		} else if vav.G > vav.Cat.Gmax {
			vav.G = vav.Cat.Gmax
			return 1
		} else {
			return 0
		}
	}
}

func vavprint(fo io.Writer, id, Nvav int, VAVs []VAV) {
	switch id {
	case 0:
		if Nvav > 0 {
			fmt.Fprintf(fo, "%s %d\n", VAV_TYPE, Nvav)
		}
		for i := 0; i < Nvav; i++ {
			vav := &VAVs[i]
			fmt.Fprintf(fo, " %s 1 2\n", vav.Name)
		}
	case 1:
		for i := 0; i < Nvav; i++ {
			vav := &VAVs[i]
			fmt.Fprintf(fo, "%s_c c c %s_G m f\n", vav.Name, vav.Name)
		}
	default:
		for i := 0; i < Nvav; i++ {
			vav := &VAVs[i]
			fmt.Fprintf(fo, "%c %6.4g\n", vav.Cmp.Elouts[0].Control, vav.Cmp.Elouts[0].G)
		}
	}
}

/* VWVシステムの流量計算 */
func FNVWVG(VWV *VAV) float64 {
	var Wa, Wwd, Tain, Twin, Q, KA, et, F, emax, emin, Gwd float64
	var A, B float64
	var h *HCC
	var i int

	h = VWV.Hcc
	Wa = h.cGa
	Q = VWV.Q
	KA = h.Cat.KA
	Tain = h.Tain
	Twin = VWV.Tin
	A = VWV.Cat.Gmin
	B = VWV.Cat.Gmax
	Gwd = (A + B) / 2.0
	Wwd = Gwd * Cw

	et = -Q / (Wa * (Tain - Twin))
	emin = FNhccet(Wa, Cw*A, KA)
	emax = FNhccet(Wa, Cw*B, KA)

	if et > emax {
		return VWV.Cat.Gmax
	} else if et < emin {
		return VWV.Cat.Gmin
	}

	for i = 0; i < 30; i++ {
		Wwd = Gwd * Cw
		F = FNhccet(Wa, Wwd, KA) - et
		if math.Abs(F) < 1.0e-5 {
			return Gwd
		} else if F > 0.0 {
			B = Gwd
		} else if F < 0.0 {
			A = Gwd
		}

		Gwd = (A + B) / 2.0
	}

	fmt.Println("xxxxxx FNVWVG  収束せず")
	fmt.Println(Gwd)

	return Gwd
}

func FNFd(Wa, Ww, KA float64) float64 {
	B := (1.0 - Wa/Ww) * KA / Wa
	exB := math.Exp(-B)
	ex2B := math.Exp(-2.0 * B)
	Ww2 := math.Pow(Ww, 2.0)

	d := Ww * math.Pow(Ww-Wa*exB, 2.0)
	n := (Ww2+Ww+KA)*ex2B - (Ww*(Ww2+Wa)-Wa*KA)*exB

	return n / d
}

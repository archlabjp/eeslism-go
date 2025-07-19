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
	"errors"
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
		vavca.Gmax = FNAN
		vavca.Gmin = FNAN
		vavca.dTset = FNAN
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

func VWVint(VAVs []*VAV, Compn []*COMPNT) {
	for _, vav := range VAVs {
		vav.Hcc = nil
		vav.Hcld = nil
		vav.Mon = '-'

		if vav.Cat.Type == VWV_PDT {
			if vav.Cmp.Hccname != "" {
				vav.Hcc = hccptr('c', vav.Cmp.Hccname, Compn, &vav.Mon).(*HCC)
			} else if vav.Cmp.Rdpnlname != "" {
				vav.Rdpnl = rdpnlptr(vav.Cmp.Rdpnlname, Compn)
				if vav.Rdpnl != nil {
					vav.Mon = 'f'
				}
			}

			if vav.Mon == '-' {
				vav.Hcld = hccptr('h', vav.Cmp.Hccname, Compn, &vav.Mon).(*HCLOAD)
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

//	+-----+ ---> [OUT 1] 空気 or 温水温度
//
// /  | VAV |
//
//	+-----+ ---> [OUT 2] 湿度 (VAV_PDTのみ)
func VAVcfv(vav []*VAV) {
	for _, v := range vav {
		Eo1 := v.Cmp.Elouts[0]

		if v.Cmp.Control != OFF_SW && Eo1.Control != OFF_SW {
			if v.Cat.Gmax < 0.0 {
				Err := fmt.Sprintf("Name=%s  Gmax=%.5g", v.Name, v.Cat.Gmax)
				Eprint("VAVcfv", Err)
			}
			if v.Cat.Gmin < 0.0 {
				Err := fmt.Sprintf("Name=%s  Gmin=%.5g", v.Name, v.Cat.Gmin)
				Eprint("VAVcfv", Err)
			}

			if v.Count == 0 {
				v.G = Eo1.G
				v.CG = Spcheat(Eo1.Fluid) * v.G

				Eo1.Coeffo = v.CG
				Eo1.Co = 0.0
				Eo1.Coeffin[0] = -v.CG
			} else {
				Eo1.Coeffo = 1.0
				Eo1.Co = 0.0
				Eo1.Coeffin[0] = -1.0
			}

			if v.Cat.Type == VAV_PDT {
				Eo2 := v.Cmp.Elouts[1]
				Eo2.Coeffo = 1.0
				Eo2.Co = 0.0
				Eo2.Coeffin[0] = -1.0
			}
		}
	}
}

/************************:/
/* ------------------------------------------ */

/* VAVコントローラ再熱部分の計算 */
/*---- Satoh Debug VAV  2000/11/27 ----*/
/*******************/
func VAVene(vav []*VAV, VAVrest *int) {
	var rest int
	var elo *ELOUT
	var Tr, Go, dTset float64

	for _, v := range vav {
		rest = 0

		elo = v.Cmp.Elouts[0]
		v.Tin = elo.Elins[0].Sysvin

		if v.Cmp.Control != OFF_SW && elo.Control != OFF_SW {
			Go = v.G
			v.Tout = elo.Sysv

			if v.Cat.Type == VAV_PDT {
				Tr = v.Cmp.Elouts[0].Emonitr.Sysv

				v.Q = Spcheat(elo.Fluid) * Go * (v.Tout - Tr)

				if math.Abs(v.Tin-Tr) > 1.0e-3 {
					v.G = (v.Tout - Tr) / (v.Tin - Tr) * Go
				} else {
					v.G = v.Cat.Gmin
				}
			} else {
				if v.Mon == 'c' && v.Count < VAVCountMAX-1 {
					v.Qrld = -v.Hcc.Qt
				} else if v.Mon == 'f' && v.Count < VAVCountMAX-1 {
					v.Qrld = -v.Rdpnl.Q
				} else if v.Mon == 'h' {
					v.Qrld = v.Hcld.Qt
				}

				v.Q = v.Qrld

				dTset = v.Cat.dTset

				if v.Mon == 'h' && dTset <= 0.0 {
					fmt.Printf("<VAVene> VWV SetDifferencialTemp=%.1f\n", v.Cat.dTset)
				}

				if v.Chmode == COOLING_SW {
					dTset = -dTset
				}

				if v.Mon == 'h' || v.Mon == 'f' {
					v.G = v.Q / (Spcheat(elo.Fluid) * dTset)
				} else if v.Mon == 'c' {
					v.G = FNVWVG(v)
				}
			}

			elo.Control = ON_SW
			elo.Sysld = 'n'

			if v.Mon != 'h' {
				elo.Emonitr.Control = ON_SW
			}

			rest = chvavswreset(v.Q, v.Chmode, v)

			if rest == 1 || math.Abs(v.G-Go) > 1.0e-5 {
				(*VAVrest)++
			}
		} else {
			v.Q = 0.0
			v.G = v.Cat.Gmin

			if v.Count == 0 {
				(*VAVrest)++
			}
		}
	}
}

func (eqsys *EQSYS) VAVcountreset() {
	for _, v := range eqsys.Vav {
		v.Count = 0
	}
}

func (eqsys *EQSYS) VAVcountinc() {
	for _, v := range eqsys.Vav {
		v.Count++
	}
}

// VAVスイッチのポインターを作成します
func vavswptr(key []string, VAV *VAV) (VPTR, error) {
	if key[1] == "chmode" {
		return VPTR{
			Ptr:  &VAV.Chmode,
			Type: SW_CTYPE,
		}, nil
	} else if key[1] == "control" {
		return VPTR{
			Ptr:  &VAV.Cmp.Elouts[0].Control,
			Type: SW_CTYPE,
		}, nil

	}
	return VPTR{}, errors.New("vavswptr error")
}

func chvavswreset(Qload float64, chmode ControlSWType, vav *VAV) int {
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

func vavprint(fo io.Writer, id int, VAVs []*VAV) {
	switch id {
	case 0:
		if len(VAVs) > 0 {
			fmt.Fprintf(fo, "%s %d\n", VAV_TYPE, len(VAVs))
		}
		for _, vav := range VAVs {
			fmt.Fprintf(fo, " %s 1 2\n", vav.Name)
		}
	case 1:
		for _, vav := range VAVs {
			fmt.Fprintf(fo, "%s_c c c %s_G m f\n", vav.Name, vav.Name)
		}
	default:
		for _, vav := range VAVs {
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

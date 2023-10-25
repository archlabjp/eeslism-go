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

/*  VALV */

package eeslism

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

/* ------------------------------------------ */

func Valvcountreset(Valv []VALV) {
	for i := range Valv {
		Valv[i].Count = 0
	}
}

/***********************************************/

func Valvcountinc(Valv []VALV) {
	for i := range Valv {
		Valv[i].Count++
	}
}

// 通常はバルブの上流の流量に比率を乗じるが、基準となるOMvavが指定されている場合には、この流量に対する比率とする
// OMvavが指定されているときだけの対応
func Valvinit(Valv []VALV, Mpath []*MPATH) {
	for k := range Valv {
		if Valv[k].Cmp.MonPlistName != "" {
			for _, mpath := range Mpath {
				for _, Plist := range mpath.Plist {
					if Valv[k].Cmp.MonPlistName == Plist.Plistname {
						Valv[k].MonPlist = Plist
						Valv[k].MGo = &Plist.G

						if Valv[k].Cmb != nil {
							CValv := Valv[k].Cmb.Eqp.(*VALV)
							CValv.MonPlist = Plist
							CValv.MGo = &Plist.G
						}
						break
					}
				}
			}
		} else {
			Pelm := Valv[k].Plist.Pelm[len(Valv[k].Plist.Pelm)-1]
			Valv[k].MGo = Pelm.Cmp.Elouts[0].Lpath.Go
			Valv[k].MonPlist = Pelm.Cmp.Elouts[0].Lpath

			if Valv[k].Cmb != nil {
				CValv := Valv[k].Cmb.Eqp.(*VALV)
				CValv.MonPlist = Valv[k].MonPlist
				CValv.MGo = &Valv[k].MonPlist.G
			}
		}
	}
}

func Valvene(Valv []VALV, Valvreset *int) {
	var etype EqpType
	var T1, T2 float64
	var Vcb *VALV
	var r float64

	for i := range Valv {
		etype = Valv[i].Cmp.Eqptype
		if etype == TVALV_TYPE && Valv[i].Org == 'y' {
			if Valv[i].Mon.Elouts[0].Control != OFF_SW {
				T1 = *Valv[i].Tin
				Vcb = Valv[i].Cmb.Eqp.(*VALV)
				T2 = *Vcb.Tin

				if math.Abs(*Valv[i].Tout-*Valv[i].Tset) >= 1.0e-3 && math.Abs(T1-T2) >= 1.0e-3 {
					r = (*Valv[i].Tset - T2) / (T1 - T2)
					r = math.Min(1.0, math.Max(r, 0.0))
					Valv[i].X = r
					Vcb.X = 1.0 - r

					Valv[i].Plist.Gcalc = r * *Valv[i].MGo
					Vcb.Plist.Gcalc = (1.0 - r) * *Valv[i].MGo
					(*Valvreset)++

					if DEBUG {
						fmt.Printf("<Valvene> Valvname=%s G=%f\n", Valv[i].Name, Valv[i].Plist.G)
						fmt.Printf("    T1=%.1f T2=%.1f Tset=%.1f\n", T1, T2, *Valv[i].Tset)
					}
				} else {
					Valv[i].Plist.Gcalc = Valv[i].X * *Valv[i].MGo
					Vcb.Plist.Gcalc = (1.0 - Valv[i].X) * *Valv[i].MGo
					(*Valvreset)++
				}
			}
		}
	}
}

/************************************************************************/

func ValvControl(fi *EeTokens, Compnt []*COMPNT, Schdl *SCHDL, Simc *SIMCONTL, Wd *WDAT, vptr *VPTR) {
	var s string
	var Valv, Vb *VALV
	var Vc *COMPNT
	var k, i int
	var elins *ELIN
	var Pelm *PELM
	var err error

	Vb = nil
	s = fi.GetToken()

	ad := fi.GetPos()

	Vc = Compntptr(s, Compnt)
	if Vc == nil {
		Eprint("<CONTRL>", s)
	}

	vptr.Ptr = &Vc.Control
	vptr.Type = SW_CTYPE

	Valv = Vc.Eqp.(*VALV)
	Valv.Org = 'y'
	for fi.IsEnd() == false {
		s = fi.GetToken()

		if strings.HasPrefix(s, ";") {
			fi.RestorePos(ad)
			break
		}

		if s == "-init" {
			s = fi.GetToken()
			ad = fi.GetPos()

			if k, err = idsch(s, Schdl.Sch, ""); err == nil {
				Valv.Xinit = &Schdl.Val[k]
			} else {
				Valv.Xinit = envptr(s, Simc, Compnt, Wd, nil)
			}
		} else if s == "-Tout" {
			s = fi.GetToken()
			ad = fi.GetPos()
			if k, err = idsch(s, Schdl.Sch, ""); err == nil {
				Valv.Tset = &Schdl.Val[k]
			} else {
				Valv.Tset = envptr(s, Simc, Compnt, Wd, nil)
			}

			Pelm = Valv.Plist.Pelm[len(Valv.Plist.Pelm)-1]
			Valv.Mon = Pelm.Cmp
			Valv.Tout = &Valv.Mon.Elouts[0].Sysv
			Valv.MGo = &Pelm.Cmp.Elouts[0].Lpath.G
			Vb = Valv.Cmb.Eqp.(*VALV)
			Vb.MGo = &Pelm.Cmp.Elouts[0].Lpath.G

			if Valv.Plist.Go == nil {
				Vb = Valv.Cmb.Eqp.(*VALV)
				Valv.Plist.Go = Valv.Mon.Elouts[0].Lpath.Go
				Vb.Plist.Go = Valv.Mon.Elouts[0].Lpath.Go
			}
		}
	}

	Pelm = Valv.Plist.Pelm[len(Valv.Plist.Pelm)-1]
	elins = Pelm.Cmp.Elins[0]

	for i = 0; i < Pelm.Cmp.Nin; i++ {
		if elins.Lpath.Valv.Name == Valv.Name {
			Valv.Tin = &elins.Sysvin
		} else {
			Vb.Tin = &elins.Sysvin
		}
	}
}

// バルブの内部変数へのポインタ
func valv_vptr(key []string, Valv *VALV) (VPTR, error) {
	var err error
	var vptr VPTR

	if strings.Compare(key[1], "value") == 0 {
		vptr.Ptr = &Valv.X
		vptr.Type = VAL_CTYPE
	} else {
		err = errors.New("'value' is expected")
	}

	return vptr, err
}

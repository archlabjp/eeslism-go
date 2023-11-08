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

func (eqsys *EQSYS) Valvcountreset() {
	for _, v := range eqsys.Valv {
		v.Count = 0
	}
}

/***********************************************/

func (eqsys *EQSYS) Valvcountinc() {
	for _, v := range eqsys.Valv {
		v.Count++
	}
}

// 通常はバルブの上流の流量に比率を乗じるが、基準となるOMvavが指定されている場合には、この流量に対する比率とする
// OMvavが指定されているときだけの対応
func Valvinit(Valv []*VALV, Mpath []*MPATH) {
	for _, v := range Valv {
		if v.Cmp.MonPlistName != "" {
			for _, mpath := range Mpath {
				for _, Plist := range mpath.Plist {
					if v.Cmp.MonPlistName == Plist.Plistname {
						v.MonPlist = Plist
						v.MGo = &Plist.G

						if v.Cmb != nil {
							CValv := v.Cmb.Eqp.(*VALV)
							CValv.MonPlist = Plist
							CValv.MGo = &Plist.G
						}
						break
					}
				}
			}
		} else {
			Pelm := v.Plist.Pelm[len(v.Plist.Pelm)-1]
			v.MGo = Pelm.Cmp.Elouts[0].Lpath.Go
			v.MonPlist = Pelm.Cmp.Elouts[0].Lpath

			if v.Cmb != nil {
				CValv := v.Cmb.Eqp.(*VALV)
				CValv.MonPlist = v.MonPlist
				CValv.MGo = &v.MonPlist.G
			}
		}
	}
}

func Valvene(Valv []*VALV, Valvreset *int) {
	var etype EqpType
	var T1, T2 float64
	var Vcb *VALV
	var r float64

	for _, v := range Valv {
		etype = v.Cmp.Eqptype
		if etype == TVALV_TYPE && v.Org == 'y' {
			if v.Mon.Elouts[0].Control != OFF_SW {
				T1 = *v.Tin
				Vcb = v.Cmb.Eqp.(*VALV)
				T2 = *Vcb.Tin

				if math.Abs(*v.Tout-*v.Tset) >= 1.0e-3 && math.Abs(T1-T2) >= 1.0e-3 {
					r = (*v.Tset - T2) / (T1 - T2)
					r = math.Min(1.0, math.Max(r, 0.0))
					v.X = r
					Vcb.X = 1.0 - r

					v.Plist.Gcalc = r * *v.MGo
					Vcb.Plist.Gcalc = (1.0 - r) * *v.MGo
					(*Valvreset)++

					if DEBUG {
						fmt.Printf("<Valvene> Valvname=%s G=%f\n", v.Name, v.Plist.G)
						fmt.Printf("    T1=%.1f T2=%.1f Tset=%.1f\n", T1, T2, *v.Tset)
					}
				} else {
					v.Plist.Gcalc = v.X * *v.MGo
					Vcb.Plist.Gcalc = (1.0 - v.X) * *v.MGo
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

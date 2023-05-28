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

package main

import (
	"fmt"
	"io"
	"math"
	"os"
	"strings"
)

/* ------------------------------------------ */

func Valvcountreset(Nvalv int, Valv []VALV) {
	var i int

	for i = 0; i < Nvalv; i++ {
		Valv[i].Count = 0
	}
}

/***********************************************/

func Valvcountinc(Nvalv int, Valv []VALV) {
	var i int

	for i = 0; i < Nvalv; i++ {
		Valv[i].Count++
	}
}

// 通常はバルブの上流の流量に比率を乗じるが、基準となるOMvavが指定されている場合には、この流量に対する比率とする
// OMvavが指定されているときだけの対応
func Valvinit(NValv int, Valv []VALV, NMpath int, Mpath []MPATH) {
	for k := 0; k < NValv; k++ {
		if Valv[k].Cmp.MonPlistName != "" {
			for i := 0; i < NMpath; i++ {
				for j := 0; j < Mpath[i].Nlpath; j++ {
					Plist := &Mpath[i].Plist[j]
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
			Pelm := Valv[k].Plist.Pelm[Valv[k].Plist.Nelm-1]
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

func Valvene(Nvalv int, Valv []VALV, Valvreset *int) {
	var i int
	var etype EqpType
	var T1, T2 float64
	var Vcb *VALV
	var r float64

	for i = 0; i < Nvalv; i++ {
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

func ValvControl(fi io.Reader, Ncompnt int, Compnt []COMPNT, Schdl *SCHDL, Simc *SIMCONTL, Wd *WDAT, vptr *VPTR) {
	var s string
	var Valv, Vb *VALV
	var Vc *COMPNT
	var k, i int
	var ad int64
	var elins *ELIN
	var Pelm *PELM

	Vb = nil
	_, err := fmt.Fscanf(fi, "%s", &s)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	ad, err = fi.(io.Seeker).Seek(0, io.SeekCurrent)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	Vc = Compntptr(s, Ncompnt, Compnt)
	if Vc == nil {
		Eprint("<CONTRL>", s)
	}

	vptr.Ptr = &Vc.Control
	vptr.Type = SW_CTYPE

	Valv = Vc.Eqp.(*VALV)
	Valv.Org = 'y'
	for {
		_, err := fmt.Fscanf(fi, "%s", &s)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
			os.Exit(1)
		}

		if strings.HasPrefix(s, ";") {
			_, err := fi.(io.Seeker).Seek(ad, io.SeekStart)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			break
		}

		if s == "-init" {
			_, err := fmt.Fscanf(fi, "%s", &s)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			ad, err = fi.(io.Seeker).Seek(0, io.SeekCurrent)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			if k = idsch(s, Schdl.Sch, ""); k >= 0 {
				Valv.Xinit = &Schdl.Val[k]
			} else {
				Valv.Xinit = envptr(s, Simc, Ncompnt, Compnt, Wd, nil)
			}
		} else if s == "-Tout" {
			_, err := fmt.Fscanf(fi, "%s", &s)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			ad, err = fi.(io.Seeker).Seek(0, io.SeekCurrent)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			if k = idsch(s, Schdl.Sch, ""); k >= 0 {
				Valv.Tset = &Schdl.Val[k]
			} else {
				Valv.Tset = envptr(s, Simc, Ncompnt, Compnt, Wd, nil)
			}

			Pelm = Valv.Plist.Pelm[Valv.Plist.Nelm-1]
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

	Pelm = Valv.Plist.Pelm[Valv.Plist.Nelm-1]
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
func valv_vptr(key []string, Valv *VALV, vptr *VPTR) int {
	var err int

	if strings.Compare(key[1], "value") == 0 {
		vptr.Ptr = &Valv.X
		vptr.Type = VAL_CTYPE
	} else {
		err = 1
	}

	return err
}

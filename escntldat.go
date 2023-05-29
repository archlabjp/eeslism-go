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

/* esccntldat.c */

package main

import (
	"fmt"
	"strconv"
	"strings"
)

/*  制御、スケジュール設定式の入力  */

// char Hload,  /* Hload = HEATING_LOAD */
// 	 Cload,  /* Cload = COOLING_LOAD */
// 	 HCload;  /* HCload = HEATCOOL_LOAD */

func Contrldata(fi *EeTokens, Ct *[]CONTL, Ncontl *int, Ci *[]CTLIF, Nctlif *int,
	Cs *[]CTLST, Nctlst *int,
	Simc *SIMCONTL, Ncompnt int, Compnt []COMPNT,
	Nmpath int, Mpath []MPATH, Wd *WDAT, Exsf *EXSFS, Schdl *SCHDL) {
	//loadcmp, cmp := (*COMPNT)(nil), (*COMPNT)(nil)
	// varcontl, Contl, ctl := (*CONTL)(nil), (*CONTL)(nil), (*CONTL)(nil)
	// ctlif, Ctlif, cti := (*CTLIF)(nil), (*CTLIF)(nil), (*CTLIF)(nil)
	// ctlst, Ctlst, cts := (*CTLST)(nil), (*CTLST)(nil), (*CTLST)(nil)
	vptr, vpath := VPTR{}, VPTR{}
	var load *rune
	i := 0
	Nm := 0
	var loadcmp *COMPNT = nil
	Hload := HEATING_LOAD
	Cload := COOLING_LOAD
	HCload := HEATCOOL_LOAD

	Ni, N := ContrlCount(fi)

	Nm = N
	if Nm > 0 {
		*Ct = make([]CONTL, Nm)
		for i = 0; i < Nm; i++ {
			contl := &(*Ct)[i]
			contl.Lgv = 0
			contl.Type = ' '
			contl.Cif = nil
			contl.AndCif = nil
			contl.AndAndCif = nil
			contl.OrCif = nil
			contl.Cst = nil
		}

		*Cs = make([]CTLST, Nm)
		for i = 0; i < Nm; i++ {
			ctlst := &(*Cs)[i]
			ctlst.Type = ' '
			ctlst.PathType = ' '
			ctlst.Path = nil
		}
	}

	Nm = Ni
	if Ni > 0 {
		*Ci = make([]CTLIF, Nm)
		for i = 0; i < Nm; i++ {
			ctlif := &(*Ci)[i]
			ctlif.Type = ' '
			ctlif.Op = ' '
			ctlif.Nlft = 0
		}
	}

	contlIdx := 0
	ctlifIdx := 0
	ctlstIdx := 0
	load = nil
	for fi.IsEnd() == false {
		s := fi.GetToken()
		if len(s) == 0 || s[0] == '*' {
			break
		}

		load = nil
		VPTRinit(&vptr)
		VPTRinit(&vpath)

		Contl := &(*Ct)[contlIdx]
		Contl.Type = ' '
		Contl.Cst = &(*Cs)[ctlstIdx]

		for fi.IsEnd() == false {
			if s == "if" {
				Ctlif := &(*Ci)[ctlifIdx]
				Contl.Type = 'c'
				Contl.Cif = Ctlif
				s = strings.Trim(fi.GetToken(), "()")
				ctifdecode(s, Ctlif, Simc, Ncompnt, Compnt, Nmpath, Mpath, Wd, Exsf, Schdl)
				ctlifIdx++
				*Nctlif = ctlifIdx
			} else if s == "AND" {
				Ctlif := &(*Ci)[ctlifIdx]
				Contl.Type = 'c'
				if Contl.AndCif == nil {
					Contl.AndCif = Ctlif
				} else {
					Contl.AndAndCif = Ctlif
				}
				s = strings.Trim(fi.GetToken(), "()")
				ctifdecode(s, Ctlif, Simc, Ncompnt, Compnt, Nmpath, Mpath, Wd, Exsf, Schdl)
				ctlifIdx++
				*Nctlif = ctlifIdx
			} else if s == "OR" {
				Ctlif := &(*Ci)[ctlifIdx]
				Contl.Type = 'c'
				Contl.OrCif = Ctlif
				s = strings.Trim(fi.GetToken(), "()")
				ctifdecode(s, Ctlif, Simc, Ncompnt, Compnt, Nmpath, Mpath, Wd, Exsf, Schdl)
				ctlifIdx++
				*Nctlif = ctlifIdx
			} else if strings.HasPrefix(s, "LOAD") {
				loadcmp = nil
				if strings.ContainsRune(s, ':') {
					if len(s) == 5 && s[5] == HEATING_LOAD {
						load = new(rune)
						*load = Hload
					} else if len(s) == 5 && s[5] == COOLING_LOAD {
						load = new(rune)
						*load = Cload
					} else {
						i = idscw(s[5:], Schdl.Scw, "")
						if i >= 0 {
							load = &Schdl.Isw[i]
						} else {
							Eprint("<Contrldata>", s)
						}
					}
				} else {
					load = &HCload
				}
			} else if s == "-e" {
				s = fi.GetToken()
				for i = 0; i < Ncompnt; i++ {
					cmp := &Compnt[i]
					if s == cmp.Name {
						loadcmp = cmp
					}
				}
			} else if st := strings.IndexRune(s, '='); st != -1 {
				s = s[:st]
				var err int
				if load != nil {
					err = loadptr(loadcmp, load, s, Ncompnt, Compnt, &vptr)
					load = nil
				} else {
					vpath.Type = 0
					err = ctlvptr(s, Simc, Ncompnt, Compnt, Nmpath, Mpath, Wd, Exsf, Schdl, &vptr, &vpath)
				}
				if err == 0 {
					Ctlst := &(*Cs)[ctlstIdx]
					Ctlst.Type = vptr.Type
					Ctlst.PathType = vpath.Type
					Ctlst.Path = vpath.Ptr
					if Ctlst.Type == VAL_CTYPE {
						Ctlst.Lft.V = vptr.Ptr.(*float64)
					} else {
						Ctlst.Lft.S = vptr.Ptr.(*string)
					}
					err = ctlrgtptr(s[st+1:], &Ctlst.Rgt, Simc, Ncompnt, Compnt, Nmpath, Mpath, Wd, Exsf, Schdl, Ctlst.Type)
				}

				Err := fmt.Sprintf("%s = %s", s[:st], s[st+1:])
				Errprint(err, "<Contrldata>", Err)
			} else if s == "TVALV" {
				ctlstIdx--
				contlIdx--
				ValvControl(fi, Ncompnt, Compnt, Schdl, Simc, Wd, &vptr)
			} else {
				Eprint("<Contrldata>", s)
			}

			s = fi.GetToken()
			if len(s) == 0 || s[0] == ';' {
				break
			}
		}

		ctlstIdx++
		*Nctlst = ctlstIdx

		contlIdx++
		*Ncontl = contlIdx
	}
}

func ContrlCount(fi *EeTokens) (Nif, N int) {
	ad := fi.GetPos()
	var N1, N2 int

	Nif, N = 0, 0

	for fi.IsEnd() == false {
		s := fi.GetToken()
		if s[0] == '*' {
			break
		}

		switch s {
		case "if", "AND", "OR":
			N1++
		case "=", "TVALV":
			N2++
		}
	}

	N = N2
	Nif = N1

	fi.RestorePos(ad)

	return Nif, N
}

/* ------------------------------------------------------ */

/*  制御条件式 (lft1 - lft2 ? rgt ) に関するポインター */

func ctifdecode(_s string, ctlif *CTLIF, Simc *SIMCONTL, Ncompnt int, Compnt []COMPNT,
	Nmpath int, Mpath []MPATH, Wd *WDAT, Exsf *EXSFS, Schdl *SCHDL) {
	var lft, op, rgt string // 左変数, 演算子, 右変数
	var err int
	var vptr, vpath VPTR

	s := strings.Split(_s, " ")
	lft, op, rgt = s[0], s[1], s[2]

	st := strings.IndexRune(lft, '-')
	if st != -1 {
		lft = lft[:st]
	}

	// 演算対象の変数 その1を設定
	ctlvptr(lft, Simc, Ncompnt, Compnt, Nmpath, Mpath, Wd, Exsf, Schdl, &vptr, &vpath)

	ctlif.Type = vptr.Type // 演算の種類を設定
	ctlif.Nlft = 1
	if vptr.Type == VAL_CTYPE {
		ctlif.Lft1.V = vptr.Ptr.(*float64)
	} else {
		ctlif.Lft1.S = vptr.Ptr.(*string)
	}

	// 演算対象の変数 その2を設定
	if st != -1 {
		ctlvptr(lft[st:], Simc, Ncompnt, Compnt, Nmpath, Mpath, Wd, Exsf, Schdl, &vptr, &vpath)

		if vptr.Type == VAL_CTYPE && ctlif.Type == vptr.Type {
			ctlif.Nlft = 2
			ctlif.Lft2.V = vptr.Ptr.(*float64)
		} else {
			Eprint("<ctifdecode>", lft[st+1:])
		}
	}

	// 比較演算子 Op の設定
	err = 0
	if op == ">" {
		ctlif.Op = 'g'
	} else if op == ">=" {
		ctlif.Op = 'G'
	} else if op == "<" {
		ctlif.Op = 'l'
	} else if op == "<=" {
		ctlif.Op = 'L'
	} else if op == "==" {
		ctlif.Op = 'E'
	} else if op == "!=" {
		ctlif.Op = 'N'
	} else {
		err = 1
	}

	Errprint(err, "<ctifdecode>", _s)

	ctlrgtptr(rgt, &ctlif.Rgt, Simc, Ncompnt, Compnt, Nmpath, Mpath, Wd, Exsf, Schdl, ctlif.Type)
}

/* ------------------------------------------------------ */

/*  条件式、設定式の右辺（定数、またはスケジュール設定値のポインター） */

func ctlrgtptr(s string, rgt *CTLTYP, Simc *SIMCONTL, Ncompnt int, Compnt []COMPNT, Nmpath int, Mpath []MPATH, Wd *WDAT, Exsf *EXSFS, Schdl *SCHDL, _type VPtrType) int {
	var vptr VPTR
	var err int

	if _type == VAL_CTYPE && isstrdigit(s) {
		v, err := strconv.ParseFloat(s, 64)
		if err == nil {
			rgt.V = new(float64)
			*rgt.V = v
		}
	} else {
		switch s {
		case "OFF":
			rgt.S = new(string)
			*rgt.S = string(OFF_SW)
		case "ON":
			rgt.S = new(string)
			*rgt.S = string(ON_SW)
		case "COOL":
			rgt.S = new(string)
			*rgt.S = string(COOLING_SW)
		case "HEAT":
			rgt.S = new(string)
			*rgt.S = string(HEATING_SW)
		default:
			if _type == SW_CTYPE && strings.HasPrefix(s, "'") && len(s) > 2 {
				rgt.S = new(string)
				*rgt.S = s[1:2]
			} else {
				err = ctlvptr(s, Simc, Ncompnt, Compnt, Nmpath, Mpath, Wd, Exsf, Schdl, &vptr, nil)
				if _type == vptr.Type {
					if _type == VAL_CTYPE {
						rgt.V = vptr.Ptr.(*float64)
					} else {
						rgt.S = vptr.Ptr.(*string)
					}
				} else {
					err = 1
				}
			}
		}
	}

	Errprint(err, "<ctlrgtptr>", s)
	return err
}

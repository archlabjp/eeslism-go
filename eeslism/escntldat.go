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

package eeslism

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

/*  制御、スケジュール設定式の入力  */

// char Hload,  /* Hload = HEATING_LOAD */
// 	 Cload,  /* Cload = COOLING_LOAD */
// 	 HCload;  /* HCload = HEATCOOL_LOAD */

func Contrldata(fi *EeTokens, Ct *[]*CONTL, Ci *[]*CTLIF,
	Cs *[]*CTLST,
	Simc *SIMCONTL, Compnt []*COMPNT,
	Mpath []*MPATH, Wd *WDAT, Exsf *EXSFS, Schdl *SCHDL) {
	//loadcmp, cmp := (*COMPNT)(nil), (*COMPNT)(nil)
	// varcontl, Contl, ctl := (*CONTL)(nil), (*CONTL)(nil), (*CONTL)(nil)
	// ctlif, Ctlif, cti := (*CTLIF)(nil), (*CTLIF)(nil), (*CTLIF)(nil)
	// ctlst, Ctlst, cts := (*CTLST)(nil), (*CTLST)(nil), (*CTLST)(nil)
	vptr, vpath := VPTR{}, VPTR{}
	var load *ControlSWType
	i := 0
	Nm := 0
	var loadcmp *COMPNT = nil
	Hload := HEATING_LOAD
	Cload := COOLING_LOAD
	HCload := HEATCOOL_LOAD

	*Ct = make([]*CONTL, 0)
	*Cs = make([]*CTLST, 0)
	*Ci = make([]*CTLIF, Nm)

	load = nil
	for fi.IsEnd() == false {
		s := fi.GetToken()
		if s == "\n" {
			continue
		}
		if len(s) == 0 || s[0] == '*' {
			break
		}

		load = nil
		VPTRinit(&vptr)
		VPTRinit(&vpath)

		Contl := NewCONTL()
		Contl.Type = ' '
		Contl.Cst = NewCTLST()

		flag_ignore := false
		for fi.IsEnd() == false {
			if s == "if" {
				Ctlif := NewCTLIF()
				Contl.Type = 'c'
				Contl.Cif = Ctlif
				s = strings.Trim(fi.GetToken(), "()")
				ctifdecode(s, Ctlif, Simc, Compnt, Mpath, Wd, Exsf, Schdl)
				*Ci = append(*Ci, Ctlif)
			} else if s == "AND" {
				Ctlif := NewCTLIF()
				Contl.Type = 'c'
				if Contl.AndCif == nil {
					Contl.AndCif = Ctlif
				} else {
					Contl.AndAndCif = Ctlif
				}
				s = strings.Trim(fi.GetToken(), "()")
				ctifdecode(s, Ctlif, Simc, Compnt, Mpath, Wd, Exsf, Schdl)
				*Ci = append(*Ci, Ctlif)
			} else if s == "OR" {
				Ctlif := NewCTLIF()
				Contl.Type = 'c'
				Contl.OrCif = Ctlif
				s = strings.Trim(fi.GetToken(), "()")
				ctifdecode(s, Ctlif, Simc, Compnt, Mpath, Wd, Exsf, Schdl)
				*Ci = append(*Ci, Ctlif)
			} else if strings.HasPrefix(s, "LOAD") {
				loadcmp = nil
				if strings.ContainsRune(s, ':') {
					if len(s) == 5 && ControlSWType(s[5]) == HEATING_LOAD {
						load = new(ControlSWType)
						*load = Hload
					} else if len(s) == 5 && ControlSWType(s[5]) == COOLING_LOAD {
						load = new(ControlSWType)
						*load = Cload
					} else {
						var iderr error
						i, iderr = idscw(s[5:], Schdl.Scw, "")
						if iderr == nil {
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
				for _, cmp := range Compnt {
					if s == cmp.Name {
						loadcmp = cmp
					}
				}
			} else if st := strings.IndexRune(s, '='); st != -1 {
				ss := strings.SplitN(s, "=", 2)
				key, value := ss[0], ss[1]
				var err error
				if load != nil {
					vptr, err = loadptr(loadcmp, load, key, Compnt)
					load = nil
				} else {
					vptr, vpath, err = ctlvptr(key, Simc, Compnt, Mpath, Wd, Exsf, Schdl)
				}
				if err == nil {
					Ctlst := Contl.Cst
					Ctlst.Type = vptr.Type
					Ctlst.PathType = vpath.Type
					Ctlst.Path = vpath.Ptr
					if Ctlst.Type == VAL_CTYPE {
						Ctlst.Lft.V = vptr.Ptr.(*float64)
					} else {
						Ctlst.Lft.S = vptr.Ptr.(*ControlSWType)
					}
					err = ctlrgtptr(value, &Ctlst.Rgt, Simc, Compnt, Mpath, Wd, Exsf, Schdl, Ctlst.Type)
				}

				if err != nil {
					Err := fmt.Sprintf("%s = %s", s[:st], s[st+1:])
					Eprint("<Contrldata>", Err)
				}
			} else if s == "TVALV" {
				flag_ignore = true
				ValvControl(fi, Compnt, Schdl, Simc, Wd, &vptr)
			} else {
				Eprint("<Contrldata>", s)
			}

			s = fi.GetToken()
			if len(s) == 0 || s[0] == ';' {
				break
			}
		}

		if !flag_ignore {
			*Ct = append(*Ct, Contl)
			*Cs = append(*Cs, Contl.Cst)
		}
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
		case "TVALV":
			N2++
		}

		if strings.Contains(s, "=") {
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

func ctifdecode(_s string, ctlif *CTLIF, Simc *SIMCONTL, Compnt []*COMPNT,
	Mpath []*MPATH, Wd *WDAT, Exsf *EXSFS, Schdl *SCHDL) {
	var lft, op, rgt string // 左変数, 演算子, 右変数
	var err int
	var vptr VPTR

	s := strings.Split(_s, " ")
	lft, op, rgt = s[0], s[1], s[2]

	st := strings.IndexRune(lft, '-')
	if st != -1 {
		lft = lft[:st]
	}

	// 演算対象の変数 その1を設定
	vptr, _, _ = ctlvptr(lft, Simc, Compnt, Mpath, Wd, Exsf, Schdl)

	ctlif.Type = vptr.Type // 演算の種類を設定
	ctlif.Nlft = 1
	if vptr.Type == VAL_CTYPE {
		ctlif.Lft1.V = vptr.Ptr.(*float64)
	} else {
		ctlif.Lft1.S = vptr.Ptr.(*ControlSWType)
	}

	// 演算対象の変数 その2を設定
	if st != -1 {
		vptr, _, _ = ctlvptr(lft[st:], Simc, Compnt, Mpath, Wd, Exsf, Schdl)

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

	ctlrgtptr(rgt, &ctlif.Rgt, Simc, Compnt, Mpath, Wd, Exsf, Schdl, ctlif.Type)
}

/* ------------------------------------------------------ */

/*  条件式、設定式の右辺（定数、またはスケジュール設定値のポインター） */

func ctlrgtptr(s string, rgt *CTLTYP, Simc *SIMCONTL, Compnt []*COMPNT, Mpath []*MPATH, Wd *WDAT, Exsf *EXSFS, Schdl *SCHDL, _type VPtrType) error {
	var vptr VPTR
	var err error

	if _type == VAL_CTYPE && isstrdigit(s) {
		var v float64
		v, err = strconv.ParseFloat(s, 64)
		if err == nil {
			rgt.V = CreateConstantValuePointer(v)
		}
	} else {
		switch s {
		case "OFF":
			rgt.S = new(ControlSWType)
			*rgt.S = OFF_SW
		case "ON":
			rgt.S = new(ControlSWType)
			*rgt.S = ON_SW
		case "COOL":
			rgt.S = new(ControlSWType)
			*rgt.S = COOLING_SW
		case "HEAT":
			rgt.S = new(ControlSWType)
			*rgt.S = HEATING_SW
		default:
			if _type == SW_CTYPE && strings.HasPrefix(s, "'") && len(s) > 2 {
				rgt.S = new(ControlSWType)
				*rgt.S = ControlSWType(s[1])
			} else {
				vptr, _, err = ctlvptr(s, Simc, Compnt, Mpath, Wd, Exsf, Schdl)
				if _type == vptr.Type {
					if _type == VAL_CTYPE {
						rgt.V = vptr.Ptr.(*float64)
					} else {
						rgt.S = vptr.Ptr.(*ControlSWType)
					}
				} else {
					err = errors.New("Generated pointer type is not expected")
				}
			}
		}
	}

	//Errprint(err, "<ctlrgtptr>", s)
	return err
}

func NewCONTL() *CONTL {
	contl := new(CONTL)
	contl.Lgv = 0
	contl.Type = ' '
	contl.Cif = nil
	contl.AndCif = nil
	contl.AndAndCif = nil
	contl.OrCif = nil
	contl.Cst = nil
	return contl
}

func NewCTLST() *CTLST {
	ctlst := new(CTLST)
	ctlst.Type = ' '
	ctlst.PathType = ' '
	ctlst.Path = nil
	return ctlst
}

func NewCTLIF() *CTLIF {
	ctlif := new(CTLIF)
	ctlif.Type = ' '
	ctlif.Op = ' '
	ctlif.Nlft = 0
	return ctlif
}

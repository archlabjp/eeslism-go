package eeslism

import (
	"fmt"
	"io"
	"os"
	"strings"
)

/* 境界条件・負荷仮想機器の要素機器データ入力ファイル設定 */

func Vcfdata(fi *EeTokens, simcon *SIMCONTL) {
	var (
		s      string
		errFmt = "(vcfileint)"
	)

	N := VCFcount(fi)
	if N > 0 {
		simcon.Vcfile = make([]VCFILE, N)
		for i := range simcon.Vcfile {
			simcon.Vcfile[i] = VCFILE{
				Name:  "",
				Fname: "",
				Ad:    -999,
				Ic:    0,
				Tlist: nil,
				Fi:    nil,
				Estl: ESTL{
					Flid:     "",
					Title:    "",
					Wdatfile: "",
					Timeid:   "",
					Wdloc:    "",
					Catnm:    nil,
					Ntimeid:  0,
					Ntime:    0,
					dtm:      0,
					Nunit:    0,
					Nrqlist:  0,
					Nvreq:    0,
					Npreq:    0,
					Npprd:    0,
					Ndata:    0,
					Rq:       nil,
					Prq:      nil,
					Vreq:     []rune{},
					Unit:     []string{},
				},
			}
		}
	}

	vcIdx := 0
	for fi.IsEnd() == false {
		s = fi.GetToken()
		if s == "*" {
			break
		}

		vcfile := &simcon.Vcfile[vcIdx]

		vcfile.Name = s
		for fi.IsEnd() == false {
			s = fi.GetToken()
			if s == ";" {
				break
			}
			switch s {
			case "-f":
				vcfile.Fname = fi.GetToken()
			default:
				e := fmt.Sprintf("Vcfile=%s %s %s", vcfile.Name, errFmt, s)
				Eprint("<Vcfdata>", e)
			}
		}
		vcIdx++
	}

	simcon.Nvcfile = vcIdx

	for i := 0; i < simcon.Nvcfile; i++ {
		vcfile := &simcon.Vcfile[i]

		if f, err := os.Open(vcfile.Fname); err != nil {
			Eprint("<Vcfdata>", vcfile.Fname)
			os.Exit(EXIT_VCFILE)
		} else {
			vcfile.Fi = f
		}

		esondat(vcfile.Fi, &vcfile.Estl)
		N := vcfile.Estl.Ndata
		if N > 0 {
			vcfile.Tlist = make([]TLIST, N)
			for j := range vcfile.Tlist {
				vcfile.Tlist[j] = TLIST{
					Cname: "",
					Name:  "",
					Id:    "",
					Unit:  "",
					Fval:  nil,
					Fstat: nil,
					Ival:  nil,
					Istat: nil,
					Cval:  nil,
					Cstat: nil,
					Fmt:   "",
					Pair:  nil,
				}
			}
		}

		esoint(vcfile.Fi, "esoint", 1, &vcfile.Estl, vcfile.Tlist)
		vcfile.Ad, _ = vcfile.Fi.Seek(0, io.SeekCurrent)

		if simcon.Wdtype == 'E' {
			simcon.Wfname = vcfile.Fname
			wdflinit(simcon, &vcfile.Estl, vcfile.Tlist)
		}
	}
}

/***** VCFILEの定義数を数える ******/

func VCFcount(fi *EeTokens) int {
	var N int
	ad := fi.GetPos()

	for fi.IsEnd() == false {
		s := fi.GetToken()
		if s[0] != '*' {
			words := strings.Fields(s)
			if len(words) > 0 && words[0] == "-f" {
				N++
			}
		}
	}

	fi.RestorePos(ad)

	return N
}

/* -------------------------------------------------- */

/* 境界条件・負荷仮想機器の要素機器データとしての入力処理 */

func flindat(Flin *FLIN) {
	var s string
	n := 0
	//Err := fmt.Sprintf(ERRFMT, "(flindat)")
	ss := Flin.Cmp.Tparm

	for _, err := fmt.Sscanf(ss, "%s", &s); err == nil && strings.IndexRune(s, '*') == -1; _, err = fmt.Sscanf(ss, "%s", &s) {
		ss = ss[len(s):]
		for len(ss) > 0 && ss[0] == ' ' {
			ss = ss[1:]
		}

		if st := strings.IndexRune(s, '='); st != -1 {
			name, value := string(s[:st]), string(s[st+1:])
			if name == "t" {
				Flin.Namet = value
				n++
			} else if name == "x" {
				Flin.Namex = value
				n++
			}
		} else {
			Eprint("<flindat>", string(s))
			os.Exit(EXIT_FLIN)
		}
	}

	if n == 1 {
		Flin.Awtype = 'W'
		Flin.Cmp.Idi = []ELIOType{ELIO_W}
		Flin.Cmp.Ido = []ELIOType{ELIO_W}
		Flin.Cmp.Airpathcpy = 'n'
	} else {
		Flin.Awtype = 'A'
		Flin.Cmp.Idi = []ELIOType{ELIO_t, ELIO_x}
		Flin.Cmp.Ido = []ELIOType{ELIO_t, ELIO_x}
		Flin.Cmp.Airpathcpy = 'y'
	}

	Flin.Cmp.Nin = n
	Flin.Cmp.Nout = n
}

/* -------------------------------------------------- */

/* 境界条件・負荷仮想機器の要素機器データのポインター設定 */

func Flinint(Nflin int, Flin []FLIN, Simc *SIMCONTL, Ncompnt int, Compnt []COMPNT, Wd *WDAT) {
	for i := 0; i < Nflin; i++ {
		// fmt.Printf("<<Flinint>>  i=%d  namet=%s\n", i, Flin[i].namet)

		Flin[i].Vart = envptr(Flin[i].Namet, Simc, Ncompnt, Compnt, Wd, nil)
		if Flin[i].Awtype == 'A' {
			Flin[i].Varx = envptr(Flin[i].Namex, Simc, Ncompnt, Compnt, Wd, nil)
		}
	}
}

/* -------------------------------------------------- */

/* 境界条件・負荷仮想機器のファイル入力データのポインター */

func vcfptr(key []string, Simc *SIMCONTL, vptr *VPTR) int {
	var Ndata, err int = 1, 1
	var Tlist *TLIST

	for j := 0; j < Simc.Nvcfile; j++ {
		Vcfile := &Simc.Vcfile[j]

		if key[0] == Vcfile.Name {
			Ndata = Vcfile.Estl.Ndata
			for k := 0; k < Ndata; k++ {
				Tlist = &Vcfile.Tlist[k]
				if key[1] == Tlist.Name && key[2] == Tlist.Id {
					vptr.Ptr = Tlist.Fval
					vptr.Type = VAL_CTYPE
					err = 0
					break
				}
			}
		}
	}

	return err
}

/* -------------------------------------------------- */

/* 境界条件・負荷仮想機器のデータファイル入力 */

var __Vcfinput_Mon, __Vcfinput_Day, __Vcfinput_Time int

func Vcfinput(Daytm *DAYTM, Nvcfile int, Vcfile []VCFILE, perio rune) {
	var Tmdt TMDT
	TMDTinit(&Tmdt)

	var idend int
	for i := 0; i < Nvcfile; i++ {
		iderr := 0
		vcfile := &Vcfile[i]

		if vcfile.Estl.Tid == 'M' && __Vcfinput_Mon != Daytm.Mon {
			idend = tmdata(vcfile, &Tmdt, Daytm, perio)
			if idend != 0 {
				if Daytm.Mon == Tmdt.Mon {
					esdatgt(vcfile.Fi, 0, vcfile.Estl.Ndata, vcfile.Tlist)
				} else {
					iderr = 1
				}
			}
		} else if vcfile.Estl.Tid == 'd' && __Vcfinput_Day != Daytm.Day {
			idend = tmdata(vcfile, &Tmdt, Daytm, perio)
			if idend != 0 {
				if Daytm.Day == Tmdt.Day {
					esdatgt(vcfile.Fi, 0, vcfile.Estl.Ndata, vcfile.Tlist)
				} else {
					iderr = 1
				}
			}
		} else if vcfile.Estl.Tid == 'h' && __Vcfinput_Time != Daytm.Ttmm {
			idend = tmdata(vcfile, &Tmdt, Daytm, perio)
			if idend != 0 {
				if Daytm.Mon == Tmdt.Mon && Daytm.Day == Tmdt.Day && Daytm.Ttmm == Tmdt.Time {
					esdatgt(vcfile.Fi, 0, vcfile.Estl.Ndata, vcfile.Tlist)
				} else {
					iderr = 1
				}
			}
		}
		if idend == 0 {
			E := fmt.Sprintf("Vcfinput file-end: %s\n", vcfile.Fname)
			Eprint("<Vcfinput>", E)
		}

		if iderr != 0 {
			E := fmt.Sprintf("Vcfinput xxx file=%s prog_MM/DD/TM=%d/%d/%d file_MM/DD/TM=%d/%d/%d\n",
				vcfile.Fname, Daytm.Mon, Daytm.Day, Daytm.Ttmm, Tmdt.Mon, Tmdt.Day, Tmdt.Time)
			Eprint("<Vcfinput>", E)
		}
	}

	__Vcfinput_Mon, __Vcfinput_Day, __Vcfinput_Time = Daytm.Mon, Daytm.Day, Daytm.Ttmm
}

/********************************************************************/

func Flinprt(N int, Fl []FLIN) {
	if DEBUG {
		for i, f := range Fl {
			fmt.Printf("<< Flinprt >> Flin i=%d  %s %s = %.2g\n", i, f.Name, f.Namet, *f.Vart)
		}
	}
	for i, f := range Fl {
		if Ferr != nil {
			fmt.Fprintf(Ferr, "<< Flinprt >> Flin i=%d  %s %s = %.2g\n", i, f.Name, f.Namet, *f.Vart)
		}
	}
	if Ferr != nil {
		fmt.Fprintf(Ferr, "\n\n")
	}
}

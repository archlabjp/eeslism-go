package main

import (
	"bufio"

	"fmt"
	"io"
	"strconv"
	"strings"
)

/* シミュレーション結果、標題、識別データの入力 */

func esondat(fi io.Reader, Estl *ESTL) {
	var s string
	var i, j, Nparm, Ndat int
	var catnm, C *CATNM

	Estl.Catnm = nil

	for {
		_, err := fmt.Fscanf(fi, "%s", &s)
		if err != nil {
			break
		}

		switch s {
		case "-t":
			fmt.Fscanf(fi, " %[^;];", &s)
			Estl.Title = s
		case "-w":
			fmt.Fscanf(fi, "%s", &s)
			Estl.Wdatfile = s
		case "-tid":
			fmt.Fscanf(fi, " %c", &Estl.Tid)
		case "-u":
			i := 0
			for {
				_, err := fmt.Fscanf(fi, "%s", &s)
				if err != nil || s == ";" {
					break
				}
				Estl.Unit[i] = s
				i++
				Estl.Nunit = i
			}
		case "-Ntime":
			fmt.Fscanf(fi, " %d", &Estl.Ntime)
		case "-dtm":
			fmt.Fscanf(fi, " %d", &Estl.dtm)
		case "-tmid":
			fmt.Fscanf(fi, "%s", &s)
			Estl.Timeid = s
			Estl.Ntimeid = len(Estl.Timeid)
		case "-cat":
			N := CATNMMAX
			if N > 0 {
				Estl.Catnm = make([]CATNM, N)
			}

			if Estl.Catnm != nil {
				for ss := 0; ss < N; ss++ {
					C = &Estl.Catnm[ss]
					C.Name = ""
					C.N = 0
					C.Ncdata = 0
				}
			}

			Estl.Ndata = 0
			catidx := 0
			for {
				catnm = &Estl.Catnm[catidx]

				_, err := fmt.Fscanf(fi, "%s", &s)
				if err != nil || s == "*" {
					break
				}
				catnm.Name = s
				catnm.Ncdata = 0
				fmt.Fscanf(fi, "%d", &catnm.N)
				for i = 0; i < catnm.N; i++ {
					fmt.Fscanf(fi, "%s %d %d", &s, &Nparm, &Ndat)
					for j = 0; j < Nparm-1; j++ {
						fmt.Fscanf(fi, "%s", &s)
					}
					Estl.Ndata += Ndat
					catnm.Ncdata += Ndat
				}
				catidx++
			}
		case "-wdloc":
			fmt.Fscanf(fi, "%[^;];", &s)
			s += " ;"
			Estl.Wdloc = s
		case "-Ndata":
			fmt.Fscanf(fi, " %d", &Estl.Ndata)
		default:
			if s[len(s)-1] == '#' {
				Estl.Flid = s
			} else {
				Eprint("<esondat>", s)
			}
		}
	}

	if Estl.Title != "" {
		fmt.Printf("esondat  title=%s\n", Estl.Title)
	}

	if Estl.Title != "" {
		fmt.Printf("esondat  w=%s\n", Estl.Wdatfile)
	}

	fmt.Printf("esondat  tid=%c\n", Estl.Tid)
	fmt.Printf("esondat  Ntime=%d\n", Estl.Ntime)
	fmt.Printf("esondat  tmdt=%s\n", Estl.Timeid)
	fmt.Printf("esondat  flid=%s\n", Estl.Flid)
}

/* ----------------------------------------------------------- */

/* 要素名、シミュレーション結果入力用記憶域確保 */

func esoint(fi io.Reader, err string, Ntime int, Estl *ESTL, _Tlist []TLIST) {
	var nm, id string
	var V *rune
	var st int
	var cat *CATNM
	var R *RQLIST
	var n int
	var catIdx = 0
	var rqIdx = 0

	cat = nil
	R = nil
	V = nil

	if Estl.Catnm != nil {
		cat = &Estl.Catnm[0]
	}
	// Rq = &Estl.Rq[0]

	for i := 0; i < Estl.Ndata; i++ {
		Tlist := &_Tlist[i]
		fmt.Fscanf(fi, " %[^_]_%s %c %c", &nm, &id, &Tlist.Vtype, &Tlist.Ptype)

		switch Tlist.Vtype {
		case 'H', 'Q', 'E', 'q', 'e', 'm':
			Tlist.Stype = 't'
		case 'T', 'X', 'R', 't', 'x', 'r':
			Tlist.Stype = 'a'
		case 'c':
			Tlist.Stype = 'c'
		default:
			switch id[len(id)-1] {
			case 'n', 'c':
				Tlist.Stype = 'n'
			case 'm', 'h', 'e', 'p':
				Tlist.Stype = 'm'
			default:
				s := fmt.Sprintf("xxxx %s xxx  %s %s %c %c %c\n", err, nm, id, id[len(id)-1], Tlist.Vtype, Tlist.Ptype)
				Eprint("<esoint>", s)
			}
		}

		if Estl.Catnm != nil {
			if n >= cat.Ncdata {
				catIdx++
				cat = &Estl.Catnm[catIdx]
				n = 0
			}
			Tlist.Cname = cat.Name
			n++
		} else {
			Tlist.Cname = "*"
		}

		Tlist.Name = nm
		Tlist.Id = id
		Tlist.Req = 'n'

		if Estl.Nrqlist == 0 && Estl.Nvreq == 0 {
			Tlist.Req = 'y'
		} else {
			R = &Estl.Rq[rqIdx]
			for j := 0; j < Estl.Nrqlist; j++ {
				if (Tlist.Name == R.Name || R.Name == "*") &&
					(Tlist.Id == R.Id || R.Id == "*") {
					Tlist.Req = 'y'
					break
				} else if st = strings.IndexRune(Tlist.Name, ':'); st != -1 {

					if Tlist.Name[:st] == R.Name[:st] && R.Id == "*" {
						Tlist.Req = 'y'
						break
					}
				}
			}

			for j := 0; j < Estl.Nvreq; j++ {
				V = &Estl.Vreq[j]
				if *V == Tlist.Vtype {
					Tlist.Req = 'y'
					break
				}
			}
		}

		switch Tlist.Ptype {
		case 'f':
			Tlist.Fval = make([]float64, Ntime)
		case 'd':
			Tlist.Ival = make([]int, Ntime)
		case 'c':
			Tlist.Cval = make([]rune, Ntime)
		}

		fofmt(Estl, Tlist)
	}
}

/* ----------------------------------------------------------- */

/* 出力時の書式指定 */

func fofmt(Estl *ESTL, Tlist *TLIST) {
	var fmt string

	switch Tlist.Vtype {
	case 't', 'T':
		fmt = "%8.1f"
	case 'r', 'R':
		fmt = "%8.0f"
	case 'x', 'X':
		fmt = "%8.4f"
	case 'q', 'e', 'Q', 'E':
		fmt = "%8.1f"
	case 'H':
		fmt = "%8d"
	case 'h':
		fmt = "%04d"
	case 'c':
		fmt = "%c"
	}

	Tlist.Fmt = fmt

	for j := 0; j < Estl.Nunit; j++ {
		if Estl.Unit[j][0] == byte(Tlist.Vtype) {
			Tlist.Unit = Estl.Unit[2]
		}
	}
}

/* ----------------------------------------------------------- */

/*  年、月、日、曜日、時刻の入力 */

func tmdata(Vcfile *VCFILE, Tmdt *TMDT, Daytm *DAYTM, perio rune) int {
	var err error
	fi := Vcfile.Fi
	Estl := &Vcfile.Estl
	r := 1

	scanner := bufio.NewScanner(fi)
	for scanner.Scan() {
		s := scanner.Text()

		if s == "-999" || s == "end" {
			if Vcfile.Ic != 2 || perio == 'y' {
				_, err := fi.Seek(Vcfile.Ad, io.SeekStart)
				if err != nil {
					panic(err)
				}
				Vcfile.Ic++
			} else {
				return 0
			}
		} else {
			D := 0
			for i := 0; i < Estl.Ntimeid; i++ {
				if i > 0 {
					if scanner.Scan() {
						s = scanner.Text()
					} else {
						return r
					}
				}

				switch Estl.Timeid[i] {
				case 'Y': //年
					Tmdt.CYear = s
					Tmdt.Year, err = strconv.Atoi(s)
					if err != nil {
						panic(err)
					}
					Tmdt.Dat[i] = &Tmdt.CYear

					if Tmdt.Year == Daytm.Year {
						D++
					}
				case 'M': //月
					Tmdt.CMon = s
					Tmdt.Mon, err = strconv.Atoi(s)
					if err != nil {
						panic(err)
					}
					Tmdt.Dat[i] = &Tmdt.CMon

					if Tmdt.Mon == Daytm.Mon {
						D++
					}
				case 'D': //日
					Tmdt.CDay = s
					Tmdt.Day, err = strconv.Atoi(s)
					if err != nil {
						panic(err)
					}
					Tmdt.Dat[i] = &Tmdt.CDay

					if Tmdt.Day == Daytm.Day {
						D++
					}
				case 'W': //曜日
					Tmdt.CWkday = s
					Tmdt.Dat[i] = &Tmdt.CWkday
				case 'T': //時刻
					Tmdt.CTime = s
					if st := strings.IndexByte(s, ':'); st != -1 {
						s = s[:st] + "." + s[st+1:]
					}
					fval, err := strconv.ParseFloat(s, 64)
					if err != nil {
						panic(err)
					}
					Tmdt.Time = int(fval*100 + 0.5)
					Tmdt.Dat[i] = &Tmdt.CTime

					if Tmdt.Time-int(100*Daytm.Time+0.5) == 0 {
						D++
					}
				}
			}

			if D == Estl.Ntimeid {
				return 1
			} else {
				for i := 0; i < Estl.Ndata; i++ {
					if scanner.Scan() {
						scanner.Text()
					} else {
						return r
					}
				}
			}
		}
	}

	return r
}

/* ----------------------------------------------------------- */

/* シミュレーション結果データ入力 */

func esdatgt(fi io.Reader, i int, Ndata int, Tlist []TLIST) {
	scanner := bufio.NewScanner(fi)
	for j := 0; j < Ndata; j++ {
		if scanner.Scan() {
			s := scanner.Text()
			if Tlist[j].Req == 'y' || Tlist[j].Vtype == 'h' || Tlist[j].Vtype == 'H' {
				switch Tlist[j].Ptype {
				case 'f':
					fval, err := strconv.ParseFloat(s, 64)
					if err != nil {
						fmt.Println(err)
					} else {
						Tlist[j].Fval[i] = fval
						// fmt.Printf("<<esdatgt>> j=%d (data=%s)  %s %s [%d]=%f\n",
						// j, s, Tlist[j].name, Tlist[j].id, i, Tlist[j].fval[i])
					}
				case 'd':
					ival, err := strconv.Atoi(s)
					if err != nil {
						fmt.Println(err)
					} else {
						Tlist[j].Ival[i] = ival
					}
				case 'c':
					Tlist[j].Cval[i] = rune(s[0])
				}
				if j > 0 {
					Tml := &Tlist[j-1]
					if Tml.Vtype == 'h' || Tml.Vtype == 'H' {
						Tlist[j].Pair = Tml
					}
				}
			}
		}
	}
}

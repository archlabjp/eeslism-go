﻿//This file is part of EESLISM.
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

/*   schdata.c  */

package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
)

/* 曜日の設定  */

func Dayweek(fi io.Reader, Ipath string, daywk []int, key int) {
	var s string
	var ce int
	var ds, de, dd, d, id, M, D int
	var fw *os.File
	var err error

	var DAYweek = [8]string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun", ""}

	if fw, err = os.Open(strings.Join([]string{Ipath, "week.ewk"}, "")); err != nil {
		Eprint("<Dayweek>", "week.ewk")
		os.Exit(EXIT_WEEK)
	}

	if key == 0 {
		n, err := fmt.Fscanf(fi, "%d/%d=%s", &M, &D, &s)
		if n != 3 || err != nil {
			panic(err)
		}
	} else {
		fmt.Fscanf(fi, "%*s")

		fmt.Fscanf(fw, "%d/%d=%s", &M, &D, &s)
	}

	id = 0
	for d = 0; d < 8; d++ {
		if s == DAYweek[d] {
			id = d
		}
	}
	if id == 8 {
		Eprint("<Dayweek>", s)
	}

	ds = FNNday(M, D)
	de = ds + 365

	for dd = ds; dd < de; dd++ {
		d = dd
		if dd > 365 {
			d = dd - 365
		}
		daywk[d] = id
		id++

		if id > 6 {
			id = 0
		}
	}

	for {
		_, err := fmt.Fscanf(fi, "%s", &s)
		if err != nil || s[0] == ';' {
			break
		}
		var s1 string
		if ce = strings.IndexRune(s, ';'); ce != -1 {
			s1, _ = s[:ce], s[ce+1:]
		} else {
			s1, _ = s, ""
		}

		if ce = strings.IndexRune(s1, '/'); ce != -1 {
			//var s1_1, s1_2 string
			//s1_1, s1_2 := s1[:ce], s1[ce+1:]
			_, err = fmt.Sscanf(s1, "%d/%d", &M, &D)
			if err != nil {
				panic(err)
			}
			d = FNNday(M, D)
			daywk[d] = 7
		}
	}

	fw.Close()
}

/* ------------------------------------------------------------ */

/*  スケジュ－ル表の入力          */

func Schtable(fi io.ReadSeeker, dsn string, Schdl *SCHDL) {
	var s string
	var ce int
	var code byte
	var is, js, iw, j, sc, jsc, sw, jsw, Nmod int
	var ssn, wkd, vl, swn, N, i int
	var ic int
	var ssnmx, vlmx, swmx int
	var Sn *SEASN
	var Wk *WKDY
	var Dh *DSCH
	var Dw *DSCW

	//E := fmt.Sprintf(ERRFMT, dsn)

	Sn = nil
	Wk = nil
	Dh = nil
	Dw = nil

	ssn = 0
	wkd = 0
	vl = 0
	swn = 0

	ic = 0

	if ic == 0 {
		Schdl.Dsch = nil
		Schdl.Dscw = nil
		Schdl.Seasn = nil
		Schdl.Wkdy = nil

		SchCount(fi, &ssn, &wkd, &vl, &swn, &ssnmx, &vlmx, &swmx)
		ssn++
		wkd++
		vl++
		swn++
		ssnmx++
		vlmx++
		swmx++

		N = int(math.Max(float64(1), float64(ssn)))
		if N > 0 {
			Schdl.Seasn = make([]SEASN, N)
		}

		if Schdl.Seasn != nil {
			for i = 0; i < N; i++ {
				Seasn := SEASN{
					name: "",
					N:    0,
					end:  0,
					sday: make([]int, ssnmx),
					eday: make([]int, ssnmx),
				}
				Schdl.Seasn[i] = Seasn
			}
		}

		N = int(math.Max(float64(1), float64(wkd)))
		if N > 0 {
			Schdl.Wkdy = make([]WKDY, N)
		}

		if Schdl.Wkdy != nil {
			for i = 0; i < N; i++ {
				Wkdy := WKDY{
					name: "",
					end:  0,
					//wday: make([]int, 8),
				}
				for j = 0; j < 8; j++ {
					Wkdy.wday[j] = 0
				}
				Schdl.Wkdy[i] = Wkdy
			}
		}

		N = int(math.Max(float64(1), float64(vl)))
		if N > 0 {
			Schdl.Dsch = make([]DSCH, N)
		}

		if Schdl.Dsch != nil {
			for i = 0; i < N; i++ {
				Dsch := DSCH{
					name:  "",
					N:     0,
					end:   0,
					stime: make([]int, vlmx),
					etime: make([]int, vlmx),
					val:   make([]float64, vlmx),
				}
				Schdl.Dsch[i] = Dsch
			}
		}

		N = int(math.Max(float64(1), float64(swn)))
		if N > 0 {
			Schdl.Dscw = make([]DSCW, N)
		}

		if Schdl.Dscw != nil {
			for i = 0; i < N; i++ {
				Dscw := DSCW{
					name:  "",
					N:     0,
					end:   0,
					Nmod:  0,
					stime: make([]int, swmx),
					etime: make([]int, swmx),
					mode:  make([]rune, swmx),
					//dcode: make([]rune, swmx),
				}
				Schdl.Dscw[i] = Dscw
			}
		}
		ic = 1
	}

	// Seasn := &Schdl.Seasn[0]
	// Wkdy := &Schdl.Wkdy[0]
	// Dsch := &Schdl.Dsch[0]
	// Dscw := &Schdl.Dscw[0]

	for {
		_, err := fmt.Fscanf(fi, "%s", &s)
		if err != nil || s[0] == '*' {
			break
		}
		if s == "-ssn" || s == "SSN" {
			for {
				_, err := fmt.Fscanf(fi, "%s", &s)
				if err != nil || s[0] == ';' {
					break
				}
				if ce = strings.IndexRune(s, ';'); ce != -1 {
					before := s[:ce]
					after := s[ce+1:]
					s = before + after
				}

				if strings.IndexRune(s, '-') == -1 {
					is++
					Sn = &Schdl.Seasn[is]
					Sn.name = s
					js = -1
				} else {
					var Ms, Ds, Me, De int
					js++
					fmt.Sscanf(s, "%d/%d-%d/%d", &Ms, &Ds, &Me, &De)
					Sn.sday[js] = FNNday(Ms, Ds)
					Sn.eday[js] = FNNday(Me, De)
				}
				if ce != -1 {
					break
				}
			}
			Sn.N = js + 1
		} else if s == "-wkd" || s == "WKD" {
			j = 9
			for {
				_, err := fmt.Fscanf(fi, "%s", &s)
				if err != nil || s[0] == ';' {
					break
				}
				if ce = strings.IndexRune(s, ';'); ce != -1 {
					before := s[:ce]
					after := s[ce+1:]
					s = before + after
				}

				if j == 9 {
					iw++
					Wk = &Schdl.Wkdy[iw]
					Wk.name = s
					j = 0
				} else {
					wday := Wk.wday
					for j = 0; j < 8; j++ {
						if s == DAYweek[j] {
							wday[j] = 1
							break
						}
					}
				}
				if ce != -1 {
					break
				}
			}
		} else if s == "-v" || s == "VL" {
			for {
				_, err := fmt.Fscanf(fi, "%s", &s)
				if err != nil || s[0] == ';' {
					break
				}
				if ce = strings.IndexRune(s, ';'); ce != -1 {
					before := s[:ce]
					after := s[ce+1:]
					s = before + after
				}

				if strings.IndexRune(s, '(') == -1 {
					sc++
					Dh = &Schdl.Dsch[sc]
					Dh.name = s
					jsc = -1
				} else {
					jsc++

					// if jsc > SCDAYTMMAX {
					// 	fmt.Printf("<Schtable> Name=%s  MAX=%d  jsc=%d\n", Dh.name, SCDAYTMMAX, jsc)
					// }

					fmt.Sscanf(s, "%d-(%f)-%d", &Dh.stime[jsc], &Dh.val[jsc], &Dh.etime[jsc])
				}
				if ce != -1 {
					break
				}
			}
			Dh.N = jsc + 1
		} else if s == "-s" || s == "SW" {
			Nmod = 0
			fmt.Fscanf(fi, " %s ", &s)

			sw++
			Dw = &Schdl.Dscw[sw]
			Dw.name = s
			jsw = -1

			for {
				_, err := fmt.Fscanf(fi, "%s", &s)
				if err != nil || s[0] == ';' {
					break
				}
				if ce = strings.IndexRune(s, ';'); ce != -1 {
					before := s[:ce]
					after := s[ce+1:]
					s = before + after
				}

				jsw++
				fmt.Sscanf(s, "%d-(%c)-%d", &Dw.stime[jsw], &code, &Dw.etime[jsw])
				Dw.mode[jsw] = rune(code)

				for j = 0; j < Nmod; j++ {
					if Dw.dcode[j] == rune(code) {
						break
					}
				}

				if j == Nmod {
					Dw.dcode[Nmod] = rune(code)
					Nmod++
				}

				if ce != -1 {
					break
				}
			}
			Dw.N = jsw + 1
			Dw.Nmod = Nmod
		}
	}
	Schdl.Seasn[0].end = is + 1
	Schdl.Wkdy[0].end = iw + 1
	Schdl.Dsch[0].end = sc + 1
	Schdl.Dscw[0].end = sw + 1
}

/* ------------------------------------------------------------ */

/*  季節、曜日によるスケジュ－ル表の組み合わせ    */

func Schdata(fi io.Reader, dsn string, daywk []int, Schdl *SCHDL) {
	var (
		s       string
		ss      string
		dmod    rune
		ce      *rune
		dname   string
		i, j, k int
		N, d    int
		ds, de  int
		day     int
		is, iw  int
		sc, sw  int
	)

	const dmax = 366

	Seasn := Schdl.Seasn
	Wkdy := Schdl.Wkdy
	Dsch := Schdl.Dsch
	Dscw := Schdl.Dscw
	Sch := Schdl.Sch
	Scw := Schdl.Scw

	//E = fmt.Sprintf(ERRFMT, *dsn)
	i = Sch[0].end - 1
	j = Scw[0].end - 1

	scanner := bufio.NewScanner(fi)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "*" {
			break
		}
		fields := strings.Fields(line)

		if fields[0] == "-v" || fields[0] == "VL" {
			dmod = 'v'
		} else {
			dmod = 'w'
		}

		s = fields[1]

		if dmod == 'v' {
			i++
			S := Sch[i]
			S.name = string(s)
			for d := range S.day {
				S.day[d] = -1
			}
		} else {
			j++
			S := Scw[j]
			S.name = string(s)
			for d := range S.day {
				S.day[d] = -1
			}
		}

		for _, field := range fields[2:] {
			if ce = new(rune); strings.ContainsRune(field, ';') {
				*ce = ';'
				*ce = '\000'
			}
			var sname string
			var wname string
			is = -1
			iw = -1
			sc = -1
			sw = -1

			fmt.Sscanf(string(field), "%[^:]:%s", dname, &ss)
			if !strings.ContainsRune(ss, '-') {
				fmt.Sscanf(ss, "%s", &sname)
			} else {
				if ss[0] == '-' {
					fmt.Sscanf(ss[1:], "%s", &wname)
				} else {
					fmt.Sscanf(ss, "%[^-]-%s", &sname, &wname)
				}
			}

			if sname[0] != '\000' {
				is = idssn(string(sname), Seasn, "")
			}
			if wname[0] != '\000' {
				iw = idwkd(string(wname), Wkdy, "")
			}
			if dname[0] != '\000' {
				if dmod == 'v' {
					sc = iddsc(string(dname), Dsch, "")
				} else {
					sw = iddsw(string(dname), Dscw, "")
				}
			}
			if is >= 0 {
				N = Seasn[is].N
			} else {
				N = 1
			}

			for k = 0; k < N; k++ {
				if is >= 0 {
					Sn := Seasn[is]
					ds = Sn.sday[k]
					de = Sn.eday[k]

					if ds > de {
						de += 365
					}
				} else {
					ds = 1
					de = dmax
				}

				for day = ds; day <= de; day++ {
					d = day
					if day > 365 {
						d = day - 365
					}

					if iw < 0 || Wkdy[iw].wday[daywk[d]] == 1 {
						if dmod == 'v' {
							S := Sch[i]
							S.day[d] = sc
						} else {
							S := Scw[j]
							S.day[d] = sw
						}
					}
				}
			}

			if ce != nil {
				break
			}
		}
	}

	Schdl.Sch[0].end = i + 1
	Schdl.Scw[0].end = j + 1
}

/* ------------------------------------------------------------ */

/*  季節、曜日によるスケジュ－ル表の組み合わせ名へのスケジュ－ル名の追加  */

var __Schname_ind, __Schname_sco, __Schname_swo int

func Schname(Ipath string, dsn string, schdl *SCHDL) {
	var (
		sw, i, j, N       int
		E                 string
		vl, sws           int
		Dsch              = schdl.Dsch
		Dscw              = schdl.Dscw
		ssnmx, vlmx, swmx int
	)

	if fi, err := os.Open(Ipath + "schnma.ewk"); err != nil {
		Eprint("<Schname>", "schnma.ewk")
		os.Exit(EXIT_SCHTB)
	} else {
		defer fi.Close()

		SchCount(fi, &i, &j, &vl, &sws, &ssnmx, &vlmx, &swmx)
		vl++
		sws++
		ssnmx++
		vlmx++
		swmx++
	}

	if __Schname_ind == 0 {
		schdl.Sch = nil
		schdl.Scw = nil

		N = int(math.Max(float64(Dsch[0].end+vl), 1))
		schdl.Sch = make([]SCH, N)
		for i := 0; i < N; i++ {
			schdl.Sch[i] = SCH{
				name: "",
				end:  0,
				//day:  make([]int, 366),
			}
		}

		N = int(math.Max(float64(Dscw[0].end+sws), 1))
		schdl.Scw = make([]SCH, N)
		for i := 0; i < N; i++ {
			schdl.Scw[i] = SCH{
				name: "",
				end:  0,
				//day:  make([]int, 366),
			}
		}

		__Schname_ind = 1
	}

	i = schdl.Sch[0].end
	N = Dsch[0].end

	E = fmt.Sprintf(E, ERRFMT, dsn)

	for sc := __Schname_sco; sc < N; sc++ {
		Sch := &schdl.Sch[i]
		i++
		Sch.name = Dsch[sc].name

		for d := range Sch.day {
			Sch.day[d] = sc
		}
		__Schname_sco = sc
		schdl.Sch[0].end = i

		j = schdl.Scw[0].end
		N = Dscw[0].end
		for sw := __Schname_swo; sw < N; sw++ {
			Scw := &schdl.Scw[j]
			j++
			Scw.name = Dscw[sw].name

			for d := range Scw.day {
				Scw.day[d] = sw
			}
		}
		__Schname_swo = sw
		schdl.Scw[0].end = j
	}
}

/****  スケジュールの数を数える  ****/
func SchCount(fi io.ReadSeeker, ssn, wkd, vl, sw, ssnmx, vlmx, swmx *int) {
	var (
		s   string
		a   int64
		i   int
		err error
	)

	*ssnmx, *vlmx, *swmx = 0, 0, 0

	a, err = fi.Seek(0, io.SeekCurrent)
	if err != nil {
		_, _ = fi.Seek(0, io.SeekStart)
	}

	*ssn, *wkd, *vl, *sw = 0, 0, 0, 0

	for {
		_, err := fmt.Fscanf(fi, "%s", &s)
		if err != nil || string(s[:]) == "*" {
			break
		}

		imax := func(a, b int) int {
			if a > b {
				return a
			} else {
				return b
			}
		}

		if string(s[:]) == "-ssn" || string(s[:]) == "SSN" {
			(*ssn)++
			i = Schcmpcount(fi)
			*ssnmx = imax(i, *ssnmx)
		} else if string(s[:]) == "-wkd" || string(s[:]) == "WKD" {
			(*wkd)++
		} else if string(s[:]) == "-v" || string(s[:]) == "VL" {
			(*vl)++
			i = Schcmpcount(fi)
			*vlmx = imax(i, *vlmx)
		} else if string(s[:]) == "-s" || string(s[:]) == "SW" {
			(*sw)++
			i = Schcmpcount(fi)
			*swmx = imax(i, *swmx)
		}
	}

	_, _ = fi.Seek(a, io.SeekStart)
}

/***************************************************************************/

func Schcmpcount(fi io.Reader) int {
	N := 0
	scanner := bufio.NewScanner(fi)
	for scanner.Scan() {
		s := scanner.Text()
		if s == ";" {
			break
		}
		N++
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}
	return N - 1
}

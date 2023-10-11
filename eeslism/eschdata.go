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

/*   schdata.c  */

package eeslism

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

/* 曜日の設定  */

func Dayweek(fi io.Reader, week string, daywk []int, key int) {
	var s string
	var ce int
	var ds, de, dd, d, id, M, D int

	var DAYweek = [8]string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun", ""}

	if key == 0 {
		n, err := fmt.Fscanf(fi, "%d/%d=%s", &M, &D, &s)
		if n != 3 || err != nil {
			panic(err)
		}
	} else {
		re := regexp.MustCompile(`(\d+)/(\d+)=(\S+)`)
		matches := re.FindStringSubmatch(week)

		if matches != nil && len(matches) == 4 {
			M, _ = strconv.Atoi(matches[1])
			D, _ = strconv.Atoi(matches[2])
			s = matches[3]
		} else {
			panic("WEEKの書式が想定外です")
		}
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
}

/* ------------------------------------------------------------ */

/*  スケジュ－ル表の入力          */
var __Schtable_ic int
var __Schtable_is = -1
var __Schtable_js int
var __Schtable_iw = -1
var __Schtable_j int
var __Schtable_sc = -1
var __Schtable_jsc int
var __Schtable_sw = -1
var __Schtable_jsw int
var __Schtable_Nmod int

// SCHTBデータセットの読み取り
// SCHTBデータセット=一日の設定値、切換スケジュールおよび季節、曜日の指定
// 入力文字列`schtba`を読み取って、 [eeslism.SCHDL]に書き込む
func Schtable(schtba string, Schdl *SCHDL) {
	fi := strings.NewReader(schtba)

	var s string
	var ce int
	var code byte
	var ssn, wkd, vl, swn, N, i int
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

	__Schtable_ic = 0

	if __Schtable_ic == 0 {
		Schdl.Dsch = nil
		Schdl.Dscw = nil
		Schdl.Seasn = nil
		Schdl.Wkdy = nil

		SchCount(fi, &ssn, &wkd, &vl, &swn)
		ssn++
		wkd++
		vl++
		swn++

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
					sday: make([]int, ssn),
					eday: make([]int, ssn),
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
				for __Schtable_j = 0; __Schtable_j < 8; __Schtable_j++ {
					Wkdy.wday[__Schtable_j] = 0
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
					stime: make([]int, vl),
					etime: make([]int, vl),
					val:   make([]float64, vl),
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
					stime: make([]int, 0),
					etime: make([]int, 0),
					mode:  make([]rune, 0),
					//dcode: make([]rune, swmx),
				}
				Schdl.Dscw[i] = Dscw
			}
		}
		__Schtable_ic = 1
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

		// 季節設定
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
					__Schtable_is++
					Sn = &Schdl.Seasn[__Schtable_is]
					Sn.name = s
					__Schtable_js = -1
				} else {
					var Ms, Ds, Me, De int
					__Schtable_js++
					fmt.Sscanf(s, "%d/%d-%d/%d", &Ms, &Ds, &Me, &De)
					Sn.sday[__Schtable_js] = FNNday(Ms, Ds)
					Sn.eday[__Schtable_js] = FNNday(Me, De)
				}
				if ce != -1 {
					break
				}
			}
			Sn.N = __Schtable_js + 1

		} else if s == "-wkd" || s == "WKD" {
			// 曜日設定
			__Schtable_j = 9
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

				if __Schtable_j == 9 {
					__Schtable_iw++
					Wk = &Schdl.Wkdy[__Schtable_iw]
					Wk.name = s
					__Schtable_j = 0
				} else {
					wday := Wk.wday
					for __Schtable_j = 0; __Schtable_j < 8; __Schtable_j++ {
						if s == DAYweek[__Schtable_j] {
							wday[__Schtable_j] = 1
							break
						}
					}
				}
				if ce != -1 {
					break
				}
			}
		} else if s == "-v" || s == "VL" {
			// 設定値スケジュール定義
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
					__Schtable_sc++
					Dh = &Schdl.Dsch[__Schtable_sc]
					Dh.name = s
					__Schtable_jsc = -1
				} else {
					__Schtable_jsc++

					// if jsc > SCDAYTMMAX {
					// 	fmt.Printf("<Schtable> Name=%s  MAX=%d  jsc=%d\n", Dh.name, SCDAYTMMAX, jsc)
					// }

					fmt.Sscanf(s, "%d-(%f)-%d", &Dh.stime[__Schtable_jsc], &Dh.val[__Schtable_jsc], &Dh.etime[__Schtable_jsc])
				}
				if ce != -1 {
					break
				}
			}
			Dh.N = __Schtable_jsc + 1
		} else if s == "-s" || s == "SW" {
			// 切替設定スケジュール定義
			__Schtable_Nmod = 0
			fmt.Fscanf(fi, " %s ", &s)

			__Schtable_sw++
			Dw = &Schdl.Dscw[__Schtable_sw]
			Dw.name = s
			__Schtable_jsw = -1

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

				__Schtable_jsw++
				fmt.Sscanf(s, "%d-(%c)-%d", &Dw.stime[__Schtable_jsw], &code, &Dw.etime[__Schtable_jsw])
				Dw.mode[__Schtable_jsw] = rune(code)

				for __Schtable_j = 0; __Schtable_j < __Schtable_Nmod; __Schtable_j++ {
					if Dw.dcode[__Schtable_j] == rune(code) {
						break
					}
				}

				if __Schtable_j == __Schtable_Nmod {
					Dw.dcode[__Schtable_Nmod] = rune(code)
					__Schtable_Nmod++
				}

				if ce != -1 {
					break
				}
			}
			Dw.N = __Schtable_jsw + 1
			Dw.Nmod = __Schtable_Nmod
		}
	}
	Schdl.Seasn[0].end = __Schtable_is + 1
	Schdl.Wkdy[0].end = __Schtable_iw + 1
	Schdl.Dsch[0].end = __Schtable_sc + 1
	Schdl.Dscw[0].end = __Schtable_sw + 1
}

/* ------------------------------------------------------------ */

/*  季節、曜日によるスケジュ－ル表の組み合わせ    */

// SCHNMデータセットの読み取り
// SCHNMデータセット = 季節、曜日によるスケジュ－ル表の組み合わせ
// 入力文字列`schenma`を読み取って、[eeslism.SCHDL]に書き込む
func Schdata(schnma string, dsn string, daywk []int, Schdl *SCHDL) {
	fi := strings.NewReader(schnma)

	var (
		s       string
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
			// 設定値名
			dmod = 'v'
		} else {
			// 切換設定名
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

			// 正規表現パターン
			pattern := `^(\w+)(?::(\w*))?(?:-(\w+))?`
			re := regexp.MustCompile(pattern)

			// マッチする部分を取り出す
			match := re.FindStringSubmatch(field)

			// 3つの部分を取り出し
			if len(match) >= 2 {
				dname = match[1]
			} else if len(match) >= 3 {
				sname = match[2]
			} else if len(match) >= 4 {
				sname = match[3]
			} else {
				panic("一致する部分が見つかりませんでした")
			}

			if sname != "" {
				is = idssn(string(sname), Seasn, "")
			}
			if wname != "" {
				iw = idwkd(string(wname), Wkdy, "")
			}
			if dname != "" {
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
		sw, i, j, N int
		E           string
		vl, sws     int
		Dsch        = schdl.Dsch
		Dscw        = schdl.Dscw
	)

	if fi, err := os.Open(Ipath + "schnma.ewk"); err != nil {
		Eprint("<Schname>", "schnma.ewk")
		os.Exit(EXIT_SCHTB)
	} else {
		defer fi.Close()

		SchCount(fi, &i, &j, &vl, &sws)
		vl++
		sws++
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
func SchCount(fi io.ReadSeeker, ssn, wkd, vl, sw *int) {
	scanner := bufio.NewScanner(fi)
	for scanner.Scan() {
		s := scanner.Text()
		fields := strings.Fields(s)
		if len(fields) > 0 {
			if fields[0] == "*" {
				break
			} else if fields[0] == "-ssn" || fields[0] == "SSN" {
				(*ssn)++
			} else if fields[0] == "-wkd" || fields[0] == "WKD" {
				(*wkd)++
			} else if fields[0] == "-v" || fields[0] == "VL" {
				(*vl)++
			} else if fields[0] == "-s" || fields[0] == "SW" {
				(*sw)++
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}
}

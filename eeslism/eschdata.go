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

// SCHTBデータセットの読み取り
// SCHTBデータセット=一日の設定値、切換スケジュールおよび季節、曜日の指定
// 入力文字列`schtba`を読み取って、 [eeslism.SCHDL]に書き込む
// NOTE: SCHTBデータセットと %s の両方を読み取るために無理が出ている
func Schtable(schtba string, Schdl *SCHDL) {
	var code byte

	tokens := NewEeTokens(schtba)

	for tokens.IsEnd() == false {
		s := tokens.GetToken()
		if s == "*" {
			break
		}
		if s == "\n" || s == ";" {
			continue
		}

		switch s {
		case "-ssn", "SSN":
			// 季節設定

			fields := tokens.GetLogicalLine()
			n := len(fields) - 1

			Sn := SEASN{
				name: fields[0], // 季節名
				N:    n,
				sday: make([]int, n),
				eday: make([]int, n),
			}

			// 開始日・終了日
			for i := 0; i < n; i++ {
				var Ms, Ds, Me, De int
				fmt.Sscanf(fields[i+1], "%d/%d-%d/%d", &Ms, &Ds, &Me, &De)
				Sn.sday[i] = FNNday(Ms, Ds)
				Sn.eday[i] = FNNday(Me, De)
			}
			Sn.N = n

			Schdl.Seasn = append(Schdl.Seasn, Sn)

			break
		case "-wkd", "WKD":
			// 曜日設定

			fields := tokens.GetLogicalLine()
			n := len(fields) - 1

			Wk := WKDY{
				name: fields[0], // 曜日名
			}

			// 対応する曜日のフラグを埋める
			wday := Wk.wday
			for i := 0; i < n; i++ {
				for j := 0; j < 8; j++ {
					if fields[i+1] == DAYweek[j] {
						wday[j] = 1
						break
					}
				}
			}

			Schdl.Wkdy = append(Schdl.Wkdy, Wk)

			break
		case "-v", "VL":
			// 設定値スケジュール定義
			fields := tokens.GetLogicalLine()
			n := len(fields) - 1

			Dh := DSCH{
				name:  fields[0], // 設定値名
				N:     n,
				stime: make([]int, n),
				etime: make([]int, n),
				val:   make([]float64, n),
			}

			// 開始時分, 終了時分, 設定値
			Dh.stime = make([]int, n)
			Dh.val = make([]float64, n)
			Dh.etime = make([]int, n)
			for i := 0; i < n; i++ {
				fmt.Sscanf(fields[i+1], "%d-(%f)-%d", &Dh.stime[i], &Dh.val[i], &Dh.etime[i])
			}

			Schdl.Dsch = append(Schdl.Dsch, Dh)

			break
		case "-s", "SW":
			// 切替設定スケジュール定義
			fields := tokens.GetLogicalLine()
			n := len(fields) - 1
			nmod := 0

			Dw := DSCW{
				name:  fields[0], // 切り替え設定名
				N:     n,
				Nmod:  0,
				stime: make([]int, 10),
				etime: make([]int, 10),
				mode:  make([]rune, 10),
				//dcode: make([]rune, swmx),
			}

			Dw.stime = make([]int, n)
			Dw.etime = make([]int, n)
			Dw.mode = make([]rune, n)
			for i := 0; i < n; i++ {
				fmt.Sscanf(fields[i+1], "%d-(%c)-%d", &Dw.stime[i], &code, &Dw.etime[i])
				Dw.mode[i] = rune(code)
			}

			// モード数を調べる
			var j int
			for j = 0; j < nmod; j++ {
				if Dw.dcode[j] == rune(code) {
					break
				}
			}

			if j == nmod {
				Dw.dcode[nmod] = rune(code)
				nmod++
			}

			Dw.Nmod = nmod

			Schdl.Dscw = append(Schdl.Dscw, Dw)

			break
		}
	}
}

/* ------------------------------------------------------------ */

/*  季節、曜日によるスケジュ－ル表の組み合わせ    */

// SCHNMデータセットの読み取り
// SCHNMデータセット = 季節、曜日によるスケジュ－ル表の組み合わせ
// 入力文字列`schenma`を読み取って、[eeslism.SCHDL]に書き込む
func Schdata(schnma string, dsn string, daywk []int, Schdl *SCHDL) {
	fi := strings.NewReader(schnma)

	var (
		s      string
		dmod   rune
		ce     *rune
		dname  string
		k      int
		N, d   int
		ds, de int
		day    int
		is, iw int
		sc, sw int
		err    error
	)

	const dmax = 366

	Seasn := Schdl.Seasn
	Wkdy := Schdl.Wkdy
	Dsch := Schdl.Dsch
	Dscw := Schdl.Dscw

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

		// 年間スケジュールの初期化
		S := SCH{
			name: string(s),
		}
		for d := range S.day {
			S.day[d] = -1
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
				is, err = idssn(string(sname), Seasn)
				if err != nil {
					panic(err)
				}
			}
			if wname != "" {
				iw, err = idwkd(string(wname), Wkdy)
				if err != nil {
					panic(err)
				}
			}
			if dname != "" {
				if dmod == 'v' {
					sc, err = iddsc(string(dname), Dsch)
					if err != nil {
						panic(err)
					}
				} else {
					sw, err = iddsw(string(dname), Dscw)
					if err != nil {
						panic(err)
					}
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
							S.day[d] = sc
						} else {
							S.day[d] = sw
						}
					}
				}
			}

			// 年間スケジュールに追加
			if dmod == 'v' {
				Schdl.Sch = append(Schdl.Sch, S)
			} else {
				Schdl.Scw = append(Schdl.Scw, S)
			}

			if ce != nil {
				break
			}
		}
	}
}

/* ------------------------------------------------------------ */

/*  季節、曜日によるスケジュ－ル表の組み合わせ名へのスケジュ－ル名の追加  */

func Schname(schdl *SCHDL) {
	// 年間一定のスケジュールを追加
	for i, sc := range schdl.Dsch {
		sch := SCH{
			name: sc.name,
		}
		for d := range sch.day {
			sch.day[d] = i
		}

		schdl.Sch = append(schdl.Sch, sch)
	}

	// 年間一定のスケジュールを追加
	for j, sw := range schdl.Dscw {
		scw := SCH{
			name: sw.name,
		}
		for d := range scw.day {
			scw.day[d] = j
		}

		schdl.Scw = append(schdl.Scw, scw)
	}
}

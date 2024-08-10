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
	"regexp"
	"strconv"
	"strings"
)

/* 曜日の設定  */

func Dayweek(fi string, week string, daywk []int, key int) {
	var s string
	var d, id, M, D int

	var DAYweek = [8]string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun", ""}

	if key == 0 {
		n, err := fmt.Sscanf(fi, "%d/%d=%s", &M, &D, &s)
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

	// 開始日と終了日
	ds := FNNday(M, D)
	de := ds + 365

	// 1日ごとのループ
	for dd := ds; dd < de; dd++ {
		d = dd
		if dd > 365 {
			d = dd - 365
		}
		daywk[d] = id // 曜日の記録
		id++

		if id > 6 {
			id = 0
		}
	}

	// `dayweek.efl`から祝日を読み取る
	tokens := NewEeTokens(fi)

	for {
		s = tokens.GetToken()
		if s == "" || s == ";" {
			break
		}

		re := regexp.MustCompile(`^(\d+)/(\d+)$`)
		matches := re.FindStringSubmatch(s)

		if matches != nil && len(matches) == 3 {
			M, _ = strconv.Atoi(matches[1])
			D, _ = strconv.Atoi(matches[2])
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
	var code ControlSWType

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
			n := len(fields) - 2

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
			for i := 0; i < n; i++ {
				for j := 0; j < 8; j++ {
					if fields[i+1] == DAYweek[j] {
						Wk.wday[j] = true
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
				mode:  make([]ControlSWType, 10),
				//dcode: make([]rune, swmx),
			}

			Dw.stime = make([]int, n)
			Dw.etime = make([]int, n)
			Dw.mode = make([]ControlSWType, n)
			for i := 0; i < n; i++ {
				fmt.Sscanf(fields[i+1], "%d-(%c)-%d", &Dw.stime[i], &code, &Dw.etime[i])
				Dw.mode[i] = ControlSWType(code)
			}

			// モード数を調べる
			var j int
			for j = 0; j < nmod; j++ {
				if Dw.dcode[j] == code {
					break
				}
			}

			if j == nmod {
				Dw.dcode[nmod] = code
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

	var err error
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

		// 定義の種類
		var dmod rune
		if fields[0] == "-v" || fields[0] == "VL" {
			// 設定値名
			dmod = 'v'
		} else if fields[0] == "-s" || fields[0] == "SW" {
			// 切換設定名
			dmod = 'w'
		} else {
			panic(fields[0])
		}

		// 年間スケジュールの初期化
		S := SCH{
			name: fields[1], // 設定値名 or 切替設定名
		}
		for d := range S.day {
			S.day[d] = -1
		}

		// ';' まで繰り返す
		for _, field := range fields[2:] {

			if field == ";" {
				break
			}

			var dname string
			var sname string
			var wname string
			is := -1
			var wkday *WKDY = nil
			sc := -1
			sw := -1

			// 正規表現パターン
			pattern := `^(\w+)(?::(\w*))?(?:-(\w+))?`
			re := regexp.MustCompile(pattern)

			// マッチする部分を取り出す
			match := re.FindStringSubmatch(field)

			// 3つの部分を取り出し
			if len(match) >= 2 {
				// ex) `TrsetC`
				dname = match[1] // 参照する1日の設定名
			}
			if len(match) >= 3 {
				// ex) `TrsetH:Winter`
				sname = match[2] // 季節設定名 ex)Summer
			}
			if len(match) >= 4 {
				// ex) `ACSWLDwd:Summer-Weekday`
				// ex) `ACSWLDwd:-Weekday`
				wname = match[3] // 曜日設定名 ex)Weekday
			}

			if sname != "" {
				// 季節設定の検索
				is, err = idssn(sname, Seasn)
				if err != nil {
					panic(err)
				}
			}
			if wname != "" {
				// 曜日設定の検索
				iw, err := idwkd(wname, Wkdy)
				if err != nil {
					panic(err)
				}
				wkday = &Wkdy[iw]
			}
			if dname != "" {
				if dmod == 'v' {
					// 一日の設定値スケジュ－ルの検索
					sc, err = iddsc(dname, Dsch)
					if err != nil {
						panic(err)
					}
				} else if dmod == 'w' {
					// 一日の切り替えスケジュ－ルの検索
					sw, err = iddsw(string(dname), Dscw)
					if err != nil {
						panic(err)
					}
				} else {
					panic(dmod)
				}
			}

			// ループ回数
			var N int
			if is >= 0 {
				// ** 季節設定がある場合 **
				N = Seasn[is].N
			} else {
				// ** 季節設定がない場合 **
				N = 1
			}

			// 年間スケジュールの作成ループ
			for k := 0; k < N; k++ {
				var ds, de int
				if is >= 0 {
					// ** 季節設定がある場合 **
					// ex) `TrsetC:Winter`
					// ex) `ACSWLDwd:Winter-Weekday`
					Sn := Seasn[is]
					ds = Sn.sday[k] //開始日
					de = Sn.eday[k] //終了日

					if ds > de {
						de += 365 // 年末跨ぎ
					}
				} else {
					// ** 季節設定がない場合場合 **
					// ex) `TrsetC`
					// ex) `ACSWLDwd:-Weekday`
					ds = 1    // 開始日 NOTE: 配列インデックス1-366を想定 (Fortran譲りか)
					de = dmax // 終了日
				}

				for day := ds; day <= de; day++ {
					d := day
					if day > 365 {
						d = day - 365 // NOTE: d=1に戻る条件になっている
					}

					// 曜日指定が無い or 指定曜日であることを確認
					if wkday == nil || wkday.wday[daywk[d]] {
						if dmod == 'v' {
							S.day[d] = sc
						} else if dmod == 'w' {
							S.day[d] = sw
						} else {
							panic(dmod)
						}
					}
				}
			}
		}

		// 年間スケジュールに追加
		if dmod == 'v' {
			Schdl.Sch = append(Schdl.Sch, S)
		} else if dmod == 'w' {
			Schdl.Scw = append(Schdl.Scw, S)
		} else {
			panic(dmod)
		}
	}

	// Val, Isw の領域確保
	Schdl.Val = make([]float64, len(Schdl.Sch))
	Schdl.Isw = make([]ControlSWType, len(Schdl.Scw))
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

	// Val, Isw の領域確保
	schdl.Val = make([]float64, len(schdl.Sch))
	schdl.Isw = make([]ControlSWType, len(schdl.Scw))
}

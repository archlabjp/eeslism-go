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

/*   binit.c   */
package eeslism

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

/*  壁体デ－タの入力  */

func Walldata(section *EeTokens, fbmlist string, dsn string, Wall *[]WALL, Nwall *int, dfwl *DFWL, pcm []PCM, Npcm int) {
	var s string
	var i = -1
	var j, jj, jw, Nlyr, k = 0, 0, 0, 0, -1
	var Nbm, ndiv int // ic
	var dt float64
	var W []BMLST
	var Wl *BMLST

	W = nil
	// E = fmt.Sprintf(ERRFMT, dsn)

	NwLines := wbmlistcount(fbmlist)

	s = "Walldata wbmlist.efl--"

	k = 0

	for _, line := range NwLines {
		s := strings.Fields(line)

		Wl = new(BMLST)
		Wl.Mcode = s[0]

		for j := 0; j < k-1; j++ {
			Wc := &W[j]
			if Wl.Mcode == Wc.Mcode {
				message := fmt.Sprintf("wbmlist.efl duplicate code=<%s>", Wl.Mcode)
				Eprint("<Walldata>", message)
			}
		}

		// Cond：熱伝導率［W/mK］、Cro：容積比熱［kJ/m3K］
		var err error
		Wl.Cond, err = strconv.ParseFloat(s[1], 64)
		if err != nil {
			Eprint("<Walldata>", "wbmlist.efl Cond error")
		}

		Wl.Cro, err = strconv.ParseFloat(s[2], 64)
		if err != nil {
			Eprint("<Walldata>", "wbmlist.efl Cro error")
		}

		W = append(W, *Wl)

		for j := 0; j < k+1; j++ {
			fmt.Printf("k=%d code=%s\n", k, W[j].Mcode)
		}

		fmt.Printf("Walldata>> name=%s Con=%f Cro=%f\n", W[k].Mcode, W[k].Cond, W[k].Cro)
		W[k] = *Wl
		k++
	}

	k++
	Nbm = k

	//N = Wallcount(fi)

	s = "Walldata --"

	i = 0
	for section.IsEnd() == false {
		line := section.GetLogicalLine()
		if line[0] == "*" {
			break
		}

		Wa := NewWall()

		s = line[0]

		// for jj = 0; jj < len(Wa.welm); jj++ {
		// 	w = NewWelm()
		// }

		// for jj = 0; jj < len(Wa.PCM); jj++ {
		// 	Wa.PCM[jj] = nil
		// 	Wa.PCMrate[jj] = -999.0
		// }

		// (1) 部位・壁体名の読み取り
		if strings.HasPrefix(s, "-") {
			// For `-E` or `-R:ROOF` or ...

			// 部位コードの指定
			Wa.ble = rune(s[1])

			if s[2] == ':' {
				// 部位と壁体名の指定(最初の1文字を必ず英字とする)
				Wa.name = s[3:]
			} else {
				// 部位のみ指定（既定値の定義）
				switch s[1] {
				case 'E': // 外壁
					dfwl.E = i
				case 'R': // 屋根
					dfwl.R = i
				case 'F': // 床(外部)
					dfwl.F = i
				case 'i': // 内壁
					dfwl.i = i
				case 'c': // 天井(内部)
					dfwl.c = i
				case 'f': // 床(内部)
					dfwl.f = i
				}
			}
		} else {
			// 部位コードの指定なし
			Wa.name = s
		}

		j = -1

		// (2) 部位・壁体のパラメータを読み取り
		// 例) `Eo=0.9 Ei=0.9 as=0.7 type=1 APR-100 APR-100/20 <P:80.3> ;`
		var layer []string
		for _, s = range line[1:] {
			var err error
			st := strings.IndexRune(s, '=')
			if st != -1 {
				dt, err = strconv.ParseFloat(s[st+1:], 64)
				if err != nil {
					panic(err)
				}

				switch strings.TrimSpace(s[:st]) {
				case "Ei":
					Wa.Ei = dt // 室内表面放射率
				case "Eo":
					Wa.Eo = dt // 外表面放射率
				case "as":
					Wa.as = dt // 外表面日射吸収率
				case "type":
					Wa.ColType = s[st+1:]
				case "tra":
					Wa.tra = dt // τα
				case "Ksu":
					Wa.Ksu = dt // 通気層内上側から屋外までの熱貫流率 [W/m2K]
				case "Ksd":
					Wa.Ksd = dt // 通気層内下側から裏面までの熱貫流率 [W/m2K]
				case "Ru":
					Wa.Ru = dt // 通気層から上面までの熱抵抗 [m2K/W]
				case "Rd":
					Wa.Rd = dt // 通気層から裏面までの熱抵抗 [m2K/W]
				case "fcu":
					Wa.fcu = dt // Kcu / Ksu （太陽光発電設置時のアレイ温度計算のみに使用）
				case "fcd":
					Wa.fcd = dt // Kcd / Ksd （太陽光発電設置時のアレイ温度計算のみに使用）
				case "Kc":
					Wa.Kc = dt
				case "Esu":
					Wa.dblEsu = dt // 通気層内上側の放射率
				case "Esd":
					Wa.dblEsd = dt // 通気層内下側の放射率
				case "Eg":
					Wa.Eg = dt // 透過体の中空層側表面の放射率
				case "Eb":
					Wa.Eb = dt // 集熱板の中空層側表面の放射率
				case "ag":
					Wa.ag = dt // 透過体の日射吸収率
				case "ta":
					Wa.ta = dt / 1000.0 // 中空層の厚さ [mm] -> [m]
				case "tnxt":
					Wa.tnxt = dt
				case "t":
					Wa.air_layer_t = dt / 1000.0
				case "KHD":
					Wa.PVwallcat.KHD = dt // 日射量年変動補正係数 (集熱板が太陽電池一体型のとき)
				case "KPD":
					Wa.PVwallcat.KPD = dt // 経時変化補正係数 (集熱板が太陽電池一体型のとき)
				case "KPM":
					Wa.PVwallcat.KPM = dt // アレイ負荷整合補正係数 (集熱板が太陽電池一体型のとき)
				case "KPA":
					Wa.PVwallcat.KPA = dt // アレイ回路補正係数 (集熱板が太陽電池一体型のとき)
				case "EffInv":
					Wa.PVwallcat.EffINO = dt // インバータ効率 (集熱板が太陽電池一体型のとき)
				case "apamax":
					Wa.PVwallcat.Apmax = dt // 最大出力温度係数 (集熱板が太陽電池一体型のとき)
				case "ap":
					Wa.PVwallcat.Ap = dt // 太陽電池裏面の対流熱伝達率
				case "Rcoloff":
					Wa.PVwallcat.Rcoloff = dt // 太陽電池から集熱器裏面までの熱抵抗 (集熱板が太陽電池一体型のとき)
					Wa.PVwallcat.Kcoloff = 1. / Wa.PVwallcat.Rcoloff
				default:
					Eprint("<Walldata>", s)
				}
			} else {
				layer = append(layer, s)
				j++
			}
		}

		// 層の数
		Nlyr = j + 1

		Wa.welm = []WELM{
			{
				Code: "ali", //内表面熱伝達率
				L:    -999.0,
				ND:   0,
				Cond: -999.0,
				Cro:  -999.0,
			},
		}

		// (3) 層の読み取り
		// `APR-100 APR-100/20 <P:80.3>`
		for j = 1; j <= Nlyr; j++ {
			if Wa.ble == 'R' || Wa.ble == 'c' {
				jj = Nlyr - j
			} else {
				jj = j - 1
			}

			// `APR-100`, `APR-100/20` or `<P>:80.3` or `<C>`
			var err error
			st := strings.IndexRune(layer[jj], '-')
			if st != -1 {
				// 一般材料のとき または 一般材料で分割層数を指定するとき
				// For `APR-100` or `APR-100/20`
				var code string // 材料コード
				var ND int      // 内部分割数
				var L float64   // 部材厚さ [mm]

				// 1) 材料コード
				code = layer[jj][:st]

				// 2) 同一層内の内部分割数と部材の厚さ
				ss := layer[jj][st+1:]
				st = strings.IndexRune(ss, '/')
				if st != -1 {
					// "20.0/2"のように分割数が指定されている場合
					sss := strings.SplitN(ss, "/", 2)
					L, err = strconv.ParseFloat(sss[0], 64)
					if err != nil {
						panic(err)
					}

					// 分割数
					ndiv, err = strconv.Atoi(sss[1])
					if err != nil {
						panic(err)
					}
					ND = ndiv - 1
				} else {
					// "20.0"のように分割数が指定されていない場合
					L, err = strconv.ParseFloat(ss, 64)
					if err != nil {
						panic(err)
					}

					if L >= 50. {
						ND = int((L - 50.) / 50.)
					} else {
						// 50mm未満の場合は分割しない
						ND = 0
					}
				}

				Wa.welm = append(Wa.welm, WELM{
					Code: code,
					ND:   ND,
					L:    L / 1000.0, // [mm] -> [m]
				})

			} else if strings.HasPrefix(layer[jj], "<P") || layer[jj] == "<C>" {
				// 放射暖冷房パネル発熱面位置の指定の場合 For `<P:80.3>` or `<P>` or `<C>`
				// ※ `<C>`はマニュアルに記述が見当たらないが、`<P>`と同じ扱いにする
				Wa.Ip = len(Wa.welm) - 1
				if layer[jj][2] == ':' {
					// パネル効率が指定されている場合 `<P:80.3>`
					Wa.effpnl, err = strconv.ParseFloat(layer[jj][3:], 64)
					if err != nil {
						panic(err)
					}
				} else {
					// パネル効率が指定されない場合 `<C>`
					Wa.effpnl = 0.7
				}
			} else if layer[jj] != "alo" && layer[jj] != "ali" {
				// 表面熱伝達率、中空層熱コンダクタンスのとき、材料コードのみを指定する
				Wa.welm = append(Wa.welm, WELM{
					Code: layer[jj],
					ND:   0,
					L:    0.0,
				})
			}
		}
		jw = len(Wa.welm) - 1

		// 建材一体型空気集熱器の総合熱貫流率の計算、データ入力状況のチェック
		if Wa.Ip >= 0 {
			if Wa.tra > 0. {
				// 壁種類 -> 建材一体型空気集熱器
				Wa.WallType = 'C'

				if (Wa.Ksu > 0. && Wa.Ksd > 0.) || (Wa.Rd > 0. && (Wa.Ru >= 0. || Wa.ta > 0.)) {
					if strings.HasPrefix(Wa.ColType, "A") {
						if Wa.Ksu > 0. {
							Wa.Kcu = Wa.fcu * Wa.Ksu
							Wa.Kcd = Wa.fcd * Wa.Ksd
							Wa.Kc = Wa.Kcu + Wa.Kcd
							Wa.ku = Wa.Kcu / Wa.Kc
							Wa.kd = Wa.Kcd / Wa.Kc
						} else {
							Wa.chrRinput = 'Y'
						}
					} else if strings.HasPrefix(Wa.ColType, "W") {
						Wa.Ko = Wa.Ksu + Wa.Ksd
						Wa.ku = Wa.Ksu / Wa.Ko
						Wa.kd = Wa.Ksd / Wa.Ko
					}

					if Wa.PVwallcat.KHD > 0. {
						PVwallPreCalc(&Wa.PVwallcat)
					}
				} else {
					s = fmt.Sprintf("ble=%c name=%s 建築一体型空気集熱の熱貫流率Ku、Kdが未定義です", Wa.ble, Wa.name)
					Eprint("<Walldata>", s)
				}

				if Wa.chrRinput == 'N' && (Wa.Kc < 0. || Wa.Ksu < 0. || Wa.Ksd < 0.) {
					s = fmt.Sprintf("ble=%c name=%s 建築一体型空気集熱の熱貫流率Kc、Kdd、Kudが未定義です",
						Wa.ble, Wa.name)
					Eprint("<Walldata>", s)
				}

				if Wa.Ip == -1 {
					s = fmt.Sprintf("ble=%c name=%s 建築一体型空気集熱の空気流通層<P>が未定義です",
						Wa.ble, Wa.name)
					Eprint("<Walldata>", s)
				}
			} else {
				// 壁種類 -> 床暖房等放射パネル
				Wa.WallType = 'P'
			}
		}

		Wa.N = jw + 1
		Walli(Nbm, W, Wa, pcm, Npcm)

		*Wall = append(*Wall, *Wa)
		i++
	}

	*Nwall = i
	(*Wall)[0].end = i
}

/* ------------------------------------------------ */

/*  窓デ－タの入力     */

func Windowdata(section *EeTokens, dsn string, Window *[]WINDOW, Nwindow *int) {
	E := fmt.Sprintf(ERRFMT, dsn)

	var N int
	for section.IsEnd() == false {
		line := section.GetLogicalLine()
		if line[0] != "*" {
			N++
		}
	}
	section.Reset()

	if N > 0 {
		*Window = make([]WINDOW, N)

		for j := 0; j < N; j++ {
			(*Window)[j] = *NewWINDOW()
		}
	}

	i := 0
	j := 0
	for section.IsEnd() == false {
		line := section.GetLogicalLine()

		if line[0] == "*" {
			break
		}

		W := &(*Window)[j]

		// 名称
		W.Name = line[0]

		// 名前の重複確認
		for k := 0; k < i; k++ {
			Wc := &(*Window)[k]
			if W.Name == Wc.Name {
				ss := fmt.Sprintf("<WINDOW>  WindowName Already Defined  (%s)", W.Name)
				Eprint("<Windowdata>", ss)
			}
		}

		// プロパティの設定
		for _, s := range line[1:] {
			// 室内透過日射が窓室内側への入射日射を屋外に透過する場合'y'
			if s == "-RStrans" {
				W.RStrans = 'y'
			} else {
				//キー・バリューの分離
				st := strings.IndexRune(s, '=')
				if st == -1 {
					panic("Windowdata: invalid format")
				} else {
					st++
				}
				key := s[:st-1]
				value := s[st+1:]

				// 小数読み取り
				var realValue float64
				var err error
				switch key {
				case "t", "B", "R", "Ei", "Eo":
					realValue, err = strconv.ParseFloat(value, 64)
					if err != nil {
						panic(err)
					}
				}

				// 値の設定
				switch key {
				case "t":
					W.tgtn = realValue // 日射透過率
				case "B":
					W.Bn = realValue // 吸収日射取得率
				case "R":
					W.Rwall = realValue // 窓部材熱抵抗 [m2K/W]
				case "Ei":
					W.Ei = realValue // 室内表面放射率(0.9)
				case "Eo":
					W.Eo = realValue // 外表面放射率(0.9)
				case "Cidtype":
					W.Cidtype = strings.Trim(value, "'") // 入射角特性の種類
				default:
					Err := fmt.Sprintf("%s %s\n", E, s)
					Eprint("<Windowdata>", Err)
				}

				//NOTE: 以下の項目を入力する箇所が不明
				// 窓ガラス面積 Ag, 開口面積 Ao, 幅 W, 高さ H, ??? K
			}
		}

		i++
	}

	*Nwindow = i
	(*Window)[0].end = i
}

/* --------------------------------------------------- */

func Snbkdata(section *EeTokens, dsn string, Snbk *[]SNBK) {
	// 入力チェック用パターン文字列
	typstr := []string{
		"HWDTLR.", // 庇
		"HWDTLRB", // 袖壁(その1)　(左右)
		"HWDTL.B", // 袖壁(その2)　(左のみ)
		"HWDT.RB", // 袖壁(その3)　(右のみ)
		"HWDT...", // 長い庇
		"HWD.LR.", // 長い袖壁(その1) (左右)
		"HWD.L..", // 長い袖壁(その2) (左のみ)
		"HWD..R.", // 長い袖壁(その3) (右のみ)
		"HWDTLRB", // ルーバー
	}

	Er := fmt.Sprintf(ERRFMT, dsn)
	Type := 0

	var N int
	for N = 0; section.IsEnd() == false; N++ {
		section.GetLogicalLine()
	}
	section.Reset()

	if N > 0 {
		*Snbk = make([]SNBK, N)

		for j := 0; j < N; j++ {
			S := &(*Snbk)[j]
			S.Name = ""
			S.W = 0.0
			S.H = 0.0
			S.D = 0.0
			S.W1 = 0.0
			S.W2 = 0.0
			S.H1 = 0.0
			S.H2 = 0.0
			S.end = 0
			S.Type = 0
			S.Ksi = 0
		}
	}

	S := &(*Snbk)[0]
	i := 0
	for section.IsEnd() == false {

		fields := section.GetLogicalLine()

		// 名前
		S.Name = fields[0]

		// 入力チェック用
		code := [8]rune{'.', '.', '.', '.', '.', '.', '.', '.'}

		for _, s := range fields[1:] {
			// キー・バリューの分離
			st := strings.IndexRune(s, '=')
			if st == -1 {
				panic("Snbkdata: invalid format")
			}
			key := s[:st]
			v := s[st+1:]

			var err error
			switch key {
			case "type":
				var vs string
				if v[0] == '-' {
					S.Ksi = 1
					vs = v[1:]
				} else {
					S.Ksi = 0
					vs = v[0:]
				}

				if vs == "H" {
					Type = 1 // 一般の庇
				} else if vs == "HL" {
					Type = 5 // 長い庇
				} else if vs == "S" {
					Type = 2 // 袖壁
				} else if vs == "SL" {
					Type = 6 // 長い袖壁
				} else if vs == "K" {
					Type = 9 // 格子ルーバー
				} else {
					E := fmt.Sprintf("`%s` is invalid", vs)
					Eprint("<Snbkdata>", E)
				}

			case "window":
				// For `window=HhhhxWwww`
				hw := strings.Split(v, "x")
				h, w := hw[0], hw[1]

				if h[0] != 'H' || w[0] != 'W' {
					panic(fmt.Sprintf("Invaid window format: %s", v))
				}

				// 開口部の高さ
				S.H, err = strconv.ParseFloat(h[1:], 64)
				if err != nil {
					panic(err)
				}
				code[0] = 'H'

				// 開口部の幅
				S.W, err = strconv.ParseFloat(w[1:], 64)
				if err != nil {
					panic(err)
				}
				code[1] = 'W'

			case "D":
				// 庇の付け根から先端までの長さ
				S.D, err = strconv.ParseFloat(v, 64)
				if err != nil {
					panic(err)
				}
				code[2] = 'D'

			case "T":
				// 開口部の上端から壁の上端までの距離
				S.H1, err = strconv.ParseFloat(v, 64)
				if err != nil {
					panic(err)
				}
				code[3] = 'T'

			case "L":
				// 開口部の左端から壁の左端までの距離
				S.W1, err = strconv.ParseFloat(v, 64)
				if err != nil {
					panic(err)
				}
				code[4] = 'L'

			case "R":
				// 開口部の右端から壁の右端までの距離
				S.W2, err = strconv.ParseFloat(v, 64)
				if err != nil {
					panic(err)
				}
				code[5] = 'R'

			case "B":
				// 地面から開口部の下端までの高さ
				S.H2, err = strconv.ParseFloat(v, 64)
				if err != nil {
					panic(err)
				}
				code[6] = 'B'

			default:
				panic(fmt.Sprintf("Invaid window format: %s", s))
			}
		}

		// 日除けの種類
		S.Type = Type

		// 日除けの種類ごとに入力チェック
		switch Type {
		case 1, 5, 9:
			// 庇 or ルーバー
			if string(code[:]) != typstr[Type-1] {
				E := fmt.Sprintf("%s %s  type=%d %s\n", Er, fields[0], Type, string(code[:]))
				Eprint("<Snbkdata>", E)
			}
		case 2, 6:
			// 袖壁
			if string(code[:]) != typstr[Type-1] {
				for j := 1; j < 3; j++ {
					if string(code[:]) == typstr[Type+j-1] {
						S.Type = Type + j
						break
					}
					if j == 3 {
						E := fmt.Sprintf("%s %s  type=%d %s\n", Er, fields[0], Type, string(code[:]))
						Eprint("<Snbkdata>", E)
					}
				}
			}
		}

		S = &(*Snbk)[i]
	}

	(*Snbk)[0].end = N
}

/************************************************************/

func wbmlistcount(fi string) []string {
	N := make([]string, 0)

	reader := strings.NewReader(fi)
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		s := scanner.Text()

		if strings.HasPrefix(s, "*") {
			break
		}

		//コメントを読み飛ばす
		st := strings.IndexRune(s, '!')
		if st != -1 {
			s = s[:st]
		}

		//ゴミ除去
		s = strings.TrimSpace(s)

		if s != "" {
			N = append(N, s)
		}
	}

	return N
}

/************************************************************/

func Wallcount(scanner *bufio.Scanner) []string {
	N := make([]string, 0)

	for scanner.Scan() {
		s := scanner.Text()

		if strings.HasPrefix(s, "*") {
			break
		}

		//コメントを読み飛ばす
		st := strings.IndexRune(s, '!')
		if st != -1 {
			s = s[:st]
		}

		//ゴミ除去
		s = strings.TrimSpace(s)

		if s != "" {
			N = append(N, s)
		}
	}

	return N
}

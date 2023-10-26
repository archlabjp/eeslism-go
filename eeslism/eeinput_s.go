package eeslism

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type EeTokens struct {
	tokens []string
	pos    int
}

// Get length
func (t *EeTokens) Len() int {
	return len(t.tokens)
}

// Get position
func (t *EeTokens) GetPos() int {
	return t.pos
}

// Restore position
func (t *EeTokens) RestorePos(pos int) {
	t.pos = pos
}

// Reset position
func (t *EeTokens) Reset() {
	t.pos = 0
}

// Create EeTokes from string
func NewEeTokens(s string) *EeTokens {
	reader := strings.NewReader(s)
	scanner := bufio.NewScanner(reader)
	tokens := make([]string, 0)
	for scanner.Scan() {
		//行単位の処理
		line := scanner.Text()

		// コメントの除去
		if st := strings.IndexRune(line, '!'); st != -1 {
			line = line[:st]
		}

		// 空文字の除去
		line = strings.TrimSpace(line)

		for _, s := range strings.Fields(line) {
			if strings.HasSuffix(s, ";") {
				s = s[:len(s)-1]
				if s != "" {
					tokens = append(tokens, s)
				}
				tokens = append(tokens, ";")
			} else if strings.ContainsRune(s, ';') {
				panic("Invalid position of `;`")
			} else {
				tokens = append(tokens, s)
			}
		}

		//改行
		tokens = append(tokens, "\n")
	}
	return &EeTokens{tokens: tokens, pos: 0}
}

// Return tokens from current position to `\n`
func (t *EeTokens) GetLine() []string {
	var line []string

	// find `\n`
	var found bool = false
	for i := t.pos; i < len(t.tokens); i++ {
		if t.tokens[i] == "\n" {
			line = t.tokens[t.pos:i]
			t.pos = i + 1
			found = true
			break
		}
	}
	// not found
	if found == false {
		t.pos = len(t.tokens)
		line = t.tokens[t.pos:]
	}

	return line
}

// Return tokens from current position to `;`
func (t *EeTokens) GetLogicalLine() []string {
	var logiline []string
	var filtered []string

	// find `;`
	var found bool = false
	for i := t.pos; i < len(t.tokens); i++ {
		if t.tokens[i] == ";" {
			logiline = t.tokens[t.pos:i] // `;` is not included
			t.pos = i + 1
			found = true
			break
		}
	}
	// not found
	if found == false {
		logiline = t.tokens[t.pos:]
		t.pos = len(t.tokens)
	}

	// filter `\n` token and return it
	for _, s := range logiline {
		if s != "\n" {
			filtered = append(filtered, s)
		}
	}
	return filtered
}

// Skip `;` or `\n`
func (t *EeTokens) SkipToEndOfLine() {
	for t.pos < len(t.tokens) && (t.tokens[t.pos] == "\n" || t.tokens[t.pos] == ";") {
		t.pos++
	}
}

// Return tokens from current position to `*`
func (t *EeTokens) GetSection() *EeTokens {
	t.SkipToEndOfLine()

	// find `*` at start of some line
	for i := t.pos; i < len(t.tokens); i++ {
		if i > 0 && t.tokens[i-1] == "\n" && t.tokens[i] == "*" {
			section := &EeTokens{tokens: t.tokens[t.pos : i+1], pos: 0}
			t.pos = i + 1
			return section
		}
	}
	// not found
	section := &EeTokens{tokens: t.tokens[t.pos:], pos: 0}
	t.pos = len(t.tokens)
	return section
}

// Check if pos is at the end of tokens
func (t *EeTokens) IsEnd() bool {
	return t.pos >= len(t.tokens)
}

// Peek next token
func (t *EeTokens) PeekToken() string {
	if t.pos < len(t.tokens) {
		return t.tokens[t.pos]
	}
	return ""
}

// Get next token
func (t *EeTokens) GetToken() string {
	if t.pos < len(t.tokens) {
		t.pos++
		return t.tokens[t.pos-1]
	}
	return ""
}

// Get next token as float64 value
func (t *EeTokens) GetFloat() float64 {
	var f float64
	fmt.Sscanf(t.GetToken(), "%f", &f)
	return f
}

// Get next token as int value
func (t *EeTokens) GetInt() int {
	var i int
	fmt.Sscanf(t.GetToken(), "%d", &i)
	return i
}

/*  建築・設備システムデータ入力  */

func Eeinput(Ipath string, bdata, week, schtba, schnma string, Simc *SIMCONTL,
	Exsf *EXSFS, Rmvls *RMVLS, Eqcat *EQCAT, Eqsys *EQSYS,
	Compnt *[]*COMPNT,
	Elout *[]*ELOUT, Nelout *int,
	Elin *[]*ELIN, Nelin *int,
	Mpath *[]*MPATH, Nmpath *int,
	Plist *[]*PLIST, Pelm *[]*PELM, Npelm *int,
	Contl *[]*CONTL, Ncontl *int,
	Ctlif *[]*CTLIF, Nctlif *int,
	Ctlst *[]*CTLST, Nctlst *int,
	Wd *WDAT, Daytm *DAYTM, key int, Nplist *int,
	bdpn *int, obsn *int, treen *int, shadn *int, polyn *int,
	bp *[]BBDP, obs *[]OBS, tree *[]TREE, shadtb *[]SHADTB, poly *[]POLYGN, monten *int, gpn *int, DE *float64, Noplpmp *NOPLPMP) (*SCHDL, []*FLOUT) {

	var Twallinit float64
	var j int
	dtm := 3600
	var nday int
	var Nday int
	daystartx := 0
	daystart := 0
	dayend := 0
	var Err, File string

	// 出力フラグ (GDAT.PRINT)
	// 中) 熱負荷要素の出力指定だけ変則的なことに注意
	wdpri := 0  // 気象データの出力指定
	revpri := 0 // 室内熱環境データの出力指定
	pmvpri := 0 // 室内のPMVの出力指定

	Nrmspri := 0 // 表面温度出力指定(室の数)
	Nqrmpri := 0 // 日射、室内発熱取得出力指定(室の数)
	Nwalpri := 0 // 壁体内部温度出力指定(壁体の数)
	Npcmpri := 0 // PCMの状態値出力フラグ(壁体の数)
	Nshdpri := 0 // 日よけの影面積出力 (壁体の数)

	var dfwl DFWL

	/*-------higuchi 070918---------start*/
	//RRMP *rp;
	//MADO *wp;
	//sunblk *sb;
	var smonth, sday, emonth, eday int

	//sb = bp.SBLK;
	//rp = bp.RMP;
	//wp = rp.WD;
	/*-------higuchi------------end*/

	Err = fmt.Sprintf(ERRFMT, "(Eeinput)")

	//*Nexs=0 ;

	Rmvls.Nwall = 0
	Rmvls.Nwindow = 0
	//Rmvls.Nroom = 0

	//Sdd := Rmvls.Sd
	Rmvls.Nsrf = 0
	// Rmvls.Nrdpnl = 0
	Rmvls.Nmwall = 0

	var err error

	// -------------------------------------------------------
	// 曜日設定ファイルの読み取り
	// -------------------------------------------------------
	var fi_dayweek *os.File
	if fi_dayweek, err = os.Open("dayweek.efl"); err != nil {
		Eprint("<Eeinput>", "dayweek.efl")
		os.Exit(EXIT_DAYWEK)
	}
	Dayweek(fi_dayweek, week, Simc.Daywk, key)
	fi_dayweek.Close()

	if DEBUG {
		dprdayweek(Simc.Daywk)
	}

	// -------------------------------------------------------
	// スケジュ－ル表の読み取り
	// -------------------------------------------------------
	var Schdl *SCHDL = new(SCHDL)
	Schtable(schtba, Schdl)
	Schname(Schdl)

	// -------------------------------------------------------
	//  季節、曜日によるスケジュ－ル表の組み合わせの読み取り
	// -------------------------------------------------------
	Schdata(schnma, "schnm", Simc.Daywk, Schdl)

	// 入力を正規化することで後処理を簡単にする
	tokens := NewEeTokens(bdata)

	for tokens.IsEnd() == false {
		s := tokens.GetToken()
		if s == "\n" || s == ";" || s == "*" {
			continue
		}
		fmt.Printf("=== %s\n", s)

		switch s {
		case "TITLE":
			line := tokens.GetLogicalLine()
			Simc.Title = strings.Join(line, " ")
			fmt.Printf("%s\n", Simc.Title)
		case "GDAT":
			section := tokens.GetSection()
			Wd.RNtype = 'C'
			Wd.Intgtsupw = 'N'
			Simc.Perio = 'n' // 周期定常計算フラグを'n'に初期化
			Gdata(section, s, Simc.File, &Simc.Wfname, &Simc.Ofname, &dtm, &Simc.Sttmm,
				&daystartx, &daystart, &dayend, &Twallinit, Simc.Dayprn,
				&wdpri, &revpri, &pmvpri, &Simc.Helmkey, &Simc.MaxIterate, Daytm, Wd, &Simc.Perio)

			// 気象データファイル名からファイル種別を判定
			if Simc.Wfname == "" {
				Simc.Wdtype = 'E'
			} else {
				Simc.Wdtype = 'H'
			}

			// 初期温度 (15[deg])
			Rmvls.Twallinit = Twallinit

			// 計算時間間隔 [s]
			Simc.DTm = dtm

			Simc.Unit = "t_C x_kg/kg r_% q_W e_W"
			Simc.Unitdy = "Q_kWh E_kWh"

			fmt.Printf("== File  Output=%s\n", Simc.Ofname)
		case "SCHTB":
			// SCHDBデータセットの読み取り
			Schtable(schtba, Schdl)
			Schname(Schdl)
		case "SCHNM":
			// SCHNMデータセットの読み取り
			Schdata(schnma, s, Simc.Daywk, Schdl)
		case "EXSRF":
			// EXSRFデータセットの読み取り
			section := tokens.GetSection()
			Exsfdata(section, s, Exsf, Schdl, Simc)

		case "SUNBRK":
			// 日よけの読み込み
			section := tokens.GetSection()
			Snbkdata(section, s, &Rmvls.Snbk)

		case "PCM":
			section := tokens.GetSection()
			PCMdata(section, s, &Rmvls.PCM, &Rmvls.Npcm, &Rmvls.Pcmiterate)

		case "WALL":
			if Fbmlist == "" {
				File = "wbmlist.efl"
			} else {
				File = Fbmlist
			}

			var fbmContent []byte
			if fbmContent, err = ioutil.ReadFile(File); err != nil {
				Eprint("<Eeinput>", "wbmlist.efl")
				os.Exit(EXIT_WBMLST)
			}
			/*******************/

			section := tokens.GetSection()
			Walldata(section, string(fbmContent), s, &Rmvls.Wall, &Rmvls.Nwall, &dfwl, Rmvls.PCM, Rmvls.Npcm)

		case "WINDOW":
			section := tokens.GetSection()
			Windowdata(section, s, &Rmvls.Window, &Rmvls.Nwindow)

		case "ROOM":
			Roomdata(tokens, "Roomdata", Exsf.Exs, &dfwl, Rmvls, Schdl, Simc)
			Balloc(Rmvls.Nsrf, Rmvls.Sd, Rmvls.Wall, &Rmvls.Mw, &Rmvls.Nmwall)

		case "RAICH", "VENT":
			section := tokens.GetSection()
			Ventdata(section, s, Schdl, Rmvls.Room, Simc)

		case "RESI":
			section := tokens.GetSection()
			Residata(section, s, Schdl, Rmvls.Room, &pmvpri, Simc)

		case "APPL":
			section := tokens.GetSection()
			Appldata(section, s, Schdl, Rmvls.Room, Simc)
		case "VCFILE":
			section := tokens.GetSection()
			Vcfdata(section, Simc)
		case "EQPCAT":
			section := tokens.GetSection()
			Eqcadata(section, "Eqcadata", Eqcat)

		case "SYSCMP": // 接続用のノードを設定している
			/*****Flwindata(Flwin, Nflwin,  Wd);********/
			section := tokens.GetSection()
			Compodata(section, "Compodata", Rmvls, Eqcat, Compnt, Eqsys)
			Elmalloc("Elmalloc ", *Compnt, Eqcat, Eqsys, Elout, Elin)

		case "SYSPTH": // 接続パスの設定をしている
			section := tokens.GetSection()
			Pathdata(section, "Pathdata", Simc, Wd, *Compnt, Schdl, Mpath, Nmpath, Plist, Pelm, Eqsys)
			Roomelm(Rmvls.Room, Rmvls.Rdpnl)

			// 変数の割り当て
			Hclelm(Eqsys.Hcload)
			Thexelm(Eqsys.Thex)
			Desielm(Eqsys.Desi)
			Evacelm(Eqsys.Evac)

			Qmeaselm(Eqsys.Qmeas)

		case "CONTL":
			section := tokens.GetSection()
			Contrldata(section, Contl, Ctlif, Ctlst, Simc, *Compnt, *Mpath, Wd, Exsf, Schdl)

		/*--------------higuchi add-------------------start*/

		// 20170503 higuchi add
		case "DIVID":
			section := tokens.GetSection()
			dividdata(section, monten, DE)

		/*--対象建物データ読み込み--*/
		case "COORDNT":
			section := tokens.GetSection()
			bdpdata(section, bdpn, bp, Exsf)

		/*--障害物データ読み込み--*/
		case "OBS":
			section := tokens.GetSection()
			obsdata(section, obsn, obs)

		/*--樹木データ読み込み--*/
		case "TREE":
			section := tokens.GetSection()
			treedata(section, treen, tree)

		/*--多角形障害物直接入力分の読み込み--*/
		case "POLYGON":
			section := tokens.GetSection()
			polydata(section, polyn, poly)

		/*--落葉スケジュール読み込み--*/
		case "SHDSCHTB":
			*shadn = 0
			var Nshadn int
			//var shdp *SHADTB
			// 落葉スケジュールの数を数える
			section := tokens.GetSection()
			Nshadn = InputCount(section, ";")
			section.Reset()

			if Nshadn > 0 {
				*shadtb = make([]SHADTB, Nshadn)
			}

			i := 0
			for section.IsEnd() == false {
				s = section.GetToken()
				if s[0] == '*' {
					break
				}
				shdp := (*shadtb)[i]
				shdp.lpname = s
				(*shadn)++
				shdp.indatn = 0

				for section.IsEnd() == false {
					s = section.GetToken()
					if s[0] == '*' {
						break
					}
					_, err = fmt.Sscanf(s, "%d/%d-%f-%d/%d", &smonth, &sday, &shdp.shad[shdp.indatn], &emonth, &eday)
					if err != nil {
						panic(err)
					}
					shdp.ndays[shdp.indatn] = nennkann(smonth, sday)
					shdp.ndaye[shdp.indatn] = nennkann(emonth, eday)
					shdp.indatn = shdp.indatn + 1
				}
				i++
			}

		/*----------higuchi add-----------------end-*/

		default:
			Err = Err + "  " + s
			Eprint("<Eeinput>", Err)
		}
	}

	/*--------------higuchi 070918-------------------start-*/
	if *bdpn != 0 {
		fmt.Printf("deviding of wall mm: %f\n", *DE)
		fmt.Printf("number of point in montekalro: %d\n", *monten)
	}
	/*----------------higuchi 7.11,061123------------------end*/

	// 外部障害物の数を数える
	Noplpmp.Nop = OPcount(*bdpn, *bp, *polyn, *poly)
	Noplpmp.Nlp = LPcount(*bdpn, *bp, *obsn, *obs, *treen, *polyn, *poly)
	Noplpmp.Nmp = Noplpmp.Nop + Noplpmp.Nlp

	//////////////////////////////////////

	//----------------------------------------------------
	// シミュレーション設定
	//----------------------------------------------------

	if daystart > dayend {
		dayend = dayend + 365
	}
	Nday = dayend - daystart + 1

	if daystartx > daystart {
		daystart = daystart + 365
	}

	Nday += daystart - daystartx
	Simc.Dayend = daystartx + Nday - 1
	Simc.Daystartx = daystartx
	Simc.Daystart = daystart

	Simc.Timeid = []rune{'M', 'D', 'T'}

	Simc.Ntimedyprt = Simc.Dayend - Simc.Daystart + 1
	Simc.Dayntime = 24 * 3600 / dtm
	Simc.Ntimehrprt = 0

	for nday = Simc.Daystart; nday <= Simc.Dayend; nday++ {
		// NOTE: オリジナルコードはバッファーオーバーランしているので、`%366`を追加
		if Simc.Dayprn[nday%366] != 0 {
			Simc.Ntimehrprt += Simc.Dayntime
		}
	}

	//----------------------------------------------------
	// 出力ファイルの追加
	//----------------------------------------------------

	for i := range Rmvls.Room {
		Rm := &Rmvls.Room[i]
		if Rm.sfpri {
			Nrmspri++
		}
		if Rm.eqpri {
			Nqrmpri++
		}
	}

	for i := 0; i < Rmvls.Nsrf; i++ {
		Sd := &Rmvls.Sd[i]
		if Sd.wlpri {
			Nwalpri++
		}

		if Sd.pcmpri {
			Npcmpri++
		}

		// 日よけの影面積出力
		if Sd.shdpri {
			Nshdpri++
		}
	}

	// 出力ファイルの追加手続き
	var Flout []*FLOUT = make([]*FLOUT, 0, 30) // ファイル出力設定
	addFlout := func(idn PrintType) {
		Flout = append(Flout, &FLOUT{Idn: idn})
	}

	// 必須出力ファイル
	addFlout(PRTPATH)    // 時間別計算値(システム経路の温湿度出力)
	addFlout(PRTCOMP)    // 時間別計算値(機器の出力)
	addFlout(PRTDYCOMP)  // 日別計算値(システム要素機器の日集計結果出力)
	addFlout(PRTMNCOMP)  // 月別計算値(システム要素機器の月集計結果出力)
	addFlout(PRTMTCOMP)  // 月-時刻計算値(部屋ごとの熱集計結果出力)
	addFlout(PRTHRSTANK) // 時間別計算値(蓄熱槽内温度分布の出力)
	addFlout(PRTWK)      // 計算年月日出力
	addFlout(PRTREV)     // 時間別計算値(毎時室温、MRTの出力)
	addFlout(PRTHROOM)   // 時間別計算値(放射パネルの出力)
	addFlout(PRTDYRM)    // 日別計算値(部屋ごとの熱集計結果出力)
	addFlout(PRTMNRM)    // 月別計算値(部屋ごとの熱集計結果出力)

	// 要素別熱損失・熱取得（記憶域確保）
	Helminit("Helminit", Simc.Helmkey, Rmvls.Room, &Rmvls.Qetotal)

	if Simc.Helmkey == 'y' {
		addFlout(PRTHELM)   // 時間別計算値(要素別熱損失・熱取得)
		addFlout(PRTDYHELM) // 日別計算値(要素別熱損失・熱取得)

		Simc.Nhelmsfpri = 0
		for i := range Rmvls.Room {
			Rm := &Rmvls.Room[i]
			for j = 0; j < Rm.N; j++ {
				Sdd := &Rm.rsrf[j]
				if Sdd.sfepri {
					Simc.Nhelmsfpri++
				}
			}
		}
		if Simc.Nhelmsfpri > 0 {
			addFlout(PRTHELMSF) // 時間別計算値(要素別熱損失・熱取得) 表面?
		}
	}

	if pmvpri > 0 {
		addFlout(PRTPMV) // 時間別計算値(PMV計算)
	}

	if Nqrmpri > 0 {
		addFlout(PRTQRM) // 時間別計算値(日射、室内熱取得の出力)
		addFlout(PRTDQR) // 日別計算値(日射、室内熱取得の出力)
	}

	if Nrmspri > 0 {
		addFlout(PRTRSF)  // 時間別計算値(室内表面温度の出力)
		addFlout(PRTSFQ)  // 時間別計算値(室内表面熱流の出力)
		addFlout(PRTSFA)  // 時間別計算値(室内表面熱伝達率の出力)
		addFlout(PRTDYSF) // 日別計算値(日積算壁体貫流熱取得の出力)
	}

	if Nwalpri > 0 {
		addFlout(PRTWAL) // // 時間別計算値(壁体内部温度の出力)
	}

	// 日よけの影面積出力
	if Nshdpri > 0 {
		addFlout(PRTSHD) // 時間別計算値(日よけの影面積の出力)
	}

	// 潜熱蓄熱材がある場合
	if Npcmpri > 0 {
		addFlout(PRTPCM) // 時間別計算値(潜熱蓄熱材の状態値の出力)
	}

	// 気象データの出力を追加
	if wdpri > 0 {
		addFlout(PRTHWD) // 時間別計算値(気象データ出力)
		addFlout(PRTDWD) // 日別計算値(気象データ日集計値出力)
		addFlout(PRTMWD) // 月別計算値(気象データ月集計値出力)
	}

	return Schdl, Flout
}

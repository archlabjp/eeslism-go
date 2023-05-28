package main

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

	// find `*`
	for i := t.pos; i < len(t.tokens); i++ {
		if t.tokens[i] == "*" {
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

// Get next token
func (t *EeTokens) GetToken() string {
	if t.pos < len(t.tokens) {
		t.pos++
		return t.tokens[t.pos-1]
	}
	return ""
}

/*  建築・設備システムデータ入力  */

func Eeinput(Ipath string, Simc *SIMCONTL, Schdl *SCHDL,
	Exsf *EXSFS, Rmvls *RMVLS, Eqcat *EQCAT, Eqsys *EQSYS,
	Compnt *[]COMPNT, Ncompnt *int,
	Ncmpalloc *int,
	Elout *[]*ELOUT, Nelout *int,
	Elin *[]*ELIN, Nelin *int,
	Mpath *[]MPATH, Nmpath *int,
	Plist *[]PLIST, Pelm *[]PELM, Npelm *int,
	Contl *[]CONTL, Ncontl *int,
	Ctlif *[]CTLIF, Nctlif *int,
	Ctlst *[]CTLST, Nctlst *int,
	Flout *[]*FLOUT, Nflout *int, Wd *WDAT, Daytm *DAYTM, key int, Nplist *int,
	bdpn *int, obsn *int, treen *int, shadn *int, polyn *int,
	bp *[]BBDP, obs *[]OBS, tree *[]TREE, shadtb *[]SHADTB, poly *[]POLYGN, monten *int, gpn *int, DE *float64, Noplpmp *NOPLPMP) {

	var fi *os.File
	var flo int
	var Twallinit float64
	var i, j int
	dtm := 3600
	var nday int
	var Nday int
	daystartx := 0
	daystart := 0
	dayend := 0
	var s, Err, File string
	wdpri := 0
	revpri := 0
	pmvpri := 0
	Nrmspri := 0
	Nqrmpri := 0
	Nwalpri := 0
	Npcmpri := 0

	SYSCMP_ID := 0
	SYSPTH_ID := 0

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
	Rmvls.Nroom = 0

	//Sdd := Rmvls.Sd
	Rmvls.Nsrf = 0
	Rmvls.Nrdpnl = 0
	Rmvls.Nmwall = 0

	var err error
	if fi, err = os.Open("dayweek.efl"); err != nil {
		Eprint("<Eeinput>", "dayweek.efl")
		os.Exit(EXIT_DAYWEK)
	}
	Dayweek(fi, Ipath, Simc.Daywk, key)
	fi.Close()

	if DEBUG {
		dprdayweek(Simc.Daywk)
	}

	if fi, err := os.Open(Ipath + "schtba.ewk"); err != nil {
		Eprint("<Eeinput>", "schtba.ewk")
		os.Exit(EXIT_SCHTB)
	} else {
		Schtable(fi, s, Schdl)
		fi.Close()

		Schname(Ipath, "Schname", Schdl)

		Schdl.Nsch = Schdl.Sch[0].end
		Schdl.Nscw = Schdl.Scw[0].end
	}

	if fi, err := os.Open(Ipath + "schnma.ewk"); err != nil {
		Eprint("<Eeinput>", "schnma.ewk")
		os.Exit(EXIT_SCHNM)
	} else {
		Schdata(fi, "schnm", Simc.Daywk, Schdl)
		fi.Close()

		Schdl.Nsch = Schdl.Sch[0].end
		Schdl.Nscw = Schdl.Scw[0].end

		if Schdl.Nsch > 0 {
			Schdl.Val = make([]float64, Schdl.Nsch)
		} else {
			Schdl.Val = nil
		}

		if Schdl.Nscw > 0 {
			Schdl.Isw = make([]rune, Schdl.Nscw)
		} else {
			Schdl.Isw = nil
		}
	}

	bdataBytes, err := ioutil.ReadFile(Ipath + "bdata.ewk")
	if err != nil {
		Eprint("<Eeinput>", "bdata.ewk")
		os.Exit(EXIT_BDATA)
	}

	// 入力を正規化することで後処理を簡単にする
	tokens := NewEeTokens(string(bdataBytes))

	for tokens.IsEnd() == false {
		s := tokens.GetToken()
		if s == "\n" || s == ";" || s == "*" {
			continue
		}
		fmt.Printf("=== %s\n", s)

		switch s {
		case "TITLE":
			line := tokens.GetLogicalLine()
			Simc.Title = line[0]
			fmt.Printf("%s\n", Simc.Title)
		case "GDAT":
			section := tokens.GetSection()
			Wd.RNtype = 'C'
			Wd.Intgtsupw = 'N'
			Simc.Perio = 'n' // 周期定常計算フラグを'n'に初期化
			Gdata(section, s, Simc.File, &Simc.Wfname, &Simc.Ofname, &dtm, &Simc.Sttmm,
				&daystartx, &daystart, &dayend, &Twallinit, Simc.Dayprn,
				&wdpri, &revpri, &pmvpri, &Simc.Helmkey, &Simc.MaxIterate, Daytm, Wd, &Simc.Perio)
			if Simc.Wfname == "" {
				Simc.Wdtype = 'E'
			} else {
				Simc.Wdtype = 'H'
			}
			Rmvls.Twallinit = Twallinit

			Simc.DTm = dtm

			Simc.Unit = "t_C x_kg/kg r_% q_W e_W"
			Simc.Unitdy = "Q_kWh E_kWh"

			fmt.Printf("== File  Output=%s\n", Simc.Ofname)
		case "SCHTB":
			Schtable(fi, s, Schdl)
			Schname(Ipath, "Schname", Schdl)

			Schdl.Nsch = Schdl.Sch[0].end
			Schdl.Nscw = Schdl.Scw[0].end
		case "SCHNM":
			Schdata(fi, s, Simc.Daywk, Schdl)

			Schdl.Nsch = Schdl.Sch[0].end
			Schdl.Nscw = Schdl.Scw[0].end
		case "EXSRF":
			section := tokens.GetSection()
			Exsfdata(section, s, Exsf, Schdl, Simc)

		case "SUNBRK":
			// 日よけの読み込み
			section := tokens.GetSection()
			Snbkdata(section, s, &Rmvls.Snbk)

		case "PCM":
			PCMdata(fi, s, &Rmvls.PCM, &Rmvls.Npcm, &Rmvls.Pcmiterate)

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
			Ventdata(fi, s, Schdl, Rmvls.Room, Simc)

		case "RESI":
			Residata(fi, s, Schdl, Rmvls.Room, &pmvpri, Simc)

		case "APPL":
			Appldata(fi, s, Schdl, Rmvls.Room, Simc)
		case "VCFILE":
			Vcfdata(fi, Simc)
		case "EQPCAT":
			Eqcadata(fi, "Eqcadata", Eqcat)

		case "SYSCMP":
			/*****Flwindata(Flwin, Nflwin,  Wd);********/
			Compodata(fi, "Compodata", Rmvls, Eqcat, Compnt, Ncompnt, Eqsys, Ncmpalloc, 0)
			Elmalloc("Elmalloc ", *Ncompnt, *Compnt, Eqcat, Eqsys,
				Elout, Nelout, Elin, Nelin)
			SYSCMP_ID++

		case "SYSPTH":
			if SYSCMP_ID == 0 {
				Compodata(fi, "Compodata", Rmvls, Eqcat, Compnt, Ncompnt, Eqsys, Ncmpalloc, 1)

				Elmalloc("Elmalloc ", *Ncompnt, *Compnt, Eqcat, Eqsys, Elout, Nelout, Elin, Nelin)
				SYSCMP_ID++
			}
			Pathdata(fi, "Pathdata", Simc, Wd, *Ncompnt, *Compnt, Schdl,
				Mpath, Nmpath, Plist, Pelm, Npelm, Nplist, 0, Eqsys)
			Roomelm(Rmvls.Nroom, Rmvls.Room, Rmvls.Nrdpnl, Rmvls.Rdpnl)

			// 変数の割り当て
			Hclelm(Eqsys.Nhcload, Eqsys.Hcload)
			Thexelm(Eqsys.Nthex, Eqsys.Thex)
			Desielm(Eqsys.Ndesi, Eqsys.Desi)
			Evacelm(Eqsys.Nevac, Eqsys.Evac)

			Qmeaselm(Eqsys.Nqmeas, Eqsys.Qmeas)
			SYSPTH_ID++

		case "CONTL":
			if SYSCMP_ID == 0 {

				Compodata(fi, "Compodata", Rmvls, Eqcat, Compnt, Ncompnt, Eqsys, Ncmpalloc, 1)

				Elmalloc("Elmalloc ", *Ncompnt, *Compnt, Eqcat, Eqsys,
					Elout, Nelout, Elin, Nelin)
				SYSCMP_ID++
			}

			if SYSPTH_ID == 0 {
				Pathdata(fi, "Pathdata", Simc, Wd, *Ncompnt, *Compnt, Schdl,
					Mpath, Nmpath, Plist, Pelm, Npelm, Nplist, 1, Eqsys)

				Roomelm(Rmvls.Nroom, Rmvls.Room, Rmvls.Nrdpnl, Rmvls.Rdpnl)

				Hclelm(Eqsys.Nhcload, Eqsys.Hcload)
				Thexelm(Eqsys.Nthex, Eqsys.Thex)
				Desielm(Eqsys.Ndesi, Eqsys.Desi)
				Evacelm(Eqsys.Nevac, Eqsys.Evac)

				Qmeaselm(Eqsys.Nqmeas, Eqsys.Qmeas)

				SYSPTH_ID++
			}

			Contrldata(fi, Contl, Ncontl, Ctlif, Nctlif, Ctlst, Nctlst,
				Simc, *Ncompnt, *Compnt, *Nmpath, *Mpath, Wd, Exsf, Schdl)

		/*--------------higuchi add-------------------start*/

		// 20170503 higuchi add
		case "DIVID":
			dividdata(fi, monten, DE)

		/*--対象建物データ読み込み--*/
		case "COORDNT":
			bdpdata(fi, bdpn, bp, Exsf)

		/*--障害物データ読み込み--*/
		case "OBS":
			obsdata(fi, obsn, obs)

		/*--樹木データ読み込み--*/
		case "TREE":
			treedata(fi, treen, tree)

		/*--多角形障害物直接入力分の読み込み--*/
		case "POLYGON":
			polydata(fi, polyn, poly)

		/*--落葉スケジュール読み込み--*/
		case "SHDSCHTB":
			*shadn = 0
			var Nshadn int
			//var shdp *SHADTB
			// 落葉スケジュールの数を数える
			Nshadn = InputCount(fi, ";")

			if Nshadn > 0 {
				*shadtb = make([]SHADTB, Nshadn)
			}

			i := 0
			for {
				_, err := fmt.Fscanf(fi, "%s", s)
				if err != nil {
					panic(err)
				}
				if s[0] == '*' {
					break
				}
				shdp := (*shadtb)[i]
				shdp.lpname = s
				(*shadn)++
				shdp.indatn = 0

				for {
					_, err := fmt.Fscanf(fi, "%s", &s)
					if err != nil {
						panic(err)
					}
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

	fi.Close()

	// 外部障害物の数を数える
	Noplpmp.Nop = OPcount(*bdpn, *bp, *polyn, *poly)
	Noplpmp.Nlp = LPcount(*bdpn, *bp, *obsn, *obs, *treen, *polyn, *poly)
	Noplpmp.Nmp = Noplpmp.Nop + Noplpmp.Nlp

	//////////////////////////////////////

	if SYSCMP_ID == 0 {
		Compodata(fi, "Compodata", Rmvls, Eqcat, Compnt, Ncompnt, Eqsys, Ncmpalloc, 1)

		Elmalloc("Elmalloc ", *Ncompnt, *Compnt, Eqcat, Eqsys, Elout, Nelout, Elin, Nelin)
	}

	if SYSPTH_ID == 0 {
		Pathdata(fi, "Pathdata", Simc, Wd, *Ncompnt, *Compnt, Schdl,
			Mpath, Nmpath, Plist, Pelm, Npelm, Nplist, 1, Eqsys)

		Roomelm(Rmvls.Nroom, Rmvls.Room, Rmvls.Nrdpnl, Rmvls.Rdpnl)

		Hclelm(Eqsys.Nhcload, Eqsys.Hcload)
		Thexelm(Eqsys.Nthex, Eqsys.Thex)

		Qmeaselm(Eqsys.Nqmeas, Eqsys.Qmeas)
	}

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

	for i := range Rmvls.Room {
		Rm := &Rmvls.Room[i]
		if Rm.sfpri == 'p' {
			Nrmspri++
		}
		if Rm.eqpri == 'p' {
			Nqrmpri++
		}
	}

	var Nshdpri int
	Nshdpri = 0
	Nwalpri = 0
	for i := 0; i < Rmvls.Nsrf; i++ {
		Sd := &Rmvls.Sd[i]
		if Sd.wlpri == 'p' {
			Nwalpri++
		}

		if Sd.pcmpri == 'y' {
			Npcmpri++
		}

		// 日よけの影面積出力
		if Sd.shdpri == 'p' {
			Nshdpri++
		}
	}

	flo = 0

	(*Flout)[flo].Idn = PRTPATH
	flo++
	(*Flout)[flo].Idn = PRTCOMP
	flo++
	(*Flout)[flo].Idn = PRTDYCOMP
	flo++
	(*Flout)[flo].Idn = PRTMNCOMP
	flo++
	(*Flout)[flo].Idn = PRTMTCOMP
	flo++
	(*Flout)[flo].Idn = PRTHRSTANK
	flo++
	(*Flout)[flo].Idn = PRTWK
	flo++
	(*Flout)[flo].Idn = PRTREV
	flo++
	(*Flout)[flo].Idn = PRTHROOM
	flo++
	(*Flout)[flo].Idn = PRTDYRM
	flo++
	(*Flout)[flo].Idn = PRTMNRM
	flo++

	Helminit("Helminit", Simc.Helmkey,
		Rmvls.Nroom, Rmvls.Room, &Rmvls.Qetotal)

	if Simc.Helmkey == 'y' {
		(*Flout)[flo].Idn = PRTHELM
		flo++

		(*Flout)[flo].Idn = PRTDYHELM
		flo++

		Simc.Nhelmsfpri = 0
		for i = 0; i < Rmvls.Nroom; i++ {
			Rm := &Rmvls.Room[i]
			for j = 0; j < Rm.N; j++ {
				Sdd := &Rm.rsrf[j]
				if Sdd.sfepri == 'y' {
					Simc.Nhelmsfpri++
				}
			}
		}
		if Simc.Nhelmsfpri > 0 {
			(*Flout)[flo].Idn = PRTHELMSF
			flo++
		}
	}

	if pmvpri > 0 {
		(*Flout)[flo].Idn = PRTPMV
		flo++
	}

	if Nqrmpri > 0 {
		(*Flout)[flo].Idn = PRTQRM
		flo++

		(*Flout)[flo].Idn = PRTDQR
		flo++
	}

	if Nrmspri > 0 {
		(*Flout)[flo].Idn = PRTRSF
		flo++
		(*Flout)[flo].Idn = PRTSFQ
		flo++
		(*Flout)[flo].Idn = PRTSFA
		flo++
		(*Flout)[flo].Idn = PRTDYSF
		flo++
	}

	if Nwalpri > 0 {
		(*Flout)[flo].Idn = PRTWAL
		flo++
	}

	// 日よけの影面積出力
	if Nshdpri > 0 {
		(*Flout)[flo].Idn = PRTSHD
		flo++
	}

	if Npcmpri > 0 {
		(*Flout)[flo].Idn = PRTPCM
		flo++
	}

	if wdpri > 0 {
		(*Flout)[flo].Idn = PRTHWD
		flo++

		(*Flout)[flo].Idn = PRTDWD
		flo++

		(*Flout)[flo].Idn = PRTMWD
		flo++
	}

	*Nflout = flo
}

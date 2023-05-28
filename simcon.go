package main

import "os"

const EEVERSION = "ES4.6"

type PrintType string

const (
	PRTREV   PrintType = "_re"
	PRTHROOM PrintType = "_rm"

	PRTPMV  PrintType = "_pm"
	PRTQRM  PrintType = "_rq"
	PRTRSF  PrintType = "_sf"
	PRTSFQ  PrintType = "_sfq"
	PRTSFA  PrintType = "_sfa"
	PRTDYSF PrintType = "_dsf"
	PRTWAL  PrintType = "_wl"
	PRTSHD  PrintType = "_shd"
	PRTPCM  PrintType = "_pcm"

	PRTPATH PrintType = "_sp"
	PRTCOMP PrintType = "_sc"

	PRTHRSTANK PrintType = "_tk"

	PRTHWD PrintType = "_wd"

	PRTDYRM   PrintType = "_dr"
	PRTMNRM   PrintType = "_mr"
	PRTDYCOMP PrintType = "_dc"
	PRTMNCOMP PrintType = "_mc"
	PRTMTCOMP PrintType = "_mt"
	PRTDQR    PrintType = "_dqr"
	PRTDWD    PrintType = "_dwd"
	PRTMWD    PrintType = "_mwd"

	PRTWK = "_wk"

	PRTHELM   = "_rqe"
	PRTHELMSF = "_sfe"
	PRTDYHELM = "_dqe"

	// SYSV_EQV = 'v'
	// LOAD_EQV = 'L'
)

type SIMCONTL struct {
	File       string
	Title      string
	Wfname     string
	Ofname     string
	Unit       string
	Unitdy     string
	Timeid     []rune
	Helmkey    rune // 要素別熱取得、熱損失計算 'y'
	Wdtype     rune // 気象データファイル種別 'H':HASP標準形式　'E':VCFILE入力形式 */
	Perio      rune // 周期定常計算の時'y'
	Fwdata     *os.File
	Fwdata2    *os.File
	Ftsupw     []byte
	Daystartx  int
	Daystart   int
	Dayend     int
	Daywk      []int
	Dayprn     []int
	Dayntime   int
	Ntimehrprt int
	Ntimedyprt int
	Nhelmsfpri int // 要素別壁体表面温度出力壁体数
	Nvcfile    int // 境界条件、負荷入力用ファイル
	Vcfile     []VCFILE
	Loc        *LOCAT
	Wdpt       WDPT
	DTm        int
	Sttmm      int
	MaxIterate int // 最大収束回数
}

type FLOUT struct {
	F   *os.File
	Idn PrintType
}

type VCFILE struct {
	Fi    *os.File
	Ad    int64
	Ic    int
	Name  string
	Fname string
	Estl  ESTL
	Tlist []TLIST
}

type DAYTM struct {
	day   int // 通日 (day)
	Year  int // 年
	Mon   int // 月
	Day   int // 日
	Ddpri int // 日積算値出力
	Time  float64
	Ttmm  int
	Tt    int
}

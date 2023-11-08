package eeslism

/*****  SIMCONTL の初期化  *****/
func NewSIMCONTL() *SIMCONTL {
	S := new(SIMCONTL)
	S.Title = ""
	S.File = ""
	S.Wfname = ""
	S.Ofname = ""
	S.Unit = ""
	S.Unitdy = ""
	S.Fwdata = nil
	S.Fwdata2 = nil
	S.Ftsupw = nil
	S.Timeid = []rune{' ', ' ', ' ', ' ', ' '}
	S.Helmkey = ' '
	S.Wdtype = ' '
	S.Daystartx = 0
	S.Daystart = 0
	S.Dayend = 0
	S.Dayntime = 0
	S.Ntimedyprt = 0
	S.Ntimehrprt = 0
	S.Nhelmsfpri = 0
	S.Nvcfile = 0
	S.DTm = 0
	S.Sttmm = 0
	S.Vcfile = nil
	S.Loc = nil
	S.Wdpt.Ta = nil
	S.Wdpt.Xa = nil
	S.Wdpt.Rh = nil
	S.Wdpt.Idn = nil
	S.Wdpt.Isky = nil
	S.Wdpt.Ihor = nil
	S.Wdpt.Cc = nil
	S.Wdpt.Rn = nil
	S.Wdpt.Wv = nil
	S.Wdpt.Wdre = nil
	S.MaxIterate = 5 // 最大収束回数のデフォルト値
	S.Daywk = make([]int, 366)
	S.Dayprn = make([]int, 366)

	for i := 0; i < 366; i++ {
		S.Daywk[i] = 0
		S.Dayprn[i] = 0
	}

	return S
}

/*****  COMPNT の初期化  *****/
func NewCOMPNT() *COMPNT {
	C := new(COMPNT)
	C.Name = ""
	C.Roomname = ""
	C.Eqptype = ""
	C.Envname = ""
	C.Exsname = ""
	C.Hccname = ""
	C.Idi = nil
	C.Ido = nil
	C.Tparm = ""
	C.Wetparm = ""
	C.Eqp = nil
	C.Ivparm = nil
	C.Elouts = nil
	C.Elins = nil
	C.Neqp, C.Ncat, C.Nivar = 0, 0, 0
	C.Eqpeff = 0.0
	C.Airpathcpy = true
	C.Control = ' '
	C.Nout, C.Nin = 3, 3
	C.Valvcmp = nil
	C.Rdpnlname = ""
	C.Omparm = ""
	C.PVcap = -999.0
	C.Ac, C.Area = -999.0, -999.0
	// C.x = 1.0
	// C.xinit = -999.0
	// C.org = 'n'
	// C.OMfanName = nil
	C.MonPlistName = ""
	C.MPCM = -999.0
	return C
}

/*****  ELOUT の初期化  *****/
func NewEloutSlice(n int) []*ELOUT {
	s := make([]*ELOUT, n)
	for i := 0; i < n; i++ {
		s[i] = NewElout()
	}
	return s
}

func NewElout() *ELOUT {
	Eo := new(ELOUT)
	Eo.Ni = 0
	Eo.G = 0.0
	Eo.Co = 0.0
	Eo.Coeffo = 0.0
	Eo.Control = ' '
	Eo.Id = ' '
	Eo.Fluid, Eo.Sysld = ' ', ' '
	Eo.Q, Eo.Sysv, Eo.Load = 0.0, 0.0, 0.0
	Eo.Sv, Eo.Sld = 0, 0

	Eo.Cmp = nil
	Eo.Elins = nil
	Eo.Eldobj, Eo.Emonitr = nil, nil

	Eo.Coeffin = nil
	Eo.Lpath = nil
	Eo.Pelmoid = 'x'

	return Eo
}

/*****  ELIN の初期化  *****/
func NewElinSlice(n int) []*ELIN {
	s := make([]*ELIN, n)
	for i := 0; i < n; i++ {
		s[i] = NewElin()
	}
	return s
}

func NewElin() *ELIN {
	Ei := new(ELIN)
	Ei.Id = ' '
	Ei.Sysvin = 0.0
	Ei.Upo, Ei.Upv = nil, nil
	Ei.Lpath = nil
	return Ei
}

func NewPLIST() *PLIST {
	Pl := new(PLIST)
	Pl.Name = ""
	Pl.Type, Pl.Control = ' ', ' '
	Pl.Batch = false
	Pl.Org = true
	Pl.Plmvb, Pl.Pelm = nil, nil
	Pl.Lpair = nil
	Pl.Go = nil
	Pl.G = -999.0
	Pl.Lvc, Pl.Nvav, Pl.Nvalv = 0, 0, 0
	// Pl.Npump = 0
	Pl.N = -999
	Pl.Valv = nil
	Pl.Mpath = nil
	Pl.Plistt, Pl.Plistx = nil, nil
	Pl.Rate = nil
	Pl.Upplist, Pl.Dnplist = nil, nil
	Pl.NOMVAV = 0
	// Pl.Pump = nil
	Pl.OMvav = nil
	Pl.UnknownFlow = 1
	Pl.Plistname = ""
	Pl.Gcalc = 0.0
	return Pl
}

/*****  PELM の初期化  ******/
func NewPELM() *PELM {
	return &PELM{
		Co:  ELIO_SPACE,
		Ci:  ELIO_SPACE,
		Cmp: nil,
		Out: nil,
		In:  nil,
		//Pelmx: nil,
	}
}

/*****  MPATH の初期化  *****/
func NewMPATH() *MPATH {
	return &MPATH{
		Name:    "",
		NGv:     0,
		NGv2:    0,
		Ncv:     0,
		Lvcmx:   0,
		Plist:   nil,
		Mpair:   nil,
		Sys:     ' ',
		Type:    ' ',
		Fluid:   ' ',
		Control: ' ',
		Pl:      nil,
		Cbcmp:   nil,
		G0:      nil,
		Rate:    false,
	}
}

/*****  EXSF の初期化  ******/
func Exsfinit(e *EXSF) {
	e.Name = ""
	e.Typ = 'S'
	e.Wa, e.Wb = 0.0, 0.0
	e.Rg, e.Fs, e.Wz, e.Ww, e.Ws = 0.0, 0.0, 0.0, 0.0, 0.0
	e.Swb, e.CbSa, e.CbCa, e.Cwa = 0.0, 0.0, 0.0, 0.0
	e.Swa, e.Z, e.Erdff, e.Cinc = 0.0, 0.0, 0.0, 0.0
	e.Tazm, e.Tprof, e.Idre, e.Idf = 0.0, 0.0, 0.0, 0.0
	e.Iw, e.Rn, e.Tearth = 0.0, 0.0, 0.0
	e.Erdff = 0.36e-6
	e.Alo = CreateConstantValuePointer(0.0)
	// e.alosch = nil
	e.Alotype = Alotype_Fix
}

/*****  SYSEQ の初期化  *****/
func Syseqinit(S *SYSEQ) {
	S.A = ' '
}

/*****  EQSYS の初期化  *****/
func NewEQSYS() *EQSYS {
	E := new(EQSYS)
	E.Cnvrg = make([]*COMPNT, 0)
	E.Hcc = make([]*HCC, 0)
	E.Boi = make([]*BOI, 0)
	E.Refa = make([]*REFA, 0)
	E.Coll = make([]*COLL, 0)
	E.Pipe = make([]*PIPE, 0)
	E.Stank = make([]*STANK, 0)
	E.Hex = make([]*HEX, 0)
	E.Pump = make([]*PUMP, 0)
	E.Flin = make([]*FLIN, 0)
	E.Hcload = make([]*HCLOAD, 0)
	E.Vav = make([]*VAV, 0)
	E.Stheat = make([]*STHEAT, 0)
	E.Thex = make([]*THEX, 0)
	E.Valv = make([]*VALV, 0)
	E.Qmeas = make([]*QMEAS, 0)
	E.PVcmp = make([]*PV, 0)
	E.OMvav = make([]*OMVAV, 0)

	// 使用されていなかった:
	//E.Ngload = 0
	//E.Gload = nil

	return E
}

/*****  RMVLS の初期化  *****/
func NewRMVLS() *RMVLS {
	R := new(RMVLS)
	R.Twallinit = 0.0
	R.Emrk = nil
	R.Wall = nil
	R.Window = nil
	R.Snbk = nil
	R.Rdpnl = nil
	R.Qrm, R.Qrmd = nil, nil
	R.Trdav = nil
	R.Sd = nil
	R.PCM = nil
	// R.airflow = nil
	R.Pcmiterate = 'n'

	R.Mw = nil
	R.Room = nil
	R.Qetotal.Name = ""

	return R
}

func VPTRinit(v *VPTR) {
	v.Type = ' '
	v.Ptr = nil
}

func TMDTinit(t *TMDT) {
	for i := 0; i < 5; i++ {
		t.Dat[i] = nil
	}

	t.CYear = ""
	t.CMon = ""
	t.CDay = ""
	t.CWkday = ""
	t.CTime = ""
	t.Year, t.Mon, t.Day, t.Time = 0, 0, 0, 0
}

func NewLOCAT() *LOCAT {
	L := new(LOCAT)
	L.Name = ""
	L.Lat, L.Lon, L.Ls, L.Tgrav, L.DTgr = -999.0, -999.0, -999.0, -999.0, -999.0
	L.Daymxert = -999

	matinitx(L.Twsup[:], 12, -999.0)

	return L
}

func NewEQCAT() *EQCAT {
	E := new(EQCAT)
	E.Rfcmp = nil
	E.Hccca = nil
	E.Boica = nil
	E.Refaca = nil
	E.Collca = nil
	E.Pipeca = nil
	E.Stankca = nil
	E.Hexca = nil
	E.Pumpca = nil
	E.Vavca = nil
	E.Stheatca = nil
	E.Thexca = nil
	E.Pfcmp = nil
	E.PVca = nil
	E.OMvavca = nil
	return E
}

func MtEdayinit(mtEday *[12][24]EDAY) {
	for i := 0; i < 12; i++ {
		for j := 0; j < 24; j++ {
			mtEday[i][j].D = 0.0
			mtEday[i][j].Hrs = 0
			mtEday[i][j].Mx = 0.0
			mtEday[i][j].Mxtime = 0
		}
	}
}

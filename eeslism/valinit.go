package eeslism

/*****  SIMCONTL の初期化  *****/
func Simcinit(S *SIMCONTL) {
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
}

/*****  COMPNT の初期化  *****/

func Compinit(N int, Clist []COMPNT) {
	for i := 0; i < N; i++ {
		C := &Clist[i]
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
		C.Airpathcpy, C.Control = ' ', ' '
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
	}
}

/*****  ELOUT の初期化  *****/
func Eloutinit(EoList []*ELOUT, N int) {
	for i := 0; i < N; i++ {
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

		EoList[i] = Eo
	}
}

/*****  ELIN の初期化  *****/
func Elininit(N int, EiList []*ELIN) {
	for i := 0; i < N; i++ {
		Ei := new(ELIN)
		Ei.Id = ' '
		Ei.Sysvin = 0.0
		Ei.Upo, Ei.Upv = nil, nil
		Ei.Lpath = nil
		EiList[i] = Ei
	}
}

/*****  PLIST の初期化  *****/
func Plistinit(N int, PlList []PLIST) {
	for i := 0; i < N; i++ {
		Pl := &PlList[i]
		Pl.Name = ""
		Pl.Type, Pl.Control = ' ', ' '
		Pl.Batch = 'n'
		Pl.Org = 'y'
		Pl.Plmvb, Pl.Pelm = nil, nil
		Pl.Lpair = nil
		Pl.Go = nil
		Pl.G = -999.0
		Pl.Nelm, Pl.Lvc, Pl.Nvav, Pl.Nvalv = 0, 0, 0, 0
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
	}
}

/*****  PELM の初期化  ******/
func Pelminit(N int, PeList []PELM) {
	for i := 0; i < N; i++ {
		Pe := &PeList[i]
		Pe.Co, Pe.Ci = ' ', ' '
		Pe.Cmp = nil
		Pe.Out = nil
		Pe.In = nil
		// Pe.Pelmx = nil
	}
}

/*****  MPATH の初期化  *****/
func Mpathinit(N int, MList []MPATH) {
	for i := 0; i < N; i++ {
		M := &MList[i]
		M.Name = ""
		M.Nlpath, M.NGv, M.NGv2, M.Ncv, M.Lvcmx = 0, 0, 0, 0, 0
		M.Plist = nil
		M.Mpair = nil
		M.Sys, M.Type, M.Fluid, M.Control = ' ', ' ', ' ', ' '
		M.Pl = nil
		M.Cbcmp = nil
		M.G0 = nil
		M.Rate = 'N'
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
	e.Alo = new(float64)
	*e.Alo = 0.0
	// e.alosch = nil
	e.Alotype = Alotype_Fix
}

/*****  SYSEQ の初期化  *****/
func Syseqinit(S *SYSEQ) {
	S.A = ' '
}

/*****  EQSYS の初期化  *****/
func Eqsysinit(E *EQSYS) {
	E.Ncnvrg, E.Nhcc, E.Nboi, E.Nrefa, E.Ncoll = 0, 0, 0, 0, 0
	E.Npipe, E.Nstank, E.Nhex, E.Npump, E.Nflin = 0, 0, 0, 0, 0
	E.Nhcload, E.Ngload, E.Nvav, E.Nstheat, E.Ndesi, E.Nevac = 0, 0, 0, 0, 0, 0
	E.Nthex, E.Nvalv, E.Nqmeas = 0, 0, 0
	E.Npv = 0
	E.Nomvav = 0

	E.Cnvrg = nil
	E.Hcc = nil
	E.Boi = nil
	E.Refa = nil
	E.Coll = nil
	E.Pipe = nil
	E.Stank = nil
	E.Hex = nil
	E.Pump = nil
	E.Flin = nil
	E.Hcload = nil
	E.Gload = nil
	E.Vav = nil
	E.Stheat = nil
	E.Thex = nil
	E.Valv = nil
	E.Qmeas = nil
	E.PVcmp = nil
	E.OMvav = nil
}

/*****  RMVLS の初期化  *****/
func Rmvlsinit(R *RMVLS) {
	R.Twallinit = 0.0
	R.Nwindow, R.Nmwall, R.Nsrf = 0, 0, 0
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
	R.Npcm = 0

	R.Mw = nil
	R.Room = nil
	R.Qetotal.Name = ""
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

func Locinit(L *LOCAT) {
	L.Name = ""
	L.Lat, L.Lon, L.Ls, L.Tgrav, L.DTgr = -999.0, -999.0, -999.0, -999.0, -999.0
	L.Daymxert = -999

	matinitx(L.Twsup[:], 12, -999.0)
}

func Eqcatinit(E *EQCAT) {
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

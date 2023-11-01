// Package eeslism は C言語による Open EESLISM をGo言語に移植したものです。
package eeslism

import (
	"fmt"
	"path/filepath"
	"strings"

	"os"
)

func Entry(InFile string) {
	var s string

	var Daytm DAYTM
	var Simc SIMCONTL
	var Loc LOCAT
	var Wd, Wdd, Wdm WDAT
	var dminute int
	var Rmvls RMVLS
	var Eqcat EQCAT
	var Eqsys EQSYS
	var i int

	/* ============================ */

	var Nelout, Nelin int
	var Compnt []*COMPNT
	var Elout []*ELOUT
	var Elin []*ELIN
	var Syseq SYSEQ
	var Nmpath, Npelm, Nplist int
	var Mpath []*MPATH
	var Plist []*PLIST
	var Pelm []*PELM
	var Ncontl, Nctlif, Nctlst int
	var Contl []*CONTL
	var Ctlif []*CTLIF
	var Ctlst []*CTLST
	var key int
	var Exsf EXSFS
	var Soldy, Solmon []float64

	var uop, ulp []bekt
	var ullp, ulmp *bekt

	/*---------------higuchi add-------------------start*/

	var bdpn, obsn, lpn, opn, mpn, monten, polyn, treen, shadn int

	var DE, co float64

	var wap []float64
	var wip [][]float64

	//var uop, ulp, ullp, ulmp *bekt

	var gp [][]XYZ
	var gpn int

	var fp1, fp2, fp3, fp4 *os.File

	var BDP []BBDP
	var obs []OBS
	var tree []TREE     /*-樹木データ-*/
	var poly []POLYGN   /*--POLYGON--*/
	var shadtb []SHADTB /*-LP面の日射遮蔽率スケジュール-*/
	var op, lp, mp []P_MENN
	var Noplpmp NOPLPMP // OP、LP、MPの定義数

	var Datintvl int
	var dcnt int
	//var Sdstr []SHADSTR

	//var st *C.char
	//	var Ipath string

	BDP = nil
	obs = nil
	tree = nil
	poly = nil
	shadtb = nil
	op = nil
	lp = nil
	mp = nil

	// ここまで修正　satoh 2008/11/8

	/*---------------higuchi add--------------------end*/

	Contl = nil
	Ctlif = nil
	Ctlst = nil

	Mpath = nil
	Plist = nil
	Pelm = nil
	Elout = nil
	Elin = nil
	//Ferr = nil

	Wd.EarthSurface = nil
	Exsf.EarthSrfFlg = false
	Exsf.Exs = nil
	Soldy = nil
	Loc.Name = ""
	Fbmlist = ""

	Rmvlsinit(&Rmvls)
	Simcinit(&Simc)

	Eqsysinit(&Eqsys)
	Locinit(&Loc)
	Eqcatinit(&Eqcat)

	/* ------------------------------------------------------ */

	Psyint()

	Ifile := InFile
	s = Ifile

	// 入力されたパスが"で始まる場合に除去する
	if strings.HasPrefix(Ifile, "\"") {
		fmt.Sscanf(Ifile, "\"%~[\"]\"", s)
	}

	// ディレクトリ名のみ
	Ifile = filepath.Dir(Ifile)

	// 注釈文の除去
	bdata0 := Eesprera(s)

	// スケジュ－ルデ－タの作成
	EWKFile := strings.TrimSuffix(s, filepath.Ext(s))
	bdata, schtba, schnma, week := Eespre(bdata0, EWKFile, &key) //key=`WEEK`が含まれているかどうか

	Simc.File = InFile
	Simc.Loc = &Loc

	Schdl, Flout := Eeinput(
		EWKFile,
		bdata, week, schtba, schnma,
		&Simc, &Exsf, &Rmvls, &Eqcat, &Eqsys,
		&Compnt,
		&Elout, &Nelout,
		&Elin, &Nelin,
		&Mpath, &Nmpath,
		&Plist,
		&Pelm, &Npelm,
		&Contl, &Ncontl,
		&Ctlif, &Nctlif,
		&Ctlst, &Nctlst,
		&Wd,
		&Daytm, key, &Nplist,
		&bdpn, &obsn, &treen, &shadn, &polyn, &BDP, &obs, &tree, &shadtb, &poly, &monten, &gpn, &DE, &Noplpmp)

	// 外部障害物のメモリを確保
	op = make([]P_MENN, Noplpmp.Nop)
	lp = make([]P_MENN, Noplpmp.Nlp)
	P_MENNinit(op, Noplpmp.Nop)
	P_MENNinit(lp, Noplpmp.Nlp)

	// 最大収束回数のセット
	LOOP_MAX := Simc.MaxIterate
	VAV_Count_MAX := Simc.MaxIterate

	// 動的カーテンの展開
	for i := 0; i < Rmvls.Nsrf; i++ {
		Sd := &Rmvls.Sd[i]
		if Sd.DynamicCode != "" {
			ctifdecode(Sd.DynamicCode, Sd.Ctlif, &Simc, Compnt, Mpath, &Wd, &Exsf, Schdl)
		}
	}

	if bdpn != 0 {

		RET := STRCUT(s, ".")
		RET1 := RET
		RET3 := RET
		RET14 := RET
		RET15 := RET

		RET += "_shadow.gchi"
		RET1 += "_I.gchi"
		RET3 += "_ffactor.gchi"
		RET14 += "_lwr.gchi"

		var err error
		if fp1, err = os.Create(RET); err != nil {
			fmt.Println("File not open _shadow.gchi")
			os.Exit(1)
		}

		if fp2, err = os.Create(RET1); err != nil {
			fmt.Println("File not open _I.gchi")
			os.Exit(1)
		}

		if fp3, err = os.Create(RET14); err != nil {
			fmt.Println("File not open _lwr.gchi")
			os.Exit(1)
		}

		if fp4, err = os.Create(RET3); err != nil {
			fmt.Println("File not open _ffactor.gchi")
			os.Exit(1)
		}

		// 座標の変換
		LP_COORDNT(&lpn, bdpn, obsn, treen, polyn, poly, tree, obs, BDP, lp)

		OP_COORDNT(&opn, bdpn, BDP, op, polyn, poly)

		// LPの構造体に日毎の日射遮蔽率を代入
		for i := 0; i < lpn; i++ {
			for j := 0; j < shadn; j++ {
				if lp[i].opname == shadtb[j].lpname {
					for k := 1; k < 366; k++ {
						for l := 0; l < shadtb[j].indatn; l++ {
							if k >= shadtb[j].ndays[l] && k <= shadtb[j].ndaye[l] {
								lp[i].shad[k] = shadtb[j].shad[l]
								break
							}
						}
					}
				}
			}
		}

		//---- mpの総数をカウント mpは、OP面+OPW面 ---------------
		mpn := 0
		for i := 0; i < opn; i++ {
			mpn += 1
			for j := 0; j < op[i].wd; j++ {
				mpn += 1
			}
		}

		//---窓壁のカウンター変数の初期化---
		//wap := make([]float64, opn)
		wip := make([][]float64, opn)
		for i := 0; i < opn; i++ {
			if op[i].wd != 0 {
				wip[i] = make([]float64, op[i].wd)
			}
		}

		//---領域の確保   gp 地面の座標(X,Y,Z)---
		gp := make([][]XYZ, mpn)
		for i := 0; i < mpn; i++ {
			gp[i] = make([]XYZ, gpn+1)
		}

		//---領域の確保 mp---
		mp := make([]P_MENN, Noplpmp.Nmp)
		P_MENNinit(mp, mpn)

		//----OP,OPWの構造体をMPへ代入する----
		DAINYUU_MP(&mp, op, opn, mpn)

		for i := 0; i < mpn; i++ {
			fmt.Fprintf(fp1, "%s\n", mp[i].opname)
		}

		//---ベクトルの向きを判別する変数の初期化---
		//---opから見たopの位置---
		uop := make([]bekt, opn)
		for i := 0; i < opn; i++ {
			uop[i].ps = make([][]float64, opn)
			for j := 0; j < opn; j++ {
				uop[i].ps[j] = make([]float64, op[j].polyd)
			}
		}

		//---opから見たlpの位置---
		ulp := make([]bekt, opn)
		for i := 0; i < opn; i++ {
			ulp[i].ps = make([][]float64, lpn)
			for j := 0; j < lpn; j++ {
				ulp[i].ps[j] = make([]float64, lp[j].polyd)
			}
		}

		//---lpから見たlpの位置---
		ullp := make([]bekt, lpn)
		for i := 0; i < lpn; i++ {
			ullp[i].ps = make([][]float64, lpn)
			for j := 0; j < lpn; j++ {
				ullp[i].ps[j] = make([]float64, lp[j].polyd)
			}
		}

		//---lpから見たmpの位置---
		ulmp := make([]bekt, lpn)
		for i := 0; i < lpn; i++ {
			ulmp[i].ps = make([][]float64, mpn)
			for j := 0; j < mpn; j++ {
				ulmp[i].ps[j] = make([]float64, mp[j].polyd)
			}
		}

		//------CG確認用データ作成-------
		HOUSING_PLACE(lpn, mpn, lp, mp, RET15)

		//----前面地面代表点および壁面の中心点を求める--------
		GRGPOINT(mp, mpn)
		for i := 0; i < lpn; i++ {
			GDATA(&lp[i], &lp[i].G)
		}

		// 20170426 higuchi add 条件追加　形態係数を計算しないパターンを組み込んだ
		if monten > 0 {
			//---LPから見た天空に対する形態係数算出------
			FFACTOR_LP(lpn, mpn, monten, lp, mp)
		}

		for i := 0; i < mpn; i++ {
			for j := range Rmvls.Sd {
				if Rmvls.Sd[j].Sname == mp[i].opname {
					mp[i].exs = Rmvls.Sd[j].exs
					mp[i].as = Rmvls.Sd[j].as
					mp[i].alo = Rmvls.Sd[j].alo
					mp[i].Eo = Rmvls.Sd[j].Eo
					break
				}
			}
		}

		for i := 0; i < mpn; i++ {
			mp[i].refg = Exsf.Exs[mp[i].exs].Rg
			//fmt.Printf("mp[%d].refg=%f\n", i, mp[i].refg)
		}

		//---面の裏か表かの判断をするためのベクトル値の算出--
		URA(opn, opn, op, uop, op)  //--opから見たopの位置--
		URA(lpn, lpn, lp, ullp, lp) //--lpから見たlpの位置--
		URA(lpn, mpn, mp, ulmp, lp) //--lpから見たmpの位置--
		URA(opn, lpn, lp, ulp, op)  //--opから見たlpの位置--

		// if test {
		// 	/*---op,lp座標の確認-------*/
		// 	ZPRINT(lp, op, lpn, opn, RET13)
		// 	ZPRINT(mp, op, mpn, opn, RET6)
		// 	mp_printf(mpn, mp, RET7)
		// 	lp_printf(lpn, lp, RET8)
		// 	e_printf(lpn, lp, RET9)
		// 	e_printf(mpn, mp, RET10)
		// 	errbekt_printf(lpn, lpn, ullp, RET11)
		// 	errbekt_printf(lpn, mpn, ulmp, RET12)
		// }

		fmt.Fprintf(fp2, "M\nD\nmt\nname\ngl_shadow\nIsky\nIg\nIb\nIdf\nIdre\n")
		fmt.Fprintf(fp3, "M\nD\nmt\nname\nRsky\nreff\nreffg\nReff\n")

	}

	/*-----------------higuchi add-----------------------------end*/

	if DEBUG {
		fmt.Println("eeinput end")

		for i, Pe := range Pelm {
			fmt.Printf("[%3d] Pelm=%s\n", i, Pe.Cmp.Name)
		}

		for i, Eo := range Elout {
			fmt.Printf("[%3d] Eo_cmp=%s\n", i, Eo.Cmp.Name)
		}

		fmt.Printf("Npelm=%d Ncompnt=%d Nelout=%d Nelin=%d\n",
			Npelm, len(Compnt), Nelout, Nelin)
	}

	Soldy = make([]float64, len(Exsf.Exs))
	Solmon = make([]float64, len(Exsf.Exs))

	DTM = float64(Simc.DTm)
	dminute = int(float64(Simc.DTm) / 60.0)
	Cff_kWh = DTM / 3600.0 / 1000.0

	for rm := range Rmvls.Room {
		Rm := &Rmvls.Room[rm]
		Rm.Qeqp = 0.0
	}

	dprschtable(Schdl.Seasn, Schdl.Wkdy, Schdl.Dsch, Schdl.Dscw)
	//dprschdata ( Schdl.Sch, Schdl.Scw ) ;
	//dprachv ( Rmvls.Nroom, Rmvls.Room ) ;
	dprexsf(Exsf.Exs)
	dprwwdata(Rmvls.Wall, Rmvls.Window)
	dprroomdata(Rmvls.Room, Rmvls.Sd)
	dprballoc(Rmvls.Mw, Rmvls.Sd)

	Simc.eeflopen(Flout)

	if DEBUG {
		fmt.Println("<<main>> eeflopen ")
	}

	if Ferr != nil {
		fmt.Fprintln(Ferr, "\n<<main>> eeflopen end")
	}

	Tinit(Rmvls.Twallinit, Rmvls.Room,
		Rmvls.Nsrf, Rmvls.Sd, Rmvls.Nmwall, Rmvls.Mw)

	if DEBUG {
		fmt.Println("<<main>> Tinit")
	}

	if Ferr != nil {
		fmt.Fprintln(Ferr, "\n<<main>> Tinit")
	}

	// ボイラ機器仕様の初期化
	Boicaint(Eqcat.Boica, &Simc, Compnt, &Wd, &Exsf, Schdl)
	Mecsinit(&Eqsys, &Simc, Compnt, Exsf.Exs, &Wd, &Rmvls)

	if DEBUG {
		fmt.Println("<<main>> Mecsinit")
	}

	if Ferr != nil {
		fmt.Fprintln(Ferr, "\n<<main>> Mecsinit")
	}

	/*******************

	1997.11.18	熱損失係数計算用改訂

	*******************/

	bdhpri(Simc.Ofname, Rmvls, &Exsf)

	// xprtwallinit (Rmvls.Nmwall, Rmvls.Mw);

	/* --------------------------------------------------------- */

	Daytm.Ddpri = 0

	if Simc.Sttmm < 0 {
		Simc.Sttmm = dminute
	}

	tt := Simc.Sttmm / 100
	mm := Simc.Sttmm % 100
	mta := (tt*60 + mm) / dminute
	mtb := (24 * 60) / dminute
	//mtb = 12;
	// 110413 higuchi add  影面積をストアして、影計算を10日おきにする
	Sdstr := make([]SHADSTR, mpn)
	for i := 0; i < mpn; i++ {
		Sdstr[i].sdsum = make([]float64, mtb)
		for jj := 0; jj < mtb; jj++ {
			Sdstr[i].sdsum[jj] = 0.0
		}
	}

	dcnt = 0

	for nday := Simc.Daystartx; nday <= Simc.Dayend; nday++ {
		if dcnt == Datintvl {
			dcnt = 0
			MATINIT_sdstr(mpn, mtb, Sdstr)
		}
		dcnt++

		if dayprn && Ferr != nil {
			fmt.Fprintf(Ferr, "\n\n\t===== Dayly Loop =====\n\n")
		}

		day := ((nday - 1) % 365) + 1
		if Simc.Perio == 'y' {
			day = Simc.Daystart
		}
		Daytm.DayOfYear = day

		dayprn = Simc.Dayprn[day] != 0

		if Simc.Perio != 'y' && nday > Simc.Daystartx {
			Daytm.Mon, Daytm.Day = monthday(Daytm.Mon, Daytm.Day)
		}

		if nday >= Simc.Daystart {
			Daytm.Ddpri = 1
		}
		if Simc.Perio == 'y' && nday != Simc.Dayend {
			Daytm.Ddpri = 0
		}

		if nday > Simc.Daystartx {
			mta = 1
		}

		for mt := mta; mt <= mtb; mt++ {
			if dayprn && Ferr != nil {
				fmt.Fprintf(Ferr, "\n\n\t===== Timely Loop =====\n\n")
			}

			if mm >= 60 {
				mm -= 60
				tt++
			}

			if tt > 24 || (tt == 24 && mm > 0) {
				tt -= 24
			}

			Daytm.Tt = tt
			Daytm.Ttmm = tt*100 + mm
			Daytm.Time = float64(Daytm.Ttmm) / 100.0

			if DEBUG {
				fmt.Printf("<< main >> nday=%d mm=%d mt=%d  tt=%d mm=%d\n", nday, mm, mt, tt, mm)
			}

			//if (day == 16 && Daytm.ttmm == 800)
			//  fmt.Printf("xxxxxx\n")

			Vcfinput(&Daytm, Simc.Nvcfile, Simc.Vcfile, Simc.Perio)
			Weatherdt(&Simc, &Daytm, &Loc, &Wd, Exsf.Exs, Exsf.EarthSrfFlg)

			if dayprn && Ferr != nil {
				fmt.Fprintf(Ferr, "\n\n\n---- date=%2d/%2d nday=%d day=%d time=%5.2f ----\n",
					Daytm.Mon, Daytm.Day, nday, day, Daytm.Time)
			}

			if DEBUG {
				fmt.Printf("---- date=%2d %2d nday=%d day=%d time=%5.2f ----\n",
					Daytm.Mon, Daytm.Day, nday, day, Daytm.Time)
			}

			if dayprn && Ferr != nil {
				Flinprt(Eqsys.Flin)
			}

			/***   if (Daytm.ttmm == 100 )****/
			if mt == mta {
				fmt.Printf("%d/%d", Daytm.Mon, Daytm.Day)
				if nday < Simc.Daystart {
					fmt.Printf(")")
				}
				if Daytm.Ddpri != 0 && Simc.Dayprn[day] != 0 {
					fmt.Printf(" *")
				}
				fmt.Printf("\n")

				/*------------------------higuchi add---形態係数の算出---------start*/
				//fmt.Printf("nday=%d,day=%d\n",nday,day) ;
				//fmt.Printf("bdpn=%d\n",bdpn) ;

				// 20170426 higuchi add 形態係数を計算しない処理の追加
				if bdpn != 0 && monten > 0 {
					if nday == Simc.Daystartx {
						fmt.Printf("form_factor calcuration start\n")
						GR_MONTE_CARLO(mp, mpn, lp, lpn, monten, day)
						MONTE_CARLO(mpn, lpn, monten, mp, lp, gp, gpn, day, Simc.Daystartx)
						ffactor_printf(fp4, mpn, lpn, mp, lp, Daytm.Mon, Daytm.Day)
						fmt.Printf("form_factor calcuration end\n")
					} else {
						for i := 0; i < lpn; i++ {
							k := day - 1
							if k == 0 {
								k = 365
							}
							if lp[i].shad[day] != lp[i].shad[k] {
								fmt.Printf("form_factor calcuration start:shad[%d]=%f,shad[%d]=%f\n", nday, lp[i].shad[day], k, lp[i].shad[k])
								GR_MONTE_CARLO(mp, mpn, lp, lpn, monten, day)
								MONTE_CARLO(mpn, lpn, monten, mp, lp, gp, gpn, day, Simc.Daystartx)
								ffactor_printf(fp4, mpn, lpn, mp, lp, Daytm.Mon, Daytm.Day)
								fmt.Printf("form_factor calcuration end\n")
								break
							}
						}
					}

				}
				/*------------------------higuchi add-----------------------end*/

				if DEBUG {
					fmt.Printf(" ** daymx=%d  Tgrav=%f  DT=%f  Tsupw=%f\n",
						Loc.Daymxert, Loc.Tgrav, Loc.DTgr, Wd.Twsup)
				}
			}

			// 傾斜面日射量の計算
			Exsfsol(&Wd, Exsf.Exs)

			/*==transplantation to eeslism from KAGExSUN by higuchi 070918==start*/
			if bdpn != 0 {

				MATINIT_sum(opn, op)
				MATINIT_sum(mpn, mp)
				MATINIT_sum(lpn, lp)

				ls := -Wd.Sw
				ms := -Wd.Ss
				ns := Wd.Sh

				if Wd.Sh > 0.0 {

					// 110413 higuchi add 下の条件
					if dcnt == 1 {

						for j := 0; j < opn; j++ {

							wap[j] = 0.0
							for i := 0; i < op[j].wd; i++ {
								wip[j][i] = 0.0
							}
							CINC(op[j], ls, ms, ns, &co)
							if co > 0.0 {
								SHADOW(j, DE, opn, lpn, ls, ms, ns, &uop[j], &ulp[j], &op[j], op, lp, &wap[j], wip[j], day)
							} else {
								op[j].sum = 1.0
								for i := 0; i < op[j].wd; i++ {
									op[j].opw[i].sumw = 1.0
								}
							}
						}
						//fmt.Printf("dcnt1=%d\n",dcnt) ;
						DAINYUU_SMO2(opn, mpn, op, mp, Sdstr, dcnt, mt)

						// 20170426 higuchi add 条件追加
						if dayprn {
							shadow_printf(fp1, Daytm.Mon, Daytm.Day, Daytm.Time, mpn, mp)
						}

					} else {
						//fmt.Printf("dcnt2=%d\n",dcnt) ;
						DAINYUU_SMO2(opn, mpn, op, mp, Sdstr, dcnt, mt)
						// 20170426 higuchi add 条件追加
						if dayprn {
							shadow_printf(fp1, Daytm.Mon, Daytm.Day, Daytm.Time, mpn, mp)
						}
					}

					//SHADSTR *Sdstrd;
					//Sdstrd = Sdstr;
					//for (i = 0; i < mpn; i++, Sdstrd++)
					//{
					//  int m;
					//  for (m = 0; m < mtb; m++)
					//      printf("Sdstr[%d].sdsum[%d]=%f\n", i, m, Sdstrd->sdsum[m]);
					//}
				}

				// 20170426 higuchi add 条件追加
				if dayprn {
					fmt.Fprintf(fp2, "%d %d %5.2f\n", Daytm.Mon, Daytm.Day, Daytm.Time)
					fmt.Fprintf(fp3, "%d %d %5.2f\n", Daytm.Mon, Daytm.Day, Daytm.Time)
				}

				// 20170426 higuchi add 引数追加 dayprn,monten
				OPIhor(fp2, fp3, lpn, mpn, mp, lp, &Wd, ullp, ulmp, gp, day, monten)
				for i := range Rmvls.Sd {
					if Rmvls.Sd[i].Sname != "" {
						for j := 0; j < mpn; j++ {
							if Rmvls.Sd[i].Sname == mp[j].opname {
								Rmvls.Sd[i].Fsdw = mp[j].sum
								//fmt.Printf("Sd->Fswd=%f\n", Rmvls.Sd[i].Fsdw)
								Rmvls.Sd[i].Idre = mp[j].Idre
								Rmvls.Sd[i].Idf = mp[j].Idf
								Rmvls.Sd[i].Iw = mp[j].Iw
								Rmvls.Sd[i].rn = mp[j].Reff
								//fmt.Printf("Sd->ali=%f\n", Rmvls.Sd[i].ali)
								break
							}
						}
					}
				}
			}

			/*===============higuchi 070918============================end*/
			if dayprn && Ferr != nil {
				xprsolrd(Exsf.Exs)
			}

			if DEBUG {
				xprsolrd(Exsf.Exs)
				fmt.Println("<<main>> Exsfsol")
			}

			// Create schedule for current time step
			Eeschdlr(day, Daytm.Ttmm, Schdl, &Rmvls)

			if DEBUG {
				fmt.Println("<<main>>  Eeschdlr")
			}

			// Calculate control status values for use in control (e.g. collector equivalent outdoor temperature)
			CalcControlStatus(&Eqsys, &Rmvls, &Wd, &Exsf)

			// Update control information
			Contlschdlr(Ncontl, Contl, Mpath, Compnt)

			// Recalculate internal heat gains after setting air conditioning on/off schedule
			Qischdlr(Rmvls.Room)

			if DEBUG {
				fmt.Println("<<main>> Contlschdlr")
			}

			/***
			eloutprint(0, Nelout, Elout, Compnt);
			*****/

			VAVcountreset(Eqsys.Vav)
			Valvcountreset(Eqsys.Valv)
			Evaccountreset(Eqsys.Evac)

			/*---- Satoh Debug VAV  2000/12/6 ----*/
			/* VAV 計算繰り返しループの開始地点 */
			for j := 0; j < VAV_Count_MAX; j++ {
				if DEBUG {
					fmt.Printf("\n\n====== VAV LOOP Count=%d ======\n\n\n", j)
				}
				if dayprn && Ferr != nil {
					fmt.Fprintf(Ferr, "\n\n====== VAV LOOP Count=%d ======\n\n\n", j)
				}

				VAVreset := 0
				Valvreset := 0

				Pumpflow(Eqsys.Pump)

				if DEBUG {
					fmt.Println("<<main>> Pumpflow")
				}

				if Simc.Dayprn[day] != 0 && Ferr != nil {
					fmt.Fprintln(Ferr, "<<main>> Pumpflow")
				}

				Pflow(Mpath, &Wd)

				if DEBUG {
					fmt.Println("<<main>> Pflow")
				}

				if dayprn && Ferr != nil {
					fmt.Fprintln(Ferr, "<<main>> Pflow")
				}

				/************
				eloutprint(0, Nelout, Elout, Compnt);
				***********/

				Sysupv(Mpath, &Rmvls)

				if DEBUG {
					fmt.Println("<<main>> Sysupv")
				}

				if dayprn && Ferr != nil {
					fmt.Fprintln(Ferr, "<<main>> Sysupv")
				}

				/*****
				elinprint(0, Compnt, Elout, Elin);
				***********/

				for i := range Rmvls.Room {
					Rmvls.Emrk[i] = '!'
				}

				for n := 0; n < Rmvls.Nsrf; n++ {
					Rmvls.Sd[n].mrk = '!'
				}

				Mecscf(&Eqsys)

				if DEBUG {
					fmt.Println("<<main>> Mecscf")
				}

				/*======higuchi update 070918==========*/
				eeroomcf(&Wd, &Exsf, &Rmvls, nday, mt)
				/*=====================================*/

				if DEBUG {
					fmt.Println("<<main>> eeroomcf")
				}

				/*   作用温度制御時の設定室内空気温度  */
				Rmotset(Rmvls.Room)
				if DEBUG {
					fmt.Println("<<main>> Rmotset End")
				}

				/* 室、放射パネルのシステム方程式作成 */
				Roomvar(Rmvls.Room, Rmvls.Rdpnl)

				if DEBUG {
					fmt.Println("<<main>>  Roomvar")
					eloutprint(1, Elout, Compnt)
					elinprint(1, Compnt, Elout, Elin)
				}

				if dayprn && Ferr != nil {
					fmt.Fprintf(Ferr, "<<main>> Roomvar\n")
					eloutfprint(1, Elout, Compnt)
					elinfprint(1, Compnt, Elout, Elin)
				}
				//eloutprint(1, Nelout, Elout, Compnt);

				//hcldmodeinit(&Eqsys);

				for i := 0; i < LOOP_MAX; i++ {
					//s := fmt.Sprintf("Loop Start %d", i)

					if i == 0 {
						hcldwetmdreset(&Eqsys)
					}

					if DEBUG {
						fmt.Printf("再計算が必要な機器のループ %d\n", i)
					}

					if dayprn && Ferr != nil {
						fmt.Fprintf(Ferr, "再計算が必要な機器のループ %d\n\n\n", i)
					}

					LDreset := 0
					DWreset := 0
					TKreset := 0
					BOIreset := 0
					Evacreset := 0
					PCMfunreset := 0

					/********************************
					if ( TKreset > 0 )
					fmt.Printf("<< main >> nday=%d mt=%d  tt=%d mm=%d TKreset=%d\n",
					nday, mt, tt, mm, TKreset );
					****************************/

					Stankcfv(Eqsys.Stank)

					Hcldcfv(Eqsys.Hcload)

					Syseqv(Nelout, Elout, &Syseq)

					Sysvar(Compnt)

					Roomene(&Rmvls, Rmvls.Room, Rmvls.Rdpnl, &Exsf, &Wd)

					Roomload(Rmvls.Room, &LDreset)

					// PCM家具の収束判定
					PCMfunchk(Rmvls.Room, &Wd, &PCMfunreset)

					// 壁体内部温度の計算と収束計算のチェック
					if Rmvls.Pcmiterate == 'y' {
						PCMwlchk(i, &Rmvls, &Exsf, &Wd, &LDreset)
					}

					Boiene(Eqsys.Boi, &BOIreset)

					Refaene(Eqsys.Refa, &LDreset)

					Hcldene(Eqsys.Hcload, &LDreset, &Wd)

					Hccdwreset(Eqsys.Hcc, &DWreset)

					Stanktss(Eqsys.Stank, &TKreset)

					Evacene(Eqsys.Evac, &Evacreset)

					if BOIreset+LDreset+DWreset+TKreset+Evacreset+PCMfunreset == 0 {
						break
					}
				}

				if i == LOOP_MAX {
					fmt.Printf("収束しませんでした。 MAX=%d\n", LOOP_MAX)
				}

				//fmt.Printf("Loop=%d\n", i)
				Hccene(Eqsys.Hcc)
				// 風量の計算は最初だけ
				//if i == 0 {
				VAVene(Eqsys.Vav, &VAVreset)
				//}
				Valvene(Eqsys.Valv, &Valvreset)

				//fmt.Printf("\n\nVAVreset=%d\n", VAVreset)
				/***************/
				if VAVreset == 0 && Valvreset == 0 {
					break
				}
				VAVcountinc(Eqsys.Vav)
				Valvcountinc(Eqsys.Valv)

				// 風量が変わったら電気蓄熱暖房器の係数を再計算
				Stheatcfv(Eqsys.Stheat)

				/*****************/

				/*---- Satoh Debug VAV  2000/12/6 ----*/
				/* VAV 計算繰り返しループの終了地点 */
			}

			// 太陽電池内蔵壁体の発電量計算
			CalcPowerOutput(Rmvls.Nsrf, Rmvls.Sd, &Wd, &Exsf)

			if Simc.Helmkey == 'y' {
				Helmroom(Rmvls.Room, Rmvls.Qrm, &Rmvls.Qetotal, Wd.T, Wd.X)
			}

			/*************
			fmt.Printf("xxxmain Pathheat\n")
			Pathheat(Nmpath, Mpath)
			************************/

			// 室の熱取得要素の計算
			Qrmsim(Rmvls.Room, &Wd, Rmvls.Qrm)

			for rm := range Rmvls.Room {
				Rmvls.Room[rm].Qeqp = 0.0
			}

			if DEBUG {
				fmt.Printf("Mecsene st\n")
			}

			/*  システム使用機器の供給熱量、エネルギーの計算  */
			Mecsene(&Eqsys)

			/***********************
			fmt.Printf("Mecsene en\n")
			/***********************/

			if DEBUG {
				mecsxprint(&Eqsys)
			}

			/* ------------------------------------------------ */
			if DEBUG {
				fmt.Printf("xxxmain 2\n")
			}

			// 前時刻の室温の入れ替え、OT、MRTの計算
			Rmsurft(Rmvls.Room, Rmvls.Sd)

			if DEBUG {
				fmt.Printf("xxxmain 3\n")
			}

			//if (Daytm.Mon == 1 && Daytm.Day == 5 && fabs(Daytm.Time - 23.15) < 1.e-5)
			//	printf("debug\n");

			// 壁体内部温度の計算（ヒステリシス考慮PCMの状態値もここで設定）
			RMwlt(Rmvls.Nmwall, Rmvls.Mw)

			if DEBUG {
				fmt.Printf("xxxmain 4\n")
			}

			// PMV、SET*の計算
			Rmcomfrt(Rmvls.Room)

			if DEBUG {
				fmt.Printf("xxxmain 5\n")
			}

			//xprsolrd (Exsf.Nexs, Exsf.Exs);

			// 代表日の毎時計算結果のファイル出力
			Eeprinth(&Daytm, &Simc, Flout, &Rmvls, &Exsf, Nmpath, Mpath, &Eqsys, &Wd)

			if DEBUG {
				fmt.Printf("xxxmain 6\n")
			}

			if Daytm.Ddpri != 0 {
				// 室の日集計、月集計
				Roomday(Daytm.Mon, Daytm.Day, day, Daytm.Ttmm, Rmvls.Room, Rmvls.Rdpnl, Simc.Dayend)
				if Simc.Helmkey == 'y' {
					Helmdy(day, Rmvls.Room, &Rmvls.Qetotal)
				}

				Compoday(Daytm.Mon, Daytm.Day, day, Daytm.Ttmm, &Eqsys, Simc.Dayend)
				/**   if (Nqrmpri > 0)  **/
				Qrmsum(Daytm.Day, Rmvls.Room, Rmvls.Qrm, Rmvls.Trdav, Rmvls.Qrmd)

				if DEBUG {
					fmt.Printf("xxxmain 7\n")
				}

				// 気象データの日集計、月集計
				Wdtsum(Daytm.Mon, Daytm.Day, day, Daytm.Ttmm, &Wd, Exsf.Exs, &Wdd, &Wdm, Soldy, Solmon, &Simc)
			}
			if DEBUG {
				fmt.Printf("xxxmain 8\n")
			}

			if DEBUG {
				//xprtwpanel (Rmvls.Nroom, Rmvls.Room, Twp, Sd, Mw);
				xprtwsrf(Rmvls.Nsrf, Rmvls.Sd)
				xprrmsrf(Rmvls.Nsrf, Rmvls.Sd)
				xprtwall(Rmvls.Nmwall, Rmvls.Mw)
			}

			mm += dminute

			// 時刻ループの最後
		}

		// 日集計の出力
		Eeprintd(&Daytm, &Simc, Flout, &Rmvls, Exsf.Exs, Soldy, &Eqsys, &Wdd)
		/*****fmt.Printf("xxxmain 9\n")*****/

		//if (Daytm.Mon == 4 && Daytm.Day == 25)
		//	printf("debug\n");

		// 月集計の出力
		if IsEndDay(Daytm.Mon, Daytm.Day, Daytm.DayOfYear, Simc.Dayend) && Daytm.Ddpri != 0 {
			//fmt.Printf("月集計出力\n")
			Eeprintm(&Daytm, &Simc, Flout, &Rmvls, Exsf.Exs, Solmon, &Eqsys, &Wdm)
		}

	}
	// 月－時刻別集計値の出力
	Eeprintmt(&Simc, Flout, &Eqsys, Rmvls.Rdpnl)

	if DEBUG {
		fmt.Printf("メモリ領域の解放\n")
	}

	Eeflclose(Flout)

	/*------------------higuchi add---------------------start*/
	if bdpn != 0 {

		defer fp1.Close()
		defer fp2.Close()
		defer fp3.Close()
		defer fp4.Close()
	}

	/*---------------------higuchi 1999.7.21-----------end*/

	// return (0);
	// In Go, you don't need to return an integer value like in C
	// since the compiler will automatically assume a zero exit status
}

func matiniti(A []int, N int) {
	for i := 0; i < N; i++ {
		A[i] = 0
	}
}

func P_MENNinit(_pm []P_MENN, N int) {
	const pmax = 200
	for i := 0; i < N; i++ {
		pm := &_pm[i]
		pm.Nopw = 0
		pm.opname = ""
		matinit(pm.rgb[:], 3)
		matinit(pm.faiwall[:], pmax)
		pm.wd = 0
		pm.exs = 0
		pm.grpx, pm.faia, pm.faig, pm.grpfaia, pm.sum, pm.ref, pm.refg, pm.wa, pm.wb = 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0
		pm.Ihor, pm.Idre, pm.Idf, pm.Iw, pm.Reff, pm.rn, pm.Te, pm.Teg = 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0
		matinit(pm.shad[:], 365)
		pm.alo, pm.as, pm.Eo = 0.0, 0.0, 0.0
		pm.polyd, pm.sbflg = 0, 0
		pm.P, pm.opw = nil, nil
	}
}

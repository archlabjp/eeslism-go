// Package eeslism は C言語による Open EESLISM をGo言語に移植したものです。
package eeslism

import (
	"fmt"
	"path/filepath"
	"strings"

	"os"
)

/*
Entry (Main Entry Point for Building Energy Simulation)

この関数は、建物のエネルギーシミュレーションのメインエントリポイントであり、
入力データの読み込み、モデルの初期化、時間ステップごとの計算ループ、
および結果の出力といった一連のプロセスを統括します。

建築環境工学的な観点:
- **シミュレーションの全体フロー**: この関数は、
  建物のエネルギーシミュレーションの全体的な流れを定義します。
  - **入力データの処理**: 建物形状、材料特性、設備機器、スケジュール、気象データなど、
    シミュレーションに必要なあらゆる入力データを読み込み、解析します。
  - **モデルの初期化**: 読み込んだデータに基づいて、
    室、壁体、窓、設備機器などの熱的モデルを初期化します。
  - **時間ステップ計算ループ**: シミュレーション期間（年間など）を小さな時間ステップに分割し、
    各時間ステップで以下の計算を繰り返します。
    - 太陽位置と日射量の計算。
    - 外部日射面と壁体内部温度の計算。
    - 室の熱収支計算と熱負荷の算出。
    - 空調システムや熱源設備の運転シミュレーション。
    - 蓄熱槽やPCMの熱的挙動の計算。
    - 室内温湿度環境の評価。
  - **結果の出力**: 各時間ステップや日、月ごとの計算結果をファイルに出力します。
- **日影計算と形態係数**: `LP_COORDNT`, `OP_COORDNT`, `GR_MONTE_CARLO`, `MONTE_CARLO`などの関数を呼び出し、
  建物や周囲の障害物による日影の形状と面積、
  および表面間の形態係数を計算します。
  これは、日射熱取得量や放射熱伝達を正確にモデル化するために不可欠です。
- **収束計算**: `LOOP_MAX`や`VAV_Count_MAX`といったループは、
  非線形な熱的挙動（PCMの相変化、VAVシステムの流量制御など）を持つシステムにおいて、
  計算の収束を確保するために用いられます。
- **デバッグと検証**: `DEBUG`フラグや`Ferr`ファイルへの出力は、
  シミュレーションの途中の計算結果を確認し、
  モデルの妥当性を検証したり、問題の特定を行ったりするために用いられます。

この関数は、建物のエネルギー性能を総合的に評価し、
省エネルギー設計、快適性評価、
および最適な設備システム設計を行うための中心的な役割を果たします。
*/
func Entry(InFile string, efl_path string) {
	var s string

	var Daytm DAYTM
	var Simc *SIMCONTL
	var Loc *LOCAT
	var Wd, Wdd, Wdm WDAT
	var dminute int
	var Rmvls *RMVLS
	var Eqcat *EQCAT
	var Eqsys *EQSYS
	var i int

	/* ============================ */

	var Compnt []*COMPNT
	var Elout []*ELOUT
	var Elin []*ELIN
	var Syseq SYSEQ
	var Mpath []*MPATH
	var Plist []*PLIST
	var Pelm []*PELM
	var Contl []*CONTL
	var Ctlif []*CTLIF
	var Ctlst []*CTLST
	var key int
	var Exsf EXSFS       // 外皮
	var Soldy []float64  // 日集計データ
	var Solmon []float64 // 月集計データ

	var uop []*bekt  // uop: opから見たopの位置
	var ulp []*bekt  // ulp: opから見たlpの位置
	var ullp []*bekt // ullp: lpから見たlpの位置
	var ulmp []*bekt // ulmp: lpから見たmpの位置

	/*---------------higuchi add-------------------start*/

	var obsn int = 0      // OBSの総数
	var lpn int = 0       // LP(被受照面)の総数
	var opn int = 0       // OP(受照面)の総数
	var mpn int = 0       // MP(受照面)の総数　(OP+OPW)
	var monten int = 1000 // モンテカルロ法の際の射出数

	var DE float64 = 100.0 // 壁面の分割による微小四角形の辺の長さ
	var co float64 = 0.0   // 壁面への太陽光線の入射角

	var wap []float64
	var wip [][]float64

	//var uop, ulp, ullp, ulmp *bekt

	var gp [][]XYZ
	var gpn int

	var fp1 *os.File // _shadow.gchi : MPの影面積の出力
	var fp2 *os.File // _I.gchi : MPの日射量の出力
	var fp3 *os.File // _lwr.gchi : MPの長波長放射量の出力
	var fp4 *os.File // _ffactor.gchi : MPの形態係数の出力

	var BDP []*BBDP
	var obs []*OBS
	var tree []*TREE     /*-樹木データ-*/
	var poly []*POLYGN   /*--POLYGON--*/
	var shadtb []*SHADTB /*-LP面の日射遮蔽率スケジュール-*/
	var op []*P_MENN     // OP面(受光面)
	var lp []*P_MENN     // LP面(被受光面)
	var mp []*P_MENN     // MP面(OP+OPW)
	var Noplpmp NOPLPMP  // OP、LP、MPの定義数

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
	Fbmlist = ""

	Rmvls = NewRMVLS()
	Simc = NewSIMCONTL()

	Eqsys = NewEQSYS()
	Loc = NewLOCAT()
	Eqcat = NewEQCAT()

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
	Simc.Loc = Loc

	// 建築・設備システムデータ入力
	Schdl, Flout := Eeinput(
		EWKFile,
		efl_path,
		bdata, week, schtba, schnma,
		Simc, &Exsf, Rmvls, Eqcat, Eqsys,
		&Compnt,
		&Elout,
		&Elin,
		&Mpath,
		&Plist,
		&Pelm,
		&Contl,
		&Ctlif,
		&Ctlst,
		&Wd,
		&Daytm, key,
		&obsn, &BDP, &obs, &tree, &shadtb, &poly, &monten, &gpn, &DE, &Noplpmp)

	// 最大収束回数のセット
	LOOP_MAX := Simc.MaxIterate
	VAV_Count_MAX := Simc.MaxIterate

	// 動的カーテンの展開
	for i := range Rmvls.Sd {
		Sd := Rmvls.Sd[i]
		if Sd.DynamicCode != "" {
			ctifdecode(Sd.DynamicCode, Sd.Ctlif, Simc, Compnt, Mpath, &Wd, &Exsf, Schdl)
		}
	}

	if len(BDP) != 0 {

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
		// 多面体、樹木、障害物、BDPの座標変換し、LP, OPに集約する
		lp = LP_COORDNT(poly, tree, obs, BDP)
		op = OP_COORDNT(BDP, poly)

		lpn = len(lp)
		opn = len(op)

		// LPの構造体に日毎の日射遮蔽率を代入
		for _, _lp := range lp {
			for _, _shadtb := range shadtb {
				if _lp.opname == _shadtb.lpname {
					for k := 1; k < 366; k++ {
						for l := 0; l < _shadtb.indatn; l++ {
							if k >= _shadtb.ndays[l] && k <= _shadtb.ndaye[l] {
								_lp.shad[k] = _shadtb.shad[l]
								break
							}
						}
					}
				}
			}
		}

		//---- mpの総数をカウント mpは、OP面+OPW面 ---------------
		// OP面 = 授照面、OPW面 = 受照窓面
		mpn = 0
		for i := 0; i < opn; i++ {
			mpn += 1
			for j := 0; j < op[i].wd; j++ {
				mpn += 1
			}
		}

		//---窓壁のカウンター変数の初期化---
		wap = make([]float64, opn)
		wip = make([][]float64, opn)
		for i := 0; i < opn; i++ {
			if op[i].wd != 0 {
				wip[i] = make([]float64, op[i].wd)
			}
		}

		//---領域の確保   gp 地面の座標(X,Y,Z)---
		gp = make([][]XYZ, mpn)
		for i := 0; i < mpn; i++ {
			gp[i] = make([]XYZ, gpn+1)
		}

		// //---領域の確保 mp---
		// mp = make([]*P_MENN, Noplpmp.Nmp)
		// P_MENNinit(mp, mpn)

		//----OP,OPWの構造体をMPへ代入する----
		mp = DAINYUU_MP(op)

		for i := 0; i < mpn; i++ {
			fmt.Fprintf(fp1, "%s\n", mp[i].opname)
		}

		//---ベクトルの向きを判別する変数の初期化---
		//---opから見たopの位置---
		uop = make([]*bekt, opn)
		for i := range op {
			uop[i] = Newbekt(op)
		}

		//---opから見たlpの位置---
		ulp = make([]*bekt, opn)
		for i := range op {
			ulp[i] = Newbekt(lp)
		}

		//---lpから見たlpの位置---
		ullp = make([]*bekt, lpn)
		for i := range lp {
			ullp[i] = Newbekt(lp)
		}

		//---lpから見たmpの位置---
		ulmp = make([]*bekt, lpn)
		for i := range lp {
			ulmp[i] = Newbekt(mp)
		}

		//------CG確認用データ作成-------
		HOUSING_PLACE(lpn, mpn, lp, mp, RET15)

		//----前面地面代表点および壁面の中心点Gを求める--------
		GRGPOINT(mp, mpn)
		for _, _lp := range lp {
			_lp.G = GDATA(_lp)
		}

		// 20170426 higuchi add 条件追加　形態係数を計算しないパターンを組み込んだ
		if monten > 0 {
			//---LPから見た天空に対する形態係数faia算出------
			FFACTOR_LP(monten, lp, mp)
		}

		// 各MP面に方位、日射吸収率等を壁体情報から取得
		for _, _mp := range mp {
			for j := range Rmvls.Sd {
				if Rmvls.Sd[j].Sname == _mp.opname {
					_mp.exs = Rmvls.Sd[j].exs
					_mp.as = Rmvls.Sd[j].as
					_mp.alo = Rmvls.Sd[j].alo
					_mp.Eo = Rmvls.Sd[j].Eo
					break
				}
			}
		}

		// 前面地面の反射率を取得
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

		fmt.Printf("Npelm=%d Ncmalloc=%d Ncompnt=%d Nelout=%d Nelin=%d\n",
			len(Pelm), len(Compnt), len(Compnt), len(Elout), len(Elin))
	}

	Soldy = make([]float64, len(Exsf.Exs))
	Solmon = make([]float64, len(Exsf.Exs))

	DTM = float64(Simc.DTm)
	dminute = int(float64(Simc.DTm) / 60.0)
	Cff_kWh = DTM / 3600.0 / 1000.0

	for rm := range Rmvls.Room {
		Rm := Rmvls.Room[rm]
		Rm.Qeqp = 0.0
	}

	// スケジュール設定のデバッグ出力
	Schdl.dprschtable()

	// 外表面方位データのデバッグ出力
	Exsf.dprexsf()

	// 壁・窓のデバッグ出力
	Rmvls.dprwwdata()

	// 室のデバッグ出力
	Rmvls.dprroomdata()

	// 重量壁体のデバッグ出力
	Rmvls.dprballoc()

	Simc.eeflopen(Flout, efl_path)

	if DEBUG {
		fmt.Println("<<main>> eeflopen ")
	}

	if Ferr != nil {
		fmt.Fprintln(Ferr, "\n<<main>> eeflopen end")
	}

	// 壁体内部温度の初期値設定
	Rmvls.Tinit()

	if DEBUG {
		fmt.Println("<<main>> Tinit")
	}

	if Ferr != nil {
		fmt.Fprintln(Ferr, "\n<<main>> Tinit")
	}

	// ボイラ機器仕様の初期化
	Eqcat.Boicaint(Simc, Compnt, &Wd, &Exsf, Schdl)

	// システム使用機器の初期設定
	Eqsys.Mecsinit(Simc, Compnt, Exsf.Exs, &Wd, Rmvls)

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
	Sdstr := make([]*SHADSTR, mpn)
	for i := 0; i < mpn; i++ {
		Sdstr[i] = new(SHADSTR)
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
			Weatherdt(Simc, &Daytm, Loc, &Wd, Exsf.Exs, Exsf.EarthSrfFlg)

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
				if len(BDP) != 0 && monten > 0 {
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
			Exsf.Exsfsol(&Wd)

			/*==transplantation to eeslism from KAGExSUN by higuchi 070918==start*/
			if len(BDP) != 0 {

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
								SHADOW(j, DE, opn, lpn, ls, ms, ns, uop[j], ulp[j], op[j], op, lp, &wap[j], wip[j], day)
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
							// 陰面積の出力
							shadow_printf(fp1, Daytm.Mon, Daytm.Day, Daytm.Time, mpn, mp)
						}

					} else {
						//fmt.Printf("dcnt2=%d\n",dcnt) ;
						DAINYUU_SMO2(opn, mpn, op, mp, Sdstr, dcnt, mt)
						// 20170426 higuchi add 条件追加
						if dayprn {
							// 陰面積の出力
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

			// 現時刻ステップのスケジュール作成
			Eeschdlr(day, Daytm.Ttmm, Schdl, Rmvls)

			if DEBUG {
				fmt.Println("<<main>>  Eeschdlr")
			}

			// 制御で使用する状態値を計算する（集熱器の相当外気温度）
			CalcControlStatus(Eqsys, Rmvls, &Wd, &Exsf)

			// 制御情報の更新
			Contlschdlr(Contl, Mpath, Compnt)

			// 空調発停スケジュール設定が完了したら人体発熱を再計算
			for _, rm := range Rmvls.Room {
				rm.Qischdlr()
			}

			if DEBUG {
				fmt.Println("<<main>> Contlschdlr")
			}

			/***
			eloutprint(0, Nelout, Elout, Compnt);
			*****/

			// カウンターリセット
			Eqsys.VAVcountreset()
			Eqsys.Valvcountreset()
			Eqsys.Evaccountreset()

			/*---- Satoh Debug VAV  2000/12/6 ----*/
			// ここから: VAV 計算繰り返しループ
			for j := 0; j < VAV_Count_MAX; j++ {
				if DEBUG {
					fmt.Printf("\n\n====== VAV LOOP Count=%d ======\n\n\n", j)
				}
				if dayprn && Ferr != nil {
					fmt.Fprintf(Ferr, "\n\n====== VAV LOOP Count=%d ======\n\n\n", j)
				}

				VAVreset := 0
				Valvreset := 0

				// ポンプ流量設定（太陽電池ポンプのみ
				Eqsys.Pumpflow()

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

				Sysupv(Mpath, Rmvls)

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

				for n := range Rmvls.Sd {
					Rmvls.Sd[n].mrk = '!'
				}

				// システム使用機器特性式係数の計算
				Eqsys.Mecscf()

				if DEBUG {
					fmt.Println("<<main>> Mecscf")
				}

				/*======higuchi update 070918==========*/
				eeroomcf(&Wd, &Exsf, Rmvls, nday, mt)
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
					fmt.Println("<<main>> Roomvar")
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

				// 収束計算
				for i := 0; i < LOOP_MAX; i++ {
					if i == 0 {
						hcldwetmdreset(Eqsys)
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

					// 蓄熱槽特性式係数
					Stankcfv(Eqsys.Stank)

					// 特性式の係数
					Hcldcfv(Eqsys.Hcload)

					// システム方程式の作成およびシステム変数の計算
					Syseqv(Elout, &Syseq)

					Sysvar(Compnt)

					// 室温・湿度計算結果代入、室供給熱量計算
					// およびパネル入口温度代入、パネル供給熱量計算
					Roomene(Rmvls, Rmvls.Room, Rmvls.Rdpnl, &Exsf, &Wd)

					// 室負荷の計算
					Roomload(Rmvls.Room, &LDreset)

					// PCM家具の収束判定
					PCMfunchk(Rmvls.Room, &Wd, &PCMfunreset)

					// 壁体内部温度の計算と収束計算のチェック
					if Rmvls.Pcmiterate == 'y' {
						PCMwlchk(i, Rmvls, &Exsf, &Wd, &LDreset)
					}

					// 供給熱量、エネルギーの計算
					Boiene(Eqsys.Boi, &BOIreset)

					// 冷却熱量/加熱量、エネルギーの計算
					Refaene(Eqsys.Refa, &LDreset)

					// 空調負荷の計算
					Hcldene(Eqsys.Hcload, &LDreset, &Wd)

					// 供給熱量の計算
					Hccdwreset(Eqsys.Hcc, &DWreset)

					// 槽内水温、水温分布逆転の検討
					Stanktss(Eqsys.Stank, &TKreset)

					// 内部温度、熱量の計算
					Evacene(Eqsys.Evac, &Evacreset)

					if BOIreset+LDreset+DWreset+TKreset+Evacreset+PCMfunreset == 0 {
						break
					}
				}

				if i == LOOP_MAX {
					fmt.Printf("収束しませんでした。 MAX=%d\n", LOOP_MAX)
				}

				// 供給熱量の計算
				Hccene(Eqsys.Hcc)

				// 風量の計算
				VAVene(Eqsys.Vav, &VAVreset)
				Valvene(Eqsys.Valv, &Valvreset)

				if VAVreset == 0 && Valvreset == 0 {
					break
				}

				// カウントアップ
				Eqsys.VAVcountinc()
				Eqsys.Valvcountinc()

				// 風量が変わったら電気蓄熱暖房器の係数を再計算
				Stheatcfv(Eqsys.Stheat)
			}
			// ここまで: VAV 計算繰り返しループ

			// 太陽電池内蔵壁体の発電量計算
			CalcPowerOutput(Rmvls.Sd, &Wd, &Exsf)

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
			Eqsys.Mecsene()

			/***********************
			fmt.Printf("Mecsene en\n")
			/***********************/

			if DEBUG {
				mecsxprint(Eqsys)
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
			RMwlt(Rmvls.Mw)

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
			Eeprinth(&Daytm, Simc, Flout, Rmvls, &Exsf, Mpath, Eqsys, &Wd)

			if DEBUG {
				fmt.Printf("xxxmain 6\n")
			}

			if Daytm.Ddpri != 0 {
				// 室の日集計、月集計
				Roomday(Daytm.Mon, Daytm.Day, day, Daytm.Ttmm, Rmvls.Room, Rmvls.Rdpnl, Simc.Dayend)
				if Simc.Helmkey == 'y' {
					Helmdy(day, Rmvls.Room, &Rmvls.Qetotal)
				}

				Compoday(Daytm.Mon, Daytm.Day, day, Daytm.Ttmm, Eqsys, Simc.Dayend)
				/**   if (Nqrmpri > 0)  **/
				Qrmsum(Daytm.Day, Rmvls.Room, Rmvls.Qrm, Rmvls.Trdav, Rmvls.Qrmd)

				if DEBUG {
					fmt.Printf("xxxmain 7\n")
				}

				// 気象データの日集計、月集計
				Wdtsum(Daytm.Mon, Daytm.Day, day, Daytm.Ttmm, &Wd, Exsf.Exs, &Wdd, &Wdm, Soldy, Solmon, Simc)
			}
			if DEBUG {
				fmt.Printf("xxxmain 8\n")
			}

			if DEBUG {
				Rmvls.xprtwsrf()
				Rmvls.xprrmsrf()
				Rmvls.xprtwall()
			}

			mm += dminute

			// 時刻ループの最後
		}

		// 日集計の出力
		Eeprintd(&Daytm, Simc, Flout, Rmvls, Exsf.Exs, Soldy, Eqsys, &Wdd)
		/*****fmt.Printf("xxxmain 9\n")*****/

		//if (Daytm.Mon == 4 && Daytm.Day == 25)
		//	printf("debug\n");

		// 月集計の出力
		if IsEndDay(Daytm.Mon, Daytm.Day, Daytm.DayOfYear, Simc.Dayend) && Daytm.Ddpri != 0 {
			//fmt.Printf("月集計出力\n")
			Eeprintm(&Daytm, Simc, Flout, Rmvls, Exsf.Exs, Solmon, Eqsys, &Wdm)
		}

	}
	// 月－時刻別集計値の出力
	Eeprintmt(Simc, Flout, Eqsys, Rmvls.Rdpnl)

	if DEBUG {
		fmt.Printf("メモリ領域の解放\n")
	}

	Eeflclose(Flout)

	/*------------------higuchi add---------------------start*/
	if len(BDP) != 0 {

		defer fp1.Close()
		defer fp2.Close()
		defer fp3.Close()
		defer fp4.Close()
	}

	/*---------------------higuchi 1999.7.21-----------end*/
}

/*
Newbekt (New Bekt Object)

この関数は、日影計算や日射量計算で用いられる「bekt」構造体を初期化します。
「bekt」構造体は、ある面から見た他の面の相対的な位置関係を格納するために用いられます。

建築環境工学的な観点:
- **幾何学的関係のモデル化**: 建物の日影や日射は、
  建物や周囲の障害物の幾何学的関係によって決まります。
  この関数は、`ps`（位置関係を示す行列）を初期化することで、
  ある面から見た他の面の相対的な位置関係をモデル化します。
- **日影計算の基礎**: `uop`（opから見たopの位置）、`ulp`（opから見たlpの位置）、
  `ullp`（lpから見たlpの位置）、`ulmp`（lpから見たmpの位置）といった`bekt`構造体は、
  日影計算において、太陽光線が障害物によって遮られるかどうかを判定するために用いられます。

この関数は、建物の日射環境を正確にモデル化し、
日射熱取得の抑制、冷房負荷の軽減、
昼光利用の最適化、および日影計算を行うための重要な役割を果たします。
*/
func Newbekt(op []*P_MENN) *bekt {
	uop := new(bekt)
	uop.ps = make([][]float64, len(op))
	for j := range op {
		uop.ps[j] = make([]float64, op[j].polyd)
	}
	return uop
}

/*
matiniti (Matrix Initialization for Integer Array)

この関数は、整数型の配列`A`の全ての要素をゼロに初期化します。

建築環境工学的な観点:
- **データ初期化の重要性**: シミュレーションでは、
  計算の開始前に全ての変数を既知の状態に初期化することが重要です。
  これにより、過去の計算結果が現在の計算に影響を与えたり、
  未定義の値によるエラーが発生したりするのを防ぎます。
- **カウンタやフラグのリセット**: この関数は、
  例えば、日影計算における日影の有無を示すフラグや、
  特定のイベントの発生回数をカウントするカウンタなどをリセットする際に用いられます。

この関数は、シミュレーションの正確性と安定性を確保するための基本的な役割を果たします。
*/
func matiniti(A []int, N int) {
	for i := 0; i < N; i++ {
		A[i] = 0
	}
}

/*
P_MENNinit (P_MENN Initialization)

この関数は、`P_MENN`構造体の配列を初期化します。
`P_MENN`構造体は、日影計算や日射量計算で用いられる「受光面（Opening Plane, OP）」や
「被受照面（Light-Receiving Plane, LP）」の幾何学的情報を格納するために用いられます。

建築環境工学的な観点:
- **幾何学的モデルの準備**: 建物の日射環境をシミュレーションする前に、
  窓、壁、日よけ、障害物などの幾何学的情報を格納するデータ構造を準備します。
  この関数は、`NewP_MENN`関数を呼び出すことで、
  各`P_MENN`構造体をデフォルト値で初期化します。
- **日影計算と日射量計算の基礎**: これらの構造体は、
  太陽位置と組み合わせて、
  窓面への影の形状と面積、
  および各表面への日射入射量を計算するために用いられます。

この関数は、建物の日射環境を正確にモデル化し、
日射熱取得の抑制、冷房負荷の軽減、
昼光利用の最適化、および日影計算を行うための重要な役割を果たします。
*/
func P_MENNinit(_pm []*P_MENN, N int) {
	for i := 0; i < N; i++ {
		_pm[i] = NewP_MENN()
	}
}

/*
NewP_MENN (New P_MENN Object)

この関数は、新しい`P_MENN`構造体を初期化します。
`P_MENN`構造体は、日影計算や日射量計算で用いられる「受光面（Opening Plane, OP）」や
「被受照面（Light-Receiving Plane, LP）」の幾何学的情報を格納するために用いられます。

建築環境工学的な観点:
- **幾何学的パラメータの初期化**: `P_MENN`構造体は、
  - `opname`: 面の名称。
  - `rgb`: 面の色（RGB値）。
  - `faiwall`: 壁面への日射入射角特性。
  - `wd`: 窓の数。
  - `exs`: 外部日射面の種類。
  - `grpx`: 前面地面の代表点までの距離。
  - `faia`, `faig`: 天空に対する形態係数、地面に対する形態係数。
  - `sum`: 影面積。
  - `ref`, `refg`: 面の反射率、前面地面の反射率。
  - `wa`, `wb`: 面の方位角、傾斜角。
  - `Ihor`, `Idre`, `Idf`, `Iw`: 水平、直達、拡散、全天日射量。
  - `Reff`, `rn`, `Te`, `Teg`: 有効放射率、夜間放射、相当外気温度、相当地面温度。
  - `shad`: 年間を通じた日影の有無。
  - `alo`, `as`, `Eo`: 外表面熱伝達率、日射吸収率、放射率。
  - `polyd`: 頂点数。
  - `sbflg`: 日よけフラグ。
  - `P`: 頂点座標。
  - `opw`: 窓のデータ。
  といった、日影計算や日射量計算に必要な様々な幾何学的・熱的・光学的パラメータを格納します。
  これらのパラメータをデフォルト値で初期化することで、
  後続のデータ入力や計算が正しく行われるための準備が整います。

この関数は、建物の日射環境を正確にモデル化し、
日射熱取得の抑制、冷房負荷の軽減、
昼光利用の最適化、および日影計算を行うための重要な役割を果たします。
*/
func NewP_MENN() *P_MENN {
	const pmax = 200

	pm := new(P_MENN)
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
	return pm
}

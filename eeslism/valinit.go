package eeslism

/*
NewSIMCONTL (New Simulation Control Object)

この関数は、建物のエネルギーシミュレーションの実行を制御するための
`SIMCONTL`構造体を初期化します。

建築環境工学的な観点:
- **シミュレーションの基本設定の初期化**: シミュレーションの期間、時間間隔、出力設定、
  および気象データや入力ファイルに関する情報などを管理する`SIMCONTL`構造体の各フィールドを、
  デフォルト値で初期化します。
  これにより、後続のデータ入力や計算が正しく行われるための準備が整います。
- **計算精度と効率の制御**: `DTm`（計算時間間隔）や`MaxIterate`（最大収束回数）といったパラメータは、
  シミュレーションの精度と計算効率に直接影響します。
  これらの初期値は、一般的なシミュレーションの要件を満たすように設定されます。
- **出力設定の準備**: `Dayprn`（データ出力日）配列を初期化することで、
  特定の日のみ詳細な結果を出力する設定を準備します。
  これにより、ユーザーは必要な情報を効率的に取得し、
  分析や検証を容易にします。

この関数は、建物のエネルギーシミュレーションの実行を制御し、
シミュレーションの正確性、効率性、および出力内容を決定するための重要な役割を果たします。
*/
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

/*
NewCOMPNT (New Component Object)

この関数は、建物のエネルギーシミュレーションにおける各機器（コンポーネント）の
`COMPNT`構造体を初期化します。

建築環境工学的な観点:
- **機器の基本情報の初期化**: 機器名称（`Name`）、設置室名称（`Roomname`）、
  機器タイプ（`Eqptype`）、周囲条件名称（`Envname`）、
  入出力の識別記号（`Idi`, `Ido`）など、
  機器の基本情報をデフォルト値で初期化します。
  これにより、後続のデータ入力や計算が正しく行われるための準備が整います。
- **熱湿気同時交換の考慮**: `Airpathcpy`は、
  空気経路の場合に湿度経路用にパスをコピーするかどうかを示すフラグです。
  これにより、熱湿気同時交換を行う機器のモデル化を適切に行えます。
- **制御情報の初期化**: `Control`は、
  機器の運転制御情報（ON/OFFなど）を初期化します。
  これにより、機器の運転状態をモデル化できます。

この関数は、建物のエネルギーシミュレーションにおいて、
多様なコンポーネントを統合的にモデル化し、
システム全体のエネルギー消費量予測、省エネルギー対策の検討、
および最適な設備システム設計を行うための重要な役割を果たします。
*/
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

/*
NewEloutSlice (New ELOUT Slice)

この関数は、指定された数（`n`）の`ELOUT`構造体のスライス（配列）を作成し、
各要素を初期化します。
`ELOUT`構造体は、機器の出口における熱媒の状態（温度、流量など）を格納するために用いられます。

建築環境工学的な観点:
- **機器の出力ポートのモデル化**: 建物のエネルギーシミュレーションでは、
  各機器が複数の出力ポートを持つ場合があります。
  この関数は、これらの出力ポートをモデル化するためのデータ構造を効率的に作成します。
- **熱媒の状態の追跡**: 各`ELOUT`構造体は、
  機器の出口における熱媒の温度、流量、熱量、連立方程式の解などを格納します。
  これにより、熱媒がシステム内をどのように流れ、
  その状態がどのように変化しているかを追跡できます。

この関数は、建物のエネルギーシミュレーションにおいて、
熱搬送システムや空調システムの構成と運転を正確にモデル化するための重要な役割を果たします。
*/
func NewEloutSlice(n int) []*ELOUT {
	s := make([]*ELOUT, n)
	for i := 0; i < n; i++ {
		s[i] = NewElout()
	}
	return s
}

/*
NewElout (New ELOUT Object)

この関数は、新しい`ELOUT`構造体を初期化します。
`ELOUT`構造体は、機器の出口における熱媒の状態（温度、流量など）を格納するために用いられます。

建築環境工学的な観点:
- **機器の出力ポートのモデル化**: 各機器の出口における熱媒の状態を定義します。
  - `Ni`: 入口の数。
  - `G`: 流量。
  - `Co`, `Coeffo`, `Coeffin`: 連立方程式の係数。
  - `Control`: 経路の制御情報。
  - `Id`: 出口の識別番号。
  - `Fluid`, `Sysld`: 流体の種類、負荷計算フラグ。
  - `Q`, `Sysv`, `Load`: 熱量、連立方程式の答え、負荷。
  - `Cmp`: 機器出口が属する機器への逆参照。
  - `Elins`: 機器出口が関連する機器入口。
  - `Lpath`: 機器出口が属する末端経路。
  - `Eldobj`, `Emonitr`: 関連するELOUTオブジェクト。
  - `Pelmoid`: 終端の割り当て完了フラグ。
  これらのパラメータをデフォルト値で初期化することで、
  後続のデータ入力や計算が正しく行われるための準備が整います。

この関数は、建物のエネルギーシミュレーションにおいて、
熱搬送システムや空調システムの構成と運転を正確にモデル化するための重要な役割を果たします。
*/
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

/*
NewElinSlice (New ELIN Slice)

この関数は、指定された数（`n`）の`ELIN`構造体のスライス（配列）を作成し、
各要素を初期化します。
`ELIN`構造体は、機器の入口における熱媒の状態（温度、流量など）を格納するために用いられます。

建築環境工学的な観点:
- **機器の入力ポートのモデル化**: 建物のエネルギーシミュレーションでは、
  各機器が複数の入力ポートを持つ場合があります。
  この関数は、これらの入力ポートをモデル化するためのデータ構造を効率的に作成します。
- **熱媒の状態の追跡**: 各`ELIN`構造体は、
  機器の入口における熱媒の温度、流量、熱量、連立方程式の解などを格納します。
  これにより、熱媒がシステム内をどのように流れ、
  その状態がどのように変化しているかを追跡できます。

この関数は、建物のエネルギーシミュレーションにおいて、
熱搬送システムや空調システムの構成と運転を正確にモデル化するための重要な役割を果たします。
*/
func NewElinSlice(n int) []*ELIN {
	s := make([]*ELIN, n)
	for i := 0; i < n; i++ {
		s[i] = NewElin()
	}
	return s
}

/*
NewElin (New ELIN Object)

この関数は、新しい`ELIN`構造体を初期化します。
`ELIN`構造体は、機器の入口における熱媒の状態（温度、流量など）を格納するために用いられます。

建築環境工学的な観点:
- **機器の入力ポートのモデル化**: 各機器の入口における熱媒の状態を定義します。
  - `Id`: 入口の識別番号。
  - `Sysvin`: 連立方程式の答え。
  - `Upo`, `Upv`: 上流の機器の出口へのポインター。
  - `Lpath`: 機器入口が属する末端経路。
  これらのパラメータをデフォルト値で初期化することで、
  後続のデータ入力や計算が正しく行われるための準備が整います。

この関数は、建物のエネルギーシミュレーションにおいて、
熱搬送システムや空調システムの構成と運転を正確にモデル化するための重要な役割を果たします。
*/
func NewElin() *ELIN {
	Ei := new(ELIN)
	Ei.Id = ' '
	Ei.Sysvin = 0.0
	Ei.Upo, Ei.Upv = nil, nil
	Ei.Lpath = nil
	return Ei
}

/*
NewPLIST (New Path List Object)

この関数は、新しい末端経路（`PLIST`）のデータ構造を初期化します。
末端経路は、システム経路（`MPATH`）を構成する個々の経路であり、
経路内の機器や流量に関する情報を定義します。

建築環境工学的な観点:
- **熱搬送経路のモデル化**: 建物のエネルギーシミュレーションでは、
  熱媒が流れる経路をモデル化し、
  熱供給量や熱回収量を計算します。
  この関数は、末端経路の名称（`Name`）、種別（`Type`）、
  制御情報（`Control`）、流量（`G`）、流量分配比（`Rate`）など、
  経路に関する情報をデフォルト値で初期化します。
- **流量制御の考慮**: `UnknownFlow`は、
  末端経路の流量が未知であるか既知であるかを示すフラグです。
  `Go`は、流量の計算に使用される係数へのポインターです。
  これにより、様々な流量制御方式をモデル化できます。
- **バッチ運転の考慮**: `Batch`は、
  バッチ運転を行う蓄熱槽のある経路であるかを示すフラグです。
  これにより、特殊な運転モードをモデル化できます。
- **空気系統の考慮**: `Plistt`と`Plistx`は、
  空気系統の場合に温度系統と湿度系統を関連付けるためのポインターです。
  これにより、熱湿気同時交換を正確にモデル化できます。

この関数は、建物のエネルギーシミュレーションにおいて、
熱搬送システムや空調システムの構成と運転を正確にモデル化するための重要な役割を果たします。
*/
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

/*
NewPELM (New Path Element Object)

この関数は、新しい経路要素（`PELM`）のデータ構造を初期化します。
経路要素は、末端経路内の個々の機器（コンポーネント）を定義します。

建築環境工学的な観点:
- **経路内の機器のモデル化**: 建物のエネルギーシミュレーションでは、
  熱媒が流れる経路内に複数の機器が存在します。
  この関数は、経路要素の入出力の識別番号（`Co`, `Ci`）、
  関連するコンポーネント（`Cmp`）、
  機器の出口（`Out`）と入口（`In`）へのポインターなど、
  経路内の機器に関する情報をデフォルト値で初期化します。

この関数は、建物のエネルギーシミュレーションにおいて、
熱搬送システムや空調システムの構成と運転を正確にモデル化するための重要な役割を果たします。
*/
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

/*
NewMPATH (New Main Path Object)

この関数は、新しいシステム経路（`MPATH`）のデータ構造を初期化します。
システム経路は、建物のエネルギーシステムにおける熱媒の流れる主要な経路を定義します。

建築環境工学的な観点:
- **システム経路のモデル化**: 建物のエネルギーシミュレーションでは、
  熱源設備、熱搬送設備、空調設備など、
  様々な機器が配管やダクトで接続されて構成されます。
  この関数は、システム経路の名称（`Name`）、
  システムの種類（`Sys`）、流体種別（`Fluid`）、
  制御情報（`Control`）など、
  システム経路に関する情報をデフォルト値で初期化します。
- **熱湿気同時交換の考慮**: `Mpair`は、
  空気系統の場合に温度経路と湿度経路を関連付けるためのポインターです。
  これにより、熱湿気同時交換を正確にモデル化できます。
- **流量計算の考慮**: `NGv`（ガス導管数）、`NGv2`（開口率が2%未満のガス導管数）、
  `Ncv`（制御弁数）、`Lvcmx`（制御弁の接続数の最大値）、
  `G0`（流量比率設定時の既知流量へのポインタ）、
  `Rate`（流量比率フラグ）は、
  システム経路内の流量計算をモデル化します。

この関数は、建物のエネルギーシミュレーションにおいて、
熱搬送システムや空調システムの構成と運転を正確にモデル化するための重要な役割を果たします。
*/
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

/*
Exsfinit (External Surface Initialization)

この関数は、外部日射面（壁、屋根、窓、地盤など）の`EXSF`構造体を初期化します。

建築環境工学的な観点:
- **外部日射面の基本情報の初期化**: 外部日射面の名称（`Name`）、
  タイプ（`Typ`）、方位角（`Wa`）、傾斜角（`Wb`）、
  地盤反射率（`Rg`）、天空形態係数（`Fs`）、
  地中深さ（`Z`）、土の熱拡散率（`Erdff`）など、
  外部日射面の基本情報をデフォルト値で初期化します。
  これにより、後続のデータ入力や計算が正しく行われるための準備が整います。
- **熱伝達率の初期化**: `Alo`（外表面総合熱伝達率）を初期化し、
  `Alotype`（外表面熱伝達率の設定方法）を`Alotype_Fix`（固定値）に設定します。
  これにより、外表面からの熱伝達をモデル化できます。
- **日射関連パラメータの初期化**: `Cinc`（入射角のcos）、
  `Tazm`（見掛けの方位角のtan）、
  `Tprof`（プロファイル角のtan）、
  `Idre`（直達日射）、`Idf`（拡散日射）、
  `Iw`（全日射）、`Rn`（夜間輻射）、
  `Tearth`（地中温度）など、
  日射関連のパラメータを初期化します。
  これにより、日射熱取得量をモデル化できます。

この関数は、建物のエネルギーシミュレーションにおいて、
外部日射面からの熱損失・熱取得を正確にモデル化し、
日射熱取得の抑制、冷房負荷の軽減、
昼光利用の最適化、および省エネルギー対策の検討を行うための重要な役割を果たします。
*/
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

/*
Syseqinit (System Equation Initialization)

この関数は、システム方程式の`SYSEQ`構造体を初期化します。

建築環境工学的な観点:
- **システム方程式の準備**: 建物のエネルギーシミュレーションでは、
  各機器や室の熱収支を連立方程式として解くことで、
  温度や流量などの未知数を計算します。
  この関数は、システム方程式の係数行列や定数項を格納する`SYSEQ`構造体を初期化し、
  計算の準備をします。

この関数は、建物のエネルギーシミュレーションにおいて、
熱負荷計算、エネルギー消費量予測、
および省エネルギー対策の検討を行うための重要な役割を果たします。
*/
func Syseqinit(S *SYSEQ) {
	S.A = ' '
}

/*
NewEQSYS (New Equipment System Object)

この関数は、建物のエネルギーシミュレーションにおける設備機器システム全体の
`EQSYS`構造体を初期化します。

建築環境工学的な観点:
- **設備機器システムの統合**: 建物のエネルギーシステムは、
  熱源設備、熱搬送設備、空調設備など、様々な機器から構成されます。
  この関数は、これらの機器を種類ごとにリストとして保持する`EQSYS`構造体を初期化します。
  これにより、システム全体を統合的にモデル化し、
  機器間の相互作用やエネルギーフローを分析できます。
- **機器リストの初期化**: 冷温水コイル（`Hcc`）、ボイラー（`Boi`）、
  冷凍機（`Refa`）、太陽熱集熱器（`Coll`）、配管（`Pipe`）、
  蓄熱槽（`Stank`）、熱交換器（`Hex`）、ポンプ（`Pump`）、
  流入境界条件（`Flin`）、空調負荷（`Hcload`）、VAV（`Vav`）、
  顕熱蓄熱器（`Stheat`）、全熱交換器（`Thex`）、弁（`Valv`）、
  カロリーメータ（`Qmeas`）、太陽光発電（`PVcmp`）、
  外気処理VAV（`OMvav`）など、
  様々な機器のリストを空のスライスとして初期化します。

この関数は、建物のエネルギーシミュレーションにおいて、
設備機器システム全体をモデル化し、
エネルギー消費量予測、省エネルギー対策の検討、
および最適な設備システム設計を行うための重要な役割を果たします。
*/
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

/*
NewRMVLS (New Room and Wall List Object)

この関数は、建物のエネルギーシミュレーションにおける室と壁体に関する
`RMVLS`構造体を初期化します。

建築環境工学的な観点:
- **室と壁体の統合管理**: 建物の熱負荷計算では、
  室の熱的挙動と壁体の熱的挙動を統合的に扱う必要があります。
  この関数は、室（`Room`）、壁体（`Wall`）、窓（`Window`）、
  日よけ（`Snbk`）、放射パネル（`Rdpnl`）、
  要素別熱損失・熱取得（`Qrm`, `Qrmd`, `Qetotal`）、
  壁体内部温度（`Mw`）、PCM（`PCM`）など、
  室と壁体に関する様々なデータをリストとして保持する`RMVLS`構造体を初期化します。
- **熱負荷計算の基礎**: これらのデータは、
  室の熱負荷計算、壁体内部温度計算、
  日射熱取得計算などの基礎となります。

この関数は、建物のエネルギーシミュレーションにおいて、
室と壁体の熱的挙動を正確にモデル化し、
熱負荷計算、エネルギー消費量予測、
および省エネルギー対策の検討を行うための重要な役割を果たします。
*/
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

/*
VPTRinit (Value Pointer Initialization)

この関数は、`VPTR`構造体を初期化します。
`VPTR`構造体は、シミュレーション中の様々な変数へのポインターと、
そのポインターが指す値の型を格納するために用いられます。

建築環境工学的な観点:
- **変数管理の柔軟性**: シミュレーションモデルでは、
  温度、流量、熱量など、様々な種類の変数を扱います。
  この関数は、`VPTR`構造体の`Type`（値の型）と`Ptr`（値へのポインター）を初期化することで、
  これらの変数を統一的に管理し、
  制御ロジックやデータ処理において柔軟にアクセスできるようにします。

この関数は、建物のエネルギーシミュレーションにおいて、
変数管理を効率的に行い、
シミュレーションモデルの構築を支援するための重要な役割を果たします。
*/
func VPTRinit(v *VPTR) {
	v.Type = ' '
	v.Ptr = nil
}

/*
TMDTinit (Time Data Initialization)

この関数は、時間データ（年、月、日、曜日、時刻など）を格納する`TMDT`構造体を初期化します。

建築環境工学的な観点:
- **時間管理の準備**: シミュレーションでは、
  各時間ステップの日付と時刻を正確に管理する必要があります。
  この関数は、`TMDT`構造体の各フィールドをデフォルト値で初期化することで、
  時間管理の準備をします。
- **出力の準備**: シミュレーション結果の出力において、
  日付と時刻を正確に表示するために、
  この構造体が用いられます。

この関数は、建物のエネルギーシミュレーションにおいて、
時間管理を正確に行い、
シミュレーション結果の出力の品質を向上させるための重要な役割を果たします。
*/
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

/*
NewLOCAT (New Location Object)

この関数は、地域データ（緯度、経度、標準子午線、地盤温度など）を格納する`LOCAT`構造体を初期化します。

建築環境工学的な観点:
- **地域情報の重要性**: 建物のエネルギーシミュレーションでは、
  その建物の位置する地域の気象条件が、
  熱負荷やエネルギー消費量に大きく影響します。
  この関数は、緯度（`Lat`）、経度（`Lon`）、
  標準子午線（`Ls`）などの地理的情報をデフォルト値で初期化し、
  太陽位置計算や日射量計算の基礎とします。
- **地盤温度の考慮**: `Tgrav`（地盤温度）や`DTgr`（地盤温度の時定数）は、
  地盤からの熱伝達をモデル化する際に用いられます。
  これは、地下室や基礎からの熱損失・熱取得を評価する上で重要です。
- **給水温度の考慮**: `Twsup`は、
  給水温度の月別データであり、
  給湯負荷計算などに用いられます。

この関数は、建物のエネルギーシミュレーションにおいて、
外部環境条件を正確にモデル化し、
熱負荷計算、エネルギー消費量予測、
および省エネルギー対策の検討を行うための重要な役割を果たします。
*/
func NewLOCAT() *LOCAT {
	L := new(LOCAT)
	L.Name = ""
	L.Lat, L.Lon, L.Ls, L.Tgrav, L.DTgr = -999.0, -999.0, -999.0, -999.0, -999.0
	L.Daymxert = -999

	matinitx(L.Twsup[:], 12, -999.0)

	return L
}

/*
NewEQCAT (New Equipment Catalog Object)

この関数は、設備機器のカタログデータを格納する`EQCAT`構造体を初期化します。

建築環境工学的な観点:
- **機器カタログの準備**: 建物のエネルギーシミュレーションでは、
  様々な種類の設備機器の性能データを参照する必要があります。
  この関数は、冷温水コイル、ボイラー、冷凍機、太陽熱集熱器、配管、蓄熱槽、熱交換器、ポンプ、VAV、顕熱蓄熱器、全熱交換器、太陽光発電、外気処理VAVなど、
  様々な機器のカタログデータをリストとして保持する`EQCAT`構造体を初期化します。
- **機器性能のデータベース**: この構造体は、
  実質的に建物の設備機器に関するデータベースとして機能します。
  各機器の定格能力、効率、部分負荷特性、制御方法などの情報が格納され、
  シミュレーションの際に参照されます。

この関数は、建物のエネルギーシミュレーションにおいて、
多様な設備機器の性能を正確にモデル化し、
システム全体のエネルギー消費量予測、省エネルギー対策の検討、
および最適な設備システム設計を行うための重要な役割を果たします。
*/
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

/*
MtEdayinit (Monthly-Time-of-Day Energy Data Initialization)

この関数は、月・時刻別で集計されるエネルギー量（電力消費量など）のデータ構造（`EDAY`）を初期化します。

建築環境工学的な観点:
- **月・時刻別のエネルギー集計の準備**: 建物のエネルギー消費量を月・時刻別で評価するためには、
  各月、各時刻の集計値をゼロにリセットする必要があります。
  この関数は、積算値（`D`）、最大値（`Mx`）、
  および運転時間回数（`Hrs`）を初期化します。
- **デマンドサイドマネジメント**: 月・時刻別のエネルギー消費量データは、
  デマンドサイドマネジメント（DSM）戦略を策定する上で非常に有用です。
  例えば、ピーク時間帯の電力消費量を削減するための運転戦略を検討したり、
  蓄熱システムや再生可能エネルギーの導入効果を評価したりする際に役立ちます。

この関数は、建物のエネルギー消費量を月・時刻別で詳細に分析し、
運用改善や省エネルギー対策の効果評価を行うための重要な役割を果たします。
*/
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

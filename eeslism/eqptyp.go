package eeslism

type QMEAS struct {
	Fluid   FliudType // 流体種別
	Id      int
	Name    string
	Cmp     *COMPNT
	Th      *float64
	Tc      *float64
	G       *float64
	Xh      *float64
	Xc      *float64
	PlistG  *PLIST // 接続している末端経路への参照 for `G`
	PlistTh *PLIST // 接続している末端経路への参照 for `H`
	Nelmh   int    // 接続している末端経路への参照時のインデックス番号 for `H`
	PlistTc *PLIST // 接続している末端経路への参照 for `C`
	Nelmc   int    // 接続している末端経路への参照時のインデックス番号 for `C`
	Plistxh *PLIST
	Plistxc *PLIST
	Qs      float64
	Ql      float64
	Qt      float64
	Thdy    SVDAY
	Tcdy    SVDAY
	xhdy    SVDAY
	xcdy    SVDAY
	Qdys    QDAY
	Qdyl    QDAY
	Qdyt    QDAY
	mThdy   SVDAY
	mTcdy   SVDAY
	mxhdy   SVDAY
	mxcdy   SVDAY
	mQdys   QDAY
	mQdyl   QDAY
	mQdyt   QDAY
	// Pelmh  *PELM
	// Pelmc  *PELM
	// idh rune
	// idc rune
	// cmph  *COMPNT
	// cmpc  *COMPNT
	// Plist  *PLIST
}

func NewQMEAS() QMEAS {
	return QMEAS{
		Name:    "",
		Cmp:     nil,
		Th:      nil,
		Tc:      nil,
		G:       nil,
		PlistG:  nil,
		PlistTc: nil,
		PlistTh: nil,
		Plistxc: nil,
		Plistxh: nil,
		Xc:      nil,
		Xh:      nil,
		Id:      0,
		Nelmc:   -999,
		Nelmh:   -999,
	}
}

type ACS struct {
	W, T, X, C float64
}

type EVACCA struct {
	Name       string  // カタログ名
	N          int     // 気化冷却器分割数
	Awet, Adry float64 // Wet側、Dry側の境界壁接触面積[m2]
	hwet, hdry float64 // Wet側、Dry側の境界壁の対流熱伝達率[W/m2K]
	Nlayer     int     // 積層数
}

type EVAC struct {
	Name                string  // 機器名称
	Cat                 *EVACCA // 機器仕様
	Cmp                 *COMPNT
	Tdryi, Tdryo        float64   // Dry側出入口温度[℃]
	Tweti, Tweto        float64   // Wet側出入口温度[℃]
	Xdryi, Xdryo        float64   // Dry側出入口絶対湿度[kg/kg']
	Xweti, Xweto        float64   // Wet側出入口絶対湿度[kg/kg']
	RHdryi, RHdryo      float64   // Dri側出入口相対湿度[%]
	RHweti, RHweto      float64   // Wet側出入口相対湿度[%]
	Gdry, Gwet          float64   // Dry側、Wet側風量[kg/s]
	M, Kx               []float64 // i層蒸発量[kg/m2 s]、物質移動係数[kg/m2 s (kg/kg')]
	Tdry, Twet          []float64 // Dry側、Wet側の空気温度[℃]
	Xdry, Xwet          []float64 // Dry側、Wet側の空気絶対湿度[kg/kg']
	Xs                  []float64 // 濡れ面近傍の絶対湿度（境界面温度における飽和絶対湿度）[kg/kg']
	Ts                  []float64 // 境界面の表面温度[℃]（Wet側、Dry側は同じ温度）
	RHwet, RHdry        []float64 // Dry側、Wet側の相対湿度[%]
	Qsdry, Qldry, Qtdry float64   // Dry側顕熱、潜熱、全熱[W]
	Qswet, Qlwet, Qtwet float64   // Wet側顕熱、潜熱、全熱[W]
	UX, UXC             []float64
	Count               int // 計算回数。１ステップで最低２回は計算する
	//UXdry, UXwet, UXC []float64 // 状態値計算用行列
	Tdryidy, Tdryody, Twetidy, Twetody                   SVDAY
	Qsdrydy, Qldrydy, Qtdrydy, Qswetdy, Qlwetdy, Qtwetdy QDAY
}

// Satoh追加　バッチ式デシカント空調機 2013/10/20
type DESICA struct {
	name string  // カタログ名
	r    float64 // シリカゲル平均直径[cm]
	ms   float64 // シリカゲル質量[g]
	rows float64 // シリカゲル充填密度[g/cm3]
	Uad  float64 // シリカゲル槽壁面の熱貫流率[W/m2K]
	A    float64 // シリカゲル槽表面積[m2]
	Vm   float64 // モル容量[cm3/mol]
	eps  float64 // シリカゲルの空隙率
	P0   float64 // シリカゲルの飽和吸湿量[kg(water)/kg(silica gel)]
	kp   float64 // Polanyi DR 定数[cm3/(mol K)2
	cps  float64 // シリカゲルの比熱[J/kgK]
}

type DESI struct {
	Name        string
	Cat         *DESICA
	Cmp         *COMPNT
	Tain, Taout float64 // 空気の出入口温度[℃]
	Xain, Xaout float64 // 空気の出入口絶対湿度[kg/kg']
	UA          float64 // デシカント槽の熱損失係数[W/K]
	Asa         float64 // シリカゲルと槽内空気の熱伝達面積[m2]
	//Ts float64				// シリカゲル温度[℃]
	//Xs float64				// シリカゲル表面の絶対湿度[kg/kg']
	Tsold, Xsold            float64 // 前時刻の状態値
	Ta                      float64 // デシカント槽内空気温度[℃]
	Xa                      float64 // デシカント槽内絶対湿度[kg/kg']
	RHold                   float64 // 前時刻の相対湿度状態値
	Pold                    float64 // 前時刻の吸湿量[kg(water)/kg(silica gel)]
	CG                      float64
	Qloss                   float64  // デシカント槽表面からの熱損失[W]
	Qs, Ql, Qt              float64  // 顕熱、潜熱、全熱[W]
	Tenv                    *float64 // 周囲温度のポインタ[℃]
	UX                      []float64
	UXC                     []float64
	Room                    *ROOM // デシカント槽設置室構造体
	Tidy, xidy              SVDAY // 入口温度日集計
	Tody, xody              SVDAY // 出口温度日集計
	Tsdy, xsdy              SVDAY // 蓄熱体温度日集計
	Qsdy, Qldy, Qtdy, Qlsdy QDAY  // 熱量日集計
}

type THEXCA struct {
	Name string  // カタログ名
	et   float64 // 温度効率
	eh   float64 // エンタルピ効率
}

type THEX struct {
	Name     string // 機器名称
	Type     rune   // t:顕熱交換型　h：全熱交換型
	Cat      *THEXCA
	Cmp      *COMPNT
	ET       float64 // 温度効率
	EH       float64 // エンタルピ効率
	CGe      float64 // 還気側熱容量流量 [W/K]
	Ge       float64 // 還気側流量 [kg/s]
	CGo      float64 // 外気側熱容量流量 [W/K]
	Go       float64 // 外気側流量 [kg/s]
	Tein     float64 // 還気側入口温度 [℃]
	Teout    float64 // 還気側出口温度 [℃]
	Toin     float64 // 外気側入口温度 [℃]
	Toout    float64 // 外気側出口温度 [℃]
	Xein     float64 // 還気側入口絶対湿度 [kg/kg']
	Xeout    float64 // 還気側出口絶対湿度 [kg/kg']
	Xoin     float64 // 外気側入口絶対湿度 [kg/kg']
	Xoout    float64 // 外気側出口絶対湿度 [kg/kg']
	Hein     float64 // 還気側入口エンタルピー [J/kg]
	Heout    float64 // 還気側出口エンタルピー [J/kg]
	Hoin     float64 // 外気側入口エンタルピー [J/kg]
	Hoout    float64 // 外気側出口エンタルピー [J/kg]
	Xeinold  float64
	Xeoutold float64
	Xoinold  float64
	Xooutold float64
	Qes      float64 // 交換顕熱 [W]
	Qel      float64 // 交換潜熱 [W]
	Qet      float64 // 交換全熱 [W]
	Qos      float64 // 交換顕熱 [W]
	Qol      float64 // 交換潜熱 [W]
	Qot      float64 // 交換全熱 [W]
	Teidy    SVDAY   // 還気側入口温度日集計
	Teody    SVDAY   // 還気側出口温度日集計
	Xeidy    SVDAY
	Xeody    SVDAY
	Toidy    SVDAY
	Toody    SVDAY
	Xoidy    SVDAY
	Xoody    SVDAY
	Heidy    SVDAY
	Heody    SVDAY
	Hoidy    SVDAY
	Hoody    SVDAY
	Qdyes    QDAY
	Qdyel    QDAY
	Qdyet    QDAY
	Qdyos    QDAY
	Qdyol    QDAY
	Qdyot    QDAY
	MTeidy   SVDAY // 還気側入口温度日集計
	MTeody   SVDAY // 還気側出口温度日集計
	MXeidy   SVDAY
	MXeody   SVDAY
	MToidy   SVDAY
	MToody   SVDAY
	MXoidy   SVDAY
	MXoody   SVDAY
	MHeidy   SVDAY
	MHeody   SVDAY
	MHoidy   SVDAY
	MHoody   SVDAY
	MQdyes   QDAY
	MQdyel   QDAY
	MQdyet   QDAY
	MQdyos   QDAY
	MQdyol   QDAY
	MQdyot   QDAY
}

type VAVCA struct {
	Name  string  // カタログ名
	Type  VAVType // A:VAV  W:VWV
	Gmax  float64 // 最大風量 [kg/s]
	Gmin  float64 // 最小風量 [kg/s]
	dTset float64 // VWV用設定温度差　[℃]
}

type OMVAVCA struct {
	Name string
	Gmax float64 // 最大風量[kg/s]
	Gmin float64 // 最小風量[kg/s]
}

type STHEATCA struct {
	Name    string  // 機器名
	Q       float64 // 電気ヒーター容量 [W]
	Hcap    float64 // 熱容量 [J/K]
	KA      float64 // 熱損失係数 [W/K]
	Eff     float64 // 温風吹出温度効率 [-]
	PCMName string  // 電気蓄熱暖房器内臓PCMのスペック名称
}

type STHEAT struct {
	Name    string
	Cat     *STHEATCA
	Cmp     *COMPNT
	Pcm     *PCM     // 電気蓄熱暖房器内臓PCMのスペック構造体
	CG      float64  /* 熱容量流量 [W/K] */
	Ts      float64  /* 蓄熱体温度 [℃] */
	Tsold   float64  /* 前時間砕石温度 [℃] */
	Tin     float64  /* 入口（吸込）温度 [℃] */
	Tout    float64  /* 出口（吹出）温度 [℃] */
	Tenv    *float64 /* 周囲温度 [℃] */
	Xin     float64  /* 入口絶対湿度 [kg/kg'] */
	Xout    float64  /* 出口絶対湿度 [kg/kg'] */
	Q       float64  /* 供給熱量 [W] */
	E       float64  /* 電気ヒーター消費電力 [W] */
	Qls     float64  /* 熱損失 [W] */
	Qsto    float64  /* 蓄熱量 [W] */
	Qlossdy float64  /* 日積算熱損失 [kWh] */
	Qstody  float64  /* 日積算蓄熱量 [kWh] */
	MPCM    float64  // 電気蓄熱暖房器内臓PCMの容量[m3]
	Hcap    float64  // 熱容量（PCM潜熱も含む）
	Room    *ROOM    /* 蓄熱暖房器設置室構造体 */
	Tidy    SVDAY    /* 入口温度日集計 */
	Tody    SVDAY    /* 出口温度日集計 */
	Tsdy    SVDAY    /* 蓄熱体温度日集計 */
	Qdy     QDAY     /* 室供給熱量日集計 */
	Edy     EDAY
	//mtEdy [12][24]EDAY
	MTidy    SVDAY /* 入口温度日集計 */
	MTody    SVDAY /* 出口温度日集計 */
	MTsdy    SVDAY /* 蓄熱体温度日集計 */
	MQdy     QDAY  /* 室供給熱量日集計 */
	MEdy     EDAY
	MQlossdy float64      /* 日積算熱損失 [kWh] */
	MQstody  float64      /* 日積算蓄熱量 [kWh] */
	MtEdy    [12][24]EDAY // 月別時刻別消費電力[kWh]
}

/*---- Satoh Debug VAV  2000/10/30 ----*/
type VAV struct {
	Chmode rune   /* 冷房用、暖房用の設定 */
	Name   string /* 機器名 */
	Mon    rune   /* 制御対象が
	　　コイルの時：c
	　　仮想空調機の時：h
	　　床暖房の時：f
	**************************/
	Cat   *VAVCA  /* VAVカタログ構造体 */
	Hcc   *HCC    /* VWVの時の制御対象コイル */
	Hcld  *HCLOAD /* VWVの時の制御対象仮想空調機 */
	Rdpnl *RDPNL  /* VWVの時の制御対象放射パネル */
	//room []ROOM			/* 制御室構造体 */
	G         float64 /* 風量 [kg/s] */
	CG        float64 /* 熱容量流量 [W/K] */
	Q         float64 /* 再熱計算時の熱量 [W] */
	Qrld      float64
	Tin, Tout float64 /* 入口、出口空気温度 */
	Count     int     /* 計算回数 */
	Cmp       *COMPNT
}

// Satoh OMVAV 2010/12/16
type OMVAV struct {
	Name   string
	Cat    *OMVAVCA
	Omwall *RMSRF // 制御対象とする集熱屋根
	Cmp    *COMPNT
	Plist  *PLIST // 接続している末端経路への参照
	G      float64
	Rdpnl  [4]*RDPNL
	Nrdpnl int
}

// 冷温水コイル機器仕様
type HCCCA struct {
	name string
	et   float64 // 定格温度効率 [-]
	KA   float64
	eh   float64 // 定格エンタルピ効率 [-]
}

// システム使用冷温水コイル
type HCC struct {
	Name                   string
	Wet                    rune   // w:湿りコイル, d:乾きコイル
	Etype                  rune   // 温度効率の入力方法 e:et (定格(温度効率固定タイプ)) k:KA (変動タイプ)
	Cat                    *HCCCA // 冷温水コイル機器仕様
	Cmp                    *COMPNT
	et                     float64 // 温度効率 [-]
	eh                     float64 // エンタルピ効率 [-]
	Et                     ACS     // 処理熱量(温度?)
	Ex                     ACS     // 処理熱量(湿度?)
	Ew                     ACS     // 処理熱量(水?)
	cGa                    float64
	Ga                     float64
	cGw                    float64
	Gw                     float64
	Tain                   float64 // IN空気温度??
	Taout                  float64 // OUT空気温度??
	Xain                   float64 // IN空気湿度??
	Twin                   float64 // IN水湿度??
	Twout                  float64 // OUT水湿度??
	Qs                     float64
	Ql                     float64
	Qt                     float64
	Taidy, xaidy, Twidy    SVDAY
	Qdys, Qdyl, Qdyt       QDAY
	mTaidy, mxaidy, mTwidy SVDAY
	mQdys, mQdyl, mQdyt    QDAY
}

type BOICA struct /*ボイラ－機器仕様*/
{
	name     string        /*名称　　　　　　　　　　*/
	ene      rune          /*使用燃料　G:ガス、O:灯油、E:電気*/
	unlimcap rune          /*エネルギー計算で機器容量上限無いとき 'y' */
	belowmin ControlSWType /* 最小出力以下の時にOFFかONかを指示 */
	/*      ON : ON_SW    OFF : OFF_SW   */
	/* ただし、Qmin > 0 の時のみ有効 */
	plf rune /*部分負荷特性コ－ド　　　*/
	//mode rune		// 温熱源の時は 'H'、冷熱源の時は 'C'
	Qostr string   // 定格能力条件
	Qo    *float64 /*定格加熱能力　　　　　　*/
	// Qo<0 の場合は冷水チラー
	Qmin float64
	eff  float64 /*ボイラ－効率　　　　　　*/
	Ph   float64 /*温水ポンプ動力 [W] */
}

// システム使用ボイラ－
type BOI struct {
	Name string
	Mode rune /* 負荷制御以外での運転モード
	最大能力：M
	最小能力：m        */
	HCmode      rune // 冷房モート゛、暖房モード
	Load        *rune
	Cat         *BOICA
	Cmp         *COMPNT
	Do, D1      float64
	cG          float64
	Tin         float64
	Toset       float64
	Q, E, Ph    float64
	Tidy        SVDAY
	Qdy         QDAY
	Edy, Phdy   EDAY
	MtEdy       [12][24]EDAY
	MtPhdy      [12][24]EDAY
	mTidy       SVDAY
	mQdy        QDAY
	mEdy, mPhdy EDAY
}

type RFCMP struct /*標準圧縮機特性　*/
{
	name  string     /*名称　　　　　　　　　　*/
	cname string     /*圧縮機タイプ説明　　　　*/
	e     [4]float64 /*蒸発器係数　　　　　　　*/
	d     [4]float64 /*凝縮器係数　　　　　　　*/
	w     [4]float64 /*軸動力係数　　　　　　　*/
	Teo   [2]float64 /*蒸発温度範囲　　　　　　*/
	Tco   [2]float64 /*凝縮温度範囲　　　　　　*/
	Meff  float64    /*モ－タ－効率　　　　　　*/
}

type HPCH struct /* ヒートポンプ定格能力　*/
{
	Qo  float64 /*定格冷却能力（加熱能力）*/
	Go  float64 /*定格冷（温）水量、風量　　　*/
	Two float64 /*定格冷（温）水出口温度（チラ－）*/
	eo  float64 /*定格水冷却（加熱）器、空調機コイル温度効率*/

	Qex float64 /*定格排出（採取）熱量　　*/
	Gex float64 /*定格冷却風量、水量　　　　　*/
	Tex float64 /*定格外気温（冷却水入口水温）*/
	eex float64 /*定格凝縮器（蒸発器）温度効率*/

	Wo float64 /*定格軸動力　　　　　　　*/
}

type REFACA struct /*ヒートポンプ（圧縮式冷凍機）機器仕様*/
{
	name  string /*名称　　　　　　　　　　*/
	awtyp rune   /*空冷（空気熱源）=a、冷却塔使用=w */

	plf      rune    /*部分負荷特性コ－ド　　　*/
	unlimcap rune    /*エネルギー計算で機器容量上限無いとき 'y' */
	mode     [2]rune /*冷房運転: C、暖房運転時: H */
	Nmode    int     /*mode[]の数。1 のとき冷房専用または暖房専用　*/
	/*            2 のとき冷・暖 切換運転　*/
	rfc *RFCMP
	Ph  float64 /*定格冷温水ポンプ動力 [W] */

	cool *HPCH /* 冷房運転時定格能力　*/
	heat *HPCH /* 暖房運転時定格能力　*/
}

// システム使用ヒートポンプ
type REFA struct {
	Name   string /*名称　　　　　　　　　　*/
	Load   *rune
	Chmode rune /*冷房運転: C、暖房運転時: H */
	Cat    *REFACA
	Cmp    *COMPNT
	Room   *ROOM
	c_e    [4]float64 /*冷房運転時蒸発器係数　　*/
	c_d    [4]float64 /*冷房運転時熱源側（凝縮器係数）*/
	c_w    [4]float64 /*冷房運転時軸動力係数　　*/

	h_e [4]float64 /*暖房運転時凝縮器係数　　*/
	h_d [4]float64 /*暖房運転時熱源側（蒸発器係数）*/
	h_w [4]float64 /*暖房運転時軸動力係数　　*/

	Ho, He float64  /*運転時能力特性式係数　　*/
	Ta     *float64 /*外気温度 */
	Do, D1 float64
	cG     float64
	Te     float64 /*運転時蒸発温度　　　　　*/
	Tc     float64 /*運転時凝縮温度　　　　　*/

	Tin         float64
	Toset       float64
	Q           float64
	Qmax        float64
	E           float64
	Ph          float64 /*冷温水ポンプ動力 [W] */
	Tidy        SVDAY
	Qdy         QDAY
	Edy, Phdy   EDAY
	mtEdy       [12][24]EDAY
	mtPhdy      [12][24]EDAY
	mTidy       SVDAY
	mQdy        QDAY
	mEdy, mPhdy EDAY
}

type COLLCA struct /*太陽熱集熱器機器仕様*/
{
	name string
	Type rune // 水熱源：w、空気熱源：a

	b0, b1 float64
	Fd     float64 // 集熱器効率係数（=Kc / Ko）
	Ko     float64 // 総合熱損失係数[W/(m2･K)]
	Ac     float64
	Ag     float64
}

// システム使用太陽熱集熱器
type COLL struct {
	Name string

	Cat    *COLLCA
	sol    *EXSF
	Cmp    *COMPNT
	Ta     *float64
	Do, D1 float64
	ec     float64
	Te     float64 // 相当外気温度
	Tcb    float64 // 集熱板温度
	//Ko float64					// 総合熱損失係数[W/(m2･K)]
	//Fd float64					// 集熱器効率係数（=Kc / Ko）
	Tin    float64 // 入口温度
	Q      float64 // 集熱量[W]
	Ac     float64 // 集熱器面積
	Sol    float64 // 集熱面日射量[W]（短波のみ）
	Tidy   SVDAY
	Qdy    QDAY
	Soldy  EDAY
	mTidy  SVDAY
	mQdy   QDAY
	mSoldy EDAY
}

type PIPECA struct /*配管・ダクト仕様*/
{
	name string
	Type rune /*配管のとき P、ダクトのときD */
	Ko   float64
}

// システム使用配管・ダクト
type PIPE struct {
	Name  string
	Loadt *rune
	Loadx *rune
	//Type rune
	Cat    *PIPECA
	Cmp    *COMPNT
	Room   *ROOM
	L      float64
	Ko     float64
	Tenv   *float64
	Ep     float64
	Do, D1 float64
	Tin    float64
	Q      float64
	Tout   float64
	Hout   float64
	Xout   float64
	RHout  float64

	Toset float64
	Xoset float64
	Tidy  SVDAY
	Qdy   QDAY
	MTidy SVDAY
	MQdy  QDAY
}

type STANKCA struct /* 蓄熱槽機器仕様 */
{
	name   string
	Type   rune   /* 形状  既定値 'C': 縦型 */
	tparm  string /* 槽分割、流入口、流出口入力データ */
	Vol    float64
	KAside float64
	KAtop  float64
	KAbtm  float64
	gxr    float64
}

// システム使用蓄熱槽
type STANK struct {
	Name      string
	Batchop   rune /* バッチ操作有　給水:'F'  排出:'D'  停止:'-'  バッチ操作無:'n' */
	Cat       *STANKCA
	Cmp       *COMPNT
	Ndiv      int /* 分割層数 */
	Nin       int /* 流入口、流出口数 */
	Jin       []int
	Jout      []int
	Jva       int
	Jvb       int
	Ncalcihex int // 内径と長さから計算される内蔵熱交のモデルの数
	Pthcon    []ELIOType
	Batchcon  []rune /* バッチ給水、排出スケジュール　'F':給水 'D':排出 */
	Ihex      []rune /* 内蔵熱交換器のある経路のとき ihex[i]='y' */
	Cfcalc    rune   /* cfcalc = 'y':要素モデル係数の計算する。
							'n':要素モデル係数の計算しない。
	(温度分布の逆転時再計算指定のときに使用*/
	B   []float64
	R   []float64
	D   []float64
	Fg  []float64 /* Fg 要素機器の係数 [Ndiv x Nin] */
	Tss []float64

	DtankF []rune /* 分割した槽内の状態　'F':満水　'E':空 */

	// 内蔵熱交換器の温度効率が入力されていたら'N'
	// KAが入力されていたら'Y'
	// 内径と長さが入力されていたら'C'
	KAinput []rune

	Dbld0  float64 // 内蔵熱交の内径[m]
	DblL   float64 // 内蔵熱交の長さ[m]
	DblTw  float64 // 熱伝達率計算用の配管内温度[℃]
	DblTa  float64 // 熱伝達率計算用タンク温度[℃]
	Tssold []float64
	Dvol   []float64
	Mdt    []float64
	KS     []float64

	KA     []float64 // 内蔵熱交換器のKA[W/K]
	Ihxeff []float64 /* 内蔵熱交換器の熱交換器有効率　サイズは[Nin] */
	CGwin  []float64 /* cGwin, *EGwin, Twin, Q のサイズは[Nin] */
	EGwin  []float64 /* EGwin = eff * cGwin  */
	Twin   []float64
	Q      []float64

	Qloss float64 /* 槽熱損失　*/
	Qsto  float64 /*  槽蓄熱量 */

	Tenv     *float64 /* 周囲温度のアドレス */
	Stkdy    []STKDAY
	Mstkdy   []STKDAY
	Qlossdy  float64
	Qstody   float64
	MQlossdy float64
	MQstody  float64
}

type STKDAY struct {
	Tidy, Tsdy SVDAY
	Qdy        QDAY
}

type HEXCA struct /* 熱交換器機器仕様 */
{
	Name string
	eff  float64 /* 熱交換器有効率 */
	KA   float64
}

// システム使用熱交換器
type HEX struct {
	Id    int
	Name  string
	Etype rune /* 温度効率の入力方法
	　　e:et
		k:KA	*/
	Cat            *HEXCA
	Cmp            *COMPNT
	Eff            float64
	ECGmin         float64
	CGc, CGh       float64
	Tcin           float64 // 流入温度?
	Thin           float64 // 流入温度?
	Qci, Qhi       float64 // 交換熱量
	Tcidy, Thidy   SVDAY
	Qcidy, Qhidy   QDAY
	MTcidy, MThidy SVDAY
	MQcidy, MQhidy QDAY
}

type PFCMP struct /* ポンプ・ファンの部分負荷特性の近似式係数 */
{
	pftype   rune   /* 'P' ポンプ  'F' ファン */
	Type     string /* ポンプ・ファンのタイプ */
	dblcoeff [5]float64
}

type PUMPCA struct /* ポンプ・ファン機器仕様 */
{
	name   string
	pftype rune   /* 'P' ポンプ  'F' ファン */
	Type   string /* 'C' 定流量　　'P' 太陽電池駆動　*/

	Wo    float64   /* モーター入力 */
	Go    float64   /* 定格流量 */
	qef   float64   /* 発熱比率（流体加熱量= gef * Wo）*/
	val   []float64 /* 特性式係数など */
	pfcmp *PFCMP
}

// システム使用ポンプ・ファン
type PUMP struct {
	Name string
	Cat  *PUMPCA
	Cmp  *COMPNT
	//pfcmp *PFCMP
	Sol              *EXSF
	Q                float64
	G                float64
	CG               float64
	Tin              float64
	E                float64
	PLC              float64 // 部分負荷特性を考慮した入力率
	Qdy, Gdy, Edy    EDAY
	MtEdy            [12][24]EDAY
	MQdy, MGdy, MEdy EDAY
}

//  境界条件設定用仮想機器
type FLIN struct {
	Name   string
	Namet  string   /* 変数名（温度、顕熱） */
	Namex  string   /* 変数名（湿度、潜熱） */
	Awtype rune     /* 'W':１変数のとき（nametの変数名のみ使用）、 'A':２変数のとき（namexの変数も使用） */
	Vart   *float64 /* nametで示された変数の値 */
	Varx   *float64 /* namexで示された変数の値 */

	Cmp *COMPNT
}

type HCLoadType rune

const (
	HCLoadType_D HCLoadType = 'D' // 直膨コイル想定
	HCLoadType_W HCLoadType = 'W' // 冷温水コイル想定
)

// 空調機負荷仮想機器
type HCLOAD struct {
	Name    string
	Loadt   *rune
	Loadx   *rune
	RMACFlg rune // ルームエアコンなら'Y'
	Chmode  rune // スケジュール等によって設定されている運転モード
	//		opmode rune			// 実際の運転時のモード
	Type    HCLoadType /* 'D':直膨コイル想定 'W':冷温水コイル想定　*/
	Wetmode rune       /* 実際のコイル状態 */
	Wet     rune       /*'y': wet coil  'n':dry coil */

	CGa   float64
	Ga    float64
	Tain  float64
	Xain  float64
	Toset float64
	Xoset float64

	/*---- Roh Debug for a constant outlet humidity model of wet coil  2003/4/25 ----*/
	RHout float64

	CGw   float64
	Gw    float64
	Twin  float64
	Twout float64

	Qfusoku float64
	Ele     float64
	COP     float64

	Qs                                 float64
	Ql                                 float64
	Qt                                 float64
	Qcmax, Qhmax, Qc, Qh, Qcmin, Qhmin float64
	COPc, COPh                         float64 // COP（定格）
	Ec, Eh, Ecmax, Ecmin               float64 // 消費電力[W]
	COPcmax, COPcmin                   float64 // COP（最大能力時、最小能力時
	Gi, Go                             float64 // 室内機、室外機風量[kg/s]
	COPhmax, COPhmin, Ehmin, Ehmax     float64
	Rc, Rh                             [3]float64 // 理論COPと実働COPの比の2次式回帰係数
	Pcc, Pch                           float64    // ファン等消費電力[W]
	BFi, BFo                           float64    // 室内機、室外機のバイパスファクタ
	rh, rc                             float64    // 定格能力と最大能力の比
	Taidy, xaidy                       SVDAY
	Qdys, Qdyl, Qdyt                   QDAY
	Qdyfusoku, Edy                     QDAY
	mtEdy                              [12][24]EDAY
	mTaidy, mxaidy                     SVDAY
	mQdys, mQdyl, mQdyt                QDAY
	mQdyfusoku, mEdy                   QDAY

	Cmp *COMPNT
}

// // 入力負荷仮想機器
// type GLOAD struct {
// 	name   string
// 	nameqs string
// 	nameql string
// 	nameQt string
// 	Qs     []float64
// 	Ql     []float64
// 	Qt     []float64

// 	cmp *COMPNT
// }

// 太陽電池のカタログデータ
type PVCA struct {
	Name        string  // 名称
	PVcap       float64 // 太陽電池容量[W]
	Area        float64 // アレイ面積[m2]
	KHD         float64 // 日射量年変動補正係数[-]
	KPD         float64 // 経時変化補正係数[-]
	KPM         float64 // アレイ負荷整合補正係数[-]
	KPA         float64 // アレイ回路補正係数[-]
	effINO      float64 // インバータ実行効率[-]
	apmax       float64 // 最大出力温度係数[-]
	ap          float64 // 太陽電池裏面の熱伝達率[W/m2K]
	Type        rune    // 結晶系：'C'  アモルファス系：'A'
	A, B        float64 // 設置方式別の太陽電池アレイ温度計算係数
	InstallType rune    // 太陽電池パネル設置方法 'A':架台設置形、'B':屋根置き形、'C':屋根材形（裏面通風構造があるタイプ）
}

// 太陽電池
type PV struct {
	Name     string //名称
	Cmp      *COMPNT
	Cat      *PVCA    // カタログデータ
	KTotal   float64  // 太陽電池の総合設計係数[-]
	KConst   float64  // 温度補正係数以外の補正係数の積（温度補正係数以外は時々刻々変化しない）
	KPT      float64  // 温度補正係数[-]
	TPV      float64  // 太陽電池温度[℃]
	Power    float64  // 発電量[W]
	Eff      float64  // 発電効率[-]
	Iarea    float64  // 太陽電池入射日射量[W]
	PVcap    float64  // 太陽電池設置容量[W]
	Area     float64  // アレイ面積[m2]
	Ta, V, I *float64 // 外気温、風速、日射量[W/m2]
	Sol      *EXSF    // 設置方位
	Edy      QDAY     // 日積算発電量[kWh]
	Soldy    EDAY
	mEdy     QDAY // 日積算発電量[kWh]
	mtEdy    [12][24]EDAY
	mSoldy   EDAY
}

// カタログデータ（機器仕様データ一覧）
type EQCAT struct {
	Hccca    []HCCCA    // <カタログ>冷温水コイル
	Boica    []BOICA    // <カタログ>ボイラー
	Refaca   []REFACA   // <カタログ>冷温水方式の圧縮式電動ヒートポンプ,仮想熱源
	Rfcmp    []RFCMP    // <カタログ>標準圧縮機特性 (for REFACA)
	Pfcmp    []PFCMP    // <カタログ>ポンプ・ファンの部分負荷特性の近似式係数  (for REFACA)
	Collca   []COLLCA   // <カタログ>架台設置型太陽熱集熱器
	Pipeca   []PIPECA   // <カタログ>配管
	Stankca  []STANKCA  // <カタログ>蓄熱槽(熱交換型内蔵型含む)
	Hexca    []HEXCA    // <カタログ>熱交換器
	Pumpca   []PUMPCA   // <カタログ>ポンプ
	Vavca    []VAVCA    // <カタログ>VAVユニット
	Stheatca []STHEATCA // <カタログ>電気蓄熱式暖房器
	Thexca   []THEXCA   // <カタログ>全熱交換器
	PVca     []PVCA     // <カタログ>架台設置型太陽電池
	OMvavca  []OMVAVCA  // <カタログ>OMVAV
	Desica   []DESICA   // <カタログ>デシカント槽
	Evacca   []EVACCA   // <カタログ>気化冷却器
}

// 「実際に」システムを構成する機器(システム使用機器データ一覧)
type EQSYS struct {
	Cnvrg []*COMPNT // 機器

	Hcc    []HCC    // システム使用冷温水コイル
	Boi    []BOI    // システム使用ボイラ－
	Refa   []REFA   // システム使用ヒートポンプ
	Coll   []COLL   // システム使用太陽熱集熱器
	Pipe   []PIPE   // システム使用配管・ダクト
	Stank  []STANK  // システム使用蓄熱槽
	Hex    []HEX    // システム使用熱交換器
	Pump   []PUMP   // システム使用ポンプ・ファン
	Flin   []FLIN   // 境界条件設定用仮想機器
	Hcload []HCLOAD // 空調機負荷仮想機器
	Vav    []VAV    // VAVユニット
	Stheat []STHEAT // 電気蓄熱式暖房器
	Thex   []THEX   // 全熱交換器
	Valv   []VALV   // VAV
	Qmeas  []QMEAS  // カロリーメータ
	PVcmp  []PV     // 太陽電池
	OMvav  []OMVAV  // OMVAV
	Desi   []DESI   // デシカント槽
	Evac   []EVAC   // 気化冷却器

	// 使用されていなかった:
	// Ngload int
	// Gload  []GLOAD // 入力負荷仮想機器
}

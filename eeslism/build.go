package eeslism

import (
	"io"
)

// 壁体の材料定義
// wbmlist.efl から読み取った材料定義リストの要素
type BMLST struct {
	Mcode string  // 材料名
	Cond  float64 // 熱伝導率 [W/mK]
	Cro   float64 // 容積比熱 [kJ/m3K]
}

type PVWALLCAT struct {
	// PVcap float64 // 太陽電池容量[W]
	KHD     float64 // 日射量年変動補正係数(安全率) [-]
	KPD     float64 // 経時変化補正係数[-]
	KPM     float64 // アレイ負荷整合補正係数[-]
	KPA     float64 // アレイ回路補正係数[-]
	KConst  float64 // 温度補正係数以外の補正係数の積（温度補正係数以外は時々刻々変化しない）
	EffINO  float64 // インバータ実行効率[-]
	Apmax   float64 // 最大出力温度係数[-]
	Ap      float64 // 太陽電池裏面の熱伝達率[W/m2K]
	Type    byte    // 結晶系：'C' アモルファス系：'A'
	Rcoloff float64 // 集熱ファン停止時の太陽電池から集熱器裏面までの熱抵抗[m2K/W]
	Kcoloff float64
}

type PVWALL struct {
	KTotal float64 // 太陽電池の総合設計係数[-]
	KPT    float64 // 温度補正係数[-]
	TPV    float64 // 太陽電池温度[℃]
	Power  float64 // 発電量[W]
	Eff    float64 // 発電効率[-]
	PVcap  float64 // 太陽電池設置容量[W]
}

type PCMSTATE struct {
	Name         *string // name
	TempPCMNodeL float64 // PCM温度（左側節点）
	TempPCMNodeR float64 // 同（右）
	TempPCMave   float64 // PCM温度（平均温度）
	//capm       float64 // PCM比熱[J/kgK]
	//lamda      float64 // PCM熱伝導率[W/mK]
	CapmL     float64 // PCM見かけの比熱（左側）[J/kgK]
	CapmR     float64 // PCM見かけの比熱（右側）[J/kgK]
	LamdaL    float64 // PCM熱伝導率（左側）[W/mK]
	LamdaR    float64 // PCM熱伝導率（右側）[W/mK]
	OldCapmL  float64 // 前時刻PCM見かけの比熱（左側）
	OldCapmR  float64 // 前時刻PCM見かけの比熱（右側）
	OldLamdaL float64 // 前時刻PCM熱伝導率（左側）
	OldLamdaR float64 // 前時刻PCM熱伝導率（右側）
}

type RMSRFType rune

const (
	RMSRFType_None RMSRFType = 0
	RMSRFType_H    RMSRFType = 'H' // 壁
	RMSRFType_E    RMSRFType = 'E' // 地中壁
	RMSRFType_W    RMSRFType = 'W' // 窓
	RMSRFType_e    RMSRFType = 'e' // 地表面境界
)

type RMSRFMwType rune

const (
	RMSRFMwType_None RMSRFMwType = 0
	RMSRFMwType_I    RMSRFMwType = 'I' // 専用壁
	RMSRFMwType_C    RMSRFMwType = 'C' // 共同壁
)

type RMSRFMwSideType rune

const (
	RMSRFMwSideType_None RMSRFMwSideType = 0
	RMSRFMwSideType_i    RMSRFMwSideType = 'i' // 壁体0側
	RMSRFMwSideType_M    RMSRFMwSideType = 'M' // 壁体M側
)

// 壁体　固定デ－タ
type RMSRF struct {
	Name  string // 壁体名
	Sname string //RMP名 higuchi 070918

	// ---- 出力指定 ---- //

	wlpri  bool // 壁体内部温度出力指定 (ROOM *p)
	shdpri bool // 日よけの影面積の出力指定 (ROOM *shd)
	sfepri bool // 要素別壁体表面温度出力指定 (ROOM *sfe)

	// 部位コ－ド
	// 'B' | 'W' | 'E'(外壁) | 'R'(屋根)  | 'F'床（外気、地中に接する） |
	// 'i'(内壁) | 'f'(床（隣室に接する）) | 'c'(天井)  |'d'
	ble      BLEType
	typ      RMSRFType       // 壁のとき'H', 地中壁のとき'E', 窓のとき'W', 地表面境界の時'e'
	mwtype   RMSRFMwType     // 専用壁のとき'I',共用壁のとき'C'
	mwside   RMSRFMwSideType // 壁体 0 側のとき'i', M 側のとき'M'
	mrk      rune            // '*' or '!'
	ffix_flg rune            // 表面への短波長放射基本吸収比率が定義されている場合'*' 未定義の場合'!'
	PCMflg   bool            // PCMの有無の判定フラグ　毎時係数行列を計算するかどうかを決める
	pcmpri   bool            // PCMの状態値出力フラグ
	fnmrk    [10]rune        // 窓変更設定用窓コ－ド

	room     *ROOM
	nextroom *ROOM
	nxsd     *RMSRF

	mw         *MWALL  // 重量壁（差分による非定常計算）
	rpnl       *RDPNL  // 輻射パネル用
	window     *WINDOW // 窓
	dynamicwin int     // 動的に窓を切り替える場合'Y'
	ifwin      *WINDOW // Trueの時の窓
	//falsewin *WINDOW	// Falseの時の窓
	Ctlif       *CTLIF // 動的窓の制御
	DynamicCode string // 動的窓 ex) "A > B" のような評価式

	rm       int    // 室番号
	n        int    // 室壁体番号
	exs      int    // 方位番号
	nxrmname string //隣室名
	nxrm     int    // 隣室番号
	nxn      int    // 隣室室壁体番号
	wd       int    // 壁体定義番号
	rmw      int    // 重量壁体番号
	lwd      int
	dr       int // ドア定義番号
	//		drd [4]int
	fn               int     // 選択窓定義番号
	Nfn              int     // 窓種類数
	fnd              [10]int // 窓定義番号
	direct_heat_gain [10]int // 日射熱取得、窓部材熱抵抗を直接指定する場合の番号
	fnsw             int     // 窓変更設定番号

	sb int // 日除け定義番号

	A    float64 // 面積 [m2]
	Eo   float64 // 外表面輻射率 [-]
	as   float64 // 外表面日射吸収率 [-]
	c    float64 // 隣室温度係数 [-]
	tgtn float64 // 日射総合透過率 [-]
	Bn   float64 // 吸収日射取得率 [-]

	// ---- 壁体、窓熱抵抗 ----

	fsol *float64 // 部位室内表面の日射吸収比率 [-]

	/*窓透過日射の吸収比率     */
	srg     float64 // 1次入射比率（隣接室への透過やガラスから屋外への放熱は無視）
	srg2    float64 // 最終的に室の日射熱取得の内の吸収比率
	srh     float64 // 人体よりの輻射の吸収比率
	srl     float64 // 照明よりの輻射の吸収比率
	sra     float64 // 機器よりの輻射の吸収比率
	alo     float64 // 外表面熱伝達率 [W/m2K]
	ali     float64 // 内表面熱伝達率 [W/m2K]
	alic    float64
	alir    float64
	Qc      float64 // 対流による熱取得 [W]
	Qr      float64 // 放射による熱取得 [W]
	Qi      float64 // 壁体貫流熱取得 [W]
	Qgt     float64 // 透過日射熱取得 [W]
	Qga     float64 // 吸収日射熱取得 [W]
	Qrn     float64 // 夜間放射熱取得 [W]
	K       float64 // 熱伝達率 [W/m2K] K = 1/(1/alo + Rwall + 1/ali)
	Rwall   float64 // 熱抵抗 [m2K/W] 表面熱伝達抵抗(1/alo+1/ali)は除く
	CAPwall float64 // 単位面積当たり熱容量[J/m2K]
	alicsch *float64
	alirsch *float64
	FI      float64
	FO      float64
	FP      float64
	CF      float64
	WSR     float64
	WSRN    []float64
	WSPL    []float64
	WSC     float64
	Fsdworg float64 // SUNBRKで定義した日よけの影面積率
	Fsdw    float64 /* 影面積  higuchi 070918 */ // KAGE-SUN用
	Ihor    float64 /*higuchi 070918 */
	Idre    float64 /*higuchi 070918 */
	Idf     float64 /*higuchi 070918 */
	Iw      float64 /*higuchi 070918 */
	rn      float64 /*higuchi 070918 */

	// ---- 室内表面に吸収される日射と内部発熱 ----

	RS     float64 // 室内表面に吸収される短波長輻射 [W/m2]
	RSsol  float64 // 室内表面に吸収される日射（短波長）[W/m2]
	RSsold float64 // 室内表面に入射する日射（短波長）（隣接室への透過を考慮する前）
	RSli   float64 // 室内表面に吸収される照明（短波長）[W/m2]
	RSin   float64 // 室内表面に吸収される人体・照明・設備・設備機器（短波長） [W/m2]

	TeEsol float64
	TeErn  float64
	Te     float64 // 外表面の相当外気温
	Tmrt   float64 // 室内表面の平均輻射温度
	Ei     float64
	Ts     float64 // 室内表面の温度 ????

	/*設備機器発熱*/
	eqrd float64
	end  int
	//ColTe  float64		// 建材一体型空気集熱器の相当外気温度[℃]
	Tcole          float64 // 建材一体型空気集熱器の相当外気温度[℃] 記号変更
	Tcoleu, Tcoled float64
	Tf             float64 // 建材一体型空気集熱器の熱媒平均温度[℃]
	SQi            QDAY    // 日積算壁体貫流熱取得
	Tsdy           SVDAY
	PVwall         PVWALL // 太陽電池一体型壁体
	mSQi           QDAY
	mTsdy          SVDAY

	// ---- 集熱器一体型壁体用計算結果 ----

	Ndiv     int //  空気式集熱器のときの流れ方向（入口から出口）の分割数
	ColCoeff float64
	// 建材一体型集熱器計算時の後退差分要素URM
	Tc []float64
	//Scol float64 // 放射熱取得量[W/m2]
	oldTx float64 // 前時刻の集熱器と躯体の境界温度（集熱器相当外気温度計算用）
	// 太陽電池一体型
	Iwall                                                 float64
	PVwallFlg                                             bool    // 太陽電池一体型の場合はtrue
	dblWsu                                                float64 // <入力値> 屋根一体型空気集熱器(集熱屋根)の通気層上側の幅 [m]
	dblWsd                                                float64 // <入力値> 屋根一体型空気集熱器(集熱屋根)の通気層下側の幅 [m]
	dblKsu, dblKsd, dblKc, dblfcu, dblfcd, dblKcu, dblKcd float64
	dblb11, dblb12, dblb21, dblb22                        float64
	dblacr, dblacc, dblao                                 float64
	dblTsu, dblTsd, dblTf                                 float64
	dblSG                                                 float64
	ku, kd                                                float64
	ras, Tg                                               float64

	pcmstate []*PCMSTATE // PCM状態値収録構造体
	Npcm     int         // PCM設置レイヤー数

	tnxt    float64 // 当該部位への入射日射の隣接空間への日射分配（連続空間の隣室への日射分配）
	RStrans bool    // 室内透過日射が窓室内側への入射日射を屋外に透過する場合'y'
}

// 壁体各層の熱抵抗と熱容量
type WELM struct {
	Code string  // <入力値> 材料コード
	L    float64 // <入力値> 各層の材料厚さ［m］（分割前）
	ND   int     // <入力値> 各層の内部分割数
	Cond float64 // <材料定義リストから読み取り> 熱伝導率  [W/mK]
	Cro  float64 // <材料定義リストから読み取り> 容積比熱  [J/m3K]
}

func NewWelm() *WELM {
	return &WELM{
		Code: "",
		L:    -999.0,
		ND:   0,
		Cond: -999.0,
		Cro:  -999.0,
	}
}

type CHARTABLE struct {
	filename             string        // テーブル形式ファイルのファイル名
	fp                   io.ReadCloser // `filename`の読み込みファイルポインタ
	PCMchar              rune          // E:エンタルピー、C:熱伝導率
	T                    []float64     // PCM温度[℃]
	Chara                []float64     // 特性値（エンタルピー、熱伝導率）
	tabletype            rune          // h:見かけの比熱、e:エンタルピー
	minTemp, maxTemp     float64       // テーブルの下限温度、上限温度
	itablerow            int           // テーブル形式の入力行数
	lowA, lowB, upA, upB float64       // 上下限温度範囲外の特性値計算用線形回帰式の傾きと切片
	minTempChng          float64       // 最低温度変動幅　前時刻からの温度変化がminTempChng以下の場合はminTempChngとして特性値を計算
}

// 潜熱蓄熱材
type PCM struct {
	Name         string       // PCM名称
	Spctype      rune         // 見かけの比熱　m:モデルで設定、t:テーブル形式
	Condtype     rune         // 熱伝導率　m:モデルで設定、t:テーブル形式
	Ql           float64      // 潜熱量[J/m3]
	Condl        float64      // 液相の熱伝導率[W/mK]
	Conds        float64      // 固相の熱伝導率[W/mK]
	Crol         float64      // 液相の容積比熱[J/m3K]
	Cros         float64      // 固相の容積比熱[J/m3K]
	Ts           float64      // 固体から融解が始まる温度[℃]
	Tl           float64      // 液体から凝固が始まる温度[℃]
	Tp           float64      // 見かけの比熱のピーク温度
	Iterate      bool         // PCM状態値を収束計算させるかどうか
	IterateTemp  bool         // 収束条件に温度も加えるかどうか（通常は比熱のみ）
	DivTemp      int          // 比熱の数値積分時の温度分割数
	Ctype        int          // 比熱
	PCMp         PCMPARAM     // 見かけの比熱計算用パラメータ
	AveTemp      rune         // PCM温度を両側の節点温度の平均で計算する場合は'y'（デフォルト）
	NWeight      float64      // 収束計算時の現在ステップの重み係数
	IterateJudge float64      // 収束計算時の前ステップ見かけの比熱の収束判定[%]
	Chartable    [2]CHARTABLE // 0:見かけの比熱またはエンタルピー、1:熱伝導率
}

// PCM見かけの比熱計算用パラメータ
type PCMPARAM struct {
	T     float64
	B     float64
	bs    float64
	bl    float64
	skew  float64
	omega float64
	a     float64
	b     float64
	c     float64
	d     float64
	e     float64
	f     float64
}

type BLEType rune

const (
	BLE_None         BLEType = 0
	BLE_ExternalWall BLEType = 'E' // 外壁
	BLE_Roof         BLEType = 'R' // 屋根
	BLE_Floor        BLEType = 'F' // 外部に接する床
	BLE_InnerWall    BLEType = 'i' // 内壁
	BLE_Ceil         BLEType = 'c' // 天井(内部)
	BLE_InnerFloor   BLEType = 'f' // 床(内部)
	BLE_d            BLEType = 'd'
	BLE_Window       BLEType = 'W' // 窓
)

type WALLType rune

const (
	WallType_None WALLType = 0
	WallType_C    WALLType = 'C' // 建材一体型空気集熱器
	WallType_P    WALLType = 'P' // 床暖房等放射パネル(通常の床暖房パネル)
	WallType_N    WALLType = 'N' // 一般壁体
)

// 壁体　定義デ－タ
type WALL struct {
	ble    BLEType // <入力値> 部位コ－ド = E,R,F,i,c,f,R
	name   string  // <入力値> 壁体名 最初の1文字は英字 省略時は既定値とみなす
	PCMflg bool    // 部材構成にPCMが含まれる場合は毎時係数行列を作成するので
	// PCMが含まれるかどうかのフラグ

	N  int /*材料層数≠節点数        */
	Ip int /*発熱面のある層の番号  */
	//	code [12][5]rune; /*各層の材料コ－ド      */
	L []float64 /*節点間の材料厚さ        */
	//	ND []int;      /*各層の内部分割数      */
	Ei          float64   // <入力値> 室内表面放射率
	Eo          float64   // <入力値> 外表面輻射率
	as          float64   // <入力値> 外表面日射吸収率
	Rwall       float64   // <内部計算値> 壁体熱抵抗(表面熱伝達抵抗(1/alo+1/ali)は除く) [m2K/W]
	CAPwall     float64   // <内部計算値> 単位面積当たりの熱容量[J/m2K]
	CAP, R      []float64 // <入力値> 登録された材料（≠節点）ごとの熱容量、熱抵抗
	effpnl      float64   // <入力値> 放射暖冷房パネルのパネル効率
	tnxt        float64   // <入力値> 当該部位への入射日射の隣接空間への日射分配（連続空間の隣室への日射分配）
	M           int       // 節点数
	mp          int       // <入力値> 放射暖冷房パネルの発熱面のある節点番号
	res         []float64 // 節点間の熱抵抗 [m2K/W]
	cap         []float64 // 節点間の熱容量 [J/m2K]
	end         int       // 要素数(インデックス0にのみ設定)
	welm        []WELM    // <入力値> 層構成(layer)
	tra         float64   // <入力値> τα
	Ko          float64   // <内部計算値> Ksu + Ksd
	Ksu         float64   // <入力値> 通気層上部から屋外までの熱貫流率[W/m2K]
	Ksd         float64   // <入力値> 通気層下部から集熱器裏面までの熱貫流率[W/m2K]
	fcu         float64   // <入力値> Kcu / Ksu
	fcd         float64   // <入力値> Kcd / Ksd
	ku          float64   // <内部計算値> Kcu / Kc
	kd          float64   // <内部計算値> Ksu / Ko
	Ru          float64   // <入力値> 通気層上部から屋外までの熱抵抗 [m2K/W]
	Rd          float64   // <入力値> 通気層下部から集熱器裏面までの熱抵抗 [m2K/W]
	Kc          float64   // <入力値> Kcu + Kcd
	Kcu         float64   // <入力値> 通気層内上側から屋外までの熱貫流率 [W/m2K]
	Kcd         float64   // <入力値> 通気層内下側から裏面までの熱貫流率 [W/m2K]
	air_layer_t float64
	dblEsu      float64 // 通気層の厚さ[m]
	dblEsd      float64 // 通気層の厚さ[m]
	ta          float64 // <入力値> 中空層の厚さ [mm]
	Eg          float64 // <入力値> 透過体の中空層側表面の放射率
	Eb          float64 // <入力値> 集熱版の中空層側表面の放射率
	ag          float64 // <入力値> 透過体の日射吸収率

	chrRinput bool     // 熱抵抗が入力されている場合は'Y', 熱貫流率が入力されている場合は'N'
	WallType  WALLType // <内部判定値> 建材一体型空気集熱器の場合：'C', 床暖房等放射パネルの場合：'P', 一般壁体の場合：'N'

	//char	PVwall ;
	// 太陽電池一体型建材（裏面通気）：'Y'
	ColType string // <入力値> 集熱器のタイプ: 'A1'=ガラス付集熱器, 'A2'=ガラス無し集熱器 or 'A2P'=太陽電池一体型集熱器

	// 集熱器タイプ　JSES2009大会論文（宇田川先生発表）のタイプ
	PVwallcat  PVWALLCAT
	PCM        []*PCM
	PCMLyr     []*PCM    // 潜熱蓄熱材
	PCMrate    []float64 // PCM含有率（Vol）
	PCMrateLyr []float64 // PCM体積含有率
}

func NewWall() *WALL {
	Wa := new(WALL)

	Wa.name = ""
	Wa.ble = ' '
	Wa.N = 0
	Wa.M = 0
	Wa.mp = 0
	Wa.end = 0
	Wa.Ip = -1
	Wa.Ei = 0.9
	Wa.Eo = 0.9
	Wa.as = 0.7
	Wa.Rwall = -999.0
	Wa.effpnl = -999.0
	Wa.CAPwall = -999.
	Wa.res = nil
	Wa.cap = nil
	// Wa.welm = nil ;
	Wa.tra = -999.0
	Wa.Ksu = -999.
	Wa.Ksd = -999.
	Wa.Kc = -999.
	Wa.fcu = -999.
	Wa.fcd = -999.
	Wa.ku = -999.
	Wa.kd = -999.
	Wa.Ko = -999.
	Wa.Rd = -999.
	Wa.Ru = -999.
	Wa.Kcu = -999.
	Wa.Kcd = -999.
	Wa.air_layer_t = -999.0
	Wa.dblEsd = 0.9
	Wa.dblEsu = 0.9
	Wa.chrRinput = false
	Wa.ColType = ""
	Wa.WallType = WallType_N //一般壁体
	//Wa.PVwall = 'N' ;

	// 太陽電池一体型空気集熱器のパラメータ初期化
	PVwallcatinit(&Wa.PVwallcat)

	Wa.ta = -999.
	Wa.ag = -999.0
	Wa.Eg = 0.9
	Wa.Eb = 0.9

	Wa.PCMLyr = nil
	Wa.PCMrateLyr = nil
	Wa.L = nil

	Wa.PCMflg = false

	Wa.tnxt = -999.0

	return Wa
}

// 壁体定義番号既定値
type DFWL struct {
	E int // 外壁(壁体定義番号既定値)
	R int // 屋根(壁体定義番号既定値)
	F int // 外部に接する床(壁体定義番号既定値)
	i int // 内壁(壁体定義番号既定値)
	c int // 天井(壁体定義番号既定値)
	f int // 隣室に接する床(壁体定義番号既定値)
}

// 重量壁体デ－タ
type MWALL struct {
	sd, nxsd *RMSRF
	wall     *WALL
	ns       int // 壁体通し番号
	rm       int // 室番号
	n        int // 室壁体番号
	nxrm     int // 隣室番号
	nxn      int // 隣室室壁体番号
	/* [UX]の先頭位置       */
	UX []float64

	M, mp    int
	res, cap []float64
	uo       float64   // 室内表面のuo
	um       float64   // 外表面のum
	Pc       float64   // 床パネル用係数
	Tw       []float64 // 壁体温度
	Told     []float64 // 以前の壁体温度
	Twd      []float64 // 現ステップの壁体内部温度
	Toldd    []float64 // PCM温度に関する収束計算過程における前ステップの壁体内温度
	end      int
}

// 窓およびドア定義デ－タ
type WINDOW struct {
	Name    string  // 名称
	Cidtype string  // 入射角特性のタイプ。 'N':一般ガラス
	K       float64 // !入力されてる？!
	Rwall   float64 // 窓部材熱抵抗 [m2K/W]
	Ei      float64 // 室内表面放射率(0.9) [-]
	Eo      float64 // 外表面放射率(0.9)（ドアのみ） [-]
	tgtn    float64 // 日射透過率 [-]
	Bn      float64 // 吸収日射取得率 [-]
	As      float64 // 日射吸収率（ドアのみ）[-]
	Ag      float64 // 窓ガラス面積 !入力されてる？!
	Ao      float64 // 開口面積 !入力されてる？!
	W       float64 // 巾 !入力されてる？!
	H       float64 // 高さ !入力されてる？!
	RStrans bool    // 室内透過日射が窓室内側への入射日射を屋外に透過する場合はtrue
	end     int     // 要素数(インデックス0に設定される)
}

func NewWINDOW() *WINDOW {
	W := new(WINDOW)
	W.Name = ""
	W.Cidtype = "N" // 入射角特性のタイプは一般ガラスとする
	W.K = 0.0
	W.Rwall = 0.0
	W.Ei = 0.9 // 室内表面放射率 デフォルト値
	W.Eo = 0.9 // 外表面放射率 デフォルト値
	W.tgtn = 0.0
	W.Bn = 0.0
	W.As = 0.0
	W.Ag = 0.0
	W.Ao = 0.0
	W.W = 0.0
	W.H = 0.0
	W.RStrans = false // 室内透過日射が窓室内側への入射日射を屋外に透過しない
	W.end = 0
	return W
}

// 日除け
type SNBK struct {
	Name string  // 名称
	Type int     // 日除けタイプ 1: 一般の庇(H), 2: 袖壁(HL), 5: 長い庇(S), 6: 長い袖壁(SL), 9: 格子ルーバー(K)
	Ksi  int     // 日影部分と日照部分の反転 0: 反転なし, 1: 反転あり
	W    float64 // 開口部の高さ (W=Width)
	H    float64 // 開口部の幅 (H=Height)
	D    float64 // 庇の付け根から先端までの長さ (D=Depth)
	W1   float64 // 開口部の左端から壁の左端までの距離 (L=Left)
	W2   float64 // 開口部の右端から壁の右端までの距離 (R=Right)
	H1   float64 // 開口部の上端から壁の上端までの距離 (T=Top)
	H2   float64 // 地面から開口部の下端までの高さ (B=Bottom)
	end  int     // 要素数(インデックス0に設定)
}

// 日射、室内発熱熱取得
type QRM struct {
	Tsol float64 // 透過日射 [W]
	Asol float64 // 外表面吸収日射室内熱取得 [W]
	Arn  float64 // 外表面吸収長波長輻射熱損失 [W]

	// --- 人体・照明・機器の顕熱 [W] ---

	Hums  float64 // 人体顕熱 [W]
	Light float64 // 照明 [W]
	Apls  float64 // 機器顕熱 [W]
	Hgins float64 //  室内発熱（顕熱）の合計 [W]

	// --- 人体・機器の潜熱 [W] ---

	Huml float64 // 人体潜熱 [W]
	Apll float64 // 機器潜熱 [W]

	// --- 熱負荷 ---

	Qinfs float64 // 換気顕熱負荷[W]
	Qinfl float64 // 換気潜熱負荷[W]

	Qsto  float64 // 室内の顕熱蓄熱量[W]
	Qstol float64 // 室内の潜熱蓄熱量[W]

	Qeqp float64 // 室内設置の配管、ボイラからの熱取得[W]

	Solo float64 //  外壁面入射日射量[W]
	Solw float64 //  窓面入射日射量[W]
	Asl  float64 // 外表面吸収日射[W]

	AE float64 // 消費電力[W]
	AG float64 // 消費ガス[W]
}

// 室間相互換気
type ACHIR struct {
	rm   int
	sch  int
	room *ROOM
	Gvr  float64
}

// 隣室
type TRNX struct {
	nextroom *ROOM
	sd       *RMSRF // room側からみたnextroomと接する表面の壁体への参照
}

// 室についての輻射パネル
type RPANEL struct {
	pnl     *RDPNL
	sd      *RMSRF
	elinpnl int // 放射パネルの入力要素の先頭位置
}

// 輻射パネル
type RDPNL struct {
	Name    string
	Loadt   *rune
	Type    rune // 建材一体型空気集熱器の場合：'C', それ以外：'P'
	rm      [2]*ROOM
	sd      [2]*RMSRF
	cmp     *COMPNT
	MC      int // 専用壁のとき MC=1, 共用壁のとき MC=2
	Ntrm    [2]int
	Nrp     [2]int
	elinpnl [2]int
	eprmnx  int // 隣室EPRN[]の位置
	epwtw   int // EPWの当該パネル入口水温の位置
	control rune

	effpnl float64
	Toset  float64
	Wp     float64
	Wpold  float64
	Tpi    float64
	Tpo    float64
	FIp    [2]float64
	FOp    [2]float64
	FPp    float64
	Epw    float64
	EPt    [2]float64
	EPR    [2][]float64
	EPW    [2][]float64
	EPC    float64
	Q      float64
	// 2009/01/26 Satoh Debug
	cG             float64 // 比熱×流量
	Ec, FI, FO, FP float64
	Ew             float64

	/* 日集計 */
	Tpody  SVDAY
	Tpidy  SVDAY
	Qdy    QDAY
	Scoldy QDAY
	PVdy   QDAY
	TPVdy  SVDAY

	// 月集計
	mTpody, mTpidy, mTPVdy SVDAY
	mQdy, mScoldy, mPVdy   QDAY
	mtPVdy                 [12][24]EDAY

	OMvav *OMVAV // 吹出を制御する変風量ユニット

}

// 室への冷温風供給熱量
type AIRSUP struct {
	Qs                  float64
	Ql                  float64
	Qt                  float64
	G                   float64
	Tin                 float64
	Xin                 float64
	Qdys                QDAY // 日積算暖冷房
	Qdyl                QDAY
	Qdyt                QDAY
	mQdys, mQdyl, mQdyt QDAY
}

// 室負荷
type RMLOAD struct {
	loadt *rune
	loadx *rune
	tropt rune /* 室温制御方法  'o': OT制御、'a': 空気温度制御 */
	hmopt rune /* 湿度制御設定値 'x': 絶対湿度、'r': 相対湿度、 'd': 露点温度 */
	Tset  float64
	Xset  float64
	Qs    float64
	Ql    float64
	Qt    float64

	FOTr float64
	FOTN []float64
	FOPL []float64
	FOC  float64

	Qdys                QDAY /* 日積算暖冷房 */
	Qdyl                QDAY
	Qdyt                QDAY
	mQdys, mQdyl, mQdyt QDAY
}

// ゾーン集計
type RZONE struct {
	name   string  // ゾーン名
	Nroom  int     // ゾーンに属する室の数
	rm     []*ROOM //ゾーンに属する室のポインター
	Afloor float64
	Tr     float64
	xr     float64
	RH     float64
	Tsav   float64
	Qhs    float64
	Qhl    float64
	Qht    float64
	Qcs    float64
	Qcl    float64
	Qct    float64
	Trdy   SVDAY
	xrdy   SVDAY
	RHdy   SVDAY
	Tsavdy SVDAY
	Qsdy   QDAY
	Qldy   QDAY
	Qtdy   QDAY
}

/* ---------------------------------------------------------- */
/* 要素別熱損失・熱取得計算用 */

type BHELM struct {
	trs float64 // 貫流
	so  float64 // 外壁入射日射
	sg  float64 // 窓入射日射
	rn  float64 // 大気放射
	in  float64 // 室内発熱
	pnl float64 // 放射暖・冷房パネル
}

type QHELM struct {
	qe   BHELM
	slo  float64 // 外壁面入射日射量 [W]
	slw  float64 // 窓面入射日射量 [W]
	asl  float64 // 外壁面吸収日射量 [W]
	tsol float64 // 窓透過日射量 [W]
	hins float64 // 室内発熱（顕熱） [W]
	hinl float64 // 室内発熱（潜熱） [W]

	nx     float64 // 隣室熱損失
	gd     float64 // 窓熱損失
	ew     float64 // 外壁熱損失
	wn     float64 // 窓熱損失
	i      float64 // 内壁熱損失
	c      float64 // 天井、屋根
	f      float64 // 床（内・外）
	vo     float64 // 換気
	vol    float64 // 換気（潜熱）
	vr     float64 // 室間換気
	vrl    float64 // 室間換気（潜熱）
	sto    float64 // 室内空気蓄熱
	stol   float64 // 室内空気蓄熱（潜熱）
	loadh  float64
	loadhl float64
	loadc  float64
	loadcl float64
}

type RMQELM struct {
	rmsb   []RMSB
	WSCwk  []BHELM
	qelm   QHELM
	qelmdy QHELM
}

type RMSBType rune

const (
	RMSBType_None RMSBType = 0
	RMSBType_E    RMSBType = 'E' // 外気に接する面
	RMSBType_G    RMSBType = 'G' // 地面に接する面
	RMSBType_i    RMSBType = 'i' // 内壁
)

type RMSB struct {
	Type RMSBType // 'E':外気に接する面、'G':地面に接する面、'i':内壁
	Ts   BHELM
	Tw   []BHELM
	Told []BHELM
}

type QETOTAL struct {
	Name   string
	Qelm   QHELM
	Qelmdy QHELM
}

// 室デ－タ
type ROOM struct {
	Name        string //室名
	N           int    //周壁数
	Brs         int    //Sd[],S[]の先頭位置
	MCAP        *float64
	CM          *float64
	TM          float64
	oldTM       float64
	HM          float64
	QM          float64
	fsolm       *float64 // 家具の日射吸収割合
	Srgm2       float64  // 家具の最終的な日射吸収割合
	Qsolm       float64  // 家具の日射吸収量[W]
	PCMfurnname string   // PCM内臓家具の部材名称
	PCM         *PCM     // PCM内臓家具の場合
	mPCM        float64  // PCM内臓家具の容積[m3]
	PCMQl       float64  // PCM内臓家具の見かけの比熱[J/m3K]
	FunHcap     float64  // 家具の熱容量（顕熱とPCM内臓家具の合計）
	// 室空気に加算したものは除く
	FMT, FMC float64

	// --- 家具の計算用パラメータ ----

	Qgt         float64 // 透過日射熱取得 [W]
	Nachr       int     //`achr`の数
	Ntr         int     //内壁を共有する隣室数
	Nrp         int     //輻射パネル数
	Nflr        int     //床の部位数
	Nfsolfix    int     //短波長放射吸収比率が定義されている面数
	Nisidermpnl int
	Nasup       int

	rsrf  []RMSRF  // 壁体
	achr  []ACHIR  // 室間相互換気
	trnx  []TRNX   // 隣室
	rmpnl []RPANEL // 室についての輻射パネル
	//rairflow []RAIRFLOW
	Arsp      []AIRSUP // 室への冷温風供給熱量
	cmp       *COMPNT
	elinasup  []*ELIN // 流入経路
	elinasupx []*ELIN // 流入経路？？？
	rmld      *RMLOAD // 室負荷
	rmqe      *RMQELM

	F     []float64
	alr   []float64
	XA    []float64
	Wradx []float64

	rsrnx bool //隣室裏面の短波長放射考慮のとき true　（床、天井のみ）
	fij   rune //形態係数 'F':外部入力、'A':面積率
	sfpri bool //表面温度出力指定
	eqpri bool //日射、室内発熱取得出力指定
	mrk   rune // '*', 'C', '!'

	VRM     float64  //室容積 [m3]
	GRM     float64  //室内空気質量
	MRM     float64  //室空気熱容量
	Area    float64  //室内表面総面積
	FArea   float64  //床面積
	flrsr   *float64 //床に直接吸収される短波長放射の比率
	tfsol   float64  //部位に直接吸収される短波長放射比率の既定値合計（Sd->fsol、Rm->flrsr、Rm->fsolmの合計）
	alrbold float64  //

	Hcap  float64 // 室内熱容量   J/K
	Mxcap float64 // 室内湿気容量 kg/(kg/kg)

	Ltyp     rune     // 照明器具形式
	Nhm      float64  // 人数
	Light    float64  // 照明器具容量
	Apsc     float64  // 機器対流放熱容量
	Apsr     float64  // 機器輻射放熱容量
	Apl      float64  // 機器潜熱放熱容量
	Gve      float64  // 換気量
	Gvi      float64  // 隙間風量
	AE       float64  // 消費電力容量[W]
	AG       float64  // 消費ガス容量[W]
	AEsch    *float64 // 消費電力スケジュール (未設定時はnil)
	AGsch    *float64 // 消費ガススケジュール (未設定時はnil)
	Lightsch *float64 // 照明スケジュール (未設定時はnil)
	Assch    *float64 // 機器顕熱スケジュール (未設定時はnil)
	Alsch    *float64 // 機器潜熱スケジュール (未設定時はnil)
	Hmsch    *float64 // 在室人数スケジュール (未設定時はnil)
	Metsch   *float64 // Met値スケジュール (未設定時はnil)
	Closch   *float64 // Clo値スケジュール (未設定時はnil)
	Wvsch    *float64 // 室内風速設定値名 (未設定時はnil)
	Hmwksch  *float64 // 作業強度設定値名 (未設定時はnil)
	Vesc     *float64 // 換気スケジュール (未設定時はnil)
	Visc     *float64 // 隙間風スケジュール (未設定時はnil)
	alc      *float64 // 室内側対流熱伝達率のスケジュール設定値  (未設定時はnil)

	// --- 室内発熱 ---

	Hc float64 //人体よりの対流  [W]
	Hr float64 //人体よりの輻射  [W]
	HL float64 //人体よりの潜熱　[W]
	Lc float64 //照明よりの対流  [W]
	Lr float64 //照明よりの輻射  [W]
	Ac float64 //機器よりの対流  [W]
	Ar float64 //機器よりの輻射  [W]
	AL float64 //機器よりの潜熱  [W]

	/*設備機器発熱*/
	eqcv  float64 //設備機器発熱の対流成分比率
	Qeqp  float64 //設備機器からの発熱[W]
	Gvent float64

	RMt  float64
	ARN  []float64
	RMP  []float64
	RMC  float64
	RMx  float64
	RMXC float64

	Tr     float64 //室内温度
	Trold  float64 //室内温度
	xr     float64 //室内絶対湿度 [kg/kg]
	xrold  float64 //室内絶対湿度
	RH     float64
	Tsav   float64 // 平均表面温度
	Tot    float64 // 作用温度
	hr     float64 // エンタルピー
	PMV    float64
	SET    float64 // SET(体感温度)
	setpri bool    // SET(体感温度)の出力フラグ

	Trdy   SVDAY
	xrdy   SVDAY
	RHdy   SVDAY
	Tsavdy SVDAY

	mTrdy, mxrdy, mRHdy, mTsavdy SVDAY
	VAVcontrl                    *VAV
	OTsetCwgt                    *float64 // 作用温度設定時の対流成分重み係数
	// デフォルトは0.5
	HGc, CA, AR float64
	Qsab        float64 // 吸収日射熱取得 [W]
	Qrnab       float64 // 夜間放射による熱損失 [W]
}

/* --------------------------------
室計算用データ
--------------------------------*/

type RMVLS struct {
	Twallinit  float64 // 初期温度 (GDAT.RUN.Tinit)
	Nwall      int
	Nwindow    int
	Nmwall     int
	Nsrf       int
	Npcm       int
	Emrk       []rune
	Wall       []WALL   // 壁
	Window     []WINDOW // 窓
	Snbk       []SNBK   // 日よけ
	PCM        []PCM    // 潜熱蓄熱材
	Sd         []RMSRF  // 壁体
	Mw         []MWALL  // 重量壁体
	Room       []ROOM   // 室
	Rdpnl      []RDPNL  // 輻射パネル
	Qrm, Qrmd  []QRM    // 日射、室内発熱熱取得
	Qetotal    QETOTAL
	Trdav      []float64
	Pcmiterate rune // PCM建材を使用し、かつ収束計算を行う場合は'y'
}

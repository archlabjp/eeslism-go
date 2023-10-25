package eeslism

type EqpType string

const (
	ROOM_TYPE       EqpType = "ROOM"       // 室
	RDPANEL_TYPE    EqpType = "RPNL"       // 放射パネル
	HCCOIL_TYPE     EqpType = "HCC"        // 冷温水コイル
	BOILER_TYPE     EqpType = "BOI"        // ボイラ－
	COLLECTOR_TYPE  EqpType = "COL"        // 太陽熱集熱器
	ACOLLECTOR_TYPE EqpType = "ACOL"       // 太陽熱集熱器
	REFACOMP_TYPE   EqpType = "REFA"       // ヒートポンプ（圧縮式冷凍機）
	STANK_TYPE      EqpType = "STANK"      // 蓄熱槽
	HEXCHANGR_TYPE  EqpType = "HEX"        // 熱交換器
	STHEAT_TYPE     EqpType = "STHEAT"     // 電気蓄熱暖房器
	THEX_TYPE       EqpType = "THEX"       // 全熱交換器
	DESI_TYPE       EqpType = "DESICCANT"  // デシカント槽
	EVAC_TYPE       EqpType = "EVPCOOLING" // 気化冷却器
	VAV_TYPE        EqpType = "VAV"
	VWV_TYPE        EqpType = "VWV"
)

const (
	COLLECTOR_PDT  = 'w'
	ACOLLECTOR_PDT = 'a'
)

type CATType string
type VAVType rune

const (
	// ---- Satoh Debug VAV  2000/10/30 ----
	VAV_PDT VAVType = 'A' // 空気
	VWV_PDT VAVType = 'W' // 温水

	PIPEDUCT_TYPE = "PIPE"
	DUCT_TYPE     = "DUCT"
	PIPE_PDT      = 'P'
	DUCT_PDT      = 'D'

	PUMP_TYPE = "PUMP"

	FAN_TYPE = "FAN"
	PUMP_PF  = 'P'
	FAN_PF   = 'F'

	PUMP_C  = "C"
	PUMP_Vv = "Vv"
	PUMP_Vr = "Vr"

	FAN_C  = "C"
	FAN_Vd = "Vd"
	FAN_Vs = "Vs"
	FAN_Vp = "Vp"
	FAN_Vr = "Vr"

	PV_TYPE = "PV"

	DIVERG_TYPE   = "B"  // 通過流体が水の分岐要素
	CONVRG_TYPE   = "C"  // 通過流体が水の合流要素
	DIVGAIR_TYPE  = "BA" // 通過流体が空気の分岐要素
	CVRGAIR_TYPE  = "CA" // 通過流体が空気の合流要素
	DIVERGCA_TYPE = "_B"
	CONVRGCA_TYPE = "_C"

	FLIN_TYPE  = "FLI" // 流入境界条件(システム経路への流入条件)
	GLOAD_TYPE = "GLD"

	HCLOAD_TYPE  = "HCLD"  // 仮想空調機コイル(直膨コイル)
	HCLOADW_TYPE = "HCLDW" // 仮想空調機コイル(冷・温水コイル)
	RMAC_TYPE    = "RMAC"  // ルームエアコン
	RMACD_TYPE   = "RMACD"

	QMEAS_TYPE = "QMEAS" // カロリーメータ
	VALV_TYPE  = "V"     // 弁およびダンパー
	TVALV_TYPE = "VT"    // 温調弁（水系統のみ）

	OMVAV_TYPE = "OMVAV"
	OAVAV_TYPE = "OAVAV"

	OUTDRAIR_NAME = "_OA"
	OUTDRAIR_PARM = "t=Ta x=xa *"

	CITYWATER_NAME = "_CW"
	CITYWATER_PARM = "t=Twsup *"
)

type ControlSWType rune

// 通過する流体の種類（a:空気（温度）、x:空気（湿度）、W:水））
type FliudType rune

const (
	AIR_FLD   FliudType = 'A' // 空気??
	AIRa_FLD  FliudType = 'a' // 空気（温度）
	AIRx_FLD  FliudType = 'x' // 空気（湿度）
	WATER_FLD FliudType = 'W' // 水

	HEATING_SYS = 'a'
	HVAC_SYS    = 'A'
	DHW_SYS     = 'W'

	THR_PTYP = 'T'
	CIR_PTYP = 'C'
	BRC_PTYP = 'B'

	DIVERG_LPTP = 'b'
	CONVRG_LPTP = 'c'
	IN_LPTP     = 'i' // 流入境界条件
	OUT_LPTP    = 'o'

	OFF_SW   ControlSWType = 'x' // 経路が停止中
	ON_SW    ControlSWType = '-' // 経路が動作中
	LOAD_SW  ControlSWType = 'F'
	FLWIN_SW ControlSWType = 'I'
	BATCH_SW ControlSWType = 'B'

	COOLING_LOAD  ControlSWType = 'C'
	HEATING_LOAD  ControlSWType = 'H'
	HEATCOOL_LOAD ControlSWType = 'L'

	COOLING_SW ControlSWType = 'C'
	HEATING_SW ControlSWType = 'H'

	TANK_FULL                 = 'F'
	TANK_EMPTY                = 'E'
	TANK_EMPTMP               = -777.0
	BTFILL      ControlSWType = 'F'
	BTDRAW      ControlSWType = 'D'

	SYSV_EQV = 'v'
	LOAD_EQV = 'L'
)

type COMPNT struct {
	Name       string     // 機器名称
	Roomname   string     // 機器の設置室名称（-room）
	Eqptype    EqpType    // 機器タイプ（"PIPE"など）
	Envname    string     // 配管等の周囲条件名称（-env）
	Exsname    string     // 方位名称
	Hccname    string     // VWV制御するときの制御対象熱交換器名称
	Rdpnlname  string     // VWV制御するときの制御対象床暖房（未完成）
	Idi        []ELIOType // 入口の識別記号 (len(Idi) == Nin)
	Ido        []ELIOType // 出口の識別記号（熱交換器の'C'、'H'や全熱交換器の'E'、'O'など）(len(Ido) == Nout)
	Tparm      string     // SYSCMPで定義された"-S"や"-V"以降の文字列を収録する
	Wetparm    string     // 湿りコイルの除湿時出口相対湿度の文字列を収録
	Omparm     string     // 集熱器が直列接続の場合に流れ方向に記載する
	Airpathcpy bool       // 空気経路の場合はtrue（湿度経路用にpathをコピーする）
	Control    ControlSWType
	Eqp        interface{} // 機器特有の構造体へのポインタ
	Neqp       int
	Ncat       int
	Nout       int // 出口の数
	Nin        int // 入口の数
	Nivar      int
	Ac         float64 // 集熱器面積[m2]
	PVcap      float64 // 太陽電池容量[W]
	Area       float64 // 太陽電池アレイ面積[m2]
	Ivparm     *float64
	Eqpeff     float64  // ボイラ室内置き時の室内供給熱量率 [-]
	Elouts     []*ELOUT // 機器出口の構造体へのポインタ（Nout個）
	Elins      []*ELIN  // 機器入口の構造体へのポインタ（Nin個）
	//	valv	*Valv

	Valvcmp *COMPNT // 三方弁の対となるValvのComptへのポインタ
	//	x,			/* バルブ開度 */
	//	xinit ;
	//	char	org ;		/* CONTRLで指定されているとき'y' それ以外は'n' */
	//	char	*OMfanName ;	// Valvが参照するファン風量
	MonPlistName string  // VALVで分岐などを流量比率で行う場合の観測対象のPlist名称
	MPCM         float64 // 電気蓄熱暖房器内臓PCMの容量[m3]
}

type ELOUT struct {
	Id      ELIOType      // 出口の識別番号（熱交換器の'C'、'H'や全熱交換器の'E'、'O'など）
	Pelmoid rune          // 終端の割り当てが完了していれば '-', そうでなければ 'x'
	Fluid   FliudType     // 通過する流体の種類（a:空気（温度）、x:空気（湿度）、W:水））
	Control ControlSWType // 経路の制御
	Sysld   rune          // 負荷を計算する場合は'y'、成り行きの場合は'n'
	G       float64       // 流量
	Q       float64       // 熱量
	Sysv    float64       // 連立方程式の答え
	Load    float64
	Co      float64   // 連立方程式の定数
	Coeffo  float64   // 出口の係数
	Coeffin []float64 // 入口の係数（入口複数の場合はそれぞれの係数）
	Ni      int       // 入口の数
	Sv      int
	Sld     int
	Cmp     *COMPNT // 機器出口の構造体が属する機器
	Elins   []*ELIN // 機器出口の構造体が関連する機器入口
	Lpath   *PLIST  // 機器出口が属する末端経路
	Eldobj  *ELOUT
	Emonitr *ELOUT
}

// 経路識別子
type ELIOType rune

const (
	ELIO_None  ELIOType = 0
	ELIO_G     ELIOType = 'G'
	ELIO_C     ELIOType = 'C' // 冷風?
	ELIO_H     ELIOType = 'H' // 温風?
	ELIO_D     ELIOType = 'D' // Tdry
	ELIO_d     ELIOType = 'd' // xdry
	ELIO_V     ELIOType = 'V' // Twet
	ELIO_v     ELIOType = 'v' // xwet
	ELIO_e     ELIOType = 'e' // 排気系統（エンタルピー） ?
	ELIO_E     ELIOType = 'E' // 排気系統（温度）?
	ELIO_O     ELIOType = 'O' // 給気系統（温度） ?
	ELIO_o     ELIOType = 'o' // 給気系統（エンタルピー） ?
	ELIO_x     ELIOType = 'x' // 空気湿度
	ELIO_f     ELIOType = 'f'
	ELIO_r     ELIOType = 'r'
	ELIO_W     ELIOType = 'W' // 温水温度
	ELIO_w     ELIOType = 'w'
	ELIO_t     ELIOType = 't' // 空気温度
	ELIO_ASTER ELIOType = '*'
	ELIO_SPACE ELIOType = ' '
)

type ELIN struct {
	Id     ELIOType // 入口の識別番号（熱交換器の'C'、'H'や全熱交換器の'E'、'O'など）
	Sysvin float64  // 連立方程式の答え
	Upo    *ELOUT   // 上流の機器の出口
	Upv    *ELOUT
	Lpath  *PLIST // 機器入口が属する末端経路
}

// SYSPTHに記載の機器
type PELM struct {
	Co  ELIOType // SYSPTHに記載の機器の出口の識別番号（熱交換器の'C'、'H'や全熱交換器の'E'、'O'など）
	Ci  ELIOType // SYSPTHに記載の機器の入口の識別番号（熱交換器の'C'、'H'や全熱交換器の'E'、'O'など）
	Cmp *COMPNT  // SYSPTH記載の機器の構造体
	//  Pelmx *PELM      // PELM構造体へのポインタ（コメントアウトされているため、Goのコードでは除外）
	Out *ELOUT // 機器の出口
	In  *ELIN  // 機器の入口
}

// 末端経路(主経路・または部分経路)
type PLIST struct {
	UnknownFlow int           // 末端経路が流量未知なら1、既知なら0
	Name        string        // 末端経路の名前
	Type        rune          // 貫流経路か循環経路かの判定
	Control     ControlSWType // 経路の制御情報
	Batch       bool          // バッチ運転を行う蓄熱槽のあるときtrue
	Org         bool          // 入力された経路のときtrue、複写された経路（空気系統の湿度経路）のとき false
	Plistname   string        // 末端経路の名前
	Lvc         int
	Nvalv       int      // 経路中のバルブ数
	Nvav        int      // 経路中のVAVユニットの数
	NOMVAV      int      // OM用変風量制御ユニット数
	N           int      // 流量計算の時の番号
	Go          *float64 // 流量の計算に使用される係数
	Gcalc       float64  // 温調弁によって計算された流量を記憶する変数
	G           float64  // 流量
	Rate        *float64 // 流量分配比
	Pelm        []*PELM  // 末端経路内の機器（バルブ、カロリーメータを除く)　OMVAVも除くべき？
	Plmvb       *PELM    // ??
	Lpair       *PLIST
	Plistt      *PLIST // 空気系当時の温度系統
	Plistx      *PLIST // 空気系当時の湿度系統
	Valv        *VALV  // 弁・ダンパーへの参照 (V,VT用)
	Mpath       *MPATH // システム経路 MPATH への逆参照
	Upplist     *PLIST
	Dnplist     *PLIST
	OMvav       *OMVAV // OMVAVへの参照 (OMVAV用)
}

// SYSPTHにおける';'で区切られる経路
// SYSPTH (1)--(N) MPATH (1)--(N) PLIST (1) -- (N) PELM
type MPATH struct {
	Name    string        // 経路名称
	Sys     byte          // 系統番号
	Type    byte          // 貫流経路か循環経路かの判定
	Fluid   FliudType     // 流体種別
	Control ControlSWType // 経路の制御情報
	NGv     int           // ガス導管数
	NGv2    int           // 開口率が2%未満のガス導管数
	Ncv     int           // 制御弁数
	Lvcmx   int           // 制御弁の接続数の最大値
	Plist   []*PLIST      // 末端経路
	Pl      []*PLIST      // 末端経路を格納する配列へのポインタ
	Rate    bool          // 流量比率(Plist[x].Rate)が入力されている経路ならtrue
	G0      *float64      // 流量比率設定時の既知流量へのポインタ
	Mpair   *MPATH        // 温度経路から湿度経路への参照
	Cbcmp   []*COMPNT     // 流量連立方程式を解くときに使用する分岐・合流機器
}

type SYSEQ struct {
	A byte
}

type VALV struct {
	Name     string
	Count    int
	X        float64  // バルブ開度
	Xinit    *float64 // バルブ開度の初期値
	Org      byte     // CONTRLで指定されているとき'y' それ以外は'n'
	Cmp      *COMPNT
	Cmb      *COMPNT
	Mon      *COMPNT
	Tin      *float64
	Tset     *float64
	Tout     *float64
	MGo      *float64
	Plist    *PLIST // 接続している末端経路への参照
	MonPlist *PLIST
}

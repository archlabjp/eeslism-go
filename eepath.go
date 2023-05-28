package main

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
	VAV_PDT VAVType = 'A'
	VWV_PDT VAVType = 'W'

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

	DIVERG_TYPE   = "B"
	CONVRG_TYPE   = "C"
	DIVGAIR_TYPE  = "BA"
	CVRGAIR_TYPE  = "CA"
	DIVERGCA_TYPE = "_B"
	CONVRGCA_TYPE = "_C"

	FLIN_TYPE  = "FLI"
	GLOAD_TYPE = "GLD"

	HCLOAD_TYPE  = "HCLD"
	HCLOADW_TYPE = "HCLDW"
	RMAC_TYPE    = "RMAC"
	RMACD_TYPE   = "RMACD"

	QMEAS_TYPE = "QMEAS"
	VALV_TYPE  = "V"
	TVALV_TYPE = "VT"

	OMVAV_TYPE = "OMVAV"
	OAVAV_TYPE = "OAVAV"

	OUTDRAIR_NAME = "_OA"
	OUTDRAIR_PARM = "t=Ta x=xa *"

	CITYWATER_NAME = "_CW"
	CITYWATER_PARM = "t=Twsup *"
)

type ControlSWType rune

const (
	AIR_FLD   = 'A'
	AIRa_FLD  = 'a'
	AIRx_FLD  = 'x'
	WATER_FLD = 'W'

	HEATING_SYS = 'a'
	HVAC_SYS    = 'A'
	DHW_SYS     = 'W'

	THR_PTYP = 'T'
	CIR_PTYP = 'C'
	BRC_PTYP = 'B'

	DIVERG_LPTP = 'b'
	CONVRG_LPTP = 'c'
	IN_LPTP     = 'i'
	OUT_LPTP    = 'o'

	OFF_SW   ControlSWType = 'x'
	ON_SW    ControlSWType = '-'
	LOAD_SW  ControlSWType = 'F'
	FLWIN_SW ControlSWType = 'I'
	BATCH_SW ControlSWType = 'B'

	COOLING_LOAD  = 'C'
	HEATING_LOAD  = 'H'
	HEATCOOL_LOAD = 'L'

	COOLING_SW = 'C'
	HEATING_SW = 'H'

	TANK_FULL   = 'F'
	TANK_EMPTY  = 'E'
	TANK_EMPTMP = -777.0
	BTFILL      = 'F'
	BTDRAW      = 'D'

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
	Idi        []ELIOType // 入口の識別記号
	Ido        []ELIOType // 出口の識別記号（熱交換器の'C'、'H'や全熱交換器の'E'、'O'など）
	Tparm      string     // SYSCMPで定義された"-S"や"-V"以降の文字列を収録する
	Wetparm    string     // 湿りコイルの除湿時出口相対湿度の文字列を収録
	Omparm     string     //
	Airpathcpy rune       // 空気経路の場合は'Y'（湿度経路用にpathをコピーする）
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
	Fluid   rune          // 通過する流体の種類（a:空気（温度）、x:空気（湿度）、W:水））
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

type ELIOType rune

const (
	ELIO_None  ELIOType = 0
	ELIO_C     ELIOType = 'C'
	ELIO_H     ELIOType = 'H'
	ELIO_D     ELIOType = 'D' // Tdry
	ELIO_d     ELIOType = 'd' // xdry
	ELIO_V     ELIOType = 'V' // Twet
	ELIO_v     ELIOType = 'v' // xwet
	ELIO_e     ELIOType = 'e'
	ELIO_E     ELIOType = 'E'
	ELIO_O     ELIOType = 'O'
	ELIO_o     ELIOType = 'o'
	ELIO_x     ELIOType = 'x'
	ELIO_f     ELIOType = 'f'
	ELIO_r     ELIOType = 'r'
	ELIO_W     ELIOType = 'W'
	ELIO_w     ELIOType = 'w'
	ELIO_t     ELIOType = 't'
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

type PELM struct {
	Co  ELIOType // SYSPTHに記載の機器の出口の識別番号（熱交換器の'C'、'H'や全熱交換器の'E'、'O'など）
	Ci  ELIOType // SYSPTHに記載の機器の入口の識別番号（熱交換器の'C'、'H'や全熱交換器の'E'、'O'など）
	Cmp *COMPNT  // SYSPTH記載の機器の構造体
	//  Pelmx *PELM      // PELM構造体へのポインタ（コメントアウトされているため、Goのコードでは除外）
	Out *ELOUT // 機器の出口
	In  *ELIN  // 機器の入口
}

type PLIST struct {
	UnknownFlow int           // 末端経路が流量未知なら1、既知なら0
	Name        string        // 末端経路の名前
	Type        rune          // 貫流経路か循環経路かの判定
	Control     ControlSWType // 経路の制御情報
	Batch       rune          // バッチ運転を行う蓄熱槽のあるとき'y'、無いとき 'n'
	Org         rune          // 入力された経路のとき'y'、複写された経路（空気系統の湿度経路）のとき'n'
	Plistname   string        // 末端経路の名前
	Nelm        int           // 末端経路内の機器の数
	Lvc         int
	Nvalv       int      // 経路中のバルブ数
	Nvav        int      // 経路中のVAVユニットの数
	NOMVAV      int      // OM用変風量制御ユニット数
	N           int      // 流量計算の時の番号
	Go          *float64 // 流量の計算に使用される係数
	Gcalc       float64  // 温調弁によって計算された流量を記憶する変数
	G           float64  // 流量
	Rate        *float64 // 流量分配比
	Pelm        []*PELM  // 末端経路内の機器
	Plmvb       *PELM
	Lpair       *PLIST
	Plistt      *PLIST // 空気系当時の温度系統
	Plistx      *PLIST // 空気系当時の湿度系統
	Valv        *VALV
	Mpath       *MPATH
	Upplist     *PLIST
	Dnplist     *PLIST
	OMvav       *OMVAV
}

type MPATH struct {
	Name    string        // 経路名称
	Sys     byte          // 系統番号
	Type    byte          // 貫流経路か循環経路かの判定
	Fluid   rune          // 流体
	Control ControlSWType // 経路の制御情報
	Nlpath  int           // 末端経路数
	NGv     int           // ガス導管数
	NGv2    int           // 開口率が2%未満のガス導管数
	Ncv     int           // 制御弁数
	Lvcmx   int           // 制御弁の接続数の最大値
	Plist   []PLIST       // 末端経路
	Pl      []*PLIST      // 末端経路を格納する配列へのポインタ
	Rate    rune          // 流量比率が入力されている経路なら'Y'
	G0      *float64      // 流量比率設定時の既知流量へのポインタ
	Mpair   *MPATH
	Cbcmp   []*COMPNT // 流量連立方程式を解くときに使用する分岐・合流機器
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
	Plist    *PLIST
	MonPlist *PLIST
}

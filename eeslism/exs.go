package eeslism

// 外表面熱伝達率の設定方法
type AloType rune

const (
	Alotype_None     AloType = 0
	Alotype_V        AloType = 'V' // 外表面熱伝達率の設定方法: 風速から計算
	Alotype_Fix      AloType = 'F' // 外表面熱伝達率の設定方法: 23.0固定
	Alotype_Schedule AloType = 'S' // 外表面熱伝達率の設定方法: スケジュール
)

// 外表面種別
type EXSFType rune

const (
	EXSFType_None EXSFType = 0
	EXSFType_S    EXSFType = 'S' // 外表面種別: 一般外表面
	EXSFType_E    EXSFType = 'E' // 外表面種別: 地下
	EXSFType_e    EXSFType = 'e' // 外表面種別: 地表面
)

// 外表面方位デ－タ
type EXSF struct {
	Name    string
	Alotype AloType  // 外表面熱伝達率の設定方法 V:風速から計算、F:23.0固定、S:スケジュール
	Typ     EXSFType // 一般外表面'S',地下'E', 地表面'e'

	// --- 事前計算する日射関連のパラメータ群 ---

	Wa    float64  // 方位角 [deg]
	Wb    float64  // 傾斜角 [deg]
	Rg    float64  // 前面の日射反射率 [-]
	Fs    float64  // 天空を見る形態係数 [-]
	Wz    float64  // 傾斜角Wbのcos
	Ww    float64  // 傾斜角Wbのsin ×  方位角Waのsin
	Ws    float64  // 傾斜角Wbのsin ×  方位角Waのcos
	Swb   float64  // 傾斜角Wbのsin
	CbSa  float64  // 傾斜角Wbのcos ×  方位角Waのsin
	CbCa  float64  // 傾斜角Wbのcos ×  方位角Wbのsin
	Cwa   float64  // 方位角Waのcos
	Swa   float64  // 方位角Wbのsin
	Alo   *float64 // 外表面総合熱伝達率 [-] (Alotype が Sの場合のみ)
	Z     float64  // 地中深さ
	Erdff float64  // 土の熱拡散率 [m2/s]

	// --- 時々刻々の計算値 ---

	Cinc   float64 // 入射角のcos
	Tazm   float64 // 見掛けの方位角のtan
	Tprof  float64 // プロファイル角のtan
	Gamma  float64 // 見かけの方位角 [rad]
	Prof   float64 // プロファイル角 [rad]
	Idre   float64 // 直逹日射  [W/m2]
	Idf    float64 // 拡散日射  [W/m2]
	Iw     float64 // 全日射    [W/m2]
	Rn     float64 // 夜間輻射  [W/m2]
	Tearth float64 // 地中温度
}

// 外表面方位デ－タ
type EXSFS struct {
	Nexs int
	Exs  []EXSF // 外表面方位デ－タ

	// ---- 外表面熱伝達率 ----

	Alotype AloType  // 外表面熱伝達率の設定方法 'V':風速から計算、'F':23.0固定、'S':スケジュール
	Alosch  *float64 // 外表面熱伝達率 [-]  (Alotype が Sの場合のみ)

	// 地表面境界
	EarthSrfFlg bool // 地表面境界がある場合はtrue
}

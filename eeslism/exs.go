/*
exs.go (External Surface Data Structures)

このファイルは、建物の外部日射面（壁、屋根、窓、地盤など）に関するデータ構造を定義します。
これらの構造体は、日射熱取得、熱損失、および周囲環境との熱交換をモデル化するために用いられます。

建築環境工学的な観点:
- **外部日射面の分類**: 建物の外部日射面は、
  その向き（方位角、傾斜角）、熱的特性、および周囲環境との関係によって分類されます。
  - `AloType`: 外表面熱伝達率の設定方法（風速から計算、固定値、スケジュール）。
  - `EXSFType`: 外表面種別（一般外表面、地下、地表面）。
  これらの分類は、各表面の熱的挙動を正確にモデル化するために重要です。
- **EXSF構造体**: 個々の外部日射面に関する詳細な情報を格納します。
  - `Name`: 外部日射面の名称。
  - `Wa`, `Wb`: 方位角、傾斜角。太陽位置との相対関係を定義します。
  - `Rg`: 前面の日射反射率。地面からの反射日射を考慮します。
  - `Fs`: 天空を見る形態係数。天空からの放射熱交換を考慮します。
  - `Wz`, `Ww`, `Ws`, `Swb`, `CbSa`, `CbCa`, `Cwa`, `Swa`: 方位角と傾斜角から導出される三角関数値。
  - `Alo`: 外表面総合熱伝達率。外表面からの熱伝達のしやすさを示します。
  - `Z`: 地中深さ。地下に接する表面の熱伝達をモデル化します。
  - `Erdff`: 土の熱拡散率。地中熱交換をモデル化します。
  - `Cinc`, `Tazm`, `Tprof`, `Gamma`, `Prof`: 日射入射角に関するパラメータ。
  - `Idre`, `Idf`, `Iw`: 直達日射、拡散日射、全日射。日射熱取得量をモデル化します。
  - `Rn`: 夜間輻射。夜間の放射熱損失をモデル化します。
  - `Tearth`: 地中温度。地中熱交換をモデル化します。
- **EXSFS構造体**: 複数の外部日射面を管理するための構造体です。
  - `Exs`: `EXSF`構造体のスライス。
  - `Alotype`, `Alosch`: 外表面熱伝達率の設定方法とスケジュール。
  - `EarthSrfFlg`: 地表面境界の有無。

このファイルは、建物の外部日射面からの熱損失・熱取得を正確にモデル化し、
日射熱取得の抑制、冷房負荷の軽減、
昼光利用の最適化、および省エネルギー対策の検討を行うための重要な役割を果たします。
*/
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
	Exs []*EXSF // 外表面方位デ－タ

	// ---- 外表面熱伝達率 ----

	Alotype AloType  // 外表面熱伝達率の設定方法 'V':風速から計算、'F':23.0固定、'S':スケジュール
	Alosch  *float64 // 外表面熱伝達率 [-]  (Alotype が Sの場合のみ)

	// 地表面境界
	EarthSrfFlg bool // 地表面境界がある場合はtrue
}
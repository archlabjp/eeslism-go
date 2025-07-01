/*
wthrd.go (Weather Data Structures)

このファイルは、建物のエネルギーシミュレーションにおける気象データに関するデータ構造を定義します。
これらの構造体は、地域情報、時刻ごとの気象要素、および気象データファイルの読み込み方法などを管理するために用いられます。

建築環境工学的な観点:
- **地域情報 (LOCAT)**:
  `LOCAT`構造体は、シミュレーション対象の建物の地理的情報（地名、緯度、経度、標準子午線）を格納します。
  - `Lat`, `Lon`, `Ls`: 緯度、経度、標準子午線。太陽位置計算の基礎となります。
  - `Daymxert`, `Tgrav`, `DTgr`: 地中温度計算に関するパラメータ。地盤からの熱伝達をモデル化します。
  - `Twsup`: 月ごとの給水温度。給湯負荷計算などに用いられます。
- **時刻ごとの気象データ (WDAT)**:
  `WDAT`構造体は、シミュレーションの各時間ステップにおける気象要素を格納します。
  - `T`, `X`, `RH`, `H`: 気温、絶対湿度、相対湿度、エンタルピー。室内の温湿度環境や熱負荷計算に用いられます。
  - `Idn`, `Isky`, `Ihor`: 法線面直達日射、水平面天空日射、水平面全日射。日射熱取得量をモデル化します。
  - `sunalt`, `sunazm`, `Sh`, `Sw`, `Ss`, `Solh`, `SolA`: 太陽高度、方位角、方向余弦。太陽位置をモデル化します。
  - `CC`: 雲量。日射量や夜間放射量の計算に用いられます。
  - `RN`, `Rsky`: 夜間輻射、大気放射量。放射熱損失をモデル化します。
  - `Wv`, `Wdre`: 風速、風向。外表面熱伝達率や換気量に影響します。
  - `RNtype`: 夜間放射量の計算方法。
  - `Intgtsupw`: 給水温度の補間フラグ。
  - `EarthSurface`: 地表面温度。地盤からの熱伝達をモデル化します。
- **気象データ項目のポインター (WDPT)**:
  `WDPT`構造体は、VCFILE形式の気象データ入力時に、
  各気象要素へのポインターを格納します。
  これにより、気象データを効率的に読み込み、
  シミュレーションに利用できます。

このファイルは、建物のエネルギーシミュレーションにおいて、
気象データを正確にモデル化し、
熱負荷計算、エネルギー消費量予測、
および省エネルギー対策の検討を行うための重要な役割を果たします。
*/
package eeslism

type LOCAT struct {
	Name string  // 地名
	Lat  float64 // 緯度[deg]
	Lon  float64 // 経度[deg]
	Ls   float64 // 標準子午線[deg]

	// 地中温度計算用
	Daymxert int
	Tgrav    float64
	DTgr     float64

	// 月毎の給水温度
	Twsup [12]float64
}

// 気象デ－タ
type WDAT struct {
	T              float64 // 気温
	X              float64 // 絶対湿度  [kg/kg]
	RH             float64 // 相対湿度 [%]
	H              float64 // エンタルピ [J/kg]
	Idn            float64 // 法線面直逹日射 [W/m2]
	Isky           float64 // 水平面天空日射 [W/m2]
	Ihor           float64 // 水平面全日射   [W/m2]
	sunalt, sunazm float64
	Sh, Sw, Ss     float64 // 太陽光線の方向余弦
	Solh, SolA     float64 // 太陽位置
	CC             float64 // 雲量
	RN             float64 // 夜間輻射 [W/m2]
	Rsky           float64 // 大気放射量[W/m2] higuchi 070918
	Wv             float64 // 風速 [m/s]
	Wdre           float64 // 風向　１６方位

	RNtype rune // 気象データ項目  C:雲量　R:夜間放射量[W/m2]

	Intgtsupw    rune      // 給水温度を補完する場合は'Y'、しない場合は'N'  デフォルトは'N'
	Twsup        float64   // 給水温度
	EarthSurface []float64 // 地表面温度[℃]
}

// 気象データ項目のポインター  VCFILEからの入力時
type WDPT struct {
	Ta   []float64 //気温
	Xa   []float64 //絶対湿度
	Rh   []float64 //相対湿度
	Idn  []float64 //法線面直逹日射
	Isky []float64 //水平面天空日射
	Ihor []float64 //水平面全日射
	Cc   []float64 //雲量
	Rn   []float64 //夜間輻射
	Wv   []float64 //風速
	Wdre []float64 //風向
}
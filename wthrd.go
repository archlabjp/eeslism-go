package main

type LOCAT struct {
	Name string  /*地名*/
	Lat  float64 /*緯度[deg] */
	Lon  float64 /*経度[deg] */
	Ls   float64 /*標準子午線[deg] */

	/*地中温度計算用*/
	Daymxert int
	Tgrav    float64
	DTgr     float64

	/*月毎の給水温度*/
	Twsup [12]float64
}

type WDAT struct /*気象デ－タ         */
{
	T              float64 /*気温                  */
	X              float64 /*絶対湿度  [kg/kg]     */
	RH             float64
	H              float64 /*エンタルピ [J/kg]     */
	Idn            float64 /*法線面直逹日射 [W/m2] */
	Isky           float64 /*水平面天空日射 [W/m2] */
	Ihor           float64 /*水平面全日射   [W/m2] */
	sunalt, sunazm float64
	Sh, Sw, Ss     float64 /*太陽光線の方向余弦    */
	Solh, SolA     float64 // 太陽位置
	CC             float64 /*雲量                  */
	RN             float64 /*夜間輻射 [W/m2]       */
	Rsky           float64 /*大気放射量[W/m2] higuchi 070918 */
	Wv             float64 /*風速 [m/s]            */
	Wdre           float64 /*風向　１６方位        */

	RNtype rune /*気象データ項目  C:雲量　R:夜間放射量[W/m2] */

	Intgtsupw rune
	// 給水温度を補完する場合は'Y'、しない場合は'N'
	// デフォルトは'N'
	Twsup        float64 /*給水温度              */
	EarthSurface []float64
	// 地表面温度[℃]
}

type WDPT struct /*気象データ項目のポインター  VCFILEからの入力時 */
{
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

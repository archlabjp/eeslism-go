package main

type EXSF struct /*外表面方位デ－タ                */
{
	Name    string
	Alotype rune // 外表面熱伝達率の設定方法
	// V:風速から計算、F:23.0固定、S:スケジュール
	Typ  rune    /*一般外表面'S',地下'E', 地表面'e'  */
	Wa   float64 /*方位角                 */
	Wb   float64 /*傾斜角                 */
	Rg   float64 /*前面の日射反射率       */
	Fs   float64 /*天空を見る形態係数     */
	Wz   float64
	Ww   float64
	Ws   float64 /*傾斜面の方向余弦       */
	Swb  float64
	Cbsa float64
	Cbca float64
	//alosch []float64		// 外表面の熱伝達率スケジュール
	Cwa float64
	Swa float64
	/* tprof float64 tazm計算用係数  */
	Alo   *float64 /*外表面総合熱伝達率　　　*/
	Z     float64  /*地中深さ　　　　　　　　*/
	Erdff float64  /*土の熱拡散率m2/s　　　　*/

	/*方位別日射関連デ－タ */
	Cinc   float64 /*入射角のcos             */
	Tazm   float64 /*見掛けの方位角のtan     */
	Tprof  float64 /*プロファイル角のtan     */
	Gamma  float64 // 見かけの方位角
	Prof   float64 // プロファイル角
	Idre   float64 /*直逹日射  [W/m2]        */
	Idf    float64 /*拡散日射  [W/m2]        */
	Iw     float64 /*全日射    [W/m2]        */
	Rn     float64 /*夜間輻射  [W/m2]        */
	Tearth float64 /*地中温度　　　　　　　　*/
	End    int     // 要素数(インデックス0のみに設定)
}

type EXSFS struct {
	Nexs    int
	Exs     []EXSF
	Alosch  *float64
	Alotype rune // 外表面熱伝達率の設定方法
	// V:風速から計算、F:23.0固定、S:スケジュール
	EarthSrfFlg rune
	// 地表面境界がある場合'Y'
}

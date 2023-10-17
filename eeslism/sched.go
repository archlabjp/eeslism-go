package eeslism

// 季節設定
type SEASN struct {
	name       string // 季節名 (sname)
	N          int
	sday, eday []int // 開始日・終了日(通日)
}

// 曜日設定
type WKDY struct {
	name string  // 曜日名 (wname)
	wday [15]int //TODO:たぶん長さ8で十分
}

// 一日の設定量スケジュ－ル
type DSCH struct {
	name         string // 設定値名 (vdname)
	N            int
	stime, etime []int     // 開始時分, 終了時分
	val          []float64 // 設定値
}

// 一日の切り替えスケジュ－ル
type DSCW struct {
	name         string   // 切替設定名 (wdname)
	dcode        [10]rune // 切替名 (mode)
	N            int      // 切替時間帯の数(stime,mode,etimeのスライスの長さ)
	stime, etime []int    // 切替開始時分, 切替終了時分
	Nmod         int      // 切替モードの種類の数 (modeの重複を除いた数)
	mode         []rune   // 切替モード
}

type SCH struct /*スケジュ－ル*/
{
	name string
	Type rune
	day  [366]int
}

// 一日の設定値、切換スケジュールおよび季節、曜日の指定
// See: [eeslism.]
type SCHDL struct {
	Seasn []SEASN // SCHTBデータセット:季節設定 (-wkd or WKD)
	Wkdy  []WKDY  // SCHTBデータセット:曜日設定 (-wkd)
	Dsch  []DSCH  // SCHTBデータセット:1日の設定値スケジュール定義(-v or VL)
	Dscw  []DSCW  // SCHTBデータセット:1日の切替設定スケジュール定義(-s or SW)

	Sch []SCH // SCHNMデータセット: 年間の設定値スケジュール
	Scw []SCH // SCHNMデータセット: 年間の切替スケジュール

	Val []float64 // `Sch`の要素数と同数
	Isw []rune    // `Scw`の要素数と同数
}

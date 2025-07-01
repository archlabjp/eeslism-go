/*
sched.go (Schedule Data Structures for Building Energy Simulation)

このファイルは、建物のエネルギーシミュレーションにおける運用スケジュールを定義するためのデータ構造を提供します。
これらの構造体は、空調システム、照明、換気、在室状況などの時間的変化をモデル化するために用いられます。

建築環境工学的な観点:
- **運用パターンのモデル化**: 建物のエネルギー消費量や室内環境は、
  居住者の活動、機器の運転、照明の使用など、
  様々な運用パターンによって大きく左右されます。
  このファイルで定義される`SEASN`（季節設定）、`WKDY`（曜日設定）、
  `DSCH`（1日の設定値スケジュール）、`DSCW`（1日の切替スケジュール）、
  `SCH`（年間スケジュール）、`SCHDL`（スケジュール全体）などの構造体は、
  これらの運用パターンを柔軟に記述することを可能にします。
- **時間的変化の考慮**: 
  - `SEASN`: 年間を複数の季節に分割し、季節ごとに異なる運用パターンを適用できます。
  - `WKDY`: 曜日ごとに異なる運用パターンを適用できます（平日、週末、祝日など）。
  - `DSCH`, `DSCW`: 1日を複数の時間帯に分割し、時間帯ごとに設定値（温度、湿度など）や
    運転モード（ON/OFF、暖房、冷房など）を変化させることができます。
  - `SCH`: 年間を通じたスケジュールを定義し、
    各日の運用パターンを決定します。
- **省エネルギー運転の実現**: 適切なスケジュール設定をモデル化することで、
  - **デマンド制御**: 実際の熱負荷に応じて機器の運転を調整し、無駄なエネルギー消費を削減します。
  - **最適制御**: 快適性を維持しつつ、エネルギー消費を最小化する運転戦略を検討します。
  - **ピークカット**: 電力料金の安い時間帯に蓄熱を行うなど、
    エネルギーコストの削減に貢献します。
- **快適性の維持**: 室内温度や湿度などの環境変数を目標値に維持するためのスケジュールをモデル化することで、
  居住者の快適性を確保できます。

このファイルは、建物のエネルギーシミュレーションにおいて、
複雑な運用パターンを正確にモデル化し、
省エネルギー、快適性、および運用効率の向上を図るための重要な役割を果たします。
*/
package eeslism

// 季節設定
type SEASN struct {
	name       string // 季節名 (sname)
	N          int    // sday, edayの配列の長さ
	sday, eday []int  // 開始日・終了日(通日)
}

// 曜日設定
type WKDY struct {
	name string // 曜日名 (wname)
	wday [8]bool
}

// 一日の設定値スケジュ－ル
type DSCH struct {
	name         string    // 設定値名 (vdname)
	N            int       // stime, etimeの配列の長さ
	stime, etime []int     // 開始時分, 終了時分
	val          []float64 // 設定値
}

// 一日の切り替えスケジュ－ル
type DSCW struct {
	name         string            // 切替設定名 (wdname)
	dcode        [10]ControlSWType // 切替名 (mode)
	N            int               // 切替時間帯の数(stime,mode,etimeのスライスの長さ)
	stime, etime []int             // 切替開始時分, 切替終了時分
	Nmod         int               // 切替モードの種類の数 (modeの重複を除いた数)
	mode         []ControlSWType   // 切替モード
}

type SCH struct /*スケジュ－ル*/
{
	name string
	Type rune
	day  [366]int //インデックス0は使用しない
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

	Val []float64       // 設定値? (`Sch`の要素数と同数)
	Isw []ControlSWType // 切替状態? (`Scw`の要素数と同数)
}
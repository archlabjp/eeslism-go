/*
simcon.go (Simulation Control Data Structures)

このファイルは、建物のエネルギーシミュレーションの実行を制御するためのデータ構造を定義します。
これらの構造体は、シミュレーションの期間、時間間隔、出力設定、
および気象データや入力ファイルに関する情報などを管理するために用いられます。

建築環境工学的な観点:
- **シミュレーションの基本設定 (SIMCONTL)**:
  `SIMCONTL`構造体は、シミュレーションの全体的な設定を格納します。
  - `File`, `Title`: 入力ファイル名とシミュレーションのタイトル。
  - `Wfname`, `Wdtype`: 気象データファイル名とその種別。
  - `Ofname`: 出力ファイル名。
  - `DTm`: 計算時間間隔 [s]。シミュレーションの細かさを決定し、計算精度と計算時間に影響します。
  - `Daystartx`, `Daystart`, `Dayend`: 助走期間、本計算開始日、計算終了日。シミュレーション期間を定義します。
  - `MaxIterate`: 最大収束計算回数。非線形な熱的挙動を持つシステムにおいて、
    計算の収束を確保するための反復回数を定義します。
  - `Perio`: 周期定常計算のフラグ。
  - `Dayprn`: データ出力日。特定の日のみ詳細な結果を出力する際に用いられます。
  - `Loc`: 地域データ（緯度、経度など）。
  - `Wdpt`: 気象データ。
- **出力設定の制御 (PrintType, FLOUT)**:
  `PrintType`は、出力するデータの種類を識別するための定数です。
  `FLOUT`構造体は、各出力ファイルの設定情報（ファイル名、ファイルポインター、出力タイプ）を格納します。
  これにより、ユーザーは必要な情報を効率的に取得し、
  分析や検証を容易にします。
- **時間管理 (DAYTM)**:
  `DAYTM`構造体は、通日、年、月、日、時刻などの時間情報を格納します。
  これにより、シミュレーション結果を時間軸に沿って分析できます。
- **外部ファイルとの連携 (VCFILE)**:
  `VCFILE`構造体は、境界条件や負荷入力用ファイルに関する情報を格納します。
  これにより、外部ファイルから動的にデータを読み込むことが可能になります。

このファイルは、建物のエネルギーシミュレーションの実行を制御し、
シミュレーションの正確性、効率性、および出力内容を決定するための重要な役割を果たします。
*/
package eeslism

import (
	"io"
)

const EEVERSION = "ES4.6"

// 出力種別
type PrintType string

const (
	// --- 時刻別計算値 ---

	PRTHWD     PrintType = "_wd"  // 時間別計算値(気象データ出力)
	PRTREV     PrintType = "_re"  // 時間別計算値(毎時室温、MRTの出力)
	PRTHROOM   PrintType = "_rm"  // 時間別計算値(放射パネルの出力)
	PRTHELM    PrintType = "_rqe" // 時間別計算値(要素別熱損失・熱取得)
	PRTHELMSF  PrintType = "_sfe" // 時間別計算値(要素別熱損失・熱取得) 表面?
	PRTPATH    PrintType = "_sp"  // 時間別計算値(システム経路の温湿度出力)
	PRTCOMP    PrintType = "_sc"  // 時間別計算値(機器の出力)
	PRTHRSTANK PrintType = "_tk"  // 時間別計算値(蓄熱槽内温度分布の出力)

	PRTPMV PrintType = "_pm"  // 時間別計算値(PMV計算)
	PRTQRM PrintType = "_rq"  // 時間別計算値(日射、室内熱取得の出力)
	PRTRSF PrintType = "_sf"  // 時間別計算値(室内表面温度の出力)
	PRTSFQ PrintType = "_sfq" // 時間別計算値(室内表面熱流の出力)
	PRTSFA PrintType = "_sfa" // 時間別計算値(室内表面熱伝達率の出力)
	PRTWAL PrintType = "_wl"  // 時間別計算値(壁体内部温度の出力)
	PRTSHD PrintType = "_shd" // 時間別計算値(日よけの影面積の出力)
	PRTPCM PrintType = "_pcm" // 時間別計算値(潜熱蓄熱材の状態値の出力)

	// --- 日別計算値 ---

	PRTWK     PrintType = "_wk"  // 計算年月日出力
	PRTDYRM   PrintType = "_dr"  // 日別計算値(部屋ごとの熱集計結果出力)
	PRTDYHELM PrintType = "_dqe" // 日別計算値(要素別熱損失・熱取得)
	PRTDQR    PrintType = "_dqr" // 日別計算値(日射、室内熱取得の出力)
	PRTDYSF   PrintType = "_dsf" // 日別計算値(日積算壁体貫流熱取得の出力)
	PRTDYCOMP PrintType = "_dc"  // 日別計算値(システム要素機器の日集計結果出力)
	PRTDWD    PrintType = "_dwd" // 日別計算値(気象データ日集計値出力)

	// --- 月別計算値 ---

	PRTMNRM   PrintType = "_mr"  // 月別計算値(部屋ごとの熱集計結果出力)
	PRTMNCOMP PrintType = "_mc"  // 月別計算値(システム要素機器の月集計結果出力)
	PRTMWD    PrintType = "_mwd" // 月別計算値(気象データ月集計値出力)

	// --- 月-時刻計算値 ---

	PRTMTCOMP PrintType = "_mt" // 月-時刻計算値(部屋ごとの熱集計結果出力)

	// SYSV_EQV = 'v'
	// LOAD_EQV = 'L'
)

type SIMCONTL struct {
	File       string        // 入力ファイル名
	Title      string        // 題目、注釈
	Wfname     string        // 気象データファイル名 (GDAT.FILE.w)
	Ofname     string        // 出力ファイル名	(GDAT.FILE.out)
	Unit       string        // 単位系
	Unitdy     string        //
	Timeid     []rune        // 時間別計算値出力識別子 ?
	Helmkey    rune          // 要素別熱取得、熱損失計算 'y'
	Wdtype     rune          // 気象データファイル種別 'H':HASP標準形式　'E':VCFILE入力形式 */
	Perio      rune          // 周期定常計算の時'y'
	Fwdata     io.ReadSeeker // 気象データファイルのファイルポインタ
	Fwdata2    io.ReadSeeker // 気象データファイルのファイルポインタ(なぜ2つあるのか?)
	Ftsupw     []byte        // 給水温度データのファイル(バイナリ)
	Daystartx  int           // 助走計算開始日
	Daystart   int           // 本計算開始日
	Dayend     int           // 計算終了日
	Daywk      []int         // 計算日 ??
	Dayprn     []int         // データ出力日
	Dayntime   int           // 1日あたりの計算回数
	Ntimehrprt int           // 時間別計算値出力回数
	Ntimedyprt int           // 日別計算値出力回数
	Nhelmsfpri int           // 要素別壁体表面温度出力回数
	Nvcfile    int           // 境界条件、負荷入力用ファイルの数
	Vcfile     []VCFILE      // 境界条件、負荷入力用ファイル等々???
	Loc        *LOCAT        // 地域データ
	Wdpt       WDPT          // 気象データ
	DTm        int           // 計算時間間隔 [s] (GDAT.RUN.dTime)
	Sttmm      int           // 計算開始時刻 (GDAT.RUN.Stime)
	MaxIterate int           // 最大収束回数 (GDAT.RUN.MaxIterate)
}

// 出力ファイルの設定情報
type FLOUT struct {
	Fname string    // 出力ファイル名
	F     io.Writer // 出力ファイルのファイルポインタ
	Idn   PrintType // 出力ファイルの種類
}

type VCFILE struct {
	Fi    io.ReadSeeker // ファイルポインタ
	Ad    int64         // ファイルの先頭アドレス
	Ic    int           // ファイルの種類??
	Name  string        // ファイル名
	Fname string        // ファイル名
	Estl  ESTL          // 要素データ??
	Tlist []TLIST       // 時刻データ??
}

type DAYTM struct {
	DayOfYear int     // 通日 (day)
	Year      int     // 年
	Mon       int     // 月
	Day       int     // 日
	Ddpri     int     // 日積算値出力
	Time      float64 // 時刻??
	Ttmm      int     // 時刻??
	Tt        int     // 時刻??
}
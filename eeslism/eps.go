/*
eps.go (Energy Performance Simulation Data Structures)

このファイルは、建物のエネルギー性能シミュレーションの結果を整理し、
報告書やグラフ作成に利用するためのデータ構造を定義します。
これらの構造体は、シミュレーション結果の解釈、分析、および可視化に不可欠です。

建築環境工学的な観点:
  - **シミュレーション結果のメタデータ (ESTL)**:
    `ESTL`構造体は、シミュレーション結果に関する注釈やメタデータを格納します。
  - `Flid`: ファイル種別記号。
  - `Title`: シミュレーションの標題。
  - `Wdatfile`: 使用した気象データファイル名。
  - `Tid`: 入力データ種別（時刻別、日別など）。
  - `Unit`: 出力データの単位。
  - `Timeid`: 時刻データの表示形式。
  - `Wdloc`: 地域情報（地名、緯度、経度など）。
  - `Catnm`: 要素カタログ名データ。
  - `Ntime`, `dtm`: データ数、計算時間間隔。
    これらの情報は、シミュレーション結果の再現性や、
    他のシミュレーションとの比較可能性を確保するために重要です。
  - **要素カタログ名データ (CATNM)**:
    `CATNM`構造体は、シミュレーションで使用された各機器カテゴリの名称、
    機器数、およびデータ項目数を格納します。
    これにより、シミュレーションモデルの構成を把握できます。
  - **時刻データ (TMDT)**:
    `TMDT`構造体は、年、月、日、曜日、時刻などの時間情報を格納します。
    これにより、シミュレーション結果を時間軸に沿って分析できます。
  - **シミュレーション結果の項目定義 (TLIST)**:
    `TLIST`構造体は、シミュレーション結果の各項目（温度、湿度、熱量、エネルギーなど）の定義を格納します。
  - `Vtype`: データ種別（温度、絶対湿度、熱量など）。
  - `Stype`: データ処理種別（積算値、平均値、最小値、最大値など）。
  - `Ptype`: データ型（文字型、整数型、実数型）。
  - `Fval`, `Ival`, `Cval`: 実際のデータ値。
    これにより、シミュレーション結果を柔軟に抽出、整形し、
    様々な形式の報告書やグラフを作成できます。
  - **選択項目 (RQLIST)**:
    `RQLIST`構造体は、シミュレーション結果から選択された項目を格納します。
    これにより、ユーザーが関心のある特定のデータのみを抽出できます。
  - **作表項目・期間指定 (PRQLIST)**:
    `PRQLIST`構造体は、作表項目や期間指定に関する情報を格納します。
    これにより、シミュレーション結果を特定の期間や項目に絞って分析できます。

このファイルは、建物のエネルギー性能シミュレーションの結果を効果的に分析し、
省エネルギー設計、快適性評価、
および運用改善のための意思決定を支援するための重要な役割を果たします。
*/
package eeslism

const VTYPEMAX = 50
const CATNMMAX = 50

type ESTL struct /* シミュレーション結果に関する注釈 */
{
	Flid     string   /* ファイル種別記号 */
	Title    string   /* 標題 */
	Wdatfile string   /* 気象データファイル名 */
	Tid      rune     /* 入力データ種別  h:時刻別  d:日別  M:月別 */
	Unit     []string /* 単位 */
	Timeid   string   /* 時刻データ表示  [Y]MD[W]T  *******/
	Wdloc    string   /* 地域情報　地名　緯度　経度など */
	Catnm    []CATNM

	Ntimeid        int /* 時刻データ表示字数 */
	Ntime          int /* 項目ごとの全データ数 */
	dtm            int /* 計算時間間隔[s] */
	Nunit          int /* 単位の定義数 */
	Nrqlist, Nvreq int
	Npreq, Npprd   int
	Ndata          int // VCFILEにおける各時刻ごとのデータ個数

	Rq   []RQLIST
	Prq  []PRQLIST
	Vreq []rune
}

type CATNM struct /* 要素カタログ名データ */
{
	Name   string
	N      int /* 機器数 */
	Ncdata int /* 全データ項目数 = 機器数 x 機器データ項目数 */
}

type TMDT struct /* 年、月、日、曜日、時刻データ */
{
	CYear  string     // 年(文字列)
	CMon   string     // 月(文字列)
	CDay   string     // 日(文字列)
	CWkday string     // 曜日(文字列)
	CTime  string     // 時刻(文字列)
	Dat    [5]*string /* 年、月、日、曜日、時刻のポインター */

	Year int // 年(数値)
	Mon  int // 月(数値)
	Day  int // 日 // 日(数値)
	Time int // 時刻(数値)
}

type TLIST struct /* シミュレーション結果 */
{
	Cname string
	Name  string // "Wd"以外はサポートされない
	Id    string /* データ項目名
	T:温度 x:絶対湿度 Idn:法線面直達日射量 Isky: 水平面天空日射量
	Ihor: 水平面全天日射量 CC:雲量 Wdre:風向 Wv:風速
	RH: 相対湿度 RN:夜間放射量
	*/
	Unit  string
	Vtype rune /* データ種別
					 t:温度  x:絶対湿度  r:相対湿度
					 T:平均温度  X:平均絶対湿度  R:平均相対湿度
					 h:発生時刻  H:積算時間
	q:熱量      Q:積算熱量   e:エネルギー E:積算エネルギー量 */
	Stype rune /* データ処理種別 t:積算値  a:平均値  n:最小値  m;最大値  */

	Ptype rune /* データ型  c:文字型  d:整数型  f:実数型 */
	Req   rune

	Fval, Fstat []float64
	Ival, Istat []int
	Cval, Cstat []rune
	Fmt         string

	Pair *TLIST
}

type RQLIST struct /* 選択項目 */
{
	Rname string
	Name  string
	Id    string
}

// type STATC struct /* 集計期間 */
// {
// 	Name     string
// 	Yrstart  int
// 	Mostart  int
// 	Daystart int
// 	Yrend    int
// 	Moend    int
// 	Dayend   int
// 	Nday     int
// 	Dymrk    [366]rune
// }

type PRQLIST struct /* 作表項目・期間指定 */
{
	Mark         rune
	Prname       []string
	Prid         []string
	Npname, Npid int
}

package main

const VTYPEMAX = 50
const CATNMMAX = 50

type ESTL struct /* シミュレーション結果に関する注釈 */
{
	Flid     string /* ファイル種別記号 */
	Title    string /* 標題 */
	Wdatfile string /* 気象データファイル名 */
	Tid      rune   /* 入力データ種別  h:時刻別  d:日別 */
	/* M:月別 　****/
	Unit   []string
	Timeid string /* 時刻データ表示  [Y]MD[W]T  *******/
	Wdloc  string /* 地域情報　地名　緯度　経度など */
	Catnm  []CATNM

	Ntimeid        int /* 時刻データ表示字数 */
	Ntime          int /* 項目ごとの全データ数 */
	dtm            int /* 計算時間間隔[s] */
	Nunit          int
	Nrqlist, Nvreq int
	Npreq, Npprd   int
	Ndata          int

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
	Day  int // 日(数値)
	Time int // 時刻(数値)
}

type TLIST struct /* シミュレーション結果 */
{
	Cname string
	Name  string
	Id    string
	Unit  string
	Vtype rune /* データ種別
					 t:温度  x:絶対湿度  r:相対湿度
					 T:平均温度  X:平均絶対湿度  R:平均相対湿度
					 h:発生時刻  H:積算時間
	q:熱量      Q:積算熱量   e:エネルギー E:積算エネルギー量 */
	Stype rune /* データ処理種別
	t:積算値  a:平均値  n:最小値  m;最大値  */

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

type STATC struct /* 集計期間 */
{
	Name     string
	Yrstart  int
	Mostart  int
	Daystart int
	Yrend    int
	Moend    int
	Dayend   int
	Nday     int
	Dymrk    [366]rune
}

type PRQLIST struct /* 作表項目・期間指定 */
{
	Mark         rune
	Prname       []string
	Prid         []string
	Npname, Npid int
}

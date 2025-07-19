//This file is part of EESLISM.
//
//Foobar is free software : you can redistribute itand /or modify
//it under the terms of the GNU General Public License as published by
//the Free Software Foundation, either version 3 of the License, or
//(at your option) any later version.
//
//Foobar is distributed in the hope that it will be useful,
//but WITHOUT ANY WARRANTY; without even the implied warranty of
//MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.See the
//GNU General Public License for more details.
//
//You should have received a copy of the GNU General Public License
//along with Foobar.If not, see < https://www.gnu.org/licenses/>.

/* roomday.c */

package eeslism

import (
	"fmt"
	"io"
)

var __Roomday_oldday = INAN
var __Roomday_oldMon = INAN

var __Rmdyprint_id int = 0
var __Rmmonprint_id int = 0

/*
Roomday (Room Daily and Monthly Data Aggregation)

この関数は、各室の熱的データ（室温、絶対湿度、相対湿度、平均表面温度、熱負荷など）を、
日単位および月単位で集計します。
これにより、室の熱的挙動、快適性、およびエネルギー消費量の詳細な分析が可能になります。

建築環境工学的な観点:
  - **日次集計**: 日次集計は、室の熱的挙動を日単位で詳細に把握するために重要です。
    例えば、日中のピーク負荷時の室温や湿度、
    あるいは暖房・冷房負荷の変動などを分析できます。
    これにより、日ごとの運用改善点を見つけ出すことが可能になります。
  - `Trdy`, `xrdy`, `RHdy`, `Tsavdy`: 室温、絶対湿度、相対湿度、平均表面温度の日平均値や最大・最小値。
  - `Qdys`, `Qdyl`, `Qdyt`: 顕熱、潜熱、全熱の日積算値。
  - `SQi`, `Tsdy`: 各表面からの熱量、表面温度の日積算値。
  - **月次集計**: 月次集計は、季節ごとの熱負荷変動や、
    室の年間を通じたエネルギー消費量の傾向を把握するために重要です。
    これにより、年間を通じた省エネルギー対策の効果を評価したり、
    熱負荷の予測精度を向上させたりすることが可能になります。
  - `mTrdy`, `mxrdy`, `mRHdy`, `mTsavdy`: 月平均値。
  - `mQdys`, `mQdyl`, `mQdyt`: 月積算値。
  - **太陽光発電の集計**: `Rdpnl`（放射パネル）に関連するデータとして、
    太陽電池パネルの温度（`TPVdy`, `mTPVdy`）や発電量（`PVdy`, `mPVdy`）も集計されます。
    これにより、太陽光発電システムの性能評価や、
    建物全体のエネルギー収支への貢献度を評価できます。
  - **データ分析の基礎**: この関数で集計されるデータは、
    室の熱的性能評価、熱負荷のベンチマーキング、
    省エネルギー対策の効果検証、そして運用改善のための意思決定の基礎となります。

この関数は、室の熱的挙動とエネルギー消費量を多角的に分析し、
快適性向上や省エネルギー対策の効果評価を行うための重要なデータ集計機能を提供します。
*/
func Roomday(Mon int, Day int, Nday int, ttmm int, Rm []*ROOM, Rdp []*RDPNL, Simdayend int) {
	Mo := Mon - 1
	tt := ConvertHour(ttmm)

	// 日集計
	if Nday != __Roomday_oldday {
		for i := range Rm {
			Room := Rm[i]

			svdyint(&Room.Trdy)
			svdyint(&Room.xrdy)
			svdyint(&Room.RHdy)
			svdyint(&Room.Tsavdy)

			R := Room.rmld
			if R != nil {
				qdyint(&R.Qdys)
				qdyint(&R.Qdyl)
				qdyint(&R.Qdyt)
			}

			for j := 0; j < Room.Nasup; j++ {
				A := Room.Arsp[j]
				qdyint(&A.Qdys)
				qdyint(&A.Qdyl)
				qdyint(&A.Qdyt)
			}

			for j := 0; j < Room.N; j++ {
				Sd := Room.rsrf[j]
				qdyint(&Sd.SQi)
				svdyint(&Sd.Tsdy)
			}
		}

		for i := range Rdp {
			Rdpnl := Rdp[i]
			svdyint(&Rdpnl.Tpody)
			svdyint(&Rdpnl.Tpidy)
			qdyint(&Rdpnl.Qdy)
			qdyint(&Rdpnl.Scoldy)
			svdyint(&Rdpnl.TPVdy)
			qdyint(&Rdpnl.PVdy)
		}

		__Roomday_oldday = Nday
	}

	// 月集計
	if Mon != __Roomday_oldMon {
		//printf("リセット\n") ;
		for i := range Rm {
			Room := Rm[i]

			svdyint(&Room.mTrdy)
			svdyint(&Room.mxrdy)
			svdyint(&Room.mRHdy)
			svdyint(&Room.mTsavdy)

			R := Room.rmld
			if R != nil {
				qdyint(&R.mQdys)
				qdyint(&R.mQdyl)
				qdyint(&R.mQdyt)
			}

			for j := 0; j < Room.Nasup; j++ {
				A := Room.Arsp[j]
				qdyint(&A.mQdys)
				qdyint(&A.mQdyl)
				qdyint(&A.mQdyt)
			}

			for j := 0; j < Room.N; j++ {
				Sd := Room.rsrf[j]
				qdyint(&Sd.mSQi)
				svdyint(&Sd.mTsdy)
			}
		}

		for i := range Rdp {
			Rdpnl := Rdp[i]
			svdyint(&Rdpnl.mTpody)
			svdyint(&Rdpnl.mTpidy)
			qdyint(&Rdpnl.mQdy)
			qdyint(&Rdpnl.mScoldy)
			svdyint(&Rdpnl.mTPVdy)
			qdyint(&Rdpnl.mPVdy)
		}

		__Roomday_oldMon = Mon
	}

	// 日集計
	for i := range Rm {
		Room := Rm[i]
		svdaysum(int64(ttmm), ON_SW, Room.Tr, &Room.Trdy)
		svdaysum(int64(ttmm), ON_SW, Room.xr, &Room.xrdy)
		svdaysum(int64(ttmm), ON_SW, Room.RH, &Room.RHdy)
		svdaysum(int64(ttmm), ON_SW, Room.Tsav, &Room.Tsavdy)

		R := Room.rmld
		if R != nil {
			qdaysum(int64(ttmm), ON_SW, R.Qs, &R.Qdys)
			qdaysum(int64(ttmm), ON_SW, R.Ql, &R.Qdyl)
			qdaysum(int64(ttmm), ON_SW, R.Qt, &R.Qdyt)
		}
		for j := 0; j < Room.Nasup; j++ {
			A := Room.Arsp[j]
			qdaysum(int64(ttmm), ON_SW, A.Qs, &A.Qdys)
			qdaysum(int64(ttmm), ON_SW, A.Ql, &A.Qdyl)
			qdaysum(int64(ttmm), ON_SW, A.Qt, &A.Qdyt)
		}

		for j := 0; j < Room.N; j++ {
			Sd := Room.rsrf[j]
			svdaysum(int64(ttmm), ON_SW, Sd.Ts, &Sd.Tsdy)
			qdaysum(int64(ttmm), ON_SW, Sd.Qi, &Sd.SQi)
		}
	}

	for i := range Rdp {
		Rdpnl := Rdp[i]

		svdaysum(int64(ttmm), Rdpnl.cmp.Control, Rdpnl.Tpo, &Rdpnl.Tpody)
		svdaysum(int64(ttmm), Rdpnl.cmp.Control, Rdpnl.Tpi, &Rdpnl.Tpidy)
		qdaysum(int64(ttmm), Rdpnl.cmp.Control, Rdpnl.Q, &Rdpnl.Qdy)
		qdaysumNotOpe(int64(ttmm), Rdpnl.sd[0].Iwall*Rdpnl.sd[0].A, &Rdpnl.Scoldy)

		control := OFF_SW
		if Rdpnl.sd[0].PVwall.Power > 0. {
			control = ON_SW
		}

		svdaysum(int64(ttmm), control, Rdpnl.sd[0].PVwall.TPV, &Rdpnl.TPVdy)
		qdaysumNotOpe(int64(ttmm), Rdpnl.sd[0].PVwall.Power, &Rdpnl.PVdy)
	}

	// 月集計
	//printf("Mon=%d Day=%d ttmm=%d\n", Mon, Day, ttmm ) ;
	for i := range Rm {
		Room := Rm[i]

		svmonsum(Mon, Day, ttmm, ON_SW, Room.Tr, &Room.mTrdy, Nday, Simdayend)
		svmonsum(Mon, Day, ttmm, ON_SW, Room.xr, &Room.mxrdy, Nday, Simdayend)
		svmonsum(Mon, Day, ttmm, ON_SW, Room.RH, &Room.mRHdy, Nday, Simdayend)
		svmonsum(Mon, Day, ttmm, ON_SW, Room.Tsav, &Room.mTsavdy, Nday, Simdayend)

		R := Room.rmld
		if R != nil {
			qmonsum(Mon, Day, ttmm, ON_SW, R.Qs, &R.mQdys, Nday, Simdayend)
			qmonsum(Mon, Day, ttmm, ON_SW, R.Ql, &R.mQdyl, Nday, Simdayend)
			qmonsum(Mon, Day, ttmm, ON_SW, R.Qt, &R.mQdyt, Nday, Simdayend)
		}
		for j := 0; j < Room.Nasup; j++ {
			A := Room.Arsp[j]
			qmonsum(Mon, Day, ttmm, ON_SW, A.Qs, &A.mQdys, Nday, Simdayend)
			qmonsum(Mon, Day, ttmm, ON_SW, A.Ql, &A.mQdyl, Nday, Simdayend)
			qmonsum(Mon, Day, ttmm, ON_SW, A.Qt, &A.mQdyt, Nday, Simdayend)
		}

		for j := 0; j < Room.N; j++ {
			Sd := Room.rsrf[j]
			svmonsum(Mon, Day, ttmm, ON_SW, Sd.Ts, &Sd.mTsdy, Nday, Simdayend)
			qmonsum(Mon, Day, ttmm, ON_SW, Sd.Qi, &Sd.mSQi, Nday, Simdayend)
		}
	}

	for i := range Rdp {
		Rdpnl := Rdp[i]

		svmonsum(Mon, Day, ttmm, Rdpnl.cmp.Control, Rdpnl.Tpo, &Rdpnl.mTpody, Nday, Simdayend)
		svmonsum(Mon, Day, ttmm, Rdpnl.cmp.Control, Rdpnl.Tpi, &Rdpnl.mTpidy, Nday, Simdayend)
		qmonsum(Mon, Day, ttmm, Rdpnl.cmp.Control, Rdpnl.Q, &Rdpnl.mQdy, Nday, Simdayend)
		qmonsumNotOpe(Mon, Day, ttmm, Rdpnl.sd[0].Iwall*Rdpnl.sd[0].A, &Rdpnl.mScoldy, Nday, Simdayend)

		control := OFF_SW
		if Rdpnl.sd[0].PVwall.Power > 0. {
			control = ON_SW
		}

		svmonsum(Mon, Day, ttmm, control, Rdpnl.sd[0].PVwall.TPV, &Rdpnl.mTPVdy, Nday, Simdayend)
		qmonsumNotOpe(Mon, Day, ttmm, Rdpnl.sd[0].PVwall.Power, &Rdpnl.mPVdy, Nday, Simdayend)

		// 月・時刻のクロス集計
		emtsum(Mon, Day, ttmm, control, Rdpnl.sd[0].PVwall.Power, &Rdpnl.mtPVdy[Mo][tt])
	}
}

/* ------------------------------------------------------- */

/*
Rmdyprint (Room Daily Output)

この関数は、各室の熱的データ（室温、絶対湿度、相対湿度、平均表面温度、熱負荷など）の
日積算値を整形して出力します。
これにより、日単位での室の熱的挙動を詳細に分析できます。

建築環境工学的な観点:
  - **日単位の熱的挙動の把握**: 日積算値を出力することで、
    日ごとの室温や湿度の変動、暖房・冷房負荷の推移などを把握できます。
    これにより、特定の日の熱負荷特性を分析したり、
    空調システムの運転状況を評価したりすることが可能になります。
  - **出力形式の制御**: `__Rmdyprint_id`によって出力形式を制御し、
    ヘッダー情報（`tttldyprint`）やカテゴリ情報（`-cat`）を出力します。
    これにより、出力データを解析ツールなどで利用しやすくなります。
  - **熱負荷の分類**: 出力されるデータには、
    室温、絶対湿度、相対湿度、平均表面温度の日平均値や最大・最小値、
    顕熱、潜熱、全熱の日積算値、
    各表面からの熱量、表面温度の日積算値などが含まれます。
    これにより、熱負荷の発生源を詳細に分析できます。
  - **放射パネルと太陽光発電のデータ**: 放射パネルからの熱量や、
    太陽電池パネルの温度・発電量の日積算値も出力されます。
    これにより、これらのシステムの性能評価や、
    室の熱負荷への貢献度を評価できます。

この関数は、室の熱的挙動とエネルギー消費量を日単位で詳細に分析し、
快適性向上や省エネルギー対策の効果評価を行うための重要なデータ出力機能を提供します。
*/
func Rmdyprint(fo io.Writer, mrk string, Simc *SIMCONTL, mon, day int, Rm []*ROOM) {
	if __Rmdyprint_id == 0 && len(Rm) > 0 {
		__Rmdyprint_id++

		ttldyprint(fo, mrk, Simc)
		fmt.Fprintf(fo, "-cat\n")
		fmt.Fprintf(fo, "%s %d\n", ROOM_TYPE, len(Rm))

		for i := range Rm {
			Room := Rm[i]

			var Nload int
			if Room.rmld != nil {
				Nload = 24
			} else {
				Nload = 0
			}

			fmt.Fprintf(fo, " %s 5 %d 24 %d %d %d\n", Room.Name,
				24+Nload+6*Room.Nasup+2*Room.Nrp,
				Nload, 6*Room.Nasup, 2*Room.Nrp)
		}
		fmt.Fprintf(fo, "*\n#\n")
	}

	if __Rmdyprint_id == 1 && len(Rm) > 0 {
		__Rmdyprint_id++

		for i := range Rm {
			Room := Rm[i]

			fmt.Fprintf(fo, "%s_Ht H d %s_Tr T f %s_ttn h d %s_Trn t f %s_ttm h d %s_Trm t f\n",
				Room.Name, Room.Name, Room.Name, Room.Name, Room.Name, Room.Name)
			fmt.Fprintf(fo, "%s_Hx H d %s_xr X f %s_txn h d %s_xrn x f %s_txm h d %s_xrm x f\n",
				Room.Name, Room.Name, Room.Name, Room.Name, Room.Name, Room.Name)
			fmt.Fprintf(fo, "%s_Hr H d %s_RH R f %s_trn h d %s_RHn r f %s_trm h d %s_RHm r f\n",
				Room.Name, Room.Name, Room.Name, Room.Name, Room.Name, Room.Name)
			fmt.Fprintf(fo, "%s_Hs H d %s_Ts T f %s_tsn h d %s_Tsn t f %s_tsm h d %s_Tsm t f\n",
				Room.Name, Room.Name, Room.Name, Room.Name, Room.Name, Room.Name)

			if Room.rmld != nil {
				fmt.Fprintf(fo, "%s_Hsh H d %s_Lsh Q f %s_Hsc H d %s_Lsc Q f ",
					Room.Name, Room.Name, Room.Name, Room.Name)
				fmt.Fprintf(fo, "%s_tsh h d %s_Lqsh q f %s_tsc h d %s_Lqsc q f\n",
					Room.Name, Room.Name, Room.Name, Room.Name)

				fmt.Fprintf(fo, "%s_Hlh H d %s_Llh Q f %s_Hlc H d %s_Llc Q f ",
					Room.Name, Room.Name, Room.Name, Room.Name)
				fmt.Fprintf(fo, "%s_tlh h d %s_Lqlh q f %s_tlc h d %s_Lqlc q f\n",
					Room.Name, Room.Name, Room.Name, Room.Name)

				fmt.Fprintf(fo, "%s_Hth H d %s_Lth Q f %s_Htc H d %s_Ltc Q f ",
					Room.Name, Room.Name, Room.Name, Room.Name)
				fmt.Fprintf(fo, "%s_tth h d %s_Lqth q f %s_ttc h d %s_Lqtc q f\n",
					Room.Name, Room.Name, Room.Name, Room.Name)
			}

			if Room.Nasup > 0 {
				for j := 0; j < Room.Nasup; j++ {
					Ei := Room.cmp.Elins[Room.Nachr+Room.Nrp+j]

					if Ei.Lpath == nil {
						fmt.Fprintf(fo, "%s:%d_Qash Q f %s:%d_Qasc Q f ",
							Room.Name, j, Room.Name, j)
						fmt.Fprintf(fo, "%s:%d_Qalh Q f %s:%d_Qalc Q f ",
							Room.Name, j, Room.Name, j)
						fmt.Fprintf(fo, "%s:%d_Qath Q f %s:%d_Qatc Q f\n",
							Room.Name, j, Room.Name, j)
					} else {
						fmt.Fprintf(fo, "%s:%s_Qash Q f %s:%s_Qasc Q f ",
							Room.Name, Ei.Lpath.Name, Room.Name, Ei.Lpath.Name)
						fmt.Fprintf(fo, "%s:%s_Qalh Q f %s:%s_Qalc Q f ",
							Room.Name, Ei.Lpath.Name, Room.Name, Ei.Lpath.Name)
						fmt.Fprintf(fo, "%s:%s_Qath Q f %s:%s_Qatc Q f\n",
							Room.Name, Ei.Lpath.Name, Room.Name, Ei.Lpath.Name)
					}
				}
			}
			for j := 0; j < Room.Nrp; j++ {
				rpnl := Room.rmpnl[j]
				fmt.Fprintf(fo, "%s:%s_Qh Q f %s:%s_Qc Q f ", Room.Name, rpnl.pnl.Name,
					Room.Name, rpnl.pnl.Name)
			}
			fmt.Fprintf(fo, "\n")
		}
	}

	fmt.Fprintf(fo, "%02d %02d\n", mon, day)

	for i := range Rm {
		Room := Rm[i]

		fmt.Fprintf(fo, "%1d %4.2f %1d %4.2f %1d %4.2f ",
			Room.Trdy.Hrs, Room.Trdy.M, Room.Trdy.Mntime,
			Room.Trdy.Mn, Room.Trdy.Mxtime, Room.Trdy.Mx)
		fmt.Fprintf(fo, "%1d %5.4f %1d %5.4f %1d %5.4f\n",
			Room.xrdy.Hrs, Room.xrdy.M, Room.xrdy.Mntime,
			Room.xrdy.Mn, Room.xrdy.Mxtime, Room.xrdy.Mx)
		fmt.Fprintf(fo, "%1d %2.0f %1d %2.0f %1d %2.0f ",
			Room.RHdy.Hrs, Room.RHdy.M, Room.RHdy.Mntime,
			Room.RHdy.Mn, Room.RHdy.Mxtime, Room.RHdy.Mx)
		fmt.Fprintf(fo, "%1d %4.2f %1d %4.2f %1d %4.2f\n",
			Room.Tsavdy.Hrs, Room.Tsavdy.M, Room.Tsavdy.Mntime,
			Room.Tsavdy.Mn, Room.Tsavdy.Mxtime, Room.Tsavdy.Mx)

		R := Room.rmld
		if R != nil {
			fmt.Fprintf(fo, "%1d %.2f ", R.Qdys.Hhr, R.Qdys.H)
			fmt.Fprintf(fo, "%1d %.2f ", R.Qdys.Chr, R.Qdys.C)
			fmt.Fprintf(fo, "%4d %2.0f ", R.Qdys.Hmxtime, R.Qdys.Hmx)
			fmt.Fprintf(fo, "%4d %2.0f ", R.Qdys.Cmxtime, R.Qdys.Cmx)

			fmt.Fprintf(fo, "%1d %.2f ", R.Qdyl.Hhr, R.Qdyl.H)
			fmt.Fprintf(fo, "%1d %.2f ", R.Qdyl.Chr, R.Qdyl.C)
			fmt.Fprintf(fo, "%4d %2.0f ", R.Qdyl.Hmxtime, R.Qdyl.Hmx)
			fmt.Fprintf(fo, "%4d %2.0f ", R.Qdyl.Cmxtime, R.Qdyl.Cmx)

			fmt.Fprintf(fo, "%1d %.2f ", R.Qdyt.Hhr, R.Qdyt.H)
			fmt.Fprintf(fo, "%1d %.2f ", R.Qdyt.Chr, R.Qdyt.C)
			fmt.Fprintf(fo, "%4d %2.0f ", R.Qdyt.Hmxtime, R.Qdyt.Hmx)
			fmt.Fprintf(fo, "%4d %2.0f\n", R.Qdyt.Cmxtime, R.Qdyt.Cmx)
		}
		if Room.Nasup > 0 {
			for j := 0; j < Room.Nasup; j++ {
				A := Room.Arsp[j]
				fmt.Fprintf(fo, "%3.1f %.2f ", A.Qdys.H, A.Qdys.C)
				fmt.Fprintf(fo, "%3.1f %.2f ", A.Qdyl.H, A.Qdyl.C)
				fmt.Fprintf(fo, "%3.1f %.2f ", A.Qdyt.H, A.Qdyt.C)
			}
			fmt.Fprintf(fo, "\n")
		}
		if Room.Nrp > 0 {
			for j := 0; j < Room.Nrp; j++ {
				rpnl := Room.rmpnl[j]
				fmt.Fprintf(fo, "%.2f %.2f\n", -rpnl.pnl.Qdy.C, -rpnl.pnl.Qdy.H)
			}
		}
		fmt.Fprintf(fo, "\n")
	}
}

/*
Rmmonprint (Room Monthly Output)

この関数は、各室の熱的データ（室温、絶対湿度、相対湿度、平均表面温度、熱負荷など）の
月積算値を整形して出力します。
これにより、月単位での室の熱的挙動を詳細に分析できます。

建築環境工学的な観点:
  - **月単位の熱的挙動の把握**: 月積算値を出力することで、
    月ごとの室温や湿度の変動、暖房・冷房負荷の推移などを把握できます。
    これにより、特定の月の熱負荷特性を分析したり、
    空調システムの運転状況を評価したりすることが可能になります。
  - **出力形式の制御**: `__Rmmonprint_id`によって出力形式を制御し、
    ヘッダー情報（`tttldyprint`）やカテゴリ情報（`-cat`）を出力します。
    これにより、出力データを解析ツールなどで利用しやすくなります。
  - **熱負荷の分類**: 出力されるデータには、
    室温、絶対湿度、相対湿度、平均表面温度の月平均値や最大・最小値、
    顕熱、潜熱、全熱の月積算値、
    各表面からの熱量、表面温度の月積算値などが含まれます。
    これにより、熱負荷の発生源を詳細に分析できます。
  - **放射パネルと太陽光発電のデータ**: 放射パネルからの熱量や、
    太陽電池パネルの温度・発電量の月積算値も出力されます。
    これにより、これらのシステムの性能評価や、
    室の熱負荷への貢献度を評価できます。

この関数は、室の熱的挙動とエネルギー消費量を月単位で詳細に分析し、
快適性向上や省エネルギー対策の効果評価を行うための重要なデータ出力機能を提供します。
*/
func Rmmonprint(fo io.Writer, mrk string, Simc *SIMCONTL, mon, day int, Rm []*ROOM) {

	Nroom := len(Rm)

	if __Rmmonprint_id == 0 && Nroom > 0 {
		__Rmmonprint_id++

		ttldyprint(fo, mrk, Simc)
		fmt.Fprintf(fo, "-cat\n")
		fmt.Fprintf(fo, "%s %d\n", ROOM_TYPE, Nroom)

		for i := 0; i < Nroom; i++ {
			Room := Rm[i]

			var Nload int
			if Room.rmld != nil {
				Nload = 24
			} else {
				Nload = 0
			}
			fmt.Fprintf(fo, " %s 5 %d 24 %d %d %d\n", Room.Name,
				24+Nload+6*Room.Nasup+2*Room.Nrp,
				Nload, 6*Room.Nasup, 2*Room.Nrp)
		}
		fmt.Fprintf(fo, "*\n#\n")
	}

	if __Rmmonprint_id == 1 && Nroom > 0 {
		__Rmmonprint_id++

		for i := 0; i < Nroom; i++ {
			Room := Rm[i]

			fmt.Fprintf(fo, "%s_Ht H d %s_Tr T f %s_ttn h d %s_Trn t f %s_ttm h d %s_Trm t f\n",
				Room.Name, Room.Name, Room.Name, Room.Name, Room.Name, Room.Name)
			fmt.Fprintf(fo, "%s_Hx H d %s_xr X f %s_txn h d %s_xrn x f %s_txm h d %s_xrm x f\n",
				Room.Name, Room.Name, Room.Name, Room.Name, Room.Name, Room.Name)
			fmt.Fprintf(fo, "%s_Hr H d %s_RH R f %s_trn h d %s_RHn r f %s_trm h d %s_RHm r f\n",
				Room.Name, Room.Name, Room.Name, Room.Name, Room.Name, Room.Name)
			fmt.Fprintf(fo, "%s_Hs H d %s_Ts T f %s_tsn h d %s_Tsn t f %s_tsm h d %s_Tsm t f\n",
				Room.Name, Room.Name, Room.Name, Room.Name, Room.Name, Room.Name)

			if Room.rmld != nil {
				fmt.Fprintf(fo, "%s_Hsh H d %s_Lsh Q f %s_Hsc H d %s_Lsc Q f ",
					Room.Name, Room.Name, Room.Name, Room.Name)
				fmt.Fprintf(fo, "%s_tsh h d %s_Lqsh q f %s_tsc h d %s_Lqsc q f\n",
					Room.Name, Room.Name, Room.Name, Room.Name)

				fmt.Fprintf(fo, "%s_Hlh H d %s_Llh Q f %s_Hlc H d %s_Llc Q f ",
					Room.Name, Room.Name, Room.Name, Room.Name)
				fmt.Fprintf(fo, "%s_tlh h d %s_Lqlh q f %s_tlc h d %s_Lqlc q f\n",
					Room.Name, Room.Name, Room.Name, Room.Name)

				fmt.Fprintf(fo, "%s_Hth H d %s_Lth Q f %s_Htc H d %s_Ltc Q f ",
					Room.Name, Room.Name, Room.Name, Room.Name)
				fmt.Fprintf(fo, "%s_tth h d %s_Lqth q f %s_ttc h d %s_Lqtc q f\n",
					Room.Name, Room.Name, Room.Name, Room.Name)
			}

			if Room.Nasup > 0 {
				for j := 0; j < Room.Nasup; j++ {
					Ei := Room.cmp.Elins[Room.Nachr+Room.Nrp+j]

					if Ei.Lpath == nil {
						fmt.Fprintf(fo, "%s:%d_Qash Q f %s:%d_Qasc Q f ",
							Room.Name, j, Room.Name, j)
						fmt.Fprintf(fo, "%s:%d_Qalh Q f %s:%d_Qalc Q f ",
							Room.Name, j, Room.Name, j)
						fmt.Fprintf(fo, "%s:%d_Qath Q f %s:%d_Qatc Q f\n",
							Room.Name, j, Room.Name, j)
					} else {
						fmt.Fprintf(fo, "%s:%s_Qash Q f %s:%s_Qasc Q f ",
							Room.Name, Ei.Lpath.Name, Room.Name, Ei.Lpath.Name)
						fmt.Fprintf(fo, "%s:%s_Qalh Q f %s:%s_Qalc Q f ",
							Room.Name, Ei.Lpath.Name, Room.Name, Ei.Lpath.Name)
						fmt.Fprintf(fo, "%s:%s_Qath Q f %s:%s_Qatc Q f\n",
							Room.Name, Ei.Lpath.Name, Room.Name, Ei.Lpath.Name)
					}
				}
			}
			for j := 0; j < Room.Nrp; j++ {
				rpnl := Room.rmpnl[j]
				fmt.Fprintf(fo, "%s:%s_Qh Q f %s:%s_Qc Q f ", Room.Name, rpnl.pnl.Name,
					Room.Name, rpnl.pnl.Name)
			}
			fmt.Fprintf(fo, "\n")
		}
	}

	fmt.Fprintf(fo, "%02d %02d\n", mon, day)

	for i := 0; i < Nroom; i++ {
		Room := Rm[i]

		fmt.Fprintf(fo, "%1d %4.2f %1d %4.2f %1d %4.2f ",
			Room.mTrdy.Hrs, Room.mTrdy.M, Room.mTrdy.Mntime,
			Room.mTrdy.Mn, Room.mTrdy.Mxtime, Room.mTrdy.Mx)
		fmt.Fprintf(fo, "%1d %5.4f %1d %5.4f %1d %5.4f\n",
			Room.mxrdy.Hrs, Room.mxrdy.M, Room.mxrdy.Mntime,
			Room.mxrdy.Mn, Room.mxrdy.Mxtime, Room.mxrdy.Mx)
		fmt.Fprintf(fo, "%1d %2.0f %1d %2.0f %1d %2.0f ",
			Room.mRHdy.Hrs, Room.mRHdy.M, Room.mRHdy.Mntime,
			Room.mRHdy.Mn, Room.mRHdy.Mxtime, Room.mRHdy.Mx)
		fmt.Fprintf(fo, "%1d %4.2f %1d %4.2f %1d %4.2f\n",
			Room.mTsavdy.Hrs, Room.mTsavdy.M, Room.mTsavdy.Mntime,
			Room.mTsavdy.Mn, Room.mTsavdy.Mxtime, Room.mTsavdy.Mx)

		R := Room.rmld
		if R != nil {
			fmt.Fprintf(fo, "%1d %3.1f ", R.mQdys.Hhr, R.mQdys.H)
			fmt.Fprintf(fo, "%1d %3.1f ", R.mQdys.Chr, R.mQdys.C)
			fmt.Fprintf(fo, "%4d %2.0f ", R.mQdys.Hmxtime, R.mQdys.Hmx)
			fmt.Fprintf(fo, "%4d %2.0f ", R.mQdys.Cmxtime, R.mQdys.Cmx)

			fmt.Fprintf(fo, "%1d %3.1f ", R.mQdyl.Hhr, R.mQdyl.H)
			fmt.Fprintf(fo, "%1d %3.1f ", R.mQdyl.Chr, R.mQdyl.C)
			fmt.Fprintf(fo, "%4d %2.0f ", R.mQdyl.Hmxtime, R.mQdyl.Hmx)
			fmt.Fprintf(fo, "%4d %2.0f ", R.mQdyl.Cmxtime, R.mQdyl.Cmx)

			fmt.Fprintf(fo, "%1d %3.1f ", R.mQdyt.Hhr, R.mQdyt.H)
			fmt.Fprintf(fo, "%1d %3.1f ", R.mQdyt.Chr, R.mQdyt.C)
			fmt.Fprintf(fo, "%4d %2.0f ", R.mQdyt.Hmxtime, R.mQdyt.Hmx)
			fmt.Fprintf(fo, "%4d %2.0f\n", R.mQdyt.Cmxtime, R.mQdyt.Cmx)
		}
		if Room.Nasup > 0 {
			for j := 0; j < Room.Nasup; j++ {
				A := Room.Arsp[j]
				fmt.Fprintf(fo, "%3.1f %3.1f ", A.mQdys.H, A.mQdys.C)
				fmt.Fprintf(fo, "%3.1f %3.1f ", A.mQdyl.H, A.mQdyl.C)
				fmt.Fprintf(fo, "%3.1f %3.1f ", A.mQdyt.H, A.mQdyt.C)
			}
			fmt.Fprintf(fo, "\n")
		}
		if Room.Nrp > 0 {
			for j := 0; j < Room.Nrp; j++ {
				rpnl := Room.rmpnl[j]
				fmt.Fprintf(fo, "%3.1f %3.1f\n", -rpnl.pnl.mQdy.C, -rpnl.pnl.mQdy.H)
			}
		}
		fmt.Fprintf(fo, "\n")
	}
}

/*
paneldyprt (Panel Daily Output)

この関数は、放射パネル（床暖房など）および太陽電池パネルの
日積算値を整形して出力します。
これにより、日単位でのパネルの熱的挙動や発電性能を詳細に分析できます。

建築環境工学的な観点:
  - **放射パネルの性能評価**: 放射パネルは、
    輻射熱によって室内を暖めたり冷やしたりするシステムです。
    日積算値を出力することで、日ごとの熱供給量や熱除去量、
    パネル表面温度の推移などを把握できます。
    これにより、放射パネルの快適性や省エネルギー効果を評価できます。
  - **太陽電池パネルの性能評価**: 太陽電池パネルは、
    太陽光エネルギーを電力に変換するシステムです。
    日積算値を出力することで、日ごとの発電量、
    パネル表面温度の推移などを把握できます。
    これにより、太陽光発電システムの性能や、
    建物全体のエネルギー収支への貢献度を評価できます。
  - **出力形式の制御**: `id`によって出力形式を制御し、
    パネルの種類（床暖房パネル、太陽電池一体型など）に応じた適切な項目を出力します。
    これにより、出力データを解析ツールなどで利用しやすくなります。

この関数は、放射パネルおよび太陽電池パネルの熱的挙動や発電性能を日単位で詳細に分析し、
快適性向上や省エネルギー対策の効果評価を行うための重要なデータ出力機能を提供します。
*/
func paneldyprt(fo io.Writer, id int, _Rdpnl []*RDPNL) {
	switch id {
	case 0:
		if len(_Rdpnl) > 0 {
			fmt.Fprintf(fo, "%s %d\n", RDPANEL_TYPE, len(_Rdpnl))
		}

		for i := range _Rdpnl {
			Rdpnl := _Rdpnl[i]
			Wall := Rdpnl.sd[0].mw.wall

			if Wall.WallType == WallType_P {
				// 床暖房パネルの場合
				fmt.Fprintf(fo, " %s 1 20\n", Rdpnl.Name)
			} else if Rdpnl.sd[0].PVwallFlg {
				//太陽電池一体型の場合
				fmt.Fprintf(fo, " %s 1 36\n", Rdpnl.Name)
			} else {
				// その他
				fmt.Fprintf(fo, " %s 1 28\n", Rdpnl.Name)
			}
		}
		break

	case 1:
		for i := range _Rdpnl {
			Rdpnl := _Rdpnl[i]
			Wall := Rdpnl.sd[0].mw.wall

			fmt.Fprintf(fo, "%s_Hto H d %s_To T f ", Rdpnl.Name, Rdpnl.Name)
			fmt.Fprintf(fo, "%s_ttno h d %s_Ton t f %s_ttm h d %s_Tom t f ",
				Rdpnl.Name, Rdpnl.Name, Rdpnl.Name, Rdpnl.Name)
			fmt.Fprintf(fo, "%s_Ht H d %s_T T f ", Rdpnl.Name, Rdpnl.Name)
			fmt.Fprintf(fo, "%s_ttn h d %s_Tn t f %s_ttm h d %s_Tm t f ",
				Rdpnl.Name, Rdpnl.Name, Rdpnl.Name, Rdpnl.Name)
			fmt.Fprintf(fo, "%s_Hh H d %s_Qh Q f %s_Hc H d %s_Qc Q f",
				Rdpnl.Name, Rdpnl.Name, Rdpnl.Name, Rdpnl.Name)
			fmt.Fprintf(fo, "%s_th h d %s_qh q f %s_tc h d %s_qc q f\n",
				Rdpnl.Name, Rdpnl.Name, Rdpnl.Name, Rdpnl.Name)

			if Wall.WallType == WallType_C {
				fmt.Fprintf(fo, "%s_ScolHh H d %s_ScolQh Q f %s_ScolHc H d %s_ScolQc Q f",
					Rdpnl.Name, Rdpnl.Name, Rdpnl.Name, Rdpnl.Name)
				fmt.Fprintf(fo, "%s_Scolth h d %s_Scolqh q f %s_Scoltc h d %s_Scolqc q f\n",
					Rdpnl.Name, Rdpnl.Name, Rdpnl.Name, Rdpnl.Name)

				if Rdpnl.sd[0].PVwallFlg {
					fmt.Fprintf(fo, "%s_PVHt H d %s_TPV T f ", Rdpnl.Name, Rdpnl.Name)
					fmt.Fprintf(fo, "%s_PVttn h d %s_TPVn t f %s_PVttm h d %s_TPVm t f ",
						Rdpnl.Name, Rdpnl.Name, Rdpnl.Name, Rdpnl.Name)
					fmt.Fprintf(fo, "%s_PVH h d %s_E E f\n", Rdpnl.Name, Rdpnl.Name)
				}
			}
		}
		break

	default:
		for i := range _Rdpnl {
			Rdpnl := _Rdpnl[i]
			Wall := Rdpnl.sd[0].mw.wall

			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ",
				Rdpnl.Tpody.Hrs, Rdpnl.Tpody.M, Rdpnl.Tpody.Mntime,
				Rdpnl.Tpody.Mn, Rdpnl.Tpody.Mxtime, Rdpnl.Tpody.Mx)
			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ",
				Rdpnl.Tpidy.Hrs, Rdpnl.Tpidy.M, Rdpnl.Tpidy.Mntime,
				Rdpnl.Tpidy.Mn, Rdpnl.Tpidy.Mxtime, Rdpnl.Tpidy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f ", Rdpnl.Qdy.Hhr, Rdpnl.Qdy.H)
			fmt.Fprintf(fo, "%1d %3.1f ", Rdpnl.Qdy.Chr, Rdpnl.Qdy.C)
			fmt.Fprintf(fo, "%1d %2.0f ", Rdpnl.Qdy.Hmxtime, Rdpnl.Qdy.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f\n", Rdpnl.Qdy.Cmxtime, Rdpnl.Qdy.Cmx)

			if Wall.WallType == WallType_C {
				fmt.Fprintf(fo, "%1d %3.1f ", Rdpnl.Scoldy.Hhr, Rdpnl.Scoldy.H)
				fmt.Fprintf(fo, "%1d %3.1f ", Rdpnl.Scoldy.Chr, Rdpnl.Scoldy.C)
				fmt.Fprintf(fo, "%1d %2.0f ", Rdpnl.Scoldy.Hmxtime, Rdpnl.Scoldy.Hmx)
				fmt.Fprintf(fo, "%1d %2.0f\n", Rdpnl.Scoldy.Cmxtime, Rdpnl.Scoldy.Cmx)

				if Rdpnl.sd[0].PVwallFlg {
					fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ",
						Rdpnl.TPVdy.Hrs, Rdpnl.TPVdy.M, Rdpnl.TPVdy.Mntime,
						Rdpnl.TPVdy.Mn, Rdpnl.TPVdy.Mxtime, Rdpnl.TPVdy.Mx)
					fmt.Fprintf(fo, "%1d %.1f\n", Rdpnl.PVdy.Hhr, Rdpnl.PVdy.H)
				}
			}
		}
		break
	}
}

/*
panelmonprt (Panel Monthly Output)

この関数は、放射パネル（床暖房など）および太陽電池パネルの
月積算値を整形して出力します。
これにより、月単位でのパネルの熱的挙動や発電性能を詳細に分析できます。

建築環境工学的な観点:
  - **放射パネルの性能評価**: 放射パネルは、
    輻射熱によって室内を暖めたり冷やしたりするシステムです。
    月積算値を出力することで、月ごとの熱供給量や熱除去量、
    パネル表面温度の推移などを把握できます。
    これにより、放射パネルの快適性や省エネルギー効果を評価できます。
  - **太陽電池パネルの性能評価**: 太陽電池パネルは、
    太陽光エネルギーを電力に変換するシステムです。
    月積算値を出力することで、月ごとの発電量、
    パネル表面温度の推移などを把握できます。
    これにより、太陽光発電システムの性能や、
    建物全体のエネルギー収支への貢献度を評価できます。
  - **出力形式の制御**: `id`によって出力形式を制御し、
    パネルの種類（床暖房パネル、太陽電池一体型など）に応じた適切な項目を出力します。
    これにより、出力データを解析ツールなどで利用しやすくなります。

この関数は、放射パネルおよび太陽電池パネルの熱的挙動や発電性能を月単位で詳細に分析し、
快適性向上や省エネルギー対策の効果評価を行うための重要なデータ出力機能を提供します。
*/
func panelmonprt(fo io.Writer, id int, _Rdpnl []*RDPNL) {

	switch id {
	case 0:
		if len(_Rdpnl) > 0 {
			fmt.Fprintf(fo, "%s %d\n", RDPANEL_TYPE, len(_Rdpnl))
		}

		for i := range _Rdpnl {
			Rdpnl := _Rdpnl[i]
			Wall := Rdpnl.sd[0].mw.wall

			if Wall.WallType == WallType_P {
				fmt.Fprintf(fo, " %s 1 20\n", Rdpnl.Name)
			} else {
				if Rdpnl.sd[0].PVwallFlg {
					fmt.Fprintf(fo, " %s 1 36\n", Rdpnl.Name)
				} else {
					fmt.Fprintf(fo, " %s 1 28\n", Rdpnl.Name)
				}
			}
		}
		break

	case 1:
		for i := range _Rdpnl {
			Rdpnl := _Rdpnl[i]
			Wall := Rdpnl.sd[0].mw.wall

			fmt.Fprintf(fo, "%s_Hto H d %s_To T f ", Rdpnl.Name, Rdpnl.Name)
			fmt.Fprintf(fo, "%s_ttno h d %s_Ton t f %s_ttm h d %s_Tom t f ",
				Rdpnl.Name, Rdpnl.Name, Rdpnl.Name, Rdpnl.Name)
			fmt.Fprintf(fo, "%s_Ht H d %s_T T f ", Rdpnl.Name, Rdpnl.Name)
			fmt.Fprintf(fo, "%s_ttn h d %s_Tn t f %s_ttm h d %s_Tm t f ",
				Rdpnl.Name, Rdpnl.Name, Rdpnl.Name, Rdpnl.Name)
			fmt.Fprintf(fo, "%s_Hh H d %s_Qh Q f %s_Hc H d %s_Qc Q f",
				Rdpnl.Name, Rdpnl.Name, Rdpnl.Name, Rdpnl.Name)
			fmt.Fprintf(fo, "%s_th h d %s_qh q f %s_tc h d %s_qc q f\n",
				Rdpnl.Name, Rdpnl.Name, Rdpnl.Name, Rdpnl.Name)

			if Wall.WallType == WallType_C {
				fmt.Fprintf(fo, "%s_ScolHh H d %s_ScolQh Q f %s_ScolHc H d %s_ScolQc Q f",
					Rdpnl.Name, Rdpnl.Name, Rdpnl.Name, Rdpnl.Name)
				fmt.Fprintf(fo, "%s_Scolth h d %s_Scolqh q f %s_Scoltc h d %s_Scolqc q f\n",
					Rdpnl.Name, Rdpnl.Name, Rdpnl.Name, Rdpnl.Name)

				if Rdpnl.sd[0].PVwallFlg {
					fmt.Fprintf(fo, "%s_PVHt H d %s_TPV T f ", Rdpnl.Name, Rdpnl.Name)
					fmt.Fprintf(fo, "%s_PVttn h d %s_TPVn t f %s_PVttm h d %s_TPVm t f ",
						Rdpnl.Name, Rdpnl.Name, Rdpnl.Name, Rdpnl.Name)
					fmt.Fprintf(fo, "%s_PVH h d %s_E E f\n", Rdpnl.Name, Rdpnl.Name)
				}
			}
		}
		break

	default:
		for i := range _Rdpnl {
			Rdpnl := _Rdpnl[i]
			Wall := Rdpnl.sd[0].mw.wall

			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ",
				Rdpnl.mTpody.Hrs, Rdpnl.mTpody.M, Rdpnl.mTpody.Mntime,
				Rdpnl.mTpody.Mn, Rdpnl.mTpody.Mxtime, Rdpnl.mTpody.Mx)
			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ",
				Rdpnl.mTpidy.Hrs, Rdpnl.mTpidy.M, Rdpnl.mTpidy.Mntime,
				Rdpnl.mTpidy.Mn, Rdpnl.mTpidy.Mxtime, Rdpnl.mTpidy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f ", Rdpnl.mQdy.Hhr, Rdpnl.mQdy.H)
			fmt.Fprintf(fo, "%1d %3.1f ", Rdpnl.mQdy.Chr, Rdpnl.mQdy.C)
			fmt.Fprintf(fo, "%1d %2.0f ", Rdpnl.mQdy.Hmxtime, Rdpnl.mQdy.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f\n", Rdpnl.mQdy.Cmxtime, Rdpnl.mQdy.Cmx)

			if Wall.WallType == WallType_C {
				fmt.Fprintf(fo, "%1d %3.1f ", Rdpnl.mScoldy.Hhr, Rdpnl.mScoldy.H)
				fmt.Fprintf(fo, "%1d %3.1f ", Rdpnl.mScoldy.Chr, Rdpnl.mScoldy.C)
				fmt.Fprintf(fo, "%1d %2.0f ", Rdpnl.mScoldy.Hmxtime, Rdpnl.mScoldy.Hmx)
				fmt.Fprintf(fo, "%1d %2.0f\n", Rdpnl.mScoldy.Cmxtime, Rdpnl.mScoldy.Cmx)

				if Rdpnl.sd[0].PVwallFlg {
					fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ",
						Rdpnl.mTPVdy.Hrs, Rdpnl.mTPVdy.M, Rdpnl.mTPVdy.Mntime,
						Rdpnl.mTPVdy.Mn, Rdpnl.mTPVdy.Mxtime, Rdpnl.mTPVdy.Mx)
					fmt.Fprintf(fo, "%1d %.1f\n", Rdpnl.mPVdy.Hhr, Rdpnl.mPVdy.H)
				}
			}
		}
		break
	}
}

func panelmtprt(fo io.Writer, id int, Rdpnl []*RDPNL, Mo int, tt int) {
	switch id {
	case 0:
		if len(Rdpnl) > 0 {
			fmt.Fprintf(fo, "%s %d\n", RDPANEL_TYPE, len(Rdpnl))
		}
		for i := range Rdpnl {
			Rdpnl := Rdpnl[i]
			fmt.Fprintf(fo, " %s 1 1\n", Rdpnl.Name)
		}
	case 1:
		for i := range Rdpnl {
			Rdpnl := Rdpnl[i]
			fmt.Fprintf(fo, "%s_E E f \n", Rdpnl.Name)
		}
	default:
		for i := range Rdpnl {
			Rdpnl := Rdpnl[i]
			fmt.Fprintf(fo, " %.2f \n", Rdpnl.mtPVdy[Mo-1][tt-1].D*Cff_kWh)
		}
	}
}

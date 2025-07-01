package eeslism

import (
	"fmt"
)

var __Eeprintd_ic int = 0
var __Wdtdprint_ic int = 0
var __Wdtprint_ic int = 0
var __Wdtmprint_ic int = 0

/*
Eeprinth (Energy Simulation Hourly Print)

この関数は、建物のエネルギーシミュレーションにおける各時間ステップの計算結果を、
時刻ごとの時系列データとして出力します。
これにより、建物の熱的挙動の動的な変化を詳細に分析できます。

建築環境工学的な観点:
- **時系列データの重要性**: 建物の熱負荷は、日射、外気温度、内部発熱など、
  様々な要因によって刻々と変化します。
  時系列データとして出力することで、
  - **ピーク負荷の把握**: 一日のうちで最も熱負荷が高くなる時間帯を特定し、
    空調設備の容量設計に役立てることができます。
  - **熱的挙動の分析**: 日射の侵入による室温上昇、
    夜間の放熱による室温低下など、
    建物の熱的挙動の動的な変化を詳細に分析できます。
  - **運用改善の検討**: 実際の運用におけるエネルギー消費量と、
    シミュレーション結果を比較することで、
    運用改善のためのヒントを得ることができます。
- **出力内容の選択**: `flo.Idn`（出力タイプ識別子）によって、
  気象データ（`PRTHWD`）、機器の運転データ（`PRTCOMP`）、
  システム経路の温湿度（`PRTPATH`）、蓄熱槽内温度分布（`PRTHRSTANK`）、
  室温・MRT（`PRTREV`）、放射パネル（`PRTHROOM`）、
  要素別熱損失・熱取得（`PRTHELM`, `PRTHELMSF`）、
  PMV（`PRTPMV`）、日射・室内熱取得（`PRTQRM`）、
  室内表面温度（`PRTRSF`）、日よけの影面積（`PRTSHD`）、
  壁体内部温度（`PRTWAL`）、潜熱蓄熱材の状態値（`PRTPCM`）、
  室内表面熱流（`PRTSFQ`）、室内表面熱伝達率（`PRTSFA`）など、
  多岐にわたるシミュレーション結果を選択的に出力できます。
  これにより、ユーザーは必要な情報を効率的に取得し、
  分析や検証を容易にします。

この関数は、建物の熱的挙動を詳細に分析し、
空調システムの設計、運用、省エネルギー対策の検討、
および快適性評価を行うための重要なデータ出力機能を提供します。
*/
func Eeprinth(Daytm *DAYTM, Simc *SIMCONTL, flout []*FLOUT, Rmvls *RMVLS, Exsfst *EXSFS, Mpath []*MPATH, Eqsys *EQSYS, Wd *WDAT) {
	if Daytm.Ddpri != 0 && Simc.Dayprn[Daytm.DayOfYear] != 0 {
		title := Simc.Title
		Mon := Daytm.Mon
		Day := Daytm.Day
		time := Daytm.Time

		for i, flo := range flout {

			if DEBUG {
				fmt.Printf("Eeprinth MAX=%d flo[%d]=%s\n", len(flout), i, flo.Idn)
			}

			switch flo.Idn {
			case PRTHWD:
				if DEBUG {
					fmt.Println("<Eeprinth> xprsolrd")
				}
				// 気象データの出力
				Wdtprint(flo.F, title, Mon, Day, time, Wd, Exsfst)
			case PRTCOMP: // 毎時機器の出力
				Hcmpprint(flo.F, string(PRTCOMP), Simc, Mon, Day, time, Eqsys, Rmvls.Rdpnl)
			case PRTPATH: // システム経路の温湿度出力
				Pathprint(flo.F, title, Mon, Day, time, Mpath)
			case PRTHRSTANK: // 蓄熱槽内温度分布の出力
				Hstkprint(flo.F, title, Mon, Day, time, Eqsys)
			default:
				if SIMUL_BUILDG { // these blocks are only compiled in debug builds
					switch flo.Idn {
					case PRTREV:
						// 毎時室温、MRTの出力
						Rmevprint(flo.F, title, Rmvls.Room, Mon, Day, time)
					case PRTHROOM:
						// 放射パネルの出力
						Rmpnlprint(flo.F, string(PRTHROOM), Simc, Mon, Day, time, Rmvls.Room)
					case PRTHELM:
						// 要素別熱損失・熱取得
						Helmprint(flo.F, string(PRTHELM), Simc, Mon, Day, time, Rmvls.Room, &Rmvls.Qetotal)
					case PRTHELMSF:
						// 要素別熱損失・熱取得
						Helmsurfprint(flo.F, string(PRTHELMSF), Simc, Mon, Day, time, Rmvls.Room)
					case PRTPMV:
						// PMV計算
						Pmvprint(flo.F, title, Rmvls.Room, Mon, Day, time)
					case PRTQRM:
						// 日射、室内熱取得の出力
						Qrmprint(flo.F, title, Mon, Day, time, Rmvls.Room, Rmvls.Qrm)
					case PRTRSF:
						// 室内表面温度の出力
						Rmsfprint(flo.F, title, Mon, Day, time, Rmvls.Room, Rmvls.Sd)
					case PRTSHD:
						// 日よけの影面積の出力
						Shdprint(flo.F, title, Mon, Day, time, Rmvls.Sd)
					case PRTWAL:
						// 壁体内部温度の出力
						Wallprint(flo.F, title, Mon, Day, time, Rmvls.Sd)
					case PRTPCM:
						// 潜熱蓄熱材の状態値の出力
						PCMprint(flo.F, title, Mon, Day, time, Rmvls.Sd)
					case PRTSFQ:
						// 室内表面熱流の出力
						Rmsfqprint(flo.F, title, Mon, Day, time, Rmvls.Room, Rmvls.Sd)
					case PRTSFA:
						// 室内表面熱伝達率の出力
						Rmsfaprint(flo.F, title, Mon, Day, time, Rmvls.Room, Rmvls.Sd)
					}
				}
			}
		}
	}
}

/*
Eeprintd (Energy Simulation Daily Print)

この関数は、建物のエネルギーシミュレーションにおける日ごとの集計結果を整形して出力します。
これにより、日単位での建物の熱的挙動やエネルギー消費量を詳細に分析できます。

建築環境工学的な観点:
- **日単位のエネルギー評価**: 建物のエネルギー消費量は、日単位で変動します。
  日積算値を出力することで、日ごとの熱負荷変動や、
  各要素からの熱損失・熱取得の割合を把握できます。
  これにより、特定の日のエネルギー消費が多かった原因を分析したり、
  省エネルギー対策の効果を日単位で評価したりすることが可能になります。
- **出力内容の選択**: `flo.Idn`（出力タイプ識別子）によって、
  気象データ日集計値（`PRTDWD`）、計算年月日（`PRTWK`）、
  システム要素機器の日集計結果（`PRTDYCOMP`）、
  部屋ごとの熱集計結果（`PRTDYRM`）、
  要素別熱損失・熱取得の日積算値（`PRTDYHELM`）、
  日射・室内熱取得（`PRTDQR`）、
  日積算壁体貫流熱取得（`PRTDYSF`）など、
  多岐にわたるシミュレーション結果を選択的に出力できます。
  これにより、ユーザーは必要な情報を効率的に取得し、
  分析や検証を容易にします。

この関数は、建物の熱的挙動を詳細に分析し、
空調システムの設計、運用、省エネルギー対策の検討、
および快適性評価を行うための重要なデータ出力機能を提供します。
*/
func Eeprintd(Daytm *DAYTM, Simc *SIMCONTL, flout []*FLOUT, Rmvls *RMVLS, Exs []*EXSF, Soldy []float64, Eqsys *EQSYS, Wdd *WDAT) {
	if Daytm.Ddpri != 0 {
		title := Simc.Title
		Mon := int(Daytm.Mon)
		Day := int(Daytm.Day)

		for _, flo := range flout {
			switch flo.Idn {
			case PRTDWD:
				// 気象データ日集計値出力
				Wdtdprint(flo.F, title, Mon, Day, Wdd, Exs, Soldy)
			case PRTWK:
				// 計算年月日出力
				if __Eeprintd_ic == 0 {
					fmt.Fprintf(flo.F, "Mo Nd Day Week\n")
					__Eeprintd_ic = 1
				}

				fmt.Fprintf(flo.F, "%2d %2d %3d %s\n", Mon, Day, Daytm.DayOfYear, DAYweek[Simc.Daywk[Daytm.DayOfYear]])
			case PRTDYCOMP:
				// システム要素機器の日集計結果出力
				Compodyprt(flo.F, string(PRTDYCOMP), Simc, Mon, Day, Eqsys, Rmvls.Rdpnl)
			case PRTDYRM:
				// 部屋ごとの熱集計結果出力
				Rmdyprint(flo.F, string(PRTDYRM), Simc, Mon, Day, Rmvls.Room)
			case PRTDYHELM:
				// 要素別熱損失・熱取得（日積算値出力）
				Helmdyprint(flo.F, string(PRTDYHELM), Simc, Mon, Day, Rmvls.Room, &Rmvls.Qetotal)
			case PRTDQR:
				// 日射、室内熱取得の出力
				Dyqrmprint(flo.F, title, Mon, Day, Rmvls.Room, Rmvls.Trdav, Rmvls.Qrmd)
			case PRTDYSF:
				// 日積算壁体貫流熱取得の出力
				Dysfprint(flo.F, title, Mon, Day, Rmvls.Room)
			}
		}
	}
}

/*
Eeprintm (Energy Simulation Monthly Print)

この関数は、建物のエネルギーシミュレーションにおける月ごとの集計結果を整形して出力します。
これにより、月単位での建物の熱的挙動やエネルギー消費量を詳細に分析できます。

建築環境工学的な観点:
- **月単位のエネルギー評価**: 建物のエネルギー消費量は、月単位で変動します。
  月積算値を出力することで、月ごとの熱負荷変動や、
  各要素からの熱損失・熱取得の割合を把握できます。
  これにより、特定の月のエネルギー消費が多かった原因を分析したり、
  省エネルギー対策の効果を月単位で評価したりすることが可能になります。
- **出力内容の選択**: `flo.Idn`（出力タイプ識別子）によって、
  気象データ月集計値（`PRTMWD`）、システム要素機器の月集計結果（`PRTMNCOMP`）、
  部屋ごとの熱集計結果（`PRTMNRM`）など、
  多岐にわたるシミュレーション結果を選択的に出力できます。
  これにより、ユーザーは必要な情報を効率的に取得し、
  分析や検証を容易にします。

この関数は、建物の熱的挙動を詳細に分析し、
空調システムの設計、運用、省エネルギー対策の検討、
および快適性評価を行うための重要なデータ出力機能を提供します。
*/
func Eeprintm(daytm *DAYTM, simc *SIMCONTL, flout []*FLOUT, rmvls *RMVLS, exs []*EXSF, solmon []float64, eqsys *EQSYS, wdm *WDAT) {
	var title string
	var mon, day int
	title = simc.Title
	mon = daytm.Mon
	day = daytm.Day
	if daytm.Ddpri != 0 {
		for _, flo := range flout {
			switch flo.Idn {
			case PRTMWD:
				// 気象データ月集計値出力
				Wdtmprint(flo.F, title, mon, day, wdm, exs, solmon)
			case PRTMNCOMP:
				// システム要素機器の月集計結果出力
				Compomonprt(flo.F, string(PRTMNCOMP), simc, mon, day, eqsys, rmvls.Rdpnl)
			case PRTMNRM:
				// 部屋ごとの熱集計結果出力
				Rmmonprint(flo.F, string(PRTMNRM), simc, mon, day, rmvls.Room)
			}
		}
	}
}

/*
Eeprintmt (Energy Simulation Monthly-Time-of-Day Print)

この関数は、建物のエネルギーシミュレーションにおける月・時刻別の集計結果を整形して出力します。
これにより、月ごとの時間帯別エネルギー消費量を詳細に分析できます。

建築環境工学的な観点:
- **月・時刻別のエネルギー評価**: シミュレーションの各時間ステップで計算されたエネルギー量を、
  月と時刻の組み合わせで集計します。
  これにより、特定の月における時間帯ごとのエネルギー消費量の傾向を把握できます。
- **デマンドサイドマネジメント**: 月・時刻別のエネルギー消費量データは、
  デマンドサイドマネジメント（DSM）戦略を策定する上で非常に有用です。
  例えば、ピーク時間帯の電力消費量を削減するための運転戦略を検討したり、
  蓄熱システムや再生可能エネルギーの導入効果を評価したりする際に役立ちます。
- **出力内容の選択**: `flo.Idn == PRTMTCOMP` の条件によって、
  月・時刻別の機器の集計結果（`Compomtprt`）を出力します。
  これにより、ユーザーは必要な情報を効率的に取得し、
  分析や検証を容易にします。

この関数は、建物のエネルギー消費量を月・時刻別で詳細に分析し、
運用改善や省エネルギー対策の効果評価を行うための重要なデータ出力機能を提供します。
*/
func Eeprintmt(simc *SIMCONTL, flout []*FLOUT, eqsys *EQSYS, rdpnl []*RDPNL) {
	for _, flo := range flout {
		if flo.Idn == PRTMTCOMP {
			Compomtprt(flo.F, string(PRTMNCOMP), simc, eqsys, rdpnl)
		}
	}
}

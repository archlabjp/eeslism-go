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

/* helm.c */

package eeslism

import (
	"fmt"
	"io"
)

var __Helmprint_id int = 0
var __Helmsurfprint_id int = 0
var __Helmdy_oldday int = -1
var __Helmdyprint_id int = 0

/*
Helminit (Building Element Heat Loss/Gain Initialization)

この関数は、建物の各要素（壁、屋根、床、窓など）からの熱損失・熱取得を計算するための
データ構造を初期化します。
これは、建物の熱負荷を詳細に分析し、省エネルギー対策の効果を評価するために不可欠です。

建築環境工学的な観点:
- **要素別熱負荷の分析**: 建物の熱負荷は、外皮（壁、屋根、窓など）からの熱損失・熱取得、
  内部発熱、換気など、様々な要因によって構成されます。
  この関数は、各要素からの熱損失・熱取得を個別に計算するためのデータ構造（`RMQELM`, `RMSB`）を準備します。
  これにより、どの要素が熱負荷に最も寄与しているかを特定し、
  効果的な省エネルギー対策を検討できます。
- **壁体内部温度の記憶**: `Rs.Tw`と`Rs.Told`は、
  壁体内部の温度分布を記憶するためのデータ構造です。
  壁体は熱容量を持つため、その内部温度は時間遅れを伴って変化します。
  これらの温度を追跡することで、壁体の蓄熱効果や熱貫流特性を正確にモデル化できます。
- **要素タイプの識別**: `Sd.ble`（建物の要素タイプ）や`Sd.typ`（表面タイプ）に基づいて、
  各要素が外部に面しているか（`RMSBType_E`, `RMSBType_G`）、
  内部に面しているか（`RMSBType_i`）を識別します。
  これにより、各要素の熱的特性に応じた適切な計算ロジックが適用されます。
- **熱損失・熱取得の分類**: `Qetotal`は、建物全体の熱損失・熱取得を統合的に集計するための構造体です。
  これにより、建物全体のエネルギー性能を評価できます。

この関数は、建物の熱負荷を詳細に分析し、
省エネルギー対策の検討、熱的性能の評価、
および快適性評価を行うための重要な初期設定機能を提供します。
*/
func Helminit(errkey string, helmkey rune, _Room []*ROOM, Qetotal *QETOTAL) {
	var Nmax, k int

	if helmkey != 'y' {
		for i := range _Room {
			Room := _Room[i]
			Room.rmqe = nil
		}
		return
	}

	for i := range _Room {
		Room := _Room[i]

		Room.rmqe = &RMQELM{}

		if Room.rmqe != nil {
			Rq := Room.rmqe
			Rq.rmsb = nil
			Rq.WSCwk = nil
		}

		N := Room.N
		if N > 0 {
			Room.rmqe.rmsb = make([]*RMSB, N)
		}

		if Room.rmqe.rmsb != nil {
			for k = 0; k < N; k++ {
				Rs := Room.rmqe.rmsb[k]
				Rs.Told = nil
				Rs.Tw = nil
			}
		}

		for j := 0; j < Room.N; j++ {
			Sd := Room.rsrf[j]
			Rs := Room.rmqe.rmsb[j]

			if Sd.mw != nil {
				N := Sd.mw.M
				if N > 0 {
					Rs.Tw = make([]*BHELM, N)
					Rs.Told = make([]*BHELM, N)
				}
			} else {
				Rs.Tw = nil
				Rs.Told = nil
			}

			switch Sd.ble {
			case BLE_ExternalWall, BLE_Roof, BLE_Floor, BLE_Window:
				if Sd.typ != RMSRFType_E && Sd.typ != RMSRFType_e {
					Rs.Type = RMSBType_E
				} else {
					Rs.Type = RMSBType_G
				}
				break
			case BLE_InnerWall, BLE_InnerFloor, BLE_Ceil, BLE_d:
				Rs.Type = RMSBType_i
				break
			}
		}
		if Room.N > Nmax {
			Nmax = Room.N
		}
	}

	for i := range _Room {
		Room := _Room[i]
		if i == 0 {
			if Nmax > 0 {
				Room.rmqe.WSCwk = make([]*BHELM, Nmax)

				Bh := Room.rmqe.WSCwk[0]
				Bh.trs = 0.0
				Bh.so = 0.0
				Bh.sg = 0.0
				Bh.rn = 0.0
				Bh.in = 0.0
				Bh.pnl = 0.0
			}
		} else {
			Room.rmqe.WSCwk = _Room[0].rmqe.WSCwk
		}
	}
	Qetotal.Name = "Qetotal"
}

/*
Helmroom (Building Element Heat Loss/Gain Calculation for Rooms)

この関数は、各室における建物の要素別熱損失・熱取得を計算します。
これは、室ごとの熱負荷を詳細に分析し、
空調システムの設計や運用、省エネルギー対策の検討に不可欠です。

建築環境工学的な観点:
- **室ごとの熱負荷分析**: 建物の熱負荷は、室ごとに異なる特性を持ちます。
  この関数は、各室の熱負荷を構成する要素（日射熱取得、壁体からの熱伝達、内部発熱など）を個別に計算し、
  `Rm.rmqe.qelm`に格納します。
  これにより、室ごとの熱負荷の要因を特定し、
  適切な空調ゾーン設定や、室ごとの省エネルギー対策を検討できます。
- **熱損失・熱取得の計算 (helmrmsrt, helmq)**:
  - `helmrmsrt(Rm, Ta)`: 室内の表面温度を計算し、
    表面からの熱伝達（対流、放射）を評価します。
  - `helmq(Room, Ta, xa)`: 室内の熱収支を計算し、
    日射熱取得、内部発熱、換気による熱損失・熱取得などを評価します。
- **熱損失・熱取得の集計**: 各室で計算された要素別熱損失・熱取得は、
  `qelmsum`関数によって建物全体の熱損失・熱取得（`Qetotal.Qelm`）に集計されます。
  これにより、建物全体のエネルギー性能を評価できます。
- **壁体内部温度の更新 (helmwall)**:
  `helmwall(Rm, Ta)`関数は、壁体内部の温度を更新します。
  壁体は熱容量を持つため、その内部温度は時間遅れを伴って変化します。
  これらの温度を追跡することで、壁体の蓄熱効果や熱貫流特性を正確にモデル化できます。

この関数は、建物の熱負荷を詳細に分析し、
空調システムの設計、運用、省エネルギー対策の検討、
および快適性評価を行うための重要な役割を果たします。
*/
func Helmroom(Room []*ROOM, Qrm []*QRM, Qetotal *QETOTAL, Ta, xa float64) {
	qelmclear(&Qetotal.Qelm)

	for i := range Room {
		Rm := Room[i]
		Qr := Qrm[i]
		qe := &Rm.rmqe.qelm

		helmrmsrt(Rm, Ta)
		helmq(Room, Ta, xa)

		qe.slo = Qr.Solo
		qe.slw = Qr.Solw
		qe.asl = Qr.Asl
		qe.tsol = Qr.Tsol
		qe.hins = Qr.Hgins

		qelmsum(qe, &Qetotal.Qelm)
	}

	for i := range Room {
		Rm := Room[i]
		helmwall(Rm, Ta)
	}
}

/*
Helmprint (Building Element Heat Loss/Gain Time-Series Output)

この関数は、建物の要素別熱損失・熱取得の計算結果を、
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
- **出力形式の制御**: `__Helmprint_id`によって出力形式を制御し、
  ヘッダー情報（`ttlprint`）やカテゴリ情報（`-cat`）を出力します。
  これにより、出力データを解析ツールなどで利用しやすくなります。
- **室ごとの詳細データ**: `helmrmprint`関数を呼び出すことで、
  各室ごとの要素別熱損失・熱取得データを出力します。
  これにより、室ごとの熱的特性を詳細に分析できます。

この関数は、建物の熱的挙動を詳細に分析し、
空調システムの設計、運用、省エネルギー対策の検討、
および快適性評価を行うための重要なデータ出力機能を提供します。
*/
func Helmprint(fo io.Writer, mrk string, Simc *SIMCONTL, mon, day int, time float64,
	Room []*ROOM, Qetotal *QETOTAL) {
	var j int

	if __Helmprint_id == 0 {
		ttlprint(fo, mrk, Simc)

		for j = 0; j < 2; j++ {
			if j == 0 {
				fmt.Fprintf(fo, "-cat\n")
			}
			helmrmprint(fo, __Helmprint_id, Room, Qetotal)
			if j == 0 {
				fmt.Fprintf(fo, "*\n#\n")
			}
			__Helmprint_id++
		}
	}

	fmt.Fprintf(fo, "%02d %02d %5.2f\n", mon, day, time)
	helmrmprint(fo, __Helmprint_id, Room, Qetotal)
}

/* ----------------------------------------------------- */

/*
helmrmprint (Building Element Heat Loss/Gain Room-wise Output)

この関数は、各室および建物全体の要素別熱損失・熱取得の計算結果を整形して出力します。
これにより、室ごとの熱負荷の構成要素を詳細に把握できます。

建築環境工学的な観点:
- **室ごとの熱負荷構成**: 建物の熱負荷は、外皮（壁、屋根、窓など）からの熱損失・熱取得、
  内部発熱、換気など、様々な要因によって構成されます。
  この関数は、各室の熱負荷を構成するこれらの要素を個別に表示することで、
  どの要素が熱負荷に最も寄与しているかを特定し、
  効果的な省エネルギー対策を検討できます。
- **熱損失・熱取得の分類**: 出力されるデータには、
  - `qldh`, `qldc`: 暖房・冷房負荷
  - `slo`, `slw`: 太陽光による顕熱・潜熱取得
  - `asl`, `tsol`: 透過日射による熱取得
  - `hins`: 内部発熱
  - `so`, `sg`, `rn`, `in`, `pnl`: 表面からの熱伝達（日射、地中、放射、侵入、パネル）
  - `trs`: 透過熱損失
  - `qew`, `qwn`, `qgd`, `qnx`: 外壁、窓、地盤、隣室からの熱伝達
  - `qi`, `qc`, `qf`: 内部発熱、換気、ファンによる熱取得
  - `vo`, `vr`, `sto`: 換気量、換気による熱回収、蓄熱量
  など、多岐にわたる熱損失・熱取得の項目が含まれています。
  これにより、熱負荷の発生源を詳細に分析できます。
- **建物全体の集計**: 各室のデータに加えて、
  `Qetotal.Qelm`として建物全体の熱損失・熱取得の合計も出力されます。
  これにより、個別の室の熱負荷と建物全体のエネルギー性能を関連付けて評価できます。

この関数は、建物の熱負荷を詳細に分析し、
空調システムの設計、運用、省エネルギー対策の検討、
および快適性評価を行うための重要なデータ出力機能を提供します。
*/
func helmrmprint(fo io.Writer, id int, _Room []*ROOM, Qetotal *QETOTAL) {
	var q *BHELM
	var qh *QHELM
	var name string

	Nroom := len(_Room)

	switch id {
	case 0:
		if Nroom > 0 {
			fmt.Fprintf(fo, "%s %d\n", ROOM_TYPE, Nroom+1)
		}

		for i := 0; i < Nroom; i++ {
			Room := _Room[i]
			if Room.rmqe != nil {
				fmt.Fprintf(fo, "%s 1 %d\n", Room.Name, 29)
			}
		}
		fmt.Fprintf(fo, "%s 1 %d\n", Qetotal.Name, 29)
		break

	case 1:
		for i := 0; i < Nroom+1; i++ {
			if i < Nroom {
				name = _Room[i].Name
			} else {
				name = Qetotal.Name
			}

			fmt.Fprintf(fo, "%s_qldh q f %s_qldc q f ", name, name)
			fmt.Fprintf(fo, "%s_slo q f %s_slw q f %s_asl q f %s_tsol q f %s_hins q f\n",
				name, name, name, name, name)
			fmt.Fprintf(fo, "%s_so q f %s_sw q f %s_rn q f %s_in q f %s_pnl q f\n",
				name, name, name, name, name)
			fmt.Fprintf(fo, "%s_trs q f %s_qew q f %s_qwn q f %s_qgd q f %s_qnx q f ",
				name, name, name, name, name)
			fmt.Fprintf(fo, "%s_qi q f %s_qc q f %s_qf q f\n",
				name, name, name)
			fmt.Fprintf(fo, "%s_vo q f %s_vr q f %s_sto q f\n", name, name, name)
			fmt.Fprintf(fo, "%s_qldhl q f %s_qldcl q f %s_hinl q f\n", name, name, name)
			fmt.Fprintf(fo, "%s_vol q f %s_vrl q f %s_stol q f\n", name, name, name)
		}
		break

	default:
		for i := 0; i < Nroom+1; i++ {
			if i < Nroom {
				Room := _Room[i]
				q = &(Room.rmqe.qelm.qe)
				qh = &Room.rmqe.qelm

				fmt.Fprintf(fo, "%3.0f %3.0f ", qh.loadh, qh.loadc)
				fmt.Fprintf(fo, "%3.0f %3.0f %3.0f %3.0f %3.0f ",
					qh.slo, qh.slw, qh.asl, qh.tsol, qh.hins)

				fmt.Fprintf(fo, "%3.0f %3.0f %3.0f %3.0f %3.0f ",
					q.so, q.sg, q.rn, q.in, q.pnl)
				fmt.Fprintf(fo, "%3.0f %3.0f %3.0f %3.0f %3.0f ",
					q.trs, qh.ew, qh.wn, qh.gd, qh.nx)
				fmt.Fprintf(fo, "%3.0f %3.0f %3.0f %3.0f %3.0f %3.0f ",
					qh.i, qh.c, qh.f, qh.vo, qh.vr, qh.sto)
				fmt.Fprintf(fo, "%3.0f %3.0f %3.0f %3.0f %3.0f %3.0f\n",
					qh.loadhl, qh.loadcl, qh.hinl, qh.vol, qh.vrl, qh.stol)
			} else {
				q = &Qetotal.Qelm.qe
				qh = &Qetotal.Qelm
				fmt.Fprintf(fo, "%3.0f %3.0f ", qh.loadh, qh.loadc)
				fmt.Fprintf(fo, "%3.0f %3.0f %3.0f %3.0f %3.0f ",
					qh.slo, qh.slw, qh.asl, qh.tsol, qh.hins)

				fmt.Fprintf(fo, "%3.0f %3.0f %3.0f %3.0f %3.0f ",
					q.so, q.sg, q.rn, q.in, q.pnl)
				fmt.Fprintf(fo, "%3.0f %3.0f %3.0f %3.0f %3.0f ",
					q.trs, qh.ew, qh.wn, qh.gd, qh.nx)
				fmt.Fprintf(fo, "%3.0f %3.0f %3.0f %3.0f %3.0f %3.0f ",
					qh.i, qh.c, qh.f, qh.vo, qh.vr, qh.sto)
				fmt.Fprintf(fo, "%3.0f %3.0f %3.0f %3.0f %3.0f %3.0f\n",
					qh.loadhl, qh.loadcl, qh.hinl, qh.vol, qh.vrl, qh.stol)
			}
		}
		break
	}
}

/* ----------------------------------------------------- */

/*
Helmsurfprint (Building Element Surface Heat Loss/Gain Time-Series Output)

この関数は、建物の各表面（壁、窓など）における熱損失・熱取得の計算結果を、
時刻ごとの時系列データとして出力します。
これにより、各表面の熱的挙動の動的な変化を詳細に分析できます。

建築環境工学的な観点:
- **表面熱収支の分析**: 建物の表面は、日射、外気温度、室内温度、放射熱交換など、
  様々な要因によって熱を交換します。
  この関数は、各表面からの熱損失・熱取得を個別に計算し、
  時系列データとして出力することで、
  - **日射熱取得の評価**: 窓や壁面が受ける日射熱が、
    どのように表面温度や熱取得量に影響するかを詳細に分析できます。
  - **結露の発生予測**: 表面温度が露点温度を下回るかどうかを監視することで、
    結露の発生リスクを評価できます。
  - **快適性評価**: 表面温度は、居住者の放射快適性に直接影響します。
    このデータは、快適性評価の基礎となります。
- **出力形式の制御**: `__Helmsurfprint_id`によって出力形式を制御し、
  ヘッダー情報（`ttlprint`）やカテゴリ情報（`-cat`）を出力します。
  これにより、出力データを解析ツールなどで利用しやすくなります。
- **室ごとの詳細データ**: `helmsfprint`関数を呼び出すことで、
  各室ごとの表面熱損失・熱取得データを出力します。
  これにより、室ごとの熱的特性を詳細に分析できます。

この関数は、建物の熱的挙動を詳細に分析し、
空調システムの設計、運用、省エネルギー対策の検討、
および快適性評価を行うための重要なデータ出力機能を提供します。
*/
func Helmsurfprint(fo io.Writer, mrk string, Simc *SIMCONTL, mon, day int, time float64, Room []*ROOM) {
	var j int

	if __Helmsurfprint_id == 0 {
		ttlprint(fo, mrk, Simc)

		for j = 0; j < 2; j++ {
			if j == 0 {
				fmt.Fprintf(fo, "-cat\n")
			}
			helmsfprint(fo, __Helmsurfprint_id, Room)
			if j == 0 {
				fmt.Fprintf(fo, "*\n#\n")
			}
			__Helmsurfprint_id++
		}
	}

	fmt.Fprintf(fo, "%02d %02d %5.2f\n", mon, day, time)
	helmsfprint(fo, __Helmsurfprint_id, Room)
}

/* ----------------------------------------------------- */

/*
helmsfprint (Building Element Surface Heat Loss/Gain Surface-wise Output)

この関数は、各室の各表面（壁、窓など）における熱損失・熱取得の計算結果を整形して出力します。
これにより、各表面の熱的挙動を詳細に把握できます。

建築環境工学的な観点:
- **表面熱収支の構成**: 建物の表面は、日射、外気温度、室内温度、放射熱交換など、
  様々な要因によって熱を交換します。
  この関数は、各表面からの熱損失・熱取得を構成する要素を個別に表示することで、
  どの表面が熱負荷に最も寄与しているかを特定し、
  効果的な省エネルギー対策を検討できます。
- **熱損失・熱取得の分類**: 出力されるデータには、
  - `trs`: 透過熱損失
  - `so`: 太陽光による熱取得
  - `sg`: 地盤からの熱取得
  - `rn`: 夜間放射による熱損失
  - `in`: 侵入空気による熱損失
  - `pnl`: パネルからの熱取得
  など、多岐にわたる熱損失・熱取得の項目が含まれています。
  これにより、熱負荷の発生源を詳細に分析できます。
- **表面ごとの詳細データ**: `Sd.sfepri`が`true`の場合にのみ出力されることで、
  特定の表面に絞って詳細な分析を行うことができます。
  これにより、例えば、日射熱取得が大きい窓や、熱損失が大きい壁などを特定し、
  集中的な対策を検討できます。

この関数は、建物の熱的挙動を詳細に分析し、
空調システムの設計、運用、省エネルギー対策の検討、
および快適性評価を行うための重要なデータ出力機能を提供します。
*/
func helmsfprint(fo io.Writer, id int, _Room []*ROOM) {
	switch id {
	case 0:
		if len(_Room) > 0 {
			fmt.Fprintf(fo, "%s %d\n", ROOM_TYPE, len(_Room))
		}

		for i := range _Room {
			Room := _Room[i]
			Nsf := 0
			for j := 0; j < Room.N; j++ {
				Sd := Room.rsrf[j]
				if Sd.sfepri {
					Nsf++
				}
			}
			fmt.Fprintf(fo, "%s 1 %d\n", Room.Name, 6*Nsf)

		}
		break

	case 1:
		for i := range _Room {
			Room := _Room[i]
			for j := 0; j < Room.N; j++ {
				Sd := Room.rsrf[j]
				if Sd.sfepri {
					var s string
					if len(Sd.Name) == 0 {
						s = fmt.Sprintf(s, "%s-%d-%c", Room.Name, j, Sd.ble)
					} else {
						s = fmt.Sprintf(s, "%s-%s", Room.Name, Sd.Name)
					}

					fmt.Fprintf(fo, "%s_trs t f %s_so f %s_sg t f ", s, s, s)
					fmt.Fprintf(fo, "%s_rn t f %s_in t f %s_pnl t f\n", s, s, s)
				}
			}
		}
		break

	default:
		for i := range _Room {
			Room := _Room[i]
			for j := 0; j < Room.N; j++ {
				Sd := Room.rsrf[j]
				rmsb := Room.rmqe.rmsb[j]
				if Sd.sfepri {
					Ts := &rmsb.Ts
					fmt.Fprintf(fo, "%5.2f %5.2f %5.2f ", Ts.trs, Ts.so, Ts.sg)
					fmt.Fprintf(fo, "%5.2f %5.2f %5.2f\n", Ts.rn, Ts.in, Ts.pnl)
				}
			}

		}
		break
	}
}

/* ----------------------------------------------------- */

/*
Helmdy (Building Element Heat Loss/Gain Daily Aggregation)

この関数は、建物の要素別熱損失・熱取得の日積算値を集計します。
これは、日単位での建物のエネルギー消費量を評価し、
省エネルギー対策の効果を分析するために用いられます。

建築環境工学的な観点:
- **日単位のエネルギー評価**: 建物のエネルギー消費量は、日単位で変動します。
  日積算値を集計することで、日ごとの熱負荷変動や、
  各要素からの熱損失・熱取得の割合を把握できます。
  これにより、特定の日のエネルギー消費が多かった原因を分析したり、
  省エネルギー対策の効果を日単位で評価したりすることが可能になります。
- **運用改善の指標**: 日積算データは、建物の運用改善のための重要な指標となります。
  例えば、休日や夜間のエネルギー消費量が過剰でないかを確認したり、
  外気温度や日射量などの気象条件とエネルギー消費量の関係を分析したりすることで、
  より効率的な運用方法を見つけることができます。
- **データ集計の準備**: `if day != __Helmdy_oldday` の条件は、
  新しい日になった場合に日積算値をリセットする（`helmdyint`）ことを意味します。
  これにより、日ごとの正確な集計が可能になります。
- **建物全体の集計**: 各室で計算された要素別熱損失・熱取得の日積算値は、
  `qelmsum`関数によって建物全体の熱損失・熱取得の日積算値（`Qetotal.Qelmdy`）に集計されます。
  これにより、建物全体のエネルギー性能を評価できます。

この関数は、建物のエネルギー消費量を日単位で詳細に分析し、
運用改善や省エネルギー対策の効果評価を行うための基礎的な役割を果たします。
*/
func Helmdy(day int, Room []*ROOM, Qetotal *QETOTAL) {
	if day != __Helmdy_oldday {
		helmdyint(Room, Qetotal)
		__Helmdy_oldday = day
	}

	for i := range Room {
		rmq := Room[i].rmqe

		if rmq != nil {
			qelmsum(&rmq.qelm, &rmq.qelmdy)
		}
	}

	qelmsum(&Qetotal.Qelm, &Qetotal.Qelmdy)
}

/* ----------------------------------------------------- */

/*
helmdyint (Building Element Heat Loss/Gain Daily Integration Initialization)

この関数は、建物の要素別熱損失・熱取得の日積算値をリセットします。
これは、新しい日の集計を開始する前に、
前日のデータをクリアするために用いられます。

建築環境工学的な観点:
- **日単位の集計の準備**: 建物のエネルギー消費量を日単位で評価するためには、
  各日の開始時に集計値をゼロにリセットする必要があります。
  この関数は、`qelmclear`関数を呼び出すことで、
  各室および建物全体の要素別熱損失・熱取得の日積算値を初期化します。
- **正確なデータ分析の確保**: 日積算値が適切にリセットされることで、
  日ごとのエネルギー消費量を正確に比較分析することが可能になります。
  これにより、特定の日のエネルギー消費が多かった原因を特定したり、
  省エネルギー対策の効果を日単位で評価したりする際の信頼性が向上します。

この関数は、建物のエネルギー消費量を日単位で詳細に分析し、
運用改善や省エネルギー対策の効果評価を行うための基礎的な役割を果たします。
*/
func helmdyint(Room []*ROOM, Qetotal *QETOTAL) {
	for i := range Room {
		if Room[i].rmqe != nil {
			qelmclear(&Room[i].rmqe.qelmdy)
		}
	}

	qelmclear(&Qetotal.Qelmdy)
}

/* ----------------------------------------------------- */

/*
Helmdyprint (Building Element Heat Loss/Gain Daily Output)

この関数は、建物の要素別熱損失・熱取得の日積算値を整形して出力します。
これにより、日単位での建物の熱的挙動を詳細に分析できます。

建築環境工学的な観点:
- **日単位のエネルギー評価**: 建物のエネルギー消費量は、日単位で変動します。
  日積算値を出力することで、日ごとの熱負荷変動や、
  各要素からの熱損失・熱取得の割合を把握できます。
  これにより、特定の日のエネルギー消費が多かった原因を分析したり、
  省エネルギー対策の効果を日単位で評価したりすることが可能になります。
- **出力形式の制御**: `__Helmdyprint_id`によって出力形式を制御し、
  ヘッダー情報（`tttldyprint`）やカテゴリ情報（`-cat`）を出力します。
  これにより、出力データを解析ツールなどで利用しやすくなります。
- **室ごとの詳細データ**: `helmrmdyprint`関数を呼び出すことで、
  各室ごとの要素別熱損失・熱取得の日積算データを出力します。
  これにより、室ごとの熱的特性を詳細に分析できます。

この関数は、建物の熱的挙動を詳細に分析し、
空調システムの設計、運用、省エネルギー対策の検討、
および快適性評価を行うための重要なデータ出力機能を提供します。
*/
func Helmdyprint(fo io.Writer, mrk string, Simc *SIMCONTL, mon, day int, Room []*ROOM, Qetotal *QETOTAL) {
	var j int

	if __Helmdyprint_id == 0 {
		ttldyprint(fo, mrk, Simc)

		for j = 0; j < 2; j++ {
			if j == 0 {
				fmt.Fprintf(fo, "-cat\n")
			}
			helmrmdyprint(fo, __Helmdyprint_id, Room, Qetotal)
			if j == 0 {
				fmt.Fprintf(fo, "*\n#\n")
			}
			__Helmdyprint_id++
		}
	}

	fmt.Fprintf(fo, "%02d %02d\n", mon, day)
	helmrmdyprint(fo, __Helmdyprint_id, Room, Qetotal)
}

/* ----------------------------------------------------- */

/*
helmrmdyprint (Building Element Heat Loss/Gain Room-wise Daily Output)

この関数は、各室および建物全体の要素別熱損失・熱取得の日積算値を整形して出力します。
これにより、日単位での室ごとの熱負荷の構成要素を詳細に把握できます。

建築環境工学的な観点:
- **日単位の熱負荷構成**: 建物の熱負荷は、外皮（壁、屋根、窓など）からの熱損失・熱取得、
  内部発熱、換気など、様々な要因によって構成されます。
  この関数は、各室の熱負荷を構成するこれらの要素の日積算値を個別に表示することで、
  どの要素が熱負荷に最も寄与しているかを特定し、
  効果的な省エネルギー対策を検討できます。
- **熱損失・熱取得の分類**: 出力されるデータには、
  `qldh`, `qldc`（暖房・冷房負荷）、`slo`, `slw`（太陽光による顕熱・潜熱取得）など、
  多岐にわたる熱損失・熱取得の項目が含まれています。
  これにより、熱負荷の発生源を詳細に分析できます。
- **建物全体の集計**: 各室のデータに加えて、
  `Qetotal.Qelmdy`として建物全体の熱損失・熱取得の日積算値の合計も出力されます。
  これにより、個別の室の熱負荷と建物全体のエネルギー性能を関連付けて評価できます。

この関数は、建物の熱負荷を詳細に分析し、
空調システムの設計、運用、省エネルギー対策の検討、
および快適性評価を行うための重要なデータ出力機能を提供します。
*/
func helmrmdyprint(fo io.Writer, id int, _Room []*ROOM, Qetotal *QETOTAL) {
	var i int
	var q *BHELM
	var qh *QHELM

	Nroom := len(_Room)

	switch id {
	case 0:
		if Nroom > 0 {
			fmt.Fprintf(fo, "%s %d\n", ROOM_TYPE, Nroom+1)
		}

		for i = 0; i < Nroom; i++ {
			Room := _Room[i]
			if Room.rmqe != nil {
				fmt.Fprintf(fo, "%s 1 %d\n", Room.Name, 29)
			}
		}
		fmt.Fprintf(fo, "%s 1 %d\n", Qetotal.Name, 29)
		break

	case 1:
		for i = 0; i < Nroom+1; i++ {
			var name string
			if i < Nroom {
				Room := _Room[i]
				name = Room.Name
			} else {
				name = Qetotal.Name
			}

			fmt.Fprintf(fo, "%s_qldh Q f %s_qldc Q f ", name, name)
			fmt.Fprintf(fo, "%s_slo Q f %s_slw Q f %s_asl Q f %s_tsol Q f %s_hins Q f\n",
				name, name, name, name, name)
			fmt.Fprintf(fo, "%s_so Q f %s_sw Q f %s_rn Q f %s_in Q f %s_pnl Q f\n",
				name, name, name, name, name)
			fmt.Fprintf(fo, "%s_trs Q f %s_qew Q f %s_qwn Q f %s_qgd Q f %s_qnx Q f ",
				name, name, name, name, name)
			fmt.Fprintf(fo, "%s_qi Q f %s_qc Q f %s_qf Q f\n",
				name, name, name)
			fmt.Fprintf(fo, "%s_qvo Q f %s_qvr Q f %s_sto Q f\n", name, name, name)
			fmt.Fprintf(fo, "%s_qldhl Q f %s_qldcl Q f %s_hinl Q f\n", name, name, name)
			fmt.Fprintf(fo, "%s_vol Q f %s_vrl Q f %s_stol Q f\n", name, name, name)
		}
		break

	default:
		for i = 0; i < Nroom+1; i++ {
			if i < Nroom {
				Room := _Room[i]
				q = &Room.rmqe.qelmdy.qe
				qh = &Room.rmqe.qelmdy
				fmt.Fprintf(fo, "%3.1f %3.1f ",
					qh.loadh*Cff_kWh, qh.loadc*Cff_kWh)
				fmt.Fprintf(fo, "%3.1f %3.1f %3.1f %3.1f %3.1f  ",
					qh.slo*Cff_kWh, qh.slw*Cff_kWh, qh.asl*Cff_kWh,
					qh.tsol*Cff_kWh, qh.hins*Cff_kWh)
				fmt.Fprintf(fo, "%3.1f %3.1f %3.1f %3.1f %3.1f  ",
					q.so*Cff_kWh, q.sg*Cff_kWh, q.rn*Cff_kWh,
					q.in*Cff_kWh, q.pnl*Cff_kWh)
				fmt.Fprintf(fo, "%3.1f %3.1f %3.1f %3.1f %3.1f ",
					q.trs*Cff_kWh, qh.ew*Cff_kWh,
					qh.wn*Cff_kWh, qh.gd*Cff_kWh, qh.nx*Cff_kWh)
				fmt.Fprintf(fo, "%3.1f %3.1f %3.1f %3.1f %3.1f %3.1f ",
					qh.i*Cff_kWh, qh.c*Cff_kWh, qh.f*Cff_kWh,
					qh.vo*Cff_kWh, qh.vr*Cff_kWh, qh.sto*Cff_kWh)
				fmt.Fprintf(fo, "%3.1f %3.1f %3.1f %3.1f %3.1f %3.1f\n",
					qh.loadhl*Cff_kWh, qh.loadcl*Cff_kWh, qh.hinl*Cff_kWh,
					qh.vol*Cff_kWh, qh.vrl*Cff_kWh, qh.stol*Cff_kWh)
			} else {
				q = &Qetotal.Qelmdy.qe
				qh = &Qetotal.Qelmdy
				fmt.Fprintf(fo, "%3.1f %3.1f ",
					qh.loadh*Cff_kWh, qh.loadc*Cff_kWh)
				fmt.Fprintf(fo, "%3.1f %3.1f %3.1f %3.1f %3.1f  ",
					qh.slo*Cff_kWh, qh.slw*Cff_kWh, qh.asl*Cff_kWh,
					qh.tsol*Cff_kWh, qh.hins*Cff_kWh)

				fmt.Fprintf(fo, "%3.1f %3.1f %3.1f %3.1f %3.1f  ",
					q.so*Cff_kWh, q.sg*Cff_kWh, q.rn*Cff_kWh,
					q.in*Cff_kWh, q.pnl*Cff_kWh)
				fmt.Fprintf(fo, "%3.1f %3.1f %3.1f %3.1f %3.1f ",
					q.trs*Cff_kWh, qh.ew*Cff_kWh,
					qh.wn*Cff_kWh, qh.gd*Cff_kWh, qh.nx*Cff_kWh)
				fmt.Fprintf(fo, "%3.1f %3.1f %3.1f %3.1f %3.1f %3.1f ",
					qh.i*Cff_kWh, qh.c*Cff_kWh, qh.f*Cff_kWh,
					qh.vo*Cff_kWh, qh.vr*Cff_kWh, qh.sto*Cff_kWh)
				fmt.Fprintf(fo, "%3.1f %3.1f %3.1f %3.1f %3.1f %3.1f\n",
					qh.loadhl*Cff_kWh, qh.loadcl*Cff_kWh, qh.hinl*Cff_kWh,
					qh.vol*Cff_kWh, qh.vrl*Cff_kWh, qh.stol*Cff_kWh)
			}
		}
		break
	}
}

/* ----------------------------------------------------- */

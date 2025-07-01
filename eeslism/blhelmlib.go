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

/* helmlib.c */

package eeslism

import "fmt"

/*
helmrmsrt (Building Element Heat Loss/Gain for Room Surface Temperature)

この関数は、室内の各表面（壁、窓など）における要素別の熱損失・熱取得を計算し、
その結果を表面温度の計算に反映させます。
これは、室内の熱環境を詳細に分析し、快適性評価を行うために不可欠です。

建築環境工学的な観点:
- **表面熱収支の構成**: 建物の表面は、日射、外気温度、室内温度、放射熱交換など、
  様々な要因によって熱を交換します。
  この関数は、各表面からの熱損失・熱取得を構成する要素を個別に計算し、
  `WSC`（要素別熱損失・熱取得の作業領域）に集計します。
- **熱損失・熱取得の分類**: `WSC`には、以下の熱損失・熱取得の項目が加算されます。
  - `WSC.trs`: 透過熱伝達による熱量（室内空気温度、外気温度、隣室温度などからの影響）。
  - `WSC.rn`: 夜間放射による熱損失。
  - `WSC.so`: 太陽光による熱取得（不透明部）。
  - `WSC.sg`: 太陽光による熱取得（透明部）。
  - `WSC.in`: 侵入空気による熱量。
  - `WSC.pnl`: 放射パネルからの熱量。
  これらの項目は、各表面の熱的挙動を詳細に分析し、
  熱負荷の発生源を特定するために用いられます。
- **壁体内部温度の考慮**: `helmsumpd(Mw.M, Mw.UX, Rmsb.Told, WSC)` は、
  壁体内部の温度履歴が表面からの熱伝達に与える影響を考慮しています。
  壁体は熱容量を持つため、その内部温度は時間遅れを伴って変化し、
  表面温度に影響を与えます。
- **表面温度の計算**: 最終的に、これらの要素別熱損失・熱取得を総合して、
  各表面の温度（`Rmsb.Ts`）が計算されます。
  この表面温度は、居住者の放射快適性や、結露の発生予測に直接影響します。

この関数は、建物の熱的挙動を詳細に分析し、
空調システムの設計、運用、省エネルギー対策の検討、
および快適性評価を行うための重要な役割を果たします。
*/
func helmrmsrt(Room *ROOM, Ta float64) {
	if Room.rmqe == nil {
		return
	}

	for i := 0; i < Room.N; i++ {
		Sd := Room.rsrf[i]
		Rmsb := Room.rmqe.rmsb[i]
		WSC := Room.rmqe.WSCwk[i]

		helmclear(WSC)

		if Mw := Sd.mw; Mw != nil {
			helmsumpd(Mw.M, Mw.UX, Rmsb.Told, WSC)
		}
		WSC.trs += Sd.FI * (Sd.alic / Sd.ali) * Room.Tr

		if Sd.rpnl != nil {
			Twmp := Sd.mw.Tw[Sd.mw.mp]
			WSC.pnl += Sd.FP * (Sd.rpnl.Tpi - Twmp)
			WSC.trs += Sd.FP * Twmp
		}

		switch Rmsb.Type {
		case RMSBType_E: // 外気に接する壁
			WSC.trs += Sd.FO * Ta
			WSC.rn += Sd.FO * Sd.TeErn
			if Sd.ble == BLE_ExternalWall {
				WSC.so += Sd.FO * Sd.TeEsol
			} else if Sd.ble == BLE_Window {
				WSC.sg += Sd.FO * Sd.TeEsol
			}
		case RMSBType_G: // 地盤に接する壁
			WSC.trs += Sd.FO * Sd.Te
		case RMSBType_i: // 内壁
			WSC.trs += Sd.FO * Sd.nextroom.Trold
		}

		WSC.sg += Sd.FI * Sd.RSsol / Sd.ali
		WSC.in += Sd.FI * Sd.RSin / Sd.ali
	}

	for i := 0; i < Room.N; i++ {
		Rmsb := Room.rmqe.rmsb[i]
		XA := Room.XA[Room.N*i : Room.N*(i+1)]
		Ts := &Rmsb.Ts
		helmclear(Ts)
		helmsumpd(Room.N, XA, Room.rmqe.WSCwk, Ts)
	}
}

/*
helmwall (Building Element Heat Loss/Gain for Walls)

この関数は、建物の壁体（不透明部）における要素別の熱損失・熱取得を計算し、
壁体内部の温度を更新します。
これは、壁体の蓄熱効果や熱貫流特性を詳細にモデル化するために不可欠です。

建築環境工学的な観点:
- **壁体内部温度の動的挙動**: 壁体は熱容量を持つため、
  その内部温度は外気温度や室内温度の変化に対して時間遅れを伴って応答します。
  この関数は、壁体内部の各層の温度を計算し、
  その結果を`rmsb.Tw`（現在の温度）と`rmsb.Told`（前時刻の温度）に格納します。
- **境界条件の考慮**: 壁体の熱伝達は、室内側と室外側の境界条件に大きく依存します。
  - `Tie`: 室内側相当温度。室内空気温度、室内表面からの放射熱、
    および室内側表面での日射吸収を考慮した熱的境界条件です。
  - `Te`: 室外側相当温度。外気温度、日射、夜間放射などを考慮した熱的境界条件です。
  - `Tpe`: パネルからの熱影響。壁体に放射パネルなどが組み込まれている場合に考慮されます。
- **熱損失・熱取得の分類**: `Tie.sg`（日射吸収）、`Tie.in`（内部発熱）など、
  壁体表面での熱損失・熱取得の項目が考慮されます。
- **壁体内部の熱伝導モデル (helmwlt)**:
  `helmwlt`関数は、壁体内部の熱伝導方程式を解き、
  各層の温度を計算します。
  これにより、壁体内部の温度分布や熱流を詳細にモデル化できます。
- **温度履歴の更新**: 計算された壁体内部温度は、
  次の時間ステップの計算のために`rmsb.Told`に更新されます。
  これにより、壁体の熱的履歴が考慮され、
  より正確な動的熱応答のシミュレーションが可能となります。

この関数は、建物の熱的挙動を詳細に分析し、
空調システムの設計、運用、省エネルギー対策の検討、
および快適性評価を行うための重要な役割を果たします。
*/
func helmwall(Room *ROOM, Ta float64) {
	if Room.rmqe == nil {
		return
	}

	for i := 0; i < Room.N; i++ {
		alr := Room.alr[Room.N*i : Room.N*(i+1)]
		Sd := Room.rsrf[i]
		rmsb := Room.rmqe.rmsb[i]

		if Mw := Sd.mw; Mw != nil {
			var Tie, Te, Tpe BHELM
			var Tm BHELM

			helmclear(&Tie)
			helmclear(&Te)
			helmclear(&Tpe)

			helmwlsft(i, Room.N, alr, Room.rmqe.rmsb, &Tm)

			helmsumpf(1, Sd.alir, &Tm, &Tie)

			Tie.trs += Sd.alic * Room.Tr

			if Sd.rpnl != nil {
				Twp := Mw.Tw
				Twmp := Twp[Mw.mp]
				Tpe.pnl = Sd.rpnl.Tpi - Twmp
				Tpe.trs = Twmp
			}

			switch rmsb.Type {
			case 'E': // 外気に接する壁
				Te.trs = Ta
				Te.so = Sd.TeEsol
				Te.rn = Sd.TeErn
			case 'G': // 地盤に接する壁
				Te.trs = Sd.Te
			case 'i': // 内壁
				Te.trs = Sd.nextroom.Trold
			}

			Tie.sg += Sd.RSsol
			Tie.in += Sd.RSin
			helmdiv(&Tie, Sd.ali)

			helmwlt(Mw.M, Mw.mp, Mw.UX, Mw.uo, Mw.um, Mw.Pc, []*BHELM{&Tie}, []*BHELM{&Te}, []*BHELM{&Tpe}, rmsb.Told, rmsb.Tw)

			for m := 0; m < Mw.M; m++ {
				Told := rmsb.Told[m]
				Tw := rmsb.Tw[m]
				helmcpy(Tw, Told)
			}
		}
	}
}

/* ---------------------------------------------- */

/*
helmwlsft (Building Element Wall Surface Temperature Calculation for Mean Radiant Temperature)

この関数は、室内の各表面における平均放射温度（MRT）の計算に用いられる、
他の表面からの放射熱伝達の影響を考慮した表面温度を計算します。

建築環境工学的な観点:
- **平均放射温度 (MRT) の重要性**: MRTは、居住者が感じる放射熱環境を代表する温度であり、
  作用温度やPMV（予測平均申告）などの快適性指標の算出に不可欠です。
  MRTは、室内の各表面温度とその表面に対する形態係数を考慮して計算されます。
- **形態係数の考慮**: `alr`は形態係数を含む配列であり、
  各表面が他の表面から受ける放射の影響を考慮します。
  `alr[0]`は、対象表面が他の表面から受ける放射の割合を示唆します。
- **表面間の放射熱交換**: この関数は、対象表面以外の全ての表面からの放射熱伝達を合計し、
  それを形態係数で割ることで、対象表面が他の表面から受ける放射の影響を平均的な温度として表現します。
  `helmsumpf`は、各表面の温度に形態係数を乗じて合計する処理を行い、
  `helmdiv`は、合計された値を正規化します。

この関数は、室内の熱環境を詳細に把握し、
居住者の快適性評価や、放射冷暖房システムの効果を評価するために不可欠な役割を果たします。
*/
func helmwlsft(i, N int, alr []float64, rmsb []*RMSB, Tm *BHELM) {
	Ralr := alr[i]

	helmclear(Tm)

	for j := 0; j < N; j++ {
		if j != i {
			helmsumpf(1, alr[0], &rmsb[j].Ts, Tm)
		}
	}

	helmdiv(Tm, Ralr)
}

/* ---------------------------------------------- */

/*
helmwlt (Building Element Wall Temperature Calculation)

この関数は、壁体内部の各層の温度を計算します。
これは、壁体の熱的挙動を詳細にモデル化し、
蓄熱効果や熱貫流特性を評価するために不可欠です。

建築環境工学的な観点:
- **壁体内部の熱伝導**: 壁体は熱容量を持つため、
  その内部温度は外気温度や室内温度の変化に対して時間遅れを伴って応答します。
  この関数は、壁体内部の熱伝導方程式を解き、
  各層の温度（`Tw`）を計算します。
- **境界条件の考慮**: 壁体の熱伝達は、室内側（`Tie`）と室外側（`Te`）の境界条件に大きく依存します。
  `helmsumpd`関数は、これらの境界条件と壁体の熱応答係数（`uo`, `um`）を用いて、
  壁体内部の温度変化を計算します。
- **パネルからの熱影響**: `Pc`が`0.0`より大きい場合、
  壁体に放射パネルなどが組み込まれており、
  そのパネルからの熱影響（`Tpe`）が壁体内部の温度に考慮されます。
- **温度履歴の考慮**: `Told`は前時刻の壁体内部温度であり、
  現在の温度計算にその履歴が考慮されます。
  これにより、壁体の熱的履歴が考慮され、
  より正確な動的熱応答のシミュレーションが可能となります。

この関数は、建物の熱的挙動を詳細にモデル化し、
空調システムの設計、運用、省エネルギー対策の検討、
および快適性評価を行うための重要な役割を果たします。
*/
func helmwlt(M, mp int, UX []float64, uo, um, Pc float64, Tie, Te, Tpe, Told, Tw []*BHELM) {
	helmsumpd(1, []float64{uo}, Tie, Told[0])
	helmsumpd(1, []float64{um}, Te, Told[M-1])

	if Pc > 0.0 {
		helmsumpd(1, []float64{Pc}, Tpe, Told[mp])
	}

	for m := 0; m < M; m++ {
		helmclear(Tw[m])
		helmsumpd(M, UX, Told, Tw[m])
		UX = UX[M:]
	}
}

/* ---------------------------------------------- */

/*
helmq (Building Element Heat Loss/Gain for Rooms)

この関数は、各室における建物の要素別熱損失・熱取得を計算します。
これは、室ごとの熱負荷を詳細に分析し、
空調システムの設計や運用、省エネルギー対策の検討に不可欠です。

建築環境工学的な観点:
- **室ごとの熱負荷分析**: 建物の熱負荷は、室ごとに異なる特性を持ちます。
  この関数は、各室の熱負荷を構成する要素（日射熱取得、壁体からの熱伝達、内部発熱、換気など）を個別に計算し、
  `Room.rmqe.qelm`に格納します。
  これにより、室ごとの熱負荷の要因を特定し、
  適切な空調ゾーン設定や、室ごとの省エネルギー対策を検討できます。
- **熱損失・熱取得の分類**: `q`（要素別熱損失・熱取得）には、以下の項目が加算されます。
  - `q.trs`: 透過熱伝達による熱量。
  - `q.so`: 太陽光による熱取得。
  - `q.sg`: 地盤からの熱取得。
  - `q.rn`: 夜間放射による熱損失。
  - `q.in`: 侵入空気による熱量。
  - `q.pnl`: 放射パネルからの熱量。
  これらの項目は、各表面の熱的挙動を詳細に分析し、
  熱負荷の発生源を特定するために用いられます。
- **内部発熱の考慮**: `q.in += Room.Hc + Room.Lc + Room.Ac` のように、
  人体発熱、照明、機器発熱などの内部発熱が熱負荷に加算されます。
- **換気による熱損失・熱取得**: `qh.vo`（顕熱）、`qh.vol`（潜熱）は、
  換気による熱損失・熱取得を表します。
  `Room.Gvent`は換気量、`Ta`は外気温度、`xa`は外気絶対湿度です。
- **室間相互換気**: `qh.vr`（顕熱）、`qh.vrl`（潜熱）は、
  隣接する室間での空気の移動による熱損失・熱取得を表します。
  `achr.Gvr`は室間相互換気量、`_Room[achr.rm].Tr`は隣室の温度です。
- **蓄熱量の考慮**: `qh.sto`（顕熱）、`qh.stol`（潜熱）は、
  室内の熱容量による蓄熱量の変化を表します。

この関数は、建物の熱負荷を詳細に分析し、
空調システムの設計、運用、省エネルギー対策の検討、
および快適性評価を行うための重要な役割を果たします。
*/
func helmq(_Room []*ROOM, Ta, xa float64) {
	var q, Ts *BHELM
	var qh *QHELM
	var Sd *RMSRF
	var rmsb *RMSB
	var achr *ACHIR
	var Aalc, qloss float64

	Room := _Room[0]

	qelmclear(&Room.rmqe.qelm)
	q = &Room.rmqe.qelm.qe
	qh = &Room.rmqe.qelm

	qh.loadh = 0.0
	qh.loadc = 0.0
	qh.loadcl = 0.0
	qh.loadhl = 0.0
	if Room.rmld != nil {
		if Room.rmld.Qs > 0.0 {
			qh.loadh = Room.rmld.Qs
		} else {
			qh.loadc = Room.rmld.Qs
		}

		if Room.rmld.Ql > 0.0 {
			qh.loadhl = Room.rmld.Ql
		} else {
			qh.loadcl = Room.rmld.Ql
		}
	}

	for i := 0; i < Room.N; i++ {
		Sd = Room.rsrf[i]
		rmsb = Room.rmqe.rmsb[i]

		Aalc = Sd.A * Sd.alic
		Ts = &rmsb.Ts
		qloss = Aalc * (Room.Tr - Ts.trs)
		q.trs -= qloss

		if rmsb.Type == RMSBType_E {
			if Sd.ble == BLE_ExternalWall {
				qh.ew -= qloss
			} else if Sd.ble == BLE_Window {
				qh.wn -= qloss
			}
		} else if rmsb.Type == RMSBType_G {
			qh.gd -= qloss
		} else if rmsb.Type == RMSBType_i {
			qh.nx -= qloss
		}

		if Sd.ble == BLE_Ceil || Sd.ble == BLE_Roof {
			qh.c -= qloss
		} else if Sd.ble == BLE_InnerFloor || Sd.ble == BLE_Floor {
			qh.f -= qloss
		} else if Sd.ble == BLE_InnerWall || Sd.ble == BLE_d {
			qh.i -= qloss
		}

		q.so += Aalc * Ts.so
		q.sg += Aalc * Ts.sg
		q.rn += Aalc * Ts.rn
		q.in += Aalc * Ts.in
		q.pnl += Aalc * Ts.pnl
	}

	q.in += Room.Hc + Room.Lc + Room.Ac

	qh.hinl = Room.AL + Room.HL

	qh.sto = Room.MRM * (Room.Trold - Room.Tr) / DTM
	qh.stol = Room.GRM * Ro * (Room.xrold - Room.xr) / DTM
	qh.vo = Ca * Room.Gvent * (Ta - Room.Tr)
	qh.vol = Ro * Room.Gvent * (xa - Room.xr)

	qh.vr = 0.0
	qh.vrl = 0.0
	for j := 0; j < Room.Nachr; j++ {
		achr = Room.achr[j]
		qh.vr += Ca * achr.Gvr * (_Room[achr.rm].Tr - Room.Tr)
		qh.vrl += Ro * achr.Gvr * (_Room[achr.rm].xr - Room.xr)
	}
}

/*
qelmclear (Building Element Heat Loss/Gain Data Clear)

この関数は、建物の要素別熱損失・熱取得のデータをゼロにリセットします。
これは、新しい時間ステップの計算を開始する前に、
前回の計算結果をクリアするために用いられます。

建築環境工学的な観点:
- **熱収支計算の準備**: 建物の熱負荷計算は、
  各時間ステップで熱収支方程式を解くことで行われます。
  この関数は、熱収支方程式の各項をゼロに初期化することで、
  新しい時間ステップでの正確な計算を可能にします。
- **データ集計の準備**: 日積算値や月積算値を計算する際に、
  各時間ステップの計算結果を正確に集計するために、
  この関数で各項目をリセットします。
- **熱損失・熱取得の分類**: `q.slo`, `q.slw`, `q.asl`, `q.tsol`, `q.hins`（日射熱取得、内部発熱）、
  `q.nx`, `q.gd`, `q.ew`, `q.wn`（壁体からの熱伝達）、
  `q.i`, `q.c`, `q.f`（内部熱伝達）、
  `q.vo`, `q.vr`, `q.sto`（換気、蓄熱）など、
  多岐にわたる熱損失・熱取得の項目がゼロに初期化されます。

この関数は、建物の熱負荷計算を正確に行い、
省エネルギー対策の効果を評価するための基礎的な役割を果たします。
*/
func qelmclear(q *QHELM) {
	helmclear(&q.qe)
	q.slo = 0.0
	q.slw = 0.0
	q.asl = 0.0
	q.tsol = 0.0
	q.hins = 0.0
	q.nx = 0.0
	q.gd = 0.0
	q.ew = 0.0
	q.wn = 0.0
	q.i = 0.0
	q.c = 0.0
	q.f = 0.0
	q.vo = 0.0
	q.vr = 0.0
	q.sto = 0.0
	q.loadh = 0.0
	q.loadc = 0.0
	q.hinl = 0.0
	q.vol = 0.0
	q.vrl = 0.0
	q.stol = 0.0
	q.loadcl = 0.0
	q.loadhl = 0.0
}

/*
qelmsum (Building Element Heat Loss/Gain Data Summation)

この関数は、建物の要素別熱損失・熱取得のデータを合計します。
これは、各室の熱負荷を建物全体の熱負荷に集計したり、
日積算値や月積算値を計算したりするために用いられます。

建築環境工学的な観点:
- **熱負荷の集計**: 建物の熱負荷は、各室や各要素からの熱損失・熱取得の合計として評価されます。
  この関数は、`a`（加算されるデータ）の各項目を`b`（合計されるデータ）に加算することで、
  熱負荷の集計を行います。
- **データ集計の柔軟性**: この関数は、
  - 各室の熱負荷を建物全体の熱負荷に集計する。
  - 各時間ステップの熱負荷を日積算値や月積算値に集計する。
  など、様々なレベルでのデータ集計に利用できます。
- **熱損失・熱取得の分類**: `b.slo += a.slo` のように、
  日射熱取得、内部発熱、壁体からの熱伝達、換気、蓄熱など、
  多岐にわたる熱損失・熱取得の項目が個別に合計されます。
  これにより、熱負荷の発生源を詳細に分析し、
  省エネルギー対策の効果を評価できます。

この関数は、建物の熱負荷計算を正確に行い、
省エネルギー対策の効果を評価するための基礎的な役割を果たします。
*/
func qelmsum(a, b *QHELM) {
	helmsum(&a.qe, &b.qe)

	b.slo += a.slo
	b.slw += a.slw
	b.asl += a.asl
	b.tsol += a.tsol
	b.hins += a.hins

	b.nx += a.nx
	b.gd += a.gd
	b.ew += a.ew
	b.wn += a.wn

	b.i += a.i
	b.c += a.c
	b.f += a.f
	b.vo += a.vo
	b.vr += a.vr
	b.sto += a.sto
	b.loadh += a.loadh
	b.loadc += a.loadc

	b.hinl += a.hinl
	b.vol += a.vol
	b.stol += a.stol
	b.vrl += a.vrl
	b.loadcl += a.loadcl
	b.loadhl += a.loadhl
}

/*
helmclear (Building Element Heat Loss/Gain Data Clear for BHELM)

この関数は、`BHELM`構造体の各項目をゼロにリセットします。
`BHELM`構造体は、建物の要素別熱損失・熱取得の各成分（透過熱伝達、日射、地盤など）を格納するために用いられます。

建築環境工学的な観点:
- **熱収支計算の準備**: 建物の熱負荷計算は、
  各時間ステップで熱収支方程式を解くことで行われます。
  この関数は、熱収支方程式の各項をゼロに初期化することで、
  新しい時間ステップでの正確な計算を可能にします。
- **データ集計の準備**: 各成分を個別に計算し、
  最終的に合計することで熱負荷を算出するため、
  この関数で各項目をリセットします。
- **熱損失・熱取得の分類**: `trs`（透過熱伝達）、`so`（太陽光）、`sg`（地盤）、
  `rn`（夜間放射）、`in`（侵入空気）、`pnl`（パネル）など、
  多岐にわたる熱損失・熱取得の項目がゼロに初期化されます。

この関数は、建物の熱負荷計算を正確に行い、
省エネルギー対策の効果を評価するための基礎的な役割を果たします。
*/
func helmclear(b *BHELM) {
	b.trs = 0.0
	b.so = 0.0
	b.sg = 0.0
	b.rn = 0.0
	b.in = 0.0
	b.pnl = 0.0
}

/*
helmsumpd (Building Element Heat Loss/Gain Data Summation with Product)

この関数は、`BHELM`構造体の配列`a`の各項目に、
対応する`u`（係数）を乗じて`b`（合計されるデータ）に加算します。
これは、熱応答係数などを用いて、
壁体内部の温度履歴が表面からの熱伝達に与える影響を計算する際に用いられます。

建築環境工学的な観点:
- **熱応答係数の適用**: 建物の熱的挙動は、
  過去の温度履歴や日射履歴に依存します。
  熱応答係数（`u`）は、これらの履歴が現在の熱伝達にどのように影響するかを示す係数です。
  この関数は、熱応答係数を各熱損失・熱取得成分に乗じることで、
  壁体などの熱容量を持つ要素の動的な熱的挙動をモデル化します。
- **熱損失・熱取得の分類**: `trs`（透過熱伝達）、`so`（太陽光）、`sg`（地盤）、
  `rn`（夜間放射）、`in`（侵入空気）、`pnl`（パネル）など、
  多岐にわたる熱損失・熱取得の項目が個別に計算され、合計されます。

この関数は、建物の熱負荷計算を正確に行い、
省エネルギー対策の効果を評価するための基礎的な役割を果たします。
*/
func helmsumpd(N int, u []float64, a []*BHELM, b *BHELM) {
	for i := 0; i < N; i++ {
		b.trs += u[i] * a[i].trs
		b.so += u[i] * a[i].so
		b.sg += u[i] * a[i].sg
		b.rn += u[i] * a[i].rn
		b.in += u[i] * a[i].in
		b.pnl += u[i] * a[i].pnl
	}
}

/*
helmsumpf (Building Element Heat Loss/Gain Data Summation with Scalar Product)

この関数は、`BHELM`構造体`a`の各項目に、
スカラー値`u`を乗じて`b`（合計されるデータ）に加算します。
これは、主に形態係数などを用いて、
表面間の放射熱伝達の影響を計算する際に用いられます。

建築環境工学的な観点:
- **形態係数の適用**: 形態係数（`u`）は、
  ある表面から別の表面へ放射される熱の割合を示す無次元数です。
  この関数は、形態係数を各熱損失・熱取得成分に乗じることで、
  表面間の放射熱伝達の影響をモデル化します。
- **熱損失・熱取得の分類**: `trs`（透過熱伝達）、`so`（太陽光）、`sg`（地盤）、
  `rn`（夜間放射）、`in`（侵入空気）、`pnl`（パネル）など、
  多岐にわたる熱損失・熱取得の項目が個別に計算され、合計されます。

この関数は、建物の熱負荷計算を正確に行い、
省エネルギー対策の効果を評価するための基礎的な役割を果たします。
*/
func helmsumpf(N int, u float64, a *BHELM, b *BHELM) {
	if N != 1 {
		panic("N != 1")
	}

	b.trs += u * a.trs
	b.so += u * a.so
	b.sg += u * a.sg
	b.rn += u * a.rn
	b.in += u * a.in
	b.pnl += u * a.pnl
}

/*
helmdiv (Building Element Heat Loss/Gain Data Division)

この関数は、`BHELM`構造体`a`の各項目を、
スカラー値`c`で除算します。
これは、平均表面温度の計算など、
熱損失・熱取得の値を正規化する際に用いられます。

建築環境工学的な観点:
- **平均値の算出**: 熱損失・熱取得の合計値を面積や形態係数などで除算することで、
  平均的な熱的挙動を評価できます。
  例えば、平均放射温度の計算において、
  各表面からの放射熱伝達の合計を形態係数で除算することで、
  平均的な放射環境を評価できます。
- **熱損失・熱取得の分類**: `trs`（透過熱伝達）、`so`（太陽光）、`sg`（地盤）、
  `rn`（夜間放射）、`in`（侵入空気）、`pnl`（パネル）など、
  多岐にわたる熱損失・熱取得の項目が個別に除算されます。

この関数は、建物の熱負荷計算を正確に行い、
省エネルギー対策の効果を評価するための基礎的な役割を果たします。
*/
func helmdiv(a *BHELM, c float64) {
	a.trs /= c
	a.so /= c
	a.sg /= c
	a.rn /= c
	a.in /= c
	a.pnl /= c
}

/*
helmsum (Building Element Heat Loss/Gain Data Summation)

この関数は、`BHELM`構造体`a`の各項目を、
`b`（合計されるデータ）に加算します。
これは、各要素からの熱損失・熱取得を合計する際に用いられます。

建築環境工学的な観点:
- **熱負荷の集計**: 建物の熱負荷は、各要素からの熱損失・熱取得の合計として評価されます。
  この関数は、`a`の各項目を`b`に加算することで、
  熱負荷の集計を行います。
- **熱損失・熱取得の分類**: `trs`（透過熱伝達）、`so`（太陽光）、`sg`（地盤）、
  `rn`（夜間放射）、`in`（侵入空気）、`pnl`（パネル）など、
  多岐にわたる熱損失・熱取得の項目が個別に合計されます。

この関数は、建物の熱負荷計算を正確に行い、
省エネルギー対策の効果を評価するための基礎的な役割を果たします。
*/
func helmsum(a, b *BHELM) {
	b.trs += a.trs
	b.so += a.so
	b.sg += a.sg
	b.rn += a.rn
	b.in += a.in
	b.pnl += a.pnl
}

/*
helmcpy (Building Element Heat Loss/Gain Data Copy)

この関数は、`BHELM`構造体`a`の各項目を、
`b`にコピーします。
これは、熱損失・熱取得のデータを一時的に保存したり、
前時刻のデータを更新したりする際に用いられます。

建築環境工学的な観点:
- **データ管理**: シミュレーションでは、
  各時間ステップで計算された熱損失・熱取得のデータを一時的に保存したり、
  次の時間ステップの計算のために前時刻のデータを更新したりする必要があります。
  この関数は、そのようなデータ管理を効率的に行います。
- **熱損失・熱取得の分類**: `trs`（透過熱伝達）、`so`（太陽光）、`sg`（地盤）、
  `rn`（夜間放射）、`in`（侵入空気）、`pnl`（パネル）など、
  多岐にわたる熱損失・熱取得の項目が個別にコピーされます。

この関数は、建物の熱負荷計算を正確に行い、
省エネルギー対策の効果を評価するための基礎的な役割を果たします。
*/
func helmcpy(a, b *BHELM) {
	b.trs = a.trs
	b.so = a.so
	b.sg = a.sg
	b.rn = a.rn
	b.in = a.in
	b.pnl = a.pnl
}

/* ========================================== */

func helmxxprint(s string, a *BHELM) {
	fmt.Printf("xxx helmprint xxx %s  trs so sg rn in pnl\n", s)
	fmt.Printf("%6.2f %6.2f %6.2f %6.2f %6.2f %6.2f\n", a.trs, a.so, a.sg, a.rn, a.in, a.pnl)
}

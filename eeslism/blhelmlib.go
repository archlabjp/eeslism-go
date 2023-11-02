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

// 外乱要素別表面温度の計算
// 入力値:
//  外気温度 Ta [C]
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

/* ---------------------------------------------- */

/* 壁体内温度 */
// 入力値:
//  外気温度 Ta [C]
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

/* 面 i についての平均表面温度 */

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

/* 要素別熱損失・熱取得 */

// 入力値:
//  外気温度 Ta [C]
//  絶対湿度 xa [kg/kg]
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

/* ---------------------------------------------- */

// Reset q to zero
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

/* ---------------------------------------------- */

// Add a to b
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

/* ---------------------------------------------- */

// Reset b to zero
func helmclear(b *BHELM) {
	b.trs = 0.0
	b.so = 0.0
	b.sg = 0.0
	b.rn = 0.0
	b.in = 0.0
	b.pnl = 0.0
}

/* ---------------------------------------------- */

// Mutiply a by u(vector) and add to b
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

/* ---------------------------------------------- */

// Mutiply a by u(scalar) and add to b
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

/* ---------------------------------------------- */

// Divide a by c
func helmdiv(a *BHELM, c float64) {
	a.trs /= c
	a.so /= c
	a.sg /= c
	a.rn /= c
	a.in /= c
	a.pnl /= c
}

/* ---------------------------------------------- */

// Add a to b
func helmsum(a, b *BHELM) {
	b.trs += a.trs
	b.so += a.so
	b.sg += a.sg
	b.rn += a.rn
	b.in += a.in
	b.pnl += a.pnl
}

/* ---------------------------------------------- */

// Copy a to b
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

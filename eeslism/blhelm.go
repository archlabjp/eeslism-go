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

/* 要素別熱損失・熱取得（記憶域確保） */

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

/* ----------------------------------------------------- */

// 要素別熱損失・熱取得（計算）
// 入力値:
//  外気温度 Ta [C]
//  絶対湿度 xa [kg/kg]
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

/* ----------------------------------------------------- */

/* 要素別熱損失・熱取得（時刻別出力） */

var __Helmprint_id int = 0

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

/* 要素別熱損失・熱取得（時刻別出力） */

var __Helmsurfprint_id int = 0

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

/* 要素別熱損失・熱取得（日積算値） */

var __Helmdy_oldday int = -1

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

func helmdyint(Room []*ROOM, Qetotal *QETOTAL) {
	for i := range Room {
		if Room[i].rmqe != nil {
			qelmclear(&Room[i].rmqe.qelmdy)
		}
	}

	qelmclear(&Qetotal.Qelmdy)
}

/* ----------------------------------------------------- */

/* 要素別熱損失・熱取得（日積算値出力） */

var __Helmdyprint_id int

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

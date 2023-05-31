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

/*  room.c       */
package main

import "fmt"

/* ----------------------------------------------------------- */

/*  室内サブシステムに関する計算         */

func RMcf(Room *ROOM) {
	N := Room.N
	for n := 0; n < N; n++ {
		Sdn := &Room.rsrf[n]

		if Sdn.mrk == '*' || Sdn.PCMflg == 'Y' {

			// 壁体（窓以外の場合）
			if Sdn.typ == 'H' || Sdn.typ == 'E' || Sdn.typ == 'e' {
				Mw := Sdn.mw
				M := Mw.M
				mp := Mw.mp

				if Sdn.mwside == 'i' {
					Sdn.FI = Mw.uo * Mw.UX[0]

					if Sdn.mw.wall.WallType != 'C' {
						Sdn.FO = Mw.um * Mw.UX[M-1]
					} else {
						Sdn.FO = Sdn.ColCoeff * Mw.UX[M-1]
					}

					if Sdn.rpnl != nil {
						Sdn.FP = Mw.Pc * Mw.UX[mp] * Sdn.rpnl.Wp
					} else {
						Sdn.FP = 0.0
					}
				} else {
					MM := (M - 1) * M

					if Sdn.mw.wall.WallType != 'C' {
						Sdn.FI = Mw.um * Mw.UX[MM+M-1]
					} else {
						Sdn.FI = Sdn.ColCoeff * Mw.UX[MM+M-1]
					}

					Sdn.FO = Mw.uo * Mw.UX[MM]
					if Sdn.rpnl != nil {
						Sdn.FP = Mw.Pc * Mw.UX[MM+mp] * Sdn.rpnl.Wp
					} else {
						Sdn.FP = 0.0
					}
				}
			} else {
				// 窓の場合
				/***      K = Sdn.K = 1.0/(Sdn.Rwall + rai + rao); ***/
				K := Sdn.K
				ali := Sdn.ali
				Sdn.FI = 1.0 - K/ali
				Sdn.FO = K / ali
				Sdn.FP = 0.0
			}
		}
	}

	alr := Room.alr
	XA := Room.XA
	for n := 0; n < N; n++ {
		Sdn := &Room.rsrf[n]

		for j := 0; j < N; j++ {
			XA[n*N+j] = -Sdn.FI * alr[n*N+j] / Sdn.ali
		}

		XA[n*N+n] = 1.0
	}

	E := fmt.Sprintf("<RMcf> name=%s", Room.Name)
	Matinv(XA, N, N, E)

	for n := 0; n < N; n++ {
		Sdn := &Room.rsrf[n]

		Sdn.WSR = 0.0

		for j := 0; j < N; j++ {
			sdj := &Room.rsrf[j]
			kc := sdj.alic / sdj.ali
			Sdn.WSR += XA[n*N+j] * sdj.FI * kc
		}

		for j := 0; j < Room.Ntr; j++ {
			wrn := &Sdn.WSRN[j]
			trn := &Room.trnx[j]
			sdk := trn.sd

			// Find the index of sdk in Room.rsrf
			var kk int
			for kk = 0; kk < Room.N; kk++ {
				if sdk == &Room.rsrf[kk] {
					break
				}
			}
			*wrn = XA[n*N+kk] * sdk.FO * sdk.nxsd.alic / sdk.nxsd.ali
		}

		for j := 0; j < Room.Nrp; j++ {
			sdk := Room.rmpnl[j].sd

			// Find the index of sdk in Room.rsrf
			var kk int
			for kk = 0; kk < Room.N; kk++ {
				if sdk == &Room.rsrf[kk] {
					break
				}
			}

			// XA：室内表面温度計算のためのマトリックス
			// FP：パネルの係数
			Sdn.WSPL[j] = XA[n*N+kk] * sdk.FP
		}
	}

	Room.AR = 0.0
	for n := 0; n < N; n++ {
		Sdn := &Room.rsrf[n]
		Room.AR += Sdn.A * Sdn.alic * (1.0 - Sdn.WSR)
	}

	// 室内空気の総合熱収支式の係数
	for j := 0; j < Room.Ntr; j++ {
		arn := 0.0
		for n := 0; n < N; n++ {
			sdk := &Room.rsrf[n]
			arn += sdk.A * sdk.alic * sdk.WSRN[j]
		}
		Room.ARN[j] = arn
	}

	for j := 0; j < Room.Nrp; j++ { // 室のパネル総数
		rpnl := 0.0
		for n := 0; n < N; n++ {
			sdk := &Room.rsrf[n]
			rpnl += sdk.A * sdk.alic * sdk.WSPL[j] // WSPL：パネルに関する係数
		}
		Room.RMP[j] = rpnl
	}

	// 室温の係数
	// 家具の熱容量の計算
	FunCoeff(Room)
}

// 家具内蔵PCMの係数計算
func FunCoeff(Room *ROOM) {
	// 室温の係数
	// 家具の熱容量の計算
	Room.FunHcap = 0.0
	if Room.CM != nil && *Room.CM > 0.0 {
		if Room.MCAP != nil && *Room.MCAP > 0.0 {
			Room.FunHcap += *Room.MCAP
		}
		if Room.PCM != nil {
			if Room.PCM.Spctype == 'm' {
				Room.PCMQl = FNPCMStatefun(Room.PCM.Ctype, Room.PCM.Cros, Room.PCM.Crol, Room.PCM.Ql,
					Room.PCM.Ts, Room.PCM.Tl, Room.PCM.Tp, Room.oldTM, Room.TM, Room.PCM.DivTemp, &Room.PCM.PCMp)
			} else {
				Room.PCMQl = FNPCMstate_table(&Room.PCM.Chartable[0], Room.oldTM, Room.TM, Room.PCM.DivTemp)
			}
			Room.FunHcap += Room.mPCM * Room.PCMQl
		}
	}
	if Room.FunHcap > 0.0 {
		Room.FMT = 1.0 / (Room.FunHcap/DTM/(*Room.CM) + 1.0)
	} else {
		Room.FMT = 1.0
	}

	Room.RMt = Room.MRM/DTM + Room.AR

	if Room.FunHcap > 0.0 {
		Room.RMt -= *Room.CM * (Room.FMT - 1.0)
	}
}

func RMrc(Room *ROOM) {
	N := Room.N
	XA := Room.XA
	CRX := make([]float64, N)

	for n := 0; n < N; n++ { // N：表面総数
		Sdn := &Room.rsrf[n]
		Sdn.CF = 0.0
		if Sdn.typ == 'H' || Sdn.typ == 'E' || Sdn.typ == 'e' { // 壁の場合
			Mw := Sdn.mw
			M := Mw.M
			if Sdn.mwside != 'M' { // 室内側
				for j := 0; j < M; j++ {
					Sdn.CF += Mw.UX[j] * Mw.Told[j]
				}
			} else {
				MM := M * (M - 1)
				UX := Mw.UX[MM:]
				for j := 0; j < M; j++ {
					Sdn.CF += UX[j] * Mw.Told[j]
				}
			}
		}
	}

	Room.HGc = Room.Hc + Room.Lc + Room.Ac + Room.Qeqp*Room.eqcv

	// 表面熱収支に関係する係数の計算
	for n := 0; n < N; n++ {
		Sdn := &Room.rsrf[n]
		CRX[n] = Sdn.CF + Sdn.FO*Sdn.Te + Sdn.FI*Sdn.RS/Sdn.ali
	}

	// 相互放射の計算
	for n := 0; n < N; n++ {
		Sdn := &Room.rsrf[n]
		Sdn.WSC = 0.0
		for j := 0; j < N; j++ {
			Sdn.WSC += XA[n*N+j] * CRX[j]
		}
	}

	Room.CA = 0.0
	for n := 0; n < N; n++ {
		Sdn := &Room.rsrf[n]
		Room.CA += Sdn.A * Sdn.alic * Sdn.WSC
	}

	// 室空気の熱収支の係数計算
	// 家具の影響項の追加
	if Room.FunHcap > 0.0 {
		dblTemp := DTM / Room.FunHcap
		Room.FMC = 1.0 / (dblTemp**Room.CM + 1.0) * (Room.oldTM + dblTemp*Room.Qsolm)
	} else {
		Room.FMC = 0.0
	}

	Room.RMC = Room.MRM/DTM*Room.Trold + Room.HGc + Room.CA
	if Room.FunHcap > 0.0 {
		Room.RMC += *Room.CM * Room.FMC
	}

}

/* ----------------------------------------------------- */
// 室Roomの壁体の表面温度の計算 -- RooM's SuRface Temperature
func RMsrt(Room *ROOM) {
	N := Room.N

	for n := 0; n < N; n++ {
		Sdn := &Room.rsrf[n]

		Sdn.Ts = Sdn.WSR*Room.Tr + Sdn.WSC

		for j := 0; j < Room.Ntr; j++ {
			trn := &Room.trnx[j]
			Sdn.Ts += Sdn.WSRN[j] * trn.nextroom.Tr
		}

		for j := 0; j < Room.Nrp; j++ {
			rmpnl := &Room.rmpnl[j]
			Sdn.Ts += Sdn.WSPL[j] * rmpnl.pnl.Tpi
		}
	}

	alr := Room.alr
	for n := 0; n < N; n++ {
		Sdn := &Room.rsrf[n]
		Sdn.Tmrt = 0.0

		for j := 0; j < N; j++ {
			Sd := &Room.rsrf[j]
			if j != n {
				Sdn.Tmrt += Sd.Ts * alr[n*N+j]
			}
		}
		Sdn.Tmrt /= alr[n*N+n]
	}

	for n := 0; n < N; n++ {
		Sd := &Room.rsrf[n]
		Sd.Qc = Sd.alic * Sd.A * (Sd.Ts - Room.Tr)
		Sd.Qr = Sd.alir * Sd.A * (Sd.Ts - Sd.Tmrt)
		Sd.Qi = Sd.Qc + Sd.Qr - Sd.RS*Sd.A
	}
}

/* ----------------------------------------------------- */

// 重量壁（後退差分）の係数行列の作成
func RMwlc(Nmwall int, Mw []MWALL, Exsfs *EXSFS, Wd *WDAT) {
	for i := 0; i < Nmwall; i++ {
		var Mw *MWALL = &Mw[i]
		var Wall *WALL = Mw.wall

		var Sd *RMSRF = Mw.sd
		rai := 1.0 / Sd.ali // 室内側表面熱抵抗
		rao := 1.0 / Sd.alo // 室外側表面熱抵抗

		Mw.res[0] = rai
		if Sd.typ == 'H' {
			Mw.res[Mw.M] = rao
		}

		// 壁体にパネルがある場合
		var Wp float64
		if Sd.rpnl != nil {
			Wp = Sd.rpnl.Wp
		} else {
			Wp = 0.0
		}

		// 行列作成
		Wallfdc(Mw.M, Mw.mp, Mw.res, Mw.cap, Wp, Mw.UX,
			&Mw.uo, &Mw.um, &Mw.Pc, Wall.WallType, Sd, Wd, Exsfs, Wall,
			Mw.Told, Mw.Toldd, Mw.sd.pcmstate)
	}
}

/* ----------------------------------------------------- */

// 壁体内部温度の計算
func RMwlt(Nmwall int, Mw []MWALL) {
	for i := 0; i < Nmwall; i++ {
		Mw := &Mw[i]
		Sd := Mw.sd

		// 壁体の反対側の表面温度 ?
		var Tee float64
		if Sd.mwtype == 'C' {
			// 共用壁の場合
			nxsd := Sd.nxsd
			Tee = (nxsd.alic*nxsd.room.Tr + nxsd.alir*nxsd.Tmrt + nxsd.RS) / nxsd.ali
		} else {
			// 専用壁の場合 => 外表面の相当外気温度
			Tee = Sd.Te
		}

		Room := Sd.room
		Tie := (Sd.alic*Room.Tr + Sd.alir*Sd.Tmrt + Sd.RS) / Sd.ali

		if DEBUG {
			fmt.Printf("----- RMwlt i=%d room=%s ble=%c %s  Tie=%f Tee=%f\n", i, Sd.room.Name, Sd.ble, Sd.Name, Tie, Tee)
		}

		var WTp float64
		if Sd.rpnl != nil {
			WTp = Sd.rpnl.Wp * Sd.rpnl.Tpi
		} else {
			WTp = 0.0
		}

		// 壁体表面、壁体内部温度の計算
		Twall(Mw.M, Mw.mp, Mw.UX, Mw.uo, Mw.um, Mw.Pc, Tie, Tee, WTp, Mw.Told, Mw.Tw, Sd, Mw.wall.PCMLyr)

		// 壁体表面温度、壁体内部温度の更新
		for m := 0; m < Mw.M; m++ {
			// 前時刻の壁体内部温度を更新
			Mw.Told[m] = Mw.Tw[m]
			// 収束過程初期値の壁体内部温度を更新
			Mw.Twd[m] = Mw.Tw[m]
			Mw.Told[m] = Mw.Tw[m]
		}
	}
}

// 壁体内部温度の仮計算
func RMwltd(Nmwall int, Mw []MWALL) {
	for i := 0; i < Nmwall; i++ {
		var Mw *MWALL = &Mw[i]
		var Sd *RMSRF = Mw.sd
		var nxsd *RMSRF = Sd.nxsd
		var Room *ROOM = Sd.room

		if Sd.PCMflg == 'Y' {
			// Tee
			var Tee float64
			if Sd.mwtype == 'C' {
				Tee = (nxsd.alic*nxsd.room.Tr + nxsd.alir*nxsd.Tmrt + nxsd.RS) /
					nxsd.ali
			} else {
				Tee = Sd.Te
			}

			// Tie
			Tie := (Sd.alic*Room.Tr + Sd.alir*Sd.Tmrt + Sd.RS) / Sd.ali

			if DEBUG {
				fmt.Printf("----- RMwlt i=%d room=%s ble=%c %s  Tie=%f Tee=%f\n",
					i, Sd.room.Name, Sd.ble, Sd.Name, Tie, Tee)
			}

			// WTp
			var WTp float64
			if Sd.rpnl != nil {
				WTp = Sd.rpnl.Wp * Sd.rpnl.Tpi
			} else {
				WTp = 0.0
			}

			// 壁体内部温度の仮計算
			Twalld(Mw.M, Mw.mp, Mw.UX, Mw.uo, Mw.um, Mw.Pc,
				Tie, Tee, WTp, Mw.Told, Mw.Twd, Sd)
		}
	}
}

/* ------------------------------------------------------ */

// 室内表面 Sd における平均表面温度の計算 (Room's Temperature of Surface - AVerage)
func RTsav(N int, Sd []RMSRF) float64 {
	var Tav, Aroom float64
	for n := 0; n < N; n++ {
		Tav += Sd[n].Ts * Sd[n].A
		Aroom += Sd[n].A
	}
	return Tav / Aroom
}

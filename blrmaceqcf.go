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

package main

/* -------------------------------------- */

var __Rmhtrcf_count int

func Rmhtrcf(exs *EXSFS, emrk []rune, rooms []ROOM, sds []RMSRF, wd *WDAT) {
	if rooms != nil {
		for _, room := range rooms {
			n := room.N
			brs := room.Brs
			sds := sds[brs : brs+n]

			// 放射熱交換係数の計算
			if __Rmhtrcf_count == 0 || emrk[0] == '*' {
				radex(n, sds, room.F, room.Wradx)
			}

			// 放射熱伝達率の入れ替え
			Radcf0(room.Tsav, &room.alrbold, n, sds, room.Wradx, room.alr)

			Htrcf(room.alc, exs.Alosch, exs.Alotype, exs.Exs, room.Tr, n, room.alr, sds, &room.mrk, wd)
		}

		__Rmhtrcf_count++
	}

	if sds != nil {
		for _, sd := range sds {
			if sd.mwtype == 'C' && sd.mwside == 'i' {
				// 内壁の場合は裏面室の熱伝達率を入れ替える
				nxsd := sd.nxsd
				sd.alo = nxsd.ali
				nxsd.alo = sd.ali
			}
		}
	}
}

/* ----------------------------------------------------------------- */

func Rmrdshfc(_Room []ROOM, Sd []RMSRF) {
	if len(_Room) == 0 {
		return
	}

	N := _Room[0].end

	for i := 0; i < N; i++ {
		Room := &_Room[i]
		brs := Room.Brs
		sd := Sd[brs:]

		Radshfc(Room.N, Room.FArea, Room.Area, sd, Room.tfsol, Room.eqcv, Room.Name, Room.fsolm)
	}
}

/* ----------------------------------------------------------------- */
func Rmhtrsmcf(Nsrf int, _Sd []RMSRF) {
	for n := 0; n < Nsrf; n++ {
		Sd := &_Sd[n]
		Sd.K = 1.0 / (Sd.Rwall + 1.0/Sd.ali + 1.0/Sd.alo)
	}
}

/* ----------------------------------------------------------------- */
// 透過日射、相当外気温度の計算
func Rmexct(Room []ROOM, Nsrf int, Sd []RMSRF, Wd *WDAT, Exs []EXSF, Snbk []SNBK, Qrm []QRM, nday, mt int) {
	var n, nn, ed, Nrm int
	var Fsdw, Idre, Idf, RN float64
	var Qgtn, Qga, Sab, Rab float64
	var Sdn, Sdnx *RMSRF
	var rm *ROOM
	var e *EXSF
	var S *SNBK
	var Tr float64
	var Eo *ELOUT

	if len(Room) == 0 {
		return
	}

	Nrm = Room[0].end

	// 部位ごとの日射吸収比率のスケジュール対応（比率入力部位の日射入射比率初期化）
	for i := 0; i < Nrm; i++ {
		rm = &Room[i]

		// 室内部位の日射吸収比率の計算
		// 2017/12/25毎時計算へ変更
		// 家具の日射吸収割合
		rm.tfsol = 0.0
		if rm.fsolm != nil {
			rm.tfsol = *(rm.fsolm)
		}

		for j := 0; j < rm.N; j++ {
			rsd := &rm.rsrf[j]

			// 床の場合
			if rsd.ble == 'F' || rsd.ble == 'f' {
				// どの部位も日射吸収比率が定義されていない場合
				if rm.Nfsolfix == 0 {
					// 床の日射吸収比率は固定
					rsd.ffix_flg = '*'
					rsd.fsol = new(float64)
					*rsd.fsol = *rm.flrsr * rsd.A / rm.FArea
				}
			}

			// fsolが規定されている部位についてfsolを合計する
			if rsd.ffix_flg == '*' {
				rm.tfsol += *rsd.fsol // fsolの合計値計算
			}
		}
	}
	// 室内部位の日射吸収比率の計算（毎計算ステップへ変更）2017/12/25
	Rmrdshfc(Room, Sd)
	for i := 0; i < Nrm; i++ {
		Q := &Qrm[i]
		rm := &Room[i]

		rm.Qgt = 0.0
		rm.Qsolm = 0.0
		rm.Qsab = 0.0
		rm.Qrnab = 0.0
		Q.Solo = 0.0
		Q.Solw = 0.0
		Q.Asl = 0.0
		Sdn = &Sd[rm.Brs]

		n := rm.Brs
		for nn := 0; nn < rm.N; nn++ {
			Sdn = &Sd[n]

			ed = Sdn.exs
			e = &Exs[ed]
			Sdn.RSsol = 0.0
			Sdn.RSsold = 0.
			Fsdw = 0.0
			Qgtn = 0.0
			Qga = 0.0
			if Sdn.Sname == "" { /*---higuchi 070918---start-*/
				Sdn.Fsdworg = 0.

				sb := Sdn.sb
				if sb >= 0 && e.Cinc > 0.0 {
					S = &Snbk[sb]
					Fsdw = FNFsdw(S.Type, S.ksi, e.Tazm, e.Tprof, S.D, S.W, S.H, S.W1, S.H1, S.W2, S.H2)
					Sdn.Fsdworg = Fsdw
				} else {
					Fsdw = 0.0
				}

				Idre = e.Idre
				Idf = e.Idf
				RN = e.Rn
			} else { /*---higuchi 070918 end--*/ /*--higuchi 070918 start--*/
				Fsdw = Sdn.Fsdw
				//                  Idre = Sdn.Idre ;  090131 higuchi Sdn.Idre が影をすでに考慮していたため、下に変更
				Idre = e.Idre /*--090131 higuchi  --*/
				Idf = Sdn.Idf
				RN = Sdn.rn
			} /*---higuchi 070918 end --*/

			switch Sdn.ble {
			case 'W':
				// 通常窓の場合
				/*--higuchi add--*/
				Glasstga(Sdn.A, Sdn.tgtn, Sdn.Bn,
					e.Cinc, Fsdw, Idre, Idf, &Qgtn, &Qga, Sdn.window.Cidtype, e.Prof, e.Gamma)
				Rab = Sdn.Eo * RN / Sdn.alo

				Sab = Qga / Sdn.A
				Sdn.TeEsol = Sab
				Sdn.TeErn = -Rab
				Sdn.TeEsol = Sab / Sdn.K

				Sdn.Te = Sab/Sdn.K - Rab + Wd.T

				// 開口部の透過日射熱取得
				Sdn.Qgt = Qgtn
				// 開口部の吸収日射熱取得
				Sdn.Qga = Qga
				// 開口部の夜間放射熱取得
				Sdn.Qrn = -Rab

				rm.Qgt += Qgtn
				rm.Qsab += Sab * Sdn.A
				rm.Qrnab += Rab * Sdn.A * Sdn.K

				Q.Solw += Sdn.A * (Idre + Idf) /*--higuchi add  --*/
				break

			case 'E', 'F', 'R': // このあたりを参考に修正（相当外気温度の計算）
				if Sdn.typ != 'E' && Sdn.typ != 'e' {
					/*---higuchi add---*/
					Sab = Sdn.as * (Idre*(1.0-Fsdw) + Idf) / Sdn.alo
					Rab = Sdn.Eo * RN / Sdn.alo // 長波長
					/*------------------*/

					Sdn.TeEsol = Sab
					Sdn.TeErn = -Rab

					// 建材一体型空気集熱器のための相当外気温度修正
					if Sdn.rpnl != nil && Sdn.rpnl.Type == 'C' {
						//wall := Sdn.mw.wall
						Sdn.Te = Sdn.Tcole
						Sdn.Iwall = Idre*(1.0-Fsdw) + Idf
					} else {
						Sdn.Te = Sab - Rab + Wd.T
					}

					rm.Qsab += Sab * Sdn.A * Sdn.K
					Sdn.Qga = Sab * Sdn.A * Sdn.K
					rm.Qrnab += Rab * Sdn.A * Sdn.K
					Q.Solo += Sdn.A * (Idre + Idf)
					Q.Asl += Sdn.as * Sdn.A * (Idre + Idf)
				} else {
					Sdn.Te = e.Tearth
					Sdn.TeEsol = 0.0
					Sdn.TeErn = 0.0
				}
				break

			case 'i', 'f', 'c', 'd':
				if Sdn.nxrm < 0 {
					Tr = Sdn.room.Trold
					Eo = Sdn.room.cmp.Elouts[0]
					if Eo.Control == LOAD_SW {
						Tr = Sdn.room.rmld.Tset
					}
					Sdn.Te = Sdn.c*Tr + (1.0-Sdn.c)*Wd.T
				} else {
					Tr = Sdn.nextroom.Trold
					Eo = Sdn.nextroom.cmp.Elouts[0]
					if Eo.Control == LOAD_SW {
						Tr = Sdn.nextroom.rmld.Tset
					}
					Sdn.Te = Sdn.c*Tr + (1.0-Sdn.c)*Wd.T
				}
				Sdn.TeEsol = 0.0
				Sdn.TeErn = 0.0
				break
			}

			n++
		} // 表面ループ

		// 室内部位への入射日射の計算（吸収日射ではない）
		for nn = 0; nn < rm.N; nn++ {
			Sdn = &Sd[rm.Brs+nn]

			// 室内部位への入射日射量の計算
			Sdn.RSsold = rm.Qgt * Sdn.srg
		}
	} // 室ループ終了

	Nrm = Room[0].end
	// 透過日射の室内部位の最終計算（隣接室への日射分配、透過日射のうちガラスから屋外に放熱される分も考慮）
	for i := 0; i < Nrm; i++ {
		rm := &Room[i]

		// 透過間仕切りなど、隣接空間への透過日射分配の計算
		n = rm.Brs
		for nn := 0; nn < rm.N; nn++ {
			Sdn := &Sd[n]
			if Sdn.tnxt > 0. && Sdn.RSsold > 0. {
				Rmnxt := Room[Sdn.nxrm]
				RSsol := Sdn.RSsold * Sdn.tnxt

				// 入射日射×透過率が当該室の透過日射熱取得より減ずる
				rm.Qgt -= RSsol
				if Sdn.nextroom != nil {
					// 外皮でない場合は隣室の透過日射熱取得に透過分を加算
					Rmnxt.Qgt += RSsol
				}
			}

			// 透過日射が入射したときに屋外に放熱されるときには、表面吸収日射はゼロとする
			if Sdn.RStrans == 'y' {
				rm.Qgt -= Sdn.RSsold
				Sdn.RSsol = 0.
			}

			n++
		}
	}

	Nrm = Room[0].end
	for i := 0; i < Nrm; i++ {
		rm := &Room[i]
		Q := &Qrm[i]

		// 室内部位の短波長吸収量の計算
		n = rm.Brs
		for nn := 0; nn < rm.N; nn++ {
			Sdn := &Sd[n]

			Sdn.RS = (Sdn.RSsol*Sdn.A + rm.Hr*Sdn.srh +
				rm.Lr*Sdn.srl + rm.Ar*Sdn.sra + rm.Qeqp*Sdn.eqrd) / Sdn.A

			Sdn.RSin = (rm.Hr*Sdn.srh +
				rm.Lr*Sdn.srl + rm.Ar*Sdn.sra + rm.Qeqp*Sdn.eqrd) / Sdn.A
			Sdn.RSli = rm.Lr * Sdn.srl / Sdn.A

			n++
		}

		// 室の透過日射熱取得を再度積算（透明間仕切りによる隣接空間からの透過日射を考慮するため）
		if rm.rsrnx == 'y' {
			for nn := 0; nn < rm.N; nn++ {
				Sdn := &Sd[rm.Brs+n]
				if Sdn.ble == 'c' || Sdn.ble == 'f' {
					if Sdn.nxn >= 0 {
						Sdn.Te += Sd[Sdn.nxn].RS / Sdn.alo
					}
				}
			}
		}

		Q.Tsol = rm.Qgt
		Q.Asol = rm.Qsab
		Q.Arn = rm.Qrnab
		// 家具の日射吸収量の計算
		rm.Qsolm = 0.
		if rm.fsolm != nil {
			rm.Qsolm = rm.Qgt * rm.Srgm2
		}

	} // 室ループ

	for n := 0; n < Nsrf; n++ {
		Sdn := &Sd[n]
		if Sdn.mwtype == 'C' {
			Sdnx = Sdn.nxsd
			Sdn.Te = (Sdnx.alir*Sdnx.Tmrt + Sdnx.RS) / Sdnx.ali
		}
	}
}

/* ----------------------------------------------------------------- */

// 室の係数、定数項の計算
func Roomcf(nmwall int, mw []MWALL, Nroom int, rooms []ROOM, nrdpnl int, rdpnl []RDPNL, wd *WDAT, exsf *EXSFS) {
	for _, rdpnl := range rdpnl {
		panelwp(&rdpnl)
	}

	// 壁体係数行列の作成（壁体数RMSRF分だけループ）
	RMwlc(nmwall, mw, exsf, wd)

	for i := 0; i < Nroom; i++ {
		room := &rooms[i]

		RMcf(room)
		RMrc(room) // 室の定数項の計算

		room.RMx = room.GRM / DTM
		room.RMXC = room.RMx*room.xrold + (room.HL+room.AL)/Ro

		room.RMt += Ca * room.Gvent
		room.RMC += Ca * room.Gvent * wd.T
		room.RMx += room.Gvent
		room.RMXC += room.Gvent * wd.X
	}

	for _, rdpnl := range rdpnl {
		Panelcf(&rdpnl)
		rdpnl.EPC = Panelce(&rdpnl)
	}
}

/* ----------------------------------------------------------------- */
// 前時刻の室温の入れ替え、OT、MRTの計算
func Rmsurft(nroom int, rooms []ROOM, sd []RMSRF) {
	if rooms == nil {
		return
	}

	// 重み係数が未定義もしくは不適切な数値の場合の対処
	r := 0.5
	if rooms[0].OTsetCwgt != nil && *(rooms[0].OTsetCwgt) >= 0.0 && *(rooms[0].OTsetCwgt) <= 1.0 {
		r = *(rooms[0].OTsetCwgt)
	}

	for _, room := range rooms {
		n := room.N
		brs := room.Brs
		sdr := sd[brs:]

		// 前時刻の温度の入れ替え
		room.mrk = 'C'
		room.Trold = room.Tr
		room.xrold = room.xr

		if room.FunHcap > 0 {
			// 家具の温度の計算
			room.TM = room.FMT*room.Tr + room.FMC
			// 家具の吸放熱量の計算
			if room.CM != nil {
				room.QM = *room.CM * (room.TM - room.Tr)
			}

			room.oldTM = room.TM
		}

		// 室内表面温度の計算
		RMsrt(&room)

		room.Tsav = RTsav(n, sdr)
		room.Tot = r*room.Tr + (1.0-r)*room.Tsav
	}
}

/* ----------------------------------------------------------------- */
// PCM収束計算過程における部位表面温度の計算
func Rmsurftd(Nroom int, _Room []ROOM, Sd []RMSRF) {
	var r float64

	if _Room == nil {
		return
	}

	Room := &_Room[0]

	if Room.OTsetCwgt == nil || *(Room.OTsetCwgt) < 0.0 || *(Room.OTsetCwgt) > 1.0 {
		r = 0.5
	} else {
		r = *(Room.OTsetCwgt)
	}

	for i := 0; i < Nroom; i++ {
		Room := &_Room[i]

		N := Room.N
		brs := Room.Brs
		sd := Sd[brs:]

		// 室内表面温度の計算
		RMsrt(Room)

		Room.Tsav = RTsav(N, sd)
		Room.Tot = r*Room.Tr + (1.0-r)*Room.Tsav
	}
}

/*--------------------------------------------------------------------------------------------*/

// 室の熱取得要素の計算
func Qrmsim(Room []ROOM, Wd *WDAT, Qrm []QRM) {
	var dTM, ro, ca float64
	Nrm := Room[0].end

	for i := 0; i < Nrm; i++ {
		Q := &Qrm[i]
		rm := &Room[i]

		Q.Hums = rm.Hc + rm.Hr
		Q.Light = rm.Lc + rm.Lr
		Q.Apls = rm.Ac + rm.Ar

		Q.Huml = rm.HL
		Q.Apll = rm.AL

		Q.Hgins = Q.Hums + Q.Light + Q.Apls

		Q.Qinfs = ca * rm.Gvent * (Wd.T - rm.Tr)
		Q.Qinfl = ro * rm.Gvent * (Wd.X - rm.xr)
		Q.Qeqp = rm.Qeqp
		Q.Qsto = rm.MRM * (rm.Trold - rm.Tr) / dTM
		Q.Qstol = rm.GRM * ro * (rm.xrold - rm.xr) / dTM

		if rm.AEsch != nil {
			Q.AE = rm.AE * *rm.AEsch
		} else {
			Q.AE = 0.0
		}

		if rm.AGsch != nil {
			Q.AG = rm.AG * *rm.AGsch
		} else {
			Q.AG = 0.0
		}
	}
}

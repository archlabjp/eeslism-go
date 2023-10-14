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

/*   bl_panel.c  */

package eeslism

import "math"

const WPTOLE = 1.0e-10

/*  輻射パネル有効熱容量流量  */

func panelwp(rdpnl *RDPNL) {
	sd := rdpnl.sd[0]
	eo := rdpnl.cmp.Elouts[0]
	wall := sd.mw.wall

	var Kc, Kcd float64
	if wall.chrRinput {
		Kc = sd.dblKc
		Kcd = sd.dblKcd
	} else {
		Kc = wall.Kc
		Kcd = wall.Kcd
	}

	if eo.Control != OFF_SW && rdpnl.cmp.Elins[0].Upv != nil {
		rdpnl.cG = eo.G * Spcheat(eo.Fluid)

		if wall.WallType == WallType_P {
			rdpnl.Wp = rdpnl.cG * rdpnl.effpnl / sd.A
		} else {
			rdpnl.Ec = 1.0 - math.Exp(-Kc*sd.A/rdpnl.cG)
			rdpnl.Wp = Kcd * rdpnl.cG * rdpnl.Ec / (Kc * sd.A)
		}
	} else {
		rdpnl.cG = 0.0
		rdpnl.Ec = 0.0
		rdpnl.Wp = 0.0
	}

	// 流量が前時刻から変化していれば係数行列を作りなおす
	if math.Abs(rdpnl.Wp-rdpnl.Wpold) >= WPTOLE || sd.PCMflg {
		rdpnl.Wpold = rdpnl.Wp

		for i := 0; i < rdpnl.MC; i++ {
			rdpnl.sd[i].mrk = '*' // 表面の係数行列再作成
			rdpnl.rm[i].mrk = '*' // 室の係数行列再作成
		}
	}
}

/* -------------------------------------- */

/*  輻射パネル計算用係数    */

func Panelcf(rdpnl *RDPNL) {
	var j, nn, m, mp, M, iup, nrp, n, N int
	var alr, epr, epw *float64
	var ew, kd float64
	var rm *ROOM
	var Sd, Sdd *RMSRF
	var wall *WALL
	var Mw *MWALL
	var C1 float64

	if rdpnl.Wp > 0.0 {
		for m = 0; m < rdpnl.MC; m++ {
			Sd = rdpnl.sd[m]
			rm = rdpnl.rm[m]
			N = rm.N
			nrp = m
			nn = N * nrp

			if Sd.mrk == '*' || Sd.PCMflg {
				if m == 0 {
					Mw = Sd.mw
					mp = Mw.mp
					M = Mw.M

					iup = mp * M

					rdpnl.FIp[m] = Mw.UX[iup] * Mw.uo
					if Mw.wall.WallType == WallType_P { // 通常の床暖房パネル
						rdpnl.FOp[m] = Mw.UX[iup+M-1] * Mw.um
					} else if Mw.wall.WallType == WallType_C { // 屋根一体型空気集熱器
						rdpnl.FOp[m] = Mw.UX[iup+M-1] * Sd.ColCoeff
					}
					rdpnl.FPp = Mw.UX[iup+mp] * Mw.Pc * rdpnl.Wp
				} else {
					Mw = Sd.mw
					rdpnl.FIp[1] = rdpnl.FOp[0]
					rdpnl.FOp[1] = rdpnl.FIp[0]
				}

				wall = Mw.wall
				C1 = Sd.alic
				for j = 0; j < N; j++ {
					alr = &rm.alr[nn+j]
					Sdd = &rm.rsrf[j]
					if j != nrp {
						C1 += *alr * Sdd.WSR
					}
				}
				C1 *= rdpnl.FIp[m] / Sd.ali

				if wall.WallType == WallType_P { // 床暖房パネル
					rdpnl.EPt[m] = C1 * rdpnl.Wp * Sd.A
				} else { // 屋根一体型空気集熱器
					if wall.chrRinput { // 集熱器の特性が熱抵抗で入力されている場合
						kd = Sd.kd
					} else {
						kd = wall.kd
					}
					rdpnl.EPt[m] = C1 * rdpnl.cG * rdpnl.Ec * kd
				}

				for j = 0; j < rm.Ntr; j++ {
					epr = &rdpnl.EPR[m][j]

					*epr = 0.0
					for n = 0; n < N; n++ {
						alr = &rm.alr[nn+n]
						Sdd = &rm.rsrf[n]

						if n != nrp {
							*epr += *alr * Sdd.WSRN[j]
						}
					}
					if wall.WallType == WallType_P {
						*epr *= rdpnl.FIp[m] / Sd.ali * rdpnl.Wp * Sd.A
					} else {
						if wall.chrRinput {
							kd = Sd.kd
						} else {
							kd = wall.kd
						}

						*epr *= rdpnl.FIp[m] / Sd.ali * rdpnl.cG * rdpnl.Ec * kd
						//*epr *= rdpnl.FIp[m] / Sd.ali * rdpnl.cG * rdpnl.Ec * wall.KdKo ;
					}

					/*********
					*epr += rdpnl.FOp[m] * Sd.nxsd.alic / Sd.nxsd.ali;
					***********/
				}
				if wall.WallType == WallType_P { // 通常の床暖房パネル
					rdpnl.Epw = rdpnl.Wp * Sd.A * (1.0 - rdpnl.FPp)
				} else { // 屋根一体型空気集熱器
					if wall.chrRinput {
						kd = Sd.kd
					} else {
						kd = wall.kd
					}

					rdpnl.Epw = rdpnl.cG * (1.0 - rdpnl.Ec*(1.-kd*rdpnl.FPp))
				}
				//}
				//else
				//	rdpnl.Epw = 1. - rdpnl.Ec * ( 1. + wall.KdKo * Sd.FP ) ;

				for j = 0; j < rm.Nrp; j++ {
					epw = &rdpnl.EPW[m][j]

					ew = 0.0
					for n = 0; n < N; n++ {
						alr = &rm.alr[nn+n]
						Sdd = &rm.rsrf[n]

						if n != nrp {
							ew += *alr * Sdd.WSPL[j]
						}
					}

					if wall.WallType == WallType_P {
						*epw = rdpnl.Wp * Sd.A * rdpnl.FIp[m] * ew / Sd.ali
					} else {
						if wall.chrRinput {
							kd = Sd.kd
						} else {
							kd = wall.kd
						}

						*epw = rdpnl.cG * rdpnl.FIp[m] * ew / Sd.ali * rdpnl.Ec * kd
					}
				}
			}
		}
	} else {
		rdpnl.Epw = 0.0
		for m = 0; m < rdpnl.MC; m++ {
			rm = rdpnl.rm[m]
			rdpnl.EPt[m] = 0.0

			for j = 0; j < rm.Ntr; j++ {
				epr = &rdpnl.EPR[m][j]
				*epr = 0.0
			}

			for j = 0; j < rm.Nrp; j++ {
				epw = &rdpnl.EPW[m][j]
				*epw = 0.0
			}
		}
	}
}

/* -------------------------------------------- */

/*  輻射パネルの外乱に関する項の計算     */

func Panelce(rdpnl *RDPNL) float64 {
	var N int
	var rm *ROOM
	var Sd, Sdd *RMSRF
	var Mw *MWALL
	var j, nn, m, mp, M, iup, nrp int
	var CFp, C, CC, kd, ku float64
	var alr *float64
	var wall *WALL

	Sd = nil
	Mw = nil
	wall = nil
	CC = 0.0

	if rdpnl.Wp > 0.0 {
		for m = 0; m < rdpnl.MC; m++ {
			Sd = rdpnl.sd[m]

			if m == 0 {
				Mw = Sd.mw
				mp = Mw.mp
				M = Mw.M
				wall = Mw.wall

				iup = mp * M
				CFp = 0.0
				for j = 0; j < M; j++ {
					CFp += Mw.UX[iup+j] * Mw.Told[j]
				}

				CC = CFp
				if Mw.wall.WallType == WallType_C {
					if Mw.wall.chrRinput {
						kd = Sd.kd
					} else {
						kd = Mw.wall.kd
					}
					CC = CFp * kd
				}
				if rdpnl.MC == 1 {
					if Mw.wall.WallType == WallType_P {
						CC += rdpnl.FOp[m] * Sd.Te
					} else {
						if wall.chrRinput {
							kd = Sd.kd
							ku = Sd.ku
						} else {
							kd = Mw.wall.kd
							ku = Mw.wall.ku
						}
						CC += (ku + kd*rdpnl.FOp[m]) * Sd.Tcoleu
					}
				}
			}

			rm = rdpnl.rm[m]
			N = rm.N
			nrp = m
			nn = N * nrp

			C = 0.0
			for j = 0; j < N; j++ {
				alr = &rm.alr[nn+j]
				Sdd = &rm.rsrf[j]
				if j != nrp {
					C += *alr * Sdd.WSC
				}
			}

			if Mw.wall.WallType == WallType_P {
				CC += rdpnl.FIp[m] * (Sd.RS + C) / Sd.ali
			} else {
				if wall.chrRinput {
					kd = Sd.kd
					ku = Sd.ku
				} else {
					kd = Mw.wall.kd
					ku = Mw.wall.ku
				}
				CC += kd * rdpnl.FIp[m] * (Sd.RS + C) / Sd.ali
			}
		}

		if Mw.wall.WallType == WallType_P {
			return (CC * rdpnl.Wp * Sd.A)
		} else {
			return (CC * rdpnl.cG * rdpnl.Ec)
		}
	} else {
		return (0.0)
	}
}

/* --------------------------- */

/* 負荷計算用設定値のポインター */

func rdpnlldsptr(load *rune, key []string, Rdpnl *RDPNL, vptr *VPTR, idmrk *byte) int {
	err := 0

	if key[1] == "Tout" {
		vptr.Ptr = &Rdpnl.Toset
		vptr.Type = VAL_CTYPE
		Rdpnl.Loadt = load
		*idmrk = 't'
	} else {
		err = 1
	}

	return err
}

/* ------------------------------------------ */

/* 負荷計算用設定値のスケジュール設定 */

func rdpnlldsschd(Rdpnl *RDPNL) {
	Eo := Rdpnl.cmp.Elouts[0]

	if Rdpnl.Loadt != nil {
		if Eo.Control != OFF_SW {
			if Rdpnl.Toset > TEMPLIMIT {
				Eo.Control = ON_SW
				//Eo.Control = LOAD_SW
				//Eo.Sysv = Rdpnl.Toset
			} else {
				Eo.Control = OFF_SW
			}
		}
	}
}

/* ------------------------------------- */

/*  屋根一体型集熱器内部変数のポインター  */

func rdpnlvptr(key []string, Rdpnl *RDPNL, vptr *VPTR) int {
	err := 0

	if Rdpnl.sd[0].mw.wall.WallType == WallType_C && key[1] == "Te" {
		vptr.Ptr = &Rdpnl.sd[0].Tcole
		vptr.Type = VAL_CTYPE
	} else {
		err = 1
	}

	return err
}

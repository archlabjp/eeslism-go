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

/* bl_roomene.c */

package eeslism

import (
	"fmt"
	"io"
	"math"
)

// 室温・湿度計算結果代入、室供給熱量計算
// およびパネル入口温度代入、パネル供給熱量計算
func Roomene(Rmvls *RMVLS, Room []*ROOM, Rdpnl []*RDPNL, Exsfs *EXSFS, Wd *WDAT) {
	var j int
	var E *ELOUT

	for _, rm := range Room {
		E = rm.cmp.Elouts[0]
		rm.Tr = E.Sysv
		E = rm.cmp.Elouts[1]
		rm.xr = E.Sysv
		rm.RH = FNRhtx(rm.Tr, rm.xr)
		rm.hr = FNH(rm.Tr, rm.xr)

		if rm.Arsp != nil {

			for j = 0; j < rm.Nasup; j++ {
				A := rm.Arsp[j]
				elin := rm.elinasup[j]
				elix := rm.elinasupx[j]

				if elin.Lpath.Control != 0 {
					A.G = elin.Lpath.G
					A.Tin = elin.Sysvin
					A.Xin = elix.Sysvin
				} else {
					A.G = 0.0
					A.Tin = 0.0
					A.Xin = 0.0
				}

				if elin.Lpath.Control != 0 {
					A.Qs = Ca * elin.Lpath.G * (elin.Sysvin - rm.Tr)
				} else {
					A.Qs = 0.0
				}

				A.Ql = 0.0
				if elix.Lpath != nil && elix.Lpath.Control != 0 {
					A.Ql = Ro * elix.Lpath.G * (elix.Sysvin - rm.xr)
				}

				A.Qt = A.Qs + A.Ql
			}
		}
	}

	for i := range Rdpnl {
		var Sd *RMSRF
		var Wall *WALL
		Sd = Rdpnl[i].sd[0]
		Wall = Sd.mw.wall
		if Rdpnl[i].cmp.Control != 0 {
			E := Rdpnl[i].cmp.Elouts[0]
			Rdpnl[i].Tpi = Rdpnl[i].cmp.Elins[0].Sysvin
			Rdpnl[i].Tpo = E.Sysv
			cG := Rdpnl[i].cG
			Rdpnl[i].Q = cG * (E.Sysv - Rdpnl[i].Tpi)

			if Wall.WallType == WallType_C {
				var Kc float64
				if Wall.chrRinput {
					Kc = Sd.dblKc
				} else {
					Kc = Wall.Kc
				}

				ECG := cG * Rdpnl[i].Ec / (Kc * Sd.A)
				Rdpnl[i].sd[0].Tf = (1.0-ECG)*Sd.Tcole + ECG*Rdpnl[i].Tpi

				if Sd.Ndiv > 0 {
					for k := 0; k < Sd.Ndiv; k++ {
						Ec := 1.0 - math.Exp(-Kc*Sd.A*float64(k+1)/float64(Sd.Ndiv)/cG)
						Sd.Tc[k] = (1.0-Ec)*Rdpnl[i].Tpi + Ec*Sd.Tcole
					}
				}
			}
		} else {
			Rdpnl[i].Q = 0.0
			Rdpnl[i].Tpi = 0.0
			Rdpnl[i].sd[0].Tf = 0.0

			if Sd.Ndiv > 0 {
				for k := 0; k < Sd.Ndiv; k++ {
					Sd.Tc[k] = 0.0
				}
			}
		}
	}
}

// PCM内蔵壁体の収束判定
func PCMwlchk(counter int, Rmvls *RMVLS, Exsfs *EXSFS, Wd *WDAT, LDreset *int) {
	var Rmwlcreset int

	Rmwlcreset = 0
	// 室温の仮計算
	for i := range Rmvls.Room {
		Rm := Rmvls.Room[i]
		Eo := Rm.cmp.Elouts[0]
		Rm.Tr = Eo.Sysv
	}

	// 部位の表面温度の計算
	Rmsurftd(Rmvls.Room, Rmvls.Sd)

	// 壁体内部温度の仮計算
	RMwltd(Rmvls.Mw)

	// PCM温度の収束判定
	for i := range Rmvls.Room {
		Rm := Rmvls.Room[i]

		// 部位でのループ
		for j := 0; j < Rm.N; j++ {
			Sd := Rm.rsrf[j]
			if Sd.PCMflg {
				mw := Sd.mw
				Wall := mw.wall
				Told := mw.Told
				Toldd := mw.Toldd
				Twd := mw.Twd
				if Sd.mwside == RMSRFMwSideType_i {
					for m := 0; m < mw.M; m++ {
						pcmstate := Sd.pcmstate[m]
						PCM := Wall.PCMLyr[m]
						PCM1 := Wall.PCMLyr[m+1]
						T := 0.0
						Toldn := 0.0
						PCMresetR := 0
						PCMresetL := 0
						nWeightR := -999.0
						nWeightL := -999.0
						if PCM != nil && PCM.Iterate {
							pcmstate.TempPCMave = (Twd[m-1] + Twd[m]) * 0.5
							pcmstate.TempPCMNodeL = Twd[m-1]
							pcmstate.TempPCMNodeR = Twd[m]
							ToldPCMave := (Told[m-1] + Told[m]) * 0.5
							//ToldPCMNodeL := Told[m-1]
							ToldPCMNodeR := Told[m]
							if PCM.AveTemp == 'y' {
								T = pcmstate.TempPCMave
								Toldn = ToldPCMave
							} else {
								T = pcmstate.TempPCMNodeR
								Toldn = ToldPCMNodeR
							}
							if PCM.Spctype == 'm' {
								pcmstate.CapmR = FNPCMStatefun(PCM.Ctype, PCM.Cros, PCM.Crol,
									PCM.Ql, PCM.Ts, PCM.Tl, PCM.Tp, Toldn, T, PCM.DivTemp, &PCM.PCMp)
							} else {
								pcmstate.CapmR = FNPCMstate_table(&PCM.Chartable[0], Toldn, T, PCM.DivTemp)
							}
							if math.Abs(pcmstate.CapmR-pcmstate.OldCapmR) > pcmstate.OldCapmR*PCM.IterateJudge ||
								(PCM.IterateTemp && math.Abs(Twd[m]-Toldd[m]) > 1e-2) {
								nWeightR = PCM.NWeight
								PCMresetR = 1
							}
							pcmstate.OldCapmR = pcmstate.CapmR
						}

						if PCM1 != nil && PCM1.Iterate {
							pcmstate_p1 := Sd.pcmstate[m+1]
							pcmstate_p1.TempPCMave = (Twd[m] + Twd[m+1]) * 0.5
							pcmstate_p1.TempPCMNodeL = Twd[m]
							pcmstate_p1.TempPCMNodeR = Twd[m+1]
							ToldPCMave := (Told[m] + Told[m+1]) * 0.5
							ToldPCMNodeL := Told[m]
							//ToldPCMNodeR := Told[m+1]
							if PCM1.AveTemp == 'y' {
								T = pcmstate_p1.TempPCMave
								Toldn = ToldPCMave
							} else {
								T = pcmstate_p1.TempPCMNodeL
								Toldn = ToldPCMNodeL
							}
							if PCM1.Spctype == 'm' {
								pcmstate_p1.CapmL = FNPCMStatefun(PCM1.Ctype, PCM1.Cros, PCM1.Crol,
									PCM1.Ql, PCM1.Ts, PCM1.Tl, PCM1.Tp, Toldn, T, PCM1.DivTemp, &PCM1.PCMp)
							} else {
								pcmstate_p1.CapmL = FNPCMstate_table(&PCM1.Chartable[0], Toldn, T, PCM1.DivTemp)
							}
							if math.Abs(pcmstate_p1.CapmL-pcmstate_p1.OldCapmL) > pcmstate_p1.OldCapmL*PCM1.IterateJudge ||
								(PCM1.IterateTemp && math.Abs(Twd[m]-Toldd[m]) > 1e-2) {
								nWeightL = PCM1.NWeight
								PCMresetL = 1
							}
							pcmstate_p1.OldCapmL = pcmstate_p1.CapmL
						}

						if PCMresetR+PCMresetL != 0 {
							var nWeight float64
							if nWeightR > 0.0 && nWeightL > 0.0 {
								nWeight = (nWeightR + nWeightL) / 2.0
							} else {
								nWeight = math.Max(nWeightR, nWeightL)
							}
							// Update Toldd[m] with the average value between the previous and current step
							Toldd[m] = ((1.0 - nWeight) * Toldd[m]) + (nWeight * Twd[m])
							(*LDreset)++
							Rmwlcreset++
						}
					}
				}
			}
		}
	}

	if Rmwlcreset > 0 {
		Roomcf(Rmvls.Mw, Rmvls.Room, Rmvls.Rdpnl, Wd, Exsfs)
	}
}

// PCM内蔵家具のPCM温度収束判定
func PCMfunchk(Room []*ROOM, Wd *WDAT, LDreset *int) {
	//var intI int
	var tempTM float64

	for intI := range Room {
		if Room[intI].PCM != nil {
			tempTM = Room[intI].TM
			Room[intI].TM = Room[intI].FMT*Room[intI].Tr + Room[intI].FMC
			if math.Abs(tempTM-Room[intI].TM) > 1e-2 && Room[intI].PCM.Iterate {
				(*LDreset)++
				if Room[intI].PCM.NWeight > 0.0 {
					Room[intI].TM = tempTM*(1.0-Room[intI].PCM.NWeight) + Room[intI].TM*Room[intI].PCM.NWeight
				} else {
					Room[intI].TM = (tempTM + Room[intI].TM) / 2.0
				}
				FunCoeff(Room[intI])

				if Room[intI].FunHcap > 0.0 {
					dblTemp := DTM / Room[intI].FunHcap
					Room[intI].FMC = 1.0 / (dblTemp**Room[intI].CM + 1.0) * (Room[intI].oldTM + dblTemp*Room[intI].Qsolm)
				} else {
					Room[intI].FMC = 0.0
				}

				Room[intI].RMC = Room[intI].MRM/DTM*Room[intI].Trold + Room[intI].HGc + Room[intI].CA
				if Room[intI].FunHcap > 0.0 {
					Room[intI].RMC += *Room[intI].CM * Room[intI].FMC
				}
				Room[intI].RMt += Ca * Room[intI].Gvent
				Room[intI].RMC += Ca * Room[intI].Gvent * Wd.T
			}
		}
	}
}

/* ------------------------------------------------------ */

// 室負荷の計算
func Roomload(Room []*ROOM, LDreset *int) {
	var reset, resetl int

	for i := range Room {
		rm := Room[i]
		if rm.rmld != nil {
			rmld := rm.rmld
			rmld.Qs = 0.0
			rmld.Ql = 0.0
			rmld.Qt = 0.0
			Eo := rm.cmp.Elouts[0]

			if Eo.Control == LOAD_SW {
				rmld.Qs = rm.RMt*rm.Tr - rm.RMC

				for j := 0; j < rm.Ntr; j++ {
					trnx := rm.trnx[j]
					arn := rm.ARN[j]
					rmld.Qs -= arn * trnx.nextroom.Tr
				}
				for j := 0; j < rm.Nrp; j++ {
					rmpnl := rm.rmpnl[j]
					rmp := rm.RMP[j]
					rmld.Qs -= rmp * rmpnl.pnl.Tpi
				}

				if rm.Arsp != nil {
					for j := 0; j < rm.Nasup; j++ {
						A := rm.Arsp[j]
						rmld.Qs -= A.Qs
					}
				}

				for j := 0; j < rm.Nachr; j++ {
					achr := rm.achr[j]
					rmld.Qs -= Ca * achr.Gvr * (Room[achr.rm].Tr - rm.Tr)
				}

				reset = rmloadreset(rmld.Qs, *rmld.loadt, Eo, ON_SW)
				if reset != 0 {
					(*LDreset)++
					//fmt.Printf("resetting...\n")
				}
			}

			Eo = rm.cmp.Elouts[1]
			if Eo.Control == LOAD_SW {
				rmld.Ql = Ro * (rm.RMx*rm.xr - rm.RMXC)
				if rm.Arsp != nil {
					for j := 0; j < rm.Nasup; j++ {
						A := rm.Arsp[j]
						rmld.Ql -= A.Ql
					}
				}

				for j := 0; j < rm.Nachr; j++ {
					achr := rm.achr[j]
					rmld.Ql -= Ro * achr.Gvr * (Room[achr.rm].xr - rm.xr)
				}

				resetl = rmloadreset(rmld.Ql, *rmld.loadt, Eo, ON_SW)
				if reset != 0 || resetl != 0 {
					Eo.Control = ON_SW
					Eo.Eldobj.Sysld = 'n'
					(*LDreset)++
					//fmt.Printf("resetting...\n")
				}
			}

			rmld.Qt = rmld.Qs + rmld.Ql
		}
	}
}

/* ------------------------------------------------------ */

/* 室供給熱量の出力 */

func rmqaprint(fo io.Writer, id int, Rooms []*ROOM) {
	var Nload, Nfnt int
	//var rpnl *RPANEL

	switch id {
	case 0:
		if len(Rooms) > 0 {
			fmt.Fprintf(fo, "%s %d\n", ROOM_TYPE, len(Rooms))
		}
		for _, Room := range Rooms {
			if Room.rmld != nil {
				Nload = 2
			} else {
				Nload = 0
			}

			Nfnt = 0
			if Room.FunHcap > 0.0 {
				Nfnt = 4
			}

			Nset := 0
			if Room.setpri {
				Nset = 1
			}
			fmt.Fprintf(fo, " %s 5 %d 4 %d %d %d\n", Room.Name,
				4+Nload+Room.Nasup*5+Room.Nrp+Nfnt+Nset,
				Nload, Room.Nasup*5, Room.Nrp)
		}
	case 1:
		for _, Room := range Rooms {
			fmt.Fprintf(fo, "%s_Tr t f %s_xr x f %s_RH r f %s_Ts t f ",
				Room.Name, Room.Name, Room.Name, Room.Name)

			if Room.setpri {
				fmt.Fprintf(fo, "%s_SET* t f ", Room.Name)
			}

			if Room.FunHcap > 0.0 {
				fmt.Fprintf(fo, "%s_TM t f %s_QM q f %s_QMsol q f %s_PCMQl q f ", Room.Name, Room.Name, Room.Name, Room.Name)
			}

			if Room.rmld != nil {
				fmt.Fprintf(fo, "%s_Ls q f %s_Lt q f ", Room.Name, Room.Name)
			}

			if Room.Nasup > 0 {
				Ei := Room.cmp.Elins[Room.Nachr+Room.Nrp]
				for j := 0; j < Room.Nasup; j++ {
					if Ei.Lpath == nil {
						fmt.Fprintf(fo, "%s:%1d_G m f %s:%1d_Tin t f %s:%1d_Xin x f %s:%1d_Qas q f %s:%1d_Qat q f ",
							Room.Name, j, Room.Name, j, Room.Name, j, Room.Name, j, Room.Name, j)
					} else {
						fmt.Fprintf(fo, "%s:%s_G m f %s:%s_Tin t f %s:%s_Xin x f %s:%s_Qas q f %s:%s_Qat q f ",
							Room.Name, Ei.Lpath.Name, Room.Name, Ei.Lpath.Name, Room.Name, Ei.Lpath.Name, Room.Name, Ei.Lpath.Name,
							Room.Name, Ei.Lpath.Name)
					}
				}
			}

			for j := 0; j < Room.Nrp; j++ {
				rpnl := Room.rmpnl[j]
				fmt.Fprintf(fo, "%s:%s_Qp q f ", Room.Name, rpnl.pnl.Name)
			}

			fmt.Fprintf(fo, "\n")
		}
	default:
		for _, Room := range Rooms {
			fmt.Fprintf(fo, "%.2f %5.4f %2.0f %.2f ",
				Room.Tr, Room.xr, Room.RH, Room.Tsav)

			if Room.setpri {
				fmt.Fprintf(fo, "%.2f ", Room.SET)
			}

			if Room.FunHcap > 0.0 {
				fmt.Fprintf(fo, "%.2f %.2f %.2f %.2f ", Room.TM, Room.QM, Room.Qsolm, Room.PCMQl)
			}

			if Room.rmld != nil {
				fmt.Fprintf(fo, "%.2f %.2f ", Room.rmld.Qs, Room.rmld.Qt)
			}

			if Room.Nasup > 0 {
				for j := 0; j < Room.Nasup; j++ {
					A := Room.Arsp[j]
					fmt.Fprintf(fo, "%.4g %4.1f %.4f %.2f %.2f ", A.G, A.Tin, A.Xin, A.Qs, A.Qt)
				}
			}

			for j := 0; j < Room.Nrp; j++ {
				rpnl := Room.rmpnl[j]
				fmt.Fprintf(fo, " %.2f", -rpnl.pnl.Q)
			}

			fmt.Fprintf(fo, "\n")
		}
	}
}

/* ------------------------------------------------------ */

/* 放射パネルに関する出力 */

func panelprint(fo io.Writer, id int, Rdpnl []*RDPNL) {
	var ld int
	var Wall *WALL

	switch id {
	case 0:
		if len(Rdpnl) > 0 {
			fmt.Fprintf(fo, "%s %d\n", RDPANEL_TYPE, len(Rdpnl))
		}
		for i := range Rdpnl {
			Sd := Rdpnl[i].sd[0]
			Wall = Sd.mw.wall
			if Sd.mw.wall.WallType == WallType_P {
				fmt.Fprintf(fo, " %s 1 5\n", Rdpnl[i].Name)
			} else {
				ld = 0
				if Wall.chrRinput {
					ld = 5
				}
				if Rdpnl[i].sd[0].PVwallFlg {
					fmt.Fprintf(fo, " %s 1 %d\n", Rdpnl[i].Name, Sd.Ndiv+11+ld)
				} else {
					fmt.Fprintf(fo, " %s 1 %d\n", Rdpnl[i].Name, Sd.Ndiv+8+ld)
				}
			}
		}
	case 1:
		for i := range Rdpnl {
			Sd := Rdpnl[i].sd[0]
			Wall = Sd.mw.wall
			if Sd.mw.wall.WallType == WallType_P {
				fmt.Fprintf(fo, "%s_c c c %s_G m f %s_Ti t f %s_To t f %s_Q q f\n",
					Rdpnl[i].Name, Rdpnl[i].Name, Rdpnl[i].Name, Rdpnl[i].Name, Rdpnl[i].Name)
			} else {
				if !Sd.PVwallFlg {
					fmt.Fprintf(fo, "%s_c c c %s_G m f %s_Ti t f %s_To t f %s_Te t f %s_Tf t f %s_Q q f %s_S q f\n",
						Rdpnl[i].Name, Rdpnl[i].Name, Rdpnl[i].Name, Rdpnl[i].Name, Rdpnl[i].Name, Rdpnl[i].Name, Rdpnl[i].Name, Rdpnl[i].Name)
				} else {
					fmt.Fprintf(fo, "%s_c c c %s_G m f %s_Ti t f %s_To t f %s_Te t f %s_Tf t f %s_Q q f %s_S q f  %s_TPV t f  %s_Iw  q f  %s_P q f\n",
						Rdpnl[i].Name, Rdpnl[i].Name, Rdpnl[i].Name, Rdpnl[i].Name, Rdpnl[i].Name, Rdpnl[i].Name, Rdpnl[i].Name,
						Rdpnl[i].Name, Rdpnl[i].Name, Rdpnl[i].Name, Rdpnl[i].Name)

					if Wall.chrRinput {
						fmt.Fprintf(fo, "%s_Ksu q f %s_Ksd q f %s_Kc q f %s_Tsu t f %s_Tsd t f\n", Rdpnl[i].Name, Rdpnl[i].Name, Rdpnl[i].Name, Rdpnl[i].Name, Rdpnl[i].Name)
					}

					if Sd.Ndiv > 0 {
						for k := 0; k < Sd.Ndiv; k++ {
							fmt.Fprintf(fo, "%s_Tc[%d] t f  ", Rdpnl[i].Name, k+1)
						}
						fmt.Fprintf(fo, "\n")
					}
				}
			}
		}
	default:
		for i := range Rdpnl {
			Sd := Rdpnl[i].sd[0]
			Wall = Sd.mw.wall
			if Sd.mw.wall.WallType == WallType_P {
				fmt.Fprintf(fo, "%c %g  %4.1f %4.1f %3.0f\n", Rdpnl[i].cmp.Elouts[0].Control,
					Rdpnl[i].cmp.Elouts[0].G, Rdpnl[i].Tpi, Rdpnl[i].cmp.Elouts[0].Sysv, Rdpnl[i].Q)
			} else {
				Eo := Rdpnl[i].cmp.Elouts[0]
				G := 0.0
				Wall = Sd.mw.wall
				if Eo.Control != OFF_SW {
					G = Eo.G
				}

				fmt.Fprintf(fo, "%c %g  %4.1f %4.1f %4.1f %4.1f %3.0f %3.0f  ", Eo.Control,
					G, Rdpnl[i].Tpi, Eo.Sysv, Sd.Tcole, Sd.Tf, Rdpnl[i].Q, Sd.Iwall*Sd.A)

				if Sd.PVwallFlg {
					fmt.Fprintf(fo, "%4.1f %4.0f %3.0f\n", Sd.PVwall.TPV, Sd.Iwall, Sd.PVwall.Power)
				} else {
					fmt.Fprintf(fo, "\n")
				}

				if Wall.chrRinput {
					fmt.Fprintf(fo, "%.3f %.3f %.3f %.1f %.1f\n", Sd.dblKsu, Sd.dblKsd, Sd.dblKc, Sd.dblTsu, Sd.dblTsd)
				}

				if Sd.Ndiv > 0 {
					for k := 0; k < Sd.Ndiv; k++ {
						fmt.Fprintf(fo, "%4.1f ", Sd.Tc[k])
					}
					fmt.Fprintf(fo, "\n")
				}
			}
		}
	}
}

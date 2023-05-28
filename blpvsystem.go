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

package main

import "math"

// 太陽電池の温度補正係数計算式
func FNKPT(TPV, apmax float64) float64 {
	return 1.0 + apmax*(TPV-25.0)/100.0
}

// 太陽電池パラメータの初期化
func PVwallcatinit(PVwallcat *PVWALLCAT) {
	PVwallcat.Type = 'C'
	PVwallcat.Apmax = -0.41
	PVwallcat.KHD = 1.0
	PVwallcat.KPD = 0.95
	PVwallcat.KPM = 0.94
	PVwallcat.KPA = 0.97
	PVwallcat.EffINO = 0.9
	PVwallcat.Ap = 10.0
	PVwallcat.Rcoloff = -999.0
	PVwallcat.Kcoloff = -999.0
}

// 温度補正係数以外は時々刻々変化しないので、最初に１度計算しておく
func PVwallPreCalc(PVwallcat *PVWALLCAT) {
	PVwallcat.KConst = PVwallcat.KHD * PVwallcat.KPD * PVwallcat.KPM * PVwallcat.KPA * PVwallcat.EffINO
}

// 太陽電池温度の計算
func FNTPV(Sd *RMSRF, Wd *WDAT, Exsfs *EXSFS) float64 {
	wall := Sd.mw.wall
	Exs := &Exsfs.Exs[Sd.exs]
	Ipv := (wall.tra - Sd.PVwall.Eff) * Sd.Iwall

	var TPV float64
	if Sd.rpnl != nil && Sd.rpnl.cG > 0.0 {
		TPV = (wall.PVwallcat.Ap*Sd.Tf + *Exs.Alo*Wd.T + Ipv) / (wall.PVwallcat.Ap + *Exs.Alo)
	} else {
		TPV = (wall.PVwallcat.Kcoloff*Sd.oldTx + *Exs.Alo*Wd.T + Ipv) / (wall.PVwallcat.Kcoloff + *Exs.Alo)
	}

	return TPV
}

func CalcPowerOutput(Nsrf int, Sd []RMSRF, Wd *WDAT, Exsfs *EXSFS) {
	for i := 0; i < Nsrf; i++ {
		if Sd[i].mw != nil {
			wall := Sd[i].mw.wall

			/// 太陽電池が設置されているときのみ
			if Sd[i].PVwallFlg == 'Y' {
				pvwall := &Sd[i].PVwall

				pvwall.TPV = FNTPV(&Sd[i], Wd, Exsfs)
				pvwall.KPT = FNKPT(pvwall.TPV, wall.PVwallcat.Apmax)
				pvwall.KTotal = wall.PVwallcat.KConst * pvwall.KPT

				pvwall.Power = math.Max(Sd[i].PVwall.PVcap*pvwall.KTotal*Sd[i].Iwall/1000.0, 0.0)

				// 発電効率の計算
				if Sd[i].Iwall > 0.0 {
					pvwall.Eff = pvwall.Power / (Sd[i].Iwall * Sd[i].A)
				} else {
					pvwall.Eff = 0.0
				}
			}
		}
	}
}

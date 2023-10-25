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

/*   eschdlr_s.c   */

package eeslism

/* -------------------------------------------------------------------------- */

func Eeschdlr(day, ttmm int, Schdl *SCHDL, Rmvls *RMVLS) {
	//r := Rmvls.Room

	for j := range Schdl.Sch {
		Schdl.Val[j] = schval(day, ttmm, &Schdl.Sch[j], Schdl.Dsch)
	}

	for j := range Schdl.Scw {
		Schdl.Isw[j] = scwmode(day, ttmm, &Schdl.Scw[j], Schdl.Dscw)
	}

	if SIMUL_BUILDG {
		if DEBUG {
			xprschval(Schdl.Val, Schdl.Isw)
		}

		Windowschdlr(Schdl.Isw, Rmvls.Window, Rmvls.Nsrf, Rmvls.Sd)
		Vtschdlr(Rmvls.Room)
		Aichschdlr(Schdl.Val, Rmvls.Room)

		if DEBUG {
			xprqin(Rmvls.Room)
			xprvent(Rmvls.Room)
		}
	}
}

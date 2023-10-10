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

	for j := 0; j < Schdl.Nsch; j++ {
		val := &Schdl.Val[j]
		Sch := &Schdl.Sch[j]
		*val = schval(day, ttmm, Sch, Schdl.Dsch)
	}

	for j := 0; j < Schdl.Nscw; j++ {
		Scw := &Schdl.Scw[j]
		isw := &Schdl.Isw[j]
		*isw = rune(scwmode(day, ttmm, Scw, Schdl.Dscw))
	}

	if SIMUL_BUILDG {
		if DEBUG {
			xprschval(Schdl.Nsch, Schdl.Val, Schdl.Nscw, Schdl.Isw)
		}

		Windowschdlr(Schdl.Isw, Rmvls.Window, Rmvls.Nsrf, Rmvls.Sd)
		Vtschdlr(Rmvls.Nroom, Rmvls.Room)
		Aichschdlr(Schdl.Val, Rmvls.Nroom, Rmvls.Room)

		if DEBUG {
			xprqin(Rmvls.Nroom, Rmvls.Room)
			xprvent(Rmvls.Nroom, Rmvls.Room)
		}
	}
}

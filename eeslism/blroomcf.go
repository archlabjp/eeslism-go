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

/*   bl_roomcf.c  */

package eeslism

// 熱伝達率の計算

func eeroomcf(Wd *WDAT, Exs *EXSFS, Rmvls *RMVLS, nday int, mt int) {
	// 熱伝達率の計算

	// 表面熱伝達率（対流・放射））の計算
	Rmhtrcf(Exs, Rmvls.Emrk, Rmvls.Room, Rmvls.Sd, Wd)

	if DEBUG {
		// 表面熱伝達率の表示
		xpralph(Rmvls.Room, Rmvls.Sd)
	}

	// 熱貫流率の計算
	Rmhtrsmcf(Rmvls.Sd)

	// 透過日射、相当外気温度の計算
	Rmexct(Rmvls.Room, Rmvls.Sd, Wd, Exs.Exs, Rmvls.Snbk, Rmvls.Qrm, nday, mt)

	// 室の係数（壁体熱伝導等））、定数項の計算
	Roomcf(Rmvls.Mw, Rmvls.Room, Rmvls.Rdpnl, Wd, Exs)
}

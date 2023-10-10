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

/* cnvrg.c  */

package eeslism

/* 合流要素 */

func Cnvrgcfv(Ncnvrg int, Cnvrg []*COMPNT) {
	for i := 0; i < Ncnvrg; i++ {
		C := Cnvrg[i]
		E := C.Elouts[0]

		// 経路が停止していなければ
		if E.Control != OFF_SW {
			E.Coeffo = E.G
			E.Co = 0.0

			if C.Elins[0].Lpath != nil {
				for j := 0; j < C.Nin; j++ {
					cfin := &E.Coeffin[j]
					I := C.Elins[j]

					*cfin = -I.Lpath.G
				}
			}
		}
	}
}

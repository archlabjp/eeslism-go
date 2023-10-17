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

// --------------------------------------------------------------
// 合流要素
//
//  [IN] ---> +---+
//            | C +---->[OUT]
//  [IN] ---> +---+
// --------------------------------------------------------------

// 合流要素 Cnvrg の 出口の係数 Coeffo, Co と入口の係数 Coeffin を計算
func Cnvrgcfv(Cnvrg []*COMPNT) {
	for i := range Cnvrg {
		C := Cnvrg[i]
		E := C.Elouts[0]

		// 経路が停止している場合
		if E.Control == OFF_SW {
			continue
		}

		// 出口係数の処理
		E.Coeffo = E.G
		E.Co = 0.0

		// 入口係数の処理
		if C.Elins[0].Lpath != nil {
			for j := 0; j < C.Nin; j++ {
				E.Coeffin[j] = -C.Elins[j].Lpath.G
			}
		}
	}
}

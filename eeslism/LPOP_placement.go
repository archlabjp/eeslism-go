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

/*

 LPとOPの位置をCGで確認するための入力ファイルを作成する
 FILE=LPOP_placement.c
 Create Date=2006.11.4
*/
package eeslism

import (
	"fmt"
	"os"
)

func HOUSING_PLACE(lpn, mpn int, lp, mp []P_MENN, RET string) {

	mlpn := lpn + mpn

	NAMAE1 := RET + "_placeLP.gchi"
	fp1, err := os.Create(NAMAE1)
	if err != nil {
		fmt.Println("File not open _placeLP.gchi")
		os.Exit(1)
	}
	defer fp1.Close()

	NAMAE2 := RET + "_placeOP.gchi"
	fp2, err := os.Create(NAMAE2)
	if err != nil {
		fmt.Println("File not open _placeOP.gchi")
		os.Exit(1)
	}
	defer fp2.Close()

	NAMAE3 := RET + "_placeALL.gchi"
	fp3, err := os.Create(NAMAE3)
	if err != nil {
		fmt.Println("File not open _placeALL.gchi")
		os.Exit(1)
	}
	defer fp3.Close()

	fmt.Fprintf(fp1, "%d ", lpn)
	fmt.Fprintf(fp3, "%d ", mlpn)
	for i := 0; i < lpn; i++ {
		fmt.Fprintf(fp1, "%s %d\n", lp[i].opname, lp[i].polyd)
		fmt.Fprintf(fp3, "%s %d\n", lp[i].opname, lp[i].polyd)
		fmt.Fprintf(fp1, "%f %f %f\n", lp[i].rgb[0], lp[i].rgb[1], lp[i].rgb[2])
		fmt.Fprintf(fp3, "%f %f %f\n", lp[i].rgb[0], lp[i].rgb[1], lp[i].rgb[2])
		for j := 0; j < lp[i].polyd; j++ {
			fmt.Fprintf(fp1, "%f %f %f\n", lp[i].P[j].X, lp[i].P[j].Y, lp[i].P[j].Z)
			fmt.Fprintf(fp3, "%f %f %f\n", lp[i].P[j].X, lp[i].P[j].Y, lp[i].P[j].Z)
		}
	}

	fmt.Fprintf(fp2, "%d ", mpn)
	for i := 0; i < mpn; i++ {
		fmt.Fprintf(fp2, "%s %d\n", mp[i].opname, mp[i].polyd)
		fmt.Fprintf(fp3, "%s %d\n", mp[i].opname, mp[i].polyd)
		fmt.Fprintf(fp2, "%f %f %f\n", mp[i].rgb[0], mp[i].rgb[1], mp[i].rgb[2])
		fmt.Fprintf(fp3, "%f %f %f\n", mp[i].rgb[0], mp[i].rgb[1], mp[i].rgb[2])
		for j := 0; j < mp[i].polyd; j++ {
			fmt.Fprintf(fp2, "%f %f %f\n", mp[i].P[j].X, mp[i].P[j].Y, mp[i].P[j].Z)
			fmt.Fprintf(fp3, "%f %f %f\n", mp[i].P[j].X, mp[i].P[j].Y, mp[i].P[j].Z)
		}
	}
}

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

   出力ファイル作成
                FILE=ERRPRINT.c
                Create Date 1999.10.26

*/

package main

import (
	"fmt"
	"io"
	"os"
)

/*------計算中の面から見える面と見えない面を判別する際に必要となる
  ベクトルの値を出力する-------------------------*/

func errbekt_print(n, m int, a []bekt, name string) {
	var i, j int
	var fp *os.File

	name += ".gchi"

	fp, err := os.Create(name)
	if err != nil {
		fmt.Println("File not open errbektPrintf")
		os.Exit(1)
	}

	fmt.Fprintf(fp, "%s ", name)

	for i = 0; i < n; i++ {
		for j = 0; j < m; j++ {
			fmt.Fprintf(fp, "%f %f %f %f\n", a[i].ps[j][0], a[i].ps[j][1], a[i].ps[j][2], a[i].ps[j][3])
		}
	}
	fmt.Fprintf(fp, "\n")
	fp.Close()
}

/*--形態係数の出力--*/
func ffactor_printf(fp4 io.Writer, mpn, lpn int, mp, lp []P_MENN, Mon, Day int) {
	var i, j, k int

	k = mpn + lpn

	fmt.Fprintf(fp4, "%d/%d\n", Mon, Day)
	for i = 0; i < mpn; i++ {
		fmt.Fprintf(fp4, "%s\nFAIA %f\nFAIG %f\n", mp[i].opname, mp[i].faia, mp[i].faig)
		for j = 0; j < k; j++ {
			if j < mpn && mp[i].faiwall[j] != 0.0 {
				fmt.Fprintf(fp4, "%s %f\n", mp[j].opname, mp[i].faiwall[j])
			} else if j >= mpn && mp[i].faiwall[j] != 0.0 {
				fmt.Fprintf(fp4, "%s %f\n", lp[j-mpn].opname, mp[i].faiwall[j])
			}
		}
	}
}

/*-----------法線ベクトルの出力-----------*/
func e_printf(n int, p []P_MENN, name string) {
	var i int
	var fp *os.File

	name += ".gchi"

	fp, err := os.Create(name)
	if err != nil {
		fmt.Println("File not open ePrintf")
		os.Exit(1)
	}

	for i = 0; i < n; i++ {
		fmt.Fprintf(fp, "%d x=%f y=%f z=%f\n", i, p[i].e.X, p[i].e.Y, p[i].e.Z)
	}

	fp.Close()
}

/*---------------影面積の出力--------------------*/
func shadow_printf(fp *os.File, M, D int, mt float64, mpn int, mp []P_MENN) {
	var i int

	fmt.Fprintf(fp, "%d %d %5.2f", M, D, mt)

	for i = 0; i < mpn; i++ {
		fmt.Fprintf(fp, " %1.2f", mp[i].sum)
	}

	fmt.Fprintf(fp, "\n")
}

/*------MP面の情報出力----------------------*/
func mp_printf(n int, mp []P_MENN, name string) {
	var i int
	var fp *os.File

	name += ".gchi"

	fp, err := os.Create(name)
	if err != nil {
		fmt.Println("File not open mpPrintf")
		os.Exit(1)
	}

	for i = 0; i < n; i++ {
		fmt.Fprintf(fp, "name=%s wlflg=%d\n", mp[i].opname, mp[i].wlflg)
		fmt.Fprintf(fp, "mp[%d]    wb=%f    wa=%f   ref=%f\n", i, mp[i].wb, mp[i].wa, mp[i].ref)
		fmt.Fprintf(fp, "         e.X=%f   e.Y=%f   e.Z=%f\n", mp[i].e.X, mp[i].e.Y, mp[i].e.Z)
		fmt.Fprintf(fp, "       grp.X=%f grp.Y=%f grp.Z=%f\n", mp[i].grp.X, mp[i].grp.Y, mp[i].grp.Z)
		fmt.Fprintf(fp, "         G.X=%f   G.Y=%f   G.Z=%f\n", mp[i].G.X, mp[i].G.Y, mp[i].G.Z)
		fmt.Fprintf(fp, "\n")
	}

	fp.Close()
}

/*------MP面の前面地面のポイント座標出力---------------------*/
func gp_printf(gp [][]XYZ, mp []P_MENN, mpn, lpn int, name string) {
	var i, k int
	var fp *os.File

	name += ".gchi"

	fp, err := os.Create(name)
	if err != nil {
		fmt.Println("File not open gpPrintf")
		os.Exit(1)
	}

	for i = 0; i < mpn; i++ {
		fmt.Fprintf(fp, "mp[%d] %s\n", i, mp[i].opname)
		k = 0
		for gp[i][k].X != -999 {
			fmt.Fprintf(fp, "%d %f %f %f\n", k, gp[i][k].X, gp[i][k].Y, gp[i][k].Z)
			k++
		}
	}

	fp.Close()
}

/*------LP面の情報出力---------------------------*/
func lp_printf(n int, lp []P_MENN, name string) {
	var i int
	var fp *os.File

	name += ".gchi"

	fp, err := os.Create(name)
	if err != nil {
		fmt.Println("File not open lpPrintf")
		os.Exit(1)
	}

	for i = 0; i < n; i++ {
		fmt.Fprintf(fp, "mp[%d]    wb=%f    wa=%f   ref=%f\n", i, lp[i].wb, lp[i].wa, lp[i].ref)
		fmt.Fprintf(fp, "         e.X=%f   e.Y=%f   e.Z=%f\n", lp[i].e.X, lp[i].e.Y, lp[i].e.Z)
		fmt.Fprintf(fp, "         G.X=%f   G.Y=%f   G.Z=%f\n", lp[i].G.X, lp[i].G.Y, lp[i].G.Z)
		fmt.Fprintf(fp, "\n")
	}

	fp.Close()
}

/*--------------LP面毎の日射遮蔽率出力--------------*/
func lp_shad_printf(lpn int, lp []P_MENN, name string) {
	var fp *os.File
	var i, j int

	name += ".gchi"

	fp, err := os.Create(name)
	if err != nil {
		fmt.Println("File not open lpShadPrintf")
		os.Exit(1)
	}

	for i = 0; i < lpn; i++ {
		fmt.Fprintf(fp, "%s ", lp[i].opname)
	}
	fmt.Fprintf(fp, "\n")

	for i = 1; i <= 365; i++ {
		for j = 0; j < lpn; j++ {
			fmt.Fprintf(fp, "%f ", lp[j].shad[i])
		}
		fmt.Fprintf(fp, "\n")
	}

	fp.Close()
}

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

/*
HOUSING_PLACE (Housing Placement for Visualization)

この関数は、日影計算や日射量計算で用いられる「被受照面（LP）」と「主面（MP）」の幾何学的情報を、
CG（コンピュータグラフィックス）ソフトウェアで可視化するためのファイルに出力します。

建築環境工学的な観点:
- **シミュレーションモデルの可視化**: 建物のエネルギーシミュレーションでは、
  複雑な建物形状や周囲の障害物をモデル化します。
  この関数は、LPとMPの座標データを出力することで、
  シミュレーションモデルが正しく構築されているかを視覚的に確認できます。
  これにより、入力データの誤りやモデルの不整合を早期に発見し、
  シミュレーションの信頼性を向上させます。
- **日影の可視化**: LPとMPのデータを用いて、
  太陽位置に応じた日影の形状や範囲をCGで可視化できます。
  これにより、日射遮蔽部材の効果や、
  周囲の建物による日影の影響を直感的に理解できます。
- **日射量分布の可視化**: 各面の法線ベクトルや日射量データと組み合わせることで、
  建物表面への日射量分布をCGで可視化できます。
  これにより、日射熱取得の多い箇所や、
  太陽光発電パネルの最適な配置などを検討できます。
- **データ形式の課題**: コメントに記載されているように、
  この関数は独自形式で出力しているため、
  汎用的なCGソフトウェアで利用するためには、
  Wavefront OBJ形式などへの変換が必要です。

この関数は、建物のエネルギーシミュレーションにおいて、
モデルの検証、日影・日射量分布の分析、
および設計検討のための重要な可視化機能を提供します。
*/
func HOUSING_PLACE(lpn, mpn int, lp, mp []*P_MENN, RET string) {

	mlpn := lpn + mpn

	// LPの位置データ用ファイル fp1
	NAMAE1 := RET + "_placeLP.gchi"
	fp1, err := os.Create(NAMAE1)
	if err != nil {
		fmt.Println("File not open _placeLP.gchi")
		os.Exit(1)
	}
	defer fp1.Close()

	// OPの位置データ用ファイル fp2
	NAMAE2 := RET + "_placeOP.gchi"
	fp2, err := os.Create(NAMAE2)
	if err != nil {
		fmt.Println("File not open _placeOP.gchi")
		os.Exit(1)
	}
	defer fp2.Close()

	// LPとOPの位置データ用ファイル fp3
	NAMAE3 := RET + "_placeALL.gchi"
	fp3, err := os.Create(NAMAE3)
	if err != nil {
		fmt.Println("File not open _placeALL.gchi")
		os.Exit(1)
	}
	defer fp3.Close()

	// LPの位置データの書き込み => fp1, fp3
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

	// OPの位置データの書き込み => fp2, fp3
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

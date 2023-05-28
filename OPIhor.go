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
   日射量、長波長放射量の計算
   FILE=OPIhor.c
   Create Date=1999.10.10

*/

package main

import (
	"fmt"
	"io"
	"math"
)

func OPIhor(
	fp io.Writer,
	fp1 io.Writer,
	lpn int,
	mpn int,
	mp []P_MENN,
	lp []P_MENN,
	Wd *WDAT,
	ullp *bekt,
	ulmp *bekt,
	gp [][]XYZ,
	nday int,
	dayprn bool,
	monten int,
) {
	var ls,
		ms,
		ns,
		//DE,
		Isky,
		Rsky,
		Igref,
		sg,
		co,
		sum,
		sumg,
		numg,
		reff,
		Ig, // Igの初期化を追加　2017/7/25 佐藤
		reffg,
		Esky,
		Px,
		Py,
		Pz,
		S,
		T float64

	//DE = 2.0

	var i, j, k, h int
	var rp int /*--多角形ループ--*/
	//var Fs float64 // 20170426 higuchi add 形態係数

	ls = -Wd.Sw
	ms = -Wd.Ss
	ns = Wd.Sh
	Esky = Sgm * math.Pow((Wd.T+273.15), 4.0)

	//昼
	if ns > 0.0 {
		for i = 0; i < mpn; i++ {
			sum = 0.0
			reff = 0.0
			CINC(mp[i], ls, ms, ns, &co)
			if co < 0.0 {
				co = 0.0
			}
			mp[i].Idre = Wd.Idn * co * (1.0 - mp[i].sum)

			//形態係数考慮
			if monten > 0 {

				/*------mp面からの拡散日射を求める-------------*/
				k = 0
				for j = 0; j < mpn; j++ {
					if mp[i].faiwall[k] > 0.0 {
						CINC(mp[j], ls, ms, ns, &co)
						if co < 0.0 {
							co = 0.0
						}

						mp[j].Ihor = Wd.Idn*co*(1.0-mp[j].sum) + mp[j].faia*Wd.Isky

						//mp[j].Te=Wd.T+mp[j].as*mp[j].Ihor/mp[j].alo ;
						mp[j].Te = Wd.T

						sum = sum + mp[j].ref*mp[j].Ihor*mp[i].faiwall[k]
						reff = reff + mp[j].Eo*Sgm*mp[i].faiwall[k]*math.Pow((mp[j].Te+273.15), 4.0)
					}
					k++
				}

				/*------lp面からの拡散日射を求める--------------*/
				for j = 0; j < lpn; j++ {
					if mp[i].faiwall[k] > 0.0 {
						CINC(lp[j], ls, ms, ns, &co)
						if co > 0.0 {

							/*--2018.1.26 higuchi 付設障害物、外部障害物からの反射日射をやめるため、削除
							SHADOWlp(j, DE, lpn, mpn, ls, ms, ns, &ullp[j],
								&ulmp[j], &lp[j], lp, mp);
							*/
							/*2018.1.26  higuchi add 上記変更より、以下追加*/
							lp[j].sum = 1

							//printf("lp[%d].sum=%f\n", j, lp[j].sum);
						} else {
							lp[j].sum = 1.0
							co = 0.0
						}

						lp[j].Ihor = Wd.Idn*co*(1.0-lp[j].sum) + lp[j].faia*Wd.Isky

						lp[j].Te = Wd.T

						sum = sum + lp[j].ref*lp[j].Ihor*mp[i].faiwall[k]

						reff = reff + 0.9*Sgm*mp[i].faiwall[k]*math.Pow((lp[j].Te+273.15), 4.0)
					}
					k++
				}

				/*------地面からの拡散日射を求める--------------*/

				if mp[i].faig > 0.0 {
					k = 0
					sumg = 0.0
					numg = 0.0

					for {
						if gp[i][k].X == -999 {
							break
						}

						numg = numg + 1.0

						/*--建物自身による地面の影を求める--*/
						for j = 0; j < mpn; j++ {
							KOUTEN(gp[i][k].X, gp[i][k].Y, gp[i][k].Z, ls, ms, ns,
								&Px, &Py, &Pz, mp[j].P[0], mp[j].e)
							rp = mp[j].polyd - 2
							for h = 0; h < rp; h++ { /*--多角形ループ---*/
								INOROUT(Px, Py, Pz, mp[j].P[0], mp[j].P[h+1], mp[j].P[h+2], &S, &T)
								if (S >= 0.0 && T >= 0.0) && ((S + T) <= 1.0) {
									sumg = sumg + 1.0
									goto koko1
								}
							}
						}

						/*--外部障害物による地面の影を求める--*/
						for j = 0; j < lpn; j++ {
							KOUTEN(gp[i][k].X, gp[i][k].Y, gp[i][k].Z, ls, ms, ns,
								&Px, &Py, &Pz, lp[j].P[0], lp[j].e)
							rp = lp[j].polyd - 2
							for h = 0; h < rp; h++ { /*--多角形ループ---*/
								INOROUT(Px, Py, Pz, lp[j].P[0], lp[j].P[h+1], lp[j].P[h+2], &S, &T)
								if (S >= 0.0 && T >= 0.0) && ((S + T) <= 1.0) {
									sumg = sumg + lp[j].shad[nday]
									goto koko1
								}

							}
						}
					koko1:
						k++
					}

					if numg == 0.0 {
						sg = 0.0
					} else {
						sg = sumg / numg
					}

					Ig = Wd.Idn*Wd.Sh*(1.0-sg) + mp[i].grpfaia*Wd.Isky

					mp[i].Teg = Wd.T + 0.7*Ig/23.0

					reffg = 0.9 * Sgm * mp[i].faig * math.Pow((mp[i].Teg+273.15), 4.0)
				} else {
					Ig = 0.0
					reffg = 0.0
				}

				/*---------------------------------------------------------*/
				Isky = Wd.Isky * mp[i].faia

				Igref = mp[i].faig * mp[i].refg * Ig

				mp[i].Idf = Isky + Igref + sum

				mp[i].Iw = mp[i].Idre + mp[i].Idf

				mp[i].rn = Wd.RN * mp[i].faia

				Rsky = mp[i].faia * Wd.Rsky

				mp[i].Reff = Esky - Rsky - reff - reffg
			} else {
				// ↓20170426 higuchi add 条件追加　形態係数を計算しないパターン

				Isky = Wd.Isky * 0.5

				Igref = 0.5 * mp[i].refg * Ig

				mp[i].Idf = Isky + Igref

				mp[i].Iw = mp[i].Idre + mp[i].Idf

				mp[i].rn = Wd.RN * 0.5

				Rsky = 0.5 * Wd.Rsky

				// higuchi add 20170915 地面が漏れていた
				reffg = 0.9 * Sgm * 0.5 * math.Pow((Wd.T+273.15), 4.0)
				mp[i].Reff = Esky - Rsky - reffg
			}

			// ↓ 20170426 higuchi add 条件追加
			if dayprn {
				fmt.Fprintf(fp, "%s %f %f %f %f %f %f\n",
					mp[i].opname, sg, Isky, Igref, sum, mp[i].Idf, mp[i].Idre)
				fmt.Fprintf(fp1, "%s %f %f %f %f %f %f\n", mp[i].opname, Esky, Rsky, reff, reffg, mp[i].Reff, mp[i].rn)
			}

		}
	} else {
		//夜
		for i = 0; i < mpn; i++ {

			reff = 0.0
			mp[i].Idre = 0.0
			mp[i].Idf = 0.0
			mp[i].Iw = 0.0

			// 20170426 higuchi add 形態係数を考慮しないパターン追加 start
			if monten > 0 {
				mp[i].rn = Wd.RN * mp[i].faia

				/*------mp面からの長波長放射を求める---------*/
				k = 0
				for j = 0; j < mpn; j++ {

					if mp[i].faiwall[k] > 0.0 {
						reff = reff + mp[j].Eo*Sgm*mp[i].faiwall[k]*math.Pow((Wd.T+273.15), 4.0)
					}

					k++
				}

				/*------lp面からの長波長放射を求める-----------*/
				for j = 0; j < lpn; j++ {
					if mp[i].faiwall[k] > 0.0 {
						reff = reff + 0.9*Sgm*mp[i].faiwall[k]*math.Pow((Wd.T+273.15), 4.0)
					}

					k++
				}
				/*--------地面からの長波長放射を求める---------------*/

				if mp[i].faig > 0.0 {
					reffg = 0.9 * Sgm * mp[i].faig * math.Pow((Wd.T+273.15), 4.0)
				} else {
					reffg = 0.0
				}

				Rsky = mp[i].faia * Wd.Rsky
				mp[i].Reff = Esky - Rsky - reff - reffg
			} else {
				mp[i].rn = Wd.RN * 0.5
				Rsky = 0.5 * Wd.Rsky
				// higuchi add 20170915 地面が漏れていた
				reffg = 0.9 * Sgm * 0.5 * math.Pow((Wd.T+273.15), 4.0)
				mp[i].Reff = Esky - Rsky - reffg
			}
			// 20170426 higuchi add 形態係数を考慮しないパターン追加 end

			if dayprn {
				fmt.Fprintf(fp1, "%s %f %f %f %f %f %f\n", mp[i].opname, Esky, Rsky, reff, reffg, mp[i].Reff, mp[i].rn)
			}
		}
	}
}

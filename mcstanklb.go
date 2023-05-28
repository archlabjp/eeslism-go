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

/* mcstanklb.c  */

package main

import "math"

const TSTOLE = 0.04

/*  蓄熱槽仮想分割　*/

func stoint(dTM float64, N int, Vol float64, KAside float64, KAtop float64, KAbtm float64,
	dvol, Mdt, KS, Tss, Tssold []float64, Jva, Jvb *int) {

	for i := 0; i < N; i++ {
		dvol[i] = Vol / float64(N)
		Mdt[i] = (Cw * Row * Vol / float64(N)) / dTM
		KS[i] = KAside / float64(N)

		Tss[i] = Tssold[i]
	}

	KS[0] += KAtop
	KS[N-1] += KAbtm

	*Jva = 0
	*Jvb = 0
}

/* ----------------------------------------------------------- */

func stofc(N, Nin int, Jcin, Jcout []int,
	ihex []rune, ihxeff []float64, Jva, Jvb int, Mdt, KS []float64,
	gxr float64, Tenv *float64, Tssold, cGwin, EGwin, B, R, d, fg []float64) {
	N2 := N * N
	for j := 0; j < N2; j++ {
		B[j] = 0.0
	}

	for j := 0; j < N; j++ {
		B[j*N+j] = Mdt[j] + KS[j]
		R[j] = Mdt[j]*Tssold[j] + KS[j]**Tenv
	}

	for j := 0; j < N-1; j++ {
		B[j*N+j+1] = -Mdt[j] * gxr
		B[(j+1)*N+j] = -Mdt[j+1] * gxr
	}

	if Jva >= 0 {
		for j := Jva; j <= Jvb; j++ {
			B[j*N+j+1] = -Mdt[j] * 1.0e6
			B[(j+1)*N+j] = -Mdt[j] * 1.0e6
		}
	}

	for i := 0; i < Nin; i++ {
		Jin := Jcin[i]
		if cGwin[i] > 0.0 {
			B[Jin*N+Jin] += EGwin[i]

			if Jin < Jcout[i] {
				for j := Jin + 1; j <= Jcout[i]; j++ {
					B[j*N+j-1] -= cGwin[i]
				}
			} else if Jin > Jcout[i] {
				for j := Jcout[i]; j < Jin; j++ {
					B[j*N+j+1] -= cGwin[i]
				}
			}
		}
	}

	for j := 1; j < N-1; j++ {
		B[j*N+j] += math.Abs(B[j*N+j-1]) + math.Abs(B[j*N+j+1])
	}

	B[0] += math.Abs(B[1])
	B[N*N-1] += math.Abs(B[N*N-2])

	Matinv(B, N, N, "<stofc>")
	Matmalv(B, R, N, N, d)

	fgIndex := 0
	for k := 0; k < Nin; k++ {
		Jo := Jcout[k]
		if ihex[k] == 'y' {
			d[Jo] *= ihxeff[k]
			for i := 0; i < Nin; i++ {
				Jin := Jcin[i]
				fg[fgIndex] = B[Jo*N+Jin] * EGwin[i] * ihxeff[k]
				if k == i {
					fg[fgIndex] += (1.0 - ihxeff[k])
				}
				fgIndex++
			}
		} else {
			for i := 0; i < Nin; i++ {
				Jin := Jcin[i]
				fg[fgIndex] = B[Jo*N+Jin] * EGwin[i]
				fgIndex++
			}
		}
	}
}

/* -------------------------------------------------------------- */

/*  蓄熱槽水温の計算　*/

func stotss(N, Nin int, Jcin []int, B, R, EGwin, Twin, Tss []float64) {
	for i := 0; i < Nin; i++ {
		Jin := Jcin[i]
		R[Jin] += EGwin[i] * Twin[i]
	}

	Matmalv(B, R, N, N, Tss)
}

/* -------------------------------------------------------------- */

/*  蓄熱槽水温分布の検討　*/

func stotsexm(N int, Tss []float64, Jva, Jvb *int, dtankF []rune, cfcalc *rune) {
	*Jvb = -1
	*Jva = -1

	for j := N - 2; j >= 0; j-- {
		if dtankF[j] == TANK_FULL {
			if Tss[j+1] > (Tss[j] + TSTOLE) {
				*Jvb = j
			}
			if *Jvb >= 0 {
				break
			}
		}
	}

	if *Jvb >= 0 {
		for j := *Jvb - 1; j >= 0; j-- {
			if dtankF[j] == TANK_FULL {
				if Tss[*Jvb+1] > (Tss[j] + TSTOLE) {
					*Jva = j
				}
			}
		}
		if *Jva == -1 {
			*Jva = *Jvb
		}
	}

	if *Jva < 0 {
		*cfcalc = 'n'
	} else {
		*cfcalc = 'y'
	}
}

/*-----------------------------------------------------------------*/

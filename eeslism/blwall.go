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

/*   bl_wall.c   */
package eeslism

import (
	"fmt"
	"strconv"
	"strings"
)

/* ------------------------------------------ */

/*  壁体後退差分計算準備   */

func Walli(Nbm int, W []BMLST, Wl *WALL, pcm []*PCM, Npcm int) {
	// int     i, j, k, m, N, M;
	// double  Rwall, *C, *Rw, CAPwall;
	var BM *BMLST
	var err error
	// WELM	*Welm;
	// PCM		*PCMcat = nil;
	// char	E[SCHAR], *st, *s, *PCMcode;
	// double	*PCMrate;

	Wl.R = nil
	Wl.CAP = nil
	BM = nil

	Rwall := 0.0
	CAPwall := 0.0

	Wl.R = make([]float64, Wl.N)
	Wl.CAP = make([]float64, Wl.N)
	Wl.PCM = make([]*PCM, Wl.N)
	Wl.PCMrate = make([]float64, Wl.N)

	for jj := 0; jj < Wl.N; jj++ {
		Wl.PCM[jj] = nil
		Wl.PCMrate[jj] = -999.0
	}

	var i int
	for i = 0; i < Wl.N; i++ {
		Welm := &Wl.welm[i]
		Rw := &Wl.R[i]
		C := &Wl.CAP[i]
		PCMrate := &Wl.PCMrate[i]

		// PCM内蔵建材かどうかをチェックする
		st := strings.IndexRune(Welm.Code, '(')
		if st != -1 {
			// フォーマット　　code(PCMcode_VolRate)
			code := Welm.Code[:st]
			PCMcode := Welm.Code[st+1:]
			st = strings.IndexRune(PCMcode, '_')
			if st == -1 {
				panic("PCMの含有率が指定されていません。")
			}
			PCMname := PCMcode[:st]                                              // PCM名称
			*PCMrate, err = strconv.ParseFloat(PCMcode[st+1:len(PCMcode)-1], 64) // PCM含有率（Vol）
			if err != nil {
				panic(err)
			}

			// PCMフラグ
			Wl.PCMflg = true

			// codeのコピー
			Welm.Code = code // 基材のコード名のコピー

			// PCMの検索
			for k := 0; k < Npcm; k++ {
				PCMcat := pcm[k]
				if PCMname == PCMcat.Name {
					Wl.PCM[i] = PCMcat
					break
				}
			}
		}

		var k int
		for k = 0; k < Nbm; k++ {
			BM = &W[k]
			if Welm.Code == BM.Mcode {
				break
			}
		}

		if k == Nbm {
			E := fmt.Sprintf("%s", Welm.Code)
			Eprint("<Walli>", E)
		}

		// 熱伝導率、容積比熱のコピー
		Welm.Cond = BM.Cond       // 熱伝導率［W/mK］
		Welm.Cro = BM.Cro * 1000. // 容積比熱［J/m3K］
		if Welm.L > 0.0 {
			*Rw = Welm.L / BM.Cond // Rw：熱抵抗[m2K/W]、L：厚さ［m］
		} else {
			*Rw = 1.0 / BM.Cond
		}

		// 熱容量［J/m2K
		*C = BM.Cro * 1000. * Welm.L

		if BM.Mcode != "ali" && BM.Mcode != "alo" {
			Rwall += *Rw
			CAPwall += *C
		}
	}

	N := i
	Wl.Rwall = Rwall
	Wl.CAPwall = CAPwall

	m := 0
	M := 0

	for i := 0; i < N; i++ {
		Welm := &Wl.welm[i]
		for j := 0; j <= Welm.ND; j++ {
			M++
		}
	}

	Wl.res = make([]float64, M+2)
	Wl.cap = make([]float64, M+2)
	Wl.L = make([]float64, M+2)

	// PCM構造体へのポインタのメモリ確保
	Wl.PCMLyr = make([]*PCM, M+2)
	Wl.PCMrateLyr = make([]float64, M+2)

	for i := 0; i < N; i++ {
		C := &Wl.CAP[i]
		Rw := &Wl.R[i]
		Welm := &Wl.welm[i]
		PCMrate := &Wl.PCMrate[i]

		for j := 0; j <= Welm.ND; j++ {
			Wl.res[m] = *Rw / (float64(Welm.ND) + 1.)
			Wl.cap[m] = *C / (float64(Welm.ND) + 1.)
			if Welm.L > 0.0 {
				Wl.L[m] = Welm.L / (float64(Welm.ND) + 1.)
			}
			if Wl.PCM[i] != nil {
				Wl.PCMLyr[m] = Wl.PCM[i]
				Wl.PCMrateLyr[m] = *PCMrate
			}

			m++
		}
	}
	Wl.M = m - 1

	if Wl.Ip > 0 {
		Wl.mp = Wl.Ip
		for i := 1; i <= Wl.Ip; i++ {
			Welm := &Wl.welm[i]
			Wl.mp += Welm.ND
		}
	} else {
		Wl.mp = -1
	}
}

/* ------------------------------------------------- */

/*  壁体後退差分計算用係数   */

func Wallfdc(M int, mp int, res []float64, cap []float64,
	Wp float64, UX []float64,
	uo *float64, um *float64, Pc *float64, WallType WALLType,
	Sd *RMSRF, Wd *WDAT,
	Exsf *EXSFS, Wall *WALL, Told []float64, Twd []float64, _pcmstate []*PCMSTATE) {
	var PCMf = 0
	// double	Croa;				// 見かけの比熱
	var ToldPCMave, ToldPCMNodeL, ToldPCMNodeR float64

	Ul := make([]float64, M)
	Ur := make([]float64, M)
	captempL := make([]float64, M+1)
	captempR := make([]float64, M+1)

	// 層構成
	for m := 0; m < M; m++ {
		// PCM内蔵床暖房の計算に活用するためcapをコピーして保持
		captempL[m] = cap[m]
		captempR[m] = cap[m+1]
	}

	for m := 0; m < M; m++ {
		PCMrate := Wall.PCMrateLyr[m] // PCM体積含有率
		//Welm := &Wall.welm[m]
		pcmstate := _pcmstate[m]

		capm, capm1, resm, resm1 := 0.0, 0.0, 0.0, 0.0
		PCM := Wall.PCMLyr[m]
		PCM1 := Wall.PCMLyr[m+1]

		if PCM == nil && PCM1 == nil {
			// PCMなしの層
			C := 0.5 * (cap[m] + cap[m+1])
			Ul[m] = DTM / (C * res[m])
			Ur[m] = DTM / (C * res[m+1])
		} else {
			// どちらかにPCMがある場合
			PCMf = 1

			// 相変化温度を考慮した物性値の計算

			// m点の左にPCMがある場合
			if PCM != nil {
				pcmstate.TempPCMave = (Twd[m-1] + Twd[m]) * 0.5
				pcmstate.TempPCMNodeL = Twd[m-1]
				pcmstate.TempPCMNodeR = Twd[m]

				// PCM温度
				var T, Toldn float64
				if PCM.AveTemp == 'y' {
					T = pcmstate.TempPCMave
					Toldn = ToldPCMave
				} else {
					T = pcmstate.TempPCMNodeR
					Toldn = ToldPCMNodeR
				}
				//pcmstate.tempPCM = T;
				// m層の見かけの比熱

				var Croa float64
				if PCM.Spctype == 'm' {
					Croa = FNPCMStatefun(PCM.Ctype, PCM.Cros, PCM.Crol, PCM.Ql, PCM.Ts, PCM.Tl, PCM.Tp, Toldn, T, PCM.DivTemp, &PCM.PCMp)
				} else {
					Croa = FNPCMstate_table(&PCM.Chartable[0], Toldn, T, PCM.DivTemp)
				}
				if Croa < 0.0 {
					fmt.Printf("Croa=%f\n", Croa)
				}

				pcmstate.CapmR = Croa
				capm = Croa * Wall.L[m]

				// m層の熱抵抗（見かけの比熱特性Typeはダミー値0）
				var lamda float64
				if PCM.Condtype == 'm' {
					lamda = FNPCMStatefun(0, PCM.Conds, PCM.Condl, 0., PCM.Ts, PCM.Tl, PCM.Tp, Toldn, T, PCM.DivTemp, &PCM.PCMp)
				} else {
					lamda = FNPCMstate_table(&PCM.Chartable[1], Toldn, T, PCM.DivTemp)
				}
				pcmstate.OldLamdaR = lamda
				pcmstate.LamdaR = lamda
				resm = Wall.L[m] / lamda
			}

			// m点の右にPCMがある場合
			if PCM1 != nil {
				pcmstate1 := _pcmstate[m+1]
				pcmstate1.TempPCMave = (Twd[m] + Twd[m+1]) * 0.5
				pcmstate1.TempPCMNodeL = Twd[m]
				pcmstate1.TempPCMNodeR = Twd[m+1]
				ToldPCMave = (Told[m-1] + Told[m]) * 0.5
				ToldPCMNodeL = Told[m]
				//ToldPCMNodeR := Told[m+1]

				// PCM温度
				var T, Toldn float64
				if PCM1.AveTemp == 'y' {
					T = pcmstate1.TempPCMave
					Toldn = ToldPCMave
				} else {
					T = pcmstate1.TempPCMNodeL
					Toldn = ToldPCMNodeL
				}

				// m層の見かけの比熱
				var Croa float64
				if PCM1.Spctype == 'm' {
					Croa = FNPCMStatefun(PCM1.Ctype, PCM1.Cros, PCM1.Crol, PCM1.Ql, PCM1.Ts, PCM1.Tl, PCM1.Tp, Toldn, T, PCM1.DivTemp, &PCM1.PCMp)
				} else {
					Croa = FNPCMstate_table(&PCM1.Chartable[0], Toldn, T, PCM1.DivTemp)
				}
				if Croa < 0. {
					fmt.Printf("Croa=%f\n", Croa)
				}

				pcmstate1.CapmL = Croa
				capm1 = Croa * Wall.L[m+1]

				// m層の熱抵抗（見かけの比熱特性Typeはダミー値0）
				var lamda float64
				if PCM1.Condtype == 'm' {
					lamda = FNPCMStatefun(0, PCM1.Conds, PCM1.Condl, 0., PCM1.Ts, PCM1.Tl, PCM1.Tp, Toldn, T, PCM1.DivTemp, &PCM1.PCMp)
				} else {
					lamda = FNPCMstate_table(&PCM1.Chartable[1], Toldn, T, PCM1.DivTemp)
				}
				pcmstate1.OldLamdaL = lamda
				pcmstate1.LamdaL = lamda
				resm1 = Wall.L[m+1] / lamda
			}

			// PCMと基材との含有率による重みづけ平均
			PCMrate1 := Wall.PCMrateLyr[m+1] // PCM体積含有率

			captempL[m] = cap[m]*(1.-PCMrate) + capm*PCMrate
			captempR[m] = cap[m+1]*(1.-PCMrate1) + capm1*PCMrate1
			C := 0.5 * (captempL[m] + captempR[m])
			Ul[m] = DTM / (C * (res[m]*(1.-PCMrate) + resm*PCMrate))
			Ur[m] = DTM / (C * (res[m+1]*(1.-PCMrate1) + resm1*PCMrate1))
		}
	}

	for m := 0; m < M; m++ {
		for j := 0; j < M; j++ {
			UX[m*M+j] = 0.
		}
	}

	UX[0] = 1.0 + Ul[0] + Ur[0]
	for m := 1; m < M; m++ {
		UX[m*M+m] = 1.0 + Ul[m] + Ur[m] // 対角要素
		UX[m*M+m-1] = -Ul[m]            // 対角要素の下
		UX[(m-1)*M+m] = -Ur[m-1]        // 対角要素の右
	}

	if Wp > 0.0 {
		if WallType == 'P' {
			// 床暖房等放射パネル
			*Pc = DTM / (0.5 * (captempL[mp] + captempR[mp]))
			UX[mp*M+mp] += Wp * *Pc
		} else {
			// 建材一体型空気集熱器の場合
			//double ECG, ECt, CFc;
			//WALL   *Wall;
			//Wall := Sd.mw.wall

			*Pc = DTM / (0.5 * captempL[mp])

			// 境界条件の計算
			ECG, ECt, CFc := 0.0, 0.0, 0.0
			FNBoundarySolarWall(Sd, &ECG, &ECt, &CFc)

			UX[mp*M+mp] = 1. - *Pc*ECt + Ul[mp]

			Sd.ColCoeff = *Pc * CFc
		}
	} else {
		if WallType == 'C' {
			// double ECG, ECt, CFc;
			// WALL   *Wall;

			//Wall := Sd.mw.wall
			*Pc = DTM / (0.5 * captempL[mp])

			// 境界条件の計算
			ECG, ECt, CFc := 0.0, 0.0, 0.0
			FNBoundarySolarWall(Sd, &ECG, &ECt, &CFc)
			UX[mp*M+mp] = 1. - *Pc*ECt + Ul[mp]
			Sd.ColCoeff = *Pc * CFc
		}

		*Pc = 0.0
	}

	*uo = Ul[0]
	*um = Ur[M-1]

	if PCMf == 5 {
		/*************/
		fmt.Printf(" Wallfdc -- U --\n")
		Matprint(" %12.8f", M, UX)
		fmt.Printf("\nuo=%f   um=%f\n", *uo, *um)
		fmt.Printf("mp=%d  Pc=%f\n", mp, *Pc)
		/***********/
	}

	Matinv(UX, M, M, "<Wallfdc>")

	/*********************/
	if PCMf == 5 {
		fmt.Printf("\n Wallfdc-- inv(U) --\n")
		Matprint("%12.8f", M, UX)
	}
	/*********************/
}

/* --------------------------------------------- */

/*  後退差分による壁体表面、内部温度の計算   */
func Twall(M, mp int, UX []float64, uo, um, Pc, Ti, To, WpT float64, Told, Tw []float64, Sd *RMSRF, pcm []*PCM) {
	// 前時刻の壁体内部温度のコピー
	Ttemp := make([]float64, M)
	Toldcalc := make([]float64, M)
	copy(Ttemp, Told)
	copy(Toldcalc, Told)

	Toldcalc[0] += uo * Ti

	if Sd.mw.wall.WallType != WallType_C {
		Toldcalc[M-1] += um * To
	} else {
		Toldcalc[M-1] += Sd.ColCoeff * To
	}

	if Pc > 0.0 {
		Toldcalc[mp] += Pc * WpT
	}

	for m := 0; m < M; m++ {
		Tw[m] = 0.0
		for j := 0; j < M; j++ {
			Tw[m] += UX[m*M+j] * Toldcalc[j]
		}
	}

	// 建材一体型集熱器の集熱器と建材の境界温度
	if Sd.mw.wall.WallType == WallType_C {
		Sd.oldTx = Tw[mp]
	}

	// PCMの温度飛び越えのチェック
	for m := 0; m < M; m++ {
		PCMLyr := pcm[m]
		if PCMLyr != nil {
			// 現在時刻のPCM温度
			Tpcm := (Tw[m-1] + Tw[m]) * 0.5

			// 前時刻のPCM温度
			Tpcmold := (Ttemp[m-1] + Ttemp[m]) * 0.5

			// 壁体温度が潜熱領域をまたいだかチェック
			if PCMLyr.Iterate == false && ((PCMLyr.Ts > Tpcmold && PCMLyr.Tl < Tpcm) || (PCMLyr.Tl < Tpcmold && PCMLyr.Ts > Tpcm)) {
				fmt.Printf("xxxx 壁体温度が潜熱領域をまたぎました Tpcm=%.1f Tpcmold=%.1f\n", Tpcm, Tpcmold)
			}
		}
	}
}

/* --------------------------------------------- */

/*  後退差分による壁体表面、内部温度の計算   */
// PCM収束計算過程のチェック用

func Twalld(M, mp int, UX []float64, uo, um, Pc, Ti, To, WpT float64, Told, Twd []float64, Sd *RMSRF) {
	// 収束計算過程なので、前時刻の計算結果が変わらないようにバックアップ
	Toldtemp := make([]float64, M)
	copy(Toldtemp, Told)

	Toldtemp[0] += uo * Ti

	if Sd.mw.wall.WallType != WallType_C {
		Toldtemp[M-1] += um * To
	} else {
		Toldtemp[M-1] += Sd.ColCoeff * To
	}

	if Pc > 0.0 {
		Toldtemp[mp] += Pc * WpT
	}

	for m := 0; m < M; m++ {
		Twd[m] = 0.0
		for j := 0; j < M; j++ {
			Twd[m] += UX[m*M+j] * Toldtemp[j]
		}
	}
}

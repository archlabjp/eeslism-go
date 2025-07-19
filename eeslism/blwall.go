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

/*
Walli (Wall Initialization)

この関数は、壁体の材料構成、熱的特性、および特殊な機能（PCM内蔵、放射パネルなど）に関する情報を初期化し、
壁体内部の熱伝導計算に必要なパラメータを設定します。

建築環境工学的な観点:
  - **壁体の層構成のモデル化**: 壁体は、複数の材料層で構成されます。
    この関数は、各層の材料コード（`Welm.Code`）、厚さ（`Welm.L`）、
    熱伝導率（`Welm.Cond`）、容積比熱（`Welm.Cro`）を読み込み、
    各層の熱抵抗（`Rw`）と熱容量（`C`）を計算します。
    これにより、壁体全体の熱貫流率（`Wl.Rwall`）と熱容量（`Wl.CAPwall`）を算出できます。
  - **PCM（相変化材料）の考慮**: `Welm.Code`にPCMの指定がある場合、
    そのPCMの名称（`PCMname`）と含有率（`PCMrate`）を読み込み、
    対応するPCMのデータ（`pcm`）を`Wl.PCM`に割り当てます。
    PCMは、潜熱を利用して大きな熱量を蓄えることができ、
    壁体の熱容量を向上させ、室温変動を緩和する効果があります。
  - **壁体内部の分割**: 壁体内部の温度分布を詳細にモデル化するために、
    各層をさらに複数の仮想的な層に分割します（`Welm.ND`）。
    これにより、壁体内部の温度勾配や熱流をより正確に計算できます。
  - **放射パネルの位置 (Wl.Ip, Wl.mp)**:
    壁体に放射パネルが内蔵されている場合、
    その位置（`Wl.Ip`）と、壁体内部のどの仮想層に相当するか（`Wl.mp`）を設定します。
    これにより、放射パネルからの熱供給が壁体内部の温度分布に与える影響をモデル化できます。

この関数は、建物の壁体の熱的挙動を詳細にモデル化し、
熱負荷計算、エネルギー消費量予測、
省エネルギー対策の検討、および快適性評価を行うための重要な初期設定機能を提供します。
*/
func Walli(Nbm int, W []BMLST, Wl *WALL, pcm []*PCM) {
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
		Wl.PCMrate[jj] = FNAN
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
			for k := range pcm {
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

/*
Wallfdc (Wall Finite Difference Coefficient Calculation)

この関数は、壁体内部の熱伝導計算に用いられる後退差分法（Finite Difference Method）の係数行列を作成します。
これは、壁体内部の温度分布や熱流を詳細にモデル化するために不可欠です。

建築環境工学的な観点:
  - **壁体内部の熱伝導モデル**: 壁体は熱容量を持つため、
    その内部温度は外気温度や室内温度の変化に対して時間遅れを伴って応答します。
    後退差分法は、この動的な熱的挙動を数値的に解くための手法です。
    この関数は、各層の熱容量（`cap`）、熱抵抗（`res`）、
    そして時間ステップ（`DTM`）を考慮して、
    壁体内部の熱伝導方程式を表現する係数行列（`UX`）を構築します。
  - **PCM（相変化材料）の考慮**: PCMが内蔵された壁体の場合、
    その相変化特性（見かけの比熱`Croa`、熱伝導率`lamda`）を考慮して、
    各層の熱容量や熱抵抗を動的に変化させます。
    これにより、PCMの潜熱蓄熱効果を正確にモデル化できます。
  - **放射パネルの考慮**: 壁体に放射パネルが内蔵されている場合（`Wp > 0.0`）、
    その熱供給が壁体内部の温度分布に与える影響を係数行列に組み込みます。
  - **集熱器一体型壁の考慮**: `WallType == 'C'` の場合、
    建材一体型空気集熱器の熱的特性（`FNBoundarySolarWall`で計算される境界条件）を考慮して、
    係数行列を構築します。
  - **逆行列の計算**: 構築された係数行列`UX`の逆行列を計算することで、
    壁体内部の温度を直接求めることができるようになります。

この関数は、建物の壁体の熱的挙動を詳細にモデル化し、
熱負荷計算、エネルギー消費量予測、
省エネルギー対策の検討、および快適性評価を行うための重要な役割を果たします。
*/
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

/*
Twall (Wall Temperature Calculation)

この関数は、壁体内部の各層の温度を計算します。
これは、壁体の熱的挙動を詳細にモデル化し、
蓄熱効果や熱貫流特性を評価するために不可欠です。

建築環境工学的な観点:
  - **壁体内部の熱伝導**: 壁体は熱容量を持つため、
    その内部温度は外気温度や室内温度の変化に対して時間遅れを伴って応答します。
    この関数は、壁体内部の熱伝導方程式を解き、
    各層の温度（`Tw`）を計算します。
  - **境界条件の考慮**: 壁体の熱伝達は、室内側（`Ti`）と室外側（`To`）の境界条件に大きく依存します。
    `Toldcalc`は、これらの境界条件と壁体の熱応答係数（`uo`, `um`）を用いて、
    壁体内部の温度変化を計算します。
  - **パネルからの熱影響**: `Pc`が`0.0`より大きい場合、
    壁体に放射パネルなどが組み込まれており、
    そのパネルからの熱影響（`WpT`）が壁体内部の温度に考慮されます。
  - **PCMの温度飛び越えチェック**: `PCMLyr.Iterate == false && ((PCMLyr.Ts > Tpcmold && PCMLyr.Tl < Tpcm) || (PCMLyr.Tl < Tpcmold && PCMLyr.Ts > Tpcm))` の条件は、
    PCM層の温度が潜熱領域を飛び越えて変化した場合に、
    その旨を通知します。
    これは、PCMの相変化が適切にモデル化されているかを確認するために重要です。
  - **温度履歴の更新**: 計算された壁体内部温度は、
    次の時間ステップの計算のために`Told`に更新されます。
    これにより、壁体の熱的履歴が考慮され、
    より正確な動的熱応答のシミュレーションが可能となります。

この関数は、建物の熱的挙動を詳細にモデル化し、
熱負荷計算、エネルギー消費量予測、
省エネルギー対策の検討、および快適性評価を行うための重要な役割を果たします。
*/
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

/*
Twalld (Wall Temperature Temporary Calculation for PCM Convergence)

この関数は、壁体内部の各層の温度を仮計算します。
主に、PCM（相変化材料）が内蔵された壁体において、
PCMの相変化を考慮した壁体内部温度の収束計算の初期段階で用いられます。

建築環境工学的な観点:
  - **PCMの相変化と収束計算**: PCMは、相変化する際に大量の潜熱を吸収・放出するため、
    壁体内部の温度分布に非線形な影響を与えます。
    そのため、壁体内部温度の計算には反復的な収束計算が必要となる場合があります。
    この`Twalld`関数は、その収束計算の初期値や中間段階での仮の温度分布を計算するために使用されます。
  - **境界条件の考慮**: 壁体の熱伝達は、室内側（`Ti`）と室外側（`To`）の境界条件に大きく依存します。
    `Toldtemp`は、これらの境界条件と壁体の熱応答係数（`uo`, `um`）を用いて、
    壁体内部の温度変化を計算します。
  - **パネルからの熱影響**: `Pc`が`0.0`より大きい場合、
    壁体に放射パネルなどが組み込まれており、
    そのパネルからの熱影響（`WpT`）が壁体内部の温度に考慮されます。
  - **温度履歴の考慮**: `Told`は前時刻の壁体内部温度であり、
    現在の温度計算にその履歴が考慮されます。
    これにより、壁体の熱的履歴が考慮され、
    より正確な動的熱応答のシミュレーションが可能となります。

この関数は、PCMを組み込んだ壁体の熱的挙動を正確にモデル化し、
熱負荷平準化、エネルギー消費量予測、
および省エネルギー対策の検討を行うための重要な役割を果たします。
*/
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

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

/*  mcdessicant.C  */
/*  バッチ式デシカント空調機 */

package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

/* ------------------------------------------ */

// 要素方程式の変数のためのメモリの割り当て
func Desielm(Ndesi int, Desi []DESI) {
	for i := 0; i < Ndesi; i++ {
		desi := &Desi[i]
		Eot := desi.Cmp.Elouts[0] // 空気温度出口
		Eox := desi.Cmp.Elouts[1] // 空気湿度出口

		elin := Eot.Elins[1]
		elin.Upo = elin.Upv // 出口空気温度の要素方程式の2つめの変数は絶対湿度
		elin.Upv = Eox.Elins[0].Upo

		elin = Eox.Elins[1]
		elin.Upo = elin.Upv         // 出口絶対湿度の要素方程式の2つめの変数は空気温度
		elin.Upv = Eot.Elins[0].Upo // 空気温度の要素方程式の2つ目の変数（空気入口温度）のupo、upvに空気湿度をつなげる
	}
}

/* 機器仕様入力　　　　　　*/

/*---- Satoh追加 2013/10/20 ----*/
func Desiccantdata(s string, desica *DESICA) int {
	st := strings.Index(s, "=")
	var dt float64
	var id int

	if st == -1 {
		desica.name = s
		desica.r = -999.0
		desica.rows = -999.0
		desica.Uad = -999.0
		desica.A = -999.0
		desica.Vm = 18.0
		desica.eps = 0.4764
		desica.P0 = 0.4
		desica.kp = 0.0012
		desica.cps = 710.0
		desica.ms = -999.0
	} else {
		st++
		dt, _ = strconv.ParseFloat(s[st:], 64)

		switch {
		case s == "Uad": // シリカゲル槽壁面の熱貫流率[W/m2K]
			desica.Uad = dt
		case s == "A": // シリカゲル槽表面積[m2]
			desica.A = dt
		case s == "ms": // シリカゲル質量[g]
			desica.ms = dt
		case s == "r": // シリカゲル平均直径[cm]
			desica.r = dt
		case s == "rows": // シリカゲル充填密度[g/cm3]
			desica.rows = dt
		default:
			id = 1
		}
	}
	return id
}

/* --------------------------- */

/*  管長・ダクト長、周囲温度設定 */

func Desiint(NDesi int, _Desi []DESI, Simc *SIMCONTL, Ncompnt int, Compnt []COMPNT, Wd *WDAT) {
	var Err string
	var Desica *DESICA

	for i := 0; i < NDesi; i++ {
		Desi := &_Desi[i]

		if Desi.Cmp.Envname != "" {
			Desi.Tenv = envptr(Desi.Cmp.Envname, Simc, Ncompnt, Compnt, Wd, nil)
		} else {
			Desi.Room = roomptr(Desi.Cmp.Roomname, Ncompnt, Compnt)
		}

		Desica = Desi.Cat

		if Desica.Uad < 0.0 {
			Err = fmt.Sprintf("Name=%s  Uad=%.4g", Desica.name, Desica.Uad)
			Eprint("Desiint", Err)
		}
		if Desica.A < 0.0 {
			Err = fmt.Sprintf("Name=%s  A=%.4g", Desica.name, Desica.A)
			Eprint("Desiint", Err)
		}
		if Desica.r < 0.0 {
			Err = fmt.Sprintf("Name=%s  r=%.4g", Desica.name, Desica.r)
			Eprint("Desiint", Err)
		}
		if Desica.rows < 0.0 {
			Err = fmt.Sprintf("Name=%s  rows=%.4g", Desica.name, Desica.rows)
			Eprint("Desiint", Err)
		}
		if Desica.ms < 0.0 {
			Err = fmt.Sprintf("Name=%s  ms=%.4g", Desica.name, Desica.ms)
			Eprint("Desiint", Err)
		}

		// 初期温度、出入口温度の初期化
		Desi.Tsold = 20.0
		Desi.Xsold = FNXtr(Desi.Tsold, 50.0)
		Desi.Tain = Desi.Tsold
		Desi.Taout = Desi.Tsold
		Desi.Xain = Desi.Xsold
		Desi.Xaout = Desi.Xsold

		// デシカント槽熱損失係数の計算
		Desi.UA = Desica.Uad * Desica.A

		// 吸湿量の初期化
		Desi.Pold = Desica.P0

		// シリカゲルと槽内空気の熱伝達面積[m2]
		Desi.Asa = 3.0 * Desica.ms * 1000.0 * (1.0 - Desica.eps) / (1.0e4 * (Desica.r / 10.0) * Desica.rows)

		// 逆行列
		Desi.UX = make([]float64, 5*5)
		Desi.UXC = make([]float64, 5)
	}
}

/* --------------------------- */

/*  特性式の係数  */

func Desicfv(NDesi int, Desi []DESI) {
	var Eo *ELOUT
	var h, i, j float64
	var Te, hsa, hsad, hAsa, hdAsa float64
	var Desica *DESICA
	var U, C, Cmat []float64

	N := 5
	N2 := N * N
	for inti := 0; inti < NDesi; inti++ {
		Desi := &Desi[inti]
		Desica = Desi.Cat

		// 係数行列のメモリ確保
		U = make([]float64, N2)
		// 定数行列のメモリ確保
		C = make([]float64, N)

		if Desi.Cmp.Envname != "" {
			Te = *Desi.Tenv
		} else {
			Te = Desi.Room.Tot
		}

		Eo = Desi.Cmp.Elouts[0]
		// 熱容量流量の計算
		Desi.CG = Spcheat(Eo.Fluid) * Eo.G

		// シリカゲルと槽内空気の対流熱伝達率の計算
		if Eo.Cmp.Control == OFF_SW {
			hsa = 4.614
		} else {
			hsa = 40.0
		}

		// シリカゲルと槽内空気の湿気伝達率の計算
		hsad = hsa / Ca

		hAsa = hsa * Desi.Asa
		hdAsa = hsad * Desi.Asa

		if Desi.Pold >= 0.25 {
			h = 0.001319
			i = 0.103335
			j = -0.05416
		} else {
			h = 0.001158
			i = 0.149479
			j = -0.05835
		}

		// 定数行列Cの作成
		Cmat = C
		Cmat[0] = Desica.ms * Desica.cps / DTM * Desi.Tsold
		Cmat[1] = Desi.UA * Te
		Cmat[3] = Desica.ms / DTM * Desi.Pold
		Cmat[4] = -j

		// 係数行列の作成
		U[0*N+0] = Desica.ms*Desica.cps/DTM + hAsa
		U[0*N+1] = -hAsa
		U[0*N+2] = -hdAsa * Ro
		U[0*N+3] = hdAsa * Ro
		U[1*N+0] = -hAsa
		U[1*N+1] = Ca*Eo.G + hAsa + Desi.UA
		U[2*N+2] = Eo.G + hdAsa
		U[2*N+3] = -hdAsa
		U[3*N+2] = -hdAsa
		U[3*N+3] = hdAsa
		U[3*N+4] = Desica.ms / DTM
		U[4*N+0] = h
		U[4*N+3] = -1.0
		U[4*N+4] = i

		// 逆行列の計算
		Matinv(U, N, N, "<Desicfv U>")

		// 行列のコピー
		matinit(Desi.UX, N2)
		matcpy(U, Desi.UX, N2)

		// {UXC}=[UX]*{C}の作成
		matinit(Desi.UXC, N)
		for ii := 0; ii < N; ii++ {
			for jj := 0; jj < N; jj++ {
				Desi.UXC[ii] += Desi.UX[ii*N+jj] * C[jj]
			}
		}

		// 出口温度の要素方程式
		Eo.Coeffo = -1.0
		Eo.Co = -Desi.UXC[1]
		Eo.Coeffin[0] = Desi.UX[1*N+1] * Eo.G * Ca
		Eo.Coeffin[1] = Desi.UX[1*N+2] * Eo.G

		// 出口湿度の要素方程式
		Eo = Desi.Cmp.Elouts[1]
		Eo.Coeffo = -1.0
		Eo.Co = -Desi.UXC[2]
		Eo.Coeffin[0] = Desi.UX[2*N+2] * Eo.G
		Eo.Coeffin[1] = Desi.UX[2*N+1] * Eo.G * Ca
	}
}

///* --------------------------- */
//
///* 取得熱量の計算 */
//
func Desiene(NDesi int, _Desi []DESI) {
	Sin := make([]float64, 5)
	S := make([]float64, 5)

	N := 5
	//N2 := N * N
	for i := 0; i < NDesi; i++ {
		Desi := &_Desi[i]
		matinit(Sin, N)
		matinit(S, N)
		elo := Desi.Cmp.Elouts[0]
		elox := Desi.Cmp.Elouts[1]
		elix := elo.Elins[1]
		Desi.Tain = elo.Elins[0].Sysvin
		Desi.Xain = elix.Sysvin

		var Te float64
		if Desi.Cmp.Envname != "" {
			Te = *Desi.Tenv
		} else {
			Te = Desi.Room.Tot
		}

		Desi.Taout = elo.Sysv
		Desi.Xaout = elox.Sysv

		// 入口状態行列Sinの作成
		Sin[1] = Ca * elo.G * Desi.Tain
		Sin[2] = elo.G * Desi.Xain
		// 内部状態値の計算
		for ii := 0; ii < N; ii++ {
			for jj := 0; jj < N; jj++ {
				S[ii] += Desi.UX[ii*N+jj] * Sin[jj]
			}
			S[ii] += Desi.UXC[ii]
		}
		// 変数への格納
		Desi.Tsold = S[0]
		Desi.Ta = S[1]
		Desi.Xa = S[2]
		Desi.Xsold = S[3]
		Desi.Pold = S[4]
		Desi.RHold = FNRhtx(Desi.Tsold, Desi.Xsold)

		// 顕熱の計算
		Desi.Qs = Desi.CG * (Desi.Taout - Desi.Tain)
		Desi.Ql = elo.G * Ro * (Desi.Xaout - Desi.Xain)
		Desi.Qt = Desi.Qs + Desi.Ql

		// デシカント槽からの熱損失の計算
		Desi.Qloss = Desi.UA * (Te - Desi.Ta)

		// 設置室内部発熱の計算
		if Desi.Room != nil {
			Desi.Room.Qeqp += (-Desi.Qloss)
		}
	}
}

// 制御で使用する内部変数
func Desivptr(key []string, Desi *DESI, vptr *VPTR) int {
	var err int

	switch key[1] {
	case "Ts":
		vptr.Ptr = &Desi.Tsold
		vptr.Type = VAL_CTYPE
	case "xs":
		vptr.Ptr = &Desi.Xsold
		vptr.Type = VAL_CTYPE
	case "RH":
		vptr.Ptr = &Desi.RHold
		vptr.Type = VAL_CTYPE
	default:
		err = 1
	}

	return err
}

///* ---------------------------*/
//
func Desiprint(fo *os.File, id int, Ndesi int, _Desi []DESI) {
	switch id {
	case 0:
		if Ndesi > 0 {
			fmt.Fprintf(fo, "%s %d\n", DESI_TYPE, Ndesi)
		}
		for i := 0; i < Ndesi; i++ {
			Desi := &_Desi[i]
			fmt.Fprintf(fo, " %s 1 14\n", Desi.Name)
		}
	case 1:
		for i := 0; i < Ndesi; i++ {
			Desi := &_Desi[i]
			fmt.Fprintf(fo, "%s_c c c %s_G m f %s_Ts t f %s_Ti t f %s_To t f %s_Qs q f ", Desi.Name, Desi.Name, Desi.Name, Desi.Name, Desi.Name, Desi.Name)
			fmt.Fprintf(fo, "%s_xs x f %s_RHs r f %s_xi x f %s_xo x f %s_Ql q f %s_Qt q f ", Desi.Name, Desi.Name, Desi.Name, Desi.Name, Desi.Name, Desi.Name)
			fmt.Fprintf(fo, "%s_Qls q f %s_P m f\n", Desi.Name, Desi.Name)
		}
	default:
		for i := 0; i < Ndesi; i++ {
			Desi := &_Desi[i]
			fmt.Fprintf(fo, "%c %6.4g %4.1f %4.1f %4.1f %2.0f  ", Desi.Cmp.Elouts[0].Control, Desi.Cmp.Elouts[0].G, Desi.Tsold, Desi.Tain, Desi.Taout, Desi.Qs)
			fmt.Fprintf(fo, "%.3f %.0f %.3f %.3f %2.0f %2.0f  ", Desi.Xsold, Desi.RHold, Desi.Xain, Desi.Xaout, Desi.Ql, Desi.Qt)
			fmt.Fprintf(fo, "%.0f %.3f\n", Desi.Qloss, Desi.Pold)
		}
	}
}

///* --------------------------- */
//
///* 日積算値に関する処理 */
//
///*******************/
func Desidyint(Ndesi int, _Desi []DESI) {
	for i := 0; i < Ndesi; i++ {
		Desi := &_Desi[i]
		svdyint(&Desi.Tidy)
		svdyint(&Desi.Tsdy)
		svdyint(&Desi.Tody)
		svdyint(&Desi.xidy)
		svdyint(&Desi.xsdy)
		svdyint(&Desi.xody)
		qdyint(&Desi.Qsdy)
		qdyint(&Desi.Qldy)
		qdyint(&Desi.Qtdy)
		qdyint(&Desi.Qlsdy)
	}
}

func Desiday(Mon, Day, ttmm, Ndesi int, _Desi []DESI, Nday, SimDayend int) {
	// Mo := Mon - 1
	// tt := ConvertHour(ttmm)

	for i := 0; i < Ndesi; i++ {
		Desi := &_Desi[i]
		// 日集計
		svdaysum(int64(ttmm), Desi.Cmp.Control, Desi.Tain, &Desi.Tidy)
		svdaysum(int64(ttmm), Desi.Cmp.Control, Desi.Taout, &Desi.Tody)
		svdaysum(int64(ttmm), Desi.Cmp.Control, Desi.Tsold, &Desi.Tsdy)
		svdaysum(int64(ttmm), Desi.Cmp.Control, Desi.Xain, &Desi.xidy)
		svdaysum(int64(ttmm), Desi.Cmp.Control, Desi.Xaout, &Desi.xody)
		svdaysum(int64(ttmm), Desi.Cmp.Control, Desi.Xsold, &Desi.xsdy)
		qdaysum(int64(ttmm), Desi.Cmp.Control, Desi.Qs, &Desi.Qsdy)
		qdaysum(int64(ttmm), Desi.Cmp.Control, Desi.Ql, &Desi.Qldy)
		qdaysum(int64(ttmm), Desi.Cmp.Control, Desi.Qt, &Desi.Qtdy)
		qdaysum(int64(ttmm), Desi.Cmp.Control, Desi.Qloss, &Desi.Qlsdy)
	}
}

func Desidyprt(fo *os.File, id, Ndesi int, _Desi []DESI) {
	switch id {
	case 0:
		if Ndesi > 0 {
			fmt.Fprintf(fo, "%s %d\n", DESI_TYPE, Ndesi)
		}
		for i := 0; i < Ndesi; i++ {
			Desi := &_Desi[i]
			fmt.Fprintf(fo, " %s 1 68\n", Desi.Name)
		}
	case 1:
		for i := 0; i < Ndesi; i++ {
			Desi := &_Desi[i]

			fmt.Fprintf(fo, "%s_Ht H d %s_Ti T f ", Desi.Name, Desi.Name)
			fmt.Fprintf(fo, "%s_ttn h d %s_Tin t f %s_ttm h d %s_Tim t f\n", Desi.Name, Desi.Name, Desi.Name, Desi.Name)

			fmt.Fprintf(fo, "%s_Ht H d %s_To T f ", Desi.Name, Desi.Name)
			fmt.Fprintf(fo, "%s_ttn h d %s_Ton t f %s_ttm h d %s_Tom t f\n", Desi.Name, Desi.Name, Desi.Name, Desi.Name)

			fmt.Fprintf(fo, "%s_Ht H d %s_Ts T f ", Desi.Name, Desi.Name)
			fmt.Fprintf(fo, "%s_ttn h d %s_Tsn t f %s_ttm h d %s_Tsm t f\n", Desi.Name, Desi.Name, Desi.Name, Desi.Name)

			fmt.Fprintf(fo, "%s_Ht H d %s_xi T f ", Desi.Name, Desi.Name)
			fmt.Fprintf(fo, "%s_ttn h d %s_xin t f %s_ttm h d %s_xim t f\n", Desi.Name, Desi.Name, Desi.Name, Desi.Name)

			fmt.Fprintf(fo, "%s_Ht H d %s_xo T f ", Desi.Name, Desi.Name)
			fmt.Fprintf(fo, "%s_ttn h d %s_xon t f %s_ttm h d %s_xom t f\n", Desi.Name, Desi.Name, Desi.Name, Desi.Name)

			fmt.Fprintf(fo, "%s_Ht H d %s_xs T f ", Desi.Name, Desi.Name)
			fmt.Fprintf(fo, "%s_ttn h d %s_xsn t f %s_ttm h d %s_xsm t f\n", Desi.Name, Desi.Name, Desi.Name, Desi.Name)

			fmt.Fprintf(fo, "%s_Hsh H d %s_Qsh Q f %s_Hsc H d %s_Qsc Q f\n", Desi.Name, Desi.Name, Desi.Name, Desi.Name)
			fmt.Fprintf(fo, "%s_tsh h d %s_qsh q f %s_tsc h d %s_qsc q f\n", Desi.Name, Desi.Name, Desi.Name, Desi.Name)

			fmt.Fprintf(fo, "%s_Hlh H d %s_Qlh Q f %s_Hlc H d %s_Qlc Q f\n", Desi.Name, Desi.Name, Desi.Name, Desi.Name)
			fmt.Fprintf(fo, "%s_tlh h d %s_qlh q f %s_tlc h d %s_qlc q f\n", Desi.Name, Desi.Name, Desi.Name, Desi.Name)

			fmt.Fprintf(fo, "%s_Hth H d %s_Qth Q f %s_Htc H d %s_Qtc Q f\n", Desi.Name, Desi.Name, Desi.Name, Desi.Name)
			fmt.Fprintf(fo, "%s_tth h d %s_qth q f %s_ttc h d %s_qtc q f\n", Desi.Name, Desi.Name, Desi.Name, Desi.Name)

			fmt.Fprintf(fo, "%s_Hlsh H d %s_Qlsh Q f %s_Hlsc H d %s_Qlsc Q f\n", Desi.Name, Desi.Name, Desi.Name, Desi.Name)
			fmt.Fprintf(fo, "%s_tlsh h d %s_qlsh q f %s_tlsc h d %s_qlsc q f\n", Desi.Name, Desi.Name, Desi.Name, Desi.Name)
		}
	default:
		for i := 0; i < Ndesi; i++ {
			Desi := &_Desi[i]
			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ", Desi.Tidy.Hrs, Desi.Tidy.M, Desi.Tidy.Mntime, Desi.Tidy.Mn, Desi.Tidy.Mxtime, Desi.Tidy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ", Desi.Tody.Hrs, Desi.Tody.M, Desi.Tody.Mntime, Desi.Tody.Mn, Desi.Tody.Mxtime, Desi.Tody.Mx)
			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ", Desi.Tsdy.Hrs, Desi.Tsdy.M, Desi.Tsdy.Mntime, Desi.Tsdy.Mn, Desi.Tsdy.Mxtime, Desi.Tsdy.Mx)

			fmt.Fprintf(fo, "%1d %.4f %1d %.4f %1d %.4f ", Desi.xidy.Hrs, Desi.xidy.M, Desi.xidy.Mntime, Desi.xidy.Mn, Desi.xidy.Mxtime, Desi.xidy.Mx)
			fmt.Fprintf(fo, "%1d %.4f %1d %.4f %1d %.4f ", Desi.xody.Hrs, Desi.xody.M, Desi.xody.Mntime, Desi.xody.Mn, Desi.xody.Mxtime, Desi.xody.Mx)
			fmt.Fprintf(fo, "%1d %.4f %1d %.4f %1d %.4f ", Desi.xsdy.Hrs, Desi.xsdy.M, Desi.xsdy.Mntime, Desi.xsdy.Mn, Desi.xsdy.Mxtime, Desi.xsdy.Mx)

			fmt.Fprintf(fo, "%1d %3.1f ", Desi.Qsdy.Hhr, Desi.Qsdy.H)
			fmt.Fprintf(fo, "%1d %3.1f ", Desi.Qsdy.Chr, Desi.Qsdy.C)
			fmt.Fprintf(fo, "%1d %2.0f ", Desi.Qsdy.Hmxtime, Desi.Qsdy.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f ", Desi.Qsdy.Cmxtime, Desi.Qsdy.Cmx)

			fmt.Fprintf(fo, "%1d %3.1f ", Desi.Qldy.Hhr, Desi.Qldy.H)
			fmt.Fprintf(fo, "%1d %3.1f ", Desi.Qldy.Chr, Desi.Qldy.C)
			fmt.Fprintf(fo, "%1d %2.0f ", Desi.Qldy.Hmxtime, Desi.Qldy.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f ", Desi.Qldy.Cmxtime, Desi.Qldy.Cmx)

			fmt.Fprintf(fo, "%1d %3.1f ", Desi.Qtdy.Hhr, Desi.Qtdy.H)
			fmt.Fprintf(fo, "%1d %3.1f ", Desi.Qtdy.Chr, Desi.Qtdy.C)
			fmt.Fprintf(fo, "%1d %2.0f ", Desi.Qtdy.Hmxtime, Desi.Qtdy.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f ", Desi.Qtdy.Cmxtime, Desi.Qtdy.Cmx)

			fmt.Fprintf(fo, "%1d %3.1f ", Desi.Qlsdy.Hhr, Desi.Qlsdy.H)
			fmt.Fprintf(fo, "%1d %3.1f ", Desi.Qlsdy.Chr, Desi.Qlsdy.C)
			fmt.Fprintf(fo, "%1d %2.0f ", Desi.Qlsdy.Hmxtime, Desi.Qlsdy.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f\n", Desi.Qlsdy.Cmxtime, Desi.Qlsdy.Cmx)
		}
	}
}

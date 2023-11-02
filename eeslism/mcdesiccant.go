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

package eeslism

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

/* ------------------------------------------ */

// 要素方程式の変数のためのメモリの割り当て
func Desielm(Desi []*DESI) {
	for _, desi := range Desi {
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

func Desiint(Desi []*DESI, Simc *SIMCONTL, Compnt []*COMPNT, Wd *WDAT) {
	var Err string
	var Desica *DESICA

	for _, desi := range Desi {

		if desi.Cmp.Envname != "" {
			desi.Tenv = envptr(desi.Cmp.Envname, Simc, Compnt, Wd, nil)
		} else {
			desi.Room = roomptr(desi.Cmp.Roomname, Compnt)
		}

		Desica = desi.Cat

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
		desi.Tsold = 20.0
		desi.Xsold = FNXtr(desi.Tsold, 50.0)
		desi.Tain = desi.Tsold
		desi.Taout = desi.Tsold
		desi.Xain = desi.Xsold
		desi.Xaout = desi.Xsold

		// デシカント槽熱損失係数の計算
		desi.UA = Desica.Uad * Desica.A

		// 吸湿量の初期化
		desi.Pold = Desica.P0

		// シリカゲルと槽内空気の熱伝達面積[m2]
		desi.Asa = 3.0 * Desica.ms * 1000.0 * (1.0 - Desica.eps) / (1.0e4 * (Desica.r / 10.0) * Desica.rows)

		// 逆行列
		desi.UX = make([]float64, 5*5)
		desi.UXC = make([]float64, 5)
	}
}

/* --------------------------- */

/*  特性式の係数  */

//
// 温度 [IN 1] --> +------+ --> [OUT 1] 出口温度
//                 | DESI |
// 湿度 [IN 2] --> +------+ --> [OUT 2] 出口湿度
//
func Desicfv(Desi []*DESI) {
	var Eo1 *ELOUT
	var h, i, j float64
	var Te, hsa, hsad, hAsa, hdAsa float64
	var Desica *DESICA
	var U, C, Cmat []float64

	N := 5
	N2 := N * N
	for _, desi := range Desi {
		Desica = desi.Cat

		// 係数行列のメモリ確保
		U = make([]float64, N2)
		// 定数行列のメモリ確保
		C = make([]float64, N)

		if desi.Cmp.Envname != "" {
			Te = *desi.Tenv
		} else {
			Te = desi.Room.Tot
		}

		Eo1 = desi.Cmp.Elouts[0]
		// 熱容量流量の計算
		desi.CG = Spcheat(Eo1.Fluid) * Eo1.G

		// シリカゲルと槽内空気の対流熱伝達率の計算
		if Eo1.Cmp.Control == OFF_SW {
			hsa = 4.614
		} else {
			hsa = 40.0
		}

		// シリカゲルと槽内空気の湿気伝達率の計算
		hsad = hsa / Ca

		hAsa = hsa * desi.Asa
		hdAsa = hsad * desi.Asa

		if desi.Pold >= 0.25 {
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
		Cmat[0] = Desica.ms * Desica.cps / DTM * desi.Tsold
		Cmat[1] = desi.UA * Te
		Cmat[3] = Desica.ms / DTM * desi.Pold
		Cmat[4] = -j

		// 係数行列の作成
		U[0*N+0] = Desica.ms*Desica.cps/DTM + hAsa
		U[0*N+1] = -hAsa
		U[0*N+2] = -hdAsa * Ro
		U[0*N+3] = hdAsa * Ro
		U[1*N+0] = -hAsa
		U[1*N+1] = Ca*Eo1.G + hAsa + desi.UA
		U[2*N+2] = Eo1.G + hdAsa
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
		matinit(desi.UX, N2)
		matcpy(U, desi.UX, N2)

		// {UXC}=[UX]*{C}の作成
		matinit(desi.UXC, N)
		for ii := 0; ii < N; ii++ {
			for jj := 0; jj < N; jj++ {
				desi.UXC[ii] += desi.UX[ii*N+jj] * C[jj]
			}
		}

		// 出口温度の要素方程式
		Eo1.Coeffo = -1.0
		Eo1.Co = -desi.UXC[1]
		Eo1.Coeffin[0] = desi.UX[1*N+1] * Eo1.G * Ca
		Eo1.Coeffin[1] = desi.UX[1*N+2] * Eo1.G

		// 出口湿度の要素方程式
		Eo2 := desi.Cmp.Elouts[1]
		Eo2.Coeffo = -1.0
		Eo2.Co = -desi.UXC[2]
		Eo2.Coeffin[0] = desi.UX[2*N+2] * Eo2.G
		Eo2.Coeffin[1] = desi.UX[2*N+1] * Eo2.G * Ca
	}
}

///* --------------------------- */
//
///* 取得熱量の計算 */
//
func Desiene(Desi []*DESI) {
	Sin := make([]float64, 5)
	S := make([]float64, 5)

	N := 5
	//N2 := N * N
	for _, desi := range Desi {
		matinit(Sin, N)
		matinit(S, N)
		elo := desi.Cmp.Elouts[0]
		elox := desi.Cmp.Elouts[1]
		elix := elo.Elins[1]
		desi.Tain = elo.Elins[0].Sysvin
		desi.Xain = elix.Sysvin

		var Te float64
		if desi.Cmp.Envname != "" {
			Te = *desi.Tenv
		} else {
			Te = desi.Room.Tot
		}

		desi.Taout = elo.Sysv
		desi.Xaout = elox.Sysv

		// 入口状態行列Sinの作成
		Sin[1] = Ca * elo.G * desi.Tain
		Sin[2] = elo.G * desi.Xain
		// 内部状態値の計算
		for ii := 0; ii < N; ii++ {
			for jj := 0; jj < N; jj++ {
				S[ii] += desi.UX[ii*N+jj] * Sin[jj]
			}
			S[ii] += desi.UXC[ii]
		}
		// 変数への格納
		desi.Tsold = S[0]
		desi.Ta = S[1]
		desi.Xa = S[2]
		desi.Xsold = S[3]
		desi.Pold = S[4]
		desi.RHold = FNRhtx(desi.Tsold, desi.Xsold)

		// 顕熱の計算
		desi.Qs = desi.CG * (desi.Taout - desi.Tain)
		desi.Ql = elo.G * Ro * (desi.Xaout - desi.Xain)
		desi.Qt = desi.Qs + desi.Ql

		// デシカント槽からの熱損失の計算
		desi.Qloss = desi.UA * (Te - desi.Ta)

		// 設置室内部発熱の計算
		if desi.Room != nil {
			desi.Room.Qeqp += (-desi.Qloss)
		}
	}
}

// 制御で使用する内部変数
func Desivptr(key []string, Desi *DESI) (VPTR, error) {
	var err error
	var vptr VPTR

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
		err = errors.New("'Ts', 'xs' or 'RH' is expected")
	}

	return vptr, err
}

///* ---------------------------*/
//
func Desiprint(fo io.Writer, id int, Desi []*DESI) {
	switch id {
	case 0:
		if len(Desi) > 0 {
			fmt.Fprintf(fo, "%s %d\n", DESI_TYPE, len(Desi))
		}
		for _, desi := range Desi {
			fmt.Fprintf(fo, " %s 1 14\n", desi.Name)
		}
	case 1:
		for _, desi := range Desi {
			fmt.Fprintf(fo, "%s_c c c %s_G m f %s_Ts t f %s_Ti t f %s_To t f %s_Qs q f ", desi.Name, desi.Name, desi.Name, desi.Name, desi.Name, desi.Name)
			fmt.Fprintf(fo, "%s_xs x f %s_RHs r f %s_xi x f %s_xo x f %s_Ql q f %s_Qt q f ", desi.Name, desi.Name, desi.Name, desi.Name, desi.Name, desi.Name)
			fmt.Fprintf(fo, "%s_Qls q f %s_P m f\n", desi.Name, desi.Name)
		}
	default:
		for _, desi := range Desi {
			fmt.Fprintf(fo, "%c %6.4g %4.1f %4.1f %4.1f %2.0f  ", desi.Cmp.Elouts[0].Control, desi.Cmp.Elouts[0].G, desi.Tsold, desi.Tain, desi.Taout, desi.Qs)
			fmt.Fprintf(fo, "%.3f %.0f %.3f %.3f %2.0f %2.0f  ", desi.Xsold, desi.RHold, desi.Xain, desi.Xaout, desi.Ql, desi.Qt)
			fmt.Fprintf(fo, "%.0f %.3f\n", desi.Qloss, desi.Pold)
		}
	}
}

///* --------------------------- */
//
///* 日積算値に関する処理 */
//
///*******************/
func Desidyint(Desi []*DESI) {
	for _, desi := range Desi {
		svdyint(&desi.Tidy)
		svdyint(&desi.Tsdy)
		svdyint(&desi.Tody)
		svdyint(&desi.xidy)
		svdyint(&desi.xsdy)
		svdyint(&desi.xody)
		qdyint(&desi.Qsdy)
		qdyint(&desi.Qldy)
		qdyint(&desi.Qtdy)
		qdyint(&desi.Qlsdy)
	}
}

func Desiday(Mon, Day, ttmm int, Desi []*DESI, Nday, SimDayend int) {
	// Mo := Mon - 1
	// tt := ConvertHour(ttmm)

	for _, desi := range Desi {
		// 日集計
		svdaysum(int64(ttmm), desi.Cmp.Control, desi.Tain, &desi.Tidy)
		svdaysum(int64(ttmm), desi.Cmp.Control, desi.Taout, &desi.Tody)
		svdaysum(int64(ttmm), desi.Cmp.Control, desi.Tsold, &desi.Tsdy)
		svdaysum(int64(ttmm), desi.Cmp.Control, desi.Xain, &desi.xidy)
		svdaysum(int64(ttmm), desi.Cmp.Control, desi.Xaout, &desi.xody)
		svdaysum(int64(ttmm), desi.Cmp.Control, desi.Xsold, &desi.xsdy)
		qdaysum(int64(ttmm), desi.Cmp.Control, desi.Qs, &desi.Qsdy)
		qdaysum(int64(ttmm), desi.Cmp.Control, desi.Ql, &desi.Qldy)
		qdaysum(int64(ttmm), desi.Cmp.Control, desi.Qt, &desi.Qtdy)
		qdaysum(int64(ttmm), desi.Cmp.Control, desi.Qloss, &desi.Qlsdy)
	}
}

func Desidyprt(fo io.Writer, id int, Desi []*DESI) {
	switch id {
	case 0:
		if len(Desi) > 0 {
			fmt.Fprintf(fo, "%s %d\n", DESI_TYPE, len(Desi))
		}
		for _, desi := range Desi {
			fmt.Fprintf(fo, " %s 1 68\n", desi.Name)
		}
	case 1:
		for _, desi := range Desi {

			fmt.Fprintf(fo, "%s_Ht H d %s_Ti T f ", desi.Name, desi.Name)
			fmt.Fprintf(fo, "%s_ttn h d %s_Tin t f %s_ttm h d %s_Tim t f\n", desi.Name, desi.Name, desi.Name, desi.Name)

			fmt.Fprintf(fo, "%s_Ht H d %s_To T f ", desi.Name, desi.Name)
			fmt.Fprintf(fo, "%s_ttn h d %s_Ton t f %s_ttm h d %s_Tom t f\n", desi.Name, desi.Name, desi.Name, desi.Name)

			fmt.Fprintf(fo, "%s_Ht H d %s_Ts T f ", desi.Name, desi.Name)
			fmt.Fprintf(fo, "%s_ttn h d %s_Tsn t f %s_ttm h d %s_Tsm t f\n", desi.Name, desi.Name, desi.Name, desi.Name)

			fmt.Fprintf(fo, "%s_Ht H d %s_xi T f ", desi.Name, desi.Name)
			fmt.Fprintf(fo, "%s_ttn h d %s_xin t f %s_ttm h d %s_xim t f\n", desi.Name, desi.Name, desi.Name, desi.Name)

			fmt.Fprintf(fo, "%s_Ht H d %s_xo T f ", desi.Name, desi.Name)
			fmt.Fprintf(fo, "%s_ttn h d %s_xon t f %s_ttm h d %s_xom t f\n", desi.Name, desi.Name, desi.Name, desi.Name)

			fmt.Fprintf(fo, "%s_Ht H d %s_xs T f ", desi.Name, desi.Name)
			fmt.Fprintf(fo, "%s_ttn h d %s_xsn t f %s_ttm h d %s_xsm t f\n", desi.Name, desi.Name, desi.Name, desi.Name)

			fmt.Fprintf(fo, "%s_Hsh H d %s_Qsh Q f %s_Hsc H d %s_Qsc Q f\n", desi.Name, desi.Name, desi.Name, desi.Name)
			fmt.Fprintf(fo, "%s_tsh h d %s_qsh q f %s_tsc h d %s_qsc q f\n", desi.Name, desi.Name, desi.Name, desi.Name)

			fmt.Fprintf(fo, "%s_Hlh H d %s_Qlh Q f %s_Hlc H d %s_Qlc Q f\n", desi.Name, desi.Name, desi.Name, desi.Name)
			fmt.Fprintf(fo, "%s_tlh h d %s_qlh q f %s_tlc h d %s_qlc q f\n", desi.Name, desi.Name, desi.Name, desi.Name)

			fmt.Fprintf(fo, "%s_Hth H d %s_Qth Q f %s_Htc H d %s_Qtc Q f\n", desi.Name, desi.Name, desi.Name, desi.Name)
			fmt.Fprintf(fo, "%s_tth h d %s_qth q f %s_ttc h d %s_qtc q f\n", desi.Name, desi.Name, desi.Name, desi.Name)

			fmt.Fprintf(fo, "%s_Hlsh H d %s_Qlsh Q f %s_Hlsc H d %s_Qlsc Q f\n", desi.Name, desi.Name, desi.Name, desi.Name)
			fmt.Fprintf(fo, "%s_tlsh h d %s_qlsh q f %s_tlsc h d %s_qlsc q f\n", desi.Name, desi.Name, desi.Name, desi.Name)
		}
	default:
		for _, desi := range Desi {
			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ", desi.Tidy.Hrs, desi.Tidy.M, desi.Tidy.Mntime, desi.Tidy.Mn, desi.Tidy.Mxtime, desi.Tidy.Mx)
			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ", desi.Tody.Hrs, desi.Tody.M, desi.Tody.Mntime, desi.Tody.Mn, desi.Tody.Mxtime, desi.Tody.Mx)
			fmt.Fprintf(fo, "%1d %3.1f %1d %3.1f %1d %3.1f ", desi.Tsdy.Hrs, desi.Tsdy.M, desi.Tsdy.Mntime, desi.Tsdy.Mn, desi.Tsdy.Mxtime, desi.Tsdy.Mx)

			fmt.Fprintf(fo, "%1d %.4f %1d %.4f %1d %.4f ", desi.xidy.Hrs, desi.xidy.M, desi.xidy.Mntime, desi.xidy.Mn, desi.xidy.Mxtime, desi.xidy.Mx)
			fmt.Fprintf(fo, "%1d %.4f %1d %.4f %1d %.4f ", desi.xody.Hrs, desi.xody.M, desi.xody.Mntime, desi.xody.Mn, desi.xody.Mxtime, desi.xody.Mx)
			fmt.Fprintf(fo, "%1d %.4f %1d %.4f %1d %.4f ", desi.xsdy.Hrs, desi.xsdy.M, desi.xsdy.Mntime, desi.xsdy.Mn, desi.xsdy.Mxtime, desi.xsdy.Mx)

			fmt.Fprintf(fo, "%1d %3.1f ", desi.Qsdy.Hhr, desi.Qsdy.H)
			fmt.Fprintf(fo, "%1d %3.1f ", desi.Qsdy.Chr, desi.Qsdy.C)
			fmt.Fprintf(fo, "%1d %2.0f ", desi.Qsdy.Hmxtime, desi.Qsdy.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f ", desi.Qsdy.Cmxtime, desi.Qsdy.Cmx)

			fmt.Fprintf(fo, "%1d %3.1f ", desi.Qldy.Hhr, desi.Qldy.H)
			fmt.Fprintf(fo, "%1d %3.1f ", desi.Qldy.Chr, desi.Qldy.C)
			fmt.Fprintf(fo, "%1d %2.0f ", desi.Qldy.Hmxtime, desi.Qldy.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f ", desi.Qldy.Cmxtime, desi.Qldy.Cmx)

			fmt.Fprintf(fo, "%1d %3.1f ", desi.Qtdy.Hhr, desi.Qtdy.H)
			fmt.Fprintf(fo, "%1d %3.1f ", desi.Qtdy.Chr, desi.Qtdy.C)
			fmt.Fprintf(fo, "%1d %2.0f ", desi.Qtdy.Hmxtime, desi.Qtdy.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f ", desi.Qtdy.Cmxtime, desi.Qtdy.Cmx)

			fmt.Fprintf(fo, "%1d %3.1f ", desi.Qlsdy.Hhr, desi.Qlsdy.H)
			fmt.Fprintf(fo, "%1d %3.1f ", desi.Qlsdy.Chr, desi.Qlsdy.C)
			fmt.Fprintf(fo, "%1d %2.0f ", desi.Qlsdy.Hmxtime, desi.Qlsdy.Hmx)
			fmt.Fprintf(fo, "%1d %2.0f\n", desi.Qlsdy.Cmxtime, desi.Qlsdy.Cmx)
		}
	}
}

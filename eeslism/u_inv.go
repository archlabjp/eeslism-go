package eeslism

import (
	"fmt"
	"io"
	"math"
	"os"
)

///* ------------------------------------------------
//
// 逆行列の計算
//
//*/

/* ========= ガウスジョルダン法の関数====================== */

func Matinv(a []float64, n, m int, s string) {
	row := make([]int, m)
	mattemp := make([]float64, m*m)

	matcpy(a, mattemp, m*m)

	for ipv := 0; ipv < m; ipv++ {
		/* ---- 最大値探索 ---------------------------- */
		big := 0.0
		pivot_row := ipv
		for i := ipv; i < m; i++ {
			if math.Abs(a[i*n+ipv]) > big {
				big = math.Abs(a[i*n+ipv])
				pivot_row = i
			}
		}
		if big == 0.0 {
			var E string
			if s != "" {
				E = fmt.Sprintf("対角要素に０があります  matrix=%dx%d  i=%d  [%s]", m, m, ipv, s)
			} else {
				E = fmt.Sprintf("対角要素に０があります  matrix=%dx%d  i=%d", m, m, ipv)
			}
			Matprint("%.2g  ", m, mattemp)
			Eprint("<matinv>", E)
			Preexit()
			os.Exit(EXIT_MATINV)
		}
		row[ipv] = pivot_row

		/* ---- 行の入れ替え -------------------------- */
		if ipv != pivot_row {
			for i := 0; i < m; i++ {
				temp := a[ipv*n+i]
				a[ipv*n+i] = a[pivot_row*n+i]
				a[pivot_row*n+i] = temp
			}
		}

		/* ---- 対角成分=1(ピボット行の処理) ---------- */
		inv_pivot := 1.0 / a[ipv*n+ipv]
		a[ipv*n+ipv] = 1.0
		for j := 0; j < m; j++ {
			a[ipv*n+j] *= inv_pivot
		}

		/* ---- ピボット列=0(ピボット行以外の処理) ---- */
		for i := 0; i < m; i++ {
			if i != ipv {
				temp := a[i*n+ipv]
				a[i*n+ipv] = 0.0
				for j := 0; j < m; j++ {
					a[i*n+j] -= temp * a[ipv*n+j]
				}
			}
		}
	}

	/* ---- 列の入れ替え(逆行列) -------------------------- */
	for j := m - 1; j >= 0; j-- {
		if j != row[j] {
			for i := 0; i < m; i++ {
				temp := a[i*n+j]
				a[i*n+j] = a[i*n+row[j]]
				a[i*n+row[j]] = temp
			}
		}
	}
}

///* -----------------------------------------------------*/
// 連立1次方程式の解法
// ガウス・ザイデル法

// [A]{B}={C}
// [A]:係数行列
// {B}:解
// {C}:定数行列

//	 m :未知数の数
//	 n :配列の定義数

// 参考文献：C言語による科学技術計算サブルーチンライブラリ
// pp.104-106
// ----------------------------------------------------- */
func Gausei(A, C []float64, m, n int, B []float64) {
	eps := 1.0e-6
	l := m + 1
	a := make([]float64, m*l)

	for i := 0; i < m*l; i++ {
		a[i] = 0.0
	}

	for i := 0; i < m; i++ {
		for j := 0; j < m+1; j++ {
			if j < m {
				a[i*l+j] = A[i*n+j] / A[i*n+i]
			} else {
				a[i*l+j] = C[i] / A[i*n+i]
			}
		}
	}

	for i := 0; i < m; i++ {
		B[i] = 0.2
	}

	def := FNAN
	k := 0

	for def > eps {
		def = 0.0

		for i := 0; i < m; i++ {
			sum := 0.0
			s := i * l
			w := a[s+i]

			for j := 0; j < m; j++ {
				if i != j {
					sum += a[s+j] * B[j]
				}
			}

			y := (a[s+m] - sum) / w
			ay := math.Abs(y - B[i])

			def = math.Max(def, ay)
			B[i] = y
		}

		if def <= eps {
			break
		}

		if k > 100 {
			for i := 0; i < m; i++ {
				fmt.Printf("i=%d  %f\n", i, B[i])
			}

			fmt.Println("収束せず")
			Preexit()
			os.Exit(EXIT_MATINV)
		}

		k++
	}
}

//	  /* -----------------------------------------------------
//	  連立1次方程式の解法
//	  ガウスの消去法
//
//	   [A]{B}={C}
//	   [A]:係数行列
//	   {B}:解
//	   {C}:定数行列
//
//		m :未知数の数
//		n :配列の定義数
//
//		 参考文献：C言語による科学技術計算サブルーチンライブラリ
//		 pp.104-106
//	  ----------------------------------------------------- */
func Gauss(A, C, B []float64, m, n int) {
	num := make([]int, m)
	pivot := make([]int, m)
	wfs := make([]float64, m*m)

	for i := 0; i < m; i++ {
		B[i] = C[i]
		for j := 0; j < m; j++ {
			wfs[i*m+j] = A[i*n+j]
		}
	}

	for i := 0; i < m; i++ {
		num[i] = i + 1
	}

	for k := 0; k < m; k++ {
		pv := 0.0

		for i := 0; i < m; i++ {
			if num[i] != 0 {
				if math.Abs(wfs[i*m+k]) > math.Abs(pv) {
					pv = wfs[i*m+k]
					pivot[k] = i
				}
			}
		}

		if pv == 0.0 {
			Eprint("<gauss>", "対角要素に 0 があります")
			Preexit()
			os.Exit(EXIT_MATINV)
		}

		for j := k; j < m; j++ {
			wfs[pivot[k]*m+j] /= pv
		}

		B[pivot[k]] /= pv
		num[pivot[k]] = 0

		for i := 0; i < m; i++ {
			if num[i] != 0 {
				tmp := wfs[i*m+k]

				for j := k + 1; j < m; j++ {
					wfs[i*m+j] -= wfs[pivot[k]*m+j] * tmp
				}

				B[i] -= B[pivot[k]] * tmp
			}
		}
	}

	for i := m - 2; i >= 0; i-- {
		for j := i + 1; j < m; j++ {
			B[pivot[i]] -= B[pivot[i]*m+j]
		}
	}
}

// /* -----------------------------------------------------
//
// 正方行列の出力
// */
func Matprint(format string, N int, a []float64) {
	fmt.Println()
	for i := 0; i < N; i++ {
		for j := 0; j < N; j++ {
			fmt.Printf(format, a[i*N+j])
		}
		fmt.Println()
	}
}

// /* -----------------------------------------------------
//
// 正方行列のファイル出力
// */
func Matfiprint(f io.Writer, fmtStr string, N int, a []float64) {
	for i := 0; i < N; i++ {
		for j := 0; j < N; j++ {
			fmt.Fprintf(f, fmtStr, a[i*N+j])
		}
		fmt.Fprintln(f)
	}
}

// /* -----------------------------------------------------
//
// 正方行列の出力 （単精度）
// */
func Matfprint(fmtStr string, N int, a []float64) {
	fmt.Printf("\n")
	for i := 0; i < N; i++ {
		for j := 0; j < N; j++ {
			fmt.Printf(fmtStr, a[i*N+j])
		}
		fmt.Printf("\n")
	}
}

// /* -----------------------------------------------------
//
// 連立一次方程式の係数行列及び右辺の出力
//
// */
func Seqprint(fmt1 string, N int, a []float64, fmt2 string, c []float64) {
	fmt.Println("--- seqprint")
	for i := 0; i < N; i++ {
		for j := 0; j < N; j++ {
			fmt.Printf(fmt1, a[i*N+j])
		}
		fmt.Print("  c=")
		fmt.Printf(fmt2, c[i])
		fmt.Println()
	}
}

// /* ---------------------------------------------- */
//
// /*
// 行列の掛け算
//
// (T)=[A](V)
// N:宣言寸法　　n:使用寸法
// */
func Matmalv(A []float64, V []float64, N int, n int, T []float64) {
	for i := 0; i < n; i++ {
		var sum float64 = 0.0
		a := A[i*N : (i+1)*N]
		for j := 0; j < n; j++ {
			sum += a[j] * V[j]
		}
		T[i] = sum
	}
}

/****************************************************************/
//		行列の０初期化
/****************************************************************/
func matinit(A []float64, N int) {
	for i := 0; i < N; i++ {
		A[i] = 0.0
	}
}

func imatinit(A []int, N int) {
	for i := 0; i < N; i++ {
		A[i] = 0
	}
}

/****************************************************************/
//		行列の数値初期化
/****************************************************************/
func matinitx(A []float64, N int, x float64) {
	for i := 0; i < N; i++ {
		A[i] = x
	}
}

/****************************************************************/
//		行列のコピー
/****************************************************************/
func matcpy(A []float64, B []float64, N int) {
	for i := 0; i < N; i++ {
		B[i] = A[i]
	}
}

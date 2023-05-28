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

/* mcevcooling.c */

package main

/*  気化冷却器  */

/*  仕様入力  */

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func Evacdata(s string, Evacca *EVACCA) int {
	st := strings.IndexByte(s, '=')

	if st == -1 {
		Evacca.Name = s
		Evacca.N = -999
		Evacca.Nlayer = -999
		Evacca.Awet = -999.0
		Evacca.Adry = -999.0
		Evacca.hdry = -999.0
		Evacca.hwet = -999.0
	} else {
		key := s[:st]
		value := s[st+1:]

		switch key {
		case "Awet":
			if val, err := strconv.ParseFloat(value, 64); err == nil {
				Evacca.Awet = val
			} else {
				panic(err)
			}
		case "Adry":
			if val, err := strconv.ParseFloat(value, 64); err == nil {
				Evacca.Adry = val
			} else {
				panic(err)
			}
		case "hwet":
			if val, err := strconv.ParseFloat(value, 64); err == nil {
				Evacca.hwet = val
			} else {
				panic(err)
			}
		case "hdry":
			if val, err := strconv.ParseFloat(value, 64); err == nil {
				Evacca.hdry = val
			} else {
				panic(err)
			}
		case "N":
			if val, err := strconv.Atoi(value); err == nil {
				Evacca.N = val
			} else {
				panic(err)
			}
		case "Nlayer":
			if val, err := strconv.Atoi(value); err == nil {
				Evacca.Nlayer = val
			} else {
				panic(err)
			}
		default:
			return 1
		}
	}

	return 0
}

/* ------------------------------------------------------ */
// 初期設定（入力漏れのチェック、変数用メモリの確保）
func Evacint(Nevac int, Evac []EVAC) {
	for i := 0; i < Nevac; i++ {
		evac := &Evac[i]
		cat := evac.Cat

		// 入力漏れのチェック
		if cat.N < 0 {
			msg := fmt.Sprintf("Name=%s catname=%s 分割数が未定義です", evac.Name, cat.Name)
			Eprint("<Evacint>", msg)
		}
		if cat.Adry < 0.0 || cat.Awet < 0.0 || (cat.Nlayer < 0 && (cat.hdry < 0.0 || cat.hwet < 0.0)) {
			msg := fmt.Sprintf("Name=%s catname=%s Adry=%.1g Awet=%.1g hdry=%.1g hwet=%.1g\n",
				evac.Name, cat.Name, cat.Adry, cat.Awet, cat.hdry, cat.hwet)
			Eprint("<Evacint>", msg)
		}

		// 面積を分割後の面積に変更
		cat.Adry /= float64(cat.N)
		cat.Awet /= float64(cat.N)

		// 必要なメモリ領域の確保
		if cat.N > 0 {
			Temp := FNXtr(20.0, 50.0)
			evac.M = make([]float64, cat.N)    // 蒸発量のメモリ確保
			evac.Kx = make([]float64, cat.N)   // 物質移動係数のメモリ確保
			evac.Tdry = make([]float64, cat.N) // Dry側温度のメモリ確保
			for i := range evac.Tdry {
				evac.Tdry[i] = 20.0
			}
			evac.Twet = make([]float64, cat.N) // Wet側温度のメモリ確保
			for i := range evac.Twet {
				evac.Twet[i] = 20.0
			}
			evac.Xdry = make([]float64, cat.N) // Dry側絶対湿度のメモリ確保
			for i := range evac.Xdry {
				evac.Xdry[i] = Temp
			}
			evac.Xwet = make([]float64, cat.N) // Wet側絶対湿度のメモリ確保
			for i := range evac.Xwet {
				evac.Xwet[i] = Temp
			}
			evac.Ts = make([]float64, cat.N) // 境界層温度のメモリ確保
			for i := range evac.Ts {
				evac.Ts[i] = 20.0
			}
			evac.Xs = make([]float64, cat.N) // 境界層絶対湿度のメモリ確保
			xs := FNXtr(20.0, 100.0)
			for i := range evac.Xs {
				evac.Xs[i] = xs
			}
			evac.RHdry = make([]float64, cat.N) // Dry側相対湿度のメモリ確保
			for i := range evac.RHdry {
				evac.RHdry[i] = 50.0
			}
			evac.RHwet = make([]float64, cat.N) // Wet側相対湿度のメモリ確保
			for i := range evac.RHwet {
				evac.RHwet[i] = 50.0
			}

			// 状態値計算用行列
			N := cat.N * 5
			N2 := N * N
			evac.UX = make([]float64, N2)
			evac.UXC = make([]float64, N)
		}
	}
}

// 飽和絶対湿度の線形近似（Ts℃付近で線形回帰式を作成）
func LinearSatx(Ts float64) (a, b float64) {
	// 線形近似の区間の設定（Tsを中心にEPS区間）
	T1 := Ts - 0.2
	T2 := Ts + 0.2

	// T1、T2における飽和絶対湿度の計算
	x1 := FNXs(T1)
	x2 := FNXs(T2)

	// 線形回帰式の傾きの計算
	a = (x2 - x1) / (T2 - T1)

	// 線形回帰式の切片の計算
	b = x1 - a*T1

	return a, b
}

// 湿り空気の飽和絶対湿度の計算
func FNXs(T float64) float64 {
	return 4.2849e-3 * math.Exp(6.0260e-2*T)
}

/*  気化冷却器出口空気温湿度に関する変数割当  */
func Evacelm(Nevac int, Evac []EVAC) {
	for i := 0; i < Nevac; i++ {
		EoTdry := Evac[i].Cmp.Elouts[0] // Tdryoutの出口温度計算用
		Eoxdry := Evac[i].Cmp.Elouts[1] // xdryoutの出口温度計算用
		EoTwet := Evac[i].Cmp.Elouts[2] // Twetoutの出口温度計算用
		Eoxwet := Evac[i].Cmp.Elouts[3] // xwetoutの出口温度計算用

		EoTdry.Elins[0].Upo = Eoxdry.Elins[0].Upo
		EoTdry.Elins[0].Upv = Eoxdry.Elins[0].Upo
		EoTdry.Elins[1].Upo = EoTwet.Elins[0].Upo
		EoTdry.Elins[1].Upv = EoTwet.Elins[0].Upo
		EoTdry.Elins[2].Upo = Eoxwet.Elins[0].Upo
		EoTdry.Elins[2].Upv = Eoxwet.Elins[0].Upo

		Eoxdry.Elins[0].Upo = EoTdry.Elins[0].Upo
		Eoxdry.Elins[0].Upv = EoTdry.Elins[0].Upo
		Eoxdry.Elins[1].Upo = EoTwet.Elins[0].Upo
		Eoxdry.Elins[1].Upv = EoTwet.Elins[0].Upo
		Eoxdry.Elins[2].Upo = Eoxwet.Elins[0].Upo
		Eoxdry.Elins[2].Upv = Eoxwet.Elins[0].Upo

		EoTwet.Elins[0].Upo = EoTdry.Elins[0].Upo
		EoTwet.Elins[0].Upv = EoTdry.Elins[0].Upo
		EoTwet.Elins[1].Upo = Eoxdry.Elins[0].Upo
		EoTwet.Elins[1].Upv = Eoxdry.Elins[0].Upo
		EoTwet.Elins[2].Upo = Eoxwet.Elins[0].Upo
		EoTwet.Elins[2].Upv = Eoxwet.Elins[0].Upo

		Eoxwet.Elins[0].Upo = EoTdry.Elins[0].Upo
		Eoxwet.Elins[0].Upv = EoTdry.Elins[0].Upo
		Eoxwet.Elins[1].Upo = Eoxdry.Elins[0].Upo
		Eoxwet.Elins[1].Upv = Eoxdry.Elins[0].Upo
		Eoxwet.Elins[2].Upo = EoTwet.Elins[0].Upo
		Eoxwet.Elins[2].Upv = EoTwet.Elins[0].Upo

	}
}

// 風速の計算
func Evacu(G, T, H, W float64, N int) float64 {
	u := G / FNarow(T) / (float64(N) * H * W)
	return u
}

// 通気部の対流熱伝達率の計算（プログラムを解読したため詳細は不明）
func Evachcc(de, L, T, H, W float64, N, G, Flg int) float64 {
	// 流路縦横比の計算
	AR := H / (W / 5.0)

	// 通気部の風速の計算
	u := Evacu(float64(G), T, H, W, N)

	// レイノルズ数の計算
	Re := u * L / FNanew(T)

	// ヌセルト数の計算
	Nu := EvacNu(AR, Re)

	// 裕度の計算
	Mgn := 0.875
	if Flg == 'd' {
		if Re > 1000.0 {
			Mgn *= (0.0000128205*Re + 1.0859)
		} else {
			Mgn *= (0.00083333*Re + 0.18333)
		}
	}

	// 対流熱伝達率の計算
	hc := Nu / de * FNalam(T) * Mgn

	return hc
}

// 通気部のヌセルト数を計算する
func EvacNu(AR, Re float64) float64 {
	var Nu float64

	if Re <= 1000.0 {
		Nu = -13.042*AR*AR*AR + 27.063*AR*AR - 18.591*AR + 7.54
	} else if Re <= 2000.0 {
		Nu = (-0.023131*AR*AR+0.018229*AR+0.00043299)*Re +
			(46.261*AR*AR - 36.459*AR + 7.0971)
	} else {
		Nu = 0.021 * math.Pow(Re, 0.8) * math.Pow(0.7, 0.4)
	}

	return Nu
}

// 要素方程式の係数計算
func Evaccfv(Nevac int, Evac []EVAC) {
	for i := 0; i < Nevac; i++ {
		EvpFlg := make([]float64, Evac[i].Cat.N)
		if Evac[i].Cmp.Control != OFF_SW {
			EoTdry := Evac[i].Cmp.Elouts[0] // Tdryoutの出口温度計算用
			Eoxdry := Evac[i].Cmp.Elouts[1] // xdryoutの出口温度計算用
			EoTwet := Evac[i].Cmp.Elouts[2] // Twetoutの出口温度計算用
			Eoxwet := Evac[i].Cmp.Elouts[3] // xwetoutの出口温度計算用

			cat := Evac[i].Cat
			Gdry := EoTdry.G
			Gwet := EoTwet.G

			if cat.Nlayer > 0 {
				Ts := Evac[i].Ts
				Tsave := 0.0
				for ii := 0; ii < cat.N; ii++ {
					Tsave += Ts[ii] / float64(cat.N)
				}
				cat.hdry = Evachcc(4.3/1000.0, 4.3/1000.0, Tsave, 2.3/1000.0, 173.0/1000.0, cat.Nlayer, int(Gdry), 'd') // Dry側の対流熱伝達率の計算
				cat.hwet = Evachcc(6.4/1000.0, 6.4/1000.0, Tsave, 3.5/1000.0, 173.0/1000.0, cat.Nlayer, int(Gwet), 'w') // Wet側の対流熱伝達率の計算
			}

			N := cat.N * 5
			N2 := N * N

			U := make([]float64, N2) // 行列Uの作成
			C := make([]float64, N)  // 行列Cの作成
			a := make([]float64, cat.N)
			b := make([]float64, cat.N)

			PreFlg := 1.0
			for ii := cat.N - 1; ii >= 0; ii-- {
				Ts := &Evac[i].Ts[ii]
				xwet := &Evac[i].Xwet[ii]
				RHwet := &Evac[i].RHwet[ii]
				//kx := &Evac[i].Kx[ii]
				EF := &EvpFlg[ii]

				a[ii], b[ii] = LinearSatx(*Ts) // 境界層温度における飽和絶対湿度計算用係数の取得
				*EF = 1.0
				if a[ii]**Ts+b[ii]-*xwet < 0.0 || *RHwet > 99.0 || math.Abs(PreFlg) <= 1e-5 {
					*EF = 0.0
				}
				PreFlg = *EF
			}

			for ii := 0; ii < cat.N; ii++ {
				Ts := &Evac[i].Ts[ii]
				xs := &Evac[i].Xs[ii]
				// xwet := &Evac[i].Xwet[ii]
				// RHwet := &Evac[i].RHwet[ii]
				kx := &Evac[i].Kx[ii]
				EF := &EvpFlg[ii]

				*kx = cat.hwet / (Ca + Cv**xs) * *EF // 物質移動係数の計算
				A := *kx * (Ro + Cv**Ts) * *EF       // 係数の計算

				// C行列の作成
				C[ii*5+0] = 0.0                    // Twetの状態方程式には定数項はない
				C[ii*5+1] = -cat.Awet * *kx * b[0] // xwetの定数項作成
				C[ii*5+2] = A * b[0]               // Tsの定数項作成
				C[ii*5+3] = 0.0                    // Tdryの係数はゼロ
				C[ii*5+4] = 0.0                    // xdryの係数はゼロ

				// U行列の作成

				// 対角行列
				U[N*(5*ii+0)+(5*ii+0)] = -(Ca*Gwet + cat.Awet*cat.hwet)   // Twetの項
				U[N*(5*ii+1)+(5*ii+1)] = -(Gwet + cat.Awet**kx)           // xwetの項
				U[N*(5*ii+2)+(5*ii+2)] = -(cat.hwet + A*a[ii] + cat.hdry) // Tsの項
				U[N*(5*ii+2)+(5*ii+3)] = -(Ca*Gdry + cat.Adry*cat.hdry)   // Tdryの項
				U[N*(5*ii+2)+(5*ii+4)] = 1.0                              // xdryの項

				U[N*(5*ii+0)+(5*ii+2)] = cat.Awet * cat.hwet   // TwetとTsの項
				U[N*(5*ii+1)+(5*ii+2)] = cat.Awet * *kx * a[0] // xwetとTsの項
				U[N*(5*ii+2)+(5*ii+0)] = cat.hwet              // TsとTwetの項
				U[N*(5*ii+2)+(5*ii+1)] = A                     // Tsとxwetの項
				U[N*(5*ii+2)+(5*ii+3)] = cat.hdry              // TsとTdryの項
				U[N*(5*ii+3)+(5*ii+2)] = cat.Adry * cat.hdry   // TdryとTsの項

				//  Dry側上流
				if ii > 0 {
					U[N*(5*ii+3)+(5*(ii-1)+3)] = Ca * Gdry // Tdryの項
					U[N*(5*ii+4)+(5*(ii-1)+4)] = -1.0      // xdryの項
				}

				// Wet側の上流
				if ii < cat.N-1 {
					U[N*(5*ii+0)+(5*(ii+1)+0)] = Ca * Gwet // Twetの項
					U[N*(5*ii+1)+(5*(ii+1)+1)] = Gwet      // xwetの項
				}
			}

			Matinv(U, N, N, "Evaccfv U") // 行列Uの逆行列を計算
			matinit(Evac[i].UX, N2)      // 行列の初期化
			matcpy(U, Evac[i].UX, N2)    // 行列のコピー

			matinit(Evac[i].UXC, N) // 行列UXCの初期化
			for ii := 0; ii < N; ii++ {
				for jj := 0; jj < N; jj++ {
					Evac[i].UXC[ii] += Evac[i].UX[ii*N+jj] * C[jj] // 行列UXとベクトルCの積の計算
				}
			}

			Row := N * (5*(cat.N-1) + 3)

			EoTdry.Coeffo = -1.0                                               // Tdryoutの要素方程式
			EoTdry.Co = -Evac[i].UXC[5*(cat.N-1)+3]                            // 定数の項
			EoTdry.Coeffin[0] = -Evac[i].UX[Row+(5*(cat.N-1)+3)] * (Ca * Gdry) // Tdryinの項
			EoTdry.Coeffin[1] = -Evac[i].UX[Row+(5*(1-1)+4)] * -1.0            // xdryinの項
			EoTdry.Coeffin[2] = -Evac[i].UX[Row+(5*(cat.N-1)+0)] * (Ca * Gwet) // Twetinの項
			EoTdry.Coeffin[3] = -Evac[i].UX[Row+(5*(cat.N-1)+1)] * (Gwet)      // xwetinの項

			Eoxdry.Coeffo = -1.0                                               // xdryoutの要素方程式
			Eoxdry.Co = -Evac[i].UXC[5*(cat.N-1)+4]                            // 定数の項
			Eoxdry.Coeffin[0] = -Evac[i].UX[Row+(5*(1-1)+4)] * -1.0            // xdryinの項
			Eoxdry.Coeffin[1] = -Evac[i].UX[Row+(5*(1-1)+3)] * (Ca * Gdry)     // Tdryinの項
			Eoxdry.Coeffin[2] = -Evac[i].UX[Row+(5*(cat.N-1)+0)] * (Ca * Gwet) // Twetinの項
			Eoxdry.Coeffin[3] = -Evac[i].UX[Row+(5*(cat.N-1)+1)] * (Gwet)      // xwetinの項

			Row = N * (5*(1-1) + 0)
			EoTwet.Coeffo = -1.0                                               // Twetoutの要素方程式
			EoTwet.Co = -Evac[i].UXC[5*(1-1)+0]                                // 定数の項
			EoTwet.Coeffin[0] = -Evac[i].UX[Row+(5*(cat.N-1)+0)] * (Ca * Gwet) // Twetinの項
			EoTwet.Coeffin[1] = -Evac[i].UX[Row+(5*(1-1)+3)] * (Ca * Gdry)     // Tdryinの項
			EoTwet.Coeffin[2] = -Evac[i].UX[Row+(5*(1-1)+4)] * -1.0            // xdryinの項
			EoTwet.Coeffin[3] = -Evac[i].UX[Row+(5*(cat.N-1)+1)] * (Gwet)      // xwetinの項

			Row = N * (5*(1-1) + 1)
			Eoxwet.Coeffo = -1.0                                               // xwetoutの要素方程式
			Eoxwet.Co = -Evac[i].UXC[5*(1-1)+1]                                // 定数の項
			Eoxwet.Coeffin[0] = -Evac[i].UX[Row+(5*(cat.N-1)+1)] * (Gwet)      // xwetinの項
			Eoxwet.Coeffin[1] = -Evac[i].UX[Row+(5*(1-1)+3)] * (Ca * Gdry)     // Tdryinの項
			Eoxwet.Coeffin[2] = -Evac[i].UX[Row+(5*(1-1)+4)] * -1.0            // xdryinの項
			Eoxwet.Coeffin[3] = -Evac[i].UX[Row+(5*(cat.N-1)+0)] * (Ca * Gwet) // Twetinの項
		}
	}
}

// 内部温度、熱量の計算
func Evacene(Nevac int, Evac []EVAC, Evacreset *int) {
	for i := 0; i < Nevac; i++ {
		Evac := &Evac[i]
		cat := Evac.Cat
		if Evac.Cmp.Control != OFF_SW {
			var Gdry, Gwet float64
			var Tdry, Twet, xdry, xwet, Ts, xs, RHwet, RHdry, M, kx []float64
			var Sin, Stat, Scmp []float64
			//var elin *ELIN

			EoTdry := Evac.Cmp.Elouts[0] //Tdryoutの出口温度計算用
			Eoxdry := Evac.Cmp.Elouts[1] //xdryoutの出口温度計算用
			EoTwet := Evac.Cmp.Elouts[2] //Twetoutの出口温度計算用
			Eoxwet := Evac.Cmp.Elouts[3] //xwetoutの出口温度計算用

			// 出入口状態値の取得
			Evac.Tdryi = EoTdry.Elins[0].Sysvin
			Evac.Xdryi = EoTdry.Elins[1].Sysvin
			Evac.Tweti = EoTdry.Elins[2].Sysvin
			Evac.Xweti = EoTdry.Elins[3].Sysvin

			Evac.Tdryo = EoTdry.Sysv
			Evac.Xdryo = Eoxdry.Sysv
			Evac.Tweto = EoTwet.Sysv
			Evac.Xweto = Eoxwet.Sysv

			Gdry = EoTdry.G
			Gwet = EoTwet.G

			// 相対湿度の計算
			Evac.RHdryi = FNRhtx(Evac.Tdryi, Evac.Xdryi)
			Evac.RHdryo = FNRhtx(Evac.Tdryo, Evac.Xdryo)
			Evac.RHweti = FNRhtx(Evac.Tweti, Evac.Xweti)
			Evac.RHweto = FNRhtx(Evac.Tweto, Evac.Xweto)

			// 熱量の計算
			Evac.Qsdry = Ca * Gdry * (Evac.Tdryo - Evac.Tdryi)
			Evac.Qldry = Ro * Gdry * (Evac.Xdryo - Evac.Xdryi)
			Evac.Qtdry = Evac.Qsdry + Evac.Qldry
			Evac.Qswet = Ca * Gwet * (Evac.Tweto - Evac.Tweti)
			Evac.Qlwet = Ro * Gwet * (Evac.Xweto - Evac.Xweti)
			Evac.Qtwet = Evac.Qswet + Evac.Qlwet

			N := cat.N * 5
			//N2 := N * N

			Sin = make([]float64, N)
			Stat = make([]float64, N)

			Sin[5*(1-1)+3] = Ca * Gdry * Evac.Tdryi
			Sin[5*(1-1)+4] = -Evac.Xdryi
			Sin[5*(cat.N-1)+0] = Ca * Gwet * Evac.Tweti
			Sin[5*(cat.N-1)+1] = Gwet * Evac.Xweti

			for ii := 0; ii < N; ii++ {
				for jj := 0; jj < N; jj++ {
					// 内部変数の計算
					Stat[ii] += -Evac.UX[ii*N+jj] * Sin[jj]
				}
				Stat[ii] += Evac.UXC[ii]
			}

			// 内部変数計算結果の格納
			Tdry = Evac.Tdry
			xdry = Evac.Xdry
			Twet = Evac.Twet
			xwet = Evac.Xwet
			Ts = Evac.Ts
			xs = Evac.Xs
			RHdry = Evac.RHdry
			RHwet = Evac.RHwet
			M = Evac.M
			Scmp = Stat
			kx = Evac.Kx
			for ii := 0; ii < cat.N; ii++ {
				Twet[ii] = Scmp[0]
				xwet[ii] = Scmp[1]
				Ts[ii] = Scmp[2]
				Tdry[ii] = Scmp[3]
				xdry[ii] = Scmp[4]
				xs[ii] = FNXtr(Ts[ii], 100.0)

				// 相対湿度の計算
				RHdry[ii] = FNRhtx(Tdry[ii], xdry[ii])
				RHwet[ii] = FNRhtx(Twet[ii], xwet[ii])

				// 蒸発量の計算
				M[ii] = kx[ii] * math.Max(xs[ii]-xwet[ii], 0.0) * cat.Awet

				Scmp = Scmp[5:]
			}
		} else {
			Evac.Qsdry = 0.0
			Evac.Qldry = 0.0
			Evac.Qtdry = 0.0
			Evac.Qswet = 0.0
			Evac.Qlwet = 0.0
			Evac.Qtwet = 0.0
			Evac.Tdryi = 0.0
			Evac.Tdryo = 0.0
			Evac.Tweti = 0.0
			Evac.Tweto = 0.0
			Evac.Xdryi = 0.0
			Evac.Xdryo = 0.0
			Evac.Xweti = 0.0
			Evac.Xweto = 0.0
			matinit(Evac.M, cat.N)
		}

		Evac.Count++
		if Evac.Count > 0 {
			*Evacreset = 1
		}
	}
}

// カウンタのリセット
func Evaccountreset(Nevac int, Evac []EVAC) {
	for i := 0; i < Nevac; i++ {
		Evac[i].Count = 0
	}
}

// 代表日の計算結果出力
func Evacprint(fo *os.File, id int, Nevac int, Evac []EVAC) {
	switch id {
	case 0:
		if Nevac > 0 {
			fmt.Fprintf(fo, "%s %d\n", EVAC_TYPE, Nevac)
		}
		for i := 0; i < Nevac; i++ {
			fmt.Fprintf(fo, " %s 1 %d\n", Evac[i].Name, 18+8*Evac[i].Cat.N)
		}

	case 1:
		for i := 0; i < Nevac; i++ {
			Evac := &Evac[i]
			// Wet側出力
			fmt.Fprintf(fo, "%s_cw c c %s_Gw m f %s_Twi t f %s_Two t f %s_xwi x f %s_xwo x f\n",
				Evac.Name, Evac.Name, Evac.Name, Evac.Name, Evac.Name, Evac.Name)
			fmt.Fprintf(fo, "%s_Qws q f %s_Qwl q f %s_Qwt q f\n",
				Evac.Name, Evac.Name, Evac.Name)
			// Dry側出力
			fmt.Fprintf(fo, "%s_cd c c %s_Gd m f %s_Tdi t f %s_Tdo t f %s_xdi x f %s_xdo x f\n",
				Evac.Name, Evac.Name, Evac.Name, Evac.Name, Evac.Name, Evac.Name)
			fmt.Fprintf(fo, "%s_Qds q f %s_Qdl q f %s_Qdt q f\n",
				Evac.Name, Evac.Name, Evac.Name)
			// 内部変数
			for ii := 0; ii < Evac.Cat.N; ii++ {
				fmt.Fprintf(fo, "%s_Tw[%d] t f %s_xw[%d] x f %s_RHw[%d] r f %s_Ts[%d] t f %s_Td[%d] t f %s_xd[%d] x f %s_RHd[%d] r f %s_M[%d] m f\n",
					Evac.Name, ii, Evac.Name, ii, Evac.Name, ii, Evac.Name, ii, Evac.Name, ii, Evac.Name, ii, Evac.Name, ii, Evac.Name, ii)
			}
		}

	default:
		for i := 0; i < Nevac; i++ {
			Evac := &Evac[i]
			// Wet側出力
			elo := Evac.Cmp.Elouts[2]
			fmt.Fprintf(fo, "%c %g %.1f %.1f %.3f %.3f %.1f %.1f %.1f\n",
				elo.Control, elo.G, Evac.Tweti, Evac.Tweto, Evac.Xweti, Evac.Xweto, Evac.Qswet, Evac.Qlwet, Evac.Qtwet)
			// Dry側出力
			elo = Evac.Cmp.Elouts[0]
			fmt.Fprintf(fo, "%c %g %.1f %.1f %.3f %.3f %.1f %.1f %.1f\n",
				elo.Control, elo.G, Evac.Tdryi, Evac.Tdryo, Evac.Xdryi, Evac.Xdryo, Evac.Qsdry, Evac.Qldry, Evac.Qtdry)
			// 内部変数
			Tdry := Evac.Tdry
			xdry := Evac.Xdry
			Twet := Evac.Twet
			xwet := Evac.Xwet
			Ts := Evac.Ts
			RHdry := Evac.RHdry
			RHwet := Evac.RHwet
			M := Evac.M
			for ii := 0; ii < Evac.Cat.N; ii++ {
				fmt.Fprintf(fo, "%.1f %.3f %.0f %.1f %.1f %.3f %.0f %.3e\n",
					Twet[ii], xwet[ii], RHwet[ii], Ts[ii], Tdry[ii], xdry[ii], RHdry[ii], M[ii])
			}
		}
	}
}

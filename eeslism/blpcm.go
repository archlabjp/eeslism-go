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

/*   binit.c   */
package eeslism

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

/*  壁体デ－タの入力  */

func PCMdata(fi *EeTokens, dsn string, pcm *[]*PCM, pcmiterate *rune) {
	N := PCMcount(fi)

	s := "PCMdata --"

	if N > 0 {
		*pcm = make([]*PCM, N)

		for j := 0; j < N; j++ {
			var PCMa = new(PCM)

			PCMa.Name = ""
			PCMa.Condl = FNAN
			PCMa.Conds = FNAN
			PCMa.Crol = FNAN
			PCMa.Cros = FNAN
			PCMa.Ql = FNAN
			PCMa.Tl = FNAN
			PCMa.Ts = FNAN
			PCMa.Tp = FNAN

			PCMa.Iterate = false // 収束なしがデフォルト
			//PCMa.iterateTemp = false	// デフォルトは見かけの比熱だけで収束
			PCMa.IterateTemp = true
			PCMa.NWeight = 0.5 // 収束計算時の現在ステップ温度の重み係数
			PCMa.AveTemp = 'y'
			PCMa.DivTemp = 1
			PCMa.Ctype = 2 // 二等辺三角形がデフォルト
			// パラメータの初期化
			PCMa.PCMp.a = FNAN
			PCMa.PCMp.B = FNAN
			PCMa.PCMp.b = FNAN
			PCMa.PCMp.bl = FNAN
			PCMa.PCMp.bs = FNAN
			PCMa.PCMp.c = FNAN
			PCMa.PCMp.d = FNAN
			PCMa.PCMp.e = FNAN
			PCMa.PCMp.f = FNAN
			PCMa.PCMp.omega = FNAN
			PCMa.PCMp.skew = FNAN
			PCMa.PCMp.T = FNAN
			PCMa.IterateJudge = 0.05 //	収束判定条件は前ステップ見かけの比熱の5%
			PCMa.Spctype = 'm'       // モデル形式をデフォルトとする
			PCMa.Condtype = 'm'

			ct := &PCMa.Chartable
			ct[0].PCMchar = 'E' // 引数0はエンタルピー
			ct[1].PCMchar = 'C' // 引数1は熱伝導率
			for i := 0; i < 2; i++ {
				ct[i].itablerow = 0
				ct[i].T = nil
				ct[i].Chara = nil
				ct[i].filename = ""
				ct[i].tabletype = 'e'
				ct[i].fp = nil
				ct[i].lowA = 0.0
				ct[i].lowB = 0.0
				ct[i].upA = 0.0
				ct[i].upB = 0.0
				ct[i].minTempChng = 0.5 // 最低温度変動幅
			}

			(*pcm)[j] = PCMa
		}
	}

	for i := 0; i < N; i++ {

		PCMa := (*pcm)[i]

		s = fi.GetToken()
		if s == "\n" {
			s = fi.GetToken()
		}

		if s[0] == '*' {
			break
		}

		PCMa.Name = s // PCM名称

		for {
			s = fi.GetToken()

			ce := strings.IndexByte(s, ';')
			if ce >= 0 {
				s = s[:ce]
			}
			if st := strings.IndexByte(s, '='); st >= 0 {
				dt, _ := strconv.ParseFloat(s[st+1:], 64)
				if strings.HasPrefix(s, "Ql") { // 潜熱量[J/m3]
					PCMa.Ql = dt
				} else if strings.HasPrefix(s, "spcheattable") {
					PCMa.Chartable[0].filename = s[st+1:]
					PCMa.Spctype = 't'
				} else if strings.HasPrefix(s, "table") {
					if s[st+1] == 'e' {
						PCMa.Chartable[0].tabletype = 'e'
					} else if s[st+1] == 'h' {
						PCMa.Chartable[0].tabletype = 'h'
					}
				} else if strings.HasPrefix(s, "conducttable") {
					PCMa.Chartable[1].filename = s[st+1:]
					PCMa.Condtype = 't'
				} else if strings.HasPrefix(s, "minTempChange") {
					PCMa.Chartable[0].minTempChng = dt
					PCMa.Chartable[1].minTemp = dt
				} else if strings.HasPrefix(s, "Condl") { // 液相熱伝導率[W/mK]
					PCMa.Condl = dt
				} else if strings.HasPrefix(s, "Conds") { // 固相熱伝導率[W/mK]
					PCMa.Conds = dt
				} else if strings.HasPrefix(s, "Crol") { // 液相容積比熱[J/m3K]
					PCMa.Crol = dt
				} else if strings.HasPrefix(s, "Cros") { // 固相容積比熱[J/m3K]
					PCMa.Cros = dt
				} else if strings.HasPrefix(s, "Tl") { // 液体から凝固が始まる温度[℃]
					PCMa.Tl = dt
				} else if strings.HasPrefix(s, "Ts") { // 固体から融解が始まる温度[℃]
					PCMa.Ts = dt
				} else if strings.HasPrefix(s, "Tp") { // 見かけの比熱のピーク温度[℃]
					PCMa.Tp = dt
				} else if strings.HasPrefix(s, "DivTemp") { // 比熱数値積分時の温度分割数
					PCMa.DivTemp = int(dt)
				} else if strings.HasPrefix(s, "Ctype") { // 見かけの比熱の特性曲線番号
					PCMa.Ctype = int(dt)
				} else if strings.HasPrefix(s, "a") { // これ以降は見かけの比熱計算のためのパラメータ
					PCMa.PCMp.a = dt
				} else if strings.HasPrefix(s, "b=") {
					PCMa.PCMp.b = dt
				} else if strings.HasPrefix(s, "c") {
					PCMa.PCMp.c = dt
				} else if strings.HasPrefix(s, "d") {
					PCMa.PCMp.d = dt
				} else if strings.HasPrefix(s, "e") {
					PCMa.PCMp.e = dt
				} else if strings.HasPrefix(s, "f") {
					PCMa.PCMp.f = dt
				} else if strings.HasPrefix(s, "B") {
					PCMa.PCMp.B = dt
				} else if strings.HasPrefix(s, "T") {
					PCMa.PCMp.T = dt
				} else if strings.HasPrefix(s, "bs") {
					PCMa.PCMp.bs = dt
				} else if strings.HasPrefix(s, "bl") {
					PCMa.PCMp.bl = dt
				} else if strings.HasPrefix(s, "skew") {
					PCMa.PCMp.skew = dt
				} else if strings.HasPrefix(s, "omega") {
					PCMa.PCMp.omega = dt
				} else if strings.HasPrefix(s, "nWieght") {
					PCMa.NWeight = dt
				} else if strings.HasPrefix(s, "IterateJudge") {
					PCMa.IterateJudge = dt
				} else {
					Eprint("<PCMdata>", s)
				}
			} else {
				if s == "-iterate" {
					PCMa.Iterate = true
					*pcmiterate = 'y'
				} else if s == "-pcmnode" {
					PCMa.AveTemp = 'n'
				} else {
					Eprint("<PCMdata>", s)
				}
			}

			if ce != -1 {
				break
			}
		}

		// テーブルの読み込み（見かけの比熱）
		if PCMa.Spctype == 't' {
			TableRead(&PCMa.Chartable[0])
		}

		// テーブルの読み込み（熱伝導率）
		if PCMa.Condtype == 't' {
			TableRead(&PCMa.Chartable[1])
		}

		// 入力情報のチェック
		if PCMa.Spctype == 'm' {
			var Tin, Bin, bsin, blin, skewin, omegain, ain, bin, cin, din, ein, fin int
			var Qlin, Condlin, Condsin, Crolin, Crosin, Tlin, Tsin, Tpin int

			Tin = dparaminit(PCMa.PCMp.T)
			Bin = dparaminit(PCMa.PCMp.B)
			bsin = dparaminit(PCMa.PCMp.bs)
			blin = dparaminit(PCMa.PCMp.bl)
			skewin = dparaminit(PCMa.PCMp.skew)
			omegain = dparaminit(PCMa.PCMp.omega)
			ain = dparaminit(PCMa.PCMp.a)
			bin = dparaminit(PCMa.PCMp.b)
			cin = dparaminit(PCMa.PCMp.c)
			din = dparaminit(PCMa.PCMp.d)
			ein = dparaminit(PCMa.PCMp.e)
			fin = dparaminit(PCMa.PCMp.f)
			Qlin = dparaminit(PCMa.Ql)
			Condlin = dparaminit(PCMa.Condl)
			Condsin = dparaminit(PCMa.Conds)
			Crolin = dparaminit(PCMa.Crol)
			Crosin = dparaminit(PCMa.Cros)
			Tlin = dparaminit(PCMa.Tl)
			Tsin = dparaminit(PCMa.Ts)
			Tpin = dparaminit(PCMa.Tp)

			// 必須入力項目のチェック
			if Condlin+Condsin+Crolin+Crosin+Tlin+Tsin != 0 {
				fmt.Printf("<PCMdata> name=%s Condl=%f Conds=%f Crol=%f Cros=%f Tl=%f Ts=%f\n",
					PCMa.Name, PCMa.Condl, PCMa.Conds, PCMa.Crol, PCMa.Cros, PCMa.Tl, PCMa.Ts)
			}

			// モデルごとに入力値をチェックする
			if PCMa.Ctype == 1 || PCMa.Ctype == 2 {
				if Qlin+Tsin+Tlin != 0 {
					fmt.Printf("<PCMdata> name=%s Ql=%f Ts=%f Tl=%f\n",
						PCMa.Name, PCMa.Ql, PCMa.Ts, PCMa.Tl)
				}
			} else if PCMa.Ctype == 3 {
				if Qlin+Tpin+Tin+Bin != 0 {
					fmt.Printf("<PCMdata> name=%s Ql=%f Tp=%f T=%f B=%f\n",
						PCMa.Name, PCMa.Ql, PCMa.Tp, PCMa.PCMp.T, PCMa.PCMp.B)
				}
			} else if PCMa.Ctype == 4 {
				if Tpin+ain+bin != 0 {
					fmt.Printf("<PCMdata> name=%s Tp=%f a=%f b=%f\n",
						PCMa.Name, PCMa.Tp, PCMa.PCMp.a, PCMa.PCMp.b)
				}
			} else if PCMa.Ctype == 5 {
				if Tpin+bsin+blin+ain != 0 {
					fmt.Printf("<PCMdata> name=%s Tp=%f bs=%f bl=%f a=%f\n",
						PCMa.Name, PCMa.Tp, PCMa.PCMp.bs, PCMa.PCMp.bl, PCMa.PCMp.a)
				}
			} else if PCMa.Ctype == 6 {
				if Qlin+Tpin+skewin+omegain != 0 {
					fmt.Printf("<PCMdata> name=%s Ql=%f Tp=%f skew=%f omega=%f\n",
						PCMa.Name, PCMa.Ql, PCMa.Tp, PCMa.PCMp.skew, PCMa.PCMp.omega)
				}
			} else if PCMa.Ctype == 7 {
				if ain+bin+cin+din+ein+fin != 0 {
					fmt.Printf("<PCMdata> name=%s a=%f b=%f c=%f d=%f e=%f f=%f\n",
						PCMa.Name, PCMa.PCMp.a, PCMa.PCMp.b, PCMa.PCMp.c, PCMa.PCMp.d, PCMa.PCMp.e, PCMa.PCMp.f)
				}
			}
		}
	}
}

// PCMの物性値テーブルの読み込み
func TableRead(ct *CHARTABLE) {
	if ct.filename == "" {
		return
	}

	fp, err := os.Open(ct.filename)
	if err != nil {
		fmt.Printf("<PCMdata> xxxx file not found %s xxxx\n", ct.filename)
		return
	}
	ct.fp = fp

	// 設定されている行数を数える
	var c byte
	var row int
	row = 0

	// Count the number of rows
	for {
		_, err := fmt.Fscanf(ct.fp, "%c", &c)
		if err != nil {
			break
		}
		if c == '\n' {
			row++
		}
	}

	// Close the file temporarily
	err = ct.fp.Close()
	if err != nil {
		fmt.Println("<PCMdata> ファイルのクローズに失敗")
		return
	}

	// Allocate memory
	ct.T = make([]float64, row)
	ct.Chara = make([]float64, row)
	if ct.T == nil || ct.Chara == nil {
		fmt.Println("<PCMdata> メモリの確保に失敗")
	}

	// Reopen the file
	fp, err = os.Open(ct.filename)
	if err != nil {
		fmt.Println("<PCMdata> ファイルのオープンに失敗")
		return
	}
	ct.fp = fp

	var st, tt int
	var spheat, prevheat, prevTemp float64
	prevheat = 0.0
	prevTemp = FNAN
	spheat = 0.0

	for i := 0; i < row; i++ {
		T := &ct.T[i]
		Char := &ct.Chara[i]

		// 温度の読み込み
		_, err := fmt.Fscanf(ct.fp, " %f ", T)
		if err != nil {
			fmt.Println("<PCMdata> 温度の読み込みに失敗")
			break
		}
		// テーブルの下限温度
		if st == 0 {
			ct.minTemp = *T
			prevTemp = *T
			st = 1
		}
		// 昇順に並んでいないときのエラー表示
		if *T <= prevTemp && tt == 1 {
			fmt.Printf("i=%d 温度データが昇順に並んでいません T=%f preT=%f\n", i, *T, prevTemp)
		}
		var dblTemp float64
		// 特性値の読み込み
		_, err = fmt.Fscanf(ct.fp, " %f ", &dblTemp)
		if err != nil {
			fmt.Println("<PCMdata> 特性値の読み込みに失敗")
			break
		}
		*Char = dblTemp
		// 見かけの比熱の場合
		if ct.tabletype == 'h' {
			*Char = prevheat + spheat*(*T-prevTemp)
		}
		prevheat = *Char
		prevTemp = *T
		spheat = dblTemp
		// テーブルの上限温度
		ct.maxTemp = *T
		ct.itablerow++
		tt = 1
	}

	// Close the file
	err = ct.fp.Close()
	if err != nil {
		fmt.Println("<PCMdata> ファイルのクローズに失敗")
	}

	// 上下限温度範囲以外の線形回帰式を作成
	// 下限以下
	ct.lowA = (ct.Chara[0] - ct.Chara[1]) / (ct.T[0] - ct.T[1])
	ct.lowB = ct.Chara[0] - ct.lowA*ct.T[0]
	// 上限以上
	ct.upA = (ct.Chara[ct.itablerow-1] - ct.Chara[ct.itablerow-2]) / (ct.T[ct.itablerow-1] - ct.T[ct.itablerow-2])
	ct.upB = ct.Chara[ct.itablerow-1] - ct.upA*ct.T[ct.itablerow-1]
}

// 初期化されているかをチェックする
func dparaminit(A float64) int {
	if math.Abs(A-(FNAN)) < 1e-5 {
		return 1
	} else {
		return 0
	}
}

func PCMcount(fi *EeTokens) int {
	N := 0
	add := fi.GetPos()

	for fi.IsEnd() == false {
		s := fi.GetToken()
		if strings.HasPrefix(s, "*") {
			break
		}
		if strings.HasPrefix(s, ";") {
			N++
		}
	}

	fi.RestorePos(add)
	return N
}

// 固相、液相物性値と潜熱量からPCM温度の物性値を計算する（比熱、熱伝導率共通）
// 熱伝導率の計算時はQl=0とする
func FNPCMState(Ctype int, Ss, Sl, Ql, Ts, Tl, Tp, T float64, PCMp *PCMPARAM) float64 {
	var Qse, Qla, Tls float64

	Tls = Tl - Ts
	// 顕熱分の補間
	if T > Ts && T < Tl {
		Qse = Ss + (Sl-Ss)/Tls*(T-Ts)
	} else if T <= Ts {
		Qse = Ss
	} else {
		Qse = Sl
	}

	// 潜熱分の補間
	Qla = 0.0

	Ql = 0.0

	// 熱伝導率計算の場合は潜熱ゼロ
	if Ctype == 0 {
		Qla = 0.0
	} else if Ctype == 1 { // 潜熱変化域潜熱比熱一定
		if T > Ts && T < Tl {
			Qla = Ql / (Tl - Ts)
		}
	} else if Ctype == 2 { // 二等辺三角形
		if T > Ts && T < Tl {
			if T < (Tl+Ts)/2.0 {
				Qla = 4.0 * Ql / (Tls * Tls) * (T - Ts)
			} else {
				Qla = 4.0 * Ql / (Tls * Tls) * (Tl - T)
			}
		}
	} else if Ctype == 3 { // 双曲線関数
		Temp := math.Cosh((2.0 * PCMp.B / PCMp.T) * (T - Tp))
		Qla = 0.5 * Ql * (2.0 * PCMp.B / PCMp.T) / (Temp * Temp)
	} else if Ctype == 4 { // ガウス関数（対象）
		Temp := (T - Tp) / PCMp.B
		Qla = PCMp.a * math.Exp(-0.5*Temp*Temp)
	} else if Ctype == 5 { // ガウス関数（非対称）
		Temp := (T - Tp) / PCMp.B
		if T <= Tp {
			Temp = PCMp.bs
		} else {
			Temp = PCMp.bl
		}
		Qla = PCMp.a * math.Exp(-((T-Tp)/Temp)*((T-Tp)/Temp))
	} else if Ctype == 6 { // 誤差関数歪度
		//Temp := math.Exp(-(T - Tp) * (T - Tp) / ((2.0 * PCMp.omega) * PCMp.omega))
		Qla = Ql / math.Sqrt(2.0*math.Pi) * math.Exp(-(T-Tp)*(T-Tp)/((2.0*PCMp.omega)*PCMp.omega)) * (1.0 + math.Erf(PCMp.skew*(T-Tp)/(math.Sqrt(2.0)*PCMp.omega)))
	} else if Ctype == 7 { // 有理関数
		if T < 0 {
			Qla = 0.0
		} else {
			Qla = math.Pow(T, PCMp.f) * (PCMp.a*T*T + PCMp.B*T + PCMp.c) / (PCMp.d*T*T + PCMp.e*T + 1.0)
		}
	}

	return (Qse + Qla)
}

/* ------------------------------------------------------ */

// PCMの状態値計算（家具用）

func FNPCMStatefun(Ctype int, Ss, Sl, Ql, Ts, Tl, Tp, oldT, T float64, DivTemp int, PCMp *PCMPARAM) float64 {
	var dblTemp, dTemp, TPCM, Qld, Qllst float64
	dTemp = (T - oldT) / float64(DivTemp)

	Qllst = 0.0
	// 前時刻からの温度変化が小さい場合は時間積分はしない
	if math.Abs(T-oldT) < 1e-4 {
		Qllst = FNPCMState(Ctype, Ss, Sl, Ql, Ts, Tl, Tp, (T+oldT)*0.5, PCMp)
	} else {
		// 見かけの比熱の時間積分
		for i := 0; i < DivTemp+1; i++ {
			TPCM = oldT + dTemp*float64(i)
			Qld = FNPCMState(Ctype, Ss, Sl, Ql, Ts, Tl, Tp, TPCM, PCMp)
			dblTemp += Qld
		}

		Qllst = dblTemp / float64(DivTemp+1)
	}

	return Qllst
}

// PCMの温度から見かけの比熱を求める（テーブル形式）
// Told:前時刻のPCM温度、T:暫定現在時刻PCM温度
func FNPCMstate_table(ct *CHARTABLE, Told, T float64, Ndiv int) float64 {
	var oldEn, En, Chara float64

	// 前時刻と暫定現在時刻のPCM温度が同温の場合
	if math.Abs(T-Told) < ct.minTempChng {
		Tave := 0.5 * (T + Told)
		// 見かけの比熱の計算
		if ct.PCMchar == 'E' {
			oldEn = FNPCMenthalpy_table_lib(ct, Tave-0.5*ct.minTempChng)
			En = FNPCMenthalpy_table_lib(ct, Tave+0.5*ct.minTempChng)
			// 見かけの比熱の計算
			Chara = (En - oldEn) / ct.minTempChng
		} else {
			// 熱伝導率の計算
			Chara = FNPCMenthalpy_table_lib(ct, Tave)
		}
	} else {
		if ct.PCMchar == 'E' {
			oldEn = FNPCMenthalpy_table_lib(ct, Told)
			En = FNPCMenthalpy_table_lib(ct, T)
			// 見かけの比熱の計算
			Chara = (En - oldEn) / (T - Told)
		} else {
			// 熱伝導率の計算
			var dblTemp, Tpcm, dTemp float64
			dblTemp = 0.0
			dTemp = (T - Told) / float64(Ndiv)
			// ToldからTまでを積分する
			for i := 0; i < Ndiv; i++ {
				Tpcm = Told + dTemp*float64(i)
				dblTemp += FNPCMenthalpy_table_lib(ct, Tpcm)
			}
			Chara = dblTemp / float64(Ndiv+1)
		}
	}

	return Chara
}

func FNPCMenthalpy_table_lib(ct *CHARTABLE, T float64) float64 {
	var prevTpcm, Tpcm, enthalpy, preventhalpy *float64
	var retVal float64

	if T < ct.minTemp {
		// テーブルの最低温度よりPCM温度が低い場合は端部の特性値を線形で外挿
		retVal = ct.lowA*T + ct.lowB
	} else if T > ct.maxTemp {
		//テーブルの最高温度よりPCM温度が高い場合は端部の特性値を線形で外挿
		retVal = ct.upA*T + ct.upB
	} else {
		for i := 0; i < ct.itablerow; i++ {

			Tpcm := &ct.T[i]
			enthalpy := &ct.Chara[i]

			if *Tpcm > T {
				break
			}
			prevTpcm = Tpcm
			preventhalpy = enthalpy
		}
		// 線形補間
		retVal = *preventhalpy + (*enthalpy-*preventhalpy)*(T-*prevTpcm)/(*Tpcm-*prevTpcm)
	}

	return retVal
}

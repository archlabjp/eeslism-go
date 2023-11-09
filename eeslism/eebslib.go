package eeslism

import (
	"fmt"
	"math"
	"os"
	"strings"
)

// 外表面方位デ－タの入力
func Exsfdata(section *EeTokens, dsn string, Exsf *EXSFS, Schdl *SCHDL, Simc *SIMCONTL) {
	var ename string
	//var st *string
	var dt, wa, wb, swa, cwa, swb, cwb float64
	var k int
	var err error

	// 外表面総合伝達率のデフォルト値
	Exsf.Alosch = envptr(fmt.Sprintf("%f", ALO), Simc, nil, nil, nil)
	Exsf.Alotype = Alotype_Fix // 固定値
	Exsf.Exs = make([]*EXSF, 0)

	var dfrg float64 // 全面地物の日射反射率（デフォルト値）

	// 最初の行: 全体に対する設定
	line := section.GetLogicalLine()
	for _, s := range line {
		if strings.HasPrefix(s, "alo=") {
			// 外表面総合伝達率[W/m2K]
			//
			value := s[4:]
			if value == "Calc" {
				// 風速から計算する
				Exsf.Alotype = Alotype_V
			} else if k, err = idsch(value, Schdl.Sch, ""); err == nil {
				// スケジュールに基づいて値を変化させる
				Exsf.Alosch = &Schdl.Val[k]
				Exsf.Alotype = Alotype_Schedule
			} else {
				// 数値 or 内部変数名
				Exsf.Alosch = envptr(value, Simc, nil, nil, nil)
				if Exsf.Alosch != nil {
					Exsf.Alotype = Alotype_Schedule
				}
			}
		} else if strings.HasPrefix(s, "r=") {
			// 全面地物の日射反射率[-]
			//
			value := s[2:]
			dfrg, err = readFloat(value)
			if err != nil || dfrg < 0.0 || dfrg > 1.0 {
				fmt.Fprintf(os.Stderr, "%s の設置値が不適切です", s)
				os.Exit(1)
			}
		}
	}

	for section.IsEnd() == false {
		// 論理行を取得
		line = section.GetLogicalLine()

		// "*"が出てきたら終了
		if line[0] == "*" {
			break
		}

		// 外表面データの初期化
		ex := new(EXSF)
		Exsfinit(ex)
		ex.Name = line[0]         // 外表面名
		ex.Alotype = Exsf.Alotype // 外表面総合伝達率の設定方法=デフォルト値
		ex.Alo = Exsf.Alosch      // 外表面総合伝達率=デフォルト値

		// 特殊名 Hor, EarchSf への対応
		if line[0] == "Hor" {
			// 水平面なので、傾斜角は0
			ex.Wb = 0.0
		} else if line[0] == "EarthSf" {
			// 地表面境界を含むことをフラグに書き込む
			Exsf.EarthSrfFlg = true

			// 種別: 地表面
			ex.Typ = EXSFType_e
		} else {
			// 通常の表面の定義
			ex.Wb = 90.0 // 傾斜角=90°
			ex.Rg = dfrg // 日射反射率=デフォルト値
		}

		for _, s := range line[1:] {
			st := strings.IndexRune(s, '=')

			key := s[:st]
			value := s[st+1:]

			if key == "a" {
				// *** 方位角 a***
				//
				var err error
				if dt, err = readFloat(value); err == nil {
					// 数値指定
					ex.Wa = dt
				} else {
					var dir rune = ' '
					if strings.Contains(s, "+") {
						st := strings.IndexRune(s, '+')
						dir = '+'
						ename = s[2:st]
						offvalue := s[st+1:]
						dt, err = readFloat(offvalue)
						if err != nil {
							panic(err)
						}
					} else if strings.Contains(s, "-") {
						st := strings.IndexRune(s, '-')
						dir = '-'
						ename = s[2:st]
						offvalue := s[st+1:]
						dt, err = readFloat(offvalue)
						if err != nil {
							panic(err)
						}
					} else {
						ename = s[2:]
					}

					var found_flag = false
					for _, exj := range Exsf.Exs {
						if exj.Name == ename {
							if dir == '+' {
								ex.Wa = exj.Wa + dt
							} else if dir == '-' {
								ex.Wa = exj.Wa - dt
							} else {
								ex.Wa = exj.Wa
							}
							found_flag = true
							break
						}
					}
					if !found_flag {
						Eprint("<Exsfdata>", s)
					}
				}
			} else if key == "alo" {
				// *** 外表面熱伝達率 alo ***
				//
				if value == "Calc" {
					// 風速から計算
					ex.Alotype = Alotype_V
				} else {
					// スケジュール
					ex.Alotype = Alotype_Schedule
					if k, err = idsch(value, Schdl.Sch, ""); err == nil {
						ex.Alo = &Schdl.Val[k]
					} else {
						ex.Alo = envptr(value, Simc, nil, nil, nil)
					}
				}
			} else {
				// *** 傾斜角 t,日射反射率 r,地中深さ Z,土の熱拡散率 d ***
				//
				dt, err = readFloat(value)
				if err != nil {
					panic(s)
				}
				switch key {
				case "t":
					// 傾斜角[°]
					ex.Wb = dt
				case "r":
					// 全面地物の日射反射率[-]
					ex.Rg = dt
				case "Z":
					// 地中深さ [m]
					ex.Z = dt
					ex.Typ = EXSFType_E // 地下扱いにする
				case "d":
					// 土の熱拡散率 [m2/s]
					ex.Erdff = dt
				default:
					Eprint("<Exsfdata>", s)
				}
			}
		}

		Exsf.Exs = append(Exsf.Exs, ex)
	}

	// 外表面熱伝達率の設定
	// if len(Exsf.Exs) > 0 {
	// 	s = strconv.FormatFloat(ALO, 'f', -1, 64)
	// 	Exsf.Alosch = envptr(s, Simc, nil, nil, nil)
	// 	Exsf.Alotype = Alotype_Fix // 固定値
	// 	Exsf.Exs = exs
	// }

	for _, ex := range Exsf.Exs {

		// 一般外表面 の場合は、日射に関するパラメータを計算する
		if ex.Typ == EXSFType_S {
			const rad = PI / 180.
			wa = ex.Wa * rad          // 方位角	[rad]
			wb = ex.Wb * rad          // 傾斜角	[rad]
			cwa = math.Cos(wa)        // 方位角の余弦
			swa = math.Sin(wa)        // 方位角の正弦
			cwb = math.Cos(wb)        // 傾斜角の余弦
			swb = math.Sin(wb)        // 傾斜角の正弦
			ex.Cwa = cwa              // = 方位角の余弦
			ex.Swa = swa              // = 方位角の正弦
			ex.Swb = swb              // = 傾斜角の正弦
			ex.Wz = cwb               // = 傾斜角の余弦
			ex.Ww = swb * swa         // = 傾斜角の正弦 ×  方位角の正弦
			ex.Ws = swb * cwa         // = 傾斜角の正弦 ×  方位角の余弦
			ex.CbSa = cwb * swa       // = 傾斜角の余弦 ×  方位角の正弦
			ex.CbCa = cwb * cwa       // = 傾斜角の余弦 ×  方位角の正弦
			ex.Fs = 0.5 * (1.0 + cwb) // 天空を見る形態係数
		}
	}
}

/*  外表面入射日射量の計算    */
func (exsf *EXSFS) Exsfsol(Wd *WDAT) {
	for _, ex := range exsf.Exs {

		if ex.Typ == EXSFType_S {
			// 入射角のcos
			cinc := Wd.Sh*ex.Wz + Wd.Sw*ex.Ww + Wd.Ss*ex.Ws

			if cinc > 0.0 {
				// 太陽が出ている場合

				// プロファイル角の計算
				ex.Tprof = (Wd.Sh*ex.Swb - Wd.Sw*ex.CbSa - Wd.Ss*ex.CbCa) / cinc
				ex.Prof = math.Atan(ex.Tprof)

				// 見かけの方位角の計算
				ex.Tazm = (Wd.Sw*ex.Cwa - Wd.Ss*ex.Swa) / cinc
				ex.Gamma = math.Atan(ex.Tazm)
				ex.Cinc = cinc
			} else {
				// 太陽が出ていない場合
				ex.Prof = 0.0
				ex.Gamma = 0.0
				ex.Cinc = 0.0
			}

			// 日射量の計算
			ex.Idre = Wd.Idn * ex.Cinc                         // 直逹日射  [W/m2]
			ex.Idf = Wd.Isky*ex.Fs + ex.Rg*Wd.Ihor*(1.0-ex.Fs) // 拡散日射  [W/m2]
			ex.Iw = ex.Idre + ex.Idf                           // 全日射    [W/m2]
			ex.Rn = Wd.RN * ex.Fs                              // 夜間輻射  [W/m2]
		}
	}
}

// ガラス日射熱取得の計算
// 入力:
//   面積 Ag [m2]
//   日射総合取得率 tgtn [-]
//   吸収日射取得率 Bn [-]
//   入射角のcos cinc [-]
//   ********** Fsdw
//   直逹日射 Idr [W/m2]
//   拡散日射 Idf [W/m2]
// 出力:
//   透過日射熱取得 Qgt [W]
//   吸収日射熱取得 Qga [W]
func Glasstga(Ag, tgtn, Bn, cinc, Fsdw, Idr, Idf float64, Cidtype string, Profile, Gamma float64) (Qgt, Qga float64) {
	var Cid, Cidf, Bid, Bidf float64

	Cid = 0.0
	Bid = 0.0
	Cidf = 0.01
	Bidf = 0.0

	// 標準
	if Cidtype == "N" {
		Cid = Glscid(cinc)
		Cidf = 0.91

		Bid = Cid
		Bidf = Cidf
	} else {
		fmt.Printf("xxxxx <eebslib.c  CidType=%s\n", Cidtype)
	}

	// 透過日射量の計算
	Qt := Ag * (Cid*Idr*(1.0-Fsdw) + Cidf*Idf)

	// 吸収日射量の計算
	Qb := Ag * (Bid*Idr*(1.0-Fsdw) + Bidf*Idf)

	Qgt = Qt * tgtn
	Qga = Qb * Bn

	return Qgt, Qga
}

// ガラスの直達日射透過率標準特性
// 入力:
//   入射角のcos cinc [-]
// 出力:
//   ガラスの直達日射透過率標準特性 Cid [-]
func Glscid(cinc float64) float64 {
	return math.Max(0, cinc*(3.4167+cinc*(-4.389+cinc*(2.4948-0.5224*cinc))))
}

// ガラスの直達日射透過率標準特性(普通複層ガラス用)
// 入力:
//   入射角のcos cinc [-]
// 出力:
//   ガラスの直達日射透過率標準特性 Cid [-]
func GlscidDG(cinc float64) float64 {
	return math.Max(0, cinc*(0.341819+cinc*(6.070709+cinc*(-9.899236+4.495774*cinc))))
}

func ExsfCount(section *EeTokens) int {
	var N int
	for section.IsEnd() == false {
		if section.GetToken() == ";" {
			N++
			break
		}
	}
	section.Reset()
	return N
}

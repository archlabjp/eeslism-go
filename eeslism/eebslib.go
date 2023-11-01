package eeslism

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

// 外表面方位デ－タの入力
func Exsfdata(section *EeTokens, dsn string, Exsf *EXSFS, Schdl *SCHDL, Simc *SIMCONTL) {
	var s, ename string
	//var st *string
	var dt, dfrg, wa, wb, swa, cwa, swb, cwb float64
	var i, j, k int
	var err error

	s = dsn
	Nd := ExsfCount(section)
	if Nd == 0 {
		Nd = 1
	}

	var exs []*EXSF
	if Nd > 0 {
		exs = make([]*EXSF, 0, Nd+1)

		s = fmt.Sprintf("%f", ALO)
		Exsf.Alosch = envptr(s, Simc, nil, nil, nil)
		Exsf.Alotype = Alotype_Fix // 固定値
	}

	line := section.GetLogicalLine()

	for _, s := range line {
		if strings.HasPrefix(s, "alo=") {
			if s[4:] == "Calc" {
				Exsf.Alotype = Alotype_V
			} else if k, err = idsch(s[4:], Schdl.Sch, ""); err == nil {
				Exsf.Alosch = &Schdl.Val[k]
				Exsf.Alotype = Alotype_Schedule
			} else {
				Exsf.Alosch = envptr(s[4:], Simc, nil, nil, nil)
				if Exsf.Alosch != nil {
					Exsf.Alotype = Alotype_Schedule
				}
			}
		} else if strings.HasPrefix(s, "r=") {
			dfrg, _ = strconv.ParseFloat(s[2:], 64)
			if dfrg < 0.0 || dfrg > 1.0 {
				fmt.Fprintf(os.Stderr, "%s の設置値が不適切です", s)
				os.Exit(1)
			}
		}
	}

	for section.IsEnd() == false {

		line = section.GetLogicalLine()
		if line[0] == "*" {
			break
		}

		ex := new(EXSF)
		Exsfinit(ex)

		i++
		ex.Name = line[0]
		ex.Alotype = Exsf.Alotype
		ex.Alo = Exsf.Alosch
		if s == "Hor" {
			ex.Wb = 0.0
		} else if s == "EarthSf" {
			Exsf.EarthSrfFlg = true
			ex.Typ = 'e'
		} else {
			ex.Wb = 90.0
			ex.Rg = dfrg
		}

		for _, s := range line[1:] {
			if strings.HasPrefix(s, "a=") {
				var err error
				if dt, err = strconv.ParseFloat(s[2:], 64); err == nil {
					ex.Wa = dt
				} else {
					var dir rune = ' '
					if strings.Contains(s, "+") {
						st := strings.IndexRune(s, '+')
						dir = '+'
						ename = s[2:st]
						dt, err = strconv.ParseFloat(s[st+1:], 64)
						if err != nil {
							panic(err)
						}
					} else if strings.Contains(s, "-") {
						st := strings.IndexRune(s, '-')
						dir = '-'
						ename = s[2:st]
						dt, err = strconv.ParseFloat(s[st+1:], 64)
						if err != nil {
							panic(err)
						}
					} else {
						ename = s[2:]
					}

					for _, exj := range exs {
						if exj.Name == ename {
							if dir == '+' {
								ex.Wa = exj.Wa + dt
							} else if dir == '-' {
								ex.Wa = exj.Wa - dt
							} else {
								ex.Wa = exj.Wa
							}
							break
						}
					}
					if j == i+1 {
						Eprint("<Exsfdata>", s)
					}
				}
			} else {
				st := strings.IndexRune(s, '=')

				if strings.HasPrefix(s, "alo") {
					// *** 外表面熱伝達率 alo ***
					if s[st+1:] == "Calc" {
						// 風速から計算
						ex.Alotype = Alotype_V
					} else {
						// スケジュール
						ex.Alotype = Alotype_Schedule
						if k, err = idsch(s[st+1:], Schdl.Sch, ""); err == nil {
							ex.Alo = &Schdl.Val[k]
						} else {
							ex.Alo = envptr(s[st+1:], Simc, nil, nil, nil)
						}
					}
				} else {
					dt, _ = readFloat(s[st+1:])
					switch s[0] {
					case 't':
						ex.Wb = dt
					case 'r':
						ex.Rg = dt
					case 'Z':
						ex.Z = dt
						ex.Typ = 'E'
					case 'd':
						ex.Erdff = dt
					default:
						Eprint("<Exsfdata>", s)
					}
				}
			}
		}

		exs = append(exs, ex)
	}

	// 外表面熱伝達率の設定
	if Nd > 0 {
		s = strconv.FormatFloat(ALO, 'f', -1, 64)
		Exsf.Alosch = envptr(s, Simc, nil, nil, nil)
		Exsf.Alotype = Alotype_Fix // 固定値
		Exsf.Exs = exs
	}

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
func Exsfsol(Wd *WDAT, Exs []*EXSF) {
	for _, ex := range Exs {

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
			ex.Iw = ex.Idre + ex.Idf                           // 全日射    [W/m2]w
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
	var Cid, Cidf, Bid, Bidf, Qt, Qb float64

	Cid = 0.0
	Bid = 0.0
	Cidf = 0.01
	Bidf = 0.0
	Qt = 0.0
	Qb = 0.0

	if Cidtype == "N" {
		Cid = Glscid(cinc)
		Cidf = 0.91

		Bid = Cid
		Bidf = Cidf
	} else {
		fmt.Printf("xxxxx <eebslib.c  CidType=%s\n", Cidtype)
	}

	Qt = Ag * (Cid*Idr*(1.0-Fsdw) + Cidf*Idf)
	Qb = Ag * (Bid*Idr*(1.0-Fsdw) + Bidf*Idf)

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

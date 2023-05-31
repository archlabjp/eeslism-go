package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

/*   ファイル、計算期間、出力日の入力        */

func Gdata(section *EeTokens, dsn string, File string, wfname *string, ofname *string, dtm *int, sttmm *int, dayxs *int, days *int, daye *int, Tini *float64, pday []int, wdpri *int, revpri *int, pmvpri *int, helmkey *rune, MaxIterate *int, Daytm *DAYTM, Wd *WDAT, perio *rune) {
	var s, ss, ce, dd string
	var st int
	var Ms, Ds, Mxs, Dxs, Me, De int
	var logprn int = 0

	*dtm = 3600
	*sttmm = -1

	*wfname = ""

	//E := fmt.Sprintf(ERRFMT, dsn)

	*ofname = File

	for i := 1; i < 366; i++ {
		pday[i] = 0
	}

	if st = strings.LastIndex(*ofname, "."); st != -1 {
		*ofname = (*ofname)[:st]
	}

	for section.IsEnd() == false {
		line := section.GetLogicalLine()

		if line[0] == "FILE" {
			for _, s := range line[1:] {
				if s == "-skyrd" { // 気象データは夜間放射量で定義されている
					Wd.RNtype = 'R'
				} else if s == "-intgtsupw" { // 給水温度を補間する
					Wd.Intgtsupw = 'Y'
				} else {
					if st := strings.IndexRune(s, '='); st != -1 {
						s1, s2 := s[:st], s[st+1:]

						var err error
						if s1 == "w" {
							_, err = fmt.Sscanf(s2, "%s", &dd)
							if err != nil {
								panic(err)
							}
							*wfname = dd
						} else if s1 == "out" {
							_, err = fmt.Sscanf(s2, "%s", &ss)
							if err != nil {
								panic(err)
							}

							//NOTE: 何をやりたいのか不明(UDA)
							const FLDELIM = "/"
							if strings.LastIndex(ss, FLDELIM) == -1 {
								*ofname = File

								if st := strings.LastIndex(*ofname, FLDELIM); st != -1 {
									s2 = ss
								}
							} else {
								*ofname = ss
							}
						} else {
							Eprint("<Gdata>", s)
						}
					}
				}

				if ce != "" {
					break
				}
			}
		} else if line[0] == "RUN" {
			*Tini = 15.0

			var err error
			for i := 1; i < len(line); i++ {
				s = line[i]
				if strings.HasPrefix(s, "Tinit") {
					kv := strings.SplitN(s, "=", 2)
					*Tini, err = strconv.ParseFloat(kv[1], 64)
					if err != nil {
						panic(err)
					}
				} else if st := strings.IndexRune(s, '='); st != -1 {
					s = s[st+1:]

					if s == "dTime" {
						// For `dTime=3600`
						*dtm, err = strconv.Atoi(s)
						if err != nil {
							panic(err)
						}
					} else if s == "Stime" {
						// For `Stime=0`
						*sttmm, err = strconv.Atoi(s)
						if err != nil {
							panic(err)
						}
						*sttmm *= 100
					} else if s == "MaxIterate" {
						// For `MaxIterate=100`
						*MaxIterate, err = strconv.Atoi(s)
						if err != nil {
							panic(err)
						}
					} else if s == "RepeatDays" { // 周期定常計算の繰り返し日数の取得
						// For `RepeatDays=365`
						var Ndays int
						Ndays, err = strconv.Atoi(s)
						if err != nil {
							panic(err)
						}

						if *perio != 'y' {
							fmt.Println("周期定常計算の指定がされていません")
						}

						*daye = *days + Ndays - 1
					}
				} else if s[0] == '(' {
					// For `(1/1)`
					_, err = fmt.Sscanf(s, "(%d/%d)", &Mxs, &Dxs)
					if err != nil {
						panic(err)
					}
					*dayxs = FNNday(Mxs, Dxs) // 助走計算開始日
				} else if s == "-periodic" { // 周期定常計算への対応
					// For `-periodic 1/1`
					*perio = 'y'  // 周期定常計算フラグの変更
					s = line[i+1] // 計算する日付の読み込み
					i++
					_, err = fmt.Sscanf(s, "%d/%d", &Ms, &Ds) // 計算する日付の取得
					if err != nil {
						panic(err)
					}
					*days = FNNday(Ms, Ds)
					*dayxs = *days // 助走計算開始日
					Daytm.Mon = Ms
					Daytm.Day = Ds
				} else if strings.IndexRune(s, '-') != -1 {
					// For `1/1-12/31`
					_, err = fmt.Sscanf(s, "%d/%d-%d/%d", &Ms, &Ds, &Me, &De)
					if err != nil {
						panic(err)
					}
					*days = FNNday(Ms, Ds)
					*daye = FNNday(Me, De)

					if Mxs == 0 {
						*dayxs = *days // 助走計算開始日
						Daytm.Mon = Ms
						Daytm.Day = Ds
					} else {
						Daytm.Mon = Mxs
						Daytm.Day = Dxs
					}
				} else {
					Eprint("<Gdata>", s)
				}

				if ce != "" {
					break
				}
			}
			fmt.Printf("<<Gdata>> dtm=%d\n", *dtm)
		} else if line[0] == "PRINT" {
			for _, s := range line[1:] {
				switch s {
				case "*wd":
					*wdpri = 1
				case "*rev":
					*revpri = 1
				case "*pmv":
					*pmvpri = 1
				case "*helm":
					*helmkey = 'y'
				case "*log":
					logprn = 1
				case "*debug":
					DEBUG = true
				default:
					if strings.IndexRune(s, '-') == -1 {
						var Ms, Ds int
						fmt.Sscanf(s, "%d/%d", &Ms, &Ds)
						pday[FNNday(Ms, Ds)] = 1
					} else {
						var Ms, Ds, Me, De, ns, ne, n int
						fmt.Sscanf(s, "%d/%d-%d/%d", &Ms, &Ds, &Me, &De)
						ns = FNNday(Ms, Ds)
						n = FNNday(Me, De)
						if ns < n {
							ne = n
						} else {
							ne = n + 365
						}
						for n = ns; n <= ne; n++ {
							if n > 365 {
								pday[n-365] = 1
							} else {
								pday[n] = 1
							}
						}
					}
				}
			}
		} else if line[0] == "*" {
			break
		} else {
			Eprint("<Gdata>", s)
		}
	}

	// Concatenate ".log" to the end of *ofname and copy to s
	s = filepath.Join(*ofname + ".log")

	// Open the file for writing
	var err error
	Ferr, err = os.Create(s)
	if err != nil {
		// Handle error
	}

	if logprn == 0 {
		// Close the file and set ferr to nil if logprn is 0
		Ferr.Close()
		Ferr = nil
	}
}

package eeslism

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

/*
Gdata (General Data Input for Simulation Control)

この関数は、建物のエネルギーシミュレーションの実行に関する全般的な設定データ（ファイル名、計算期間、出力設定など）を
入力ファイルから読み込み、対応する構造体に格納します。

建築環境工学的な観点:
- **シミュレーションの基本設定**: 建物のエネルギーシミュレーションは、
  特定の期間（年間など）にわたる建物の熱的挙動を予測するために行われます。
  この関数は、以下の基本設定を読み込みます。
  - `wfname`: 気象データファイル名。シミュレーションの外部環境条件を定義します。
  - `ofname`: 出力ファイル名。シミュレーション結果の保存先を定義します。
  - `dtm`: 計算時間間隔 [s]。シミュレーションの細かさを決定し、計算精度と計算時間に影響します。
  - `sttmm`: 計算開始時刻。シミュレーションの開始時点を定義します。
  - `dayxs`, `days`, `daye`: 助走期間、計算開始日、計算終了日。シミュレーション期間を定義します。
  - `Tini`: 初期温度。シミュレーション開始時の建物の初期状態を定義します。
  - `MaxIterate`: 最大収束計算回数。非線形な熱的挙動を持つシステムにおいて、
    計算の収束を確保するための反復回数を定義します。
- **気象データ処理のオプション**: `skyrd`（夜間放射量で定義された気象データ）、
  `intgtsupw`（給水温度の補間）などのオプションは、
  気象データの特性や、特定の熱負荷計算に必要な補助データを処理する方法を定義します。
- **出力設定の制御**: `PRINT`セクションでは、
  - `*wd`: 気象データの出力。
  - `*rev`: 熱損失係数の出力。
  - `*pmv`: PMV（予測平均申告）の出力。
  - `*helm`: 要素別熱損失・熱取得の出力。
  - `*log`: ログファイルの出力。
  - `*debug`: デバッグモードのON/OFF。
  など、シミュレーション結果の出力内容を詳細に制御できます。
  これにより、必要な情報を効率的に取得し、
  分析や検証を容易にします。
- **周期定常計算 (periodic)**:
  `periodic`オプションは、周期定常計算を行うかどうかを定義します。
  周期定常計算は、建物の熱的挙動が日単位で繰り返されると仮定し、
  短期間のシミュレーションで年間を通じた熱的挙動を推定する際に用いられます。

この関数は、建物のエネルギーシミュレーションの実行を制御し、
シミュレーションの正確性、効率性、および出力内容を決定するための重要な役割を果たします。
*/
func Gdata(section *EeTokens, File string, wfname *string,
	ofname *string, dtm *int, sttmm *int, dayxs *int, days *int, daye *int,
	Tini *float64, pday []int, wdpri *int, revpri *int, pmvpri *int,
	helmkey *rune, MaxIterate *int, Daytm *DAYTM, Wd *WDAT, perio *rune) {
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
			for _, s := range line[1 : len(line)-1] {
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
					key := s[:st]
					value := s[st+1:]

					if key == "dTime" {
						// For `dTime=3600`
						*dtm, err = strconv.Atoi(value)
						if err != nil {
							panic(err)
						}
					} else if key == "Stime" {
						// For `Stime=0`
						*sttmm, err = strconv.Atoi(value)
						if err != nil {
							panic(err)
						}
						*sttmm *= 100
					} else if key == "MaxIterate" {
						// For `MaxIterate=100`
						*MaxIterate, err = strconv.Atoi(value)
						if err != nil {
							panic(err)
						}
					} else if key == "RepeatDays" { // 周期定常計算の繰り返し日数の取得
						// For `RepeatDays=365`
						var Ndays int
						Ndays, err = strconv.Atoi(value)
						if err != nil {
							panic(err)
						}

						if *perio != 'y' {
							fmt.Println("周期定常計算の指定がされていません")
						}

						*daye = *days + Ndays - 1
					} else {
						panic(s)
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
			for _, s := range line[1 : len(line)-1] {
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

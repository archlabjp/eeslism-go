package eeslism

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

/*
Pathdata (System Path Data Input)

この関数は、建物のエネルギーシミュレーションにおける熱媒（空気、水など）の
流れる経路（パス）に関するデータを入力ファイルから読み込み、
対応する構造体（`MPATH`, `PLIST`, `PELM`）に格納します。
これにより、熱搬送システムや空調システムのシステム構成を詳細にモデル化できます。

建築環境工学的な観点:
- **システム構成のモデル化**: 建物のエネルギーシステムは、
  熱源設備、熱搬送設備、空調設備など、様々な機器が配管やダクトで接続されて構成されます。
  この関数は、これらの複雑なシステム構成を階層的にモデル化します。
  - `MPATH`: システム経路。熱媒の種類（水、空気）、システムの種類（空調、給湯）などを定義します。
  - `PLIST`: 末端経路。MPATHを構成する個々の経路であり、
    経路内の機器（`PELM`）や流量に関する情報（`Go`, `Rate`）を定義します。
  - `PELM`: 経路要素。経路内の個々の機器（コンポーネント）を定義します。
- **熱媒の種類と状態**: `Mpath.Fluid`は、
  経路を流れる熱媒の種類（空気の温度、空気の湿度、水など）を定義します。
  これにより、熱媒の種類に応じた適切な熱計算や、
  顕熱と潜熱の分離処理を正確に行うことができます。
- **流量制御と運転モード**: `Plist.Go`（流量の計算に使用される係数）や`Plist.Rate`（流量分配比）は、
  経路内の流量制御をモデル化します。
  また、`Plist.Batch`（バッチ運転）や`stank.Batchcon`（蓄熱槽のバッチ運転制御）は、
  蓄熱槽などの特殊な運転モードをモデル化します。
- **システム統合とエネルギーフロー**: この関数は、
  個々の機器の性能だけでなく、
  システム全体での熱媒の流れやエネルギーフローを追跡するために用いられます。
  これにより、システム全体のエネルギー効率を評価し、
  エネルギーの無駄を特定し、
  省エネルギー対策の効果を定量的に把握できます。
- **エラーハンドリング**: 入力データの不備や、
  コンポーネントが見つからない場合などにエラーメッセージを出力し、
  プログラムを終了します。

この関数は、建物のエネルギーシミュレーションにおいて、
熱搬送システムや空調システムの複雑な構成と運転を正確にモデル化し、
エネルギー消費量予測、省エネルギー対策の検討、
および最適な設備システム設計を行うための重要な役割を果たします。
*/
func Pathdata(
	f *EeTokens,
	Simc *SIMCONTL,
	Wd *WDAT,
	Compnt []*COMPNT,
	Schdl *SCHDL,
	M *[]*MPATH,
	Plst *[]*PLIST,
	Plm *[]*PELM,
	Eqsys *EQSYS,
	Elout *[]*ELOUT,
	Elin *[]*ELIN,
) {
	var C *COMPNT
	var stank *STANK
	var Qmeas *QMEAS
	var s, ss string
	var elm string
	var id int

	var i, j, m, ncv, iswc int

	errkey := "Pathdata"
	if DEBUG {
		fmt.Printf("\n")
		for i := range Compnt {

			C = Compnt[i]
			fmt.Printf("name=%s Nin=%d  Nout=%d\n", C.Name, C.Nin, C.Nout)
		}
	}

	// 最外ループ: MPATHの読み込み 開始
	for f.IsEnd() == false {
		ss = f.GetToken()
		if ss[0] == '*' {
			break
		}
		if ss[0] == '\n' {
			continue
		}

		if DEBUG {
			fmt.Printf("eepathdat.c  ss=%s\n", ss)
		}

		Mpath := NewMPATH()

		// if *Nplst > 0 {
		// 	// 経路系統 -> 末端経路 への参照ポインタの確保
		// 	Mpath.Pl = make([]*PLIST, *Nplst)
		// 	for i := 0; i < *Nplst; i++ {
		// 		Mpath.Pl[i] = nil
		// 	}
		// }

		// if *Nplst > 0 {
		// 	// 流量連立方程式を解くときに使用する分岐・合流機器 への参照ポインタの確保
		// 	Mpath.Cbcmp = make([]*COMPNT, *Nplst*2)
		// 	for i := 0; i < *Nplst*2; i++ {
		// 		Mpath.Cbcmp[i] = nil
		// 	}
		// }

		// SYSPTHにおける';'で区切られる経路の基本設定
		// 入力例: `LDAC -sys A -f A`
		Mpath.Name = ss                 // 経路名称
		Mpath.NGv = 0                   // ガス導管数
		Mpath.NGv2 = 0                  // 開口率が2%未満のガス導管数
		Mpath.Ncv = 0                   // 制御弁数
		Mpath.Plist = make([]*PLIST, 0) // 末端経路
		//Pelmpre = nil

		// PLIST読み込み用ループ 開始
		iPlist := 0 //デフォルト名の設定用カウンタ
		for f.IsEnd() == false {
			s = f.GetToken()
			if s[0] == ';' {
				break
			}

			if DEBUG {
				fmt.Printf("eepathdat.c  s=%s\n", s)
			}

			if s[0] == '-' {
				//
				// MPATH用の属性を読み込む -sys or -f
				//
				ss = f.GetToken()

				if s[1:] == "sys" {
					// システムの分類（A：空調・暖房システム、D：給湯システム）
					Mpath.Sys = ss[0]
				} else if s[1:] == "f" {
					// 循環、通過流体の種類（水系統、空気系統の別）
					//（W：水系統、A：空気系統で温・湿度とも計算、a：空気系統で温度のみ計算。
					Mpath.Fluid = FliudType(ss[0])
				} else {
					Errprint(1, errkey, s)
				}
			} else if s == ">" {
				//
				// PLIST用の作成　">"から">"の塊を読みとる
				//
				Plist := NewPLIST()
				Plist.Plistname = fmt.Sprintf("Path%d", iPlist) // 末端経路名
				Plist.Pelm = make([]*PELM, 0, 10)               // 機器のリスト
				iPlist++

				// PELM読み込み用ループ 開始
				for f.IsEnd() == false {
					s = f.GetToken()

					// ">" が来ると終了
					if s == ">" {
						break
					}

					if DEBUG {
						fmt.Printf("eepathdat.c  s=%s\n", s)
					}

					if s[0] == '(' {
						//
						// 流量の計算に使用される係数を読み取る - `(<flwvol>)` のような文字列で定義される
						//

						// "(value)" のような文字列から vaue を取り出す正規表現
						re := regexp.MustCompile(`\((.*?)\)`)
						matches := re.FindStringSubmatch(s)

						if len(matches) >= 2 {
							ss := matches[1]

							if Go, err := readFloat(ss); err == nil {
								Plist.Go = CreateConstantValuePointer(Go)
								if DEBUG {
									fmt.Printf("Go=%f\n", *Plist.Go)
								}
							} else {
								if DEBUG {
									fmt.Printf("s=%s ss=%s\n", s, ss)
								}

								if j, err = idsch(ss, Schdl.Sch, ""); err == nil {
									Plist.Go = &Schdl.Val[j]
								} else {
									Plist.Go = envptr(ss, Simc, Compnt, Wd, nil)
								}

								if DEBUG {
									fmt.Printf("Go=%f\n", *Plist.Go)
								}
							}
						}
					} else if s[0] == '[' {
						//
						// 流量比率を読み取る - `[<flwrate>]` のような文字列で定義される
						//

						// 流量比率設定フラグのセット
						Mpath.Rate = true

						var Go float64
						i, err := fmt.Scanf(s[1:], "%f", &Go)
						if err != nil {
							panic(err)
						}
						if i == 1 {
							// 指定されていたのが数値なので、流量比率に入力された固定値を設定する
							Plist.Rate = CreateConstantValuePointer(Go)
							if DEBUG {
								fmt.Printf("rate=%f\n", *Plist.Rate)
							}
						} else {
							// 指定されたいたのは数値ではなかったので、設定値スケジュール名である。
							// 、) 文字が見つかるまでの全ての文字を読み取ります.
							_, err := fmt.Sscanf(s[1:], "%[^]]", &ss)
							if err != nil {
								panic(err)
							}

							if DEBUG {
								fmt.Printf("s=%s ss=%s\n", s, ss)
							}

							// 設定値スケジュールの検索
							if j, err := idsch(ss, Schdl.Sch, ""); err == nil {
								// 設定値スケジュールが見つかったので、スケジュール設定値へのポインタを指定
								Plist.Rate = &Schdl.Val[j]
							} else {
								Plist.Rate = envptr(ss, Simc, Compnt, Wd, nil)
							}

							if DEBUG {
								fmt.Printf("rate=%f\n", *Plist.Rate)
							}
						}
					} else if strings.HasPrefix(s, "name=") {
						// 末端経路名称の指定
						// デフォルトでは、 "Path<No.>"のような名前で指定されるので、これを上書きする
						_, err := fmt.Sscanf(s, "%*[^=]=%s", &ss)
						if err != nil {
							panic(err)
						}
						Plist.Plistname = ss
					} else {

						// 要素名を読み取る 4パターンある. 読み取った要素名は elm に格納する。
						//
						var stv string = "" // stv=蓄熱槽のスケジュール
						var ci ELIOType = ELIO_SPACE
						var co ELIOType = ELIO_SPACE

						if idx := strings.IndexRune(s, '/'); idx >= 0 {
							// ex: `xxx/LD`
							s = s[idx+1:] // 蓄熱槽の機器名
							stv = s[:idx] // stv=蓄熱槽のスケジュール
						}
						if idx := strings.IndexRune(s, ':'); idx >= 0 {
							// ex: `LD:xxx`
							Plist.Name = s
							s = s[:idx]
							elm = s
						} else {
							if idx := strings.IndexRune(s, '['); idx >= 0 {
								// ex: `LD[r]`
								// rは経路識別子
								// 2流体式熱交換器や蓄熱槽のように構成要素が複数の経路にまたがる場合の書式
								co, ci = ELIOType(s[idx+1]), ELIOType(s[idx+1])
								s = s[:idx]
								elm = s
							} else {
								// ex: `LD`
								elm = s
								co = ELIOType(0)
								ci = ELIOType(0)
							}
						}

						// SYSCMPで定義したシステム要素名から elm を探す
						err := 1
						_, cmp, er := FindComponent(elm, Compnt)
						if er != nil {
							panic(er)
						}
						err = 0

						// 機器の種類に応じた処理分け
						//
						if cmp.Eqptype == FLIN_TYPE && len(Plist.Pelm) == 0 {
							// 経路の先頭が流入境界条件である場合、経路の種類は流入境界条件である。
							Plist.Type = IN_LPTP
						} else if cmp.Eqptype == VALV_TYPE || cmp.Eqptype == TVALV_TYPE {
							// バルブが見つかったので、後の処理のために経路と機器に相互参照を付与する
							Plist.Nvalv++
							Plist.Valv = cmp.Eqp.(*VALV)
							Plist.Valv.Plist = Plist // 逆参照
						} else if cmp.Eqptype == OMVAV_TYPE {
							// OMバルブが見つかったので、後の処理のために経路と機器に相互参照を付与する
							// Satoh OMVAV 2010/12/16
							Plist.NOMVAV++
							Plist.OMvav = cmp.Eqp.(*OMVAV)
							Plist.OMvav.Plist = Plist // 逆参照
						} else if cmp.Eqptype == VAV_TYPE || cmp.Eqptype == VWV_TYPE {
							/*---- Satoh Debug VAV  2000/12/6 ----*/
							// VAVユニットが見つかったが、見つかった数だけ記録する。
							Plist.Nvav++
						} else if cmp.Eqptype == QMEAS_TYPE {
							/*---- Satoh Debug QMEAS  2003/6/2 ----*/
							// カロリーメータが見つかったので、後の処理のために経路と機器に相互参照を付与する
							// ただし、単一のカロリーメータは、3種類の値を同時に参照できるため注意する。
							Qmeas = cmp.Eqp.(*QMEAS)
							if co == ELIO_G {
								Qmeas.G = &Plist.G
								Qmeas.PlistG = Plist
								Qmeas.Fluid = Mpath.Fluid
							} else if co == ELIO_H {
								Qmeas.PlistTh = Plist
								Qmeas.Nelmh = id
							} else if co == ELIO_C {
								Qmeas.PlistTc = Plist
								Qmeas.Nelmc = id
							} else {
								// NOTE: オリジナルコードではこのelseはない。念のため導入してみた。
								panic(co)
							}
						} else if cmp.Eqptype == STANK_TYPE {
							// 蓄熱槽が見つかった
							if stv != "" {
								Plist.Batch = true
								stank = cmp.Eqp.(*STANK)
								for i := 0; i < stank.Nin; i++ {
									if stank.Pthcon[i] == co {
										var err error
										if iswc, err = idscw(stv, Schdl.Scw, ""); err == nil {
											stank.Batchcon[i] = Schdl.Isw[iswc]
										}
									}
								}
							}
						}

						// バルブ、カロリーメータは末端経路ごとに1つまでのようだ
						// その他の要素は複数存在しても良い。
						if cmp.Eqptype != VALV_TYPE && cmp.Eqptype != TVALV_TYPE &&
							cmp.Eqptype != QMEAS_TYPE && cmp.Eqptype != OMVAV_TYPE {

							Pelm := NewPELM()
							Pelm.Out = nil
							Pelm.Cmp = cmp
							Pelm.Ci = ci
							Pelm.Co = co
							//Pelmpre = Pelm

							Plist.Pelm = append(Plist.Pelm, Pelm)
							*Plm = append(*Plm, Pelm)
						}

						if cmp.Eqptype != VALV_TYPE && cmp.Eqptype != TVALV_TYPE &&
							cmp.Eqptype != QMEAS_TYPE && cmp.Eqptype != DIVERG_TYPE &&

							cmp.Eqptype != CONVRG_TYPE && cmp.Eqptype != DIVGAIR_TYPE &&
							cmp.Eqptype != CVRGAIR_TYPE && cmp.Eqptype != OMVAV_TYPE {
							id++
						}

						Errprint(err, errkey, elm)

						if DEBUG {
							fmt.Printf("<<Pathdata>> Mp=%s  elm=%s Npelm=%d\n", Mpath.Name, elm, len(*Plm))
						}
					}
				}
				// PELM読み込み用ループ 終了

				*Plst = append(*Plst, Plist)
				Mpath.Plist = append(Mpath.Plist, Plist)

				Plist.Mpath = Mpath // 子→親の逆参照
				//Pelmpre = nil
				id = 0
			} else {
				Errprint(1, errkey, s)
			}
		}
		// PLIST読み込み用ループ 終了

		if DEBUG {
			// 読み取ったMPATHの順番と流体種別を表示
			fmt.Printf("<<Pathdata>>  Mpath=%d fliud=%c\n", len(*M), Mpath.Fluid)
		}

		// 流体種別が空気の場合: 空気系統用の絶対湿度経路へのコピーを行う必要がある
		if Mpath.Fluid == AIR_FLD {
			if DEBUG {
				fmt.Printf("<<Pathdata  a>> Mp=%s  Npelm=%d\n", Mpath.Name, len(*Plm))
			}

			// 空気温度用の経路を追加
			*M = append(*M, Mpath)

			// if *Nplst > 0 {
			// 	Mpath.Pl = make([]*PLIST, *Nplst)
			// 	for k := 0; k < *Nplst; k++ {
			// 		Mpath.Pl[k] = nil
			// 	}
			// }

			// if *Nplst > 0 {
			// 	Mpath.Cbcmp = make([]*COMPNT, *Nplst*2)
			// 	for k := 0; k < *Nplst*2; k++ {
			// 		Mpath.Cbcmp[k] = nil
			// 	}
			// }

			// 空気系統用の絶対湿度経路へのコピー
			Mpath_x := plistcpy(Mpath, Plm, Plst, Compnt)
			*M = append(*M, Mpath_x)
		} else {
			*M = append(*M, Mpath)
		}
	}
	// 最外ループ: MPATHの読み込み 完了

	if DEBUG {
		if len(*M) > 0 {
			plistprint(*M, *Plm, *Elout, *Elin)
		}

		fmt.Printf("SYSPTH  Data Read end\n")
		fmt.Printf("Nmpath=%d\n", len(*M))
	}

	/* ============================================================================ */

	// すべてのMPATHについてのループ
	for i, Mpath := range *M {
		if DEBUG {
			fmt.Printf("1----- MAX=%d  i=%d\n", len(*M), i)
		}

		ncv = 0

		//
		// --- pelmci or pelmco を呼び出す----
		//
		for j, Plist := range Mpath.Plist {

			if DEBUG {
				fmt.Printf("eepath.c  Mpath->Nlpath=%d\n", len(Mpath.Plist))
				fmt.Printf("<<Pathdata>>  i=%d Mpath=%d  j=%d Plist=%d\n", i, i, j, j)
			}

			for m, Pelm := range Plist.Pelm {

				if DEBUG {
					fmt.Printf("<<Pathdata>>  m=%d  pelm=%d  %s\n", m, m, Pelm.Cmp.Name)
					fmt.Printf("MAX=%d  m=%d\n", len(Plist.Pelm), m)
				}

				//
				// --- システム要素入出力端割当 ---
				//

				idci := true // システム要素入力端割当を行うか？
				idco := true // システム要素出力端割当を行うか？
				etyp := Pelm.Cmp.Eqptype

				if m == 0 && etyp == FLIN_TYPE {
					// 末端経路の先頭要素が*流入境界条件*である場合
					idci = false // 流入境界条件
				}

				if m == 0 && (etyp == CONVRG_TYPE || etyp == DIVERG_TYPE) {
					// 末端経路の先頭要素が*水*の合流または分岐である場合
					idci = false
				} else if m == 0 && (etyp == CVRGAIR_TYPE || etyp == DIVGAIR_TYPE) {
					// 末端経路の先頭要素が*空気*の合流または分岐である場合
					idci = false
				}

				if m == len(Plist.Pelm)-1 && (etyp == CONVRG_TYPE || etyp == DIVERG_TYPE) {
					// 末端経路の最後尾要素が*水*の合流または分岐である場合
					idco = false
				} else if m == len(Plist.Pelm)-1 && (etyp == CVRGAIR_TYPE || etyp == DIVGAIR_TYPE) {
					// 末端経路の最後尾要素が*空気*の合流または分岐である場合
					idco = false
				}

				if idci {
					// システム要素入力端割当
					pelmci(Mpath.Fluid, Pelm, errkey)
					Pelm.In.Lpath = Plist
				}

				if idco {
					// システム要素出力端割当
					pelmco(Mpath.Fluid, Pelm, errkey)

					Pelm.Out.Lpath = Plist
					Pelm.Out.Fluid = Mpath.Fluid
				}
			}
		}

		if DEBUG {
			plistprint((*M)[i:i+1], *Plm, *Elout, *Elin)
			fmt.Printf("i=%d\n", i)
		}

		//
		// --- 貫流経路か循環経路かの判定 + 要素間の接続 ---
		//
		if len(Mpath.Plist) == 1 {
			//
			// 末端経路の数が1の場合
			//

			Plist := Mpath.Plist[0]

			if DEBUG {
				fmt.Printf("<<Pathdata>>   Plist->type=%c\n", Plist.Type)
			}

			if Plist.Type == IN_LPTP {
				Mpath.Type = THR_PTYP

				if DEBUG {
					fmt.Printf("<<Pathdata>>   Mpath->type=%c\n", Mpath.Type)
				}
			} else {
				Mpath.Type = CIR_PTYP
				Plist.Type = CIR_PTYP
				Plist.Pelm[0].In.Upo = Plist.Pelm[len(Plist.Pelm)-1].Out
			}

			if DEBUG {
				fmt.Printf("<<Pathdata>> test end\n")
			}

			// 2番目以降の要素について
			for m = 1; m < len(Plist.Pelm); m++ {
				Pelm := Plist.Pelm[m]
				PelmPrev := Plist.Pelm[m-1]

				// 要素間の接続: 1つ前の要素の出力への参照を設定
				Pelm.In.Upo = PelmPrev.Out
			}
		} else {
			//
			// 末端経路の数が2以上の場合
			//

			Mpath.Type = BRC_PTYP

			if DEBUG {
				fmt.Printf("<<Pathdata>> Mpath i=%d  type=%c\n", i, Mpath.Type)
			}

			for j, Plist := range Mpath.Plist {
				// 1. 先頭要素による判定
				//
				Pelm_0 := Plist.Pelm[0]
				etyp_0 := Pelm_0.Cmp.Eqptype

				if DEBUG {
					fmt.Printf("<<Pathdata>> Plist j=%d name=%s eqptype=%s\n", j, Pelm_0.Cmp.Name, etyp_0)
				}

				if etyp_0 == DIVERG_TYPE || etyp_0 == DIVGAIR_TYPE {
					// 先頭要素が水または空気の*分岐*である場合
					Plist.Type = DIVERG_LPTP
				}

				if etyp_0 == CONVRG_TYPE || etyp_0 == CVRGAIR_TYPE {
					// 先頭要素が水または空気の*合流*である場合
					Plist.Type = CONVRG_LPTP
					ncv++
				}

				// 2. 最後尾要素による判定
				//
				etyp_fin := Plist.Pelm[len(Plist.Pelm)-1].Cmp.Eqptype
				if etyp_fin != DIVERG_TYPE && etyp_fin != CONVRG_TYPE &&
					etyp_fin != DIVGAIR_TYPE && etyp_fin != CVRGAIR_TYPE {
					// 最後尾要素が水または空気の分岐または合流である場合
					Plist.Type = OUT_LPTP
				}

				// 2番目以降の要素について
				for m := 1; m < len(Plist.Pelm); m++ {
					Pelm := Plist.Pelm[m]
					PelmPrev := Plist.Pelm[m-1]

					// 要素間の接続: 1つ前の要素の出力への参照を設定
					Pelm.In.Upo = PelmPrev.Out
				}

				if DEBUG {
					fmt.Printf("<<Pathdata>> Plist MAX=%d  j=%d  type=%c\n", len(Mpath.Plist), j, Plist.Type)
				}
			}
		}
		Mpath.Ncv = ncv

		if DEBUG {
			fmt.Printf("2----- MAX=%d  i=%d\n", len(*M), i)
		}
	}

	if DEBUG {
		if len(*M) > 0 {
			mpi := *M
			plistprint(mpi, *Plm, *Elout, *Elin)
		}
	}

	// バルブがモニターするPlistの検索
	Valvinit(Eqsys.Valv, *M)

	// 未知流量等の構造解析
	pflowstrct(*M)

	if DEBUG {
		if len(*M) > 0 {
			plistprint(*M, *Plm, *Elout, *Elin)
		}
	}

	if DEBUG {
		fmt.Printf("\n")
		for i = range Compnt {
			C := Compnt[i]
			fmt.Printf("name=%s Nin=%d  Nout=%d\n", C.Name, C.Nin, C.Nout)
		}
	}
}

/*
Mpathcount (Main Path Count)

この関数は、入力ファイルから読み込まれるシステム経路（`MPATH`）の総数をカウントします。

建築環境工学的な観点:
- **システム経路の総数把握**: 建物のエネルギーシミュレーションでは、
  熱媒が流れる複数のシステム経路（空調系統、給湯系統など）をモデル化します。
  この関数は、入力ファイル内の区切り文字（`;`）を数えることで、
  システム経路の総数を把握します。
- **シミュレーションの準備**: システム経路の総数を事前に把握することで、
  シミュレーションに必要なメモリ領域を確保したり、
  計算ループの回数を決定したりすることができます。

この関数は、建物のエネルギーシミュレーションにおいて、
システム構成を正確にモデル化し、
エネルギー消費量予測、省エネルギー対策の検討、
および最適な設備システム設計を行うための重要な役割を果たします。
*/
func Mpathcount(fi *EeTokens, Pl *int) int {
	var N int
	var ad int
	var s string

	ad = fi.GetPos()
	*Pl = 0

	for fi.IsEnd() == false {
		s = fi.GetToken()

		if s == "*" {
			break
		}

		if s == ";" {
			N++
		}

		if s == ">" {
			*Pl++
		}
	}

	*Pl /= 2

	fi.RestorePos(ad)

	return N
}

/*
Plcount (Path List Count)

この関数は、入力ファイルから読み込まれる末端経路（`PLIST`）の総数をカウントします。

建築環境工学的な観点:
- **末端経路の総数把握**: 建物のエネルギーシミュレーションでは、
  システム経路を構成する複数の末端経路をモデル化します。
  この関数は、入力ファイル内の区切り文字（`;`）と（`>`）を数えることで、
  末端経路の総数を把握します。
- **シミュレーションの準備**: 末端経路の総数を事前に把握することで、
  シミュレーションに必要なメモリ領域を確保したり、
  計算ループの回数を決定したりすることができます。

この関数は、建物のエネルギーシミュレーションにおいて、
システム構成を正確にモデル化し、
エネルギー消費量予測、省エネルギー対策の検討、
および最適な設備システム設計を行うための重要な役割を果たします。
*/
func Plcount(fi *EeTokens, N []int) {
	i := 0
	M := 0
	ad := fi.GetPos()

	for fi.IsEnd() == false {
		s := fi.GetToken()

		if s == "*" {
			break
		}

		if s == ";" {
			N[M] = i
			M++
			i = 0
		}

		if s == ">" {
			i++
			fi.GetToken() // skip next token
		}
	}

	// Print the contents of the N slice for debugging purposes
	// for i := 0; i < len(*N); i++ {
	// 	fmt.Printf("i=%d pl=%d\n", i, (*N)[i])
	// }

	fi.RestorePos(ad)
}

/*
Pelmcount (Path Element Count)

この関数は、入力ファイルから読み込まれる経路要素（`PELM`）の総数をカウントします。
特に、流体が空気であって温・湿度の両方を計算する場合には、
経路を複製して1経路を2経路分としてカウントします。

建築環境工学的な観点:
- **経路要素の総数把握**: 建物のエネルギーシミュレーションでは、
  末端経路を構成する複数の経路要素（機器）をモデル化します。
  この関数は、入力ファイル内の経路要素の総数を把握します。
- **熱湿気同時交換の考慮**: 流体が空気であって温・湿度の両方を計算する場合には、
  温度経路と湿度経路を別々にモデル化する必要があるため、
  経路要素の数を2倍としてカウントします。
  これにより、熱湿気同時交換を正確にモデル化できます。
- **シミュレーションの準備**: 経路要素の総数を事前に把握することで、
  シミュレーションに必要なメモリ領域を確保したり、
  計算ループの回数を決定したりすることができます。

この関数は、建物のエネルギーシミュレーションにおいて、
システム構成を正確にモデル化し、
エネルギー消費量予測、省エネルギー対策の検討、
および最適な設備システム設計を行うための重要な役割を果たします。
*/
func Pelmcount(fi *EeTokens) int {
	ad := fi.GetPos()
	i := 1
	N := 0

	for fi.IsEnd() == false {
		s := fi.GetToken()
		i = 1

		if s == "*" {
			break
		}

		for fi.IsEnd() == false {
			s = fi.GetToken()

			if s == ";" {
				break
			}

			if s == "-f" {
				// 循環、通過流体の種類（水系統、空気系統の別）
				t := fi.GetToken()

				if t == "W" || t == "a" {
					// W：水系統、a：空気系統で温度のみ計算
					i = 1
				} else {
					// A：空気系統で温・湿度とも計算
					i = 2 // 温・湿度に分けるので2経路としてカウントする
				}
			}

			if s == "-sys" {
				// `A` or `D`. システムの分類（A：空調・暖房システム、D：給湯システム）
				fi.GetToken()
			}

			if s != ">" && s[:1] != "(" && s[:1] != "-" && s[:1] != ";" {
				N += i
			}
		}
	}

	fi.RestorePos(ad)
	return N
}

/***********************************************************************/

/*
Elcount (Element Count for Inputs and Outputs)

この関数は、全てのコンポーネント（機器）の入出力要素（`ELIN`, `ELOUT`）の総数をカウントします。

建築環境工学的な観点:
- **入出力要素の総数把握**: 建物のエネルギーシミュレーションでは、
  各機器の入出力要素をモデル化し、
  機器間の接続や熱媒の流れを追跡します。
  この関数は、全ての機器の入出力要素の総数を把握します。
- **特殊な機器の考慮**: `HCLOADW_TYPE`（冷・温水コイル）や`THEX_TYPE`（全熱交換器）のように、
  入出力要素の数が特殊な機器の場合、
  その特性を考慮してカウントを行います。
- **シミュレーションの準備**: 入出力要素の総数を事前に把握することで、
  シミュレーションに必要なメモリ領域を確保したり、
  計算ループの回数を決定したりすることができます。

この関数は、建物のエネルギーシミュレーションにおいて、
システム構成を正確にモデル化し、
エネルギー消費量予測、省エネルギー対策の検討、
および最適な設備システム設計を行うための重要な役割を果たします。
*/
func Elcount(C []*COMPNT) (int, int) {
	var Nelout, Nelin int = 0, 0

	for i := range C {
		e := C[i].Eqptype
		Nelout += C[i].Nout
		Nelin += C[i].Nin

		if e == HCLOADW_TYPE {
			Nelin += 8
		} else if e == THEX_TYPE {
			Nelout += 4
			Nelin += 14
		}
	}

	Nelout *= 4
	Nelin *= 4

	return Nelout, Nelin
}

/*
FindComponent (Find Component by Name)

この関数は、与えられた名称（`name`）に基づいて、
コンポーネントのリスト（`Compnt`）から該当するコンポーネントを検索します。

建築環境工学的な観点:
- **コンポーネントの参照**: 建物のエネルギーシミュレーションでは、
  様々なコンポーネントが相互に接続され、熱や空気、水などをやり取りします。
  この関数は、あるコンポーネントが別のコンポーネントを参照する際に、
  その名称に基づいて対象のコンポーネントを効率的に見つけ出すために用いられます。
- **システム構成の動的な構築**: シミュレーションモデルを構築する際、
  コンポーネント間の接続関係は、入力ファイルから読み込まれる名称に基づいて動的に設定されます。
  この関数は、その動的な接続を可能にするための基本的な検索機能を提供します。

この関数は、建物のエネルギーシミュレーションにおいて、
コンポーネント間の接続関係を確立し、
システム全体の熱・空気・水の流れをモデル化するための重要な役割を果たします。
*/
func FindComponent(name string, Compnt []*COMPNT) (int, *COMPNT, error) {
	for i := range Compnt {
		cmp := Compnt[i]
		if cmp.Name == name {
			return i, cmp, nil
		}
	}
	return -1, nil, errors.New(fmt.Sprintf("Component [%s] not found", name))
}

/*
CreateConstantValuePointer (Create Pointer to Constant Value)

この関数は、与えられた定数値（`constValue`）へのポインターを作成します。

建築環境工学的な観点:
- **定数値の参照**: シミュレーションモデルでは、
  定数として扱われる値（例: 固定流量、固定温度）が多数存在します。
  この関数は、これらの定数値をポインターとして扱うことで、
  他の変数と同じように参照できるようにします。
- **柔軟なデータ管理**: 定数値をポインターとして扱うことで、
  後からその値を変更する必要が生じた場合でも、
  ポインターが指す先の値を変更するだけで済み、
  コードの修正を最小限に抑えることができます。

この関数は、建物のエネルギーシミュレーションにおいて、
データの管理を柔軟に行い、
シミュレーションモデルの構築を支援するための重要な役割を果たします。
*/
func CreateConstantValuePointer(constValue float64) *float64 {
	val := new(float64)
	*val = constValue
	return val
}

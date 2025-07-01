package eeslism

import (
	"fmt"
	"math"
	"os"
)

/*
Pflow (Path Flow Calculation)

この関数は、建物のエネルギーシミュレーションにおける熱媒（空気、水など）の
流れる経路（パス）の流量を計算します。
これには、流量が既知の末端経路、バルブやVAVユニットによって制御される経路、
および分岐・合流点における流量の連立方程式の解法が含まれます。

建築環境工学的な観点:
- **流量計算の重要性**: 熱搬送システムや空調システムにおける流量は、
  熱供給量や熱回収量、そして機器のエネルギー消費量を決定する上で不可欠です。
  正確な流量計算は、システム全体のエネルギー効率を評価し、
  熱負荷への対応能力を予測するために重要です。
- **流量が既知の末端経路**: `Plist.Go != nil && Plist.Nvalv == 0` の条件は、
  流量が固定値として与えられている末端経路を示します。
  これは、例えば、定流量ポンプやファンによって駆動される経路に適用されます。
- **バルブ・VAVユニットによる流量制御**: `Plist.Nvalv > 0` や `Plist.Nvav > 0` の条件は、
  バルブやVAVユニットによって流量が制御される末端経路を示します。
  - **二方弁**: `Plist.Valv.X * *Plist.Go` のように、
    バルブ開度（`vc.X`）に応じて流量を調整します。
  - **三方弁**: `vc.X * *vc.MGo` や `(1.0 - vcmb.X) * *vc.MGo` のように、
    バルブ開度と連動する別のバルブの開度に応じて流量を調整します。
  - **VAVユニット**: `vav.Cat.Gmax` や `vav.G` を用いて、
    VAVユニットの最大流量や現在の流量を考慮します。
  - **OMVAVユニット**: `OMflowcalc`関数を呼び出して、
    外気処理VAVユニットの流量を計算します。
  これらの制御は、室の熱負荷変動に応じて熱媒の供給量を調整し、
  省エネルギーと快適性を両立させるために重要です。
- **分岐・合流点における流量の連立方程式**: `Mpath.NGv`（ガス導管数）が`0`より大きい場合、
  分岐・合流点における流量の連立方程式を解きます。
  これは、各分岐・合流点での質量保存則に基づいて流量を決定するもので、
  複雑な配管・ダクトネットワークにおける流量分布を正確に計算するために不可欠です。
- **流量の妥当性チェック**: 計算された流量が負の値になった場合、
  `Plist.G = 0.0` と設定し、経路を停止させます。
  これは、物理的に不可能な流量を防ぎ、シミュレーションの安定性を確保するために重要です。

この関数は、建物のエネルギーシミュレーションにおいて、
熱搬送システムや空調システムの流量を正確にモデル化し、
エネルギー消費量予測、省エネルギー対策の検討、
および最適な設備システム設計を行うための重要な役割を果たします。
*/
func Pflow(_Mpath []*MPATH, Wd *WDAT) {
	var i, j, n, NG int
	//var mpi *MPATH
	var pl *PLIST
	var eli *ELIN
	var elo *ELOUT
	var cmp *COMPNT
	var vc, vcmb *VALV
	var Go float64
	var G float64
	var Err, s string
	/*---- Satoh Debug VAV  2000/12/6 ----*/
	var vav *VAV
	var G0 float64

	if len(_Mpath) > 0 {
		for m, Mpath := range _Mpath {

			if DEBUG {
				fmt.Printf("m=%d mMAX=%d name=%s\n", m, len(_Mpath), Mpath.Name)
			}

			// 流量が既知の末端流量の初期化
			for _, Plist := range Mpath.Plist {
				if Plist.Go != nil && Plist.Nvalv == 0 {
					Plist.G = *Plist.Go
				}
			}

			for i, Plist := range Mpath.Plist {
				Plist.G = 0.0

				if DEBUG {
					fmt.Printf("i=%d iMAX=%d name=%s\n", i, len(Mpath.Plist), get_string_or_null(Plist.Name))
				}

				// 流量が既知の末端経路
				if Plist.Go != nil && Plist.Nvalv == 0 {
					Plist.G = *Plist.Go
				} else if Plist.Go != nil && Plist.Nvalv > 0 ||
					Plist.NOMVAV > 0 ||
					(Plist.Go == nil && Plist.Nvalv > 0 && Plist.UnknownFlow == 1) {
					if Plist.Go != nil && Plist.Valv != nil &&
						Plist.Valv.Cmp.Eqptype == VALV_TYPE {
						// 二方弁の計算
						Plist.G = *Plist.Go
						vc = Plist.Valv
						if vc == nil || vc.Org == 'y' {
							if vc.X < 0.0 {
								s = fmt.Sprintf("%s のバルブ開度 %f が不正です。", vc.Name, vc.X)
								Eprint("<Pflow>", s)
							}
							Plist.G = vc.X * *Plist.Go
						} else {
							vcmb = vc.Cmb.Eqp.(*VALV)
							Plist.G = (1.0 - vcmb.X) * *Plist.Go
						}
					} else if Plist.Valv != nil && Plist.Valv.MGo != nil &&
						*Plist.Valv.MGo > 0.0 && Plist.Control != OFF_SW {
						// 三方弁の計算

						vc = Plist.Valv
						vcmb = vc.Cmb.Eqp.(*VALV)

						if vc.Org == 'y' {
							Plist.G = vc.X * *vc.MGo
						} else {
							Plist.G = (1.0 - vcmb.X) * *vc.MGo
						}

						if Plist.G > 0. {
							Plist.Control = ON_SW
						}
					} else if Plist.Valv != nil && Plist.Valv.MGo != nil && *Plist.Valv.MGo <= 0.0 {
						Plist.G = 0.0
					} else if Plist.Valv != nil && Plist.Valv.Count > 0 {
						Plist.G = Plist.Gcalc
					} else if Plist.NOMVAV > 0 {
						Plist.G = OMflowcalc(Plist.OMvav, Wd)
					}

					if Plist.G <= 0.0 {
						Plist.lpathscdd(OFF_SW)
					}

					if Plist.G > 0. {
						Plist.Control = ON_SW
					}
				} else if Plist.Nvav > 0 {
					/*---- Satoh Debug VAV  2000/12/6 ----*/

					/* VAVユニット時の流量 */

					G = -999.0
					for _, Pelm := range Plist.Pelm {
						if Pelm.Cmp.Eqptype == VAV_TYPE ||
							Pelm.Cmp.Eqptype == VWV_TYPE {
							vav = Pelm.Cmp.Eqp.(*VAV)

							if vav.Count == 0 {
								G = math.Max(G, vav.Cat.Gmax)
							} else {
								G = math.Max(G, vav.G)
							}
						}
					}
					Plist.G = G
				} else if Plist.Rate != nil {
					Plist.G = *Mpath.G0 * *Plist.Rate
				} else if !Plist.Batch {
					if Plist.Go != nil {
						Go = *Plist.Go
					} else {
						Go = 0.0
					}

					if Plist.Pelm != nil {
						Err = fmt.Sprintf("Mpath=%s  lpath=%d  elm=%s  Go=%f\n", Mpath.Name, 0, Plist.Pelm[0].Cmp.Name, Go)
					}
				}
			}

			NG = Mpath.NGv

			X := make([]float64, NG)
			Y := make([]float64, NG)
			A := make([]float64, NG*NG)

			for i = 0; i < NG; i++ {
				if DEBUG {
					fmt.Printf("i=%d iMAX=%d\n", i, NG)
				}

				cmp = Mpath.Cbcmp[i]

				if DEBUG {
					fmt.Printf("<Pflow> Name=%s\n", cmp.Name)
				}

				for j = 0; j < cmp.Nin; j++ {
					eli = cmp.Elins[j]

					if DEBUG {
						fmt.Printf("j=%d jMAX=%d\n", j, cmp.Nin)
					}

					if eli.Lpath.Go != nil ||
						eli.Lpath.Nvav != 0 ||
						eli.Lpath.Nvalv != 0 ||
						eli.Lpath.Rate != nil ||
						eli.Lpath.NOMVAV != 0 {
						Y[i] -= eli.Lpath.G
					} else {
						n = eli.Lpath.N

						if n < 0 || n >= NG {
							Err = fmt.Sprintf("n=%d", n)
							Eprint("<Pflow>", Err)
							os.Exit(EXIT_PFLOW)
						}

						A[i*NG+n] = 1.0
					}
				}

				////////

				for j = 0; j < cmp.Nout; j++ {
					elo = cmp.Elouts[j]

					if elo.Lpath.Go != nil ||
						elo.Lpath.Nvav != 0 ||
						elo.Lpath.Nvalv != 0 ||
						elo.Lpath.Rate != nil {
						Y[i] += elo.Lpath.G
					} else {
						n = elo.Lpath.N

						if n < 0 || n >= NG {
							Err = fmt.Sprintf(Err, "n=%d", n)
							Eprint("<Pflow>", Err)
							os.Exit(EXIT_PFLOW)
						}

						A[i*NG+n] = -1.0
					}
				}
			}

			if NG > 0 {

				if DEBUG {
					for i = 0; i < NG; i++ {
						fmt.Printf("%s\t", Mpath.Cbcmp[i].Name)

						for j = 0; j < NG; j++ {
							fmt.Printf("%6.1f", A[i*NG+j])
						}

						fmt.Printf("\t%.5f\n", Y[i])
					}
				}

				if dayprn && Ferr != nil {
					for i = 0; i < NG; i++ {
						fmt.Fprintf(Ferr, "%s\t", Mpath.Cbcmp[i].Name)

						for j = 0; j < NG; j++ {
							fmt.Fprintf(Ferr, "\t%.1g", A[i*NG+j])
						}

						fmt.Fprintf(Ferr, "\t\t%.2g\n", Y[i])
					}
				}

				if NG > 1 {
					Matinv(A, NG, NG, "<Pflow>")
					Matmalv(A, Y, NG, NG, X)
				} else {
					X[0] = Y[0] / A[0]
				}

				if DEBUG {
					fmt.Printf("<Pflow>  Flow Rate\n")
					for i = 0; i < NG; i++ {
						fmt.Printf("\t%6.2f\n", X[i])
					}
				}

				if dayprn && Ferr != nil {
					for i = 0; i < NG; i++ {
						fmt.Fprintf(Ferr, "\t\t%.2g\n", X[i])
					}
				}
			}

			for i := 0; i < NG; i++ {
				pl = Mpath.Pl[i]
				pl.G = X[i]
			}

			for i, Plist := range Mpath.Plist {

				if DEBUG {
					fmt.Printf("<< Pflow >> e i=%d iMAX=%d control=%c G=%.5f\n",
						i, len(Mpath.Plist), Plist.Control, Plist.G)
				}

				if Plist.Control == OFF_SW {
					Plist.G = 0.0
				} else if Plist.G <= 0.0 {
					// 負であればエラーを表示する
					//if (Plist.G < 0. )
					//	fmt.Printf("<%s>  流量が負になっています %g\n", Mpath.Name, Plist.G ) ;

					Plist.G = 0.0
					Plist.Control = OFF_SW
					Plist.lpathscdd(Plist.Control)
				}

				for j, Pelm := range Plist.Pelm {
					if Pelm.Out != nil {
						Pelm.Out.G = Plist.G
					}

					if DEBUG {
						if Pelm.Out != nil {
							G0 = Pelm.Out.G
						} else {
							G0 = 0.0
						}

						fmt.Printf("< Pflow > j=%d\tjMAX=%d\tPelm-G=%.5f\tPlist->G=%.5f\n",
							j, len(Plist.Pelm), G0, Plist.G)
					}
				}
			}
		}
	}
}

/*
get_string_or_null (Get String or Null Representation)

この関数は、与えられた文字列`s`が空文字列の場合に`"(null)"`を返し、
それ以外の場合には元の文字列をそのまま返します。

建築環境工学的な観点:
- **デバッグ出力の可読性向上**: シミュレーションのデバッグ出力において、
  文字列が空の場合に`"(null)"`と表示することで、
  出力の可読性を向上させ、
  データが欠落している箇所を視覚的に分かりやすくします。
  これは、シミュレーションモデルの検証や問題の特定に役立ちます。

この関数は、建物のエネルギーシミュレーションにおいて、
デバッグ出力の品質を向上させるための補助的な役割を果たします。
*/
func get_string_or_null(s string) string {
	if s == "" {
		return "(null)"
	}
	return s
}

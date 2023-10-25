package eeslism

import (
	"errors"
	"fmt"
	"strings"
)

/* 経路の定義用関数 */

/* ----------------------------------------------- */

// システム要素出力端割当
func pelmco(pflow FliudType, Pelm *PELM, errkey string) {
	var Nout int
	err := 0
	var cmp *COMPNT
	var elmo *ELOUT

	cmp = Pelm.Cmp
	Nout = cmp.Nout

	elmo_idx := 0
	elmo = cmp.Elouts[elmo_idx]

	if cmp.Eqptype == DIVERG_TYPE || cmp.Eqptype == DIVGAIR_TYPE {
		for i := 0; i < Nout; i++ {
			elmo = cmp.Elouts[elmo_idx]
			if elmo.Id != '*' {
				Pelm.Co = '*'
				elmo.Id = '*'
				Pelm.Out = elmo
				err = 0
				break
			}
			elmo_idx++
		}
	} else if cmp.Eqptype == STANK_TYPE {
		stank := cmp.Eqp.(*STANK)
		var i, ii int
		for i = 0; i < stank.Nin; i++ {
			if stank.Pthcon[i] == Pelm.Co {
				break
			}
		}
		ii = i
		for i = 0; i < Nout; i++ {
			elmo = cmp.Elouts[elmo_idx]
			if elmo.Id == stank.Pthcon[ii] {
				Pelm.Out = elmo
				err = 0
				break
			}
			elmo_idx++
		}
	} else if Nout == 1 {
		Pelm.Out = elmo
		Pelm.Co = elmo.Id
	} else if cmp.Eqptype == RDPANEL_TYPE {
		if pflow == WATER_FLD || pflow == AIRa_FLD {
			Pelm.Out = elmo
			Pelm.Co = elmo.Id
		} else if pflow == AIRx_FLD {
			elmo_idx++
			Pelm.Out = elmo
			Pelm.Co = elmo.Id
		}
	} else {
		err = 1
		for i := 0; i < Nout; i++ {
			elmo = cmp.Elouts[i]
			if Pelm.Co == elmo.Id {
				Pelm.Out = elmo
				err = 0
				break
			}
		}
	}

	if err != 0 {
		for i := 0; i < Nout; i++ {
			elmo = cmp.Elouts[i]

			if (pflow == AIRa_FLD && elmo.Id == 't') ||
				(pflow == AIRx_FLD && elmo.Id == 'x') ||
				(pflow == WATER_FLD && elmo.Id == 'W') {
				Pelm.Out = elmo
				Pelm.Co = elmo.Id
				err = 0
				break
			}
		}
	}

	// Satoh追加　気化冷却器　2013/10/31
	if err != 0 {
		if cmp.Eqptype == EVAC_TYPE {
			Nout = 4
		}

		for i := 0; i < Nout; i++ {
			elmo = cmp.Elouts[i]
			if Pelm.Co == elmo.Id ||
				(Pelm.Co == 'W' && elmo.Id == 'V') ||
				(Pelm.Co == 'w' && elmo.Id == 'v') {
				Pelm.Out = elmo
				err = 0
				break
			}
		}
	}

	Errprint(err, errkey+" <pelmco>", cmp.Name)
}

/* ----------------------------------------------- */

// システム要素入力端割当
func pelmci(pflow FliudType, Pelm *PELM, errkey string) {
	err := 0
	cmp := Pelm.Cmp
	elmi := cmp.Elins[0]
	Nin := cmp.Nin

	// 入口の数が0の場合は処理を行わない
	if Nin <= 0 {
		return
	}

	if cmp.Eqptype == CONVRG_TYPE || cmp.Eqptype == CVRGAIR_TYPE {
		// 合流要素の場合
		//
		for i := 0; i < Nin; i++ {
			elmi := cmp.Elins[i]
			if elmi.Id != ELIO_ASTER {
				elmi.Id = ELIO_ASTER
				Pelm.Ci = ELIO_ASTER
				Pelm.In = elmi
				break
			}
		}
	} else if SIMUL_BUILDG && cmp.Eqptype == ROOM_TYPE {
		// 室要素の場合
		//
		room := cmp.Eqp.(*ROOM)

		var elins *[]*ELIN
		if pflow == AIRa_FLD {
			// 流体が空気(温度)の場合:
			elins = &room.elinasup
		} else if pflow == AIRx_FLD {
			// 流体が空気(湿度)の場合:
			elins = &room.elinasupx
		} else {
			panic(pflow)
		}

		for i := 0; i < room.Nasup; i++ {
			elmi := (*elins)[i]
			if elmi.Id != ELIO_ASTER {
				Pelm.Ci = ELIO_ASTER
				elmi.Id = ELIO_ASTER
				Pelm.In = elmi
				break
			}
		}
	} else if SIMUL_BUILDG && cmp.Eqptype == RDPANEL_TYPE {
		// 放射パネル要素の場合
		//
		if pflow == WATER_FLD || pflow == AIRa_FLD {
			Elout_a := cmp.Elouts[0] // 温水または空気温度
			for i := 0; i < Elout_a.Ni; i++ {
				elmi := cmp.Elins[i]
				if elmi.Id == ELIO_f {
					Pelm.Ci = elmi.Id
					Pelm.In = elmi
					break
				}
			}
		} else if pflow == AIRx_FLD {
			Elout_x := cmp.Elouts[1] // 空気湿度
			elmi := Elout_x.Elins[0]
			Pelm.Ci = elmi.Id
			Pelm.In = elmi
		}
	} else if Nin == 1 {
		// 入口の数が1の場合
		//
		Pelm.In = elmi
		Pelm.Ci = elmi.Id
	} else {
		err = 1

		for i := 0; i < Nin; i++ {
			elmi := cmp.Elins[i]

			// ACの絶対湿度はここに入った
			if Pelm.Ci == elmi.Id {
				Pelm.In = elmi
				err = 0
				break
			}
		}
	}

	if err != 0 {

		if cmp.Eqptype == HCLOADW_TYPE {
			Hcload := cmp.Eqp.(*HCLOAD)
			if Hcload.Wet {
				Nin = 4
			}
		} else if cmp.Eqptype == DESI_TYPE {
			Nin = 4
		}

		for i := 0; i < Nin; i++ {
			elmi := cmp.Elins[i]
			if (pflow == AIRa_FLD && elmi.Id == ELIO_t) ||
				(pflow == AIRx_FLD && elmi.Id == ELIO_x) ||
				(pflow == WATER_FLD && elmi.Id == ELIO_W) {
				Pelm.In = elmi
				Pelm.Ci = elmi.Id
				err = 0
				break
			}
		}

		//printf("\n") ;
	}

	if err != 0 {

		if cmp.Eqptype == THEX_TYPE {
			N := 0
			for i := 0; i < cmp.Nout; i++ {
				Eo := cmp.Elouts[i]
				N += Eo.Ni
			}

			Nin = N
		}

		for i := 0; i < Nin; i++ {
			elmi := cmp.Elins[i]
			if Pelm.Ci == elmi.Id {
				Pelm.In = elmi
				err = 0
				break
			}
		}
	}

	if err != 0 {
		if cmp.Eqptype == EVAC_TYPE {
			N := 0
			for i := 0; i < cmp.Nout; i++ {
				Eo := cmp.Elouts[i]
				N += Eo.Ni
			}

			Nin = N
		}

		for i := 0; i < Nin; i++ {
			elmi := cmp.Elins[i]
			if Pelm.Ci == elmi.Id ||
				(Pelm.Ci == 'W' && elmi.Id == 'V') ||
				(Pelm.Ci == 'w' && elmi.Id == 'v') {
				Pelm.In = elmi
				err = 0
				break
			}
		}
	}

	Errprint(err, errkey+" <pelmci>", cmp.Name)
}

/* ----------------------------------------------- */

// システム要素接続データのコピー（空気系統湿度経路用
// 空気経路の場合は湿度経路用にpathをコピーする
//
// - 空気温度用のシステム経路 mpath_t の設定を 空湿度用のシステム経路 にコピーする。
// - コピーに際して要素(PELM)を追加する。
// - 要素(PELM)は _Pelm 配列に追加するものとし、 Npelm 番目から詰め込むとする。 Npelmは上書きする。
// ##- _Plist は システム経路 mpath_t に属するすべての 末端経路の配列である。
// - Compntには SYSCMPデータセットで読み込んだすべての機器情報が保持されている。
func plistcpy(mpath_t *MPATH, _Pelm *[]*PELM, _Plist *[]*PLIST, Compnt []*COMPNT) *MPATH {
	// 空気湿度用経路
	var mpath_x *MPATH = NewMPATH()
	mpath_x.Name = mpath_t.Name + ".x" // 湿度用経路の名前 = 温度経路用の名前 + ".x"
	//mpath_x.Plist = _Plist             // 末端経路 (要確認)
	mpath_x.Fluid = AIRx_FLD    // 流体種別 = 空気湿度
	mpath_x.G0 = mpath_t.G0     // 流量比率
	mpath_x.Rate = mpath_t.Rate // 流量比率フラグ

	// 空気温度用経路
	mpath_t.Fluid = AIRa_FLD // 流体種別を念のため上書き?
	mpath_t.Mpair = mpath_x  // 空気湿度経路への参照(Mpair)を設定

	// 末端経路についてループ
	for i := range mpath_t.Plist {
		pli := mpath_t.Plist[i]
		Plist := NewPLIST()

		// ターゲットの末端経路
		//pli := &mpath_t.Plist[i]

		// 相互参照設定
		pli.Lpair = Plist
		pli.Plistx = Plist
		Plist.Plistt = pli

		// コピー
		Plist.Pelm = nil
		Plist.Org = false
		Plist.Type = pli.Type
		Plist.Go = pli.Go
		Plist.Nvav = pli.Nvav
		Plist.Nvalv = pli.Nvalv
		Plist.NOMVAV = pli.NOMVAV
		Plist.OMvav = pli.OMvav
		Plist.Valv = pli.Valv
		Plist.Rate = pli.Rate
		Plist.UnknownFlow = pli.UnknownFlow

		// 名前のコピー: ".x"を付与しながらコピー
		if pli.Name != "" {
			Plist.Name = pli.Name + ".x"
		} else {
			Plist.Name = ".x"
		}

		// 要素のコピー
		nelm := 0
		Plist.Pelm = make([]*PELM, 0, len(pli.Pelm))
		for _, peli := range pli.Pelm {

			// コピー対象は空気経路のみ
			if !peli.Cmp.Airpathcpy {
				continue
			}

			var Pelm *PELM = NewPELM()

			*_Pelm = append(*_Pelm, Pelm)
			Plist.Pelm = append(Plist.Pelm, Pelm)

			if peli.Cmp.Eqptype == CVRGAIR_TYPE || peli.Cmp.Eqptype == DIVGAIR_TYPE {
				// ** 合流要素の場合 **

				// Find index
				k, err := FindComponentRef(peli.Cmp, Compnt)
				if err != nil {
					panic(err)
				}

				// k+1番目位以降のコンポーネントのみ検索している: 理由？？
				var cmp *COMPNT
				for k++; k < len(Compnt); k++ {
					cmp = Compnt[k]
					s := cmp.Name

					// "name.xxx" のうち name だけで一致する機器を探す。
					if idx := strings.IndexRune(s, '.'); idx >= 0 {
						s = s[:idx]

						if peli.Cmp.Name == s {
							break
						}
					}
				}
				Pelm.Cmp = cmp // 検索で見つけた機器参照
			} else if peli.Cmp.Eqptype == THEX_TYPE {
				// ** 全熱交換器の場合 **
				Pelm.Cmp = peli.Cmp
				if peli.Ci == ELIO_E {
					Pelm.Ci = ELIO_e
					Pelm.Co = ELIO_e
				} else {
					Pelm.Ci = ELIO_o
					Pelm.Co = ELIO_o
				}
			} else if peli.Cmp.Eqptype == EVAC_TYPE {
				// Satoh追加　気化冷却器　2013/10/31
				Pelm.Cmp = peli.Cmp
				if peli.Ci == ELIO_D {
					Pelm.Ci = ELIO_d
					Pelm.Co = ELIO_d
				} else if peli.Ci == ELIO_W {
					Pelm.Ci = ELIO_w
					Pelm.Co = ELIO_w
				}
			} else {
				Pelm.Cmp = peli.Cmp
				Pelm.Ci = peli.Ci
				Pelm.Co = peli.Co
			}

			Pelm.Out = peli.Out
			nelm++
		}

		mpath_x.Plist = append(mpath_x.Plist, Plist)
	}

	return mpath_x
}

/* ----------------------------------------------- */

/*  合流レベルの設定  */

func plevel(Nmpath int, Mpath []MPATH, Ncnvrg int, Cnvrg []*COMPNT) {

	var i, j int
	lvc := 0
	var lvcmx, lvcf int
	var Plist *PLIST
	var cmp *COMPNT
	var elin *ELIN

	for i = 0; i < Ncnvrg; i++ {
		cmp = Cnvrg[i]
		cmp.Elouts[0].Lpath.Lvc = -1
	}

	lvcf = Ncnvrg

	for lvcf > 0 {
		for i = 0; i < Ncnvrg; i++ {
			cmp = Cnvrg[i]
			if cmp.Elouts[0].Lpath.Lvc <= 0 {
				for j = 0; j < cmp.Nin; j++ {
					elin = cmp.Elins[j]
					Plist = elin.Lpath
					if Plist.Type != CONVRG_LPTP {
						lvc = 0
					} else {
						if Plist.Lvc > 0 {
							if lvc <= Plist.Lvc {
								lvc = Plist.Lvc
							}
						} else {
							break
						}
					}
				}

				if j == cmp.Nin {
					lvc++
					cmp.Elouts[0].Lpath.Lvc = lvc

					lvcf--
				}
			}
		}
	}

	for i = 0; i < Nmpath; i++ {
		_Mpath := &Mpath[i]

		lvcmx = 0

		for _, Plist := range _Mpath.Plist {
			if Plist.Lvc > lvcmx {
				lvcmx = Plist.Lvc
			}
		}
		_Mpath.Lvcmx = lvcmx
	}
}

/* ----------------------------------------------- */

func pflowstrct(_Mpath []*MPATH) {
	var j, n, M, MM, k int
	var etype EqpType
	var Elout *ELOUT
	var Elin *ELIN

	var nplist int = 0
	for _, Mpath := range _Mpath {
		nplist += len(Mpath.Plist)
	}

	for _, Mpath := range _Mpath {
		Mpath.Pl = make([]*PLIST, nplist)
		Mpath.Cbcmp = make([]*COMPNT, nplist)
	}

	for _, Mpath := range _Mpath {
		n = 0

		for _, Plist := range Mpath.Plist {
			// 流量未設定の末端経路を検索
			if Plist.Go == nil && Plist.Nvav == 0 &&
				Plist.Rate == nil && Plist.NOMVAV == 0 &&
				(Plist.Nvalv == 0 || (Plist.Nvalv > 0 && Plist.Valv.MonPlist == nil && Plist.Valv.MGo == nil)) {
				Mpath.Pl[n] = Plist
				Plist.N = n
				// 末端経路未知フラグの変更
				Plist.UnknownFlow = 0
				n++
			}

			if Mpath.Rate && (Plist.Go != nil || Plist.Nvav > 0 || Plist.Nvalv > 0 || Plist.NOMVAV > 0) {
				Mpath.G0 = &Plist.G // 流量比率設定時の既知流量へのポインタをセット
			}
		}

		Mpath.NGv = n
		Mpath.NGv2 = n * n

		n = 0

		for _, Plist := range Mpath.Plist {

			// 末端経路の先頭機器
			etype = Plist.Pelm[0].Cmp.Eqptype

			// 末端経路の先頭が分岐・合流の場合
			if etype == CONVRG_TYPE ||
				etype == CVRGAIR_TYPE ||
				etype == DIVERG_TYPE ||
				etype == DIVGAIR_TYPE {
				MM = 0

				for k = 0; k < n; k++ {
					if Plist.Pelm[0].Cmp == Mpath.Cbcmp[k] {
						MM++
						break
					}
				}

				if MM == 0 {
					// 末端経路の先頭機器（分岐・合流）の入口、出口経路の未知流量の数(M)を数える
					M = 0

					// 末端経路の先頭機器の入口
					for j = 0; j < Plist.Pelm[0].Cmp.Nin; j++ {
						Elin = Plist.Pelm[0].Cmp.Elins[j]

						Plist.Upplist = Elin.Lpath
						// 末端経路の先頭機器の上流の流量が未定義
						if Elin.Lpath.UnknownFlow == 0 {
							M++
							break
						}
					}

					if M == 0 {
						// 末端経路の先頭機器の出口
						Elout = Plist.Pelm[0].Cmp.Elouts[0]

						for j = 0; j < Plist.Pelm[0].Cmp.Nout; j++ {
							Elout = Plist.Pelm[0].Cmp.Elouts[j]

							Plist.Dnplist = Elout.Lpath
							// 末端経路の先頭機器の出口経路の流量が未知なら
							if Elout.Lpath.UnknownFlow == 0 {
								M++
								break
							}
						}
					}

					// 末端経路先頭にある分岐・合流の入口もしくは出口経路の未知流量の数(M)
					if M > 0 {
						// 流量連立方程式を解くときに使用する分岐・合流機器と数(n)
						Mpath.Cbcmp[n] = Plist.Pelm[0].Cmp
						n++
					}
				}
			}
		}

		// 既知末端流量数のチェック
		if n > 0 && (n-1 != Mpath.NGv) {
			fmt.Printf("<%s> 末端流量の与えすぎ、もしくは少なすぎです n=%d NGv=%d\n", Mpath.Name, n, Mpath.NGv)
		}
	}
}

func FindComponentRef(target *COMPNT, Compnt []*COMPNT) (int, error) {
	for k := range Compnt {
		cmp := Compnt[k]
		if cmp == target {
			return k, nil
		}
	}
	return -1, errors.New("Not Found")
}

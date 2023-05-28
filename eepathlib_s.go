package main

import (
	"fmt"
	"strings"
)

/* 経路の定義用関数 */

/* ----------------------------------------------- */

/*  システム要素出力端割当  */

func pelmco(pflow rune, Pelm *PELM, errkey string) {
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

/*  システム要素入力端割当  */

func pelmci(pflow rune, Pelm *PELM, errkey string) {
	var Nin int
	err := 0
	var N int
	var cmp *COMPNT
	var elmi *ELIN
	var Elout, Eo *ELOUT

	var room *ROOM
	var Hcload *HCLOAD

	cmp = Pelm.Cmp
	elmi = cmp.Elins[0]
	Nin = cmp.Nin

	if Nin <= 0 {
		return
	} else if cmp.Eqptype == CONVRG_TYPE || cmp.Eqptype == CVRGAIR_TYPE {
		for i := 0; i < Nin; i++ {
			elmi = cmp.Elins[i]
			if elmi.Id != '*' {
				Pelm.Ci = '*'
				elmi.Id = '*'
				Pelm.In = elmi
				break
			}
		}
	} else if SIMUL_BUILDG && cmp.Eqptype == ROOM_TYPE {
		room = cmp.Eqp.(*ROOM)
		/*************
		if (pflow == AIRa_FLD)
		ii = room.Nachr + room.Nrp;
		else if (pflow == AIRx_FLD)
		ii = Nin - room.Nasup;
		elmi += ii;
		**************/

		if pflow == AIRa_FLD {
			elmi = room.elinasup[0]
		} else if pflow == AIRx_FLD {
			elmi = room.elinasupx[0]
		}

		/***************
		printf("<<pelmci>>  room=%s  ii=%d  Nin=%d\n", room.name, ii, Nin);
		***********************/

		for i := 0; i < room.Nasup; i++ {
			elmi = cmp.Elins[i]
			if elmi.Id != '*' {
				Pelm.Ci = '*'
				elmi.Id = '*'
				Pelm.In = elmi
				break
			}
		}
	} else if SIMUL_BUILDG && cmp.Eqptype == RDPANEL_TYPE {

		Elout = cmp.Elouts[0]
		if pflow == WATER_FLD || pflow == AIRa_FLD {
			for i := 0; i < Elout.Ni; i++ {
				elmi = cmp.Elins[i]
				if elmi.Id == 'f' {
					Pelm.Ci = elmi.Id
					Pelm.In = elmi
					break
				}
			}
		} else if pflow == AIRx_FLD {
			Elout = cmp.Elouts[1]
			Pelm.In = Elout.Elins[0]
			Pelm.Ci = Elout.Elins[0].Id
		}
	} else if Nin == 1 {
		Pelm.In = elmi
		Pelm.Ci = elmi.Id
	} else {
		err = 1

		for i := 0; i < Nin; i++ {
			elmi = cmp.Elins[i]

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
			Hcload = cmp.Eqp.(*HCLOAD)
			if Hcload.Wet == 'y' {
				Nin = 4
			}
		} else if cmp.Eqptype == DESI_TYPE {
			Nin = 4
		}

		for i := 0; i < Nin; i++ {
			elmi = cmp.Elins[i]
			if (pflow == AIRa_FLD && elmi.Id == 't') ||
				(pflow == AIRx_FLD && elmi.Id == 'x') ||
				(pflow == WATER_FLD && elmi.Id == 'W') {
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
			N = 0
			for i := 0; i < cmp.Nout; i++ {
				Eo = cmp.Elouts[i]
				N += Eo.Ni
			}

			Nin = N
		}

		for i := 0; i < Nin; i++ {
			elmi = cmp.Elins[i]
			if Pelm.Ci == elmi.Id {
				Pelm.In = elmi
				err = 0
				break
			}
		}
	}

	if err != 0 {
		if cmp.Eqptype == EVAC_TYPE {
			N = 0
			for i := 0; i < cmp.Nout; i++ {
				Eo = cmp.Elouts[i]
				N += Eo.Ni
			}

			Nin = N
		}

		for i := 0; i < Nin; i++ {
			elmi = cmp.Elins[i]
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

/* システム要素接続データのコピー（空気系統湿度経路用） */

func plistcpy(Mpath *MPATH, Mpath_prev *MPATH, Npelm *int, _Pelm []PELM, _Plist []PLIST,
	Ncompnt int, Compnt []COMPNT) {
	var mpi *MPATH
	var cmp *COMPNT
	var i, j, nelm int
	var s string

	mpi = Mpath_prev

	mpi.Mpair = Mpath

	Mpath.Name = mpi.Name + ".x"

	Mpath.Nlpath = mpi.Nlpath
	Mpath.Plist = _Plist
	mpi.Fluid = AIRa_FLD
	Mpath.Fluid = AIRx_FLD
	Mpath.G0 = mpi.G0
	Mpath.Rate = mpi.Rate

	for i = 0; i < mpi.Nlpath; i++ {
		pli := &mpi.Plist[i]
		Plist := &_Plist[i]

		pli.Lpair = Plist
		pli.Plistx = Plist
		Plist.Plistt = pli

		nelm = 0
		Plist.Pelm = nil
		Plist.Org = 'n'
		Plist.Type = pli.Type
		Plist.Go = pli.Go
		Plist.Nvav = pli.Nvav
		Plist.Nvalv = pli.Nvalv
		Plist.NOMVAV = pli.NOMVAV
		Plist.OMvav = pli.OMvav
		Plist.Valv = pli.Valv
		Plist.Rate = pli.Rate
		Plist.UnknownFlow = pli.UnknownFlow

		if pli.Name != "" {
			Plist.Name = pli.Name + ".x"
		} else {
			Plist.Name = ".x"
		}

		for j = 0; j < pli.Nelm; j++ {
			peli := pli.Pelm[j]

			if peli.Cmp.Airpathcpy == 'y' {
				Pelm := &_Pelm[nelm]

				if Plist.Pelm == nil {
					Plist.Pelm = []*PELM{Pelm}
				}

				(*Npelm)++

				if peli.Cmp.Eqptype == CVRGAIR_TYPE ||
					peli.Cmp.Eqptype == DIVGAIR_TYPE {

					// Find index
					var k int
					for k = 0; k < Ncompnt; k++ {
						cmp = &Compnt[k]
						if cmp == peli.Cmp {
							break
						}
					}

					for ; k < Ncompnt; k++ {
						cmp = &Compnt[k]
						s = cmp.Name

						if idx := strings.IndexRune(s, '.'); idx >= 0 {
							s = s[:idx]

							if peli.Cmp.Name == s {
								break
							}
						}
					}
					Pelm.Cmp = cmp
				} else if peli.Cmp.Eqptype == THEX_TYPE {
					Pelm.Cmp = peli.Cmp
					if peli.Ci == 'E' {
						Pelm.Ci = 'e'
						Pelm.Co = 'e'
					} else {
						Pelm.Ci = 'o'
						Pelm.Co = 'o'
					}
				} else if peli.Cmp.Eqptype == EVAC_TYPE {
					// Satoh追加　気化冷却器　2013/10/31
					Pelm.Cmp = peli.Cmp
					if peli.Ci == 'D' {
						Pelm.Ci = 'd'
						Pelm.Co = 'd'
					} else if peli.Ci == 'W' {
						Pelm.Ci = 'w'
						Pelm.Co = 'w'
					}
				} else {
					Pelm.Cmp = peli.Cmp
					Pelm.Ci = peli.Ci
					Pelm.Co = peli.Co
				}

				Pelm.Out = peli.Out
				nelm++
			}
		}
		Plist.Nelm = nelm
	}
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

		for j = 0; j < _Mpath.Nlpath; j++ {
			Plist = &_Mpath.Plist[j]
			if Plist.Lvc > lvcmx {
				lvcmx = Plist.Lvc
			}
		}
		_Mpath.Lvcmx = lvcmx
	}
}

/* ----------------------------------------------- */

func pflowstrct(Nmpath int, _Mpath []MPATH) {
	var m, i, j, n, M, MM, k int
	var Plist *PLIST
	var etype EqpType
	var Elout *ELOUT
	var Elin *ELIN

	for m = 0; m < Nmpath; m++ {
		Mpath := _Mpath[m]

		n = 0

		for i = 0; i < Mpath.Nlpath; i++ {
			Plist = &Mpath.Plist[i]

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

			if Mpath.Rate == 'Y' && (Plist.Go != nil || Plist.Nvav > 0 || Plist.Nvalv > 0 || Plist.NOMVAV > 0) {
				Mpath.G0 = &Plist.G // 流量比率設定時の既知流量へのポインタをセット
			}
		}

		Mpath.NGv = n
		Mpath.NGv2 = n * n

		n = 0

		for i = 0; i < Mpath.Nlpath; i++ {
			Plist = &Mpath.Plist[i]

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

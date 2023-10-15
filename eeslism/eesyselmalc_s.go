package eeslism

/* 機器使用データの割り付けおよびシステム要素から入力、出力要素の割り付け */

const idmrkc = "txW"

func Elmalloc(
	errkey string,
	_Compnt []COMPNT,
	Eqcat *EQCAT,
	Eqsys *EQSYS,
	Elo *[]*ELOUT,
	Nelout *int,
	Eli *[]*ELIN,
	Nelin *int,
) {
	var cmp []COMPNT
	var elop, elo *ELOUT
	var room *ROOM
	var rdpnl *RDPNL
	var Hcc *HCC
	var Boi *BOI
	var Refa *REFA
	var Coll *COLL
	var Pipe *PIPE
	var Stank *STANK
	var Pump *PUMP
	var Cnvrg []*COMPNT
	var Flin []FLIN
	var Hcload []HCLOAD
	var Stheat *STHEAT
	// var Elout []*ELOUT
	// var Elin []*ELIN
	var Thex *THEX
	var PV *PV
	var Desi *DESI
	var Evac *EVAC

	var i, ii, mm, neqp, ncat int
	icv := 0
	var name string

	var id string
	var Nvalv, NQmeas, NOMvav int
	idTe := "EO"
	idTo := "OE"
	var idxe, idxo []ELIOType

	cmp = _Compnt
	Nout, Nin := Elcount(_Compnt)

	if Nout > 0 {
		*Elo = make([]*ELOUT, Nout)
		Eloutinit(*Elo, Nout)
	}

	if Nin > 0 {
		*Eli = make([]*ELIN, Nin)
		Elininit(Nin, *Eli)
	}

	eloIdx := 0
	elinIdx := 0

	flinIdx := 0
	hcloadIdx := 0
	Cnvrg = Eqsys.Cnvrg
	Flin = Eqsys.Flin
	Hcload = Eqsys.Hcload

	for m := range _Compnt {
		Compnt := &_Compnt[m]

		if Compnt.Eqptype != PV_TYPE {
			Compnt.Elouts = (*Elo)[eloIdx : eloIdx+Compnt.Nout]
			Compnt.Elins = (*Eli)[elinIdx : elinIdx+Compnt.Nin]
		}

		name = Compnt.Name
		neqp = Compnt.Neqp
		ncat = Compnt.Ncat

		c := Compnt.Eqptype

		if SIMUL_BUILDG && c == ROOM_TYPE {
			room = Compnt.Eqp.(*ROOM)
			room.cmp = Compnt

			id = idmrkc
			for i = 0; i < 2; i++ {
				Elout := (*Elo)[eloIdx]
				Elout.Cmp = Compnt
				Elout.Id = ELIOType(id[i])
				Elout.Fluid = rune(id[i])
				if i == 0 {
					Elout.Ni = room.Nachr + room.Ntr + room.Nrp + room.Nasup
				} else if i == 1 {
					Elout.Ni = room.Nachr + room.Nasup
				}
				Elout.Elins = (*Eli)[elinIdx : elinIdx+Elout.Ni]

				elinIdx += Elout.Ni
				eloIdx++
			}
			elo = Compnt.Elouts[0]
			room.elinasup = elo.Elins[room.Nachr+room.Ntr+room.Nrp:]
			elo = Compnt.Elouts[1]
			room.elinasupx = elo.Elins[room.Nachr:]
		} else if SIMUL_BUILDG && c == RDPANEL_TYPE {
			rdpnl = Compnt.Eqp.(*RDPNL)
			rdpnl.cmp = Compnt
			rdpnl.Tpi = 15.0

			Elout := (*Elo)[eloIdx]
			Elout.Cmp = Compnt
			Elout.Id = 'f'
			Elout.Ni = 1 + 1 + rdpnl.Ntrm[0] + rdpnl.Nrp[0]
			if rdpnl.MC == 2 {
				Elout.Ni += 1 + rdpnl.Ntrm[1] + rdpnl.Nrp[1]
			}
			Elout.Elins = (*Eli)[elinIdx : elinIdx+Elout.Ni]

			eloIdx++

			Elin := (*Eli)[elinIdx]
			Elin.Id = ELIO_f
			elinIdx++

			for mm = 0; mm < rdpnl.MC; mm++ {

				Elin = (*Eli)[elinIdx]
				Elin.Id = ELIO_r
				elinIdx++

				for ii = 0; ii < rdpnl.Ntrm[mm]; ii++ {
					Elin = (*Eli)[elinIdx]
					Elin.Id = ELIO_r
					elinIdx++
				}
				for ii = 0; ii < rdpnl.Nrp[mm]; ii++ {
					Elin = (*Eli)[elinIdx]
					Elin.Id = ELIO_f
					elinIdx++
				}
			}

			/* 空気経路湿度用 */
			Elout = (*Elo)[eloIdx]
			Elout.Ni = 1
			Elout.Cmp = Compnt
			Elout.Elins = (*Eli)[elinIdx : elinIdx+1]
			Elout.Id = ELIO_x
			Elin.Id = ELIO_x
			eloIdx++
			elinIdx++
		} else if c == DIVERG_TYPE || c == DIVGAIR_TYPE {
			if c == DIVGAIR_TYPE {
				Compnt.Airpathcpy = 'y'
			} else {
				Compnt.Airpathcpy = 'n'
			}

			Compnt.Nin = 1

			for i = 0; i < Compnt.Nout; i++ {
				Elout := (*Elo)[eloIdx]
				Elout.Cmp = Compnt

				Elout.Ni = Compnt.Nin
				Elout.Elins = (*Eli)[elinIdx : elinIdx+Elout.Ni]

				Elout.Id = Compnt.Ido[i]
				eloIdx++
			}

			Elin := (*Eli)[elinIdx]
			Elin.Id = 'i'
			elinIdx++
		} else if c == CONVRG_TYPE || c == CVRGAIR_TYPE {
			if c == CVRGAIR_TYPE {
				Compnt.Airpathcpy = 'y'
			} else {
				Compnt.Airpathcpy = 'n'
			}

			Compnt.Nout = 1

			Cnvrg[icv] = Compnt
			icv++

			Elout := (*Elo)[eloIdx]
			Elout.Id = 'o'
			Elout.Cmp = Compnt

			Elout.Ni = Compnt.Nin
			Elout.Elins = (*Eli)[elinIdx : elinIdx+Elout.Ni]

			eloIdx++
			for i = 0; i < Compnt.Nin; i++ {
				Elin := (*Eli)[elinIdx]
				Elin.Id = Compnt.Idi[i]
				elinIdx++
			}
		} else if c == HCCOIL_TYPE {
			Hcc = &Eqsys.Hcc[neqp]
			Compnt.Eqp = Hcc
			Hcc.Name = name
			Hcc.Cmp = Compnt
			Hcc.Cat = &Eqcat.Hccca[ncat]

			for i = 0; i < Compnt.Nout; i++ {
				Elin := (*Eli)[elinIdx]
				Elout := (*Elo)[eloIdx]

				Elout.Cmp = Compnt

				Elout.Ni = Compnt.Nin
				Elout.Elins = Compnt.Elins

				Elout.Id = Compnt.Ido[i]
				Elin.Id = Compnt.Idi[i]

				eloIdx++
				elinIdx++
			}
		} else if c == HEXCHANGR_TYPE {
			Hex := &Eqsys.Hex[neqp]
			Compnt.Eqp = Hex
			Hex.Name = name
			Hex.Cmp = Compnt
			Hex.Cat = &Eqcat.Hexca[ncat]

			for i = 0; i < Compnt.Nout; i++ {
				Elin := (*Eli)[elinIdx]
				Elout := (*Elo)[eloIdx]

				Elout.Cmp = Compnt

				Elout.Ni = Compnt.Nin
				Elout.Elins = Compnt.Elins

				Elout.Id = Compnt.Ido[i]
				Elin.Id = Compnt.Idi[i]

				eloIdx++
				elinIdx++
			}
		} else if c == BOILER_TYPE {
			Boi = &Eqsys.Boi[neqp]
			Compnt.Eqp = Boi
			Boi.Name = name
			Boi.Cmp = Compnt
			Boi.Cat = &Eqcat.Boica[ncat]

			Elout := (*Elo)[eloIdx]
			Elout.Cmp = Compnt
			Elout.Ni = Compnt.Nin
			Elout.Elins = (*Eli)[elinIdx : elinIdx+Elout.Ni]

			eloIdx++
			elinIdx++
		} else if c == COLLECTOR_TYPE || c == ACOLLECTOR_TYPE {
			Coll = &Eqsys.Coll[neqp]
			Compnt.Eqp = Coll
			Coll.Name = name
			Coll.Cmp = Compnt
			Coll.Cat = &Eqcat.Collca[ncat]
			Coll.Ac = Compnt.Ac

			if Coll.Cat.Type == COLLECTOR_PDT {
				Elout := (*Elo)[eloIdx]
				Elout.Cmp = Compnt

				Elout.Ni = 1
				Elout.Elins = (*Eli)[elinIdx : elinIdx+1]

				eloIdx++
				elinIdx++
			} else {
				id = idmrkc
				for i = 0; i < Compnt.Nout; i++ {
					Elin := (*Eli)[elinIdx]
					Elout := (*Elo)[eloIdx]

					Elout.Cmp = Compnt
					Elout.Id = ELIOType(idmrkc[i])
					Elin.Id = ELIOType(idmrkc[i])
					Elout.Ni = 1
					Elout.Elins = (*Eli)[elinIdx : elinIdx+1]

					eloIdx++
					elinIdx++
				}
			}
		} else if c == PV_TYPE {
			PV = &Eqsys.PVcmp[neqp]
			Compnt.Eqp = PV
			PV.Name = name
			PV.Cmp = Compnt
			PV.Cat = &Eqcat.PVca[ncat]
			PV.PVcap = Compnt.PVcap
			PV.Area = Compnt.Area
		} else if c == REFACOMP_TYPE {
			Refa = &Eqsys.Refa[neqp]
			Compnt.Eqp = Refa
			Refa.Name = name
			Refa.Cmp = Compnt
			Refa.Cat = &Eqcat.Refaca[ncat]

			Elout := (*Elo)[eloIdx]
			Elout.Cmp = Compnt
			Elout.Ni = Compnt.Nin
			Elout.Elins = (*Eli)[elinIdx : elinIdx+Elout.Ni]

			eloIdx++
			elinIdx++
		} else if c == PUMP_TYPE {
			Pump = &Eqsys.Pump[neqp]
			Compnt.Eqp = Pump
			Pump.Name = name
			Pump.Cmp = Compnt
			Pump.Cat = &Eqcat.Pumpca[ncat]

			if Pump.Cat.pftype == PUMP_PF {
				Elout := (*Elo)[eloIdx]
				Elout.Cmp = Compnt
				Elout.Elins = (*Eli)[elinIdx : elinIdx+1]
				Elout.Ni = 1
				eloIdx++
				elinIdx++
			} else {
				for i = 0; i < Compnt.Nout; i++ {
					Elin := (*Eli)[elinIdx]
					Elout := (*Elo)[eloIdx]

					Elout.Cmp = Compnt
					Elout.Id = ELIOType(idmrkc[i])
					Elin.Id = ELIOType(idmrkc[i])
					Elout.Elins = (*Eli)[elinIdx : elinIdx+1]
					Elout.Ni = 1

					eloIdx++
					elinIdx++
				}
			}
		} else if c == PIPEDUCT_TYPE {
			Pipe = &Eqsys.Pipe[neqp]
			Compnt.Eqp = Pipe
			Pipe.Name = name
			Pipe.Cmp = Compnt
			Pipe.Cat = &Eqcat.Pipeca[ncat]

			if Pipe.Cat.Type == PIPE_PDT {
				Elout := (*Elo)[eloIdx]

				Elout.Cmp = Compnt
				Elout.Elins = (*Eli)[elinIdx : elinIdx+1]
				Elout.Ni = 1

				eloIdx++
				elinIdx++
			} else {
				id = idmrkc
				for i = 0; i < Compnt.Nout; i++ {
					Elin := (*Eli)[elinIdx]
					Elout := (*Elo)[eloIdx]

					Elout.Cmp = Compnt
					Elout.Id = ELIOType(idmrkc[i])
					Elin.Id = ELIOType(idmrkc[i])
					Elout.Elins = (*Eli)[elinIdx : elinIdx+1]
					Elout.Ni = 1

					eloIdx++
					elinIdx++
				}
			}
		} else if c == STANK_TYPE {
			Stank = &Eqsys.Stank[neqp]
			Compnt.Eqp = Stank
			Stank.Name = name
			Stank.Cmp = Compnt
			Stank.Cat = &Eqcat.Stankca[ncat]

			Stankmemloc("Stankmemloc", Stank)

			Compnt.Nin = Stank.Nin
			Compnt.Nout = Stank.Nin
			Compnt.Idi = make([]ELIOType, Stank.Nin)
			Compnt.Ido = make([]ELIOType, Stank.Nin)

			for i = 0; i < Compnt.Nout; i++ {
				Elout := (*Elo)[eloIdx]

				Elout.Cmp = Compnt
				Elout.Ni = Compnt.Nin
				Elout.Elins = (*Eli)[elinIdx : elinIdx+Elout.Ni]
				Elout.Id = Stank.Pthcon[i]
				Compnt.Ido[i] = Stank.Pthcon[i]

				eloIdx++
			}

			for i = 0; i < Compnt.Nin; i++ {
				Elin := (*Eli)[elinIdx]

				Elin.Id = Stank.Pthcon[i]
				Compnt.Idi[i] = Stank.Pthcon[i]

				elinIdx++
			}
		} else if c == FLIN_TYPE {
			Compnt.Eqp = &Flin[flinIdx]
			Flin[flinIdx].Cmp = Compnt
			Flin[flinIdx].Name = name

			flindat(&Flin[flinIdx])

			for i = 0; i < Compnt.Nout; i++ {
				//Elin := (*Eli)[elinIdx]
				Elout := (*Elo)[eloIdx]

				Elout.Cmp = Compnt

				Elout.Ni = Compnt.Nin
				Elout.Elins = Compnt.Elins

				Elout.Id = Compnt.Ido[i]

				eloIdx++
				elinIdx++
			}

			flinIdx++
		} else if c == HCLOAD_TYPE ||
			c == HCLOADW_TYPE ||
			c == RMAC_TYPE ||
			c == RMACD_TYPE {
			// 仮想空調機

			Compnt.Eqp = Hcload
			Hcload[hcloadIdx].Cmp = Compnt
			Hcload[hcloadIdx].Name = name

			// ルームエアコンの場合
			if c == RMAC_TYPE {
				Hcload[hcloadIdx].RMACFlg = 'Y'

				// エアコンの機器スペックを読み込む
				rmacdat(&Hcload[hcloadIdx])
			} else if c == RMACD_TYPE {
				Hcload[hcloadIdx].RMACFlg = 'y'

				// エアコンの機器スペックを読み込む
				rmacddat(&Hcload[hcloadIdx])
			}

			/*---- Roh Debug for a constant outlet humidity model of wet coil  2003/4/25 ----*/
			if Compnt.Ivparm != nil {
				Hcload[hcloadIdx].RHout = *(Compnt.Ivparm)
			}

			if Compnt.Wetparm == "wet" {
				Hcload[hcloadIdx].Wet = 'y' // 湿りコイル（吹出相対湿度一定
			} else {
				Hcload[hcloadIdx].Wet = 'n' // 吹出相対湿度は成り行き
			}

			// 空気のみの流入、流出
			if c == HCLOAD_TYPE || c == RMAC_TYPE || c == RMACD_TYPE {
				Hcload[hcloadIdx].Type = HCLoadType_D
				Compnt.Nout = 2
				Compnt.Nin = 2
			} else {
				// 空気＋水の流入、流出
				Hcload[hcloadIdx].Type = 'W'
				Compnt.Nout = 3
				Compnt.Nin = 3
			}

			// 空気の絶対湿度用経路コピーを行う
			Compnt.Airpathcpy = 'y'
			id = idmrkc
			for i = 0; i < Compnt.Nout; i++ {
				Elin := (*Eli)[elinIdx]
				Elout := (*Elo)[eloIdx]

				Elout.Cmp = Compnt
				Elout.Id = ELIOType(idmrkc[i])
				Elin.Id = ELIOType(idmrkc[i])

				Elout.Ni = 1
				//Elout.Ni = 2;
				// 湿りコイル（吹出相対湿度一定）で出口絶対湿度の経路の場合
				// 要素方程式の未知数は2つ（入口絶対湿度と出口温度）
				if i == 1 && Hcload[hcloadIdx].Wet == 'y' {
					Elout.Ni = 2
				} else if i == 2 && Hcload[hcloadIdx].Type == 'W' {
					// 冷温水コイルで水側系統の場合
					// 要素方程式の未知数は5個
					// 水入口温度、空気入口温度、空気入口湿度
					// 空気出口温度、空気出口湿度
					Elout.Ni = 5
				}
				Elout.Elins = (*Eli)[elinIdx : elinIdx+Elout.Ni]

				for ii = 0; ii < Elout.Ni; ii++ {
					// 空気出口絶対湿度の計算の2つ目の変数は空気出口温度
					if i == 1 && ii == 1 {
						Elin.Id = ELIO_ASTER
					}
					elinIdx++
				}

				/***** printf("xxx Elmalloc xxx   %s  i=%d  Elout.Ni=%d\n",
				Hcload.name, i, Elout.Ni); *****/

				eloIdx++
			}
			hcloadIdx++
		} else if c == VAV_TYPE || c == VWV_TYPE {
			/*---- Satoh Debug VAV  2000/12/5 ----*/

			Compnt.Eqp = &Eqsys.Vav[neqp]
			Eqsys.Vav[neqp].Name = name
			Eqsys.Vav[neqp].Cmp = Compnt
			Eqsys.Vav[neqp].Cat = &Eqcat.Vavca[ncat]
			Compnt.Nin = 2
			Compnt.Nout = 2

			if Eqsys.Vav[neqp].Cat.Type == VAV_PDT {
				Compnt.Airpathcpy = 'y'
				for i = 0; i < Compnt.Nout; i++ {
					Elin := (*Eli)[elinIdx]
					Elout := (*Elo)[eloIdx]

					Elout.Cmp = Compnt
					Elout.Id = ELIOType(idmrkc[i])
					Elin.Id = ELIOType(idmrkc[i])
					Elout.Elins = (*Eli)[elinIdx : elinIdx+1]
					Elout.Ni = 1

					eloIdx++
					elinIdx++
				}
			} else {
				//Elin := (*Eli)[elinIdx]
				Elout := (*Elo)[eloIdx]
				Elout.Cmp = Compnt
				Elout.Elins = (*Eli)[elinIdx : elinIdx+1]
				Elout.Ni = 1
				eloIdx++
				elinIdx++
			}
		} else if c == STHEAT_TYPE {
			// 電気蓄熱暖房器

			Stheat = &Eqsys.Stheat[neqp]
			Compnt.Eqp = Stheat
			Stheat.Name = name
			Stheat.Cmp = Compnt
			Stheat.Cat = &Eqcat.Stheatca[ncat]
			Compnt.Airpathcpy = 'y'
			Compnt.Nin = 2
			Compnt.Nout = 2

			id = idmrkc
			for i = 0; i < Compnt.Nout; i++ {
				Elin := (*Eli)[elinIdx]
				Elout := (*Elo)[eloIdx]

				Elout.Cmp = Compnt
				Elout.Id = ELIOType(idmrkc[i])
				Elin.Id = ELIOType(idmrkc[i])
				Elout.Elins = (*Eli)[elinIdx : elinIdx+1]
				Elout.Ni = 1

				eloIdx++
				elinIdx++
			}
		} else if c == DESI_TYPE {
			// Satoh追加　デシカント槽　2013/10/23

			Desi = &Eqsys.Desi[neqp]
			Compnt.Eqp = Desi
			Desi.Name = name
			Desi.Cmp = Compnt
			Desi.Cat = &Eqcat.Desica[ncat]

			// 絶対湿度経路のコピー
			Compnt.Airpathcpy = 'y'

			id = idmrkc
			for i = 0; i < Compnt.Nout; i++ {
				Elin := (*Eli)[elinIdx]
				Elout := (*Elo)[eloIdx]

				Elout.Cmp = Compnt
				Elout.Id = ELIOType(idmrkc[i])
				Elin.Id = ELIOType(idmrkc[i])
				Elout.Elins = (*Eli)[elinIdx : elinIdx+2]

				// すべての出口状態計算のための変数は2つ（温度と湿度）
				Elout.Ni = 2

				eloIdx++
			}
		} else if c == EVAC_TYPE {
			// Satoh追加　気化冷却器 2013/10/26

			Evac = &Eqsys.Evac[neqp]
			Compnt.Eqp = Evac
			Evac.Name = name
			Evac.Cmp = Compnt
			Evac.Cat = &Eqcat.Evacca[ncat]

			// 機器の出入口数（Tdry, xdry, Twet, xwet）
			Compnt.Nout = 4
			Compnt.Nin = 4

			// 絶対湿度経路のコピー
			Compnt.Airpathcpy = 'y'

			// D:Tdry d:xdry V:Twet v:xwet
			idd := [4]ELIOType{
				ELIO_D, ELIO_d, ELIO_V, ELIO_v,
			}
			for i = 0; i < Compnt.Nout; i++ {
				Elin := (*Eli)[elinIdx]
				Elout := (*Elo)[eloIdx]

				Elout.Cmp = Compnt
				Elout.Id = idd[i]
				Elin.Id = idd[i]
				Elout.Elins = (*Eli)[elinIdx : elinIdx+4]
				// すべての出口状態計算のための変数は4つ（Wet、Dryの温度と湿度）
				Elout.Ni = 4
				// 出口状態計算のための変数分だけメモリを確保する
				elinIdx += Elout.Ni

				eloIdx++
			}
		} else if c == VALV_TYPE || c == TVALV_TYPE {
			Valv := &Eqsys.Valv[Nvalv]
			Compnt.Eqp = Valv
			Valv.Name = name
			Valv.Cmp = Compnt

			if Valv.Cmp.Valvcmp != nil {
				Valv.Cmb = Valv.Cmp.Valvcmp
			}

			Nvalv++
		} else if c == OMVAV_TYPE {
			// Satoh OMVAV  2010/12/16

			OMvav := &Eqsys.OMvav[NOMvav]
			Compnt.Eqp = OMvav
			OMvav.Name = name
			OMvav.Cmp = Compnt
			OMvav.Cat = &Eqcat.OMvavca[ncat]

			if OMvav.Cmp.Omparm != "" {
				OMvavControl(OMvav, cmp)
			}

			NOMvav++
		} else if c == QMEAS_TYPE {
			Qmeas := &Eqsys.Qmeas[NQmeas]
			Compnt.Eqp = Qmeas
			Qmeas.Name = name
			Qmeas.Cmp = Compnt

			NQmeas++
		} else if c == THEX_TYPE {
			Thex = &Eqsys.Thex[neqp]
			Compnt.Eqp = Thex
			Thex.Name = name
			Thex.Cmp = Compnt
			Thex.Cat = &Eqcat.Thexca[ncat]
			Compnt.Airpathcpy = 'y'
			Compnt.Nout = 4

			for i = 0; i < Compnt.Nout; i++ {
				Elin := (*Eli)[elinIdx]
				Elout := (*Elo)[eloIdx]

				Elout.Cmp = Compnt
				Elout.Id = Compnt.Ido[i]
				//Elin.id = Compnt.idi[i] ;

				if Thex.Cat.eh > 0.0 {
					Elout.Ni = 5
					idxe = []ELIOType{
						//"eE*Oo"
						ELIO_e, ELIO_E, ELIO_ASTER, ELIO_O, ELIO_o,
					}
					idxo = []ELIOType{
						//"oO*Ee"
						ELIO_o, ELIO_O, ELIO_ASTER, ELIO_E, ELIO_e,
					}
				} else {
					Elout.Ni = 1
					idxe = []ELIOType{ELIO_e}
					idxo = []ELIOType{ELIO_o}
				}

				//Elout.fluid = AIR_FLD ;

				if i == 0 || i == 2 { // 温度の経路（要素方程式の変数は2つ）
					Elout.Ni = 2
				}

				Elout.Elins = (*Eli)[elinIdx : elinIdx+Elout.Ni]

				for ii = 0; ii < Elout.Ni; ii++ {
					if i == 0 {
						Elin.Id = ELIOType(idTe[ii])
					} else if i == 2 {
						Elin.Id = ELIOType(idTo[ii])
					} else if i == 1 {
						Elin.Id = idxe[ii]
					} else if i == 3 {
						Elin.Id = idxo[ii]
					}

					elinIdx++
				}

				eloIdx++
			}
		} else {
			Errprint(1, errkey, string(c))
		}

		for i = 0; i < Compnt.Nout; i++ {
			elop = Compnt.Elouts[i]
			elop.Coeffin = make([]float64, elop.Ni)
		}
	}

	*Nelout = eloIdx
	*Nelin = elinIdx

	for i = 0; i < *Nelin; i++ {
		Elin := (*Eli)[i]
		Elin.Upo = nil
		Elin.Upv = nil
	}
}

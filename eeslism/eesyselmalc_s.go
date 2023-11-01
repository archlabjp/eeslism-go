package eeslism

/* 機器使用データの割り付けおよびシステム要素から入力、出力要素の割り付け */

var idmrkc = []FliudType{
	AIRt_FLD,  //'t' 空気（温度）
	AIRx_FLD,  //'x' 空気（湿度）
	WATER_FLD, //'W' 水
}

func Elmalloc(
	errkey string,
	_Compnt []*COMPNT,
	Eqcat *EQCAT,
	Eqsys *EQSYS,
	Elo *[]*ELOUT,
	Eli *[]*ELIN,
) {
	var cmp []*COMPNT
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

	var Nvalv, NQmeas, NOMvav int
	idTe := "EO"
	idTo := "OE"
	var idxe, idxo []ELIOType

	cmp = _Compnt

	*Elo = make([]*ELOUT, 0)
	*Eli = make([]*ELIN, 0)

	// eloIdx := 0
	// elinIdx := 0

	flinIdx := 0
	hcloadIdx := 0
	Cnvrg = Eqsys.Cnvrg
	Flin = Eqsys.Flin
	Hcload = Eqsys.Hcload

	for _, Compnt := range _Compnt {

		if Compnt.Eqptype != PV_TYPE {
			Compnt.Elouts = make([]*ELOUT, 0)
			Compnt.Elins = make([]*ELIN, 0)
		}

		name = Compnt.Name
		neqp = Compnt.Neqp
		ncat = Compnt.Ncat

		c := Compnt.Eqptype

		if SIMUL_BUILDG && c == ROOM_TYPE {
			room := Compnt.Eqp.(*ROOM)
			room.cmp = Compnt //逆参照の設定

			id := idmrkc
			for i := 0; i < 2; i++ {
				Elout := NewElout()
				Elout.Cmp = Compnt //逆参照の設定
				Elout.Id = ELIOType(id[i])
				Elout.Fluid = FliudType(id[i])
				if i == 0 {
					// 空気温度の流入経路の数
					Elout.Ni = room.Nachr + room.Ntr + room.Nrp + room.Nasup
				} else if i == 1 {
					// 空気湿度の流入経路の数
					Elout.Ni = room.Nachr + room.Nasup
				}
				Elout.Elins = NewElinSlice(Elout.Ni)

				*Elo = append(*Elo, Elout)
				*Eli = append(*Eli, Elout.Elins...)

				Compnt.Elouts = append(Compnt.Elouts, Elout)
				Compnt.Elins = append(Compnt.Elins, Elout.Elins...)
			}

			// 空気温度・湿度それぞれの流入経路の設定
			room.elinasup = Compnt.Elouts[0].Elins[room.Nachr+room.Ntr+room.Nrp:]
			room.elinasupx = Compnt.Elouts[1].Elins[room.Nachr:]
		} else if SIMUL_BUILDG && c == RDPANEL_TYPE {
			rdpnl := Compnt.Eqp.(*RDPNL)
			rdpnl.cmp = Compnt
			rdpnl.Tpi = 15.0

			// 空気経路温度用
			Elout_t := NewElout()
			Elout_t.Cmp = Compnt
			Elout_t.Id = ELIO_f
			Elout_t.Ni = 1 + 1 + rdpnl.Ntrm[0] + rdpnl.Nrp[0]
			if rdpnl.MC == 2 {
				// 共用壁の場合
				Elout_t.Ni += 1 + rdpnl.Ntrm[1] + rdpnl.Nrp[1]
			}
			Elout_t.Elins = make([]*ELIN, 0, Elout_t.Ni)

			Elin := NewElin()
			Elin.Id = ELIO_f
			Elout_t.Elins = append(Elout_t.Elins, Elin)

			for mm = 0; mm < rdpnl.MC; mm++ {
				Elin := NewElin()
				Elin.Id = ELIO_r
				Elout_t.Elins = append(Elout_t.Elins, Elin)

				for ii = 0; ii < rdpnl.Ntrm[mm]; ii++ {
					Elin := NewElin()
					Elin.Id = ELIO_r
					Elout_t.Elins = append(Elout_t.Elins, Elin)
				}
				for ii = 0; ii < rdpnl.Nrp[mm]; ii++ {
					Elin := NewElin()
					Elin.Id = ELIO_f
					Elout_t.Elins = append(Elout_t.Elins, Elin)
				}
			}

			*Elo = append(*Elo, Elout_t)
			*Eli = append(*Eli, Elout_t.Elins...)

			// 空気経路湿度用
			Elout_x := NewElout()
			Elout_x.Cmp = Compnt
			Elout_x.Id = ELIO_x
			Elout_x.Ni = 1
			Elout_x.Elins = make([]*ELIN, 0, Elout_x.Ni)
			Elin.Id = ELIO_x

			*Elo = append(*Elo, Elout_x)
			*Eli = append(*Eli, Elout_x.Elins...)

			Compnt.Elouts = append(Compnt.Elouts, Elout_t, Elout_x)
			Compnt.Elins = append(Compnt.Elins, Elout_t.Elins...)
			Compnt.Elins = append(Compnt.Elins, Elout_x.Elins...)
		} else if c == DIVERG_TYPE || c == DIVGAIR_TYPE {
			if c == DIVGAIR_TYPE {
				Compnt.Airpathcpy = true
			} else {
				Compnt.Airpathcpy = false
			}

			Compnt.Nin = 1

			// 分岐なので、入口は共通で1つ
			Elin := NewElin()
			Elin.Id = ELIO_i
			*Eli = append(*Eli, Elin)
			Compnt.Elins = append(Compnt.Elins, Elin)

			// 出口は複数
			for i := 0; i < Compnt.Nout; i++ {
				Elout := NewElout()
				Elout.Cmp = Compnt
				Elout.Ni = Compnt.Nin
				Elout.Elins = []*ELIN{Elin} //共通
				Elout.Id = Compnt.Ido[i]
				*Elo = append(*Elo, Elout)
				Compnt.Elouts = append(Compnt.Elouts, Elout)
			}
		} else if c == CONVRG_TYPE || c == CVRGAIR_TYPE {
			if c == CVRGAIR_TYPE {
				Compnt.Airpathcpy = true
			} else {
				Compnt.Airpathcpy = false
			}

			Compnt.Nout = 1

			Cnvrg[icv] = Compnt
			icv++

			// 合流なので、出口は1つ
			Elout := NewElout()
			Elout.Id = ELIO_o
			Elout.Cmp = Compnt
			Elout.Ni = Compnt.Nin
			Elout.Elins = make([]*ELIN, 0, Elout.Ni)

			// 入口は複数
			for i := 0; i < Compnt.Nin; i++ {
				Elin := NewElin()
				Elin.Id = Compnt.Idi[i]
				Elout.Elins = append(Elout.Elins, Elin)
			}

			*Elo = append(*Elo, Elout)
			*Eli = append(*Eli, Elout.Elins...)
			Compnt.Elouts = append(Compnt.Elouts, Elout)
			Compnt.Elins = append(Compnt.Elins, Elout.Elins...)
		} else if c == HCCOIL_TYPE {
			Hcc = &Eqsys.Hcc[neqp]
			Compnt.Eqp = Hcc
			Hcc.Name = name
			Hcc.Cmp = Compnt
			Hcc.Cat = &Eqcat.Hccca[ncat]

			// 入口の数=出口の数
			for i = 0; i < Compnt.Nout; i++ {
				Elout := NewElout()
				Elout.Cmp = Compnt
				Elout.Ni = Compnt.Nin
				Elout.Elins = Compnt.Elins
				Elout.Id = Compnt.Ido[i]

				Elin := NewElin()
				Elin.Id = Compnt.Idi[i]

				*Elo = append(*Elo, Elout)
				*Eli = append(*Eli, Elout.Elins...)
				Compnt.Elouts = append(Compnt.Elouts, Elout)
				Compnt.Elins = append(Compnt.Elins, Elout.Elins...)
			}
		} else if c == HEXCHANGR_TYPE {
			Hex := &Eqsys.Hex[neqp]
			Compnt.Eqp = Hex
			Hex.Name = name
			Hex.Cmp = Compnt
			Hex.Cat = &Eqcat.Hexca[ncat]

			for i = 0; i < Compnt.Nout; i++ {
				Elout := NewElout()
				Elout.Cmp = Compnt
				Elout.Ni = Compnt.Nin
				Elout.Elins = Compnt.Elins
				Elout.Id = Compnt.Ido[i]

				Elin := NewElin()
				Elin.Id = Compnt.Idi[i]

				*Elo = append(*Elo, Elout)
				*Eli = append(*Eli, Elout.Elins...)
				Compnt.Elouts = append(Compnt.Elouts, Elout)
				Compnt.Elins = append(Compnt.Elins, Elout.Elins...)
			}
		} else if c == BOILER_TYPE {
			Boi = &Eqsys.Boi[neqp]
			Compnt.Eqp = Boi
			Boi.Name = name
			Boi.Cmp = Compnt
			Boi.Cat = &Eqcat.Boica[ncat]

			Elout := NewElout()
			Elout.Cmp = Compnt
			Elout.Ni = Compnt.Nin
			Elout.Elins = NewElinSlice(Elout.Ni)

			*Elo = append(*Elo, Elout)
			*Eli = append(*Eli, Elout.Elins...)
			Compnt.Elouts = append(Compnt.Elouts, Elout)
			Compnt.Elins = append(Compnt.Elins, Elout.Elins...)
		} else if c == COLLECTOR_TYPE || c == ACOLLECTOR_TYPE {
			Coll = &Eqsys.Coll[neqp]
			Compnt.Eqp = Coll
			Coll.Name = name
			Coll.Cmp = Compnt
			Coll.Cat = &Eqcat.Collca[ncat]
			Coll.Ac = Compnt.Ac

			if Coll.Cat.Type == COLLECTOR_PDT {
				Elout := NewElout()
				Elout.Cmp = Compnt

				Elout.Ni = 1
				Elout.Elins = NewElinSlice(Elout.Ni)

				*Elo = append(*Elo, Elout)
				*Eli = append(*Eli, Elout.Elins...)
				Compnt.Elouts = append(Compnt.Elouts, Elout)
				Compnt.Elins = append(Compnt.Elins, Elout.Elins...)
			} else {
				for i = 0; i < Compnt.Nout; i++ {
					Elin := NewElin()
					Elout := NewElout()

					Elout.Cmp = Compnt
					Elout.Id = ELIOType(idmrkc[i])
					Elin.Id = ELIOType(idmrkc[i])
					Elout.Ni = 1
					Elout.Elins = []*ELIN{Elin}

					*Elo = append(*Elo, Elout)
					*Eli = append(*Eli, Elout.Elins...)
					Compnt.Elouts = append(Compnt.Elouts, Elout)
					Compnt.Elins = append(Compnt.Elins, Elout.Elins...)
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

			Elout := NewElout()
			Elout.Cmp = Compnt
			Elout.Ni = Compnt.Nin
			Elout.Elins = NewElinSlice(Elout.Ni)

			*Elo = append(*Elo, Elout)
			*Eli = append(*Eli, Elout.Elins...)
			Compnt.Elouts = append(Compnt.Elouts, Elout)
			Compnt.Elins = append(Compnt.Elins, Elout.Elins...)
		} else if c == PUMP_TYPE {
			Pump = &Eqsys.Pump[neqp]
			Compnt.Eqp = Pump
			Pump.Name = name
			Pump.Cmp = Compnt
			Pump.Cat = &Eqcat.Pumpca[ncat]

			if Pump.Cat.pftype == PUMP_PF {
				Elout := NewElout()
				Elout.Cmp = Compnt
				Elout.Elins = NewElinSlice(1)
				Elout.Ni = 1
				*Elo = append(*Elo, Elout)
				*Eli = append(*Eli, Elout.Elins...)
				Compnt.Elouts = append(Compnt.Elouts, Elout)
				Compnt.Elins = append(Compnt.Elins, Elout.Elins...)
			} else {
				for i = 0; i < Compnt.Nout; i++ {
					Elin := NewElin()
					Elout := NewElout()

					Elout.Cmp = Compnt
					Elout.Id = ELIOType(idmrkc[i])
					Elin.Id = ELIOType(idmrkc[i])
					Elout.Elins = NewElinSlice(1)
					Elout.Ni = 1

					*Elo = append(*Elo, Elout)
					*Eli = append(*Eli, Elout.Elins...)
					Compnt.Elouts = append(Compnt.Elouts, Elout)
					Compnt.Elins = append(Compnt.Elins, Elout.Elins...)
				}
			}
		} else if c == PIPEDUCT_TYPE {
			Pipe = &Eqsys.Pipe[neqp]
			Compnt.Eqp = Pipe
			Pipe.Name = name
			Pipe.Cmp = Compnt
			Pipe.Cat = &Eqcat.Pipeca[ncat]

			if Pipe.Cat.Type == PIPE_PDT {
				Elout := NewElout()

				Elout.Cmp = Compnt
				Elout.Elins = NewElinSlice(1)
				Elout.Ni = 1

				*Elo = append(*Elo, Elout)
				*Eli = append(*Eli, Elout.Elins...)
				Compnt.Elouts = append(Compnt.Elouts, Elout)
				Compnt.Elins = append(Compnt.Elins, Elout.Elins...)
			} else {
				for i = 0; i < Compnt.Nout; i++ {
					Elin := NewElin()
					Elout := NewElout()

					Elout.Cmp = Compnt
					Elout.Id = ELIOType(idmrkc[i])
					Elin.Id = ELIOType(idmrkc[i])
					Elout.Elins = []*ELIN{Elin}
					Elout.Ni = 1

					*Elo = append(*Elo, Elout)
					*Eli = append(*Eli, Elout.Elins...)
					Compnt.Elouts = append(Compnt.Elouts, Elout)
					Compnt.Elins = append(Compnt.Elins, Elout.Elins...)
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

			// 入力は共通
			Elins := NewElinSlice(Compnt.Nin)
			for i := 0; i < Compnt.Nin; i++ {
				Elins[i].Id = Stank.Pthcon[i]
				Compnt.Idi[i] = Stank.Pthcon[i]
			}
			*Eli = append(*Eli, Elins...)
			Compnt.Elins = append(Compnt.Elins, Elins...)

			for i = 0; i < Compnt.Nout; i++ {
				Elout := NewElout()

				Elout.Cmp = Compnt
				Elout.Ni = Compnt.Nin
				Elout.Elins = Elins
				Elout.Id = Stank.Pthcon[i]
				Compnt.Ido[i] = Stank.Pthcon[i]

				*Elo = append(*Elo, Elout)
				Compnt.Elouts = append(Compnt.Elouts, Elout)
			}
		} else if c == FLIN_TYPE {
			// 流入境界条件
			Compnt.Eqp = &Flin[flinIdx]
			Flin[flinIdx].Cmp = Compnt
			Flin[flinIdx].Name = name

			flindat(&Flin[flinIdx])

			for i = 0; i < Compnt.Nout; i++ {
				Elout := NewElout()
				Elout.Cmp = Compnt
				Elout.Ni = Compnt.Nin
				Elout.Elins = NewElinSlice(Elout.Ni)
				Elout.Id = Compnt.Ido[i]

				*Elo = append(*Elo, Elout)
				*Eli = append(*Eli, Elout.Elins...)
				Compnt.Elouts = append(Compnt.Elouts, Elout)
				Compnt.Elins = append(Compnt.Elins, Elout.Elins...)
			}

			flinIdx++
		} else if c == HCLOAD_TYPE ||
			c == HCLOADW_TYPE ||
			c == RMAC_TYPE ||
			c == RMACD_TYPE {
			// 仮想空調機

			Compnt.Eqp = &Hcload[hcloadIdx]
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
				Hcload[hcloadIdx].Wet = true // 湿りコイル（吹出相対湿度一定
			} else {
				Hcload[hcloadIdx].Wet = false // 吹出相対湿度は成り行き
			}

			// 空気のみの流入、流出
			if c == HCLOAD_TYPE || c == RMAC_TYPE || c == RMACD_TYPE {
				Hcload[hcloadIdx].Type = HCLoadType_D
				Compnt.Nout = 2
				Compnt.Nin = 2
			} else {
				// 空気＋水の流入、流出
				Hcload[hcloadIdx].Type = HCLoadType_W
				Compnt.Nout = 3
				Compnt.Nin = 3
			}

			// 空気の絶対湿度用経路コピーを行う
			Compnt.Airpathcpy = true
			for i := 0; i < Compnt.Nout; i++ {
				Elout := NewElout()

				Elout.Cmp = Compnt
				Elout.Id = ELIOType(idmrkc[i])

				Elout.Ni = 1
				//Elout.Ni = 2;
				// 湿りコイル（吹出相対湿度一定）で出口絶対湿度の経路の場合
				// 要素方程式の未知数は2つ（入口絶対湿度と出口温度）
				if i == 1 && Hcload[hcloadIdx].Wet {
					Elout.Ni = 2
				} else if i == 2 && Hcload[hcloadIdx].Type == 'W' {
					// 冷温水コイルで水側系統の場合
					// 要素方程式の未知数は5個
					// 水入口温度、空気入口温度、空気入口湿度
					// 空気出口温度、空気出口湿度
					Elout.Ni = 5
				}
				Elout.Elins = NewElinSlice(Elout.Ni)

				for ii := 0; ii < Elout.Ni; ii++ {
					Elout.Elins[ii].Id = Elout.Id

					// 空気出口絶対湿度の計算の2つ目の変数は空気出口温度
					if i == 1 && ii == 1 {
						Elout.Elins[ii].Id = ELIO_ASTER
					}
				}
				*Eli = append(*Eli, Elout.Elins...)

				/***** printf("xxx Elmalloc xxx   %s  i=%d  Elout.Ni=%d\n",
				Hcload.name, i, Elout.Ni); *****/

				*Elo = append(*Elo, Elout)

				Compnt.Elouts = append(Compnt.Elouts, Elout)
				Compnt.Elins = append(Compnt.Elins, Elout.Elins...)
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
				Compnt.Airpathcpy = true
				for i = 0; i < Compnt.Nout; i++ {
					Elin := NewElin()
					Elout := NewElout()

					Elout.Cmp = Compnt
					Elout.Id = ELIOType(idmrkc[i])
					Elin.Id = ELIOType(idmrkc[i])
					Elout.Elins = []*ELIN{Elin}
					Elout.Ni = 1

					*Elo = append(*Elo, Elout)
					*Eli = append(*Eli, Elout.Elins...)
					Compnt.Elouts = append(Compnt.Elouts, Elout)
					Compnt.Elins = append(Compnt.Elins, Elout.Elins...)
				}
			} else {
				Elin := NewElin()
				Elout := NewElout()
				Elout.Cmp = Compnt
				Elout.Elins = []*ELIN{Elin}
				Elout.Ni = 1

				*Elo = append(*Elo, Elout)
				*Eli = append(*Eli, Elout.Elins...)
				Compnt.Elouts = append(Compnt.Elouts, Elout)
				Compnt.Elins = append(Compnt.Elins, Elout.Elins...)
			}
		} else if c == STHEAT_TYPE {
			// 電気蓄熱暖房器

			Stheat = &Eqsys.Stheat[neqp]
			Compnt.Eqp = Stheat
			Stheat.Name = name
			Stheat.Cmp = Compnt
			Stheat.Cat = &Eqcat.Stheatca[ncat]
			Compnt.Airpathcpy = true
			Compnt.Nin = 2
			Compnt.Nout = 2

			for i = 0; i < Compnt.Nout; i++ {
				Elin := NewElin()
				Elout := NewElout()

				Elout.Cmp = Compnt
				Elout.Id = ELIOType(idmrkc[i])
				Elin.Id = ELIOType(idmrkc[i])
				Elout.Elins = []*ELIN{Elin}
				Elout.Ni = 1

				*Elo = append(*Elo, Elout)
				*Eli = append(*Eli, Elout.Elins...)
				Compnt.Elouts = append(Compnt.Elouts, Elout)
				Compnt.Elins = append(Compnt.Elins, Elout.Elins...)
			}
		} else if c == DESI_TYPE {
			// Satoh追加　デシカント槽　2013/10/23

			Desi = &Eqsys.Desi[neqp]
			Compnt.Eqp = Desi
			Desi.Name = name
			Desi.Cmp = Compnt
			Desi.Cat = &Eqcat.Desica[ncat]

			// 絶対湿度経路のコピー
			Compnt.Airpathcpy = true

			for i := 0; i < Compnt.Nout; i++ {
				Elout := NewElout()

				Elout.Cmp = Compnt
				Elout.Id = ELIOType(idmrkc[i])
				Elout.Elins = NewElinSlice(2)
				Elout.Elins[0].Id = ELIOType(idmrkc[i])

				// すべての出口状態計算のための変数は2つ（温度と湿度）
				Elout.Ni = 2

				*Elo = append(*Elo, Elout)
				*Eli = append(*Eli, Elout.Elins...)
				Compnt.Elouts = append(Compnt.Elouts, Elout)
				Compnt.Elins = append(Compnt.Elins, Elout.Elins...)
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
			Compnt.Airpathcpy = true

			// D:Tdry d:xdry V:Twet v:xwet
			idd := [4]ELIOType{
				ELIO_D, ELIO_d, ELIO_V, ELIO_v,
			}
			for i = 0; i < Compnt.Nout; i++ {
				Elout := NewElout()
				Elout.Cmp = Compnt
				Elout.Id = idd[i]
				Elout.Elins = NewElinSlice(4)
				// すべての出口状態計算のための変数は4つ（Wet、Dryの温度と湿度）
				Elout.Ni = 4
				// 出口状態計算のための変数分だけメモリを確保する
				Elout.Elins[0].Id = idd[i]

				*Elo = append(*Elo, Elout)
				*Eli = append(*Eli, Elout.Elins...)
				Compnt.Elouts = append(Compnt.Elouts, Elout)
				Compnt.Elins = append(Compnt.Elins, Elout.Elins...)
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
			Compnt.Airpathcpy = true
			Compnt.Nout = 4

			for i = 0; i < Compnt.Nout; i++ {
				Elout := NewElout()

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

				Elout.Elins = NewElinSlice(Elout.Ni)

				for ii = 0; ii < Elout.Ni; ii++ {
					if i == 0 {
						Elout.Elins[ii].Id = ELIOType(idTe[ii])
					} else if i == 2 {
						Elout.Elins[ii].Id = ELIOType(idTo[ii])
					} else if i == 1 {
						Elout.Elins[ii].Id = idxe[ii]
					} else if i == 3 {
						Elout.Elins[ii].Id = idxo[ii]
					}
				}

				*Elo = append(*Elo, Elout)
				*Eli = append(*Eli, Elout.Elins...)
				Compnt.Elouts = append(Compnt.Elouts, Elout)
				Compnt.Elins = append(Compnt.Elins, Elout.Elins...)
			}
		} else {
			Errprint(1, errkey, string(c))
		}

		for i = 0; i < Compnt.Nout; i++ {
			elop := Compnt.Elouts[i]
			elop.Coeffin = make([]float64, elop.Ni)
		}
	}

	// 上流の機器の出口の参照をクリアしておく
	for i := range *Eli {
		Elin := (*Eli)[i]
		Elin.Upo = nil
		Elin.Upv = nil
	}
}

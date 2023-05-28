//This file is part of EESLISM.
//
//Foobar is free software : you can redistribute itand /or modify
//it under the terms of the GNU General Public License as published by
//the Free Software Foundation, either version 3 of the License, or
//(at your option) any later version.
//
//Foobar is distributed in the hope that it will be useful,
//but WITHOUT ANY WARRANTY; without even the implied warranty of
//MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.See the
//GNU General Public License for more details.
//
//You should have received a copy of the GNU General Public License
//along with Foobar.If not, see < https://www.gnu.org/licenses/>.

/* escntllb_s.c */

package main

import (
	"strings"
)

/*  システム変数名、内部変数名、スケジュール名のポインター  */

func ctlvptr(s string, Simc *SIMCONTL, Ncompnt int, Compnt []COMPNT, Nmpath int, Mpath []MPATH, Wd *WDAT, Exsf *EXSFS, Schdl *SCHDL, vptr *VPTR, vpath *VPTR) int {
	var err int

	if i := idsch(s, Schdl.Sch, ""); i >= 0 {
		vptr.Ptr = &Schdl.Val[i]
		vptr.Type = VAL_CTYPE
	} else if i := idscw(s, Schdl.Scw, ""); i >= 0 {
		vptr.Ptr = &Schdl.Isw[i]
		vptr.Type = SW_CTYPE
	} else {
		err = kynameptr(s, Simc, Ncompnt, Compnt, Nmpath, Mpath, Wd, Exsf, vptr, vpath)
	}

	Errprint(err, "<ctlvptr>", s)
	return err
}

/* ----------------------------------------------------------------- */

/* システム経路名、要素名、システム変数名、内部変数名の分離 */

func strkey(s string) ([]string, int) {
	if len(s) == 0 {
		return nil, 0
	}

	key := strings.Split(s, "_")
	return key, len(key)
}

/* ----------------------------------------------------------------- */

/*  経路名、システム変数名、内部変数名のポインター  */

func kynameptr(s string, Simc *SIMCONTL, Ncompnt int, _Compnt []COMPNT,
	Nmpath int, Mpath []MPATH, Wd *WDAT, Exsf *EXSFS, vptr *VPTR, vpath *VPTR) int {
	var err int

	key := strings.Split(s, "_")
	nk := len(key)

	switch key[0] {
	case "Ta":
		vptr.Type = VAL_CTYPE
		vptr.Ptr = &Wd.T
	case "xa":
		vptr.Type = VAL_CTYPE
		vptr.Ptr = &Wd.X
	case "RHa":
		vptr.Type = VAL_CTYPE
		vptr.Ptr = &Wd.RH
	case "ha":
		vptr.Type = VAL_CTYPE
		vptr.Ptr = &Wd.H
	case "Twsup":
		vptr.Type = VAL_CTYPE
		vptr.Ptr = &Wd.Twsup
	case "Ihol":
		vptr.Type = VAL_CTYPE
		vptr.Ptr = &Wd.Ihor
	default:
		// // 傾斜面名称の検索
		if Exsf != nil {
			for i := 0; i < Exsf.Nexs; i++ {
				Exs := &Exsf.Exs[i]
				if key[0] == Exs.Name {
					switch key[1] {
					case "Idre":
						// 傾斜面への入射直達日射量
						vptr.Type = VAL_CTYPE
						vptr.Ptr = &Exs.Idre
						return 0
					case "Idf":
						// 傾斜面への入射拡散日射量
						vptr.Type = VAL_CTYPE
						vptr.Ptr = &Exs.Idf
						return 0
					case "Iw":
						// 傾斜面への入射全日射量
						vptr.Type = VAL_CTYPE
						vptr.Ptr = &Exs.Iw
						return 0
					}
				}
			}
		}

		if Nmpath > 0 {
			err = pathvptr(nk, key, Nmpath, Mpath, vptr, vpath)
		} else {
			err = 1
		}

		if err != 0 {
			if Simc.Nvcfile > 0 {
				err = vcfptr(key, Simc, vptr)
			} else {
				err = 1
			}
		}

		if err != 0 {
			for i := 0; i < Ncompnt; i++ {
				Compnt := &_Compnt[i]
				if key[0] == Compnt.Name {
					err = compntvptr(nk, key, Compnt, vptr)
					if err != 0 {
						e := Compnt.Eqptype
						switch e {
						case ROOM_TYPE:
							if SIMUL_BUILDG {
								err = roomvptr(nk, key, Compnt.Eqp.(*ROOM), vptr)
							}
						case REFACOMP_TYPE:
							err = refaswptr(key, Compnt.Eqp.(*REFA), vptr)
						case HCLOAD_TYPE, HCLOADW_TYPE, RMAC_TYPE, RMACD_TYPE:
							err = hcldswptr(key, Compnt.Eqp.(*HCLOAD), vptr)
						case VAV_TYPE, VWV_TYPE:
							/* VAV Satoh Debug 2001/1/19 */
							err = vavswptr(key, Compnt.Eqp.(*VAV), vptr)
						case COLLECTOR_TYPE:
							err = collvptr(key, Compnt.Eqp.(*COLL), vptr)
						case STANK_TYPE:
							err = stankvptr(key, Compnt.Eqp.(*STANK), vptr)
						case STHEAT_TYPE:
							err = stheatvptr(key, Compnt.Eqp.(*STHEAT), vptr, vpath)
						case DESI_TYPE:
							// Satoh追加　デシカント槽　2013/10/23
							err = Desivptr(key, Compnt.Eqp.(*DESI), vptr)
						case PIPEDUCT_TYPE:
							err = pipevptr(key, Compnt.Eqp.(*PIPE), vptr)
						case RDPANEL_TYPE:
							err = rdpnlvptr(key, Compnt.Eqp.(*RDPNL), vptr)
						case VALV_TYPE, TVALV_TYPE:
							err = valv_vptr(key, Compnt.Eqp.(*VALV), vptr)
						default:
							Eprint("CONTL", Compnt.Name)
						}
					}
					break
				}
			}
		}
	}

	Errprint(err, "<kynameptr>", s)

	return err
}

/* ----------------------------------------------------------------- */

/*  経路名のポインター  */

func pathvptr(nk int, key []string, Nmpath int, Mpath []MPATH, vptr *VPTR, vpath *VPTR) int {
	var i, err int
	var Mp, Mpe *MPATH
	var Plist, Plie *PLIST

	Mp = &Mpath[0]

	for i = 0; i < Nmpath; i++ {
		if string(key[0]) == Mpath[i].Name {
			vpath.Type = MAIN_CPTYPE
			vpath.Ptr = Mpath[i]

			if nk == 1 || string(key[1]) == "control" {
				vptr.Type = SW_CTYPE
				vptr.Ptr = &Mpath[i].Control
			}
			break
		}
	}

	err = 0
	if i == Nmpath {
		err = 1
		Mpe = &Mpath[Nmpath-1]

		for j := 0; j < Mpe.Nlpath; j++ {
			Plist = &Mp.Plist[j]
			if Plist.Name != "" {
				if key[0] == Plist.Name {
					vpath.Type = LOCAL_CPTYPE
					vpath.Ptr = Plist

					if nk == 1 || key[1] == "control" {
						vptr.Type = SW_CTYPE
						vptr.Ptr = &Plist.Control
					} else if key[1] == "G" {
						vptr.Type = VAL_CTYPE
						vptr.Ptr = &Plist.G
					}
					break
				}
			}
		}
		if Plist == Plie {
			err = 1
		}
	}
	return err
}

func Compntptr(name string, N int, Compnt []COMPNT) *COMPNT {
	for i := 0; i < N; i++ {
		if name == Compnt[i].Name {
			return &Compnt[i]
		}
	}

	return nil
}

/* ----------------------------------------------------------------- */

/*  システム要素出口温度、湿度のポインター  */

func compntvptr(nk int, key []string, Compnt *COMPNT, vptr *VPTR) int {
	var i, err int

	if nk == 1 || key[1] == "control" {
		etype := Compnt.Eqptype
		if etype != VALV_TYPE && etype != TVALV_TYPE {
			// ボイラなど機器自体の停止ではなく、燃焼の停止とする
			Eo := Compnt.Elouts[0]
			if etype == STHEAT_TYPE {
				vptr.Ptr = &Compnt.Control
			} else {
				vptr.Ptr = &Eo.Control
			}
			vptr.Type = SW_CTYPE
		} else {
			v := Compnt.Eqp.(*VALV)
			vptr.Ptr = &v.X
			vptr.Type = VAL_CTYPE
			v.Org = 'y'
		}
	} else {
		for i = 0; i < Compnt.Nout; i++ {
			Eo := Compnt.Elouts[i]
			if (Eo.Fluid == AIRa_FLD && string(key[1]) == "Taout") ||
				(Eo.Fluid == AIRx_FLD && string(key[1]) == "xout") ||
				(Eo.Fluid == WATER_FLD && string(key[1]) == "Twout") {
				vptr.Ptr = &Eo.Sysv
				vptr.Type = VAL_CTYPE
				break
			}
		}
		if i == Compnt.Nout {
			err = 1
		}
	}
	return err
}

/* ----------------------------------------------------------------- */

/* 負荷計算を行うシステム要素の設定システム変数のポインター */

func loadptr(loadcmp *COMPNT, load *rune, s string, Ncompnt int, _Compnt []COMPNT, vptr *VPTR) int {
	var Room *ROOM
	var key []string
	var idmrk byte = ' '
	var err int

	key = strings.Split(s, "_")
	nk := len(key)

	if nk != 0 {
		for i := 0; i < Ncompnt; i++ {
			Compnt := &_Compnt[i]
			if key[0] == Compnt.Name {
				switch Compnt.Eqptype {
				case BOILER_TYPE:
					err = boildptr(load, key, Compnt.Eqp.(*BOI), vptr)
					idmrk = 't'
				case REFACOMP_TYPE:
					err = refaldptr(load, key, Compnt.Eqp.(*REFA), vptr)
					idmrk = 't'
				case HCLOAD_TYPE, RMAC_TYPE, RMACD_TYPE:
					if SIMUL_BUILDG {
						err = hcldptr(load, key, Compnt.Eqp.(*HCLOAD), vptr, &idmrk)
					}
				case PIPEDUCT_TYPE:
					if SIMUL_BUILDG {
						err = pipeldsptr(load, key, Compnt.Eqp.(*PIPE), vptr, &idmrk)
					}
				case RDPANEL_TYPE:
					if SIMUL_BUILDG {
						Rdpnl := Compnt.Eqp.(*RDPNL)
						err = rdpnlldsptr(load, key, Rdpnl, vptr, &idmrk)
						if loadcmp != nil && loadcmp.Eqptype == OMVAV_TYPE {
							Rdpnl.OMvav = loadcmp.Eqp.(*OMVAV)
							Rdpnl.OMvav.Omwall = Rdpnl.sd[0]
						}
					}
				case ROOM_TYPE:
					if SIMUL_BUILDG {
						Room = Compnt.Eqp.(*ROOM)
						if Room.rmld == nil {
							Room.rmld = new(RMLOAD)
							if Room.rmld == nil {
								Ercalloc(1, "roomldptr")
							}

							key = strings.Split(s, "_")

							R := Room.rmld
							R.loadt = nil
							R.loadx = nil
							R.FOTN = nil
							R.FOPL = nil

							if loadcmp != nil && loadcmp.Eqptype == VAV_TYPE {
								Room.VAVcontrl = loadcmp.Eqp.(*VAV)
							}
						}
						err = roomldptr(load, key, Room, vptr, &idmrk)
					}
				}

				if err == 0 {
					if loadcmp == nil {
						loadcmp = Compnt
					}

					Eo := Compnt.Elouts[0]
					eold := loadcmp.Elouts[0]

					if idmrk == 'x' {
						Eo = Compnt.Elouts[1]
						eold = loadcmp.Elouts[1]
					}

					Eo.Eldobj = eold
					eold.Emonitr = Eo

					break
				}
			} else {
				err = 1
			}
		}
		return err
	} else {
		return 1
	}
}

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

package eeslism

import (
	"errors"
	"strings"
)

/*  システム変数名、内部変数名、スケジュール名のポインター  */

func ctlvptr(s string, Simc *SIMCONTL, Compnt []*COMPNT, Mpath []*MPATH, Wd *WDAT, Exsf *EXSFS, Schdl *SCHDL) (VPTR, VPTR, error) {
	var err error
	var vptr, vpath VPTR

	if i, err2 := idsch(s, Schdl.Sch, ""); err2 == nil {
		// 年間の設定値スケジュールへのポインターを作成する
		vptr = VPTR{
			Ptr:  &Schdl.Val[i],
			Type: VAL_CTYPE,
		}
	} else if i, iderr := idscw(s, Schdl.Scw, ""); iderr == nil {
		// 年間の切替スケジュールへのポインターを作成する
		vptr = VPTR{
			Ptr:  &Schdl.Isw[i],
			Type: SW_CTYPE,
		}
	} else {
		// 経路名、システム変数名、内部変数名のポインターを作成する
		vptr, vpath, err = kynameptr(s, Simc, Compnt, Mpath, Wd, Exsf)
	}

	//Errprint(1, "<ctlvptr>", s)
	return vptr, vpath, err
}

/* ----------------------------------------------------------------- */

// システム経路名、要素名、システム変数名、内部変数名の分離
func strkey(s string) ([]string, int) {
	if len(s) == 0 {
		return nil, 0
	}

	key := strings.Split(s, "_")
	return key, len(key)
}

/* ----------------------------------------------------------------- */

// 経路名、システム変数名、内部変数名のポインターを作成する
func kynameptr(s string, Simc *SIMCONTL, _Compnt []*COMPNT,
	Mpath []*MPATH, Wd *WDAT, Exsf *EXSFS) (VPTR, VPTR, error) {
	var err error
	var vptr, vpath VPTR

	key := strings.Split(s, "_")
	nk := len(key)

	if nk > 0 {
		switch key[0] {
		case "Ta":
			vptr = VPTR{
				Type: VAL_CTYPE,
				Ptr:  &Wd.T,
			}
		case "xa":
			vptr = VPTR{
				Type: VAL_CTYPE,
				Ptr:  &Wd.X,
			}
		case "RHa":
			vptr = VPTR{
				Type: VAL_CTYPE,
				Ptr:  &Wd.RH,
			}
		case "ha":
			vptr = VPTR{
				Type: VAL_CTYPE,
				Ptr:  &Wd.H,
			}
		case "Twsup":
			vptr = VPTR{
				Type: VAL_CTYPE,
				Ptr:  &Wd.Twsup,
			}
		case "Ihol":
			vptr = VPTR{
				Type: VAL_CTYPE,
				Ptr:  &Wd.Ihor,
			}
		default:
			// // 傾斜面名称の検索
			if Exsf != nil {
				for _, Exs := range Exsf.Exs {
					if key[0] == Exs.Name {
						switch key[1] {
						case "Idre":
							// 傾斜面への入射直達日射量
							vptr = VPTR{
								Type: VAL_CTYPE,
								Ptr:  &Exs.Idre,
							}
							return vptr, vpath, nil
						case "Idf":
							// 傾斜面への入射拡散日射量
							vptr = VPTR{
								Type: VAL_CTYPE,
								Ptr:  &Exs.Idf,
							}
							return vptr, vpath, nil
						case "Iw":
							// 傾斜面への入射全日射量
							vptr = VPTR{
								Type: VAL_CTYPE,
								Ptr:  &Exs.Iw,
							}
							return vptr, vpath, nil
						}
					}
				}
			}

			if len(Mpath) > 0 {
				vptr, vpath, err = pathvptr(nk, key, Mpath)
			} else {
				err = errors.New("Nmpath == 0")
			}

			if err != nil {
				if Simc.Nvcfile > 0 {
					vptr, err = vcfptr(key, Simc)
				} else {
					err = errors.New("Simc.Nvcfile == 0")
				}
			}

			if err != nil {
				for i := range _Compnt {
					Compnt := _Compnt[i]
					if key[0] == Compnt.Name {
						vptr, err = compntvptr(nk, key, Compnt)
						if err != nil {
							e := Compnt.Eqptype
							switch e {
							case ROOM_TYPE:
								if SIMUL_BUILDG {
									vptr, err = roomvptr(nk, key, Compnt.Eqp.(*ROOM))
								}
							case REFACOMP_TYPE:
								vptr, err = refaswptr(key, Compnt.Eqp.(*REFA))
							case HCLOAD_TYPE, HCLOADW_TYPE, RMAC_TYPE, RMACD_TYPE:
								vptr, err = hcldswptr(key, Compnt.Eqp.(*HCLOAD))
							case VAV_TYPE, VWV_TYPE:
								/* VAV Satoh Debug 2001/1/19 */
								vptr, err = vavswptr(key, Compnt.Eqp.(*VAV))
							case COLLECTOR_TYPE:
								vptr, err = collvptr(key, Compnt.Eqp.(*COLL))
							case STANK_TYPE:
								vptr, err = stankvptr(key, Compnt.Eqp.(*STANK))
							case STHEAT_TYPE:
								vptr, vpath, err = stheatvptr(key, Compnt.Eqp.(*STHEAT))
							case DESI_TYPE:
								// Satoh追加　デシカント槽　2013/10/23
								vptr, err = Desivptr(key, Compnt.Eqp.(*DESI))
							case PIPEDUCT_TYPE:
								vptr, err = pipevptr(key, Compnt.Eqp.(*PIPE))
							case RDPANEL_TYPE:
								vptr, err = rdpnlvptr(key, Compnt.Eqp.(*RDPNL))
							case VALV_TYPE, TVALV_TYPE:
								vptr, err = valv_vptr(key, Compnt.Eqp.(*VALV))
							default:
								Eprint("CONTL", Compnt.Name)
							}
						}
						break
					}
				}
			}
		}
	} else {
		err = errors.New("Some error")
	}

	if err != nil {
		Eprint("<kynameptr>", s)
	}

	return vptr, vpath, err
}

/* ----------------------------------------------------------------- */

// 経路名のポインター
func pathvptr(nk int, key []string, Mpath []*MPATH) (VPTR, VPTR, error) {
	var err error
	var Plist, Plie *PLIST
	var vptr, vpath VPTR

	found := false
	for _, Mp := range Mpath {
		if string(key[0]) == Mp.Name {
			vpath = VPTR{
				Type: MAIN_CPTYPE,
				Ptr:  Mp,
			}

			if nk == 1 || string(key[1]) == "control" {
				vptr = VPTR{
					Type: SW_CTYPE,
					Ptr:  &Mp.Control,
				}
			}
			found = true
			break
		}
	}

	if found == false {
		err = errors.New("i == Nmpath")
		Mpe := Mpath[len(Mpath)-1]

		for _, Plist := range Mpe.Plist {
			if Plist.Name != "" {
				if key[0] == Plist.Name {
					vpath = VPTR{
						Type: LOCAL_CPTYPE,
						Ptr:  Plist,
					}

					if nk == 1 || key[1] == "control" {
						vptr = VPTR{
							Type: SW_CTYPE,
							Ptr:  &Plist.Control,
						}
					} else if key[1] == "G" {
						vptr = VPTR{
							Type: VAL_CTYPE,
							Ptr:  &Plist.G,
						}
					}
					break
				}
			}
		}
		if Plist == Plie {
			err = errors.New("Plist == Plie")
		}
	}
	return vptr, vpath, err
}

func Compntptr(name string, Compnt []*COMPNT) *COMPNT {
	for i := range Compnt {
		if name == Compnt[i].Name {
			return Compnt[i]
		}
	}

	return nil
}

/* ----------------------------------------------------------------- */

// システム要素出口温度、湿度のポインター
func compntvptr(nk int, key []string, Compnt *COMPNT) (VPTR, error) {
	var i int
	var err error
	var vptr VPTR

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
			vptr = VPTR{
				Ptr:  &v.X,
				Type: VAL_CTYPE,
			}
			v.Org = 'y'
		}
	} else {
		for i = 0; i < Compnt.Nout; i++ {
			Eo := Compnt.Elouts[i]
			if (Eo.Fluid == AIRa_FLD && string(key[1]) == "Taout") ||
				(Eo.Fluid == AIRx_FLD && string(key[1]) == "xout") ||
				(Eo.Fluid == WATER_FLD && string(key[1]) == "Twout") {
				vptr = VPTR{
					Ptr:  &Eo.Sysv,
					Type: VAL_CTYPE,
				}
				break
			}
		}
		if i == Compnt.Nout {
			err = errors.New("i == Compnt.Nout")
		}
	}
	return vptr, err
}

/* ----------------------------------------------------------------- */

// 負荷計算を行うシステム要素の設定システム変数のポインターを作成します。
// 負荷計算を行うシステム要素の設定システム変数のポインターを作成し、 vtr に保存します。
// 内部では、 boildptr, refaldptr, hcldptr, pipeldsptr, rdpnlldsptr,roomldptr に処理を委譲します。
func loadptr(loadcmp *COMPNT, load *ControlSWType, s string, _Compnt []*COMPNT) (VPTR, error) {
	var Room *ROOM
	var key []string
	var idmrk byte = ' '
	var err error
	var vptr VPTR

	key = strings.Split(s, "_")
	nk := len(key)

	if nk != 0 {
		for i := range _Compnt {
			Compnt := _Compnt[i]
			if key[0] == Compnt.Name {
				switch Compnt.Eqptype {
				case BOILER_TYPE:
					vptr, err = boildptr(load, key, Compnt.Eqp.(*BOI))
					idmrk = 't'
				case REFACOMP_TYPE:
					vptr, err = refaldptr(load, key, Compnt.Eqp.(*REFA))
					idmrk = 't'
				case HCLOAD_TYPE, RMAC_TYPE, RMACD_TYPE:
					if SIMUL_BUILDG {
						vptr, err = hcldptr(load, key, Compnt.Eqp.(*HCLOAD), &idmrk)
					}
				case PIPEDUCT_TYPE:
					if SIMUL_BUILDG {
						vptr, err = pipeldsptr(load, key, Compnt.Eqp.(*PIPE), &idmrk)
					}
				case RDPANEL_TYPE:
					if SIMUL_BUILDG {
						Rdpnl := Compnt.Eqp.(*RDPNL)
						vptr, err = rdpnlldsptr(load, key, Rdpnl, &idmrk)
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
						vptr, err = roomldptr(load, key, Room, &idmrk)
					}
				}

				if err == nil {
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
				err = errors.New("")
			}
		}
		return vptr, err
	} else {
		return vptr, errors.New("s is empty")
	}
}

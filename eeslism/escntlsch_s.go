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

/*  es_cntlsch_s.c */
package eeslism

import (
	"fmt"
	"io"
)

func Contlschdlr(Ncontl int, _Contl []CONTL, Mpath []*MPATH, _Compnt []COMPNT) {

	// 全ての経路、機器を停止で初期化
	for _, Mp := range Mpath {
		Mp.Control = OFF_SW
		mpathschd(OFF_SW, Mp.Plist)
	}

	// 機器の制御情報を「停止」で初期化
	for i := range _Compnt {
		Compnt := &_Compnt[i]
		Compnt.Control = OFF_SW

		if Compnt.Eqptype == VALV_TYPE || Compnt.Eqptype == TVALV_TYPE {
			v := Compnt.Eqp.(*VALV)
			if v.Xinit != nil {
				v.X = *v.Xinit
			} else {
				v.X = 1.0
			}
		}

		for m := 0; m < Compnt.Nout; m++ {
			Eo := Compnt.Elouts[m]
			Eo.Sysld = 'n'
			Eo.Control = OFF_SW
		}
	}

	// CONTLの制御情報を反映
	for i := 0; i < Ncontl; i++ {
		Contl := &_Contl[i]
		Contl.Lgv = 1
		// True:1、False:0
		// if分で制御される場合
		if Contl.Type == 'c' {
			Contl.Lgv = contrlif(Contl.Cif)

			// ANDで2条件の場合
			if Contl.AndCif != nil {
				Contl.Lgv = Contl.Lgv * contrlif(Contl.AndCif)
			}

			// ANDで3条件の場合
			if Contl.AndAndCif != nil {
				Contl.Lgv = Contl.Lgv * contrlif(Contl.AndAndCif)
			}

			// or条件の場合
			if Contl.OrCif != nil {
				Contl.Lgv = Contl.Lgv + contrlif(Contl.OrCif)
			}
		}

		if Contl.Lgv != 0 {
			if Contl.Cst != nil {
				if Contl.Cst.Type == VAL_CTYPE {
					*Contl.Cst.Lft.V = *Contl.Cst.Rgt.V
				} else {
					*Contl.Cst.Lft.S = *Contl.Cst.Rgt.S

					if Contl.Cst.PathType == MAIN_CPTYPE {
						Mp := Contl.Cst.Path.(*MPATH)
						Mp.Control = ControlSWType((*Contl.Cst.Lft.S)[0])
						mpathschd(Mp.Control, Mp.Plist)
					} else if Contl.Cst.PathType == LOCAL_CPTYPE {
						Pli := Contl.Cst.Path.(*PLIST)
						lpathscdd(ControlSWType((*Contl.Cst.Lft.S)[0]), Pli)
					}
				}
			}
		}
	}

	for i := range _Compnt {
		Compnt := &_Compnt[i]

		switch Compnt.Eqptype {
		case ROOM_TYPE:
			if SIMUL_BUILDG {
				Compnt.Control = ON_SW
				Eo := Compnt.Elouts[0]
				Eo.Control = ON_SW
				Eo = Compnt.Elouts[1]
				Eo.Control = ON_SW
				roomldschd(Compnt.Eqp.(*ROOM))
			}
		case BOILER_TYPE:
			boildschd(Compnt.Eqp.(*BOI))
		case REFACOMP_TYPE:
			refaldschd(Compnt.Eqp.(*REFA))
		case HCLOAD_TYPE, HCLOADW_TYPE, RMAC_TYPE, RMACD_TYPE:
			hcldschd(Compnt.Eqp.(*HCLOAD))
		case PIPEDUCT_TYPE:
			pipeldsschd(Compnt.Eqp.(*PIPE))
		case RDPANEL_TYPE:
			rdpnlldsschd(Compnt.Eqp.(*RDPNL))
		case HCCOIL_TYPE:
			Eo := Compnt.Elouts[1]
			if Eo.Lpath == nil {
				Eo.Control = OFF_SW
				Hcc := Compnt.Eqp.(*HCC)
				Hcc.Wet = 'd'
			}
		case FLIN_TYPE:
			Compnt.Control = FLWIN_SW
			Flin := Compnt.Eqp.(*FLIN)
			Eo := Compnt.Elouts[0]
			Eo.Control = FLWIN_SW
			Eo.Sysv = *Flin.Vart
			if Flin.Awtype == 'A' {
				Eo := Compnt.Elouts[1]
				Eo.Control = FLWIN_SW
				Eo.Sysv = *Flin.Varx
			}
		}

		for m := 0; m < Compnt.Nout; m++ {
			Eo := Compnt.Elouts[m]
			if Eo.Control == LOAD_SW {
				Eo.Eldobj.Sysld = 'y'
			}
		}
	}

	for _, Mp := range Mpath {
		for _, Pli := range Mp.Plist {
			if Pli.Batch {
				lpathschbat(Pli)
			}
		}
	}
}

/* --------------------------------------------------- */

func contrlif(ctlif *CTLIF) int {
	var id int
	var a, b float64

	boolToInt := func(b bool) int {
		if b {
			return 1
		}
		return 0
	}

	if ctlif.Type == VAL_CTYPE {
		a = *ctlif.Lft1.V
		if ctlif.Nlft == 2 {
			a -= *ctlif.Lft2.V
		}
		b = *ctlif.Rgt.V

		switch ctlif.Op {
		case 'l':
			id = boolToInt(a < b)
		case 'g':
			id = boolToInt(a > b)
		case 'E':
			id = boolToInt(a == b)
		case 'L':
			id = boolToInt(a <= b)
		case 'G':
			id = boolToInt(a >= b)
		case 'N':
			id = boolToInt(a != b)
		default:
			id = 0
		}
	} else {
		switch ctlif.Op {
		case 'E':
			id = boolToInt(*ctlif.Lft1.S == *ctlif.Rgt.S)
		case 'N':
			id = boolToInt(*ctlif.Lft1.S != *ctlif.Rgt.S)
		default:
			id = 0
		}
	}

	return id
}

/* --------------------------------------------------- */

func mpathschd(control ControlSWType, Plist []*PLIST) {
	for j := range Plist {
		Plist[j].Control = control
		lpathscdd(control, Plist[j])
	}
}

/* --------------------------------------------------- */

func lpathscdd(control ControlSWType, plist *PLIST) {
	if plist.Org {
		lpathschd(control, plist.Pelm)

		if plist.Lpair != nil {
			plist.Lpair.Control = control
			lpathschd(control, plist.Lpair.Pelm)
		}
	}
}

/* --------------------------------------------------- */

func lpathschd(control ControlSWType, pelm []*PELM) {
	for _, Pelm := range pelm {

		// 電気蓄熱暖房器の機器制御は変更しない
		if Pelm.Cmp.Eqptype != STHEAT_TYPE {
			Pelm.Cmp.Control = control
		}

		if Pelm.Out != nil {
			// 一時的に室の場合は止められないようにした。電気蓄熱暖房器も同様
			// 順次対応予定
			if Pelm.Cmp.Eqptype != ROOM_TYPE { // && strcmp(Pelm.cmp.eqptype, STHEAT_TYPE) != 0)
				Pelm.Out.Control = control
			}
			//else if (strcmp(Pelm.cmp.eqptype, STHEAT_TYPE) != 0)
			//	Pelm.out.control = control;
		}
	}
}

/* --------------------------------------------------- */

/* 蓄熱槽のバッチ給水、排水時の設定  */

func lpathschbat(Plist *PLIST) {
	var j, k, i, jt, ifl int
	var batop rune
	var Tsout, Gbat float64
	var Stank *STANK
	var Pelm *PELM

	Stank = nil
	for j, Pe := range Plist.Pelm {

		if Pe.Cmp.Eqptype == STANK_TYPE {
			Plist.G = Gbat
			Gbat = 0.0
			for k = 0; k < len(Pe.Cmp.Elouts); k++ {
				if Pe.Out == Pe.Cmp.Elouts[k] {
					break
				}
			}
			Stank = Pe.Cmp.Eqp.(*STANK)
			Stank.Batchop = rune(OFF_SW)
			batop = Stank.Batchcon[k]

			if batop == BTFILL {
				for i = 0; i < Stank.Ndiv; i++ {
					if Stank.DtankF[i] == TANK_EMPTY {
						jt = j
						Gbat += Stank.Dvol[i]
						Stank.Batchop = BTFILL
					}
				}
			} else {
				ifl = 0
				for i = 0; i < Stank.Ndiv; i++ {
					if Stank.DtankF[i] == TANK_FULL {
						ifl++
					}
				}
				if Stank.Ndiv == ifl {
					Stank.Batchop = BTDRAW
					jt = j
				}
			}
			break
		}
	}

	peiIdx := 0
	if Stank.Batchop == BTFILL || Stank.Batchop == BTDRAW {
		if Pelm.Out.Control == FLWIN_SW {
			peiIdx++
		}

		lpathschd(ON_SW, Plist.Pelm[peiIdx:])

		if Stank.Batchop == BTFILL {
			Plist.G = Gbat * Row / DTM
			lpathschd(OFF_SW, Plist.Pelm[jt:])
			Plist.Pelm[jt].Out.Control = BATCH_SW
		} else if Stank.Batchop == BTDRAW {
			lpathschd(OFF_SW, Plist.Pelm[:jt])
			Plist.Pelm[jt].Out.Control = BATCH_SW

			for j = 0; j <= Stank.Jout[k]; j++ {
				Stank.DtankF[j] = TANK_EMPTY
				Tsout += Stank.Tssold[j]
				Gbat += Stank.Dvol[j]
			}
			Tsout /= float64(Stank.Jout[k] + 1)
			Plist.G = Gbat * Row / DTM

			Plist.Pelm[jt].Out.Sysv = Tsout
		}
	}
}

/* --------------------------------------------------- */

func contlxprint(Ncontl int, C *CONTL, out io.Writer) {
	var i int
	var cif *CTLIF
	var cst *CTLST
	var V float64
	var Contl *CONTL

	Contl = C
	if DEBUG {
		fmt.Fprintln(out, "contlxprint --- Contlschdlr")

		for i = 0; i < Ncontl; i++ {
			fmt.Fprintf(out, "[%d]  type=%c  lgv=%d\n",
				i, Contl.Type, Contl.Lgv)
			cif = Contl.Cif
			cst = Contl.Cst

			if cif != nil {
				fmt.Fprintf(out, "  iftype=%c [ %c ]", cif.Type, cif.Op)
				if cif.Type == VAL_CTYPE {
					if cif.Nlft > 1 {
						V = *cif.Lft2.V
					} else {
						V = 0.0
					}

					fmt.Fprintf(out, "  lft1=%.2f  lft2=%.2f  rgt=%.2f\n", *cif.Lft1.V, V, *cif.Rgt.V)
				} else {
					fmt.Fprintf(out, "  lft1=%s  lft2=%s  rgt=%s\n", *cif.Lft1.S, *cif.Lft2.S, *cif.Rgt.S)
				}
			}

			fmt.Fprintf(out, "  sttype=%c  pathtype=%c", cst.Type, cst.PathType)
			if cst.Type == VAL_CTYPE {
				fmt.Fprintf(out, "  lft=%.2f    rgt=%.2f\n", *cst.Lft.V, *cst.Rgt.V)
			} else {
				fmt.Fprintf(out, "  lft=%s    rgt=%s\n", *cst.Lft.S, *cst.Rgt.S)
			}
		}
	}

	Contl = C
	if Ferr != nil {
		fmt.Fprintln(out, "contlxprint --- Contlschdlr")

		for i = 0; i < Ncontl; i++ {
			fmt.Fprintf(out, "[%d]\ttype=%c\tlgv=%d\n", i, Contl.Type, Contl.Lgv)
			cif = Contl.Cif
			cst = Contl.Cst

			if cif != nil {
				fmt.Fprintf(out, "\tiftype=%c\t[%c]", cif.Type, cif.Op)
				if cif.Type == VAL_CTYPE {
					if cif.Nlft > 1 {
						V = *cif.Lft2.V
					} else {
						V = 0.0
					}

					fmt.Fprintf(out, "\tlft1=%.2f\tlft2=%.2f\trgt=%.2f\n", *cif.Lft1.V, V, *cif.Rgt.V)
				} else {
					fmt.Fprintf(out, "\tlft1=%s\tlft2=%s\trgt=%s\n", *cif.Lft1.S, *cif.Lft2.S, *cif.Rgt.S)
				}
			}

			fmt.Fprintf(out, "\tsttype=%c\tpathtype=%c", cst.Type, cst.PathType)
			if cst.Type == VAL_CTYPE {
				fmt.Fprintf(out, "\tlft=%.2f\trgt=%.2f\n", *cst.Lft.V, *cst.Rgt.V)
			} else {
				fmt.Fprintf(out, "\tlft=%s\trgt=%s\n", *cst.Lft.S, *cst.Rgt.S)
			}
		}
	}
}

/* --------------------------------------------------- */

/* 負荷計算モードの再設定、
暖房時冷房負荷または冷房時暖房負荷のときは自然室温計算　*/

func rmloadreset(Qload float64, loadsw rune, Eo *ELOUT, SWITCH ControlSWType) int {
	if Eo.Sysld == 'y' {
		if (loadsw == HEATING_LOAD && Qload < 0.0) ||
			(loadsw == COOLING_LOAD && Qload > 0.0) {
			Eo.Control = SWITCH
			Eo.Sysld = 'n'
			return 1
		} else {
			return 0
		}
	} else {
		return 0
	}
}

/* --------------------------------------------------- */

/* 加熱、冷却モードの再設定、
加熱房時冷房負荷または冷却時暖房負荷のときは機器停止　*/

func chswreset(Qload float64, chmode rune, Eo *ELOUT) int {
	var Elo *ELOUT

	if (chmode == HEATING_SW && Qload < 0.0) ||
		(chmode == COOLING_SW && Qload > 0.0) {
		Elo = Eo
		Elo.Control = ON_SW
		Elo.Sysld = 'n'
		Elo.Emonitr.Control = ON_SW
		return 1
	} else {
		return 0
	}
}

/* --------------------------------------------------- */

/* 仮想空調機の再設定
湿りコイルの非除湿時乾きコイルへの変更 */

func chqlreset(Hcload *HCLOAD) int {

	Ql := Hcload.Ql
	Qs := Hcload.Qs
	wet := Hcload.Wetmode
	chmode := Hcload.Chmode
	Elo := Hcload.Cmp.Elouts[1]
	//Elos := Hcload.Cmp.Elouts[0]

	if (wet && Ql > 1.e-6) || (wet && Qs >= 0.0) {
		Hcload.Wetmode = false
		Elo.Control = ON_SW
		Elo.Sysld = 'n'

		if Elo.Emonitr != nil {
			Elo.Emonitr.Control = ON_SW
		}

		return 1
	}

	if chmode == COOLING_SW && (Ql > 0.0 || Qs >= 0.0) {
		Elo.Control = ON_SW
		Hcload.Wetmode = false
		Elo.Sysld = 'n'
		if Elo.Emonitr != nil {
			Elo.Emonitr.Control = ON_SW
		}

		return 1
	}

	return 0
}

/* 過負荷運転ための再設定 */

func maxcapreset(Qload, Qmax float64, chmode rune, Eo *ELOUT) int {
	var Boi *BOI

	var Eosysld rune
	var Boimode rune
	var Eocontrol, Eoemonitrcontrol, Boicmpcntrol ControlSWType

	Boi = Eo.Cmp.Eqp.(*BOI)

	Eocontrol = Eo.Control
	Eosysld = Eo.Sysld
	if Eo.Emonitr != nil {
		Eoemonitrcontrol = Eo.Emonitr.Control
	} else {
		return 0
	}

	Boicmpcntrol = Boi.Cmp.Control
	Boimode = Boi.Mode

	if (chmode == HEATING_SW && Qload > Qmax) ||
		(chmode == COOLING_SW && Qload < Qmax) {
		// 過負荷なので最大能力で再計算する
		Eo.Control = ON_SW
		Eo.Sysld = 'n'
		Eo.Emonitr.Control = ON_SW
		Boi.Cmp.Control = ON_SW

		//最大能力運転フラグ
		Boi.Mode = 'M'
	}

	if Eo.Control == Eocontrol &&
		Eo.Sysld == Eosysld &&
		Eo.Emonitr.Control == Eoemonitrcontrol &&
		Boi.Cmp.Control == Boicmpcntrol &&
		Boi.Mode == Boimode {
		return 0
	} else {
		return 1
	}
}

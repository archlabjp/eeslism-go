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

/*   tcomfrt.c    */

package eeslism

import "math"

/*   作用温度制御時の設定室内空気温度  */

var __Rmotset_Pint int

func Rmotset(Nroom int, _Room []ROOM) {
	Fotinit(Nroom, _Room)

	for i := 0; i < Nroom; i++ {
		Room := &_Room[i]

		Eo := Room.cmp.Elouts[0]
		if Eo.Control == LOAD_SW {
			rmld := Room.rmld

			if rmld.tropt == 'o' {
				Fotf(Room)

				OT := rmld.Tset
				a := OT - rmld.FOC

				for j := 0; j < Room.Ntr; j++ {
					trnx := Room.trnx[j]
					a -= rmld.FOTN[j] * trnx.nextroom.Tr
				}

				for j := 0; j < Room.Nrp; j++ {
					rmpnl := &Room.rmpnl[j]

					var Twi float64
					if __Rmotset_Pint == 0 {
						Twi = rmpnl.sd.mw.Tw[rmpnl.sd.mw.mp]
						__Rmotset_Pint = 1
					} else {
						Twi = rmpnl.pnl.Tpi
					}

					a -= rmld.FOPL[j] * Twi
				}

				Trset := a / rmld.FOTr
				Eo.Sysv = Trset
				Room.Tr = Trset

				if rmld.loadx != nil {
					Eo = Room.cmp.Elouts[1]
					if Eo.Control == LOAD_SW && rmld.hmopt == 'r' {
						Eo.Sysv = FNXtr(Trset, rmld.Xset)
					}
				}
			}
		}
	}
}

/* -------------------------------------- */

var __Fotinit_init int = 'i'

func Fotinit(Nroom int, _Room []ROOM) {
	if __Fotinit_init == 'i' {
		for i := 0; i < Nroom; i++ {
			Room := &_Room[i]
			if Room.rmld != nil {
				Room.rmld.FOTN = nil
				Room.rmld.FOPL = nil

				Room.rmld.FOTN = make([]float64, Room.Ntr)
				Room.rmld.FOPL = make([]float64, Room.Nrp)
			}
		}
		__Fotinit_init = 'x'
	}
}

/* -------------------------------------- */
func Fotf(Room *ROOM) {
	var r float64
	if Room.OTsetCwgt == nil || *(Room.OTsetCwgt) < 0.0 || *(Room.OTsetCwgt) > 1.0 {
		r = 0.5
	} else {
		r = *(Room.OTsetCwgt)
	}

	{
		var a, C float64
		for i := 0; i < Room.N; i++ {
			Sd := &Room.rsrf[i]
			a += Sd.A * Sd.WSR
			C += Sd.A * Sd.WSC
		}
		Room.rmld.FOTr = r + (1.0-r)*a/Room.Area
		Room.rmld.FOC = (1.0 - r) * C / Room.Area
	}

	for k := 0; k < Room.Ntr; k++ {
		ft := &Room.rmld.FOTN[k]

		var a float64
		for i := 0; i < Room.N; i++ {
			Sd := &Room.rsrf[i]
			a += Sd.A * Sd.WSRN[k]
		}

		*ft = (1.0 - r) * a / Room.Area
	}

	for k := 0; k < Room.Nrp; k++ {
		ft := &Room.rmld.FOPL[k]

		var a float64
		for i := 0; i < Room.N; i++ {
			Sd := &Room.rsrf[i]
			a += Sd.A * Sd.WSPL[k]
		}

		*ft = (1.0 - r) * a / Room.Area
	}
}

/* -------------------------------------- */

/*   各室の温熱環境指標計算　　*/

func Rmcomfrt(Nroom int, _Room []ROOM) {
	met := 0.0
	Icl := 0.0
	v := 0.0

	for i := 0; i < Nroom; i++ {
		id := 0
		Room := &_Room[i]
		if Room.Metsch != nil && *Room.Metsch > 0.0 {
			met = *Room.Metsch
			id++
		}
		if Room.Closch != nil && *Room.Closch > 0.0 {
			Icl = *Room.Closch
			id++
		}
		if Room.Wvsch != nil && *Room.Wvsch > 0.0 {
			v = *Room.Wvsch
			id++
		}

		if id == 3 {
			Room.PMV = Pmv0(met, Icl, Room.Tr, Room.xr, Room.Tsav, v)
			Room.SET = SET_star(Room.Tr, Room.Tsav, v, Room.RH, met, Icl, 0.0, 101.3)
		} else {
			Room.PMV = -999.0
		}
		/*******************
		fmt.Printf("**** Rmcomfrt  met=%.1f Icl=%.1f v=%.2f  Tr=%.1f  xr=%.4f Tmrt=%.1f  PMV=%.2f\n",
			met, Icl, v, Room.Tr, Room.xr, Room.Tsav, Room.PMV)
		*******************/
	}
}

/* ----------------------------------------------------- */

/*   PMVの計算     */

func Pmv0(met, Icl, Tr, xr, Tmrt, v float64) float64 {
	/* m [kcal.m2h] */

	Po := 760.0
	eta := 0.0

	m := met * 50.0
	Pa := xr * Po / (xr + 0.62198)
	hc := 10.4 * math.Sqrt(v)
	Tm := 0.5*(37.0+0.5*(Tr+Tmrt)) + 273.15
	hr := 13.6e-8 * Tm * Tm * Tm
	fcl := 1.0
	if Icl < 0.5 {
		fcl = 1.0 + 0.2*Icl
	} else {
		fcl = 1.05 + 0.1*Icl
	}
	Ifcl := 0.18 * Icl * fcl

	tcl := (35.7 - 0.032*m*(1.0-eta) + Ifcl*(hr*Tmrt+hc*Tr)) / (1.0 + Ifcl*(hr+hc))

	L := m*(0.60135-0.0023*(44.0-Pa)-0.0014*(34.0-Tr)) + 0.35*Pa + 5.95 - 0.6013*eta - fcl*(hr*(tcl-Tmrt)+hc*(tcl-Tr))

	return (0.352*math.Exp(-0.042*m) + 0.032) * L
}

//	SET*の計算
func SET_star(TA, TR, VEL, RH, MET, CLO, WME, PATM float64) float64 {
	//Input doubleiables ? TA (air temperature): °C, TR (mean radiant temperature): °C, VEL (air velocity): m/s,
	//RH (relative humidity): %, MET: met unit, CLO: clo unit, WME (external work): W/m2, PATM (atmospheric pressure): kPa
	const KCLO = 0.25
	const BODYWEIGHT = 69.9        //kg
	const BODYSURFACEAREA = 1.8258 //m2
	const METFACTOR = 58.2         //W/m2
	const SBC = 0.000000056697     //Stefan-Boltzmann constant (W/m2K4)
	const CSW = 170.0
	const CDIL = 120.0
	const CSTR = 0.5
	const LTIME = 60
	var PS = FindSaturatedVaporPressureTorr(TA)
	var VaporPressure = RH * PS / 100.0
	var AirVelocity = math.Max(VEL, 0.1)
	const TempSkinNeutral = 33.7
	const TempCoreNeutral = 36.49
	const TempBodyNeutral = 36.49
	const SkinBloodFlowNeutral = 6.3
	var TempSkin = TempSkinNeutral //Initial values
	var TempCore = TempCoreNeutral
	var SkinBloodFlow = SkinBloodFlowNeutral
	var MSHIV = 0.0
	var ALFA = 0.1
	var ESK = 0.1 * MET
	// 桁落ちによる誤差を避けるため換算係数を変更
	//double PressureInAtmospheres = PATM * 0.009869;
	var PressureInAtmospheres = PATM / 101.3
	var RCL = 0.155 * CLO
	var FACL = 1.0 + 0.15*CLO
	var LR = 2.2 / PressureInAtmospheres //Lewis Relation is 2.2 at sea level
	var RM = MET * METFACTOR
	var M = MET * METFACTOR
	var ICL, WCRIT float64
	if CLO <= 0 {
		WCRIT = 0.38 * math.Pow(AirVelocity, -0.29)
		ICL = 1.0
	} else {
		WCRIT = 0.59 * math.Pow(AirVelocity, -0.08)
		ICL = 0.45
	}
	var CHC = 3.0 * math.Pow(PressureInAtmospheres, 0.53)
	var CHCV = 8.600001 * math.Pow((AirVelocity*PressureInAtmospheres), 0.53)
	CHC = math.Max(CHC, CHCV)
	var CHR = 4.7
	var CTC = CHR + CHC
	var RA = 1.0 / (FACL * CTC) //Resistance of air layer to dry heat transfer
	var TOP = (CHR*TR + CHC*TA) / CTC
	var TCL = TOP + (TempSkin-TOP)/(CTC*(RA+RCL))
	//TCL and CHR are solved iteratively using: H(Tsk - TOP) = CTC(TCL - TOP),
	//where H = 1/(RA + RCL) and RA = 1/FACL*CTC
	var TCL_OLD = TCL
	var flag = true
	var DRY, HFCS, ERES, CRES, SCR, SSK, TCSK, TCCR, DTSK, DTCR, TB, SKSIG, WARMS, COLDS, CRSIG, WARMC,
		COLDC, BDSIG, WARMB, REGSW, EMAX float64
	var PWET = 0.0

	//Begin iteration
	for TIM := 1; TIM <= LTIME; TIM++ {
		for i := 0; i < 100; i++ {
			if flag {
				TCL_OLD = TCL
				CHR = 4.0 * SBC * math.Pow(((TCL+TR)/2.0+273.15), 3.0) * 0.72
				CTC = CHR + CHC
				RA = 1.0 / (FACL * CTC) //Resistance of air layer to dry heat transfer
				TOP = (CHR*TR + CHC*TA) / CTC
			}
			TCL = (RA*TempSkin + RCL*TOP) / (RA + RCL)
			flag = true
			if math.Abs(TCL-TCL_OLD) <= 0.01 {
				break
			}
		}
		flag = false
		DRY = (TempSkin - TOP) / (RA + RCL)
		HFCS = (TempCore - TempSkin) * (5.28 + 1.163*SkinBloodFlow)
		ERES = 0.0023 * M * (44.0 - VaporPressure)
		CRES = 0.0014 * M * (34.0 - TA)
		SCR = M - HFCS - ERES - CRES - WME
		SSK = HFCS - DRY - ESK
		TCSK = 0.97 * ALFA * BODYWEIGHT
		TCCR = 0.97 * (1. - ALFA) * BODYWEIGHT
		DTSK = (SSK * BODYSURFACEAREA) / (TCSK * 60.0) //°C/min
		DTCR = SCR * BODYSURFACEAREA / (TCCR * 60.0)   //°C/min
		TempSkin = TempSkin + DTSK
		TempCore = TempCore + DTCR
		TB = ALFA*TempSkin + (1.-ALFA)*TempCore
		SKSIG = TempSkin - TempSkinNeutral
		if SKSIG > 0 {
			WARMS = SKSIG
		} else {
			WARMS = 0
		}
		if (-1.0 * SKSIG) > 0 {
			COLDS = (-1.0 * SKSIG)
		} else {
			COLDS = 0
		}
		CRSIG = (TempCore - TempCoreNeutral)
		if CRSIG > 0 {
			WARMC = CRSIG
		} else {
			WARMC = 0
		}
		if (-1.0 * CRSIG) > 0 {
			COLDC = (-1.0 * CRSIG)
		} else {
			COLDC = 0
		}
		BDSIG = TB - TempBodyNeutral
		if BDSIG > 0 {
			WARMB = BDSIG
		} else {
			WARMB = 0
		}
		SkinBloodFlow = (SkinBloodFlowNeutral + CDIL*WARMC) / (1. + CSTR*COLDS)
		SkinBloodFlow = math.Max(0.5, math.Min(90.0, SkinBloodFlow))
		REGSW = CSW * WARMB * math.Exp(WARMS/10.7)
		REGSW = math.Min(REGSW, 500.0)
		var ERSW = 0.68 * REGSW
		var REA = 1.0 / (LR * FACL * CHC) //Evaporative resistance of air layer
		var RECL = RCL / (LR * ICL)       //Evaporative resistance of clothing (icl=.45)
		EMAX = (FindSaturatedVaporPressureTorr(TempSkin) - VaporPressure) / (REA + RECL)
		var PRSW = ERSW / EMAX
		PWET = 0.06 + 0.94*PRSW
		var EDIF = PWET*EMAX - ERSW
		ESK = ERSW + EDIF
		//if (TIM == 60)
		//	printf("aqa\n");
		if PWET > WCRIT {
			PWET = WCRIT
			PRSW = WCRIT / 0.94
			ERSW = PRSW * EMAX
			EDIF = 0.06 * (1.0 - PRSW) * EMAX
			ESK = ERSW + EDIF
		}
		if EMAX < 0 {
			EDIF = 0
			ERSW = 0
			PWET = WCRIT
			PRSW = WCRIT
			ESK = EMAX
		}
		ESK = ERSW + EDIF
		//printf("%d\t%f\n", TIM, ESK);
		MSHIV = 19.4 * COLDS * COLDC
		M = RM + MSHIV
		ALFA = 0.0417737 + 0.7451833/(SkinBloodFlow+0.585417)
	} //End iteration

	var HSK = DRY + ESK //Total heat loss from skin
	var RN = M - WME    //Net metabolic heat production
	var ECOMF = 0.42 * (RN - (1. * METFACTOR))
	if ECOMF < 0.0 {
		ECOMF = 0.0 //From Fanger
	}
	EMAX = EMAX * WCRIT
	var W = PWET
	var PSSK = FindSaturatedVaporPressureTorr(TempSkin)
	var CHRS = CHR //Definition of ASHRAE standard environment
	//... denoted “S”
	var CHCS float64
	if MET < 0.85 {
		CHCS = 3.0
	} else {
		CHCS = 5.66 * math.Pow((MET-0.85), 0.39)
		CHCS = math.Max(CHCS, 3.0)
	}
	var CTCS = CHCS + CHRS
	var RCLOS = 1.52/((MET-WME/METFACTOR)+0.6944) - 0.1835
	var RCLS = 0.155 * RCLOS
	var FACLS = 1.0 + KCLO*RCLOS
	var FCLS = 1.0 / (1.0 + 0.155*FACLS*CTCS*RCLOS)
	const IMS = 0.45
	var ICLS = IMS * CHCS / CTCS * (1. - FCLS) / (CHCS/CTCS - FCLS*IMS)
	var RAS = 1.0 / (FACLS * CTCS)
	var REAS = 1.0 / (LR * FACLS * CHCS)
	var RECLS = RCLS / (LR * ICLS)
	var HD_S = 1.0 / (RAS + RCLS)
	var HE_S = 1.0 / (REAS + RECLS)
	//SET determined using Newton’s iterative solution
	const DELTA = .0001
	var dx = 100.0
	var SET, ERR1, ERR2 float64
	var SET_OLD = TempSkin - HSK/HD_S //Lower bound for SET
	for math.Abs(dx) > .01 {
		ERR1 = (HSK - HD_S*(TempSkin-SET_OLD) - W*HE_S*(PSSK-0.5*
			FindSaturatedVaporPressureTorr(SET_OLD)))
		ERR2 = (HSK - HD_S*(TempSkin-(SET_OLD+DELTA)) - W*HE_S*(PSSK-0.5*
			FindSaturatedVaporPressureTorr((SET_OLD+DELTA))))
		SET = SET_OLD - DELTA*ERR1/(ERR2-ERR1)
		dx = SET - SET_OLD
		SET_OLD = SET
	}
	return SET
}

func FindSaturatedVaporPressureTorr(T float64) float64 {
	//Helper function for pierceSET calculates Saturated Vapor Pressure (Torr) at Temperature T (°C)
	return math.Exp(18.6686 - 4030.183/(T+235.0))
}

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

/*  eqpcat.c  */

package eeslism

// 名称がcatname の設備をEcatから探して、 C に格納する
// Esys には、設備の種類ごとの個数を格納する
func eqpcat(catname string, C *COMPNT, Ecat *EQCAT, Esys *EQSYS) bool {
	C.Airpathcpy = false
	C.Idi = nil
	C.Ido = nil

	for i := range Ecat.Hccca {
		Hccca := &Ecat.Hccca[i]
		if catname == Hccca.name {
			C.Eqptype = HCCOIL_TYPE
			C.Ncat = i
			C.Neqp = len(Esys.Hcc)
			Esys.Hcc = append(Esys.Hcc, NewHCC())
			C.Nout = 3
			C.Nin = 3
			C.Idi = []ELIOType{ELIO_t, ELIO_x, ELIO_W} // txW
			C.Ido = []ELIOType{ELIO_t, ELIO_x, ELIO_W} // txW
			C.Airpathcpy = true
			return true
		}
	}

	for i := range Ecat.Boica {
		Boica := &Ecat.Boica[i]
		if catname == Boica.name {
			C.Eqptype = BOILER_TYPE
			C.Ncat = i
			C.Neqp = len(Esys.Boi)
			Esys.Boi = append(Esys.Boi, NewBOI())
			C.Nout = 1
			C.Nin = 1
			return true
		}
	}

	for i := range Ecat.Collca {
		Collca := &Ecat.Collca[i]
		if catname == Collca.name {
			C.Eqptype = COLLECTOR_TYPE
			C.Ncat = i
			C.Neqp = len(Esys.Coll)
			Esys.Coll = append(Esys.Coll, NewCOLL())
			C.Ac = Collca.Ac

			if Collca.Type == COLLECTOR_PDT {
				C.Nout = 1
				C.Nin = 1
			} else {
				C.Nout = 2
				C.Nin = 2
				C.Airpathcpy = true
			}
			return true
		}
	}

	for i := range Ecat.PVca {
		PVca := &Ecat.PVca[i]
		if catname == PVca.Name {
			C.Eqptype = PV_TYPE
			C.Ncat = i
			C.Neqp = len(Esys.PVcmp)
			Esys.PVcmp = append(Esys.PVcmp, NewPV())
			C.PVcap = PVca.PVcap
			C.Nout = 0
			C.Nin = 0
			C.Area = PVca.Area

			return true
		}
	}

	for i := range Ecat.Refaca {
		Refaca := &Ecat.Refaca[i]
		if catname == Refaca.name {
			C.Eqptype = REFACOMP_TYPE
			C.Ncat = i
			C.Neqp = len(Esys.Refa)
			Esys.Refa = append(Esys.Refa, NewREFA())
			C.Nout = 1
			C.Nin = 1
			return true
		}
	}

	for i := range Ecat.Pipeca {
		Pipeca := &Ecat.Pipeca[i]
		if catname == Pipeca.name {
			C.Eqptype = PIPEDUCT_TYPE
			C.Ncat = i
			C.Neqp = len(Esys.Pipe)
			Esys.Pipe = append(Esys.Pipe, NewPIPE())

			if Pipeca.Type == PIPE_PDT {
				C.Nout = 1
				C.Nin = 1
			} else {
				C.Nout = 2
				C.Nin = 2
				C.Airpathcpy = true
			}

			return true
		}
	}

	for i := range Ecat.Stankca {
		Stankca := &Ecat.Stankca[i]
		if catname == Stankca.name {
			C.Eqptype = STANK_TYPE
			C.Ncat = i
			C.Neqp = len(Esys.Stank)
			Esys.Stank = append(Esys.Stank, NewSTANK())

			return true
		}
	}

	for i := range Ecat.Hexca {
		Hexca := &Ecat.Hexca[i]
		if catname == Hexca.Name {
			C.Eqptype = HEXCHANGR_TYPE
			C.Ncat = i
			C.Neqp = len(Esys.Hex)
			Esys.Hex = append(Esys.Hex, NewHEX())

			C.Nout = 2
			C.Nin = 2
			C.Idi = []ELIOType{ELIO_C, ELIO_H} // CH
			C.Ido = []ELIOType{ELIO_C, ELIO_H} // CH

			return true
		}
	}

	for i := range Ecat.Pumpca {
		Pumpca := &Ecat.Pumpca[i]
		if catname == Pumpca.name {
			C.Eqptype = PUMP_TYPE
			C.Ncat = i
			C.Neqp = len(Esys.Pump)
			Esys.Pump = append(Esys.Pump, NewPUMP())

			if Pumpca.pftype == PUMP_PF {
				C.Nout = 1
				C.Nin = 1
			} else {
				C.Nout = 2
				C.Nin = 2
				C.Airpathcpy = true
			}

			return true
		}
	}

	/*---- Satoh Debug VAV  2000/12/5 ----*/
	for i := range Ecat.Vavca {
		Vavca := &Ecat.Vavca[i]
		if catname == Vavca.Name {
			if Vavca.Type == VAV_PDT {
				C.Eqptype = VAV_TYPE
			} else {
				C.Eqptype = VWV_TYPE
			}

			C.Ncat = i
			C.Neqp = len(Esys.Vav)
			Esys.Vav = append(Esys.Vav, NewVAV())

			if Vavca.Type == VAV_PDT {
				C.Nout = 2
				C.Nin = 2
				// 温湿度計算のために出入り口数は2
				C.Airpathcpy = true
			} else {
				C.Nout = 1
				C.Nin = 1
			}

			return true
		}
	}

	// Satoh OMVAV  2010/12/16
	for i := range Ecat.OMvavca {
		OMvavca := &Ecat.OMvavca[i]
		if catname == OMvavca.Name {
			C.Eqptype = OMVAV_TYPE
			C.Ncat = i
			C.Neqp = len(Esys.OMvav)
			Esys.OMvav = append(Esys.OMvav, NewOMVAV())
			C.Nout = 0
			C.Nin = 0

			return true
		}
	}

	for i := range Ecat.Stheatca {
		Stheatca := &Ecat.Stheatca[i]
		if catname == Stheatca.Name {
			C.Eqptype = STHEAT_TYPE
			C.Ncat = i
			C.Neqp = len(Esys.Stheat)
			Esys.Stheat = append(Esys.Stheat, NewSTHEAT())

			// NOTE: たぶんここは 2が正しいのでは
			C.Nout = 3
			C.Nin = 3
			// 温湿度計算のために出入り口数は2
			C.Airpathcpy = true

			return true
		}
	}

	// Satoh追加　デシカント槽　2013/10/23
	for i := range Ecat.Desica {
		Desica := &Ecat.Desica[i]
		if catname == Desica.name {
			C.Eqptype = DESI_TYPE
			C.Ncat = i
			C.Neqp = len(Esys.Desi)
			Esys.Desi = append(Esys.Desi, NewDESI())

			C.Nout = 2
			C.Nin = 2
			// 温湿度計算のために出入り口数は2
			C.Airpathcpy = true

			return true
		}
	}

	// Satoh追加　気化冷却器　2013/10/27
	for i := range Ecat.Evacca {
		Evcaca := &Ecat.Evacca[i]
		if catname == Evcaca.Name {
			C.Eqptype = EVAC_TYPE
			C.Ncat = i
			C.Neqp = len(Esys.Evac)
			Esys.Evac = append(Esys.Evac, NewEVAC())
			C.Airpathcpy = true
			C.Nout = 4
			C.Nin = 4
			C.Idi = []ELIOType{ELIO_D, ELIO_d, ELIO_W, ELIO_w} // DdWw
			C.Ido = []ELIOType{ELIO_D, ELIO_d, ELIO_W, ELIO_w} // DdWw

			return true
		}
	}

	for i := range Ecat.Thexca {
		Thexca := &Ecat.Thexca[i]
		if catname == Thexca.Name {
			C.Eqptype = THEX_TYPE
			C.Ncat = i
			C.Neqp = len(Esys.Thex)
			Esys.Thex = append(Esys.Thex, NewTHEX())
			C.Airpathcpy = true
			C.Nout = 4
			C.Nin = 4
			C.Idi = []ELIOType{ELIO_E, ELIO_e, ELIO_O, ELIO_o} // EeOo
			C.Ido = []ELIOType{ELIO_E, ELIO_e, ELIO_O, ELIO_o} // EeOo

			return true
		}
	}

	return false
}

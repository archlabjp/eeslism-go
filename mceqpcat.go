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

package main

// 名称がcatname の設備をEcatから探して、 C に格納する
// Esys には、設備の種類ごとの個数を格納する
func eqpcat(catname string, C *COMPNT, Ecat *EQCAT, Esys *EQSYS) bool {
	C.Airpathcpy = 'n'
	C.Idi = nil
	C.Ido = nil

	for i := 0; i < Ecat.Nhccca; i++ {
		Hccca := &Ecat.Hccca[i]
		if catname == Hccca.name {
			C.Eqptype = HCCOIL_TYPE
			C.Ncat = i
			C.Neqp = Esys.Nhcc
			Esys.Nhcc++
			C.Nout = 3
			C.Nin = 3
			C.Idi = []ELIOType{ELIO_t, ELIO_x, ELIO_W} // txW
			C.Ido = []ELIOType{ELIO_t, ELIO_x, ELIO_W} // txW
			C.Airpathcpy = 'y'
			return true
		}
	}

	for i := 0; i < Ecat.Nboica; i++ {
		Boica := &Ecat.Boica[i]
		if catname == Boica.name {
			C.Eqptype = BOILER_TYPE
			C.Ncat = i
			C.Neqp = Esys.Nboi
			Esys.Nboi++
			C.Nout = 1
			C.Nin = 1
			return true
		}
	}

	for i := 0; i < Ecat.Ncollca; i++ {
		Collca := &Ecat.Collca[i]
		if catname == Collca.name {
			C.Eqptype = COLLECTOR_TYPE
			C.Ncat = i
			C.Neqp = Esys.Ncoll
			Esys.Ncoll++
			C.Ac = Collca.Ac

			if Collca.Type == COLLECTOR_PDT {
				C.Nout = 1
				C.Nin = 1
			} else {
				C.Nout = 2
				C.Nin = 2
				C.Airpathcpy = 'y'
			}
			return true
		}
	}

	for i := 0; i < Ecat.Npvca; i++ {
		PVca := &Ecat.PVca[i]
		if catname == PVca.name {
			C.Eqptype = PV_TYPE
			C.Ncat = i
			C.Neqp = Esys.Npv
			Esys.Npv++
			C.PVcap = PVca.PVcap
			C.Nout = 0
			C.Nin = 0
			C.Area = PVca.Area

			return true
		}
	}

	for i := 0; i < Ecat.Nrefaca; i++ {
		Refaca := &Ecat.Refaca[i]
		if catname == Refaca.name {
			C.Eqptype = REFACOMP_TYPE
			C.Ncat = i
			C.Neqp = Esys.Nrefa
			Esys.Nrefa++
			C.Nout = 1
			C.Nin = 1
			return true
		}
	}

	for i := 0; i < Ecat.Npipeca; i++ {
		Pipeca := &Ecat.Pipeca[i]
		if catname == Pipeca.name {
			C.Eqptype = PIPEDUCT_TYPE
			C.Ncat = i
			C.Neqp = Esys.Npipe
			Esys.Npipe++

			if Pipeca.Type == PIPE_PDT {
				C.Nout = 1
				C.Nin = 1
			} else {
				C.Nout = 2
				C.Nin = 2
				C.Airpathcpy = 'y'
			}

			return true
		}
	}

	for i := 0; i < Ecat.Nstankca; i++ {
		Stankca := Ecat.Stankca[i]
		if catname == Stankca.name {
			C.Eqptype = STANK_TYPE
			C.Ncat = i
			C.Neqp = Esys.Nstank
			Esys.Nstank++

			return true
		}
	}

	for i := 0; i < Ecat.Nhexca; i++ {
		Hexca := &Ecat.Hexca[i]
		if catname == Hexca.Name {
			C.Eqptype = HEXCHANGR_TYPE
			C.Ncat = i
			C.Neqp = Esys.Nhex
			Esys.Nhex++

			C.Nout = 2
			C.Nin = 2
			C.Idi = []ELIOType{ELIO_C, ELIO_H} // CH
			C.Ido = []ELIOType{ELIO_C, ELIO_H} // CH

			return true
		}
	}

	for i := 0; i < Ecat.Npumpca; i++ {
		Pumpca := &Ecat.Pumpca[i]
		if catname == Pumpca.name {
			C.Eqptype = PUMP_TYPE
			C.Ncat = i
			C.Neqp = Esys.Npump
			Esys.Npump++

			if Pumpca.pftype == PUMP_PF {
				C.Nout = 1
				C.Nin = 1
			} else {
				C.Nout = 2
				C.Nin = 2
				C.Airpathcpy = 'y'
			}

			return true
		}
	}

	/*---- Satoh Debug VAV  2000/12/5 ----*/
	for i := 0; i < Ecat.Nvavca; i++ {
		Vavca := &Ecat.Vavca[i]
		if catname == Vavca.Name {
			if Vavca.Type == VAV_PDT {
				C.Eqptype = VAV_TYPE
			} else {
				C.Eqptype = VWV_TYPE
			}

			C.Ncat = i
			C.Neqp = Esys.Nvav
			Esys.Nvav++

			if Vavca.Type == VAV_PDT {
				C.Nout = 2
				C.Nin = 2
				// 温湿度計算のために出入り口数は2
				C.Airpathcpy = 'y'
			} else {
				C.Nout = 1
				C.Nin = 1
			}

			return true
		}
	}

	// Satoh OMVAV  2010/12/16
	for i := 0; i < Ecat.Nomvavca; i++ {
		OMvavca := &Ecat.OMvavca[i]
		if catname == OMvavca.Name {
			C.Eqptype = OMVAV_TYPE
			C.Ncat = i
			C.Neqp = Esys.Nomvav
			Esys.Nomvav++
			C.Nout = 0
			C.Nin = 0

			return true
		}
	}

	for i := 0; i < Ecat.Nstheatca; i++ {
		Stheatca := &Ecat.Stheatca[i]
		if catname == Stheatca.Name {
			C.Eqptype = STHEAT_TYPE
			C.Ncat = i
			C.Neqp = Esys.Nstheat
			Esys.Nstheat++

			C.Nout = 3
			C.Nin = 3
			// 温湿度計算のために出入り口数は2
			C.Airpathcpy = 'y'

			return true
		}
	}

	// Satoh追加　デシカント槽　2013/10/23
	for i := 0; i < Ecat.Ndesica; i++ {
		Desica := &Ecat.Desica[i]
		if catname == Desica.name {
			C.Eqptype = DESI_TYPE
			C.Ncat = i
			C.Neqp = Esys.Ndesi
			Esys.Ndesi++

			C.Nout = 2
			C.Nin = 2
			// 温湿度計算のために出入り口数は2
			C.Airpathcpy = 'y'

			return true
		}
	}

	// Satoh追加　気化冷却器　2013/10/27
	for i := 0; i < Ecat.Nevacca; i++ {
		Evcaca := &Ecat.Evacca[i]
		if catname == Evcaca.name {
			C.Eqptype = EVAC_TYPE
			C.Ncat = i
			C.Neqp = Esys.Nevac
			Esys.Nevac++
			C.Airpathcpy = 'y'
			C.Nout = 4
			C.Nin = 4
			C.Idi = []ELIOType{ELIO_D, ELIO_d, ELIO_W, ELIO_w} // DdWw
			C.Ido = []ELIOType{ELIO_D, ELIO_d, ELIO_W, ELIO_w} // DdWw

			return true
		}
	}

	for i := 0; i < Ecat.Nthexca; i++ {
		Thexca := &Ecat.Thexca[i]
		if catname == Thexca.Name {
			C.Eqptype = THEX_TYPE
			C.Ncat = i
			C.Neqp = Esys.Nthex
			Esys.Nthex++
			C.Airpathcpy = 'y'
			C.Nout = 4
			C.Nin = 4
			C.Idi = []ELIOType{ELIO_E, ELIO_e, ELIO_O, ELIO_o} // EeOo
			C.Ido = []ELIOType{ELIO_E, ELIO_e, ELIO_O, ELIO_o} // EeOo

			return true
		}
	}

	return false
}

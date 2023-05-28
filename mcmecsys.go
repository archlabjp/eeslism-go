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

/*  mecsys.c  */

package main

/*  システム使用機器の初期設定  */

func Mecsinit(dTM float64, Eqsys *EQSYS, Simc *SIMCONTL, Ncompnt int, Compnt []COMPNT, Nexsf int, Exsf []EXSF, Wd *WDAT, Rmvls *RMVLS) {
	Refaint(Eqsys.Nrefa, Eqsys.Refa, Wd, Ncompnt, Compnt)
	Collint(Eqsys.Ncoll, Eqsys.Coll, Nexsf, Exsf, Wd)
	Pipeint(Eqsys.Npipe, Eqsys.Pipe, Simc, Ncompnt, Compnt, Wd)
	Stankint(dTM, Eqsys.Nstank, Eqsys.Stank, Simc, Ncompnt, Compnt, Wd)
	Pumpint(Eqsys.Npump, Eqsys.Pump, Nexsf, Exsf)
	Stheatint(Eqsys.Nstheat, Eqsys.Stheat, Simc, Ncompnt, Compnt, Wd, Rmvls.Npcm, Rmvls.PCM)
	Flinint(Eqsys.Nflin, Eqsys.Flin, Simc, Ncompnt, Compnt, Wd)
	VWVint(Eqsys.Nvav, Eqsys.Vav, Ncompnt, Compnt)
	Thexint(Eqsys.Nthex, Eqsys.Thex)
	PVint(Eqsys.Npv, Eqsys.PVcmp, Nexsf, Exsf, Wd)

	// Satoh追加　デシカント槽 2013/10/23
	Desiint(Eqsys.Ndesi, Eqsys.Desi, Simc, Ncompnt, Compnt, Wd)

	// Satoh追加　気化冷却器　2013/10/31
	Evacint(Eqsys.Nevac, Eqsys.Evac)
}

/*  システム使用機器特性式係数の計算  */

func Mecscf(Eqsys *EQSYS) {
	Cnvrgcfv(Eqsys.Ncnvrg, Eqsys.Cnvrg)
	Hccdwint(Eqsys.Nhcc, Eqsys.Hcc)
	Hcccfv(Eqsys.Nhcc, Eqsys.Hcc)
	Boicfv(Eqsys.Nboi, Eqsys.Boi)
	Collcfv(Eqsys.Ncoll, Eqsys.Coll)
	Refacfv(Eqsys.Nrefa, Eqsys.Refa)
	Pipecfv(Eqsys.Npipe, Eqsys.Pipe)
	Hexcfv(Eqsys.Nhex, Eqsys.Hex)
	Pumpcfv(Eqsys.Npump, Eqsys.Pump)
	/*---- Satoh Debug VAV  2000/12/5 ----*/
	VAVcfv(Eqsys.Nvav, Eqsys.Vav)
	Stheatcfv(Eqsys.Nstheat, Eqsys.Stheat)
	Thexcfv(Eqsys.Nthex, Eqsys.Thex)
	// Satoh追加　デシカント槽　2013/10/23
	Desicfv(Eqsys.Ndesi, Eqsys.Desi)
	// Satoh追加　気化冷却器　2013/10/26
	Evaccfv(Eqsys.Nevac, Eqsys.Evac)
}

/*  システム使用機器の供給熱量、エネルギーの計算  */

func Mecsene(Eqsys *EQSYS) {
	Hccene(Eqsys.Nhcc, Eqsys.Hcc)
	Collene(Eqsys.Ncoll, Eqsys.Coll)
	Refaene2(Eqsys.Nrefa, Eqsys.Refa)
	Pipeene(Eqsys.Npipe, Eqsys.Pipe)
	Hexene(Eqsys.Nhex, Eqsys.Hex)
	Stankene(Eqsys.Nstank, Eqsys.Stank)
	Pumpene(Eqsys.Npump, Eqsys.Pump)
	Stheatene(Eqsys.Nstheat, Eqsys.Stheat)
	// Satoh追加　デシカント槽　2013/10/23
	Desiene(Eqsys.Ndesi, Eqsys.Desi)
	Thexene(Eqsys.Nthex, Eqsys.Thex)
	Qmeasene(Eqsys.Nqmeas, Eqsys.Qmeas)
	PVene(Eqsys.Npv, Eqsys.PVcmp)
}

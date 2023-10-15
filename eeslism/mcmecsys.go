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

package eeslism

// システム使用機器の初期設定
func Mecsinit(Eqsys *EQSYS, Simc *SIMCONTL, Compnt []COMPNT, Nexsf int, Exsf []EXSF, Wd *WDAT, Rmvls *RMVLS) {
	// ヒートポンプ
	Refaint(Eqsys.Nrefa, Eqsys.Refa, Wd, Compnt)

	// 太陽熱集熱器
	Collint(Eqsys.Ncoll, Eqsys.Coll, Nexsf, Exsf, Wd)

	// 配管・ダクト
	Pipeint(Eqsys.Npipe, Eqsys.Pipe, Simc, Compnt, Wd)

	// 蓄熱槽
	Stankint(Eqsys.Nstank, Eqsys.Stank, Simc, Compnt, Wd)

	// 定流量ポンプ、変流量ポンプおよび太陽電池駆動ポンプ
	Pumpint(Eqsys.Npump, Eqsys.Pump, Nexsf, Exsf)

	// 電気蓄熱暖房器
	Stheatint(Eqsys.Nstheat, Eqsys.Stheat, Simc, Compnt, Wd, Rmvls.Npcm, Rmvls.PCM)

	// 境界条件設定用仮想機器
	Flinint(Eqsys.Nflin, Eqsys.Flin, Simc, Compnt, Wd)

	// VAVユニット
	VWVint(Eqsys.Nvav, Eqsys.Vav, Compnt)

	// 全熱交換器
	Thexint(Eqsys.Nthex, Eqsys.Thex)

	// 太陽電池
	PVint(Eqsys.Npv, Eqsys.PVcmp, Nexsf, Exsf, Wd)

	// デシカント槽
	Desiint(Eqsys.Ndesi, Eqsys.Desi, Simc, Compnt, Wd)

	// 気化冷却器
	Evacint(Eqsys.Nevac, Eqsys.Evac)
}

// システム使用機器特性式係数の計算
func Mecscf(Eqsys *EQSYS) {
	Cnvrgcfv(Eqsys.Ncnvrg, Eqsys.Cnvrg)

	// 冷温水コイル
	Hccdwint(Eqsys.Nhcc, Eqsys.Hcc)
	Hcccfv(Eqsys.Nhcc, Eqsys.Hcc)

	// ボイラー
	Boicfv(Eqsys.Nboi, Eqsys.Boi)

	// 太陽熱集熱器
	Collcfv(Eqsys.Ncoll, Eqsys.Coll)

	// ヒートポンプ
	Refacfv(Eqsys.Nrefa, Eqsys.Refa)

	// 配管
	Pipecfv(Eqsys.Npipe, Eqsys.Pipe)

	// 熱交換器
	Hexcfv(Eqsys.Nhex, Eqsys.Hex)

	// 定流量ポンプ、変流量ポンプおよび太陽電池駆動ポンプ
	Pumpcfv(Eqsys.Npump, Eqsys.Pump)

	// VAVユニット
	VAVcfv(Eqsys.Nvav, Eqsys.Vav)

	// 蓄熱槽
	Stheatcfv(Eqsys.Nstheat, Eqsys.Stheat)

	// 全熱交換器
	Thexcfv(Eqsys.Nthex, Eqsys.Thex)

	// デシカント槽
	Desicfv(Eqsys.Ndesi, Eqsys.Desi)

	// 気化冷却器
	Evaccfv(Eqsys.Nevac, Eqsys.Evac)
}

// システム使用機器の供給熱量、エネルギーの計算
func Mecsene(Eqsys *EQSYS) {
	// 冷温水コイル
	Hccene(Eqsys.Nhcc, Eqsys.Hcc)

	// 太陽熱集熱器
	Collene(Eqsys.Ncoll, Eqsys.Coll)

	// ヒートポンプ
	Refaene2(Eqsys.Nrefa, Eqsys.Refa)

	// 配管
	Pipeene(Eqsys.Npipe, Eqsys.Pipe)

	// 熱交換器
	Hexene(Eqsys.Nhex, Eqsys.Hex)

	// 蓄熱槽
	Stankene(Eqsys.Nstank, Eqsys.Stank)

	// ポンプ
	Pumpene(Eqsys.Npump, Eqsys.Pump)

	// 電気蓄熱暖房器
	Stheatene(Eqsys.Nstheat, Eqsys.Stheat)

	// デシカント槽
	Desiene(Eqsys.Ndesi, Eqsys.Desi)

	// 全熱交換器
	Thexene(Eqsys.Nthex, Eqsys.Thex)

	Qmeasene(Eqsys.Nqmeas, Eqsys.Qmeas)

	// 太陽電池
	PVene(Eqsys.Npv, Eqsys.PVcmp)
}

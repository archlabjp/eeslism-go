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
func Mecsinit(Eqsys *EQSYS, Simc *SIMCONTL, Compnt []*COMPNT, Nexsf int, Exsf []EXSF, Wd *WDAT, Rmvls *RMVLS) {
	// ヒートポンプ
	Refaint(Eqsys.Refa, Wd, Compnt)

	// 太陽熱集熱器
	Collint(Eqsys.Coll, Nexsf, Exsf, Wd)

	// 配管・ダクト
	Pipeint(Eqsys.Pipe, Simc, Compnt, Wd)

	// 蓄熱槽
	Stankint(Eqsys.Stank, Simc, Compnt, Wd)

	// 定流量ポンプ、変流量ポンプおよび太陽電池駆動ポンプ
	Pumpint(Eqsys.Pump, Nexsf, Exsf)

	// 電気蓄熱暖房器
	Stheatint(Eqsys.Stheat, Simc, Compnt, Wd, Rmvls.Npcm, Rmvls.PCM)

	// 境界条件設定用仮想機器
	Flinint(Eqsys.Flin, Simc, Compnt, Wd)

	// VAVユニット
	VWVint(Eqsys.Vav, Compnt)

	// 全熱交換器
	Thexint(Eqsys.Thex)

	// 太陽電池
	PVint(Eqsys.PVcmp, Nexsf, Exsf, Wd)

	// デシカント槽
	Desiint(Eqsys.Desi, Simc, Compnt, Wd)

	// 気化冷却器
	Evacint(Eqsys.Evac)
}

// システム使用機器特性式係数の計算
func Mecscf(Eqsys *EQSYS) {
	// 合流要素
	Cnvrgcfv(Eqsys.Cnvrg)

	// 冷温水コイル
	Hccdwint(Eqsys.Hcc)
	Hcccfv(Eqsys.Hcc)

	// ボイラー
	Boicfv(Eqsys.Boi)

	// 太陽熱集熱器
	Collcfv(Eqsys.Coll)

	// ヒートポンプ
	Refacfv(Eqsys.Refa)

	// 配管
	Pipecfv(Eqsys.Pipe)

	// 熱交換器
	Hexcfv(Eqsys.Hex)

	// 定流量ポンプ、変流量ポンプおよび太陽電池駆動ポンプ
	Pumpcfv(Eqsys.Pump)

	// VAVユニット
	VAVcfv(Eqsys.Vav)

	// 蓄熱槽
	Stheatcfv(Eqsys.Stheat)

	// 全熱交換器
	Thexcfv(Eqsys.Thex)

	// デシカント槽
	Desicfv(Eqsys.Desi)

	// 気化冷却器
	Evaccfv(Eqsys.Evac)
}

// システム使用機器の供給熱量、エネルギーの計算
func Mecsene(Eqsys *EQSYS) {
	// 冷温水コイル
	Hccene(Eqsys.Hcc)

	// 太陽熱集熱器
	Collene(Eqsys.Coll)

	// ヒートポンプ
	Refaene2(Eqsys.Refa)

	// 配管
	Pipeene(Eqsys.Pipe)

	// 熱交換器
	Hexene(Eqsys.Hex)

	// 蓄熱槽
	Stankene(Eqsys.Stank)

	// ポンプ
	Pumpene(Eqsys.Pump)

	// 電気蓄熱暖房器
	Stheatene(Eqsys.Stheat)

	// デシカント槽
	Desiene(Eqsys.Desi)

	// 全熱交換器
	Thexene(Eqsys.Thex)

	// カロリーメータ
	Qmeasene(Eqsys.Qmeas)

	// 太陽電池
	PVene(Eqsys.PVcmp)
}

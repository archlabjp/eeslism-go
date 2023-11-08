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

/*  eqcadat.c  */

package eeslism

import (
	"fmt"
	"strings"
)

/* ----------------------------------------- */

// 機器仕様入力
func Eqcadata(f *EeTokens, Eqcat *EQCAT) {
	var (
		s  string
		ce int
		E  string
	)

	dsn := "Eqcadata"

	// 各種別のカタログの領域を初期化
	Eqcat.Hccca = make([]*HCCCA, 0)
	Eqcat.Boica = make([]*BOICA, 0)
	Eqcat.Refaca = make([]*REFACA, 0)
	Eqcat.Collca = make([]*COLLCA, 0)
	Eqcat.PVca = make([]*PVCA, 0)
	Eqcat.Pipeca = make([]*PIPECA, 0)
	Eqcat.Stankca = make([]*STANKCA, 0)
	Eqcat.Hexca = make([]*HEXCA, 0)
	Eqcat.Pumpca = make([]*PUMPCA, 0)
	Eqcat.Vavca = make([]*VAVCA, 0)
	Eqcat.Stheatca = make([]*STHEATCA, 0)
	Eqcat.Thexca = make([]*THEXCA, 0)
	Eqcat.OMvavca = make([]*OMVAVCA, 0)
	Eqcat.Desica = make([]*DESICA, 0)
	Eqcat.Evacca = make([]*EVACCA, 0)

	// 圧縮機特性リストを reflist.efl から読み取る
	Eqcat.Rfcmp = Refcmpdat()

	// ポンプ・ファンの部分負荷特性の近似式係数 を pumpfanlst.efl から読み取る
	Eqcat.Pfcmp = PFcmpdata()

	E = fmt.Sprintf(ERRFMT, dsn)

	for f.IsEnd() == false {
		s = f.GetToken()
		if s[0] == '*' {
			break
		}

		eqpType := EqpType(s)

		if eqpType == HCCOIL_TYPE {
			Hccca := new(HCCCA)
			Hccca.name = ""
			for f.IsEnd() == false {
				s = f.GetToken()
				if s[0] == ';' {
					break
				}
				if ce = strings.IndexRune(s, ';'); ce != -1 {
					s = s[:ce]
				}
				if Hccdata(s, Hccca) != 0 {
					fmt.Printf("%s %s\n", E, s)
				}
				if ce != -1 {
					break
				}
			}

			Eqcat.Hccca = append(Eqcat.Hccca, Hccca)
		} else if eqpType == BOILER_TYPE {
			Boica := new(BOICA)
			Boica.name = ""
			for f.IsEnd() == false {
				s = f.GetToken()
				if s[0] == ';' {
					break
				}
				if ce = strings.IndexRune(s, ';'); ce != -1 {
					s = s[:ce]
				}
				if Boidata(s, Boica) != 0 {
					fmt.Printf("%s %s\n", E, s)
				}
				if ce != -1 {
					break
				}
			}
			Eqcat.Boica = append(Eqcat.Boica, Boica)
		} else if eqpType == COLLECTOR_TYPE || eqpType == ACOLLECTOR_TYPE {
			Collca := new(COLLCA)
			Collca.name = ""
			Collca.Fd = 0.9
			ReadCatalogData(f, func(ss string) int { return Colldata(eqpType, ss, Collca) }, s, E)
			Eqcat.Collca = append(Eqcat.Collca, Collca)
		} else if eqpType == PV_TYPE {
			PVca := new(PVCA)
			PVca.Name = ""
			ReadCatalogData(f, func(ss string) int { return PVcadata(ss, PVca) }, s, E)
			Eqcat.PVca = append(Eqcat.PVca, PVca)
		} else if eqpType == REFACOMP_TYPE {
			Refaca := new(REFACA)
			Refaca.name = ""
			ReadCatalogData(f, func(ss string) int { return Refadata(ss, Refaca, Eqcat.Rfcmp) }, s, E)
			Eqcat.Refaca = append(Eqcat.Refaca, Refaca)
		} else if eqpType == PIPEDUCT_TYPE || eqpType == DUCT_TYPE {
			Pipeca := new(PIPECA)
			Pipeca.name = ""
			ReadCatalogData(f, func(ss string) int { return Pipedata(eqpType, ss, Pipeca) }, s, E)
			Eqcat.Pipeca = append(Eqcat.Pipeca, Pipeca)
		} else if eqpType == STANK_TYPE {
			Stankca := new(STANKCA)
			Stankca.name = ""
			ReadCatalogData(f, func(ss string) int { return Stankdata(f, ss, Stankca) }, s, E)
			Eqcat.Stankca = append(Eqcat.Stankca, Stankca)
		} else if eqpType == HEXCHANGR_TYPE {
			Hexca := new(HEXCA)
			Hexca.Name = ""
			ReadCatalogData(f, func(ss string) int { return Hexdata(ss, Hexca) }, s, E)
			Eqcat.Hexca = append(Eqcat.Hexca, Hexca)
		} else if eqpType == PUMP_TYPE || eqpType == FAN_TYPE {
			Pumpca := new(PUMPCA)
			Pumpca.name = ""
			Pumpca.Type = ""
			Pumpca.val = nil
			Pumpca.pfcmp = nil
			ReadCatalogData(f, func(ss string) int { return Pumpdata(eqpType, ss, Pumpca, Eqcat.Pfcmp) }, s, E)
			Eqcat.Pumpca = append(Eqcat.Pumpca, Pumpca)
		} else if eqpType == VAV_TYPE || eqpType == VWV_TYPE {
			vavca := new(VAVCA)
			vavca.dTset = -999.0
			vavca.Name = ""
			ReadCatalogData(f, func(ss string) int { return VAVdata(eqpType, ss, vavca) }, s, E)
			Eqcat.Vavca = append(Eqcat.Vavca, vavca)
		} else if eqpType == OMVAV_TYPE || eqpType == OAVAV_TYPE {
			OMvavca := new(OMVAVCA)
			OMvavca.Name = ""
			OMvavca.Gmax = -999.0
			OMvavca.Gmin = -999.0
			ReadCatalogData(f, func(ss string) int { return OMVAVdata(ss, OMvavca) }, s, E)
			Eqcat.OMvavca = append(Eqcat.OMvavca, OMvavca)
		} else if eqpType == STHEAT_TYPE {
			stheatca := new(STHEATCA)
			stheatca.Name = ""
			stheatca.PCMName = ""
			ReadCatalogData(f, func(ss string) int { return Stheatdata(ss, stheatca) }, s, E)
			Eqcat.Stheatca = append(Eqcat.Stheatca, stheatca)
		} else if eqpType == THEX_TYPE {
			Thexca := new(THEXCA)
			Thexca.Name = ""
			Thexca.et = -999.0
			ReadCatalogData(f, func(ss string) int { return Thexdata(ss, Thexca) }, s, E)
			Eqcat.Thexca = append(Eqcat.Thexca, Thexca)
		} else if eqpType == DESI_TYPE {
			Desica := new(DESICA)
			Desica.name = ""
			ReadCatalogData(f, func(ss string) int { return Desiccantdata(ss, Desica) }, s, E)
			Eqcat.Desica = append(Eqcat.Desica, Desica)
		} else if eqpType == EVAC_TYPE {
			Evacca := new(EVACCA)
			Evacca.Name = ""
			ReadCatalogData(f, func(ss string) int { return Evacdata(ss, Evacca) }, s, E)
			Eqcat.Evacca = append(Eqcat.Evacca, Evacca)
		} else {
			fmt.Printf("%s %s\n", E, s)
		}
	}
}

func ReadCatalogData(f *EeTokens, reader func(string) int, s string, E string) {
	var ce int
	for f.IsEnd() == false {
		ss := f.GetToken()
		if ss[0] == ';' {
			break
		}
		if ce = strings.IndexRune(ss, ';'); ce != -1 {
			s = s[:ce]
		}
		if reader(ss) != 0 {
			fmt.Printf("%s %s\n", E, s)
		}
		if ce != -1 {
			break
		}
	}
}

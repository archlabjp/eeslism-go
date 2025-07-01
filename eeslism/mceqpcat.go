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

/*
eqpcat (Equipment Category Assignment)

この関数は、与えられた機器名称（`catname`）に基づいて、
その機器のタイプを識別し、対応する機器カテゴリ情報（`C`）と
システム全体の機器リスト（`Esys`）に登録します。
これにより、建物のエネルギーシミュレーションにおいて、
各機器の特性に応じた適切なモデルが適用されるようになります。

建築環境工学的な観点:
- **機器の分類と識別**: 建物には、冷温水コイル、ボイラー、太陽熱集熱器、ポンプ、ファン、デシカント空調機など、
  様々な種類の設備機器が導入されます。
  この関数は、機器名称からその種類を特定し、
  シミュレーションモデル内で一意に識別するための`Eqptype`（機器タイプ）と`Ncat`（カタログ番号）を設定します。
- **システム構成の定義**: 各機器がシステム内でどのような役割を果たすか（入出力ポートの数`Nout`, `Nin`、
  入出力データの種類`Idi`, `Ido`など）を定義します。
  例えば、冷温水コイルは空気の温度と湿度、水の温度を入出力として持ち、
  ボイラーは熱媒の温度を入出力として持つといった違いをモデル化します。
- **機器インスタンスの生成**: `Esys`は、シミュレーション対象となる建物に実際に設置される機器のインスタンスを管理します。
  この関数は、新しい機器が定義されるたびに、
  対応する機器タイプの新しいインスタンス（例: `NewHCC()`, `NewBOI()`）を生成し、
  `Esys`リストに追加します。
- **熱湿気同時交換の考慮 (Airpathcpy)**:
  `C.Airpathcpy = true` と設定される機器（冷温水コイル、集熱器の一部、パイプ・ダクトの一部、VAV、デシカント、気化冷却器、熱交換器など）は、
  空気の温度だけでなく、湿度も同時に変化させる熱湿気同時交換を行う機器であることを示唆します。
  これは、室内空気質や潜熱負荷の計算において重要な情報となります。
- **エネルギーシステム全体の統合**: この関数は、
  個々の機器を建物全体のエネルギーシステムモデルに統合するための重要なステップです。
  各機器の特性が適切に定義されることで、
  システム全体のエネルギー消費量、熱負荷、および室内環境の予測精度が向上します。

この関数は、建物のエネルギーシミュレーションにおいて、
多様な設備機器を正確にモデル化し、
システム全体のエネルギー消費量予測、省エネルギー対策の検討、
および最適な設備システム設計を行うための重要な役割を果たします。
*/
func eqpcat(catname string, C *COMPNT, Ecat *EQCAT, Esys *EQSYS) bool {
	C.Airpathcpy = false
	C.Idi = nil
	C.Ido = nil

	for i, Hccca := range Ecat.Hccca {
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

	for i, Boica := range Ecat.Boica {
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

	for i, Collca := range Ecat.Collca {
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

	for i, PVca := range Ecat.PVca {
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

	for i, Refaca := range Ecat.Refaca {
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

	for i, Pipeca := range Ecat.Pipeca {
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

	for i, Stankca := range Ecat.Stankca {
		if catname == Stankca.name {
			C.Eqptype = STANK_TYPE
			C.Ncat = i
			C.Neqp = len(Esys.Stank)
			Esys.Stank = append(Esys.Stank, NewSTANK())

			return true
		}
	}

	for i, Hexca := range Ecat.Hexca {
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

	for i, Pumpca := range Ecat.Pumpca {
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
	for i, Vavca := range Ecat.Vavca {
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
	for i, OMvavca := range Ecat.OMvavca {
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

	for i, Stheatca := range Ecat.Stheatca {
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
	for i, Desica := range Ecat.Desica {
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
	for i, Evcaca := range Ecat.Evacca {
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

	for i, Thexca := range Ecat.Thexca {
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

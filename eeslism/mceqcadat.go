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

/*
Eqcadata (Equipment Catalog Data Input)

この関数は、様々な種類の設備機器（冷温水コイル、ボイラー、太陽熱集熱器、ポンプ、ファンなど）の
カタログデータを読み込み、それぞれの機器タイプに応じた構造体に格納します。
これは、建物のエネルギーシミュレーションにおいて、
多様な設備機器の性能をモデル化し、システム全体のエネルギー消費量を評価するために不可欠です。

建築環境工学的な観点:
  - **多様な設備機器のモデル化**: 建物には、空調、給湯、換気、照明、再生可能エネルギーなど、
    様々な目的の設備機器が導入されます。
    この関数は、これらの多様な機器の性能特性を統一的に管理し、
    シミュレーションモデルに組み込むための基盤を提供します。
    各機器タイプ（`HCCOIL_TYPE`, `BOILER_TYPE`, `COLLECTOR_TYPE`など）ごとに、
    対応するデータ読み込み関数（`Hccdata`, `Boidata`, `Colldata`など）を呼び出しています。
  - **機器性能のデータベース**: この関数は、実質的に建物の設備機器に関するデータベースを構築します。
    各機器の定格能力、効率、部分負荷特性、制御方法などの情報が格納され、
    シミュレーションの際に参照されます。
    これにより、特定の機器を選定した場合のエネルギー消費量や、
    システム全体の性能を評価できます。
  - **システム設計の柔軟性**: 様々な機器タイプに対応できることで、
    建物の用途や規模、地域の気候条件に応じた最適な設備システムを設計する際の柔軟性が高まります。
    例えば、高効率な機器の導入効果を検証したり、
    異なる熱源や熱搬送方式を比較検討したりすることが可能になります。
  - **部分負荷特性の考慮**: `Eqcat.Rfcmp`（圧縮機特性）や`Eqcat.Pfcmp`（ポンプ・ファンの部分負荷特性）を読み込むことで、
    機器が定格運転時だけでなく、部分負荷時にもどのように性能が変化するかをモデル化できます。
    実際の建物では、機器が定格能力で運転される時間は限られているため、
    部分負荷特性の考慮はエネルギー消費量予測の精度向上に不可欠です。

この関数は、建物のエネルギーシミュレーションにおいて、
多様な設備機器の性能を正確にモデル化し、
システム全体のエネルギー消費量予測、省エネルギー対策の検討、
および最適な設備システム設計を行うための重要な役割を果たします。
*/
func Eqcadata(f *EeTokens, Eqcat *EQCAT) {
	if Eqcat == nil {
		panic("Eqcat is nil")
	}

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
			vavca.dTset = FNAN
			vavca.Name = ""
			ReadCatalogData(f, func(ss string) int { return VAVdata(eqpType, ss, vavca) }, s, E)
			Eqcat.Vavca = append(Eqcat.Vavca, vavca)
		} else if eqpType == OMVAV_TYPE || eqpType == OAVAV_TYPE {
			OMvavca := new(OMVAVCA)
			OMvavca.Name = ""
			OMvavca.Gmax = FNAN
			OMvavca.Gmin = FNAN
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
			Thexca.et = FNAN
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

/*
ReadCatalogData (Read Equipment Catalog Data)

この関数は、汎用的なカタログデータ読み込み処理を提供します。
特定の機器タイプに依存しない共通のロジックで、
入力ファイルから機器の仕様データを解析し、対応するデータ構造に格納します。

建築環境工学的な観点:
  - **データ入力の標準化**: 建物のエネルギーシミュレーションでは、
    様々な種類の設備機器のデータを取り扱う必要があります。
    この関数は、各機器のデータ形式が類似している場合に、
    共通の読み込みロジックを適用することで、
    データ入力処理の効率化とコードの再利用性を高めます。
  - **柔軟なデータ形式への対応**: `reader func(string) int` という引数により、
    具体的なデータの解析処理を外部から注入できる（コールバック関数）ため、
    機器タイプごとに異なるデータフォーマットや解析ロジックに柔軟に対応できます。
    これにより、新しい機器タイプが追加された場合でも、
    この共通関数を変更することなく対応が可能となります。
  - **エラーハンドリング**: データの読み込み中にエラーが発生した場合（`reader(ss) != 0`）、
    エラーメッセージを出力する機能が含まれています。
    これは、入力データの不備を早期に発見し、
    シミュレーションの信頼性を確保するために重要です。
  - **入力ファイルの構造化**: この関数は、
    入力ファイルが特定の区切り文字（`;`）で機器のデータブロックを区切る構造になっていることを前提としています。
    このような構造化された入力ファイルは、
    大量の機器データを効率的に管理し、シミュレーションモデルに組み込むために有用です。

この関数は、建物のエネルギーシミュレーションにおいて、
多様な設備機器のデータを効率的かつ柔軟に読み込み、
シミュレーションモデルを構築するための重要な基盤機能を提供します。
*/
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

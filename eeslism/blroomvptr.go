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

package eeslism

import "errors"

/*
roomvptr (Room Variable Pointer Setting)

この関数は、室および関連するシステム変数、内部変数へのポインターを設定します。
これにより、シミュレーション中にこれらの変数にアクセスし、
室の熱的状態を監視したり、制御したりすることが可能になります。

建築環境工学的な観点:
- **室の熱的状態の監視**: シミュレーション中に室温（`Tr`）、絶対湿度（`xr`）、
  相対湿度（`RH`）、平均表面温度（`Tsav`）、作用温度（`Tot`）、
  エンタルピー（`hr`）などの変数にアクセスすることで、
  室内の熱的状態をリアルタイムで監視できます。
  これらの値は、居住者の快適性評価や、熱負荷計算の基礎となります。
- **快適性指標の監視**: `PMV`（予測平均申告）のような快適性指標へのポインターを設定することで、
  室内の快適性が目標範囲内にあるかを監視できます。
  PMVは、温熱環境が人間に与える快適感の度合いを数値化したもので、
  空調システムの設計や運用において重要な指標です。
- **表面温度の監視**: 各表面（壁、窓など）の温度（`Ts`）、
  平均放射温度（`Tmrt`）、相当外気温度（`Te`）へのポインターを設定することで、
  各表面の熱的挙動を詳細に監視できます。
  これにより、結露の発生予測や、日射熱取得の影響などを評価できます。
- **制御への応用**: これらの変数へのポインターは、
  室温制御や湿度制御などの空調システムの制御ロジックに利用されます。
  例えば、室温が設定値を超えた場合に冷房運転を開始するなど、
  フィードバック制御の基礎となります。

この関数は、室の熱的挙動を詳細に分析し、
快適性向上や省エネルギー対策の効果評価を行うための重要な役割を果たします。
*/
func roomvptr(Nk int, key []string, Room *ROOM) (VPTR, error) {
	var vptr VPTR
	vptr.Ptr = nil

	if Nk == 2 {
		switch string(key[1]) {
		case "Tr":
			vptr.Ptr = &Room.Tr
			vptr.Type = VAL_CTYPE
		case "xr":
			vptr.Ptr = &Room.xr
			vptr.Type = VAL_CTYPE
		case "RH":
			vptr.Ptr = &Room.RH
			vptr.Type = VAL_CTYPE
		case "PMV":
			vptr.Ptr = &Room.PMV
			vptr.Type = VAL_CTYPE
		case "Tsav":
			vptr.Ptr = &Room.Tsav
			vptr.Type = VAL_CTYPE
		case "Tot":
			vptr.Ptr = &Room.Tot
			vptr.Type = VAL_CTYPE
		case "hr":
			vptr.Ptr = &Room.hr
			vptr.Type = VAL_CTYPE
		}
	} else if Nk == 3 {
		for i := 0; i < Room.N; i++ {
			Sd := Room.rsrf[i]
			if string(key[1]) == Sd.Name {
				switch string(key[2]) {
				case "Ts":
					vptr.Ptr = &Sd.Ts
					vptr.Type = VAL_CTYPE
				case "Tmrt":
					vptr.Ptr = &Sd.Tmrt
					vptr.Type = VAL_CTYPE
				case "Te":
					vptr.Ptr = &Sd.Tcole
					vptr.Type = VAL_CTYPE
				}
			}
		}
	}

	if vptr.Ptr == nil {
		return vptr, errors.New("roomvptr error")
	}

	return vptr, nil
}

/* ------------------------------------------- */

/*
roomldptr (Room Load Pointer Setting)

この関数は、室の負荷計算において、
制御対象となるパラメータ（室温、絶対湿度、相対湿度、表面温度など）へのポインターを設定します。
これにより、室の熱負荷を特定の目標値に追従させる制御をモデル化できます。

建築環境工学的な観点:
- **熱負荷追従制御**: 室の熱負荷は常に変動するため、
  空調システムは、その負荷変動に応じて熱供給量や除湿量を調整する必要があります。
  この調整は、室温や湿度を制御したり、表面温度を制御したりすることで行われます。
- **制御対象の指定**: `key[1]`によって、
  - `Tr`: 室温
  - `Tot`: 作用温度
  - `RH`: 相対湿度
  - `Tdp`: 露点温度
  - `xr`: 絶対湿度
  - `<roomname> Ts`: 特定の表面温度
  などを制御対象とすることを意味します。
  `vptr.Ptr`は、これらの変数へのポインターを設定し、
  `vptr.Type = VAL_CTYPE`は、そのポインターが制御値であることを示します。
- **フィードバック制御の基礎**: このポインター設定は、
  室のフィードバック制御の基礎となります。
  シミュレーションの各時間ステップで、
  現在の室温や湿度と目標値を比較し、その差に基づいて空調システムの運転を調整します。
  これにより、室内温湿度環境の安定化や、熱供給の効率化を図ることができます。

この関数は、室の制御ロジックをモデル化し、
熱負荷変動に対する空調システムの応答をシミュレーションするために不可欠な役割を果たします。
*/
func roomldptr(load *ControlSWType, key []string, Room *ROOM, idmrk *byte) (VPTR, error) {
	var err error
	var i int
	var Sd *RMSRF
	var vptr VPTR

	if key[1] == "Tr" {
		vptr.Ptr = &Room.rmld.Tset
		vptr.Type = VAL_CTYPE
		Room.rmld.loadt = load
		Room.rmld.tropt = 'a'
		*idmrk = 't'
	} else if key[1] == "Tot" {
		vptr.Ptr = &Room.rmld.Tset
		vptr.Type = VAL_CTYPE
		Room.rmld.loadt = load
		Room.rmld.tropt = 'o'
		*idmrk = 't'
	} else if key[1] == "RH" {
		vptr.Ptr = &Room.rmld.Xset
		vptr.Type = VAL_CTYPE
		Room.rmld.loadx = load
		Room.rmld.hmopt = 'r'
		*idmrk = 'x'
	} else if key[1] == "Tdp" {
		vptr.Ptr = &Room.rmld.Xset
		vptr.Type = VAL_CTYPE
		Room.rmld.loadx = load
		Room.rmld.hmopt = 'd'
		*idmrk = 'x'
	} else if key[1] == "xr" {
		vptr.Ptr = &Room.rmld.Xset
		vptr.Type = VAL_CTYPE
		Room.rmld.loadx = load
		Room.rmld.hmopt = 'x'
		*idmrk = 'x'
	} else if len(key) > 2 && key[2] == "Ts" {

		for i = 0; i < Room.N; i++ {
			Sd = Room.rsrf[i]

			if Sd.Name == key[1] {
				vptr.Ptr = &Sd.Ts
				vptr.Type = VAL_CTYPE
				Room.rmld.loadt = load
				*idmrk = 't'
				err = nil
				break
			}
			err = errors.New("Surface not found: " + Sd.Name)
		}
	} else {
		err = errors.New("'Tr', 'Tot', 'RH', 'Tdp', 'xr' or '<roomname> Ts' are expected")
	}

	return vptr, err
}

/* ------------------------------------------- */

/*
roomldschd (Room Load Schedule Setting)

この関数は、室の負荷計算において、
スケジュールに基づいて運転を制御するための設定を行います。
特に、目標室温や目標湿度（相対湿度、露点温度、絶対湿度）が設定されている場合に、
その目標値に応じて空調システムの運転をON/OFFしたり、
目標値を設定したりするロジックを実装します。

建築環境工学的な観点:
- **スケジュール運転**: 室の熱負荷は、時間帯や曜日、季節によって変動します。
  空調システムをスケジュールに基づいて運転することで、
  不要な時間帯の運転を停止したり、熱需要に応じて運転モードを切り替えたりすることができ、
  エネルギー消費量の削減に貢献します。
- **目標室温・湿度の制御**: `rmld.Tset`は目標室温、`rmld.Xset`は目標湿度を示します。
  - `rmld.Tset > TEMPLIMIT` の条件は、
    目標室温が有効な範囲内にある場合に空調システムを運転することを意味します。
  - `rmld.hmopt`によって、相対湿度（`'r'`）、露点温度（`'d'`）、絶対湿度（`'x'`）のいずれかで湿度を制御します。
  `Eo.Control = LOAD_SW` は、空調システムが負荷追従運転モードであることを示し、
  `Eo.Sysv`は、空調システムの出口温度や湿度を目標値に設定します。
- **省エネルギー運転**: 目標室温・湿度を適切に設定することで、
  過剰な冷暖房や除湿を防ぎ、エネルギーの無駄を削減できます。
  例えば、外気温度が快適な時期には空調を停止したり、
  設定温度を緩和したりすることで、エネルギー消費を削減できます。
- **システム連携**: この関数は、室の運転制御が、
  空調システム（VAVシステムなど）の運転と連携して行われることを示唆します。
  これにより、建物全体の空調システムを統合的にモデル化し、
  エネルギーマネジメント戦略を評価できます。

この関数は、室のスケジュール運転と目標温湿度制御をモデル化し、
熱負荷変動に対する空調システムの応答、
およびエネルギー消費量をシミュレーションするために不可欠な役割を果たします。
*/
func roomldschd(Room *ROOM) {
	var Eo *ELOUT
	var rmld *RMLOAD

	if rmld = Room.rmld; rmld != nil {
		Eo = Room.cmp.Elouts[0]
		if rmld.loadt != nil {
			if Eo == Eo.Eldobj || Eo.Eldobj.Control != OFF_SW {
				if rmld.Tset > TEMPLIMIT {
					Eo.Sysv = rmld.Tset
					Room.Tr = rmld.Tset
					Eo.Control = LOAD_SW
				} else {
					if Room.VAVcontrl != nil {
						Room.VAVcontrl.Cmp.Control = OFF_SW
						Room.VAVcontrl.Cmp.Elouts[0].Control = OFF_SW
					}
				}
			}
		}

		Eo = Room.cmp.Elouts[1]
		if rmld.loadx != nil {
			if Eo == Eo.Eldobj || Eo.Eldobj.Control != OFF_SW {
				switch rmld.hmopt {
				case 'r':
					// 相対湿度
					if rmld.Xset > 0.0 {
						Eo.Sysv = FNXtr(Room.Tr, rmld.Xset)
						Eo.Control = LOAD_SW
					}
				case 'd':
					// 露点温度
					if rmld.Xset > TEMPLIMIT {
						Eo.Sysv = FNXp(FNPws(rmld.Xset))
						Eo.Control = LOAD_SW
					}
				case 'x':
					// 絶対湿度
					if rmld.Xset > 0.0 {
						Eo.Sysv = rmld.Xset
						Eo.Control = LOAD_SW
					}
				}
			}
		}
	}
}

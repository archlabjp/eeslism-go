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


/*
nennkann (Day of Year Calculation)

この関数は、与えられた月（`M`）と日（`D`）から、
その日が1月1日から数えて何日目にあたるか（通日）を計算します。

建築環境工学的な観点:
- **時間管理の標準化**: 建物のエネルギーシミュレーションでは、
  年間を通じて変化する気象データやスケジュールデータを扱う際に、
  日付を数値的に表現する「通日」が広く用いられます。
  この関数は、月日を通日に変換することで、
  時間管理を標準化し、シミュレーションのデータ処理を容易にします。
- **季節変化のモデル化**: 太陽位置、外気温度、日射量などの気象要素は、
  年間を通じて周期的に変化します。
  通日を用いることで、これらの季節変化を正確にモデル化し、
  建物の熱負荷やエネルギー消費量の季節変動を予測できます。
- **スケジュール制御の基礎**: 空調システムや照明システムなどの運転スケジュールは、
  通常、月日や曜日によって定義されます。
  通日を用いることで、これらのスケジュールをシミュレーションに組み込み、
  建物の運用パターンを正確に再現できます。

この関数は、建物のエネルギーシミュレーションにおいて、
時間管理を標準化し、
気象データやスケジュールデータの処理を効率化するための基礎的な役割を果たします。
*/
func nennkann(M, D int) int {
	var n int

	if M == 1 {
		n = D
	} else if M == 2 {
		n = 31 + D
	} else if M == 3 {
		n = 31 + 28 + D
	} else if M == 4 {
		n = 31 + 28 + 31 + D
	} else if M == 5 {
		n = 31 + 28 + 31 + 30 + D
	} else if M == 6 {
		n = 31 + 28 + 31 + 30 + 31 + D
	} else if M == 7 {
		n = 31 + 28 + 31 + 30 + 31 + 30 + D
	} else if M == 8 {
		n = 31 + 28 + 31 + 30 + 31 + 30 + 31 + D
	} else if M == 9 {
		n = 31 + 28 + 31 + 30 + 31 + 30 + 31 + 31 + D
	} else if M == 10 {
		n = 31 + 28 + 31 + 30 + 31 + 30 + 31 + 31 + 30 + D
	} else if M == 11 {
		n = 31 + 28 + 31 + 30 + 31 + 30 + 31 + 31 + 30 + 31 + D
	} else if M == 12 {
		n = 31 + 28 + 31 + 30 + 31 + 30 + 31 + 31 + 30 + 31 + 30 + D
	}

	return n
}

/*-----------------------------------------------------*/

/* ================================================================

 SUNLIB

  太陽位置および日射量計算用ライブラリ－
  （宇田川、パソコンによる空気調和計算法、プログラム4.1の C 言語版, ANSI C 版）

---------------------------------------------------------------- */

package eeslism

/*
FNNday (Function for Day Number Calculation)

この関数は、与えられた月日（`Mo`月`Nd`日）が1月1日から数えて何日目（通日）にあたるかを計算します。
通日は、太陽位置や日射量など、年間を通じて変化する気象データを扱う際に、
日付を数値的に表現するための基本的な情報となります。

建築環境工学的な観点:
- **年間を通じた日射・日照の変動**: 建物の日射熱取得量や日照時間は、
  季節によって大きく変動します。
  例えば、夏至（通日約172日）の頃は日射量が最大となり、
  冬至（通日約355日）の頃は日照時間が短くなります。
  これらの季節変化をモデル化する上で、通日は日付を統一的に扱うための重要なインデックスとなります。
- **気象データの同期**: エネルギーシミュレーションでは、
  日射量、外気温度、湿度などの気象データを時間単位で用います。
  これらの気象データは通常、通日や時刻と関連付けられており、
  この関数で計算される通日は、気象データとシミュレーション時刻を正確に同期させるために利用されます。
- **設計条件の設定**: 建物の設計段階では、特定の季節や日付における日射条件を考慮して、
  窓の配置、日射遮蔽部材の設計、太陽光発電システムの容量などを検討します。
  例えば、「夏至の正午における日射熱取得量」や「冬至の午前中の日照時間」といった条件を設定する際に、
  この通日計算が基礎となります。

この関数は、建物のエネルギーシミュレーションにおいて、
年間を通じた日射環境や気象条件の変化を正確にモデル化し、
パッシブソーラー設計や省エネルギー設計を行うための基礎的な役割を果たします。
*/
func FNNday(Mo int, Nd int) int {
	if Mo < 1 || Nd < 1 || Mo > 13 || Nd > 31 {
		panic("適切な月日を指定してください。")
	}
	var Mo2, Mo3 int
	if Mo < 3 {
		Mo2 = 1
		Mo3 = 0
	} else {
		Mo2 = 0
		Mo3 = 1
	}
	return int((153*(Mo-1)+2*(Mo2)-9*(Mo3))/5 + Nd)
}

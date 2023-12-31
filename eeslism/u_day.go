/* ================================================================

 SUNLIB

  太陽位置および日射量計算用ライブラリ－
  （宇田川、パソコンによる空気調和計算法、プログラム4.1の C 言語版, ANSI C 版）

---------------------------------------------------------------- */

package eeslism

// Mo月No日が1月1日から数えて何日目か(通日)を返す
// 1月1日は1である。
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

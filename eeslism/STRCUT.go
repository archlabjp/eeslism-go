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

import "strings"

/*
STRCUT (String Cut from Last Index)

この関数は、与えられた文字列`DATA`の中から、
指定された部分文字列`a`が最後に出現する位置を検索し、
その位置までの部分文字列を返します。
*/
func STRCUT(DATA string, a string) string {
	idx := strings.LastIndex(DATA, a)
	if idx == -1 {
		return ""
	}
	return DATA[:idx]
}

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

// DATAの中にあるaを後ろから探して、先頭からaまでの文字列を返す
func STRCUT(DATA string, a string) string {
	idx := strings.LastIndex(DATA, a)
	if idx == -1 {
		return ""
	}
	return DATA[:idx]
}

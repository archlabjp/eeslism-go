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

/* rzone.c */
package main

import (
	"bufio"
	"io"
	"strings"
)

/* ゾーン集計実施室の指定  */

func Rzonedata(fi io.Reader, dsn string, Nroom int, Room []ROOM, Nrzone *int, _Rzone []RZONE) {
	scanner := bufio.NewScanner(fi)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "*") {
			break
		}

		fields := strings.Fields(line)

		Rzone := &_Rzone[*Nrzone]
		Rzone.name = fields[0]
		Rzone.Nroom = 0
		Rzone.Afloor = 0.0

		for _, s := range fields[1:] {
			if i := idroom(s, Room, ""); i < Nroom {
				Rm := &Room[i]
				Rzone.rm = append(Rzone.rm, Rm)
				Rzone.Nroom++
				Rzone.Afloor += Rm.FArea
			} else {
				Eprint(dsn, "<Rzinedata> room name")
			}
		}

		(*Nrzone)++
	}
}

/* -------------------------------------------------------- */

/* 室内熱環境、室負荷のゾーン集計  */

func Rzonetotal(Nrzone int, Rzone *RZONE) {
	for i := 0; i < Nrzone; i++ {
		Rzone.Tr = 0.0
		Rzone.xr = 0.0
		Rzone.RH = 0.0
		Rzone.Tsav = 0.0
		Rzone.Qhs = 0.0
		Rzone.Qhl = 0.0
		Rzone.Qht = 0.0
		Rzone.Qcs = 0.0
		Rzone.Qcl = 0.0
		Rzone.Qct = 0.0

		for j := 0; j < Rzone.Nroom; j++ {
			R := Rzone.rm[j]

			Rzone.Tr += R.Tr * R.FArea
			Rzone.xr += R.xr * R.FArea
			Rzone.RH += R.RH * R.FArea
			Rzone.Tsav += R.Tsav * R.FArea
			if R.rmld != nil {
				if R.rmld.Qs > 0.0 {
					Rzone.Qhs += R.rmld.Qs
				} else {
					Rzone.Qcs += R.rmld.Qs
				}

				if R.rmld.Ql > 0.0 {
					Rzone.Qhl += R.rmld.Ql
				} else {
					Rzone.Qcl += R.rmld.Ql
				}

				if R.rmld.Qt > 0.0 {
					Rzone.Qht += R.rmld.Qt
				} else {
					Rzone.Qct += R.rmld.Qt
				}
			}
		}

		Rzone.Tr /= Rzone.Afloor
		Rzone.xr /= Rzone.Afloor
		Rzone.RH /= Rzone.Afloor
		Rzone.Tsav /= Rzone.Afloor
	}
}

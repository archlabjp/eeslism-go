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

/*   rmprint.c   */

package eeslism

import (
	"fmt"
	"io"
)

/* ---------------------------------------------------------------- */
/* 室内表面温度の出力 */

var __Rmsfprint_ic int

func Rmsfprint(fo io.Writer, title string, Mon, Day int, time float64, Room []ROOM, Sd []RMSRF) {
	N := Room[0].end

	if __Rmsfprint_ic == 0 {
		__Rmsfprint_ic++

		var n int
		for i := 0; i < N; i++ {
			Rm := &Room[i]
			if Rm.sfpri {
				n++
			}
		}

		fmt.Fprintf(fo, "%s;\n %d\n", title, n)
		fmt.Fprint(fo, "Mo\tNd\ttime\t")

		for i := 0; i < N; i++ {
			Rm := &Room[i]
			if Rm.sfpri {
				fmt.Fprintf(fo, "%s\t", Rm.Name)

				for n := 0; n < Rm.N; n++ {
					S := &Sd[Rm.Brs+n]
					if S.Name == "" {
						fmt.Fprintf(fo, "%d-%c_Ts\t", n-Rm.Brs, S.ble)
					} else {
						fmt.Fprintf(fo, "%s_Ts\t", S.Name)
					}
				}
			}
		}
		fmt.Fprint(fo, "\n")
	}
	fmt.Fprintf(fo, "%d\t%d\t%.2f\t", Mon, Day, time)

	for i := 0; i < N; i++ {
		Rm := &Room[i]
		if Rm.sfpri {
			fmt.Fprint(fo, "\t")

			for n := 0; n < Rm.N; n++ {
				S := &Sd[Rm.Brs+n]
				fmt.Fprintf(fo, "%.1f\t", S.Ts)
			}
		}
	}
	fmt.Fprint(fo, "\n")
}

/* ---------------------------------------------------------------- */
/* 室内表面熱流の出力 */

var __Rmsfqprint_ic int

func Rmsfqprint(fo io.Writer, title string, Mon, Day int, time float64, Room []ROOM, Sd []RMSRF) {
	N := Room[0].end

	if __Rmsfqprint_ic == 0 {
		__Rmsfqprint_ic++

		var n int
		for i := 0; i < N; i++ {
			Rm := &Room[i]
			if Rm.sfpri {
				n++
			}
		}

		fmt.Fprintf(fo, "%s;\n %d\n", title, n)
		fmt.Fprint(fo, "Mo\tNd\ttime\t")

		for i := 0; i < N; i++ {
			Rm := &Room[i]
			if Rm.sfpri {
				fmt.Fprintf(fo, "%s\t", Rm.Name)

				for n := 0; n < Rm.N; n++ {
					S := &Sd[Rm.Brs+n]
					if S.Name == "" {
						fmt.Fprintf(fo, "%d-%c_Qc\t%d-%c_Qr\t%d-%c_RS\t%d-%c_Qi\t%d-%c_RSsol\t%d-%c_RSli\t%d-%c_tsol\t%d-%c_asol\t%d-%c_rn\t",
							n, S.ble, n, S.ble, n, S.ble,
							n, S.ble, n, S.ble, n, S.ble, n, S.ble, n, S.ble, n, S.ble)
					} else {
						fmt.Fprintf(fo, "%s_Qc\t%s_Qr\t%s_RS\t%s_Qi\t%s_RSsol\t%s_RSli\t%s_tsol\t%s_asol\t%s_rn\t",
							S.Name, S.Name, S.Name, S.Name, S.Name, S.Name, S.Name, S.Name, S.Name)
					}
				}
			}
		}
		fmt.Fprint(fo, "\n")
	}
	fmt.Fprintf(fo, "%d\t%d\t%.2f\t", Mon, Day, time)

	for i := 0; i < N; i++ {
		Rm := &Room[i]
		if Rm.sfpri {
			fmt.Fprint(fo, "\t")

			// 2003/9/10　表面熱取得を負とするために短波長成分RSの符号を変更した。
			for n := 0; n < Rm.N; n++ {
				S := &Sd[Rm.Brs+n]
				fmt.Fprintf(fo, "%.4e\t%.4e\t%.4e\t%.4e\t%.4e\t%.4e\t%.4e\t%.4e\t%.4e\t", S.Qc, S.Qr,
					-S.RS*S.A, S.Qi, -S.RSsol*S.A, -S.RSli*S.A, S.Qgt, S.Qga, S.Qrn)
			}
		}
	}
	fmt.Fprint(fo, "\n")
}

/* ---------------------------------------------------------------- */
/* 室内表面熱伝達率の出力 */

var __Rmsfaprint_ic int

func Rmsfaprint(fo io.Writer, title string, Mon, Day int, time float64, Room []ROOM, Sd []RMSRF) {
	N := Room[0].end

	if __Rmsfaprint_ic == 0 {
		__Rmsfaprint_ic++

		var n int
		for i := 0; i < N; i++ {
			Rm := &Room[i]
			if Rm.sfpri {
				n++
			}
		}

		fmt.Fprintf(fo, "%s;\n %d\n", title, n)
		fmt.Fprint(fo, "Mo\tNd\ttime\t")

		for i := 0; i < N; i++ {
			Rm := &Room[i]
			if Rm.sfpri {
				fmt.Fprintf(fo, "%s\t", Rm.Name)

				for nn := 0; nn < Rm.N; nn++ {
					S := &Sd[Rm.Brs+nn]
					if S.Name == "" {
						fmt.Fprintf(fo, "%d-%c_K\t%d-%c_alc\t%d-%c_alr\t",
							n-Rm.Brs, S.ble, n-Rm.Brs, S.ble, n-Rm.Brs, S.ble)
					} else {
						fmt.Fprintf(fo, "%s_K\t%s_alc\t%s_alr\t",
							S.Name, S.Name, S.Name)
					}
				}
			}
		}
		fmt.Fprint(fo, "\n")
	}
	fmt.Fprintf(fo, "%d\t%d\t%.2f\t", Mon, Day, time)

	for i := 0; i < N; i++ {
		Rm := &Room[i]
		if Rm.sfpri {
			fmt.Fprint(fo, "\t")

			for nn := 0; nn < Rm.N; nn++ {
				S := &Sd[Rm.Brs+nn]
				fmt.Fprintf(fo, "%.3g\t%.3g\t%.3g\t", S.K, S.alic, S.alir)
			}
		}
	}
	fmt.Fprint(fo, "\n")
}

/* 日積算壁体貫流熱取得の出力 */
var __Dysfprint_ic int

func Dysfprint(fo io.Writer, title string, Mon, Day int, Room []ROOM) {
	N := Room[0].end

	if __Dysfprint_ic == 0 {
		__Dysfprint_ic++

		var n int
		for i := 0; i < N; i++ {
			Rm := &Room[i]
			if Rm.sfpri {
				n++
			}
		}

		fmt.Fprintf(fo, "%s;\n %d\n", title, n)
		fmt.Fprint(fo, "Mo\tNd\t")

		for i := 0; i < N; i++ {
			Rm := &Room[i]
			if Rm.sfpri {
				fmt.Fprintf(fo, "%s\t", Rm.Name)

				for n := 0; n < Rm.N; n++ {
					Sd := &Rm.rsrf[n]
					if Sd.Name == "" {
						fmt.Fprintf(fo, "%d-%c_Ts\t%d-%c_Tsmax\t%d-%c_Tsmin\t%d-%c_Qih\t%d-%c_Qic\t",
							n, Sd.ble, n, Sd.ble, n, Sd.ble, n, Sd.ble, n, Sd.ble)
					} else {
						fmt.Fprintf(fo, "%s_Ts\t%s_Tsmax\t%s_Tsmin\t%s_Qih\t%s_Qic\t",
							Sd.Name, Sd.Name, Sd.Name, Sd.Name, Sd.Name)
					}
				}
			}
		}
		fmt.Fprint(fo, "\n")
	}

	fmt.Fprintf(fo, "%d\t%d\t", Mon, Day)

	for i := 0; i < N; i++ {
		Rm := &Room[i]
		if Rm.sfpri {
			fmt.Fprint(fo, "\t")

			for n := 0; n < Rm.N; n++ {
				Sd := &Rm.rsrf[n]
				fmt.Fprintf(fo, "%.2f\t%.2f\t%.2f\t%.3g\t%.3g\t",
					Sd.Tsdy.M, Sd.Tsdy.Mx, Sd.Tsdy.Mn, Sd.SQi.H, Sd.SQi.C)
			}
		}
	}

	fmt.Fprint(fo, "\n")
}

/* ---------------------------------------------------------------- */

/* 日よけの影面積の出力 */

var __Shdprint_ic int

func Shdprint(fo io.Writer, title string, Mon, Day int, time float64, Nsrf int, Sd []RMSRF) {
	if __Shdprint_ic == 0 {
		__Shdprint_ic++

		var m int
		for i := 0; i < Nsrf; i++ {
			Sdd := &Sd[i]
			if Sdd.shdpri && Sdd.sb >= 0 {
				m++
			}
		}

		fmt.Fprintf(fo, "%s;\n %d\n", title, m)

		for i := 0; i < Nsrf; i++ {
			Sdd := &Sd[i]
			if Sdd.shdpri && Sdd.sb >= 0 {
				fmt.Fprintf(fo, "%s\t%d:%s\n", Sdd.room.Name, i-Sdd.room.Brs, Sdd.Name)
			}
		}
	}

	fmt.Fprintf(fo, "%d\t%d\t%.2f\t", Mon, Day, time)

	for i := 0; i < Nsrf; i++ {
		Sdd := &Sd[i]
		if Sdd.shdpri && Sdd.sb >= 0 {
			fmt.Fprintf(fo, "%.2f\t", Sdd.Fsdworg)
		}
	}

	fmt.Fprint(fo, "\n")
}

/* ---------------------------------------------------------------- */

/* 壁体内部温度の出力 */

var __Wallprint_ic int

func Wallprint(fo io.Writer, title string, Mon, Day int, time float64, Nsrf int, Sd []RMSRF) {
	if __Wallprint_ic == 0 {
		__Wallprint_ic++
		var m int
		for i := 0; i < Nsrf; i++ {
			Sdd := &Sd[i]
			if Sdd.wlpri && Sdd.wd >= 0 {
				m++
			}
		}

		fmt.Fprintf(fo, "%s;\n %d\n", title, m)

		for i := 0; i < Nsrf; i++ {
			Sdd := &Sd[i]
			if Sdd.wlpri && Sdd.wd >= 0 {
				fmt.Fprintf(fo, "%s\t%d-%c:%s\t%d\n", Sdd.room.Name, i-Sdd.room.Brs, Sdd.ble, Sdd.Name, Sdd.mw.M)
			}
		}
	}

	fmt.Fprintf(fo, "%d\t%d\t%.2f\t", Mon, Day, time)

	for i := 0; i < Nsrf; i++ {
		Sdd := &Sd[i]
		if Sdd.wlpri && Sdd.wd >= 0 {
			Mw := Sdd.mw

			// 室内が壁体０側の場合
			if Sdd.mwside == RMSRFMwSideType_i {
				for m := 0; m < Mw.M; m++ {
					fmt.Fprintf(fo, "\t%.2f", Mw.Tw[m])
				}
			} else { // 室内が壁体M側の場合
				for m := Mw.M - 1; m >= 0; m-- {
					fmt.Fprintf(fo, "\t%.2f", Mw.Tw[m])
				}
			}

			fmt.Fprint(fo, "\t")
		}
	}

	fmt.Fprint(fo, "\n")
}

/* ---------------------------------------------------------------- */

/* 潜熱蓄熱材の状態値の出力 */
var __PCMprint_ic int

func PCMprint(fo io.Writer, title string, Mon, Day int, time float64, Nsrf int, Sd []RMSRF) {
	var Sdd *RMSRF
	var pcmstate *PCMSTATE

	if __PCMprint_ic == 0 {
		__PCMprint_ic++

		Sdd = &Sd[0]
		m := 0
		for i := 0; i < Nsrf; i++ {
			if Sdd.pcmpri && Sdd.wd >= 0 {
				m += Sdd.Npcm
			}
			Sdd = &Sd[i]
		}

		fmt.Fprintf(fo, "%s;\n %d\n", title, m)

		for i := 0; i < Nsrf; i++ {
			Sdd := Sd[i]
			if Sdd.mwside == RMSRFMwSideType_i {
				if Sdd.pcmpri && Sdd.wd >= 0 {
					for m := 0; m < Sdd.mw.M; m++ {
						pcmstate = Sdd.pcmstate[m]
						if pcmstate != nil && pcmstate.Name != nil {
							fmt.Fprintf(fo, "%s\t%d-%c:%s\t%s\tTpcm\tcp\tLamda\n", Sdd.room.Name, i-Sdd.room.Brs, Sdd.ble, Sdd.Name, *pcmstate.Name)
						}
					}
				}
			}
		}
	}

	fmt.Fprintf(fo, "%d\t%d\t%.2f\t", Mon, Day, time)

	Sdd = &Sd[0]
	for i := 0; i < Nsrf; i++ {
		if Sdd.pcmpri && Sdd.wd >= 0 {
			Mw := Sdd.mw

			if Sdd.mwside == RMSRFMwSideType_i {
				for j := 0; j < Mw.M; j++ {
					pcmstate = Sdd.pcmstate[j]
					if pcmstate != nil && pcmstate.Name != nil {
						fmt.Fprintf(fo, "\t%.3f\t%.3f\t%.3f\t%.0f\t%.0f\t%.4g\t%.4g",
							pcmstate.TempPCMNodeL, pcmstate.TempPCMNodeR,
							pcmstate.TempPCMave, pcmstate.CapmL, pcmstate.CapmR,
							pcmstate.LamdaL, pcmstate.LamdaR)
					}
				}
			}

			fmt.Fprintf(fo, "\t")
		}
	}

	fmt.Fprintf(fo, "\n")
}

/* ---------------------------------------------------------------- */

/* 日射、室内熱取得の出力 */

var __Qrmprint_ic int

func Qrmprint(fo io.Writer, title string, Mon, Day int, time float64, Room []ROOM, Qrm []QRM) {
	N := Room[0].end

	if __Qrmprint_ic == 0 {
		__Qrmprint_ic++

		// 日射、室内発熱取得出力指定の部屋数を数える
		var n int
		for i := 0; i < N; i++ {
			Rm := &Room[i]
			if Rm.eqpri {
				n++
			}
		}

		fmt.Fprintf(fo, "%s;\n %d\n", title, n)
		fmt.Fprint(fo, "Mo\tNd\ttt\t")

		key := [16]string{"tsol", "asol", "arn", "hums", "light", "apls",
			"huml", "apll", "Qeqp", "Qis", "Qil", "Qsto", "Qstol", "AE", "AG"}

		for i := 0; i < N; i++ {
			Rm := &Room[i]
			if Rm.eqpri {
				for j := 0; j < 16; j++ {
					fmt.Fprintf(fo, "%s_%s\t", Rm.Name, key[j])
				}
			}
		}

		fmt.Fprint(fo, "\n")
	}

	fmt.Fprintf(fo, "%d\t%d\t%.2f\t", Mon, Day, time)

	for i := 0; i < N; i++ {
		Rm := &Room[i]
		if Rm.eqpri {
			Q := &Qrm[i]
			fmt.Fprintf(fo, "%.5g\t%.5g\t%.5g\t%.5g\t%.5g\t%.5g\t%.5g\t%.5g\t%.5g\t%.5g\t%.5g\t%.5g\t",
				Q.Tsol, Q.Asol, Q.Arn, Q.Hums, Q.Light, Q.Apls, Q.Huml, Q.Apll, Q.Qeqp, Rm.QM, Q.Qinfs, Q.Qinfl)
			fmt.Fprintf(fo, "%.5g\t%.5g\t", Q.Qsto, Q.Qstol)
			fmt.Fprintf(fo, "%.5g\t%.5g\t", Q.AE, Q.AG)
		}
	}

	fmt.Fprint(fo, "\n")
}

/* ---------------------------------------------------------------- */

/* 日射、室内熱取得の出力 */

var __Dyqrmprint_ic int

func Dyqrmprint(fo io.Writer, title string, Mon int, Day int, Room []ROOM, Trdav []float64, Qrmd []QRM) {
	N := Room[0].end

	if __Dyqrmprint_ic == 0 {
		__Dyqrmprint_ic++

		var n int

		for i := 0; i < N; i++ {
			Rm := &Room[i]
			if Rm.eqpri {
				n++
			}
		}

		fmt.Fprintf(fo, "%s;\n %d\n", title, n)
		fmt.Fprintf(fo, "Mo\tNd\t")

		key := [16]string{"Tr", "tsol", "asol", "arn", "hums", "light", "apls",
			"huml", "apll", "Qeqp", "Qis", "Qil", "Qsto", "Qstol", "AE", "AG"}

		for i := 0; i < N; i++ {
			Rm := &Room[i]
			if Rm.eqpri {
				for j := 0; j < 16; j++ {
					fmt.Fprintf(fo, "%s_%s\t", Rm.Name, key[j])
				}
			}
		}

		fmt.Fprintf(fo, "\n")
	}

	fmt.Fprintf(fo, "%d\t%d\t", Mon, Day)

	for i := 0; i < N; i++ {
		if Room[i].eqpri {
			Q := &Qrmd[i]
			fmt.Fprintf(fo,
				"%.1f\t%.4g\t%.4g\t%.4g\t%.4g\t%.4g\t%.4g\t%.4g\t%.4g\t%.4g\t%.4g\t%.4g\t%.4g\t%.4g\t%.4g\t%.4g\t",
				Trdav[i], Q.Tsol, Q.Asol, Q.Arn, Q.Hums, Q.Light, Q.Apls, Q.Huml, Q.Apll,
				Q.Qeqp, Q.Qinfs, Q.Qinfl, Q.Qsto, Q.Qstol, Q.AE, Q.AG)
		}
	}

	fmt.Fprintf(fo, "\n")
}

/* ---------------------------------------------------------------- */

var __Qrmsum_oldday int

func Qrmsum(Day int, _Room []ROOM, Qrm []QRM, Trdav []float64, Qrmd []QRM) {
	if Day != __Qrmsum_oldday {
		for i := range _Room {
			Q := &Qrmd[i]
			T := &Trdav[i]

			*T = 0.0
			Q.Tsol = 0.0
			Q.Asol = 0.0
			Q.Arn = 0.0
			Q.Hums = 0.0
			Q.Light = 0.0
			Q.Apls = 0.0
			Q.Huml = 0.0
			Q.Apll = 0.0
			Q.Qeqp = 0.0
			Q.Qinfl = 0.0
			Q.Qinfs = 0.0
			Q.Qsto = 0.0
			Q.Qstol = 0.0
			Q.AE = 0.0
			Q.AG = 0.0
		}
		__Qrmsum_oldday = Day
	}

	for i := range _Room {
		Q := &Qrmd[i]
		Qr := &Qrm[i]
		T := &Trdav[i]
		Room := &_Room[i]

		scale := DTM / 3600.0

		*T += Room.Tr * scale / 24.0
		Q.Tsol += Qr.Tsol * scale
		Q.Asol += Qr.Asol * scale
		Q.Arn += Qr.Arn * scale
		Q.Hums += Qr.Hums * scale
		Q.Light += Qr.Light * scale
		Q.Apls += Qr.Apls * scale
		Q.Huml += Qr.Huml * scale
		Q.Apll += Qr.Apll * scale
		Q.Qinfs += Qr.Qinfs * scale
		Q.Qinfl += Qr.Qinfl * scale
		Q.Qeqp += Qr.Qeqp * scale
		Q.Qsto += Qr.Qsto * scale
		Q.Qstol += Qr.Qstol * scale
		Q.AE += Qr.AE * scale
		Q.AG += Qr.AG * scale
	}
}

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

/* xprmcs.c */

package main

import "fmt"

/*  システム使用機器についての検討用出力  */

func mecsxprint(Eqsys *EQSYS) {
	hccxprint(Eqsys.Nhcc, Eqsys.Hcc)
	boixprint(Eqsys.Nboi, Eqsys.Boi)
	collxprint(Eqsys.Ncoll, Eqsys.Coll)
	refaxprint(Eqsys.Nrefa, Eqsys.Refa)
	pipexprint(Eqsys.Npipe, Eqsys.Pipe)
}

/* --------------------------- */

func boixprint(Nboi int, Boi []BOI) {
	if Nboi > 0 {
		fmt.Printf("%s N=%d\n", BOILER_TYPE, Nboi)

		for i := 0; i < Nboi; i++ {
			b := &Boi[i]
			fmt.Printf("[%d] %-10s Do=%5.3f  D1=%5.3f Tin=%5.2f Tout=%5.2f Q=%4.0f E=%4.0f\n",
				i, b.Name, b.Do, b.D1, b.Tin,
				b.Cmp.Elouts[0].Sysv, b.Q, b.E)
		}
	}
}

/* ------------------------------------------ */

func hccxprint(Nhcc int, Hcc []HCC) {
	if Nhcc > 0 {
		fmt.Printf("%s N=%d\n", HCCOIL_TYPE, Nhcc)

		for i := 0; i < Nhcc; i++ {
			h := &Hcc[i]
			fmt.Printf("[%d] %-10s et=%5.3f eh=%5.3f\n", i, h.Name, h.et, h.eh)
			E := h.Et
			fmt.Printf("     Et w=%7.3f t=%7.3f x=%7.3f C=%7.3f\n", E.W, E.T, E.X, E.C)
			E = h.Ex
			fmt.Printf("     Et w=%7.3f t=%7.3f x=%7.3f C=%7.3f\n", E.W, E.T, E.X, E.C)
			E = h.Ew
			fmt.Printf("     Et w=%7.3f t=%7.3f x=%7.3f C=%7.3f\n", E.W, E.T, E.X, E.C)
			el := h.Cmp.Elouts[0]
			fmt.Printf("     Tain=%5.2f  Taout=%5.2f  Qs=%4.0f\n", h.Tain, el.Sysv, h.Qs)
			el = h.Cmp.Elouts[1]
			fmt.Printf("     xain=%5.4f  xaout=%5.4f  Qs=%4.0f\n", h.Xain, el.Sysv, h.Ql)
			el = h.Cmp.Elouts[2]
			fmt.Printf("     Wwin=%5.2f  Twout=%5.4f  Qt=%4.0f\n", h.Twin, el.Sysv, h.Qt)
		}
	}
}

/* --------------------------- */

func pipexprint(Npipe int, Pipe []PIPE) {
	var Te float64

	if Npipe > 0 {
		fmt.Printf("%s N=%d\n", PIPEDUCT_TYPE, Npipe)

		for i := 0; i < Npipe; i++ {
			p := &Pipe[i]

			if p.Cmp.Envname != "" {
				Te = *p.Tenv
			} else {
				Te = p.Room.Tot
			}

			fmt.Printf("[%d] %-10s Do=%6.3f  D1=%6.3f Tin=%5.2f Tout=%5.2f ep=%5.3f env=%4.1f Q=%4.0f\n",
				i, p.Name, p.Do, p.D1, p.Tin,
				p.Cmp.Elouts[0].Sysv, p.Ep, Te, p.Q)
		}
	}
}

/* ------------------------------------------------------------- */

func refaxprint(Nrefa int, Refa []REFA) {
	if Nrefa > 0 {
		fmt.Printf("%s N=%d\n", REFACOMP_TYPE, Nrefa)

		for i := 0; i < Nrefa; i++ {
			r := &Refa[i]
			fmt.Printf("[%d] %-10s Do=%6.3f D1=%6.3f Tin=%5.2f Tout=%5.2f Ta=%4.1f\n",
				i, r.Name, r.Do, r.D1, r.Tin,
				r.Cmp.Elouts[0].Sysv, *r.Ta)
			fmt.Printf("     Te=%5.2f  Tc=%5.2f  Q=%6.0f E=%6.0f Ph=%3.0f\n",
				r.Te, r.Tc, r.Q, r.E, r.Ph)
		}
	}
}

/* ------------------------------------------------------------- */

func collxprint(Ncoll int, Colls []COLL) {
	if Ncoll > 0 {
		fmt.Printf("%s N=%d\n", COLLECTOR_TYPE, Ncoll)

		for i := 0; i < Ncoll; i++ {
			Coll := &Colls[i]
			fmt.Printf("[%d] %-10s Do=%6.3f  D1=%6.3f Tin=%5.2f Tout=%5.2f Q=%4.0f Sol=%4.0f Te=%5.1f\n",
				i, Coll.Name, Coll.Do, Coll.D1, Coll.Tin,
				Coll.Cmp.Elouts[0].Sysv, Coll.Q, Coll.Sol, Coll.Te)
			fmt.Printf("   exs=%s  b0=%5.3f  b1=%5.3f ec=%5.3f\n", Coll.sol.Name,
				Coll.Cat.b0, Coll.Cat.b1, Coll.ec)
		}
	}
}

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

package eeslism

import "fmt"

/*  システム使用機器についての検討用出力  */

func mecsxprint(Eqsys *EQSYS) {
	hccxprint(Eqsys.Hcc)
	boixprint(Eqsys.Boi)
	collxprint(Eqsys.Coll)
	refaxprint(Eqsys.Refa)
	pipexprint(Eqsys.Pipe)
}

/* --------------------------- */

func boixprint(Boi []*BOI) {
	if len(Boi) > 0 {
		fmt.Printf("%s N=%d\n", BOILER_TYPE, len(Boi))

		for i, b := range Boi {
			fmt.Printf("[%d] %-10s Do=%5.3f  D1=%5.3f Tin=%5.2f Tout=%5.2f Q=%4.0f E=%4.0f\n",
				i, b.Name, b.Do, b.D1, b.Tin,
				b.Cmp.Elouts[0].Sysv, b.Q, b.E)
		}
	}
}

/* ------------------------------------------ */

func hccxprint(Hcc []*HCC) {
	if len(Hcc) > 0 {
		fmt.Printf("%s N=%d\n", HCCOIL_TYPE, len(Hcc))

		for i, h := range Hcc {
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

func pipexprint(Pipe []*PIPE) {
	var Te float64

	if len(Pipe) > 0 {
		fmt.Printf("%s N=%d\n", PIPEDUCT_TYPE, len(Pipe))

		for i, p := range Pipe {
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

func refaxprint(Refa []*REFA) {
	if len(Refa) > 0 {
		fmt.Printf("%s N=%d\n", REFACOMP_TYPE, len(Refa))

		for i, r := range Refa {
			fmt.Printf("[%d] %-10s Do=%6.3f D1=%6.3f Tin=%5.2f Tout=%5.2f Ta=%4.1f\n",
				i, r.Name, r.Do, r.D1, r.Tin,
				r.Cmp.Elouts[0].Sysv, *r.Ta)
			fmt.Printf("     Te=%5.2f  Tc=%5.2f  Q=%6.0f E=%6.0f Ph=%3.0f\n",
				r.Te, r.Tc, r.Q, r.E, r.Ph)
		}
	}
}

/* ------------------------------------------------------------- */

func collxprint(Colls []*COLL) {
	if len(Colls) > 0 {
		fmt.Printf("%s N=%d\n", COLLECTOR_TYPE, len(Colls))

		for i, Coll := range Colls {
			fmt.Printf("[%d] %-10s Do=%6.3f  D1=%6.3f Tin=%5.2f Tout=%5.2f Q=%4.0f Sol=%4.0f Te=%5.1f\n",
				i, Coll.Name, Coll.Do, Coll.D1, Coll.Tin,
				Coll.Cmp.Elouts[0].Sysv, Coll.Q, Coll.Sol, Coll.Te)
			fmt.Printf("   exs=%s  b0=%5.3f  b1=%5.3f ec=%5.3f\n", Coll.sol.Name,
				Coll.Cat.b0, Coll.Cat.b1, Coll.ec)
		}
	}
}

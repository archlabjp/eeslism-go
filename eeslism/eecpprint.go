package eeslism

import (
	"errors"
	"fmt"
)

func cmpprint(id, N int, cmp []COMPNT, Elout []*ELOUT, Elin []*ELIN) {
	name, eqptype, envname := "name", "eqptype", "envname"
	if id == 1 {
		fmt.Printf("COMPNT\n n %-10s %-10s %-10s -c- nca neq Nout Nin nivr Elou Elin\n", name, eqptype, envname)
	}
	for i := 0; i < N; i++ {
		c := &cmp[i]
		var cEloutsIdx, cElinIdx int
		for cEloutsIdx = 0; cEloutsIdx < len(Elout); cEloutsIdx++ {
			if &c.Elouts[0] == &Elout[0] {
				break
			}
		}
		for cElinIdx = 0; cElinIdx < len(Elin); cElinIdx++ {
			if &c.Elins[0] == &Elin[0] {
				break
			}
		}
		fmt.Printf("%2d %-10s %-10s %-10s   %c %4d %3d %4d %3d %d %4d %4d\n", i, c.Name, c.Eqptype, c.Envname, c.Control,
			c.Ncat, c.Neqp, c.Nout, c.Nin, c.Nivar, cEloutsIdx, cElinIdx)
	}
}

func eloutprint(id int, E []*ELOUT, cmp []*COMPNT) {
	if id == 1 {
		fmt.Printf("ELOUT\n  n name            id fld contl sysld Cmp   G      cfo    cfin\n")
	}
	for i := range E {
		e := E[i]
		var eCmpIdx int
		for eCmpIdx = 0; eCmpIdx < len(cmp); eCmpIdx++ {
			if e.Cmp == cmp[eCmpIdx] {
				break
			}
		}
		fmt.Printf("%3d (%-10s)     %c   %c   %c    %c  %4d [%5.3f]  %6.3f",
			i, e.Cmp.Name, e.Id, e.Fluid, e.Control, e.Sysld, eCmpIdx, e.G, e.Coeffo)

		for j := 0; j < e.Ni; j++ {
			fmt.Printf(" %6.3f", e.Coeffin[j])
		}

		fmt.Printf(" Co=%6.4f\n", e.Co)
	}
}

func eloutfprint(id int, E []*ELOUT, cmp []*COMPNT) {
	if id == 1 {
		fmt.Fprintf(Ferr, "ELOUT\n  n         id fld contl sysld Cmp   G      cfo    cfin\n")
	}
	for i := range E {
		e := E[i]
		cmp_idx := 0
		for cmp_idx = 0; cmp_idx < len(cmp); cmp_idx++ {
			if e.Cmp == cmp[cmp_idx] {
				break
			}
		}
		fmt.Fprintf(Ferr, "%3d (%-6s) %c   %c   %c    %c  %4d [%5.3f]  %6.3f",
			i, e.Cmp.Name, e.Id, e.Fluid, e.Control, e.Sysld, cmp_idx, e.G, e.Coeffo)

		for j := 0; j < e.Ni; j++ {
			fmt.Fprintf(Ferr, " %6.3f", e.Coeffin[j])
		}

		fmt.Fprintf(Ferr, " Co=%6.4f\n", e.Co)
	}
}

func elinprint(id int, C []*COMPNT, eo []*ELOUT, ei []*ELIN) {

	if id == 1 {
		fmt.Printf("ELIN\n  n  id   upo  upv\n")
	}

	for i := range C {
		for ii := 0; ii < C[i].Nout; ii++ {
			Eo := C[i].Elouts[ii]

			for j := 0; j < Eo.Ni; j++ {
				var o, v int

				E := Eo.Elins[j]

				if E.Upo != nil && eo != nil {
					Upo_idx := 0
					for Upo_idx = 0; Upo_idx < len(eo); Upo_idx++ {
						if E.Upo == eo[Upo_idx] {
							break
						}
					}
					o = Upo_idx
				} else {
					o = -999
				}

				if E.Upv != nil && eo != nil {
					Upv_idx := 0
					for Upv_idx = 0; Upv_idx < len(eo); Upv_idx++ {
						if E.Upv == eo[Upv_idx] {
							break
						}
					}
					v = Upv_idx
				} else {
					v = -999
				}

				var l int
				for l = 0; l < len(ei); l++ {
					if E == ei[l] {
						break
					}
				}
				fmt.Printf("%3d (%-6s) %c   %3d   %3d",
					l, C[i].Name, E.Id, o, v)
				if E.Upo != nil {
					fmt.Printf(" upo=(%-6s)", E.Upo.Cmp.Name)
				}
				if E.Upv != nil {
					fmt.Printf(" upv=(%-6s)", E.Upv.Cmp.Name)
				}
				fmt.Printf("\n")
			}
		}
	}
}

func elinfprint(id int, C []*COMPNT, eo []*ELOUT, ei []*ELIN) {
	var E *ELIN
	var Eo *ELOUT
	var o, v int

	if id == 1 {
		fmt.Fprintf(Ferr, "ELIN\n  n  id   upo  upv\n")
	}

	for i := range C {
		Ci := C[i]

		for ii := 0; ii < Ci.Nout; ii++ {
			Eo = Ci.Elouts[ii]

			for j := 0; j < Eo.Ni; j++ {
				E = Eo.Elins[j]

				if E.Upo != nil && eo != nil {
					for o = 0; o < len(eo); o++ {
						if E.Upo == eo[o] {
							break
						}
					}
				} else {
					o = -999
				}

				if E.Upv != nil && eo != nil {
					for v = 0; v < len(eo); v++ {
						if E.Upv == eo[v] {
							break
						}
					}
				} else {
					v = -999
				}

				var l int = 0
				for l = 0; l < len(ei); l++ {
					if E == ei[l] {
						break
					}
				}
				fmt.Fprintf(Ferr, "%3d (%-6s) %c   %3d   %3d",
					l, Ci.Name, E.Id, o, v)
				if E.Upo != nil {
					fmt.Fprintf(Ferr, " upo=(%-6s)", E.Upo.Cmp.Name)
				}
				if E.Upv != nil {
					fmt.Fprintf(Ferr, " upv=(%-6s)", E.Upv.Cmp.Name)
				}
				fmt.Fprintf(Ferr, "\n")
			}
		}
	}
}

func plistprint(Mpath []*MPATH, Pe []*PELM, Eo []*ELOUT, Ei []*ELIN) {
	fmt.Printf("xxx plistprint\n")
	idx_0, _ := FindPELM(Pe, Mpath[0].Plist[0].Pelm[0])
	for i, Mpathi := range Mpath {
		fmt.Printf("\nMpath=[%d] %s sys=%c type=%c fluid=%c Nlpath= %d  Ncv=%d lvcmx=%d\n",
			i, Mpathi.Name, Mpathi.Sys, Mpathi.Type, Mpathi.Fluid, len(Mpathi.Plist),
			Mpathi.Ncv, Mpathi.Lvcmx)

		for j, pl := range Mpathi.Plist {
			// 要素のグローバルインデックスの取得
			idx, _ := FindPELM(Pe[idx_0:], pl.Pelm[0])

			fmt.Printf("PLIST\n  n type Nelm Npump Nvav lvc Pelm  G \n")
			fmt.Printf("%3d  %c  %3d  %3d %3d %3d %6.3f\n",
				j, pl.Type, len(pl.Pelm), pl.Nvav, pl.Lvc, idx, pl.G)

			fmt.Printf("    PELM  n  co ci elin eout\n")
			for _, p := range pl.Pelm {
				pIdx, _ := FindPELM(pl.Pelm, p)
				pInIdx, _ := FindELIN(Ei, p.In)
				pOutIdx, _ := FindELOUT(Eo, p.Out)
				fmt.Printf("        %3d   %c  %c %4d %4d  %s\n",
					pIdx, p.Co, p.Ci, pInIdx, pOutIdx, p.Cmp.Name)
			}
		}
	}
}

func FindPELM(pl []*PELM, p *PELM) (int, error) {
	for i := 0; i < len(pl); i++ {
		if p == pl[i] {
			return i, nil
		}
	}
	return -1, errors.New("Not Found")
}

func FindELIN(elins []*ELIN, elin *ELIN) (int, error) {
	for i := 0; i < len(elins); i++ {
		if elin == elins[i] {
			return i, nil
		}
	}
	return -1, errors.New("Not Found")
}

func FindELOUT(elouts []*ELOUT, elout *ELOUT) (int, error) {
	for i := 0; i < len(elouts); i++ {
		if elout == elouts[i] {
			return i, nil
		}
	}
	return -1, errors.New("Not Found")
}

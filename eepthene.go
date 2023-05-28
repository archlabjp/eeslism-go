package main

/* 経路に沿ったシステム要素の熱量計算 */

func Pathheat(Nmpath int, Mpath []*MPATH) {
	var c, cG float64
	for i := 0; i < Nmpath; i++ {
		c = Spcheat(Mpath[i].Fluid)
		for j := 0; j < Mpath[i].Nlpath; j++ {
			Pli := &Mpath[i].Plist[j]
			cG = c * Pli.G
			for k := 0; k < Pli.Nelm; k++ {
				Pelm := Pli.Pelm[k]
				if Pelm.Cmp.Eqptype == DIVERG_TYPE || Pelm.Cmp.Eqptype == CONVRG_TYPE ||
					Pelm.Cmp.Eqptype == DIVGAIR_TYPE || Pelm.Cmp.Eqptype == CVRGAIR_TYPE {
					Pelm.Out.Q = 0.0
				} else if Pelm.Out.Control == OFF_SW {
					Pelm.Out.Q = 0.0
				} else {
					Pelm.Out.Q = cG * (Pelm.Out.Sysv - Pelm.In.Sysvin)
				}
			}
		}
	}
}

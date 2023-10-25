package eeslism

/* 経路に沿ったシステム要素の熱量計算 */

func Pathheat(Mpath []*MPATH) {
	for _, mpath := range Mpath {
		c := Spcheat(mpath.Fluid)
		for _, Pli := range mpath.Plist {
			cG := c * Pli.G
			for _, Pelm := range Pli.Pelm {
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

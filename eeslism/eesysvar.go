package eeslism

func Sysvar(Compnt []*COMPNT) {
	for m := range Compnt {
		for i := 0; i < Compnt[m].Nin; i++ {
			I := Compnt[m].Elins[i]
			I.Sysvin = 0.0
			if I.Upv != nil {
				I.Sysvin = I.Upv.Sysv
			}
		}
	}
}

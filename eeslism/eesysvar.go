package eeslism

func Sysvar(Ncompnt int, Compnt []COMPNT) {
	for m := 0; m < Ncompnt; m++ {
		for i := 0; i < Compnt[m].Nin; i++ {
			I := Compnt[m].Elins[i]
			I.Sysvin = 0.0
			if I.Upv != nil {
				I.Sysvin = I.Upv.Sysv
			}
		}
	}
}

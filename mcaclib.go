package main

func wcoil(Air_SW ControlSWType, Water_SW ControlSWType, wet rune, Gaet float64, Gaeh float64, xai float64, Twi float64, Et *ACS, Ex *ACS, Ew *ACS) {
	var ca, cv, ro float64 // ca, cv, ro は外部で定義される想定
	if wet == 'd' || Water_SW == OFF_SW || Air_SW == OFF_SW {
		// 片側系統が停止していたときに対応するように修正
		// Satoh Debug 2009/1/9
		if Water_SW != OFF_SW {
			Et.W = ca * Gaet
			Et.T = ca * Gaet
			Et.X = 0.0
			Et.C = 0.0
		} else {
			Et.W = 0.0
			Et.T = 0.0
			Et.X = 0.0
			Et.C = 0.0
		}

		Ex.W = 0.0
		Ex.T = 0.0
		Ex.X = 0.0
		Ex.C = 0.0

		if Air_SW != OFF_SW {
			Ew.W = ca * Gaet
			Ew.T = ca * Gaet
			Ew.X = 0.0
			Ew.C = 0.0
		} else {
			Ew.W = 0.0
			Ew.T = 0.0
			Ew.X = 0.0
			Ew.C = 0.0
		}
	} else {
		var aw, bw, cs float64
		hstaircf(Twi, Twi+5.0, &aw, &bw)
		cs = ca + cv*xai

		Et.W = ca * Gaet
		Et.T = ca * Gaet
		Et.X = 0.0
		Et.C = 0.0

		Ex.W = (Gaeh*bw - Gaet*ca) / ro
		Ex.T = (Gaeh*cs - Gaet*ca) / ro
		Ex.X = Gaeh
		Ex.C = -Gaeh * aw / ro

		Ew.W = Gaeh * bw
		Ew.T = Gaeh * cs
		Ew.X = Gaeh * ro
		Ew.C = -Gaeh * aw
	}
}

func Qcoils(Et ACS, Tai float64, xai float64, Twi float64) float64 {
	return Et.W*Twi - Et.T*Tai - Et.X*xai - Et.C
}

func Qcoill(Ex ACS, Tai float64, xai float64, Twi float64) float64 {
	var ro float64 // ro は外部で定義される想定
	return ro * (Ex.W*Twi - Ex.T*Tai - Ex.X*xai - Ex.C)
}

func hstaircf(Tw1 float64, Tw2 float64, a *float64, b *float64) {
	var h1, h2 float64
	h1 = FNH(Tw1, FNXtr(Tw1, 100.0))
	h2 = FNH(Tw2, FNXtr(Tw2, 100.0))
	*b = (h2 - h1) / (Tw2 - Tw1)
	*a = h1 - *b*Tw1
	// fmt.Printf("== hstaircf Tw1,Tw2=%f %f a=%f b=%f\n",Tw1,Tw2,*a,*b)
}

package eeslism

// 冷温水コイルの処理熱量計算用係数
func wcoil(Air_SW ControlSWType, Water_SW ControlSWType, wet rune, Gaet float64, Gaeh float64,
	xai float64, Twi float64) (ACS, ACS, ACS) {

	var Et, Ex, Ew ACS

	if wet == 'd' || Water_SW == OFF_SW || Air_SW == OFF_SW {
		// 片側系統が停止していたときに対応するように修正
		// Satoh Debug 2009/1/9
		if Water_SW != OFF_SW {
			Et = ACS{
				W: Ca * Gaet,
				T: Ca * Gaet,
				X: 0.0,
				C: 0.0,
			}
		} else {
			Et = ACS{
				W: 0.0,
				T: 0.0,
				X: 0.0,
				C: 0.0,
			}
		}

		Ex = ACS{
			W: 0.0,
			T: 0.0,
			X: 0.0,
			C: 0.0,
		}

		if Air_SW != OFF_SW {
			Ew = ACS{
				W: Ca * Gaet,
				T: Ca * Gaet,
				X: 0.0,
				C: 0.0,
			}
		} else {
			Ew = ACS{
				W: 0.0,
				T: 0.0,
				X: 0.0,
				C: 0.0,
			}
		}
	} else {
		aw, bw := hstaircf(Twi, Twi+5.0)
		cs := Ca + Cv*xai

		Et = ACS{
			W: Ca * Gaet,
			T: Ca * Gaet,
			X: 0.0,
			C: 0.0,
		}

		Ex = ACS{
			W: (Gaeh*bw - Gaet*Ca) / Ro,
			T: (Gaeh*cs - Gaet*Ca) / Ro,
			X: Gaeh,
			C: -Gaeh * aw / Ro,
		}

		Ew = ACS{
			W: Gaeh * bw,
			T: Gaeh * cs,
			X: Gaeh * Ro,
			C: -Gaeh * aw,
		}
	}

	return Et, Ex, Ew
}

func Qcoils(Et ACS, Tai float64, xai float64, Twi float64) float64 {
	return Et.W*Twi - Et.T*Tai - Et.X*xai - Et.C
}

func Qcoill(Ex ACS, Tai float64, xai float64, Twi float64) float64 {
	return Ro * (Ex.W*Twi - Ex.T*Tai - Ex.X*xai - Ex.C)
}

func hstaircf(Tw1 float64, Tw2 float64) (float64, float64) {
	h1 := FNH(Tw1, FNXtr(Tw1, 100.0))
	h2 := FNH(Tw2, FNXtr(Tw2, 100.0))
	b := (h2 - h1) / (Tw2 - Tw1)
	a := h1 - b*Tw1
	// fmt.Printf("== hstaircf Tw1,Tw2=%f %f a=%f b=%f\n",Tw1,Tw2,*a,*b)
	return a, b
}

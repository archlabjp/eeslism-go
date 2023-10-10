package eeslism

const MAXINT_DAY = -999.0
const MININT_DAY = 999.0

func svdyint(vd *SVDAY) {
	vd.M = 0.0
	vd.Mn = MININT_DAY
	vd.Mx = MAXINT_DAY
	vd.Hrs = 0
	vd.Mntime = -1
	vd.Mxtime = -1
}

func svdaysum(time int64, control ControlSWType, v float64, vd *SVDAY) {
	if control != '0' {
		vd.M += v
		vd.Hrs++
		minmark(&vd.Mn, &vd.Mntime, v, time)
		maxmark(&vd.Mx, &vd.Mxtime, v, time)
	}
	if time == 2400 && vd.Hrs > 0 {
		vd.M /= float64(vd.Hrs)
	}
}

func svmonsum(Mon int, Day int, time int, control ControlSWType, v float64, vd *SVDAY, Dayend int, SimDayend int) {
	MoNdTt := int64(1000000*Mon + 10000*Day + time)

	if control != '0' {
		vd.M += v
		vd.Hrs++
		minmark(&vd.Mn, &vd.Mntime, v, MoNdTt)
		maxmark(&vd.Mx, &vd.Mxtime, v, MoNdTt)
	}

	if IsEndDay(Mon, Day, Dayend, SimDayend) && vd.Hrs > 0 && time == 2400 {
		vd.M /= float64(vd.Hrs)
	}
}

/* ------------------------------------ */

func qdyint(Qd *QDAY) {
	Qd.H = 0.0
	Qd.C = 0.0
	Qd.Hmx = 0.0
	Qd.Cmx = 0.0
	Qd.Hhr = 0
	Qd.Chr = 0
	Qd.Hmxtime = -1
	Qd.Cmxtime = -1
}

func qdaysum(time int64, control ControlSWType, Q float64, Qd *QDAY) {
	if control != '0' {
		if Q > 0.0 {
			Qd.H += Q
			maxmark(&Qd.Hmx, &Qd.Hmxtime, Q, time)
			Qd.Hhr++
		} else if Q < 0.0 {
			Qd.C += Q
			minmark(&Qd.Cmx, &Qd.Cmxtime, Q, time)
			Qd.Chr++
		}
	}

	if time == 2400 {
		Qd.H *= Cff_kWh
		Qd.C *= Cff_kWh
	}
}

func qmonsum(Mon int, Day int, time int, control ControlSWType, Q float64, Qd *QDAY, Dayend int, SimDayend int) {
	MoNdTt := int64(1000000*Mon + 10000*Day + time)

	if control != '0' {
		if Q > 0.0 {
			Qd.H += Q
			maxmark(&Qd.Hmx, &Qd.Hmxtime, Q, MoNdTt)
			Qd.Hhr++
		} else if Q < 0.0 {
			Qd.C += Q
			minmark(&Qd.Cmx, &Qd.Cmxtime, Q, MoNdTt)
			Qd.Chr++
		}
	}

	if IsEndDay(Mon, Day, Dayend, SimDayend) && time == 2400 {
		Qd.H *= Cff_kWh
		Qd.C *= Cff_kWh
	}
}

// 日集計関数。非運転時にも集計を行う
func qdaysumNotOpe(time int64, Q float64, Qd *QDAY) {
	if Q > 0.0 {
		Qd.H += Q
		maxmark(&Qd.Hmx, &Qd.Hmxtime, Q, time)
		Qd.Hhr++
	} else if Q < 0.0 {
		Qd.C += Q
		minmark(&Qd.Cmx, &Qd.Cmxtime, Q, time)
		Qd.Chr++
	}

	if time == 2400 {
		Qd.H *= Cff_kWh
		Qd.C *= Cff_kWh
	}
}

func qmonsumNotOpe(Mon int, Day int, time int, Q float64, Qd *QDAY, Dayend int, SimDayend int) {
	MoNdTt := int64(1000000*Mon + 10000*Day + time)

	if Q > 0.0 {
		Qd.H += Q
		maxmark(&Qd.Hmx, &Qd.Hmxtime, Q, MoNdTt)
		Qd.Hhr++
	} else if Q < 0.0 {
		Qd.C += Q
		minmark(&Qd.Cmx, &Qd.Cmxtime, Q, MoNdTt)
		Qd.Chr++
	}

	if IsEndDay(Mon, Day, Dayend, SimDayend) && time == 2400 {
		Qd.H *= Cff_kWh
		Qd.C *= Cff_kWh
	}
}

/* ------------------------------------ */

func edyint(Ed *EDAY) {
	Ed.D = 0.0
	Ed.Mx = 0.0
	Ed.Hrs = 0
	Ed.Mxtime = -1
}

func edaysum(time int, control ControlSWType, E float64, Ed *EDAY) {
	if control != '0' {
		Ed.D += E
		maxmark(&Ed.Mx, &Ed.Mxtime, E, int64(time))
		Ed.Hrs++
	}

	if time == 2400 {
		Ed.D *= Cff_kWh
	}
}

func emonsum(Mon, Day, time int, control ControlSWType, E float64, Ed *EDAY, Dayend, SimDayend int) {
	var MoNdTt int64 = int64(1000000*Mon + 10000*Day + time)

	if control != OFF_SW {
		Ed.D += E
		maxmark(&Ed.Mx, &Ed.Mxtime, E, MoNdTt)
		Ed.Hrs++
	}

	if IsEndDay(Mon, Day, Dayend, SimDayend) && time == 2400 {
		Ed.D *= Cff_kWh
	}
}

func emtsum(Mon, Day, time int, control ControlSWType, E float64, Ed *EDAY) {
	if control != OFF_SW {
		Ed.D += E
	}
}

func minmark(minval *float64, timemin *int64, v float64, time int64) {
	if v <= *minval {
		*timemin = time
		*minval = v
	}
}

func maxmark(maxval *float64, timemax *int64, v float64, time int64) {
	if v >= *maxval {
		*timemax = time
		*maxval = v
	}
}

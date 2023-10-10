package eeslism

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"
)

/* VCFILE からの気象データ入力 */

func wdflinit(Simc *SIMCONTL, Estl *ESTL, Tlist []TLIST) {
	var wp WDPT
	var s, ss, Err string
	var dt float64
	var id, N, i, m int

	if s = Estl.Wdloc; s == "" {
		return
	}

	Err = fmt.Sprintf(ERRFMT, "(wdflinit)")
	Locinit(Simc.Loc)
	loc := Simc.Loc

	m = -1

	for _, field := range strings.Fields(s) {
		N = len(field)
		if st := strings.Index(field, "="); st >= 0 {
			ss = field[:st]
			dt, _ = strconv.ParseFloat(field[st+1:], 64)
			switch ss {
			case "Lat":
				loc.Lat = dt
			case "Lon":
				loc.Lon = dt
			case "Ls":
				loc.Ls = dt
			case "Tgrav":
				loc.Tgrav = dt
			case "DTgr":
				loc.DTgr = dt
			case "daymx":
				loc.Daymxert = int(dt)
			default:
				id = 1
			}
		} else {
			if field == "-" || m >= 0 {
				switch field {
				case "-Twsup":
					m = 0
				default:
					loc.Twsup[m], _ = strconv.ParseFloat(field, 64)
					m++
					if m > 11 {
						m = -1
					}
				}
			} else {
				loc.Name = field
			}
		}
		s = s[N:]
		for len(s) > 0 && unicode.IsSpace(rune(s[0])) {
			s = s[1:]
		}
	}

	if id != 0 {
		fmt.Printf("%s %s\n", Err, ss)
	}

	wp.Ta = nil
	wp.Xa = nil
	wp.Rn = nil
	wp.Ihor = nil
	wp.Rh = nil
	wp.Cc = nil
	wp.Wv = nil
	wp.Wdre = nil

	for i = 0; i < Estl.Ndata; i++ {
		t := &Tlist[i]
		if t.Name == "Wd" {
			s = t.Id
			val := t.Fval
			switch s {
			case "T": // 温度
				wp.Ta = val
			case "x": // 絶対湿度
				wp.Xa = val
			case "Idn": // 法線面直達日射量
				wp.Idn = val
			case "Isky": // 水平面天空日射量
				wp.Isky = val
			case "Ihor": // 水平面全天日射量
				wp.Ihor = val
			case "CC": // 雲量
				wp.Cc = val
			case "Wdre": // 風向
				wp.Wdre = val
			case "Wv": // 風速
				wp.Wv = val
			case "RH": // 相対湿度
				wp.Rh = val
			case "RN": // 夜間放射量
				wp.Rn = val
			}
		}
	}

	Simc.Wdpt = wp
}

func Wdflinput(wp *WDPT, Wd *WDAT) {
	var Br float64

	Wd.T = wp.Ta[0]
	Wd.Idn = wp.Idn[0]
	Wd.Isky = wp.Isky[0]

	if wp.Ihor == nil {
		Wd.Ihor = Wd.Idn*Wd.Sh + Wd.Isky
	}

	if wp.Xa != nil {
		Wd.X = wp.Xa[0]
	} else {
		Wd.X = -999.0
	}

	if wp.Rh != nil {
		Wd.RH = wp.Rh[0]
	} else {
		Wd.RH = -999.0
	}

	if wp.Cc != nil {
		Wd.CC = wp.Cc[0]
	} else {
		Wd.CC = -999.0
	}

	if wp.Rn != nil {
		Wd.RN = wp.Rn[0]
	} else {
		Wd.RN = -999.0
	}

	if wp.Wv != nil {
		Wd.Wv = wp.Wv[0]
	} else {
		Wd.Wv = -999.0
	}

	if wp.Wdre != nil {
		Wd.Wdre = wp.Wdre[0]
	} else {
		Wd.Wdre = 0.0
	}

	if Wd.X > 0.0 && Wd.RH < 0.0 {
		Wd.RH = FNRhtx(Wd.T, Wd.X)
	} else if Wd.X < 0.0 && Wd.RH > 0.0 {
		Wd.X = FNXtr(Wd.T, Wd.RH)
	}

	if Wd.X > 0.0 {
		Wd.H = FNH(Wd.T, Wd.X)
	}

	if Wd.X > 0.0 && Wd.CC > 0.0 || Wd.RN < 0.0 {
		Br = 0.51 + 0.209*math.Sqrt(FNPwx(Wd.X))
		Wd.RN = (1.0 - 0.62*Wd.CC/10.0) * (1.0 - Br) * Sgm * math.Pow(Wd.T+273.15, 4.0)
		Wd.Rsky = ((1.0-0.62*Wd.CC/10.0)*Br + 0.62*Wd.CC/10.0) * Sgm * math.Pow(Wd.T+273.15, 4.0)
	} else {
		Wd.Rsky = Sgm*math.Pow(Wd.T+273.15, 4.0) - Wd.RN
	}
}

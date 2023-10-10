package eeslism

import "os"

const (
	ALO = 23.0

	// Uncomment these lines if you want to use these constants
	// MAXBDP  = 100
	// MAXOBS  = 100
	// MAXTREE = 10 // Maximum number of trees
	// MAXPOLY = 50

	UNIT = "SI"
	PI   = 3.141592654
)

var (
	Sgm = 5.67e-8
	Ca  = 1005.0
	Cv  = 1846.0
	Roa = 1.29
	Cw  = 4186.0

	Row = 1000.0
	Ro  = 2501000.0

	G           = 9.8
	DTM         = 0.0 // Assign the value of dTM here
	Cff_kWh     = 0.0 // Assign the value of cff_kWh here
	VAVCountMAX = 0   // Assign the value of VAV_Count_MAX here

	Fbmlist = "" // Assign the value of Fbmlist here

	DEBUG   = true
	dayprn  = false
	DAYweek = [8]string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun", "Hol"}
	Ferr    = os.Stderr // Assuming you want to write errors to standard error
	//DISPLAY_DELAY = 0 // Assign the value of DISPLAY_DELAY here
	SETprint = 0
)

// 月の末日かどうかをチェックする
func IsEndDay(Mon, Day, Dayend, SimDayend int) bool {
	Nde := []int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	nday := SimDayend
	if nday > 365 {
		nday -= 365
	}
	if Day == Nde[Mon-1] || Dayend == nday {
		return true
	}

	return false
}

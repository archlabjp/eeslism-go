package main

func dmin(a, b float64) float64 {
	if a <= b {
		return a
	} else {
		return b
	}
}

func dmax(a, b float64) float64 {
	if a >= b {
		return a
	} else {
		return b
	}
}

func imax(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

func imin(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

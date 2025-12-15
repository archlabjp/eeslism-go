//go:build cmath
// +build cmath

package eeslism

/*
#cgo LDFLAGS: -lm
#include <math.h>
*/
import "C"

// C言語版の数学関数をラップ
// ビルド時に -tags cmath を指定すると有効になる

func mathPow(x, y float64) float64 {
	return float64(C.pow(C.double(x), C.double(y)))
}

func mathExp(x float64) float64 {
	return float64(C.exp(C.double(x)))
}

func mathLog(x float64) float64 {
	return float64(C.log(C.double(x)))
}

func mathSqrt(x float64) float64 {
	return float64(C.sqrt(C.double(x)))
}

func mathAbs(x float64) float64 {
	return float64(C.fabs(C.double(x)))
}

func mathSin(x float64) float64 {
	return float64(C.sin(C.double(x)))
}

func mathCos(x float64) float64 {
	return float64(C.cos(C.double(x)))
}

func mathTan(x float64) float64 {
	return float64(C.tan(C.double(x)))
}

func mathAtan(x float64) float64 {
	return float64(C.atan(C.double(x)))
}

func mathAtan2(y, x float64) float64 {
	return float64(C.atan2(C.double(y), C.double(x)))
}

func mathAsin(x float64) float64 {
	return float64(C.asin(C.double(x)))
}

func mathAcos(x float64) float64 {
	return float64(C.acos(C.double(x)))
}

func mathFloor(x float64) float64 {
	return float64(C.floor(C.double(x)))
}

func mathCeil(x float64) float64 {
	return float64(C.ceil(C.double(x)))
}

func mathMod(x, y float64) float64 {
	return float64(C.fmod(C.double(x), C.double(y)))
}

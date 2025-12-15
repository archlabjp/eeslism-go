//go:build !cmath
// +build !cmath

package eeslism

import "math"

// Go標準ライブラリの数学関数をラップ
// デフォルトで使用される（-tags cmath を指定しない場合）

func mathPow(x, y float64) float64 {
	return math.Pow(x, y)
}

func mathExp(x float64) float64 {
	return math.Exp(x)
}

func mathLog(x float64) float64 {
	return math.Log(x)
}

func mathSqrt(x float64) float64 {
	return math.Sqrt(x)
}

func mathAbs(x float64) float64 {
	return math.Abs(x)
}

func mathSin(x float64) float64 {
	return math.Sin(x)
}

func mathCos(x float64) float64 {
	return math.Cos(x)
}

func mathTan(x float64) float64 {
	return math.Tan(x)
}

func mathAtan(x float64) float64 {
	return math.Atan(x)
}

func mathAtan2(y, x float64) float64 {
	return math.Atan2(y, x)
}

func mathAsin(x float64) float64 {
	return math.Asin(x)
}

func mathAcos(x float64) float64 {
	return math.Acos(x)
}

func mathFloor(x float64) float64 {
	return math.Floor(x)
}

func mathCeil(x float64) float64 {
	return math.Ceil(x)
}

func mathMod(x, y float64) float64 {
	return math.Mod(x, y)
}

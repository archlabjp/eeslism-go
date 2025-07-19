package eeslism

import (
	"testing"

	"gotest.tools/assert"
)

func Test_FNarow(t *testing.T) {
	assert.Equal(t, 1.293, FNarow(0))
	assert.Equal(t, 1.293/2.0, FNarow(273.15))
}

func Test_FNac(t *testing.T) {
	assert.Equal(t, 1005.0, FNac())
}

func Test_FNalam(t *testing.T) {
	assert.Equal(t, 0.0241, FNalam(0))
	assert.Equal(t, FNAN, FNalam(100))
}

// func Test_FNamew(t *testing.T) {
// 	assert.Equal(t, 0.0074237, FNamew(0))
// 	assert.Equal(t, 0.0074237/2.0, FNamew(273.15))
// }

func Test_FNanew(t *testing.T) {
	assert.Equal(t, FNamew(0)/FNarow(0), FNanew(0), 0.00001)
	assert.Equal(t, FNamew(273.15)/FNarow(273.15), FNanew(273.15), 0.00001)
}

// func Test_FNabeta(t *testing.T) {
// 	assert.Equal(t, 1.0/273.15, FNabeta(0), 0.00001)
// 	assert.Equal(t, 1.0/273.15/2.0, FNabeta(273.15), 0.00001)
// }

// func Test_FNwrow(t *testing.T) {
// 	assert.Equal(t, NAN, FNwrow(0))
// 	assert.Equal(t, 1000.5-0.068737*50.0-0.0035781*50.0*50.0, FNwrow(50), 0.001)
// 	assert.Equal(t, 1008.7-0.28735*100.0-0.0021643*100.0*100.0, FNwrow(100), 0.001)
// 	assert.Equal(t, 1008.7-0.28735*150.0-0.0021643*150.0*150.0, FNwrow(150), 0.001)
// 	assert.Equal(t, NAN, FNwrow(200))
// }

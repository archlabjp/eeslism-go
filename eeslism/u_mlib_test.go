package eeslism

import (
	"testing"

	"gotest.tools/assert"
)

func Test_Spcheat(t *testing.T) {
	assert.Equal(t, Cw, Spcheat(WATER_FLD))
	assert.Equal(t, Ca, Spcheat(AIRa_FLD))
	assert.Equal(t, -9999.0, Spcheat('x'))
}

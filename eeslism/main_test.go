package eeslism

import (
	"os"
	"testing"
)

func Test_wind16(t *testing.T) {
	os.Chdir("../Base")
	Entry("標準プラン.txt")
	// spd, dir := Wind16(1.0, 1.0)
	// assert.InDelta(t, 1.4141456, spd, 0.0001)
	// assert.Equal(t, 180.0+45.0, dir)
}

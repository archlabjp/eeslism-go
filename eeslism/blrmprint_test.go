
package eeslism

import (
	"testing"
)

func TestQrmsum(t *testing.T) {
	// Setup
	DTM = 3600 // 1 hour
	__Qrmsum_oldday = 0
	Day1 := 1
	Day2 := 2
	_Room := []*ROOM{
		{Tr: 25.0},
	}
	Qrm := []*QRM{
		{Tsol: 100, Asol: 50, Arn: 20, Hums: 10, Light: 30, Apls: 40, Huml: 5, Apll: 15, Qeqp: 5, Qinfs: 10, Qinfl: 2, Qsto: 5, Qstol: 1, AE: 10, AG: 2},
	}
	Trdav := []float64{0.0}
	Qrmd := []*QRM{
		{},
	}

	// Execute for Day 1
	Qrmsum(Day1, _Room, Qrm, Trdav, Qrmd)

	// Verify for Day 1
	if Qrmd[0].Tsol != 100 {
		t.Errorf("Qrmd[0].Tsol should be 100, but got %f", Qrmd[0].Tsol)
	}

	// Execute for Day 1 again
	Qrmsum(Day1, _Room, Qrm, Trdav, Qrmd)

	// Verify for Day 1 again
	if Qrmd[0].Tsol != 200 {
		t.Errorf("Qrmd[0].Tsol should be 200, but got %f", Qrmd[0].Tsol)
	}

	// Execute for Day 2
	Qrmsum(Day2, _Room, Qrm, Trdav, Qrmd)

	// Verify for Day 2
	if Qrmd[0].Tsol != 100 {
		t.Errorf("Qrmd[0].Tsol should be 100, but got %f", Qrmd[0].Tsol)
	}
}

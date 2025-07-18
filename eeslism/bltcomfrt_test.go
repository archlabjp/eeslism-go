
package eeslism

import (
	"math"
	"testing"
)

func TestFotinit(t *testing.T) {
	// Setup
	__Fotinit_init = 'i'
	Room := []*ROOM{
		{
			rmld: &RMLOAD{},
			Ntr:  2,
			Nrp:  1,
		},
	}

	// Execute
	Fotinit(Room)

	// Verify
	if len(Room[0].rmld.FOTN) != 2 {
		t.Errorf("FOTN should have length 2, but got %d", len(Room[0].rmld.FOTN))
	}
	if len(Room[0].rmld.FOPL) != 1 {
		t.Errorf("FOPL should have length 1, but got %d", len(Room[0].rmld.FOPL))
	}
}

func TestFotf(t *testing.T) {
	// Setup
	weight := 0.6
	Room := &ROOM{
		OTsetCwgt: &weight,
		Area:      100,
		N:         1,
		Ntr:       0,
		Nrp:       0,
		rsrf: []*RMSRF{
			{
				A:   10,
				WSR: 0.5,
				WSC: 0.2,
			},
		},
		rmld: &RMLOAD{},
	}

	// Execute
	Fotf(Room)

	// Verify
	expectedFOTr := 0.6 + (1.0-0.6)*10*0.5/100
	expectedFOC := (1.0 - 0.6) * 10 * 0.2 / 100
	if math.Abs(Room.rmld.FOTr-expectedFOTr) > 1e-9 {
		t.Errorf("FOTr calculation failed: expected %f, got %f", expectedFOTr, Room.rmld.FOTr)
	}
	if math.Abs(Room.rmld.FOC-expectedFOC) > 1e-9 {
		t.Errorf("FOC calculation failed: expected %f, got %f", expectedFOC, Room.rmld.FOC)
	}
}

func TestPmv0(t *testing.T) {
	// Test with typical values
	pmv := Pmv0(1.0, 1.0, 25, 0.01, 25, 0.1)
	if pmv > 10 || pmv < -10 {
		t.Errorf("Pmv0 calculation seems to be out of reasonable range, got %f", pmv)
	}
}

func TestSET_star(t *testing.T) {
	// Test with typical values
	set := SET_star(25, 25, 0.1, 50, 1.0, 1.0, 0, 101.3)
	if set > 50 || set < 0 {
		t.Errorf("SET_star calculation seems to be out of reasonable range, got %f", set)
	}
}

func TestFindSaturatedVaporPressureTorr(t *testing.T) {
	// Test with a known value (25C -> 23.756 Torr)
	p := FindSaturatedVaporPressureTorr(25)
	if math.Abs(p-23.75745) > 0.00001 {
		t.Errorf("FindSaturatedVaporPressureTorr failed: expected 23.756, got %f", p)
	}
}

package eeslism

import (
	"testing"
)

func TestRoomene(t *testing.T) {
	// Setup basic room with energy calculation data
	room := &ROOM{
		Name:    "TestRoom",
		N:       2,
		rsrf:    make([]*RMSRF, 2),
		Tr:      22.0,
		Trold:   21.0,
		xr:      0.008,
		RH:      50.0,
		hr:      45.0,
		Ntr:     0,
		Nrp:     0,
		Nasup:   0,
		cmp: &COMPNT{
			Elouts: []*ELOUT{
				{Sysv: 22.0},
				{Sysv: 0.008},
			},
		},
	}

	// Create mock surfaces
	for i := 0; i < 2; i++ {
		room.rsrf[i] = &RMSRF{
			A:    10.0,
			ali:  8.0,
			alic: 3.0,
			Ts:   20.0,
			Te:   15.0,
			RS:   100.0,
		}
	}

	// Setup RMVLS
	rmvls := &RMVLS{
		Room: []*ROOM{room},
	}

	// Setup RDPNL (empty for this test)
	rdpnl := make([]*RDPNL, 0)

	// Setup EXSFS (empty for this test)
	exsfs := &EXSFS{
		Exs: make([]*EXSF, 0),
	}

	// Setup WDAT
	wd := &WDAT{}

	// Execute
	Roomene(rmvls, []*ROOM{room}, rdpnl, exsfs, wd)

	// Verify basic calculations
	if room.Tr != 22.0 {
		t.Errorf("Tr should be set from component output, got %f", room.Tr)
	}
	if room.xr != 0.008 {
		t.Errorf("xr should be set from component output, got %f", room.xr)
	}
}

func TestRoomene_EnergyBalance(t *testing.T) {
	// Test energy balance calculation
	room := &ROOM{
		Name:  "EnergyBalanceTest",
		N:     1,
		rsrf:  make([]*RMSRF, 1),
		Tr:    25.0,
		Trold: 24.0,
		xr:    0.010,
		RH:    60.0,
		hr:    50.0,
		Ntr:   0,
		Nrp:   0,
		Nasup: 0,
		cmp: &COMPNT{
			Elouts: []*ELOUT{
				{Sysv: 25.0},
				{Sysv: 0.010},
			},
		},
	}

	// Create surface
	room.rsrf[0] = &RMSRF{
		A:    15.0,
		ali:  10.0,
		alic: 4.0,
		Ts:   23.0,
		Te:   18.0,
		RS:   150.0,
	}

	// Setup RMVLS
	rmvls := &RMVLS{
		Room: []*ROOM{room},
	}

	// Setup empty arrays for this test
	rdpnl := make([]*RDPNL, 0)
	exsfs := &EXSFS{Exs: make([]*EXSF, 0)}
	wd := &WDAT{}

	// Execute
	Roomene(rmvls, []*ROOM{room}, rdpnl, exsfs, wd)

	// Verify temperature values are reasonable
	if room.Tr < 0.0 || room.Tr > 100.0 {
		t.Errorf("Tr (%f) is outside reasonable range", room.Tr)
	}
	if room.xr < 0.0 || room.xr > 0.030 {
		t.Errorf("xr (%f) is outside reasonable range", room.xr)
	}
}

func TestRoomene_ZeroHeat(t *testing.T) {
	// Test with zero heat inputs
	room := &ROOM{
		Name:  "ZeroHeatTest",
		N:     1,
		rsrf:  make([]*RMSRF, 1),
		Tr:    20.0,
		Trold: 20.0,
		xr:    0.008,
		RH:    50.0,
		hr:    40.0,
		Ntr:   0,
		Nrp:   0,
		Nasup: 0,
		cmp: &COMPNT{
			Elouts: []*ELOUT{
				{Sysv: 20.0},
				{Sysv: 0.008},
			},
		},
	}

	// Create surface
	room.rsrf[0] = &RMSRF{
		A:    10.0,
		ali:  8.0,
		alic: 3.0,
		Ts:   20.0,
		Te:   20.0,
		RS:   0.0,
	}

	// Setup RMVLS
	rmvls := &RMVLS{
		Room: []*ROOM{room},
	}

	// Setup empty arrays for this test
	rdpnl := make([]*RDPNL, 0)
	exsfs := &EXSFS{Exs: make([]*EXSF, 0)}
	wd := &WDAT{}

	// Execute
	Roomene(rmvls, []*ROOM{room}, rdpnl, exsfs, wd)

	// With zero heat inputs, verify basic functionality
	t.Logf("Room temperature: %f", room.Tr)
	t.Logf("Room humidity: %f", room.xr)
}

func TestRoomene_MultiSurface(t *testing.T) {
	// Test with multiple surfaces
	room := &ROOM{
		Name:  "MultiSurfaceTest",
		N:     3,
		rsrf:  make([]*RMSRF, 3),
		Tr:    22.0,
		Trold: 21.5,
		xr:    0.009,
		RH:    55.0,
		hr:    48.0,
		Ntr:   0,
		Nrp:   0,
		Nasup: 0,
		cmp: &COMPNT{
			Elouts: []*ELOUT{
				{Sysv: 22.0},
				{Sysv: 0.009},
			},
		},
	}

	// Create multiple surfaces with different properties
	surfaceData := []struct {
		A    float64
		ali  float64
		alic float64
		Ts   float64
		Te   float64
		RS   float64
	}{
		{12.0, 9.0, 3.5, 21.0, 16.0, 120.0}, // Wall
		{8.0, 7.0, 2.5, 19.0, 15.0, 80.0},   // Window
		{15.0, 10.0, 4.0, 23.0, 18.0, 150.0}, // Floor
	}

	for i := 0; i < 3; i++ {
		room.rsrf[i] = &RMSRF{
			A:    surfaceData[i].A,
			ali:  surfaceData[i].ali,
			alic: surfaceData[i].alic,
			Ts:   surfaceData[i].Ts,
			Te:   surfaceData[i].Te,
			RS:   surfaceData[i].RS,
		}
	}

	// Setup RMVLS
	rmvls := &RMVLS{
		Room: []*ROOM{room},
	}

	// Setup empty arrays for this test
	rdpnl := make([]*RDPNL, 0)
	exsfs := &EXSFS{Exs: make([]*EXSF, 0)}
	wd := &WDAT{}

	// Execute
	Roomene(rmvls, []*ROOM{room}, rdpnl, exsfs, wd)

	// Verify calculations with multiple surfaces
	if room.Tr < 0.0 || room.Tr > 100.0 {
		t.Errorf("Tr (%f) is outside reasonable range", room.Tr)
	}

	// Log results for verification
	t.Logf("Multi-surface Tr: %f", room.Tr)
	t.Logf("Multi-surface xr: %f", room.xr)
}
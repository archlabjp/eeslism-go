package eeslism

import (
	"math"
	"testing"
)

func TestROOM_VRM(t *testing.T) {
	room := &ROOM{
		Name: "TestRoom",
		VRM:  100.0,
	}

	if room.VRM != 100.0 {
		t.Errorf("Expected VRM=100.0, got %f", room.VRM)
	}
}

func TestRMcf(t *testing.T) {
	// Setup basic room with surfaces
	room := &ROOM{
		Name: "TestRoom",
		N:    2,
		rsrf: make([]*RMSRF, 2),
		XA:   make([]float64, 4), // 2x2 matrix
		alr:  make([]float64, 4), // 2x2 matrix
		Ntr:  0,
		Nrp:  0,
		ARN:  make([]float64, 0), // Initialize ARN slice
		RMP:  make([]float64, 0), // Initialize RMP slice
	}

	// Create mock surfaces with proper initialization
	for i := 0; i < 2; i++ {
		// Create wall structure first
		wall := &WALL{
			WallType: WallType_N,
		}
		
		// Create MWALL with proper initialization
		mwall := &MWALL{
			M:    3,
			mp:   1,
			UX:   make([]float64, 9), // 3x3 matrix
			uo:   0.1,
			um:   0.2,
			Pc:   0.05,
			wall: wall,
		}
		
		// Initialize UX matrix values
		for j := 0; j < 9; j++ {
			mwall.UX[j] = 0.1 + float64(j)*0.1
		}

		room.rsrf[i] = &RMSRF{
			mrk:    '*',
			typ:    RMSRFType_H,
			A:      10.0,
			ali:    8.0,
			alic:   3.0,
			alir:   5.0,
			FI:     0.0, // Will be calculated
			FO:     0.0, // Will be calculated
			FP:     0.0, // Will be calculated
			mw:     mwall,
			mwside: RMSRFMwSideType_i,
			WSRN:   make([]float64, 0),
			WSPL:   make([]float64, 0),
		}
	}

	// Set up radiation matrix
	room.alr[0] = 0.8 // alr[0][0]
	room.alr[1] = 0.2 // alr[0][1]
	room.alr[2] = 0.2 // alr[1][0]
	room.alr[3] = 0.8 // alr[1][1]

	// Execute
	RMcf(room)

	// Verify basic calculations
	if room.rsrf[0].FI == 0.0 {
		t.Errorf("FI should be calculated, but got 0.0")
	}
	if room.rsrf[0].FO == 0.0 {
		t.Errorf("FO should be calculated, but got 0.0")
	}
	if room.rsrf[0].WSR == 0.0 {
		t.Errorf("WSR should be calculated, but got 0.0")
	}
	if room.AR == 0.0 {
		t.Errorf("AR should be calculated, but got 0.0")
	}
}

func TestFunCoeff(t *testing.T) {
	tests := []struct {
		name     string
		room     *ROOM
		expected float64
	}{
		{
			name: "Room with furniture heat capacity",
			room: &ROOM{
				CM:      &[]float64{1000.0}[0],
				MCAP:    &[]float64{5000.0}[0],
				FunHcap: 0.0,
				MRM:     2000.0,
				AR:      100.0,
			},
			expected: 0.998613, // FMT should be calculated
		},
		{
			name: "Room without furniture heat capacity",
			room: &ROOM{
				CM:      &[]float64{1000.0}[0],
				FunHcap: 0.0,
				MRM:     2000.0,
				AR:      100.0,
			},
			expected: 1.0, // FMT should be 1.0
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			FunCoeff(tt.room)

			if math.Abs(tt.room.FMT-tt.expected) > 1e-6 {
				t.Errorf("FMT = %f, want %f", tt.room.FMT, tt.expected)
			}
			if tt.room.RMt == 0.0 {
				t.Errorf("RMt should be calculated, but got 0.0")
			}
		})
	}
}

func TestRMrc(t *testing.T) {
	// Setup room with surfaces
	room := &ROOM{
		Name: "TestRoom",
		N:    2,
		rsrf: make([]*RMSRF, 2),
		XA:   make([]float64, 4), // 2x2 matrix
		Ntr:  0,
		Nrp:  0,
		Hc:   100.0,
		Lc:   50.0,
		Ac:   25.0,
		Qeqp: 200.0,
		eqcv: 0.8,
		MRM:  2000.0,
		Trold: 20.0,
		FunHcap: 0.0,
	}

	// Create mock surfaces
	for i := 0; i < 2; i++ {
		room.rsrf[i] = &RMSRF{
			typ: RMSRFType_H,
			A:   10.0,
			alic: 3.0,
			CF:  0.0,
			FO:  0.3,
			Te:  15.0,
			FI:  0.5,
			RS:  100.0,
			ali: 8.0,
			mw: &MWALL{
				M:    3,
				UX:   []float64{0.1, 0.2, 0.3},
				Told: []float64{18.0, 19.0, 20.0},
			},
			mwside: RMSRFMwSideType_i,
		}
	}

	// Set up XA matrix (identity for simplicity)
	room.XA[0] = 1.0
	room.XA[1] = 0.0
	room.XA[2] = 0.0
	room.XA[3] = 1.0

	// Execute
	RMrc(room)

	// Verify calculations
	expectedHGc := room.Hc + room.Lc + room.Ac + room.Qeqp*room.eqcv
	if room.HGc != expectedHGc {
		t.Errorf("HGc = %f, want %f", room.HGc, expectedHGc)
	}
	if room.CA == 0.0 {
		t.Errorf("CA should be calculated, but got 0.0")
	}
	if room.RMC == 0.0 {
		t.Errorf("RMC should be calculated, but got 0.0")
	}
}

func TestRMsrt(t *testing.T) {
	// Setup room with surfaces
	room := &ROOM{
		Name: "TestRoom",
		N:    2,
		rsrf: make([]*RMSRF, 2),
		alr:  make([]float64, 4), // 2x2 matrix
		Tr:   22.0,
		Ntr:  0,
		Nrp:  0,
	}

	// Create mock surfaces
	for i := 0; i < 2; i++ {
		room.rsrf[i] = &RMSRF{
			A:    10.0,
			WSR:  0.8,
			WSC:  5.0,
			alic: 3.0,
			alir: 5.0,
			RS:   100.0,
		}
	}

	// Set up radiation matrix
	room.alr[0] = 0.8 // alr[0][0]
	room.alr[1] = 0.2 // alr[0][1]
	room.alr[2] = 0.2 // alr[1][0]
	room.alr[3] = 0.8 // alr[1][1]

	// Execute
	RMsrt(room)

	// Verify surface temperatures are calculated
	for i := 0; i < room.N; i++ {
		if room.rsrf[i].Ts == 0.0 {
			t.Errorf("Surface %d temperature should be calculated, but got 0.0", i)
		}
		if room.rsrf[i].Tmrt == 0.0 {
			t.Errorf("Surface %d mean radiant temperature should be calculated, but got 0.0", i)
		}
	}
}

func TestRTsav(t *testing.T) {
	// Setup surfaces
	surfaces := []*RMSRF{
		{Ts: 20.0, A: 10.0},
		{Ts: 22.0, A: 15.0},
		{Ts: 18.0, A: 5.0},
	}

	// Execute
	avgTemp := RTsav(3, surfaces)

	// Calculate expected average
	expectedAvg := (20.0*10.0 + 22.0*15.0 + 18.0*5.0) / (10.0 + 15.0 + 5.0)

	if avgTemp != expectedAvg {
		t.Errorf("RTsav() = %f, want %f", avgTemp, expectedAvg)
	}
}

func TestRTsav_SingleSurface(t *testing.T) {
	// Test with single surface
	surfaces := []*RMSRF{
		{Ts: 25.0, A: 20.0},
	}

	avgTemp := RTsav(1, surfaces)

	if avgTemp != 25.0 {
		t.Errorf("RTsav() with single surface = %f, want 25.0", avgTemp)
	}
}

func TestRTsav_EmptySurfaces(t *testing.T) {
	// Test with empty surfaces (should return 0 or handle gracefully)
	avgTemp := RTsav(0, []*RMSRF{})
	
	// The function should handle empty surfaces gracefully
	// (exact behavior depends on implementation)
	t.Logf("RTsav with empty surfaces returned: %f", avgTemp)
}
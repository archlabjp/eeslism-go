package eeslism

import (
	"testing"
)

func TestWalli(t *testing.T) {
	// Setup basic materials
	materials := []BMLST{
		{Mcode: "concrete", Cond: 1.6, Cro: 1.8},
		{Mcode: "insulation", Cond: 0.04, Cro: 0.02},
		{Mcode: "ali", Cond: 25.0, Cro: 0.0},
		{Mcode: "alo", Cond: 25.0, Cro: 0.0},
	}

	// Setup wall with elements
	wall := &WALL{
		name:     "TestWall",
		WallType: WallType_N,
		N:        2,
		welm: []WELM{
			{Code: "concrete", L: 0.15, ND: 2},
			{Code: "insulation", L: 0.10, ND: 1},
		},
		Ip: 0,
	}

	// Setup PCM (empty for this test)
	var pcm []*PCM

	// Execute
	Walli(len(materials), materials, wall, pcm)

	// Verify basic calculations
	if wall.Rwall <= 0.0 {
		t.Errorf("Rwall should be positive, got %f", wall.Rwall)
	}
	if wall.CAPwall <= 0.0 {
		t.Errorf("CAPwall should be positive, got %f", wall.CAPwall)
	}
	if wall.M <= 0 {
		t.Errorf("M should be positive, got %d", wall.M)
	}
	if len(wall.R) != wall.N {
		t.Errorf("R array length should be %d, got %d", wall.N, len(wall.R))
	}
	if len(wall.CAP) != wall.N {
		t.Errorf("CAP array length should be %d, got %d", wall.N, len(wall.CAP))
	}
}

func TestWalli_WithPCM(t *testing.T) {
	// Setup materials
	materials := []BMLST{
		{Mcode: "concrete", Cond: 1.6, Cro: 1.8},
		{Mcode: "pcm_base", Cond: 0.2, Cro: 1.5},
	}

	// Setup PCM
	pcm := []*PCM{
		{
			Name:    "TestPCM",
			Spctype: 'm',
			Ctype:   1,
			Cros:    1500.0,
			Crol:    2000.0,
			Ql:      200000.0,
			Ts:      23.0,
			Tl:      25.0,
			Tp:      27.0,
		},
	}

	// Setup wall with PCM element
	wall := &WALL{
		name:     "TestWallPCM",
		WallType: WallType_N,
		N:        2,
		welm: []WELM{
			{Code: "concrete", L: 0.15, ND: 2},
			{Code: "pcm_base(TestPCM_0.3)", L: 0.05, ND: 1}, // 30% PCM content
		},
		Ip: 0,
	}

	// Execute
	Walli(len(materials), materials, wall, pcm)

	// Verify PCM flag is set
	if !wall.PCMflg {
		t.Errorf("PCMflg should be true when PCM is present")
	}
	if wall.PCM[1] == nil {
		t.Errorf("PCM should be assigned to layer 1")
	}
	if wall.PCMrate[1] != 0.3 {
		t.Errorf("PCM rate should be 0.3, got %f", wall.PCMrate[1])
	}
}

func TestWallBasicStructure(t *testing.T) {
	t.Run("Wall structure initialization", func(t *testing.T) {
		// Test basic wall structure creation
		wall := &WALL{
			name:     "TestWall",
			WallType: WallType_N, // Normal wall type
			Ei: 0.9,
			Eo: 0.9,
		}
		
		if wall.name != "TestWall" {
			t.Errorf("Wall name = %s, want TestWall", wall.name)
		}
		if wall.WallType != WallType_N {
			t.Errorf("WallType = %c, want %c", wall.WallType, WallType_N)
		}
		if wall.Ei != 0.9 {
			t.Errorf("Ei = %f, want 0.9", wall.Ei)
		}
		if wall.Eo != 0.9 {
			t.Errorf("Eo = %f, want 0.9", wall.Eo)
		}
	})

	t.Run("Different wall types", func(t *testing.T) {
		tests := []struct {
			name     string
			wallType WALLType
		}{
			{"Normal wall", WallType_N},
			{"Solar collector wall", WallType_C},
			{"Panel wall", WallType_P},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				wall := &WALL{
					name:     tt.name,
					WallType: tt.wallType,
				}
				if wall.WallType != tt.wallType {
					t.Errorf("WallType = %c, want %c", wall.WallType, tt.wallType)
				}
			})
		}
	})
}
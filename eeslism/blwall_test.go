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

func TestWallfdc(t *testing.T) {
	// Set global DTM for time step (seconds)
	oldDTM := DTM
	DTM = 3600.0 // 1 hour
	defer func() { DTM = oldDTM }()

	t.Run("Basic wall without PCM", func(t *testing.T) {
		M := 3 // 3 internal nodes
		mp := 1

		// Thermal resistance for each layer (m2K/W) - need M+1 elements
		res := []float64{0.05, 0.1, 0.05, 0.02}
		// Heat capacity for each layer (J/m2K) - need M+2 elements for cap[m+1] access
		cap := []float64{50000.0, 100000.0, 100000.0, 50000.0, 30000.0}

		// Output matrix
		UX := make([]float64, M*M)

		var uo, um, Pc float64
		Wp := 0.0 // No panel

		// Create wall structure
		wall := &WALL{
			WallType:    WallType_N,
			M:           M,
			L:           []float64{0.1, 0.15, 0.1},
			PCMLyr:      make([]*PCM, M+1),
			PCMrateLyr:  make([]float64, M+1),
		}

		// Create PCM state array (all nil for non-PCM wall)
		pcmstate := make([]*PCMSTATE, M+1)
		for i := range pcmstate {
			pcmstate[i] = &PCMSTATE{}
		}

		// Temperature arrays
		Told := make([]float64, M+1)
		Twd := make([]float64, M+1)
		for i := range Told {
			Told[i] = 20.0
			Twd[i] = 20.0
		}

		// Execute
		Wallfdc(M, mp, res, cap, Wp, UX, &uo, &um, &Pc, wall.WallType, nil, nil, nil, wall, Told, Twd, pcmstate)

		// Verify outputs
		if uo <= 0 {
			t.Errorf("uo should be positive, got %f", uo)
		}
		if um <= 0 {
			t.Errorf("um should be positive, got %f", um)
		}
		if Pc != 0 {
			t.Errorf("Pc should be 0 for non-panel wall, got %f", Pc)
		}

		// Check that UX matrix has reasonable values
		// Diagonal elements should be > 1
		for m := 0; m < M; m++ {
			if UX[m*M+m] <= 0 {
				t.Errorf("UX diagonal element [%d][%d] should be positive, got %f", m, m, UX[m*M+m])
			}
		}

		t.Logf("Non-PCM wall: uo=%f, um=%f, Pc=%f", uo, um, Pc)
	})

	t.Run("Panel wall (WallType P)", func(t *testing.T) {
		M := 3
		mp := 1 // Panel at layer 1

		res := []float64{0.05, 0.1, 0.05, 0.02}
		cap := []float64{50000.0, 100000.0, 100000.0, 50000.0, 30000.0}
		UX := make([]float64, M*M)

		var uo, um, Pc float64
		Wp := 100.0 // Panel coefficient > 0

		wall := &WALL{
			WallType:   WallType_P,
			M:          M,
			L:          []float64{0.1, 0.15, 0.1},
			PCMLyr:     make([]*PCM, M+1),
			PCMrateLyr: make([]float64, M+1),
		}

		pcmstate := make([]*PCMSTATE, M+1)
		for i := range pcmstate {
			pcmstate[i] = &PCMSTATE{}
		}

		Told := make([]float64, M+1)
		Twd := make([]float64, M+1)
		for i := range Told {
			Told[i] = 20.0
			Twd[i] = 20.0
		}

		Wallfdc(M, mp, res, cap, Wp, UX, &uo, &um, &Pc, wall.WallType, nil, nil, nil, wall, Told, Twd, pcmstate)

		// For panel wall, Pc should be positive
		if Pc <= 0 {
			t.Errorf("Pc should be positive for panel wall, got %f", Pc)
		}

		t.Logf("Panel wall: uo=%f, um=%f, Pc=%f", uo, um, Pc)
	})

	t.Run("Wall with PCM layer", func(t *testing.T) {
		M := 3
		mp := 1

		res := []float64{0.05, 0.1, 0.05, 0.02}
		cap := []float64{50000.0, 100000.0, 100000.0, 50000.0, 30000.0}
		UX := make([]float64, M*M)

		var uo, um, Pc float64
		Wp := 0.0

		// Create PCM
		pcm := &PCM{
			Name:     "TestPCM",
			Spctype:  'm',      // Model-based
			Ctype:    1,
			Cros:     1500000.0, // Solid heat capacity (J/m3K)
			Crol:     2000000.0, // Liquid heat capacity
			Ql:       200000.0,  // Latent heat
			Ts:       20.0,      // Solidus temperature
			Tl:       25.0,      // Liquidus temperature
			Tp:       23.0,      // Peak temperature
			DivTemp:  10,
			Conds:    0.5,       // Solid conductivity
			Condl:    0.3,       // Liquid conductivity
			Condtype: 'm',
			AveTemp:  'y',
		}

		wall := &WALL{
			WallType:   WallType_N,
			M:          M,
			L:          []float64{0.1, 0.05, 0.1},
			PCMLyr:     []*PCM{nil, nil, pcm, nil}, // PCM in layer 2 (between nodes 1 and 2)
			PCMrateLyr: []float64{0.0, 0.0, 0.3, 0.0}, // 30% PCM content
		}

		pcmstate := make([]*PCMSTATE, M+1)
		for i := range pcmstate {
			pcmstate[i] = &PCMSTATE{}
		}

		// Temperature arrays
		Told := []float64{18.0, 20.0, 22.0, 24.0}
		Twd := []float64{19.0, 21.0, 23.0, 25.0}

		Wallfdc(M, mp, res, cap, Wp, UX, &uo, &um, &Pc, wall.WallType, nil, nil, nil, wall, Told, Twd, pcmstate)

		if uo <= 0 {
			t.Errorf("uo should be positive for PCM wall, got %f", uo)
		}
		if um <= 0 {
			t.Errorf("um should be positive for PCM wall, got %f", um)
		}

		t.Logf("PCM wall: uo=%f, um=%f", uo, um)
	})

	t.Run("Multiple layer wall", func(t *testing.T) {
		M := 5 // 5 internal nodes
		mp := 2

		// Need M+1 elements for res, M+2 elements for cap
		res := make([]float64, M+1)
		cap := make([]float64, M+2)
		for i := 0; i <= M; i++ {
			res[i] = 0.05 + float64(i)*0.01
			cap[i] = 50000.0 + float64(i)*10000.0
		}
		cap[M+1] = 30000.0 // Extra element for cap[m+1] access

		UX := make([]float64, M*M)
		var uo, um, Pc float64
		Wp := 0.0

		wall := &WALL{
			WallType:   WallType_N,
			M:          M,
			L:          make([]float64, M),
			PCMLyr:     make([]*PCM, M+1),
			PCMrateLyr: make([]float64, M+1),
		}
		for i := range wall.L {
			wall.L[i] = 0.1
		}

		pcmstate := make([]*PCMSTATE, M+1)
		for i := range pcmstate {
			pcmstate[i] = &PCMSTATE{}
		}

		Told := make([]float64, M+1)
		Twd := make([]float64, M+1)
		for i := range Told {
			Told[i] = 20.0
			Twd[i] = 20.0
		}

		Wallfdc(M, mp, res, cap, Wp, UX, &uo, &um, &Pc, wall.WallType, nil, nil, nil, wall, Told, Twd, pcmstate)

		if uo <= 0 {
			t.Errorf("uo should be positive for 5-layer wall, got %f", uo)
		}
		if um <= 0 {
			t.Errorf("um should be positive for 5-layer wall, got %f", um)
		}

		t.Logf("5-layer wall: uo=%f, um=%f", uo, um)
	})
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
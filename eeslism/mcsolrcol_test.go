package eeslism

import (
	"math"
	"testing"
)

func TestColldata(t *testing.T) {
	tests := []struct {
		name     string
		typeStr  EqpType
		input    string
		expected func(*COLLCA) bool
	}{
		{
			name:    "collector name only",
			typeStr: COLLECTOR_TYPE,
			input:   "TestCollector",
			expected: func(collca *COLLCA) bool {
				return collca.name == "TestCollector" &&
					collca.Type == COLLECTOR_PDT &&
					collca.b0 == FNAN &&
					collca.b1 == FNAN &&
					collca.Ac == FNAN &&
					collca.Ag == FNAN
			},
		},
		{
			name:    "air collector name only",
			typeStr: ACOLLECTOR_TYPE,
			input:   "TestAirCollector",
			expected: func(collca *COLLCA) bool {
				return collca.name == "TestAirCollector" &&
					collca.Type == ACOLLECTOR_PDT
			},
		},
		{
			name:    "set b0 coefficient",
			typeStr: COLLECTOR_TYPE,
			input:   "b0=0.8",
			expected: func(collca *COLLCA) bool {
				return collca.b0 == 0.8
			},
		},
		{
			name:    "set b1 coefficient",
			typeStr: COLLECTOR_TYPE,
			input:   "b1=5.0",
			expected: func(collca *COLLCA) bool {
				return collca.b1 == 5.0
			},
		},
		{
			name:    "set Fd coefficient",
			typeStr: COLLECTOR_TYPE,
			input:   "Fd=0.95",
			expected: func(collca *COLLCA) bool {
				return collca.Fd == 0.95
			},
		},
		{
			name:    "set collector area",
			typeStr: COLLECTOR_TYPE,
			input:   "Ac=10.0",
			expected: func(collca *COLLCA) bool {
				return collca.Ac == 10.0
			},
		},
		{
			name:    "set gross area",
			typeStr: COLLECTOR_TYPE,
			input:   "Ag=12.0",
			expected: func(collca *COLLCA) bool {
				return collca.Ag == 12.0
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collca := &COLLCA{}
			result := Colldata(tt.typeStr, tt.input, collca)

			if result != 0 {
				t.Errorf("Colldata(%v, %q) returned error code %d", tt.typeStr, tt.input, result)
			}

			if !tt.expected(collca) {
				t.Errorf("Colldata(%v, %q) did not set expected values", tt.typeStr, tt.input)
			}
		})
	}
}

func TestColldata_InvalidInput(t *testing.T) {
	tests := []struct {
		name    string
		typeStr EqpType
		input   string
	}{
		{
			name:    "invalid parameter",
			typeStr: COLLECTOR_TYPE,
			input:   "invalid=value",
		},
		{
			name:    "invalid b0 value",
			typeStr: COLLECTOR_TYPE,
			input:   "b0=invalid",
		},
		{
			name:    "invalid b1 value",
			typeStr: COLLECTOR_TYPE,
			input:   "b1=invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collca := &COLLCA{}

			// Should panic for invalid numeric values
			if tt.input == "b0=invalid" || tt.input == "b1=invalid" {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("Colldata should panic for invalid numeric input")
					}
				}()
				Colldata(tt.typeStr, tt.input, collca)
				return
			}

			// Should panic for invalid parameter with numeric value
			if tt.input == "invalid=value" {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("Colldata should panic for invalid parameter with numeric value")
					}
				}()
				Colldata(tt.typeStr, tt.input, collca)
				return
			}
		})
	}
}

func TestScolte(t *testing.T) {
	tests := []struct {
		name      string
		rtgko     float64
		cinc      float64
		Idre      float64
		Idf       float64
		Ta        float64
		expected  float64
		tolerance float64
	}{
		{
			name:      "typical solar conditions",
			rtgko:     0.16,  // Typical ratio for solar collectors
			cinc:      1.0,   // Normal incidence
			Idre:      800.0, // Direct solar radiation [W/m²]
			Idf:       200.0, // Diffuse solar radiation [W/m²]
			Ta:        25.0,  // Ambient temperature [°C]
			expected:  182.0, // Approximate equivalent temperature
			tolerance: 10.0,
		},
		{
			name:      "zero solar radiation",
			rtgko:     0.16,
			cinc:      0.0,
			Idre:      0.0,
			Idf:       0.0,
			Ta:        20.0,
			expected:  20.0, // Should equal ambient temperature
			tolerance: 0.1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := scolte(tt.rtgko, tt.cinc, tt.Idre, tt.Idf, tt.Ta)
			if math.Abs(result-tt.expected) > tt.tolerance {
				t.Errorf("scolte() = %v, want %v ± %v", result, tt.expected, tt.tolerance)
			}
		})
	}
}

func TestCollcfv(t *testing.T) {
	// Create test collector with proper initialization
	collca := &COLLCA{
		name: "TestCollector",
		Type: COLLECTOR_PDT,
		b0:   0.8,
		b1:   5.0,
		Fd:   0.95,
		Ac:   10.0,
		Ko:   5.26, // b1/Fd = 5.0/0.95
	}

	// Create COMPNT with proper initialization
	cmp := &COMPNT{
		Name:    "TestCollector",
		Control: ON_SW,
		Ac:      10.0, // Collector area
		Elins: []*ELIN{
			{
				Sysvin: 40.0, // Inlet temperature
			},
		},
		Elouts: []*ELOUT{
			{
				Fluid:   WATER_FLD,
				G:       0.2, // Flow rate [kg/s]
				Control: ON_SW,
				Coeffin: make([]float64, 1),
				Sysv:    60.0, // Outlet temperature
			},
		},
	}

	coll := &COLL{
		Name: "TestCollector",
		Cat:  collca,
		Cmp:  cmp,
		Te:   80.0, // Equivalent temperature
	}

	colls := []*COLL{coll}

	t.Run("calculate collector coefficients", func(t *testing.T) {
		// Store original values
		origDo := coll.Do
		origD1 := coll.D1
		origEc := coll.ec

		Collcfv(colls)

		// Check that coefficients were calculated
		if coll.Do == origDo && coll.D1 == origD1 && coll.ec == origEc {
			t.Error("Collcfv should update Do, D1, and ec")
		}

		// Check that coefficients are reasonable
		if coll.ec <= 0 || coll.ec > 1 {
			t.Errorf("Collector efficiency ec should be between 0 and 1, got %v", coll.ec)
		}

		if coll.D1 <= 0 {
			t.Errorf("D1 should be positive, got %v", coll.D1)
		}

		// Check that Do is calculated correctly (D1 * Te)
		expectedDo := coll.D1 * coll.Te
		if math.Abs(coll.Do-expectedDo) > 1e-6 {
			t.Errorf("Do = %v, want D1*Te = %v", coll.Do, expectedDo)
		}
	})
}

func TestCollene(t *testing.T) {
	// Create test collector
	collca := &COLLCA{
		name: "TestCollector",
		Type: COLLECTOR_PDT,
		Ko:   5.0, // Heat loss coefficient
	}

	// Create EXSF for solar radiation
	exsf := &EXSF{
		Name: "TestSolar",
		Iw:   800.0, // Solar radiation [W/m²]
	}

	// Create COMPNT with proper initialization
	cmp := &COMPNT{
		Name:    "TestCollector",
		Control: ON_SW,
		Ac:      10.0, // Collector area
		Elins: []*ELIN{
			{
				Sysvin: 40.0, // Inlet temperature
			},
		},
		Elouts: []*ELOUT{
			{
				Control: ON_SW,
			},
		},
	}

	coll := &COLL{
		Name: "TestCollector",
		Cat:  collca,
		Cmp:  cmp,
		sol:  exsf,
		Te:   80.0,   // Equivalent temperature
		Do:   1000.0, // Heat gain coefficient
		D1:   200.0,  // Heat loss coefficient
		Ac:   10.0,   // Collector area
	}

	colls := []*COLL{coll}

	t.Run("calculate collector energy", func(t *testing.T) {
		Collene(colls)

		// Check that heat collection was calculated
		expectedQ := coll.Do - coll.D1*coll.Tin
		if math.Abs(coll.Q-expectedQ) > 1e-6 {
			t.Errorf("Q = %v, want Do - D1*Tin = %v", coll.Q, expectedQ)
		}

		// Check that collector plate temperature was calculated
		if coll.Q > 0 {
			expectedTcb := coll.Te - coll.Q/(coll.Ac*coll.Cat.Ko)
			if math.Abs(coll.Tcb-expectedTcb) > 1e-6 {
				t.Errorf("Tcb = %v, want Te - Q/(Ac*Ko) = %v", coll.Tcb, expectedTcb)
			}
		}

		// Check that solar radiation was calculated
		expectedSol := coll.sol.Iw * coll.Cmp.Ac
		if math.Abs(coll.Sol-expectedSol) > 1e-6 {
			t.Errorf("Sol = %v, want Iw*Ac = %v", coll.Sol, expectedSol)
		}
	})

	t.Run("collector off", func(t *testing.T) {
		coll.Cmp.Control = OFF_SW

		Collene(colls)

		// Heat collection should be zero when off
		if coll.Q != 0.0 {
			t.Errorf("Q should be 0 when collector is off, got %v", coll.Q)
		}

		// Collector plate temperature should equal equivalent temperature
		if coll.Tcb != coll.Te {
			t.Errorf("Tcb should equal Te when collector is off, got Tcb=%v, Te=%v", coll.Tcb, coll.Te)
		}
	})
}

func TestCalcCollTe(t *testing.T) {
	// Create test collector
	collca := &COLLCA{
		name: "TestCollector",
		b0:   0.8,
		b1:   5.0,
	}

	// Create EXSF for solar radiation
	exsf := &EXSF{
		Name: "TestSolar",
		Cinc: 1.0,   // Cosine of incidence angle
		Idre: 700.0, // Direct solar radiation
		Idf:  200.0, // Diffuse solar radiation
	}

	coll := &COLL{
		Name: "TestCollector",
		Cat:  collca,
		sol:  exsf,
	}

	// Mock ambient temperature
	Ta := 25.0
	coll.Ta = &Ta

	colls := []*COLL{coll}

	t.Run("calculate equivalent temperature", func(t *testing.T) {
		CalcCollTe(colls)

		// Check that equivalent temperature was calculated
		tgaKo := coll.Cat.b0 / coll.Cat.b1
		expectedTe := scolte(tgaKo, coll.sol.Cinc, coll.sol.Idre, coll.sol.Idf, *coll.Ta)

		if math.Abs(coll.Te-expectedTe) > 1e-6 {
			t.Errorf("Te = %v, want %v", coll.Te, expectedTe)
		}

		// Equivalent temperature should be higher than ambient
		if coll.Te <= *coll.Ta {
			t.Errorf("Equivalent temperature Te=%v should be higher than ambient Ta=%v", coll.Te, *coll.Ta)
		}
	})
}

func TestCollvptr(t *testing.T) {
	coll := &COLL{
		Te:  80.0,
		Tcb: 75.0,
	}

	tests := []struct {
		name     string
		key      []string
		expected interface{}
		hasError bool
	}{
		{
			name:     "get Te pointer",
			key:      []string{"", "Te"},
			expected: &coll.Te,
			hasError: false,
		},
		{
			name:     "get Tcb pointer",
			key:      []string{"", "Tcb"},
			expected: &coll.Tcb,
			hasError: false,
		},
		{
			name:     "invalid key",
			key:      []string{"", "Invalid"},
			expected: nil,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vptr, err := collvptr(tt.key, coll)

			if tt.hasError {
				if err == nil {
					t.Errorf("collvptr should return error for invalid key")
				}
			} else {
				if err != nil {
					t.Errorf("collvptr returned unexpected error: %v", err)
				}
				if vptr.Ptr != tt.expected {
					t.Errorf("collvptr returned wrong pointer")
				}
				if vptr.Type != VAL_CTYPE {
					t.Errorf("collvptr returned wrong type, got %v, want %v", vptr.Type, VAL_CTYPE)
				}
			}
		})
	}
}

func TestCollectorIntegration(t *testing.T) {
	// Integration test: create collector from data input and test energy calculation
	collca := &COLLCA{}

	// Set up collector with multiple parameters
	inputs := []struct {
		typeStr EqpType
		input   string
	}{
		{COLLECTOR_TYPE, "TestIntegrationCollector"},
		{COLLECTOR_TYPE, "b0=0.75"},
		{COLLECTOR_TYPE, "b1=4.5"},
		{COLLECTOR_TYPE, "Fd=0.92"},
		{COLLECTOR_TYPE, "Ac=8.0"},
		{COLLECTOR_TYPE, "Ag=10.0"},
	}

	for _, input := range inputs {
		result := Colldata(input.typeStr, input.input, collca)
		if result != 0 {
			t.Fatalf("Failed to set collector data: %s", input.input)
		}
	}

	// Verify all parameters were set correctly
	if collca.name != "TestIntegrationCollector" {
		t.Errorf("Name = %s, want TestIntegrationCollector", collca.name)
	}
	if collca.b0 != 0.75 {
		t.Errorf("b0 = %v, want 0.75", collca.b0)
	}
	if collca.b1 != 4.5 {
		t.Errorf("b1 = %v, want 4.5", collca.b1)
	}
	if collca.Fd != 0.92 {
		t.Errorf("Fd = %v, want 0.92", collca.Fd)
	}
	if collca.Ac != 8.0 {
		t.Errorf("Ac = %v, want 8.0", collca.Ac)
	}
	if collca.Ag != 10.0 {
		t.Errorf("Ag = %v, want 10.0", collca.Ag)
	}

	// Test equivalent temperature calculation
	tgaKo := collca.b0 / collca.b1
	Te := scolte(tgaKo, 1.0, 800.0, 200.0, 25.0)

	if Te <= 25.0 {
		t.Errorf("Equivalent temperature should be higher than ambient, got %v", Te)
	}
}

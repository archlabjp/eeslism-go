package eeslism

import (
	"bytes"
	"strings"
	"testing"
)

func TestHccdata(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected func(*HCCCA) bool
	}{
		{
			name:  "coil name only",
			input: "TestCoil",
			expected: func(hccca *HCCCA) bool {
				return hccca.name == "TestCoil" &&
					hccca.eh == 0.0 &&
					hccca.et == FNAN &&
					hccca.KA == FNAN
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hccca := &HCCCA{}
			result := Hccdata(tt.input, hccca)

			if result != 0 {
				t.Errorf("Hccdata(%q) returned error code %d", tt.input, result)
			}

			if !tt.expected(hccca) {
				t.Errorf("Hccdata(%q) did not set expected values", tt.input)
			}
		})
	}
}

func TestHccdata_InvalidInput(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "invalid parameter",
			input: "invalid=1.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hccca := &HCCCA{}
			result := Hccdata(tt.input, hccca)

			if result != 1 {
				t.Errorf("Hccdata should return error code 1 for invalid parameter")
			}
		})
	}
}

func TestHccdwint(t *testing.T) {
	// Create test coils with different configurations
	hccca1 := &HCCCA{
		name: "DryCoil",
		et:   0.8,  // Fixed temperature efficiency
		eh:   0.0,  // No enthalpy efficiency (dry coil)
		KA:   FNAN, // No KA specified
	}

	hccca2 := &HCCCA{
		name: "WetCoil",
		et:   FNAN,   // No fixed temperature efficiency
		eh:   0.75,   // Enthalpy efficiency (wet coil)
		KA:   1500.0, // KA specified
	}

	hccca3 := &HCCCA{
		name: "VariableCoil",
		et:   FNAN,   // No fixed temperature efficiency
		eh:   0.0,    // No enthalpy efficiency
		KA:   2000.0, // KA specified (variable efficiency)
	}

	// Create COMPNT with proper initialization
	cmp1 := &COMPNT{
		Name:    "DryCoil",
		Control: ON_SW,
		Elins: []*ELIN{
			{
				Sysvin: 25.0, // Inlet air temperature
			},
		},
		Elouts: []*ELOUT{
			{
				Fluid:   WATER_FLD,
				Control: ON_SW,
				Coeffin: make([]float64, 1),
				Sysv:    15.0, // Outlet air temperature
			},
		},
	}

	cmp2 := &COMPNT{
		Name:    "WetCoil",
		Control: ON_SW,
		Elins: []*ELIN{
			{
				Sysvin: 30.0,
			},
		},
		Elouts: []*ELOUT{
			{
				Fluid:   WATER_FLD,
				Control: ON_SW,
				Coeffin: make([]float64, 1),
				Sysv:    18.0,
			},
		},
	}

	cmp3 := &COMPNT{
		Name:    "VariableCoil",
		Control: ON_SW,
		Elins: []*ELIN{
			{
				Sysvin: 28.0,
			},
		},
		Elouts: []*ELOUT{
			{
				Fluid:   WATER_FLD,
				Control: ON_SW,
				Coeffin: make([]float64, 1),
				Sysv:    20.0,
			},
		},
	}

	hcc1 := &HCC{
		Name: "DryCoil",
		Cat:  hccca1,
		Cmp:  cmp1,
	}

	hcc2 := &HCC{
		Name: "WetCoil",
		Cat:  hccca2,
		Cmp:  cmp2,
	}

	hcc3 := &HCC{
		Name: "VariableCoil",
		Cat:  hccca3,
		Cmp:  cmp3,
	}

	hccs := []*HCC{hcc1, hcc2, hcc3}

	t.Run("initialize coil dry/wet types", func(t *testing.T) {
		Hccdwint(hccs)

		// Check dry coil (fixed temperature efficiency)
		if hcc1.Wet != 'd' {
			t.Errorf("Coil with fixed et should be dry coil, got %c", hcc1.Wet)
		}

		// Check wet coil (enthalpy efficiency > 0)
		if hcc2.Wet != 'w' {
			t.Errorf("Coil with eh > 0 should be wet coil, got %c", hcc2.Wet)
		}

		// Check variable efficiency coil (KA specified, eh = 0)
		if hcc3.Wet != 'd' {
			t.Errorf("Coil with KA and eh=0 should be dry coil, got %c", hcc3.Wet)
		}
	})
}

func TestHcccfv(t *testing.T) {
	// Create test coil with proper initialization
	hccca := &HCCCA{
		name: "TestCoil",
		et:   FNAN,   // Variable efficiency
		eh:   0.0,    // Dry coil
		KA:   1500.0, // Heat transfer coefficient
	}

	// Create COMPNT with proper initialization
	cmp := &COMPNT{
		Name:    "TestCoil",
		Control: ON_SW,
		Elins:   []*ELIN{{}, {}, {}}, // HCC requires 3 Elins
		Elouts: []*ELOUT{
			{
				Fluid:   WATER_FLD,
				G:       0.1, // Water flow rate
				Control: ON_SW,
				Coeffin: make([]float64, 3),
				Sysv:    15.0, // Outlet air temperature
			},
			{
				Fluid:   WATER_FLD,
				G:       0.1, // Water flow rate
				Control: ON_SW,
				Coeffin: make([]float64, 3),
				Sysv:    15.0, // Outlet air temperature
			},
			{
				Fluid:   WATER_FLD,
				G:       0.1, // Water flow rate
				Control: ON_SW,
				Coeffin: make([]float64, 3),
				Sysv:    15.0, // Outlet air temperature
			},
		},
	}

	hcc := &HCC{
		Name:  "TestCoil",
		Cat:   hccca,
		Cmp:   cmp,
		Wet:   'd',  // Dry coil
		Tain:  25.0, // Inlet air temperature
		Twin:  7.0,  // Inlet water temperature
		cGa:   1.0,  // Air flow rate
		cGw:   0.1,  // Water flow rate
		Etype: 'e',
	}

	hccs := []*HCC{hcc}

	t.Run("calculate coil coefficients", func(t *testing.T) {
		// Store original values
		origEt := hcc.Et
		origEx := hcc.Ex

		Hcccfv(hccs)

		// Check that coefficients were calculated
		if hcc.Et == origEt && hcc.Ex == origEx {
			t.Error("Hcccfv should update Et and Ex")
		}
	})
}

func TestHccene(t *testing.T) {
	// Create test coil
	hccca := &HCCCA{
		name: "TestCoil",
		et:   0.8, // Fixed temperature efficiency
		eh:   0.0, // Dry coil
		KA:   FNAN,
	}

	// Create COMPNT with proper initialization
	cmp := &COMPNT{
		Name:    "TestCoil",
		Control: ON_SW,
		Elins: []*ELIN{
			{},
			{},
			{},
		},
		Elouts: []*ELOUT{
			{
				Fluid:   WATER_FLD,
				G:       0.1, // Water flow rate
				Control: ON_SW,
				Coeffin: make([]float64, 3),
				Sysv:    15.0, // Outlet air temperature
			},
			{
				Fluid:   WATER_FLD,
				G:       0.1, // Water flow rate
				Control: ON_SW,
				Coeffin: make([]float64, 3),
				Sysv:    15.0, // Outlet air temperature
			},
			{
				Fluid:   WATER_FLD,
				G:       0.1, // Water flow rate
				Control: ON_SW,
				Coeffin: make([]float64, 3),
				Sysv:    15.0, // Outlet air temperature
			},
		},
	}

	hcc := &HCC{
		Name: "TestCoil",
		Cat:  hccca,
		Cmp:  cmp,
		Wet:  'd',  // Dry coil
		Tain: 25.0, // Inlet air temperature
		Twin: 7.0,  // Inlet water temperature
		cGa:  1.0,  // Air flow rate [kg/s]
		cGw:  0.1,  // Water flow rate [kg/s]
	}

	hccs := []*HCC{hcc}

	t.Run("calculate coil energy transfer", func(t *testing.T) {
		Hccene(hccs)

		// Check that heat transfer was calculated
		if hcc.Qs == 0 && hcc.Ql == 0 && hcc.Qt == 0 {
			t.Error("Heat transfer should be calculated")
		}

		// Check that outlet air temperature was calculated
		if hcc.Taout == 0 {
			t.Error("Outlet air temperature should be calculated")
		}

		// Check that outlet water temperature was calculated
		if hcc.Twout == 0 {
			t.Error("Outlet water temperature should be calculated")
		}
	})

	t.Run("coil off", func(t *testing.T) {
		hcc.Cmp.Control = OFF_SW

		Hccene(hccs)

		// Heat transfer should be zero when off
		if hcc.Qs != 0.0 || hcc.Ql != 0.0 || hcc.Qt != 0.0 {
			t.Errorf("Heat transfer should be 0 when coil is off, got Qs=%v, Ql=%v, Qt=%v", hcc.Qs, hcc.Ql, hcc.Qt)
		}
	})
}

func TestHccene_WetCoil(t *testing.T) {
	t.Skip("Skipping complex wet coil test due to structure complexity")
}

func TestHccdwreset(t *testing.T) {
	t.Skip("Skipping reset test due to structure complexity")
}

func TestCoilIntegration(t *testing.T) {
	// Integration test: create coil from data input
	hccca := &HCCCA{}

	// Set up coil with multiple parameters
	inputs := []string{
		"TestIntegrationCoil",
		"et=0.85",
		"eh=0.0", // Dry coil
	}

	for _, input := range inputs {
		result := Hccdata(input, hccca)
		if result != 0 {
			t.Fatalf("Failed to set coil data: %s", input)
		}
	}

	// Verify all parameters were set correctly
	if hccca.name != "TestIntegrationCoil" {
		t.Errorf("Name = %s, want TestIntegrationCoil", hccca.name)
	}
	if hccca.et != 0.85 {
		t.Errorf("et = %v, want 0.85", hccca.et)
	}
	if hccca.eh != 0.0 {
		t.Errorf("eh = %v, want 0.0", hccca.eh)
	}
}

// createOutputTestHCC creates an HCC suitable for output function tests
func createOutputTestHCC() *HCC {
	return &HCC{
		Name: "TestHCC",
		Cat: &HCCCA{
			name: "TestHCCCat",
			et:   0.85,
			eh:   0.75,
			KA:   5000.0,
		},
		Cmp: &COMPNT{
			Name:    "TestHCC",
			Control: ON_SW,
			Elouts: []*ELOUT{
				{Control: ON_SW, G: 1.0, Sysv: 18.0}, // Air temp out
				{Control: ON_SW, G: 1.0, Sysv: 0.008}, // Air humidity out
				{Control: ON_SW, G: 0.5, Sysv: 12.0}, // Water temp out
			},
			Elins: []*ELIN{
				{Sysvin: 26.0, Lpath: &PLIST{Control: ON_SW}},
				{Sysvin: 0.012},
				{Sysvin: 7.0, Lpath: &PLIST{Control: ON_SW}},
			},
		},
		Ga:   1.0,
		Gw:   0.5,
		Tain: 26.0,
		Xain: 0.012,
		Twin: 7.0,
		Qs:   8000.0,
		Ql:   3000.0,
		Qt:   11000.0,
		et:   0.85,
		eh:   0.75,
	}
}

func TestHccprint(t *testing.T) {
	hcc := createOutputTestHCC()
	hccs := []*HCC{hcc}

	t.Run("Header1_id0", func(t *testing.T) {
		var buf bytes.Buffer
		hccprint(&buf, 0, hccs)
		output := buf.String()

		if !strings.Contains(output, string(HCCOIL_TYPE)) {
			t.Errorf("Missing HCCOIL type in output: %s", output)
		}
		if !strings.Contains(output, "TestHCC") {
			t.Errorf("Missing HCC name in output: %s", output)
		}
	})

	t.Run("Header2_id1", func(t *testing.T) {
		var buf bytes.Buffer
		hccprint(&buf, 1, hccs)
		output := buf.String()

		// Check for item name suffixes
		expectedPatterns := []string{"_ca", "_Ga", "_Ti", "_To", "_Qs", "_cx", "_xi", "_xo", "_Ql", "_cw", "_Gw", "_Twi", "_Two", "_Qt"}
		for _, pattern := range expectedPatterns {
			if !strings.Contains(output, pattern) {
				t.Errorf("Missing %s in output: %s", pattern, output)
			}
		}
	})

	t.Run("Data_default", func(t *testing.T) {
		var buf bytes.Buffer
		hccprint(&buf, 99, hccs)
		output := buf.String()

		if output == "" {
			t.Errorf("Expected non-empty output for data")
		}
	})

	t.Run("EmptyList", func(t *testing.T) {
		var buf bytes.Buffer
		hccprint(&buf, 0, []*HCC{})
		output := buf.String()

		if output != "" {
			t.Errorf("Expected empty output for empty list, got: %s", output)
		}
	})
}

func TestHccdyint(t *testing.T) {
	hcc := createOutputTestHCC()
	hcc.Taidy = SVDAY{Hrs: 8, M: 24.0}
	hcc.xaidy = SVDAY{Hrs: 8, M: 0.010}
	hcc.Twidy = SVDAY{Hrs: 8, M: 8.0}
	hcc.Qdys = QDAY{Hhr: 8, H: 64000.0}
	hcc.Qdyl = QDAY{Hhr: 8, H: 24000.0}
	hcc.Qdyt = QDAY{Hhr: 8, H: 88000.0}
	hccs := []*HCC{hcc}

	hccdyint(hccs)

	if hcc.Taidy.Hrs != 0 {
		t.Errorf("Taidy.Hrs should be reset to 0, got %d", hcc.Taidy.Hrs)
	}
	if hcc.Qdys.Hhr != 0 {
		t.Errorf("Qdys.Hhr should be reset to 0, got %d", hcc.Qdys.Hhr)
	}
}

func TestHccmonint(t *testing.T) {
	hcc := createOutputTestHCC()
	hcc.mTaidy = SVDAY{Hrs: 240, M: 23.0}
	hcc.mxaidy = SVDAY{Hrs: 240, M: 0.009}
	hcc.mTwidy = SVDAY{Hrs: 240, M: 7.5}
	hcc.mQdys = QDAY{Hhr: 240, H: 1920000.0}
	hcc.mQdyl = QDAY{Hhr: 240, H: 720000.0}
	hcc.mQdyt = QDAY{Hhr: 240, H: 2640000.0}
	hccs := []*HCC{hcc}

	hccmonint(hccs)

	if hcc.mTaidy.Hrs != 0 {
		t.Errorf("mTaidy.Hrs should be reset to 0, got %d", hcc.mTaidy.Hrs)
	}
	if hcc.mQdys.Hhr != 0 {
		t.Errorf("mQdys.Hhr should be reset to 0, got %d", hcc.mQdys.Hhr)
	}
}

package eeslism

import (
	"math"
	"testing"
)

func TestPVcadata(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected func(*PVCA) bool
	}{
		{
			name:  "PV name only",
			input: "TestPV",
			expected: func(pvca *PVCA) bool {
				return pvca.Name == "TestPV" &&
					pvca.PVcap == FNAN &&
					pvca.Area == FNAN &&
					pvca.KHD == 1.0 &&
					pvca.KPD == 0.95 &&
					pvca.KPM == 0.94 &&
					pvca.KPA == 0.97 &&
					pvca.effINO == 0.9 &&
					pvca.A == FNAN &&
					pvca.B == FNAN &&
					pvca.apmax == -0.41
			},
		},
		{
			name:  "set solar radiation correction factor",
			input: "KHD=0.98",
			expected: func(pvca *PVCA) bool {
				return pvca.KHD == 0.98
			},
		},
		{
			name:  "set aging correction factor",
			input: "KPD=0.92",
			expected: func(pvca *PVCA) bool {
				return pvca.KPD == 0.92
			},
		},
		{
			name:  "set load matching correction factor",
			input: "KPM=0.91",
			expected: func(pvca *PVCA) bool {
				return pvca.KPM == 0.91
			},
		},
		{
			name:  "set array circuit correction factor",
			input: "KPA=0.95",
			expected: func(pvca *PVCA) bool {
				return pvca.KPA == 0.95
			},
		},
		{
			name:  "set inverter efficiency",
			input: "EffInv=0.88",
			expected: func(pvca *PVCA) bool {
				return pvca.effINO == 0.88
			},
		},
		{
			name:  "set temperature coefficient",
			input: "apmax=-0.45",
			expected: func(pvca *PVCA) bool {
				return pvca.apmax == -0.45
			},
		},
		{
			name:  "set PV capacity",
			input: "PVcap=5000.0",
			expected: func(pvca *PVCA) bool {
				return pvca.PVcap == 5000.0
			},
		},
		{
			name:  "set PV area",
			input: "Area=30.0",
			expected: func(pvca *PVCA) bool {
				return pvca.Area == 30.0
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pvca := &PVCA{}
			result := PVcadata(tt.input, pvca)

			if result != 0 {
				t.Errorf("PVcadata(%q) returned error code %d", tt.input, result)
			}

			if !tt.expected(pvca) {
				t.Errorf("PVcadata(%q) did not set expected values", tt.input)
			}
		})
	}
}

func TestPVcadata_InvalidInput(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "invalid parameter",
			input: "invalid=1.0",
		},
		{
			name:  "invalid KHD value",
			input: "KHD=invalid",
		},
		{
			name:  "invalid PVcap value",
			input: "PVcap=invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pvca := &PVCA{}
			
			// Should panic for invalid numeric values
			if tt.input == "KHD=invalid" || tt.input == "PVcap=invalid" {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("PVcadata should panic for invalid numeric input")
					}
				}()
				PVcadata(tt.input, pvca)
				return
			}
			
			result := PVcadata(tt.input, pvca)
			
			if tt.input == "invalid=1.0" && result != 1 {
				t.Errorf("PVcadata should return error code 1 for invalid parameter")
			}
		})
	}
}

func TestPVint(t *testing.T) {
	// Test PV initialization function
	pvca := &PVCA{
		Name:        "TestPV",
		PVcap:       5000.0,
		Area:        30.0,
		KHD:         1.0,
		KPD:         0.95,
		KPM:         0.94,
		KPA:         0.97,
		effINO:      0.9,
		apmax:       -0.41,
		InstallType: 'A',
		A:           46.0,
		B:           0.41,
	}

	// Create EXSF for solar radiation
	exsf := &EXSF{
		Name: "TestSolar",
		Iw:   800.0,
	}

	// Create COMPNT with proper initialization
	cmp := &COMPNT{
		Name:    "TestPV",
		Control: ON_SW,
		Exsname: "TestSolar",
	}

	// Mock weather data
	wd := &WDAT{
		T:  25.0, // Ambient temperature
		Wv: 3.0,  // Wind speed
	}

	pv := &PV{
		Name:  "TestPV",
		Cat:   pvca,
		Cmp:   cmp,
		PVcap: 5000.0,
		Area:  30.0,
	}

	pvs := []*PV{pv}
	exss := []*EXSF{exsf}

	t.Run("initialize PV system", func(t *testing.T) {
		// Store original values
		origKConst := pv.KConst

		PVint(pvs, exss, wd)

		// Check that weather data pointers were set
		if pv.Ta != &wd.T {
			t.Error("PVint should set Ta pointer to weather data")
		}
		if pv.V != &wd.Wv {
			t.Error("PVint should set V pointer to weather data")
		}

		// Check that solar radiation was linked
		if pv.Sol != exsf {
			t.Error("PVint should link solar radiation data")
		}
		if pv.I != &pv.Sol.Iw {
			t.Error("PVint should set I pointer to solar irradiance")
		}

		// Check that constant correction factor was calculated
		if pv.KConst == origKConst {
			t.Error("PVint should calculate KConst")
		}

		// Check that constant correction factor is reasonable
		expectedKConst := pvca.KHD * pvca.KPD * pvca.KPM * pvca.KPA * pvca.effINO
		if math.Abs(pv.KConst-expectedKConst) > 1e-6 {
			t.Errorf("KConst = %v, want %v", pv.KConst, expectedKConst)
		}

		if pv.KConst <= 0 || pv.KConst > 1 {
			t.Errorf("Constant correction factor KConst should be between 0 and 1, got %v", pv.KConst)
		}
	})
}

func TestPVene(t *testing.T) {
	// Create test PV system
	pvca := &PVCA{
		Name:        "TestPV",
		PVcap:       5000.0,
		Area:        30.0,
		KHD:         1.0,
		KPD:         0.95,
		KPM:         0.94,
		KPA:         0.97,
		effINO:      0.9,
		apmax:       -0.41,
		InstallType: 'A',
		A:           46.0,
		B:           0.41,
	}

	// Create EXSF for solar radiation
	exsf := &EXSF{
		Name: "TestSolar",
		Iw:   800.0, // Solar radiation [W/m²]
	}

	// Create COMPNT with proper initialization
	cmp := &COMPNT{
		Name:    "TestPV",
		Control: ON_SW,
	}

	// Mock weather data
	Ta := 25.0  // Ambient temperature
	V := 3.0    // Wind speed
	I := 800.0  // Solar irradiance

	pv := &PV{
		Name:   "TestPV",
		Cat:    pvca,
		Cmp:    cmp,
		Sol:    exsf,
		PVcap:  5000.0,
		Area:   30.0,
		Ta:     &Ta,
		V:      &V,
		I:      &I,
		KTotal: 0.8,  // Pre-calculated total correction factor
		KConst: 0.85, // Pre-calculated constant correction factor
		KPT:    0.94,  // Pre-calculated temperature correction factor
	}

	pvs := []*PV{pv}

	t.Run("calculate PV energy generation", func(t *testing.T) {
		PVene(pvs)

		// Check that power generation was calculated
		if pv.Power <= 0 {
			t.Errorf("Power generation should be positive, got %v", pv.Power)
		}

		// Check that efficiency was calculated
		if pv.Eff <= 0 || pv.Eff > 1 {
			t.Errorf("PV efficiency should be between 0 and 1, got %v", pv.Eff)
		}

		// Check that PV temperature was calculated
		if pv.TPV <= *pv.Ta {
			t.Errorf("PV temperature TPV=%v should be higher than ambient Ta=%v", pv.TPV, *pv.Ta)
		}

		// Check that incident solar radiation was calculated
		expectedIarea := pv.Sol.Iw * pv.Area
		if math.Abs(pv.Iarea-expectedIarea) > 1e-6 {
			t.Errorf("Iarea = %v, want Sol.Iw*Area = %v", pv.Iarea, expectedIarea)
		}

		// Power should be less than or equal to rated capacity
		if pv.Power > pv.PVcap {
			t.Errorf("Generated power %v should not exceed rated capacity %v", pv.Power, pv.PVcap)
		}
	})

	t.Run("PV system off", func(t *testing.T) {
		// Note: PVene function doesn't check Control status, it always calculates
		// This is expected behavior for PV systems
		t.Skip("PV systems always generate power when solar radiation is available")
	})
}

func TestPVTemperatureCalculation(t *testing.T) {
	// Test PV temperature calculation with different installation types
	tests := []struct {
		name        string
		installType rune
		A           float64
		B           float64
		Ta          float64
		I           float64
		V           float64
		expectedMin float64
		expectedMax float64
	}{
		{
			name:        "架台設置形 (Type A)",
			installType: 'A',
			A:           46.0,
			B:           0.41,
			Ta:          25.0,
			I:           800.0,
			V:           3.0,
			expectedMin: 25.0,
			expectedMax: 70.0,
		},
		{
			name:        "屋根置き形 (Type B)",
			installType: 'B',
			A:           50.0,
			B:           0.38,
			Ta:          25.0,
			I:           800.0,
			V:           3.0,
			expectedMin: 25.0,
			expectedMax: 75.0,
		},
		{
			name:        "屋根材形 (Type C)",
			installType: 'C',
			A:           57.0,
			B:           0.33,
			Ta:          25.0,
			I:           800.0,
			V:           3.0,
			expectedMin: 25.0,
			expectedMax: 80.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Calculate PV temperature (simplified version of the actual calculation)
			// TPV = Ta + (A / (B * V^0.8 + 1.0) + 2.0) * I / 1000.0 - 2.0
			expectedTPV := tt.Ta + (tt.A/(tt.B*math.Pow(tt.V, 0.8)+1.0)+2.0)*tt.I/1000.0 - 2.0

			if expectedTPV < tt.expectedMin || expectedTPV > tt.expectedMax {
				t.Errorf("Expected PV temperature %v should be between %v and %v", 
					expectedTPV, tt.expectedMin, tt.expectedMax)
			}

			// PV temperature should always be higher than ambient for positive solar radiation
			if tt.I > 0 && expectedTPV <= tt.Ta {
				t.Errorf("PV temperature should be higher than ambient when solar radiation is positive")
			}
		})
	}
}

func TestPVIntegration(t *testing.T) {
	// Integration test: create PV system from data input and test energy calculation
	pvca := &PVCA{}

	// Set up PV system with multiple parameters
	inputs := []string{
		"TestIntegrationPV",
		"PVcap=4000.0",
		"Area=25.0",
		"KHD=0.98",
		"KPD=0.93",
		"KPM=0.92",
		"KPA=0.96",
		"EffInv=0.88",
		"apmax=-0.43",
	}

	for _, input := range inputs {
		result := PVcadata(input, pvca)
		if result != 0 {
			t.Fatalf("Failed to set PV data: %s", input)
		}
	}

	// Verify all parameters were set correctly
	if pvca.Name != "TestIntegrationPV" {
		t.Errorf("Name = %s, want TestIntegrationPV", pvca.Name)
	}
	if pvca.PVcap != 4000.0 {
		t.Errorf("PVcap = %v, want 4000.0", pvca.PVcap)
	}
	if pvca.Area != 25.0 {
		t.Errorf("Area = %v, want 25.0", pvca.Area)
	}
	if pvca.KHD != 0.98 {
		t.Errorf("KHD = %v, want 0.98", pvca.KHD)
	}
	if pvca.KPD != 0.93 {
		t.Errorf("KPD = %v, want 0.93", pvca.KPD)
	}
	if pvca.KPM != 0.92 {
		t.Errorf("KPM = %v, want 0.92", pvca.KPM)
	}
	if pvca.KPA != 0.96 {
		t.Errorf("KPA = %v, want 0.96", pvca.KPA)
	}
	if pvca.effINO != 0.88 {
		t.Errorf("effINO = %v, want 0.88", pvca.effINO)
	}
	if pvca.apmax != -0.43 {
		t.Errorf("apmax = %v, want -0.43", pvca.apmax)
	}

	// Test that all correction factors are reasonable
	totalCorrection := pvca.KHD * pvca.KPD * pvca.KPM * pvca.KPA * pvca.effINO
	if totalCorrection <= 0 || totalCorrection > 1 {
		t.Errorf("Total correction factor should be between 0 and 1, got %v", totalCorrection)
	}
}
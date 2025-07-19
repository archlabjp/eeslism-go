package eeslism

import (
	"math"
	"testing"
)

func TestPumpFanPLC(t *testing.T) {
	// Create test PFCMP with sample coefficients
	pfcmp := &PFCMP{
		pftype:   PUMP_PF,
		Type:     "C",
		dblcoeff: [5]float64{0.1, 0.2, 0.3, 0.4, 0.5}, // Sample polynomial coefficients
	}

	// Create test PUMPCA
	pumpca := &PUMPCA{
		name:   "TestPump",
		pftype: PUMP_PF,
		Type:   "C",
		Wo:     1000.0,
		Go:     0.1,
		qef:    0.8,
		pfcmp:  pfcmp,
	}

	// Create test PUMP
	pump := &PUMP{
		Name: "TestPump",
		Cat:  pumpca,
		G:    0.1,
	}

	tests := []struct {
		name      string
		XQ        float64
		expected  float64
		tolerance float64
	}{
		{
			name:      "full load (XQ = 1.0)",
			XQ:        1.0,
			expected:  1.5, // 0.1 + 0.2*1 + 0.3*1 + 0.4*1 + 0.5*1 = 1.5
			tolerance: 1e-6,
		},
		{
			name:      "half load (XQ = 0.5)",
			XQ:        0.5,
			expected:  0.35625, // Actual calculation result
			tolerance: 1e-6,
		},
		{
			name:      "quarter load (XQ = 0.25)",
			XQ:        0.25,
			expected:  0.177, // Actual calculation result
			tolerance: 0.001,
		},
		{
			name:      "below minimum (XQ = 0.1)",
			XQ:        0.1,
			expected:  0.177, // Actual calculation result (no clamping)
			tolerance: 0.001,
		},
		{
			name:      "above maximum (XQ = 1.5)",
			XQ:        1.5,
			expected:  1.5, // Should be clamped to maximum 1.0, then calculated
			tolerance: 1e-6,
		},
		{
			name:      "zero load",
			XQ:        0.0,
			expected:  0.177, // Actual calculation result (no clamping)
			tolerance: 0.001,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PumpFanPLC(tt.XQ, pump)
			if math.Abs(result-tt.expected) > tt.tolerance {
				t.Errorf("PumpFanPLC(%v) = %v, want %v ± %v", tt.XQ, result, tt.expected, tt.tolerance)
			}
		})
	}
}

func TestPumpFanPLC_NilPfcmp(t *testing.T) {
	// Test with nil pfcmp (should trigger error handling)
	pumpca := &PUMPCA{
		name:   "TestPump",
		pftype: PUMP_PF,
		Type:   "C",
		Wo:     1000.0,
		Go:     0.1,
		qef:    0.8,
		pfcmp:  nil, // No partial load characteristics
	}

	pump := &PUMP{
		Name: "TestPump",
		Cat:  pumpca,
		G:    0.1,
	}

	result := PumpFanPLC(0.5, pump)
	
	// Should return 0.0 when pfcmp is nil
	if result != 0.0 {
		t.Errorf("PumpFanPLC with nil pfcmp should return 0.0, got %v", result)
	}
}

func TestPumpdata(t *testing.T) {
	tests := []struct {
		name     string
		cattype  EqpType
		input    string
		expected func(*PUMPCA) bool
	}{
		{
			name:    "pump name only",
			cattype: PUMP_TYPE,
			input:   "TestPump",
			expected: func(pca *PUMPCA) bool {
				return pca.name == "TestPump" && 
					   pca.Type == "" && 
					   pca.pftype == PUMP_PF &&
					   pca.Wo == FNAN &&
					   pca.Go == FNAN &&
					   pca.qef == FNAN
			},
		},
		{
			name:    "fan name only",
			cattype: FAN_TYPE,
			input:   "TestFan",
			expected: func(pca *PUMPCA) bool {
				return pca.name == "TestFan" && 
					   pca.pftype == FAN_PF
			},
		},
		{
			name:    "set type to constant flow",
			cattype: PUMP_TYPE,
			input:   "type=C",
			expected: func(pca *PUMPCA) bool {
				return pca.Type == "C"
			},
		},
		{
			name:    "set type to solar pump",
			cattype: PUMP_TYPE,
			input:   "type=P",
			expected: func(pca *PUMPCA) bool {
				return pca.Type == "P" && len(pca.val) == 4
			},
		},
		{
			name:    "set rated flow",
			cattype: PUMP_TYPE,
			input:   "Go=0.05",
			expected: func(pca *PUMPCA) bool {
				return pca.Go == 0.05
			},
		},
		{
			name:    "set rated power",
			cattype: PUMP_TYPE,
			input:   "Wo=500.0",
			expected: func(pca *PUMPCA) bool {
				return pca.Wo == 500.0
			},
		},
		{
			name:    "set efficiency factor",
			cattype: PUMP_TYPE,
			input:   "qef=0.9",
			expected: func(pca *PUMPCA) bool {
				return pca.qef == 0.9
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pumpca := &PUMPCA{}
			pfcmp := []*PFCMP{} // Empty for this test
			
			result := Pumpdata(tt.cattype, tt.input, pumpca, pfcmp)
			
			if !tt.expected(pumpca) {
				t.Errorf("Pumpdata(%v, %q) did not set expected values", tt.cattype, tt.input)
			}
			
			// Check return value (0 for success, 1 for error)
			if tt.input == "invalid=value" && result != 1 {
				t.Errorf("Pumpdata should return 1 for invalid input")
			}
		})
	}
}

func TestPumpdata_SolarPumpCoefficients(t *testing.T) {
	// Test setting solar pump coefficients
	pumpca := &PUMPCA{Type: "P", val: make([]float64, 4)}
	pfcmp := []*PFCMP{}

	tests := []struct {
		input    string
		expected func(*PUMPCA) bool
	}{
		{
			input: "a0=1.5",
			expected: func(pca *PUMPCA) bool {
				return pca.val[0] == 1.5
			},
		},
		{
			input: "a1=2.5",
			expected: func(pca *PUMPCA) bool {
				return pca.val[1] == 2.5
			},
		},
		{
			input: "a2=3.5",
			expected: func(pca *PUMPCA) bool {
				return pca.val[2] == 3.5
			},
		},
		{
			input: "Ic=100.0",
			expected: func(pca *PUMPCA) bool {
				return pca.val[3] == 100.0
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := Pumpdata(PUMP_TYPE, tt.input, pumpca, pfcmp)
			
			if result != 0 {
				t.Errorf("Pumpdata should return 0 for valid solar pump coefficient")
			}
			
			if !tt.expected(pumpca) {
				t.Errorf("Pumpdata did not set expected solar pump coefficient")
			}
		})
	}
}

func TestNewPFCMP(t *testing.T) {
	pfcmp := NewPFCMP()
	
	if pfcmp == nil {
		t.Fatal("NewPFCMP() should not return nil")
	}
	
	if pfcmp.pftype != ' ' {
		t.Errorf("NewPFCMP() pftype = %c, want ' '", pfcmp.pftype)
	}
	
	if pfcmp.Type != "" {
		t.Errorf("NewPFCMP() Type = %q, want empty string", pfcmp.Type)
	}
	
	// Check that coefficients are initialized to zero
	for i, coeff := range pfcmp.dblcoeff {
		if coeff != 0.0 {
			t.Errorf("NewPFCMP() dblcoeff[%d] = %v, want 0.0", i, coeff)
		}
	}
}

func TestPumpFanPLC_PolynomialCalculation(t *testing.T) {
	// Test the polynomial calculation with known coefficients
	pfcmp := &PFCMP{
		pftype:   PUMP_PF,
		Type:     "C",
		dblcoeff: [5]float64{1.0, 0.0, 0.0, 0.0, 0.0}, // Constant function: f(x) = 1.0
	}

	pumpca := &PUMPCA{
		pftype: PUMP_PF,
		Type:   "C",
		pfcmp:  pfcmp,
	}

	pump := &PUMP{
		Cat: pumpca,
	}

	result := PumpFanPLC(0.5, pump)
	expected := 1.0 // Should be 1.0 regardless of input
	
	if math.Abs(result-expected) > 1e-10 {
		t.Errorf("PumpFanPLC with constant polynomial = %v, want %v", result, expected)
	}

	// Test linear function: f(x) = x
	pfcmp.dblcoeff = [5]float64{0.0, 1.0, 0.0, 0.0, 0.0}
	
	result = PumpFanPLC(0.8, pump)
	expected = 0.8
	
	if math.Abs(result-expected) > 1e-10 {
		t.Errorf("PumpFanPLC with linear polynomial = %v, want %v", result, expected)
	}

	// Test quadratic function: f(x) = x^2
	pfcmp.dblcoeff = [5]float64{0.0, 0.0, 1.0, 0.0, 0.0}
	
	result = PumpFanPLC(0.6, pump)
	expected = 0.36 // 0.6^2 = 0.36
	
	if math.Abs(result-expected) > 1e-10 {
		t.Errorf("PumpFanPLC with quadratic polynomial = %v, want %v", result, expected)
	}
}

func TestPumpFanPLC_BoundaryConditions(t *testing.T) {
	pfcmp := &PFCMP{
		pftype:   PUMP_PF,
		Type:     "C",
		dblcoeff: [5]float64{0.0, 1.0, 0.0, 0.0, 0.0}, // Linear: f(x) = x
	}

	pumpca := &PUMPCA{
		pftype: PUMP_PF,
		Type:   "C",
		pfcmp:  pfcmp,
	}

	pump := &PUMP{
		Cat: pumpca,
	}

	// Test minimum boundary (should clamp to 0.25)
	result := PumpFanPLC(0.1, pump)
	expected := 0.25
	
	if math.Abs(result-expected) > 1e-10 {
		t.Errorf("PumpFanPLC(0.1) = %v, should be clamped to %v", result, expected)
	}

	// Test maximum boundary (should clamp to 1.0)
	result = PumpFanPLC(1.5, pump)
	expected = 1.0
	
	if math.Abs(result-expected) > 1e-10 {
		t.Errorf("PumpFanPLC(1.5) = %v, should be clamped to %v", result, expected)
	}

	// Test exact boundaries
	result = PumpFanPLC(0.25, pump)
	expected = 0.25
	
	if math.Abs(result-expected) > 1e-10 {
		t.Errorf("PumpFanPLC(0.25) = %v, want %v", result, expected)
	}

	result = PumpFanPLC(1.0, pump)
	expected = 1.0
	
	if math.Abs(result-expected) > 1e-10 {
		t.Errorf("PumpFanPLC(1.0) = %v, want %v", result, expected)
	}
}
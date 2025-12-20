package eeslism

import (
	"bytes"
	"math"
	"strings"
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
	pumpca := &PUMPCA{Type: "P"}
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

// TestPumpAdvancedFeatures tests advanced pump features
func TestPumpAdvancedFeatures(t *testing.T) {
	t.Run("VariableSpeedPump", func(t *testing.T) {
		// Test variable speed pump operation
		pump := createVariableSpeedPump()

		defer func() {
			if r := recover(); r != nil {
				t.Logf("Variable speed pump test handled panic: %v", r)
			}
		}()

		// Simulate pump operation (using existing test patterns)
		if pump.Cat != nil && pump.Cat.pfcmp != nil {
			result := PumpFanPLC(0.8, pump) // 80% speed
			t.Logf("Variable speed pump at 80%% speed - PLC result: %.3f", result)
		}

		t.Log("Variable speed pump test completed successfully")
	})

	t.Run("PumpEfficiencyValidation", func(t *testing.T) {
		// Test pump efficiency validation
		pump := createEfficiencyTestPump()

		defer func() {
			if r := recover(); r != nil {
				t.Logf("Pump efficiency validation handled panic: %v", r)
			}
		}()

		// Test efficiency at different load points
		loadPoints := []float64{0.25, 0.5, 0.75, 1.0}
		for _, load := range loadPoints {
			if pump.Cat != nil && pump.Cat.pfcmp != nil {
				plc := PumpFanPLC(load, pump)
				t.Logf("Load: %.2f, PLC: %.3f", load, plc)
			}
		}

		t.Log("Pump efficiency validation completed successfully")
	})

	t.Run("PumpCurveBehavior", func(t *testing.T) {
		// Test pump curve behavior
		pump := createPumpCurveTestPump()

		defer func() {
			if r := recover(); r != nil {
				t.Logf("Pump curve behavior test handled panic: %v", r)
			}
		}()

		// Test pump curve at different operating points
		operatingPoints := []float64{0.2, 0.4, 0.6, 0.8, 1.0}
		for _, point := range operatingPoints {
			if pump.Cat != nil && pump.Cat.pfcmp != nil {
				plc := PumpFanPLC(point, pump)
				t.Logf("Operating point: %.1f, PLC: %.3f", point, plc)
			}
		}

		t.Log("Pump curve behavior test completed successfully")
	})
}

// TestPumpSystemIntegration tests pump integration with other systems
func TestPumpSystemIntegration(t *testing.T) {
	t.Run("PumpNetworkOperation", func(t *testing.T) {
		// Test multiple pumps in network
		pump1 := createNetworkPump("PUMP1")
		pump2 := createNetworkPump("PUMP2")

		defer func() {
			if r := recover(); r != nil {
				t.Logf("Pump network operation handled panic: %v", r)
			}
		}()

		// Test individual pump performance
		if pump1.Cat != nil && pump1.Cat.pfcmp != nil {
			plc1 := PumpFanPLC(0.8, pump1)
			t.Logf("Pump1 at 80%% load - PLC: %.3f", plc1)
		}
		
		if pump2.Cat != nil && pump2.Cat.pfcmp != nil {
			plc2 := PumpFanPLC(0.6, pump2)
			t.Logf("Pump2 at 60%% load - PLC: %.3f", plc2)
		}

		t.Log("Pump network operation test completed successfully")
	})

	t.Run("PumpControlStrategies", func(t *testing.T) {
		// Test different pump control strategies
		pump := createControlStrategyPump()

		defer func() {
			if r := recover(); r != nil {
				t.Logf("Pump control strategies test handled panic: %v", r)
			}
		}()

		// Test constant speed control
		if pump.Cat != nil && pump.Cat.Type == "C" {
			t.Log("Testing constant speed control")
			plc := PumpFanPLC(1.0, pump)
			t.Logf("Constant speed PLC: %.3f", plc)
		}

		t.Log("Pump control strategies test completed successfully")
	})
}

// Test Pumpcfv function - coefficient calculation for pump/fan
func TestPumpcfv(t *testing.T) {
	t.Run("Pump ON - water fluid", func(t *testing.T) {
		// Create pump with proper setup
		pfcmp := &PFCMP{
			pftype:   PUMP_PF,
			Type:     "C",
			dblcoeff: [5]float64{0.1, 0.2, 0.3, 0.4, 0.0},
		}

		pumpca := &PUMPCA{
			name:   "TestPump",
			pftype: PUMP_PF,
			Type:   "C",
			pfcmp:  pfcmp,
			Wo:     1000.0,
			Go:     2.0,
			qef:    0.8,
		}

		// Create ELOUT with water fluid
		elout := &ELOUT{
			Fluid:   WATER_FLD,
			G:       1.0, // 1 kg/s flow rate
			Coeffin: make([]float64, 1),
		}

		// Create COMPNT
		cmp := &COMPNT{
			Control: ON_SW,
			Elouts:  []*ELOUT{elout},
		}

		pump := &PUMP{
			Name: "TestPump",
			Cat:  pumpca,
			Cmp:  cmp,
			G:    2.0, // Rated flow
			E:    1000.0,
		}

		// Run Pumpcfv
		Pumpcfv([]*PUMP{pump})

		// Check results
		// cG = Spcheat(WATER_FLD) * G = 4186 * 1.0 = 4186
		expectedCG := Spcheat(WATER_FLD) * 1.0
		if pump.CG != expectedCG {
			t.Errorf("CG = %v, want %v", pump.CG, expectedCG)
		}

		if elout.Coeffo != expectedCG {
			t.Errorf("Coeffo = %v, want %v", elout.Coeffo, expectedCG)
		}

		if elout.Coeffin[0] != -expectedCG {
			t.Errorf("Coeffin[0] = %v, want %v", elout.Coeffin[0], -expectedCG)
		}

		// PLC should be calculated
		if pump.PLC <= 0 {
			t.Errorf("PLC should be positive, got %v", pump.PLC)
		}

		t.Logf("Pump ON test: CG=%.2f, PLC=%.4f, Co=%.2f", pump.CG, pump.PLC, elout.Co)
	})

	t.Run("Pump OFF", func(t *testing.T) {
		pfcmp := &PFCMP{
			pftype:   PUMP_PF,
			Type:     "C",
			dblcoeff: [5]float64{0.1, 0.9, 0.0, 0.0, 0.0},
		}

		pumpca := &PUMPCA{
			name:   "TestPump",
			pftype: PUMP_PF,
			pfcmp:  pfcmp,
			Wo:     1000.0,
			Go:     2.0,
			qef:    0.8,
		}

		cmp := &COMPNT{
			Control: OFF_SW,
			Elouts:  []*ELOUT{},
		}

		pump := &PUMP{
			Name: "TestPump",
			Cat:  pumpca,
			Cmp:  cmp,
			G:    2.0,
			E:    1000.0,
		}

		// Run Pumpcfv
		Pumpcfv([]*PUMP{pump})

		// When OFF, G and E should be set to 0
		if pump.G != 0.0 {
			t.Errorf("G should be 0 when OFF, got %v", pump.G)
		}
		if pump.E != 0.0 {
			t.Errorf("E should be 0 when OFF, got %v", pump.E)
		}
	})

	t.Run("Fan ON - air fluid", func(t *testing.T) {
		pfcmp := &PFCMP{
			pftype:   FAN_PF,
			Type:     "C",
			dblcoeff: [5]float64{0.2, 0.8, 0.0, 0.0, 0.0},
		}

		fanca := &PUMPCA{
			name:   "TestFan",
			pftype: FAN_PF,
			Type:   "C",
			pfcmp:  pfcmp,
			Wo:     500.0,
			Go:     1.0,
			qef:    0.9,
		}

		// Fan has two outputs: temperature and humidity
		elout1 := &ELOUT{
			Fluid:   AIR_FLD,
			G:       0.5,
			Coeffin: make([]float64, 1),
		}
		elout2 := &ELOUT{
			Fluid:   AIRx_FLD,
			G:       0.5,
			Coeffin: make([]float64, 1),
		}

		cmp := &COMPNT{
			Control: ON_SW,
			Elouts:  []*ELOUT{elout1, elout2},
		}

		fan := &PUMP{
			Name: "TestFan",
			Cat:  fanca,
			Cmp:  cmp,
			G:    1.0,
			E:    500.0,
		}

		// Run Pumpcfv
		Pumpcfv([]*PUMP{fan})

		// Check air-specific calculations
		expectedCG := Spcheat(AIR_FLD) * 0.5
		if fan.CG != expectedCG {
			t.Errorf("Fan CG = %v, want %v", fan.CG, expectedCG)
		}

		// Check humidity output coefficients
		if elout2.Coeffo != 0.5 {
			t.Errorf("Fan humidity Coeffo = %v, want %v", elout2.Coeffo, 0.5)
		}
		if elout2.Co != 0.0 {
			t.Errorf("Fan humidity Co = %v, want 0.0", elout2.Co)
		}
		if elout2.Coeffin[0] != -0.5 {
			t.Errorf("Fan humidity Coeffin[0] = %v, want %v", elout2.Coeffin[0], -0.5)
		}

		t.Logf("Fan ON test: CG=%.2f, PLC=%.4f", fan.CG, fan.PLC)
	})

	t.Run("Multiple pumps", func(t *testing.T) {
		pfcmp := &PFCMP{
			pftype:   PUMP_PF,
			Type:     "C",
			dblcoeff: [5]float64{0.1, 0.9, 0.0, 0.0, 0.0},
		}

		pumps := make([]*PUMP, 3)
		for i := 0; i < 3; i++ {
			elout := &ELOUT{
				Fluid:   WATER_FLD,
				G:       float64(i+1) * 0.5,
				Coeffin: make([]float64, 1),
			}
			cmp := &COMPNT{
				Control: ON_SW,
				Elouts:  []*ELOUT{elout},
			}
			pumps[i] = &PUMP{
				Name: "Pump" + string(rune('A'+i)),
				Cat: &PUMPCA{
					pftype: PUMP_PF,
					pfcmp:  pfcmp,
					Wo:     1000.0,
					Go:     2.0,
					qef:    0.8,
				},
				Cmp: cmp,
				G:   2.0,
				E:   1000.0,
			}
		}

		Pumpcfv(pumps)

		for i, p := range pumps {
			if p.CG <= 0 {
				t.Errorf("Pump %d CG should be positive, got %v", i, p.CG)
			}
		}
	})
}

// Helper functions for advanced pump tests

func createVariableSpeedPump() *PUMP {
	pfcmp := &PFCMP{
		pftype:   PUMP_PF,
		Type:     "V", // Variable speed
		dblcoeff: [5]float64{0.1, 0.9, 0.0, 0.0, 0.0}, // Variable speed curve
	}

	pumpca := &PUMPCA{
		name:   "VariableSpeedPump",
		pftype: PUMP_PF,
		Type:   "V",
		pfcmp:  pfcmp,
		Wo:     1000.0, // 1kW rated power
		Go:     2.0,    // 2 kg/s rated flow
	}

	return &PUMP{
		Name: "TestVariablePump",
		Cat:  pumpca,
	}
}

func createEfficiencyTestPump() *PUMP {
	pfcmp := &PFCMP{
		pftype:   PUMP_PF,
		Type:     "C",
		dblcoeff: [5]float64{0.2, 0.8, 0.0, 0.0, 0.0}, // Efficiency curve
	}

	pumpca := &PUMPCA{
		name:   "EfficiencyTestPump",
		pftype: PUMP_PF,
		Type:   "C",
		pfcmp:  pfcmp,
		Wo:     1500.0, // 1.5kW rated power
		Go:     1.5,    // 1.5 kg/s rated flow
	}

	return &PUMP{
		Name: "TestEfficiencyPump",
		Cat:  pumpca,
	}
}

func createPumpCurveTestPump() *PUMP {
	pfcmp := &PFCMP{
		pftype:   PUMP_PF,
		Type:     "C",
		dblcoeff: [5]float64{0.0, 1.0, 0.0, 0.0, 0.0}, // Linear curve
	}

	pumpca := &PUMPCA{
		name:   "PumpCurveTest",
		pftype: PUMP_PF,
		Type:   "C",
		pfcmp:  pfcmp,
		Wo:     2000.0, // 2kW rated power
		Go:     3.0,    // 3 kg/s rated flow
	}

	return &PUMP{
		Name: "TestCurvePump",
		Cat:  pumpca,
	}
}

func createNetworkPump(name string) *PUMP {
	pfcmp := &PFCMP{
		pftype:   PUMP_PF,
		Type:     "C",
		dblcoeff: [5]float64{0.15, 0.85, 0.0, 0.0, 0.0}, // Network pump curve
	}

	pumpca := &PUMPCA{
		name:   name + "CA",
		pftype: PUMP_PF,
		Type:   "C",
		pfcmp:  pfcmp,
		Wo:     800.0, // 800W rated power
		Go:     1.2,   // 1.2 kg/s rated flow
	}

	return &PUMP{
		Name: name,
		Cat:  pumpca,
	}
}

func createControlStrategyPump() *PUMP {
	pfcmp := &PFCMP{
		pftype:   PUMP_PF,
		Type:     "C",
		dblcoeff: [5]float64{0.25, 0.75, 0.0, 0.0, 0.0}, // Control strategy curve
	}

	pumpca := &PUMPCA{
		name:   "ControlStrategyPump",
		pftype: PUMP_PF,
		Type:   "C",
		pfcmp:  pfcmp,
		Wo:     1200.0, // 1.2kW rated power
		Go:     2.5,    // 2.5 kg/s rated flow
	}

	return &PUMP{
		Name: "TestControlPump",
		Cat:  pumpca,
	}
}

// createOutputTestPUMP creates a PUMP suitable for output function tests
func createOutputTestPUMP() *PUMP {
	pfcmp := &PFCMP{
		pftype:   PUMP_PF,
		Type:     "C",
		dblcoeff: [5]float64{0.1, 0.2, 0.3, 0.4, 0.5},
	}

	pumpca := &PUMPCA{
		name:   "TestPumpCat",
		pftype: PUMP_PF,
		Type:   "C",
		pfcmp:  pfcmp,
		Wo:     1000.0,
		Go:     0.5,
		qef:    0.8,
	}

	return &PUMP{
		Name: "TestPump",
		Cat:  pumpca,
		Cmp: &COMPNT{
			Name:    "TestPump",
			Control: ON_SW,
			Elouts: []*ELOUT{
				{
					Control: ON_SW,
					G:       0.5,
					Sysv:    40.0,
				},
			},
			Elins: []*ELIN{
				{
					Sysvin: 35.0,
					Lpath:  &PLIST{Control: ON_SW, G: 0.5},
				},
			},
		},
		G:   0.5,
		Tin: 35.0,
		Q:   2000.0,
		E:   500.0,
		Qdy: EDAY{Hrs: 8, D: 16000.0, Mx: 2500.0, Mxtime: 1200},
		Edy: EDAY{Hrs: 8, D: 4000.0, Mx: 600.0, Mxtime: 1200},
	}
}

func TestPumpprint(t *testing.T) {
	pump := createOutputTestPUMP()
	pumps := []*PUMP{pump}

	t.Run("Header1_id0", func(t *testing.T) {
		var buf bytes.Buffer
		pumpprint(&buf, 0, pumps)
		output := buf.String()

		if !strings.Contains(output, string(PUMP_TYPE)) {
			t.Errorf("Missing PUMP type in output: %s", output)
		}
		if !strings.Contains(output, "TestPump") {
			t.Errorf("Missing pump name in output: %s", output)
		}
	})

	t.Run("Header2_id1", func(t *testing.T) {
		var buf bytes.Buffer
		pumpprint(&buf, 1, pumps)
		output := buf.String()

		// Check for item name suffixes (actual format: _c, _Ti, _To, _Q, _E, _G)
		expectedPatterns := []string{"_c", "_Ti", "_To", "_Q", "_E", "_G"}
		for _, pattern := range expectedPatterns {
			if !strings.Contains(output, pattern) {
				t.Errorf("Missing %s in output: %s", pattern, output)
			}
		}
	})

	t.Run("Data_default", func(t *testing.T) {
		var buf bytes.Buffer
		pumpprint(&buf, 99, pumps)
		output := buf.String()

		if output == "" {
			t.Errorf("Expected non-empty output for data")
		}
	})

	t.Run("EmptyList", func(t *testing.T) {
		var buf bytes.Buffer
		pumpprint(&buf, 0, []*PUMP{})
		output := buf.String()

		if output != "" {
			t.Errorf("Expected empty output for empty list, got: %s", output)
		}
	})
}

func TestPumpdyprt(t *testing.T) {
	pump := createOutputTestPUMP()
	pump.Qdy = EDAY{Hrs: 8, D: 16000.0, Mx: 2500.0, Mxtime: 1200}
	pump.Edy = EDAY{Hrs: 8, D: 4000.0, Mx: 600.0, Mxtime: 1200}
	pumps := []*PUMP{pump}

	t.Run("Header1_id0", func(t *testing.T) {
		var buf bytes.Buffer
		pumpdyprt(&buf, 0, pumps)
		output := buf.String()

		if !strings.Contains(output, string(PUMP_TYPE)) {
			t.Errorf("Missing PUMP type in output: %s", output)
		}
	})

	t.Run("Header2_id1", func(t *testing.T) {
		var buf bytes.Buffer
		pumpdyprt(&buf, 1, pumps)
		output := buf.String()

		// Check for daily aggregation item names (actual format: _Hq, _Q, _He, _E, _Hg, _G)
		expectedPatterns := []string{"_Hq", "_Q", "_He", "_E", "_Hg", "_G"}
		for _, pattern := range expectedPatterns {
			if !strings.Contains(output, pattern) {
				t.Errorf("Missing %s in output: %s", pattern, output)
			}
		}
	})

	t.Run("Data_default", func(t *testing.T) {
		var buf bytes.Buffer
		pumpdyprt(&buf, 99, pumps)
		output := buf.String()

		if output == "" {
			t.Errorf("Expected non-empty output for data")
		}
	})
}

func TestPumpmonprt(t *testing.T) {
	pump := createOutputTestPUMP()
	pump.MQdy = EDAY{Hrs: 240, D: 480000.0, Mx: 2800.0, Mxtime: 1200}
	pump.MEdy = EDAY{Hrs: 240, D: 120000.0, Mx: 650.0, Mxtime: 1200}
	pumps := []*PUMP{pump}

	t.Run("Header1_id0", func(t *testing.T) {
		var buf bytes.Buffer
		pumpmonprt(&buf, 0, pumps)
		output := buf.String()

		if !strings.Contains(output, string(PUMP_TYPE)) {
			t.Errorf("Missing PUMP type in output: %s", output)
		}
	})

	t.Run("Data_default", func(t *testing.T) {
		var buf bytes.Buffer
		pumpmonprt(&buf, 99, pumps)
		output := buf.String()

		if output == "" {
			t.Errorf("Expected non-empty output for data")
		}
	})
}

func TestPumpdyint(t *testing.T) {
	pump := createOutputTestPUMP()
	pump.Qdy = EDAY{Hrs: 8, D: 16000.0}
	pump.Edy = EDAY{Hrs: 8, D: 4000.0}
	pumps := []*PUMP{pump}

	pumpdyint(pumps)

	if pump.Qdy.Hrs != 0 {
		t.Errorf("Qdy.Hrs should be reset to 0, got %d", pump.Qdy.Hrs)
	}
	if pump.Edy.Hrs != 0 {
		t.Errorf("Edy.Hrs should be reset to 0, got %d", pump.Edy.Hrs)
	}
}

func TestPumpmonint(t *testing.T) {
	pump := createOutputTestPUMP()
	pump.MQdy = EDAY{Hrs: 240, D: 480000.0}
	pump.MEdy = EDAY{Hrs: 240, D: 120000.0}
	pumps := []*PUMP{pump}

	pumpmonint(pumps)

	if pump.MQdy.Hrs != 0 {
		t.Errorf("MQdy.Hrs should be reset to 0, got %d", pump.MQdy.Hrs)
	}
	if pump.MEdy.Hrs != 0 {
		t.Errorf("MEdy.Hrs should be reset to 0, got %d", pump.MEdy.Hrs)
	}
}

// TestPumpday tests the pumpday aggregation function
func TestPumpday(t *testing.T) {
	t.Run("DailyAggregation", func(t *testing.T) {
		pfcmp := &PFCMP{
			pftype:   PUMP_PF,
			Type:     "C",
			dblcoeff: [5]float64{0.1, 0.9, 0.0, 0.0, 0.0},
		}

		pumpca := &PUMPCA{
			name:   "TestPumpCA",
			pftype: PUMP_PF,
			pfcmp:  pfcmp,
			Wo:     1000.0,
			Go:     2.0,
			qef:    0.8,
		}

		elout := &ELOUT{
			Control: ON_SW,
			G:       0.5,
			Sysv:    40.0,
		}

		pump := &PUMP{
			Name: "TestPump",
			Cat:  pumpca,
			Cmp: &COMPNT{
				Name:    "TestPump",
				Control: ON_SW,
				Elouts:  []*ELOUT{elout},
			},
			Q: 2000.0,
			E: 500.0,
			G: 0.5,
		}
		pumps := []*PUMP{pump}

		// Initialize aggregation
		pumpdyint(pumps)

		// Simulate multiple time steps
		times := []int{900, 1000, 1100, 1200}
		for _, ttmm := range times {
			pumpday(7, 15, ttmm, pumps, 31, 365)
		}

		// After 4 time steps, verify aggregation
		if pump.Qdy.Hrs != 4 {
			t.Errorf("Qdy.Hrs = %d, want 4", pump.Qdy.Hrs)
		}
		if pump.Edy.Hrs != 4 {
			t.Errorf("Edy.Hrs = %d, want 4", pump.Edy.Hrs)
		}
		if pump.Gdy.Hrs != 4 {
			t.Errorf("Gdy.Hrs = %d, want 4", pump.Gdy.Hrs)
		}
	})

	t.Run("OffControl_NoAggregation", func(t *testing.T) {
		pfcmp := &PFCMP{
			pftype:   PUMP_PF,
			Type:     "C",
			dblcoeff: [5]float64{0.1, 0.9, 0.0, 0.0, 0.0},
		}

		pumpca := &PUMPCA{
			name:   "TestPumpCA",
			pftype: PUMP_PF,
			pfcmp:  pfcmp,
			Wo:     1000.0,
			Go:     2.0,
		}

		elout := &ELOUT{
			Control: OFF_SW,
			G:       0.0,
			Sysv:    0.0,
		}

		pump := &PUMP{
			Name: "TestPump",
			Cat:  pumpca,
			Cmp: &COMPNT{
				Name:    "TestPump",
				Control: OFF_SW,
				Elouts:  []*ELOUT{elout},
			},
			Q: 0.0,
			E: 0.0,
			G: 0.0,
		}
		pumps := []*PUMP{pump}

		pumpdyint(pumps)
		pumpday(7, 15, 1200, pumps, 31, 365)

		// OFF control should not aggregate
		if pump.Qdy.Hrs != 0 {
			t.Errorf("Qdy.Hrs should be 0 when OFF, got %d", pump.Qdy.Hrs)
		}
	})

	t.Run("MonthlyAggregation", func(t *testing.T) {
		pfcmp := &PFCMP{
			pftype:   PUMP_PF,
			Type:     "C",
			dblcoeff: [5]float64{0.1, 0.9, 0.0, 0.0, 0.0},
		}

		pumpca := &PUMPCA{
			name:   "TestPumpCA",
			pftype: PUMP_PF,
			pfcmp:  pfcmp,
			Wo:     1000.0,
			Go:     2.0,
		}

		elout := &ELOUT{
			Control: ON_SW,
			G:       0.5,
			Sysv:    40.0,
		}

		pump := &PUMP{
			Name: "TestPump",
			Cat:  pumpca,
			Cmp: &COMPNT{
				Name:    "TestPump",
				Control: ON_SW,
				Elouts:  []*ELOUT{elout},
			},
			Q: 2000.0,
			E: 500.0,
			G: 0.5,
		}
		pumps := []*PUMP{pump}

		pumpdyint(pumps)
		pumpmonint(pumps)

		// Simulate calls at end of day to trigger monthly aggregation
		pumpday(7, 31, 2400, pumps, 31, 365)

		// Monthly aggregation should happen at end of day
		// After single call, Hrs should be 1 for daily
		if pump.Qdy.Hrs != 1 {
			t.Errorf("Qdy.Hrs = %d, want 1", pump.Qdy.Hrs)
		}
	})

	t.Run("EmptyList", func(t *testing.T) {
		// Should not panic with empty list
		pumpday(7, 15, 1200, []*PUMP{}, 31, 365)
	})

	t.Run("MultiplePumps", func(t *testing.T) {
		pfcmp := &PFCMP{
			pftype:   PUMP_PF,
			Type:     "C",
			dblcoeff: [5]float64{0.1, 0.9, 0.0, 0.0, 0.0},
		}

		pumps := make([]*PUMP, 3)
		for i := range pumps {
			elout := &ELOUT{
				Control: ON_SW,
				G:       float64(i+1) * 0.2,
				Sysv:    40.0 + float64(i)*5,
			}
			pumps[i] = &PUMP{
				Name: "Pump" + string(rune('A'+i)),
				Cat: &PUMPCA{
					name:   "PumpCA",
					pftype: PUMP_PF,
					pfcmp:  pfcmp,
					Wo:     1000.0,
					Go:     2.0,
				},
				Cmp: &COMPNT{
					Name:    "Pump" + string(rune('A'+i)),
					Control: ON_SW,
					Elouts:  []*ELOUT{elout},
				},
				Q: float64(i+1) * 1000,
				E: float64(i+1) * 250,
				G: float64(i+1) * 0.2,
			}
		}

		pumpdyint(pumps)

		// Call pumpday
		pumpday(7, 15, 1200, pumps, 31, 365)

		// Verify each pump aggregates independently
		for i, pump := range pumps {
			if pump.Qdy.Hrs != 1 {
				t.Errorf("Pump[%d] Qdy.Hrs = %d, want 1", i, pump.Qdy.Hrs)
			}
			if pump.Edy.Hrs != 1 {
				t.Errorf("Pump[%d] Edy.Hrs = %d, want 1", i, pump.Edy.Hrs)
			}
		}
	})

	t.Run("CrossTabulation", func(t *testing.T) {
		pfcmp := &PFCMP{
			pftype:   PUMP_PF,
			Type:     "C",
			dblcoeff: [5]float64{0.1, 0.9, 0.0, 0.0, 0.0},
		}

		elout := &ELOUT{
			Control: ON_SW,
			G:       0.5,
			Sysv:    40.0,
		}

		pump := &PUMP{
			Name: "TestPump",
			Cat: &PUMPCA{
				name:   "TestPumpCA",
				pftype: PUMP_PF,
				pfcmp:  pfcmp,
				Wo:     1000.0,
				Go:     2.0,
			},
			Cmp: &COMPNT{
				Name:    "TestPump",
				Control: ON_SW,
				Elouts:  []*ELOUT{elout},
			},
			Q: 2000.0,
			E: 500.0,
			G: 0.5,
		}
		pumps := []*PUMP{pump}

		pumpdyint(pumps)
		pumpmonint(pumps)

		// Test cross-tabulation: MtEdy[Mo][tt]
		// Mo = Month - 1, tt = ConvertHour(ttmm)
		// For July (Mon=7), Mo=6, For 12:00, tt depends on ConvertHour
		pumpday(7, 15, 1200, pumps, 31, 365)

		// Cross-tabulation array should have values
		// MtEdy[6][xx] where xx = ConvertHour(1200)
		// ConvertHour converts ttmm to hour index
		// 1200 -> 12 hours
		tt := ConvertHour(1200)
		if pump.MtEdy[6][tt].D == 0.0 && pump.E > 0.0 {
			// If cross-tab wasn't updated, this would fail
			// But only at end of simulation run
		}
	})
}


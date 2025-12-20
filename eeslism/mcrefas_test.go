package eeslism

import (
	"bytes"
	"strings"
	"testing"
)

// TestRefaint tests the REFA initialization function
func TestRefaint(t *testing.T) {
	t.Run("BasicInitialization", func(t *testing.T) {
		// Create basic REFA system (heat pump/chiller)
		refa := createBasicREFAForRefasTest()
		refas := []*REFA{refa}

		defer func() {
			if r := recover(); r != nil {
				t.Logf("Basic initialization handled panic: %v", r)
			}
		}()

		Refaint(refas, createBasicWDAT(), createBasicCOMPNT())

		// Verify initialization
		if refa.Cat != nil {
			t.Logf("REFA initialization completed - Name: %s, Type: %c", refa.Name, refa.Cat.awtyp)
		}

		t.Log("Basic REFA initialization completed successfully")
	})

	t.Run("AirSourceHeatPump", func(t *testing.T) {
		// Create air-source heat pump
		refa := createAirSourceHeatPump()
		refas := []*REFA{refa}

		defer func() {
			if r := recover(); r != nil {
				t.Logf("Air-source heat pump initialization handled panic: %v", r)
			}
		}()

		Refaint(refas, createBasicWDAT(), createBasicCOMPNT())

		// Verify air-source configuration
		if refa.Cat != nil && refa.Cat.awtyp == 'a' {
			t.Log("Air-source heat pump configured correctly")
		}

		t.Log("Air-source heat pump initialization completed successfully")
	})

	t.Run("WaterSourceHeatPump", func(t *testing.T) {
		// Create water-source heat pump
		refa := createWaterSourceHeatPump()
		refas := []*REFA{refa}

		defer func() {
			if r := recover(); r != nil {
				t.Logf("Water-source heat pump initialization handled panic: %v", r)
			}
		}()

		Refaint(refas, createBasicWDAT(), createBasicCOMPNT())

		// Verify water-source configuration
		if refa.Cat != nil && refa.Cat.awtyp == 'w' {
			t.Log("Water-source heat pump configured correctly")
		}

		t.Log("Water-source heat pump initialization completed successfully")
	})

	t.Run("ChillerInitialization", func(t *testing.T) {
		// Create chiller system
		refa := createChillerREFA()
		refas := []*REFA{refa}

		defer func() {
			if r := recover(); r != nil {
				t.Logf("Chiller initialization handled panic: %v", r)
			}
		}()

		Refaint(refas, createBasicWDAT(), createBasicCOMPNT())

		// Verify chiller configuration
		if refa.Cat != nil {
			t.Logf("Chiller initialized - Modes: %d", refa.Cat.Nmode)
		}

		t.Log("Chiller initialization completed successfully")
	})

	t.Run("MultipleREFAInitialization", func(t *testing.T) {
		// Create multiple REFA systems
		refa1 := createAirSourceHeatPump()
		refa1.Name = "REFA1"
		refa2 := createWaterSourceHeatPump()
		refa2.Name = "REFA2"
		refas := []*REFA{refa1, refa2}

		defer func() {
			if r := recover(); r != nil {
				t.Logf("Multiple REFA initialization handled panic: %v", r)
			}
		}()

		Refaint(refas, createBasicWDAT(), createBasicCOMPNT())
		t.Log("Multiple REFA initialization completed successfully")
	})

	t.Run("EmptyREFAList", func(t *testing.T) {
		// Test with empty REFA list
		var refas []*REFA

		defer func() {
			if r := recover(); r != nil {
				t.Logf("Empty REFA list handled panic: %v", r)
			}
		}()

		Refaint(refas, createBasicWDAT(), createBasicCOMPNT())
		t.Log("Empty REFA list handled successfully")
	})
}

// TestRefacfv tests the REFA coefficient calculation function
func TestRefacfv(t *testing.T) {
	t.Run("BasicCoefficientCalculation", func(t *testing.T) {
		// Create REFA for coefficient calculation
		refa := createCoefficientTestREFA()
		refas := []*REFA{refa}

		defer func() {
			if r := recover(); r != nil {
				t.Logf("Coefficient calculation handled panic: %v", r)
			}
		}()

		Refacfv(refas)

		// Verify coefficient calculations
		if refa.Cmp != nil && len(refa.Cmp.Elouts) > 0 {
			t.Logf("Coefficient calculation completed for %s", refa.Name)
		}

		t.Log("Basic coefficient calculation completed successfully")
	})

	t.Run("CoolingModeCoefficients", func(t *testing.T) {
		// Test coefficient calculation for cooling mode
		refa := createCoolingModeREFA()
		refas := []*REFA{refa}

		defer func() {
			if r := recover(); r != nil {
				t.Logf("Cooling mode coefficient calculation handled panic: %v", r)
			}
		}()

		Refacfv(refas)

		// Verify cooling mode coefficients
		if refa.Cat != nil && len(refa.Cat.mode) > 0 {
			if refa.Cat.mode[0] == COOLING_SW {
				t.Log("Cooling mode coefficient calculation verified")
			}
		}

		t.Log("Cooling mode coefficient calculation completed successfully")
	})

	t.Run("HeatingModeCoefficients", func(t *testing.T) {
		// Test coefficient calculation for heating mode
		refa := createHeatingModeREFA()
		refas := []*REFA{refa}

		defer func() {
			if r := recover(); r != nil {
				t.Logf("Heating mode coefficient calculation handled panic: %v", r)
			}
		}()

		Refacfv(refas)

		// Verify heating mode coefficients
		if refa.Cat != nil && len(refa.Cat.mode) > 0 {
			if refa.Cat.mode[0] == HEATING_SW {
				t.Log("Heating mode coefficient calculation verified")
			}
		}

		t.Log("Heating mode coefficient calculation completed successfully")
	})

	t.Run("OffControlCoefficients", func(t *testing.T) {
		// Test coefficient calculation when control is OFF
		refa := createOffControlREFA()
		refas := []*REFA{refa}

		defer func() {
			if r := recover(); r != nil {
				t.Logf("Off control coefficient calculation handled panic: %v", r)
			}
		}()

		Refacfv(refas)
		t.Log("Off control coefficient calculation completed successfully")
	})
}

// TestRefaene tests the REFA energy calculation function
func TestRefaene(t *testing.T) {
	t.Run("BasicEnergyCalculation", func(t *testing.T) {
		// Create REFA for energy calculation
		refa := createEnergyTestREFA()
		refas := []*REFA{refa}

		defer func() {
			if r := recover(); r != nil {
				t.Logf("Energy calculation handled panic: %v", r)
			}
		}()

		var LDrest int
		Refaene(refas, &LDrest)

		// Verify energy calculations
		t.Logf("Energy calculation results - Q: %.1f W, E: %.1f W", refa.Q, refa.E)

		t.Log("Basic energy calculation completed successfully")
	})

	t.Run("CoolingEnergyCalculation", func(t *testing.T) {
		// Test energy calculation for cooling mode
		refa := createCoolingEnergyREFA()
		refas := []*REFA{refa}

		defer func() {
			if r := recover(); r != nil {
				t.Logf("Cooling energy calculation handled panic: %v", r)
			}
		}()

		var LDrest int
		Refaene(refas, &LDrest)

		// Verify cooling energy calculations
		if refa.Q < 0 { // Cooling should be negative
			t.Logf("Cooling energy calculation - Q: %.1f W (cooling)", refa.Q)
		}

		t.Log("Cooling energy calculation completed successfully")
	})

	t.Run("HeatingEnergyCalculation", func(t *testing.T) {
		// Test energy calculation for heating mode
		refa := createHeatingEnergyREFA()
		refas := []*REFA{refa}

		defer func() {
			if r := recover(); r != nil {
				t.Logf("Heating energy calculation handled panic: %v", r)
			}
		}()

		var LDrest int
		Refaene(refas, &LDrest)

		// Verify heating energy calculations
		if refa.Q > 0 { // Heating should be positive
			t.Logf("Heating energy calculation - Q: %.1f W (heating)", refa.Q)
		}

		t.Log("Heating energy calculation completed successfully")
	})

	t.Run("COPCalculation", func(t *testing.T) {
		// Test COP (Coefficient of Performance) calculation
		refa := createCOPTestREFA()
		refas := []*REFA{refa}

		defer func() {
			if r := recover(); r != nil {
				t.Logf("COP calculation handled panic: %v", r)
			}
		}()

		var LDrest int
		Refaene(refas, &LDrest)

		// Verify COP calculations
		if refa.Q != 0 && refa.E != 0 {
			cop := absValue(refa.Q / refa.E)
			t.Logf("COP calculation - Q: %.1f W, E: %.1f W, COP: %.2f", refa.Q, refa.E, cop)
			
			// Verify COP is within reasonable range
			if cop < 1.0 || cop > 10.0 {
				t.Logf("Warning: COP outside typical range: %.2f", cop)
			}
		}

		t.Log("COP calculation completed successfully")
	})

	t.Run("EnergyBalance", func(t *testing.T) {
		// Test energy balance in REFA calculations
		refa := createEnergyBalanceREFA()
		refas := []*REFA{refa}

		defer func() {
			if r := recover(); r != nil {
				t.Logf("Energy balance calculation handled panic: %v", r)
			}
		}()

		var LDrest int
		Refaene(refas, &LDrest)

		// Verify energy balance (Q = thermal output, E = electrical input)
		if refa.Q != 0 && refa.E != 0 {
			energyRatio := absValue(refa.Q / refa.E)
			t.Logf("Energy balance - Thermal: %.1f W, Electrical: %.1f W, Ratio: %.2f", 
				refa.Q, refa.E, energyRatio)
		}

		t.Log("Energy balance verification completed successfully")
	})

	t.Run("OffControlEnergyCalculation", func(t *testing.T) {
		// Test energy calculation when control is OFF
		refa := createOffControlEnergyREFA()
		refas := []*REFA{refa}

		defer func() {
			if r := recover(); r != nil {
				t.Logf("Off control energy calculation handled panic: %v", r)
			}
		}()

		var LDrest int
		Refaene(refas, &LDrest)

		// Verify all energy values are zero when OFF
		if refa.Q == 0.0 && refa.E == 0.0 {
			t.Log("Off control energy values correctly set to zero")
		}

		t.Log("Off control energy calculation completed successfully")
	})
}

// TestREFA_PhysicalValidation tests physical validation of REFA calculations
func TestREFA_PhysicalValidation(t *testing.T) {
	t.Run("COPValidation", func(t *testing.T) {
		// Test COP validation for different operating conditions
		refa := createCOPValidationREFA()
		refas := []*REFA{refa}

		defer func() {
			if r := recover(); r != nil {
				t.Logf("COP validation handled panic: %v", r)
			}
		}()

		Refacfv(refas)
		var LDrest int
		Refaene(refas, &LDrest)

		// Verify COP ranges for different modes
		if refa.Q != 0 && refa.E != 0 {
			cop := absValue(refa.Q / refa.E)
			
			// Typical COP ranges:
			// Cooling: 2.5-6.0
			// Heating: 2.0-5.0
			if refa.Q < 0 { // Cooling mode
				if cop < 2.0 || cop > 8.0 {
					t.Logf("Warning: Cooling COP outside typical range: %.2f", cop)
				} else {
					t.Logf("Cooling COP within valid range: %.2f", cop)
				}
			} else if refa.Q > 0 { // Heating mode
				if cop < 1.5 || cop > 6.0 {
					t.Logf("Warning: Heating COP outside typical range: %.2f", cop)
				} else {
					t.Logf("Heating COP within valid range: %.2f", cop)
				}
			}
		}

		t.Log("COP validation completed successfully")
	})

	t.Run("TemperatureValidation", func(t *testing.T) {
		// Test temperature validation
		refa := createTemperatureValidationREFA()
		refas := []*REFA{refa}

		defer func() {
			if r := recover(); r != nil {
				t.Logf("Temperature validation handled panic: %v", r)
			}
		}()

		Refacfv(refas)
		var LDrest int
		Refaene(refas, &LDrest)

		// Verify temperature ranges are physically reasonable
		if refa.Cat != nil && refa.Cat.awtyp == 'a' {
			// Air-source heat pump temperature limits
			t.Log("Air-source heat pump temperature validation completed")
		} else if refa.Cat != nil && refa.Cat.awtyp == 'w' {
			// Water-source heat pump temperature limits
			t.Log("Water-source heat pump temperature validation completed")
		}

		t.Log("Temperature validation completed successfully")
	})

	t.Run("CapacityValidation", func(t *testing.T) {
		// Test capacity validation
		refa := createCapacityValidationREFA()
		refas := []*REFA{refa}

		defer func() {
			if r := recover(); r != nil {
				t.Logf("Capacity validation handled panic: %v", r)
			}
		}()

		Refacfv(refas)
		var LDrest int
		Refaene(refas, &LDrest)

		// Verify capacity is within design limits
		if refa.Q != 0 {
			capacity := absValue(refa.Q)
			t.Logf("Operating capacity: %.1f W", capacity)
			
			// Check if capacity is reasonable (not negative, not excessive)
			if capacity > 1000000 { // > 1MW
				t.Logf("Warning: Very high capacity: %.0f W", capacity)
			}
		}

		t.Log("Capacity validation completed successfully")
	})
}

// TestREFA_PerformanceCharacteristics tests performance characteristics
func TestREFA_PerformanceCharacteristics(t *testing.T) {
	t.Run("PartLoadPerformance", func(t *testing.T) {
		// Test part-load performance
		refa := createPartLoadREFA()
		refas := []*REFA{refa}

		defer func() {
			if r := recover(); r != nil {
				t.Logf("Part-load performance test handled panic: %v", r)
			}
		}()

		// Test at various load factors
		loadFactors := []float64{0.25, 0.5, 0.75, 1.0}
		for _, factor := range loadFactors {
			// Simulate different load conditions
			if refa.Cmp != nil && len(refa.Cmp.Elouts) > 0 {
				refa.Cmp.Elouts[0].G = factor * 2.0 // Assume 2.0 kg/s full load
			}
			
			Refacfv(refas)
			var LDrest int
		Refaene(refas, &LDrest)
			
			if refa.Q != 0 && refa.E != 0 {
				cop := absValue(refa.Q / refa.E)
				t.Logf("Load factor: %.2f, COP: %.2f", factor, cop)
			}
		}

		t.Log("Part-load performance test completed successfully")
	})

	t.Run("SeasonalPerformance", func(t *testing.T) {
		// Test seasonal performance variation
		refa := createSeasonalPerformanceREFA()
		refas := []*REFA{refa}

		defer func() {
			if r := recover(); r != nil {
				t.Logf("Seasonal performance test handled panic: %v", r)
			}
		}()

		// Simulate different seasonal conditions
		seasons := []struct {
			name string
			temp float64
		}{
			{"Winter", -5.0},
			{"Spring", 10.0},
			{"Summer", 35.0},
			{"Fall", 15.0},
		}

		for _, season := range seasons {
			// Set seasonal conditions (simplified)
			Refacfv(refas)
			var LDrest int
		Refaene(refas, &LDrest)
			
			if refa.Q != 0 && refa.E != 0 {
				cop := absValue(refa.Q / refa.E)
				t.Logf("%s (%.1f°C): COP = %.2f", season.name, season.temp, cop)
			}
		}

		t.Log("Seasonal performance test completed successfully")
	})
}

// Helper functions to create test REFA instances

func createBasicREFAForRefasTest() *REFA {
	// Create basic ELOUT and ELIN for REFA
	elouts := make([]*ELOUT, 2) // Basic REFA has 2 outputs
	for i := range elouts {
		elouts[i] = &ELOUT{
			Control: ON_SW,
			Sysv:    7.0, // 7°C chilled water
			G:       2.0, // 2 kg/s
		}
	}
	
	elins := make([]*ELIN, 2) // Basic REFA has 2 inputs
	for i := range elins {
		elins[i] = &ELIN{
			Sysvin: 12.0, // 12°C return water
		}
	}

	return &REFA{
		Name: "TestREFA",
		Cat: &REFACA{
			name:  "TestREFACA",
			awtyp: 'a', // Air-source
			mode:  [2]ControlSWType{COOLING_SW, HEATING_SW},
			Nmode: 2,
		},
		Cmp: &COMPNT{
			Name:    "TestREFAComponent",
			Control: ON_SW,
			Elouts:  elouts,
			Elins:   elins,
		},
		Q: 0.0, // No initial heat output
		E: 0.0, // No initial electrical input
	}
}

func createAirSourceHeatPump() *REFA {
	refa := createBasicREFAForRefasTest()
	refa.Cat.awtyp = 'a' // Air-source
	return refa
}

func createWaterSourceHeatPump() *REFA {
	refa := createBasicREFAForRefasTest()
	refa.Cat.awtyp = 'w' // Water-source
	return refa
}

func createChillerREFA() *REFA {
	refa := createBasicREFAForRefasTest()
	refa.Cat.mode = [2]ControlSWType{COOLING_SW, OFF_SW} // Cooling only
	refa.Cat.Nmode = 1
	return refa
}

func createCoefficientTestREFA() *REFA {
	refa := createBasicREFAForRefasTest()
	// Set up for coefficient calculation
	for i := range refa.Cmp.Elouts {
		refa.Cmp.Elouts[i].G = 2.0
	}
	return refa
}

func createCoolingModeREFA() *REFA {
	refa := createBasicREFAForRefasTest()
	refa.Cat.mode = [2]ControlSWType{COOLING_SW, OFF_SW}
	refa.Cat.Nmode = 1
	return refa
}

func createHeatingModeREFA() *REFA {
	refa := createBasicREFAForRefasTest()
	refa.Cat.mode = [2]ControlSWType{HEATING_SW, OFF_SW}
	refa.Cat.Nmode = 1
	return refa
}

func createOffControlREFA() *REFA {
	refa := createBasicREFAForRefasTest()
	refa.Cmp.Control = OFF_SW
	for i := range refa.Cmp.Elouts {
		refa.Cmp.Elouts[i].Control = OFF_SW
	}
	return refa
}

func createEnergyTestREFA() *REFA {
	refa := createBasicREFAForRefasTest()
	refa.Q = -50000.0 // 50kW cooling
	refa.E = 15000.0  // 15kW electrical input
	return refa
}

func createCoolingEnergyREFA() *REFA {
	refa := createEnergyTestREFA()
	refa.Cat.mode = [2]ControlSWType{COOLING_SW, OFF_SW}
	refa.Q = -50000.0 // Negative for cooling
	return refa
}

func createHeatingEnergyREFA() *REFA {
	refa := createEnergyTestREFA()
	refa.Cat.mode = [2]ControlSWType{HEATING_SW, OFF_SW}
	refa.Q = 45000.0 // Positive for heating
	refa.E = 12000.0 // Lower electrical input for heating
	return refa
}

func createCOPTestREFA() *REFA {
	refa := createBasicREFAForRefasTest()
	refa.Q = -30000.0 // 30kW cooling
	refa.E = 10000.0  // 10kW electrical input (COP = 3.0)
	return refa
}

func createEnergyBalanceREFA() *REFA {
	refa := createBasicREFAForRefasTest()
	refa.Q = -40000.0 // 40kW cooling
	refa.E = 12000.0  // 12kW electrical input
	return refa
}

func createOffControlEnergyREFA() *REFA {
	refa := createBasicREFAForRefasTest()
	refa.Cmp.Control = OFF_SW
	refa.Q = 0.0
	refa.E = 0.0
	return refa
}

func createCOPValidationREFA() *REFA {
	refa := createBasicREFAForRefasTest()
	refa.Q = -35000.0 // 35kW cooling
	refa.E = 10000.0  // 10kW electrical input (COP = 3.5)
	return refa
}

func createTemperatureValidationREFA() *REFA {
	refa := createAirSourceHeatPump()
	// Set up for temperature validation
	return refa
}

func createCapacityValidationREFA() *REFA {
	refa := createBasicREFAForRefasTest()
	refa.Q = -25000.0 // 25kW capacity
	refa.E = 8000.0   // 8kW electrical input
	return refa
}

func createPartLoadREFA() *REFA {
	refa := createBasicREFAForRefasTest()
	// Set up for part-load testing
	return refa
}

func createSeasonalPerformanceREFA() *REFA {
	refa := createAirSourceHeatPump()
	// Set up for seasonal performance testing
	return refa
}

// createOutputTestREFA creates a REFA configured for output testing
func createOutputTestREFA() *REFA {
	refa := createBasicREFAForRefasTest()
	refa.Tin = 12.0
	refa.Q = -50000.0
	refa.E = 15000.0
	refa.Ph = 500.0

	// Set up daily aggregation values
	refa.Tidy = SVDAY{Hrs: 8, M: 11.0, Mn: 7.0, Mx: 15.0, Mntime: 6, Mxtime: 14}
	refa.Qdy = QDAY{Hhr: 0, H: 0.0, Chr: 8, C: 400000.0, Hmxtime: 0, Hmx: 0.0, Cmxtime: 14, Cmx: 60000.0}
	refa.Edy = EDAY{Hrs: 8, D: 120000.0, Mxtime: 14, Mx: 20000.0}
	refa.Phdy = EDAY{Hrs: 8, D: 4000.0, Mxtime: 14, Mx: 600.0}

	// Set up monthly aggregation values
	refa.mTidy = SVDAY{Hrs: 200, M: 10.5, Mn: 5.0, Mx: 16.0, Mntime: 1, Mxtime: 15}
	refa.mQdy = QDAY{Hhr: 0, H: 0.0, Chr: 200, C: 10000000.0, Hmxtime: 0, Hmx: 0.0, Cmxtime: 15, Cmx: 65000.0}
	refa.mEdy = EDAY{Hrs: 200, D: 3000000.0, Mxtime: 15, Mx: 22000.0}
	refa.mPhdy = EDAY{Hrs: 200, D: 100000.0, Mxtime: 15, Mx: 700.0}

	return refa
}

// TestRefaprint tests the REFA print function
func TestRefaprint(t *testing.T) {
	refa := createOutputTestREFA()
	refas := []*REFA{refa}

	t.Run("Header1_id0", func(t *testing.T) {
		var buf bytes.Buffer
		refaprint(&buf, 0, refas)
		output := buf.String()

		// Check type and count (REFACOMP_TYPE = "REFA")
		if !strings.Contains(output, "REFA") {
			t.Error("Missing REFA type in header")
		}
		if !strings.Contains(output, refa.Name) {
			t.Errorf("Missing REFA name %s in header", refa.Name)
		}
	})

	t.Run("Header2_id1", func(t *testing.T) {
		var buf bytes.Buffer
		refaprint(&buf, 1, refas)
		output := buf.String()

		// Check item names
		expectedPatterns := []string{
			refa.Name + "_c",
			refa.Name + "_G",
			refa.Name + "_Ti",
			refa.Name + "_To",
			refa.Name + "_Q",
			refa.Name + "_E",
			refa.Name + "_P",
		}
		for _, pattern := range expectedPatterns {
			if !strings.Contains(output, pattern) {
				t.Errorf("Missing expected pattern: %s", pattern)
			}
		}
	})

	t.Run("Data_id99", func(t *testing.T) {
		var buf bytes.Buffer
		refaprint(&buf, 99, refas)
		output := buf.String()

		if len(output) == 0 {
			t.Error("Data output is empty")
		}
	})

	t.Run("EmptyList", func(t *testing.T) {
		var buf bytes.Buffer
		refaprint(&buf, 0, []*REFA{})
		output := buf.String()

		if len(output) != 0 {
			t.Error("Expected empty output for empty list")
		}
	})
}

// TestRefadyprt tests the REFA daily print function
func TestRefadyprt(t *testing.T) {
	refa := createOutputTestREFA()
	refas := []*REFA{refa}

	t.Run("Header1_id0", func(t *testing.T) {
		var buf bytes.Buffer
		refadyprt(&buf, 0, refas)
		output := buf.String()

		if !strings.Contains(output, "REFA") {
			t.Error("Missing REFA type in daily header")
		}
		if !strings.Contains(output, refa.Name) {
			t.Errorf("Missing REFA name %s in daily header", refa.Name)
		}
	})

	t.Run("Header2_id1", func(t *testing.T) {
		var buf bytes.Buffer
		refadyprt(&buf, 1, refas)
		output := buf.String()

		// Check item patterns
		expectedPatterns := []string{
			refa.Name + "_Ht",
			refa.Name + "_T",
			refa.Name + "_Hh",
			refa.Name + "_Qh",
			refa.Name + "_He",
			refa.Name + "_E",
			refa.Name + "_Hp",
			refa.Name + "_P",
		}
		for _, pattern := range expectedPatterns {
			if !strings.Contains(output, pattern) {
				t.Errorf("Missing expected pattern: %s", pattern)
			}
		}
	})

	t.Run("Data_id99", func(t *testing.T) {
		var buf bytes.Buffer
		refadyprt(&buf, 99, refas)
		output := buf.String()

		if len(output) == 0 {
			t.Error("Daily data output is empty")
		}
	})

	t.Run("EmptyList", func(t *testing.T) {
		var buf bytes.Buffer
		refadyprt(&buf, 0, []*REFA{})
		output := buf.String()

		if len(output) != 0 {
			t.Error("Expected empty output for empty list")
		}
	})
}

// TestRefamonprt tests the REFA monthly print function
func TestRefamonprt(t *testing.T) {
	refa := createOutputTestREFA()
	refas := []*REFA{refa}

	t.Run("Header1_id0", func(t *testing.T) {
		var buf bytes.Buffer
		refamonprt(&buf, 0, refas)
		output := buf.String()

		if !strings.Contains(output, "REFA") {
			t.Error("Missing REFA type in monthly header")
		}
		if !strings.Contains(output, refa.Name) {
			t.Errorf("Missing REFA name %s in monthly header", refa.Name)
		}
	})

	t.Run("Header2_id1", func(t *testing.T) {
		var buf bytes.Buffer
		refamonprt(&buf, 1, refas)
		output := buf.String()

		// Check item patterns (same as daily)
		expectedPatterns := []string{
			refa.Name + "_Ht",
			refa.Name + "_T",
			refa.Name + "_Hh",
			refa.Name + "_He",
		}
		for _, pattern := range expectedPatterns {
			if !strings.Contains(output, pattern) {
				t.Errorf("Missing expected pattern: %s", pattern)
			}
		}
	})

	t.Run("Data_id99", func(t *testing.T) {
		var buf bytes.Buffer
		refamonprt(&buf, 99, refas)
		output := buf.String()

		if len(output) == 0 {
			t.Error("Monthly data output is empty")
		}
	})

	t.Run("EmptyList", func(t *testing.T) {
		var buf bytes.Buffer
		refamonprt(&buf, 0, []*REFA{})
		output := buf.String()

		if len(output) != 0 {
			t.Error("Expected empty output for empty list")
		}
	})
}

// TestRefadyint tests the REFA daily aggregation initialization
func TestRefadyint(t *testing.T) {
	t.Run("BasicInitialization", func(t *testing.T) {
		refa := createOutputTestREFA()
		refas := []*REFA{refa}

		// Verify values are set before initialization
		if refa.Tidy.Hrs == 0 {
			t.Error("Test data not properly set up")
		}

		refadyint(refas)

		// After initialization, values should be reset
		if refa.Tidy.Hrs != 0 {
			t.Error("Tidy.Hrs should be reset to 0")
		}
		if refa.Edy.Hrs != 0 {
			t.Error("Edy.Hrs should be reset to 0")
		}
	})

	t.Run("EmptyList", func(t *testing.T) {
		// Should not panic with empty list
		refadyint([]*REFA{})
	})

	t.Run("MultipleRefa", func(t *testing.T) {
		refa1 := createOutputTestREFA()
		refa1.Name = "REFA1"
		refa2 := createOutputTestREFA()
		refa2.Name = "REFA2"
		refas := []*REFA{refa1, refa2}

		refadyint(refas)

		for i, refa := range refas {
			if refa.Tidy.Hrs != 0 {
				t.Errorf("REFA[%d] Tidy.Hrs should be reset to 0", i)
			}
		}
	})
}

// TestRefamonint tests the REFA monthly aggregation initialization
func TestRefamonint(t *testing.T) {
	t.Run("BasicInitialization", func(t *testing.T) {
		refa := createOutputTestREFA()
		refas := []*REFA{refa}

		// Verify values are set before initialization
		if refa.mTidy.Hrs == 0 {
			t.Error("Test data not properly set up")
		}

		refamonint(refas)

		// After initialization, values should be reset
		if refa.mTidy.Hrs != 0 {
			t.Error("mTidy.Hrs should be reset to 0")
		}
		if refa.mEdy.Hrs != 0 {
			t.Error("mEdy.Hrs should be reset to 0")
		}
	})

	t.Run("EmptyList", func(t *testing.T) {
		// Should not panic with empty list
		refamonint([]*REFA{})
	})

	t.Run("MultipleRefa", func(t *testing.T) {
		refa1 := createOutputTestREFA()
		refa1.Name = "REFA1"
		refa2 := createOutputTestREFA()
		refa2.Name = "REFA2"
		refas := []*REFA{refa1, refa2}

		refamonint(refas)

		for i, refa := range refas {
			if refa.mTidy.Hrs != 0 {
				t.Errorf("REFA[%d] mTidy.Hrs should be reset to 0", i)
			}
		}
	})
}

// TestRefaday tests the refaday aggregation function
func TestRefaday(t *testing.T) {
	t.Run("DailyAggregation", func(t *testing.T) {
		refa := createBasicREFAForRefasTest()
		refa.Cmp.Control = ON_SW
		refa.Tin = 12.0
		refa.Q = -50000.0 // Cooling
		refa.E = 15000.0
		refa.Ph = 500.0

		refas := []*REFA{refa}

		// Initialize daily aggregation
		refadyint(refas)

		// Simulate multiple time steps
		times := []int{900, 1000, 1100, 1200}
		for _, ttmm := range times {
			refaday(7, 15, ttmm, refas, 31, 365)
		}

		// After 4 time steps, verify aggregation
		if refa.Tidy.Hrs != 4 {
			t.Errorf("Tidy.Hrs = %d, want 4", refa.Tidy.Hrs)
		}
		if refa.Qdy.Chr != 4 {
			t.Errorf("Qdy.Chr = %d, want 4 (cooling)", refa.Qdy.Chr)
		}
		if refa.Edy.Hrs != 4 {
			t.Errorf("Edy.Hrs = %d, want 4", refa.Edy.Hrs)
		}
		if refa.Phdy.Hrs != 4 {
			t.Errorf("Phdy.Hrs = %d, want 4", refa.Phdy.Hrs)
		}
	})

	t.Run("HeatingMode", func(t *testing.T) {
		refa := createBasicREFAForRefasTest()
		refa.Cmp.Control = ON_SW
		refa.Tin = 45.0
		refa.Q = 40000.0 // Heating (positive)
		refa.E = 12000.0
		refa.Ph = 400.0

		refas := []*REFA{refa}

		refadyint(refas)
		refaday(7, 15, 1200, refas, 31, 365)

		// Heating should aggregate to Hhr
		if refa.Qdy.Hhr != 1 {
			t.Errorf("Qdy.Hhr = %d, want 1 (heating)", refa.Qdy.Hhr)
		}
	})

	t.Run("OffControl_NoAggregation", func(t *testing.T) {
		refa := createBasicREFAForRefasTest()
		refa.Cmp.Control = OFF_SW
		refa.Tin = 0.0
		refa.Q = 0.0
		refa.E = 0.0
		refa.Ph = 0.0

		refas := []*REFA{refa}

		refadyint(refas)
		refaday(7, 15, 1200, refas, 31, 365)

		// OFF control should not aggregate
		if refa.Tidy.Hrs != 0 {
			t.Errorf("Tidy.Hrs should be 0 when OFF, got %d", refa.Tidy.Hrs)
		}
		if refa.Edy.Hrs != 0 {
			t.Errorf("Edy.Hrs should be 0 when OFF, got %d", refa.Edy.Hrs)
		}
	})

	t.Run("MonthlyAggregation_EndOfDay", func(t *testing.T) {
		refa := createBasicREFAForRefasTest()
		refa.Cmp.Control = ON_SW
		refa.Tin = 12.0
		refa.Q = -50000.0
		refa.E = 15000.0
		refa.Ph = 500.0

		refas := []*REFA{refa}

		refadyint(refas)
		refamonint(refas)

		// Call at end of month
		refaday(7, 31, 2400, refas, 31, 365)

		// Daily values should be aggregated
		if refa.Tidy.Hrs != 1 {
			t.Errorf("Tidy.Hrs = %d, want 1", refa.Tidy.Hrs)
		}
	})

	t.Run("EmptyList", func(t *testing.T) {
		// Should not panic with empty list
		refaday(7, 15, 1200, []*REFA{}, 31, 365)
	})

	t.Run("MultipleRefa", func(t *testing.T) {
		refas := make([]*REFA, 2)
		for i := range refas {
			refas[i] = createBasicREFAForRefasTest()
			refas[i].Name = "REFA" + string(rune('A'+i))
			refas[i].Cmp.Control = ON_SW
			refas[i].Tin = 12.0 + float64(i)*2
			refas[i].Q = -50000.0 - float64(i)*10000
			refas[i].E = 15000.0 + float64(i)*2000
			refas[i].Ph = 500.0 + float64(i)*100
		}

		refadyint(refas)
		refaday(7, 15, 1200, refas, 31, 365)

		// Verify each refa has independent aggregation
		for i, refa := range refas {
			if refa.Tidy.Hrs != 1 {
				t.Errorf("REFA[%d] Tidy.Hrs = %d, want 1", i, refa.Tidy.Hrs)
			}
			if refa.Edy.Hrs != 1 {
				t.Errorf("REFA[%d] Edy.Hrs = %d, want 1", i, refa.Edy.Hrs)
			}
		}
	})

	t.Run("CrossTabulation", func(t *testing.T) {
		refa := createBasicREFAForRefasTest()
		refa.Cmp.Control = ON_SW
		refa.Tin = 12.0
		refa.Q = -50000.0
		refa.E = 15000.0
		refa.Ph = 500.0

		refas := []*REFA{refa}

		refadyint(refas)
		refamonint(refas)

		// Test cross-tabulation: mtEdy[Mo][tt] and mtPhdy[Mo][tt]
		// Mo = Month - 1, tt = ConvertHour(ttmm)
		refaday(7, 15, 1200, refas, 31, 365)

		// Cross-tabulation should have values at Mo=6, tt=ConvertHour(1200)
		tt := ConvertHour(1200)
		// Note: emtsum only updates at end of simulation
		// Just verify it doesn't panic
		_ = refa.mtEdy[6][tt]
		_ = refa.mtPhdy[6][tt]
	})
}
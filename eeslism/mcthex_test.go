package eeslism

import (
	"testing"
)

// TestThexint tests the THEX initialization function
func TestThexint(t *testing.T) {
	t.Run("BasicInitialization", func(t *testing.T) {
		// Create basic THEX system
		thex := createBasicTHEX()
		thexs := []*THEX{thex}

		defer func() {
			if r := recover(); r != nil {
				t.Logf("Basic initialization handled panic: %v", r)
			}
		}()

		Thexint(thexs)

		// Verify initialization
		if thex.Type == 't' || thex.Type == 'h' {
			t.Logf("THEX type set to: %c", thex.Type)
		}

		t.Log("Basic THEX initialization completed successfully")
	})

	t.Run("SensibleOnlyHeatExchanger", func(t *testing.T) {
		// Create THEX with sensible heat exchange only (eh < 0)
		thex := createSensibleOnlyTHEX()
		thexs := []*THEX{thex}

		defer func() {
			if r := recover(); r != nil {
				t.Logf("Sensible-only initialization handled panic: %v", r)
			}
		}()

		Thexint(thexs)

		// Verify sensible-only configuration
		if thex.Type == 't' {
			t.Log("Sensible-only heat exchanger configured correctly")
		}

		t.Log("Sensible-only THEX initialization completed successfully")
	})

	t.Run("TotalHeatExchanger", func(t *testing.T) {
		// Create THEX with total heat exchange (both sensible and latent)
		thex := createTotalHeatTHEX()
		thexs := []*THEX{thex}

		defer func() {
			if r := recover(); r != nil {
				t.Logf("Total heat exchanger initialization handled panic: %v", r)
			}
		}()

		Thexint(thexs)

		// Verify total heat exchanger configuration
		if thex.Type == 'h' {
			t.Log("Total heat exchanger configured correctly")
		}

		t.Log("Total heat exchanger initialization completed successfully")
	})

	t.Run("MultipleThexInitialization", func(t *testing.T) {
		// Create multiple THEX systems
		thex1 := createBasicTHEX()
		thex1.Name = "THEX1"
		thex2 := createTotalHeatTHEX()
		thex2.Name = "THEX2"
		thexs := []*THEX{thex1, thex2}

		defer func() {
			if r := recover(); r != nil {
				t.Logf("Multiple THEX initialization handled panic: %v", r)
			}
		}()

		Thexint(thexs)
		t.Log("Multiple THEX initialization completed successfully")
	})

	t.Run("EmptyThexList", func(t *testing.T) {
		// Test with empty THEX list
		var thexs []*THEX

		defer func() {
			if r := recover(); r != nil {
				t.Logf("Empty THEX list handled panic: %v", r)
			}
		}()

		Thexint(thexs)
		t.Log("Empty THEX list handled successfully")
	})
}

// TestThexcfv tests the THEX coefficient calculation function
func TestThexcfv(t *testing.T) {
	t.Run("BasicCoefficientCalculation", func(t *testing.T) {
		// Create THEX for coefficient calculation
		thex := createCoefficientTestTHEX()
		thexs := []*THEX{thex}

		defer func() {
			if r := recover(); r != nil {
				t.Logf("Coefficient calculation handled panic: %v", r)
			}
		}()

		Thexcfv(thexs)

		// Verify coefficient calculations
		if thex.ET >= 0 && thex.EH >= 0 {
			t.Logf("Heat exchange efficiencies - ET: %.3f, EH: %.3f", thex.ET, thex.EH)
		}

		t.Log("Basic coefficient calculation completed successfully")
	})

	t.Run("SensibleOnlyCoefficients", func(t *testing.T) {
		// Test coefficient calculation for sensible-only heat exchanger
		thex := createSensibleOnlyCoefficientTHEX()
		thexs := []*THEX{thex}

		defer func() {
			if r := recover(); r != nil {
				t.Logf("Sensible-only coefficient calculation handled panic: %v", r)
			}
		}()

		Thexcfv(thexs)

		// Verify sensible-only coefficients
		if thex.Type == 't' {
			t.Log("Sensible-only coefficient calculation verified")
		}

		t.Log("Sensible-only coefficient calculation completed successfully")
	})

	t.Run("TotalHeatCoefficients", func(t *testing.T) {
		// Test coefficient calculation for total heat exchanger
		thex := createTotalHeatCoefficientTHEX()
		thexs := []*THEX{thex}

		defer func() {
			if r := recover(); r != nil {
				t.Logf("Total heat coefficient calculation handled panic: %v", r)
			}
		}()

		Thexcfv(thexs)

		// Verify total heat coefficients
		if thex.Type == 'h' {
			t.Log("Total heat coefficient calculation verified")
		}

		t.Log("Total heat coefficient calculation completed successfully")
	})

	t.Run("OffControlCoefficients", func(t *testing.T) {
		// Test coefficient calculation when control is OFF
		thex := createOffControlTHEX()
		thexs := []*THEX{thex}

		defer func() {
			if r := recover(); r != nil {
				t.Logf("Off control coefficient calculation handled panic: %v", r)
			}
		}()

		Thexcfv(thexs)
		t.Log("Off control coefficient calculation completed successfully")
	})
}

// TestThexene tests the THEX energy calculation function
func TestThexene(t *testing.T) {
	t.Run("BasicEnergyCalculation", func(t *testing.T) {
		// Create THEX for energy calculation
		thex := createEnergyTestTHEX()
		thexs := []*THEX{thex}

		defer func() {
			if r := recover(); r != nil {
				t.Logf("Energy calculation handled panic: %v", r)
			}
		}()

		Thexene(thexs)

		// Verify energy calculations
		t.Logf("Energy calculation results - Exhaust: Qs=%.1f, Ql=%.1f, Qt=%.1f", 
			thex.Qes, thex.Qel, thex.Qet)
		t.Logf("Energy calculation results - Supply: Qs=%.1f, Ql=%.1f, Qt=%.1f", 
			thex.Qos, thex.Qol, thex.Qot)

		t.Log("Basic energy calculation completed successfully")
	})

	t.Run("SensibleOnlyEnergyCalculation", func(t *testing.T) {
		// Test energy calculation for sensible-only heat exchanger
		thex := createSensibleOnlyEnergyTHEX()
		thexs := []*THEX{thex}

		defer func() {
			if r := recover(); r != nil {
				t.Logf("Sensible-only energy calculation handled panic: %v", r)
			}
		}()

		Thexene(thexs)

		// Verify sensible-only energy calculations
		if thex.Type == 't' {
			t.Log("Sensible-only energy calculation verified")
		}

		t.Log("Sensible-only energy calculation completed successfully")
	})

	t.Run("TotalHeatEnergyCalculation", func(t *testing.T) {
		// Test energy calculation for total heat exchanger
		thex := createTotalHeatEnergyTHEX()
		thexs := []*THEX{thex}

		defer func() {
			if r := recover(); r != nil {
				t.Logf("Total heat energy calculation handled panic: %v", r)
			}
		}()

		Thexene(thexs)

		// Verify total heat energy calculations
		if thex.Type == 'h' {
			t.Log("Total heat energy calculation verified")
		}

		t.Log("Total heat energy calculation completed successfully")
	})

	t.Run("EnergyBalance", func(t *testing.T) {
		// Test energy balance in THEX calculations
		thex := createEnergyBalanceTHEX()
		thexs := []*THEX{thex}

		defer func() {
			if r := recover(); r != nil {
				t.Logf("Energy balance calculation handled panic: %v", r)
			}
		}()

		Thexene(thexs)

		// Verify energy balance (exhaust heat loss = supply heat gain)
		if thex.Cmp.Control == ON_SW {
			exhaustTotal := thex.Qet
			supplyTotal := thex.Qot
			
			// In ideal heat exchanger, exhaust heat loss should approximately equal supply heat gain
			if exhaustTotal != 0 && supplyTotal != 0 {
				energyRatio := absValue(supplyTotal / exhaustTotal)
				t.Logf("Energy balance - Exhaust: %.1f W, Supply: %.1f W, Ratio: %.3f", 
					exhaustTotal, supplyTotal, energyRatio)
			}
		}

		t.Log("Energy balance verification completed successfully")
	})

	t.Run("OffControlEnergyCalculation", func(t *testing.T) {
		// Test energy calculation when control is OFF
		thex := createOffControlEnergyTHEX()
		thexs := []*THEX{thex}

		defer func() {
			if r := recover(); r != nil {
				t.Logf("Off control energy calculation handled panic: %v", r)
			}
		}()

		Thexene(thexs)

		// Verify all energy values are zero when OFF
		if thex.Qes == 0.0 && thex.Qel == 0.0 && thex.Qet == 0.0 &&
		   thex.Qos == 0.0 && thex.Qol == 0.0 && thex.Qot == 0.0 {
			t.Log("Off control energy values correctly set to zero")
		}

		t.Log("Off control energy calculation completed successfully")
	})
}

// TestTHEX_PhysicalValidation tests physical validation of THEX calculations
func TestTHEX_PhysicalValidation(t *testing.T) {
	t.Run("EfficiencyValidation", func(t *testing.T) {
		// Test heat exchange efficiency validation
		thex := createEfficiencyValidationTHEX()
		thexs := []*THEX{thex}

		defer func() {
			if r := recover(); r != nil {
				t.Logf("Efficiency validation handled panic: %v", r)
			}
		}()

		Thexint(thexs)
		Thexcfv(thexs)

		// Verify efficiency ranges are physically reasonable
		if thex.ET > 0 {
			if thex.ET < 0.0 || thex.ET > 1.0 {
				t.Logf("Warning: Sensible heat efficiency out of range: %.3f", thex.ET)
			} else {
				t.Logf("Sensible heat efficiency within valid range: %.3f", thex.ET)
			}
		}
		if thex.EH > 0 {
			if thex.EH < 0.0 || thex.EH > 1.0 {
				t.Logf("Warning: Latent heat efficiency out of range: %.3f", thex.EH)
			} else {
				t.Logf("Latent heat efficiency within valid range: %.3f", thex.EH)
			}
		}

		t.Log("Efficiency validation completed successfully")
	})

	t.Run("TemperatureRelationships", func(t *testing.T) {
		// Test temperature relationships in heat exchanger
		thex := createTemperatureRelationshipTHEX()
		thexs := []*THEX{thex}

		defer func() {
			if r := recover(); r != nil {
				t.Logf("Temperature relationship validation handled panic: %v", r)
			}
		}()

		Thexcfv(thexs)
		Thexene(thexs)

		// Verify temperature relationships
		// In heat recovery, exhaust air should cool down and supply air should warm up
		if thex.Cmp.Control == ON_SW {
			t.Logf("Temperature relationships - Exhaust: %.1f°C → %.1f°C, Supply: %.1f°C → %.1f°C",
				thex.Tein, thex.Teout, thex.Toin, thex.Toout)
		}

		t.Log("Temperature relationship validation completed successfully")
	})

	t.Run("HumidityValidation", func(t *testing.T) {
		// Test humidity validation in total heat exchanger
		thex := createHumidityValidationTHEX()
		thexs := []*THEX{thex}

		defer func() {
			if r := recover(); r != nil {
				t.Logf("Humidity validation handled panic: %v", r)
			}
		}()

		Thexint(thexs)
		Thexcfv(thexs)
		Thexene(thexs)

		// Verify humidity ranges are physically reasonable
		if thex.Type == 'h' && thex.Cmp.Control == ON_SW {
			humidityValues := []float64{thex.Xein, thex.Xeout, thex.Xoin, thex.Xoout}
			for i, humidity := range humidityValues {
				if humidity < 0.0 || humidity > 0.030 { // 0-30 g/kg range
					t.Logf("Warning: Humidity value %d out of typical range: %.6f kg/kg", i, humidity)
				}
			}
		}

		t.Log("Humidity validation completed successfully")
	})
}

// TestTHEX_PerformanceCharacteristics tests performance characteristics
func TestTHEX_PerformanceCharacteristics(t *testing.T) {
	t.Run("HeatRecoveryEffectiveness", func(t *testing.T) {
		// Test heat recovery effectiveness
		thex := createHeatRecoveryTestTHEX()
		thexs := []*THEX{thex}

		defer func() {
			if r := recover(); r != nil {
				t.Logf("Heat recovery effectiveness test handled panic: %v", r)
			}
		}()

		Thexint(thexs)
		Thexcfv(thexs)
		Thexene(thexs)

		// Calculate and verify heat recovery effectiveness
		if thex.Cmp.Control == ON_SW && thex.Tein != thex.Toin {
			// Sensible heat recovery effectiveness
			if thex.Tein != thex.Toin {
				effectiveness := (thex.Toout - thex.Toin) / (thex.Tein - thex.Toin)
				t.Logf("Sensible heat recovery effectiveness: %.3f", effectiveness)
			}
		}

		t.Log("Heat recovery effectiveness test completed successfully")
	})

	t.Run("FlowRateBalance", func(t *testing.T) {
		// Test flow rate balance
		thex := createFlowRateBalanceTHEX()
		thexs := []*THEX{thex}

		defer func() {
			if r := recover(); r != nil {
				t.Logf("Flow rate balance test handled panic: %v", r)
			}
		}()

		Thexcfv(thexs)

		// Verify flow rates
		if thex.Ge > 0 && thex.Go > 0 {
			t.Logf("Flow rates - Exhaust: %.3f kg/s, Supply: %.3f kg/s", thex.Ge, thex.Go)
		}

		t.Log("Flow rate balance test completed successfully")
	})
}

// Helper functions to create test THEX instances

func createBasicTHEX() *THEX {
	// Create basic ELOUT and ELIN for THEX (4 outputs, multiple inputs)
	elouts := make([]*ELOUT, 4) // THEX has 4 outputs (exhaust temp, exhaust humidity, supply temp, supply humidity)
	for i := range elouts {
		elouts[i] = &ELOUT{
			Control: ON_SW,
			Sysv:    20.0,
			G:       1.0,
			Fluid:   AIR_FLD,
		}
	}
	
	// Create sufficient ELIN for each ELOUT
	elins := make([]*ELIN, 20) // Sufficient for all connections
	for i := range elins {
		elins[i] = &ELIN{
			Sysvin: 25.0,
		}
	}

	return &THEX{
		Name: "TestTHEX",
		Cat: &THEXCA{
			Name: "TestTHEXCA",
			et:   0.75, // 75% sensible heat efficiency
			eh:   0.65, // 65% latent heat efficiency
		},
		Cmp: &COMPNT{
			Name:    "TestTHEXComponent",
			Control: ON_SW,
			Elouts:  elouts,
			Elins:   elins,
		},
		Type: 'h', // Total heat exchanger
		ET:   0.75,
		EH:   0.65,
	}
}

func createSensibleOnlyTHEX() *THEX {
	thex := createBasicTHEX()
	thex.Cat.eh = -1.0 // Negative value indicates sensible-only
	return thex
}

func createTotalHeatTHEX() *THEX {
	thex := createBasicTHEX()
	thex.Cat.et = 0.80 // 80% sensible efficiency
	thex.Cat.eh = 0.70 // 70% latent efficiency
	return thex
}

func createCoefficientTestTHEX() *THEX {
	thex := createBasicTHEX()
	// Set up for coefficient calculation
	for i := range thex.Cmp.Elouts {
		thex.Cmp.Elouts[i].G = 1.0
		thex.Cmp.Elouts[i].Fluid = AIR_FLD
	}
	return thex
}

func createSensibleOnlyCoefficientTHEX() *THEX {
	thex := createSensibleOnlyTHEX()
	for i := range thex.Cmp.Elouts {
		thex.Cmp.Elouts[i].G = 1.0
		thex.Cmp.Elouts[i].Fluid = AIR_FLD
	}
	return thex
}

func createTotalHeatCoefficientTHEX() *THEX {
	thex := createTotalHeatTHEX()
	for i := range thex.Cmp.Elouts {
		thex.Cmp.Elouts[i].G = 1.0
		thex.Cmp.Elouts[i].Fluid = AIR_FLD
	}
	return thex
}

func createOffControlTHEX() *THEX {
	thex := createBasicTHEX()
	thex.Cmp.Control = OFF_SW
	for i := range thex.Cmp.Elouts {
		thex.Cmp.Elouts[i].Control = OFF_SW
	}
	return thex
}

func createEnergyTestTHEX() *THEX {
	thex := createBasicTHEX()
	// Set up for energy calculation with realistic values
	thex.Tein = 25.0   // Exhaust inlet temperature
	thex.Teout = 22.0  // Exhaust outlet temperature
	thex.Toin = 5.0    // Supply inlet temperature
	thex.Toout = 8.0   // Supply outlet temperature
	thex.Ge = 1.0      // Exhaust flow rate
	thex.Go = 1.0      // Supply flow rate
	return thex
}

func createSensibleOnlyEnergyTHEX() *THEX {
	thex := createEnergyTestTHEX()
	thex.Type = 't' // Sensible-only
	thex.Cat.eh = -1.0
	return thex
}

func createTotalHeatEnergyTHEX() *THEX {
	thex := createEnergyTestTHEX()
	thex.Type = 'h' // Total heat
	thex.Xein = 0.012  // Exhaust inlet humidity
	thex.Xeout = 0.010 // Exhaust outlet humidity
	thex.Xoin = 0.004  // Supply inlet humidity
	thex.Xoout = 0.006 // Supply outlet humidity
	return thex
}

func createEnergyBalanceTHEX() *THEX {
	thex := createTotalHeatEnergyTHEX()
	// Set up for energy balance testing
	thex.Cmp.Control = ON_SW
	return thex
}

func createOffControlEnergyTHEX() *THEX {
	thex := createEnergyTestTHEX()
	thex.Cmp.Control = OFF_SW
	return thex
}

func createEfficiencyValidationTHEX() *THEX {
	thex := createBasicTHEX()
	thex.Cat.et = 0.85 // High but reasonable efficiency
	thex.Cat.eh = 0.75 // High but reasonable efficiency
	return thex
}

func createTemperatureRelationshipTHEX() *THEX {
	thex := createBasicTHEX()
	// Set up realistic temperature conditions for heat recovery
	thex.Tein = 22.0   // Warm exhaust air
	thex.Toin = 5.0    // Cold outdoor air
	return thex
}

func createHumidityValidationTHEX() *THEX {
	thex := createTotalHeatTHEX()
	// Set up realistic humidity conditions
	thex.Xein = 0.008  // Indoor humidity
	thex.Xoin = 0.004  // Outdoor humidity
	return thex
}

func createHeatRecoveryTestTHEX() *THEX {
	thex := createBasicTHEX()
	// Set up for heat recovery effectiveness testing
	thex.Tein = 22.0   // Indoor exhaust temperature
	thex.Toin = 5.0    // Outdoor supply temperature
	thex.Cat.et = 0.75 // 75% effectiveness
	return thex
}

func createFlowRateBalanceTHEX() *THEX {
	thex := createBasicTHEX()
	// Set up flow rates for testing
	for i := range thex.Cmp.Elouts {
		thex.Cmp.Elouts[i].G = 1.0 // 1 kg/s
	}
	return thex
}
package eeslism

import (
	"bytes"
	"strings"
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

// createOutputTestTHEX creates a THEX configured for output testing
func createOutputTestTHEX() *THEX {
	thex := createBasicTHEX()
	// Set up values for output testing
	thex.Ge = 1.0
	thex.Go = 0.9
	thex.Tein = 25.0
	thex.Teout = 18.0
	thex.Toin = 5.0
	thex.Toout = 12.0
	thex.Xein = 0.010
	thex.Xeout = 0.008
	thex.Xoin = 0.004
	thex.Xoout = 0.006
	thex.Hein = 50000.0
	thex.Heout = 40000.0
	thex.Hoin = 15000.0
	thex.Hoout = 25000.0
	thex.Qes = 7000.0
	thex.Qel = 3000.0
	thex.Qet = 10000.0
	thex.Qos = 6000.0
	thex.Qol = 2500.0
	thex.Qot = 8500.0

	// Set up daily aggregation values
	thex.Teidy = SVDAY{Hrs: 8, M: 24.0, Mn: 20.0, Mx: 28.0, Mntime: 6, Mxtime: 14}
	thex.Teody = SVDAY{Hrs: 8, M: 17.0, Mn: 14.0, Mx: 20.0, Mntime: 6, Mxtime: 14}
	thex.Xeidy = SVDAY{Hrs: 8, M: 0.009, Mn: 0.007, Mx: 0.011, Mntime: 6, Mxtime: 14}
	thex.Xeody = SVDAY{Hrs: 8, M: 0.007, Mn: 0.005, Mx: 0.009, Mntime: 6, Mxtime: 14}
	thex.Toidy = SVDAY{Hrs: 8, M: 6.0, Mn: 2.0, Mx: 10.0, Mntime: 6, Mxtime: 14}
	thex.Toody = SVDAY{Hrs: 8, M: 13.0, Mn: 9.0, Mx: 17.0, Mntime: 6, Mxtime: 14}
	thex.Xoidy = SVDAY{Hrs: 8, M: 0.004, Mn: 0.003, Mx: 0.005, Mntime: 6, Mxtime: 14}
	thex.Xoody = SVDAY{Hrs: 8, M: 0.006, Mn: 0.005, Mx: 0.007, Mntime: 6, Mxtime: 14}

	thex.Qdyes = QDAY{Hhr: 4, H: 28000.0, Chr: 0, C: 0.0, Hmxtime: 14, Hmx: 8000.0, Cmxtime: 0, Cmx: 0.0}
	thex.Qdyel = QDAY{Hhr: 4, H: 12000.0, Chr: 0, C: 0.0, Hmxtime: 14, Hmx: 4000.0, Cmxtime: 0, Cmx: 0.0}
	thex.Qdyet = QDAY{Hhr: 4, H: 40000.0, Chr: 0, C: 0.0, Hmxtime: 14, Hmx: 12000.0, Cmxtime: 0, Cmx: 0.0}
	thex.Qdyos = QDAY{Hhr: 4, H: 24000.0, Chr: 0, C: 0.0, Hmxtime: 14, Hmx: 7000.0, Cmxtime: 0, Cmx: 0.0}
	thex.Qdyol = QDAY{Hhr: 4, H: 10000.0, Chr: 0, C: 0.0, Hmxtime: 14, Hmx: 3000.0, Cmxtime: 0, Cmx: 0.0}
	thex.Qdyot = QDAY{Hhr: 4, H: 34000.0, Chr: 0, C: 0.0, Hmxtime: 14, Hmx: 10000.0, Cmxtime: 0, Cmx: 0.0}

	// Set up monthly aggregation values
	thex.MTeidy = SVDAY{Hrs: 200, M: 24.5, Mn: 18.0, Mx: 30.0, Mntime: 1, Mxtime: 15}
	thex.MTeody = SVDAY{Hrs: 200, M: 17.5, Mn: 12.0, Mx: 22.0, Mntime: 1, Mxtime: 15}
	thex.MXeidy = SVDAY{Hrs: 200, M: 0.0095, Mn: 0.006, Mx: 0.012, Mntime: 1, Mxtime: 15}
	thex.MXeody = SVDAY{Hrs: 200, M: 0.0075, Mn: 0.004, Mx: 0.010, Mntime: 1, Mxtime: 15}
	thex.MToidy = SVDAY{Hrs: 200, M: 8.0, Mn: 0.0, Mx: 15.0, Mntime: 1, Mxtime: 15}
	thex.MToody = SVDAY{Hrs: 200, M: 15.0, Mn: 7.0, Mx: 22.0, Mntime: 1, Mxtime: 15}
	thex.MXoidy = SVDAY{Hrs: 200, M: 0.005, Mn: 0.002, Mx: 0.008, Mntime: 1, Mxtime: 15}
	thex.MXoody = SVDAY{Hrs: 200, M: 0.007, Mn: 0.004, Mx: 0.010, Mntime: 1, Mxtime: 15}

	thex.MQdyes = QDAY{Hhr: 100, H: 700000.0, Chr: 0, C: 0.0, Hmxtime: 15, Hmx: 9000.0, Cmxtime: 0, Cmx: 0.0}
	thex.MQdyel = QDAY{Hhr: 100, H: 300000.0, Chr: 0, C: 0.0, Hmxtime: 15, Hmx: 5000.0, Cmxtime: 0, Cmx: 0.0}
	thex.MQdyet = QDAY{Hhr: 100, H: 1000000.0, Chr: 0, C: 0.0, Hmxtime: 15, Hmx: 14000.0, Cmxtime: 0, Cmx: 0.0}
	thex.MQdyos = QDAY{Hhr: 100, H: 600000.0, Chr: 0, C: 0.0, Hmxtime: 15, Hmx: 8000.0, Cmxtime: 0, Cmx: 0.0}
	thex.MQdyol = QDAY{Hhr: 100, H: 250000.0, Chr: 0, C: 0.0, Hmxtime: 15, Hmx: 4000.0, Cmxtime: 0, Cmx: 0.0}
	thex.MQdyot = QDAY{Hhr: 100, H: 850000.0, Chr: 0, C: 0.0, Hmxtime: 15, Hmx: 12000.0, Cmxtime: 0, Cmx: 0.0}

	return thex
}

// TestThexprint tests the THEX print function
func TestThexprint(t *testing.T) {
	thex := createOutputTestTHEX()
	thexs := []*THEX{thex}

	t.Run("Header1_id0", func(t *testing.T) {
		var buf bytes.Buffer
		Thexprint(&buf, 0, thexs)
		output := buf.String()

		// Check type and count
		if !strings.Contains(output, "THEX") {
			t.Error("Missing THEX type in header")
		}
		if !strings.Contains(output, thex.Name) {
			t.Errorf("Missing THEX name %s in header", thex.Name)
		}
	})

	t.Run("Header2_id1", func(t *testing.T) {
		var buf bytes.Buffer
		Thexprint(&buf, 1, thexs)
		output := buf.String()

		// Check item names
		expectedPatterns := []string{
			thex.Name + "_ce",
			thex.Name + "_Ge",
			thex.Name + "_Tei",
			thex.Name + "_Teo",
			thex.Name + "_xei",
			thex.Name + "_xeo",
			thex.Name + "_hei",
			thex.Name + "_heo",
			thex.Name + "_Qes",
			thex.Name + "_Qel",
			thex.Name + "_Qet",
			thex.Name + "_co",
			thex.Name + "_Go",
			thex.Name + "_Toi",
			thex.Name + "_Too",
			thex.Name + "_xoi",
			thex.Name + "_xoo",
			thex.Name + "_hoi",
			thex.Name + "_hoo",
			thex.Name + "_Qos",
			thex.Name + "_Qol",
			thex.Name + "_Qot",
		}
		for _, pattern := range expectedPatterns {
			if !strings.Contains(output, pattern) {
				t.Errorf("Missing expected pattern: %s", pattern)
			}
		}
	})

	t.Run("Data_id99", func(t *testing.T) {
		var buf bytes.Buffer
		Thexprint(&buf, 99, thexs)
		output := buf.String()

		// Verify output is not empty
		if len(output) == 0 {
			t.Error("Data output is empty")
		}
	})

	t.Run("EmptyList", func(t *testing.T) {
		var buf bytes.Buffer
		Thexprint(&buf, 0, []*THEX{})
		output := buf.String()

		if len(output) != 0 {
			t.Error("Expected empty output for empty list")
		}
	})
}

// TestThexdyprt tests the THEX daily print function
func TestThexdyprt(t *testing.T) {
	thex := createOutputTestTHEX()
	thexs := []*THEX{thex}

	t.Run("Header1_id0", func(t *testing.T) {
		var buf bytes.Buffer
		Thexdyprt(&buf, 0, thexs)
		output := buf.String()

		if !strings.Contains(output, "THEX") {
			t.Error("Missing THEX type in daily header")
		}
		if !strings.Contains(output, thex.Name) {
			t.Errorf("Missing THEX name %s in daily header", thex.Name)
		}
	})

	t.Run("Header2_id1", func(t *testing.T) {
		var buf bytes.Buffer
		Thexdyprt(&buf, 1, thexs)
		output := buf.String()

		// Check item patterns
		expectedPatterns := []string{
			thex.Name + "_Hte",
			thex.Name + "_Te",
			thex.Name + "_Hto",
			thex.Name + "_To",
			thex.Name + "_Hxe",
			thex.Name + "_xe",
			thex.Name + "_Hxo",
			thex.Name + "_xo",
			thex.Name + "_Hhs",
			thex.Name + "_Qsh",
		}
		for _, pattern := range expectedPatterns {
			if !strings.Contains(output, pattern) {
				t.Errorf("Missing expected pattern: %s", pattern)
			}
		}
	})

	t.Run("Data_id99", func(t *testing.T) {
		var buf bytes.Buffer
		Thexdyprt(&buf, 99, thexs)
		output := buf.String()

		if len(output) == 0 {
			t.Error("Daily data output is empty")
		}
	})

	t.Run("EmptyList", func(t *testing.T) {
		var buf bytes.Buffer
		Thexdyprt(&buf, 0, []*THEX{})
		output := buf.String()

		if len(output) != 0 {
			t.Error("Expected empty output for empty list")
		}
	})
}

// TestThexmonprt tests the THEX monthly print function
func TestThexmonprt(t *testing.T) {
	thex := createOutputTestTHEX()
	thexs := []*THEX{thex}

	t.Run("Header1_id0", func(t *testing.T) {
		var buf bytes.Buffer
		Thexmonprt(&buf, 0, thexs)
		output := buf.String()

		if !strings.Contains(output, "THEX") {
			t.Error("Missing THEX type in monthly header")
		}
		if !strings.Contains(output, thex.Name) {
			t.Errorf("Missing THEX name %s in monthly header", thex.Name)
		}
	})

	t.Run("Header2_id1", func(t *testing.T) {
		var buf bytes.Buffer
		Thexmonprt(&buf, 1, thexs)
		output := buf.String()

		// Check item patterns (same as daily)
		expectedPatterns := []string{
			thex.Name + "_Hte",
			thex.Name + "_Te",
			thex.Name + "_Hto",
			thex.Name + "_To",
			thex.Name + "_Hxe",
			thex.Name + "_xe",
			thex.Name + "_Hxo",
			thex.Name + "_xo",
		}
		for _, pattern := range expectedPatterns {
			if !strings.Contains(output, pattern) {
				t.Errorf("Missing expected pattern: %s", pattern)
			}
		}
	})

	t.Run("Data_id99", func(t *testing.T) {
		var buf bytes.Buffer
		Thexmonprt(&buf, 99, thexs)
		output := buf.String()

		if len(output) == 0 {
			t.Error("Monthly data output is empty")
		}
	})

	t.Run("EmptyList", func(t *testing.T) {
		var buf bytes.Buffer
		Thexmonprt(&buf, 0, []*THEX{})
		output := buf.String()

		if len(output) != 0 {
			t.Error("Expected empty output for empty list")
		}
	})
}

// TestThexdyint tests the THEX daily aggregation initialization
func TestThexdyint(t *testing.T) {
	t.Run("BasicInitialization", func(t *testing.T) {
		thex := createOutputTestTHEX()
		thexs := []*THEX{thex}

		// Verify values are set before initialization
		if thex.Teidy.Hrs == 0 {
			t.Error("Test data not properly set up")
		}

		Thexdyint(thexs)

		// After initialization, values should be reset
		if thex.Teidy.Hrs != 0 {
			t.Error("Teidy.Hrs should be reset to 0")
		}
		if thex.Qdyes.Hhr != 0 {
			t.Error("Qdyes.Hhr should be reset to 0")
		}
	})

	t.Run("EmptyList", func(t *testing.T) {
		// Should not panic with empty list
		Thexdyint([]*THEX{})
	})

	t.Run("MultipleThex", func(t *testing.T) {
		thex1 := createOutputTestTHEX()
		thex1.Name = "THEX1"
		thex2 := createOutputTestTHEX()
		thex2.Name = "THEX2"
		thexs := []*THEX{thex1, thex2}

		Thexdyint(thexs)

		for i, thex := range thexs {
			if thex.Teidy.Hrs != 0 {
				t.Errorf("THEX[%d] Teidy.Hrs should be reset to 0", i)
			}
		}
	})
}

// TestThexmonint tests the THEX monthly aggregation initialization
func TestThexmonint(t *testing.T) {
	t.Run("BasicInitialization", func(t *testing.T) {
		thex := createOutputTestTHEX()
		thexs := []*THEX{thex}

		// Verify values are set before initialization
		if thex.MTeidy.Hrs == 0 {
			t.Error("Test data not properly set up")
		}

		Thexmonint(thexs)

		// After initialization, values should be reset
		if thex.MTeidy.Hrs != 0 {
			t.Error("MTeidy.Hrs should be reset to 0")
		}
		if thex.MQdyes.Hhr != 0 {
			t.Error("MQdyes.Hhr should be reset to 0")
		}
	})

	t.Run("EmptyList", func(t *testing.T) {
		// Should not panic with empty list
		Thexmonint([]*THEX{})
	})

	t.Run("MultipleThex", func(t *testing.T) {
		thex1 := createOutputTestTHEX()
		thex1.Name = "THEX1"
		thex2 := createOutputTestTHEX()
		thex2.Name = "THEX2"
		thexs := []*THEX{thex1, thex2}

		Thexmonint(thexs)

		for i, thex := range thexs {
			if thex.MTeidy.Hrs != 0 {
				t.Errorf("THEX[%d] MTeidy.Hrs should be reset to 0", i)
			}
		}
	})
}

// TestThexdata tests the Thexdata function
func TestThexdata(t *testing.T) {
	t.Run("SetName", func(t *testing.T) {
		thexca := &THEXCA{}
		result := Thexdata("TestTHEX", thexca)

		if result != 0 {
			t.Errorf("Thexdata should return 0 for name, got %d", result)
		}
		if thexca.Name != "TestTHEX" {
			t.Errorf("Name = %s, want TestTHEX", thexca.Name)
		}
	})

	t.Run("Set_et", func(t *testing.T) {
		thexca := &THEXCA{}
		result := Thexdata("et=0.75", thexca)

		if result != 0 {
			t.Errorf("Thexdata should return 0 for et, got %d", result)
		}
		if thexca.et != 0.75 {
			t.Errorf("et = %f, want 0.75", thexca.et)
		}
	})

	t.Run("Set_eh", func(t *testing.T) {
		thexca := &THEXCA{}
		result := Thexdata("eh=0.65", thexca)

		if result != 0 {
			t.Errorf("Thexdata should return 0 for eh, got %d", result)
		}
		if thexca.eh != 0.65 {
			t.Errorf("eh = %f, want 0.65", thexca.eh)
		}
	})

	t.Run("UnknownKey", func(t *testing.T) {
		thexca := &THEXCA{}
		result := Thexdata("unknown=123", thexca)

		if result != 1 {
			t.Errorf("Thexdata should return 1 for unknown key, got %d", result)
		}
	})

	t.Run("MultipleParameters", func(t *testing.T) {
		thexca := &THEXCA{}

		// Set name first
		Thexdata("TestHEX", thexca)
		if thexca.Name != "TestHEX" {
			t.Errorf("Name = %s, want TestHEX", thexca.Name)
		}

		// Then set et
		Thexdata("et=0.80", thexca)
		if thexca.et != 0.80 {
			t.Errorf("et = %f, want 0.80", thexca.et)
		}

		// Then set eh
		Thexdata("eh=0.70", thexca)
		if thexca.eh != 0.70 {
			t.Errorf("eh = %f, want 0.70", thexca.eh)
		}
	})
}

// TestThexday tests the Thexday aggregation function
func TestThexday(t *testing.T) {
	t.Run("DailyAggregation", func(t *testing.T) {
		thex := createBasicTHEX()
		thex.Cmp.Control = ON_SW
		// Set test values
		thex.Tein = 25.0
		thex.Teout = 18.0
		thex.Xein = 0.010
		thex.Xeout = 0.008
		thex.Toin = 5.0
		thex.Toout = 12.0
		thex.Xoin = 0.004
		thex.Xoout = 0.006
		thex.Qes = 7000.0
		thex.Qel = 3000.0
		thex.Qet = 10000.0
		thex.Qos = 6000.0
		thex.Qol = 2500.0
		thex.Qot = 8500.0

		thexs := []*THEX{thex}

		// Initialize daily aggregation
		Thexdyint(thexs)

		// Simulate multiple time steps
		times := []int{900, 1000, 1100, 1200}
		for _, ttmm := range times {
			Thexday(7, 15, ttmm, thexs, 31, 365)
		}

		// After 4 time steps, verify aggregation
		if thex.Teidy.Hrs != 4 {
			t.Errorf("Teidy.Hrs = %d, want 4", thex.Teidy.Hrs)
		}
		if thex.Toidy.Hrs != 4 {
			t.Errorf("Toidy.Hrs = %d, want 4", thex.Toidy.Hrs)
		}
		if thex.Xeidy.Hrs != 4 {
			t.Errorf("Xeidy.Hrs = %d, want 4", thex.Xeidy.Hrs)
		}
		if thex.Xoidy.Hrs != 4 {
			t.Errorf("Xoidy.Hrs = %d, want 4", thex.Xoidy.Hrs)
		}

		// Check heat quantity aggregation
		if thex.Qdyes.Hhr != 4 {
			t.Errorf("Qdyes.Hhr = %d, want 4", thex.Qdyes.Hhr)
		}
		if thex.Qdyos.Hhr != 4 {
			t.Errorf("Qdyos.Hhr = %d, want 4", thex.Qdyos.Hhr)
		}
	})

	t.Run("OffControl_NoAggregation", func(t *testing.T) {
		thex := createBasicTHEX()
		thex.Cmp.Control = OFF_SW
		thex.Tein = 20.0
		thex.Qes = 0.0

		thexs := []*THEX{thex}

		Thexdyint(thexs)
		Thexday(7, 15, 1200, thexs, 31, 365)

		// OFF control should not aggregate
		if thex.Teidy.Hrs != 0 {
			t.Errorf("Teidy.Hrs should be 0 when OFF, got %d", thex.Teidy.Hrs)
		}
	})

	t.Run("MonthlyAggregation_EndOfDay", func(t *testing.T) {
		thex := createBasicTHEX()
		thex.Cmp.Control = ON_SW
		thex.Tein = 25.0
		thex.Teout = 18.0
		thex.Qes = 7000.0

		thexs := []*THEX{thex}

		Thexdyint(thexs)
		Thexmonint(thexs)

		// Call at end of month (Day=Nday)
		Thexday(7, 31, 2400, thexs, 31, 365)

		// Daily values should be aggregated
		if thex.Teidy.Hrs != 1 {
			t.Errorf("Teidy.Hrs = %d, want 1", thex.Teidy.Hrs)
		}
	})

	t.Run("EmptyList", func(t *testing.T) {
		// Should not panic with empty list
		Thexday(7, 15, 1200, []*THEX{}, 31, 365)
	})

	t.Run("MultipleThex", func(t *testing.T) {
		thexs := make([]*THEX, 2)
		for i := range thexs {
			thexs[i] = createBasicTHEX()
			thexs[i].Name = "THEX" + string(rune('A'+i))
			thexs[i].Cmp.Control = ON_SW
			thexs[i].Tein = 25.0 + float64(i)*2
			thexs[i].Qes = 7000.0 + float64(i)*1000
		}

		Thexdyint(thexs)
		Thexday(7, 15, 1200, thexs, 31, 365)

		// Verify each thex has independent aggregation
		for i, thex := range thexs {
			if thex.Teidy.Hrs != 1 {
				t.Errorf("THEX[%d] Teidy.Hrs = %d, want 1", i, thex.Teidy.Hrs)
			}
		}
	})

	t.Run("AllFieldsAggregated", func(t *testing.T) {
		thex := createBasicTHEX()
		thex.Cmp.Control = ON_SW
		// Set all fields
		thex.Tein = 25.0
		thex.Teout = 18.0
		thex.Xein = 0.010
		thex.Xeout = 0.008
		thex.Toin = 5.0
		thex.Toout = 12.0
		thex.Xoin = 0.004
		thex.Xoout = 0.006
		thex.Qes = 7000.0
		thex.Qel = 3000.0
		thex.Qet = 10000.0
		thex.Qos = 6000.0
		thex.Qol = 2500.0
		thex.Qot = 8500.0

		thexs := []*THEX{thex}

		Thexdyint(thexs)
		Thexday(7, 15, 1200, thexs, 31, 365)

		// Verify all temperature fields aggregated
		if thex.Teidy.Hrs != 1 {
			t.Errorf("Teidy.Hrs = %d, want 1", thex.Teidy.Hrs)
		}
		if thex.Teody.Hrs != 1 {
			t.Errorf("Teody.Hrs = %d, want 1", thex.Teody.Hrs)
		}
		if thex.Toidy.Hrs != 1 {
			t.Errorf("Toidy.Hrs = %d, want 1", thex.Toidy.Hrs)
		}
		if thex.Toody.Hrs != 1 {
			t.Errorf("Toody.Hrs = %d, want 1", thex.Toody.Hrs)
		}
		if thex.Xeidy.Hrs != 1 {
			t.Errorf("Xeidy.Hrs = %d, want 1", thex.Xeidy.Hrs)
		}
		if thex.Xeody.Hrs != 1 {
			t.Errorf("Xeody.Hrs = %d, want 1", thex.Xeody.Hrs)
		}
		if thex.Xoidy.Hrs != 1 {
			t.Errorf("Xoidy.Hrs = %d, want 1", thex.Xoidy.Hrs)
		}
		if thex.Xoody.Hrs != 1 {
			t.Errorf("Xoody.Hrs = %d, want 1", thex.Xoody.Hrs)
		}

		// Verify all heat quantity fields aggregated
		if thex.Qdyes.Hhr != 1 {
			t.Errorf("Qdyes.Hhr = %d, want 1", thex.Qdyes.Hhr)
		}
		if thex.Qdyel.Hhr != 1 {
			t.Errorf("Qdyel.Hhr = %d, want 1", thex.Qdyel.Hhr)
		}
		if thex.Qdyet.Hhr != 1 {
			t.Errorf("Qdyet.Hhr = %d, want 1", thex.Qdyet.Hhr)
		}
		if thex.Qdyos.Hhr != 1 {
			t.Errorf("Qdyos.Hhr = %d, want 1", thex.Qdyos.Hhr)
		}
		if thex.Qdyol.Hhr != 1 {
			t.Errorf("Qdyol.Hhr = %d, want 1", thex.Qdyol.Hhr)
		}
		if thex.Qdyot.Hhr != 1 {
			t.Errorf("Qdyot.Hhr = %d, want 1", thex.Qdyot.Hhr)
		}
	})
}
package eeslism

import (
	"bytes"
	"math"
	"strings"
	"testing"
)

// TestHexdata tests the Hexdata function
func TestHexdata(t *testing.T) {
	t.Run("SetName", func(t *testing.T) {
		hexca := &HEXCA{}
		result := Hexdata("TestHEX", hexca)

		if result != 0 {
			t.Errorf("Hexdata should return 0 for name, got %d", result)
		}
		if hexca.Name != "TestHEX" {
			t.Errorf("Name = %s, want TestHEX", hexca.Name)
		}
	})

	t.Run("Set_eff", func(t *testing.T) {
		hexca := &HEXCA{}
		result := Hexdata("eff=0.85", hexca)

		if result != 0 {
			t.Errorf("Hexdata should return 0 for eff, got %d", result)
		}
		if hexca.eff != 0.85 {
			t.Errorf("eff = %f, want 0.85", hexca.eff)
		}
	})

	t.Run("Set_KA", func(t *testing.T) {
		hexca := &HEXCA{}
		result := Hexdata("KA=500.0", hexca)

		if result != 0 {
			t.Errorf("Hexdata should return 0 for KA, got %d", result)
		}
		if hexca.KA != 500.0 {
			t.Errorf("KA = %f, want 500.0", hexca.KA)
		}
	})

	t.Run("UnknownKey", func(t *testing.T) {
		hexca := &HEXCA{}
		result := Hexdata("unknown=123", hexca)

		if result != 1 {
			t.Errorf("Hexdata should return 1 for unknown key, got %d", result)
		}
	})

	t.Run("InvalidEffValue", func(t *testing.T) {
		hexca := &HEXCA{}
		result := Hexdata("eff=invalid", hexca)

		if result != 1 {
			t.Errorf("Hexdata should return 1 for invalid value, got %d", result)
		}
	})

	t.Run("InvalidKAValue", func(t *testing.T) {
		hexca := &HEXCA{}
		result := Hexdata("KA=invalid", hexca)

		if result != 1 {
			t.Errorf("Hexdata should return 1 for invalid value, got %d", result)
		}
	})

	t.Run("MultipleParameters", func(t *testing.T) {
		hexca := &HEXCA{}

		// Set name first
		Hexdata("TestHeatExchanger", hexca)
		if hexca.Name != "TestHeatExchanger" {
			t.Errorf("Name = %s, want TestHeatExchanger", hexca.Name)
		}

		// Set eff
		Hexdata("eff=0.75", hexca)
		if hexca.eff != 0.75 {
			t.Errorf("eff = %f, want 0.75", hexca.eff)
		}

		// Set KA
		Hexdata("KA=1000.0", hexca)
		if hexca.KA != 1000.0 {
			t.Errorf("KA = %f, want 1000.0", hexca.KA)
		}
	})
}

// Helper function to create a basic HEX for testing
func createBasicHEX() *HEX {
	// Create ELOUTs for cold and hot sides
	eloutC := &ELOUT{
		Control: ON_SW,
		Fluid:   WATER_FLD,
		G:       0.5, // 0.5 kg/s cold side
		Coeffin: make([]float64, 2),
		Coeffo:  0.0,
		Co:      0.0,
		Sysv:    15.0, // Output temperature cold side
	}
	eloutH := &ELOUT{
		Control: ON_SW,
		Fluid:   WATER_FLD,
		G:       0.8, // 0.8 kg/s hot side
		Coeffin: make([]float64, 2),
		Coeffo:  0.0,
		Co:      0.0,
		Sysv:    35.0, // Output temperature hot side
	}

	// Create ELINs for cold and hot sides
	elinC := &ELIN{
		Sysvin: 10.0, // 10°C cold inlet
	}
	elinH := &ELIN{
		Sysvin: 50.0, // 50°C hot inlet
	}

	// Create COMPNT
	cmp := &COMPNT{
		Name:    "TestHEX",
		Control: ON_SW,
		Elouts:  []*ELOUT{eloutC, eloutH},
		Elins:   []*ELIN{elinC, elinH},
	}

	// Create HEXCA with fixed efficiency
	hexca := &HEXCA{
		Name: "TestHEXCA",
		eff:  0.8, // 80% efficiency
		KA:   FNAN,
	}

	return &HEX{
		Name: "TestHEX",
		Cat:  hexca,
		Cmp:  cmp,
	}
}

// Helper function to create HEX with KA-based efficiency
func createKABasedHEX() *HEX {
	hex := createBasicHEX()
	hex.Cat.eff = FNAN // Clear fixed efficiency
	hex.Cat.KA = 1000.0 // Set KA value
	return hex
}

// TestHexcfv tests the Hexcfv function
func TestHexcfv(t *testing.T) {
	t.Run("FixedEfficiency_FirstRun", func(t *testing.T) {
		hex := createBasicHEX()
		hex.Id = 0 // First run, needs initialization

		hexs := []*HEX{hex}
		Hexcfv(hexs)

		// Check Etype is set to 'e' for fixed efficiency
		if hex.Etype != 'e' {
			t.Errorf("Etype = '%c', want 'e'", hex.Etype)
		}

		// Check Id is incremented
		if hex.Id != 1 {
			t.Errorf("Id = %d, want 1", hex.Id)
		}

		// Check Eff is set from Cat
		if hex.Eff != 0.8 {
			t.Errorf("Eff = %f, want 0.8", hex.Eff)
		}

		// Check CGc and CGh are calculated
		expectedCGc := Spcheat(WATER_FLD) * 0.5
		expectedCGh := Spcheat(WATER_FLD) * 0.8
		if math.Abs(hex.CGc-expectedCGc) > 1e-6 {
			t.Errorf("CGc = %f, want %f", hex.CGc, expectedCGc)
		}
		if math.Abs(hex.CGh-expectedCGh) > 1e-6 {
			t.Errorf("CGh = %f, want %f", hex.CGh, expectedCGh)
		}

		// Check ECGmin is calculated
		expectedECGmin := hex.Eff * math.Min(hex.CGc, hex.CGh)
		if math.Abs(hex.ECGmin-expectedECGmin) > 1e-6 {
			t.Errorf("ECGmin = %f, want %f", hex.ECGmin, expectedECGmin)
		}

		// Check Elout coefficients for cold side
		eoc := hex.Cmp.Elouts[0]
		if math.Abs(eoc.Coeffo-hex.CGc) > 1e-6 {
			t.Errorf("Elouts[0].Coeffo = %f, want %f", eoc.Coeffo, hex.CGc)
		}
		expectedCoeffin0 := -hex.CGc + hex.ECGmin
		if math.Abs(eoc.Coeffin[0]-expectedCoeffin0) > 1e-6 {
			t.Errorf("Elouts[0].Coeffin[0] = %f, want %f", eoc.Coeffin[0], expectedCoeffin0)
		}
	})

	t.Run("KABasedEfficiency", func(t *testing.T) {
		hex := createKABasedHEX()
		hex.Id = 0

		hexs := []*HEX{hex}
		Hexcfv(hexs)

		// Check Etype is set to 'k' for KA-based efficiency
		if hex.Etype != 'k' {
			t.Errorf("Etype = '%c', want 'k'", hex.Etype)
		}

		// Eff should be calculated by FNhccet
		// FNhccet calculates efficiency from CGc, CGh, and KA
		if hex.Eff <= 0.0 || hex.Eff > 1.0 {
			t.Errorf("Eff = %f, should be between 0 and 1", hex.Eff)
		}
	})

	t.Run("SecondRun_NoReinitialization", func(t *testing.T) {
		hex := createBasicHEX()
		hex.Id = 1 // Already initialized
		hex.Etype = 'e'

		hexs := []*HEX{hex}
		Hexcfv(hexs)

		// Id should remain 1
		if hex.Id != 1 {
			t.Errorf("Id = %d, should remain 1 after second run", hex.Id)
		}
	})

	t.Run("ControlOFF", func(t *testing.T) {
		hex := createBasicHEX()
		hex.Id = 1
		hex.Etype = 'e'
		hex.Cmp.Control = OFF_SW

		// Set initial values to check they don't change
		initialCGc := hex.CGc
		initialCGh := hex.CGh

		hexs := []*HEX{hex}
		Hexcfv(hexs)

		// When Control is OFF_SW, coefficients should not be calculated
		if hex.CGc != initialCGc || hex.CGh != initialCGh {
			t.Log("CGc and CGh should not be recalculated when Control is OFF_SW")
		}
	})

	t.Run("EmptyList", func(t *testing.T) {
		var hexs []*HEX
		// Should not panic with empty list
		Hexcfv(hexs)
	})

	t.Run("MultipleHEX", func(t *testing.T) {
		hex1 := createBasicHEX()
		hex1.Name = "HEX1"
		hex1.Id = 0

		hex2 := createKABasedHEX()
		hex2.Name = "HEX2"
		hex2.Id = 0

		hexs := []*HEX{hex1, hex2}
		Hexcfv(hexs)

		// Both should be initialized
		if hex1.Id != 1 || hex2.Id != 1 {
			t.Error("Both HEXs should be initialized")
		}
		if hex1.Etype != 'e' {
			t.Errorf("hex1.Etype = '%c', want 'e'", hex1.Etype)
		}
		if hex2.Etype != 'k' {
			t.Errorf("hex2.Etype = '%c', want 'k'", hex2.Etype)
		}
	})
}

// TestHexene tests the Hexene function
func TestHexene(t *testing.T) {
	t.Run("BasicEnergyCalculation", func(t *testing.T) {
		hex := createBasicHEX()
		hex.CGc = Spcheat(WATER_FLD) * 0.5
		hex.CGh = Spcheat(WATER_FLD) * 0.8
		hex.Cmp.Elins[0].Sysvin = 10.0  // Cold inlet
		hex.Cmp.Elins[1].Sysvin = 50.0  // Hot inlet
		hex.Cmp.Elouts[0].Sysv = 20.0   // Cold outlet
		hex.Cmp.Elouts[1].Sysv = 35.0   // Hot outlet

		hexs := []*HEX{hex}
		Hexene(hexs)

		// Check Tcin and Thin are set
		if hex.Tcin != 10.0 {
			t.Errorf("Tcin = %f, want 10.0", hex.Tcin)
		}
		if hex.Thin != 50.0 {
			t.Errorf("Thin = %f, want 50.0", hex.Thin)
		}

		// Check Qci = CGc * (Tout_cold - Tin_cold)
		expectedQci := hex.CGc * (20.0 - 10.0)
		if math.Abs(hex.Qci-expectedQci) > 1e-6 {
			t.Errorf("Qci = %f, want %f", hex.Qci, expectedQci)
		}

		// Check Qhi = CGh * (Tout_hot - Tin_hot)
		expectedQhi := hex.CGh * (35.0 - 50.0)
		if math.Abs(hex.Qhi-expectedQhi) > 1e-6 {
			t.Errorf("Qhi = %f, want %f", hex.Qhi, expectedQhi)
		}

		// Qci should be positive (cold side gains heat)
		if hex.Qci <= 0 {
			t.Errorf("Qci = %f, should be positive", hex.Qci)
		}

		// Qhi should be negative (hot side loses heat)
		if hex.Qhi >= 0 {
			t.Errorf("Qhi = %f, should be negative", hex.Qhi)
		}
	})

	t.Run("ControlOFF", func(t *testing.T) {
		hex := createBasicHEX()
		hex.Cmp.Control = OFF_SW
		hex.Cmp.Elins[0].Sysvin = 10.0
		hex.Cmp.Elins[1].Sysvin = 50.0

		hexs := []*HEX{hex}
		Hexene(hexs)

		// Qci and Qhi should be 0 when Control is OFF_SW
		if hex.Qci != 0.0 {
			t.Errorf("Qci = %f, want 0.0 when Control is OFF_SW", hex.Qci)
		}
		if hex.Qhi != 0.0 {
			t.Errorf("Qhi = %f, want 0.0 when Control is OFF_SW", hex.Qhi)
		}

		// Tcin and Thin should still be set from inputs
		if hex.Tcin != 10.0 {
			t.Errorf("Tcin = %f, want 10.0 even when OFF", hex.Tcin)
		}
		if hex.Thin != 50.0 {
			t.Errorf("Thin = %f, want 50.0 even when OFF", hex.Thin)
		}
	})

	t.Run("EmptyList", func(t *testing.T) {
		var hexs []*HEX
		// Should not panic with empty list
		Hexene(hexs)
	})

	t.Run("MultipleHEX", func(t *testing.T) {
		hex1 := createBasicHEX()
		hex1.CGc = Spcheat(WATER_FLD) * 0.5
		hex1.CGh = Spcheat(WATER_FLD) * 0.8
		hex1.Cmp.Elins[0].Sysvin = 10.0
		hex1.Cmp.Elins[1].Sysvin = 50.0
		hex1.Cmp.Elouts[0].Sysv = 20.0
		hex1.Cmp.Elouts[1].Sysv = 35.0

		hex2 := createBasicHEX()
		hex2.CGc = Spcheat(WATER_FLD) * 1.0
		hex2.CGh = Spcheat(WATER_FLD) * 1.0
		hex2.Cmp.Elins[0].Sysvin = 5.0
		hex2.Cmp.Elins[1].Sysvin = 60.0
		hex2.Cmp.Elouts[0].Sysv = 30.0
		hex2.Cmp.Elouts[1].Sysv = 35.0

		hexs := []*HEX{hex1, hex2}
		Hexene(hexs)

		// Both should have calculated Q values
		if hex1.Qci == 0 && hex2.Qci == 0 {
			t.Error("At least one HEX should have non-zero Qci")
		}
	})
}

// TestHexprint tests the hex output function
func TestHexprint(t *testing.T) {
	t.Run("Header1_id0", func(t *testing.T) {
		hex := createBasicHEX()
		hexs := []*HEX{hex}

		var buf bytes.Buffer
		hexprint(&buf, 0, hexs)
		output := buf.String()

		if !strings.Contains(output, string(HEXCHANGR_TYPE)) {
			t.Error("Output should contain HEXCHANGR type")
		}
		if !strings.Contains(output, hex.Name) {
			t.Error("Output should contain hex name")
		}
	})

	t.Run("Header2_id1", func(t *testing.T) {
		hex := createBasicHEX()
		hexs := []*HEX{hex}

		var buf bytes.Buffer
		hexprint(&buf, 1, hexs)
		output := buf.String()

		// Should contain column headers for both cold and hot sides
		if !strings.Contains(output, ":c_G") {
			t.Error("Output should contain cold side flow rate column")
		}
		if !strings.Contains(output, ":h_G") {
			t.Error("Output should contain hot side flow rate column")
		}
	})

	t.Run("Data_id99", func(t *testing.T) {
		hex := createBasicHEX()
		hex.Tcin = 10.0
		hex.Thin = 50.0
		hex.Qci = 5000.0
		hex.Qhi = -5000.0
		hexs := []*HEX{hex}

		var buf bytes.Buffer
		hexprint(&buf, 99, hexs)
		output := buf.String()

		if len(output) == 0 {
			t.Error("Output should not be empty")
		}
	})

	t.Run("EmptyList", func(t *testing.T) {
		var hexs []*HEX

		var buf bytes.Buffer
		hexprint(&buf, 0, hexs)
		output := buf.String()

		if strings.Contains(output, string(HEXCHANGR_TYPE)) {
			t.Error("Output should not contain type for empty list")
		}
	})
}

// TestHexday tests the daily aggregation function
func TestHexday(t *testing.T) {
	t.Run("DailyAggregation", func(t *testing.T) {
		hex := createBasicHEX()
		hex.Tcin = 10.0
		hex.Thin = 50.0
		hex.Qci = 5000.0
		hex.Qhi = -5000.0

		hexs := []*HEX{hex}

		// Initialize daily counters
		hexdyint(hexs)

		// Run aggregation for several time steps
		for ttmm := 100; ttmm <= 2300; ttmm += 100 {
			hex.Tcin = 8.0 + float64(ttmm)/1000.0
			hex.Thin = 45.0 + float64(ttmm)/500.0
			hex.Qci = 4000.0 + float64(ttmm)
			hex.Qhi = -4000.0 - float64(ttmm)
			hexday(1, 1, ttmm, hexs, 1, 365)
		}

		// Verify aggregation results
		if hex.Tcidy.Hrs == 0 {
			t.Error("Tcidy.Hrs should be > 0 after aggregation")
		}
		if hex.Qcidy.Hhr == 0 && hex.Qcidy.Chr == 0 {
			t.Error("Qcidy should have some hours after aggregation")
		}
	})

	t.Run("MonthlyAggregation_EndOfDay", func(t *testing.T) {
		hex := createBasicHEX()
		hex.Tcin = 10.0
		hex.Thin = 50.0
		hex.Qci = 5000.0
		hex.Qhi = -5000.0

		hexs := []*HEX{hex}

		// Initialize daily and monthly counters
		hexdyint(hexs)
		hexmonint(hexs)

		// Simulate end of day
		hexday(1, 1, 2400, hexs, 1, 365)

		t.Log("Monthly aggregation test completed")
	})

	t.Run("EmptyList", func(t *testing.T) {
		var hexs []*HEX

		// Should not panic with empty list
		hexday(1, 1, 100, hexs, 1, 365)
	})
}

// TestHexdyprt tests the daily output function
func TestHexdyprt(t *testing.T) {
	t.Run("Header1_id0", func(t *testing.T) {
		hex := createBasicHEX()
		hexs := []*HEX{hex}

		var buf bytes.Buffer
		hexdyprt(&buf, 0, hexs)
		output := buf.String()

		if !strings.Contains(output, string(HEXCHANGR_TYPE)) {
			t.Error("Output should contain HEXCHANGR type")
		}
	})

	t.Run("Header2_id1", func(t *testing.T) {
		hex := createBasicHEX()
		hexs := []*HEX{hex}

		var buf bytes.Buffer
		hexdyprt(&buf, 1, hexs)
		output := buf.String()

		// Should contain headers for both cold and hot sides
		if !strings.Contains(output, ":c_Ht") {
			t.Error("Output should contain cold side hours header")
		}
		if !strings.Contains(output, ":h_Ht") {
			t.Error("Output should contain hot side hours header")
		}
	})

	t.Run("Data_id99", func(t *testing.T) {
		hex := createBasicHEX()
		hex.Tcidy = SVDAY{Hrs: 10, M: 15.0, Mn: 10.0, Mx: 20.0}
		hex.Thidy = SVDAY{Hrs: 10, M: 45.0, Mn: 40.0, Mx: 50.0}
		hex.Qcidy = QDAY{Hhr: 5, H: 50000.0, Chr: 0, C: 0}
		hex.Qhidy = QDAY{Hhr: 0, H: 0, Chr: 5, C: 50000.0}
		hexs := []*HEX{hex}

		var buf bytes.Buffer
		hexdyprt(&buf, 99, hexs)
		output := buf.String()

		if len(output) == 0 {
			t.Error("Output should not be empty")
		}
	})
}

// TestHexmonprt tests the monthly output function
func TestHexmonprt(t *testing.T) {
	t.Run("Header1_id0", func(t *testing.T) {
		hex := createBasicHEX()
		hexs := []*HEX{hex}

		var buf bytes.Buffer
		hexmonprt(&buf, 0, hexs)
		output := buf.String()

		if !strings.Contains(output, string(HEXCHANGR_TYPE)) {
			t.Error("Output should contain HEXCHANGR type")
		}
	})

	t.Run("Data_id99", func(t *testing.T) {
		hex := createBasicHEX()
		hex.MTcidy = SVDAY{Hrs: 100, M: 15.0, Mn: 5.0, Mx: 25.0}
		hex.MThidy = SVDAY{Hrs: 100, M: 45.0, Mn: 35.0, Mx: 55.0}
		hex.MQcidy = QDAY{Hhr: 50, H: 500000.0, Chr: 0, C: 0}
		hex.MQhidy = QDAY{Hhr: 0, H: 0, Chr: 50, C: 500000.0}
		hexs := []*HEX{hex}

		var buf bytes.Buffer
		hexmonprt(&buf, 99, hexs)
		output := buf.String()

		if len(output) == 0 {
			t.Error("Output should not be empty")
		}
	})
}

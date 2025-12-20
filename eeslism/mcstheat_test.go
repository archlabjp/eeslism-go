package eeslism

import (
	"bytes"
	"math"
	"strings"
	"testing"
)

// TestStheatdata tests the Stheatdata function
func TestStheatdata(t *testing.T) {
	t.Run("SetName", func(t *testing.T) {
		stheatca := &STHEATCA{}
		result := Stheatdata("TestSTHEAT", stheatca)

		if result != 0 {
			t.Errorf("Stheatdata should return 0 for name, got %d", result)
		}
		if stheatca.Name != "TestSTHEAT" {
			t.Errorf("Name = %s, want TestSTHEAT", stheatca.Name)
		}
		// Initial values should be FNAN
		if stheatca.Eff != FNAN {
			t.Errorf("Eff should be FNAN, got %f", stheatca.Eff)
		}
	})

	t.Run("Set_Q", func(t *testing.T) {
		stheatca := &STHEATCA{}
		result := Stheatdata("Q=5000.0", stheatca)

		if result != 0 {
			t.Errorf("Stheatdata should return 0 for Q, got %d", result)
		}
		if stheatca.Q != 5000.0 {
			t.Errorf("Q = %f, want 5000.0", stheatca.Q)
		}
	})

	t.Run("Set_eff", func(t *testing.T) {
		stheatca := &STHEATCA{}
		result := Stheatdata("eff=0.95", stheatca)

		if result != 0 {
			t.Errorf("Stheatdata should return 0 for eff, got %d", result)
		}
		if stheatca.Eff != 0.95 {
			t.Errorf("Eff = %f, want 0.95", stheatca.Eff)
		}
	})

	t.Run("Set_Hcap", func(t *testing.T) {
		stheatca := &STHEATCA{}
		result := Stheatdata("Hcap=50000.0", stheatca)

		if result != 0 {
			t.Errorf("Stheatdata should return 0 for Hcap, got %d", result)
		}
		if stheatca.Hcap != 50000.0 {
			t.Errorf("Hcap = %f, want 50000.0", stheatca.Hcap)
		}
	})

	t.Run("Set_KA", func(t *testing.T) {
		stheatca := &STHEATCA{}
		result := Stheatdata("KA=10.0", stheatca)

		if result != 0 {
			t.Errorf("Stheatdata should return 0 for KA, got %d", result)
		}
		if stheatca.KA != 10.0 {
			t.Errorf("KA = %f, want 10.0", stheatca.KA)
		}
	})

	t.Run("Set_PCM", func(t *testing.T) {
		stheatca := &STHEATCA{}
		result := Stheatdata("PCM=TestPCM", stheatca)

		if result != 0 {
			t.Errorf("Stheatdata should return 0 for PCM, got %d", result)
		}
		if stheatca.PCMName != "TestPCM" {
			t.Errorf("PCMName = %s, want TestPCM", stheatca.PCMName)
		}
	})

	t.Run("UnknownKey", func(t *testing.T) {
		stheatca := &STHEATCA{}
		result := Stheatdata("unknown=123", stheatca)

		if result != 1 {
			t.Errorf("Stheatdata should return 1 for unknown key, got %d", result)
		}
	})

	t.Run("MultipleParameters", func(t *testing.T) {
		stheatca := &STHEATCA{}

		// Set name first
		Stheatdata("TestHeater", stheatca)
		if stheatca.Name != "TestHeater" {
			t.Errorf("Name = %s, want TestHeater", stheatca.Name)
		}

		// Set Q
		Stheatdata("Q=3000.0", stheatca)
		if stheatca.Q != 3000.0 {
			t.Errorf("Q = %f, want 3000.0", stheatca.Q)
		}

		// Set eff
		Stheatdata("eff=0.9", stheatca)
		if stheatca.Eff != 0.9 {
			t.Errorf("Eff = %f, want 0.9", stheatca.Eff)
		}

		// Set Hcap
		Stheatdata("Hcap=20000.0", stheatca)
		if stheatca.Hcap != 20000.0 {
			t.Errorf("Hcap = %f, want 20000.0", stheatca.Hcap)
		}

		// Set KA
		Stheatdata("KA=5.0", stheatca)
		if stheatca.KA != 5.0 {
			t.Errorf("KA = %f, want 5.0", stheatca.KA)
		}
	})
}

// TestStheatint tests the Stheatint function
func TestStheatint(t *testing.T) {
	t.Run("WithEnvname_NumericValue", func(t *testing.T) {
		// Create STHEAT with Envname set to a numeric value
		elout := &ELOUT{Control: ON_SW, Fluid: AIRa_FLD}
		cmp := &COMPNT{
			Name:    "TestSTHEAT",
			Envname: "20.0", // Numeric value - envptr will create constant
			Tparm:   "25.0", // Initial temperature
			Elouts:  []*ELOUT{elout},
		}
		stheatca := &STHEATCA{
			Name: "TestCA",
			Q:    1000.0,
			Hcap: 5000.0,
			KA:   1.0,
			Eff:  0.9,
		}
		stheat := &STHEAT{Name: "TestSTHEAT", Cat: stheatca, Cmp: cmp}
		stheats := []*STHEAT{stheat}

		Stheatint(stheats, nil, nil, nil, nil)

		if stheat.Tenv == nil {
			t.Error("Tenv should be set when Envname is numeric")
		}
		if stheat.Tsold != 25.0 {
			t.Errorf("Tsold = %f, want 25.0", stheat.Tsold)
		}
	})

	t.Run("WithRoomname", func(t *testing.T) {
		// Create STHEAT with Roomname
		elout := &ELOUT{Control: ON_SW, Fluid: AIRa_FLD}
		cmp := &COMPNT{
			Name:     "TestSTHEAT",
			Roomname: "TestRoom",
			Tparm:    "22.0",
			Elouts:   []*ELOUT{elout},
		}
		stheatca := &STHEATCA{
			Name: "TestCA",
			Q:    1000.0,
			Hcap: 5000.0,
			KA:   1.0,
			Eff:  0.9,
		}
		stheat := &STHEAT{Name: "TestSTHEAT", Cat: stheatca, Cmp: cmp}
		stheats := []*STHEAT{stheat}

		// Create a room component for roomptr to find
		room := &ROOM{Name: "TestRoom", Tot: 20.0}
		roomCmp := &COMPNT{Name: "TestRoom", Eqp: room}
		compnts := []*COMPNT{roomCmp}

		Stheatint(stheats, nil, compnts, nil, nil)

		if stheat.Room == nil {
			t.Error("Room should be set when Roomname is provided")
		} else if stheat.Room.Name != "TestRoom" {
			t.Errorf("Room.Name = %s, want TestRoom", stheat.Room.Name)
		}
		if stheat.Tsold != 22.0 {
			t.Errorf("Tsold = %f, want 22.0", stheat.Tsold)
		}
	})

	t.Run("WithPCM", func(t *testing.T) {
		// Create STHEAT with PCM reference
		elout := &ELOUT{Control: ON_SW, Fluid: AIRa_FLD}
		cmp := &COMPNT{
			Name:    "TestSTHEAT",
			Envname: "15.0",
			Tparm:   "30.0",
			Elouts:  []*ELOUT{elout},
		}
		stheatca := &STHEATCA{
			Name:    "TestCA",
			PCMName: "TestPCM",
			Q:       1000.0,
			KA:      1.0,
			Eff:     0.9,
		}
		stheat := &STHEAT{Name: "TestSTHEAT", Cat: stheatca, Cmp: cmp}
		stheats := []*STHEAT{stheat}

		// Create PCM for lookup
		pcm := &PCM{Name: "TestPCM", Ql: 10000.0}
		pcms := []*PCM{pcm}

		Stheatint(stheats, nil, nil, nil, pcms)

		if stheat.Pcm == nil {
			t.Error("Pcm should be set when PCMName is provided")
		} else if stheat.Pcm.Name != "TestPCM" {
			t.Errorf("Pcm.Name = %s, want TestPCM", stheat.Pcm.Name)
		}
	})

	t.Run("WithMPCM", func(t *testing.T) {
		// Create STHEAT with MPCM
		elout := &ELOUT{Control: ON_SW, Fluid: AIRa_FLD}
		cmp := &COMPNT{
			Name:    "TestSTHEAT",
			Envname: "15.0",
			Tparm:   "25.0",
			MPCM:    100.0, // 100 kg of PCM
			Elouts:  []*ELOUT{elout},
		}
		stheatca := &STHEATCA{
			Name: "TestCA",
			Q:    1000.0,
			Hcap: 5000.0,
			KA:   1.0,
			Eff:  0.9,
		}
		stheat := &STHEAT{Name: "TestSTHEAT", Cat: stheatca, Cmp: cmp}
		stheats := []*STHEAT{stheat}

		Stheatint(stheats, nil, nil, nil, nil)

		if stheat.MPCM != 100.0 {
			t.Errorf("MPCM = %f, want 100.0", stheat.MPCM)
		}
	})

	t.Run("EmptyList", func(t *testing.T) {
		// Empty list should not panic
		stheats := []*STHEAT{}
		Stheatint(stheats, nil, nil, nil, nil)
	})

	t.Run("NegativeQ_Warning", func(t *testing.T) {
		// Test negative Q triggers error message (but continues)
		elout := &ELOUT{Control: ON_SW, Fluid: AIRa_FLD}
		cmp := &COMPNT{
			Name:    "TestSTHEAT",
			Envname: "20.0",
			Tparm:   "25.0",
			Elouts:  []*ELOUT{elout},
		}
		stheatca := &STHEATCA{
			Name: "TestCA",
			Q:    -100.0, // Negative Q
			Hcap: 5000.0,
			KA:   1.0,
			Eff:  0.9,
		}
		stheat := &STHEAT{Name: "TestSTHEAT", Cat: stheatca, Cmp: cmp}
		stheats := []*STHEAT{stheat}

		// Should not panic, just print warning
		Stheatint(stheats, nil, nil, nil, nil)
	})

	t.Run("NegativeHcap_Warning", func(t *testing.T) {
		// Test negative Hcap triggers error message (but continues)
		elout := &ELOUT{Control: ON_SW, Fluid: AIRa_FLD}
		cmp := &COMPNT{
			Name:    "TestSTHEAT",
			Envname: "20.0",
			Tparm:   "25.0",
			Elouts:  []*ELOUT{elout},
		}
		stheatca := &STHEATCA{
			Name: "TestCA",
			Q:    1000.0,
			Hcap: -5000.0, // Negative Hcap
			KA:   1.0,
			Eff:  0.9,
		}
		stheat := &STHEAT{Name: "TestSTHEAT", Cat: stheatca, Cmp: cmp}
		stheats := []*STHEAT{stheat}

		// Should not panic, just print warning
		Stheatint(stheats, nil, nil, nil, nil)
	})

	t.Run("NegativeKA_Warning", func(t *testing.T) {
		// Test negative KA triggers error message (but continues)
		elout := &ELOUT{Control: ON_SW, Fluid: AIRa_FLD}
		cmp := &COMPNT{
			Name:    "TestSTHEAT",
			Envname: "20.0",
			Tparm:   "25.0",
			Elouts:  []*ELOUT{elout},
		}
		stheatca := &STHEATCA{
			Name: "TestCA",
			Q:    1000.0,
			Hcap: 5000.0,
			KA:   -1.0, // Negative KA
			Eff:  0.9,
		}
		stheat := &STHEAT{Name: "TestSTHEAT", Cat: stheatca, Cmp: cmp}
		stheats := []*STHEAT{stheat}

		// Should not panic, just print warning
		Stheatint(stheats, nil, nil, nil, nil)
	})

	t.Run("NegativeEff_Warning", func(t *testing.T) {
		// Test negative Eff triggers error message (but continues)
		elout := &ELOUT{Control: ON_SW, Fluid: AIRa_FLD}
		cmp := &COMPNT{
			Name:    "TestSTHEAT",
			Envname: "20.0",
			Tparm:   "25.0",
			Elouts:  []*ELOUT{elout},
		}
		stheatca := &STHEATCA{
			Name: "TestCA",
			Q:    1000.0,
			Hcap: 5000.0,
			KA:   1.0,
			Eff:  -0.9, // Negative Eff
		}
		stheat := &STHEAT{Name: "TestSTHEAT", Cat: stheatca, Cmp: cmp}
		stheats := []*STHEAT{stheat}

		// Should not panic, just print warning
		Stheatint(stheats, nil, nil, nil, nil)
	})
}

// Helper function to create a basic STHEAT for testing
func createBasicSTHEAT() *STHEAT {
	// Create ELOUTs
	elout1 := &ELOUT{
		Control: ON_SW,
		Fluid:   AIRa_FLD,
		G:       0.1, // 0.1 kg/s
		Coeffin: make([]float64, 1),
		Coeffo:  0.0,
		Co:      0.0,
		Sysv:    30.0, // Output temperature
		Elins:   []*ELIN{},
	}
	elout2 := &ELOUT{
		Control: ON_SW,
		Fluid:   AIRx_FLD,
		G:       0.1,
		Coeffin: make([]float64, 1),
		Coeffo:  0.0,
		Co:      0.0,
		Sysv:    0.010, // Output humidity
	}

	// Create ELINs
	elin1 := &ELIN{
		Sysvin: 20.0, // 20°C inlet temperature
	}
	elin2 := &ELIN{
		Sysvin: 0.008, // Inlet humidity
	}

	// Link elin to elout
	elout1.Elins = []*ELIN{elin1}

	// Create COMPNT
	cmp := &COMPNT{
		Name:    "TestSTHEAT",
		Control: ON_SW,
		Elouts:  []*ELOUT{elout1, elout2},
		Elins:   []*ELIN{elin1, elin2},
		Tparm:   "25.0", // Initial storage temperature
	}

	// Create STHEATCA
	stheatca := &STHEATCA{
		Name: "TestSTHEATCA",
		Q:    3000.0,  // 3kW rated capacity
		Eff:  0.95,    // 95% efficiency
		Hcap: 50000.0, // 50kJ heat capacity
		KA:   5.0,     // Heat transfer coefficient
	}

	return &STHEAT{
		Name:  "TestSTHEAT",
		Cat:   stheatca,
		Cmp:   cmp,
		Tsold: 25.0, // Initial storage temperature
	}
}

// TestStheatcfv tests the Stheatcfv function
func TestStheatcfv(t *testing.T) {
	t.Run("BasicCalculation_WithRoom", func(t *testing.T) {
		stheat := createBasicSTHEAT()
		stheat.Room = &ROOM{
			Name: "TestRoom",
			Tot:  22.0, // Room temperature
		}
		stheat.Cmp.Envname = "" // No environment name

		stheats := []*STHEAT{stheat}
		Stheatcfv(stheats)

		// Check CG is calculated
		expectedCG := Spcheat(AIRa_FLD) * 0.1
		if math.Abs(stheat.CG-expectedCG) > 1e-6 {
			t.Errorf("CG = %f, want %f", stheat.CG, expectedCG)
		}

		// Check Hcap is set from Cat (no PCM)
		if stheat.Hcap != 50000.0 {
			t.Errorf("Hcap = %f, want 50000.0", stheat.Hcap)
		}

		// Check E is set when Control is ON
		if stheat.E != 3000.0 {
			t.Errorf("E = %f, want 3000.0", stheat.E)
		}

		// Check Elout coefficients are set
		elout := stheat.Cmp.Elouts[0]
		if elout.Coeffo != 1.0 {
			t.Errorf("Elouts[0].Coeffo = %f, want 1.0", elout.Coeffo)
		}
	})

	t.Run("BasicCalculation_WithTenv", func(t *testing.T) {
		stheat := createBasicSTHEAT()
		envTemp := 20.0
		stheat.Tenv = &envTemp
		stheat.Cmp.Envname = "outdoor"
		stheat.Room = nil

		stheats := []*STHEAT{stheat}
		Stheatcfv(stheats)

		// Check Hcap is set from Cat (no PCM)
		if stheat.Hcap != 50000.0 {
			t.Errorf("Hcap = %f, want 50000.0", stheat.Hcap)
		}
	})

	t.Run("ControlOFF_Component", func(t *testing.T) {
		stheat := createBasicSTHEAT()
		stheat.Cmp.Control = OFF_SW
		stheat.Room = &ROOM{
			Name: "TestRoom",
			Tot:  22.0,
		}

		stheats := []*STHEAT{stheat}
		Stheatcfv(stheats)

		// E should be 0 when Control is OFF_SW
		if stheat.E != 0.0 {
			t.Errorf("E = %f, want 0.0 when Control is OFF_SW", stheat.E)
		}
	})

	t.Run("ControlOFF_Elout", func(t *testing.T) {
		stheat := createBasicSTHEAT()
		stheat.Cmp.Elouts[0].Control = OFF_SW
		stheat.Room = &ROOM{
			Name: "TestRoom",
			Tot:  22.0,
		}

		stheats := []*STHEAT{stheat}
		Stheatcfv(stheats)

		// When Elout Control is OFF_SW, coefficients should be set differently
		elout := stheat.Cmp.Elouts[0]
		if elout.Coeffo != 1.0 {
			t.Errorf("Elouts[0].Coeffo = %f, want 1.0", elout.Coeffo)
		}
		if elout.Co != 0.0 {
			t.Errorf("Elouts[0].Co = %f, want 0.0", elout.Co)
		}
		if elout.Coeffin[0] != -1.0 {
			t.Errorf("Elouts[0].Coeffin[0] = %f, want -1.0", elout.Coeffin[0])
		}
	})

	t.Run("EmptyList", func(t *testing.T) {
		var stheats []*STHEAT
		// Should not panic with empty list
		Stheatcfv(stheats)
	})
}

// TestStheatene tests the Stheatene function
func TestStheatene(t *testing.T) {
	t.Run("BasicEnergyCalculation", func(t *testing.T) {
		stheat := createBasicSTHEAT()
		stheat.Room = &ROOM{
			Name: "TestRoom",
			Tot:  22.0,
			Qeqp: 0.0,
		}
		stheat.Cmp.Envname = ""
		stheat.CG = Spcheat(AIRa_FLD) * 0.1 // Pre-calculated
		stheat.Hcap = 50000.0
		stheat.E = 3000.0
		stheat.Tsold = 30.0

		stheats := []*STHEAT{stheat}
		Stheatene(stheats)

		// Check Tin is set from input
		if stheat.Tin != 20.0 {
			t.Errorf("Tin = %f, want 20.0", stheat.Tin)
		}

		// Check Tout is set
		if stheat.Tout != 30.0 {
			t.Errorf("Tout = %f, want 30.0", stheat.Tout)
		}

		// Check Q is calculated: Q = CG * (Tout - Tin)
		expectedQ := stheat.CG * (stheat.Tout - stheat.Tin)
		if math.Abs(stheat.Q-expectedQ) > 1e-6 {
			t.Errorf("Q = %f, want %f", stheat.Q, expectedQ)
		}

		// Note: Ts, Qls, Qsto calculations depend on DTM which may not be set
		// These values are only logged for debugging
		t.Logf("Ts = %f, Qls = %f, Qsto = %f", stheat.Ts, stheat.Qls, stheat.Qsto)

		// Verify Tsold is updated to match Ts
		if !math.IsNaN(stheat.Ts) && stheat.Tsold != stheat.Ts {
			t.Errorf("Tsold = %f should equal Ts = %f after update", stheat.Tsold, stheat.Ts)
		}
	})

	t.Run("RoomQeqpUpdate", func(t *testing.T) {
		stheat := createBasicSTHEAT()
		stheat.Room = &ROOM{
			Name: "TestRoom",
			Tot:  22.0,
			Qeqp: 0.0,
		}
		stheat.Cmp.Envname = ""
		stheat.CG = Spcheat(AIRa_FLD) * 0.1
		stheat.Hcap = 50000.0
		stheat.E = 3000.0
		stheat.Tsold = 30.0

		initialQeqp := stheat.Room.Qeqp
		stheats := []*STHEAT{stheat}
		Stheatene(stheats)

		// Room.Qeqp should be updated with -Qls
		if stheat.Room.Qeqp == initialQeqp && stheat.Qls != 0.0 {
			t.Error("Room.Qeqp should be updated with heat loss")
		}
	})

	t.Run("EmptyList", func(t *testing.T) {
		var stheats []*STHEAT
		// Should not panic with empty list
		Stheatene(stheats)
	})
}

// TestStheatprint tests the stheat output function
func TestStheatprint(t *testing.T) {
	t.Run("Header1_id0", func(t *testing.T) {
		stheat := createBasicSTHEAT()
		stheats := []*STHEAT{stheat}

		var buf bytes.Buffer
		stheatprint(&buf, 0, stheats)
		output := buf.String()

		if !strings.Contains(output, string(STHEAT_TYPE)) {
			t.Error("Output should contain STHEAT type")
		}
		if !strings.Contains(output, stheat.Name) {
			t.Error("Output should contain stheat name")
		}
	})

	t.Run("Header2_id1", func(t *testing.T) {
		stheat := createBasicSTHEAT()
		stheats := []*STHEAT{stheat}

		var buf bytes.Buffer
		stheatprint(&buf, 1, stheats)
		output := buf.String()

		// Should contain column headers
		if !strings.Contains(output, "_c") {
			t.Error("Output should contain control column")
		}
		if !strings.Contains(output, "_G") {
			t.Error("Output should contain flow rate column")
		}
		if !strings.Contains(output, "_Ts") {
			t.Error("Output should contain storage temp column")
		}
	})

	t.Run("Data_id99", func(t *testing.T) {
		stheat := createBasicSTHEAT()
		stheat.Ts = 35.0
		stheat.Tin = 20.0
		stheat.Tout = 30.0
		stheat.Q = 1000.0
		stheat.Qsto = 500.0
		stheat.Qls = 100.0
		stheat.E = 3000.0
		stheat.Hcap = 50000.0
		stheats := []*STHEAT{stheat}

		var buf bytes.Buffer
		stheatprint(&buf, 99, stheats)
		output := buf.String()

		if len(output) == 0 {
			t.Error("Output should not be empty")
		}
	})

	t.Run("EmptyList", func(t *testing.T) {
		var stheats []*STHEAT

		var buf bytes.Buffer
		stheatprint(&buf, 0, stheats)
		output := buf.String()

		if strings.Contains(output, string(STHEAT_TYPE)) {
			t.Error("Output should not contain type for empty list")
		}
	})
}

// TestStheatday tests the daily aggregation function
func TestStheatday(t *testing.T) {
	t.Run("DailyAggregation", func(t *testing.T) {
		stheat := createBasicSTHEAT()
		stheat.Tin = 20.0
		stheat.Tout = 30.0
		stheat.Ts = 35.0
		stheat.Q = 1000.0
		stheat.E = 3000.0
		stheat.Qls = 100.0
		stheat.Qsto = 500.0

		stheats := []*STHEAT{stheat}

		// Initialize daily counters
		stheatdyint(stheats)

		// Run aggregation for several time steps
		for ttmm := 100; ttmm <= 2300; ttmm += 100 {
			stheat.Tin = 18.0 + float64(ttmm)/1000.0
			stheat.Tout = 28.0 + float64(ttmm)/500.0
			stheat.Ts = 32.0 + float64(ttmm)/1000.0
			stheat.Q = 900.0 + float64(ttmm)/10.0
			stheat.E = 2500.0 + float64(ttmm)/5.0
			stheat.Qls = 80.0 + float64(ttmm)/100.0
			stheat.Qsto = 400.0 + float64(ttmm)/20.0
			stheatday(1, 1, ttmm, stheats, 1, 365)
		}

		// Verify aggregation results
		if stheat.Tidy.Hrs == 0 {
			t.Error("Tidy.Hrs should be > 0 after aggregation")
		}
		if stheat.Qlossdy == 0 {
			t.Error("Qlossdy should be > 0 after aggregation")
		}
		if stheat.Qstody == 0 {
			t.Error("Qstody should be > 0 after aggregation")
		}
	})

	t.Run("MonthlyAggregation_EndOfDay", func(t *testing.T) {
		stheat := createBasicSTHEAT()
		stheat.Tin = 20.0
		stheat.Tout = 30.0
		stheat.Ts = 35.0
		stheat.Q = 1000.0
		stheat.E = 3000.0
		stheat.Qls = 100.0
		stheat.Qsto = 500.0

		stheats := []*STHEAT{stheat}

		// Initialize daily and monthly counters
		stheatdyint(stheats)
		stheatmonint(stheats)

		// Simulate end of day
		stheatday(1, 1, 2400, stheats, 1, 365)

		t.Log("Monthly aggregation test completed")
	})

	t.Run("EmptyList", func(t *testing.T) {
		var stheats []*STHEAT

		// Should not panic with empty list
		stheatday(1, 1, 100, stheats, 1, 365)
	})
}

// TestStheatdyprt tests the daily output function
func TestStheatdyprt(t *testing.T) {
	t.Run("Header1_id0", func(t *testing.T) {
		stheat := createBasicSTHEAT()
		stheats := []*STHEAT{stheat}

		var buf bytes.Buffer
		stheatdyprt(&buf, 0, stheats)
		output := buf.String()

		if !strings.Contains(output, string(STHEAT_TYPE)) {
			t.Error("Output should contain STHEAT type")
		}
	})

	t.Run("Data_id99", func(t *testing.T) {
		stheat := createBasicSTHEAT()
		stheat.Tidy = SVDAY{Hrs: 10, M: 20.0, Mn: 18.0, Mx: 22.0}
		stheat.Tody = SVDAY{Hrs: 10, M: 30.0, Mn: 28.0, Mx: 32.0}
		stheat.Tsdy = SVDAY{Hrs: 10, M: 35.0, Mn: 30.0, Mx: 40.0}
		stheat.Qdy = QDAY{Hhr: 8, H: 8000.0, Chr: 0, C: 0}
		stheat.Edy = EDAY{Hrs: 10, D: 30000.0, Mx: 3500.0}
		stheat.Qlossdy = 1000.0
		stheat.Qstody = 5000.0
		stheats := []*STHEAT{stheat}

		var buf bytes.Buffer
		stheatdyprt(&buf, 99, stheats)
		output := buf.String()

		if len(output) == 0 {
			t.Error("Output should not be empty")
		}
	})
}

// TestStheatmonprt tests the monthly output function
func TestStheatmonprt(t *testing.T) {
	t.Run("Header1_id0", func(t *testing.T) {
		stheat := createBasicSTHEAT()
		stheats := []*STHEAT{stheat}

		var buf bytes.Buffer
		stheatmonprt(&buf, 0, stheats)
		output := buf.String()

		if !strings.Contains(output, string(STHEAT_TYPE)) {
			t.Error("Output should contain STHEAT type")
		}
	})

	t.Run("Data_id99", func(t *testing.T) {
		stheat := createBasicSTHEAT()
		stheat.MTidy = SVDAY{Hrs: 100, M: 20.0, Mn: 15.0, Mx: 25.0}
		stheat.MTody = SVDAY{Hrs: 100, M: 30.0, Mn: 25.0, Mx: 35.0}
		stheat.MTsdy = SVDAY{Hrs: 100, M: 35.0, Mn: 25.0, Mx: 45.0}
		stheat.MQdy = QDAY{Hhr: 80, H: 80000.0, Chr: 0, C: 0}
		stheat.MEdy = EDAY{Hrs: 100, D: 300000.0, Mx: 4000.0}
		stheat.MQlossdy = 10000.0
		stheat.MQstody = 50000.0
		stheats := []*STHEAT{stheat}

		var buf bytes.Buffer
		stheatmonprt(&buf, 99, stheats)
		output := buf.String()

		if len(output) == 0 {
			t.Error("Output should not be empty")
		}
	})
}

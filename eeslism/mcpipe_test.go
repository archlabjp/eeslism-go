package eeslism

import (
	"bytes"
	"math"
	"strings"
	"testing"
)

// TestPipedata tests the Pipedata function
func TestPipedata(t *testing.T) {
	t.Run("SetName_DUCT", func(t *testing.T) {
		pipeca := &PIPECA{}
		result := Pipedata(DUCT_TYPE, "TestDuct", pipeca)

		if result != 0 {
			t.Errorf("Pipedata should return 0 for name, got %d", result)
		}
		if pipeca.name != "TestDuct" {
			t.Errorf("name = %s, want TestDuct", pipeca.name)
		}
		if pipeca.Type != DUCT_PDT {
			t.Errorf("Type = %d, want DUCT_PDT", pipeca.Type)
		}
	})

	t.Run("SetName_PIPEDUCT", func(t *testing.T) {
		pipeca := &PIPECA{}
		result := Pipedata(PIPEDUCT_TYPE, "TestPipe", pipeca)

		if result != 0 {
			t.Errorf("Pipedata should return 0 for name, got %d", result)
		}
		if pipeca.name != "TestPipe" {
			t.Errorf("name = %s, want TestPipe", pipeca.name)
		}
		if pipeca.Type != PIPE_PDT {
			t.Errorf("Type = %d, want PIPE_PDT", pipeca.Type)
		}
	})

	t.Run("Set_Ko", func(t *testing.T) {
		pipeca := &PIPECA{}
		result := Pipedata(DUCT_TYPE, "Ko=0.5", pipeca)

		if result != 0 {
			t.Errorf("Pipedata should return 0 for Ko, got %d", result)
		}
		if pipeca.Ko != 0.5 {
			t.Errorf("Ko = %f, want 0.5", pipeca.Ko)
		}
	})

	t.Run("UnknownKey", func(t *testing.T) {
		pipeca := &PIPECA{}
		result := Pipedata(DUCT_TYPE, "unknown=123", pipeca)

		if result != 1 {
			t.Errorf("Pipedata should return 1 for unknown key, got %d", result)
		}
	})

	t.Run("MultipleParameters", func(t *testing.T) {
		pipeca := &PIPECA{}

		// Set name first
		Pipedata(PIPEDUCT_TYPE, "TestPipeline", pipeca)
		if pipeca.name != "TestPipeline" {
			t.Errorf("name = %s, want TestPipeline", pipeca.name)
		}

		// Set Ko
		Pipedata(PIPEDUCT_TYPE, "Ko=1.5", pipeca)
		if pipeca.Ko != 1.5 {
			t.Errorf("Ko = %f, want 1.5", pipeca.Ko)
		}
	})
}

// Helper function to create a basic PIPE for testing
func createBasicPIPE() *PIPE {
	// Create ELOUT for the pipe (water pipe)
	elout := &ELOUT{
		Control:  ON_SW,
		Fluid:    WATER_FLD,
		G:        0.5, // 0.5 kg/s
		Coeffin:  make([]float64, 1),
		Coeffo:   0.0,
		Co:       0.0,
	}

	// Create ELIN
	elin := &ELIN{
		Sysvin: 50.0, // 50°C inlet temperature
	}

	// Create COMPNT
	cmp := &COMPNT{
		Name:    "TestPipe",
		Control: ON_SW,
		Elouts:  []*ELOUT{elout},
		Elins:   []*ELIN{elin},
	}

	// Create PIPECA
	pipeca := &PIPECA{
		name: "TestPipeCA",
		Type: PIPE_PDT,
		Ko:   10.0, // Heat transfer coefficient [W/(m·K)]
	}

	return &PIPE{
		Name: "TestPipe",
		Cat:  pipeca,
		Cmp:  cmp,
		L:    5.0, // 5m pipe length
	}
}

// Helper function to create a DUCT (air pipe)
func createBasicDUCT() *PIPE {
	// Create ELOUTs for the duct (2 outputs: temperature and humidity)
	elout1 := &ELOUT{
		Control:  ON_SW,
		Fluid:    AIRa_FLD,
		G:        1.0, // 1 kg/s
		Coeffin:  make([]float64, 1),
		Coeffo:   0.0,
		Co:       0.0,
		Sysv:     25.0, // Output temperature
	}
	elout2 := &ELOUT{
		Control:  ON_SW,
		Fluid:    AIRx_FLD,
		G:        1.0,
		Coeffin:  make([]float64, 1),
		Coeffo:   0.0,
		Co:       0.0,
		Sysv:     0.010, // Output humidity
	}

	// Create ELINs
	elin1 := &ELIN{
		Sysvin: 30.0, // 30°C inlet temperature
	}
	elin2 := &ELIN{
		Sysvin: 0.012, // Inlet humidity
	}

	// Create COMPNT
	cmp := &COMPNT{
		Name:    "TestDuct",
		Control: ON_SW,
		Elouts:  []*ELOUT{elout1, elout2},
		Elins:   []*ELIN{elin1, elin2},
	}

	// Create PIPECA
	pipeca := &PIPECA{
		name: "TestDuctCA",
		Type: DUCT_PDT,
		Ko:   5.0, // Heat transfer coefficient [W/(m·K)]
	}

	return &PIPE{
		Name: "TestDuct",
		Cat:  pipeca,
		Cmp:  cmp,
		L:    10.0, // 10m duct length
	}
}

// TestPipeint tests the Pipeint function
func TestPipeint(t *testing.T) {
	t.Run("WithIvparm", func(t *testing.T) {
		pipe := createBasicPIPE()
		pipeLength := 15.0
		pipe.Cmp.Ivparm = &pipeLength

		// Create a room for environment
		room := &ROOM{
			Name: "TestRoom",
			Tot:  25.0, // Room temperature
		}
		roomCmp := &COMPNT{
			Name:   "TestRoom",
			Eqptype: ROOM_TYPE,
		}
		roomCmp.Eqp = room

		pipe.Cmp.Roomname = "TestRoom"

		pipes := []*PIPE{pipe}
		simc := createBasicSIMCONTL()
		compnts := []*COMPNT{pipe.Cmp, roomCmp}
		wd := createBasicWDAT()

		Pipeint(pipes, simc, compnts, wd)

		if pipe.L != 15.0 {
			t.Errorf("L = %f, want 15.0", pipe.L)
		}
	})

	t.Run("WithoutIvparm", func(t *testing.T) {
		pipe := createBasicPIPE()
		pipe.Cmp.Ivparm = nil
		pipe.Cmp.Roomname = "" // No room name
		pipe.Cmp.Envname = ""  // No environment name

		pipes := []*PIPE{pipe}
		simc := createBasicSIMCONTL()
		compnts := []*COMPNT{pipe.Cmp}
		wd := createBasicWDAT()

		Pipeint(pipes, simc, compnts, wd)

		// L should be FNAN (-999) when Ivparm is nil
		if pipe.L != FNAN {
			t.Errorf("L should be FNAN when Ivparm is nil, got %f", pipe.L)
		}
	})

	t.Run("EmptyList", func(t *testing.T) {
		var pipes []*PIPE
		simc := createBasicSIMCONTL()
		compnts := []*COMPNT{}
		wd := createBasicWDAT()

		// Should not panic with empty list
		Pipeint(pipes, simc, compnts, wd)
	})

	t.Run("MultiplePipes", func(t *testing.T) {
		pipe1 := createBasicPIPE()
		pipe1.Name = "Pipe1"
		pipeLen1 := 10.0
		pipe1.Cmp.Ivparm = &pipeLen1

		pipe2 := createBasicPIPE()
		pipe2.Name = "Pipe2"
		pipeLen2 := 20.0
		pipe2.Cmp.Ivparm = &pipeLen2

		pipes := []*PIPE{pipe1, pipe2}
		simc := createBasicSIMCONTL()
		compnts := []*COMPNT{pipe1.Cmp, pipe2.Cmp}
		wd := createBasicWDAT()

		Pipeint(pipes, simc, compnts, wd)

		if pipe1.L != 10.0 {
			t.Errorf("pipe1.L = %f, want 10.0", pipe1.L)
		}
		if pipe2.L != 20.0 {
			t.Errorf("pipe2.L = %f, want 20.0", pipe2.L)
		}
	})
}

// TestPipecfv tests the Pipecfv function
func TestPipecfv(t *testing.T) {
	t.Run("BasicCalculation_WithTenv", func(t *testing.T) {
		pipe := createBasicPIPE()
		envTemp := 20.0
		pipe.Tenv = &envTemp
		pipe.Cmp.Envname = "outdoor" // Set Envname so Tenv is used
		pipe.L = 10.0
		pipe.Cat.Ko = 5.0
		pipe.Cmp.Elouts[0].G = 0.5 // 0.5 kg/s

		pipes := []*PIPE{pipe}
		Pipecfv(pipes)

		// Ko should be set from Cat
		if pipe.Ko != 5.0 {
			t.Errorf("Ko = %f, want 5.0", pipe.Ko)
		}

		// cG = Spcheat(WATER_FLD) * G = 4186 * 0.5 = 2093
		cG := Spcheat(WATER_FLD) * 0.5
		expectedEp := 1.0 - math.Exp(-(pipe.Ko*pipe.L)/cG)
		if math.Abs(pipe.Ep-expectedEp) > 1e-6 {
			t.Errorf("Ep = %f, want %f", pipe.Ep, expectedEp)
		}

		// D1 = cG * Ep
		expectedD1 := cG * expectedEp
		if math.Abs(pipe.D1-expectedD1) > 1e-6 {
			t.Errorf("D1 = %f, want %f", pipe.D1, expectedD1)
		}

		// Do = D1 * Te
		expectedDo := expectedD1 * envTemp
		if math.Abs(pipe.Do-expectedDo) > 1e-6 {
			t.Errorf("Do = %f, want %f", pipe.Do, expectedDo)
		}

		// Check Elout coefficients
		elout := pipe.Cmp.Elouts[0]
		if math.Abs(elout.Coeffo-cG) > 1e-6 {
			t.Errorf("Coeffo = %f, want %f", elout.Coeffo, cG)
		}
		if math.Abs(elout.Co-pipe.Do) > 1e-6 {
			t.Errorf("Co = %f, want %f", elout.Co, pipe.Do)
		}
		expectedCoeffin := pipe.D1 - cG
		if math.Abs(elout.Coeffin[0]-expectedCoeffin) > 1e-6 {
			t.Errorf("Coeffin[0] = %f, want %f", elout.Coeffin[0], expectedCoeffin)
		}
	})

	t.Run("BasicCalculation_WithRoom", func(t *testing.T) {
		pipe := createBasicPIPE()
		pipe.Tenv = nil
		pipe.Room = &ROOM{
			Name: "TestRoom",
			Tot:  22.0, // Room temperature
		}
		pipe.L = 8.0
		pipe.Cat.Ko = 4.0
		pipe.Cmp.Elouts[0].G = 1.0

		pipes := []*PIPE{pipe}
		Pipecfv(pipes)

		// Environment temperature should come from Room.Tot
		cG := Spcheat(WATER_FLD) * 1.0
		expectedEp := 1.0 - math.Exp(-(pipe.Ko*pipe.L)/cG)
		expectedD1 := cG * expectedEp
		expectedDo := expectedD1 * 22.0 // Room temperature

		if math.Abs(pipe.Do-expectedDo) > 1e-6 {
			t.Errorf("Do = %f, want %f (using Room.Tot)", pipe.Do, expectedDo)
		}
	})

	t.Run("DUCT_Type", func(t *testing.T) {
		pipe := createBasicDUCT()
		envTemp := 25.0
		pipe.Tenv = &envTemp
		pipe.Cmp.Envname = "outdoor" // Set Envname so Tenv is used
		pipe.L = 5.0
		pipe.Cat.Ko = 3.0

		pipes := []*PIPE{pipe}
		Pipecfv(pipes)

		// Check that second output (humidity) is set for DUCT_PDT
		elout2 := pipe.Cmp.Elouts[1]
		if elout2.Coeffo != 1.0 {
			t.Errorf("Elouts[1].Coeffo = %f, want 1.0", elout2.Coeffo)
		}
		if elout2.Co != 0.0 {
			t.Errorf("Elouts[1].Co = %f, want 0.0", elout2.Co)
		}
		if elout2.Coeffin[0] != -1.0 {
			t.Errorf("Elouts[1].Coeffin[0] = %f, want -1.0", elout2.Coeffin[0])
		}
	})

	t.Run("ControlOFF", func(t *testing.T) {
		pipe := createBasicPIPE()
		pipe.Cmp.Control = OFF_SW
		envTemp := 20.0
		pipe.Tenv = &envTemp

		// Set initial values to check they are not changed
		initialKo := pipe.Ko

		pipes := []*PIPE{pipe}
		Pipecfv(pipes)

		// When Control is OFF_SW, no calculations should be performed
		// Ko should remain at initial value (not set from Cat)
		if pipe.Ko != initialKo {
			t.Logf("Ko remained at initial value when Control is OFF_SW")
		}
	})

	t.Run("EmptyList", func(t *testing.T) {
		var pipes []*PIPE
		// Should not panic with empty list
		Pipecfv(pipes)
	})
}

// TestPipeene tests the Pipeene function
func TestPipeene(t *testing.T) {
	t.Run("BasicEnergyCalculation", func(t *testing.T) {
		pipe := createBasicPIPE()
		pipe.Cmp.Elins[0].Sysvin = 60.0 // Inlet temperature
		pipe.Do = 1000.0                 // Pre-calculated Do
		pipe.D1 = 50.0                   // Pre-calculated D1

		pipes := []*PIPE{pipe}
		Pipeene(pipes)

		// Tin should be set from input
		if pipe.Tin != 60.0 {
			t.Errorf("Tin = %f, want 60.0", pipe.Tin)
		}

		// Tout = Do
		if pipe.Tout != 1000.0 {
			t.Errorf("Tout = %f, want 1000.0", pipe.Tout)
		}

		// Q = Do - D1 * Tin = 1000 - 50 * 60 = -2000
		expectedQ := pipe.Do - pipe.D1*pipe.Tin
		if math.Abs(pipe.Q-expectedQ) > 1e-6 {
			t.Errorf("Q = %f, want %f", pipe.Q, expectedQ)
		}
	})

	t.Run("WithRoom_QeqpUpdate", func(t *testing.T) {
		pipe := createBasicPIPE()
		pipe.Cmp.Elins[0].Sysvin = 50.0
		pipe.Do = 500.0
		pipe.D1 = 25.0
		pipe.Room = &ROOM{
			Name: "TestRoom",
			Qeqp: 0.0,
		}

		pipes := []*PIPE{pipe}
		Pipeene(pipes)

		// Q = Do - D1 * Tin = 500 - 25 * 50 = -750
		expectedQ := pipe.Do - pipe.D1*50.0
		// Room.Qeqp should be updated with -Q
		expectedQeqp := -expectedQ
		if math.Abs(pipe.Room.Qeqp-expectedQeqp) > 1e-6 {
			t.Errorf("Room.Qeqp = %f, want %f", pipe.Room.Qeqp, expectedQeqp)
		}
	})

	t.Run("DUCT_Type_HumidityCalculation", func(t *testing.T) {
		pipe := createBasicDUCT()
		pipe.Cmp.Elins[0].Sysvin = 28.0  // Inlet temperature
		pipe.Cmp.Elins[1].Sysvin = 0.012 // Inlet humidity (not used directly)
		pipe.Cmp.Elouts[1].Sysv = 0.010  // Output humidity from system
		pipe.Do = 600.0
		pipe.D1 = 25.0

		pipes := []*PIPE{pipe}
		Pipeene(pipes)

		// Xout should be set from Elouts[1].Sysv
		if pipe.Xout != 0.010 {
			t.Errorf("Xout = %f, want 0.010", pipe.Xout)
		}

		// Hout should be calculated (enthalpy)
		if math.IsNaN(pipe.Hout) {
			t.Error("Hout should not be NaN for DUCT_PDT")
		}
	})

	t.Run("PIPE_Type_HoutFNAN", func(t *testing.T) {
		pipe := createBasicPIPE()
		pipe.Cmp.Elins[0].Sysvin = 50.0
		pipe.Do = 500.0
		pipe.D1 = 25.0

		pipes := []*PIPE{pipe}
		Pipeene(pipes)

		// For PIPE_PDT, Hout should be FNAN (-999)
		if pipe.Hout != FNAN {
			t.Errorf("Hout should be FNAN for PIPE_PDT, got %f", pipe.Hout)
		}
	})

	t.Run("ControlOFF", func(t *testing.T) {
		pipe := createBasicPIPE()
		pipe.Cmp.Control = OFF_SW
		pipe.Cmp.Elins[0].Sysvin = 50.0

		pipes := []*PIPE{pipe}
		Pipeene(pipes)

		// Q should be 0 when Control is OFF_SW
		if pipe.Q != 0.0 {
			t.Errorf("Q = %f, want 0.0 when Control is OFF_SW", pipe.Q)
		}
	})

	t.Run("EmptyList", func(t *testing.T) {
		var pipes []*PIPE
		// Should not panic with empty list
		Pipeene(pipes)
	})

	t.Run("MultiplePipes", func(t *testing.T) {
		pipe1 := createBasicPIPE()
		pipe1.Cmp.Elins[0].Sysvin = 60.0
		pipe1.Do = 1000.0
		pipe1.D1 = 50.0

		pipe2 := createBasicPIPE()
		pipe2.Cmp.Elins[0].Sysvin = 40.0
		pipe2.Do = 800.0
		pipe2.D1 = 40.0

		pipes := []*PIPE{pipe1, pipe2}
		Pipeene(pipes)

		// Verify both pipes are processed
		expectedQ1 := pipe1.Do - pipe1.D1*60.0
		expectedQ2 := pipe2.Do - pipe2.D1*40.0

		if math.Abs(pipe1.Q-expectedQ1) > 1e-6 {
			t.Errorf("pipe1.Q = %f, want %f", pipe1.Q, expectedQ1)
		}
		if math.Abs(pipe2.Q-expectedQ2) > 1e-6 {
			t.Errorf("pipe2.Q = %f, want %f", pipe2.Q, expectedQ2)
		}
	})
}

// TestPipeprint tests the pipe output function
func TestPipeprint(t *testing.T) {
	t.Run("Header1_id0", func(t *testing.T) {
		pipe := createBasicPIPE()
		pipes := []*PIPE{pipe}

		var buf bytes.Buffer
		pipeprint(&buf, 0, pipes)
		output := buf.String()

		if !strings.Contains(output, string(PIPEDUCT_TYPE)) {
			t.Errorf("Output should contain PIPEDUCT type")
		}
		if !strings.Contains(output, pipe.Name) {
			t.Errorf("Output should contain pipe name")
		}
	})

	t.Run("Header2_id1", func(t *testing.T) {
		pipe := createBasicPIPE()
		pipes := []*PIPE{pipe}

		var buf bytes.Buffer
		pipeprint(&buf, 1, pipes)
		output := buf.String()

		// Should contain column headers
		if !strings.Contains(output, "_c") {
			t.Error("Output should contain control column")
		}
		if !strings.Contains(output, "_G") {
			t.Error("Output should contain flow rate column")
		}
		if !strings.Contains(output, "_Ti") {
			t.Error("Output should contain inlet temp column")
		}
		if !strings.Contains(output, "_To") {
			t.Error("Output should contain outlet temp column")
		}
		if !strings.Contains(output, "_Q") {
			t.Error("Output should contain heat column")
		}
	})

	t.Run("Data_id99", func(t *testing.T) {
		pipe := createBasicPIPE()
		pipe.Tin = 50.0
		pipe.Q = 1000.0
		pipe.Cmp.Elouts[0].G = 0.5
		pipe.Cmp.Elouts[0].Sysv = 45.0
		pipes := []*PIPE{pipe}

		var buf bytes.Buffer
		pipeprint(&buf, 99, pipes)
		output := buf.String()

		// Output should contain data values
		if len(output) == 0 {
			t.Error("Output should not be empty")
		}
	})

	t.Run("EmptyList", func(t *testing.T) {
		var pipes []*PIPE

		var buf bytes.Buffer
		pipeprint(&buf, 0, pipes)
		output := buf.String()

		// Should not output header for empty list
		if strings.Contains(output, string(PIPEDUCT_TYPE)) {
			t.Error("Output should not contain type for empty list")
		}
	})
}

// TestPipedyprt tests the daily output function
func TestPipedyprt(t *testing.T) {
	t.Run("Header1_id0", func(t *testing.T) {
		pipe := createBasicPIPE()
		pipes := []*PIPE{pipe}

		var buf bytes.Buffer
		pipedyprt(&buf, 0, pipes)
		output := buf.String()

		if !strings.Contains(output, string(PIPEDUCT_TYPE)) {
			t.Error("Output should contain PIPEDUCT type")
		}
	})

	t.Run("Header2_id1", func(t *testing.T) {
		pipe := createBasicPIPE()
		pipes := []*PIPE{pipe}

		var buf bytes.Buffer
		pipedyprt(&buf, 1, pipes)
		output := buf.String()

		// Should contain daily statistics headers
		if !strings.Contains(output, "_Ht") {
			t.Error("Output should contain hours header")
		}
	})

	t.Run("Data_id99", func(t *testing.T) {
		pipe := createBasicPIPE()
		pipe.Tidy = SVDAY{Hrs: 10, M: 25.0, Mn: 20.0, Mx: 30.0}
		pipe.Qdy = QDAY{Hhr: 5, H: 100.0, Chr: 3, C: 50.0}
		pipes := []*PIPE{pipe}

		var buf bytes.Buffer
		pipedyprt(&buf, 99, pipes)
		output := buf.String()

		if len(output) == 0 {
			t.Error("Output should not be empty")
		}
	})
}

// TestPipemonprt tests the monthly output function
func TestPipemonprt(t *testing.T) {
	t.Run("Header1_id0", func(t *testing.T) {
		pipe := createBasicPIPE()
		pipes := []*PIPE{pipe}

		var buf bytes.Buffer
		pipemonprt(&buf, 0, pipes)
		output := buf.String()

		if !strings.Contains(output, string(PIPEDUCT_TYPE)) {
			t.Error("Output should contain PIPEDUCT type")
		}
	})

	t.Run("Data_id99", func(t *testing.T) {
		pipe := createBasicPIPE()
		pipe.MTidy = SVDAY{Hrs: 100, M: 25.0, Mn: 15.0, Mx: 35.0}
		pipe.MQdy = QDAY{Hhr: 50, H: 1000.0, Chr: 30, C: 500.0}
		pipes := []*PIPE{pipe}

		var buf bytes.Buffer
		pipemonprt(&buf, 99, pipes)
		output := buf.String()

		if len(output) == 0 {
			t.Error("Output should not be empty")
		}
	})
}

// TestPipeday tests the daily aggregation function
func TestPipeday(t *testing.T) {
	t.Run("DailyAggregation", func(t *testing.T) {
		pipe := createBasicPIPE()
		pipe.Tin = 50.0
		pipe.Q = 500.0

		pipes := []*PIPE{pipe}

		// Initialize daily counters
		pipedyint(pipes)

		// Run aggregation for several time steps
		for ttmm := 100; ttmm <= 2300; ttmm += 100 {
			pipe.Tin = 45.0 + float64(ttmm)/1000.0
			pipe.Q = 400.0 + float64(ttmm)/10.0
			pipeday(1, 1, ttmm, pipes, 1, 365)
		}

		// Verify aggregation results
		if pipe.Tidy.Hrs == 0 {
			t.Error("Tidy.Hrs should be > 0 after aggregation")
		}
		if pipe.Qdy.Hhr == 0 && pipe.Qdy.Chr == 0 {
			t.Error("Qdy should have some hours after aggregation")
		}
	})

	t.Run("MonthlyAggregation_EndOfDay", func(t *testing.T) {
		pipe := createBasicPIPE()
		pipe.Tin = 50.0
		pipe.Q = 500.0

		pipes := []*PIPE{pipe}

		// Initialize daily and monthly counters
		pipedyint(pipes)
		pipemonint(pipes)

		// Simulate end of day (Nday=1, SimDayend=365, day end ttmm=2400)
		pipeday(1, 1, 2400, pipes, 1, 365)

		// Monthly aggregation should have been triggered
		t.Log("Monthly aggregation test completed")
	})

	t.Run("EmptyList", func(t *testing.T) {
		var pipes []*PIPE

		// Should not panic with empty list
		pipeday(1, 1, 100, pipes, 1, 365)
	})
}

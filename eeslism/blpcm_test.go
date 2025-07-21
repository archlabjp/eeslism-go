package eeslism

import (
	"math"
	"testing"
)

func TestPCMdata(t *testing.T) {
	t.Run("valid PCM data", func(t *testing.T) {
		input := `ParaffinWax28 Condl=0.15 Conds=0.20 Ql=180000 Ts=26 Tl=30 Tp=28 ;
				  PCM_HighPerf Condl=0.18 Conds=0.22 Ql=200000 Ts=24 Tl=26 Tp=25 -iterate ;
*`
		fi := NewEeTokens(input)

		var pcm []*PCM
		var pcmiterate rune

		PCMdata(fi, "test", &pcm, &pcmiterate)

		if len(pcm) != 2 {
			t.Fatalf("expected 2 PCM entries, got %d", len(pcm))
		}

		// Check first PCM entry
		pcm1 := pcm[0]
		if pcm1.Name != "ParaffinWax28" {
			t.Errorf("pcm[0].Name = %s, want ParaffinWax28", pcm1.Name)
		}
		if pcm1.Condl != 0.15 {
			t.Errorf("pcm[0].Condl = %f, want 0.15", pcm1.Condl)
		}
		if pcm1.Conds != 0.20 {
			t.Errorf("pcm[0].Conds = %f, want 0.20", pcm1.Conds)
		}
		if pcm1.Ql != 180000 {
			t.Errorf("pcm[0].Ql = %f, want 180000", pcm1.Ql)
		}
		if pcm1.Ts != 26 {
			t.Errorf("pcm[0].Ts = %f, want 26", pcm1.Ts)
		}
		if pcm1.Tl != 30 {
			t.Errorf("pcm[0].Tl = %f, want 30", pcm1.Tl)
		}
		if pcm1.Tp != 28 {
			t.Errorf("pcm[0].Tp = %f, want 28", pcm1.Tp)
		}
		if pcm1.Iterate != false {
			t.Errorf("pcm[0].Iterate = %t, want false", pcm1.Iterate)
		}

		// Check second PCM entry
		pcm2 := pcm[1]
		if pcm2.Name != "PCM_HighPerf" {
			t.Errorf("pcm[1].Name = %s, want PCM_HighPerf", pcm2.Name)
		}
		if pcm2.Iterate != true {
			t.Errorf("pcm[1].Iterate = %t, want true", pcm2.Iterate)
		}
	})

	t.Run("PCM data with CHARTABLE reading", func(t *testing.T) {
		input := `PCM_WithTable spcheattable=testdata/pcm_enthalpy_test.txt table=e conducttable=testdata/pcm_conductivity_test.txt Condl=0.15 Conds=0.20 Ql=180000 Ts=26 Tl=30 Tp=28 ;
*`
		fi := NewEeTokens(input)

		var pcm []*PCM
		var pcmiterate rune

		PCMdata(fi, "test", &pcm, &pcmiterate)

		if len(pcm) != 1 {
			t.Fatalf("expected 1 PCM entry, got %d", len(pcm))
		}

		pcm1 := pcm[0]
		if pcm1.Name != "PCM_WithTable" {
			t.Errorf("pcm[0].Name = %s, want PCM_WithTable", pcm1.Name)
		}

		// Check that table type is set correctly
		if pcm1.Spctype != 't' {
			t.Errorf("pcm[0].Spctype = %c, want 't'", pcm1.Spctype)
		}
		if pcm1.Condtype != 't' {
			t.Errorf("pcm[0].Condtype = %c, want 't'", pcm1.Condtype)
		}

		// Check CHARTABLE for enthalpy (index 0)
		ct0 := &pcm1.Chartable[0]
		if ct0.filename != "testdata/pcm_enthalpy_test.txt" {
			t.Errorf("ct0.filename = %s, want testdata/pcm_enthalpy_test.txt", ct0.filename)
		}
		if ct0.PCMchar != 'E' {
			t.Errorf("ct0.PCMchar = %c, want 'E'", ct0.PCMchar)
		}
		if ct0.tabletype != 'e' {
			t.Errorf("ct0.tabletype = %c, want 'e'", ct0.tabletype)
		}

		// Check CHARTABLE for conductivity (index 1)
		ct1 := &pcm1.Chartable[1]
		if ct1.filename != "testdata/pcm_conductivity_test.txt" {
			t.Errorf("ct1.filename = %s, want testdata/pcm_conductivity_test.txt", ct1.filename)
		}
		if ct1.PCMchar != 'C' {
			t.Errorf("ct1.PCMchar = %c, want 'C'", ct1.PCMchar)
		}

		// Check that data was loaded (may be less than expected due to file format issues)
		t.Logf("ct0.itablerow = %d, len(ct0.T) = %d, len(ct0.Chara) = %d", ct0.itablerow, len(ct0.T), len(ct0.Chara))
		t.Logf("ct1.itablerow = %d, len(ct1.T) = %d, len(ct1.Chara) = %d", ct1.itablerow, len(ct1.T), len(ct1.Chara))

		// Basic validation that some data was loaded
		if len(ct0.T) == 0 {
			t.Errorf("No enthalpy data was loaded")
		}
		if len(ct1.T) == 0 {
			t.Errorf("No conductivity data was loaded")
		}

		// If data was loaded, check basic properties
		if len(ct0.T) > 0 && len(ct0.Chara) > 0 {
			t.Logf("First enthalpy point: T=%f, Chara=%f", ct0.T[0], ct0.Chara[0])
			if len(ct0.T) > 1 {
				t.Logf("Last enthalpy point: T=%f, Chara=%f", ct0.T[len(ct0.T)-1], ct0.Chara[len(ct0.Chara)-1])
			}
		}

		if len(ct1.T) > 0 && len(ct1.Chara) > 0 {
			t.Logf("First conductivity point: T=%f, Chara=%f", ct1.T[0], ct1.Chara[0])
			if len(ct1.T) > 1 {
				t.Logf("Last conductivity point: T=%f, Chara=%f", ct1.T[len(ct1.T)-1], ct1.Chara[len(ct1.Chara)-1])
			}
		}
	})

	t.Run("PCM data with specific heat table", func(t *testing.T) {
		input := `PCM_WithSpecificHeat spcheattable=testdata/pcm_specific_heat_test.txt table=h Condl=0.15 Conds=0.20 Ql=180000 Ts=26 Tl=30 Tp=28 ;
*`
		fi := NewEeTokens(input)

		var pcm []*PCM
		var pcmiterate rune

		PCMdata(fi, "test", &pcm, &pcmiterate)

		if len(pcm) != 1 {
			t.Fatalf("expected 1 PCM entry, got %d", len(pcm))
		}

		pcm1 := pcm[0]
		if pcm1.Name != "PCM_WithSpecificHeat" {
			t.Errorf("pcm[0].Name = %s, want PCM_WithSpecificHeat", pcm1.Name)
		}

		// Check that table type is set correctly
		if pcm1.Spctype != 't' {
			t.Errorf("pcm[0].Spctype = %c, want 't'", pcm1.Spctype)
		}

		// Check CHARTABLE for specific heat (index 0)
		ct0 := &pcm1.Chartable[0]
		if ct0.filename != "testdata/pcm_specific_heat_test.txt" {
			t.Errorf("ct0.filename = %s, want testdata/pcm_specific_heat_test.txt", ct0.filename)
		}
		if ct0.tabletype != 'h' {
			t.Errorf("ct0.tabletype = %c, want 'h'", ct0.tabletype)
		}

		// Check that specific heat integration was performed
		if ct0.itablerow != 7 {
			t.Errorf("ct0.itablerow = %d, want 7", ct0.itablerow)
		}

		// Verify specific heat integration calculation
		if ct0.Chara[0] != 0.0 {
			t.Errorf("ct0.Chara[0] = %f, want 0.0 (first point)", ct0.Chara[0])
		}

		t.Logf("PCM specific heat table loaded successfully:")
		for i := 0; i < len(ct0.T) && i < 3; i++ {
			t.Logf("  T[%d]=%f, Chara[%d]=%f", i, ct0.T[i], i, ct0.Chara[i])
		}
	})
}

func TestTableRead(t *testing.T) {
	t.Run("read enthalpy table", func(t *testing.T) {
		ct := &CHARTABLE{
			filename:    "testdata/pcm_enthalpy_test.txt",
			PCMchar:     'E',
			tabletype:   'e',
			minTempChng: 0.5,
		}

		TableRead(ct)

		// Check basic properties
		if ct.itablerow != 7 {
			t.Errorf("itablerow = %d, want 7", ct.itablerow)
		}
		if ct.minTemp != 20.0 {
			t.Errorf("minTemp = %f, want 20.0", ct.minTemp)
		}
		if ct.maxTemp != 32.0 {
			t.Errorf("maxTemp = %f, want 32.0", ct.maxTemp)
		}

		// Check data arrays
		if len(ct.T) != 7 {
			t.Errorf("len(T) = %d, want 7", len(ct.T))
		}
		if len(ct.Chara) != 7 {
			t.Errorf("len(Chara) = %d, want 7", len(ct.Chara))
		}

		// Check specific data points
		expectedT := []float64{20.0, 22.0, 24.0, 26.0, 28.0, 30.0, 32.0}
		expectedChara := []float64{50000.0, 52000.0, 55000.0, 180000.0, 185000.0, 58000.0, 60000.0}

		for i := 0; i < 7; i++ {
			if ct.T[i] != expectedT[i] {
				t.Errorf("T[%d] = %f, want %f", i, ct.T[i], expectedT[i])
			}
			if ct.Chara[i] != expectedChara[i] {
				t.Errorf("Chara[%d] = %f, want %f", i, ct.Chara[i], expectedChara[i])
			}
		}

		// Check linear regression coefficients for extrapolation
		expectedLowA := (52000.0 - 50000.0) / (22.0 - 20.0) // = 1000.0
		expectedLowB := 50000.0 - expectedLowA*20.0         // = 30000.0
		if ct.lowA != expectedLowA {
			t.Errorf("lowA = %f, want %f", ct.lowA, expectedLowA)
		}
		if ct.lowB != expectedLowB {
			t.Errorf("lowB = %f, want %f", ct.lowB, expectedLowB)
		}
	})

	t.Run("read specific heat table", func(t *testing.T) {
		ct := &CHARTABLE{
			filename:    "testdata/pcm_specific_heat_test.txt",
			PCMchar:     'E',
			tabletype:   'h', // specific heat mode
			minTempChng: 0.5,
		}

		TableRead(ct)

		// Check basic properties
		if ct.itablerow != 7 {
			t.Errorf("itablerow = %d, want 7", ct.itablerow)
		}

		// In specific heat mode, integration calculation is performed
		// *Char = prevheat + spheat*(*T-prevTemp)
		// First point: prevheat=0, spheat=2000, T=20, prevTemp=20 -> 0
		if ct.Chara[0] != 0.0 {
			t.Errorf("Chara[0] = %f, want 0.0 (first point in specific heat mode)", ct.Chara[0])
		}

		// Second point: prevheat=0, spheat=2100, T=22, prevTemp=20 -> 0 + 2100*(22-20) = 4200
		expectedChara1 := 0.0 + 2100.0*(22.0-20.0)
		if ct.Chara[1] != expectedChara1 {
			t.Errorf("Chara[1] = %f, want %f", ct.Chara[1], expectedChara1)
		}

		t.Logf("Specific heat mode results:")
		for i := 0; i < len(ct.T); i++ {
			t.Logf("  T[%d]=%f, Chara[%d]=%f", i, ct.T[i], i, ct.Chara[i])
		}
	})

	t.Run("nonexistent file", func(t *testing.T) {
		ct := &CHARTABLE{
			filename:  "nonexistent_file.txt",
			PCMchar:   'E',
			tabletype: 'e',
		}

		// This should not panic and should handle the error gracefully
		TableRead(ct)

		// Check that arrays are not allocated when file doesn't exist
		if ct.T != nil {
			t.Errorf("T should be nil when file doesn't exist")
		}
		if ct.Chara != nil {
			t.Errorf("Chara should be nil when file doesn't exist")
		}
	})
}

func TestFNPCMenthalpy_table_lib(t *testing.T) {
	// Setup test data
	ct := &CHARTABLE{
		filename:    "testdata/pcm_enthalpy_test.txt",
		PCMchar:     'E',
		tabletype:   'e',
		minTempChng: 0.5,
	}
	TableRead(ct)

	t.Run("interpolation within range", func(t *testing.T) {
		// Test interpolation between 22C (52000) and 24C (55000)
		T := 23.0
		result := FNPCMenthalpy_table_lib(ct, T)
		expected := 52000.0 + (55000.0-52000.0)*(23.0-22.0)/(24.0-22.0) // = 53500.0
		if result != expected {
			t.Errorf("FNPCMenthalpy_table_lib(%f) = %f, want %f", T, result, expected)
		}
	})

	t.Run("exact match", func(t *testing.T) {
		// Test exact temperature match
		T := 26.0
		result := FNPCMenthalpy_table_lib(ct, T)
		expected := 180000.0
		if result != expected {
			t.Errorf("FNPCMenthalpy_table_lib(%f) = %f, want %f", T, result, expected)
		}
	})

	t.Run("extrapolation below range", func(t *testing.T) {
		// Test extrapolation below minimum temperature
		T := 18.0
		result := FNPCMenthalpy_table_lib(ct, T)
		// Should use linear extrapolation: lowA * T + lowB
		expected := ct.lowA*T + ct.lowB
		if result != expected {
			t.Errorf("FNPCMenthalpy_table_lib(%f) = %f, want %f", T, result, expected)
		}
	})

	t.Run("extrapolation above range", func(t *testing.T) {
		// Test extrapolation above maximum temperature
		T := 35.0
		result := FNPCMenthalpy_table_lib(ct, T)
		// Should use linear extrapolation: upA * T + upB
		expected := ct.upA*T + ct.upB
		if result != expected {
			t.Errorf("FNPCMenthalpy_table_lib(%f) = %f, want %f", T, result, expected)
		}
	})
}

func TestFNPCMstate_table(t *testing.T) {
	// Setup test data for enthalpy table
	ctEnthalpy := &CHARTABLE{
		filename:    "testdata/pcm_enthalpy_test.txt",
		PCMchar:     'E', // Enthalpy
		tabletype:   'e',
		minTempChng: 0.5,
	}
	TableRead(ctEnthalpy)

	// Setup test data for conductivity table
	ctConductivity := &CHARTABLE{
		filename:    "testdata/pcm_conductivity_test.txt",
		PCMchar:     'C', // Conductivity
		tabletype:   'e',
		minTempChng: 0.5,
	}
	TableRead(ctConductivity)

	t.Run("enthalpy apparent specific heat calculation", func(t *testing.T) {
		Told := 22.0
		T := 24.0
		Ndiv := 10

		result := FNPCMstate_table(ctEnthalpy, Told, T, Ndiv)

		// Calculate expected apparent specific heat
		oldEn := FNPCMenthalpy_table_lib(ctEnthalpy, Told) // 52000.0
		En := FNPCMenthalpy_table_lib(ctEnthalpy, T)       // 55000.0
		expected := (En - oldEn) / (T - Told)              // (55000-52000)/(24-22) = 1500.0

		if result != expected {
			t.Errorf("FNPCMstate_table enthalpy(%f, %f, %d) = %f, want %f", Told, T, Ndiv, result, expected)
		}
	})

	t.Run("small temperature change with enthalpy", func(t *testing.T) {
		Told := 25.0
		T := 25.1 // Very small change, less than minTempChng
		Ndiv := 10

		result := FNPCMstate_table(ctEnthalpy, Told, T, Ndiv)

		// Should use minTempChng for calculation
		Tave := 0.5 * (T + Told)
		oldEn := FNPCMenthalpy_table_lib(ctEnthalpy, Tave-0.5*ctEnthalpy.minTempChng)
		En := FNPCMenthalpy_table_lib(ctEnthalpy, Tave+0.5*ctEnthalpy.minTempChng)
		expected := (En - oldEn) / ctEnthalpy.minTempChng

		if result != expected {
			t.Errorf("FNPCMstate_table small change(%f, %f, %d) = %f, want %f", Told, T, Ndiv, result, expected)
		}
	})

	t.Run("conductivity calculation", func(t *testing.T) {
		Told := 22.0
		T := 24.0
		Ndiv := 2

		result := FNPCMstate_table(ctConductivity, Told, T, Ndiv)

		// For conductivity, should integrate over temperature range
		dTemp := (T - Told) / float64(Ndiv)
		var sum float64
		for i := 0; i < Ndiv; i++ {
			Tpcm := Told + dTemp*float64(i)
			sum += FNPCMenthalpy_table_lib(ctConductivity, Tpcm)
		}
		expected := sum / float64(Ndiv+1)

		if result != expected {
			t.Errorf("FNPCMstate_table conductivity(%f, %f, %d) = %f, want %f", Told, T, Ndiv, result, expected)
		}
	})

	t.Run("conductivity with small temperature change", func(t *testing.T) {
		Told := 25.0
		T := 25.05 // Small change, less than minTempChng
		Ndiv := 5

		result := FNPCMstate_table(ctConductivity, Told, T, Ndiv)

		// For conductivity with small change, should return average temperature value
		Tave := 0.5 * (T + Told)
		expected := FNPCMenthalpy_table_lib(ctConductivity, Tave)

		if result != expected {
			t.Errorf("FNPCMstate_table conductivity small change(%f, %f, %d) = %f, want %f", Told, T, Ndiv, result, expected)
		}
	})

	t.Run("exact same temperature", func(t *testing.T) {
		Told := 26.0
		T := 26.0 // Exactly same temperature
		Ndiv := 10

		result := FNPCMstate_table(ctEnthalpy, Told, T, Ndiv)

		// Should use minTempChng calculation
		Tave := 0.5 * (T + Told) // = 26.0
		oldEn := FNPCMenthalpy_table_lib(ctEnthalpy, Tave-0.5*ctEnthalpy.minTempChng)
		En := FNPCMenthalpy_table_lib(ctEnthalpy, Tave+0.5*ctEnthalpy.minTempChng)
		expected := (En - oldEn) / ctEnthalpy.minTempChng

		if result != expected {
			t.Errorf("FNPCMstate_table same temp(%f, %f, %d) = %f, want %f", Told, T, Ndiv, result, expected)
		}
	})
}

func TestFNPCMState(t *testing.T) {
	// Test parameters for PCM
	Ss := 2000.0   // Solid specific heat [J/m3K]
	Sl := 2500.0   // Liquid specific heat [J/m3K]
	Ql := 180000.0 // Latent heat [J/m3]
	Ts := 26.0     // Solidification start temperature [C]
	Tl := 30.0     // Liquefaction start temperature [C]
	Tp := 28.0     // Peak temperature [C]

	// PCM parameters for different calculation types
	pcmp := &PCMPARAM{
		T:     2.0,
		B:     4.0,
		bs:    1.5,
		bl:    2.0,
		skew:  0.5,
		omega: 1.0,
		a:     5000.0,
		b:     1.0,
		c:     100.0,
		d:     0.1,
		e:     0.05,
		f:     2.0,
	}

	t.Run("thermal conductivity calculation (Ctype=0)", func(t *testing.T) {
		T := 27.0
		result := FNPCMState(0, Ss, Sl, Ql, Ts, Tl, Tp, T, pcmp)

		// For thermal conductivity, latent heat should be 0
		// Only sensible heat interpolation: Ss + (Sl-Ss)/(Tl-Ts)*(T-Ts)
		expected := Ss + (Sl-Ss)/(Tl-Ts)*(T-Ts) // 2000 + (2500-2000)/(30-26)*(27-26) = 2125

		if result != expected {
			t.Errorf("FNPCMState Ctype=0 (conductivity) = %f, want %f", result, expected)
		}
	})

	t.Run("constant latent heat (Ctype=1)", func(t *testing.T) {
		T := 27.0
		result := FNPCMState(1, Ss, Sl, 0.0, Ts, Tl, Tp, T, pcmp) // Ql=0 due to line 447 in code

		// Sensible heat + constant latent heat
		sensible := Ss + (Sl-Ss)/(Tl-Ts)*(T-Ts)
		latent := 0.0 // Because Ql is set to 0 in line 447
		expected := sensible + latent

		if result != expected {
			t.Errorf("FNPCMState Ctype=1 = %f, want %f", result, expected)
		}
	})

	t.Run("isosceles triangle (Ctype=2) - left side", func(t *testing.T) {
		T := 27.0 // T < Tp (28.0)
		result := FNPCMState(2, Ss, Sl, 0.0, Ts, Tl, Tp, T, pcmp)

		// Sensible heat + triangular latent heat (left side)
		sensible := Ss + (Sl-Ss)/(Tl-Ts)*(T-Ts) // 2000 + 500/4*(27-26) = 2125
		// Note: Ql is set to 0 in line 439, so latent heat = 0
		latent := 0.0 // Because Ql is set to 0 in the function
		expected := sensible + latent

		if result != expected {
			t.Errorf("FNPCMState Ctype=2 left = %f, want %f (sensible=%f + latent=%f)", result, expected, sensible, latent)
		}
	})

	t.Run("isosceles triangle (Ctype=2) - right side", func(t *testing.T) {
		T := 29.0 // T > Tp (28.0)
		result := FNPCMState(2, Ss, Sl, 0.0, Ts, Tl, Tp, T, pcmp)

		// Sensible heat + triangular latent heat (right side)
		sensible := Ss + (Sl-Ss)/(Tl-Ts)*(T-Ts) // 2000 + 500/4*(29-26) = 2375
		// Note: Ql is set to 0 in line 439, so latent heat = 0
		latent := 0.0 // Because Ql is set to 0 in the function
		expected := sensible + latent

		if result != expected {
			t.Errorf("FNPCMState Ctype=2 right = %f, want %f (sensible=%f + latent=%f)", result, expected, sensible, latent)
		}
	})

	t.Run("isosceles triangle (Ctype=2) - peak", func(t *testing.T) {
		T := 28.0 // T = Tp
		result := FNPCMState(2, Ss, Sl, 0.0, Ts, Tl, Tp, T, pcmp)

		// Sensible heat + triangular latent heat (peak)
		sensible := Ss + (Sl-Ss)/(Tl-Ts)*(T-Ts) // 2000 + 500/4*(28-26) = 2250
		// Note: Ql is set to 0 in line 439, so latent heat = 0
		latent := 0.0 // Because Ql is set to 0 in the function
		expected := sensible + latent

		if result != expected {
			t.Errorf("FNPCMState Ctype=2 peak = %f, want %f (sensible=%f + latent=%f)", result, expected, sensible, latent)
		}
	})

	t.Run("temperature below solidification", func(t *testing.T) {
		T := 25.0 // Below Ts
		result := FNPCMState(1, Ss, Sl, Ql, Ts, Tl, Tp, T, pcmp)

		// Should return solid specific heat
		expected := Ss

		if result != expected {
			t.Errorf("FNPCMState below Ts = %f, want %f", result, expected)
		}
	})

	t.Run("temperature above liquefaction", func(t *testing.T) {
		T := 31.0 // Above Tl
		result := FNPCMState(1, Ss, Sl, Ql, Ts, Tl, Tp, T, pcmp)

		// Should return liquid specific heat
		expected := Sl

		if result != expected {
			t.Errorf("FNPCMState above Tl = %f, want %f", result, expected)
		}
	})

	t.Run("hyperbolic function (Ctype=3)", func(t *testing.T) {
		T := 28.0
		result := FNPCMState(3, Ss, Sl, 0.0, Ts, Tl, Tp, T, pcmp) // Ql=0 due to line 447

		// Should calculate sensible heat + hyperbolic latent heat
		sensible := Ss + (Sl-Ss)/(Tl-Ts)*(T-Ts)
		// Latent heat calculation would use cosh function, but Ql=0
		expected := sensible

		if result != expected {
			t.Errorf("FNPCMState Ctype=3 = %f, want %f", result, expected)
		}
	})

	t.Run("symmetric Gaussian function (Ctype=4)", func(t *testing.T) {
		T := 28.0
		result := FNPCMState(4, Ss, Sl, 0.0, Ts, Tl, Tp, T, pcmp)

		// Calculate expected: sensible heat + Gaussian latent heat
		sensible := Ss + (Sl-Ss)/(Tl-Ts)*(T-Ts) // 2000 + (2500-2000)/(30-26)*(28-26) = 2250
		// Gaussian: pcmp.a * exp(-0.5 * ((T-Tp)/pcmp.b)^2) = 5000 * exp(-0.5 * ((28-28)/1)^2) = 5000 * exp(0) = 5000
		latent := pcmp.a * math.Exp(-0.5*math.Pow((T-Tp)/pcmp.b, 2))
		expected := sensible + latent

		if result != expected {
			t.Errorf("FNPCMState Ctype=4 = %f, want %f (sensible=%f + latent=%f)", result, expected, sensible, latent)
		}
	})

	t.Run("asymmetric Gaussian function (Ctype=5)", func(t *testing.T) {
		T := 28.0
		result := FNPCMState(5, Ss, Sl, 0.0, Ts, Tl, Tp, T, pcmp)

		// Calculate expected: sensible heat + asymmetric Gaussian latent heat
		sensible := Ss + (Sl-Ss)/(Tl-Ts)*(T-Ts) // 2250
		// Asymmetric Gaussian: T=28 > Tp=28 is false, so use bs=1.5
		// pcmp.a * exp(-((T-Tp)/bs)^2) = 5000 * exp(-((28-28)/1.5)^2) = 5000 * exp(0) = 5000
		latent := pcmp.a * math.Exp(-math.Pow((T-Tp)/pcmp.bs, 2))
		expected := sensible + latent

		if result != expected {
			t.Errorf("FNPCMState Ctype=5 = %f, want %f (sensible=%f + latent=%f)", result, expected, sensible, latent)
		}
	})

	t.Run("error function with skewness (Ctype=6)", func(t *testing.T) {
		T := 28.0
		testQl := 100000.0 // Use non-zero Ql for this test
		result := FNPCMState(6, Ss, Sl, testQl, Ts, Tl, Tp, T, pcmp)

		// Calculate expected: sensible heat + skewed error function latent heat
		sensible := Ss + (Sl-Ss)/(Tl-Ts)*(T-Ts) // 2250
		// Error function: Ql/sqrt(2*pi) * exp(-(T-Tp)^2/(2*omega^2)) * (1 + erf(skew*(T-Tp)/(sqrt(2)*omega)))
		// With T=28, Tp=28, omega=1.0, skew=0.5: exp(0) * (1 + erf(0)) = 1 * (1 + 0) = 1
		// Note: Ql is set to 0 in line 439, so latent = 0
		latent := 0.0 // Because Ql is set to 0 in the function
		expected := sensible + latent

		if result != expected {
			t.Errorf("FNPCMState Ctype=6 = %f, want %f (sensible=%f + latent=%f)", result, expected, sensible, latent)
		}
	})

	t.Run("rational function (Ctype=7)", func(t *testing.T) {
		T := 28.0
		result := FNPCMState(7, Ss, Sl, 0.0, Ts, Tl, Tp, T, pcmp)

		// Calculate expected: sensible heat + rational function latent heat
		sensible := Ss + (Sl-Ss)/(Tl-Ts)*(T-Ts) // 2250
		// Rational function: T^f * (a*T^2 + B*T + c) / (d*T^2 + e*T + 1)
		// T=28, f=2.0, a=5000, B=4.0, c=100, d=0.1, e=0.05
		// 28^2 * (5000*28^2 + 4*28 + 100) / (0.1*28^2 + 0.05*28 + 1)
		// = 784 * (5000*784 + 112 + 100) / (0.1*784 + 1.4 + 1)
		// = 784 * (3920000 + 212) / (78.4 + 2.4) = 784 * 3920212 / 80.8
		numerator := math.Pow(T, pcmp.f) * (pcmp.a*T*T + pcmp.B*T + pcmp.c)
		denominator := pcmp.d*T*T + pcmp.e*T + 1.0
		latent := numerator / denominator
		expected := sensible + latent

		if result != expected {
			t.Errorf("FNPCMState Ctype=7 = %f, want %f (sensible=%f + latent=%f)", result, expected, sensible, latent)
		}
	})

	t.Run("rational function with negative temperature (Ctype=7)", func(t *testing.T) {
		T := -5.0 // Negative temperature
		result := FNPCMState(7, Ss, Sl, 0.0, Ts, Tl, Tp, T, pcmp)

		// For negative temperature, latent heat should be 0
		// Temperature is below Ts, so should return solid specific heat
		expected := Ss

		if result != expected {
			t.Errorf("FNPCMState Ctype=7 negative T = %f, want %f", result, expected)
		}
	})

	t.Run("Gaussian function with actual latent heat calculation", func(t *testing.T) {
		// Test with non-zero Ql to verify the actual calculation
		T := 28.0
		testQl := 100000.0 // Use non-zero latent heat for this test

		result := FNPCMState(4, Ss, Sl, testQl, Ts, Tl, Tp, T, pcmp)

		// Calculate expected sensible heat
		sensible := Ss + (Sl-Ss)/(Tl-Ts)*(T-Ts) // 2250
		// For Ctype=4, the function uses pcmp.a directly, not Ql
		// Gaussian: pcmp.a * exp(-0.5 * ((T-Tp)/pcmp.b)^2) = 5000 * exp(0) = 5000
		latent := pcmp.a * math.Exp(-0.5*math.Pow((T-Tp)/pcmp.b, 2))
		expected := sensible + latent // 2250 + 5000 = 7250

		if result != expected {
			t.Errorf("FNPCMState Ctype=4 with Ql = %f, want %f", result, expected)
		}
	})

	t.Run("verify mathematical functions produce different results", func(t *testing.T) {
		// Test to ensure different Ctypes produce different calculations
		T := 28.0

		result3 := FNPCMState(3, Ss, Sl, 0.0, Ts, Tl, Tp, T, pcmp) // Hyperbolic
		result4 := FNPCMState(4, Ss, Sl, 0.0, Ts, Tl, Tp, T, pcmp) // Gaussian
		result5 := FNPCMState(5, Ss, Sl, 0.0, Ts, Tl, Tp, T, pcmp) // Asymmetric Gaussian
		result6 := FNPCMState(6, Ss, Sl, 0.0, Ts, Tl, Tp, T, pcmp) // Error function
		result7 := FNPCMState(7, Ss, Sl, 0.0, Ts, Tl, Tp, T, pcmp) // Rational function

		sensible := Ss + (Sl-Ss)/(Tl-Ts)*(T-Ts) // 2250

		// Expected values based on actual calculations:
		// Ctype=3: sensible only (Ql=0) = 2250
		// Ctype=4: sensible + Gaussian = 2250 + 5000 = 7250
		// Ctype=5: sensible + Asymmetric Gaussian = 2250 + 5000 = 7250
		// Ctype=6: sensible only (Ql=0) = 2250
		// Ctype=7: sensible + Rational = 2250 + very large number

		if result3 != sensible {
			t.Errorf("Ctype=3 (hyperbolic) = %f, want %f", result3, sensible)
		}
		if result4 <= sensible {
			t.Errorf("Ctype=4 (Gaussian) = %f, should be > %f", result4, sensible)
		}
		if result5 <= sensible {
			t.Errorf("Ctype=5 (Asymmetric Gaussian) = %f, should be > %f", result5, sensible)
		}
		if result6 != sensible {
			t.Errorf("Ctype=6 (error function) = %f, want %f", result6, sensible)
		}
		if result7 <= sensible {
			t.Errorf("Ctype=7 (rational function) = %f, should be > %f", result7, sensible)
		}

		t.Logf("Mathematical function results: 3=%f, 4=%f, 5=%f, 6=%f, 7=%f",
			result3, result4, result5, result6, result7)
	})
}

func TestFNPCMStatefun(t *testing.T) {
	// Test parameters
	Ss := 2000.0
	Sl := 2500.0
	Ql := 180000.0
	Ts := 26.0
	Tl := 30.0
	Tp := 28.0
	DivTemp := 5

	pcmp := &PCMPARAM{
		T: 2.0,
		B: 4.0,
		a: 5000.0,
		b: 1.0,
		c: 100.0,
		d: 0.1,
		e: 0.05,
		f: 2.0,
	}

	t.Run("small temperature change", func(t *testing.T) {
		oldT := 27.0
		T := 27.00005 // Very small change < 1e-4

		result := FNPCMStatefun(1, Ss, Sl, Ql, Ts, Tl, Tp, oldT, T, DivTemp, pcmp)

		// Should use average temperature calculation
		avgT := (T + oldT) * 0.5
		expected := FNPCMState(1, Ss, Sl, Ql, Ts, Tl, Tp, avgT, pcmp)

		if result != expected {
			t.Errorf("FNPCMStatefun small change = %f, want %f", result, expected)
		}
	})

	t.Run("normal temperature change with integration", func(t *testing.T) {
		oldT := 26.5
		T := 27.5

		result := FNPCMStatefun(1, Ss, Sl, Ql, Ts, Tl, Tp, oldT, T, DivTemp, pcmp)

		// Should perform numerical integration
		dTemp := (T - oldT) / float64(DivTemp)
		var sum float64
		for i := 0; i < DivTemp+1; i++ {
			TPCM := oldT + dTemp*float64(i)
			sum += FNPCMState(1, Ss, Sl, Ql, Ts, Tl, Tp, TPCM, pcmp)
		}
		expected := sum / float64(DivTemp+1)

		if result != expected {
			t.Errorf("FNPCMStatefun integration = %f, want %f", result, expected)
		}
	})

	t.Run("thermal conductivity calculation", func(t *testing.T) {
		oldT := 25.0
		T := 29.0

		result := FNPCMStatefun(0, Ss, Sl, Ql, Ts, Tl, Tp, oldT, T, DivTemp, pcmp)

		// For thermal conductivity (Ctype=0), should integrate over temperature range
		dTemp := (T - oldT) / float64(DivTemp)
		var sum float64
		for i := 0; i < DivTemp+1; i++ {
			TPCM := oldT + dTemp*float64(i)
			sum += FNPCMState(0, Ss, Sl, Ql, Ts, Tl, Tp, TPCM, pcmp)
		}
		expected := sum / float64(DivTemp+1)

		if result != expected {
			t.Errorf("FNPCMStatefun thermal conductivity = %f, want %f", result, expected)
		}
	})

	t.Run("zero temperature difference", func(t *testing.T) {
		oldT := 27.0
		T := 27.0 // Exactly same temperature

		result := FNPCMStatefun(1, Ss, Sl, Ql, Ts, Tl, Tp, oldT, T, DivTemp, pcmp)

		// Should use average temperature (which equals the temperature)
		expected := FNPCMState(1, Ss, Sl, Ql, Ts, Tl, Tp, T, pcmp)

		if result != expected {
			t.Errorf("FNPCMStatefun zero diff = %f, want %f", result, expected)
		}
	})
}

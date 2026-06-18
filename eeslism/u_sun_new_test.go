package eeslism

import (
	"math"
	"testing"
)

func TestSunint(t *testing.T) {
	// Store original values
	origLat := Lat
	origSlat := Slat
	origClat := Clat
	origTlat := Tlat
	origIsc := Isc

	// Restore original values after test
	defer func() {
		Lat = origLat
		Slat = origSlat
		Clat = origClat
		Tlat = origTlat
		Isc = origIsc
	}()

	t.Run("SI unit initialization", func(t *testing.T) {
		// Set test latitude (Tokyo: approximately 35.7°N)
		Lat = 35.7

		Sunint()

		// Check trigonometric values
		expectedSlat := math.Sin(35.7 * math.Pi / 180.0)
		expectedClat := math.Cos(35.7 * math.Pi / 180.0)
		expectedTlat := math.Tan(35.7 * math.Pi / 180.0)
		expectedIsc := 1370.0 // SI unit

		tolerance := 1e-6
		if math.Abs(Slat-expectedSlat) > tolerance {
			t.Errorf("Sunint() Slat = %v, want %v", Slat, expectedSlat)
		}
		if math.Abs(Clat-expectedClat) > tolerance {
			t.Errorf("Sunint() Clat = %v, want %v", Clat, expectedClat)
		}
		if math.Abs(Tlat-expectedTlat) > tolerance {
			t.Errorf("Sunint() Tlat = %v, want %v", Tlat, expectedTlat)
		}
		if math.Abs(Isc-expectedIsc) > tolerance {
			t.Errorf("Sunint() Isc = %v, want %v", Isc, expectedIsc)
		}
	})
}

func TestFNDecl(t *testing.T) {
	tests := []struct {
		name      string
		day       int
		expected  float64
		tolerance float64
	}{
		{
			name:      "spring equinox (March 21, day 80)",
			day:       80,
			expected:  0.0, // Declination should be close to 0 at equinox
			tolerance: 0.1,
		},
		{
			name:      "summer solstice (June 21, day 172)",
			day:       172,
			expected:  0.409, // Maximum declination ~23.45° = 0.409 rad
			tolerance: 0.05,
		},
		{
			name:      "autumn equinox (September 23, day 266)",
			day:       266,
			expected:  0.0, // Declination should be close to 0 at equinox
			tolerance: 0.1,
		},
		{
			name:      "winter solstice (December 21, day 355)",
			day:       355,
			expected:  -0.409, // Minimum declination ~-23.45° = -0.409 rad
			tolerance: 0.05,
		},
		{
			name:      "January 1 (day 1)",
			day:       1,
			expected:  -0.384, // Approximate declination in early January
			tolerance: 0.05,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FNDecl(tt.day)
			if math.Abs(result-tt.expected) > tt.tolerance {
				t.Errorf("FNDecl(%d) = %v rad (%.1f°), want %v rad (%.1f°) ± %v",
					tt.day, result, result*180/math.Pi, tt.expected, tt.expected*180/math.Pi, tt.tolerance)
			}

			// Check that declination is within physical bounds [-23.45°, +23.45°]
			maxDecl := 0.41 // ~23.5° in radians
			if math.Abs(result) > maxDecl {
				t.Errorf("FNDecl(%d) = %v rad is outside physical bounds [-%v, +%v]",
					tt.day, result, maxDecl, maxDecl)
			}
		})
	}
}

func TestFNE(t *testing.T) {
	tests := []struct {
		name      string
		day       int
		expected  float64
		tolerance float64
	}{
		{
			name:      "April 15 (day 105) - near zero",
			day:       105,
			expected:  0.0, // Should be close to zero
			tolerance: 0.1,
		},
		{
			name:      "June 14 (day 165) - near zero",
			day:       165,
			expected:  0.0, // Should be close to zero
			tolerance: 0.1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FNE(tt.day)
			if math.Abs(result-tt.expected) > tt.tolerance {
				t.Errorf("FNE(%d) = %v hours (%.1f min), want %v hours (%.1f min) ± %v",
					tt.day, result, result*60, tt.expected, tt.expected*60, tt.tolerance)
			}

			// Check that equation of time is within physical bounds [-0.3, +0.3] hours
			maxE := 0.3 // ~18 minutes
			if math.Abs(result) > maxE {
				t.Errorf("FNE(%d) = %v hours is outside physical bounds [-%v, +%v]",
					tt.day, result, maxE, maxE)
			}
		})
	}
}

func TestFNSro(t *testing.T) {
	// Initialize solar constants
	Sunint()

	tests := []struct {
		name      string
		day       int
		expected  float64
		tolerance float64
	}{
		{
			name:      "January 1 (day 1) - perihelion",
			day:       1,
			expected:  1415.0, // Isc * (1 + 0.033*cos(2π*1/365)) ≈ 1370 * 1.033
			tolerance: 10.0,
		},
		{
			name:      "July 1 (day 182) - aphelion",
			day:       182,
			expected:  1325.0, // Isc * (1 + 0.033*cos(2π*182/365)) ≈ 1370 * 0.967
			tolerance: 10.0,
		},
		{
			name:      "April 1 (day 91) - intermediate",
			day:       91,
			expected:  1370.0, // Should be close to Isc
			tolerance: 20.0,
		},
		{
			name:      "October 1 (day 274) - intermediate",
			day:       274,
			expected:  1370.0, // Should be close to Isc
			tolerance: 20.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FNSro(tt.day)
			if math.Abs(result-tt.expected) > tt.tolerance {
				t.Errorf("FNSro(%d) = %v, want %v ± %v", tt.day, result, tt.expected, tt.tolerance)
			}

			// Check that result is positive and reasonable
			if result < 1300 || result > 1450 {
				t.Errorf("FNSro(%d) = %v is outside reasonable bounds [1300, 1450] W/m²", tt.day, result)
			}
		})
	}
}

func TestFNTtas(t *testing.T) {
	tests := []struct {
		name      string
		tt        float64
		e         float64
		expected  float64
		tolerance float64
	}{
		{
			name:      "noon with zero equation of time",
			tt:        12.0,
			e:         0.0,
			expected:  18.666,
			tolerance: 0.001,
		},
		{
			name:      "noon with positive equation of time",
			tt:        12.0,
			e:         0.25, // +15 minutes
			expected:  18.916,
			tolerance: 0.001,
		},
		{
			name:      "noon with negative equation of time",
			tt:        12.0,
			e:         -0.25, // -15 minutes
			expected:  18.416,
			tolerance: 0.001,
		},
		{
			name:      "morning time",
			tt:        9.0,
			e:         0.1,
			expected:  15.766,
			tolerance: 0.001,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FNTtas(tt.tt, tt.e, 135, 35)
			if math.Abs(result-tt.expected) > tt.tolerance {
				t.Errorf("FNTtas(%v, %v) = %v, want %v ± %v", tt.tt, tt.e, result, tt.expected, tt.tolerance)
			}
		})
	}
}

func TestFNTt(t *testing.T) {
	tests := []struct {
		name      string
		ttas      float64
		e         float64
		expected  float64
		tolerance float64
	}{
		{
			name:      "apparent solar noon with zero equation of time",
			ttas:      12.0,
			e:         0.0,
			expected:  5.333,
			tolerance: 0.001,
		},
		{
			name:      "apparent solar noon with positive equation of time",
			ttas:      12.25,
			e:         0.25,
			expected:  5.333,
			tolerance: 0.001,
		},
		{
			name:      "apparent solar noon with negative equation of time",
			ttas:      11.75,
			e:         -0.25,
			expected:  5.333,
			tolerance: 0.001,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FNTt(tt.ttas, tt.e, 135, 35)
			if math.Abs(result-tt.expected) > tt.tolerance {
				t.Errorf("FNTt(%v, %v) = %v, want %v ± %v", tt.ttas, tt.e, result, tt.expected, tt.tolerance)
			}
		})
	}
}

func TestSolarCalculationConsistency(t *testing.T) {
	// Test consistency between FNTtas and FNTt (they should be inverse operations)
	t.Run("FNTtas and FNTt consistency", func(t *testing.T) {
		testTimes := []float64{8.0, 10.0, 12.0, 14.0, 16.0}
		testE := []float64{-0.2, -0.1, 0.0, 0.1, 0.2}

		for _, tt := range testTimes {
			for _, e := range testE {
				ttas := FNTtas(tt, e, 135, 35)
				ttBack := FNTt(ttas, e, 135, 35)

				if math.Abs(tt-ttBack) > 1e-10 {
					t.Errorf("Inconsistency: Tt=%v, E=%v -> Ttas=%v -> Tt=%v", tt, e, ttas, ttBack)
				}
			}
		}
	})

	// Test that declination varies smoothly throughout the year
	t.Run("declination smoothness", func(t *testing.T) {
		var prevDecl float64
		for day := 1; day <= 365; day++ {
			decl := FNDecl(day)
			if day > 1 {
				// Check that daily change is reasonable (less than 0.01 radians ≈ 0.6°)
				dailyChange := math.Abs(decl - prevDecl)
				if dailyChange > 0.01 {
					t.Errorf("Large daily declination change on day %d: %v rad (%.2f°)",
						day, dailyChange, dailyChange*180/math.Pi)
				}
			}
			prevDecl = decl
		}
	})
}

func TestFNTtd(t *testing.T) {
	// FNTtd uses the global Tlat set by Sunint().
	// Save and restore globals so this test does not affect others.
	origLat := Lat
	origTlat := Tlat
	defer func() {
		Lat = origLat
		Tlat = origTlat
	}()

	t.Run("Tokyo latitude spring equinox 12h", func(t *testing.T) {
		Lat = 35.7
		Sunint()
		// Decl=0 → Cws=0 → Ttd = 7.6394 * acos(0) ≈ 12.0
		Ttd := FNTtd(0.0)
		if math.Abs(Ttd-12.0) > 0.01 {
			t.Errorf("equinox day length = %.4f, want ≈12.0", Ttd)
		}
	})

	t.Run("Tokyo latitude summer solstice ~14.42h", func(t *testing.T) {
		Lat = 35.7
		Sunint()
		Decl := 23.45 * math.Pi / 180.0
		Ttd := FNTtd(Decl)
		// Expected ≈ 14.42 hours (calculated via python)
		if math.Abs(Ttd-14.42) > 0.05 {
			t.Errorf("summer solstice day length = %.4f, want ≈14.42", Ttd)
		}
	})

	t.Run("Tokyo latitude winter solstice ~9.58h", func(t *testing.T) {
		Lat = 35.7
		Sunint()
		Decl := -23.45 * math.Pi / 180.0
		Ttd := FNTtd(Decl)
		// Expected ≈ 9.58 hours
		if math.Abs(Ttd-9.58) > 0.05 {
			t.Errorf("winter solstice day length = %.4f, want ≈9.58", Ttd)
		}
	})

	t.Run("Polar day Lat=80 Decl=20deg returns 24h", func(t *testing.T) {
		Lat = 80.0
		Sunint()
		Decl := 20.0 * math.Pi / 180.0
		Ttd := FNTtd(Decl)
		// Cws = -tan(80°)*tan(20°) ≈ -2.06 ≤ -1 → white night → 24h
		if Ttd != 24.0 {
			t.Errorf("polar day: got %.4f, want 24.0", Ttd)
		}
	})

	t.Run("Polar night Lat=80 Decl=-20deg returns 0h", func(t *testing.T) {
		Lat = 80.0
		Sunint()
		Decl := -20.0 * math.Pi / 180.0
		Ttd := FNTtd(Decl)
		// Cws = -tan(80°)*tan(-20°) ≈ +2.06 ≥ 1 → polar night → 0h
		if Ttd != 0.0 {
			t.Errorf("polar night: got %.4f, want 0.0", Ttd)
		}
	})

	t.Run("summer > winter day length", func(t *testing.T) {
		Lat = 35.7
		Sunint()
		summer := FNTtd(23.45 * math.Pi / 180.0)
		winter := FNTtd(-23.45 * math.Pi / 180.0)
		if summer <= winter {
			t.Errorf("summer (%.4f) should be longer than winter (%.4f)", summer, winter)
		}
	})
}

func TestSrdclr(t *testing.T) {
	const tol = 0.5 // W/m² tolerance

	t.Run("normal conditions Sh=sin(60deg)", func(t *testing.T) {
		Io := 1370.0
		P := 0.75
		Sh := math.Sin(60.0 * math.Pi / 180.0) // ≈ 0.8660
		var Idn, Isky float64
		Srdclr(Io, P, Sh, &Idn, &Isky)
		// Expected: Idn≈982.77, Isky≈83.66 (pre-calculated)
		if math.Abs(Idn-982.77) > tol {
			t.Errorf("Idn = %.4f, want ≈982.77", Idn)
		}
		if math.Abs(Isky-83.66) > tol {
			t.Errorf("Isky = %.4f, want ≈83.66", Isky)
		}
	})

	t.Run("Sh at threshold boundary returns zero", func(t *testing.T) {
		var Idn, Isky float64
		Srdclr(1370.0, 0.75, 0.0005, &Idn, &Isky)
		if Idn != 0.0 || Isky != 0.0 {
			t.Errorf("Sh<=0.001: got Idn=%.4f Isky=%.4f, want both 0", Idn, Isky)
		}
	})

	t.Run("Sh=0 returns zero", func(t *testing.T) {
		var Idn, Isky float64
		Srdclr(1370.0, 0.75, 0.0, &Idn, &Isky)
		if Idn != 0.0 || Isky != 0.0 {
			t.Errorf("Sh=0: got Idn=%.4f Isky=%.4f, want both 0", Idn, Isky)
		}
	})

	t.Run("physical bounds: 0 <= Idn <= Io and Isky >= 0", func(t *testing.T) {
		Io := 1370.0
		for _, Sh := range []float64{0.1, 0.3, 0.5, 0.7, 0.9, 1.0} {
			var Idn, Isky float64
			Srdclr(Io, 0.75, Sh, &Idn, &Isky)
			if Idn < 0 || Idn > Io {
				t.Errorf("Sh=%.1f: Idn=%.4f out of range [0, %.1f]", Sh, Idn, Io)
			}
			if Isky < 0 {
				t.Errorf("Sh=%.1f: Isky=%.4f < 0", Sh, Isky)
			}
		}
	})
}

func TestDnsky(t *testing.T) {
	const tol = 0.5 // W/m² tolerance

	t.Run("low Kt branch (Ihol=600 Sh=0.7)", func(t *testing.T) {
		Io := 1370.0
		Ihol := 600.0
		Sh := 0.7
		var Idn, Isky float64
		Dnsky(Io, Ihol, Sh, &Idn, &Isky)
		// Kt=0.6257 < threshold=0.7533 → low Kt branch
		// Expected: Idn≈507.91, Isky≈244.46 (pre-calculated)
		if math.Abs(Idn-507.91) > tol {
			t.Errorf("low Kt: Idn=%.4f, want ≈507.91", Idn)
		}
		if math.Abs(Isky-244.46) > tol {
			t.Errorf("low Kt: Isky=%.4f, want ≈244.46", Isky)
		}
	})

	t.Run("high Kt branch (Ihol=800 Sh=0.7)", func(t *testing.T) {
		Io := 1370.0
		Ihol := 800.0
		Sh := 0.7
		var Idn, Isky float64
		Dnsky(Io, Ihol, Sh, &Idn, &Isky)
		// Kt=0.8342 >= threshold=0.7533 → high Kt branch
		// Expected: Idn≈1045.19, Isky≈68.37 (pre-calculated)
		if math.Abs(Idn-1045.19) > tol {
			t.Errorf("high Kt: Idn=%.4f, want ≈1045.19", Idn)
		}
		if math.Abs(Isky-68.37) > tol {
			t.Errorf("high Kt: Isky=%.4f, want ≈68.37", Isky)
		}
	})

	t.Run("energy conservation: Idn*Sh + Isky == Ihol", func(t *testing.T) {
		Io := 1370.0
		for _, tc := range []struct{ Ihol, Sh float64 }{
			{600.0, 0.7},
			{800.0, 0.7},
			{100.0, 0.7},
		} {
			var Idn, Isky float64
			Dnsky(Io, tc.Ihol, tc.Sh, &Idn, &Isky)
			got := Idn*tc.Sh + Isky
			if math.Abs(got-tc.Ihol) > 0.01 {
				t.Errorf("Ihol=%.0f Sh=%.1f: Idn*Sh+Isky=%.6f, want %.6f", tc.Ihol, tc.Sh, got, tc.Ihol)
			}
		}
	})

	t.Run("Sh at threshold returns Idn=0 Isky=Ihol", func(t *testing.T) {
		var Idn, Isky float64
		Ihol := 50.0
		Dnsky(1370.0, Ihol, 0.0005, &Idn, &Isky)
		if Idn != 0.0 {
			t.Errorf("Sh<=0.001: Idn=%.4f, want 0", Idn)
		}
		if Isky != Ihol {
			t.Errorf("Sh<=0.001: Isky=%.4f, want %.4f", Isky, Ihol)
		}
	})
}

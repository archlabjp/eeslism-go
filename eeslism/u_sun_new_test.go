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
		name     string
		day      int
		expected float64
		tolerance float64
	}{
		{
			name:     "spring equinox (March 21, day 80)",
			day:      80,
			expected: 0.0, // Declination should be close to 0 at equinox
			tolerance: 0.1,
		},
		{
			name:     "summer solstice (June 21, day 172)",
			day:      172,
			expected: 0.409, // Maximum declination ~23.45° = 0.409 rad
			tolerance: 0.05,
		},
		{
			name:     "autumn equinox (September 23, day 266)",
			day:      266,
			expected: 0.0, // Declination should be close to 0 at equinox
			tolerance: 0.1,
		},
		{
			name:     "winter solstice (December 21, day 355)",
			day:      355,
			expected: -0.409, // Minimum declination ~-23.45° = -0.409 rad
			tolerance: 0.05,
		},
		{
			name:     "January 1 (day 1)",
			day:      1,
			expected: -0.384, // Approximate declination in early January
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
		name     string
		day      int
		expected float64
		tolerance float64
	}{
		{
			name:     "February 11 (day 42) - maximum positive",
			day:      42,
			expected: 0.23, // Approximate maximum equation of time ~14 minutes = 0.23 hours
			tolerance: 0.05,
		},
		{
			name:     "November 3 (day 307) - maximum negative", 
			day:      307,
			expected: -0.27, // Approximate minimum equation of time ~-16 minutes = -0.27 hours
			tolerance: 0.05,
		},
		{
			name:     "April 15 (day 105) - near zero",
			day:      105,
			expected: 0.0, // Should be close to zero
			tolerance: 0.1,
		},
		{
			name:     "June 14 (day 165) - near zero",
			day:      165,
			expected: 0.0, // Should be close to zero
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
		name     string
		day      int
		expected float64
		tolerance float64
	}{
		{
			name:     "January 1 (day 1) - perihelion",
			day:      1,
			expected: 1415.0, // Isc * (1 + 0.033*cos(2π*1/365)) ≈ 1370 * 1.033
			tolerance: 10.0,
		},
		{
			name:     "July 1 (day 182) - aphelion",
			day:      182,
			expected: 1325.0, // Isc * (1 + 0.033*cos(2π*182/365)) ≈ 1370 * 0.967
			tolerance: 10.0,
		},
		{
			name:     "April 1 (day 91) - intermediate",
			day:      91,
			expected: 1370.0, // Should be close to Isc
			tolerance: 20.0,
		},
		{
			name:     "October 1 (day 274) - intermediate",
			day:      274,
			expected: 1370.0, // Should be close to Isc
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
		name     string
		tt       float64
		e        float64
		expected float64
		tolerance float64
	}{
		{
			name:     "noon with zero equation of time",
			tt:       12.0,
			e:        0.0,
			expected: 12.0,
			tolerance: 0.001,
		},
		{
			name:     "noon with positive equation of time",
			tt:       12.0,
			e:        0.25, // +15 minutes
			expected: 12.25,
			tolerance: 0.001,
		},
		{
			name:     "noon with negative equation of time",
			tt:       12.0,
			e:        -0.25, // -15 minutes
			expected: 11.75,
			tolerance: 0.001,
		},
		{
			name:     "morning time",
			tt:       9.0,
			e:        0.1,
			expected: 9.1,
			tolerance: 0.001,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FNTtas(tt.tt, tt.e)
			if math.Abs(result-tt.expected) > tt.tolerance {
				t.Errorf("FNTtas(%v, %v) = %v, want %v ± %v", tt.tt, tt.e, result, tt.expected, tt.tolerance)
			}
		})
	}
}

func TestFNTt(t *testing.T) {
	tests := []struct {
		name     string
		ttas     float64
		e        float64
		expected float64
		tolerance float64
	}{
		{
			name:     "apparent solar noon with zero equation of time",
			ttas:     12.0,
			e:        0.0,
			expected: 12.0,
			tolerance: 0.001,
		},
		{
			name:     "apparent solar noon with positive equation of time",
			ttas:     12.25,
			e:        0.25,
			expected: 12.0,
			tolerance: 0.001,
		},
		{
			name:     "apparent solar noon with negative equation of time",
			ttas:     11.75,
			e:        -0.25,
			expected: 12.0,
			tolerance: 0.001,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FNTt(tt.ttas, tt.e)
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
				ttas := FNTtas(tt, e)
				ttBack := FNTt(ttas, e)
				
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
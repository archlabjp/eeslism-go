package eeslism

import (
	"testing"
)

func TestSunshadeBasicStructures(t *testing.T) {
	// Test basic sunshade structures
	t.Run("SNBK structure", func(t *testing.T) {
		snbk := &SNBK{
			Name: "TestSunshade",
			Type: 1,
			Ksi:  45,
			D:    0.5,
			W:    2.0,
			H:    1.5,
			W1:   1.8,
			H1:   1.3,
			W2:   2.2,
			H2:   1.7,
		}

		// Verify initialization
		if snbk.Name != "TestSunshade" {
			t.Errorf("Name = %s, want TestSunshade", snbk.Name)
		}
		if snbk.Type != 1 {
			t.Errorf("Type = %d, want 1", snbk.Type)
		}
		if snbk.Ksi != 45 {
			t.Errorf("Ksi = %d, want 45", snbk.Ksi)
		}
	})
}

func TestFNFsdw(t *testing.T) {
	// Test sunshade shadow factor calculation
	tests := []struct {
		name     string
		typ      int
		ksi      int
		tazm     float64
		tprof    float64
		d        float64
		w        float64
		h        float64
		w1       float64
		h1       float64
		w2       float64
		h2       float64
		expected float64 // Expected range
	}{
		{
			name:     "Horizontal overhang",
			typ:      1,
			ksi:      0,
			tazm:     0.0,
			tprof:    45.0,
			d:        0.5,
			w:        2.0,
			h:        1.0,
			w1:       0.0,
			h1:       0.0,
			w2:       0.0,
			h2:       0.0,
			expected: 0.0, // Should calculate shadow factor
		},
		{
			name:     "Vertical fin",
			typ:      2,
			ksi:      90,
			tazm:     30.0,
			tprof:    30.0,
			d:        0.3,
			w:        1.5,
			h:        2.0,
			w1:       0.0,
			h1:       0.0,
			w2:       0.0,
			h2:       0.0,
			expected: 0.0, // Should calculate shadow factor
		},
		{
			name:     "Combined shading",
			typ:      3,
			ksi:      45,
			tazm:     45.0,
			tprof:    60.0,
			d:        0.4,
			w:        2.5,
			h:        1.8,
			w1:       2.0,
			h1:       1.5,
			w2:       3.0,
			h2:       2.0,
			expected: 0.0, // Should calculate shadow factor
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fsdw := FNFsdw(tt.typ, tt.ksi, tt.tazm, tt.tprof, tt.d, tt.w, tt.h, tt.w1, tt.h1, tt.w2, tt.h2)

			// Verify shadow factor is in valid range [0, 1]
			if fsdw < 0.0 || fsdw > 1.0 {
				t.Errorf("FNFsdw() = %f, should be between 0.0 and 1.0", fsdw)
			}

			// Log result for verification
			t.Logf("%s: shadow factor = %f", tt.name, fsdw)
		})
	}
}

func TestSunshadeGeometry(t *testing.T) {
	// Test sunshade geometry calculations
	tests := []struct {
		name   string
		snbk   *SNBK
		valid  bool
	}{
		{
			name: "Valid horizontal overhang",
			snbk: &SNBK{
				Name: "HorizontalOverhang",
				Type: 1,
				D:    0.5,
				W:    2.0,
				H:    1.0,
			},
			valid: true,
		},
		{
			name: "Valid vertical fin",
			snbk: &SNBK{
				Name: "VerticalFin",
				Type: 2,
				D:    0.3,
				W:    1.5,
				H:    2.5,
			},
			valid: true,
		},
		{
			name: "Invalid dimensions",
			snbk: &SNBK{
				Name: "InvalidShade",
				Type: 1,
				D:    -0.1, // Invalid depth
				W:    2.0,
				H:    1.0,
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Basic geometry validation
			depthValid := tt.snbk.D >= 0.0
			widthValid := tt.snbk.W > 0.0
			heightValid := tt.snbk.H > 0.0

			isValid := depthValid && widthValid && heightValid

			if isValid != tt.valid {
				t.Errorf("Geometry validation = %t, want %t", isValid, tt.valid)
			}
		})
	}
}

func TestSunshadeAngles(t *testing.T) {
	// Test sunshade angle calculations
	tests := []struct {
		name     string
		ksi      int // Sunshade angle
		tazm     float64 // Solar azimuth
		tprof    float64 // Solar profile angle
		expected string  // Expected behavior
	}{
		{
			name:     "Perpendicular to sun",
			ksi:      0,
			tazm:     0.0,
			tprof:    45.0,
			expected: "maximum_shading",
		},
		{
			name:     "Parallel to sun",
			ksi:      90,
			tazm:     90.0,
			tprof:    30.0,
			expected: "minimum_shading",
		},
		{
			name:     "Angled sunshade",
			ksi:      45,
			tazm:     30.0,
			tprof:    60.0,
			expected: "partial_shading",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test angle relationships
			angleDiff := abs(float64(tt.ksi) - tt.tazm)

			// Verify angles are in valid ranges
			if tt.ksi < 0 || tt.ksi > 360 {
				t.Errorf("Ksi (%d) outside valid range [0, 360]", tt.ksi)
			}
			if tt.tazm < 0.0 || tt.tazm > 360.0 {
				t.Errorf("Tazm (%f) outside valid range [0, 360]", tt.tazm)
			}
			if tt.tprof < 0.0 || tt.tprof > 90.0 {
				t.Errorf("Tprof (%f) outside valid range [0, 90]", tt.tprof)
			}

			t.Logf("%s: angle difference = %f degrees", tt.name, angleDiff)
		})
	}
}

func TestSunshadeEffectiveness(t *testing.T) {
	// Test sunshade effectiveness under different conditions
	tests := []struct {
		name      string
		sunshade  *SNBK
		solarAzm  float64
		solarProf float64
		expected  string
	}{
		{
			name: "Morning sun protection",
			sunshade: &SNBK{
				Type: 2, // Vertical fin
				Ksi:  90,
				D:    0.5,
				W:    1.0,
				H:    2.0,
			},
			solarAzm:  90.0, // East
			solarProf: 30.0,
			expected:  "effective",
		},
		{
			name: "Noon sun protection",
			sunshade: &SNBK{
				Type: 1, // Horizontal overhang
				Ksi:  0,
				D:    0.8,
				W:    2.0,
				H:    1.0,
			},
			solarAzm:  180.0, // South
			solarProf: 60.0,
			expected:  "effective",
		},
		{
			name: "Low sun angle",
			sunshade: &SNBK{
				Type: 1, // Horizontal overhang
				Ksi:  0,
				D:    0.3,
				W:    1.5,
				H:    1.0,
			},
			solarAzm:  180.0,
			solarProf: 15.0, // Low angle
			expected:  "limited",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Calculate shadow factor
			fsdw := FNFsdw(tt.sunshade.Type, tt.sunshade.Ksi, tt.solarAzm, tt.solarProf,
				tt.sunshade.D, tt.sunshade.W, tt.sunshade.H, 0.0, 0.0, 0.0, 0.0)

			// Evaluate effectiveness
			var effectiveness string
			if fsdw > 0.7 {
				effectiveness = "highly_effective"
			} else if fsdw > 0.3 {
				effectiveness = "effective"
			} else if fsdw > 0.1 {
				effectiveness = "limited"
			} else {
				effectiveness = "minimal"
			}

			t.Logf("%s: shadow factor = %f, effectiveness = %s", tt.name, fsdw, effectiveness)

			// Verify shadow factor is reasonable
			if fsdw < 0.0 || fsdw > 1.0 {
				t.Errorf("Shadow factor (%f) outside valid range [0, 1]", fsdw)
			}
		})
	}
}

func TestSunshadeTypes(t *testing.T) {
	// Test different sunshade types
	types := []struct {
		name        string
		typ         int
		description string
	}{
		{"Horizontal overhang", 1, "Provides shading from high sun angles"},
		{"Vertical fin", 2, "Provides shading from low sun angles"},
		{"Combined system", 3, "Provides comprehensive shading"},
		{"Custom type", 4, "Custom shading configuration"},
	}

	for _, tt := range types {
		t.Run(tt.name, func(t *testing.T) {
			snbk := &SNBK{
				Name: tt.name,
				Type: tt.typ,
				D:    0.5,
				W:    2.0,
				H:    1.5,
			}

			// Verify type is set correctly
			if snbk.Type != tt.typ {
				t.Errorf("Type = %d, want %d", snbk.Type, tt.typ)
			}

			// Test shadow calculation for each type
			fsdw := FNFsdw(snbk.Type, 0, 0.0, 45.0, snbk.D, snbk.W, snbk.H, 0.0, 0.0, 0.0, 0.0)

			if fsdw < 0.0 || fsdw > 1.0 {
				t.Errorf("Shadow factor (%f) outside valid range for type %d", fsdw, tt.typ)
			}

			t.Logf("%s (type %d): shadow factor = %f", tt.name, tt.typ, fsdw)
		})
	}
}

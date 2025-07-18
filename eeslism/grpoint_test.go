package eeslism

import (
	"math"
	"testing"
)

func TestGRGPOINT(t *testing.T) {
	tests := []struct {
		name     string
		mp       []*P_MENN
		mpn      int
		expected []XYZ // Expected grp values
	}{
		{
			name: "Horizontal surface (ez close to zero)",
			mp: []*P_MENN{
				{
					P:    []XYZ{{X: 0, Y: 0, Z: 10}, {X: 1, Y: 0, Z: 10}}, // ez = 0
					wa:   0.0,
					grpx: 0.0,
				},
			},
			mpn:      1,
			expected: []XYZ{{X: 0, Y: 0, Z: 0}},
		},
		{
			name: "Sloped surface (positive ez)",
			// P[0] and P[1] define the vector for ez calculation
			// Let's make a simple case where P[0].Z = 0, P[1].Z = 10, so ez = 10
			// G.Z will be 5 (from mock GDATA)
			// t = -5 / 10 = -0.5
			// ex = 0, ey = 0
			// grpx = 1, wa = 90 (math.Pi/2)
			// grp.X = -1 * sin(90) = -1
			// grp.Y = -1 * cos(90) = 0
			mp: []*P_MENN{
				{
					P:    []XYZ{{X: 0, Y: 0, Z: 0}, {X: 0, Y: 0, Z: 10}}, // ex=0, ey=0, ez=10
					wa:   90.0, // 90 degrees
					grpx: 1.0,
				},
			},
			mpn:      1,
			expected: []XYZ{{X: -1.0, Y: 0.0, Z: 0.0}},
		},
		{
			name: "Sloped surface (negative ez)",
			// P[0].Z = 10, P[1].Z = 0, so ez = -10
			// G.Z will be 5 (from mock GDATA)
			// t = -5 / -10 = 0.5
			// ex = 0, ey = 0
			// grpx = 1, wa = 0
			// grp.X = -1 * sin(0) = 0
			// grp.Y = -1 * cos(0) = -1
			mp: []*P_MENN{
				{
					P:    []XYZ{{X: 0, Y: 0, Z: 10}, {X: 0, Y: 0, Z: 0}}, // ex=0, ey=0, ez=-10
					wa:   0.0, // 0 degrees
					grpx: 1.0,
				},
			},
			mpn:      1,
			expected: []XYZ{{X: 0.0, Y: -1.0, Z: 0.0}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GRGPOINT(tt.mp, tt.mpn)

			for i := 0; i < tt.mpn; i++ {
				// Use a small epsilon for float comparisons
				const epsilon = 1e-9
				if math.Abs(tt.mp[i].grp.X-tt.expected[i].X) > epsilon ||
					math.Abs(tt.mp[i].grp.Y-tt.expected[i].Y) > epsilon ||
					math.Abs(tt.mp[i].grp.Z-tt.expected[i].Z) > epsilon {
					t.Errorf("For mp[%d], expected grp %v, got %v", i, tt.expected[i], tt.mp[i].grp)
				}
			}
		})
	}
}
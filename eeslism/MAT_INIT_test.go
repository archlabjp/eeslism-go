
package eeslism

import (
	"testing"
)

func TestMATINIT(t *testing.T) {
	// Setup
	q := []*P_MENN{
		{
			e:     XYZ{1, 1, 1},
			polyd: 2,
			P:     []XYZ{{1, 1, 1}, {2, 2, 2}},
		},
		{
			e:     XYZ{3, 3, 3},
			polyd: 1,
			P:     []XYZ{{3, 3, 3}},
		},
	}

	// Execute
	MATINIT(q, len(q))

	// Verify
	for i, p := range q {
		if p.e.X != 0.0 || p.e.Y != 0.0 || p.e.Z != 0.0 {
			t.Errorf("Test Case %d Failed: p.e should be zeroed, but got %v", i+1, p.e)
		}
		for j, pt := range p.P {
			if pt.X != 0.0 || pt.Y != 0.0 || pt.Z != 0.0 {
				t.Errorf("Test Case %d, Point %d Failed: pt should be zeroed, but got %v", i+1, j+1, pt)
			}
		}
	}
}

func TestMATINIT_sum(t *testing.T) {
	// Setup
	op := []*P_MENN{
		{
			sum: 10.0,
			wd:  1,
			opw: []WD_MENN{
				{
					sumw: 5.0,
				},
			},
		},
	}

	// Execute
	MATINIT_sum(len(op), op)

	// Verify
	if op[0].sum != 0.0 {
		t.Errorf("Test Case Failed: op.sum should be zeroed, but got %f", op[0].sum)
	}
	if op[0].opw[0].sumw != 0.0 {
		t.Errorf("Test Case Failed: op.opw[0].sumw should be zeroed, but got %f", op[0].opw[0].sumw)
	}
}

func TestMATINIT_sdstr(t *testing.T) {
	// Setup
	Sdstr := []*SHADSTR{
		{sdsum: []float64{1, 2, 3}},
		{sdsum: []float64{4, 5, 6}},
	}

	// Execute
	MATINIT_sdstr(len(Sdstr), 3, Sdstr)

	// Verify
	for i, s := range Sdstr {
		for j, v := range s.sdsum {
			if v != 0.0 {
				t.Errorf("Test Case %d, Index %d Failed: value should be zeroed, but got %f", i+1, j+1, v)
			}
		}
	}
}

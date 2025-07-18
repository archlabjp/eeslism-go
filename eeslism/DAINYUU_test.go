
package eeslism

import (
	"testing"
)

func TestDAINYUU_MP(t *testing.T) {
	// Setup
	op := []*P_MENN{
		{
			opname: "op1",
			polyd:  4,
			P:      []XYZ{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}, {10, 11, 12}},
			wd:     1,
			opw:    []WD_MENN{
				{
					opwname: "opw1",
					ref:     0.5,
					rgb:     [3]float64{0.1, 0.2, 0.3},
					P:       []XYZ{{1, 1, 1}, {2, 2, 2}},
					grpx:    10.0,
				},
			},
			refg: 0.2,
			wa:   30,
			wb:   60,
			e:    XYZ{0.5, 0.5, 0.707},
		},
	}

	// Execute
	mp := DAINYUU_MP(op)

	// Verify
	if len(mp) != 2 {
		t.Fatalf("Expected 2 P_MENN objects, but got %d", len(mp))
	}

	// Verify the main plane
	if mp[0].opname != "op1" || mp[0].wlflg != 0 {
		t.Errorf("Main plane data is incorrect")
	}

	// Verify the window plane
	if mp[1].opname != "opw1" || mp[1].wlflg != 1 {
		t.Errorf("Window plane data is incorrect")
	}
	if mp[1].ref != 0.5 {
		t.Errorf("Window ref is incorrect")
	}
}

func TestDAINYUU_GP(t *testing.T) {
	// Setup
	p := &XYZ{}
	O := XYZ{1, 1, 1}
	E := XYZ{0, 0, 1} // Plane parallel to XY plane
	ls, ms, ns := 0.0, 0.0, 1.0 // Ray pointing up Z axis

	// Execute
	DAINYUU_GP(p, O, E, ls, ms, ns)

	// Verify
	expected := XYZ{0, 0, 1}
	if *p != expected {
		t.Errorf("Expected %v, but got %v", expected, *p)
	}
}

func TestDAINYUU_SMO2(t *testing.T) {
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
	mp := []*P_MENN{{}, {}}
	Sdstr := []*SHADSTR{
		{sdsum: make([]float64, 24)},
		{sdsum: make([]float64, 24)},
	}

	// Execute (dcnt == 1)
	DAINYUU_SMO2(1, 2, op, mp, Sdstr, 1, 10)

	// Verify (dcnt == 1)
	if mp[0].sum != 10.0 || Sdstr[0].sdsum[10] != 10.0 {
		t.Errorf("Incorrect sum for main plane (dcnt=1)")
	}
	if mp[1].sum != 5.0 || Sdstr[1].sdsum[10] != 5.0 {
		t.Errorf("Incorrect sum for window plane (dcnt=1)")
	}

	// Execute (dcnt != 1)
	Sdstr[0].sdsum[12] = 20.0
	Sdstr[1].sdsum[12] = 15.0
	DAINYUU_SMO2(1, 2, op, mp, Sdstr, 0, 12)

	// Verify (dcnt != 1)
	if mp[0].sum != 20.0 {
		t.Errorf("Incorrect sum for main plane (dcnt=0)")
	}
	if mp[1].sum != 15.0 {
		t.Errorf("Incorrect sum for window plane (dcnt=0)")
	}
}


package eeslism

import (
	"testing"
)

func TestSHADOW(t *testing.T) {
	// Setup a simple case for SHADOW
	g := 0
	DE := 100.0
	opn := 1
	lpn := 0
	ls, ms, ns := 0.0, 0.0, -1.0
	s := &bekt{}
	t_bekt := &bekt{}
	op := &P_MENN{
		polyd: 4,
		P:     []XYZ{{-1, -1, 0}, {1, -1, 0}, {1, 1, 0}, {-1, 1, 0}},
		e:     XYZ{0, 0, 1},
		wd:    0,
	}
	OP := []*P_MENN{op}
	LP := []*P_MENN{}
	var wap float64
	wip := make([]float64, 0)
	nday := 0

	s.ps = make([][]float64, opn)
	for i := 0; i < opn; i++ {
		s.ps[i] = make([]float64, 4)
	}
	t_bekt.ps = make([][]float64, lpn)
	for i := 0; i < lpn; i++ {
		t_bekt.ps[i] = make([]float64, 4)
	}

	// Execute
	SHADOW(g, DE, opn, lpn, ls, ms, ns, s, t_bekt, op, OP, LP, &wap, wip, nday)

	// Verify (simple check)
	if op.sum > 1.0 || op.sum < 0.0 {
		t.Errorf("SHADOW test failed: op.sum should be between 0 and 1, but got %f", op.sum)
	}
}

func TestSHADOWlp(t *testing.T) {
	// Setup a simple case for SHADOWlp
	g := 0
	DE := 100.0
	lpn := 1
	mpn := 0
	ls, ms, ns := 0.0, 0.0, -1.0
	s := &bekt{}
	t_bekt := &bekt{}
	lp := &P_MENN{
		polyd: 4,
		P:     []XYZ{{-1, -1, 1}, {1, -1, 1}, {1, 1, 1}, {-1, 1, 1}},
		e:     XYZ{0, 0, 1},
	}
	LP := []P_MENN{
		{
			polyd: 4,
			P:     []XYZ{{-1, -1, 0}, {1, -1, 0}, {1, 1, 0}, {-1, 1, 0}},
			e:     XYZ{0, 0, 1},
		},
	}
	MP := []P_MENN{}

	s.ps = make([][]float64, lpn)
	for i := 0; i < lpn; i++ {
		s.ps[i] = make([]float64, 4)
		for j := 0; j < 4; j++ {
			s.ps[i][j] = 1.0 // Assume all points are on the front side
		}
	}

	// Execute
	SHADOWlp(g, DE, lpn, mpn, ls, ms, ns, s, t_bekt, lp, LP, MP)

	// Verify (simple check)
	if lp.sum > 1.0 || lp.sum < 0.0 {
		t.Errorf("SHADOWlp test failed: lp.sum should be between 0 and 1, but got %f", lp.sum)
	}
}

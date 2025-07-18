
package eeslism

import (
	"testing"
)

func TestHelmClear(t *testing.T) {
	b := &BHELM{1, 2, 3, 4, 5, 6}
	helmclear(b)
	if *b != (BHELM{}) {
		t.Errorf("helmclear failed: expected all zeros, got %v", *b)
	}
}

func TestHelmSum(t *testing.T) {
	a := &BHELM{1, 2, 3, 4, 5, 6}
	b := &BHELM{7, 8, 9, 10, 11, 12}
	helmsum(a, b)
	expected := &BHELM{8, 10, 12, 14, 16, 18}
	if *b != *expected {
		t.Errorf("helmsum failed: expected %v, got %v", *expected, *b)
	}
}

func TestHelmCpy(t *testing.T) {
	a := &BHELM{1, 2, 3, 4, 5, 6}
	b := &BHELM{}
	helmcpy(a, b)
	if *a != *b {
		t.Errorf("helmcpy failed: expected %v, got %v", *a, *b)
	}
}

func TestHelmDiv(t *testing.T) {
	a := &BHELM{10, 20, 30, 40, 50, 60}
	helmdiv(a, 10)
	expected := &BHELM{1, 2, 3, 4, 5, 6}
	if *a != *expected {
		t.Errorf("helmdiv failed: expected %v, got %v", *expected, *a)
	}
}

func TestHelmsumpd(t *testing.T) {
	N := 2
	u := []float64{2, 3}
	a := []*BHELM{{1, 2, 3, 4, 5, 6}, {7, 8, 9, 10, 11, 12}}
	b := &BHELM{}
	helmsumpd(N, u, a, b)
	expected := &BHELM{23, 28, 33, 38, 43, 48}
	if *b != *expected {
		t.Errorf("helmsumpd failed: expected %v, got %v", *expected, *b)
	}
}

func TestHelmsumpf(t *testing.T) {
	N := 1
	u := 2.0
	a := &BHELM{1, 2, 3, 4, 5, 6}
	b := &BHELM{}
	helmsumpf(N, u, a, b)
	expected := &BHELM{2, 4, 6, 8, 10, 12}
	if *b != *expected {
		t.Errorf("helmsumpf failed: expected %v, got %v", *expected, *b)
	}
}

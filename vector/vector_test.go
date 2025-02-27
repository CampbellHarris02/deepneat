package vector

import (
	"testing"
)

func TestVectorOperations(t *testing.T) {
	v1 := NewVector([]float64{1, 2, 3})
	v2 := NewVector([]float64{4, 5, 6})

	// Test Addition
	addResult, _ := v1.Add(v2)
	expectedAdd := []float64{5, 7, 9}
	for i, val := range addResult.Values {
		if val != expectedAdd[i] {
			t.Errorf("Expected %v, got %v", expectedAdd, addResult.Values)
		}
	}

	// Test Dot Product
	dotProduct, _ := v1.Dot(v2)
	expectedDot := 32.0
	if dotProduct != expectedDot {
		t.Errorf("Expected %f, got %f", expectedDot, dotProduct)
	}

	// Test Norm
	norm := v1.Norm()
	expectedNorm := 3.7416573867739413
	if norm != expectedNorm {
		t.Errorf("Expected %f, got %f", expectedNorm, norm)
	}
}

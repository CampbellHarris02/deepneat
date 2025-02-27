package vector

import (
	"errors"
	"math"
)

// Vector represents an n-dimensional vector.
type Vector struct {
	Values []float64
}

// NewVector creates a new vector from a slice.
func NewVector(values []float64) Vector {
	return Vector{Values: values}
}

// Add adds two vectors element-wise.
func (v Vector) Add(other Vector) (Vector, error) {
	if len(v.Values) != len(other.Values) {
		return Vector{}, errors.New("vectors must have the same length")
	}

	result := make([]float64, len(v.Values))
	for i := range v.Values {
		result[i] = v.Values[i] + other.Values[i]
	}
	return NewVector(result), nil
}

// Subtract subtracts two vectors element-wise.
func (v Vector) Subtract(other Vector) (Vector, error) {
	if len(v.Values) != len(other.Values) {
		return Vector{}, errors.New("vectors must have the same length")
	}

	result := make([]float64, len(v.Values))
	for i := range v.Values {
		result[i] = v.Values[i] - other.Values[i]
	}
	return NewVector(result), nil
}

// Dot computes the dot product of two vectors.
func (v Vector) Dot(other Vector) (float64, error) {
	if len(v.Values) != len(other.Values) {
		return 0, errors.New("vectors must have the same length")
	}

	var sum float64
	for i := range v.Values {
		sum += v.Values[i] * other.Values[i]
	}
	return sum, nil
}

// Scale multiplies the vector by a scalar.
func (v Vector) Scale(scalar float64) Vector {
	result := make([]float64, len(v.Values))
	for i := range v.Values {
		result[i] = v.Values[i] * scalar
	}
	return NewVector(result)
}

// Norm computes the magnitude (Euclidean norm) of the vector.
func (v Vector) Norm() float64 {
	sum := 0.0
	for _, val := range v.Values {
		sum += val * val
	}
	return math.Sqrt(sum)
}

// Normalize returns the unit vector (vector with magnitude 1).
func (v Vector) Normalize() (Vector, error) {
	norm := v.Norm()
	if norm == 0 {
		return Vector{}, errors.New("cannot normalize a zero vector")
	}
	return v.Scale(1 / norm), nil
}

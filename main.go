package main

import (
	"deepneat"
	"deepneat/vector"
	"fmt"
	"log"
)

func main() {
	v1 := vector.NewVector([]float64{1, 2, 3})
	v2 := vector.NewVector([]float64{4, 5, 6})

	// Add vectors
	sum, err := v1.Add(v2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Sum:", sum.Values)

	// Dot product
	dot, err := v1.Dot(v2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Dot Product:", dot)

	// Normalize
	unitVec, err := v1.Normalize()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Normalized Vector:", unitVec.Values)

	genome := deepneat.Genome{GenomeID: 1, NumInputs: 2, NumOutputs: 3}
	fmt.Println(genome)
}

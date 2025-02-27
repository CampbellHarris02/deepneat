package deepneat

import (
	"math/rand"
)

type Activation int

const (
	Sigmoid Activation = iota
	ReLU
	Tanh
)

type NeuronGene struct {
	NeuronID   int
	Bias       float64
	Activation Activation
}

type LinkId struct {
	InputID  int
	OutputID int
}

type LinkGene struct {
	LinkID    LinkId
	Weight    float64
	IsEnabled bool
}

type Genome struct {
	GenomeID   int
	NumInputs  int
	NumOutputs int
}

func CrossoverNueron(a, b NeuronGene) NeuronGene {
	if a.NeuronID != b.NeuronID {
		panic("NeuronGene IDs do not match!")
	}
	bias := choose(0.5, a.Bias, b.Bias)
	activation := chooseActivation(0.5, a.Activation, b.Activation)
	return NeuronGene{
		NeuronID:   a.NeuronID,
		Bias:       bias,
		Activation: activation,
	}
}

func choose(prob float64, val1, val2 float64) float64 {
	if rand.Float64() < prob {
		return val1
	}
	return val2
}

func chooseActivation(prob float64, val1, val2 Activation) Activation {
	if rand.Float64() < prob {
		return val1
	}
	return val2
}

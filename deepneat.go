package deepneat

type Activation int

const (
	Sigmoid Activation = iota
	ReLU
	Tanh
)

type NeuronGene struct {
	neuron_id  int
	bias       float64
	activation Activation
}

type LinkId struct {
	input_id  int
	output_id int
}

type LinkGene struct {
	link_id    LinkId
	weight     float64
	is_enabled bool
}

type Genome struct {
	genome_id   int
	num_inputs  int
	num_outputs int
}

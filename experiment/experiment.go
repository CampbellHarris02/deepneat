package experiment

import (
	"deepneat/neat/genetics"
	"encoding/gob"
	"fmt"
	"github.com/sbinet/npyio/npz"
	"gonum.org/v1/gonum/mat"
	"io"
	"math"
	"sort"
	"time"
)

// An Experiment is a collection of trials for one experiment.
// It's helpful for analysis of a series of experiments.
type Experiment struct {
	Id       int
	Name     string
	RandSeed int64
	Trials
	// The maximal allowed fitness score as defined by fitness function of experiment.
	// It is used to normalize fitness score value used in efficiency score calculation. If this value
	// is not set the fitness score will not be normalized during efficiency score estimation.
	MaxFitnessScore float64
}

// AvgTrialDuration Calculates average duration of experiment's trial. Returns EmptyDuration for experiment with no trials.
// Note, that most trials finish after solution solved, so this metric can be used to represent how efficient the solvers
// was generated
func (e *Experiment) AvgTrialDuration() time.Duration {
	total := time.Duration(0)
	for _, t := range e.Trials {
		total += t.Duration
	}
	if len(e.Trials) > 0 {
		return total / time.Duration(len(e.Trials))
	} else {
		return EmptyDuration
	}
}

// AvgEpochDuration Calculates average duration of evaluations among all generations of organism populations in this experiment
func (e *Experiment) AvgEpochDuration() time.Duration {
	total := time.Duration(0)
	for _, t := range e.Trials {
		total += t.AvgEpochDuration()
	}
	if len(e.Trials) > 0 {
		return total / time.Duration(len(e.Trials))
	} else {
		return EmptyDuration
	}
}

// AvgGenerationsPerTrial Calculates average number of generations evaluated per trial during this experiment.
// This can be helpful when estimating algorithm efficiency, because when winner organism is found the trial is
// terminated, i.e. less evaluations - more fast convergence.
func (e *Experiment) AvgGenerationsPerTrial() float64 {
	total := 0.0
	for _, t := range e.Trials {
		total += float64(len(t.Generations))
	}
	if len(e.Trials) > 0 {
		return total / float64(len(e.Trials))
	} else {
		return 0
	}
}

// MostRecentTrialEvalTime Returns the time of evaluation of the most recent trial
func (e *Experiment) MostRecentTrialEvalTime() time.Time {
	if len(e.Trials) == 0 {
		return time.Time{}
	}
	var u time.Time
	for _, t := range e.Trials {
		ut := t.RecentEpochEvalTime()
		if u.Before(ut) {
			u = ut
		}
	}
	return u
}

// BestOrganism Finds the most fit organism among all trials in this experiment. It's also possible to get the best organism
// only among the ones which was able to solve the experiment's problem. Returns the best fit organism in this experiment
// among with ID of trial where it was found and boolean value to indicate if search was successful.
func (e *Experiment) BestOrganism(onlySolvers bool) (*genetics.Organism, int, bool) {
	var orgs = make(genetics.Organisms, 0, len(e.Trials))
	for i, t := range e.Trials {
		org, found := t.BestOrganism(onlySolvers)
		if found {
			orgs = append(orgs, org)
			org.Flag = i
		}

	}
	if len(orgs) > 0 {
		sort.Sort(sort.Reverse(orgs))
		return orgs[0], orgs[0].Flag, true
	} else {
		return nil, -1, false
	}

}

// Solved is to check if solution was found in at least one trial
func (e *Experiment) Solved() bool {
	for _, t := range e.Trials {
		if t.Solved() {
			return true
		}
	}
	return false
}

// BestFitness finds the fitness values of the best organisms for each trial.
func (e *Experiment) BestFitness() Floats {
	var x Floats = make([]float64, len(e.Trials))
	for i, t := range e.Trials {
		if org, ok := t.BestOrganism(false); ok {
			x[i] = org.Fitness
		}
	}
	return x
}

// BestSpeciesAge finds the age values of the species with the best organisms for each trial.
func (e *Experiment) BestSpeciesAge() Floats {
	var x Floats = make([]float64, len(e.Trials))
	for i, t := range e.Trials {
		if org, ok := t.BestOrganism(false); ok {
			if org.Species != nil {
				x[i] = float64(org.Species.Age)
			}
		}
	}
	return x
}

// BestComplexity finds the complexity values of the best organisms for each trial.
func (e *Experiment) BestComplexity() Floats {
	var x Floats = make([]float64, len(e.Trials))
	for i, t := range e.Trials {
		if org, ok := t.BestOrganism(false); ok {
			x[i] = float64(organismComplexity(org))
		}
	}
	return x
}

// AvgDiversity returns the average number of species in each trial
func (e *Experiment) AvgDiversity() Floats {
	var x Floats = make([]float64, len(e.Trials))
	for i, t := range e.Trials {
		x[i] = t.Diversity().Mean()
	}
	return x
}

// EpochsPerTrial calculates the number of epochs per trial
func (e *Experiment) EpochsPerTrial() Floats {
	var x Floats = make([]float64, len(e.Trials))
	for i, t := range e.Trials {
		x[i] = float64(len(t.Generations))
	}
	return x
}

// TrialsSolved calculates the number of trials solved
func (e *Experiment) TrialsSolved() int {
	count := 0
	for _, t := range e.Trials {
		if t.Solved() {
			count++
		}
	}
	return count
}

// SuccessRate calculates the success rate of the experiment as ratio of trials with successful solvers
// per total number of trials executed.
func (e *Experiment) SuccessRate() float64 {
	solved := float64(e.TrialsSolved())
	if len(e.Trials) > 0 {
		return solved / float64(len(e.Trials))
	} else {
		return 0
	}
}

// AvgWinnerStatistics calculates the average number of nodes, genes, organisms evaluations,
// and species diversity of winners among all trials, i.e. for all trials where winning solution was found.
func (e *Experiment) AvgWinnerStatistics() (avgNodes, avgGenes, avgEvals, avgDiversity float64) {
	totalNodes, totalGenes, totalEvals, totalDiversity := 0, 0, 0, 0
	count := 0
	for _, t := range e.Trials {
		if t.Solved() {
			nodes, genes, evals, diversity := t.WinnerStatistics()
			totalNodes += nodes
			totalGenes += genes
			totalEvals += evals
			totalDiversity += diversity

			count++
		}
	}
	if count == 0 {
		return -1, -1, -1, -1
	}
	avgNodes = float64(totalNodes) / float64(count)
	avgGenes = float64(totalGenes) / float64(count)
	avgEvals = float64(totalEvals) / float64(count)
	avgDiversity = float64(totalDiversity) / float64(count)
	return avgNodes, avgGenes, avgEvals, avgDiversity
}

// EfficiencyScore calculates the efficiency score of the found solution.
//
// We are interested in efficient solver solution that take less time per epoch, fewer generations per trial,
// and produce less complicated winner genomes. At the same time it should have maximal fitness score and maximal
// success rate among trials.
//
//	This value can only be compared against values obtained for the same type of experiments.
func (e *Experiment) EfficiencyScore() float64 {
	meanComplexity, meanFitness := 0.0, 0.0
	if len(e.Trials) > 1 {
		count := 0.0
		for _, t := range e.Trials {
			if t.Solved() {
				if t.WinnerGeneration == nil {
					// find winner
					t.WinnerStatistics()
				}

				meanComplexity += float64(t.WinnerGeneration.ChampionComplexity())
				meanFitness += t.WinnerGeneration.Champion.Fitness

				count++
			}
		}
		meanComplexity /= count
		meanFitness /= count
	}

	// normalize and scale fitness score if appropriate
	fitnessScore := meanFitness
	if e.MaxFitnessScore > 0 {
		fitnessScore = (fitnessScore / e.MaxFitnessScore) * 100
	}

	score := e.penaltyScore(meanComplexity)
	if score > 0 {
		// calculate normalized score
		successRate := e.SuccessRate()
		logPenaltyScore := math.Log(score)
		score = successRate * fitnessScore / logPenaltyScore
	}

	return score
}

func (e *Experiment) penaltyScore(meanComplexity float64) float64 {
	return e.AvgEpochDuration().Seconds() * 1000.0 * e.AvgGenerationsPerTrial() * meanComplexity
}

// PrintStatistics Prints experiment statistics
func (e *Experiment) PrintStatistics() {
	fmt.Printf("\nSolved %d trials from %d, success rate: %f\n", e.TrialsSolved(), len(e.Trials), e.SuccessRate())
	fmt.Printf("Random seed: %d\n", e.RandSeed)
	fmt.Printf("Average\n\tTrial duration:\t\t%s\n\tEpoch duration:\t\t%s\n\tGenerations/trial:\t%.1f\n",
		e.AvgTrialDuration(), e.AvgEpochDuration(), e.AvgGenerationsPerTrial())
	// Print absolute champion statistics
	if org, trid, found := e.BestOrganism(true); found {
		nodes, genes, evals, divers := e.Trials[trid].WinnerStatistics()
		fmt.Printf("\nChampion found in %d trial run\n\tNodes:\t\t\t%d\n\tGenes:\t\t\t%d\n\tEvaluations:\t\t%d\n\n\tDiversity:\t\t%d",
			trid, nodes, genes, evals, divers)
		fmt.Printf("\n\tComplexity:\t\t%d\n\tAge:\t\t\t%d\n\tFitness:\t\t%f\n",
			organismComplexity(org), org.Species.Age, org.Fitness)
	} else {
		fmt.Println("\nNo winner found in the experiment!!!")
	}

	// Print average winner statistics
	meanComplexity, meanDiversity, meanAge, meanFitness := 0.0, 0.0, 0.0, 0.0
	if len(e.Trials) > 1 {
		avgNodes, avgGenes, avgEvals, avgDivers, avgGenerations := 0.0, 0.0, 0.0, 0.0, 0.0
		count := 0.0
		for _, t := range e.Trials {
			if t.Solved() {
				nodes, genes, evals, diversity := t.WinnerStatistics()
				avgNodes += float64(nodes)
				avgGenes += float64(genes)
				avgEvals += float64(evals)
				avgDivers += float64(diversity)
				avgGenerations += float64(len(t.Generations))

				meanComplexity += float64(t.WinnerGeneration.ChampionComplexity())
				meanAge += float64(t.WinnerGeneration.Champion.Species.Age)
				meanFitness += t.WinnerGeneration.Champion.Fitness

				count++
			}
		}
		avgNodes /= count
		avgGenes /= count
		avgEvals /= count
		avgDivers /= count
		avgGenerations /= count
		fmt.Printf("\nAverages among winners (successfull solvers)\n\tAvg Nodes:\t\t%.1f\n\tAvg Genes:\t\t%.1f\n\tAvg Evaluations:\t%.1f\n\tGenerations/trial:\t%.1f\n\n\tDiversity:\t\t%f\n",
			avgNodes, avgGenes, avgEvals, avgGenerations, avgDivers)

		meanComplexity /= count
		meanAge /= count
		meanFitness /= count
		fmt.Printf("\tComplexity:\t\t%f\n\tAge:\t\t\t%f\n\tFitness:\t\t%f\n",
			meanComplexity, meanAge, meanFitness)
	}

	// Print the average values for each population of organisms evaluated
	count := float64(len(e.Trials))
	for _, t := range e.Trials {
		fitness, age, complexity := t.Average()

		meanComplexity += complexity.Mean()
		meanDiversity += t.Diversity().Mean()
		meanAge += age.Mean()
		meanFitness += fitness.Mean()
	}
	meanComplexity /= count
	meanDiversity /= count
	meanAge /= count
	meanFitness /= count
	fmt.Printf("\nAverages of the best organisms evaluated during experiment\n\tDiversity:\t\t%f\n\tComplexity:\t\t%f\n\tAge:\t\t\t%f\n\tFitness:\t\t%f\n",
		meanDiversity, meanComplexity, meanAge, meanFitness)

	score := e.EfficiencyScore()
	fmt.Printf("\nEfficiency score:\t\t%f\n\n", score)
}

// Write is to write encoded experiment data into provided writer
func (e *Experiment) Write(w io.Writer) error {
	enc := gob.NewEncoder(w)
	return e.Encode(enc)
}

// Encode Encodes experiment with GOB encoding
func (e *Experiment) Encode(enc *gob.Encoder) error {
	if err := enc.Encode(e.Id); err != nil {
		return err
	}
	if err := enc.Encode(e.Name); err != nil {
		return err
	}

	// encode trials
	if err := enc.Encode(len(e.Trials)); err != nil {
		return err
	}
	for _, t := range e.Trials {
		if err := t.Encode(enc); err != nil {
			return err
		}
	}
	return nil
}

// Read is to read experiment data from provided reader and decodes it
func (e *Experiment) Read(r io.Reader) error {
	dec := gob.NewDecoder(r)
	return e.Decode(dec)
}

// Decode Decodes experiment data
func (e *Experiment) Decode(dec *gob.Decoder) error {
	if err := dec.Decode(&e.Id); err != nil {
		return err
	}
	if err := dec.Decode(&e.Name); err != nil {
		return err
	}

	// decode Trials
	var tNum int
	if err := dec.Decode(&tNum); err != nil {
		return err
	}

	e.Trials = make([]Trial, tNum)
	for i := 0; i < tNum; i++ {
		trial := Trial{}
		if err := trial.Decode(dec); err != nil {
			return err
		}
		e.Trials[i] = trial
	}
	return nil
}

// WriteNPZ Dumps experiment results to the NPZ file.
// The file has the following structure:
// - trials_fitness - the mean, variance of fitness scores per trial
// - trials_ages - the mean, variance of species ages per trial
// - trials_complexity - the mean, variance of genome complexity of the best organisms among species per trial
// - trial_[0...n]_epoch_mean_fitnesses - the mean fitness scores per epoch per trial
// - trial_[0...n]_epoch_best_fitnesses - the best fitness scores per epoch per trial
// the same for AGE and COMPLEXITY per epoch per trial
// - trial_[0...n]_epoch_diversity - the number of species per epoch per trial
func (e *Experiment) WriteNPZ(w io.Writer) error {
	// write general statistics
	trialsFitness, trialsAges, trialsComplexity := e.fitnessAgeComplexityMat()
	out := npz.NewWriter(w)
	if err := out.Write("trials_number", Floats{float64(len(e.Trials))}); err != nil {
		return err
	}
	if err := out.Write("trials_fitness", trialsFitness); err != nil {
		return err
	}
	if err := out.Write("trials_ages", trialsAges); err != nil {
		return err
	}
	if err := out.Write("trials_complexity", trialsComplexity); err != nil {
		return err
	}
	// write statistics per epoch per trial
	//
	for i, t := range e.Trials {
		fitness, age, complexity := t.Average()
		if err := out.Write(fmt.Sprintf("trial_%d_epoch_mean_fitnesses", i), fitness); err != nil {
			return err
		}
		if err := out.Write(fmt.Sprintf("trial_%d_epoch_mean_ages", i), age); err != nil {
			return err
		}
		if err := out.Write(fmt.Sprintf("trial_%d_epoch_mean_complexities", i), complexity); err != nil {
			return err
		}
		if err := out.Write(fmt.Sprintf("trial_%d_epoch_best_fitnesses", i), t.ChampionsFitness()); err != nil {
			return err
		}
		if err := out.Write(fmt.Sprintf("trial_%d_epoch_best_ages", i), t.ChampionSpeciesAges()); err != nil {
			return err
		}
		if err := out.Write(fmt.Sprintf("trial_%d_epoch_best_complexities", i), t.ChampionsComplexities()); err != nil {
			return err
		}
		if err := out.Write(fmt.Sprintf("trial_%d_epoch_diversity", i), t.Diversity()); err != nil {
			return err
		}
	}
	return out.Close()
}

func (e *Experiment) fitnessAgeComplexityMat() (trialsFitness, trialsAges, trialsComplexity *mat.Dense) {
	trialsFitness = mat.NewDense(len(e.Trials), 2, nil)    // mean, var
	trialsAges = mat.NewDense(len(e.Trials), 2, nil)       // mean, var
	trialsComplexity = mat.NewDense(len(e.Trials), 2, nil) // mean, var
	for i, t := range e.Trials {
		fitness, age, complexity := t.Average()
		trialsFitness.SetRow(i, fitness.MeanVariance())
		trialsAges.SetRow(i, age.MeanVariance())
		trialsComplexity.SetRow(i, complexity.MeanVariance())
	}
	return trialsFitness, trialsAges, trialsComplexity
}

// Experiments is a sortable list of experiments by execution time and Id
type Experiments []Experiment

func (es Experiments) Len() int {
	return len(es)
}
func (es Experiments) Swap(i, j int) {
	es[i], es[j] = es[j], es[i]
}
func (es Experiments) Less(i, j int) bool {
	ui := es[i].MostRecentTrialEvalTime()
	uj := es[j].MostRecentTrialEvalTime()
	if ui.Equal(uj) {
		return es[i].Id < es[j].Id
	}
	return ui.Before(uj)
}

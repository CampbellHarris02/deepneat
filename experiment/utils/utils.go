// Package utils provides common utilities to be used by experiments.
package utils

import (
	"bytes"
	"deepneat/experiment"
	"deepneat/neat"
	"deepneat/neat/genetics"
	"deepneat/neat/network"
	"deepneat/neat/network/formats"
	"fmt"
	"log"
	"os"
)

// WriteGenomePlain is to write genome of the organism to the genomeFile in the outDir directory using plain encoding.
// The method return path to the file if successful or error if failed.
func WriteGenomePlain(genomeFile, outDir string, org *genetics.Organism, epoch *experiment.Generation) (string, error) {
	phenotype, err := org.Phenotype()
	if err != nil {
		return "", err
	}
	orgPath := fmt.Sprintf("%s/%s_%d-%d", CreateOutDirForTrial(outDir, epoch.TrialId),
		genomeFile, phenotype.NodeCount(), phenotype.LinkCount())
	if file, err := os.Create(orgPath); err != nil {
		return "", err
	} else if err = org.Genotype.Write(file); err != nil {
		return "", err
	}
	return orgPath, nil
}

// WriteGenomeDOT is to write genome of the organism to the genomeFile in the outDir directory using DOT encoding.
// The method return path to the file if successful or error if failed.
func WriteGenomeDOT(genomeFile, outDir string, org *genetics.Organism, epoch *experiment.Generation) (string, error) {
	phenotype, err := org.Phenotype()
	if err != nil {
		return "", err
	}
	orgPath := fmt.Sprintf("%s/%s_%d-%d.dot", CreateOutDirForTrial(outDir, epoch.TrialId),
		genomeFile, phenotype.NodeCount(), phenotype.LinkCount())
	if file, err := os.Create(orgPath); err != nil {
		return "", err
	} else if err = formats.WriteDOT(file, phenotype); err != nil {
		return "", err
	}
	return orgPath, nil
}

// WriteGenomeCytoscapeJSON is to write genome of the organism to the genomeFile in the outDir directory using Cytoscape JSON encoding.
// The method return path to the file if successful or error if failed.
func WriteGenomeCytoscapeJSON(genomeFile, outDir string, org *genetics.Organism, epoch *experiment.Generation) (string, error) {
	phenotype, err := org.Phenotype()
	if err != nil {
		return "", err
	}
	orgPath := fmt.Sprintf("%s/%s_%d-%d.cyjs", CreateOutDirForTrial(outDir, epoch.TrialId),
		genomeFile, phenotype.NodeCount(), phenotype.LinkCount())
	if file, err := os.Create(orgPath); err != nil {
		return "", err
	} else if err = formats.WriteCytoscapeJSON(file, phenotype); err != nil {
		return "", err
	}
	return orgPath, nil
}

// WritePopulationPlain is to write genomes of the entire population using plain encoding in the outDir directory.
// The methods return path to the file if successful or error if failed.
func WritePopulationPlain(outDir string, pop *genetics.Population, epoch *experiment.Generation) (string, error) {
	popPath := fmt.Sprintf("%s/gen_%d", CreateOutDirForTrial(outDir, epoch.TrialId), epoch.Id)
	if file, err := os.Create(popPath); err != nil {
		return "", err
	} else if err = pop.WriteBySpecies(file); err != nil {
		return "", err
	}
	return popPath, nil
}

// CreateOutDirForTrial allows creating the output directory for specific trial of the experiment using standard name.
func CreateOutDirForTrial(outDir string, trialID int) string {
	dir := fmt.Sprintf("%s/%d", outDir, trialID)
	if _, err := os.Stat(dir); err != nil {
		// create output dir
		if err = os.MkdirAll(dir, os.ModePerm); err != nil {
			log.Fatal("Failed to create output directory: ", err)
		}
	}
	return dir
}

// PrintActivationDepth is to print maximal activation depth of phenotype network of the organism
func PrintActivationDepth(organism *genetics.Organism, printActivationPath bool) {
	phenotype, err := organism.Phenotype()
	if err != nil {
		return
	}
	if depth, err := phenotype.MaxActivationDepthWithCap(0); err == nil {
		neat.InfoLog(fmt.Sprintf("Activation depth of the winner: %d\n", depth))
	}
	if printActivationPath {
		buf := bytes.NewBufferString("Activation paths of the winner:\n")
		if err = network.PrintAllActivationDepthPaths(phenotype, buf); err == nil {
			neat.InfoLog(buf.String())
		}
	}
}

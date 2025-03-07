package main

import (
	"context"
	"deepneat/examples/pole"
	"deepneat/examples/pole2"
	"deepneat/examples/xor"
	"deepneat/experiment"
	"deepneat/neat"
	"deepneat/neat/genetics"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// The experiment runner boilerplate code
func main() {
	var outDirPath = flag.String("out", "./out", "The output directory to store results.")
	var contextPath = flag.String("context", "./data/xor.neat", "The execution context configuration file.")
	var genomePath = flag.String("genome", "./data/xorstartgenes", "The seed genome to start with.")
	var experimentName = flag.String("experiment", "XOR", "The name of experiment to run. [XOR, cart_pole, cart_2pole_markov, cart_2pole_non-markov]")
	var trialsCount = flag.Int("trials", 0, "The number of trials for experiment. Overrides the one set in configuration.")
	var logLevel = flag.String("log_level", "", "The logger level to be used. Overrides the one set in configuration.")
	var randSeed = flag.Int64("seed", 0, "The seed for random number generator")

	flag.Parse()

	// Seed the random-number generator with current time so that
	// the numbers will be different every time we run.
	seed := time.Now().Unix()
	if randSeed != nil {
		seed = *randSeed
	}
	rand.Seed(seed)

	// Load NEAT options
	neatOptions, err := neat.ReadNeatOptionsFromFile(*contextPath)
	if err != nil {
		log.Fatal("Failed to load NEAT options: ", err)
	}

	// Load Genome
	log.Printf("Loading start genome for %s experiment from file '%s'\n", *experimentName, *genomePath)
	reader, err := genetics.NewGenomeReaderFromFile(*genomePath)
	if err != nil {
		log.Fatalf("Failed to open genome file, reason: '%s'", err)
	}
	startGenome, err := reader.Read()
	if err != nil {
		log.Fatalf("Failed to read start genome, reason: '%s'", err)
	}
	fmt.Println(startGenome)

	// Check if output dir exists
	outDir := *outDirPath
	if _, err := os.Stat(outDir); err == nil {
		// backup it
		backUpDir := fmt.Sprintf("%s-%s", outDir, time.Now().Format("2006-01-02T15_04_05"))
		// clear it
		err = os.Rename(outDir, backUpDir)
		if err != nil {
			log.Fatal("Failed to do previous results backup: ", err)
		}
	}
	// create output dir
	err = os.MkdirAll(outDir, os.ModePerm)
	if err != nil {
		log.Fatal("Failed to create output directory: ", err)
	}

	// Override neatOptions configuration parameters with ones set from command line
	if *trialsCount > 0 {
		neatOptions.NumRuns = *trialsCount
	}
	if len(*logLevel) > 0 {
		if err = neat.InitLogger(*logLevel); err != nil {
			log.Fatal("Failed to initialize logger: ", err)
		}
	}

	// create experiment
	exp := experiment.Experiment{
		Id:       0,
		Trials:   make(experiment.Trials, neatOptions.NumRuns),
		RandSeed: seed,
	}
	var generationEvaluator experiment.GenerationEvaluator
	switch *experimentName {
	case "XOR":
		exp.MaxFitnessScore = 16.0 // as given by fitness function definition
		generationEvaluator = xor.NewXORGenerationEvaluator(outDir)
	case "cart_pole":
		exp.MaxFitnessScore = 1.0 // as given by fitness function definition
		generationEvaluator = pole.NewCartPoleGenerationEvaluator(outDir, true, 1500000)
	case "cart_pole_parallel":
		exp.MaxFitnessScore = 1.0 // as given by fitness function definition
		generationEvaluator = pole.NewCartPoleParallelGenerationEvaluator(outDir, true, 1500000)
	case "cart_2pole_markov":
		exp.MaxFitnessScore = 1.0 // as given by fitness function definition
		generationEvaluator = pole2.NewCartDoublePoleGenerationEvaluator(outDir, true, pole2.ContinuousAction)
	case "cart_2pole_non-markov":
		generationEvaluator = pole2.NewCartDoublePoleGenerationEvaluator(outDir, false, pole2.ContinuousAction)
	case "cart_2pole_markov_parallel":
		exp.MaxFitnessScore = 1.0 // as given by fitness function definition
		generationEvaluator = pole2.NewCartDoublePoleParallelGenerationEvaluator(outDir, true, pole2.ContinuousAction)
	default:
		log.Fatalf("Unsupported experiment: %s", *experimentName)
	}

	// prepare to execute
	errChan := make(chan error)
	ctx, cancel := context.WithCancel(context.Background())

	// run experiment in the separate GO routine
	go func() {
		if err = exp.Execute(neat.NewContext(ctx, neatOptions), startGenome, generationEvaluator, nil); err != nil {
			errChan <- err
		} else {
			errChan <- nil
		}
	}()

	// register handler to wait for termination signals
	//
	go func(cancel context.CancelFunc) {
		fmt.Println("\nPress Ctrl+C to stop")

		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		select {
		case <-signals:
			// signal to stop test fixture
			cancel()
		case err = <-errChan:
			// stop waiting
		}
	}(cancel)

	// Wait for experiment completion
	//
	err = <-errChan
	if err != nil {
		// error during execution
		log.Fatalf("Experiment execution failed: %s", err)
	}

	// Print experiment results statistics
	//
	exp.PrintStatistics()

	fmt.Printf(">>> Start genome file:  %s\n", *genomePath)
	fmt.Printf(">>> Configuration file: %s\n", *contextPath)

	// Save experiment data in native format
	//
	expResPath := fmt.Sprintf("%s/%s.dat", outDir, *experimentName)
	if expResFile, err := os.Create(expResPath); err != nil {
		log.Fatal("Failed to create file for experiment results", err)
	} else if err = exp.Write(expResFile); err != nil {
		log.Fatal("Failed to save experiment results", err)
	}

	// Save experiment data in Numpy NPZ format if requested
	//
	npzResPath := fmt.Sprintf("%s/%s.npz", outDir, *experimentName)
	if npzResFile, err := os.Create(npzResPath); err != nil {
		log.Fatalf("Failed to create file for experiment results: [%s], reason: %s", npzResPath, err)
	} else if err = exp.WriteNPZ(npzResFile); err != nil {
		log.Fatal("Failed to save experiment results as NPZ file", err)
	}
}

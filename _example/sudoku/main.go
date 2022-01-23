package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/arl/evolve"
	"github.com/arl/evolve/condition"
	"github.com/arl/evolve/engine"
	"github.com/arl/evolve/generator"
	"github.com/arl/evolve/operator"
	"github.com/arl/evolve/operator/xover"
	"github.com/arl/evolve/pkg/mt19937"
	"github.com/arl/evolve/selection"
)

var (
	//go:embed puzzles/blank.txt
	blankPuzzle []byte
	//go:embed puzzles/easy.txt
	easyPuzzle []byte
	//go:embed puzzles/medium.txt
	mediumPuzzle []byte
	//go:embed puzzles/hard.txt
	hardPuzzle []byte
)

func check(err error, v ...interface{}) {
	if err != nil {
		if len(v) == 0 {
			log.Fatal(v, err)
		}
	}
}

func readPattern(r io.Reader) ([]string, error) {
	puzzle := []string{}

	s := bufio.NewScanner(r)
	for s.Scan() {
		puzzle = append(puzzle, s.Text())
	}
	return puzzle, s.Err()
}

func solveSudoku(pattern []string) error {
	const (
		popsize = 500
		nelites = 25
	)

	// Crossover rows between parents (so offspring is x rows from parent1 and y
	// rows from parent2).
	xover := xover.New[*sudoku](mater{})
	xover.Points = generator.Const(1)
	xover.Probability = generator.Const(1.0)

	rng := rand.New(mt19937.New(time.Now().UnixNano()))

	mutation := &rowMutation{
		Number: generator.NewPoisson[uint](generator.Const(2.0), rng),
		Amount: generator.Uniform[uint](1, 8, rng),
	}

	pipeline := operator.Pipeline[*sudoku]{xover, mutation}

	selector := &selection.Tournament[*sudoku]{Probability: generator.Const(0.85)}

	fac, err := newFactory(pattern)
	check(err)

	epocher := engine.Generational[*sudoku]{
		Operator:  pipeline,
		Evaluator: evaluator{},
		Selection: selector,
		Elites:    nelites,
	}

	eng := &engine.Engine[*sudoku]{
		Factory:   fac,
		Evaluator: evaluator{},
		Epocher:   &epocher,
		RNG:       rng,
		EndConditions: []evolve.Condition[*sudoku]{
			condition.TargetFitness[*sudoku]{Fitness: 0, Natural: false},
			condition.NewUserAbort[*sudoku](),
		},
	}

	eng.AddObserver(engine.ObserverFunc(func(stats *evolve.PopulationStats[*sudoku]) {
		// Only shows multiple of 100 generations
		if stats.Generation%100 == 0 {
			return
		}
		log.Printf("Generation %d: %s (%v)\n", stats.Generation, stats.Best, stats.BestFitness)
	}))

	bests, _, err := eng.Evolve(popsize)
	check(err)

	log.Printf("Sudoku solution:\n%v\n", bests.Candidates[0])
	return nil
}

func ui() (io.Reader, error) {
	puzzles := [][]byte{blankPuzzle, easyPuzzle, mediumPuzzle, hardPuzzle}

	fmt.Printf("\t[0] blank\n")
	fmt.Printf("\t[1] easy\n")
	fmt.Printf("\t[2] medium\n")
	fmt.Printf("\t[3] hard\n")

	fmt.Print("Choose the sudoku puzzle you want to solve? ")
	var i int
	if _, err := fmt.Scanf("%d", &i); err != nil {
		return nil, fmt.Errorf("can't read your choice: %v", err)
	}
	if i < 0 || i > 3 {
		return nil, fmt.Errorf("invalid entry")
	}

	return bytes.NewReader(puzzles[i]), nil
}

func main() {
	fpuzzle := flag.String("puzzle", "", "file with the sudoku puzzle to solve")
	flag.Parse()

	var r io.Reader // puzzle buffer

	if *fpuzzle != "" {
		f, err := os.Open(*fpuzzle)
		if err != nil {
			log.Fatalf("can't open puzzle file: %v", err)
		}
		defer f.Close()

		r = f
	} else {
		pr, err := ui()
		if err != nil {
			log.Fatal(err)
		}

		r = pr
	}

	pattern, err := readPattern(r)
	if err != nil {
		log.Fatalf("can't read sudo pattern: %v", err)
	}

	err = solveSudoku(pattern)
	if err != nil {
		log.Fatalf("couldn't solve sudoku pattern: %v\n", err)
	}
}

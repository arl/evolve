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

	selector := selection.NewTournament[*sudoku]()
	check(selector.SetProb(0.85))

	obs := engine.ObserverFunc(func(stats *evolve.PopulationStats[*sudoku]) {
		// Only shows multiple of 100 generations
		if stats.GenNumber%100 == 0 {
			return
		}
		log.Printf("Gen %d, best solution has a fitness of %v\n%v\n",
			stats.GenNumber, stats.BestFitness, stats.BestCand)
	})

	gen, err := newFactory(pattern)
	check(err)

	epocher := engine.Generational[*sudoku]{Op: pipeline, Eval: evaluator{}, Sel: selector}

	eng, err := engine.New[*sudoku](
		gen,
		evaluator{},
		&epocher,
		engine.Observe(obs),
		engine.Rand[*sudoku](rng),
	)
	check(err)

	const (
		popsize = 500
		nelites = 500 * 0.05
	)
	bests, _, err := eng.Evolve(
		popsize,
		engine.Elites[*sudoku](nelites),
		engine.EndOn[*sudoku](condition.TargetFitness[*sudoku]{Fitness: 0, Natural: false}),
		engine.EndOn[*sudoku](condition.NewUserAbort[*sudoku]()),
	)
	check(err)

	log.Printf("Sudoku solution:\n%v\n", bests[0].Candidate)
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

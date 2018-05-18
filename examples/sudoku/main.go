package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path"
	"time"

	"github.com/aurelien-rainone/evolve/pkg/api"
	"github.com/aurelien-rainone/evolve/pkg/engine"
	"github.com/aurelien-rainone/evolve/pkg/operator"
	"github.com/aurelien-rainone/evolve/pkg/operator/xover"
	"github.com/aurelien-rainone/evolve/pkg/selection"
	"github.com/aurelien-rainone/evolve/pkg/termination"
	"github.com/aurelien-rainone/evolve/random"
)

func readSudokus(dir string) ([]string, error) {
	f, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	names, err := f.Readdirnames(0)
	switch {
	case err != nil:
		return nil, err
	case len(names) == 0:
		return nil, errors.New("empty directory")
	}
	return names, err
}

func readPattern(fn string) ([]string, error) {
	f, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	puzzle := []string{}

	s := bufio.NewScanner(f)
	for s.Scan() {
		puzzle = append(puzzle, s.Text())
	}
	if s.Err() != nil {
		return nil, s.Err()
	}
	return puzzle, nil
}

func solveSudoku(pattern []string) error {
	rng := rand.New(random.NewMT19937(time.Now().UnixNano()))

	// Crossover rows between parents (so offspring is x rows from parent1 and y
	// rows from parent2).
	xover := xover.New(mater{})
	xover.SetPoints(1)

	mutation := newRowMutation()
	// TODO: use a PoissonGenerator for mutation count and a
	// DiscreteUniformGenerator for mutation amount
	mutation.SetMutationsRange(1, 2)
	mutation.SetAmountRange(1, 8)

	pipeline := operator.Pipeline{xover, mutation}

	selector := selection.NewTournament()
	err := selector.SetProb(0.85)
	if err != nil {
		return fmt.Errorf("can't create selection strategy: %v", err)
	}

	factory, err := newSudokuFactory(pattern)
	if err != nil {
		return fmt.Errorf("can't create factory strategy: %v", err)
	}

	eng := engine.NewGenerational(factory, pipeline, evaluator{}, selector, rng)

	eng.AddObserver(api.ObserverFunc(func(data *api.PopulationData) {
		if data.GenNumber%100 == 0 {
			return
		}
		// only shows multiple of 100 generations
		fmt.Printf("gen:%d, fitness:%v:\n%v\n", data.GenNumber, data.BestFitness, data.BestCand.(*sudoku))
	}))

	const (
		popsize = 500
		nelites = 500 * 0.05
	)
	solution := eng.Evolve(
		popsize, nelites,
		termination.TargetFitness{0, false},
		termination.NewUserAbort(),
	)

	fmt.Printf("solution:\n%v\n", solution.(*sudoku))
	return nil
}

func main() {
	puzdir := flag.String("puzzles", "./puzzles", "directory containing the puzzles")
	flag.Parse()

	puzzles, err := readSudokus(*puzdir)
	if err != nil {
		fmt.Printf("can't read puzzles directory: %v", err)
		return
	}

	for i, p := range puzzles {
		fmt.Printf("\t[%d] %s\n", i, p)
	}

	fmt.Print("Choose the sudoku puzzle you want to solve? ")
	var i int
	if _, err := fmt.Scanf("%d", &i); err != nil {
		fmt.Printf("can't read choice: %v", err)
		return
	}
	if i < 0 || i >= len(puzzles) {
		fmt.Println("invalid entry")
		return
	}
	fmt.Println("you chose ", i)
	pattern, err := readPattern(path.Join(*puzdir, puzzles[i]))
	if err != nil {
		fmt.Printf("can't read pattern: %v", err)
		return
	}
	fmt.Println("pattern ", pattern)
	fmt.Println("sudoku:", (&sudoku{}).String())

	err = solveSudoku(pattern)
	if err != nil {
		fmt.Printf("Couldn't solve: %v\n", err)
	}
}

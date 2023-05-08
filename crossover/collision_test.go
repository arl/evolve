package crossover

import (
	"bytes"
	_ "embed"
	"fmt"
	"math/rand"
	"sort"
	"testing"
	"time"

	"github.com/arl/evolve"
	"github.com/arl/evolve/generator"
	"github.com/arl/evolve/pkg/mt19937"
	"github.com/arl/evolve/pkg/set"
	"github.com/arl/evolve/pkg/tsp"
)

/*
0,0
                A
               /\
              /  \
             /    \
            /      \
           /        \
          /          \
         /            \
        /              \
       /                \
      /                  \
     /                    \
    /                      \
   /                        \
  /                          \
 /                            \
| F                           |  B
|                             |
|                             |
|                             |
|                             |
|                             |
|               D             |
|              /\             |
|             /  \            |
|            /    \           |
|           /      \          |
|          /        \         |
|         /          \        |
|        /            \       |
|       /              \      |
|      /                \     |
|     /                  \    |
|    /                    \   |
|   /                      \  |
|  /                        \ |
| /                          \|
|/ E                          | C


city | coord city |index
  A  | ( 50 , 0)  | 0
  B  | (100, 100) | 1
  C  | (100, 250) | 2
  D  | ( 50, 150) | 3
  E  | (  0, 250) | 4
  F  | (  0, 100) | 5

  optimum route -> ABFDCE
*/

func TestCollisionCrossover(t *testing.T) {
	optimum := []int{0, 1, 5, 3, 2, 4}
	cities := []tsp.Point2D{
		{X: 50, Y: 0},
		{X: 100, Y: 100},
		{X: 100, Y: 250},
		{X: 50, Y: 150},
		{X: 0, Y: 250},
		{X: 0, Y: 100},
	}
	eval := tsp.NewSymmetricEvaluator[int](cities)
	pop := evolve.NewPopulationOf[[]int]([][]int{
		{0, 1, 2, 3, 4, 5}, // AECFDB
		{0, 2, 4, 1, 3, 5}, // ACFDBE
		// {0, 3, 4, 2, 1, 5}, // ADECBF
	}, eval)

	xover := Collision[int]{
		Probability: generator.Const(1.),
		EdgeWeight: func(i, j int) float64 {
			return eval.Distances[i][j]
		},
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	totals := set.NewOf[float64]()
	for i := 0; i < 100000; i++ {
		xover.Apply(pop, rng)
		for c := 0; c < pop.Len(); c++ {
			f := pop.Evaluator.Fitness(pop.Candidates[c])
			pop.Fitness[c] = f
		}
		sort.Sort(pop)
		totals.Insert(pop.Fitness[0])
	}

	fmt.Println("optimum", eval.Fitness(optimum))
	fmt.Println("fitnesses", totals)
}

// TODO(arl) sice CollisionCrossover is TSP-specific we could move it into pkg/tsp (so tests could benefit from the testdata package directly)
//
//go:embed testdata/berlin52.tsp
var berlin52 []byte

func TestCollisionCrossover_collide(t *testing.T) {
	f, err := tsp.Load(bytes.NewReader(berlin52))
	if err != nil {
		t.Fatal(err)
	}
	e := tsp.NewSymmetricEvaluator[int](f.Nodes)
	c := Collision[int]{
		EdgeWeight: func(i, j int) float64 {
			return e.Distances[i][j]
		},
	}

	p1 := []int{0, 1, 2, 3, 4}
	p2 := []int{3, 2, 4, 0, 1}

	seed := 0
	mt := mt19937.New()
	mt.Seed(int64(seed))
	rng := rand.New(mt)
	cost0, cost1 := e.Fitness(p1), e.Fitness(p2)
	off0, off1 := c.collide(p1, p2, int(cost0), int(cost1), rng)
	fmt.Println(off0, off1)
}

func TestCollision2Crossover_collide(t *testing.T) {
	f, err := tsp.Load(bytes.NewReader(berlin52))
	if err != nil {
		t.Fatal(err)
	}
	e := tsp.NewSymmetricEvaluator[int](f.Nodes)
	c := Collision[int]{
		EdgeWeight: func(i, j int) float64 {
			return e.Distances[i][j]
		},
	}

	p1 := []int{0, 1, 2, 3, 4}
	p2 := []int{3, 2, 4, 0, 1}

	mt := mt19937.New()
	rng := rand.New(mt)
	cost0, cost1 := e.Fitness(p1), e.Fitness(p2)
	off0, off1 := c.collide(p1, p2, int(cost0), int(cost1), rng)
	_, _ = off0, off1
}

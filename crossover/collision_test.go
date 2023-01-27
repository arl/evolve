package crossover

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
	"time"

	"github.com/arl/evolve"
	"github.com/arl/evolve/generator"
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
	eval := tsp.NewSymmetricEvaluator(cities)
	pop := evolve.NewPopulationOf[[]int]([][]int{
		{0, 4, 2, 5, 3, 1}, // AECFDB
		{0, 2, 5, 3, 1, 4}, // ACFDBE
		// {0, 3, 4, 2, 1, 5}, // ADECBF
	}, eval)

	xover := Collision{
		Probability: generator.Const(1.),
		CalcDistance: func(i, j int) float64 {
			return eval.Distances[i][j]
		},
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	best := math.MaxFloat64
	for i := 0; i < 100; i++ {
		xover.Apply(pop, rng)
		for c := 0; c < pop.Len(); c++ {
			f := pop.Evaluator.Fitness(pop.Candidates[c])
			pop.Fitness[c] = f
			if f < best {
				best = f
			}
		}
		fmt.Println(pop.Candidates, pop.Fitness, "cur best", best)
	}

	fmt.Println("optimum", eval.Fitness(optimum))
}

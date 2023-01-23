package crossover

import "math/rand"

// The Collision
type Collision[T comparable] struct{}

// CONTINUER
// paper: https://arxiv.org/pdf/1801.02335.pdf

// Mate mates 2 parents and generates a pair of offsprings with Collision. the number
// of cut points is unused.
func (p Collision[T]) Mate(x1, x2 []T, nxpts int, rng *rand.Rand) (y1, y2 []T) {
	if len(x1) != len(x2) {
		panic("Collision cannot mate parents of different lengths")
	}

	// Create empty genotypes.
	y1 = make([]T, len(x1))
	y2 = make([]T, len(x1))

	y1[0] = x1[0]
	y2[0] = x2[0]

	// Keep track of initialized indices in offsprings
	init := make(map[int]struct{})
	init[0] = struct{}{}

	// Map x1 values to their index
	x1pos := make(map[T]int)
	for idx, val := range x1 {
		x1pos[val] = idx
	}

	// Copy cycle to children.
	i := 0
	for {
		// j is the index of x2[i] in x1
		j := x1pos[x2[i]]
		if j == 0 {
			// cycle end
			break
		}
		y1[j] = x1[j]
		y2[j] = x2[j]
		i = j
		init[i] = struct{}{}
	}

	// Copy untouched positions from p1 to c2 and from p2 to c1.
	for i := range x1 {
		if _, ok := init[i]; !ok {
			y1[i] = x2[i]
			y2[i] = x1[i]
		}
	}

	return
}

/*
C# implementation of Collision Crossover
Courtesy of Prof. Ahmad Hassanat

//do crossover from 2 given solutions and output 2 children using phisical collison

// the two chromosomes are collided in opposite direction v1, and v2 velocity of both

//m1 and m2 mass of each node, which is total distance between current node previous and next

TwoOffsprings crossoverCollision(Solution Parent1, Solution Parent2)

{

    TwoOffsprings temp = new TwoOffsprings();

    temp.ch1.Chromosome = new List<int>(CityNum);

    temp.ch2.Chromosome = new List<int>(CityNum);

    int cut = rand.Next(2, CityNum - 1);

    List<bool> visited1 = new List<bool>(CityNum);

    for (int i = 0; i < CityNum; i++)

        visited1.Add(false);

    List<bool> visited2 = new List<bool>(CityNum);

    for (int i = 0; i < CityNum; i++)

        visited2.Add(false);



    //collision

    temp.ch1.Chromosome.Add(0); visited1[0] = true;

    temp.ch2.Chromosome.Add(0); visited2[0] = true;



    double v1, v2;

    double m1, m2, v1p, v2p;

    v1 = (double)(rand.Next(1, (int)Parent1.Cost));

    v2 = (double)(-1.0 * rand.Next(1, (int)Parent2.Cost));

    for (int j = 1; j < Parent1.Chromosome.Count - 1; j++)

    {

        m1 = ED(CitiesXY[Parent1.Chromosome[j - 1]], CitiesXY[Parent1.Chromosome[j]]) +

            ED(CitiesXY[Parent1.Chromosome[j]], CitiesXY[Parent1.Chromosome[j + 1]]);



        m2 = ED(CitiesXY[Parent2.Chromosome[j - 1]], CitiesXY[Parent2.Chromosome[j]]) +

            ED(CitiesXY[Parent2.Chromosome[j]], CitiesXY[Parent2.Chromosome[j + 1]]);



        v1p = (v1 * (m1 - m2) / (m1 + m2)) + (v2 * 2.0 * m2 / (m1 + m2));

        v2p = (v1 * 2.0 * m1 / (m1 + m2)) - (v2 * (m1 - m2) / (m1 + m2));

        //add to ch1

        if (v1p <= 0)

        {

            temp.ch1.Chromosome.Add(Parent1.Chromosome[j]);

            visited1[Parent1.Chromosome[j]] = true;

        }

        else

            temp.ch1.Chromosome.Add(0);

        //add to ch1

        if (v2p >= 0)

        {

            temp.ch2.Chromosome.Add(Parent2.Chromosome[j]);

            visited2[Parent2.Chromosome[j]] = true;

        }

        else

            temp.ch2.Chromosome.Add(0);

    }

    //do it for the last gene

    int jj = Parent1.Chromosome.Count - 1;

    m1 = ED(CitiesXY[Parent1.Chromosome[jj - 1]], CitiesXY[Parent1.Chromosome[jj]]) +

        ED(CitiesXY[Parent1.Chromosome[jj]], CitiesXY[Parent1.Chromosome[0]]);



    m2 = ED(CitiesXY[Parent2.Chromosome[jj - 1]], CitiesXY[Parent2.Chromosome[jj]]) +

        ED(CitiesXY[Parent2.Chromosome[jj]], CitiesXY[Parent2.Chromosome[0]]);



    v1p = (v1 * (m1 - m2) / (m1 + m2)) + (v2 * 2.0 * m2 / (m1 + m2));

    v2p = (v1 * 2.0 * m1 / (m1 + m2)) - (v2 * (m1 - m2) / (m1 + m2));

    //add to ch1

    if (v1p <= 0)//if the city refelected or stopped in parent1

    {

        temp.ch1.Chromosome.Add(Parent1.Chromosome[jj]);

        visited1[Parent1.Chromosome[jj]] = true;

    }

    else

        temp.ch1.Chromosome.Add(0);//check fpr zeros later

    //add to ch1

    if (v2p >= 0)//if the city refelected or stopped in parent2

    {

        temp.ch2.Chromosome.Add(Parent2.Chromosome[jj]);

        visited2[Parent2.Chromosome[jj]] = true;

    }

    else

        temp.ch2.Chromosome.Add(0);//check fpr zeros later



    //fill the rest of genes the one which is not refelcted or stopped from the collision

    int indx = temp.ch1.Chromosome.IndexOf(0, 1);

   // while(indx){

    for (int j = 1; j < jj + 1; j++)

    {

        if (indx==-1) {

            j = jj + 1;

            break;

        }



        if (!visited1[Parent2.Chromosome[j]])

        {

            temp.ch1.Chromosome[indx] = Parent2.Chromosome[j];//fill the rest of cromosome1 from parent2

            visited1[Parent2.Chromosome[j]] = true;

            indx = temp.ch1.Chromosome.IndexOf(0, 1);

        }

        else

            continue;

    }

    indx = temp.ch2.Chromosome.IndexOf(0, 1);

    for (int j = 1; j < jj + 1; j++)

    {

        if (indx==-1)

        {

            j = jj + 1;

            break;

        }

        if (!visited2[Parent1.Chromosome[j]])

        {

            temp.ch2.Chromosome[indx] = Parent1.Chromosome[j];//fill the rest of cromosome2 from parent1

            visited2[Parent1.Chromosome[j]] = true;

            indx = temp.ch2.Chromosome.IndexOf(0, 1);

        }

        else

            continue;

    }

    temp.ch1.Cost = (ulong)CalcCostED(temp.ch1, jj + 1);

    temp.ch2.Cost = (ulong)CalcCostED(temp.ch2, jj + 1);

    return temp;//output 2 children

}
*/

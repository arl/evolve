package islands

import (
	"math/rand"

	"github.com/aurelien-rainone/evolve/framework"
)

// RingMigration migrates a fixed number of individuals from each island to the
// adjacent island. Operates as if the islands are arranged in a ring with
// migration occurring in a clockwise direction. The individuals to be migrated
// are chosen completely at random.
type RingMigration struct{}

// Migrate migrates a fixed number of individuals from each island to the
// adjacent island. Operates as if the islands are arranged in a ring with
// migration occurring in a clockwise direction. The individuals to be migrated
// are chosen completely at random.
//
// islandPopulations is a list of the populations of each island.
// migrantCount is the number of (randomly selected) individuals to be moved on
// from each island.
// rng is a source of randomness.
func (mig RingMigration) Migrate(islandPopulations []framework.EvaluatedPopulation, migrantCount int, rng *rand.Rand) {
	// The first batch of immigrants is from the last island to the first.
	lastIslandIdx := len(islandPopulations) - 1
	lastIsland := islandPopulations[lastIslandIdx]
	framework.ShuffleEvaluatedPopulation(lastIsland, rng)
	migrants := lastIsland[len(lastIsland)-migrantCount:]
	immigrants := make(framework.EvaluatedPopulation, len(migrants))

	for iidx, island := range islandPopulations {

		// Migrants from the last island are immigrants for this island.
		copy(immigrants, migrants)
		if iidx != lastIslandIdx {
			// We've already migrated individuals from the last island.
			// Select the migrants that will move to the next island to make
			// room for the immigrants here. Randomise the population so
			// that there is no bias concerning which individuals are
			// migrated.
			framework.ShuffleEvaluatedPopulation(island, rng)

			copy(migrants, island[len(island)-migrantCount:])
		}
		// Copy the immigrants over the last members of the population
		// (those that are themselves migrating to the next island).
		copy(island[len(island)-migrantCount:], immigrants)
	}
}

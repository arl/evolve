package islands

import (
	"math/rand"

	"github.com/aurelien-rainone/evolve/framework"
)

// RandomMigration migrates a fixed number of candidates away from each island.
// Which individuals are migrated is determined randomly and which islands they
// move to is also random. This contrasts with the more ordered migration
// offered by RingMigration.
//
// If the migration count is greater than one, it is possible (probable) that
// migrants from the same island will be moved to different islands. It is also
// possible that when a migrant's destination is randomly chosen, it gets sent
// back to the island that it came from.
type RandomMigration struct{}

// Migrate migrates a fixed number of candidates away from each island. Which
// individuals are migrated is determined randomly and which islands they move
// to is also random.
//
// If the migration count is greater than one, it is possible (probable)
// that migrants from the same island will be moved to different islands. It
// is also possible that when a migrant's destination is randomly chosen, it
// gets sent back to the island that it came from.
// - islandPopulations is a list of the populations of each island.
// - migrantCount is the number of (randomly selected) individuals to be
// moved on from each island.
// - rng is a source of randomness.
func (mig RandomMigration) Migrate(islandPopulations []framework.EvaluatedPopulation, migrantCount int, rng *rand.Rand) {
	migrants := make(framework.EvaluatedPopulation, 0, migrantCount*len(islandPopulations))

	var ind *framework.EvaluatedCandidate

	for iidx, island := range islandPopulations {
		framework.ShuffleEvaluatedPopulation(island, rng)
		for i := 0; i < migrantCount; i++ {
			ind, island = island[len(island)-1], island[:len(island)-1]
			migrants = append(migrants, ind)
		}
		// we modified the slice underlying array but the original slice still
		// has the same size because we popped from a copy of it
		islandPopulations[iidx] = island
	}
	framework.ShuffleEvaluatedPopulation(migrants, rng)

	var migrantIdx int
	for iidx, island := range islandPopulations {
		for i := 0; i < migrantCount; i++ {
			island = append(island, migrants[migrantIdx])
			migrantIdx++
		}
		// we modified the slice underlying array but the original slice still
		// has the same because we appended to a copy of it
		islandPopulations[iidx] = island
	}
}

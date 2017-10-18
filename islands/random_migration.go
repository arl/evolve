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

	var island framework.EvaluatedPopulation

	for i := 0; i < len(islandPopulations); i++ {
		island = islandPopulations[i]
		// because we shuffle island population ...
		framework.ShuffleEvaluatedPopulation(island, rng)
		// taking N migrants is the same as taking N random individuals from the
		// island. Keep those migrants in a slice
		for j := 0; j < migrantCount; j++ {
			migrants = append(migrants, island[len(island)-1])
		}
	}
	// shuffle the migrants
	framework.ShuffleEvaluatedPopulation(migrants, rng)

	var migrantIdx int
	for i := 0; i < len(islandPopulations); i++ {
		island = islandPopulations[i]
		// replace the last N individuals of an island with N random migrants
		for j := 0; j < migrantCount; j++ {
			island[len(island)-migrantCount+j] = migrants[migrantIdx]
			migrantIdx++
		}
	}
}

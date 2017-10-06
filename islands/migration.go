package islands

import (
	"math/rand"

	"github.com/aurelien-rainone/evolve/framework"
)

// Migration is a strategy interface for different ways of migrating individuals
// between islands in IslandEvolution.
type Migration interface {

	// Migrate performs the actual migration of individual. It moves candidates
	// betweens islands.
	//
	// - islandPopulations is the populations of each island in the system.
	// - migrantCount is the number of individuals to move from each island.
	// - rng is a source of randomness.
	Migrate(islandPopulations []framework.EvaluatedPopulation, migrantCount int, rng *rand.Rand)
}

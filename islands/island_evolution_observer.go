package islands

import "github.com/aurelien-rainone/evolve/framework"

// IslandEvolutionObserver is specialisation of framework.EvolutionObserver
// that, as well as receiving global population updates (at the end of each
// epoch), can receive individual island population updates (at the end of each
// generation on each island).
type IslandEvolutionObserver interface {
	framework.EvolutionObserver

	// IslandPopulationUpdate is called to notify the listener of the state of
	// the population of an individual island.
	//
	// This will be called once for each generation on each island.
	// islandIndex identifies which individual island the data comes from.
	// Indices start at zero and are sequential.
	// data is the latest data from the evolution on the specified island.
	IslandPopulationUpdate(islandIndex int, data *framework.PopulationData)
}

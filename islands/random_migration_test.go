package islands

import (
	"math/rand"
	"testing"

	"github.com/aurelien-rainone/evolve/framework"
	"github.com/stretchr/testify/assert"
)

func TestRandomMigrationZeroMigration(t *testing.T) {
	// Make sure that nothing strange happens when there is no migration.
	migration := RandomMigration{}
	rng := rand.New(rand.NewSource(99))

	islandPopulations := []framework.EvaluatedPopulation{
		createTestPopulation("A", "A", "A"),
		createTestPopulation("B", "B", "B"),
		createTestPopulation("C", "C", "C"),
	}

	migration.Migrate(islandPopulations, 0, rng)
	assert.Len(t, islandPopulations, 3, "wrong number of populations after migration")

	testPopulationContents(t, islandPopulations[0], "A", "A", "A")
	testPopulationContents(t, islandPopulations[1], "B", "B", "B")
	testPopulationContents(t, islandPopulations[2], "C", "C", "C")
}

func TestRandomMigrationNonZeroMigration(t *testing.T) {
	// Make sure that nothing strange happens when the entire island is migrated.
	migration := RandomMigration{}
	rng := rand.New(rand.NewSource(99))

	islandPopulations := []framework.EvaluatedPopulation{
		createTestPopulation("A", "A", "A"),
		createTestPopulation("B", "B", "B"),
		createTestPopulation("C", "C", "C"),
	}

	migration.Migrate(islandPopulations, 3, rng)
	assert.Len(t, islandPopulations, 3, "wrong number of populations after migration")

	// Each population should still have 3 members (but it's not sure which members).
	assert.Len(t, islandPopulations[0], 3, "wrong population size")
	assert.Len(t, islandPopulations[1], 3, "wrong population size")
	assert.Len(t, islandPopulations[2], 3, "wrong population size")
}

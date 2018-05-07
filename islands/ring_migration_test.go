package islands

import (
	"math/rand"
	"testing"

	"github.com/aurelien-rainone/evolve/internal/test"
	"github.com/aurelien-rainone/evolve/pkg/api"
	"github.com/stretchr/testify/assert"
)

func TestRingMigrationZeroMigration(t *testing.T) {
	// Make sure that nothing strange happens when there is no migration.
	migration := RingMigration{}
	rng := rand.New(rand.NewSource(99))

	islandPopulations := []api.EvaluatedPopulation{
		test.CreateTestPopulation("A", "A", "A"),
		test.CreateTestPopulation("B", "B", "B"),
		test.CreateTestPopulation("C", "C", "C"),
	}
	migration.Migrate(islandPopulations, 0, rng)
	assert.Len(t, islandPopulations, 3, "wrong number of populations after migration")

	test.AssertPopulationContents(t, islandPopulations[0], "A", "A", "A")
	test.AssertPopulationContents(t, islandPopulations[1], "B", "B", "B")
	test.AssertPopulationContents(t, islandPopulations[2], "C", "C", "C")
}

func TestRingMigrationFullMigration(t *testing.T) {
	// Make sure that nothing strange happens when the entire island is migrated.
	migration := RingMigration{}
	rng := rand.New(rand.NewSource(99))

	islandPopulations := []api.EvaluatedPopulation{
		test.CreateTestPopulation("A", "A", "A"),
		test.CreateTestPopulation("B", "B", "B"),
		test.CreateTestPopulation("C", "C", "C"),
	}

	migration.Migrate(islandPopulations, 3, rng)
	assert.Len(t, islandPopulations, 3, "wrong number of populations after migration")

	test.AssertPopulationContents(t, islandPopulations[0], "C", "C", "C")
	test.AssertPopulationContents(t, islandPopulations[1], "A", "A", "A")
	test.AssertPopulationContents(t, islandPopulations[2], "B", "B", "B")
}

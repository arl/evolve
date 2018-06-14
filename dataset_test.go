package evolve

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testDataSet []float64

func init() {
	testDataSet = []float64{1, 2, 3, 4, 5}
}

// Make sure that the data set's capacity grows correctly as
// more values are added.
func TestDataSetCapacityIncrease(t *testing.T) {
	data := NewDataset(3)
	assert.Empty(t, data.Len(), "Initial size should be 0.")
	data.AddValue(1)
	data.AddValue(2)
	data.AddValue(3)
	assert.Equal(t, data.Len(), 3)
	// Add a value to take the size beyond the initial capacity.
	data.AddValue(4)
	assert.Equal(t, data.Len(), 4)
}

func TestDataSetAggregate(t *testing.T) {
	data := NewDataset(len(testDataSet))
	data.AddValues(testDataSet...)
	assert.Equal(t, round(data.Sum()), 15)
}

func TestDataSetProduct(t *testing.T) {
	data := NewDataset(len(testDataSet))
	data.AddValues(testDataSet...)
	product := int(factorial(5))
	assert.Equal(t, round(data.Product()), product)
}

func TestDataSetMinimum(t *testing.T) {
	data := NewDataset(0)
	data.AddValue(4)
	assert.Equal(t, data.Min(), 4.0)
	data.AddValue(7)
	assert.Equal(t, data.Min(), 4.0)
	data.AddValue(2)
	assert.Equal(t, data.Min(), 2.0)
	data.AddValue(-9)
	assert.Equal(t, data.Min(), -9.0)
}

func TestDataSetMaximum(t *testing.T) {
	data := NewDataset(0)
	data.AddValue(9)
	assert.Equal(t, data.Max(), 9.0)
	data.AddValue(8)
	assert.Equal(t, data.Max(), 9.0)
	data.AddValue(-15)
	assert.Equal(t, data.Max(), 9.0)
	data.AddValue(12)
	assert.Equal(t, data.Max(), 12.0)
}

func TestDataSetMedian(t *testing.T) {
	data := NewDataset(0)
	data.AddValue(15)
	assert.Equal(t, data.Median(), 15.0)
	data.AddValue(17)
	assert.Equal(t, round(data.Median()), 16)
	data.AddValue(102)
	assert.Equal(t, round(data.Median()), 17)
}

func TestDataSetArithmeticMean(t *testing.T) {
	data := NewDataset(len(testDataSet))
	data.AddValues(testDataSet...)
	assert.Equal(t, round(data.ArithmeticMean()), 3)
}

func TestDataSetGeometricMean(t *testing.T) {
	data := NewDataset(len(testDataSet))
	data.AddValues(testDataSet...)
	product := float64(factorial(5))
	assert.Equal(t, data.GeometricMean(), math.Pow(product, 0.2))
}

func TestDataSetHarmonicMean(t *testing.T) {
	data := NewDataset(0)
	data.AddValues(1, 2, 4, 4)
	// Reciprocals are 1, 1/2, 1/4 and 1/4.
	// Sum of reciprocals is 2.  Therefore, harmonic mean is 4/2 = 2.
	assert.Equal(t, data.HarmonicMean(), 2.0)
}

func TestDataSetMeanDeviation(t *testing.T) {
	data := NewDataset(len(testDataSet))
	data.AddValues(testDataSet...)
	assert.Equal(t, data.MeanDeviation(), 1.2)
}

func TestDataSetPopulationVariance(t *testing.T) {
	data := NewDataset(len(testDataSet))
	data.AddValues(testDataSet...)
	assert.Equal(t, round(data.Variance()), 2)
}

func TestDataSetSampleVariance(t *testing.T) {
	data := NewDataset(len(testDataSet))
	data.AddValues(testDataSet...)
	assert.Equal(t, data.SampleVariance(), 2.5)
}

func TestDataSetPopulationStandardDeviation(t *testing.T) {
	data := NewDataset(len(testDataSet))
	data.AddValues(testDataSet...)
	assert.Equal(t, data.StandardDeviation(), math.Sqrt(2))
}

func TestDataSetSampleStandardDeviation(t *testing.T) {
	data := NewDataset(len(testDataSet))
	data.AddValues(testDataSet...)
	assert.Equal(t, data.SampleStandardDeviation(), math.Sqrt(2.5))
}

// Check that an appropriate exception is thrown when attempting to
// calculate stats without any data.
func TestDataSetEmptyDataSet(t *testing.T) {
	data := NewDataset(10)
	assert.Panics(t, func() { data.ArithmeticMean() })
}

// round rounds floats into integer numbers.
// FIXME: Remove this function when math.Round will exist (in Go 1.10)
func round(a float64) int {
	if a < 0 {
		return int(a - 0.5)
	}
	return int(a + 0.5)
}

func factorial(n uint64) (result uint64) {
	if n > 0 {
		result = n * factorial(n-1)
		return result
	}
	return 1
}

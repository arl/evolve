package evolve

import (
	"math"
	"testing"
)

// Make sure that the data set's capacity grows correctly as
// more values are added.
func TestDataSetCapacityIncrease(t *testing.T) {
	data := NewDataset(3)
	if data.Len() != 0 {
		t.Fatalf("Len() = %v, want %v", data.Len(), 0)
	}
	data.AddValue(1)
	data.AddValue(2)
	data.AddValue(3)
	if data.Len() != 3 {
		t.Fatalf("Len() = %v, want %v", data.Len(), 3)
	}
	// Add a value to take the size beyond the initial capacity.
	data.AddValue(4)
	if data.Len() != 4 {
		t.Fatalf("Len() = %v, want %v", data.Len(), 4)
	}
}

func TestDataSetAggregate(t *testing.T) {
	vals := []float64{1, 2, 3, 4, 5}
	data := NewDataset(len(vals))
	data.AddValues(vals...)
	if got := math.Round(data.Sum()); got != 15 {
		t.Errorf("got %v, want %v", got, 15)
	}
}

func TestDataSetProduct(t *testing.T) {
	vals := []float64{1, 2, 3, 4, 5}
	data := NewDataset(len(vals))
	data.AddValues(vals...)
	product := 1.0
	for _, v := range vals {
		product *= v
	}
	if got := math.Round(data.Product()); got != product {
		t.Errorf("got %v, want %v", got, product)
	}
}

func TestDataSetMinimum(t *testing.T) {
	data := NewDataset(0)
	data.AddValue(4)
	if data.Min() != 4 {
		t.Errorf("got %v, want %v", data.Min(), 4)
	}
	data.AddValue(7)
	if data.Min() != 4 {
		t.Errorf("got %v, want %v", data.Min(), 4)
	}
	data.AddValue(2)
	if data.Min() != 2 {
		t.Errorf("got %v, want %v", data.Min(), 2)
	}
	data.AddValue(-9)
	if data.Min() != -9 {
		t.Errorf("got %v, want %v", data.Min(), -9)
	}
}

func TestDataSetMaximum(t *testing.T) {
	data := NewDataset(0)
	data.AddValue(9)
	if data.Max() != 9 {
		t.Errorf("got %v, want %v", data.Max(), 9)
	}
	data.AddValue(8)
	if data.Max() != 9 {
		t.Errorf("got %v, want %v", data.Max(), 9)
	}
	data.AddValue(-15)
	if data.Max() != 9 {
		t.Errorf("got %v, want %v", data.Max(), 9)
	}
	data.AddValue(12)
	if data.Max() != 12 {
		t.Errorf("got %v, want %v", data.Max(), 12)
	}
}

func TestDataSetMedian(t *testing.T) {
	data := NewDataset(0)
	data.AddValue(15)
	if data.Median() != 15 {
		t.Errorf("got %v, want %v", data.Median(), 15)
	}
	data.AddValue(17)
	if data.Median() != 16 {
		t.Errorf("got %v, want %v", data.Median(), 16)
	}
	data.AddValue(102)
	if data.Median() != 17 {
		t.Errorf("got %v, want %v", data.Median(), 17)
	}
}

func TestDataSetArithmeticMean(t *testing.T) {
	vals := []float64{1, 2, 3, 4, 5}
	data := NewDataset(len(vals))
	data.AddValues(vals...)
	if data.ArithmeticMean() != 3 {
		t.Errorf("got %v, want %v", data.ArithmeticMean(), 3)
	}
}

func TestDataSetGeometricMean(t *testing.T) {
	vals := []float64{1, 2, 3, 4, 5}
	data := NewDataset(len(vals))
	data.AddValues(vals...)
	product := 1.0
	for _, v := range vals {
		product *= v
	}

	want := math.Pow(product, 0.2)
	if data.GeometricMean() != want {
		t.Errorf("got %v, want %v", data.GeometricMean(), want)
	}
}

func TestDataSetHarmonicMean(t *testing.T) {
	data := NewDataset(0)
	data.AddValues(1, 2, 4, 4)
	// Reciprocals are 1, 1/2, 1/4 and 1/4.
	// Sum of reciprocals is 2.  Therefore, harmonic mean is 4/2 = 2.
	if data.HarmonicMean() != 2 {
		t.Errorf("got %v, want %v", data.HarmonicMean(), 2)
	}
}

func TestDataSetMeanDeviation(t *testing.T) {
	vals := []float64{1, 2, 3, 4, 5}
	data := NewDataset(len(vals))
	data.AddValues(vals...)
	if data.MeanDeviation() != 1.2 {
		t.Errorf("got %v, want %v", data.MeanDeviation(), 1.2)
	}
}

func TestDataSetPopulationVariance(t *testing.T) {
	vals := []float64{1, 2, 3, 4, 5}
	data := NewDataset(len(vals))
	data.AddValues(vals...)
	if data.Variance() != 2 {
		t.Errorf("got %v, want %v", data.Variance(), 2)
	}
}

func TestDataSetSampleVariance(t *testing.T) {
	vals := []float64{1, 2, 3, 4, 5}
	data := NewDataset(len(vals))
	data.AddValues(vals...)
	if data.SampleVariance() != 2.5 {
		t.Errorf("got %v, want %v", data.SampleVariance(), 2.5)
	}
}

func TestDataSetPopulationStandardDeviation(t *testing.T) {
	vals := []float64{1, 2, 3, 4, 5}
	data := NewDataset(len(vals))
	data.AddValues(vals...)
	if data.StandardDeviation() != math.Sqrt(2) {
		t.Errorf("got %v, want %v", data.StandardDeviation(), math.Sqrt(2))
	}
}

func TestDataSetSampleStandardDeviation(t *testing.T) {
	vals := []float64{1, 2, 3, 4, 5}
	data := NewDataset(len(vals))
	data.AddValues(vals...)
	if data.SampleStandardDeviation() != math.Sqrt(2.5) {
		t.Errorf("got %v, want %v", data.SampleStandardDeviation(), math.Sqrt(2.5))
	}
}

// Check that an appropriate exception is thrown when attempting to
// calculate stats without any data.
func TestDataSetEmptyDataSet(t *testing.T) {
	didpanic := false
	defer func() {
		if r := recover(); r != nil {
			didpanic = true
		}
	}()

	NewDataset(10).ArithmeticMean()
	if !didpanic {
		t.Errorf("ArithmeticMean on empty dataset should have panicked")
	}
}

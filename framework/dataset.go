package framework

import (
	"math"
	"sort"
)

const (
	defaultCapacity = 50
)

// DataSet is a utility struct for calculating statistics for a finite data set.
type DataSet struct {
	dataSet          []float64
	total            float64
	product          float64
	reciprocalSum    float64
	minimum, maximum float64
}

// DataSetOption is the type of functions used to specify options during the
// creation of DataSet objects.
type DataSetOption func(*DataSet)

// NewDataSet creates an empty data set with a default initial capacity.
func NewDataSet(options ...DataSetOption) *DataSet {
	ds := &DataSet{
		minimum:       math.MaxFloat64,
		maximum:       math.SmallestNonzeroFloat64,
		product:       1,
		total:         0,
		reciprocalSum: 0,
	}

	// set dataset options
	for _, option := range options {
		option(ds)
	}

	if ds.dataSet == nil {
		ds.dataSet = make([]float64, 0, defaultCapacity)
	}
	return ds
}

// WithInitialCapacity allows to specify the initial capacity of the
// data set.
//
// The initial capacity for the data set corresponds to the number
// of values that can be added without needing to resize the
// internal data storage.
func WithInitialCapacity(capacity int) DataSetOption {
	return func(ds *DataSet) {
		ds.dataSet = make([]float64, 0, capacity)
	}
}

// WithPrePopulatedDataSet allows to prepopulates the data set with the
// specified values.
func WithPrePopulatedDataSet(dataSet []float64) DataSetOption {
	return func(ds *DataSet) {
		ds.dataSet = make([]float64, len(dataSet))
		copy(ds.dataSet, dataSet)

		for _, value := range ds.dataSet {
			ds.updateStatsWithNewValue(value)
		}
	}
}

// AddValue adds a single value to the data set and updates any statistics that
// are calculated cumulatively.
func (ds *DataSet) AddValue(value float64) {
	ds.dataSet = append(ds.dataSet, value)
	ds.updateStatsWithNewValue(value)
}

func (ds *DataSet) updateStatsWithNewValue(value float64) {
	ds.total += value
	ds.product *= value
	ds.reciprocalSum += 1 / value
	ds.minimum = math.Min(ds.minimum, value)
	ds.maximum = math.Max(ds.maximum, value)
}

func (ds *DataSet) assertNotEmpty() {
	if len(ds.dataSet) == 0 {
		panic("DataSet should not be empty")
	}
}

// Len returns the number of values in this data set.
func (ds *DataSet) Len() int {
	return len(ds.dataSet)
}

// Minimum returns the smallest value in the data set.
//
// panics if the data set is empty.
func (ds *DataSet) Minimum() float64 {
	ds.assertNotEmpty()
	return ds.minimum
}

// Maximum returns the biggest value in the data set.
//
// panics if the data set is empty.
func (ds *DataSet) Maximum() float64 {
	ds.assertNotEmpty()
	return ds.maximum
}

// Median determines the median value of the data set.
//
// If the number of elements is odd, returns the middle element.
// If the number of elements is even, returns the midpoint of the two
// middle elements.
//
// panics if the data set is empty.
func (ds *DataSet) Median() float64 {
	ds.assertNotEmpty()
	// Sort the data (take a copy to do this).
	dataCopy := make([]float64, len(ds.dataSet))
	copy(dataCopy, ds.dataSet)
	sort.Float64s(dataCopy)
	midPoint := len(dataCopy) / 2
	if len(dataCopy)%2 != 0 {
		return dataCopy[midPoint]
	}
	return dataCopy[midPoint-1] + (dataCopy[midPoint]-dataCopy[midPoint-1])/2
}

// Aggregate returns the sum of all values.
//
// panics if the data set is empty.
func (ds *DataSet) Aggregate() float64 {
	ds.assertNotEmpty()
	return ds.total
}

// Product returns the product of all values.
//
// panics if the data set is empty.
func (ds *DataSet) Product() float64 {
	ds.assertNotEmpty()
	return ds.product
}

// ArithmeticMean returns the arithmetic mean of all elements ion the data
// set.
//
// The artithmetic mean of an n-element set is the sum of all the elements
// divided by n. The arithmetic mean is often referred to simply as the "mean"
// or "average" of a data set.
//
// See GeometricMean()
//
// panics if the data set is empty.
func (ds *DataSet) ArithmeticMean() float64 {
	ds.assertNotEmpty()
	return ds.total / float64(len(ds.dataSet))
}

// GeometricMean returns the geometric mean of all elements in the data set.
//
// The geometric mean of an n-element set is the nth-root of the product of all
// the elements. The geometric mean is used for finding the average factor (e.g.
// an average interest rate).
//
// See ArithmeticMean() and HarmonicMean()
//
// panics if the data set is empty.
func (ds *DataSet) GeometricMean() float64 {
	ds.assertNotEmpty()
	return math.Pow(ds.product, 1.0/float64(len(ds.dataSet)))
}

// HarmonicMean returns the harmonic mean of all the elements in the data set.
//
// The harmonic mean of an n-element set is n divided by the sum of the
// reciprocals of the values (where the reciprocal of a value x is 1/x.
// The harmonic mean is used to calculate an average rate (e.g. an average
// speed).
//
// See ArithmeticMean() and GeometricMean()
//
// panics if the data set is empty.
func (ds *DataSet) HarmonicMean() float64 {
	ds.assertNotEmpty()
	return float64(len(ds.dataSet)) / ds.reciprocalSum
}

// MeanDeviation returns the mean absolute deviation of the data set.
//
// The mean deviation is the average (absolute) amount that a single value
// deviates from the arithmetic mean.
//
// See ArithmeticMean(), Variance() and StandardDeviation()
//
// panics if the data set is empty.
func (ds *DataSet) MeanDeviation() float64 {
	mean := ds.ArithmeticMean()
	var diffs float64
	for _, value := range ds.dataSet {
		diffs += math.Abs(mean - value)
	}
	return diffs / float64(len(ds.dataSet))
}

// Variance returns the population variance of the data set.
//
// The variance is a measure of statistical dispersion of the data set.
//
// There are different measures of variance depending on whether the data set is
// itself a finite population or is a sample from some larger population. For
// large data sets the difference is negligible. This method calculates the
// population variance.
//
// See SampleVariance(), StandardDeviation() and MeanDeviation()
//
// panics if the data set is empty.
func (ds *DataSet) Variance() float64 {
	return ds.sumSquaredDiffs() / float64(len(ds.dataSet))
}

// sumSquaredDiffs is an helper method for variance calculations.
//
// It returns the sum of the squares of the differences between
// each value and the arithmetic mean.
//
// panics if the data set is empty.
func (ds *DataSet) sumSquaredDiffs() float64 {
	mean := ds.ArithmeticMean()
	var squaredDiffs float64
	for _, value := range ds.dataSet {
		diff := mean - value
		squaredDiffs += (diff * diff)
	}
	return squaredDiffs
}

// StandardDeviation returns the standard deviation of the population.
//
// The standard deviation is the square root of the variance. This method
// calculates the population standard deviation as opposed to the sample
// standard deviation. For large data sets the difference is negligible.
//
// See SampleStandardDeviation(), Variance() and MeanDeviation()
//
// panics if the data set is empty.
func (ds *DataSet) StandardDeviation() float64 {
	return math.Sqrt(ds.Variance())
}

// SampleVariance returns the sample variance of the data set.
//
// Calculates the variance (a measure of statistical dispersion) of the data
// set. There are different measures of variance depending on whether the data
// set is itself a finite population or is a sample from some larger population.
// For large data sets the difference is negligible. This method calculates the
// sample variance.
//
// See Variance(), SampleStandardDeviation() and MeanDeviation()
//
// panics if the data set is empty.
func (ds *DataSet) SampleVariance() float64 {
	return ds.sumSquaredDiffs() / float64(len(ds.dataSet)-1)
}

// SampleStandardDeviation returns the sample standard deviation of the data
// set.
//
// The sample standard deviation is the square root of the sample variance. For
// large data sets the difference between sample standard deviation and
// population standard deviation is negligible.
//
// See StandardDeviation(), SampleVariance() and MeanDeviation()
//
// panics if the data set is empty.
func (ds *DataSet) SampleStandardDeviation() float64 {
	return math.Sqrt(ds.SampleVariance())
}

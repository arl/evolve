package evolve

import (
	"math"
	"sort"
)

// A Dataset stores some floating point values and compute some statistics
// about the finite dataset composed of all the values it stores.
type Dataset struct {
	values   []float64
	total    float64
	product  float64
	recipsum float64 // reciprocal sum
	min, max float64
}

// NewDataset creates an empty data set with the provided initial capacity.
//
// The initial capacity for the data set corresponds to the number of values
// that can be added without needing to resize the internal data storage.
func NewDataset(capacity int) *Dataset {
	return &Dataset{
		min:      math.MaxFloat64,
		max:      math.SmallestNonzeroFloat64,
		product:  1,
		total:    0,
		recipsum: 0,
		values:   make([]float64, 0, capacity),
	}
}

// AddValue adds a single value to the data set and updates any statistics that
// are calculated cumulatively.
func (ds *Dataset) AddValue(value float64) {
	ds.values = append(ds.values, value)
	ds.update(value)
}

// AddValues adds multiple values to the data set and updates any statistics that
// are calculated cumulatively.
func (ds *Dataset) AddValues(values ...float64) {
	ds.values = append(ds.values, values...)
	for _, value := range ds.values {
		ds.update(value)
	}
}

// Clear clears all the values previously added and resets the statistics. The
// dataset capacity remains unchanged.
func (ds *Dataset) Clear() {
	ds.min = math.MaxFloat64
	ds.max = math.SmallestNonzeroFloat64
	ds.product = 1
	ds.total = 0
	ds.recipsum = 0
	ds.values = ds.values[:0]
}

// update the dataset by considering the new value that has been added
func (ds *Dataset) update(value float64) {
	ds.total += value
	ds.product *= value
	ds.recipsum += 1 / value
	ds.min = math.Min(ds.min, value)
	ds.max = math.Max(ds.max, value)
}

func (ds *Dataset) mustNotEmpty() {
	if len(ds.values) == 0 {
		panic("DataSet should not be empty")
	}
}

// Len returns the number of values in this data set.
func (ds *Dataset) Len() int {
	return len(ds.values)
}

// Min returns the smallest value in the data set.
//
// panics if the data set is empty.
func (ds *Dataset) Min() float64 {
	ds.mustNotEmpty()
	return ds.min
}

// Max returns the biggest value in the data set.
//
// panics if the data set is empty.
func (ds *Dataset) Max() float64 {
	ds.mustNotEmpty()
	return ds.max
}

// Median determines the median value of the data set.
//
// If the number of elements is odd, returns the middle element.
// If the number of elements is even, returns the midpoint of the two
// middle elements.
//
// panics if the data set is empty.
func (ds *Dataset) Median() float64 {
	ds.mustNotEmpty()
	// Sort the data (take a copy to do this)
	// TODO: why exactly ??
	cpy := make([]float64, len(ds.values))
	copy(cpy, ds.values)
	sort.Float64s(cpy)
	middle := len(cpy) / 2
	if len(cpy)%2 != 0 {
		return cpy[middle]
	}
	return cpy[middle-1] + (cpy[middle]-cpy[middle-1])/2
}

// Sum returns the sum of all values.
//
// panics if the data set is empty.
func (ds *Dataset) Sum() float64 {
	ds.mustNotEmpty()
	return ds.total
}

// Product returns the product of all values.
//
// panics if the data set is empty.
func (ds *Dataset) Product() float64 {
	ds.mustNotEmpty()
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
func (ds *Dataset) ArithmeticMean() float64 {
	ds.mustNotEmpty()
	return ds.total / float64(len(ds.values))
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
func (ds *Dataset) GeometricMean() float64 {
	ds.mustNotEmpty()
	return math.Pow(ds.product, 1.0/float64(len(ds.values)))
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
func (ds *Dataset) HarmonicMean() float64 {
	ds.mustNotEmpty()
	return float64(len(ds.values)) / ds.recipsum
}

// MeanDeviation returns the mean absolute deviation of the data set.
//
// The mean deviation is the average (absolute) amount that a single value
// deviates from the arithmetic mean.
//
// See ArithmeticMean(), Variance() and StandardDeviation()
//
// panics if the data set is empty.
func (ds *Dataset) MeanDeviation() float64 {
	mean := ds.ArithmeticMean()
	var diffs float64
	for _, value := range ds.values {
		diffs += math.Abs(mean - value)
	}
	return diffs / float64(len(ds.values))
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
func (ds *Dataset) Variance() float64 {
	return ds.sumSquaredDiffs() / float64(len(ds.values))
}

// sumSquaredDiffs is an helper method for variance calculations.
//
// It returns the sum of the squares of the differences between
// each value and the arithmetic mean.
//
// panics if the data set is empty.
func (ds *Dataset) sumSquaredDiffs() float64 {
	mean := ds.ArithmeticMean()
	var sqdiffs float64
	for _, value := range ds.values {
		diff := mean - value
		sqdiffs += (diff * diff)
	}
	return sqdiffs
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
func (ds *Dataset) StandardDeviation() float64 {
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
func (ds *Dataset) SampleVariance() float64 {
	return ds.sumSquaredDiffs() / float64(len(ds.values)-1)
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
func (ds *Dataset) SampleStandardDeviation() float64 {
	return math.Sqrt(ds.SampleVariance())
}

package evolve

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

type StatsToCSV[T any] struct {
	EveryNth int // write 1 csv record every N generations (default 10)
	wc       io.WriteCloser
	csv      *csv.Writer
}

// NewCSVStats returns an evolution observer that saves statistics in CSV format
// to the file at the given path. When finished using, Close should be called to
// ensure that buffered data is written to disk.
func NewStatsToCSV[T any](path string) (*StatsToCSV[T], error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("csvstats: %v", err)
	}

	csv := csv.NewWriter(f)
	csv.Write([]string{"generation", "best fitness", "mean", "standard deviation", "elapsed (ms)"})

	const defaultEvery = 10
	s := &StatsToCSV[T]{
		wc:       f,
		csv:      csv,
		EveryNth: defaultEvery,
	}
	return s, nil
}

// Close flushes the latest csv data to disk and closes the file.
func (s *StatsToCSV[T]) Close() error {
	s.csv.Flush()
	if err := s.wc.Close(); err != nil {
		return fmt.Errorf("csvstats: %v", err)
	}
	return nil
}

func (s *StatsToCSV[T]) Observe(stats *PopulationStats[T]) {
	if stats.Generation%s.EveryNth != 0 {
		return
	}
	s.csv.Write([]string{
		strconv.Itoa(stats.Generation),
		strconv.FormatFloat(stats.BestFitness, 'f', 3, 64),
		strconv.FormatFloat(stats.Mean, 'f', 3, 64),
		strconv.FormatFloat(stats.StdDev, 'f', 3, 64),
		strconv.Itoa(int(stats.Elapsed.Truncate(time.Millisecond) / time.Millisecond)),
	})
}

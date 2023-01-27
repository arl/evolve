// Package tsp provides types and functions to read TSP files. TSP files
// describes a data set for an instance of the TSP (Travaling Salesperson
// Problem). Specification from:
// http://comopt.ifi.uni-heidelberg.de/software/TSPLIB95/tsp95.pdf

package tsp

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// File is a TSP file.
type File struct {
	// Name identifies the TSP file.
	Name string

	// Type specifies the type of the data contained in the file.
	Type string

	// Comment holds additional comments.
	Comment string

	// Dimension is the dimension of the dataset. Sspecficially, for a TSP or
	// ATSP, it is the dimension is the number of its nodes. For a CVRP, it is
	// the total number of nodes and depots. For a TOUR file it is the dimension
	// of the corresponding problem.
	Dimension int

	// Specifies the truck capacity in a CVRP.
	Capacity int

	// EdgeWeightType specifies how the edge weights (or distances) are given.
	EdgeWeightType string

	// EdgeWeightFormat describes the format of the edge weights if they are
	// given explicitly.
	EdgeWeightFormat string

	// Nodes are the node coordinates, indexed by their number in the file.
	Nodes []Point2D
}

// Point2D is the point coordinates in 2D space.
type Point2D struct{ X, Y float64 }

// Load fills a *File from a TSP file on the filesystem.
func LoadFromFile(path string) (*File, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed loading tsp file %v: %v", path, err)
	}
	defer f.Close()
	return Load(f)
}

// Load fills a *File from a reader reading a TSP file.
func Load(r io.Reader) (*File, error) {
	f := File{}
	s := bufio.NewScanner(r)

scanLoop:
	for s.Scan() {
		if k, v, found := strings.Cut(s.Text(), ":"); found {
			switch strings.TrimSpace(k) {
			case "NAME":
				f.Name = strings.TrimSpace(v)
			case "TYPE":
				f.Type = strings.TrimSpace(v)
			case "COMMENT":
				f.Comment = strings.TrimSpace(v)
			case "DIMENSION":
				dim, err := strconv.Atoi(strings.TrimSpace(v))
				if err != nil {
					return nil, fmt.Errorf("failed to parse line: %v", s.Text())
				}
				f.Dimension = dim
			case "EDGE_WEIGHT_TYPE":
				f.EdgeWeightType = strings.TrimSpace(v)
			default: // Skip other keywords.
			}
			continue
		}

		switch s.Text() {
		case "NODE_COORD_SECTION":
			nodes, err := readNodeCoords(s, f.Dimension)
			if err != nil {
				return nil, fmt.Errorf("reading node coord section: %v", err)
			}
			f.Nodes = nodes
		case "EOF":
			break scanLoop
		default: // Skip other sections.
			continue
		}
	}

	if err := s.Err(); err != nil {
		return nil, fmt.Errorf("error reading the TSP file: %v", err)
	}

	// Validate file with respect to the currently supported feature set.
	if f.EdgeWeightType != "EUC_2D" {
		return nil, fmt.Errorf("unsupported EDGE_WEIGHT_TYPE '%v'", f.EdgeWeightType)
	}

	return &f, nil
}

// reads the NODE_COORD_SECTION.
func readNodeCoords(s *bufio.Scanner, dim int) ([]Point2D, error) {
	if dim <= 0 {
		return nil, fmt.Errorf("dimension is unset")
	}

	nodes := make([]Point2D, dim)
	nread := 0
	for s.Scan() {
		spl := strings.Split(s.Text(), " ")
		if len(spl) < 3 {
			return nil, fmt.Errorf("malformed coord line: %v", s.Text())
		}
		idx, err := strconv.Atoi(spl[0])
		if err != nil {
			return nil, fmt.Errorf("malformed index in line: %v", s.Text())
		}
		if idx < 1 || idx > dim {
			return nil, fmt.Errorf("incorrect index %v, must be [%d %d]", idx, 1, dim)
		}
		x, err := strconv.ParseFloat(spl[1], 64)
		if err != nil {
			return nil, fmt.Errorf("malformed x coord in line: %v", s.Text())
		}
		y, err := strconv.ParseFloat(spl[2], 64)
		if err != nil {
			return nil, fmt.Errorf("malformed y coord in line: %v", s.Text())
		}
		nodes[idx-1] = Point2D{X: x, Y: y}

		if nread++; nread == dim {
			break
		}
	}

	return nodes, nil
}

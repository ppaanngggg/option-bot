package pricing

import (
	"math"

	"gonum.org/v1/gonum/optimize"
)

type MyConverge struct {
	first          bool
	lastX          []float64
	diffThreshold  float64
	count          int
	countThreshold int
}

func (m *MyConverge) Init(dim int) {
	m.first = true
	m.lastX = make([]float64, dim)
	if m.diffThreshold == 0 {
		m.diffThreshold = 1e-4
	}
	if m.countThreshold == 0 {
		m.countThreshold = 2
	}
}

func (m *MyConverge) Converged(loc *optimize.Location) optimize.Status {
	if m.first {
		m.first = false
		copy(m.lastX, loc.X)
		return optimize.NotTerminated
	}
	maxDiff := 0.0
	for i, x := range loc.X {
		diff := math.Abs(x - m.lastX[i])
		if diff > maxDiff {
			maxDiff = diff
		}
	}
	copy(m.lastX, loc.X)
	if maxDiff < m.diffThreshold {
		m.count++
	} else {
		m.count = 0
	}
	if m.count > m.countThreshold {
		return optimize.FunctionConvergence
	}
	return optimize.NotTerminated
}

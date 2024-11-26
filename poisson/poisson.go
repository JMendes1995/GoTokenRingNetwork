package poisson

import (
	"math"
	"math/rand"
)

// PoissonProcess represents a Poisson process random number generator.
type PoissonProcess struct {
	Lambda float64
	Rng    *rand.Rand
}

// PoissonRandom generates a Poisson-distributed random number of events based on a given average rate (lambda).
func (p *PoissonProcess) PoissonRandom() int {
	L := math.Exp(-p.Lambda)
	k := 0
	pVal := 1.0

	for pVal > L {
		k++
		pVal *= p.Rng.Float64()
	}
	return k - 1
}

// ExponentialRandom generates an exponentially distributed random number based on a given rate (lambda).
func (p *PoissonProcess) ExponentialRandom() float64 {
	return -math.Log(1.0-p.Rng.Float64()) / p.Lambda
}


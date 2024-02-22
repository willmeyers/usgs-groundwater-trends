package trends

import (
	"github.com/montanaflynn/stats"
	"math"
	"sort"
)

func meanFloat32(xs []float64) float64 {
	var sum float64

	for _, x := range xs {
		sum += x
	}

	return sum / float64(len(xs))
}

func ApproxMedianSlopeInMM(xs []float64, ys []float64, step int) float64 {
	slopes := make([]float64, 0)
	n := len(xs)

	if n != len(ys) || n == 0 || step < 1 {
		return math.NaN()
	}

	for i := 0; i < n-1; i += step {
		for j := i + 1; j < n; j += step {
			if xs[j] != xs[i] {
				slope := (ys[j] - ys[i]) / (xs[j] - xs[i])
				slopes = append(slopes, slope)
			}
		}
	}

	sort.Float64s(slopes)
	mid := len(slopes) / 2
	if len(slopes)%2 == 0 {
		return ((slopes[mid-1] + slopes[mid]) / 2.0) * 1000.0
	}

	return slopes[mid] * 1000.0
}

func MannKendall(xs []float64) (float64, float64) {
	// Implementation from https://www.statisticshowto.com/wp-content/uploads/2016/08/Mann-Kendall-Analysis-1.pdf
	N := len(xs)

	S := 0
	for i := 0; i < N-1; i++ {
		for j := i + 1; j < N; j++ {
			diff := xs[j] - xs[i]
			if diff > 0 {
				S++
			} else if diff < 0 {
				S--
			}
		}
	}

	unique := make([]float64, 0, N)
	counts := make(map[float64]int)
	for _, v := range xs {
		_, ok := counts[v]
		if !ok {
			unique = append(unique, v)
			counts[v] = 1
		} else {
			counts[v]++
		}
	}
	G := len(unique)

	var variance float64
	if N == G {
		variance = float64(N*(N-1)*(2*N+5)) / 18
	} else {
		var tiSum int
		for _, ti := range counts {
			tiSum += ti * (ti - 1) * (2*ti + 5)
		}
		variance = float64((N*(N-1)*(2*N+5))-tiSum) / 18
	}

	var Z float64

	if S > 0 {
		Z = (float64(S) - 1.0) / math.Pow(variance, 0.5)
	} else if S == 0 {
		Z = 0
	} else {
		Z = (float64(S) + 1.0) / math.Pow(variance, 0.5)
	}

	P := stats.NormCdf(Z, 0.0, 1.0)

	return Z, P
}

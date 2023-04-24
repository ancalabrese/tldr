package distance

import (
	"errors"
	"math"
)

func Cosine(a []float64, b []float64) (cosine float64, err error) {
	count := 0
	aLen := len(a)
	bLen := len(b)

	if aLen > bLen {
		count = aLen
	} else {
		count = bLen
	}

	sumA := 0.0
	s1 := 0.0
	s2 := 0.0

	for i := 0; i < count; i++ {
		if i >= aLen {
			s2 += math.Pow(b[i], 2)
			continue
		}
		if i >= bLen {
			s1 += math.Pow(a[i], 2)
			continue
		}
		sumA += a[i] * b[i]
		s1 += math.Pow(a[i], 2)
		s2 += math.Pow(b[i], 2)
	}
	if s1 == 0 || s2 == 0 {
		return 0.0, errors.New("Vectors should not be null (all zeros)")
	}
	return sumA / (math.Sqrt(s1) * math.Sqrt(s2)), nil
}

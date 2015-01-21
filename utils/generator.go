package utils

import (
	"errors"
	"math/rand"
)

func GeneratorInteger(min int, max int, r *rand.Rand) (int, error) {
	if min > max {
		return -1, errors.New("GeneratorInteger requires that min should be bigger than max...")
	}

	ret := r.Intn(max)

	if ret < min {
		ret = min
	}

	return ret, nil
}

func GeneratorString(rangeStr string, min int, max int, r *rand.Rand) (string, error) {
	erange, err := ExpandRanges(rangeStr)

	if err != nil {
		return "", err
	}

	ret := ""
	if max == 0 {
		max = r.Intn(256)
	}

	nChars := r.Intn(max)
	if nChars <= min {
		nChars = min
	}

	for i := 0; i < nChars; i++ {
		var c int
		n := r.Intn(len(erange))
		if n > 0 {
			c = r.Intn(n)
		} else {
			c = n
		}
		ret += string(erange[c])
	}

	return ret, nil
}

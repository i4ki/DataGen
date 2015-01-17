package utils

import (
	"errors"
	"math/rand"
)
import "time"

func GeneratorInteger(min int, max int) (int, error) {
	if min > max {
		return -1, errors.New("GeneratorInteger requires that min should be bigger than max...")
	}

	ret := rand.Intn(max)

	if ret < min {
		ret = min
	}

	return ret, nil
}

func GeneratorString(rangeStr string, min int, max int) (string, error) {
	erange, err := ExpandRanges(rangeStr)

	if err != nil {
		return "", err
	}

	ret := ""
	rand.Seed(time.Now().UTC().UnixNano())
	if max == 0 {
		max = rand.Intn(256)
	}

	nChars := rand.Intn(max)
	if nChars <= min {
		nChars = min
	}

	for i := 0; i < nChars; i++ {
		var c int
		n := rand.Intn(len(erange))
		if n > 0 {
			c = rand.Intn(n)
		} else {
			c = n
		}
		ret += string(erange[c])
	}

	return ret, nil
}

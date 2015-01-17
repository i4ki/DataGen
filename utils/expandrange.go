package utils

import "errors"

func calcRange(from rune, to rune) (string, error) {
	var retBytes []byte

	if from > to {
		return "", errors.New("Invalid range: ")
	}

	for i := from; i <= to; i++ {
		retBytes = append(retBytes, byte(i))
	}

	return string(retBytes), nil
}

func ExpandRanges(ranges string) (string, error) {
	var result, tmpStr string = "", ""
	var previous rune
	var previousPos int
	hasRange := false
	var err error

	for _, v := range ranges {
		if v != '-' {
			if hasRange {
				tmpStr, err = calcRange(previous, v)
				if err != nil {
					return "", err
				}

				result = result[0:previousPos]
				result += tmpStr

				previous = 0
				previousPos = len(result) - 1
				hasRange = false
			} else {
				if hasRange {
					return "", errors.New("Failed to parse range: " + ranges)
				}

				previous = v
				result += string(v)
				previousPos = len(result) - 1
			}
		} else {
			hasRange = true
		}
	}

	return result, nil
}

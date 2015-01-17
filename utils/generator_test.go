package utils

import (
	"strconv"
	"testing"
)

func TestGeneratorString(t *testing.T) {
	ret, err := GeneratorString("abc", 3, 256)

	if err != nil || (ret == "" || (len(ret) < 3 || len(ret) > 256)) {
		t.Error("Failed to generate correct string: "+ret+", Length: "+strconv.Itoa(len(ret)), err)
	}

	ret, err = GeneratorString("abc", 3, 3)

	if err != nil || (ret == "" || (len(ret) < 3 || len(ret) > 3)) {
		t.Error("Failed to generate correct string: " + ret + ", Length: " + strconv.Itoa(len(ret)))
	}

	ret, err = GeneratorString("a", 1, 1)

	if err != nil || (ret != "a" || (len(ret) < 1 || len(ret) > 1)) {
		t.Error("Failed to generate correct string: " + ret + ", Length: " + strconv.Itoa(len(ret)))
	}

	ret, err = GeneratorString("ab", 2, 2)

	if err != nil || (ret != "aa" || ret != "bb" && ret != "ab" && ret != "ba" && (len(ret) < 2 || len(ret) > 2)) {
		t.Error("Failed to generate correct string: " + ret + ", Length: " + strconv.Itoa(len(ret)))
	}

	ret, err = GeneratorString("0-9", 0, 10)

	if err != nil || (len(ret) < 0 || len(ret) > 10) {
		t.Error("Failed to generate correct string: " + ret + ", Length: " + strconv.Itoa(len(ret)))
	}

	ret, err = GeneratorString("a-zA-Z0-9", 0, 256)

	if err != nil || (len(ret) < 0 || len(ret) > 256) {
		t.Error("Failed to generate correct string: " + ret + ", Length: " + strconv.Itoa(len(ret)))
	}
}

func TestGeneratorInteger(t *testing.T) {
	ret, err := GeneratorInteger(1, 10)

	if err != nil || (ret < 1 || ret > 10) {
		t.Error("Failed to generate integer: ", ret, err)
	}

	ret, err = GeneratorInteger(1, 1)

	if err != nil || (ret < 1 || ret > 1) {
		t.Error("Failed to generate integer: ", ret, err)
	}

	ret, err = GeneratorInteger(1, 2)

	if err != nil || (ret < 1 || ret > 2) {
		t.Error("Failed to generate integer: ", ret, err)
	}

	ret, err = GeneratorInteger(0, 1)

	if err != nil || (ret < 0 || ret > 1) {
		t.Error("Failed to generate integer: ", ret, err)
	}

	ret, err = GeneratorInteger(0, 256)

	if err != nil || (ret < 0 || ret > 256) {
		t.Error("Failed to generate integer: ", ret, err)
	}

	ret, err = GeneratorInteger(255, 256)

	if err != nil || (ret < 255 || ret > 256) {
		t.Error("Failed to generate integer: ", ret, err)
	}
}

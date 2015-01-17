package utils

import "testing"

func TestExpandRangeRange(t *testing.T) {
	ret, err := ExpandRanges("0-9")
	if ret != "0123456789" {
		t.Error("Failed to expand: ", ret, err)
	}

	ret, err = ExpandRanges("a-z")
	if ret != "abcdefghijklmnopqrstuvwxyz" {
		t.Error("Failed to expand:", ret, err)
	}

	ret, err = ExpandRanges("0-9a-z")
	if ret != "0123456789abcdefghijklmnopqrstuvwxyz" {
		t.Error("Failed to expand:", ret, err)
	}

	ret, err = ExpandRanges("0-9a-zA-Z")
	if ret != "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		t.Error("Failed to expand:", ret, err)
	}

	ret, err = ExpandRanges("012")

	if ret != "012" {
		t.Error("Failed to expand: ", ret, err)
	}
}

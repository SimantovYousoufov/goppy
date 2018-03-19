package main

import (
	"testing"
)

func TestItCanTruncateString(t *testing.T) {
	s := `some string
with 3
newlines
here`

	trunc := truncateString(s, 22, 3)

	expect := `some string
with 3
new...`

	if trunc != expect {
		t.Fatalf("Truncated strings are not equal. Expected:\n %s \n Received:\n %s \n", expect, trunc)
	}
}

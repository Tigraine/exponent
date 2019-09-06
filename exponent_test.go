package exponent

import (
	"testing"
	"time"
)

func TestExponentialBackoff(t *testing.T) {
	e := NewExponent(13, 10*time.Second)
	var n int
	t.Log("starting")
	for e.Do() {
		n++
		t.Logf("Sleeping for %v for run %d", e.WaitFor(), e.n)
	}
	t.Log("finished")
	if n != 13 {
		t.Errorf("Expected 13 loops but only got %d", n)
	}

}

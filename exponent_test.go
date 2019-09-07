package exponent

import (
	"testing"
)

func TestExponentialBackoffRetries(t *testing.T) {
	e := NewBackoff(13)
	e.strategy = LinearBackoff
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

func TestExponentialBackoffStopsOnSuccess(t *testing.T) {
	e := NewBackoff(10)

	var n int
	for e.Do() {
		n++
		if n == 4 {
			e.Success()
		}
	}
	if n != 4 {
		t.Errorf("Expected to run only 4 times got %d", n)
	}
}

func TestSuccess(t *testing.T) {
	e := NewBackoff(10)
	e.Success()
	if e.done != true {
		t.Errorf("Expected done to be true got %v", e.done)
	}
}

func BenchmarkWaitForFullJitter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		e := NewBackoff(20)
		e.strategy = FullJitter
		for e.Do() {
			e.WaitFor()
		}
	}
}
func BenchmarkWaitForDecorrelatedJitter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		e := NewBackoff(20)
		e.strategy = DecorrelatedJitter
		for e.Do() {
			e.WaitFor()
		}
	}
}
func BenchmarkWaitForExponentialBackoff(b *testing.B) {
	for i := 0; i < b.N; i++ {
		e := NewBackoff(20)
		e.strategy = ExponentialBackoff
		for e.Do() {
			e.WaitFor()
		}
	}
}
func BenchmarkWaitForLinearBackoff(b *testing.B) {
	for i := 0; i < b.N; i++ {
		e := NewBackoff(20)
		e.strategy = ExponentialBackoff
		for e.Do() {
			e.WaitFor()
		}
	}
}

package exponent

import (
	"testing"
	"time"
)

func TestExponentialBackoffRetries(t *testing.T) {
	e := NewBackoff(13)
	e.strategy = LinearBackoff
	var n int
	t.Log("starting")
	for e.Do() {
		n++
		t.Logf("Sleeping for %v for run %d", e.WaitFor(), e.N)
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

func TestFailed(t *testing.T) {
	e := NewBackoff(2)
	if e.Failed() == true {
		t.Errorf("Expected failed to be false got true")
	}
	e.Do()
	e.Do()
	if e.Failed() == false {
		t.Errorf("Expected failed to be true after 2 Do but got false")
	}
	e.Success()
	if e.Failed() == true {
		t.Errorf("Expected failed to be false after 2 Do but call to success - got true")
	}
}

func TestWithMinimum(t *testing.T) {
	higher := func(e Exp) time.Duration { return 300 * time.Millisecond }
	lower := func(e Exp) time.Duration { return 1 * time.Millisecond }

	if delay := WithMinimum(higher, 50*time.Millisecond)(Exp{}); delay != 300*time.Millisecond {
		t.Errorf("Expected WithMinimum to use max of 300ms got: %v", delay)
	}
	if delay := WithMinimum(lower, 50*time.Millisecond)(Exp{}); delay != 50*time.Millisecond {
		t.Errorf("Expected WithMinimum to use min of 50 got: %v", delay)
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

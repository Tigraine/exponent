package exponent

import (
	"context"
	"errors"
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

func TestNewApi(t *testing.T) {
	e := NewExponent(context.TODO(), 10)
	result, err := e.Try(func() (any, error) {
		return "foo", nil
	})
	if err != nil {
		t.Errorf("expected nil got %v", err)
	}
	if result != "foo" {
		t.Errorf("expected result 'foo' got '%s'", result)
	}
}
func TestNewApiRetry(t *testing.T) {
	e := NewExponent(context.TODO(), 10).Strategy(DecorrelatedJitter)
	attempts := 0
	result, err := e.Try(func() (any, error) {
		attempts++
		if attempts < 4 {
			return nil, errors.New("fail")
		}
		return "foo", nil
	})
	if err != nil {
		t.Errorf("expected nil got %v", err)
	}
	if result != "foo" {
		t.Errorf("expected result 'foo' got '%s'", result)
	}
}

func TestNewApiMaxRetry(t *testing.T) {
	e := NewExponent(context.TODO(), 1).Strategy(DecorrelatedJitter)
	_, err := e.Try(func() (any, error) {
		return nil, errors.New("foo")
	})
	if err.Error() != errors.New("foo").Error() {
		t.Errorf("expected error foo got %v", err)
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

func TestWithContextAndTimeout(t *testing.T) {
	deadline, cancelFunc := context.WithTimeout(context.TODO(), 200*time.Millisecond)
	defer cancelFunc()
	e := NewBackoff(10).WithContext(deadline)
	for e.Do() {
		// do nothing - we keep looping
		e.Wait()
	}
	if !e.Failed() {
		t.Errorf("expected to fail due to timeout")
	}
}

func TestWithContextSuccess(t *testing.T) {
	deadline, cancelFunc := context.WithTimeout(context.TODO(), 200*time.Millisecond)
	defer cancelFunc()
	e := NewBackoff(10).WithContext(deadline)
	for e.Do() {
		// do nothing - we keep looping
		e.Wait()
		e.Success()
	}
	if e.Failed() {
		t.Errorf("expected to not fail due to timeout")
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
	higher := func(attempt int) time.Duration { return 300 * time.Millisecond }
	lower := func(attempt int) time.Duration { return 1 * time.Millisecond }

	if delay := WithMinimum(higher, 50*time.Millisecond)(0); delay != 300*time.Millisecond {
		t.Errorf("Expected WithMinimum to use max of 300ms got: %v", delay)
	}
	if delay := WithMinimum(lower, 50*time.Millisecond)(0); delay != 50*time.Millisecond {
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
		e.strategy = LinearBackoff
		for e.Do() {
			e.WaitFor()
		}
	}
}

func BenchmarkNoWait(b *testing.B) {
	for i := 0; i < b.N; i++ {
		e := NewBackoff(20)
		e.SetStrategy(ExponentialBackoff)
		for e.Do() {
			e.Success()
		}
	}
}

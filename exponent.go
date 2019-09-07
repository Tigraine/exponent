package exponent

import (
	"math"
	"math/rand"
	"time"
)

type strategy func(e *exp) time.Duration

var DecorrelatedJitter = func(e *exp) time.Duration {
	ex := math.Exp2(float64(e.n))
	jitter := rand.Int63n(int64(ex) / 2)
	return time.Duration(jitter+int64(ex/2)) * time.Millisecond
}

var FullJitter = func(e *exp) time.Duration {
	ex := math.Exp2(float64(e.n))
	jitter := rand.Int63n(int64(ex))
	return time.Duration(jitter) * time.Millisecond
}

var ExponentialBackoff = func(e *exp) time.Duration {
	return time.Duration(math.Exp2(float64(e.n))) * time.Millisecond
}

var LinearBackoff = func(e *exp) time.Duration {
	return time.Duration(e.n*100) * time.Millisecond
}

type exp struct {
	n        int
	retries  int
	max      time.Duration
	strategy strategy
	done     bool
}

func (e *exp) Do() bool {
	e.n++
	return !e.done && e.n <= e.retries
}

func (e *exp) WaitFor() time.Duration {
	return e.strategy(e)
}

func (e *exp) Wait() time.Duration {
	sleep := e.WaitFor()
	time.Sleep(sleep)
	return sleep
}

func (e *exp) Success() {
	e.done = true
}

func NewExponent(retries int, max time.Duration) *exp {
	return &exp{
		retries:  retries,
		max:      max,
		strategy: DecorrelatedJitter,
	}
}

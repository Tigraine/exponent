package exponent

import (
	"math"
	"math/rand"
	"time"
)

var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

type strategy func(e *exp) time.Duration

//https://aws.amazon.com/de/blogs/architecture/exponential-backoff-and-jitter/

// DecorrelatedJitter
var DecorrelatedJitter = func(e *exp) time.Duration {
	ex := math.Exp2(float64(e.n))
	jitter := rnd.Int63n(int64(ex) / 2)
	return time.Duration(jitter+int64(ex/2)) * time.Millisecond
}

var FullJitter = func(e *exp) time.Duration {
	ex := math.Exp2(float64(e.n))
	jitter := rnd.Int63n(int64(ex))
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

func (e *exp) SetStrategy(strat strategy) {
	e.strategy = strat
}

func (e *exp) Failed() bool {
	return e.n >= e.retries && !e.done
}

func NewBackoff(retries int) *exp {
	return &exp{
		retries:  retries,
		strategy: DecorrelatedJitter,
	}
}

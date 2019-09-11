package exponent

import (
	"math"
	"math/rand"
	"time"
)

var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

type Strategy func(e Exp) time.Duration

//https://aws.amazon.com/de/blogs/architecture/exponential-backoff-and-jitter/

// DecorrelatedJitter
var DecorrelatedJitter = func(e Exp) time.Duration {
	ex := math.Exp2(float64(e.N))
	jitter := rnd.Int63n(int64(ex) / 2)
	return time.Duration(jitter+int64(ex/2)) * time.Millisecond
}

var FullJitter = func(e Exp) time.Duration {
	ex := math.Exp2(float64(e.N))
	jitter := rnd.Int63n(int64(ex))
	return time.Duration(jitter) * time.Millisecond
}

var ExponentialBackoff = func(e Exp) time.Duration {
	return time.Duration(math.Exp2(float64(e.N))) * time.Millisecond
}

var LinearBackoff = func(e Exp) time.Duration {
	return time.Duration(e.N*100) * time.Millisecond
}

type Exp struct {
	N        int
	retries  int
	strategy Strategy
	done     bool
}

func (e *Exp) Do() bool {
	e.N++
	return !e.done && e.N <= e.retries
}

func (e *Exp) WaitFor() time.Duration {
	return e.strategy(*e)
}

func (e *Exp) Wait() time.Duration {
	sleep := e.WaitFor()
	time.Sleep(sleep)
	return sleep
}

func (e *Exp) Success() {
	e.done = true
}

func (e *Exp) SetStrategy(strat Strategy) {
	e.strategy = strat
}

func (e *Exp) Failed() bool {
	return e.N >= e.retries && !e.done
}

func NewBackoff(retries int) *Exp {
	return &Exp{
		retries:  retries,
		strategy: DecorrelatedJitter,
	}
}

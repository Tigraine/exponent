package exponent

import (
	"context"
	"math"
	"math/rand"
	"time"
)

var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

type Strategy func(attempt int) time.Duration

//https://aws.amazon.com/de/blogs/architecture/exponential-backoff-and-jitter/

// DecorrelatedJitter
var DecorrelatedJitter = func(attempt int) time.Duration {
	ex := math.Exp2(float64(attempt))
	jitter := rnd.Int63n(int64(ex) / 2)
	return time.Duration(jitter+int64(ex/2)) * time.Millisecond
}

var FullJitter = func(attempt int) time.Duration {
	ex := math.Exp2(float64(attempt))
	jitter := rnd.Int63n(int64(ex))
	return time.Duration(jitter) * time.Millisecond
}

var ExponentialBackoff = func(attempt int) time.Duration {
	return time.Duration(math.Exp2(float64(attempt))) * time.Millisecond
}

var LinearBackoff = func(attempt int) time.Duration {
	return time.Duration(attempt*100) * time.Millisecond
}

var WithMinimum = func(strat Strategy, min time.Duration) func(attempt int) time.Duration {
	return func(attempt int) time.Duration {
		x := strat(attempt)
		if x < min {
			return min
		}
		return x
	}
}

type Exp struct {
	N        int
	retries  int
	strategy Strategy
	done     bool
	ctx      context.Context
}

func (e *Exp) Do() bool {
	e.N++
	return !e.done && e.N <= e.retries && e.ctx.Err() == nil
}

func (e *Exp) WaitFor() time.Duration {
	return e.strategy(e.N)
}

func (e *Exp) Wait() time.Duration {
	sleep := e.WaitFor()
	select {
	case <-time.After(sleep):
		return sleep
	case <-e.ctx.Done():
		return sleep
	}
}

func (e *Exp) Success() {
	e.done = true
}

func (e *Exp) SetStrategy(strat Strategy) {
	e.strategy = strat
}

func (e *Exp) Failed() bool {
	if e.ctx.Err() != nil {
		return !e.done
	}
	return e.N >= e.retries && !e.done
}

func (e *Exp) WithContext(ctx context.Context) *Exp {
	e.ctx = ctx
	return e
}

const DEFAULT_MIN_DELAY = 150 * time.Millisecond

func NewBackoff(retries int) *Exp {
	return &Exp{
		retries:  retries,
		strategy: WithMinimum(DecorrelatedJitter, DEFAULT_MIN_DELAY),
		ctx:      context.TODO(),
	}
}

package exponent

import (
	"math"
	"math/rand"
	"time"
)

type exp struct {
	n       int
	retries int
	max     time.Duration
}

func (e *exp) Do() bool {
	e.n++
	return e.n <= e.retries
}

func (e *exp) WaitFor() time.Duration {
	ex := math.Exp2(float64(e.n))
	jitter := rand.Int63n(int64(ex) / 2)
	return time.Duration(jitter+int64(ex/2)) * time.Millisecond
}

func (e *exp) Wait() time.Duration {
	sleep := e.WaitFor()
	time.Sleep(sleep)
	return sleep
}

func NewExponent(retries int, max time.Duration) *exp {
	return &exp{retries: retries, max: max}
}

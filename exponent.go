package exponent

import (
	"context"
)

type Exponent struct {
	exp *Exp
}

func NewExponent(ctx context.Context, maxAttempts int) *Exponent {
	return &Exponent{
		exp: NewBackoff(maxAttempts).WithContext(ctx),
	}
}

func (e *Exponent) Try(fn func() (any, error)) (any, error) {
	var err error
	for e.exp.Do() {
		var result any
		result, err = fn()
		if err != nil {
			e.exp.Wait()
			continue
		}
		e.exp.Success()
		return result, nil
	}
	return nil, err
}

func (e *Exponent) Strategy(s Strategy) *Exponent {
	e.exp.SetStrategy(s)
	return e
}

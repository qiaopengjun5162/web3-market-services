package retry

import (
	"context"
	"fmt"
	"time"
)

type ErrFailedPermanently struct {
	attempts int
	LastErr  error
}

func (e *ErrFailedPermanently) Error() string {
	return fmt.Sprintf("operation failed permanently after %d attempts: %v", e.attempts, e.LastErr)
}

func (e *ErrFailedPermanently) Unwrap() error {
	return e.LastErr
}

type pair[T, U any] struct {
	a T
	b U
}

func Do2[T, U any](ctx context.Context, maxAttempts int, strategy Strategy, op func() (T, U, error)) (T, U, error) {
	f := func() (pair[T, U], error) {
		a, b, err := op()
		return pair[T, U]{a, b}, err
	}
	res, err := Do(ctx, maxAttempts, strategy, f)
	return res.a, res.b, err
}

func Do[T any](ctx context.Context, maxAttempts int, strategy Strategy, op func() (T, error)) (T, error) {
	var empty, ret T
	var err error
	if maxAttempts < 1 {
		return empty, fmt.Errorf("need at least one attempt to run op, but got %d", maxAttempts)
	}

	for i := 0; i < maxAttempts; i++ {
		if ctx.Err() != nil {
			return empty, ctx.Err()
		}
		ret, err = op()
		if err == nil {
			return ret, nil
		}
		if i != maxAttempts-1 {
			time.Sleep(strategy.Duration(i))
		}
	}
	return empty, &ErrFailedPermanently{attempts: maxAttempts, LastErr: err}
}

package retry

import (
	"context"
	"errors"
	"time"
)

var ErrRetryable = errors.New("retryable error")

type Operation func(context.Context) error

type Config struct {
	MaxRetries int
	BaseDelay  time.Duration
	MaxDelay   time.Duration
}

func Do(ctx context.Context, cfg Config, op Operation) error {
	var err error

	for attempt := 0; attempt <= cfg.MaxRetries; attempt++ {
		err = op(ctx)

		// успешно 200, ничего не делаем
		if err == nil {
			return nil
		}

		if !errors.Is(err, ErrRetryable) {
			return err
		}

		if attempt == cfg.MaxRetries {
			return err
		}

		delay := CalculateDelay(attempt, cfg.BaseDelay, cfg.MaxDelay)
		timer := time.NewTimer(delay)

		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()

		case <-timer.C:
		}
	}

	return err
}

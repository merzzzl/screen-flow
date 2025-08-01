package actions

import (
	"context"
	"time"

	"github.com/merzzzl/screen-flow/device"
)

type ActionWait struct {
	Duration time.Duration
}

func (s *ActionWait) Handle(ctx context.Context, _ *device.Conn) error {
	t := time.NewTimer(s.Duration)

	select {
	case <-t.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

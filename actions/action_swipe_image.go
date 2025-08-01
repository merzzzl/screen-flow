package actions

import (
	"context"
	"fmt"
	"image"
	"time"

	"github.com/merzzzl/screen-flow/device"
)

type ActionSwipeImage struct {
	ImageTemplate image.Image
	H             int
	W             int
	Duration      time.Duration
	SearchArea    *image.Rectangle
}

func (s *ActionSwipeImage) Handle(ctx context.Context, conn *device.Conn) error {
	if err := conn.CheckVision(); err != nil {
		return fmt.Errorf("need vision: %w, %w", ErrNoClints, err)
	}

	point, err := conn.GetVision().Find(ctx, s.ImageTemplate)
	if err != nil {
		return fmt.Errorf("find point: %w", err)
	}

	if s.SearchArea == nil {
		nextStep := ActionSwipe{
			X1:       point.X,
			Y1:       point.Y,
			X2:       point.X + s.H,
			Y2:       point.Y + s.W,
			Duration: s.Duration,
		}

		return nextStep.Handle(ctx, conn)
	}

	if point.In(*s.SearchArea) {
		nextStep := ActionSwipe{
			X1:       point.X,
			Y1:       point.Y,
			X2:       point.X + s.H,
			Y2:       point.Y + s.W,
			Duration: s.Duration,
		}

		return nextStep.Handle(ctx, conn)
	}

	return ErrImageNotFound
}

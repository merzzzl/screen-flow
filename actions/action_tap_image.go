package actions

import (
	"context"
	"fmt"
	"image"
	"time"

	"github.com/merzzzl/screen-flow/device"
)

type ActionTapImage struct {
	ImageTemplate image.Image
	Duration      time.Duration
	SearchArea    *image.Rectangle
}

func (s *ActionTapImage) Handle(ctx context.Context, conn *device.Conn) error {
	if err := conn.CheckVision(); err != nil {
		return fmt.Errorf("need vision: %w, %w", ErrNoClints, err)
	}

	point, err := conn.GetVision().Find(ctx, s.ImageTemplate)
	if err != nil {
		return fmt.Errorf("find point: %w", err)
	}

	if s.SearchArea == nil {
		nextStep := ActionTap{
			X: point.X,
			Y: point.Y,
		}

		return nextStep.Handle(ctx, conn)
	}

	if point.In(*s.SearchArea) {
		nextStep := ActionTap{
			X:        point.X,
			Y:        point.Y,
			Duration: s.Duration,
		}

		return nextStep.Handle(ctx, conn)
	}

	return ErrImageNotFound
}

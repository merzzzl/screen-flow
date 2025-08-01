package actions

import (
	"context"
	"fmt"
	"image"
	"time"

	"github.com/merzzzl/screen-flow/device"
)

type ActionWaitImage struct {
	ImageTemplate image.Image
	Duration      *time.Duration
	SearchArea    *image.Rectangle
}

func (s *ActionWaitImage) Handle(ctx context.Context, conn *device.Conn) error {
	if err := conn.CheckVision(); err != nil {
		return fmt.Errorf("need vision: %w, %w", ErrNoClints, err)
	}

	startAt := time.Now()

	for {
		point, err := conn.GetVision().Find(ctx, s.ImageTemplate)
		if err != nil {
			return fmt.Errorf("find point: %w", err)
		}

		if s.SearchArea == nil {
			return nil
		}

		if point.In(*s.SearchArea) {
			return nil
		}

		time.Sleep(time.Microsecond * 200)

		if s.Duration != nil {
			if time.Now().After(startAt.Add(*s.Duration)) {
				return fmt.Errorf("find point: %w", ErrImageNotFound)
			}
		}
	}
}

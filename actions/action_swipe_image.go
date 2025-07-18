package actions

import (
	"fmt"
	"image"
	"time"

	"github.com/merzzzl/screen-flow/device"
	"github.com/merzzzl/screen-flow/events"
	"github.com/merzzzl/screen-flow/vision"
)

type StepActionSwipeImage struct {
	ImageTemplate image.Image
	H             int
	W             int
	Duration      *time.Duration
	SearchArea    *image.Rectangle
	Wait          bool
}

func (s *StepActionSwipeImage) Handle(conn *device.Conn) error {
	startAt := time.Now()

	for {
		source := conn.GetScreenImage()

		point, ok, err := vision.FindPoint(source, s.ImageTemplate)
		if err != nil {
			return fmt.Errorf("find point: %w", err)
		}

		if ok {
			if s.SearchArea == nil {
				conn.BroadcastEvent(events.NewFoundImageEvent(point.X, point.Y))

				nextStep := ActionSwipe{
					X1: point.X,
					Y1: point.Y,
					X2: point.X + s.H,
					Y2: point.Y + s.W,
				}

				return nextStep.Handle(conn)
			}

			if point.In(*s.SearchArea) {
				conn.BroadcastEvent(events.NewFoundImageEvent(point.X, point.Y))

				nextStep := ActionSwipe{
					X1: point.X,
					Y1: point.Y,
					X2: point.X + s.H,
					Y2: point.Y + s.W,
				}

				return nextStep.Handle(conn)
			}
		}

		if s.Duration != nil {
			if time.Now().After(startAt.Add(*s.Duration)) {
				return fmt.Errorf("find point: %w", err)
			}
		}

		if !s.Wait {
			return fmt.Errorf("find point: %w", ErrImageNotFound)
		}
	}
}

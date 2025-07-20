package actions

import (
	"fmt"
	"image"
	"time"

	"github.com/merzzzl/screen-flow/device"
)

type ActionTapImage struct {
	ImageTemplate image.Image
	Duration      *time.Duration
	SearchArea    *image.Rectangle
	Wait          bool
}

func (s *ActionTapImage) Handle(conn *device.Conn) error {
	startAt := time.Now()

	for {
		point, err := conn.FindPoint(s.ImageTemplate)
		if err != nil {
			return fmt.Errorf("find point: %w", err)
		}

		if s.SearchArea == nil {
			nextStep := ActionTapXY{
				X: point.X,
				Y: point.Y,
			}

			return nextStep.Handle(conn)
		}

		if point.In(*s.SearchArea) {
			nextStep := ActionTapXY{
				X: point.X,
				Y: point.Y,
			}

			return nextStep.Handle(conn)
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

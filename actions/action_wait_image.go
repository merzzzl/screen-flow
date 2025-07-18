package actions

import (
	"errors"
	"fmt"
	"image"
	"time"

	"github.com/merzzzl/screen-flow/device"
	"github.com/merzzzl/screen-flow/events"
	"github.com/merzzzl/screen-flow/vision"
)

var ErrImageNotFound = errors.New("image not found")

type ActionWait struct {
	ImageTemplate image.Image
	Duration      *time.Duration
	SearchArea    *image.Rectangle
}

func (s *ActionWait) Handle(conn *device.Conn) error {
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

				return nil
			}

			if point.In(*s.SearchArea) {
				conn.BroadcastEvent(events.NewFoundImageEvent(point.X, point.Y))

				return nil
			}
		}

		time.Sleep(time.Microsecond * 200)

		if s.Duration != nil {
			if time.Now().After(startAt.Add(*s.Duration)) {
				return fmt.Errorf("find point: %w", ErrImageNotFound)
			}
		}
	}
}

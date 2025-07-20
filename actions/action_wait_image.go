package actions

import (
	"errors"
	"fmt"
	"image"
	"time"

	"github.com/merzzzl/screen-flow/device"
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
		point, err := conn.FindPoint(s.ImageTemplate)
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

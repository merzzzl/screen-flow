package actions

import (
	"fmt"
	"time"

	scrcpy "github.com/merzzzl/scrcpy-go"
	"github.com/merzzzl/screen-flow/device"
	"github.com/merzzzl/screen-flow/events"
)

type ActionSwipe struct {
	X1       int
	Y1       int
	X2       int
	Y2       int
	Duration time.Duration
}

func (s *ActionSwipe) Handle(conn *device.Conn) error {
	if err := conn.InjectTouch(
		scrcpy.ActionDown, 1,
		uint32(s.X1), uint32(s.Y1),
		65535, scrcpy.ButtonPrimary, scrcpy.ButtonPrimary); err != nil {
		return fmt.Errorf("inject action: %w", err)
	}

	conn.BroadcastEvent(events.NewTouchEvent(s.X1, s.Y1))

	if s.Duration == 0 {
		s.Duration = 200 * time.Millisecond
	}

	steps := s.Duration / 10
	currentX := s.X1
	currentY := s.Y1
	moveX := (s.X2 - s.X1) / int(steps)
	moveY := (s.Y2 - s.Y1) / int(steps)

	for range steps {
		currentX += moveX
		currentY += moveY

		if err := conn.InjectTouch(
			scrcpy.ActionMove, 1,
			uint32(currentX), uint32(currentY),
			65535, 0, scrcpy.ButtonPrimary); err != nil {
			return fmt.Errorf("inject action: %w", err)
		}

		conn.BroadcastEvent(events.NewTouchEvent(currentX, currentY))
	}

	if err := conn.InjectTouch(
		scrcpy.ActionUp, 1,
		uint32(s.X2), uint32(s.Y2),
		65535, scrcpy.ButtonPrimary, 0); err != nil {
		return fmt.Errorf("inject action: %w", err)
	}

	conn.BroadcastEvent(events.NewTouchEvent(s.X2, s.Y2))

	return nil
}

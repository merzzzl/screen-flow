package actions

import (
	"context"
	"fmt"
	"time"

	abg "github.com/merzzzl/accessibility-bridge-go"
	scrcpy "github.com/merzzzl/scrcpy-go"
	"github.com/merzzzl/screen-flow/device"
)

type ActionSwipe struct {
	X1       int
	Y1       int
	X2       int
	Y2       int
	Duration time.Duration
}

func (s *ActionSwipe) Handle(ctx context.Context, conn *device.Conn) error {
	if conn.CheckABG() == nil {
		return s.abg(ctx, conn)
	}

	if conn.CheckSCRCPY() == nil {
		return s.scrcpy(conn)
	}

	return fmt.Errorf("need accessibility-bridge or scrcpy: %w", ErrNoClints)
}

func (s *ActionSwipe) scrcpy(conn *device.Conn) error {
	if err := conn.GetSCRCPY().InjectTouch(
		scrcpy.ActionDown, 1,
		uint32(s.X1), uint32(s.Y1),
		65535, scrcpy.ButtonPrimary, scrcpy.ButtonPrimary); err != nil {
		return fmt.Errorf("inject action: %w", err)
	}

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

		if err := conn.GetSCRCPY().InjectTouch(
			scrcpy.ActionMove, 1,
			uint32(currentX), uint32(currentY),
			65535, 0, scrcpy.ButtonPrimary); err != nil {
			return fmt.Errorf("inject action: %w", err)
		}
	}

	if err := conn.GetSCRCPY().InjectTouch(
		scrcpy.ActionUp, 1,
		uint32(s.X2), uint32(s.Y2),
		65535, scrcpy.ButtonPrimary, 0); err != nil {
		return fmt.Errorf("inject action: %w", err)
	}

	return nil
}

func (s *ActionSwipe) abg(ctx context.Context, conn *device.Conn) error {
	if err := conn.GetABG().PerformSwipe(ctx, &abg.ActionSwipe{
		Finger: &abg.Finger{
			Duration: int32(s.Duration),
			Start: &abg.Finger_StartPoint{
				StartPoint: &abg.Point{
					X: int32(s.X1),
					Y: int32(s.Y1),
				},
			},
			Width:  int32(s.X2) - int32(s.X1),
			Height: int32(s.Y2) - int32(s.Y1),
		},
	}); err != nil {
		return fmt.Errorf("inject action: %w", err)
	}

	return nil
}

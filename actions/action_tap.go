package actions

import (
	"context"
	"fmt"
	"time"

	abg "github.com/merzzzl/accessibility-bridge-go"
	scrcpy "github.com/merzzzl/scrcpy-go"
	"github.com/merzzzl/screen-flow/device"
)

type ActionTap struct {
	X        int
	Y        int
	Duration time.Duration
}

func (s *ActionTap) Handle(ctx context.Context, conn *device.Conn) error {
	if conn.CheckABG() == nil {
		return s.abg(ctx, conn)
	}

	if conn.CheckSCRCPY() == nil {
		return s.scrcpy(conn)
	}

	return fmt.Errorf("need accessibility-bridge or scrcpy: %w", ErrNoClints)
}

func (s *ActionTap) scrcpy(conn *device.Conn) error {
	if err := conn.GetSCRCPY().InjectTouch(
		scrcpy.ActionDown, 1,
		uint32(s.X), uint32(s.Y),
		65535, scrcpy.ButtonPrimary, scrcpy.ButtonPrimary); err != nil {
		return fmt.Errorf("inject action: %w", err)
	}

	if s.Duration == 0 {
		s.Duration = 20 * time.Millisecond
	}

	time.Sleep(s.Duration)

	if err := conn.GetSCRCPY().InjectTouch(
		scrcpy.ActionUp, 1,
		uint32(s.X), uint32(s.Y),
		65535, scrcpy.ButtonPrimary, 0); err != nil {
		return fmt.Errorf("inject action: %w", err)
	}

	return nil
}

func (s *ActionTap) abg(ctx context.Context, conn *device.Conn) error {
	if err := conn.GetABG().PerformClick(ctx, &abg.ActionClick{
		Duration: int32(s.Duration),
		Click: &abg.ActionClick_ClickPoint{
			ClickPoint: &abg.Point{
				X: int32(s.X),
				Y: int32(s.Y),
			},
		},
	}); err != nil {
		return fmt.Errorf("inject action: %w", err)
	}

	return nil
}
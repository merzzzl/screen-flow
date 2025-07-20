package actions

import (
	"fmt"
	"time"

	scrcpy "github.com/merzzzl/scrcpy-go"
	"github.com/merzzzl/screen-flow/device"
)

type ActionTapXY struct {
	X        int
	Y        int
	Duration time.Duration
}

func (s *ActionTapXY) Handle(conn *device.Conn) error {
	if err := conn.InjectTouch(
		scrcpy.ActionDown, 1,
		uint32(s.X), uint32(s.Y),
		65535, scrcpy.ButtonPrimary, scrcpy.ButtonPrimary); err != nil {
		return fmt.Errorf("inject action: %w", err)
	}

	if s.Duration == 0 {
		s.Duration = 20 * time.Millisecond
	}

	time.Sleep(s.Duration)

	if err := conn.InjectTouch(
		scrcpy.ActionUp, 1,
		uint32(s.X), uint32(s.Y),
		65535, scrcpy.ButtonPrimary, 0); err != nil {
		return fmt.Errorf("inject action: %w", err)
	}

	return nil
}

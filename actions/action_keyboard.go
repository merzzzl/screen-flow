package actions

import (
	"fmt"
	"time"

	scrcpy "github.com/merzzzl/scrcpy-go"
	"github.com/merzzzl/screen-flow/device"
)

type ActionKeyboard struct {
	Press    []int
	Duration time.Duration
}

func (s *ActionKeyboard) Handle(conn *device.Conn) error {
	for _, key := range s.Press {
		if err := conn.InjectKeycode(uint32(key), scrcpy.ActionDown, 0, 0); err != nil {
			return fmt.Errorf("inject action: %w", err)
		}

		time.Sleep(s.Duration)

		for _, key := range s.Press {
			if err := conn.InjectKeycode(uint32(key), scrcpy.ActionUp, 0, 0); err != nil {
				return fmt.Errorf("inject action: %w", err)
			}
		}
	}

	return nil
}

package actions

import (
	"fmt"

	"github.com/merzzzl/screen-flow/device"
)

type ActionSetClipboard struct {
	Payload string
	Past    bool
}

func (s *ActionSetClipboard) Handle(conn *device.Conn) error {
	if err := conn.SetClipboard(0, s.Payload, s.Past); err != nil {
		return fmt.Errorf("inject action: %w", err)
	}

	return nil
}

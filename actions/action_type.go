package actions

import (
	"fmt"

	"github.com/merzzzl/screen-flow/device"
	"github.com/merzzzl/screen-flow/events"
)

type ActionType struct {
	Payload string
}

func (s *ActionType) Handle(conn *device.Conn) error {
	if err := conn.InjectText(s.Payload); err != nil {
		return fmt.Errorf("inject action: %w", err)
	}

	conn.BroadcastEvent(events.NewTypeTextEvent(s.Payload))

	return nil
}

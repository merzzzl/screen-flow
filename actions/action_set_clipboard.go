package actions

import (
	"github.com/merzzzl/screen-flow/device"
	"github.com/merzzzl/screen-flow/events"
)

type ActionSetClipboard struct {
	Payload string
	Past    bool
}

func (s *ActionSetClipboard) Handle(conn *device.Conn) error {
	if err := conn.SetClipboard(0, s.Payload, s.Past); err != nil {
		return err
	}

	if s.Past {
		conn.BroadcastEvent(events.NewPastTextEvent(s.Payload))
	}

	return nil
}

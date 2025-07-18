package actions

import (
	"time"

	"github.com/merzzzl/screen-flow/device"
)

type ActionDelay struct {
	Duration time.Duration
}

func (s *ActionDelay) Handle(*device.Conn) error {
	t := time.NewTimer(s.Duration)

	<-t.C

	return nil
}

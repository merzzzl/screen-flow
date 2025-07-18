package actions

import (
	"fmt"
	"time"

	"github.com/merzzzl/screen-flow/device"
	"github.com/merzzzl/screen-flow/vision"
)

type ActionWaitStaticFrame struct {
	Threshold float64
}

func (s *ActionWaitStaticFrame) Handle(conn *device.Conn) error {
	first := conn.GetScreenImage()

	for {
		time.Sleep(time.Millisecond * 250)

		second := conn.GetScreenImage()

		ok, err := vision.IsFrameStatic(first, second, s.Threshold)
		if err != nil {
			return fmt.Errorf("compare frames: %w", err)
		}

		if ok {
			break
		}

		first = second
	}

	return nil
}

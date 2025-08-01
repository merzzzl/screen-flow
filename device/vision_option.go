package device

import (
	"context"
	"fmt"

	"github.com/merzzzl/screen-flow/vision"
)

type OptionVision struct {
	algo vision.Algorithm
}

func WithVision(algo vision.Algorithm) *OptionVision {
	return &OptionVision{
		algo: algo,
	}
}

func (o *OptionVision) apply(_ context.Context, conn *Conn) error {
	handshake := conn.scrcpy.GetHandshake()

	if handshake.Height == 0 || handshake.Width == 0 {
		return fmt.Errorf("need scrcpy: %w", ErrNoSCRCPY)
	}

	conn.vision = vision.NewPipe(
		conn.decoder,
		int(handshake.Width),
		int(handshake.Height),
		o.algo,
	)

	return nil
}

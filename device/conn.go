package device

import (
	"context"
	"time"

	abg "github.com/merzzzl/accessibility-bridge-go"
	scrcpy "github.com/merzzzl/scrcpy-go"
	"github.com/merzzzl/screen-flow/vision"
)

type Option interface {
	apply(ctx context.Context, conn *Conn) error
}

type Conn struct {
	abg       abg.ActionManagerClient
	scrcpy    *scrcpy.Client
	decoder   *scrcpy.FFmpeg
	clipboard chan string
	vision    *vision.Pipe
}

func Connect(ctx context.Context, options ...Option) (*Conn, error) {
	ctx, cancel := context.WithCancel(ctx)
	conn := &Conn{}

	for _, op := range options {
		if err := op.apply(ctx, conn); err != nil {
			cancel()

			return nil, err
		}
	}

	go func() {
		if conn.scrcpy != nil {
			_ = conn.scrcpy.Serve(ctx)

			cancel()
		}
	}()

	go func() {
		if conn.vision != nil {
			_ = conn.vision.Process(ctx)

			cancel()
		}
	}()

	t := time.NewTimer(time.Second)

	select {
	case <-t.C:
		return conn, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (c *Conn) CheckABG() error {
	if c != nil && c.abg != nil {
		return nil
	}

	return ErrNoABG
}

func (c *Conn) CheckSCRCPY() error {
	if c != nil && c.scrcpy != nil {
		return nil
	}

	return ErrNoSCRCPY
}

func (c *Conn) CheckVision() error {
	if c != nil && c.vision != nil {
		return nil
	}

	return ErrNoVision
}

func (c *Conn) GetABG() *ABG {
	return &ABG{
		conn: c,
	}
}

func (c *Conn) GetSCRCPY() *SCRCPY {
	return &SCRCPY{
		conn: c,
	}
}

func (c *Conn) GetVision() *Vision {
	return &Vision{
		conn: c,
	}
}

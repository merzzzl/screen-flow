package device

import (
	"context"
	"fmt"
	"image"
	"time"

	scrcpy "github.com/merzzzl/scrcpy-go"
	"github.com/merzzzl/screen-flow/vision"
)

type Conn struct {
	*scrcpy.Client

	decoder   *scrcpy.FFmpeg
	clipboard chan string
	vision    *vision.Pipe
}

func Connect(ctx context.Context, addr string, window vision.Window) (*Conn, error) {
	conn := &Conn{
		clipboard: make(chan string, 1),
	}

	client, err := scrcpy.Dial(ctx, addr)
	if err != nil {
		return nil, fmt.Errorf("init scrcpy: %w", err)
	}

	dec, err := scrcpy.NewDecoder(ctx)
	if err != nil {
		return nil, fmt.Errorf("init decoder: %w", err)
	}

	ctx, cancel := context.WithCancel(ctx)

	client.SetVideoHandler(dec.VideoHandler)

	client.SetControlHandler(func(_ context.Context, cm scrcpy.ControlMessage) error {
		if cm.Type == 0 {
			for len(conn.clipboard) > 0 {
				<-conn.clipboard
			}

			select {
			case conn.clipboard <- string(cm.Payload):
			default:
			}
		}

		return nil
	})

	handshake := client.GetHandshake()
	vpipe := vision.NewPipe(dec, int(handshake.Width), int(handshake.Height), vision.AlgorithmTM, window)

	go func() {
		_ = client.Serve(ctx)

		cancel()
	}()

	go func() {
		_ = vpipe.Process(ctx)

		cancel()
	}()

	conn.Client = client
	conn.decoder = dec
	conn.vision = vpipe

	t := time.NewTimer(time.Second)

	select {
	case <-t.C:
		return conn, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (c *Conn) WaitClipboard() string {
	return <-c.clipboard
}

func (c *Conn) FindPoint(img image.Image) (image.Point, error) {
	pt, err := c.vision.Found(img)
	if err != nil {
		return pt, fmt.Errorf("vision: %w", err)
	}

	return pt, nil
}

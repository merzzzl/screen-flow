package device

import (
	"context"
	"fmt"
	"image"
	"time"

	scrcpy "github.com/merzzzl/scrcpy-go"
	"github.com/merzzzl/screen-flow/events"
	"github.com/merzzzl/screen-flow/utils"
)

type Conn struct {
	*scrcpy.Client

	decoder   *Decoder
	clipboard chan string
	screenCh  <-chan image.Image
	events    *utils.Broadcaster[*events.Base]
}

func Connect(ctx context.Context, addr string) (*Conn, error) {
	conn := &Conn{
		clipboard: make(chan string, 1),
		events:    utils.NewBroadcaster[*events.Base](),
	}

	client, err := scrcpy.Dial(ctx, addr)
	if err != nil {
		return nil, fmt.Errorf("init scrcpy: %w", err)
	}

	dec, err := NewDecoder()
	if err != nil {
		return nil, fmt.Errorf("init decoder: %w", err)
	}

	ctx, cancel := context.WithCancel(ctx)

	client.SetVideoHandler(func(frame []byte) {
		err := dec.Decode(frame)
		if err != nil {
			cancel()
		}
	})

	client.SetControlHandler(func(cm scrcpy.ControlMessage) {
		if cm.Type == 0 {
			for len(conn.clipboard) > 0 {
				<-conn.clipboard
			}

			select {
			case conn.clipboard <- string(cm.Payload):
			default:
			}
		}
	})

	go func() {
		_ = client.Serve(ctx)

		cancel()
	}()

	go func() {
		_ = dec.Looper(ctx)

		cancel()
	}()

	conn.Client = client
	conn.decoder = dec
	conn.screenCh = dec.Subscribe()

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

func (c *Conn) GetScreenImage() image.Image {
	return <-c.screenCh
}

func (c *Conn) GetScreenStream() <-chan image.Image {
	return c.decoder.Subscribe()
}

func (c *Conn) BroadcastEvent(e *events.Base) {
	c.events.Broadcast(e)
}

func (c *Conn) GetEventStream() <-chan *events.Base {
	return c.events.Subscribe()
}

package device

import (
	"context"
	"fmt"

	scrcpy "github.com/merzzzl/scrcpy-go"
)

type OptionSCRCPY struct {
	addr string
}

func WithSCRCPY(addr string) *OptionSCRCPY {
	return &OptionSCRCPY{
		addr: addr,
	}
}

func (o *OptionSCRCPY) apply(ctx context.Context, conn *Conn) error {
	client, err := scrcpy.Dial(ctx, o.addr)
	if err != nil {
		return fmt.Errorf("init scrcpy: %w", err)
	}

	dec, err := scrcpy.NewDecoder(ctx)
	if err != nil {
		return fmt.Errorf("init decoder: %w", err)
	}

	conn.scrcpy = client
	conn.decoder = dec
	conn.clipboard = make(chan string, 1)

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

	return nil
}

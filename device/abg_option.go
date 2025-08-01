package device

import (
	"context"
	"fmt"

	abg "github.com/merzzzl/accessibility-bridge-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type OptionABG struct {
	addr string
}

func WithABG(addr string) *OptionABG {
	return &OptionABG{
		addr: addr,
	}
}

func (o *OptionABG) apply(ctx context.Context, conn *Conn) error {
	dial, err := grpc.NewClient(o.addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("init accessibility-bridge-go: %w", err)
	}

	go func() {
		<-ctx.Done()

		_ = dial.Close()
	}()

	conn.abg = abg.NewActionManagerClient(dial)

	return nil
}

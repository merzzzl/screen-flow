package device

import (
	"context"
	"fmt"

	abg "github.com/merzzzl/accessibility-bridge-go"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ABG struct {
	conn *Conn
}

func (c *ABG) ScreenDump(ctx context.Context) (*abg.ScreenView, error) {
	if err := c.conn.CheckABG(); err != nil {
		return nil, fmt.Errorf("conn: %w", err)
	}

	out, err := c.conn.abg.ScreenDump(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, fmt.Errorf("abg: %w", err)
	}

	return out, nil
}

func (c *ABG) PerformSwipe(ctx context.Context, action *abg.ActionSwipe) error {
	if err := c.conn.CheckABG(); err != nil {
		return fmt.Errorf("conn: %w", err)
	}

	_, err := c.conn.abg.PerformSwipe(ctx, action)
	if err != nil {
		return fmt.Errorf("abg: %w", err)
	}

	return nil
}

func (c *ABG) PerformMultiTouch(ctx context.Context, fingers []*abg.Finger) error {
	if err := c.conn.CheckABG(); err != nil {
		return fmt.Errorf("conn: %w", err)
	}

	_, err := c.conn.abg.PerformMultiTouch(ctx, &abg.ActionMultiTouch{Finger: fingers})
	if err != nil {
		return fmt.Errorf("abg: %w", err)
	}

	return nil
}

func (c *ABG) PerformClick(ctx context.Context, action *abg.ActionClick) error {
	if err := c.conn.CheckABG(); err != nil {
		return fmt.Errorf("conn: %w", err)
	}

	_, err := c.conn.abg.PerformClick(ctx, action)
	if err != nil {
		return fmt.Errorf("abg: %w", err)
	}

	return nil
}

func (c *ABG) TypeText(ctx context.Context, text string) error {
	if err := c.conn.CheckABG(); err != nil {
		return fmt.Errorf("conn: %w", err)
	}

	_, err := c.conn.abg.TypeText(ctx, &abg.ActionTypeText{Text: text})
	if err != nil {
		return fmt.Errorf("abg: %w", err)
	}

	return nil
}

func (c *ABG) PerformAction(ctx context.Context, key abg.ActionKey_Key) error {
	if err := c.conn.CheckABG(); err != nil {
		return fmt.Errorf("conn: %w", err)
	}

	_, err := c.conn.abg.PerformAction(ctx, &abg.ActionKey{Key: key})
	if err != nil {
		return fmt.Errorf("abg: %w", err)
	}

	return nil
}

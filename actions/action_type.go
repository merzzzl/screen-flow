package actions

import (
	"context"
	"fmt"

	"github.com/merzzzl/screen-flow/device"
)

type ActionType struct {
	Payload string
}

func (s *ActionType) Handle(ctx context.Context, conn *device.Conn) error {
	if conn.CheckABG() == nil {
		return s.abg(ctx, conn)
	}

	if conn.CheckSCRCPY() == nil {
		return s.scrcpy(conn)
	}

	return fmt.Errorf("need accessibility-bridge or scrcpy: %w", ErrNoClints)
}

func (s *ActionType) scrcpy(conn *device.Conn) error {
	if err := conn.GetSCRCPY().InjectText(s.Payload); err != nil {
		return fmt.Errorf("inject action: %w", err)
	}

	return nil
}

func (s *ActionType) abg(ctx context.Context, conn *device.Conn) error {
	if err := conn.GetABG().TypeText(ctx, s.Payload); err != nil {
		return fmt.Errorf("inject action: %w", err)
	}

	return nil
}

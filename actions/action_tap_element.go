package actions

import (
	"context"
	"fmt"
	"time"

	abg "github.com/merzzzl/accessibility-bridge-go"
	"github.com/merzzzl/screen-flow/device"
)

type ActionTapElement struct {
	Regexp   string
	UniqueID string
	Duration time.Duration
}

func (s *ActionTapElement) Handle(ctx context.Context, conn *device.Conn) error {
	if err := conn.CheckABG(); err != nil {
		return fmt.Errorf("need accessibility-bridge: %w, %w", ErrNoClints, err)
	}

	err := conn.GetABG().PerformClick(ctx, &abg.ActionClick{
		Click: &abg.ActionClick_ClickElement{
			ClickElement: &abg.ElementSelector{
				UniqueId: s.UniqueID,
				Regex:    s.Regexp,
			},
		},
		Duration: int32(s.Duration),
	})
	if err != nil {
		return fmt.Errorf("inject action: %w", err)
	}

	return nil
}

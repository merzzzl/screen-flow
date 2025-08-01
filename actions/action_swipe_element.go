package actions

import (
	"context"
	"fmt"
	"time"

	abg "github.com/merzzzl/accessibility-bridge-go"
	"github.com/merzzzl/screen-flow/device"
)

type ActionSwipeElement struct {
	Regexp   string
	UniqueID string
	H        int
	W        int
	Duration time.Duration
}

func (s *ActionSwipeElement) Handle(ctx context.Context, conn *device.Conn) error {
	if err := conn.CheckABG(); err != nil {
		return fmt.Errorf("need accessibility-bridge: %w, %w", ErrNoClints, err)
	}

	err := conn.GetABG().PerformSwipe(ctx, &abg.ActionSwipe{
		Finger: &abg.Finger{
			FingerId: 1,
			Start: &abg.Finger_StartElement{
				StartElement: &abg.ElementSelector{
					UniqueId: s.UniqueID,
					Regex:    s.Regexp,
				},
			},
			Width:    int32(s.W),
			Height:   int32(s.H),
			Duration: int32(s.Duration),
		},
	})
	if err != nil {
		return fmt.Errorf("inject action: %w", err)
	}

	return nil
}

package device

import (
	"context"
	"fmt"
	"image"
)

type Vision struct {
	conn *Conn
}

func (c *Vision) Find(ctx context.Context, img image.Image) (image.Point, error) {
	if err := c.conn.CheckVision(); err != nil {
		return image.Pt(0, 0), fmt.Errorf("conn: %w", err)
	}

	out, err := c.conn.vision.Find(ctx, img)
	if err != nil {
		return image.Pt(0, 0), fmt.Errorf("abg: %w", err)
	}

	return out, nil
}

package device

import (
	"context"
	"fmt"
)

type SCRCPY struct {
	conn *Conn
}

func (c *SCRCPY) StartApp(name string) error {
	if err := c.conn.CheckSCRCPY(); err != nil {
		return fmt.Errorf("conn: %w", err)
	}

	err := c.conn.scrcpy.StartApp(name)
	if err != nil {
		return fmt.Errorf("scrcpy: %w", err)
	}

	return nil
}

func (c *SCRCPY) SetDisplayPower(on bool) error {
	if err := c.conn.CheckSCRCPY(); err != nil {
		return fmt.Errorf("conn: %w", err)
	}

	if err := c.conn.scrcpy.SetDisplayPower(on); err != nil {
		return fmt.Errorf("scrcpy: %w", err)
	}

	return nil
}

func (c *SCRCPY) RotateDevice() error {
	if err := c.conn.CheckSCRCPY(); err != nil {
		return fmt.Errorf("conn: %w", err)
	}

	if err := c.conn.scrcpy.RotateDevice(); err != nil {
		return fmt.Errorf("scrcpy: %w", err)
	}

	return nil
}

func (c *SCRCPY) InjectText(text string) error {
	if err := c.conn.CheckSCRCPY(); err != nil {
		return fmt.Errorf("conn: %w", err)
	}

	if err := c.conn.scrcpy.InjectText(text); err != nil {
		return fmt.Errorf("scrcpy: %w", err)
	}

	return nil
}

func (c *SCRCPY) GetClipboard(ctx context.Context, copyKey byte) (string, error) {
	if err := c.conn.CheckSCRCPY(); err != nil {
		return "", fmt.Errorf("conn: %w", err)
	}

	if err := c.conn.scrcpy.GetClipboard(copyKey); err != nil {
		return "", fmt.Errorf("scrcpy: %w", err)
	}

	str, err := c.waitClipboard(ctx)
	if err != nil {
		return "", err
	}

	return str, nil
}

func (c *SCRCPY) SetClipboard(sequence uint64, text string, paste bool) error {
	if err := c.conn.CheckSCRCPY(); err != nil {
		return fmt.Errorf("conn: %w", err)
	}

	if err := c.conn.scrcpy.SetClipboard(sequence, text, paste); err != nil {
		return fmt.Errorf("scrcpy: %w", err)
	}

	return nil
}

func (c *SCRCPY) UhidCreate(id, vendorID, productID uint16, name string, data []byte) error {
	if err := c.conn.CheckSCRCPY(); err != nil {
		return fmt.Errorf("conn: %w", err)
	}

	if err := c.conn.scrcpy.UhidCreate(id, vendorID, productID, name, data); err != nil {
		return fmt.Errorf("scrcpy: %w", err)
	}

	return nil
}

func (c *SCRCPY) UhidDestroy(id uint16) error {
	if err := c.conn.CheckSCRCPY(); err != nil {
		return fmt.Errorf("conn: %w", err)
	}

	if err := c.conn.scrcpy.UhidDestroy(id); err != nil {
		return fmt.Errorf("scrcpy: %w", err)
	}

	return nil
}

func (c *SCRCPY) UhidInput(id uint16, data []byte) error {
	if err := c.conn.CheckSCRCPY(); err != nil {
		return fmt.Errorf("conn: %w", err)
	}

	if err := c.conn.scrcpy.UhidInput(id, data); err != nil {
		return fmt.Errorf("scrcpy: %w", err)
	}

	return nil
}

func (c *SCRCPY) BackOrScreenOn(action byte) error {
	if err := c.conn.CheckSCRCPY(); err != nil {
		return fmt.Errorf("conn: %w", err)
	}

	if err := c.conn.scrcpy.BackOrScreenOn(action); err != nil {
		return fmt.Errorf("scrcpy: %w", err)
	}

	return nil
}

func (c *SCRCPY) CollapsePanels() error {
	if err := c.conn.CheckSCRCPY(); err != nil {
		return fmt.Errorf("conn: %w", err)
	}

	if err := c.conn.scrcpy.CollapsePanels(); err != nil {
		return fmt.Errorf("scrcpy: %w", err)
	}

	return nil
}

func (c *SCRCPY) ExpandNotificationPanel() error {
	if err := c.conn.CheckSCRCPY(); err != nil {
		return fmt.Errorf("conn: %w", err)
	}

	if err := c.conn.scrcpy.ExpandNotificationPanel(); err != nil {
		return fmt.Errorf("scrcpy: %w", err)
	}

	return nil
}

func (c *SCRCPY) ExpandSettingsPanel() error {
	if err := c.conn.CheckSCRCPY(); err != nil {
		return fmt.Errorf("conn: %w", err)
	}

	if err := c.conn.scrcpy.ExpandSettingsPanel(); err != nil {
		return fmt.Errorf("scrcpy: %w", err)
	}

	return nil
}

func (c *SCRCPY) OpenHardKeyboardSettings() error {
	if err := c.conn.CheckSCRCPY(); err != nil {
		return fmt.Errorf("conn: %w", err)
	}

	if err := c.conn.scrcpy.OpenHardKeyboardSettings(); err != nil {
		return fmt.Errorf("scrcpy: %w", err)
	}

	return nil
}

func (c *SCRCPY) InjectKeycode(keycode uint32, action byte, repeat, meta uint32) error {
	if err := c.conn.CheckSCRCPY(); err != nil {
		return fmt.Errorf("conn: %w", err)
	}

	if err := c.conn.scrcpy.InjectKeycode(keycode, action, repeat, meta); err != nil {
		return fmt.Errorf("scrcpy: %w", err)
	}

	return nil
}

func (c *SCRCPY) InjectScroll(x, y int32, hscroll, vscroll int16, buttons uint32) error {
	if err := c.conn.CheckSCRCPY(); err != nil {
		return fmt.Errorf("conn: %w", err)
	}

	if err := c.conn.scrcpy.InjectScroll(x, y, hscroll, vscroll, buttons); err != nil {
		return fmt.Errorf("scrcpy: %w", err)
	}

	return nil
}

func (c *SCRCPY) InjectTouch(action byte, pointerID uint64, x, y uint32, pressure uint16, actionButton, buttons uint32) error {
	if err := c.conn.CheckSCRCPY(); err != nil {
		return fmt.Errorf("conn: %w", err)
	}

	if err := c.conn.scrcpy.InjectTouch(action, pointerID, x, y, pressure, actionButton, buttons); err != nil {
		return fmt.Errorf("scrcpy: %w", err)
	}

	return nil
}

func (c *SCRCPY) waitClipboard(ctx context.Context) (string, error) {
	if err := c.conn.CheckSCRCPY(); err != nil {
		return "", fmt.Errorf("conn: %w", err)
	}

	select {
	case <-ctx.Done():
		return "", fmt.Errorf("waiter: %w", ctx.Err())
	case text := <-c.conn.clipboard:
		return text, nil
	}
}

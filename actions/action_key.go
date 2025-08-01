package actions

import (
	"context"
	"fmt"
	"time"

	abg "github.com/merzzzl/accessibility-bridge-go"
	scrcpy "github.com/merzzzl/scrcpy-go"
	"github.com/merzzzl/screen-flow/device"
)

const (
	ActionKeyHome = iota
	ActionKeyBack
	ActionKeyRecent
)

type ActionKey struct {
	Press int
}

func (s *ActionKey) Handle(ctx context.Context, conn *device.Conn) error {
	if conn.CheckABG() == nil {
		return s.abg(ctx, conn)
	}

	if conn.CheckSCRCPY() == nil {
		return s.scrcpy(conn)
	}

	return fmt.Errorf("need accessibility-bridge or scrcpy: %w", ErrNoClints)
}

func (s *ActionKey) scrcpy(conn *device.Conn) error {
	if err := conn.GetSCRCPY().InjectKeycode(actionKeyToSCRCPY(s.Press), scrcpy.ActionDown, scrcpy.AndroidKey0, 0); err != nil {
		return fmt.Errorf("inject action: %w", err)
	}

	time.Sleep(time.Millisecond * 50)

	if err := conn.GetSCRCPY().InjectKeycode(actionKeyToSCRCPY(s.Press), scrcpy.ActionUp, 0, 0); err != nil {
		return fmt.Errorf("inject action: %w", err)
	}

	return nil
}

func (s *ActionKey) abg(ctx context.Context, conn *device.Conn) error {
	if err := conn.GetABG().PerformAction(ctx, actionKeyToABG(s.Press)); err != nil {
		return fmt.Errorf("inject action: %w", err)
	}

	return nil
}

func actionKeyToSCRCPY(key int) uint32 {
	switch key {
	case ActionKeyHome:
		return scrcpy.AndroidKeyHome
	case ActionKeyBack:
		return scrcpy.AndroidKeyBack
	case ActionKeyRecent:
		return scrcpy.AndroidKeyRecentApps
	}

	return scrcpy.AndroidKeyHome
}

func actionKeyToABG(key int) abg.ActionKey_Key {
	switch key {
	case ActionKeyHome:
		return abg.ActionKey_GLOBAL_ACTION_HOME
	case ActionKeyBack:
		return abg.ActionKey_GLOBAL_ACTION_BACK
	case ActionKeyRecent:
		return abg.ActionKey_GLOBAL_ACTION_RECENTS
	}

	return abg.ActionKey_GLOBAL_ACTION_HOME
}

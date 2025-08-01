package actions

import (
	"context"
	"fmt"
	"regexp"
	"time"

	abg "github.com/merzzzl/accessibility-bridge-go"
	"github.com/merzzzl/screen-flow/device"
)

type ActionWaitElement struct {
	Regexp   string
	UniqueID string
	Duration *time.Duration
}

func (s *ActionWaitElement) Handle(ctx context.Context, conn *device.Conn) error {
	if err := conn.CheckABG(); err != nil {
		return fmt.Errorf("need accessibility-bridge: %w, %w", ErrNoClints, err)
	}

	startAt := time.Now()

	for {
		dump, err := conn.GetABG().ScreenDump(ctx)
		if err != nil {
			return fmt.Errorf("inject action: %w", err)
		}

		if s.check(dump) {
			return nil
		}

		time.Sleep(time.Microsecond * 200)

		if s.Duration != nil {
			if time.Now().After(startAt.Add(*s.Duration)) {
				return fmt.Errorf("find point: %w", ErrImageNotFound)
			}
		}
	}
}

func (s *ActionWaitElement) check(dump *abg.ScreenView) bool {
	ok := true

	if s.UniqueID != "" && ok {
		ok = dump.GetUniqueId() == s.UniqueID
	}

	if s.Regexp != "" && ok {
		rx, err := regexp.Compile(s.Regexp)
		if err != nil {
			return false
		}

		ok = rx.MatchString(dump.GetText())
	}

	if ok {
		return true
	}

	for _, c := range dump.GetChildren() {
		if s.check(c) {
			return true
		}
	}

	return false
}

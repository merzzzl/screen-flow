package control

import (
	"context"
	"fmt"
	"image"
	"time"

	scrcpy "github.com/merzzzl/scrcpy-go"
	"github.com/merzzzl/screen-flow/vision"
)

type StepActionTapXY struct {
	X        int
	Y        int
	Duration time.Duration
}

type StepActionSwipe struct {
	X1       int
	Y1       int
	X2       int
	Y2       int
	Duration time.Duration
}

type StepActionSetClipboard struct {
	Payload string
	Past    bool
}

type StepActionType struct {
	Payload string
}

type StepActionKeyboard struct {
	Press    []int
	Duration time.Duration
}

type StepActionTapImage struct {
	ImageTemplate image.Image
}

type StepActionSwipeImage struct {
	ImageTemplate image.Image
	H             int
	W             int
}

type StepActionStaticFrame struct {
	Threshold float64
}

type StepAction interface {
	Handle(dev *Device) error
}

type ImageTrigger struct {
	ImageTemplate image.Image
	SearchArea    *image.Rectangle
}

type FlowStep struct {
	StepAction   StepAction
	DelayBefor   time.Duration
	DelayAfter   time.Duration
	TriggerBefor *ImageTrigger
	TriggerAfer  *ImageTrigger
}

func (s *StepActionTapXY) Handle(dev *Device) error {
	if err := dev.client.InjectTouch(
		scrcpy.ActionDown, 1,
		uint32(s.X), uint32(s.Y),
		65535, scrcpy.ButtonPrimary, scrcpy.ButtonPrimary); err != nil {
		return fmt.Errorf("inject action: %w", err)
	}

	if s.Duration == 0 {
		s.Duration = 20 * time.Millisecond
	}

	time.Sleep(s.Duration)

	if err := dev.client.InjectTouch(
		scrcpy.ActionUp, 1,
		uint32(s.X), uint32(s.Y),
		65535, scrcpy.ButtonPrimary, 0); err != nil {
		return fmt.Errorf("inject action: %w", err)
	}

	return nil
}

func (s *StepActionSwipe) Handle(dev *Device) error {
	if err := dev.client.InjectTouch(
		scrcpy.ActionDown, 1,
		uint32(s.X1), uint32(s.Y1),
		65535, scrcpy.ButtonPrimary, scrcpy.ButtonPrimary); err != nil {
		return fmt.Errorf("inject action: %w", err)
	}

	if s.Duration == 0 {
		s.Duration = 200 * time.Millisecond
	}

	steps := s.Duration / 10
	currentX := s.X1
	currentY := s.Y1
	moveX := (s.X2 - s.X1) / int(steps)
	moveY := (s.Y2 - s.Y1) / int(steps)

	for range steps {
		currentX += moveX
		currentY += moveY

		if err := dev.client.InjectTouch(
			scrcpy.ActionMove, 1,
			uint32(currentX), uint32(currentY),
			65535, 0, scrcpy.ButtonPrimary); err != nil {
			return fmt.Errorf("inject action: %w", err)
		}
	}

	if err := dev.client.InjectTouch(
		scrcpy.ActionUp, 1,
		uint32(s.X2), uint32(s.Y2),
		65535, scrcpy.ButtonPrimary, 0); err != nil {
		return fmt.Errorf("inject action: %w", err)
	}

	return nil
}

func (s *StepActionSetClipboard) Handle(dev *Device) error {
	if err := dev.client.SetClipboard(0, s.Payload, s.Past); err != nil {
		return fmt.Errorf("inject action: %w", err)
	}

	return nil
}

func (s *StepActionType) Handle(dev *Device) error {
	if err := dev.client.InjectText(s.Payload); err != nil {
		return fmt.Errorf("inject action: %w", err)
	}

	return nil
}

func (s *StepActionKeyboard) Handle(dev *Device) error {
	for _, key := range s.Press {
		if err := dev.client.InjectKeycode(uint32(key), scrcpy.ActionDown, 0, 0); err != nil {
			return fmt.Errorf("inject action: %w", err)
		}
	}

	time.Sleep(s.Duration)

	for _, key := range s.Press {
		if err := dev.client.InjectKeycode(uint32(key), scrcpy.ActionUp, 0, 0); err != nil {
			return fmt.Errorf("inject action: %w", err)
		}
	}

	return nil
}

func (s *StepActionTapImage) Handle(dev *Device) error {
	source := dev.decoder.Image()

	point, _, err := vision.FindPoint(source, s.ImageTemplate)
	if err != nil {
		return fmt.Errorf("find point: %w", err)
	}

	nextStep := StepActionTapXY{
		X: point.X,
		Y: point.Y,
	}

	return nextStep.Handle(dev)
}

func (s *StepActionSwipeImage) Handle(dev *Device) error {
	source := dev.decoder.Image()

	point, _, err := vision.FindPoint(source, s.ImageTemplate)
	if err != nil {
		return fmt.Errorf("find point: %w", err)
	}

	nextStep := StepActionSwipe{
		X1: point.X,
		Y1: point.Y,
		X2: point.X + s.H,
		Y2: point.Y + s.W,
	}

	return nextStep.Handle(dev)
}

func (s *StepActionStaticFrame) Handle(dev *Device) error {
	first := dev.decoder.Image()

	for {
		time.Sleep(time.Millisecond * 250)

		second := dev.decoder.Image()

		ok, err := vision.IsFrameStatic(first, second, s.Threshold)
		if err != nil {
			return fmt.Errorf("compare frames: %w", err)
		}

		if ok {
			break
		}

		first = second
	}

	return nil
}

func (f *FlowStep) HandleDelayBefor(ctx context.Context) error {
	t := time.NewTimer(f.DelayBefor)

	select {
	case <-t.C:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("context: %w", ctx.Err())
	}
}

func (f *FlowStep) HandleDelayAfter(ctx context.Context) error {
	t := time.NewTimer(f.DelayAfter)

	select {
	case <-t.C:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("context: %w", ctx.Err())
	}
}

func (f *FlowStep) HandleTriggerBefor(ctx context.Context, dev *Device) error {
	if f.TriggerBefor == nil {
		return nil
	}

	for {
		source := dev.decoder.Image()

		point, ok, err := vision.FindPoint(source, f.TriggerBefor.ImageTemplate)
		if err != nil {
			return fmt.Errorf("find point: %w", err)
		}

		if ok {
			if f.TriggerBefor.SearchArea == nil {
				return nil
			}

			if point.In(*f.TriggerBefor.SearchArea) {
				return nil
			}
		}

		t := time.NewTimer(time.Millisecond * 200)

		select {
		case <-t.C:
			continue
		case <-ctx.Done():
			return fmt.Errorf("context: %w", ctx.Err())
		}
	}
}

func (f *FlowStep) HandleTriggerAfter(ctx context.Context, dev *Device) error {
	if f.TriggerAfer == nil {
		return nil
	}

	for {
		source := dev.decoder.Image()

		point, ok, err := vision.FindPoint(source, f.TriggerAfer.ImageTemplate)
		if err != nil {
			return fmt.Errorf("find point: %w", err)
		}

		if ok {
			if f.TriggerBefor.SearchArea == nil {
				return nil
			}

			if point.In(*f.TriggerBefor.SearchArea) {
				return nil
			}
		}

		t := time.NewTimer(time.Millisecond * 200)

		select {
		case <-t.C:
			continue
		case <-ctx.Done():
			return fmt.Errorf("context: %w", ctx.Err())
		}
	}
}

func (f *FlowStep) HandleAction(dev *Device) error {
	return f.StepAction.Handle(dev)
}

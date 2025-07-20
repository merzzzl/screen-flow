package screenflow

import (
	"context"
	"fmt"
	"image"
	"time"

	"github.com/merzzzl/scrcpy-go"
	"github.com/merzzzl/screen-flow/actions"
	"github.com/merzzzl/screen-flow/device"
	"github.com/merzzzl/screen-flow/vision"
)

type Flow struct {
	address string
	steps   []FlowStep
}

type FlowState struct {
	StartAt        time.Time
	EndAt          time.Time
	StepsCount     int
	CompletedSteps int
	Handshake      scrcpy.Handshake
}

type FlowStep interface {
	Handle(conn *device.Conn) error
}

type customAction struct {
	handler func(conn *device.Conn) error
}

func NewFlow(address string) *Flow {
	return &Flow{
		address: address,
		steps:   make([]FlowStep, 0),
	}
}

func (f *Flow) Load(steps []FlowStep) *Flow {
	f.steps = append(f.steps, steps...)

	return f
}

func (f *Flow) Run(ctx context.Context, alg vision.Algorithm, window vision.Window) (*FlowState, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	conn, err := device.Connect(ctx, f.address, alg, window)
	if err != nil {
		return nil, fmt.Errorf("connect to device: %w", err)
	}

	state := &FlowState{
		StartAt:    time.Now(),
		StepsCount: len(f.steps),
		Handshake:  conn.GetHandshake(),
	}

	defer func() {
		state.EndAt = time.Now()
	}()

	for i, step := range f.steps {
		if ctx.Err() != nil {
			return state, ctx.Err()
		}

		if err := step.Handle(conn); err != nil {
			return state, fmt.Errorf("step %d failed: %w", i, err)
		}

		state.CompletedSteps++
	}

	return state, nil
}

func (a *customAction) Handle(conn *device.Conn) error {
	err := a.handler(conn)
	if err != nil {
		return fmt.Errorf("custom action: %w", err)
	}

	return nil
}

func ActionTapXY(x, y int) FlowStep {
	return &actions.ActionTapXY{
		X: x,
		Y: y,
	}
}

func (f *Flow) ActionTapXY(x, y int) *Flow {
	f.steps = append(f.steps, ActionTapXY(x, y))

	return f
}

func ActionSwipe(x1, y1, x2, y2 int) FlowStep {
	return &actions.ActionSwipe{
		X1: x1,
		Y1: y1,
		X2: x2,
		Y2: y2,
	}
}

func (f *Flow) ActionSwipe(x1, y1, x2, y2 int) *Flow {
	f.steps = append(f.steps, ActionSwipe(x1, y1, x2, y2))

	return f
}

func ActionSetClipboard(payload string, past bool) FlowStep {
	return &actions.ActionSetClipboard{
		Payload: payload,
		Past:    past,
	}
}

func (f *Flow) ActionSetClipboard(payload string, past bool) *Flow {
	f.steps = append(f.steps, ActionSetClipboard(payload, past))

	return f
}

func ActionType(payload string) FlowStep {
	return &actions.ActionType{
		Payload: payload,
	}
}

func (f *Flow) ActionType(payload string) *Flow {
	f.steps = append(f.steps, ActionType(payload))

	return f
}

func ActionKeyboard(keys []int, duration time.Duration) FlowStep {
	return &actions.ActionKeyboard{
		Press:    keys,
		Duration: duration,
	}
}

func (f *Flow) ActionKeyboard(keys []int, duration time.Duration) *Flow {
	f.steps = append(f.steps, ActionKeyboard(keys, duration))

	return f
}

func ActionTapImage(img image.Image, wait bool, area *image.Rectangle, dur *time.Duration) FlowStep {
	return &actions.ActionTapImage{
		ImageTemplate: img,
		Duration:      dur,
		SearchArea:    area,
		Wait:          wait,
	}
}

func (f *Flow) ActionTapImage(img image.Image, wait bool, area *image.Rectangle, dur *time.Duration) *Flow {
	f.steps = append(f.steps, ActionTapImage(img, wait, area, dur))

	return f
}

func ActionSwipeImage(img image.Image, h, w int, wait bool, area *image.Rectangle, dur *time.Duration) FlowStep {
	return &actions.StepActionSwipeImage{
		ImageTemplate: img,
		H:             h,
		W:             w,
		Duration:      dur,
		SearchArea:    area,
		Wait:          wait,
	}
}

func (f *Flow) ActionSwipeImage(img image.Image, h, w int, wait bool, area *image.Rectangle, dur *time.Duration) *Flow {
	f.steps = append(f.steps, ActionSwipeImage(img, h, w, wait, area, dur))

	return f
}

func ActionFunc(fn func(conn *device.Conn) error) FlowStep {
	return &customAction{handler: fn}
}

func (f *Flow) ActionFunc(fn func(conn *device.Conn) error) *Flow {
	f.steps = append(f.steps, ActionFunc(fn))

	return f
}

func ActionWait(img image.Image, area *image.Rectangle, dur *time.Duration) FlowStep {
	return &actions.ActionWait{
		ImageTemplate: img,
		SearchArea:    area,
		Duration:      dur,
	}
}

func (f *Flow) ActionWait(img image.Image, area *image.Rectangle, dur *time.Duration) *Flow {
	f.steps = append(f.steps, ActionWait(img, area, dur))

	return f
}

func ActionDelay(dur time.Duration) FlowStep {
	return &actions.ActionDelay{
		Duration: dur,
	}
}

func (f *Flow) ActionDelay(dur time.Duration) *Flow {
	f.steps = append(f.steps, ActionDelay(dur))

	return f
}

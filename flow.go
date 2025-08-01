package screenflow

import (
	"context"
	"fmt"
	"image"
	"time"

	"github.com/merzzzl/screen-flow/actions"
	"github.com/merzzzl/screen-flow/device"
)

type Flow struct {
	steps []FlowStep
}

type FlowState struct {
	StartAt        time.Time
	EndAt          time.Time
	StepsCount     int
	CompletedSteps int
}

type FlowStep interface {
	Handle(ctx context.Context, conn *device.Conn) error
}

type customAction struct {
	handler func(ctx context.Context, conn *device.Conn) error
}

func NewFlow() *Flow {
	return &Flow{
		steps: make([]FlowStep, 0),
	}
}

func (f *Flow) Load(steps []FlowStep) *Flow {
	f.steps = append(f.steps, steps...)

	return f
}

func (f *Flow) Run(ctx context.Context, options ...device.Option) (*FlowState, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	conn, err := device.Connect(ctx, options...)
	if err != nil {
		return nil, fmt.Errorf("connect to device: %w", err)
	}

	state := &FlowState{
		StartAt:    time.Now(),
		StepsCount: len(f.steps),
	}

	defer func() {
		state.EndAt = time.Now()
	}()

	for i, step := range f.steps {
		if ctx.Err() != nil {
			return state, ctx.Err()
		}

		if err := step.Handle(ctx, conn); err != nil {
			return state, fmt.Errorf("step %d failed: %w", i, err)
		}

		state.CompletedSteps++
	}

	return state, nil
}

func (a *customAction) Handle(ctx context.Context, conn *device.Conn) error {
	err := a.handler(ctx, conn)
	if err != nil {
		return fmt.Errorf("custom action: %w", err)
	}

	return nil
}

func ActionTap(x, y int) FlowStep {
	return &actions.ActionTap{
		X: x,
		Y: y,
	}
}

func (f *Flow) ActionTap(x, y int) *Flow {
	f.steps = append(f.steps, ActionTap(x, y))

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

func ActionType(payload string) FlowStep {
	return &actions.ActionType{
		Payload: payload,
	}
}

func (f *Flow) ActionType(payload string) *Flow {
	f.steps = append(f.steps, ActionType(payload))

	return f
}

func ActionKey(key int) FlowStep {
	return &actions.ActionKey{
		Press: key,
	}
}

func (f *Flow) ActionKey(key int) *Flow {
	f.steps = append(f.steps, ActionKey(key))

	return f
}

func ActionTapImage(img image.Image, area *image.Rectangle, dur time.Duration) FlowStep {
	return &actions.ActionTapImage{
		ImageTemplate: img,
		Duration:      dur,
		SearchArea:    area,
	}
}

func (f *Flow) ActionTapImage(img image.Image, area *image.Rectangle, dur time.Duration) *Flow {
	f.steps = append(f.steps, ActionTapImage(img, area, dur))

	return f
}

func ActionTapElement(regexp, uniqid string, dur time.Duration) FlowStep {
	return &actions.ActionTapElement{
		Regexp:   regexp,
		UniqueID: uniqid,
		Duration: dur,
	}
}

func (f *Flow) ActionTapElement(regexp, uniqid string, dur time.Duration) *Flow {
	f.steps = append(f.steps, ActionTapElement(regexp, uniqid, dur))

	return f
}

func ActionSwipeImage(img image.Image, h, w int, area *image.Rectangle, dur time.Duration) FlowStep {
	return &actions.ActionSwipeImage{
		ImageTemplate: img,
		H:             h,
		W:             w,
		Duration:      dur,
		SearchArea:    area,
	}
}

func (f *Flow) ActionSwipeImage(img image.Image, h, w int, area *image.Rectangle, dur time.Duration) *Flow {
	f.steps = append(f.steps, ActionSwipeImage(img, h, w, area, dur))

	return f
}

func ActionSwipeElement(regexp, uniqid string, h, w int, dur time.Duration) FlowStep {
	return &actions.ActionSwipeElement{
		Regexp:   regexp,
		UniqueID: uniqid,
		H:        h,
		W:        w,
		Duration: dur,
	}
}

func (f *Flow) ActionSwipeElement(regexp, uniqid string, h, w int, dur time.Duration) *Flow {
	f.steps = append(f.steps, ActionSwipeElement(regexp, uniqid, h, w, dur))

	return f
}

func ActionFunc(fn func(ctx context.Context, conn *device.Conn) error) FlowStep {
	return &customAction{handler: fn}
}

func (f *Flow) ActionFunc(fn func(ctx context.Context, conn *device.Conn) error) *Flow {
	f.steps = append(f.steps, ActionFunc(fn))

	return f
}

func ActionWaitImage(img image.Image, area *image.Rectangle, dur *time.Duration) FlowStep {
	return &actions.ActionWaitImage{
		ImageTemplate: img,
		SearchArea:    area,
		Duration:      dur,
	}
}

func (f *Flow) ActionWaitImage(img image.Image, area *image.Rectangle, dur *time.Duration) *Flow {
	f.steps = append(f.steps, ActionWaitImage(img, area, dur))

	return f
}

func ActionWaitElement(regexp, uniqid string, dur *time.Duration) FlowStep {
	return &actions.ActionWaitElement{
		Regexp:   regexp,
		UniqueID: uniqid,
		Duration: dur,
	}
}

func (f *Flow) ActionWaitElement(regexp, uniqid string, dur *time.Duration) *Flow {
	f.steps = append(f.steps, ActionWaitElement(regexp, uniqid, dur))

	return f
}

func ActionWait(dur time.Duration) FlowStep {
	return &actions.ActionWait{
		Duration: dur,
	}
}

func (f *Flow) ActionWait(dur time.Duration) *Flow {
	f.steps = append(f.steps, ActionWait(dur))

	return f
}

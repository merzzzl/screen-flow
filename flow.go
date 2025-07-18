package screenflow

import (
	"context"
	"fmt"
	"image"
	"time"

	"github.com/merzzzl/scrcpy-go"
	"github.com/merzzzl/screen-flow/actions"
	"github.com/merzzzl/screen-flow/device"
	"github.com/merzzzl/screen-flow/events"
)

type Flow struct {
	address string
	steps   []FlowStep
	stream  chan<- image.Image
	events  chan<- *events.Base
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

func NewFlow(address string, s chan<- image.Image, e chan<- *events.Base) *Flow {
	return &Flow{
		address: address,
		steps:   make([]FlowStep, 0),
		stream:  s,
		events:  e,
	}
}

func (f *Flow) Run(ctx context.Context) (*FlowState, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	conn, err := device.Connect(ctx, f.address)
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

	if f.stream != nil {
		streamCh := conn.GetScreenStream()

		go func() {
			for ctx.Err() == nil {
				select {
				case <-ctx.Done():
					return
				case img := <-streamCh:
					f.stream <- img
				}
			}
		}()
	}

	if f.stream != nil {
		eventCh := conn.GetEventStream()

		go func() {
			for ctx.Err() == nil {
				select {
				case <-ctx.Done():
					return
				case ev := <-eventCh:
					f.events <- ev
				}
			}
		}()
	}

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

func (f *Flow) ActionTapXY(x, y int) *Flow {
	action := actions.ActionTapXY{
		X: x,
		Y: y,
	}

	f.steps = append(f.steps, &action)

	return f
}

func (f *Flow) ActionSwipe(x1, y1, x2, y2 int) *Flow {
	action := actions.ActionSwipe{
		X1: x1,
		Y1: y1,
		X2: x2,
		Y2: y2,
	}

	f.steps = append(f.steps, &action)

	return f
}

func (f *Flow) ActionSetClipboard(payload string, past bool) *Flow {
	action := actions.ActionSetClipboard{
		Payload: payload,
		Past:    past,
	}

	f.steps = append(f.steps, &action)

	return f
}

func (f *Flow) ActionType(payload string) *Flow {
	action := actions.ActionType{
		Payload: payload,
	}

	f.steps = append(f.steps, &action)

	return f
}

func (f *Flow) ActionKeyboard(keys []int, duration time.Duration) *Flow {
	action := actions.ActionKeyboard{
		Press:    keys,
		Duration: duration,
	}

	f.steps = append(f.steps, &action)

	return f
}

func (f *Flow) ActionTapImage(img image.Image, wait bool, area *image.Rectangle, dur *time.Duration) *Flow {
	action := actions.ActionTapImage{
		ImageTemplate: img,
		Duration:      dur,
		SearchArea:    area,
		Wait:          wait,
	}

	f.steps = append(f.steps, &action)

	return f
}

func (f *Flow) ActionSwipeImage(img image.Image, h, w int, wait bool, area *image.Rectangle, dur *time.Duration) *Flow {
	action := actions.StepActionSwipeImage{
		ImageTemplate: img,
		H:             h,
		W:             w,
		Duration:      dur,
		SearchArea:    area,
		Wait:          wait,
	}

	f.steps = append(f.steps, &action)

	return f
}

func (f *Flow) ActionFunc(fn func(conn *device.Conn) error) *Flow {
	f.steps = append(f.steps, &customAction{handler: fn})

	return f
}

func (f *Flow) WaitStaticFrame(threshold float64) *Flow {
	action := actions.ActionWaitStaticFrame{
		Threshold: threshold,
	}

	f.steps = append(f.steps, &action)

	return f
}

func (f *Flow) ActionWait(img image.Image, area *image.Rectangle, dur *time.Duration) *Flow {
	action := actions.ActionWait{
		ImageTemplate: img,
		SearchArea:    area,
		Duration:      dur,
	}

	f.steps = append(f.steps, &action)

	return f
}

func (f *Flow) ActionDelay(dur time.Duration) *Flow {
	action := actions.ActionDelay{
		Duration: dur,
	}

	f.steps = append(f.steps, &action)

	return f
}

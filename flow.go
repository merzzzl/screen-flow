package screenflow

import (
	"context"
	"fmt"
	"image"
	"time"

	"github.com/merzzzl/screen-flow/control"
)

type Flow struct {
	address string
	steps   []*control.FlowStep
}

type FlowOption interface {
	apply(step *control.FlowStep)
}

type flowOption struct {
	fn func(step *control.FlowStep)
}

type customAction struct {
	handler func(dev *control.Device) error
}

func NewFlow(address string) *Flow {
	return &Flow{
		address: address,
		steps:   make([]*control.FlowStep, 0),
	}
}

func (f *Flow) Run(ctx context.Context) (*control.FlowState, error) {
	dev, err := control.Connect(ctx, f.address)
	if err != nil {
		return nil, fmt.Errorf("connect to device: %w", err)
	}

	state, err := dev.RunFlow(ctx, f.steps)
	if err != nil {
		return state, fmt.Errorf("exec flow: %w", err)
	}

	return state, nil
}

func (o *flowOption) apply(step *control.FlowStep) {
	o.fn(step)
}

func OptionDelayBefor(t time.Duration) FlowOption {
	return &flowOption{
		fn: func(step *control.FlowStep) {
			step.DelayBefor = t
		},
	}
}

func OptionDelayAfter(t time.Duration) FlowOption {
	return &flowOption{
		fn: func(step *control.FlowStep) {
			step.DelayBefor = t
		},
	}
}

func OptionTriggerBefor(img image.Image, area *image.Rectangle) FlowOption {
	return &flowOption{
		fn: func(step *control.FlowStep) {
			step.TriggerBefor = &control.ImageTrigger{
				ImageTemplate: img,
				SearchArea:    area,
			}
		},
	}
}

func OptionTriggerAfer(img image.Image, area *image.Rectangle) FlowOption {
	return &flowOption{
		fn: func(step *control.FlowStep) {
			step.TriggerAfer = &control.ImageTrigger{
				ImageTemplate: img,
				SearchArea:    area,
			}
		},
	}
}

func newStep(action control.StepAction, opts []FlowOption) *control.FlowStep {
	step := control.FlowStep{
		StepAction: action,
	}

	for _, o := range opts {
		o.apply(&step)
	}

	return &step
}

func (a *customAction) Handle(dev *control.Device) error {
	if err := a.handler(dev); err != nil {
		return fmt.Errorf("custom action: %w", err)
	}

	return nil
}

func (f *Flow) ActionTapXY(x, y int, opts ...FlowOption) *Flow {
	action := control.StepActionTapXY{
		X: x,
		Y: y,
	}

	f.steps = append(f.steps, newStep(&action, opts))

	return f
}

func (f *Flow) ActionSwipe(x1, y1, x2, y2 int, opts ...FlowOption) *Flow {
	action := control.StepActionSwipe{
		X1: x1,
		Y1: y1,
		X2: x2,
		Y2: y2,
	}

	f.steps = append(f.steps, newStep(&action, opts))

	return f
}

func (f *Flow) ActionSetClipboard(payload string, past bool, opts ...FlowOption) *Flow {
	action := control.StepActionSetClipboard{
		Payload: payload,
		Past:    past,
	}

	f.steps = append(f.steps, newStep(&action, opts))

	return f
}

func (f *Flow) ActionType(payload string, opts ...FlowOption) *Flow {
	action := control.StepActionType{
		Payload: payload,
	}

	f.steps = append(f.steps, newStep(&action, opts))

	return f
}

func (f *Flow) ActionKeyboard(keys []int, duration time.Duration, opts ...FlowOption) *Flow {
	action := control.StepActionKeyboard{
		Press:    keys,
		Duration: duration,
	}

	f.steps = append(f.steps, newStep(&action, opts))

	return f
}

func (f *Flow) ActionTapImage(img image.Image, opts ...FlowOption) *Flow {
	action := control.StepActionTapImage{
		ImageTemplate: img,
	}

	f.steps = append(f.steps, newStep(&action, opts))

	return f
}

func (f *Flow) ActionSwipeImage(img image.Image, h, w int, opts ...FlowOption) *Flow {
	action := control.StepActionSwipeImage{
		ImageTemplate: img,
		H:             h,
		W:             w,
	}

	f.steps = append(f.steps, newStep(&action, opts))

	return f
}

func (f *Flow) ActionFunc(fn func(dev *control.Device) error, opts ...FlowOption) *Flow {
	f.steps = append(f.steps, newStep(&customAction{handler: fn}, opts))

	return f
}

func (f *Flow) WaitStaticFrame(threshold float64, opts ...FlowOption) *Flow {
	action := control.StepActionStaticFrame{
		Threshold: threshold,
	}

	f.steps = append(f.steps, newStep(&action, opts))

	return f
}

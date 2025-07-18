package control

import (
	"context"
	"fmt"
	"time"

	scrcpy "github.com/merzzzl/scrcpy-go"
)

type Device struct {
	err       chan error
	client    *scrcpy.Client
	decoder   *Decoder
	clipboard chan string
}

type FlowState struct {
	StartAt        time.Time
	EndAt          time.Time
	StepsCount     int
	CompletedSteps int
	Handshake      scrcpy.Handshake
}

func Connect(ctx context.Context, addr string) (*Device, error) {
	dev := &Device{
		err:       make(chan error, 1),
		clipboard: make(chan string, 1),
	}

	client, err := scrcpy.Dial(ctx, addr)
	if err != nil {
		return nil, fmt.Errorf("init scrcpy: %w", err)
	}

	dec, err := NewDecoder()
	if err != nil {
		return nil, fmt.Errorf("init decoder: %w", err)
	}

	client.SetVideoHandler(func(frame []byte) {
		err := dec.Decode(frame)
		if err != nil {
			dev.err <- err
		}
	})

	client.SetControlHandler(func(cm scrcpy.ControlMessage) {
		if cm.Type == 0 {
			dev.clipboard <- string(cm.Payload)
		}
	})

	ctx, cancel := context.WithCancel(ctx)

	go func() {
		err := client.Serve(ctx)
		if err != nil {
			dev.err <- err
		}

		cancel()
	}()

	go func() {
		err := dec.Looper(ctx)
		if err != nil {
			dev.err <- err
		}

		cancel()
	}()

	dev.client = client
	dev.decoder = dec

	return dev, nil
}

func (d *Device) GetClipboard() (string, error) {
	if err := d.client.GetClipboard(0); err != nil {
		return "", err
	}

	t := time.NewTimer(time.Second)

	select {
	case payload := <-d.clipboard:
		return payload, nil
	case <-t.C:
		return "", nil
	}
}

func (d *Device) RunFlow(ctx context.Context, steps []*FlowStep) (*FlowState, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var flowErr error

	go func() {
		select {
		case err := <-d.err:
			flowErr = err
		case <-ctx.Done():
			flowErr = ctx.Err()
		}
	}()

	state := &FlowState{
		StartAt:    time.Now(),
		StepsCount: len(steps),
		Handshake:  d.client.GetHandshake(),
	}

	defer func() {
		state.EndAt = time.Now()
	}()

	for flowErr == nil {
		if d.decoder.Image() != nil {
			break
		}

		time.Sleep(time.Second)
	}

	for i, step := range steps {
		if flowErr != nil {
			return state, fmt.Errorf("context: %w", flowErr)
		}

		if err := d.runFlowStep(ctx, step); err != nil {
			return state, fmt.Errorf("sted %d failed: %w", i, err)
		}

		state.CompletedSteps++
	}

	return state, nil
}

func (d *Device) runFlowStep(ctx context.Context, step *FlowStep) error {
	if err := step.HandleDelayBefor(ctx); err != nil {
		return fmt.Errorf("delay: %w", err)
	}

	if err := step.HandleTriggerBefor(ctx, d); err != nil {
		return fmt.Errorf("trigger: %w", err)
	}

	if err := step.HandleAction(d); err != nil {
		return fmt.Errorf("action: %w", err)
	}

	if err := step.HandleDelayAfter(ctx); err != nil {
		return fmt.Errorf("delay: %w", err)
	}

	if err := step.HandleTriggerAfter(ctx, d); err != nil {
		return fmt.Errorf("trigger: %w", err)
	}

	return nil
}

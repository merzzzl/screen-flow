package events

import (
	"time"
)

type Action uint32

const (
	EventActionTouch Action = iota
	EventActionTypeText
	EventActionPastText
	EventActionPressKey
	EventActionFoundImage
)

type Base struct {
	action  Action
	ts      time.Time
	payload string
	keys    []int
	x       int
	y       int
}

func NewTouchEvent(x, y int) *Base {
	return &Base{
		action: EventActionTouch,
		ts:     time.Now(),
		x:      x,
		y:      y,
	}
}

func NewTypeTextEvent(s string) *Base {
	return &Base{
		action:  EventActionTypeText,
		ts:      time.Now(),
		payload: s,
	}
}

func NewPastTextEvent(s string) *Base {
	return &Base{
		action:  EventActionPastText,
		ts:      time.Now(),
		payload: s,
	}
}

func NewPressKeyEvent(keys []int) *Base {
	return &Base{
		action: EventActionPressKey,
		ts:     time.Now(),
		keys:   keys,
	}
}

func NewFoundImageEvent(x, y int) *Base {
	return &Base{
		action: EventActionFoundImage,
		ts:     time.Now(),
		x:      x,
		y:      y,
	}
}

func (e *Base) Time() time.Time {
	return e.ts
}

func (e *Base) Action() Action {
	return e.action
}

func (e *Base) GetPayload() (string, bool) {
	return e.payload, (e.action == EventActionPastText || e.action == EventActionTypeText)
}

func (e *Base) GetXY() (int, int, bool) {
	return e.x, e.y, (e.action == EventActionTouch || e.action == EventActionFoundImage)
}

func (e *Base) GetKeys() ([]int, bool) {
	return e.keys, (e.action == EventActionPressKey)
}

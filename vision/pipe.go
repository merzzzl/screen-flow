package vision

import (
	"context"
	"fmt"
	"image"
	"io"
	"sync/atomic"

	"gocv.io/x/gocv"
)

type Pipe struct {
	window  Window
	tmpl    atomic.Pointer[gocv.Mat]
	success atomic.Uint32
	point   chan image.Point
	algo    Algorithm
	r       io.Reader
	h       int
	w       int
}

type Window interface {
	Resize(width int, height int)
	Show(img gocv.Mat)
}

func NewPipe(stream io.Reader, w, h int, algo Algorithm, window Window) *Pipe {
	if window != nil {
		window.Resize(w, h)
	}

	return &Pipe{
		r:       stream,
		w:       w,
		h:       h,
		point:   make(chan image.Point),
		tmpl:    atomic.Pointer[gocv.Mat]{},
		success: atomic.Uint32{},
		window:  window,
		algo:    algo,
	}
}

func (p *Pipe) Process(ctx context.Context) error {
	var (
		static    uint32
		lastPoint *image.Point
		prev      *gocv.Mat
		frame     = make([]byte, p.w*p.h*3)
	)

	defer func() {
		if prev != nil {
			_ = prev.Close()
		}

		close(p.point)
	}()

	for ctx.Err() == nil {
		if _, err := io.ReadFull(p.r, frame); err != nil {
			return fmt.Errorf("read frame: %w", err)
		}

		next, err := gocv.NewMatFromBytes(p.h, p.w, gocv.MatTypeCV8UC3, frame)
		if err != nil {
			return fmt.Errorf("convert to mat: %w", err)
		}

		if prev != nil {
			change := calcChangeRatio(*prev, next)

			if change < 0.10 {
				static++

				if static > 120 {
					static = 120
				}
			} else {
				static = 0
			}

			_ = prev.Close()
		}

		prev = &next

		var (
			goodPoint image.Point
			allPoints []image.Point
		)

		if tmpl := p.tmpl.Load(); static > 30 && tmpl != nil {
			var ok bool

			goodPoint, allPoints, ok = findPoint(next, *tmpl, p.algo)

			if ok && lastPoint != nil && lastPoint.X == goodPoint.X && lastPoint.Y == goodPoint.Y {
				if p.success.Load() >= 5 {
					select {
					case p.point <- goodPoint:
					default:
					}
				} else {
					p.success.Add(1)
				}
			} else {
				p.success.Store(0)
			}

			if ok {
				lastPoint = &goodPoint
			}
		}

		if p.window != nil {
			p.showImage(next, goodPoint, allPoints)
		}
	}

	return nil
}

func (p *Pipe) Found(img image.Image) (image.Point, error) {
	obj, err := toMat(img)
	if err != nil {
		return image.Pt(0, 0), fmt.Errorf("convert to mat: %w", err)
	}

	p.tmpl.Store(&obj)
	p.success.Store(0)

	defer func() {
		p.tmpl.Store(nil)

		_ = obj.Close()
	}()

	return <-p.point, nil
}

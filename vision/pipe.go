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
	tmpl    atomic.Pointer[gocv.Mat]
	success atomic.Uint32
	point   chan image.Point
	algo    Algorithm
	r       io.Reader
	h       int
	w       int
}

func NewPipe(stream io.Reader, w, h int, algo Algorithm) *Pipe {
	return &Pipe{
		r:       stream,
		w:       w,
		h:       h,
		point:   make(chan image.Point),
		tmpl:    atomic.Pointer[gocv.Mat]{},
		success: atomic.Uint32{},
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

		nextSrc, scale := resizeSrc(next)
		_ = next.Close()

		if prev != nil {
			change := calcChangeRatio(*prev, nextSrc)

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

		prev = &nextSrc

		var res *Result

		if tmpl := p.tmpl.Load(); static > 30 && tmpl != nil {
			var (
				ok    bool
				point image.Point
			)

			tpl := resizeTpl(*tmpl, scale)
			res, ok = findPoint(nextSrc, tpl, p.algo)
			_ = tpl.Close()

			if ok {
				point = restorePoint(res.best, scale)
			}

			if ok && lastPoint != nil && lastPoint.X == point.X && lastPoint.Y == point.Y {
				if p.success.Load() >= 5 {
					select {
					case p.point <- point:
					default:
					}
				} else {
					p.success.Add(1)
				}
			} else {
				p.success.Store(0)
			}

			if ok {
				lastPoint = &point
			}
		}

		if res != nil {
			res.Close()
		}
	}

	return nil
}

func (p *Pipe) Find(ctx context.Context, img image.Image) (image.Point, error) {
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

	select {
	case <-ctx.Done():
		return image.Pt(0, 0), ctx.Err()
	case point := <-p.point:
		return point, nil
	}
}

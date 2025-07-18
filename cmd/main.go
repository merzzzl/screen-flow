package main

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"time"

	screenflow "github.com/merzzzl/screen-flow"
	"github.com/merzzzl/screen-flow/events"
	"gocv.io/x/gocv"
)

func main() {
	chromeImage, err := loadimage("chrome.png")
	if err != nil {
		panic(err)
	}

	searchImage, err := loadimage("search.png")
	if err != nil {
		panic(err)
	}

	enterImage, err := loadimage("enter.png")
	if err != nil {
		panic(err)
	}

	sCh := make(chan image.Image, 1)
	eCh := make(chan *events.Base, 1)

	go func() {
		flow := screenflow.NewFlow("127.0.0.1:10000", sCh, eCh).
			ActionTapImage(chromeImage, true, nil, nil).
			WaitStaticFrame(0.05).
			ActionTapImage(searchImage, true, nil, nil).
			WaitStaticFrame(0.05).
			ActionWait(enterImage, nil, nil).
			ActionType("Hello, world!").
			ActionTapImage(enterImage, true, nil, nil).
			ActionDelay(time.Second * 2)

		s, _ := flow.ToYAML()
		_ = flow.FromYAML(s)

		log.Println(s)

		state, err := flow.Run(context.Background())
		if err != nil {
			log.Printf("run flow: %v", err)
		}

		log.Printf("%+v", state)

		close(sCh)
		close(eCh)
	}()

	showImage(sCh, eCh)
}

func showImage(s <-chan image.Image, e <-chan *events.Base) {
	window := gocv.NewWindow("screen-flow")
	tpoints := make([]tpoint, 0)

	go func() {
		for event := range e {
			log.Printf("%+v\n", event)

			x, y, ok := event.GetXY()
			if !ok {
				continue
			}

			tpoints = append(tpoints, tpoint{
				p: image.Pt(x, y),
				t: event.Time(),
			})
		}
	}()

	for img := range s {
		if img == nil {
			return
		}

		bounds := img.Bounds()
		x := bounds.Dx()
		y := bounds.Dy()
		bytes := make([]byte, 0, x*y*3)

		for j := bounds.Min.Y; j < bounds.Max.Y; j++ {
			for i := bounds.Min.X; i < bounds.Max.X; i++ {
				r, g, b, _ := img.At(i, j).RGBA()
				bytes = append(bytes, byte(b>>8), byte(g>>8), byte(r>>8))
			}
		}

		rgb, err := gocv.NewMatFromBytes(y, x, gocv.MatTypeCV8UC3, bytes)
		if err != nil {
			return
		}

		if !window.IsOpen() || rgb.Empty() {
			return
		}

		for i := 0; i < len(tpoints); i++ {
			if tpoints[i].t.Add(time.Second * 2).Before(time.Now()) {
				tpoints = append(tpoints[:i], tpoints[i+1:]...)

				continue
			}

			drowpoint(&rgb, tpoints[i].p, color.RGBA{0, 255, 0, 255})
		}

		window.ResizeWindow(x, y)
		window.IMShow(rgb)
		rgb.Close()

		key := window.WaitKey(10)

		if key == 27 {
			break
		}
	}
}

func loadimage(file string) (image.Image, error) {
	imageFile, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("open image: %w", err)
	}

	img, err := png.Decode(imageFile)
	if err != nil {
		return nil, fmt.Errorf("decode image: %w", err)
	}

	return img, nil
}

func drowpoint(mat *gocv.Mat, point image.Point, cl color.RGBA) {
	if point.X > 0 && point.Y > 0 {
		startH := image.Pt(point.X-40, point.Y)
		endH := image.Pt(point.X+40, point.Y)

		_ = gocv.Line(mat, startH, endH, cl, 8)

		startV := image.Pt(point.X, point.Y-40)
		endV := image.Pt(point.X, point.Y+40)

		_ = gocv.Line(mat, startV, endV, cl, 8)
	}
}

type tpoint struct {
	p image.Point
	t time.Time
}

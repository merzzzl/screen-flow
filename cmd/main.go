package main

import (
	"context"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"runtime"
	"time"

	screenflow "github.com/merzzzl/screen-flow"
	"github.com/merzzzl/screen-flow/vision"
	"gocv.io/x/gocv"
)

func main() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	window := initWindow()

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

	flow := screenflow.NewFlow("127.0.0.1:10000").
		ActionTapImage(chromeImage, true, nil, nil).
		ActionTapImage(searchImage, true, nil, nil).
		ActionWait(enterImage, nil, nil).
		ActionType("Hello, world!").
		ActionTapImage(enterImage, true, nil, nil).
		ActionDelay(time.Second * 2)

	s, _ := flow.ToYAML()
	_ = flow.FromYAML(s)

	log.Println(s)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		state, err := flow.Run(ctx, vision.AlgorithmSIFT, window)
		if err != nil {
			log.Printf("run flow: %v", err)
		}

		log.Printf("%+v", state)

		cancel()
	}()

	if err := window.Handler(ctx); err != nil {
		log.Printf("window: %v", err)
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

type faceWindow struct {
	resize chan any
	show   chan gocv.Mat
	window gocv.Window
	w      int
	h      int
}

func initWindow() *faceWindow {
	win := gocv.NewWindow("screen-flow")

	return &faceWindow{
		resize: make(chan any, 1),
		show:   make(chan gocv.Mat, 1),
		window: *win,
	}
}

func (fw *faceWindow) Resize(x, y int) {
	fw.w = x
	fw.h = y

	fw.resize <- struct{}{}
}

func (fw *faceWindow) Show(img gocv.Mat) {
	dstSrc := gocv.NewMat()
	defer dstSrc.Close()

	gocv.Resize(img, &dstSrc, image.Pt(fw.w, fw.h), 0, 0, gocv.InterpolationLinear)

	fw.show <- dstSrc.Clone()
}

func (fw *faceWindow) Handler(ctx context.Context) error {
	defer func() {
		_ = fw.window.Close()
	}()

	for ctx.Err() == nil {
		select {
		case <-fw.resize:
			if err := fw.window.ResizeWindow(fw.w, fw.h); err != nil {
				return err
			}
		case img := <-fw.show:
			if err := fw.window.IMShow(img); err != nil {
				return err
			}

			_ = img.Close()
		default:
			key := fw.window.WaitKey(1)
			if key == 'q' {
				return nil
			}
		}
	}

	return nil
}

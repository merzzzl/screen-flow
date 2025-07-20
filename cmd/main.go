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
		state, err := flow.Run(ctx, window)
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
	resize chan image.Rectangle
	show   chan gocv.Mat
	window gocv.Window
}

func initWindow() *faceWindow {
	return &faceWindow{
		resize: make(chan image.Rectangle, 1),
		show:   make(chan gocv.Mat, 1),
		window: *gocv.NewWindow("screen-flow"),
	}
}

func (fw *faceWindow) Resize(x, y int) {
	fw.resize <- image.Rect(0, 0, x, y)
}

func (fw *faceWindow) Show(img gocv.Mat) {
	fw.show <- img.Clone()
}

func (fw *faceWindow) Handler(ctx context.Context) error {
	defer func() {
		_ = fw.window.Close()
	}()

	for ctx.Err() == nil {
		select {
		case rec := <-fw.resize:
			if err := fw.window.ResizeWindow(rec.Max.X, rec.Max.Y); err != nil {
				return err
			}
		case img := <-fw.show:
			if err := fw.window.IMShow(img); err != nil {
				return err
			}
		default:
			key := fw.window.WaitKey(1)
			if key == 'q' {
				return nil
			}
		}
	}

	return nil
}

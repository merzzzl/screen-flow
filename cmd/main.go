package main

import (
	"context"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"time"

	screenflow "github.com/merzzzl/screen-flow"
)

func main() {
	chromeImage, err := loadImage("chrome.png")
	if err != nil {
		panic(err)
	}

	searchImage, err := loadImage("search.png")
	if err != nil {
		panic(err)
	}

	enterImage, err := loadImage("enter.png")
	if err != nil {
		panic(err)
	}

	state, err := screenflow.NewFlow("127.0.0.1:10000").
		ActionTapImage(chromeImage, screenflow.OptionTriggerBefor(chromeImage, nil)).
		WaitStaticFrame(0.05).
		ActionTapImage(searchImage, screenflow.OptionTriggerBefor(searchImage, nil)).
		WaitStaticFrame(0.05).
		ActionType("Hello, world!", screenflow.OptionTriggerBefor(enterImage, nil)).
		ActionTapImage(enterImage, screenflow.OptionDelayAfter(time.Second)).
		Run(context.Background())
	if err != nil {
		log.Printf("run flow: %v", err)

		return
	}

	log.Printf("%+v", state)
}

func loadImage(file string) (image.Image, error) {
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

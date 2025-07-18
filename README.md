# screen‑flow

screen‑flow is a **UI automation & testing helper** built on top of [`scrcpy‑go`](https://github.com/merzzzl/scrcpy-go).  

It lets you describe interaction flows as a chain of steps driven by
computer‑vision template matching: *tap that button, wait for the dialog,
type text, verify that the screen stopped changing …* — all from pure Go.

## 📚  Features

- **Declarative flows**  
  Build a test as a list of steps (`TapXY`, `Swipe`, `Type`, `SetClipboard`, …).

- **Image‑based actions & triggers**  
  Tap or swipe **relative to a template image**;  
  wait until an image *appears*/**disappears** before/after a step.

- **Static‑frame wait**  
  Pause the flow until two consecutive video frames differ less than
  `threshold` — great for “wait until loading stops”.

- **Custom Go callbacks**  
  Insert `ActionFunc()` to run arbitrary logic on the connected device.

- **Fine‑grained timing**  
  Per‑step `DelayBefore` / `DelayAfter`.

- **Powered by**  
  - [scrcpy-go](https://github.com/merzzzl/scrcpy-go) for control & video  
  - [GoCV](https://gocv.io/) for template matching  
  - [FFmpeg](https://ffmpeg.org/) for H.264 decoding

## 🗂  Step catalogue

| Step                                           | Description                                |
| ---------------------------------------------- | ------------------------------------------ |
| `ActionTapXY(x, y)`                            | Tap at absolute coordinates                |
| `ActionSwipe(x1, y1, x2, y2)`                  | Swipe from point A to B                    |
| `ActionKeyboard(keys, dur)`                    | Press keycodes with duration               |
| `ActionSetClipboard(str, paste)`               | Set clipboard text and optionally paste it |
| `ActionType(str)`                              | Type UTF‑8 string text                     |
| `ActionTapImage(img, wait, area, dur)`         | Wait for image on screen and tap center    |
| `ActionSwipeImage(img, h, w, wait, area, dur)` | Swipe from image anchor (H,W offset)       |
| `WaitStaticFrame(threshold)`                   | Wait until screen becomes still            |
| `ActionWait(img, area, dur)`                   | Wait until image appears on screen         |
| `ActionDelay(dur)`                             | Sleep for duration                         |
| `ActionFunc(fn)`                               | Execute custom Go callback                 |

## 🚀  Getting started

1) Install OpenCV & FFmpeg (system package manager)
2) Clone project
3) Run exampe flow ```go run cmd/*```

> Don’t forget to launch `scrcpy-server` (see `scrcpy-go make run`)
> so `screen‑flow` can connect to **tcp:10000**.

## 🛠  Requirements

| Tool / Library     | Version / Notes                                         |
| ------------------ | ------------------------------------------------------- |
| **Go**             | 1.22+                                                   |
| **OpenCV**         | 4.x (headers + libs)                                    |
| **FFmpeg CLI**     | `ffmpeg` available in `PATH`                            |

## 📄  License

MIT

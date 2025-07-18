package screenflow

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/png"
	"time"

	"github.com/merzzzl/screen-flow/actions"
	"gopkg.in/yaml.v2"
)

var ErrUnknownStep = errors.New("unknown step")

type stepYAML struct {
	Type      string         `yaml:"type"`
	X         int            `yaml:"x,omitempty"`
	Y         int            `yaml:"y,omitempty"`
	X2        int            `yaml:"x2,omitempty"`
	Y2        int            `yaml:"y2,omitempty"`
	Keys      []int          `yaml:"keys,omitempty"`
	Duration  *time.Duration `yaml:"duration,omitempty"`
	Payload   string         `yaml:"payload,omitempty"`
	Paste     bool           `yaml:"paste,omitempty"`
	Wait      bool           `yaml:"wait,omitempty"`
	H         int            `yaml:"h,omitempty"`
	W         int            `yaml:"w,omitempty"`
	Threshold float64        `yaml:"threshold,omitempty"`
	Area      []int          `yaml:"area,omitempty"`
	Img       string         `yaml:"img,omitempty"`
}

func encodeImage(img image.Image) (string, error) {
	var b bytes.Buffer

	if err := png.Encode(&b, img); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(b.Bytes()), nil
}

func decodeImage(s string) (image.Image, error) {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}

	return png.Decode(bytes.NewReader(data))
}

func (f *Flow) ToYAML() (string, error) {
	out := make([]stepYAML, len(f.steps))

	for i, s := range f.steps {
		switch v := s.(type) {
		case *actions.ActionTapXY:
			out[i] = stepYAML{Type: "tap_xy", X: v.X, Y: v.Y}

		case *actions.ActionSwipe:
			out[i] = stepYAML{Type: "swipe", X: v.X1, Y: v.Y1, X2: v.X2, Y2: v.Y2}

		case *actions.ActionKeyboard:
			out[i] = stepYAML{Type: "keyboard", Keys: v.Press, Duration: &v.Duration}

		case *actions.ActionSetClipboard:
			out[i] = stepYAML{Type: "set_clipboard", Payload: v.Payload, Paste: v.Past}

		case *actions.ActionType:
			out[i] = stepYAML{Type: "type", Payload: v.Payload}

		case *actions.ActionTapImage:
			img, err := encodeImage(v.ImageTemplate)
			if err != nil {
				return "", err
			}

			area := make([]int, 0)

			if v.SearchArea != nil {
				area = []int{v.SearchArea.Min.X, v.SearchArea.Min.Y, v.SearchArea.Max.X, v.SearchArea.Max.Y}
			}

			out[i] = stepYAML{
				Type:     "tap_image",
				Img:      img,
				Duration: v.Duration,
				Wait:     v.Wait,
				Area:     area,
			}

		case *actions.StepActionSwipeImage:
			img, err := encodeImage(v.ImageTemplate)
			if err != nil {
				return "", err
			}

			area := make([]int, 0)

			if v.SearchArea != nil {
				area = []int{v.SearchArea.Min.X, v.SearchArea.Min.Y, v.SearchArea.Max.X, v.SearchArea.Max.Y}
			}

			out[i] = stepYAML{
				Type:     "swipe_image",
				Img:      img,
				H:        v.H,
				W:        v.W,
				Duration: v.Duration,
				Wait:     v.Wait,
				Area:     area,
			}

		case *actions.ActionWait:
			img, err := encodeImage(v.ImageTemplate)
			if err != nil {
				return "", err
			}

			area := make([]int, 0)

			if v.SearchArea != nil {
				area = []int{v.SearchArea.Min.X, v.SearchArea.Min.Y, v.SearchArea.Max.X, v.SearchArea.Max.Y}
			}

			out[i] = stepYAML{
				Type:     "wait",
				Img:      img,
				Duration: v.Duration,
				Area:     area,
			}

		case *actions.ActionWaitStaticFrame:
			out[i] = stepYAML{Type: "wait_static_frame", Threshold: v.Threshold}

		case *actions.ActionDelay:
			out[i] = stepYAML{Type: "delay", Duration: &v.Duration}

		default:
			return "", fmt.Errorf("%w: %q", ErrUnknownStep, s)
		}
	}

	b, err := yaml.Marshal(out)

	return string(b), err
}

func (f *Flow) FromYAML(s string) error {
	var raw []stepYAML
	if err := yaml.Unmarshal([]byte(s), &raw); err != nil {
		return err
	}

	steps := make([]FlowStep, len(raw))

	for i := range raw {
		r := raw[i]

		switch r.Type {
		case "tap_xy":
			steps[i] = ActionTapXY(r.X, r.Y)

		case "swipe":
			steps[i] = ActionSwipe(r.X, r.Y, r.X2, r.Y2)

		case "keyboard":
			steps[i] = ActionKeyboard(r.Keys, *r.Duration)

		case "set_clipboard":
			steps[i] = ActionSetClipboard(r.Payload, r.Paste)

		case "type":
			steps[i] = ActionType(r.Payload)

		case "tap_image":
			img, err := decodeImage(r.Img)
			if err != nil {
				return err
			}

			var area *image.Rectangle

			if len(r.Area) == 4 {
				rect := image.Rect(r.Area[0], r.Area[1], r.Area[2], r.Area[3])
				area = &rect
			}

			var dur *time.Duration

			steps[i] = ActionTapImage(img, r.Wait, area, dur)

		case "swipe_image":
			img, err := decodeImage(r.Img)
			if err != nil {
				return err
			}

			var area *image.Rectangle

			if len(r.Area) == 4 {
				rect := image.Rect(r.Area[0], r.Area[1], r.Area[2], r.Area[3])
				area = &rect
			}

			var dur *time.Duration

			steps[i] = ActionSwipeImage(img, r.H, r.W, r.Wait, area, dur)

		case "wait":
			img, err := decodeImage(r.Img)
			if err != nil {
				return err
			}

			var area *image.Rectangle

			if len(r.Area) == 4 {
				rect := image.Rect(r.Area[0], r.Area[1], r.Area[2], r.Area[3])
				area = &rect
			}

			dur := r.Duration
			steps[i] = ActionWait(img, area, dur)

		case "wait_static_frame":
			steps[i] = WaitStaticFrame(r.Threshold)

		case "delay":
			steps[i] = ActionDelay(*r.Duration)

		default:
			return fmt.Errorf("%w: %q", ErrUnknownStep, r.Type)
		}
	}

	f.steps = steps

	return nil
}

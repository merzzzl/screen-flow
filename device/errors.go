package device

import "errors"

var (
	ErrNoABG    = errors.New("accessibility bridge not initialize")
	ErrNoSCRCPY = errors.New("scrcpy not initialize")
	ErrNoVision = errors.New("vision not initialize")
)

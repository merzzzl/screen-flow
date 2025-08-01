package actions

import "errors"

var ErrNoClints = errors.New("action supported client not found")
var ErrImageNotFound = errors.New("image not found")

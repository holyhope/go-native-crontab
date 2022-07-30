//go:build !darwin && linux && !freebsd && !netbsd && !openbsd && !windows && !js
// +build !darwin,linux,!freebsd,!netbsd,!openbsd,!windows,!js

package systemd

import (
	"github.com/holyhope/god"
)

func init() {
	god.New = New
}

func Support() bool {
	return true
}

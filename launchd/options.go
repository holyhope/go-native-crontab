package launchd

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"time"

	"github.com/holyhope/god"
)

func IsValidName(value string) bool {
	return fs.ValidPath(value) && len(value) > 0 && value != "." && value != ".."
}

func IsValidType(value string) bool {
	return true
}
func IsValidProgram(value string) bool {
	return true
}

func IsValidProgramArguments(value []string) bool {
	return true
}

func IsValidInterval(value time.Duration) bool {
	return true
}

func Path(scope god.ScopeValue, filename string) (string, error) {
	switch scope {
	case god.ScopeUser:
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("get user home: %w", err)
		}

		return path.Join(home, "Library/LaunchAgents", filename), nil
	case god.ScopeSystem:
		return path.Join("/Library/LaunchAgents", filename), nil
	default:
		return "", &god.InvalidOptionError{
			Key:   god.Scope,
			Value: scope,
		}
	}
}

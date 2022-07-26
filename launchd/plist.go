//go:build darwin && !linux && !freebsd && !netbsd && !openbsd && !windows && !js
// +build darwin,!linux,!freebsd,!netbsd,!openbsd,!windows,!js

package launchd

// // See file:///Applications/Xcode.app/Contents/Developer//Platforms/MacOSX.platform/Developer/SDKs/MacOSX.sdk/usr/include/launch.h
// #include <launch.h>
import "C"

type resourceLimits struct {
	Core          int64 `plist:"Core"`
	Memory        int64 `plist:"MemoryLimit"`
	NumberOfFiles int64 `plist:"NumberOfFiles"`
}

// See `man 5 launchd.plist` for more information
type launchUnit map[string]interface{}

func (u launchUnit) Label(value string) {
	u[C.LAUNCH_JOBKEY_LABEL] = value
}

func (u launchUnit) Program(value string) {
	u[C.LAUNCH_JOBKEY_PROGRAM] = value
}

func (u launchUnit) ProgramArguments(value ...string) {
	if len(value) == 0 {
		delete(u, C.LAUNCH_JOBKEY_PROGRAMARGUMENTS)
		return
	}

	u[C.LAUNCH_JOBKEY_PROGRAMARGUMENTS] = value
}

func (u launchUnit) RunAtLoad(value bool) {
	u[C.LAUNCH_JOBKEY_RUNATLOAD] = value
}

func (u launchUnit) StartInterval(value int) {
	u[C.LAUNCH_JOBKEY_STARTINTERVAL] = value
}

func (u launchUnit) KeepAlive(value bool) {
	u[C.LAUNCH_JOBKEY_KEEPALIVE] = value
}

func (u launchUnit) SoftResourceLimits(value *resourceLimits) {
	if value == nil {
		delete(u, C.LAUNCH_JOBKEY_SOFTRESOURCELIMITS)
		return
	}

	u[C.LAUNCH_JOBKEY_SOFTRESOURCELIMITS] = value
}

func (u launchUnit) Username(value string) {
	u[C.LAUNCH_JOBKEY_USERNAME] = value
}

func (u launchUnit) Groupname(value string) {
	u[C.LAUNCH_JOBKEY_GROUPNAME] = value
}

func (u launchUnit) HardResourceLimits(value *resourceLimits) {
	if value == nil {
		delete(u, C.LAUNCH_JOBKEY_HARDRESOURCELIMITS)
		return
	}

	u[C.LAUNCH_JOBKEY_HARDRESOURCELIMITS] = value
}

func (u launchUnit) StandardOutPath(value string) {
	u[C.LAUNCH_JOBKEY_STANDARDOUTPATH] = value
}

func (u launchUnit) StandardErrorPath(value string) {
	u[C.LAUNCH_JOBKEY_STANDARDERRORPATH] = value
}

func (u launchUnit) WorkingDirectory(value string) {
	u[C.LAUNCH_JOBKEY_WORKINGDIRECTORY] = value
}

func (u launchUnit) EnvironmentVariables(value map[string]string) {
	if value == nil {
		delete(u, C.LAUNCH_JOBKEY_ENVIRONMENTVARIABLES)
		return
	}

	u[C.LAUNCH_JOBKEY_ENVIRONMENTVARIABLES] = value
}

func (u launchUnit) LimitLoadSessionType(value string) {
	u[C.LAUNCH_JOBKEY_LIMITLOADTOSESSIONTYPE] = value
}

func (u launchUnit) ThrottleInterval(value int) {
	u[C.LAUNCH_JOBKEY_THROTTLEINTERVAL] = value
}

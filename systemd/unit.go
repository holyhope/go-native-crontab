package systemd

// // See file:///Applications/Xcode.app/Contents/Developer//Platforms/MacOSX.platform/Developer/SDKs/MacOSX.sdk/usr/include/launch.h
// #include <launch.h>
import "C"
import "strings"

// See `man 5 launchd.plist` for more information
type unit map[string]map[string]interface{}

func (u unit) AddToSection(section string, key string, value interface{}) {
	if _, ok := u[section]; !ok {
		u[section] = make(map[string]interface{})
	}

	u[section][key] = value
}

func (u unit) Description(value string) {
	u.AddToSection("Unit", "Description", value)
}

func (u unit) ExecStart(program string, arguments ...string) {
	b := new(strings.Builder)
	b.WriteString(program)
	for _, a := range arguments {
		b.WriteByte(' ')
		b.WriteString(a)
	}

	u.AddToSection("Service", "ExecStart", b.String())
}

func (u unit) RunAtLoad(value bool) {
	u[C.LAUNCH_JOBKEY_RUNATLOAD] = value
}

func (u unit) OnCalendar(value string) {
	u.AddToSection("Timer", "OnCalendar", value)
}

func (u unit) Name(names ...string) {
	u.AddToSection("Install", "Alias", strings.Join(names, ","))
}

func (u unit) SoftResourceLimits(value *resourceLimits) {
	if value == nil {
		delete(u, C.LAUNCH_JOBKEY_SOFTRESOURCELIMITS)
		return
	}

	u[C.LAUNCH_JOBKEY_SOFTRESOURCELIMITS] = value
}

func (u unit) Username(value string) {
	u[C.LAUNCH_JOBKEY_USERNAME] = value
}

func (u unit) Groupname(value string) {
	u[C.LAUNCH_JOBKEY_GROUPNAME] = value
}

func (u unit) HardResourceLimits(value *resourceLimits) {
	if value == nil {
		delete(u, C.LAUNCH_JOBKEY_HARDRESOURCELIMITS)
		return
	}

	u[C.LAUNCH_JOBKEY_HARDRESOURCELIMITS] = value
}

func (u unit) StandardOutPath(value string) {
	u[C.LAUNCH_JOBKEY_STANDARDOUTPATH] = value
}

func (u unit) StandardErrorPath(value string) {
	u[C.LAUNCH_JOBKEY_STANDARDERRORPATH] = value
}

func (u unit) WorkingDirectory(value string) {
	u[C.LAUNCH_JOBKEY_WORKINGDIRECTORY] = value
}

func (u unit) EnvironmentVariables(value map[string]string) {
	if value == nil {
		delete(u, C.LAUNCH_JOBKEY_ENVIRONMENTVARIABLES)
		return
	}

	u[C.LAUNCH_JOBKEY_ENVIRONMENTVARIABLES] = value
}

func (u unit) LimitLoadSessionType(value string) {
	u[C.LAUNCH_JOBKEY_LIMITLOADTOSESSIONTYPE] = value
}

func (u unit) ThrottleInterval(value int) {
	u.AddToSection("Service", "StartLimitInterval", value)
}

func (u unit) WatchPaths(value ...string) {
	u[C.LAUNCH_JOBKEY_WATCHPATHS] = value
}

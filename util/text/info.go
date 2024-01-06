package text

import "runtime/debug"

var Commit = func() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				return setting.Value[len(setting.Value)-7:]
			}
		}
	}
	return ""
}()

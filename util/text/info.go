package text

import (
	"runtime/debug"
)

var Commit = func() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		hash := "none"
		modified := "/u" // Unknown build
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				hash = setting.Value[:7]
			}
			if setting.Key == "vcs.modified" {
				if setting.Value == "false" {
					modified = "/c" // Clean build
				} else {
					modified = "/d" // Dirty build
				}
			}
		}
		return hash + modified + "/" + Build
	}
	return ""
}()

package common

import "fmt"

func FormatVersion(version uint32) string {
	if version == 0 {
		return "0.0.0"
	}
	major := version << 8
	major = major >> 24

	minor := version << 16
	minor = minor >> 24

	patch := version << 24
	patch = patch >> 24

	return fmt.Sprintf("%d.%d.%d", major, minor, patch)
}

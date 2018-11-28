package utils

import "runtime"

// GetPathSeparator returns the path separator based on the OS type
func GetPathSeparator() string {
	if runtime.GOOS == WindowsOs {
		return WindowsPathSeparator
	}
	return UnixPathSeparator
}

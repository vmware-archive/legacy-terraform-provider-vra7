package utils

import "runtime"

func GetPathSeparator() string {
	if runtime.GOOS == WINDOWS_OS {
		return WINDOWS_PATH_SEPARATOR
	} else {
		return UNIX_PATH_SEPARATOR
	}
}

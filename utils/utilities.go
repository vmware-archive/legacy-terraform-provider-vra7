package utils

import (
	"bytes"
	"encoding/json"
	"runtime"
)

// GetPathSeparator returns the path separator based on the OS type
func GetPathSeparator() string {
	if runtime.GOOS == WindowsOs {
		return WindowsPathSeparator
	}
	return UnixPathSeparator
}

// Unmarshal - decodes json
func UnmarshalJSON(data []byte, v interface{}) error {
	err := json.Unmarshal(data, v)
	if err != nil {
		return err
	}
	return nil
}

// Marshal the object to JSON and convert to *bytes.Buffer
func MarshalToJSON(v interface{}) (*bytes.Buffer, error) {
	buffer := new(bytes.Buffer)
	err := json.NewEncoder(buffer).Encode(v)
	if err != nil {
		return nil, err
	}
	return buffer, nil
}

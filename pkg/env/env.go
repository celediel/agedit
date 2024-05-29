package env

import (
	"os"
	"runtime"
)

// GetEditor gets the configured editor by checking environmental
// variables EDITOR and VISUAL
func GetEditor() string {
	var editor string
	if v, ok := os.LookupEnv("EDITOR"); ok {
		editor = v
	} else if v, ok := os.LookupEnv("VISUAL"); ok {
		editor = v
	} /* else {
		// TODO: maybe pick something based on the OS
	} */

	return editor
}

// GetTempDirectory returns the systems temporary directory
//
// returns %TEMP% on Windows, /tmp on UNIX-like systems
func GetTempDirectory() string {
	var tmp string
	switch runtime.GOOS {
	case "windows":
		tmp = os.Getenv("TEMP")
	case "android":
		if t := os.Getenv("TMPDIR"); t != "" {
			tmp = t
		} else if t = os.Getenv("PREFIX"); t != "" {
			tmp = t + "/tmp"
		}
	default:
		fallthrough
	case "darwin":
		fallthrough
	case "linux":
		tmp = "/tmp"
	}
	return tmp
}

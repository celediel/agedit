package env

import (
	"os"
	"regexp"
	"runtime"
	"strings"

	"github.com/adrg/xdg"
)

// GetEditor gets the configured editor by checking environmental
// variables EDITOR and VISUAL
func GetEditor() string {
	var editor string
	if os.Getenv("EDITOR") != "" {
		editor = os.Getenv("EDITOR")
	} else if os.Getenv("VISUAL") != "" {
		editor = os.Getenv("VISUAL")
	} /* else {
		// TODO: maybe pick something based on the OS
	} */

	return editor
}

// GetConfigDir gets a config directory based from environmental variables + the app name
//
// On Windows, %APPDATA%\agedit is used
//
// On UNIX-like systems, $XDG_CONFIG_HOME/agedit is tried, if it isn't defined, $HOME/.config/agedit is used
func GetConfigDir(appname string) string {
	return make_path(xdg.ConfigHome, appname)
}

// GetConfigDir gets a config directory based from environmental variables + the app name
//
// On Windows, %LOCALAPPDATA%\agedit is used
//
// On UNIX-like systems, $XDG_DATA_HOME/agedit is tried, if it isn't defined, $HOME/.local/share/agedit is used
func GetDataDir(appname string) string {
	return make_path(xdg.DataHome, appname)
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

func make_path(paths ...string) string {
	sep := string(os.PathSeparator)
	output := strings.Builder{}

	// add / to the start if it's not already there and we're not on Windows
	if match, err := regexp.Match("^\\w", []byte(paths[0])); err == nil && match && runtime.GOOS != "windows" {
		output.WriteString(sep)
	}

	for _, path := range paths {
		// don't add / to the end if it's there
		if strings.HasSuffix(path, sep) {
			output.WriteString(path)
		} else {
			output.WriteString(path + sep)
		}
	}

	return output.String()
}

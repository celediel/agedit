package env

import (
	"os"
	"regexp"
	"runtime"
	"strings"
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
	var configdir string
	switch runtime.GOOS {
	case "windows":
		configdir = os.Getenv("APPDATA")
	default:
		fallthrough
	case "darwin":
		// TODO: figure out the proper Mac OS local directories
		fallthrough
	case "linux":
		if confighome := os.Getenv("XDG_CONFIG_HOME"); confighome != "" {
			configdir = confighome
		} else {
			configdir = make_path(os.Getenv("HOME"), ".config")
		}
	}

	return make_path(configdir, appname)
}

// GetConfigDir gets a config directory based from environmental variables + the app name
//
// On Windows, %LOCALAPPDATA%\agedit is used
//
// On UNIX-like systems, $XDG_DATA_HOME/agedit is tried, if it isn't defined, $HOME/.local/share/agedit is used
func GetDataDir(appname string) string {
	var datadir string
	switch runtime.GOOS {
	case "windows":
		datadir = os.Getenv("LOCALAPPDATA")
	default:
		fallthrough
	case "darwin":
		// TODO: also here
		fallthrough
	case "linux":
		if datahome := os.Getenv("XDG_DATA_HOME"); datahome != "" {
			datadir = datahome
		} else {
			datadir = make_path(os.Getenv("HOME"), "local", "share")
		}
	}

	return make_path(datadir, appname)
}

// GetTempDirectory returns the systems temporary directory
//
// returns %TEMP% on Windows, /tmp on UNIX-like systems
func GetTempDirectory() string {
	switch runtime.GOOS {
	case "windows":
		return os.Getenv("TEMP")
	default:
		fallthrough
	case "darwin":
		fallthrough
	case "linux":
		return "/tmp"
	}
}

func make_path(paths ...string) string {
	sep := string(os.PathSeparator)
	output := strings.Builder{}

	// add / to the start if it's not already there and we're not on Windows
	if match, err := regexp.Match("^\\w", []byte(paths[0])); err == nil && match && runtime.GOOS != "windows" {
		output.WriteString(sep)
	}

	for _, path := range paths {
		output.WriteString(path + sep)
	}

	return output.String()
}

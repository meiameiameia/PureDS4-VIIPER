package configpaths

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
)

// DefaultConfigDir returns the platform-specific configuration directory for VIIPER.
func DefaultConfigDir() (string, error) {
	switch runtime.GOOS {
	case "windows":
		if appdata := os.Getenv("AppData"); appdata != "" {
			return filepath.Join(appdata, "VIIPER"), nil
		}
		return "", errors.New("AppData not set")
	default:
		if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
			return filepath.Join(xdg, "github.com/Alia5/viiper"), nil
		}
		if home := os.Getenv("HOME"); home != "" {
			return filepath.Join(home, ".config", "github.com/Alia5/viiper"), nil
		}
		return "", errors.New("HOME not set")
	}
}

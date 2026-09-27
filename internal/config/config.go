// Package config defines the CLI structure and configuration for VIIPER.
package config

import (
	"github.com/Alia5/VIIPER/internal/cmd"
)

type UpdateNotify string

const (
	UpdateNotifyNone       UpdateNotify = "none"
	UpdateNotifyStable     UpdateNotify = "stable"
	UpdateNotifyPrerelease UpdateNotify = "prerelease"
)

type Log struct {
	Level   string `aliases:"l" help:"Log level: trace, debug, info, warn, error" default:"info" env:"VIIPER_LOG_LEVEL"`
	File    string `help:"Log file path (default: none; logs only to console)" env:"VIIPER_LOG_FILE"`
	RawFile string `help:"Raw packet log file path (default: none)" env:"VIIPER_LOG_RAW_FILE"`
}

// CLI is the root command structure for Kong CLI parsing.
type CLI struct {
	// Global
	UpdateNotify UpdateNotify `help:"Deprecated and ignored; PureDS4 controls backend updates" default:"none" env:"VIIPER_UPDATE_NOTIFY"`
	Log          `embed:"" prefix:"log."`

	Server cmd.Server `cmd:"" help:"Start the VIIPER USB-IP server" default:""`
}

// Preserve compatibility with old configuration while making all update
// notification values inert. PureDS4 owns replacement of the backend binary.
func (c *CLI) AfterApply() error {
	c.UpdateNotify = UpdateNotifyNone
	return nil
}

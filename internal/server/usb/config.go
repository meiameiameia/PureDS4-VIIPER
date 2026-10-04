package usb

import "time"

// ServerConfig represents the server subcommand configuration.
type ServerConfig struct {
	Addr                    string        `help:"USB-IP listen address (127.0.0.1 only)" default:"127.0.0.1:3241" env:"VIIPER_USB_ADDR"`
	ConnectionTimeout       time.Duration `kong:"-"`
	BusCleanupTimeout       time.Duration `help:"-"`
	WriteBatchFlushInterval time.Duration `default:"0" help:"Interval to flush write batches to clients; default: disabled / immediate updates" env:"VIIPER_USB_WRITE_BATCH_FLUSH_INTERVAL"`
}

package main

import (
	"context"
	"io"
	"log/slog"
	"os"
	"testing"

	"github.com/Alia5/VIIPER/internal/config"
	"github.com/Alia5/VIIPER/internal/server/api"
	"github.com/alecthomas/kong"
	"github.com/stretchr/testify/require"
)

func TestPureDS4BackendCannotEnableInheritedUpdateNotifications(t *testing.T) {
	for _, mode := range []string{"stable", "prerelease"} {
		t.Run(mode, func(t *testing.T) {
			var cli config.CLI
			parser, err := kong.New(&cli)
			require.NoError(t, err)
			_, err = parser.Parse([]string{"--update-notify=" + mode, "server"})
			require.NoError(t, err)
			require.Equal(t, config.UpdateNotifyNone, cli.UpdateNotify)
		})
	}
}

func TestPureDS4BackendDefaultsToLocalListeners(t *testing.T) {
	for _, variable := range []string{"VIIPER_USB_ADDR", "VIIPER_API_ADDR"} {
		value, exists := os.LookupEnv(variable)
		require.NoError(t, os.Unsetenv(variable))
		t.Cleanup(func() {
			if exists {
				_ = os.Setenv(variable, value)
			} else {
				_ = os.Unsetenv(variable)
			}
		})
	}
	var cli config.CLI
	parser, err := kong.New(&cli)
	require.NoError(t, err)
	_, err = parser.Parse([]string{"server"})
	require.NoError(t, err)
	require.Equal(t, "127.0.0.1:3241", cli.Server.USBServerConfig.Addr)
	require.Equal(t, "127.0.0.1:3242", cli.Server.APIServerConfig.Addr)
}

func TestPureDS4BackendRejectsUnsafeListenerOverrides(t *testing.T) {
	for _, endpoint := range []string{"usb", "api"} {
		for _, addr := range []string{":0", "0.0.0.0:0", "[::]:0", "192.0.2.1:0", "example.com:0"} {
			t.Run(endpoint+"/"+addr, func(t *testing.T) {
				t.Setenv("VIIPER_USB_ADDR", "127.0.0.1:0")
				t.Setenv("VIIPER_API_ADDR", "127.0.0.1:0")
				var cli config.CLI
				parser, err := kong.New(&cli)
				require.NoError(t, err)
				_, err = parser.Parse([]string{"server", "--" + endpoint + ".addr=" + addr})
				require.ErrorContains(t, err, "requires listen host")
				// Programmatic startup also rejects the same configuration before
				// probing USB/IP, writing keys, or listening on a socket.
				err = cli.Server.StartServer(context.Background(), slog.New(slog.NewTextHandler(io.Discard, nil)), nil)
				require.ErrorContains(t, err, "requires listen host")
			})
		}
	}
}

func TestPureDS4BackendRejectsUnsafeListenerEnvironment(t *testing.T) {
	for _, variable := range []string{"VIIPER_USB_ADDR", "VIIPER_API_ADDR"} {
		t.Run(variable, func(t *testing.T) {
			t.Setenv("VIIPER_USB_ADDR", "127.0.0.1:0")
			t.Setenv("VIIPER_API_ADDR", "127.0.0.1:0")
			t.Setenv(variable, "0.0.0.0:3241")
			var cli config.CLI
			parser, err := kong.New(&cli)
			require.NoError(t, err)
			_, err = parser.Parse([]string{"server"})
			require.ErrorContains(t, err, "requires listen host")
		})
	}
}

func TestPureDS4BackendRegistersOnlyRequiredGamepads(t *testing.T) {
	for _, deviceType := range []string{
		"dualshock4", "dualshock4audioduplexv3",
		"dualshock4audioonlyduplexv3", "xbox360",
	} {
		require.NotNil(t, api.GetRegistration(deviceType), deviceType)
	}
	for _, deviceType := range []string{
		"dualsense", "dualsenseedge", "keyboard", "mouse", "ns2pro",
	} {
		require.Nil(t, api.GetRegistration(deviceType), deviceType)
	}
}

func TestPureDS4BackendDoesNotOfferHostManagementCommands(t *testing.T) {
	for _, command := range []string{"install", "uninstall", "proxy", "config", "codegen"} {
		t.Run(command, func(t *testing.T) {
			var cli config.CLI
			parser, err := kong.New(&cli)
			require.NoError(t, err)
			_, err = parser.Parse([]string{command})
			require.Error(t, err)
		})
	}
}

func TestPureDS4BackendDoesNotReadInheritedConfigFiles(t *testing.T) {
	var cli config.CLI
	parser, err := kong.New(&cli)
	require.NoError(t, err)
	_, err = parser.Parse([]string{"--config=server.yaml", "server"})
	require.Error(t, err)
}

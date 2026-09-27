package main

import (
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

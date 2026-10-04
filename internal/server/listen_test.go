package server

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLocalListenAddress(t *testing.T) {
	for _, addr := range []string{"127.0.0.1:3241", "127.0.0.1:0", "localhost:3242"} {
		t.Run(addr, func(t *testing.T) {
			got, err := LocalListenAddress(addr)
			require.NoError(t, err)
			if addr == "localhost:3242" {
				require.Equal(t, "127.0.0.1:3242", got)
			} else {
				require.Equal(t, addr, got)
			}
		})
	}
	for _, addr := range []string{
		"", ":3241", "0.0.0.0:3241", "[::]:3241", "[::1]:3241",
		"192.0.2.1:3241", "example.com:3241", "127.0.0.2:3241",
		"127.0.0.1:http", "127.0.0.1:-1", "127.0.0.1:65536", "127.0.0.1:",
	} {
		t.Run(addr, func(t *testing.T) {
			_, err := LocalListenAddress(addr)
			require.Error(t, err)
		})
	}
}

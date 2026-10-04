//go:build windows

package api

import (
	"bytes"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/require"
)

func TestUsbipAttachABISizes(t *testing.T) {
	require.Equal(t, uintptr(1100), unsafe.Sizeof(attachIOCTL{}),
		"usbip-win2 0.9.7.7 plugin_hardware ABI changed")
}

func TestNativeAttachTargetsTheSameIPv4LoopbackAsTheListener(t *testing.T) {
	data := newLocalAttachIOCTL("1-2", "3241")
	require.Equal(t, uint32(1100), data.Size)
	require.Equal(t, "127.0.0.1", string(bytes.TrimRight(data.Host[:], "\x00")))
	require.Equal(t, "3241", string(bytes.TrimRight(data.Service[:], "\x00")))
	require.Equal(t, "1-2", string(bytes.TrimRight(data.BusID[:], "\x00")))
}

func TestNativeAutoAttachResultCarriesExactPort(t *testing.T) {
	got := AutoAttachResult{
		USBIPPort: 7,
	}

	require.Equal(t, int32(7), got.USBIPPort)
	require.Empty(t, got.USBIPOwnerSerial)
}

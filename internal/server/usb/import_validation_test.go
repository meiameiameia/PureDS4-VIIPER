package usb

import (
	"bytes"
	"encoding/binary"
	"io"
	"log/slog"
	"net"
	"testing"
	"time"

	"github.com/Alia5/VIIPER/usbip"
	"github.com/stretchr/testify/require"
)

func TestImportRejectsInvalidBusIDWithoutPanic(t *testing.T) {
	for _, tc := range []struct {
		name    string
		payload []byte
		want    string
	}{
		{"unterminated", bytes.Repeat([]byte{'x'}, busIDSize), "invalid import busid"},
		{"empty", make([]byte, busIDSize), "invalid import busid"},
		{"truncated", []byte("1-1"), "read import busid"},
		{"unknown", append([]byte("1-1"), make([]byte, busIDSize-3)...), "no device matches busid 1-1"},
		{"maximum terminated length", append(bytes.Repeat([]byte{'x'}, busIDSize-1), 0), "no device matches busid"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			serverConn, clientConn := net.Pipe()
			defer serverConn.Close()
			defer clientConn.Close()
			require.NoError(t, serverConn.SetDeadline(time.Now().Add(2*time.Second)))
			require.NoError(t, clientConn.SetDeadline(time.Now().Add(2*time.Second)))
			written := make(chan error, 1)
			go func() {
				_, err := clientConn.Write(tc.payload)
				_ = clientConn.Close()
				written <- err
			}()
			s := New(ServerConfig{}, slog.New(slog.NewTextHandler(io.Discard, nil)), nil)
			require.NotPanics(t, func() {
				dev, err := s.handleImport(serverConn)
				require.Nil(t, dev)
				require.ErrorContains(t, err, tc.want)
			})
			require.NoError(t, <-written)
		})
	}
}

func TestUSBServerServesValidRequestAfterMalformedImport(t *testing.T) {
	s := New(ServerConfig{Addr: "localhost:0", ConnectionTimeout: 2 * time.Second},
		slog.New(slog.NewTextHandler(io.Discard, nil)), nil)
	done := make(chan error, 1)
	go func() { done <- s.ListenAndServe() }()
	t.Cleanup(func() {
		_ = s.Close()
		select {
		case err := <-done:
			require.NoError(t, err)
		case <-time.After(3 * time.Second):
			t.Error("USB test listener did not stop")
		}
	})
	select {
	case <-s.Ready():
	case err := <-done:
		t.Fatalf("USB listener failed: %v", err)
	case <-time.After(3 * time.Second):
		t.Fatal("USB listener did not start")
	}
	host, _, err := net.SplitHostPort(s.Addr())
	require.NoError(t, err)
	require.Equal(t, "127.0.0.1", host)
	for _, malformed := range [][]byte{
		bytes.Repeat([]byte{'x'}, busIDSize), make([]byte, busIDSize), []byte("1-1"),
	} {
		conn, err := net.DialTimeout("tcp4", s.Addr(), 2*time.Second)
		require.NoError(t, err)
		require.NoError(t, conn.SetDeadline(time.Now().Add(2*time.Second)))
		var request bytes.Buffer
		require.NoError(t, (&usbip.MgmtHeader{Version: usbip.Version, Command: usbip.OpReqImport}).Write(&request))
		request.Write(malformed)
		_, err = conn.Write(request.Bytes())
		require.NoError(t, err)
		// A partial request must signal EOF; no driver or live backend is involved.
		require.NoError(t, conn.(*net.TCPConn).CloseWrite())
		var response [1]byte
		_, err = conn.Read(response[:])
		require.Error(t, err)
		if netErr, ok := err.(net.Error); ok {
			require.False(t, netErr.Timeout(), "malformed connection must close, not hang")
		}
		require.NoError(t, conn.Close())
	}
	conn, err := net.DialTimeout("tcp4", s.Addr(), 2*time.Second)
	require.NoError(t, err)
	defer conn.Close()
	require.NoError(t, conn.SetDeadline(time.Now().Add(2*time.Second)))
	require.NoError(t, (&usbip.MgmtHeader{Version: usbip.Version, Command: usbip.OpReqDevlist}).Write(conn))
	var response [12]byte // management header plus zero-device count
	_, err = io.ReadFull(conn, response[:])
	require.NoError(t, err)
	require.Equal(t, uint16(usbip.OpRepDevlist), binary.BigEndian.Uint16(response[2:4]))
	require.Equal(t, uint32(0), binary.BigEndian.Uint32(response[4:8]))
	require.Equal(t, uint32(0), binary.BigEndian.Uint32(response[8:12]))
}

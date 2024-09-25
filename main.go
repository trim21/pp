//go:build linux

package main

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"os"
	"syscall"

	"github.com/things-go/go-socks5"
	"golang.org/x/sys/unix"
)

func main() {
	log := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	port := os.Getenv("PP_SOCKS5_PORT")
	if port == "" {
		port = "1080"
	}

	lc := net.ListenConfig{}

	l, err := lc.Listen(context.Background(), "tcp", ":"+port)
	if err != nil {
		panic(err)
	}

	log.Info("start server")

	dialer := net.Dialer{
		ControlContext: func(ctx context.Context, network, address string, c syscall.RawConn) error {
			log.Debug("new connection", "network", network, "address", address)
			var innerErr error
			outerErr := c.Control(func(fd uintptr) {
				innerErr = syscall.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_MARK, int(uint32(0x163)))
			})

			err := errors.Join(outerErr, innerErr)
			if err != nil {
				log.Error("failed to set socket options", "err", err)
			}

			return err
		},
	}

	s := socks5.NewServer(
		socks5.WithDial(dialer.DialContext),
	)

	// Create SOCKS5 proxy on localhost port 8000
	for {
		if err := s.Serve(l); err != nil {
			log.Error("failed to accept new connection", "error", err)
		}
	}
}

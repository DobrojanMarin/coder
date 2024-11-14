//go:build darwin
// +build darwin

package main

import "C"

import (
	"context"
	"io"
	"os"

	"golang.org/x/sys/unix"
	"golang.org/x/xerrors"

	"cdr.dev/slog"
	"cdr.dev/slog/sloggers/sloghuman"
	"github.com/coder/coder/v2/vpn"
)

// OpenTunnel creates a new VPN tunnel by `dup`ing the provided 'PIPE'
// file descriptors for reading, writing, and logging.
//
//export OpenTunnel
func OpenTunnel(cReadFD, cWriteFD, cLogFD int32) int32 {
	ctx := context.Background()

	conn, err := newFdConn(cReadFD, cWriteFD)
	if err != nil {
		return -1
	}

	logger, err := newFdLogger(cLogFD)
	if err != nil {
		return -1
	}

	_, err = vpn.NewTunnel(ctx, logger, conn)
	if err != nil {
		return -1
	}

	return 0
}

type pipeConn struct {
	r *os.File
	w *os.File
}

func (f *pipeConn) Read(p []byte) (n int, err error) {
	return f.r.Read(p)
}

func (f *pipeConn) Write(p []byte) (n int, err error) {
	return f.w.Write(p)
}

func (f *pipeConn) Close() error {
	_ = f.r.Close()
	_ = f.w.Close()
	return nil
}

func newFdConn(cReadFD, cWriteFD int32) (io.ReadWriteCloser, error) {
	readFD, err := unix.Dup(int(cReadFD))
	if err != nil {
		return nil, xerrors.Errorf("dup readFD: %w", err)
	}
	reader := os.NewFile(uintptr(readFD), "PIPE")
	if reader == nil {
		unix.Close(readFD)
		return nil, xerrors.New("failed to create reader")
	}

	writeFD, err := unix.Dup(int(cWriteFD))
	if err != nil {
		return nil, xerrors.Errorf("dup writeFD: %w", err)
	}
	writer := os.NewFile(uintptr(writeFD), "PIPE")
	if writer == nil {
		unix.Close(readFD)
		unix.Close(writeFD)
		return nil, xerrors.New("failed to create writer")
	}

	return &pipeConn{
		r: reader,
		w: writer,
	}, nil
}

func newFdLogger(cLogFD int32) (slog.Logger, error) {
	logFD, err := unix.Dup(int(cLogFD))
	if err != nil {
		return slog.Logger{}, xerrors.New("failed to dup logFD")
	}
	logFile := os.NewFile(uintptr(logFD), "PIPE")
	if logFile == nil {
		unix.Close(logFD)
		return slog.Logger{}, xerrors.New("failed to create log file")
	}
	return slog.Make(sloghuman.Sink(logFile)).Leveled(slog.LevelDebug), nil
}

func main() {}

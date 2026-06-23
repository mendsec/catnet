package cli

import (
	"context"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

const (
	ExitCodeSuccess     = 0
	ExitCodeInputError  = 1
	ExitCodeRuntimeError = 2
	ExitCodeInterrupted = 130
)

func WithCancelOnSignal(parent context.Context) (context.Context, context.CancelFunc) {
	signals := []os.Signal{os.Interrupt}
	if runtime.GOOS != "windows" {
		signals = append(signals, syscall.SIGTERM)
	}
	return signal.NotifyContext(parent, signals...)
}

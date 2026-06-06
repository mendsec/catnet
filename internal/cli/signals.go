package cli

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

const (
	// ExitCodeSuccess indicates successful execution
	ExitCodeSuccess = 0
	// ExitCodeInputError indicates an invalid input or flag
	ExitCodeInputError = 1
	// ExitCodeRuntimeError indicates a failure during execution
	ExitCodeRuntimeError = 2
	// ExitCodeInterrupted indicates the process was cancelled by a signal
	ExitCodeInterrupted = 130
)

// WithCancelOnSignal retorna um contexto que é cancelado quando SIGINT ou SIGTERM é recebido.
// O caller deve checar ctx.Err() após o retorno da operação para distinguir cancelamento por
// sinal de cancelamento por timeout.
func WithCancelOnSignal(parent context.Context) (context.Context, context.CancelFunc) {
	return signal.NotifyContext(parent, os.Interrupt, syscall.SIGTERM)
}

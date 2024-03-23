package context

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/bdlm/log/v2"
)

// New returns a signal aware root context for services that listens for syscall
// signals. Additional signals can be passed in.
func New(additionalSignals ...os.Signal) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())

	// Create the application context and start a signal handler.
	go func() {
		interrupt := make(chan os.Signal, 1)
		signal.Notify(
			interrupt,
			append(additionalSignals, syscall.SIGTERM)..., // docker, docker-compose, k8s, etc.
		)

		defer cancel()
		sig := <-interrupt
		log.WithField("signal", sig.String()).Info("signal received, shutting down")
	}()

	return ctx, cancel
}

package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
	"os"
	"os/signal"
	"syscall"
)

func (app *application) serve() error {
	srv := &http.Server{
		Addr:			fmt.Sprintf(":%d", app.config.port),
		Handler:		app.routes(),
		IdleTimeout:	time.Minute,
		ReadTimeout:	10 * time.Second,
		WriteTimeout:	30 * time.Second,
	}

	shutdownError := make(chan error)

	// background goroutine
	go func() {
		// create a quit channel which carries os.Signal values
		quit := make(chan os.Signal, 1)

		// signal.Notify() to listen for incoming SIGINT and SIGTERM signals
		// and relay them to the quit channel. Any other signals will not be caught by
		// signal.Notify() and will retain their default behavior
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		// read the signal from the quit channel
		s := <-quit

		app.logger.PrintInfo("shutting down server", map[string]string{
			"signal": s.String(),
		})

		ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
		defer cancel()

		// Exit the application with a 0 (success) status code
		//os.Exit(0)

		shutdownError <- srv.Shutdown(ctx)
	}()

	app.logger.PrintInfo("starting server", map[string]string{
		"addr": srv.Addr,
		"env": app.config.env,
	})

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	app.logger.PrintInfo("stopped server", map[string]string{
		"add": srv.Addr,
	})

	return nil
}
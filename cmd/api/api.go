package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/janakerman/flux-event-store/internal/server"
	"github.com/sirupsen/logrus"
)

var serverAddr string

func SetupSignalHandler() context.Context {
	ctx, cancel := context.WithCancel(context.Background())

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cancel()
		<-c
		os.Exit(1) // second signal. Exit directly.
	}()

	return ctx
}

func main() {
	flag.StringVar(&serverAddr, "serverAddr", ":8080", "The address the event server binds to.")

	ctx := SetupSignalHandler()
	srv := server.NewEventServer(serverAddr, logrus.New())
	srv.ListenAndServe(ctx.Done())
}

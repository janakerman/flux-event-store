package server

import (
    "context"
    "net/http"
    "os"
    "time"

    "github.com/sirupsen/logrus"
)

type EventServer struct {
    port       string
    logger     *logrus.Logger
}

func NewEventServer(port string, logger *logrus.Logger) *EventServer {
    return &EventServer{
        port:       port,
        logger:     logger,
    }
}

// ListenAndServe starts the HTTP server on the specified port
func (s *EventServer) ListenAndServe(stopCh <-chan struct{}) {
    mux := http.DefaultServeMux
    mux.Handle("/health", http.HandlerFunc(s.handleHealth))
    mux.Handle("/notification", http.HandlerFunc(s.handleNotification()))
    srv := &http.Server{
        Addr:    s.port,
        Handler: http.HandlerFunc(mux.ServeHTTP),
    }

    go func() {
        if err := srv.ListenAndServe(); err != http.ErrServerClosed {
            s.logger.Error(err, "Event server crashed")
            os.Exit(1)
        }
    }()

    // wait for SIGTERM or SIGINT
    <-stopCh
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
        s.logger.Error(err, "Event server graceful shutdown failed")
    } else {
        s.logger.Info("Event server stopped")
    }
}

func (s *EventServer) handleNotification() func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {

    }
}

func (s *EventServer) handleHealth(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(200)
}

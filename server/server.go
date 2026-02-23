package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/kofuk/spaghettini/server/backend"
)

type Server struct {
	logger *slog.Logger
	server *http.Server
}

type LogLevel string

const (
	Debug LogLevel = "debug"
	Info  LogLevel = "info"
	Warn  LogLevel = "warn"
	Error LogLevel = "error"
)

func (l LogLevel) ToSlogLevel() slog.Level {
	switch l {
	case Debug:
		return slog.LevelDebug
	case Info:
		return slog.LevelInfo
	case Warn:
		return slog.LevelWarn
	case Error:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

type ServerOptions struct {
	LogLevel LogLevel
	Addr     string
	Source   string
}

func NewServer(options ServerOptions) (*Server, error) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: options.LogLevel.ToSlogLevel(),
	}))

	backend, err := backend.NewGoTemplateBackend(options.Source)
	if err != nil {
		return nil, err
	}

	server := &http.Server{
		Addr: options.Addr,
		Handler: &Handler{
			logger: logger,
			template: &Evaluator{
				backend:       backend,
				printResponse: options.LogLevel == Debug,
			},
		},
	}

	return &Server{
		logger: logger,
		server: server,
	}, nil
}

func (s *Server) Start(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		s.logger.Info("shutting down server")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := s.server.Shutdown(ctx); err != nil {
			s.logger.Error("failed to shutdown server", "error", err)
		}
	}()

	s.logger.Info(fmt.Sprintf("Listening on %s", s.server.Addr))
	if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}

	return nil
}

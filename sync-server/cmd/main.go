package main

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"time"

	"sync-server/internal/config"
	"sync-server/internal/handler"
)

func main() {
	cfg := config.MustLoad()
	slog.Info("Config loaded", slog.String("config", fmt.Sprintf("%+v", cfg)))

	err := setupLogging(cfg)
	if err != nil {
		slog.Error("Error setting up logging", slog.String("err", err.Error()))

		return
	}

	slog.Info("Logging setup")

	listener, err := net.Listen("tcp", net.JoinHostPort(cfg.Host, cfg.Port))
	if err != nil {
		slog.Error("Error starting server", slog.String("err", err.Error()))

		return
	}
	defer listener.Close()

	slog.Info("Server started", slog.String("addr", listener.Addr().String()))

	for {
		conn, err := listener.Accept()
		if err != nil {
			slog.Warn("Error accepting connection", slog.String("err", err.Error()))
		}

		slog.Info("Client connected", slog.String("addr", conn.RemoteAddr().String()))

		err = conn.SetDeadline(time.Now().Add(cfg.ReadWriteTimeoutDuration))
		if err != nil {
			slog.Warn("Error setting read and write deadline", slog.String("err", err.Error()))
		}

		handler.HandleConnection(conn, *cfg)
	}
}

func setupLogging(cfg *config.Config) error {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{})))

	switch cfg.LogLevel {
	case "error":
		slog.SetLogLoggerLevel(slog.LevelError)
	case "warn":
		slog.SetLogLoggerLevel(slog.LevelWarn)
	case "info":
		slog.SetLogLoggerLevel(slog.LevelInfo)
	case "debug":
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	return nil
}

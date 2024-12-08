package main

import (
	"fmt"
	"log/slog"
	"net"
	"os"

	"sync-server/internal/config"
	"sync-server/internal/handler"
)

func main() {
	cfg := config.MustLoad()

	if err := setupLogging(cfg); err != nil {
		slog.Error("Error setting up logging", slog.String("err", err.Error()))

		return
	}

	slog.Info("Config loaded", slog.String("config", fmt.Sprintf("%+v", cfg)))

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
		handler.HandleConnection(conn, *cfg)
	}
}

func setupLogging(cfg *config.Config) error {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{})))

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

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
			slog.Error("Error accepting connection", slog.String("err", err.Error()))
		}

		slog.Info("Client connected", slog.String("addr", conn.RemoteAddr().String()))

		err = conn.SetReadDeadline(time.Now().Add(cfg.ReadTimeoutDuration))
		if err != nil {
			slog.Error("Error setting read deadline", slog.String("err", err.Error()))
		}

		err = conn.SetWriteDeadline(time.Now().Add(cfg.WriteTimeoutDuration))
		if err != nil {
			slog.Error("Error setting write deadline", slog.String("err", err.Error()))
		}

		go handler.HandleConnection(conn, *cfg)
	}
}

func setupLogging(cfg *config.Config) error {
	fileName := fmt.Sprintf("logs/log_%s.txt", time.Now().Format("2006-01-02_15-04-05"))

	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("cannot create log file: %w", err)
	}

	slog.SetDefault(slog.New(slog.NewTextHandler(file, &slog.HandlerOptions{})))

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

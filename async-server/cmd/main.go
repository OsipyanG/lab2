package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"net"
	"os"
	"sync-server/internal/config"
	"sync-server/internal/handler"
	"time"
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
		defer conn.Close()

		slog.Info("Client connected", slog.String("addr", conn.RemoteAddr().String()))

		conn.SetReadDeadline(time.Now().Add(cfg.ReadTimeoutDuration))
		conn.SetWriteDeadline(time.Now().Add(cfg.WriteTimeoutDuration))

		go handler.HandleConnection(conn, *cfg)

		slog.Info("Client disconnected", slog.String("addr", conn.RemoteAddr().String()))
	}

}

func setupLogging(cfg *config.Config) error {
	fileName := fmt.Sprintf("logs/log_%s.txt", time.Now().Format("2006-01-02_15-04-05"))

	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("cannot create log file: %w", err)
	}

	fileBufferWriter := bufio.NewWriter(file)

	slog.SetDefault(slog.New(slog.NewTextHandler(fileBufferWriter, &slog.HandlerOptions{})))

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

package handler

import (
	"bufio"
	"errors"
	"log/slog"
	"net"
	"os"
	"strings"
	"time"

	"sync-server/internal/config"
)

func HandleConnection(conn net.Conn, cfg config.Config) {
	defer func() {
		err := conn.Close()
		if err != nil {
			slog.Warn("Error closing connection", slog.String("err", err.Error()), slog.String("addr", conn.RemoteAddr().String()))
		} else {
			slog.Info("Connection closed", slog.String("addr", conn.RemoteAddr().String()))
		}
	}()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		msg := scanner.Text()

		slog.Info("Received message", slog.String("msg", msg), slog.String("addr", conn.RemoteAddr().String()))

		reversedMsg := invertString(msg)
		sentMsg := reversedMsg + "\tСервер написан Осипяном Г. В. М3О-425Бк-21\n"

		time.Sleep(cfg.ImitationTimeDuration)
		slog.Info("Imitated work", slog.String("addr", conn.RemoteAddr().String()))

		_, err := conn.Write([]byte(sentMsg))
		if err != nil && errors.Is(err, os.ErrDeadlineExceeded) {
			slog.Info("Connection timeout", slog.String("addr", conn.RemoteAddr().String()))

			return
		} else if err != nil {
			slog.Warn("Error sending message", slog.String("err", err.Error()), slog.String("addr", conn.RemoteAddr().String()))

			return
		}

		slog.Info("Sent message", slog.String("msg", sentMsg), slog.String("addr", conn.RemoteAddr().String()))
	}

	if err := scanner.Err(); err != nil {
		slog.Warn("Error reading message", slog.String("err", err.Error()), slog.String("addr", conn.RemoteAddr().String()))
	}
}

func invertString(str string) string {
	runes := []rune(str)
	builder := strings.Builder{}
	builder.Grow(len(runes))

	for i := len(runes) - 1; i >= 0; i-- {
		builder.WriteRune(runes[i])
	}

	return builder.String()
}

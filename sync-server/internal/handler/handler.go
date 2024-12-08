package handler

import (
	"errors"
	"log/slog"
	"net"
	"os"
	"strings"
	"time"

	"sync-server/internal/config"
)

const bufferSize = 4096

func HandleConnection(conn net.Conn, cfg config.Config) {
	slog.Info("Handle connection", slog.String("addr", conn.RemoteAddr().String()))

	defer func() {
		err := conn.Close()
		if err != nil {
			slog.Warn("Error closing connection", slog.String("err", err.Error()), slog.String("addr", conn.RemoteAddr().String()))
		} else {
			slog.Info("Connection closed", slog.String("addr", conn.RemoteAddr().String()))
		}
	}()

	receiveMsgBuilder := strings.Builder{}
	buffer := make([]byte, bufferSize)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if errors.Is(err, os.ErrDeadlineExceeded) {
				slog.Info("Connection timeout", slog.String("addr", conn.RemoteAddr().String()))

				return
			}
		}

		receiveMsgBuilder.WriteString(string(buffer))
		if n < bufferSize {
			break
		}
	}

	slog.Info("Received message", slog.String("msg", receiveMsgBuilder.String()), slog.String("addr", conn.RemoteAddr().String()))

	reversedMsg := invertString(receiveMsgBuilder.String())
	sentMsg := reversedMsg + "\nСервер написан Осипяном Г. В. М3О-425Бк-21\n"

	time.Sleep(cfg.ImitationTimeDuration)
	slog.Info("Imitated work", slog.String("addr", conn.RemoteAddr().String()))

	_, err = conn.Write([]byte(sentMsg))
	if err != nil && errors.Is(err, os.ErrDeadlineExceeded) {
		slog.Info("Connection timeout", slog.String("addr", conn.RemoteAddr().String()))

		return
	} else if err != nil {
		slog.Warn("Error sending message", slog.String("err", err.Error()), slog.String("addr", conn.RemoteAddr().String()))

		return
	}

	slog.Info("Sent message", slog.String("msg", sentMsg), slog.String("addr", conn.RemoteAddr().String()))
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

package handler

import (
	"async-server/internal/config"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"strings"
	"time"
)

const bufferSize = 4096

var ErrConnectionClosedByServer = errors.New("connection closed by server")

func HandleConnection(conn net.Conn, cfg config.Config) {
	closeFlag := false
	defer func() {
		if closeFlag {
			if err := conn.Close(); err != nil {
				slog.Warn("Error closing connection", slog.String("err", err.Error()), slog.String("addr", conn.RemoteAddr().String()))
			} else {
				slog.Info("Connection closed", slog.String("addr", conn.RemoteAddr().String()))
			}
		}
	}()

	for {
		err := processConnection(conn, cfg)
		if errors.Is(err, ErrConnectionClosedByServer) {
			slog.Info("Server closed connection on request", slog.String("addr", conn.RemoteAddr().String()))

			break
		} else if errors.Is(err, net.ErrClosed) {
			slog.Info("Client closed connection", slog.String("addr", conn.RemoteAddr().String()))

			break
		} else if errors.Is(err, io.EOF) {
			slog.Info("Client finished sending data", slog.String("addr", conn.RemoteAddr().String()))
			closeFlag = true

			break
		} else if err != nil {
			slog.Warn("Error processing message", slog.String("err", err.Error()), slog.String("addr", conn.RemoteAddr().String()))
			closeFlag = true

			break
		}
	}
}

func processConnection(conn net.Conn, cfg config.Config) error {
	receiveMsgBuilder := strings.Builder{}
	buffer := make([]byte, bufferSize)

	for {
		conn.SetReadDeadline(time.Now().Add(cfg.ReadTimeoutDuration))

		n, err := conn.Read(buffer)
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return net.ErrClosed
			}

			return fmt.Errorf("error reading from connection: %w", err)
		}

		receiveMsgBuilder.WriteString(string(buffer[:n]))

		if n < bufferSize {
			break
		}
	}

	message := strings.TrimSpace(receiveMsgBuilder.String())

	slog.Info("Received message", slog.String("msg", message), slog.String("addr", conn.RemoteAddr().String()))

	if message == "exit" {
		slog.Info("Client requested to close connection", slog.String("addr", conn.RemoteAddr().String()))

		conn.Close()

		return ErrConnectionClosedByServer
	}

	response := invertString(message) + "\nСервер написан Осипяном Г. В. М3О-425Бк-21\n"

	slog.Info("Imitating work", slog.String("addr", conn.RemoteAddr().String()))
	time.Sleep(cfg.ImitationTimeDuration)

	conn.SetWriteDeadline(time.Now().Add(cfg.WriteTimeoutDuration))

	_, err := conn.Write([]byte(response))
	if err != nil {
		return fmt.Errorf("error writing to connection: %w", err)
	}

	slog.Info("Sent message", slog.String("msg", response), slog.String("addr", conn.RemoteAddr().String()))

	return nil
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

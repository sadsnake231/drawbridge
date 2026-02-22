package logging

import (
	"log/slog"
	"os"
)

func InitLogger(logPath string) (*os.File, error) {
	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return nil, err
	}

	handler := slog.NewTextHandler(file, nil)
	logger := slog.New(handler)

	slog.SetDefault(logger)

	return file, nil
}

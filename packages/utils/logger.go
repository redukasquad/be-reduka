package utils

import (
	"io"
	"log/slog"
	"os"
)

var (
	Log       *slog.Logger
	ErrorLog  *slog.Logger
	Important *slog.Logger
)

func openLogFile(path string) *os.File {
	file, err := os.OpenFile(
		path,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		panic("failed to open log file: " + err.Error())
	}
	return file
}

func InitLogger() {
	_ = os.MkdirAll("logs", 0755)

	// files
	appFile := openLogFile("logs/app.log")
	errorFile := openLogFile("logs/error.log")
	importantFile := openLogFile("logs/important.log")

	// stdout + file
	appWriter := io.MultiWriter(os.Stdout, appFile)

	Log = slog.New(
		slog.NewJSONHandler(appWriter, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}),
	)

	ErrorLog = slog.New(
		slog.NewJSONHandler(errorFile, &slog.HandlerOptions{
			Level: slog.LevelError,
		}),
	)

	Important = slog.New(
		slog.NewJSONHandler(importantFile, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}),
	)

	Log.Info("Logger initialized")
}

package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

func JsonLog() {
	logger := slog.New( //	Tạo một logger
		slog.NewJSONHandler(os.Stdout, nil), //	Tạo handler định dạng log thành JSON, 	Ghi log ra standard output
	)

	logger.Info(
		"User logged in",
		"user_id", 123,
		"email", "user@example.com",
	)

	logger.Error(
		"Database connection failed",
		"database", "postgres",
		"error", "connection refused",
	)

	//Tái sử dụng thuộc tính với With()
	userLogger := logger.With(
		"user_id", 123,
		"request_id", "req-abc-001",
	)
	userLogger.Info("User started checkout")
	userLogger.Info("Payment processing started")
}
func LogBasic() {
	userID := 123
	action := "login"

	log.Printf("User %d performed action %s", userID, action)

	slog.Info("User performed an action", "user_id", userID, "action", action)

	slog.Debug("Debugging application")
	slog.Info("Server started", "port", 8080)
	slog.Warn("Request is taking too long", "duration_ms", 2500)
	slog.Error("Failed to connect to database", "error", "connection refused")

	slog.Info(
		"User logged in",
		"user_id", 123,
		"email", "user@example.com",
		"ip", "192.168.1.10",
	)
}

func LogLevel() {
	var level slog.LevelVar

	level.Set(slog.LevelInfo)
	logger := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: &level,
		}),
	)
	logger.Debug("Debug message")
	logger.Info("Info message")
	logger.Warn("Warning message")
	logger.Error("Error message")

	level.Set(slog.LevelDebug)

	logger.Debug("Debug message after changing level")
}
func LogProduction() {
	ctx := context.Background()
	logger := slog.New(
		slog.NewJSONHandler(os.Stdout, nil),
	)
	logger.LogAttrs(
		ctx,
		slog.LevelInfo,
		"User logged in",
		slog.Int("user_id", 123),
		slog.String("action", "login"),
	)
}
func main() {
	LogProduction()

	logger := slog.New(
		tint.NewTextHandler(os.Stdout, &tint.Options{
			Level: slog.LevelDebug,
		}),
	)

	logger.Debug("Debugging application", "step", 1)
	logger.Info("Server started", "port", 8080)
	logger.Warn("Slow request", "duration_ms", 2500)
	logger.Error("Database connection failed", "database", "postgres")
}

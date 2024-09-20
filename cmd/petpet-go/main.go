package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/wavy-cat/petpet-go/http/handler"
	"github.com/wavy-cat/petpet-go/http/middleware"
	"github.com/wavy-cat/petpet-go/internal/config"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Настраиваем логгер
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			fmt.Println("Error syncing logger:", err)
		}
	}(logger)

	// Настраиваем роутер
	router := mux.NewRouter()

	handle := middleware.Logging{
		Logger: logger,
		Next:   handler.Handler{},
	}
	router.Handle("/ds/{user_id}.gif", &handle).Methods(http.MethodGet)
	router.Handle("/ds/{user_id}", &handle).Methods(http.MethodGet)

	router.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		_, err := w.Write([]byte("Ты думал тут что-то будет?"))
		if err != nil {
			logger.Error("Error sending response", zap.Error(err))
		}
	})

	// Настраиваем сервер
	srv := &http.Server{
		Addr:    config.HTTPAddress,
		Handler: router,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Запускаем сервер
	go func() {
		logger.Info("Starting the HTTP server...")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("Server failed:", zap.Error(err))
		}
	}()

	// Ожидаем сигнал завершения
	<-stop

	// Создаём контекст с таймаутом для корректного завершения сервера
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.ShutdownTimeout)*time.Second)
	defer cancel()

	logger.Info("Shutting down the server...")
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown:", zap.Error(err))
	}

	logger.Info("Server exited properly")
}

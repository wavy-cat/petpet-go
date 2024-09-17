package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/wavy-cat/petpet-go/http/handler"
	"github.com/wavy-cat/petpet-go/http/middleware"
	"github.com/wavy-cat/petpet-go/internal/config"
	"go.uber.org/zap"
	"net/http"
)

func main() {
	// Настраиваем логгер
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Println("Error initializing logger:", err)
		return
	}
	defer logger.Sync()

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

	// Запускаем сервер
	logger.Info("Starting the HTTP server...")
	if err := http.ListenAndServe(config.HTTPAddress, router); err != nil {
		logger.Fatal(err.Error())
	}
}

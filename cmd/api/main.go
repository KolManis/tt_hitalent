package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/KolManis/tt_hitalent/internal/config"
	"github.com/KolManis/tt_hitalent/internal/database"
	"github.com/KolManis/tt_hitalent/internal/handlers"
	"github.com/KolManis/tt_hitalent/pkg/middleware"
)

const (
	httpPort          = "8080"
	readheaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

func main() {

	cfg := config.Load()

	// Подключение к БД
	db, err := database.Connect(cfg)

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := database.RunMigrations(db); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	chatHandler := handlers.NewChatHandler(db)
	messageHandler := handlers.NewMessageHandler(db)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /chats", chatHandler.CreateChat)
	mux.HandleFunc("GET /chats/{id}", chatHandler.GetChat)
	mux.HandleFunc("DELETE /chats/{id}", chatHandler.DeleteChat)

	mux.HandleFunc("POST /chats/{id}/messages", messageHandler.CreateMessage)

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	handler := middleware.Logger(mux)

	server := &http.Server{
		Addr:              ":" + httpPort,
		Handler:           handler,
		ReadHeaderTimeout: readheaderTimeout, // Защита от Slowloris
		WriteTimeout:      10 * time.Second,  // Защита от медленной отправки
	}

	go func() {
		log.Printf("HTTP-сервер запущен на порту %s\n", httpPort)
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("Ошибка запуска сервера: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Printf("Завершение работы сервера")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	// новые соединения не будут создаваться, а старые мы завершаем в пределат timeout
	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("Ошибка при остановке сервера %v\n", err)
	}

	log.Println("Сервер остановлен")
}

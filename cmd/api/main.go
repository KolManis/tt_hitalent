package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	httpPort = "8080"
	//Таймауты для HTTP сервера
	readheaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

func foo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Привет! Это тестовый сервер.\n")
	fmt.Fprintf(w, "Параметры запроса: %s\n", r.URL.Query())
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", foo)

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	server := &http.Server{
		Addr:              net.JoinHostPort("localhost", httpPort),
		Handler:           mux,
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
	err := server.Shutdown(ctx)
	if err != nil {
		log.Printf("Ошибка при остановке сервера %v\n", err)
	}

	log.Println("Сервер остановлен")
}

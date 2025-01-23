package run

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/jwtauth"
	"go.uber.org/zap"
	"studentgit.kata.academy/ponomarenko.100299/go-petstore/internal/db"
	"studentgit.kata.academy/ponomarenko.100299/go-petstore/internal/modules"
	"studentgit.kata.academy/ponomarenko.100299/go-petstore/internal/responder"
	"studentgit.kata.academy/ponomarenko.100299/go-petstore/internal/router"
)

type Server interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}

type Runner interface {
	Run()
}

// Bootstraper - интерфейс инициализации приложения
type Bootstraper interface {
	Bootstrap() Runner
}

type App struct {
	srv Server
}

func NewApp() *App {
	return &App{}
}

func (a *App) Run() {
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Println("Shutdown signal received...")

		// Таймаут для завершения сервера
		shutdownCtx, cancel := context.WithTimeout(serverCtx, 5*time.Second)
		defer cancel()

		// Завершение работы сервера
		if err := a.srv.Shutdown(shutdownCtx); err != nil {
			log.Printf("Graceful shutdown failed: %v", err)
		} else {
			log.Println("Server stopped gracefully")
		}

		// Завершаем основной контекст сервера
		serverStopCtx()
	}()

	// Запуск сервера в отдельной горутине
	go func() {
		log.Println("Starting server...")
		err := a.srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	<-serverCtx.Done()
}

func (a *App) Bootstrap() Runner {
	var tokenAuth *jwtauth.JWTAuth
	const Secret = "rixha9-kabtuJ-xizpej"
	tokenAuth = jwtauth.New("HS256", []byte(Secret), nil)

	logger, _ := zap.NewProduction()

	responder := responder.NewResponder(logger)

	dbConf := db.NewDbConf()

	dbRaw, err := db.NewDB(*dbConf)
	if err != nil {
		log.Fatal(err)
	}

	err = db.MigrateDB(dbRaw)
	if err != nil {
		log.Fatal(err)
	}

	storages := modules.NewStorages(dbRaw)

	services := modules.NewServices(*storages, tokenAuth)

	controllers := modules.NewControllers(services, responder)

	r := router.NewRouter(controllers, tokenAuth)

	a.srv = &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return a
}

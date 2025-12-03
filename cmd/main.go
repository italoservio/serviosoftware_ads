package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/italoservio/serviosoftware_ads/internal/api"
	"github.com/italoservio/serviosoftware_ads/internal/deps"
	"github.com/italoservio/serviosoftware_ads/pkg/env"
)

func main() {
	envVars := env.Load()
	log.SetOutput(os.Stdout)

	container := deps.NewContainer(envVars)

	router := mux.NewRouter()
	router.MethodNotAllowedHandler = http.HandlerFunc(api.MethodNotAllowed)

	api.RegisterInfraRoutes(router)
	api.RegisterCloakersRoutes(router, container)

	corsOrigins := handlers.AllowedOrigins([]string{
		"http://localhost:5173",
		"http://localhost:5174",
		"https://users.serviosoftware.com",
		"https://ads.serviosoftware.com",
	})
	corsHeaders := handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})
	corsMethods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"})
	corsCredentials := handlers.AllowCredentials()
	withCors := handlers.CORS(corsOrigins, corsHeaders, corsMethods, corsCredentials)

	svr := &http.Server{
		Addr:    ":8080",
		Handler: withCors(router),
	}

	wg := sync.WaitGroup{}

	wg.Go(func() {
		log.Println("servidor escutando na porta :8080")

		err := svr.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}

		log.Println("servidor encerrado")
	})

	exitCode := 0
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT)
	signal.Notify(sigCh, syscall.SIGTERM)

	wg.Go(func() {
		sig := <-sigCh
		log.Println("sinal recebido:", strings.ToUpper(sig.String()))

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := svr.Shutdown(ctx); err != nil {
			println("falha ao encerrar o servidor:", err.Error())
			exitCode = 1
		}

		if err := container.DB.Disconnect(); err != nil {
			println("falha ao desconectar do banco de dados:", err.Error())
			exitCode = 1
		}
	})

	wg.Wait()
	os.Exit(exitCode)
}

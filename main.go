package main

import (
	"context"
	_ "embed"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/flowchartsman/swaggerui"

	"github.com/ziemowit141/payment_api/database"
	"github.com/ziemowit141/payment_api/database/seed"
	"github.com/ziemowit141/payment_api/handlers"
)

//go:embed swagger.yaml
var spec []byte

func main() {
	log.Println("Starting server")
	db := database.SetupDatabase("myapp")
	seed.LoadTestCreditCards(db)

	sm := http.NewServeMux()
	sm.Handle("/authorize", handlers.NewAuthorizeHandler(db))
	sm.Handle("/void", handlers.NewVoidHandler(db))
	sm.Handle("/capture", handlers.NewCaptureHandler(db))
	sm.Handle("/refund", handlers.NewRefundHandler(db))

	sm.Handle("/swagger/", http.StripPrefix("/swagger", swaggerui.Handler(spec)))

	server := http.Server{
		Addr:         ":3000",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, syscall.SIGTERM)

	sig := <-sigChan
	log.Println("Graceful shutdown", sig)

	tc, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	err := server.Shutdown(tc)
	if err != nil {
		log.Printf("Error occured during server shuttdown: %s", err)
		log.Println("Canceling context manually")
		cancel()
	}
}

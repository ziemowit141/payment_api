package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ziemowit141/payment_api/database"
	"github.com/ziemowit141/payment_api/database/seed"
	"github.com/ziemowit141/payment_api/handlers"
)

func main() {
	log.Println("Starting server")
	db := database.SetupDatabase("myapp")
	seed.LoadTestCreditCards(db)

	sm := http.NewServeMux()
	sm.Handle("/authorize", handlers.NewAuthorizeHandler(db))
	sm.Handle("/void", handlers.NewVoidHandler(db))
	sm.Handle("/capture", handlers.NewCaptureHandler(db))
	sm.Handle("/refund", handlers.NewRefundHandler(db))

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

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	server.Shutdown(tc)
}

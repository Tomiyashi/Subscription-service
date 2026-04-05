package main

import (
	"log"
	"net/http"
	"os"
	"subscription-service/internal/handlers"
	"subscription-service/internal/repository"
	"subscription-service/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	pool, err := repository.InitDB()
	if err != nil {
		log.Fatalf("❌ DB error: %v", err)
	}
	defer pool.Close()

	repo := repository.NewSubscriptionRepository(pool)
	svc := service.NewSubscriptionService(repo)
	h := handlers.NewHandler(svc)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", h.HealthCheck)
	r.Post("/subscriptions", h.CreateSubscription)
	r.Get("/subscriptions/{id}", h.GetSubscription)
	r.Get("/subscriptions", h.ListSubscriptions)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server starting on :%s", port)

	log.Fatal(http.ListenAndServe(":"+port, r))

}

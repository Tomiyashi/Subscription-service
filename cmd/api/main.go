package main

import (
	"log"
	"net/http"
	"os"
	"subscription-service/internal/handlers"
	"subscription-service/internal/repository"
	"subscription-service/internal/service"

	_ "subscription-service/docs"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
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

	r.Get("/subscriptions/total-cost", h.GetTotalCost)
	r.Get("/health", h.HealthCheck)
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Get("/subscriptions/{id}", h.GetSubscription)
	r.Put("/subscriptions/{id}", h.UpdateSubscription)
	r.Delete("/subscriptions/{id}", h.DeleteSubscription)

	r.Get("/subscriptions", h.ListSubscriptions)
	r.Post("/subscriptions", h.CreateSubscription)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server starting on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

package service

import (
	"context"
	"fmt"
	"subscription-service/internal/models"
	"subscription-service/internal/repository"
	"time"

	"github.com/google/uuid"
)

type SubscriptionService struct {
	repo repository.SubscriptionRepository
}

func NewSubscriptionService(repo repository.SubscriptionRepository) *SubscriptionService {
	return &SubscriptionService{repo: repo}
}

func (s *SubscriptionService) CreateSubscription(ctx context.Context, req *models.CreateSubscriptionRequest) (*models.Subscription, error) {
	if req.Price <= 0 {
		return nil, fmt.Errorf("validation: price must be positive, got %d", req.Price)
	}
	if req.ServiceName == "" {
		return nil, fmt.Errorf("validation: service_name cannot be empty")
	}
	today := time.Now().UTC().Truncate(24 * time.Hour)
	if req.StartDate.Before(today) {
		return nil, fmt.Errorf("validation: start_date cannot be in the past")
	}
	if req.EndDate != nil && req.EndDate.Before(req.StartDate) {
		return nil, fmt.Errorf("validation: end_date must be after start_date")
	}
	sub := &models.Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      req.UserID,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
	}

	if err := s.repo.Create(ctx, sub); err != nil {
		return nil, fmt.Errorf("service: create subscription failed: %w", err)
	}

	return sub, nil
}

func (s *SubscriptionService) GetSubscription(ctx context.Context, id uuid.UUID) (*models.Subscription, error) {
	sub, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("service: get subscription failed: %w", err)
	}
	return sub, nil
}

func (s *SubscriptionService) ListSubscriptions(ctx context.Context, userID uuid.UUID) ([]*models.Subscription, error) {
	subs, err := s.repo.List(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("service: list subscriptions failed: %w", err)
	}
	return subs, nil
}

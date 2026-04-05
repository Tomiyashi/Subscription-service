package repository

import (
	"context"
	"fmt"
	"subscription-service/internal/models"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type subscriptionRepo struct {
	pool *pgxpool.Pool
}

type SubscriptionRepository interface {
	Create(ctx context.Context, sub *models.Subscription) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Subscription, error)
	List(ctx context.Context, userID uuid.UUID) ([]*models.Subscription, error)
	Update(ctx context.Context, sub *models.Subscription) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetTotalCost(ctx context.Context, userID uuid.UUID, serviceName *string, from, to time.Time) (int, error)
}

func NewSubscriptionRepository(pool *pgxpool.Pool) SubscriptionRepository {
	return &subscriptionRepo{pool: pool}
}

func (r *subscriptionRepo) Create(ctx context.Context, sub *models.Subscription) error {
	sql := `
	INSERT INTO subscriptions(service_name, price, user_id, start_date, end_date)
	VALUES($1, $2, $3, $4, $5) RETURNING id
	`

	err := r.pool.QueryRow(ctx, sql,
		sub.ServiceName,
		sub.Price,
		sub.UserID,
		sub.StartDate,
		sub.EndDate).Scan(&sub.ID)

	if err != nil {
		return fmt.Errorf("repository: create subscription failed: %w", err)
	}
	return nil
}

func (r *subscriptionRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Subscription, error) {
	sub := &models.Subscription{}
	sql := `
	SELECT id, service_name, price, user_id, start_date, end_date
	FROM subscriptions
	WHERE id = $1
	`
	err := r.pool.QueryRow(ctx, sql, id).Scan(
		&sub.ID,
		&sub.ServiceName,
		&sub.Price,
		&sub.UserID,
		&sub.StartDate,
		&sub.EndDate)
	if err != nil {
		return nil, fmt.Errorf("repository: get by id failed: %w", err)
	}
	return sub, nil
}

func (r *subscriptionRepo) List(ctx context.Context, userID uuid.UUID) ([]*models.Subscription, error) {
	var subscriptions []*models.Subscription

	sql := `
	SELECT id, service_name, price, user_id, start_date, end_date 
	FROM subscriptions 
	WHERE user_id = $1
	ORDER BY start_date DESC
	`
	rows, err := r.pool.Query(ctx, sql, userID)
	if err != nil {
		return nil, fmt.Errorf("repository: list failed: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		sub := &models.Subscription{}
		err := rows.Scan(
			&sub.ID,
			&sub.ServiceName,
			&sub.Price,
			&sub.UserID,
			&sub.StartDate,
			&sub.EndDate,
		)
		if err != nil {
			return nil, fmt.Errorf("repository: list scan failed: %w", err)
		}
		subscriptions = append(subscriptions, sub)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repository: list rows error: %w", err)
	}

	return subscriptions, nil
}

func (r *subscriptionRepo) Update(ctx context.Context, sub *models.Subscription) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE subscriptions 
		SET service_name=$1, price=$2, start_date=$3, end_date=$4 
		WHERE id=$5`,
		sub.ServiceName, sub.Price, sub.StartDate, sub.EndDate, sub.ID)
	if err != nil {
		return fmt.Errorf("repository: update failed: %w", err)
	}
	return nil
}

func (r *subscriptionRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM subscriptions WHERE id=$1`, id)
	if err != nil {
		return fmt.Errorf("repository: delete failed: %w", err)
	}
	return nil
}

func (r *subscriptionRepo) GetTotalCost(ctx context.Context, userID uuid.UUID, serviceName *string, from, to time.Time) (int, error) {
	var total int
	query := `SELECT COALESCE(SUM(price), 0) FROM subscriptions WHERE user_id=$1 AND start_date >= $2 AND start_date <= $3`
	args := []interface{}{userID, from, to}
	if serviceName != nil && *serviceName != "" {
		query += ` AND service_name ILIKE $4`
		args = append(args, "%"+*serviceName+"%")
	}
	err := r.pool.QueryRow(ctx, query, args...).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("repository: get total cost failed: %w", err)
	}
	return total, nil
}

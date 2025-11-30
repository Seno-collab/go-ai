package restaurantrepo

import (
	"context"
	"errors"
	"go-ai/internal/domain/restaurant"
	sqlc "go-ai/internal/infra/sqlc/restaurant"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RestaurantRepo struct {
	q *sqlc.Queries
}

func NewRestaurantRepo(pool *pgxpool.Pool) *RestaurantRepo {
	return &RestaurantRepo{
		q: sqlc.New(pool),
	}
}

func (rr *RestaurantRepo) CreateRestaurant(r *restaurant.Entity) (int32, error) {
	_, err := rr.q.GetByName(context.Background(), r.Name)
	if err != nil || !errors.Is(err, pgx.ErrNoRows) {
		return 0, err
	}
	id, err := rr.q.CreateRestaurant(context.Background(), sqlc.CreateRestaurantParams{
		Name:        r.Name,
		Description: r.Description,
		Address:     r.Address,
		Category:    r.Category,
		City:        r.City,
		District:    r.District,
		LogoUrl:     r.LogoUrl,
		BannerUrl:   r.BannerUrl,
		PhoneNumber: r.PhoneNumber,
		WebsiteUrl:  r.WebsiteUrl,
		Email:       r.Email,
		UserID:      r.UserID,
	})
	if err != nil {
		return 0, err
	}
	return id, nil
}

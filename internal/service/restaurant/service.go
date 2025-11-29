package restaurant_service

import "go-ai/internal/domain/restaurant"

type Service struct {
	repo restaurant.Repository
}

func NewService(r restaurant.Repository) *Service {
	return &Service{
		repo: r,
	}
}

// func (s *Service) CreateRestaurant(r *restaurant.Entity) (int32, error) {

// }

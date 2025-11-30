package restaurant

type Repository interface {
	CreateRestaurant(r *Entity) (int32, error)
}

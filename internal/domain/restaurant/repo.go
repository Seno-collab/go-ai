package restaurant

type Repository interface {
	CreateRestaurant(*Entity) (int, error)
}

package pkg

type Repository[T any] interface {
	Create(t *T) (*T, error)
	DeleteById(id int) error
	Exists(t *T) (bool, error)
	GetAll() (*[]T, error)
	GetBySlug(slug string) (*T, error)
	Update(t *T) (*T, error)
}

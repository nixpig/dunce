package pkg

type Repository[T any, D any] interface {
	Create(t *D) (*T, error)
	DeleteById(id int) error
	Exists(t *D) (bool, error)
	GetAll() (*[]T, error)
	GetBySlug(slug string) (*T, error)
	Update(t *T) (*T, error)
}

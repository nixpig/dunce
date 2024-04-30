package pkg

type Service[T any] interface {
	Create(t *T) (*T, error)
	DeleteById(id int) error
	GetAll() (*[]T, error)
	GetBySlug(slug string) (*T, error)
	Update(t *T) (*T, error)
}

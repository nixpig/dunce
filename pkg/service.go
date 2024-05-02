package pkg

type Service[T any, N any] interface {
	Create(n *N) (*T, error)
	DeleteById(id int) error
	GetAll() (*[]T, error)
	GetBySlug(slug string) (*T, error)
	Update(t *T) (*T, error)
}

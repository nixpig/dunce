package pkg

type Service[T any, N any] interface {
	Create(n *N) (*T, error)
	GetAll() (*[]T, error)
	GetByAttribute(attr, value string) (*T, error)
	Update(t *T) (*T, error)
	DeleteById(id int) error
}

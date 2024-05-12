package tag

type Tag struct {
	Id   int    `validate:"omitempty"`
	Name string `validate:"required,min=2,max=30"`
	Slug string `validate:"required,slug,min=2,max=50,lowercase"`
}

type CreateTagRequestDto struct {
	Name string `validate:"required,min=2,max=30"`
	Slug string `validate:"required,slug,min=2,max=50"`
}

type UpdateTagRequestDto struct {
	Id   int    `validate:"required"`
	Name string `validate:"required,min=2,max=30"`
	Slug string `validate:"required,slug,min=2,max=50"`
}

type TagResponseDto struct {
	Id   int
	Name string
	Slug string
}

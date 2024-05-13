package tag

type Tag struct {
	Id   int    `validate:"omitempty"`
	Name string `validate:"required,min=2,max=30"`
	Slug string `validate:"required,slug,min=2,max=50,lowercase"`
}

type TagNewRequestDto struct {
	Name string `validate:"required,min=2,max=30"`
	Slug string `validate:"required,slug,min=2,max=50"`
}

type TagUpdateRequestDto struct {
	Id   int    `validate:"required"`
	Name string `validate:"required,min=2,max=30"`
	Slug string `validate:"required,slug,min=2,max=50"`
}

type TagResponseDto struct {
	Id   int    `validate:"required"`
	Name string `validate:"required,min=2,max=30"`
	Slug string `validate:"required,slug,min=2,max=50"`
}

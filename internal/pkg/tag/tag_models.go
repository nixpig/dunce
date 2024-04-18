package tag

type Tag struct {
	Id   int
	Name string `validate:"required,min=5,max=30"`
	Slug string `validate:"required,slug,min=5,max=50"`
}

type TagNew struct {
	Name string `validate:"required,min=5,max=30"`
	Slug string `validate:"required,slug,min=5,max=50"`
}

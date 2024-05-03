package tag

type Tag struct {
	Id int
	TagData
}

type TagData struct {
	Name string `validate:"required,min=2,max=30"`
	Slug string `validate:"required,slug,min=2,max=50,lowercase"`
}

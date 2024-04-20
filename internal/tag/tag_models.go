package tag

type Tag struct {
	Id   int
	Name string `validate:"required,min=5,max=30"`
	Slug string `validate:"required,slug,min=5,max=50"`
}

func NewTag(name, slug string) Tag {
	return Tag{Name: name, Slug: slug}
}

func NewTagWithId(id int, name, slug string) Tag {
	return Tag{Id: id, Name: name, Slug: slug}
}

package tag

type Tag struct {
	Id   int
	Name string `validate:"required,tagname,min=2,max=30"`
	Slug string `validate:"required,slug,min=2,max=50,lowercase"`
}

func NewTag(name, slug string) Tag {
	return Tag{Name: name, Slug: slug}
}

func NewTagWithId(id int, name, slug string) Tag {
	return Tag{Id: id, Name: name, Slug: slug}
}

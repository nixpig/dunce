package tag

type TagService struct {
	data TagDataInterface
}

func NewTagService(data TagDataInterface) TagService {
	return TagService{data}
}

func (ts *TagService) Create(tag *Tag) (*Tag, error) {
	// TODO: data validation and such

	createdTag, err := ts.data.create(tag)
	if err != nil {
		return nil, err
	}

	return createdTag, nil
}

func (ts *TagService) DeleteById(id int) error {
	if err := ts.data.deleteById(id); err != nil {
		return err
	}

	return nil
}

package tag

import (
	"github.com/go-playground/validator/v10"
)

type ITagService interface {
	Create(tag *TagData) (*Tag, error)
	DeleteById(id int) error
	GetAll() (*[]Tag, error)
	GetByAttribute(attr, value string) (*Tag, error)
	Update(tag *Tag) (*Tag, error)
}

type TagService struct {
	repo     ITagRepository
	validate *validator.Validate
}

func NewTagService(
	repo ITagRepository,
	validate *validator.Validate,
) TagService {
	return TagService{
		repo:     repo,
		validate: validate,
	}
}

func (t TagService) Create(tag *TagData) (*Tag, error) {
	// TODO: make slug lowercase
	// TODO: custom validator for tag name

	if err := t.validate.Struct(tag); err != nil {
		return nil, err
	}

	createdTag, err := t.repo.Create(tag)
	if err != nil {
		return nil, err
	}

	return createdTag, nil
}

func (t TagService) DeleteById(id int) error {
	return t.repo.DeleteById(id)
}

func (t TagService) GetAll() (*[]Tag, error) {
	tags, err := t.repo.GetAll()
	if err != nil {
		return nil, err
	}

	return tags, nil
}

func (t TagService) GetManyByAttribute(attr, value string) (*[]Tag, error) {
	return nil, nil
}

func (t TagService) GetByAttribute(attr, value string) (*Tag, error) {
	tag, err := t.repo.GetByAttribute(attr, value)
	if err != nil {
		return nil, err
	}

	return tag, nil
}

func (t TagService) Update(tag *Tag) (*Tag, error) {
	if err := t.validate.Struct(tag); err != nil {
		return nil, err
	}

	updatedTag, err := t.repo.Update(tag)
	if err != nil {
		return nil, err
	}

	return updatedTag, nil
}

package tags

import (
	"github.com/go-playground/validator/v10"
	"github.com/nixpig/dunce/pkg/logging"
)

type TagService struct {
	data     TagDataInterface
	validate *validator.Validate
	log      logging.Logger
}

type TagServiceInterface interface {
	Create(tag *Tag) (*Tag, error)
	DeleteById(id int) error
	GetAll() (*[]Tag, error)
	GetBySlug(slug string) (*Tag, error)
	Update(tag *Tag) (*Tag, error)
}

func NewTagService(
	data TagDataInterface,
	validate *validator.Validate,
	log logging.Logger,
) TagService {
	return TagService{
		data:     data,
		validate: validate,
		log:      log,
	}
}

func (ts TagService) Create(tag *Tag) (*Tag, error) {
	// TODO: make slug lowercase
	// TODO: custom validator for tag name

	if err := ts.validate.Struct(tag); err != nil {
		ts.log.Error(err.Error())
		return nil, err
	}

	createdTag, err := ts.data.Create(tag)
	if err != nil {
		ts.log.Error(err.Error())
		return nil, err
	}

	return createdTag, nil
}

func (ts TagService) DeleteById(id int) error {
	if err := ts.data.DeleteById(id); err != nil {
		ts.log.Error(err.Error())
		return err
	}

	return nil
}

func (ts TagService) GetAll() (*[]Tag, error) {
	tags, err := ts.data.GetAll()
	if err != nil {
		ts.log.Error(err.Error())
		return nil, err
	}

	return tags, nil
}

func (ts TagService) GetBySlug(slug string) (*Tag, error) {
	tag, err := ts.data.GetBySlug(slug)
	if err != nil {
		ts.log.Error(err.Error())
		return nil, err
	}

	return tag, nil
}

func (ts TagService) Update(tag *Tag) (*Tag, error) {
	if err := ts.validate.Struct(tag); err != nil {
		ts.log.Error(err.Error())
		return nil, err
	}

	updatedTag, err := ts.data.Update(tag)
	if err != nil {
		ts.log.Error(err.Error())
		return nil, err
	}

	return updatedTag, nil
}

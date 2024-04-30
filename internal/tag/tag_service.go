package tag

import (
	"github.com/go-playground/validator/v10"
	"github.com/nixpig/dunce/pkg"
	"github.com/nixpig/dunce/pkg/logging"
)

type TagService struct {
	repo     pkg.Repository[Tag]
	validate *validator.Validate
	log      logging.Logger
}

func NewTagService(
	repo pkg.Repository[Tag],
	validate *validator.Validate,
	log logging.Logger,
) TagService {
	return TagService{
		repo:     repo,
		validate: validate,
		log:      log,
	}
}

func (t TagService) Create(tag *Tag) (*Tag, error) {
	// TODO: make slug lowercase
	// TODO: custom validator for tag name

	if err := t.validate.Struct(tag); err != nil {
		t.log.Error(err.Error())
		return nil, err
	}

	createdTag, err := t.repo.Create(tag)
	if err != nil {
		t.log.Error(err.Error())
		return nil, err
	}

	return createdTag, nil
}

func (t TagService) DeleteById(id int) error {
	if err := t.repo.DeleteById(id); err != nil {
		t.log.Error(err.Error())
		return err
	}

	return nil
}

func (t TagService) GetAll() (*[]Tag, error) {
	tags, err := t.repo.GetAll()
	if err != nil {
		t.log.Error(err.Error())
		return nil, err
	}

	return tags, nil
}

func (t TagService) GetBySlug(slug string) (*Tag, error) {
	tag, err := t.repo.GetBySlug(slug)
	if err != nil {
		t.log.Error(err.Error())
		return nil, err
	}

	return tag, nil
}

func (t TagService) Update(tag *Tag) (*Tag, error) {
	if err := t.validate.Struct(tag); err != nil {
		t.log.Error(err.Error())
		return nil, err
	}

	updatedTag, err := t.repo.Update(tag)
	if err != nil {
		t.log.Error(err.Error())
		return nil, err
	}

	return updatedTag, nil
}

package tag

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

type TagService interface {
	Create(tag *CreateTagRequestDto) (*TagResponseDto, error)
	DeleteById(id int) error
	GetAll() (*[]TagResponseDto, error)
	GetByAttribute(attr, value string) (*TagResponseDto, error)
	Update(tag *UpdateTagRequestDto) (*TagResponseDto, error)
}

type TagServiceImpl struct {
	repo     TagRepository
	validate *validator.Validate
}

func NewTagService(
	repo TagRepository,
	validate *validator.Validate,
) TagServiceImpl {
	return TagServiceImpl{
		repo:     repo,
		validate: validate,
	}
}

func (t TagServiceImpl) Create(tag *CreateTagRequestDto) (*TagResponseDto, error) {
	tagToCreate := Tag{
		Name: tag.Name,
		Slug: strings.ToLower(tag.Slug),
	}

	if err := t.validate.Struct(tagToCreate); err != nil {
		return nil, err
	}

	createdTag, err := t.repo.Create(&tagToCreate)
	if err != nil {
		return nil, err
	}

	return &TagResponseDto{
		Id:   createdTag.Id,
		Name: createdTag.Name,
		Slug: createdTag.Slug,
	}, nil
}

func (t TagServiceImpl) DeleteById(id int) error {
	return t.repo.DeleteById(id)
}

func (t TagServiceImpl) GetAll() (*[]TagResponseDto, error) {
	tags, err := t.repo.GetAll()
	if err != nil {
		return nil, err
	}

	allTags := make([]TagResponseDto, len(*tags))

	for index, tag := range *tags {
		allTags[index] = TagResponseDto{
			Id:   tag.Id,
			Name: tag.Name,
			Slug: tag.Slug,
		}
	}

	return &allTags, nil
}

func (t TagServiceImpl) GetByAttribute(attr, value string) (*TagResponseDto, error) {
	tag, err := t.repo.GetByAttribute(attr, value)
	if err != nil {
		return nil, err
	}

	return &TagResponseDto{
		Id:   tag.Id,
		Name: tag.Name,
		Slug: tag.Slug,
	}, nil
}

func (t TagServiceImpl) Update(tag *UpdateTagRequestDto) (*TagResponseDto, error) {
	tagToUpdate := Tag{
		Id:   tag.Id,
		Name: tag.Name,
		Slug: strings.ToLower(tag.Slug),
	}

	if err := t.validate.Struct(tagToUpdate); err != nil {
		return nil, err
	}

	updatedTag, err := t.repo.Update(&tagToUpdate)
	if err != nil {
		return nil, err
	}

	return &TagResponseDto{
		Id:   updatedTag.Id,
		Name: updatedTag.Name,
		Slug: updatedTag.Slug,
	}, nil
}

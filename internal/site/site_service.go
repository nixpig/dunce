package site

type SiteService interface {
	Create(key, value string) (*SiteItemResponseDto, error)
}

type SiteServiceImpl struct {
	repo SiteRepository
}

func NewSiteService(repo SiteRepository) SiteServiceImpl {
	return SiteServiceImpl{repo: repo}
}

func (s SiteServiceImpl) Create(key, value string) (*SiteItemResponseDto, error) {
	item, err := s.repo.Create(key, value)
	if err != nil {
		return nil, err
	}

	return &SiteItemResponseDto{
		Id:    item.Id,
		Key:   item.Key,
		Value: item.Value,
	}, nil
}

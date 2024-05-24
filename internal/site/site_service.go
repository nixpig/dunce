package site

type SiteService interface {
	Create(key, value string) (*SiteKv, error)
}

type SiteServiceImpl struct {
	repo SiteRepository
}

func NewSiteService(repo SiteRepository) SiteServiceImpl {
	return SiteServiceImpl{repo: repo}
}

func (s SiteServiceImpl) Create(key, value string) (*SiteKv, error) {
	return s.repo.Create(key, value)
}

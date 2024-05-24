package site

type Site struct {
	Name    string
	Tagline string
	Domain  string
}

type SiteKv struct {
	Id    uint
	Key   string
	Value string
}

type SiteItemResponseDto struct {
	Id    uint
	Key   string
	Value string
}

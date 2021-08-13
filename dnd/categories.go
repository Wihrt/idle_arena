package dnd

type CategoryIndex struct {
	Count   int            `json:"count"`
	Results []CategoryItem `json:"results"`
}

type CategoryList struct {
	Index      string         `json:"index"`
	Name       string         `json:"name"`
	Equipement []CategoryItem `json:"equipment"`
}

type CategoryItem struct {
	Index string `json:"index"`
	Name  string `json:"name"`
	URL   string `json:"url"`
}

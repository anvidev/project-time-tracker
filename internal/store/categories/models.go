package categories

type Category struct {
	Id        int64  `json:"id"`
	Title     string `json:"title"`
	RootTitle string `json:"rootTitle"`
}

type CategoryTree struct {
	Id         int64           `json:"id"`
	ParentId   *int64          `json:"parentId"`
	Title      string          `json:"title"`
	IsRetired  bool            `json:"isRetired"`
	IsFollowed bool            `json:"isFollowed"`
	Children   []*CategoryTree `json:"children"`
}

type CreateCategoryInput struct {
	Title    string `json:"title"`
	ParentId *int64 `json:"parentId"`
}

type UpdateCategoryInput struct {
	Title string `json:"title"`
}

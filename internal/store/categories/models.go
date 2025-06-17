package categories

type Category struct {
	Id        int64  `json:"id"`
	Title     string `json:"title"`
	RootTitle string `json:"rootTitle"`
}

package news

type Post struct {
	ID      int    `json:"ID,omitempty"`
	Title   string `json:"title,omitempty"`
	Content string `json:"content,omitempty"`
	PubTime int64  `json:"pubTime,omitempty"`
	Link    string `json:"link,omitempty"`
}

type Pagination struct {
	NumOfPages int `json:"numOfPages,omitempty"`
	Page       int `json:"page,omitempty"`
	Limit      int `json:"limit,omitempty"`
}

// news contract
type Interface interface {
	Posts(limit, offset int) ([]Post, error)
	AddPost(p Post) error
	PostSearchILIKE(keyWord string, limit, offset int) ([]Post, Pagination, error)
	PostsCreation([]Post) error
	PostDetails(id int) (Post, error)
	CreateGonewsTable() error
	DropGonewsTable() error
}

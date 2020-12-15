package model


// ListSearch 博客搜索
type ListSearch struct {
	SortKey   string   `json:"sortKey"`
	Tags      []string `json:"tags"`
	Page      int      `json:"page"`
	PageSize  int      `json:"pageSize"`
	SearchKey string   `json:"searchKey"`
	CreateAt  struct {
		Begin int64 `json:"begin"`
		End   int64 `json:"end"`
	}
}

// BlogList 博客列表
type BlogList struct {
	Author string `json:"author"`
	Uv int64 `json:"uv"`
	Pv int64 `json:"pv"`
	CreatedAt int64 `json:"createdTime"`
	Title string `json:"title"`
	ID string `json:"id"`
	Tags []string `json:"tags"`
}

// BlogDetail 博客详情
type BlogDetail struct {
	BlogList
	Content string `json:"content"`
}

package model

// StateOverViewResp 博客统计返回
type StateOverViewResp struct {
	ID string `json:"id"`
	Title string `json:"title"`
	Tags []string `json:"tags"`
	Author string `json:"author"`
	Pv int64 `json:"pv"`
	Uv int64 `json:"uv"`
}

// StatDetailResp 博客趋势详情
type StatDetailResp struct {
	Date string `json:"date"`
	Pv int64 `json:"pv"`
	Uv int64 `json:"uv"`
}
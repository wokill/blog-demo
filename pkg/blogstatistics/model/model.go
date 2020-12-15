package model

const (
	PV = "pv"
	UV = "uv"
)

// Statistics 统计数据
type Statistics struct {
	ID        string `json:"id" bson:"_id"`
	BlogID    string `json:"blogId" bson:"blogId"`       // 博客ID
	Typ       string `json:"type" bson:"type"`           // type 类型 pv 统计为pv的数据， uv统计uv的数据
	Author    string `json:"author" bson:"author"`       // 作者ID
	View      int    `json:"view" bson:"view"`           // 查看次数
	CreatedAt int64  `json:"createdAt" bson:"createdAt"` // 统计时间
}

// ViewLogs 阅读记录
type ViewLogs struct {
	ID        string `json:"id" bson:"_id"`
	BlogID    string `json:"blogId" bson:"blogId"`     // 博客ID
	UserID    string `json:"userId" bson:"userId"`     // 用户ID 为 "" 时代表访客
	ClientIP  string `json:"clientIp" bson:"clientIp"` // 客户端IP
	Year      int    `json:"year" bson:"year"`         // 年份
	Month     int    `json:"month" bson:"month"`       // 月份
	Day       int    `json:"day" json:"day"`
	LookAt    int64  `json:"lookAt" bson:"lookAt"` // 阅读时间
	CreatedAt int64  `json:"createdAt" bson:"createdAt"`
}

// StatBlog 博客统计信息
type StatBlog struct {
	BlogID string `bson:"blogId"` // 博客ID
	Typ string `bson:"typ"` // 类型
	View int64 `bson:"view"` //浏览量
}

// StatOverview 统计概览数据
type StatOverview struct {
	BlogID string `bson:"_id"`
	Stat []struct{
		Uv int64 `bson:"uv"`
		Pv int64 `bson:"pv"`
		Typ string `bson:"typ"`
	} `bson:"stat"`
}

// StatBlogDetail 统计不同时间维度博客的趋势
type StatBlogDetail struct {
	Date string `bson:"_id"`
	Stat []struct{
		Uv int64 `bson:"uv"`
		Pv int64 `bson:"pv"`
		Typ string `bson:"typ"`
	} `bson:"stat"`
}
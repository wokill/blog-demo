package model

// Blog 博文
type Blog struct {
	ID        string   `bson:"_id" json:"id"`
	Title     string   `json:"title" bson:"title"`         // 标题
	Content   string   `json:"content" bson:"content"`     //内容
	Author    string   `json:"author" bson:"author"`       //作者ID数据
	Tags      []string `json:"tags" bson:"tags"`           // tags
	CreatedAt int64    `json:"createdAt" bson:"createdAt"` // 创建时间
	Uv        int64    `json:"uv" bson:"-"`
	Pv        int64    `json:"pv" bson:"-"`
}


type IncBlogStat struct {
	ID string
	Type string
	View int
}
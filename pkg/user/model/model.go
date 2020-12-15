package model

// User 用户数据
type User struct {
	ID        string `json:"id" bson:"_id"`
	Account   string `json:"account" bson:"account"`
	Name      string `json:"name" bson:"name"`
	Password  string `json:"password" bson:"password"`
	Avatar    string `json:"avatar" bson:"avatar"`
	LoginAt   int64  `json:"loginAt" bson:"loginAt"`
	LogoutAt  int64  `json:"logoutAt" bson:"logoutAt"`
	CreatedAt int64  `json:"createdAt" bson:"createdAt"`
}

// UserSession 用户登录信息
type UserSession struct {
	ID      string
	Account string
	Token   string
	LoginAt int64
}

// Visitor 浏览信息
type Visitor struct {
	UserID    string `json:"user_id"`    // 用户ID
	Client    string `json:"client"`     // 客户端
	ArticleID string `json:"article_id"` // 文章ID
	ViewAt    int64  `json:"view_at"`    // 浏览时间
}

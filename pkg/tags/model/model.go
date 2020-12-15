package model

// Tag 标签
type Tag struct {
	ID   string `json:"id" bson:"_id"`
	Name string `json:"name" bson:"name"`
}

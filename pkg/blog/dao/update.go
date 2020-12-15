package dao

import (
	"git.dustess.com/mk-training/mk-blog-svc/pkg/blog/model"
	"go.mongodb.org/mongo-driver/bson"
)

// UpdateOne 更新
func (m *BlogDao) UpdateOne(filter, data interface{}) (int64, error) {
	result, err :=m.dao.UpdateOne(m.ctx, filter, bson.M{"$set": data})
	if err != nil {
		return 0, err
	}
	return result.ModifiedCount, nil
}

// UpdateBlog 更新博客
func (m *BlogDao) UpdateBlog (id string, blog model.Blog) error {
	filter := bson.M{"_id": id}
	_, err := m.UpdateOne(filter, blog)
	return err
}
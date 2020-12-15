package dao

import (
	"git.dustess.com/mk-base/util/crypto"
	"git.dustess.com/mk-training/mk-blog-svc/pkg/blog/model"
	"time"
)

// InsertBlog 插入博客
func (m *BlogDao) InsertOne(data model.Blog) (string, error) {
	if data.ID == "" {
		data.ID = crypto.UUID()
	}
	data.CreatedAt = time.Now().Unix()
	r, err := m.dao.InsertOne(m.ctx, data)
	var id string
	if r != nil {
		result := r.InsertedID
		id = result.(string)
	}
	return id, err
}

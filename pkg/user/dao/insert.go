package dao

import (
	"git.dustess.com/mk-base/util/crypto"
	"git.dustess.com/mk-training/mk-blog-svc/pkg/user/model"
)

// InsertTags 插入用户数据
func (m *UserDao) InsertUser(data model.User) (string, error) {
	if data.ID == "" {
		data.ID = crypto.UUID()
	}
	r, err := m.dao.InsertOne(m.ctx, data)
	var id string
	if r != nil {
		result := r.InsertedID
		id = result.(string)
	}
	return id, err
}

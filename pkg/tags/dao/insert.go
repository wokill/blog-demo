package dao

import (
	"errors"
	"git.dustess.com/mk-base/util/crypto"
	"git.dustess.com/mk-training/mk-blog-svc/pkg/tags/model"
	"go.mongodb.org/mongo-driver/bson"
)

// InsertTags 插入标签数据
func (m *TagDao) InsertTags(data []string) ([]string, error) {
	if len(data) == 0 {
		return nil, errors.New("tags is empty")
	}
	var (
		updateData []model.Tag
		filter bson.M
		tags []string
		insertData []interface{}
		id []string
	)
	for _, v := range data {
		d := model.Tag{
			ID:   crypto.UUID(),
			Name: v,
		}
		updateData = append(updateData,d)
		tags = append(tags, v)
	}
	filter = bson.M{"name": bson.M{"$in": tags}}
	t , _ := m.FindMany(filter)
	for _, v := range updateData {
		var b bool
		for _, vv := range t {
			if v.Name == vv.Name {
				b = true
				id = append(id, vv.ID)
				break
			}
		}
		if !b {
			insertData = append(insertData, v)
		}

	}
	if len(insertData) < 1 {
		return nil, errors.New("tags is empty")
	}

	r, err := m.dao.InsertMany(m.ctx, insertData)

	if r != nil {
		for _, v := range r.InsertedIDs {
			id = append(id, v.(string))
		}
	}
	return id, err
}

package dao

import (
	"errors"
	"git.dustess.com/mk-base/mongo-driver/mongo"
	"git.dustess.com/mk-training/mk-blog-svc/pkg/blog/model"
	"go.mongodb.org/mongo-driver/bson"
)

// FindOne 查询单条数据
func (m *BlogDao) FindOne(filter interface{}) (blog model.Blog, err error) {
	err = m.dao.FindOne(m.ctx, filter).Decode(&blog)
	return
}

// 通过ID查询
func (m *BlogDao) FindByID(ID string) (blog model.Blog, err error) {
	if ID == "" {
		err = errors.New("blog ID can`t empty")
		return
	}
	filter := bson.M{"_id": ID}
	return m.FindOne(filter)
}

// FindAll 查找全部数据
func (m *BlogDao) FindAll(offset, limit int64, filter interface{}) (result []model.Blog, err error) {
	var opt = &mongo.FindOptions{
		Sort: bson.M{"createdAt": -1},
	}
	opt.SetSkip(offset).SetLimit(limit)
	cursor, err := m.dao.Find(m.ctx, filter, opt)
	if err != nil {
		return nil, err
	}
	if err := cursor.All(m.ctx, &result); err != nil {
		return nil, err
	}
	return
}

// FindCount 查询总数
func (m *BlogDao) FindCount(filter interface{}) int64 {
	count, _ := m.dao.CountDocuments(m.ctx, filter)
	return count
}

// FindByIds 批量根据ID查询
func (m *BlogDao) FindByIds(blogID []string) (result []model.Blog, err error) {
	if len(blogID) == 0 {
		return
	}
	filter := bson.M{"_id": bson.M{"$in": blogID}}
	return m.FindAll(0, int64(len(blogID)), filter)
}
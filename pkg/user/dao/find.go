package dao

import (
	"errors"
	"git.dustess.com/mk-training/mk-blog-svc/pkg/user/model"
	"go.mongodb.org/mongo-driver/bson"
)

// FindOne 查询一条数据
func (m *UserDao) FindOne(filter interface{}) (result model.User, err error) {
	err = m.dao.FindOne(m.ctx, filter).Decode(&result)
	return
}

// FindByAccount 通过账号查询
func (m *UserDao) FindByAccount(account string) (result model.User, err error) {
	q := bson.M{"account": account}
	return m.FindOne(q)
}

// FindByID 通过用户ID查询用户
func (m *UserDao) FindByID(id string) (result model.User, err error) {
	q := bson.M{"_id": id}
	return m.FindOne(q)
}

// FindManyByID 批量查询
func (m *UserDao) FindManyByID(ids []string) (result []model.User, err error) {
	if len(ids) < 1 {
		return nil,errors.New("ids is empty")
	}
	filter := bson.M{"_id": bson.M{"$in": ids}}
	cursor, e := m.dao.Find(m.ctx, filter)
	if e != nil {
		return nil, err
	}
	if err := cursor.All(m.ctx, &result); err != nil {
		return nil, err
	}
	return result, nil
}

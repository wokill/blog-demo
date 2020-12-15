package dao

import (
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

// UpdateUser 修改用户数据
func (m *UserDao) UpdateUser(filter, data interface{}) (int64, error) {
	result, err :=m.dao.UpdateOne(m.ctx, filter, bson.M{"$set": data})
	if err != nil {
		return 0, err
	}
	return result.ModifiedCount, nil
}

// Login 登录
func (m *UserDao) Login(userID string) (err error)  {
	filter := bson.M{"_id": userID}
	data := bson.M{"loginAt": time.Now().Unix()}
	_, err = m.UpdateUser(filter, data)
	return
}

// Logout 登录
func (m *UserDao) Logout(userID string) (err error)  {
	filter := bson.M{"_id": userID}
	data := bson.M{"logoutAt": time.Now().Unix()}
	_, err = m.UpdateUser(filter, data)
	return
}
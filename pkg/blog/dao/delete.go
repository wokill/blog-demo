package dao

import "go.mongodb.org/mongo-driver/bson"

// Delete 删除
func (m *BlogDao) Delete(filter interface{}) error {
	_, err := m.dao.DeleteOne(m.ctx, filter)
	return err
}

// DeleteBuID 通过ID删除
func (m *BlogDao) DeleteByID (ID string) error  {
	f := bson.M{"_id": ID}
	return m.Delete(f)
}

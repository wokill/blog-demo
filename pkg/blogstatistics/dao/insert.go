package dao

import (
	"errors"
	"git.dustess.com/mk-training/mk-blog-svc/pkg/blogstatistics/model"
)

// ViewLogInsertOne 插入单条数据
func (m *BlogLogDao) ViewLogInsertOne(data model.ViewLogs) (string, error) {
	result, err := m.dao.InsertOne(m.ctx, data)
	if err != nil {
		return "", err
	}
	return result.InsertedID.(string), nil
}

// InsertMany 批量插入
func (m *BlogStatDao) InsertMany(insertData []model.Statistics) error {
	if len(insertData) == 0 {
		return  errors.New("data is empty")
	}
	var _d = make([]interface{}, len(insertData))
	for k, v := range insertData {
		_d[k] = v
	}
	_, err := m.dao.InsertMany(m.ctx, _d)
	return err
}
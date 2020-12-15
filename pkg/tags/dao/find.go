package dao

import "git.dustess.com/mk-training/mk-blog-svc/pkg/tags/model"

// FindOne 查询一条数据
func (m *TagDao) FindOne(filter interface{}) (result model.Tag, err error) {
	err = m.dao.FindOne(m.ctx, filter).Decode(&result)
	return
}

// FindMany 批量查询
func (m *TagDao) FindMany(filter interface{})(result []model.Tag, err error) {
	cursor, e := m.dao.Find(m.ctx, filter)
	if e != nil {
		return nil, err
	}
	if err := cursor.All(m.ctx, &result); err != nil {
		return nil, err
	}
	return result, nil
}
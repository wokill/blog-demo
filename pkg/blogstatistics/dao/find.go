package dao

import (
	"git.dustess.com/mk-base/mongo-driver/mongo"
	"git.dustess.com/mk-training/mk-blog-svc/pkg/blogstatistics/model"
	"go.mongodb.org/mongo-driver/bson"
)

// AsynStaticsByBlog 统计数据（非实时）
func (m *BlogStatDao) AsynStaticsByBlog(filter interface{}) (result []model.StatBlog) {
	pipe := []bson.M{
		{"$match": filter},
		{"$group": bson.M{"_id": bson.M{"blogID": "$blogId", "type": "$type"}, "view": bson.M{"$sum": "$view"}}},
		{"$project": bson.M{"blogId": "$_id.blogID", "typ": "$_id.type", "view": 1}},
	}
	cursor, err := m.dao.Aggregate(m.ctx, pipe)
	if err != nil {
		return nil
	}

	err = cursor.All(m.ctx, &result)
	if err != nil {
		return nil
	}
	return
}

// BlogStatQuery 统计条件
func (m *BlogStatDao) BlogStatQuery(blogID []string, start, end int64, typ string) interface{} {
	q := bson.M{}
	if len(blogID) > 0 {
		q["blogId"] = bson.M{"$in": blogID}
	}
	if typ != "" {
		q["type"] = typ
	}
	tq := bson.M{}
	if start > 0 {
		tq["$gte"] = start
	}
	if end > 0 {
		tq["$lte"] = end
	}
	if len(tq) > 0 {
		q["createdAt"] = tq
	}
	return q
}

// FindLast 查找最近统计数据
func (m *BlogStatDao) FindLast(typ string) (stat model.Statistics, err error) {
	filter := bson.M{"type": typ}
	op := &mongo.FindOneOptions{
		Sort: bson.M{"createdAt": -1},
	}
	err = m.dao.FindOne(m.ctx, filter, op).Decode(&stat)
	return
}

// StatisticsByBlog 统计博客信息
func (m *BlogStatDao) StatisticsByBlog(filter interface{}, typ string, limit int64) (result []model.StatOverview) {
	sort := bson.M{}
	if typ == model.PV {
		sort["stat.pv"] = -1
	} else {
		sort["stat.uv"] = -1
	}
	sort["_id"] = 1
	pipe := []bson.M{
		{"$match": filter},
		{"$group": bson.M{"_id": bson.M{"blogID": "$blogId", "type": "$type"}, "uv": bson.M{"$sum": "$view"},  "pv": bson.M{"$sum": "$view"}}},
		{"$project": bson.M{"blogId": "$_id.blogID", "typ": "$_id.type", "uv": 1, "pv": 1, "_id": 0}},
		{"$group": bson.M{"_id": "$blogId", "stat": bson.M{"$push": bson.M{"uv": "$uv", "pv":"$pv","typ": "$typ"}}}},
		{"$sort": sort},
		{"$limit": limit},
	}
	cursor, err := m.dao.Aggregate(m.ctx, pipe)
	if err != nil {
		return nil
	}

	err = cursor.All(m.ctx, &result)
	if err != nil {
		return nil
	}
	return
}

// StaticsDetail 统计详情
func (m *BlogStatDao) StaticsDetail(blogID string, start,end int64) (result []model.StatBlogDetail) {
	filter := bson.M{"blogId": blogID, "createdAt": bson.M{"$gte": start, "$lt": end}}
	dateBson := bson.M{"$dateToString": bson.M{"format": "%Y-%m-%d", "date": bson.M{"$toDate": bson.M{"$add": bson.A{28800000, bson.M{"$multiply": bson.A{"$createdAt", 1000}}}}}}}
	pipe := []bson.M{
		{"$match": filter},
		{"$project": bson.M{"date": dateBson, "blogId": 1, "type":1,"view":1}},
		{"$group": bson.M{"_id": bson.M{"date": "$date", "type": "$type"}, "uv": bson.M{"$sum": "$view"},  "pv": bson.M{"$sum": "$view"}}},
		{"$project": bson.M{"date": "$_id.date", "typ": "$_id.type", "uv": 1, "pv": 1, "_id": 0}},
		{"$group": bson.M{"_id": "$date", "stat": bson.M{"$push": bson.M{"uv": "$uv", "pv":"$pv","typ": "$typ"}}}},
		{"$sort": bson.M{"_id": 1}},
	}
	cursor, err := m.dao.Aggregate(m.ctx, pipe)
	if err != nil {
		return nil
	}

	err = cursor.All(m.ctx, &result)
	if err != nil {
		return nil
	}
	return
}


// FindLogsByTime 统计数据
func (m *BlogLogDao) FindLogsByTime(at int64) (result []model.ViewLogs, err error) {
	filter := bson.M{"createdAt": bson.M{"$gte": at}}
	cursor, err := m.dao.Find(m.ctx, filter)
	if err != nil {
		return nil, err
	}
	err = cursor.All(m.ctx, &result)
	return
}

// FindLogs 批量查询
func (m *BlogLogDao) FindLogs(filter interface{}) (result []model.ViewLogs, err error) {
	cursor, err := m.dao.Find(m.ctx, filter)
	if err != nil {
		return nil, err
	}
	err = cursor.All(m.ctx, &result)
	return
}

// FindLogsByUserIDs 通过UserID查询
func (m *BlogLogDao) FindLogsByUserIDs(userIDs []string, blogID string, at int64) (result []model.ViewLogs, err error) {
	if len(userIDs) < 1 {
		return
	}
	return m.FindLogs(bson.M{"userId": bson.M{"$in": userIDs}, "blogId": blogID, "createdAt": bson.M{"$lt": at}})
}

// FindLogsByIps 通过IP查询
func (m *BlogLogDao) FindLogsByIps(ips []string, blogID string, at int64) (result []model.ViewLogs, err error) {
	if len(ips) < 1 {
		return
	}
	return m.FindLogs(bson.M{"clientIp": bson.M{"$in": ips}, "blogId": blogID, "createdAt": bson.M{"$lt": at}})
}

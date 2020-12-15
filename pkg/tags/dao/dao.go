package dao

import (
	"context"
	"git.dustess.com/mk-base/mongo-driver/mongo"
	"git.dustess.com/mk-training/mk-blog-svc/config"
)

const collName = "tags"

// TagDao 客数据可连接
type TagDao struct {
	dao       *mongo.Dao
	ctx       context.Context
}

// NewTagDao 创建对象
func NewTagDao(ctx context.Context) *TagDao {
	return &TagDao{
		dao:       mongo.NewDao(mongo.MKBiz, config.Get().Mongo.MongoMK.MKDB.Name, collName),
		ctx:       ctx,
	}
}

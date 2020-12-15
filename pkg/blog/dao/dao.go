package dao

import (
	"context"
	"git.dustess.com/mk-base/mongo-driver/mongo"
	"git.dustess.com/mk-training/mk-blog-svc/config"
)

const collName = "blog"

// BlogDao 客数据可连接
type BlogDao struct {
	dao *mongo.Dao
	ctx context.Context
}

// NewBlogDao 创建对象
func NewBlogDao(ctx context.Context) *BlogDao {
	return &BlogDao{
		dao: mongo.NewDao(mongo.MKBiz, config.Get().Mongo.MongoMK.MKDB.Name, collName),
		ctx: ctx,
	}
}

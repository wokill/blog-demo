package dao

import (
"context"
"git.dustess.com/mk-base/mongo-driver/mongo"
"git.dustess.com/mk-training/mk-blog-svc/config"
)

const collName = "blog_statistics"
const logCollName = "blog_log"

// BlogDao 客数据可连接
type BlogStatDao struct {
	dao *mongo.Dao
	ctx context.Context
}

// NewBlogStatDao 创建对象
func NewBlogStatDao(ctx context.Context) *BlogStatDao {
	return &BlogStatDao{
		dao: mongo.NewDao(mongo.MKBiz, config.Get().Mongo.MongoMK.MKDB.Name, collName),
		ctx: ctx,
	}
}


// BlogLogDao 客数据可连接
type BlogLogDao struct {
	dao *mongo.Dao
	ctx context.Context
}

// NewBlogStatDao 创建对象
func NewBlogLogDao(ctx context.Context) *BlogLogDao {
	return &BlogLogDao{
		dao: mongo.NewDao(mongo.MKBiz, config.Get().Mongo.MongoMK.MKDB.Name, logCollName),
		ctx: ctx,
	}
}

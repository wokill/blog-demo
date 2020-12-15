package dao

import (
	"context"
	"git.dustess.com/mk-base/mongo-driver/mongo"
	"git.dustess.com/mk-training/mk-blog-svc/config"
)

const collName = "user"

// UserDao 客数据可连接
type UserDao struct {
	dao       *mongo.Dao
	ctx       context.Context
}

// NewUserDao 创建对象
func NewUserDao(ctx context.Context) *UserDao {
	return &UserDao{
		dao:       mongo.NewDao(mongo.MKBiz, config.Get().Mongo.MongoMK.MKDB.Name, collName),
		ctx:       ctx,
	}
}

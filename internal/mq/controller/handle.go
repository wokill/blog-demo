package controller

import (
	"context"
	"git.dustess.com/mk-base/kafka-driver/consumer"
	"git.dustess.com/mk-base/log"
	"git.dustess.com/mk-training/mk-blog-svc/config"
	"git.dustess.com/mk-training/mk-blog-svc/pkg/blogstatistics/dao"
)

// NewConsumer [...]
func NewConsumer() *Consumer {
	return &Consumer{
		dao: dao.NewBlogLogDao(context.Background()),
	}
}

// CustomerView 消费浏览数据
func CustomerView() {
	conf := config.Get().Kafka
	handler := NewConsumer()
	customer := consumer.NewConsumer(conf.Addrs, conf.Topic.MKWat.Name, conf.Group.MKCof, handler)
	err := customer.Consume().Error()
	log.Error("init kafka consumer fail. error : ", err)
}

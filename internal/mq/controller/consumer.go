package controller

import (
	"encoding/json"
	"git.dustess.com/mk-base/util/crypto"
	"git.dustess.com/mk-training/mk-blog-svc/pkg/blogstatistics/dao"
	statModel "git.dustess.com/mk-training/mk-blog-svc/pkg/blogstatistics/model"
	"git.dustess.com/mk-training/mk-blog-svc/pkg/user/model"
	"github.com/Shopify/sarama"
	"time"
)

// Consumer 消费
type Consumer struct {
	dao *dao.BlogLogDao
}

func (c *Consumer) Setup( sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup 清理数据
func (c *Consumer) Cleanup( sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim 消费数据
func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		session.MarkMessage(message, "")
		var visitor model.Visitor
		err := json.Unmarshal(message.Value, &visitor)
		if err != nil {
			continue
		}
		c.Flush(visitor)
	}
	return nil
}

// Flush 写入数据库
func (c *Consumer) Flush(visitor model.Visitor) {
	if visitor.ArticleID == "" {
		return
	}
	date := time.Unix(visitor.ViewAt, 0)
	insertData := statModel.ViewLogs{
		ID:        crypto.UUID(),
		BlogID:    visitor.ArticleID,
		UserID:    visitor.UserID,
		ClientIP:  visitor.Client,
		Year:      date.Year(),
		Month:     int(date.Month()),
		Day:       date.Day(),
		LookAt:    visitor.ViewAt,
		CreatedAt: visitor.ViewAt,
	}
	_, _ = c.dao.ViewLogInsertOne(insertData)
}


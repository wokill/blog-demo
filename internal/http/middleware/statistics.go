package middleware

import (
	"git.dustess.com/mk-training/mk-blog-svc/pkg/common"
	"git.dustess.com/mk-training/mk-blog-svc/pkg/mq/producer"
	"github.com/gin-gonic/gin"
)

// StatisticsVisitor 统计访问量
func StatisticsVisitor(c *gin.Context) {
	c.Next()
	var (
		uid string
		ip  = c.ClientIP()
	)
	articleID := c.Param("article_id")
	if articleID == "" {
		return
	}
	sess, err := common.GetSession(c)
	if err == nil {
		uid = sess.ID
	}
	go func() {
		producer.SendVisitorLog(uid, articleID, ip)
	}()
}

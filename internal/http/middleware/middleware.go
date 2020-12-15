package middleware

import (
	"git.dustess.com/mk-base/gin-ext/extend"
	"git.dustess.com/mk-training/mk-blog-svc/pkg/common"
	"git.dustess.com/mk-training/mk-blog-svc/pkg/user/dao"
	"git.dustess.com/mk-training/mk-blog-svc/pkg/user/model"
	"github.com/gin-gonic/gin"
	"strings"
)

// Authorization 身份认证
func Authorization(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.Abort()
		extend.SendUnauthorized(c)
		return
	}
	cache := common.NewCache(common.SessionKey, int64(common.SessionExpired))
	var session  model.UserSession
	err := cache.FindJSON(strings.TrimSpace(token), &session)
	if err != nil {
		c.Abort()
		extend.SendUnauthorized(c)
		return
	}
	m := dao.NewUserDao(c)
	user, err := m.FindByID(session.ID)
	if err != nil {
		c.Abort()
		extend.SendUnauthorized(c)
		return
	}
	if user.LogoutAt > session.LoginAt {
		c.Abort()
		extend.SendUnauthorized(c)
		return
	}
	c.Set(common.SessionHeaderKey, session)
	c.Next()
}

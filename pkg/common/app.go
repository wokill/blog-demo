package common

import (
	"errors"
	"git.dustess.com/mk-base/gin-ext/api"
	"git.dustess.com/mk-training/mk-blog-svc/pkg/user/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

const (
	SessionHeaderKey = "session"
)

type SuccessRes api.SuccessRes

type FailedRes api.FailedRes

// SendOK 成功
func SendOK(ctx *gin.Context, data interface{}) {
	res := api.DefaultRes()
	res.Code = 200
	if data == nil {
		data = make(map[string]struct{}, 0)
	}
	res.Data = data
	ctx.JSON(http.StatusOK, res)
}

// Session 获取 session
func Session(ctx *gin.Context) (session model.UserSession, err error) {
	s, ok := ctx.Get(SessionHeaderKey)
	if !ok {
		err = errors.New("not found user")
		return
	}
	session = s.(model.UserSession)
	return
}

//limit 操作
func Limit(ctx *gin.Context) int64 {
	limit := ctx.DefaultQuery("pageSize", "20")
	num, _ := strconv.Atoi(limit)
	return int64(num)
}

//offset 操作
func Offset(ctx *gin.Context) int64 {
	offset := ctx.DefaultQuery("page", "1")
	page, _ := strconv.Atoi(offset)
	if page < 1 {
		page = 1
	}
	return int64(page-1) * Limit(ctx)
}

// 获取session
func GetSession(ctx *gin.Context) (session model.UserSession, err error) {
	token := ctx.GetHeader("Authorization")
	if token == "" {
		return
	}
	cache := NewCache(SessionKey, int64(SessionExpired))
	err = cache.FindJSON(strings.TrimSpace(token), &session)
	return
}

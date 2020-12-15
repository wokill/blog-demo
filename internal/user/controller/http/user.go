package http

import (
	"errors"
	"git.dustess.com/mk-base/gin-ext/extend"
	"git.dustess.com/mk-base/log"
	"git.dustess.com/mk-training/mk-blog-svc/pkg/common"
	"git.dustess.com/mk-training/mk-blog-svc/pkg/user/service"
	"github.com/gin-gonic/gin"
)

// LoginReq 登录数据
type LoginReq struct {
	Account string `json:"account" form:"account" binding:"required"`
	Pass    string `json:"pass" form:"pass" binding:"required"`
}

// LoginResp 登录返回参数
type LoginResp struct {
	Token  string `json:"token"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

// Login 登录
// @Summary 新增批量修改客户类型
// @Description 新增批量修改客户类型
// @Tags user
// @Accept  json
// @Produce json
// @Param body body LoginReq true "params"
// @Success 200 {object} common.SuccessRes"请求成功"
// @Failure 400 {object} common.FailedRes "参数解析错误"
// @Failure 500 {object} common.FailedRes "服务器内部错误"
// @Router /visitor/login  [POST]
func Login(c *gin.Context) {
	var q LoginReq
	if err := c.ShouldBind(&q); err != nil {
		extend.SendParamParseError(c)
		return
	}
	svc := service.NewUserService(c.Request.Context())
	token, user, err := svc.Login(q.Account, q.Pass)
	if err != nil {
		extend.SendData(c, nil, err)
		return
	}
	common.SendOK(c, LoginResp{
		Token:  token,
		Name:   user.Name,
		Avatar: user.Avatar,
	})
}

// 注册账号
type RegisterReq struct {
	Account string `json:"account" form:"account" binding:"required"`
	Pass    string `json:"pass" form:"pass" binding:"required"`
	Name    string `json:"name" form:"name" binding:"required"`
}

// Register 注册
// @Summary 注册账号
// @Description 注册账号
// @Tags user
// @Accept  json
// @Produce json
// @Param body body RegisterReq true "params"
// @Success 200 {object} common.SuccessRes"请求成功"
// @Failure 400 {object} common.FailedRes "参数解析错误"
// @Failure 500 {object} common.FailedRes "服务器内部错误"
// @Router /visitor/register  [POST]
func Register(c *gin.Context) {
	var q RegisterReq
	if err := c.ShouldBind(&q); err != nil {
		extend.SendParamParseError(c)
		return
	}
	svc := service.NewUserService(c.Request.Context())
	_, err := svc.Register(q.Account, q.Name, q.Pass)
	if err != nil {
		log.Error("register user error: ", err)
		extend.SendData(c, nil, errors.New("注册用户失败"))
		return
	}
	common.SendOK(c, nil)
}

// Logout 退出登录
// @Summary 退出登录
// @Description 退出登录
// @Tags user
// @Accept  json
// @Produce json
// @Param Authorization header string true "认证信息 eg:xxxx-xxxx-xxxx-xxx"
// @Success 200 {object} common.SuccessRes"请求成功"
// @Failure 400 {object} common.FailedRes "参数解析错误"
// @Failure 500 {object} common.FailedRes "服务器内部错误"
// @Router /visitor/logout  [PUT]
func Logout(c *gin.Context) {
	session, err := common.Session(c)
	if err != nil {
		extend.SendUnauthorized(c)
		return
	}
	svc := service.NewUserService(c)
	err = svc.Logout(session)
	if err != nil {
		log.Error("user logout error: ", err)
		extend.SendData(c, nil, errors.New("退出登录失败"))
		return
	}
	common.SendOK(c, nil)
}

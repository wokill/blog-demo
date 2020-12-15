package http

import (
	"git.dustess.com/mk-base/gin-ext/extend"
	"git.dustess.com/mk-training/mk-blog-svc/pkg/blog/model"
	"git.dustess.com/mk-training/mk-blog-svc/pkg/blog/service"
	"git.dustess.com/mk-training/mk-blog-svc/pkg/common"
	"github.com/gin-gonic/gin"
)

type CreateBlogReq struct {
	Title   string   `json:"title" form:"title" binding:"required"`
	Content string   `json:"content" form:"content" binding:"required"`
	Tags    []string `json:"tags"`
}

// CreateBlog 创建博客
// @Summary 创建博客
// @Description 创建博客
// @Tags blog
// @Accept  json
// @Produce json
// @Param body body CreateBlogReq true "params"
// @Param Authorization header string true "认证信息 eg:xxxx-xxxx-xxxx-xxx"
// @Success 200 {object} common.SuccessRes"请求成功"
// @Failure 400 {object} common.FailedRes "参数解析错误"
// @Failure 500 {object} common.FailedRes "服务器内部错误"
// @Router /blogger/article  [POST]
func CreateBlog(c *gin.Context) {
	session, err := common.Session(c)
	if err != nil {
		extend.SendUnauthorized(c)
		return
	}
	var q CreateBlogReq
	if err := c.Bind(&q); err != nil {
		extend.SendParamParseError(c)
		return
	}
	data := model.Blog{
		Title:   q.Title,
		Content: q.Content,
		Author:  session.ID,
	}
	svc := service.NewBlogService(c)
	err = svc.CreateBlog(data, q.Tags)
	if err != nil {
		extend.SendData(c, nil, err)
		return
	}
	common.SendOK(c, nil)
}

// CreateBlog 修改博客
// @Summary 修改博客
// @Description 修改博客
// @Tags blog
// @Accept  json
// @Produce json
// @Param body body CreateBlogReq true "params"
// @Param article_id path string true "博客ID"
// @Param Authorization header string true "认证信息 eg:xxxx-xxxx-xxxx-xxx"
// @Success 200 {object} common.SuccessRes"请求成功"
// @Failure 400 {object} common.FailedRes "参数解析错误"
// @Failure 500 {object} common.FailedRes "服务器内部错误"
// @Router /blogger/article/{article_id} [PUT]
func UpdateBlog(c *gin.Context) {
	session, err := common.Session(c)
	if err != nil {
		extend.SendUnauthorized(c)
		return
	}
	id := c.Param("article_id")
	if id == "" {
		extend.SendParamParseError(c)
		return
	}
	var q CreateBlogReq
	if err := c.Bind(&q); err != nil {
		extend.SendParamParseError(c)
		return
	}
	data := model.Blog{
		Title:   q.Title,
		Content: q.Content,
		Author:  session.ID,
		Tags:    q.Tags,
	}
	svc := service.NewBlogService(c)
	err = svc.UpdateBlog(session, id, data)
	if err != nil {
		extend.SendData(c, nil, err)
		return
	}
	common.SendOK(c, nil)
}

// CreateBlog 删除博客
// @Summary 删除博客
// @Description 删除博客
// @Tags blog
// @Accept  json
// @Produce json
// @Param body body CreateBlogReq true "params"
// @Param article_id path string true "博客ID"
// @Param Authorization header string true "认证信息 eg:xxxx-xxxx-xxxx-xxx"
// @Success 200 {object} common.SuccessRes"请求成功"
// @Failure 400 {object} common.FailedRes "参数解析错误"
// @Failure 500 {object} common.FailedRes "服务器内部错误"
// @Router /blogger/article/{article_id} [DELETE]
func DeleteBlog(c *gin.Context) {
	session, err := common.Session(c)
	if err != nil {
		extend.SendUnauthorized(c)
		return
	}
	id := c.Param("article_id")
	if id == "" {
		extend.SendParamParseError(c)
		return
	}
	svc := service.NewBlogService(c)
	err = svc.Delete(session, id)
	if err != nil {
		extend.SendData(c, nil, err)
		return
	}
	common.SendOK(c, nil)
}


type BlogListResp struct {
	Count int64 `json:"count"`
	Items []model.BlogList `json:"items"`
}

// ListsBlog 查看我的博客
// @Summary 查看我的博客
// @Description 查看我的博客
// @Tags blog
// @Accept  json
// @Produce json
// @Param body body CreateBlogReq true "params"
// @Param Authorization header string true "认证信息 eg:xxxx-xxxx-xxxx-xxx"
// @Success 200 {object} BlogListResp "请求成功"
// @Failure 400 {object} common.FailedRes "参数解析错误"
// @Failure 500 {object} common.FailedRes "服务器内部错误"
// @Router /blogger/article [GET]
func ListsBlog(c *gin.Context) {
	var (
		limit = common.Limit(c)
		offset = common.Offset(c)
	)
	session, err := common.Session(c)
	if err != nil {
		extend.SendUnauthorized(c)
		return
	}
	svc := service.NewBlogService(c)
	resp, count := svc.Lists(session, offset, limit)
	common.SendOK(c, BlogListResp{
		Count: count,
		Items: resp,
	})
}

package http

import (
	"git.dustess.com/mk-base/gin-ext/extend"
	"git.dustess.com/mk-training/mk-blog-svc/pkg/blog/model"
	"git.dustess.com/mk-training/mk-blog-svc/pkg/blog/service"
	"git.dustess.com/mk-training/mk-blog-svc/pkg/common"
	"github.com/gin-gonic/gin"
)

// BlogDetail 博客详情
// @Summary 博客详情
// @Description 博客详情
// @Tags blog
// @Accept  json
// @Produce json
// @Param article_id path string true "博客ID"
// @Param Authorization header string false "认证信息 eg:xxxx-xxxx-xxxx-xxx"
// @Success 200 {object} model.BlogDetail "请求成功"
// @Failure 400 {object} common.FailedRes "参数解析错误"
// @Failure 500 {object} common.FailedRes "服务器内部错误"
// @Router /visitor/article/{article_id} [GET]
func BlogDetail(ctx *gin.Context) {
	id := ctx.Param("article_id")
	if id == "" {
		extend.SendParamParseError(ctx)
		return
	}
	svc := service.NewBlogService(ctx)
	resp, err := svc.Detail(id)
	if err != nil {
		extend.SendData(ctx, nil, err)
	}
	common.SendOK(ctx, resp)
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
// @Param body body model.ListSearch true "params"
// @Success 200 {object} BlogListResp "请求成功"
// @Failure 400 {object} common.FailedRes "参数解析错误"
// @Failure 500 {object} common.FailedRes "服务器内部错误"
// @Router /visitor/articles [POST]
func BlogLists(ctx *gin.Context) {
	var q model.ListSearch
	_ = ctx.Bind(&q)
	svc := service.NewBlogService(ctx)
	result ,count := svc.VisitorLists(q)
	common.SendOK(ctx, BlogListResp{
		Count: count,
		Items: result,
	})
}

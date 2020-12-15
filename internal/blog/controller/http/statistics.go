package http

import (
	"git.dustess.com/mk-base/gin-ext/extend"
	"git.dustess.com/mk-training/mk-blog-svc/pkg/blogstatistics/model"
	"git.dustess.com/mk-training/mk-blog-svc/pkg/blogstatistics/service"
	"git.dustess.com/mk-training/mk-blog-svc/pkg/common"
	"github.com/gin-gonic/gin"
	"time"
)

// StatReq 统计参数
type StatReq struct {
	Start int64 `json:"start"`
	End int64 `json:"end"`
	Sort string `json:"sort"`
}

// Statistics 统计数据
// @Summary 统计数据
// @Description 统计数据
// @Tags statistics
// @Accept  json
// @Produce json
// @Param start query integer false "开始时间"
// @Param end query integer false "结束时间时间"
// @Param sort query string false "排序，uv按UV排， pv按PV排序，默认 pv"
// @Success 200 {object} model.StateOverViewResp "请求成功"
// @Failure 400 {object} common.FailedRes "参数解析错误"
// @Failure 500 {object} common.FailedRes "服务器内部错误"
// @Router /blogger/article [GET]
func Statistics(c *gin.Context) {
	var q StatReq
	_ = c.BindQuery(&q)
	if q.Sort != model.UV {
		q.Sort = model.PV
	}
	n := time.Now()
	if q.Start == 0 {
		q.Start = n.AddDate(0,0,-7).Unix()
	}
	if q.End == 0 {
		q.End = n.Unix()
	}
	svc := service.NewBlogService(c)
	resp := svc.BlogHot(q.Sort,q.Start,q.End)
	common.SendOK(c, resp)
}

// BlogTrend 博客趋势
// @Summary 博客趋势
// @Description 博客趋势
// @Tags statistics
// @Accept  json
// @Produce json
// @Param article_id path string true "文章ID"
// @Param start query integer false "开始时间"
// @Param end query integer false "结束时间时间"
// @Success 200 {object} model.StatDetailResp "请求成功"
// @Failure 400 {object} common.FailedRes "参数解析错误"
// @Failure 500 {object} common.FailedRes "服务器内部错误"
// @Router /blogger/article [GET]
func BlogTrend(c *gin.Context) {
	article := c.Param("article_id")
	if article == "" {
		extend.SendParamParseError(c)
		return
	}
	var q StatReq
	_ = c.BindQuery(&q)
	n := time.Now()
	if q.Start == 0 {
		q.Start = n.AddDate(0,0,-7).Unix()
	}
	if q.End == 0 {
		q.End = n.Unix()
	}
	svc := service.NewBlogService(c)
	resp := svc.BlogDetail(article, q.Start, q.End)
	common.SendOK(c, resp)
}

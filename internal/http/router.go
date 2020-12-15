package http

import (
	blogCtl "git.dustess.com/mk-training/mk-blog-svc/internal/blog/controller/http"
	"git.dustess.com/mk-training/mk-blog-svc/internal/http/middleware"
	visitorCtl "git.dustess.com/mk-training/mk-blog-svc/internal/user/controller/http"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	v1prefix         = "/demo/v1"
	v1InternalPrefix = "/internal/demo/v1"
)

// initRouter 初始化路由
func initRouter(router *gin.Engine) error {

	router.GET("/ready", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "ok")
	})

	router.GET("/healthy", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "ok")
	})

	router.GET(v1InternalPrefix+"/ping", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "pong")
	})
	//router.Use(middleware.Cache)
	visitor := router.Group("visitor")
	{
		visitor.POST("register", visitorCtl.Register)                                           // 	注册用户
		visitor.POST("login", visitorCtl.Login)                                                 //登录
		visitor.GET("article/:article_id", middleware.StatisticsVisitor, visitorCtl.BlogDetail) // 博客详情
		visitor.POST("articles",visitorCtl.BlogLists)                                          // 博客列表
		visitor.Use(middleware.Authorization)
		visitor.POST("logout", visitorCtl.Logout) // 登出
	}

	blog := router.Group("blogger", middleware.Authorization)
	{
		blog.POST("article", blogCtl.CreateBlog)               // 创建博客
		blog.PUT("article/:article_id", blogCtl.UpdateBlog)    // 修改博客
		blog.DELETE("article/:article_id", blogCtl.DeleteBlog) // 删除博客
		blog.GET("articles", blogCtl.ListsBlog)                 // 博客列表
	}

	article := router.Group("article")
	{
		article.GET("statistics/:article_id", blogCtl.BlogTrend) // 博客趋势
		article.GET("statistics", blogCtl.Statistics)             // 热度统计
	}

	return nil
}

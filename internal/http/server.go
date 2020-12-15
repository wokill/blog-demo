package http

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"git.dustess.com/mk-training/mk-blog-svc/config"

	"git.dustess.com/mk-base/gin-ext/constant"
	"git.dustess.com/mk-base/gin-ext/middleware"
	"git.dustess.com/mk-base/log"
	ginSwaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "git.dustess.com/mk-training/mk-blog-svc/api/swagger" // for swagger
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gopkg.in/go-playground/validator.v8"
)

var httpServer *http.Server

func router() http.Handler {
	r := gin.New()

	// 跨域中间件
	r.Use(middleware.CORS())

	// 请求日志
	r.Use(middleware.Logger())
	if gin.Mode() == gin.DebugMode {
		r.Use(gin.Logger())
	}

	// validator 信息翻译
	uni := middleware.NewZHUNI()
	binding.Validator = middleware.NewV9Validator(uni)
	if _, ok := binding.Validator.Engine().(*validator.Validate); ok {
	}
	// validator 错误处理
	r.Use(middleware.NewErrorHandler(uni).HandleErrors)

	// 错误恢复
	r.Use(gin.Recovery())

	err := initRouter(r)
	if err != nil {
		log.Error("init router", err)
		return r
	}

	// 错误恢复
	r.Use(gin.Recovery())

	// release 模式下不提供接口文档访问
	if string(constant.ReleaseMode) != os.Getenv("GIN_MODE") {
		r.GET(v1prefix+"/swagger/*any", ginSwagger.WrapHandler(ginSwaggerFiles.Handler))
	}

	return r
}

// Serve 启动服务
// @title 账户服务(mk-blog-svc)
// @version 1.0
// @description 示例服务
// @host mk-dev.dustess.com
// @BasePath /demo
func Serve() error {
	conf := *config.Get()

	log.Info("正在启动http服务，监听端口", conf.Server.Port)

	httpServer = &http.Server{
		Addr:         fmt.Sprintf("%s:%s", conf.Server.Host, conf.Server.Port),
		Handler:      router(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return httpServer.ListenAndServe()
}

// Shutdown 关闭服务
func Shutdown() error {
	log.Info("正在关闭http服务")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		return err
	}
	log.Info("http服务成功关闭")
	return nil
}

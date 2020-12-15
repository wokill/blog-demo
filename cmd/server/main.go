package main

import (
	"git.dustess.com/mk-training/mk-blog-svc/internal/mq/controller"
	"git.dustess.com/mk-training/mk-blog-svc/pkg/blogstatistics/service"
	"github.com/go-errors/errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"git.dustess.com/mk-base/es-driver/es"
	"git.dustess.com/mk-base/log"
	"git.dustess.com/mk-base/mongo-driver/mongo"
	"git.dustess.com/mk-base/redis-driver/redis"
	"git.dustess.com/mk-base/redis-driver/session"
	"git.dustess.com/mk-training/mk-blog-svc/config"
	"golang.org/x/sync/errgroup"

	// grpcServer "git.dustess.com/mk-training/mk-blog-svc/internal/grpc"
	httpServer "git.dustess.com/mk-training/mk-blog-svc/internal/http"
)

var (
	g errgroup.Group
)

func init() {

	err := config.Init()
	if err != nil {
		panic(errors.New(err))
	}
	conf := *config.Get()

	// 初始化mongo连接
	c, err := mongo.InitClient(conf.ToMongoMKBizConfig())
	if err != nil {
		panic(err)
	}
	if c == nil {
		panic(errors.New("mk_biz mongo 连接失败"))
	}

	c, err = mongo.InitClient(conf.ToMongoMKWatConfig())
	if err != nil {
		panic(err)
	}
	if c == nil {
		panic(errors.New("mk_wat mongo 连接失败"))
	}

	c, err = mongo.InitClient(conf.ToMongoWPConfig())
	if err != nil {
		panic(err)
	}
	if c == nil {
		panic(errors.New("wp mongo 连接失败"))
	}

	// 初始化redis
	_, err = redis.InitClient(conf.ToSessionConfig())
	if err != nil {
		panic(err)
	}

	_, err = redis.InitClient(conf.ToCacheConfig())
	if err != nil {
		panic(err)
	}

	_, err = session.InitSession()
	if err != nil {
		panic(err)
	}

	// 初始化es连接
	err = es.Init(conf.ToESConfig())
	if err != nil {
		panic(err)
	}
}

func main() {

	g.Go(func() error {
		if err := httpServer.Serve(); err != nil && err != http.ErrServerClosed {
			return err
		}
		return nil
	})

	// g.Go(func() error {
	// 	if err := grpcServer.Serve(); err != nil {
	// 		return err
	// 	}
	// 	return nil
	// })

	g.Go(func() error {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
		for {
			si := <-c
			switch si {
			case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
				return shutdown()
			case syscall.SIGHUP:
			default:
				return nil
			}
		}
	})
	go func() {
		controller.CustomerView()
	}()

	go func() {
		timer()
	}()

	if err := g.Wait(); err != nil {
		log.Error("服务运行失败: ", err)
		panic(err)
	}
}

func timer()  {
	t :=time.Minute
	timer := time.NewTicker(t)
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			svc := service.NewStat()
			_ = svc.StatisticsView()
		}
	}
}

// shutdown 关闭服务
func shutdown() error {
	// 关闭http服务
	if err := httpServer.Shutdown(); err != nil {
		return err
	}

	// // 关闭grpc服务
	// if err := grpcServer.Shutdown(); err != nil {
	// 	return err
	// }
	return nil
}

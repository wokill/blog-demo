package common

import (
	"git.dustess.com/mk-base/redis-driver/redis"
	"git.dustess.com/mk-training/mk-blog-svc/config"
	"testing"
	"time"
)

func TestLook(t *testing.T) {
	config.Init()
	conf := *config.Get()
	conf.Redis.Addr = "127.0.0.1:6379"
	conf.Redis.Password = "abcd1234"
	conf.Redis.CacheDB = 6
	conf.Redis.PoolSize = 5
	cli , _ :=redis.InitClient(conf.ToSessionConfig())
	d := redis.NewCao(cli)
	t1 := time.Minute * 3
	t2 := time.Hour * 5
	v2, err := d.AcqLock("tangping", true, t1, t2)
	t.Log(v2, err)
}

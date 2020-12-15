package common

import (
	"encoding/json"
	"errors"
	"git.dustess.com/mk-base/redis-driver/redis"
	"sync"
	"time"
)

const (
	SessionKey     = "bl.session"
	SessionExpired = time.Minute * 30
)

var _onceCache sync.Once
var redisClient *redis.Cao

// getCache 获取缓存句柄
func getCache() *redis.Cao {
	if redisClient != nil {
		return redisClient
	}

	_onceCache.Do(func() {
		redisClient = redis.NewCao(redis.Client(redis.MKCache))
	})

	return redisClient
}

type Cache struct {
	cao     *redis.Cao
	prefix  string
	expired int64
}

func NewCache(prefix string, expired int64) *Cache {
	return &Cache{
		cao:     getCache(),
		prefix:  prefix,
		expired: expired,
	}
}

// SetSession 存储session
func (c *Cache) SetSession(key, userID string) error {
	return c.cao.SetByTTL(c.prefix+key, userID, c.expired)
}

// Renew 续期
func (c *Cache) Renew(key string, t int64) error {
	return c.cao.SetTTL(key, t)
}

// RenewSession 为 session续期
func (c *Cache) RenewSession(key string) error {
	return c.Renew(c.prefix+key, c.expired)
}

// DelSession 删除 session
func (c *Cache) DelSession(key string) error {
	return c.cao.Del(c.prefix + key)
}

// CreateJSON 添加json缓存
func (c *Cache) CreateJSON(key string, v interface{}) error {
	cs, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return c.cao.SetByTTL(c.prefix+key, string(cs), c.expired)
}

// FindJSON 获取cache缓存
func (c *Cache) FindJSON(key string, v interface{}) error {
	id := c.prefix + key
	if id == "" {
		return errors.New("Cache Find id is required")
	}
	if v == nil {
		return errors.New("FIFO.Find v is nil")
	}
	cs, err := c.cao.Get(id)
	if err != nil {
		if err.Error() == "redis: nil" {
			return nil
		}
		return err
	}
	err = json.Unmarshal([]byte(cs), v)
	return err
}

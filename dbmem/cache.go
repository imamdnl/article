package dbmem

import (
	"time"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

type cacheCapsule struct {
	cache     *redis.Client
	cacheclst *redis.ClusterClient
	isclst    bool
	logger    *zap.Logger
}

type Cacher interface {
	Set(key string, pair string, value interface{}, duration time.Duration)
	HMSet(key string, pair string, value map[string]interface{})
	Delete(key string, pair string)
	Get(key string, pair string) string
	Ttl(key string, pair string) time.Duration
}

func NewCache(cache *redis.Client, cacheclst *redis.ClusterClient, isclst bool, logger *zap.Logger) Cacher {
	return &cacheCapsule{
		cache:     cache,
		cacheclst: cacheclst,
		isclst:    isclst,
		logger:    logger,
	}
}

func clstTtl(c *cacheCapsule, k string, p string) (out time.Duration) {
	cmd := c.cacheclst.TTL(k + ":" + p)
	if cmd.Err() != nil {
		return out
	} else {
		return cmd.Val()
	}
}

func nonClstTtl(c *cacheCapsule, k string, p string) (out time.Duration) {
	cmd := c.cache.TTL(k + ":" + p)
	if cmd.Err() != nil {
		c.logger.Error("error on redis TTL", zap.Error(cmd.Err()))
		return out
	} else {
		return cmd.Val()
	}
}

func (c *cacheCapsule) Ttl(k string, p string) time.Duration {
	if c.isclst {
		return clstTtl(c, k, p)
	} else {
		return nonClstTtl(c, k, p)
	}
}

func clstSet(c *cacheCapsule, k string, p string, v interface{}, d time.Duration) {
	c.cacheclst.Del(k + ":" + p)
	if d != 0*time.Second {
		_, err := c.cacheclst.SetNX(k+":"+p, v, d).Result()
		if err != nil {
			c.logger.Error("check error on clst setnx", zap.Error(err))
		}
	} else {
		_, err := c.cacheclst.Set(k+":"+p, v, 0).Result()
		if err != nil {
			c.logger.Error("check error on clst set", zap.Error(err))
		}
	}
}

func clstHMSet(c *cacheCapsule, k string, p string, v map[string]interface{}) {
	_, err := c.cacheclst.HMSet(k+":"+p, v).Result()
	if err != nil {
		c.logger.Error("check error on clst hmset", zap.Error(err))
	}
}

func nonClstSet(c *cacheCapsule, k string, p string, v interface{}, d time.Duration) {
	c.cache.Del(k + ":" + p)
	if d != 0*time.Second {
		c.cache.SetNX(k+":"+p, v, d)
	} else {
		c.cache.Set(k+":"+p, v, 0)
	}
}

func nonClstHMSet(c *cacheCapsule, k string, p string, v map[string]interface{}) {
	_, err := c.cache.HMSet(k+":"+p, v).Result()
	if err != nil {
		c.logger.Error("check error on non clst hmset", zap.Error(err))
	}
}

func (c *cacheCapsule) Set(k string, p string, v interface{}, d time.Duration) {
	if c.isclst {
		clstSet(c, k, p, v, d)
	} else {
		nonClstSet(c, k, p, v, d)
	}
}

func (c *cacheCapsule) HMSet(k string, p string, v map[string]interface{}) {
	if c.isclst {
		clstHMSet(c, k, p, v)
	} else {
		nonClstHMSet(c, k, p, v)
	}
}

func clstDel(c *cacheCapsule, k string, p string) {
	_, err := c.cacheclst.Del(k + ":" + p).Result()
	if err != nil {
		c.logger.Error("check error clst del", zap.Error(err))
	}
}

func nonClstDel(c *cacheCapsule, k string, p string) {
	c.cache.Del(k + ":" + p)
}

func (c *cacheCapsule) Delete(k string, p string) {
	if c.isclst {
		clstDel(c, k, p)
	} else {
		nonClstDel(c, k, p)
	}
}

func clstGet(c *cacheCapsule, k string, p string) (string, error) {
	return c.cacheclst.Get(k + ":" + p).Result()
}

func nonClstGet(c *cacheCapsule, k string, p string) (string, error) {
	return c.cache.Get(k + ":" + p).Result()
}

func (c *cacheCapsule) Get(k string, p string) string {
	var val string
	var err error
	if c.isclst {
		val, err = clstGet(c, k, p)
	} else {
		val, err = nonClstGet(c, k, p)
	}
	if err != nil {
		c.logger.Error("check error on get", zap.Error(err))
		return ""
	}
	return val
}

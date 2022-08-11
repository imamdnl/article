package configuration

import (
	"github.com/go-redis/redis"
)

func ConfigCache(addr string, pswd string) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     pswd,
		DB:           0,
		PoolSize:     50,
		MinIdleConns: 10,
	})
	return rdb
}

func ConfigClusterCache(adrs []string, pswd string) *redis.ClusterClient {
	rdbc := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:        adrs,
		Password:     pswd,
		PoolSize:     10,
		MinIdleConns: 10,
	})
	return rdbc
}

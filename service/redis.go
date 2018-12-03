package service

import (
	"github.com/go-redis/redis"
	"log"
	"time"
)

type RedisConfig struct{
	Addrs []string
	PoolSize int
	MinIdleConns int
	IdleTimeout time.Duration
	DialTimeout time.Duration
}

var (
	redisCluster *redis.ClusterClient
)

func (r* RedisConfig) Connect() error {

	rediClusterClient := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: r.Addrs,
		DialTimeout: r.DialTimeout,
		PoolSize: r.PoolSize,
		MinIdleConns:r.MinIdleConns,
		IdleTimeout: r.IdleTimeout,
	})

	err := rediClusterClient.Ping().Err()
	if err != nil {
		return err
	} else {
		log.Printf("connect redis cluster success")
	}

	redisCluster = rediClusterClient

	return nil
}

func Pub (channel string, message interface{}) (int64, error){
	cmd := redisCluster.Publish(channel, message)
	return cmd.Result()
}

func Sub (chanel string) {
	pubSub := redisCluster.Subscribe(chanel)
	for {
		message, err := pubSub.ReceiveMessage()
		if err!= nil {
			log.Printf("sub redis message err -> %s", err)
		} else {
			log.Printf("sub redis message -> %s", message.Payload)
		}
	}

}

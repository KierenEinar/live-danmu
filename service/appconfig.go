package service

import (
	"github.com/go-yaml/yaml"
	"io/ioutil"
	"log"
	"path"
	"strings"
	"time"
)

type AppConfig struct {
	RedisClusterNodes string `yaml:"redisClusterNodes"`
	RedisPoolSize int `yaml:"redisPoolSize"`
	RedisMinIdleConns int `yaml:"redisMinIdleConns"`
	RedisIdleTimeout string `yaml:"redisIdleTimeout"`
	RedisDialTimeout string `yaml:"redisDialTimeout"`
	ServerAddr string `yaml:"serverAddr"`
}

func LoadConfig () AppConfig{
	fileName := path.Join( "conf", "app.yaml")
	log.Printf("config file name -> %s", fileName)
	c, e := ioutil.ReadFile(fileName)
	if e != nil {
		log.Printf("locd config error-> %s", e)
	}
	log.Printf("load config -> %s", c)
	var appConfig AppConfig
	yaml.Unmarshal(c[:], &appConfig)
	return appConfig
}

func (this *AppConfig) GetRedisConfig () *RedisConfig {
	var addrs []string = strings.Split(this.RedisClusterNodes, ",")
	var poolSize int = this.RedisPoolSize
	var minIdleConns int = this.RedisMinIdleConns
	idleTimeout,_  := time.ParseDuration(this.RedisIdleTimeout)
	dialTimeout,_ := time.ParseDuration(this.RedisDialTimeout)
	return &RedisConfig{
		addrs,
		poolSize,
		minIdleConns,
		idleTimeout,
		dialTimeout,
	}
}

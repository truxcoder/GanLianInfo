package dao

import (
	"context"
	"os"
	"path/filepath"

	"github.com/go-redis/redis/v8"

	"github.com/ilyakaznacheev/cleanenv"
	dm8 "github.com/truxcoder/gorm-dm8"
	log "github.com/truxcoder/truxlog"
	"gorm.io/gorm"
)

type Config struct {
	Host          string `yaml:"host"`
	Port          string `yaml:"port"`
	UserName      string `yaml:"userName"`
	Password      string `yaml:"password"`
	RedisHost     string `yaml:"redisHost"`
	RedisPort     string `yaml:"redisPort"`
	RedisPassword string `yaml:"redisPassword"`
}

var cfg Config

func init() {
	var path string
	dir, _ := os.Getwd()
	file := filepath.Join(dir, "config.yaml")
	if pathExists(file) {
		path = file
	} else {
		//path = "D:\\server\\ganlian\\config\\config.yaml"
		path = filepath.Join(dir, "config/config.yaml")
	}
	err := cleanenv.ReadConfig(path, &cfg)
	if err != nil {
		log.Errorf("读取配置文件错误: %v\n", err)
	}
}

func Connect() *gorm.DB {
	dsn := "dm://" + cfg.UserName + ":" + cfg.Password + "@" + cfg.Host + ":" + cfg.Port + "?autoCommit=false"
	db, _ := gorm.Open(dm8.Open(dsn), &gorm.Config{})
	return db
}

func RedisConnect() (*redis.Client, error) {
	var ctx = context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisHost + ":" + cfg.RedisPort,
		Password: cfg.RedisPassword, // no password set
		DB:       0,                 // use default DB
	})
	if _, err := rdb.Ping(ctx).Result(); err != nil {
		log.Error("Redis服务器连接异常")
		return nil, err
	}
	return rdb, nil
}

func pathExists(path string) bool {
	//log.Infof("path:%s", path)
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

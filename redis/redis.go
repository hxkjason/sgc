package redis

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"sgc/utils"
	"time"
)

var Rdb *redis.Client

func Init(addr, password string, enableLTS bool) {
	options := redis.Options{
		Addr:     addr,
		Password: password,
	}

	if enableLTS {
		options.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}
	Rdb = redis.NewClient(&options)

	// 测试连接
	if _, err := Rdb.Set(ctx, "testTmp-"+time.Now().Format("2006-01-02T15:04:05.000Z"), 1, 2*time.Second).Result(); err != nil {
		panic("redis conn err:" + err.Error())
	} else {
		fmt.Println("redis conn success!")
	}

}

func CloseRedis() {
	Rdb.Close()
}

// SetKeyValue 设置缓存
func SetKeyValue(keyName string, value interface{}, duration time.Duration) {

	err := Rdb.Set(ctx, keyName, value, duration).Err()

	if err != nil {
		fmt.Println(utils.SplicingStr("设置 redis key:", keyName, "失败:", err.Error()))
		panic(errors.New("设置缓存服务繁忙，请稍后重试"))
	}
}

// GetKeyValue 获取缓存
func GetKeyValue(keyName string) (string, error) {

	value, err := Rdb.Get(ctx, keyName).Result()

	if err != nil {
		fmt.Println(utils.SplicingStr("获取 redis key:", keyName, "失败:", err.Error()))
	}
	return value, err
}

func Exists(keyName string) (exists int64, err error) {
	return Rdb.Exists(ctx, keyName).Result()
}

// SetKeyNX 当 keyName 不存在时设置缓存
func SetKeyNX(keyName string, value interface{}, duration time.Duration) (hasSet bool, err error) {
	return Rdb.SetNX(ctx, keyName, value, duration).Result()
}

func DelKey(keyName string) (delKeyCount int64, err error) {
	return Rdb.Del(ctx, keyName).Result()
}

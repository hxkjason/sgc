package redis

import (
	"context"
	"errors"
	uuid "github.com/satori/go.uuid"
	"time"
)

type Locker struct {
	Key       string
	RequestId string // 客户端唯一id 用来指定锁不被其他线程(协程)删除
}

var ctx = context.Background()

// AddLock 加锁
func AddLock(key string, expiration time.Duration, tryTimes uint) (Locker, error) {

	locker := Locker{
		Key:       key,
		RequestId: uuid.NewV4().String(),
	}

	triedTimes := uint(0)
	for {
		if triedTimes >= tryTimes {
			break
		}
		triedTimes++
		value, err := Rdb.SetNX(ctx, locker.Key, locker.RequestId, expiration).Result()

		if err == nil {
			if value == true {
				return locker, nil
			} else {
				return locker, errors.New("请求频繁")
			}
		}
	}
	return Locker{}, errors.New("TryLock fail")
}

// Unlock 解锁
func (locker *Locker) Unlock() error {

	//查key的值
	lockerRequestId, err := Rdb.Get(ctx, locker.Key).Result()
	if err != nil {
		return err
	}

	//校验客户端请求ID
	if lockerRequestId != locker.RequestId {
		return errors.New("锁已失效")
	}

	//解锁
	delVal, err := Rdb.Del(context.Background(), locker.Key).Result()
	if err != nil {
		return err
	}
	if delVal != 1 {
		return errors.New("解锁失败")
	}
	return nil
}

func LockKey(key string, expiredDuration, sleepDuration time.Duration, maxTimes int) (Locker, error) {

	if maxTimes == 0 || maxTimes > 10 {
		maxTimes = 1
	}

	locker := Locker{
		Key:       key,
		RequestId: uuid.NewV4().String(),
	}

	tryTimes := 0

	for tryTimes < maxTimes {
		tryTimes++
		currentTimeSet, err := SetKeyNX(locker.Key, locker.RequestId, expiredDuration)
		if err != nil {
			continue
		}
		if currentTimeSet {
			return locker, nil
		}
		time.Sleep(sleepDuration)
	}

	return locker, errors.New("获取锁[" + key + "]失败")
}

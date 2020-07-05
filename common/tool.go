package common

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/go-redis/redis"
)

// 生成随机字符串
var src = rand.NewSource(time.Now().UnixNano())

const letterBytes = "0123456abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6
	letterIdxMask = 1<<letterIdxBits - 1
	letterIdxMax  = 63 / letterIdxBits
)

func RandStringBytesMaskImprSrc(n int) string {
	b := make([]byte, n)
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func DeleteRedisRecodeFromRecode(recode string) {
	tokenRedisClient := GetUserTokenRedisClient()
	reCodeRedisClient := GetReCodeRedisClient()
	if !KeyISExistInRedis(recode, reCodeRedisClient) {
		log.Println(fmt.Errorf("key %s is not in reCodeRedis", recode))
	}
	userToken := reCodeRedisClient.Get(recode).Val()
	if !KeyISExistInRedis(userToken, tokenRedisClient) {
		log.Println(fmt.Errorf("recode %s is in reCodeRedis,But key %s is not in userTokenRedis", recode, userToken))
	}
	reCodeRedisClient.Del(recode)
	tokenRedisClient.Del(userToken)
}

func KeyISExistInRedis(str string, client *redis.Client) bool {
	if client.Exists(str).Val() == 0 {
		return false
	}
	return true
}

func SetShadowKeyInRedis(key string, value interface{}, expiration time.Duration, client *redis.Client) error {
	err := client.Set(key, value, 0).Err()
	if err != nil {
		log.Println(fmt.Sprintf("WriteShadowKeyInRedis key %s value %v fail in realKey Set", key, value))
		return err
	}
	err = client.Set(SHADOWKEYPREFIX+key, "", expiration).Err()
	if err != nil {
		log.Println(fmt.Sprintf("WriteShadowKeyInRedis key %s value %v fail in shadowKey Set", key, value))
		return err
	}
	return nil
}

func ReplaceShadowKeyInRedis(key string, expiration time.Duration, client *redis.Client) error {
	err = client.Expire(SHADOWKEYPREFIX+key, expiration).Err()
	if err != nil {
		log.Println(fmt.Sprintf("ReplaceShadowKeyInRedis key %s fail in shadowKey Set", key))
		return err
	}
	return nil
}

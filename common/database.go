package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Db struct {
	User         string
	Pass         string
	Host         string
	Port         string
	Database     string
	MaxIdleConns int
	MaxOpenConns int
	Debug        bool
}

type RedisDb struct {
	Host string
	Port string
	Pass string
}

type Dbs struct {
	Master *Db
	Redis  *RedisDb
	Slave  *Db
}

type Database struct {
	*gorm.DB
}

var (
	DB                   *gorm.DB
	DbConfig             *Dbs
	UserTokenRedisClient *redis.Client
	ReCodeRedisClient    *redis.Client
)

func getDBConfigFromFile() (*Dbs, error) {
	var config Dbs
	if conf, err := ioutil.ReadFile(DBConfigFile); err == nil {
		e := json.Unmarshal(conf, &config)
		return &config, e
	} else {
		return &config, err
	}
}

func InitDB() *gorm.DB {
	DBConf, err := getDBConfigFromFile()
	if err != nil {
		fmt.Println("Get DBConfig From File Err:", err)
	}
	DbConfig = DBConf
	var (
		db *gorm.DB
		e  error
	)
	for {
		// 重试连接
		db, e = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&collation=utf8mb4_bin&parseTime=True&loc=%s", DBConf.Master.User, DBConf.Master.Pass, DBConf.Master.Host, DBConf.Master.Port, DBConf.Master.Database, MYSQLTIMEZONE))
		if e != nil {
			fmt.Println("Gorm Open DB Err: ", err)
			time.Sleep(20 * time.Second)
		} else {
			break
		}
	}
	log.Println("Connected to database ", DBConf.Master.User, " ", DBConf.Master.Pass, " ", DBConf.Master.Host, ":", DBConf.Master.Port, " ", DBConf.Master.Database)
	db.DB().SetMaxIdleConns(DBConf.Master.MaxIdleConns)
	DB = db
	return DB
}

func GetDB() *gorm.DB {
	return DB
}

func GetUserTokenRedisClient() *redis.Client {
	return UserTokenRedisClient
}

func GetReCodeRedisClient() *redis.Client {
	return ReCodeRedisClient
}

func InitRedis(redisDBCode int) *redis.Client {
	addr := fmt.Sprintf("%s:%s", DbConfig.Redis.Host, DbConfig.Redis.Port)
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: DbConfig.Redis.Pass,
		DB:       redisDBCode,
	})
	_, err := client.Ping().Result()
	if err != nil {
		log.Println("Ping Redis DB:", redisDBCode, "Get err", err)
	}
	log.Println("Connected to Redis:", addr, "DB Number:", redisDBCode)

	return client
}

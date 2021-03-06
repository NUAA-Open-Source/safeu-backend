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
	User            string
	Pass            string
	Host            string
	Port            string
	Database        string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime int
	Debug           bool
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
	connectString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&collation=utf8mb4_bin&parseTime=True&loc=%s&timeout=10s", DBConf.Master.User, DBConf.Master.Pass, DBConf.Master.Host, DBConf.Master.Port, DBConf.Master.Database, MYSQLTIMEZONE)
	// 重试连接
	for db, e = gorm.Open("mysql", connectString); e != nil; {
		fmt.Println("Gorm Open DB Err: ", e)
		log.Println(fmt.Sprintf("GORM cannot connect to database, retry in %d seconds...", DB_CONNECT_FAIL_RETRY_INTERVAL))
		time.Sleep(DB_CONNECT_FAIL_RETRY_INTERVAL * time.Second)
	}

	log.Println("Connected to database ", DBConf.Master.User, " ", DBConf.Master.Pass, " ", DBConf.Master.Host, ":", DBConf.Master.Port, " ", DBConf.Master.Database)
	db.DB().SetMaxIdleConns(DBConf.Master.MaxIdleConns)
	db.DB().SetMaxOpenConns(DBConf.Master.MaxOpenConns)
	db.DB().SetConnMaxLifetime(time.Duration(DBConf.Master.ConnMaxLifetime) * time.Second)
	DB = db
	DB.LogMode(true)
	return DB
}

func GetDB() *gorm.DB {
	// Ping
	err := DB.DB().Ping()
	if err != nil {
		log.Println("Cannot access the database (PING FAILED)")
	}
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

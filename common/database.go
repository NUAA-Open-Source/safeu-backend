package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

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

type Dbs struct {
	Master *Db
	Slave  *Db
}

type Database struct {
	*gorm.DB
}

var (
	DB       *gorm.DB
	DbConfig *Dbs
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
	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@/%s?charset=utf8mb4&parseTime=True&loc=Local", DBConf.Master.User, DBConf.Master.Pass, DBConf.Master.Database))
	if err != nil {
		fmt.Println("Gorm Open DB Err: ", err)
	}
	db.DB().SetMaxIdleConns(DBConf.Master.MaxIdleConns)
	DB = db
	return DB
}

func GetDB() *gorm.DB {
	return DB
}

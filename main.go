//  Copyright 2019 A2OS SafeU Dev Team
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package main

import (
	"log"

	"a2os/safeu-backend/common"
	"a2os/safeu-backend/item"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func Migrate(db *gorm.DB) {
	db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4 auto_increment=1").AutoMigrate(&item.Item{})
	db.AutoMigrate(&common.Config{})
	db.AutoMigrate(&item.Token{})
}
func init() {
	//DB init
	db := common.InitDB()
	Migrate(db)
	//Redis init
	//初始化 UserToken Redis连接
	common.UserTokenRedisClient = common.InitRedis(common.USER_TOKEN)
	//defer db.Close()
	//Read Config
	conf, err := common.GetCloudConfig()
	if err != nil {
		log.Println("GetCloudConfig Err", err)
	}
	common.CloudConfig = conf
}

func main() {

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	api := r.Group("/api")
	v1 := api.Group("/v1")
	{
		item.UploadRegister(v1.Group("/upload"))
		v1.POST("/password/:retrieveCode", item.ChangePassword)
		v1.POST("/recode/:retrieveCode", item.ChangeRecode)
		v1.POST("/downCount/:retrieveCode", item.ChangeDownCount)
		v1.POST("/item/:retrieveCode", item.DownloadItems)
		v1.POST("/validation/:retrieveCode", item.Validation)
	}

	r.Run() // listen and serve on 0.0.0.0:8080
}

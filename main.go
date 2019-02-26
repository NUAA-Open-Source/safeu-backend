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
	"fmt"
	"log"
	"time"

	"a2os/safeu-backend/common"
	"a2os/safeu-backend/item"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"io"
)

func Migrate(db *gorm.DB) {
	db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4 COLLATE=utf8mb4_bin auto_increment=1").AutoMigrate(&item.Item{})
	db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4 COLLATE=utf8mb4_bin auto_increment=1").AutoMigrate(&common.Config{})
	db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4 COLLATE=utf8mb4_bin auto_increment=1").AutoMigrate(&item.Token{})
}
func init() {

	//_ = os.Mkdir("log", os.ModePerm)
	//logFile, err := os.OpenFile("log/safeu-backend.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//log.SetOutput(logFile)

	// Logger init
	common.InitLogger()

	//DB init
	db := common.InitDB()
	Migrate(db)
	//Redis init
	//初始化 UserToken Redis连接
	common.UserTokenRedisClient = common.InitRedis(common.USER_TOKEN)
	common.ReCodeRedisClient = common.InitRedis(common.RECODE)
	//defer db.Close()
	//Read Config
	conf, err := common.GetCloudConfig()
	if err != nil {
		log.Println("GetCloudConfig Err", err)
	}
	common.CloudConfig = conf
	log.Println(fmt.Sprintf("Read Aliyun Config :%v", conf.AliyunConfig))
	log.Println(fmt.Sprintf("Read Server Config :%v", conf.Server))
	log.Println(fmt.Sprintf("Read FaaS Config: %v", conf.FaaS))

	// 初始化阿里云对象存储客户端对象
	common.InitAliyunOSSClient()

}

func main() {

	// Before init router
	if common.DEBUG {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
		// Redirect log to file
		gin.DisableConsoleColor()
		logFile := common.GetLogFile()
		defer logFile.Close()
		gin.DefaultWriter = io.MultiWriter(logFile)
	}

	r := gin.Default()

	// After init router
	if common.DEBUG {
		r.Use(cors.New(cors.Config{
			AllowAllOrigins:  true,
			AllowMethods:     common.CORS_ALLOW_METHODS,
			AllowHeaders:     common.CORS_ALLOW_HEADERS,
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}))
		//r.Use(CORS())
	} else {
		// RELEASE Mode
		r.Use(cors.New(cors.Config{
			AllowOrigins:     common.CORS_ALLOW_ORIGINS,
			AllowMethods:     common.CORS_ALLOW_METHODS,
			AllowHeaders:     common.CORS_ALLOW_HEADERS,
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}))
	}

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	v1 := r.Group("/v1")
	{
		item.UploadRegister(v1.Group("/upload"))
		v1.POST("/password/:retrieveCode", item.ChangePassword)
		v1.POST("/recode/:retrieveCode", item.ChangeRecode)
		v1.POST("/delete/:retrieveCode", item.DeleteManual)
		v1.GET("/downCount/:retrieveCode", item.DownloadCount)
		v1.POST("/downCount/:retrieveCode", item.ChangeDownCount)
		v1.POST("/expireTime/:retrieveCode", item.ChangeExpireTime)
		v1.POST("/item/:retrieveCode", item.DownloadItems)
		v1.POST("/validation/:retrieveCode", item.Validation)
	}

	r.Run(":" + common.PORT) // listen and serve on 0.0.0.0:PORT
}

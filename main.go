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
	"io"
	"log"
	"net/http"
	"time"

	"a2os/safeu-backend/common"
	"a2os/safeu-backend/item"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/utrack/gin-csrf"
)

func Migrate(db *gorm.DB) {
	db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4 COLLATE=utf8mb4_bin auto_increment=1").AutoMigrate(&item.Item{})
	db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4 COLLATE=utf8mb4_bin auto_increment=1").AutoMigrate(&common.Config{})
	db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4 COLLATE=utf8mb4_bin auto_increment=1").AutoMigrate(&item.Token{})
}

// 系统启动后的任务
func Tasks() {
	// 主动删除
	go item.ActiveDelete(common.GetReCodeRedisClient())
}
func init() {

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
	// 系统启动后的任务
	Tasks()

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
	// 错误处理
	r.Use(common.ErrorHandling())
	// After init router
	// CORS
	if common.DEBUG {
		r.Use(cors.New(cors.Config{
			// The value of the 'Access-Control-Allow-Origin' header in the
			// response must not be the wildcard '*' when the request's
			// credentials mode is 'include'.
			AllowOrigins:     common.CORS_ALLOW_DEBUG_ORIGINS,
			AllowMethods:     common.CORS_ALLOW_METHODS,
			AllowHeaders:     common.CORS_ALLOW_HEADERS,
			ExposeHeaders:    common.CORS_EXPOSE_HEADERS,
			AllowCredentials: true,
			AllowWildcard:    true,
			MaxAge:           12 * time.Hour,
		}))
		//r.Use(CORS())
	} else {
		// RELEASE Mode
		r.Use(cors.New(cors.Config{
			AllowOrigins:     common.CORS_ALLOW_ORIGINS,
			AllowMethods:     common.CORS_ALLOW_METHODS,
			AllowHeaders:     common.CORS_ALLOW_HEADERS,
			ExposeHeaders:    common.CORS_EXPOSE_HEADERS,
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}))
	}

	// CSRF
	store := cookie.NewStore(common.CSRF_COOKIE_SECRET)
	r.Use(sessions.Sessions(common.CSRF_SESSION_NAME, store))
	CSRF := csrf.Middleware(csrf.Options{
		Secret: common.CSRF_SECRET,
		ErrorFunc: func(c *gin.Context) {
			//c.String(http.StatusBadRequest, "CSRF token mismatch")
			c.JSON(http.StatusBadRequest, gin.H{
				"err_code": 10007,
				"message":  common.Errors[10007],
			})
			log.Println(c.ClientIP(), "CSRF token mismatch")
			c.Abort()
		},
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/csrf", CSRF, func(c *gin.Context) {
		c.Header("X-CSRF-TOKEN", csrf.GetToken(c))
		c.String(http.StatusOK, "IN HEADER")
		log.Println(c.ClientIP(), "response CSRF token", csrf.GetToken(c))
	})

	// the API without CSRF middleware
	v1 := r.Group("/v1")
	{
		v1.POST("/upload/callback", item.UploadCallBack) //回调
	}

	// the API with CSRF middleware
	v1_csrf := r.Group("/v1", CSRF)
	{
		v1_csrf.GET("/upload/policy", item.GetPolicyToken) //鉴权
		v1_csrf.POST("/upload/finish", item.FinishUpload)  //结束
		v1_csrf.POST("/password/:retrieveCode", item.ChangePassword)
		v1_csrf.POST("/recode/:retrieveCode", item.ChangeRecode)
		v1_csrf.POST("/delete/:retrieveCode", item.DeleteManual)
		v1_csrf.POST("/info/:retrieveCode", item.GetItemInfo)
		v1_csrf.POST("/minusDownCount/:retrieveCode", item.MinusDownloadCount)
		v1_csrf.POST("/downCount/:retrieveCode", item.ChangeDownCount)
		v1_csrf.POST("/expireTime/:retrieveCode", item.ChangeExpireTime)
		v1_csrf.POST("/item/:retrieveCode", item.DownloadItems)
		v1_csrf.POST("/validation/:retrieveCode", item.Validation)
	}

	r.Run(":" + common.PORT) // listen and serve on 0.0.0.0:PORT
}

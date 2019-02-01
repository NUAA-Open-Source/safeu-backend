package item

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"a2os/safeu-backend/common"

	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
)

func UploadRegister(router *gin.RouterGroup) {
	router.GET("/policy", GetPolicyToken)    //鉴权
	router.POST("/callback", UploadCallBack) //回调
}

func GetPolicyToken(c *gin.Context) {
	//TODO:错误处理
	response := get_policy_token()
	c.String(http.StatusOK, response)
}

type FileInfo struct {
	Bucket   string `json:"bucket"`
	Object   string `json:"object"`
	Etag     string `json:"etag"`
	Size     int    `json:"size"`
	MimeType string `json:"mimeType"`
	Height   int    `json:"height"`
	Width    int    `json:"width"`
	Format   string `json:"format"`
}

func UploadCallBack(c *gin.Context) {
	r := c.Request
	// 双拷贝http Request流 避免读写偏移
	buf, _ := ioutil.ReadAll(r.Body)
	rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
	rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf))
	bufTemp := new(bytes.Buffer)
	bufTemp.ReadFrom(rdr1)
	s := fmt.Sprintf("{%s}", bufTemp.String())
	var fileInfo FileInfo
	err := json.Unmarshal([]byte(s), &fileInfo)
	if err != nil {
		fmt.Println("Json Unmarshal Fail", err)
	}
	bytePublicKey, err := getPublicKey(r)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err,
		})
		return
	}

	byteAuthorization, err := getAuthorization(r)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err,
		})
		return
	}

	r.Body = rdr2
	byteMD5, err := getMD5FromNewAuthString(r)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err,
		})
		return
	}

	if verifySignature(bytePublicKey, byteMD5, byteAuthorization) {
		// TODO:完善此处
		host := fmt.Sprintf("https://%s.%s/%s", common.CloudConfig.Aliyun[0].EndPoint[0].Bucket[0].Name, common.CloudConfig.Aliyun[0].EndPoint[0].Base, fileInfo.Object)
		u := uuid.Must(uuid.NewV4())
		item := Item{Name: u.String(), Host: host, Status: 0, OriginalName: fileInfo.Object}
		db := common.GetDB()
		db.NewRecord(item)
		db.Create(&item)
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
		return
	} else {
		log.Println("Fail")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err,
		})
		return
	}
}

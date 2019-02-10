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

type FinishedFiles struct {
	Files []uuid.UUID `json:"files"`
}

func UploadRegister(router *gin.RouterGroup) {
	router.GET("/policy", GetPolicyToken)    //鉴权
	router.POST("/callback", UploadCallBack) //回调
	router.POST("/finish", FinishUpload)     //结束
}

func GetPolicyToken(c *gin.Context) {
	//TODO:错误处理
	response := get_policy_token()
	c.String(http.StatusOK, response)
}

func FinishUpload(c *gin.Context) {
	//TODO：本函数待优化
	var finishedFiles FinishedFiles
	db := common.GetDB()
	err := c.BindJSON(&finishedFiles)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
		return
	}
	if finishedFiles.Files == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Parameter error",
		})
		return
	}

	// 数据库中存在满足uuid且状态为"上传完成"，返回原来生成的提取码
	var item Item
	if !db.Where("name = ? AND status = ?", finishedFiles.Files[0], common.UPLOAD_FINISHED).First(&item).RecordNotFound() {
		c.JSON(http.StatusOK, gin.H{
			"recode": item.ReCode,
		})
		return
	}
	// 存在满足条件uuid且状态为"上传阶段"，生成新的提取码
	reCode := common.RandStringBytesMaskImprSrc(common.ReCodeLength)
	var files []string
	for _, value := range finishedFiles.Files {
		fmt.Println(value)
		files = append(files, value.String())
		db.Model(&Item{}).Where("name = ? AND status = ?", value, common.UPLOAD_BEGIN).Update(map[string]interface{}{"re_code": reCode, "status": common.UPLOAD_FINISHED})
	}
	// 将用户识别码推入Redis
	tokenRedisClient := common.GetUserTokenRedisClient()
	owner := common.RandStringBytesMaskImprSrc(common.UserTokenLength)
	tokenRedisClient.SAdd(owner, files)
	c.JSON(http.StatusOK, gin.H{
		"recode": reCode,
		"owner":  owner,
	})
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
			"uuid": u,
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

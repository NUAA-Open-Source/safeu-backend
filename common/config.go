package common

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type CloudConfig struct {
	Aliyun []AliyunAccount
}
type AliyunAccount struct {
	AccountId       string
	AccessKey       string
	AccessKeySecret string
	EndPoint        []EndPoint
}
type EndPoint struct {
	EndPointId string
	Base       string
	Bucket     []Bucket
}
type Bucket struct {
	BucketId string
	Name     string
}

type Config struct {
	Name    string
	Content string `sql:"type:text"`
}

func (c *CloudConfig) GetCloudConfigFromFile() *CloudConfig {
	yamlFile, err := ioutil.ReadFile(CloudConfigFile)
	if err != nil {
		log.Println("GetCloudConfigFromFile YamlFile Read Fail", err)
	}
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		log.Println("Yaml Unmarshal Fail", err)
	}
	return c
}

func (c *CloudConfig) GetCloudConfigFromDB() *CloudConfig {
	db := GetDB()
	var conf Config
	db.Where("name = ?", CloudConfigDBName).First(&conf)
	err := yaml.Unmarshal([]byte(conf.Content), &c)
	if err != nil {
		log.Println("GetCloudConfigFromDB Yaml Unmarshal Fail", err)
	}
	return c
}

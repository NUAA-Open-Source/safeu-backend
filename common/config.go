package common

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type CloudConfiguration struct {
	AliyunConfig AliyunConfig
	Server       []ServerConfig
	FaaS         []FaaSConfig
}
type AliyunConfig struct {
	Accounts []AliyunAccount
	Retry    int
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

type ServerConfig struct {
	ServerId       string
	ServerCallBack string
}

type FaaSConfig struct {
	Name     string
	Endpoint string
}

func GetCloudConfig() (c *CloudConfiguration, err error) {
	c, err = GetCloudConfigFromDB()
	if err != nil || c == nil {
		log.Println("GetCloudConfigFromDB Fail", err)
		c, err = GetCloudConfigFromFile()
		if err != nil || c == nil {
			log.Println("GetCloudConfigFromFile Fail", err)
			return c, err
		}
		log.Println("Get Cloud Config From File Success!")
		return c, nil
	}
	log.Println("Get Cloud Config From DB Success!")
	return c, nil
}
func GetCloudConfigFromFile() (*CloudConfiguration, error) {
	var cloudConfig CloudConfiguration
	yamlFile, err := ioutil.ReadFile(CloudConfigFile)
	if err != nil {
		log.Println("GetCloudConfigFromFile YamlFile Read Fail", err)
		return &cloudConfig, err
	}
	err = yaml.Unmarshal(yamlFile, &cloudConfig)
	if err != nil {
		log.Println("Yaml Unmarshal Fail", err)
		return &cloudConfig, err
	}
	return &cloudConfig, nil
}

func GetCloudConfigFromDB() (*CloudConfiguration, error) {
	var cloudConfig CloudConfiguration

	db := GetDB()
	var conf Config
	db.Where("name = ?", CloudConfigDBName).First(&conf)
	err := yaml.Unmarshal([]byte(conf.Content), &cloudConfig)
	if err != nil {
		log.Println("GetCloudConfigFromDB Yaml Unmarshal Fail", err)
		return &cloudConfig, err
	}
	return &cloudConfig, nil
}

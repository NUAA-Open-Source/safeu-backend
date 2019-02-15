package common

import (
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type CloudConfiguration struct {
	AliyunConfig AliyunConfig   `yaml:"aliyun"`
	Server       []ServerConfig `yaml:"server"`
	FaaS         []FaaSConfig   `yaml:"faas"`
}
type AliyunConfig struct {
	Accounts []AliyunAccount `yaml:"account"`
	Retry    int             `yaml:"retry"`
}
type AliyunAccount struct {
	AccountId       string     `yaml:"accountid"`
	AccessKey       string     `yaml:"accesskey"`
	AccessKeySecret string     `yaml:"accesskeysecret"`
	EndPoint        []EndPoint `yaml:"endpoint"`
}
type EndPoint struct {
	EndPointId string   `yaml:"endpointid"`
	Base       string   `yaml:"base"`
	Bucket     []Bucket `yaml:"bucket"`
}
type Bucket struct {
	BucketId string `yaml:"bucketid"`
	Name     string `yaml:"name"`
}

type Config struct {
	Name    string
	Content string `sql:"type:text"`
}

type ServerConfig struct {
	ServerId       string `yaml:"serverid"`
	ServerCallBack string `yaml:"servercallback"`
}

type FaaSConfig struct {
	Name     string `yaml:"name"`
	Endpoint string `yaml:"endpoint"`
}

func GetCloudConfig() (c *CloudConfiguration, err error) {
	c, err = GetCloudConfigFromDB()
	if err != nil || isAliyunConfigEmpty(c) != nil {
		log.Println("GetCloudConfigFromDB Fail", err)
		c, err = GetCloudConfigFromFile()
		if err != nil || isAliyunConfigEmpty(c) != nil {
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

func isAliyunConfigEmpty(c *CloudConfiguration) error {
	if c.AliyunConfig.Accounts == nil ||
		c.AliyunConfig.Retry == 0 {
		return fmt.Errorf("aliyun config is empty")
	}
	return nil
}

package common

type Config struct {
	Name    string
	Content string `sql:"TYPE:json"`
}

//TODO:读取数据库中阿里云配置文件
//type AliyunConfig struct {
//}
//
//func InitConfig(db *gorm.DB) {
//
//	db.Where("name = ?", "aliyun").First(AliyunConfig{})
//
//}

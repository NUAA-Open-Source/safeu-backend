package validation

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Token struct {
	gorm.Model
	Token			string	`gorm:"NOT NULL;INDEX"`
	RetrieveCode	string	`gorm:"NOT NULL"`
	Valid			bool
	ExpiredAt		time.Time	`gorm:"NOT NULL"`
}

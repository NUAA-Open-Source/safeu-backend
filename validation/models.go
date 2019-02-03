package validation

import (
	"github.com/jinzhu/gorm"
)

type Token struct {
	gorm.Model
	Token	string	`gorm:"NOT NULL;INDEX"`
	Valid	bool
}

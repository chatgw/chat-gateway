package schema

import "gorm.io/gorm"

type User struct {
	gorm.Model

	LoginName string
	Password  string
}

func (User) TableName() string {
	return "tbl_users"
}

package schema

import "gorm.io/gorm"

type User struct {
	gorm.Model

	LoginName string
	Password  string
}

// TableName overrides the table name used by User to `profiles`
func (User) TableName() string {
	return "tbl_users"
}

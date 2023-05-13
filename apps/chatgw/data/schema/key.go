package schema

import "gorm.io/gorm"

type Key struct {
	gorm.Model

	Title  string
	Vendor string // openai azure
	Token  string
}

func (Key) TableName() string {
	return "tbl_keys"
}

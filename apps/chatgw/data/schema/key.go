package schema

import "gorm.io/gorm"

type Key struct {
	gorm.Model

	Title  string
	Vendor string // openai azure
	Token  string
}

// TableName overrides the table name used by Key to `profiles`
func (Key) TableName() string {
	return "tbl_keys"
}

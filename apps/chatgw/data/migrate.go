package data

import (
	"github.com/airdb/chat-gateway/apps/chatgw/data/schema"
	"gorm.io/gorm"
)

func Migrate(conn *gorm.DB) {
	conn.Migrator().AutoMigrate(
		&schema.Key{},
		&schema.User{},
	)
}

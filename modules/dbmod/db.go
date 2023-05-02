package dbmod

import (
	"os"

	"github.com/xo/dburl"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewConn() (*gorm.DB, error) {
	dsn, err := dburl.Parse(os.Getenv("CHATGW_DSN"))
	if err != nil {
		return nil, err
	}

	return gorm.Open(mysql.Open(dsn.DSN), &gorm.Config{})
}

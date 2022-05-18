package dal

import (
	"github.com/AlaricGilbert/argos-core/master/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDatabase() {
	var err error
	db, err = gorm.Open(mysql.Open(config.DbConnDsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
}

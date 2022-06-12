package gorm_db

import (
	"dbtogo/defines"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
)

var db *gorm.DB
var mu sync.Mutex

func NewGormDB(dbConfig *defines.DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/information_schema?charset=%s&parseTime=True&loc=Local&allowAllFiles=true", dbConfig.Username, dbConfig.Password, dbConfig.IP, dbConfig.Port, dbConfig.Charset)
	var err error
	mu.Lock()
	defer mu.Unlock()
	if db != nil {
		return db, err
	}
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	return db, err
}

func ClearGormDB() {
	mu.Lock()
	defer mu.Unlock()
	if db != nil {
		db = nil
	}
}

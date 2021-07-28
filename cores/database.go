package cores

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"strings"
)

const (
	DriverMySQL = "MYSQL"
	DriverSQLite = "SQLITE"
)

func InitDataBase(conf *ConfigGlobal) error {

	// 初始化数据库引擎
	dataBaseDriver := strings.ToUpper(conf.Database.Driver)

	if (dataBaseDriver != DriverMySQL) && (dataBaseDriver != DriverSQLite) {
		return fmt.Errorf("does not support the current drive type `%s`", conf.Database.Driver)
	}

	var (
		db *gorm.DB
		err error
	)

	switch dataBaseDriver {
	case DriverSQLite:

	default:
		db, err = gorm.Open(strings.ToLower(dataBaseDriver),
			fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
				conf.Database.Username,
				conf.Database.Password,
				conf.Database.Host,
				conf.Database.Port,
				conf.Database.DbName))
	}

	if err != nil {
		return err
	}

	db.LogMode(conf.Database.SqlMode)
	db.DB().SetMaxIdleConns(conf.Database.MaxIdelConnections)
	db.DB().SetMaxOpenConns(conf.Database.MaxOpenConnections)

	conf.Runtime.DB = db

	return nil
}

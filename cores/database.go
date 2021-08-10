package cores

import (
	"apioak-admin/app/packages"
	"fmt"
	"gorm.io/driver/mysql"
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"strings"
)

const (
	DriverMySQL  = "MYSQL"
	DriverSQLite = "SQLITE"
)

func InitDataBase(conf *ConfigGlobal) error {

	dataBaseDriver := strings.ToUpper(conf.Database.Driver)
	if (dataBaseDriver != DriverMySQL) && (dataBaseDriver != DriverSQLite) {
		return fmt.Errorf("does not support the current drive type `%s`", conf.Database.Driver)
	}

	var (
		db  *gorm.DB
		err error
	)

	switch dataBaseDriver {
	case DriverSQLite:

	default:
		db, err = gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			conf.Database.Username,
			conf.Database.Password,
			conf.Database.Host,
			conf.Database.Port,
			conf.Database.DbName)), &gorm.Config{
			Logger:logger.Default.LogMode(logger.Info),
		})
	}

	if err != nil {
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("sqlDB init error: `%s`", err)
	}
	sqlDB.SetMaxIdleConns(conf.Database.MaxIdelConnections)
	sqlDB.SetMaxOpenConns(conf.Database.MaxOpenConnections)

	conf.Runtime.DB = db
	packages.SetDb(db)
	return nil
}

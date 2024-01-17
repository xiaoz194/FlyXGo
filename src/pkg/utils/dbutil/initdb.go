package dbutil

import (
	"github.com/xiaoz194/FlyXGo/src/pkg/utils/logutil"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"gorm.io/gorm/schema"
	"time"
)

var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB(dsn string, enableDBLogMode bool, dbMaxIdleConns int, dbMaxOpenConns int, values ...interface{}) error {
	// 打开数据库连接
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{ // 使用单数表名
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		logutil.LogrusObj.Infof("db connect error: %v", err)
		panic("db connect error")
	}
	// 设置日志模式
	if enableDBLogMode {
		db.Logger = db.Logger.LogMode(logger.Info)
	} else {
		db.Logger = db.Logger.LogMode(logger.Silent)
	}
	// 获取通用数据库对象 sql.DB ，然后使用其提供的功能
	sqlDB, err := db.DB()
	if err != nil {
		logutil.LogrusObj.Infof("get object DB error: %v", err)
		return err
	}
	// 设置连接池
	sqlDB.SetMaxIdleConns(dbMaxIdleConns)      //最大空闲连接数
	sqlDB.SetMaxOpenConns(dbMaxOpenConns)      // 最大打开连接数
	sqlDB.SetConnMaxLifetime(time.Second * 30) // 连接的最大存活时间
	if len(values) > 0 {
		migration(db, values...)
	}
	DB = db
	return nil
}

// migration 迁移函数，
// 参数可变参数 values ...interface{}
// 返回值 无
// 示例：调用函数时，可以传入任意数量的参数，它们将被打包为切片并传递给函数。例如：
// migration(&InternalNotifyMessage{}, &MessageRecord{}, &WechatSubscribe{})
func migration(db *gorm.DB, values ...interface{}) {
	//自动迁移模式
	logutil.LogrusObj.Info(values...)
	db.Set("gorm:table_options", "charset=utf8mb4").AutoMigrate(values...)
}

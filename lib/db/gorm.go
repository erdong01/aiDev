package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

// InitMySQL 初始化 MySQL 连接
// dsn 格式: user:password@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local
func InitMySQL(dsn string) {
	fmt.Println("dsn:", dsn)
	var err error

	// 自定义 logger
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond, // 慢 SQL 阈值
			LogLevel:                  logger.Info,            // 日志级别
			IgnoreRecordNotFoundError: true,                   // 忽略 RecordNotFound 错误
			Colorful:                  true,                   // 彩色输出
		},
	)

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: "pq_",
		},
	})
	if err != nil {
		panic(fmt.Sprintf("连接 MySQL 失败: %v", err))
	}

	// 获取底层 *sql.DB 以配置连接池
	sqlDB, err := DB.DB()
	if err != nil {
		panic(fmt.Sprintf("获取底层 sql.DB 失败: %v", err))
	}

	sqlDB.SetMaxIdleConns(10)           // 最大空闲连接数
	sqlDB.SetMaxOpenConns(100)          // 最大打开连接数
	sqlDB.SetConnMaxLifetime(time.Hour) // 连接最大存活时间

	log.Println("MySQL 连接初始化成功")
}

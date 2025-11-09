package config

import (
	"fmt"
	"log"
	"os"

	"github.com/NCUHOME-Y/25-HACK-1-Leaflet-BE/model"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 数据库连接实例(导入数据库时使用)
var DB *gorm.DB

func ConnectDatabase() {
	// 加载 .env 文件
	err := godotenv.Load()
	if err != nil {
		log.Println("警告: 无法加载 .env 文件")
	}

	// 从环境变量中获取数据库配置
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbCharset := os.Getenv("DB_CHARSET")
	dbParseTime := os.Getenv("DB_PARSE_TIME")
	dbLoc := os.Getenv("DB_LOC")

	// 构建 DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%s&loc=%s",
		dbUser, dbPassword, dbHost, dbPort, dbName, dbCharset, dbParseTime, dbLoc)

	// 连接数据库
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("数据库连接失败: " + err.Error())
	}
	// 自动迁移数据库模式 (创建表)
	// 先创建被引用的表，再创建有外键的表
	err = DB.AutoMigrate(&model.User{}, &model.Problem{}, &model.Solve{}, &model.Status{}, &model.EncouragementMorning{}, &model.EncouragementAfternoon{}, &model.EncouragementEvening{}, &model.Myself{})
	if err != nil {
		log.Fatal("数据库迁移失败:", err)
	}

	log.Println("数据库连接成功")
}

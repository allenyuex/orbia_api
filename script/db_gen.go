package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
	"gopkg.in/yaml.v3"
)

// Config 配置结构体
type Config struct {
	Database struct {
		MySQL struct {
			Host     string `yaml:"host"`
			Port     int    `yaml:"port"`
			Database string `yaml:"database"`
			Username string `yaml:"username"`
			Password string `yaml:"password"`
			Charset  string `yaml:"charset"`
		} `yaml:"mysql"`
	} `yaml:"database"`
}

func main() {
	// 获取当前工作目录作为项目根目录
	projectRoot, err := os.Getwd()
	if err != nil {
		log.Fatal("获取当前目录失败:", err)
	}

	// 读取配置文件
	configPath := filepath.Join(projectRoot, "conf", "config.yaml")
	config, err := loadConfig(configPath)
	if err != nil {
		log.Fatal("读取配置文件失败:", err)
	}

	// 构建数据库连接字符串
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		config.Database.MySQL.Username,
		config.Database.MySQL.Password,
		config.Database.MySQL.Host,
		config.Database.MySQL.Port,
		config.Database.MySQL.Database,
		config.Database.MySQL.Charset,
	)

	// 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("连接数据库失败:", err)
	}

	fmt.Println("数据库连接成功，开始生成代码...")

	// 创建生成器
	g := gen.NewGenerator(gen.Config{
		OutPath:           filepath.Join(projectRoot, "biz", "dal", "mysql"),
		OutFile:           "gen.go",
		ModelPkgPath:      "",  // 留空让gorm自动推断
		WithUnitTest:      false,
		FieldNullable:     true,
		FieldCoverable:    false,
		FieldSignable:     false,
		FieldWithIndexTag: false,
		FieldWithTypeTag:  true,
	})

	// 使用数据库连接
	g.UseDB(db)

	// 获取所有表名
	tables, err := getTableNames(db)
	if err != nil {
		log.Fatal("获取表名失败:", err)
	}

	fmt.Printf("发现 %d 个表: %v\n", len(tables), tables)

	// 为每个表生成模型
	for _, tableName := range tables {
		fmt.Printf("正在生成表 %s 的模型...\n", tableName)
		g.GenerateModel(tableName)
	}

	// 生成代码
	g.Execute()

	fmt.Println("代码生成完成！")
	fmt.Printf("生成的文件位置: %s\n", filepath.Join(projectRoot, "biz", "dal", "mysql"))
}

// loadConfig 加载配置文件
func loadConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	return &config, nil
}

// getTableNames 获取数据库中所有以 orbia_ 开头的表名
func getTableNames(db *gorm.DB) ([]string, error) {
	var tables []string
	
	// 查询所有以 orbia_ 开头的表
	rows, err := db.Raw("SHOW TABLES LIKE 'orbia_%'").Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, err
		}
		tables = append(tables, tableName)
	}

	return tables, nil
}
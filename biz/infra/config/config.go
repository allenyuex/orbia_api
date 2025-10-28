package config

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

var GlobalConfig *Config

type Config struct {
	Server           ServerConfig           `yaml:"server"`
	Database         DatabaseConfig         `yaml:"database"`
	Redis            RedisConfig            `yaml:"redis"`
	JWT              JWTConfig              `yaml:"jwt"`
	Log              LogConfig              `yaml:"log"`
	R2               R2Config               `yaml:"r2"`
	SMTP             SMTPConfig             `yaml:"smtp"`
	VerificationCode VerificationCodeConfig `yaml:"verification_code"`
}

type ServerConfig struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	ReadTimeout  string `yaml:"read_timeout"`
	WriteTimeout string `yaml:"write_timeout"`
}

type DatabaseConfig struct {
	MySQL MySQLConfig `yaml:"mysql"`
}

type MySQLConfig struct {
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	Database        string `yaml:"database"`
	Username        string `yaml:"username"`
	Password        string `yaml:"password"`
	Charset         string `yaml:"charset"`
	ParseTime       bool   `yaml:"parse_time"`
	Loc             string `yaml:"loc"`
	MaxIdleConns    int    `yaml:"max_idle_conns"`
	MaxOpenConns    int    `yaml:"max_open_conns"`
	ConnMaxLifetime int    `yaml:"conn_max_lifetime"`
}

type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
	PoolSize int    `yaml:"pool_size"`
}

type JWTConfig struct {
	Secret      string `yaml:"secret"`
	ExpireHours int    `yaml:"expire_hours"`
}

type LogConfig struct {
	Level      string `yaml:"level"`
	FilePath   string `yaml:"file_path"`
	MaxSize    int    `yaml:"max_size"`
	MaxBackups int    `yaml:"max_backups"`
	MaxAge     int    `yaml:"max_age"`
}

type R2Config struct {
	Endpoint                 string                     `yaml:"endpoint"`
	AccessKey                string                     `yaml:"access_key"`
	SecretKey                string                     `yaml:"secret_key"`
	Bucket                   string                     `yaml:"bucket"`
	PublicURL                string                     `yaml:"public_url"`
	UploadTokenExpireMinutes int                        `yaml:"upload_token_expire_minutes"`
	MaxFileSize              int64                      `yaml:"max_file_size"`      // 默认最大文件大小
	AllowedExtensions        map[string]ExtensionConfig `yaml:"allowed_extensions"` // 按扩展名配置
}

// ExtensionConfig 文件扩展名配置
type ExtensionConfig struct {
	MaxSize     int64  `yaml:"max_size"`     // 该扩展名的最大文件大小（字节）
	DefaultPath string `yaml:"default_path"` // 默认存储路径，如 avatars, documents
}

// SMTPConfig SMTP邮件配置
type SMTPConfig struct {
	Server   string `yaml:"server"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Email    string `yaml:"email"`
	FromName string `yaml:"from_name"`
}

// VerificationCodeConfig 验证码配置
type VerificationCodeConfig struct {
	ExpireMinutes int `yaml:"expire_minutes"`
	Length        int `yaml:"length"`
}

// LoadConfig 加载配置文件
// 根据环境变量 ORBIA_ENV 来决定加载哪个环境的配置
// 可选值: dev, prod，默认为 dev
func LoadConfig() error {
	// 获取环境变量，默认为 dev
	env := os.Getenv("ORBIA_ENV")
	if env == "" {
		env = "dev"
	}

	// 构建配置文件路径
	configPath := fmt.Sprintf("conf/%s/config.yaml", env)

	// 读取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("读取配置文件失败 (%s): %v", configPath, err)
	}

	// 替换环境变量
	content := expandEnvVars(string(data))

	// 解析配置
	config := &Config{}
	if err := yaml.Unmarshal([]byte(content), config); err != nil {
		return fmt.Errorf("解析配置文件失败: %v", err)
	}

	GlobalConfig = config
	return nil
}

// expandEnvVars 替换配置文件中的环境变量
// 支持格式: ${VAR_NAME:default_value} 或 ${VAR_NAME}
func expandEnvVars(content string) string {
	// 匹配 ${VAR_NAME:default} 或 ${VAR_NAME}
	re := regexp.MustCompile(`\$\{([^:}]+)(?::([^}]*))?\}`)

	return re.ReplaceAllStringFunc(content, func(match string) string {
		// 提取变量名和默认值
		parts := re.FindStringSubmatch(match)
		if len(parts) < 2 {
			return match
		}

		varName := strings.TrimSpace(parts[1])
		defaultValue := ""
		if len(parts) > 2 {
			defaultValue = parts[2]
		}

		// 获取环境变量值
		value := os.Getenv(varName)
		if value == "" {
			return defaultValue
		}
		return value
	})
}

// GetDSN 获取 MySQL DSN 连接字符串
func (c *MySQLConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.Database,
		c.Charset,
		c.ParseTime,
		c.Loc,
	)
}

package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"strings"
)

var globalConfig *Config

// Load 加载配置文件
func Load() (*Config, error) {
	v := viper.New()

	// 设置配置文件路径
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./configs")

	// 环境变量
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}
	fmt.Println("Using config file:", v.ConfigFileUsed())

	// 解析配置
	config := &Config{}
	if err := v.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}
	// 环境注入 > 配置文件、从环境变量覆盖敏感配置文件信息
	if val := os.Getenv("MYSQL_PASSWORD"); val != "" {
		config.Database.MySQL.Password = val
	}
	if val := os.Getenv("REDIS_PASSWORD"); val != "" {
		config.Database.Redis.Password = val
	}

	globalConfig = config
	return config, nil
}

// Get 获取全局配置
func Get() *Config {
	if globalConfig == nil {
		panic("配置未初始化，请先调用 Load() 加载配置")
	}
	return globalConfig
}

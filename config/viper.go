package config

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	ModelInfo struct {
		ModelName   string  `mapstructure:"model_name"`
		BaseUrl     string  `mapstructure:"base_url"`
		ApiKey      string  `mapstructure:"api_key"`
		Temperature float32 `mapstructure:"temperature"`
	} `mapstructure:"model_info"`
	VKDB struct {
		Ak     string `mapstructure:"ak"`
		Sk     string `mapstructure:"sk"`
		Region string `mapstructure:"region"`
		Host   string `mapstructure:"host"`
		Scheme string `mapstructure:"scheme"`
	} `mapstructure:"vkdb"`
}

var Cfg Config

func LoadConfig() {
	// 获取可执行文件的目录

	// 让 viper 在可执行文件所在目录查找 config.yaml
	viper.SetConfigName("config") // 不要加 .yaml
	viper.SetConfigType("yaml")
	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("⚠️  未找到配置文件: %v", err)
	}
	// 解析到结构体
	if err := viper.Unmarshal(&Cfg); err != nil {
		log.Fatalf("❌ 解析配置失败: %v", err)
	}

}

package main

import (
	"encoding/json"
	"os"
)

// configFile 定义了保存 AppCode 的文件名。
const configFile = "config.json"

// Config 定义了配置文件的结构。
type Config struct {
	AppCode string `json:"appcode"`
}

// LoadConfig 从 config.json 文件读取配置。
// 如果文件不存在或解析失败，它会返回一个错误。
func LoadConfig() (Config, error) {
	var config Config
	file, err := os.ReadFile(configFile)
	if err != nil {
		return config, err
	}
	err = json.Unmarshal(file, &config)
	return config, err
}

// SaveConfig 将给定的配置写入 config.json 文件。
func SaveConfig(config Config) error {
	file, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configFile, file, 0644)
}

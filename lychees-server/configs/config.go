package configs

import (
	"flag"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type ConfigType struct {
	Debug bool `yaml:"debug"`

	DomainName string `yaml:"domainName"`
	Port       int    `yaml:"port"`

	JwtSigningKey string `yaml:"jwtSigningKey"`

	PostgreSQL struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`

		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbName"`
	} `yaml:"postgresql"`

	Mongo struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`

		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbName"`
	} `yaml:"mongo"`

	Redis struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`

		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBName   int    `yaml:"dbName"`
	} `yaml:"redis"`
}

var Config *ConfigType

func init() {
	err := InitConfig()
	if err != nil {
		log.Fatal(err)
	}
}

func parseFlags() (string, error) {
	var configPath string
	// 读取配置的所在路径，如果没有就读取默认配置文件的
	flag.StringVar(&configPath, "config", "./config.yaml", "Path to config file")

	flag.Parse()

	//if err := validateConfigPath(configPath); err != nil {
	//	return "", err
	//}

	return configPath, nil
}

func InitConfig() error {
	// 获取文件配置的路径
	configPath, err := parseFlags()
	if err != nil {
		log.Fatal(err)
	}
	// 读取文件配置
	//if err != nil {
	//	log.Fatal(err)
	//}

	// 读取配置文件的
	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(data, &Config); err != nil {
		return err
	}

	return nil
}

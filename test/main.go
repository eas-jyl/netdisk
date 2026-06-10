package main


import (
	"fmt"
	"strings"
	"github.com/spf13/viper"
)

// 把配置映射到结构体
type Config struct {
    Server ServerConfig `mapstructure:"server"`
    MySQL  MySQLConfig  `mapstructure:"mysql"`
}

type ServerConfig struct {
    Host string `mapstructure:"host"`
    Port int    `mapstructure:"port"`
    Mode string `mapstructure:"mode"`
}

type MySQLConfig struct {
    Host     string `mapstructure:"host"`
    Port     int    `mapstructure:"port"`
    Username string `mapstructure:"username"`
    Password string `mapstructure:"password"`
}


// 接口：加载配置文件
func LoadConfig(path string) (*Config, error) {
	v := viper.New()

	v.SetConfigFile(path)
    v.SetDefault("server.host", "0.0.0.0")
    v.SetDefault("server.port", 8080)
    v.SetDefault("server.mode", "release")

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	// 映射到结构体
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil 

}


func main() {
	cfg, err := LoadConfig("./config.yaml")
	if err != nil {
		panic(err)
	}

	fmt.Println("配置读取成功")
	fmt.Println("server.host =", cfg.Server.Host)
	fmt.Println("server.port =", cfg.Server.Port)
	fmt.Println("server.mode =", cfg.Server.Mode)
	fmt.Println("mysql.host =", cfg.MySQL.Host)
	fmt.Println("mysql.port =", cfg.MySQL.Port)
	fmt.Println("mysql.username =", cfg.MySQL.Username)
	fmt.Println("mysql.password =", cfg.MySQL.Password)
}

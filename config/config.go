package config

import (
	"fmt"
	"os"

	"github.com/prynnekey/gin-mall/dao"
	"github.com/spf13/viper"
)

var AppConfig *Config

type Config struct {
	viper        *viper.Viper
	ServerConfig *serverConfig
	MysqlConfig  *mysqlConfig
	RedisConfig  *redisConfig
}

type serverConfig struct {
	Mode string
	Port string
}

type mysqlConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

type redisConfig struct {
	addr     string
	password string
	db       int
}

func newConfig() *Config {
	config := &Config{}
	config.init()

	return config
}

func Init() {
	AppConfig = newConfig()
}

func (c *Config) init() {
	c.initViper()
	c.readServerConfig()
	c.readMysqlConfig()
	c.readRedisConfig()

	// 配置mysql主从复制
	mysql(c.MysqlConfig)
}

func (c *Config) readRedisConfig() {
	rc := &redisConfig{}
	rc.addr = c.viper.GetString("redis.addr")
	rc.password = c.viper.GetString("redis.password")
	rc.db = c.viper.GetInt("redis.db")
	c.RedisConfig = rc
}

// 配置mysql读写分离
func mysql(m *mysqlConfig) {
	// 读(8) 主
	pathRead := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", m.Username, m.Password, m.Host, m.Port, m.Database)
	// 写(2) 从
	pathWrite := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", m.Username, m.Password, m.Host, m.Port, m.Database)
	dao.Database(pathRead, pathWrite)
}

func (c *Config) readMysqlConfig() {
	mc := &mysqlConfig{}
	mc.Host = c.viper.GetString("mysql.host")
	mc.Port = c.viper.GetString("mysql.port")
	mc.Username = c.viper.GetString("mysql.username")
	mc.Password = c.viper.GetString("mysql.password")
	mc.Database = c.viper.GetString("mysql.database")
	c.MysqlConfig = mc
}

func (c *Config) readServerConfig() {
	sc := &serverConfig{}
	sc.Mode = c.viper.GetString("server.mode")
	sc.Port = ":" + c.viper.GetString("server.port")
	c.ServerConfig = sc
}

func (c *Config) initViper() {
	c.viper = viper.New()
	dir, _ := os.Getwd()                   // 当前项目目录
	c.viper.SetConfigName("config")        // name of config file (without extension)
	c.viper.SetConfigType("yml")           // REQUIRED if the config file does not have the extension in the name
	c.viper.AddConfigPath(dir + "/config") // optionally look for config in the working directory
	err := c.viper.ReadInConfig()          // Find and read the config file
	if err != nil {                        // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
}

package config

import (
	"crypto/tls"
	"errors"
	"github.com/go-redis/redis"
	"github.com/mkideal/log"
	robfigconf "github.com/robfig/config"
	"github.com/weikaishio/redis_orm"
	"github.com/weikaishio/redis_orm_workbench/common"
	"strings"
	"time"
)

var Cfg *ConfigTools

type ConfigTools struct {
	Location string
	Dir      string

	LogDir   string
	LogLevel string

	HttpPort int

	EncryptKey string

	RedisClient  *redis.Client
	RedisORM     *redis_orm.Engine
	IsShowORMLog bool

	UserMap map[string]string
}

func NewConfig(dir, location string) (*ConfigTools, error) {
	cfg := &ConfigTools{
		UserMap: make(map[string]string),
	}

	cfg.Location = location
	cfg.Dir = dir

	return cfg, nil
}
func (cfg *ConfigTools) Reload() error {

	basic := cfg.Dir + "basic.conf"
	err := cfg.loadBasicConfig(basic)
	if err != nil {
		return err
	}

	advanced := cfg.Dir + cfg.Location + ".conf"
	err = cfg.loadAdvancedConfig(advanced)
	return err
}
func (cfg *ConfigTools) loadBasicConfig(conf string) error {
	c, err := robfigconf.ReadDefault(conf)
	if err != nil {
		return err
	}

	section := "log"
	cfg.LogDir, _ = c.String(section, "path")
	cfg.LogLevel, _ = c.String(section, "level")

	section = "web"
	cfg.HttpPort, _ = c.Int(section, "port")
	cfg.EncryptKey, _ = c.String(section, "encryptKey")
	isShowORMLog, _ := c.Int(section, "isShowORMLog")
	if isShowORMLog == 1 {
		cfg.IsShowORMLog = true
	}

	return nil
}

func (cfg *ConfigTools) loadAdvancedConfig(conf string) error {
	log.Info("loadAdvancedConfig:%s", conf)
	c, err := robfigconf.ReadDefault(conf)
	if err != nil {
		return err
	}

	// Redis配置
	cfg.RedisClient, err = readRedisConfig("redis", cfg.EncryptKey, c)
	if err != nil {
		return err
	}

	cfg.RedisORM = redis_orm.NewEngine(cfg.RedisClient)
	cfg.RedisORM.IsShowLog(cfg.IsShowORMLog)

	section := "login"
	opts, err := c.SectionOptions(section)
	if err != nil {
		log.Error("SectionOptions(login) err:%v", err)
	}else {
		for _, opt := range opts {
			pwd, _ := c.String(section, opt)
			cfg.UserMap[opt] = pwd
		}
	}
	return nil
}

func readRedisConfig(section, rc4key string, c *robfigconf.Config) (*redis.Client, error) {
	host, _ := c.String(section, "host")
	if host == "" || len(strings.Split(host, ":")) != 2 {
		return nil, errors.New("rd.host is null")
	}
	portStr := strings.Split(host, ":")[1]
	var port int32
	err := redis_orm.SetInt32FromStr(&port, portStr)
	if err != nil || port <= 0 {
		return nil, errors.New("rd.host port error")
	}
	password, _ := c.String(section, "password")
	if rc4key != "" {
		passwordArray := common.DescryptRC4Base64(password, rc4key)
		password = string(passwordArray)
	}
	dbIndex, _ := c.Int(section, "database")
	options := redis.Options{
		Addr:               host,
		Password:           password,
		DB:                 dbIndex,
		DialTimeout:        10 * time.Second,
		ReadTimeout:        10 * time.Second,
		WriteTimeout:       10 * time.Second,
		IdleTimeout:        60 * time.Second,
		IdleCheckFrequency: 15 * time.Second,
	}
	if port == 6380 { //azure的规则，就把6380当ssl处理
		options.TLSConfig = &tls.Config{
			ServerName: strings.Split(host, ":")[0],
		}
	}
	client := redis.NewClient(&options)
	ping, err := client.Ping().Result()
	if err != nil {
		log.Error("Redis failed to ping err: %v", err)
		client.Close()
		return nil, err
	}
	if strings.ToLower(ping) != "pong" {
		log.Warn("Redis unexpected ping response, pong:%s", ping)
		return nil, errors.New("redis failed to ping")
	}
	return client, nil
}

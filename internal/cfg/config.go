package cfg

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	gormLogger "gorm.io/gorm/logger"
)

var Cfg *Config

type CfgEnvironment string

const (
	Local       CfgEnvironment = "local"
	Development CfgEnvironment = "development"
	Staging     CfgEnvironment = "staging"
	Production  CfgEnvironment = "production"
)

type Config struct {
	Server      ServerCfg      `mapstructure:"server"`
	Environment CfgEnvironment `mapstructure:"environment" `
	Log         LogCfg         `mapstructure:"log"`
	DB          DatabaseConfig `mapstructure:"db"`
	Cache       CacheConfig    `mapstructure:"cache"`
	Verbose     bool           `mapstructure:"verbose" `
}

type ServerCfg struct {
	Port      int     `mapstructure:"port"`
	Host      string  `mapstructure:"host"`
	AccessLog bool    `mapstructure:"accesslog"`
	Cors      CorsCfg `mapstructure:"cors"`
	JWT       JWTCfg  `mapstructure:"jwt"`
}

type JWTCfg struct {
	Secret         []byte        `mapstructure:"secret"`
	Timeout        time.Duration `mapstructure:"timeout"`
	RefreshTimeout time.Duration `mapstructure:"refresh_timeout"`
}

type CorsCfg struct {
	Origins        []string `mapstructure:"origins"`
	Methods        []string `mapstructure:"methods"`
	AllowedHeaders []string `mapstructure:"allowed-headers"`
}

type LogCfg struct {
	Level     slog.Level `mapstructure:"level" `
	ErrorPath string     `mapstructure:"error-path" `
	InfoPath  string     `mapstructure:"info-path" `
	DebugPath string     `mapstructure:"debug-path" `
}

type CacheConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	DB       int    `mapstructure:"db"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}

type DatabaseConfig struct {
	Host                   string              `mapstructure:"host"`
	Port                   int                 `mapstructure:"port"`
	Name                   string              `mapstructure:"name"`
	User                   string              `mapstructure:"user"`
	Password               string              `mapstructure:"password"`
	MaxOpenConnectionCount int                 `mapstructure:"max_open_connection_count"`
	GormLogLevel           gormLogger.LogLevel `mapstructure:"gorm-log-level"`
}

func LoadConfig(cfgPath string, parsedFlags *pflag.FlagSet) error {
	var err error
	sync.OnceFunc(func() { err = loadConfig(cfgPath, parsedFlags) })()
	return err
}

func loadConfig(cfgPath string, parsedFlags *pflag.FlagSet) error {

	v := viper.NewWithOptions()
	v.SetConfigFile(cfgPath)
	v.AddConfigPath(".")
	err := v.ReadInConfig()
	if err != nil {
		return err
	}

	v.AutomaticEnv()
	if parsedFlags != nil {
		if err := v.BindPFlags(parsedFlags); err != nil {
			return fmt.Errorf("failed binding flags: %w", err)
		}
	}

	var C Config
	err = v.Unmarshal(&C)
	if err != nil {
		return err
	}
	Cfg = &C
	return nil
}

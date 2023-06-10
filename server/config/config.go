package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm/logger"
	"os"
	"path/filepath"
	"runtime"
	"server/util"
	"strings"
)

// 系统配置，对应yml
// viper内置了mapstructure, yml文件用"-"区分单词, 转为驼峰方便

// 全局配置变量
var Conf = new(config)

type config struct {
	System    *SystemConfig    `mapstructure:"system" json:"system"`
	Database  *DatabaseConfig  `mapstructure:"database" json:"database"`
	Logs      *LogsConfig      `mapstructure:"logs" json:"logs"`
	Casbin    *CasbinConfig    `mapstructure:"casbin" json:"casbin"`
	Jwt       *JwtConfig       `mapstructure:"jwt" json:"jwt"`
	RateLimit *RateLimitConfig `mapstructure:"rate-limit" json:"rateLimit"`
	Ssl       *Ssl             `mapstructure:"ssl" json:"ssl"`
	Mdns      *Mdns            `mapstructure:"mdns" json:"mdns"`
	Env       *Env             `json:"env"`
}

// 设置读取配置信息
func InitConfig() {
	workDir, err := os.Getwd()
	if err != nil {
		panic(fmt.Errorf("读取应用目录失败:%s \n", err))
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(workDir)
	// 读取配置信息
	err = viper.ReadInConfig()

	// 热更新配置
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		// 将读取的配置信息保存至全局变量Conf
		if err := viper.Unmarshal(Conf); err != nil {
			panic(fmt.Errorf("初始化配置文件失败:%s \n", err))
		}
		// 读取rsa key
		Conf.System.RSAPublicBytes = util.RSAReadKeyFromFile(Conf.System.RSAPublicKey)
		Conf.System.RSAPrivateBytes = util.RSAReadKeyFromFile(Conf.System.RSAPrivateKey)
	})

	if err != nil {
		panic(fmt.Errorf("读取配置文件失败:%s \n", err))
	}
	// 将读取的配置信息保存至全局变量Conf
	if err := viper.Unmarshal(Conf); err != nil {
		panic(fmt.Errorf("初始化配置文件失败:%s \n", err))
	}
	// 读取rsa key
	Conf.System.RSAPublicBytes = util.RSAReadKeyFromFile(Conf.System.RSAPublicKey)
	Conf.System.RSAPrivateBytes = util.RSAReadKeyFromFile(Conf.System.RSAPrivateKey)
	// 处理特殊的url-path-prefix情况
	if Conf.System.UrlPathPrefix == "" || Conf.System.UrlPathPrefix == "/" {
		Conf.System.UrlPathPrefix = "/"
	} else {
		Conf.System.UrlPathPrefix = strings.Trim(Conf.System.UrlPathPrefix, "/")
	}
	// 读取全局环境变量
	Conf.Env = &Env{}
	Conf.Env.WorkDir = workDir
	Conf.Env.Ip = util.GetDefaultIp()
	Conf.Env.HostName, _ = os.Hostname()

	// 可执行文件路径
	execPath, err := filepath.Abs(os.Args[0])
	if err != nil {
		panic(fmt.Errorf("获取可执行文件路径失败, error: %v", err))
	} else {
		// md5 值
		md5Str, err := util.GetFileMd5(execPath)
		if err != nil {
			panic(fmt.Errorf("获取可执行文件md5失败, error: %v", err))
		}
		Conf.Env.BinFileHash = md5Str
	}
	// 获取OS类型
	Conf.Env.OsType = runtime.GOOS
}

type SystemConfig struct {
	Mode            string `mapstructure:"mode" json:"mode"`
	UrlPathPrefix   string `mapstructure:"url-path-prefix" json:"urlPathPrefix"`
	Port            int    `mapstructure:"port" json:"port"`
	InitData        bool   `mapstructure:"init-data" json:"initData"`
	RSAPublicKey    string `mapstructure:"rsa-public-key" json:"rsaPublicKey"`
	RSAPrivateKey   string `mapstructure:"rsa-private-key" json:"rsaPrivateKey"`
	RSAPublicBytes  []byte `mapstructure:"-" json:"-"`
	RSAPrivateBytes []byte `mapstructure:"-" json:"-"`
}

type LogsConfig struct {
	Level           zapcore.Level `mapstructure:"level" json:"level"`
	Path            string        `mapstructure:"path" json:"path"`
	MaxSize         int           `mapstructure:"max-size" json:"maxSize"`
	MaxBackups      int           `mapstructure:"max-backups" json:"maxBackups"`
	MaxAge          int           `mapstructure:"max-age" json:"maxAge"`
	Compress        bool          `mapstructure:"compress" json:"compress"`
	EnableReportLog bool          `mapstructure:"enable-report-log" json:"enableReportLog"`
}

type CasbinConfig struct {
	ModelPath string `mapstructure:"model-path" json:"modelPath"`
}

type JwtConfig struct {
	Realm      string `mapstructure:"realm" json:"realm"`
	Key        string `mapstructure:"key" json:"key"`
	Timeout    int    `mapstructure:"timeout" json:"timeout"`
	MaxRefresh int    `mapstructure:"max-refresh" json:"maxRefresh"`
}

type RateLimitConfig struct {
	Enable       bool  `mapstructure:"enable" json:"enable"`
	FillInterval int64 `mapstructure:"fill-interval" json:"fillInterval"`
	Capacity     int64 `mapstructure:"capacity" json:"capacity"`
}

type Ssl struct {
	Enable   bool
	KeyFile  string
	CertFile string
}

type Mdns struct {
	Enable bool
}

type Env struct {
	Ip          string
	HostName    string
	WorkDir     string
	BinFileHash string
	OsType      string
}

type DatabaseConfig struct {
	Driver  string          `mapstructure:"driver" json:"driver"`
	Source  string          `mapstructure:"source" json:"source"`
	LogMode logger.LogLevel `mapstructure:"log-mode" json:"logMode"`
}

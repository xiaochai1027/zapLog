package zlog

import (
	"github.com/spf13/viper"
	"time"
)

type ZlogCfg struct {
	Level         string
	Compress      bool
	MaxAge        int
	MaxSize       int
	MaxBackups    int
	FileName      string
	FlushInterval time.Duration
	BuffSize      int
	AddCaller     bool
	AddSkip       int
}

func InitConfig(confPath, confName, confType string) (error, *ZlogCfg) {
	//设置默认值
	setDefault()
	viper.AddConfigPath(confPath)
	viper.SetConfigName(confName)
	viper.SetConfigType(confType)

	err := viper.ReadInConfig()
	if err != nil {
		return err, nil
	}

	return err, getZlogCnf()
}

func InitConfigByFilePath(configFile string) (error, *ZlogCfg) {
	//设置默认值
	setDefault()
	viper.SetConfigFile(configFile)

	err := viper.ReadInConfig()
	if err != nil {
		return err, nil
	}

	return err, getZlogCnf()
}

func getZlogCnf() *ZlogCfg {
	zlogcnf := &ZlogCfg{
		Level:         viper.GetString("zlog.level"),
		Compress:      viper.GetBool("zlog.compress"),
		MaxAge:        viper.GetInt("zlog.max_age"),
		MaxSize:       viper.GetInt("zlog.max_size"),
		MaxBackups:    viper.GetInt("zlog.max_backups"),
		FileName:      viper.GetString("zlog.file_name"),
		FlushInterval: viper.GetDuration("zlog.flush_interval"),
		BuffSize:      viper.GetInt("zlog.buff_size"),
		AddCaller:     viper.GetBool("zlog.add_caller"),
		AddSkip:       viper.GetInt("zlog.add_skip"),
	}
	return zlogcnf
}

func setDefault() {
	viper.SetDefault("zlog.file_name", "./zlog")
	viper.SetDefault("zlog.compress", false)
	viper.SetDefault("zlog.max_age", 15)
	viper.SetDefault("zlog.max_size", 500)
	viper.SetDefault("zlog.max_backups", 30)
	viper.SetDefault("zlog.flush_interval", 300)
	viper.SetDefault("zlog.buff_size", 10)
	viper.SetDefault("zlog.add_caller", true)
	viper.SetDefault("zlog.add_skip", 1)
}

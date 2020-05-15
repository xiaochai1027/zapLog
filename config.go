package zlog

import (
	"flag"
	"time"

	"gitlab.xuanke.com/live/kconf"
)

type ZlogCfg struct {
	Level         string        `kconf:"level"`
	Compress      bool          `kconf:"compress"`
	MaxAge        int           `kconf:"max_age"`
	MaxSize       int           `kconf:"max_size"`
	MaxBackups    int           `kconf:"max_backups"`
	FileName      string        `kconf:"file_name"`
	FlushInterval time.Duration `kconf:"flush_interval"`
	BuffSize      int           `kconf:"buff_size"`
	AddCaller     bool          `kconf:"add_caller"`
	AddSkip       int           `kconf:"add_skip"`
}

type Config struct {
	Zlog ZlogCfg `kconf:"zlog"`
}

func InitConfig(path string) (error, ZlogCfg) {
	flag.Parse()

	var config Config
	err := kconf.LoadConfig(path, &config)
	if err != nil {
		return err, ZlogCfg{}
	}

	return nil, config.Zlog
}

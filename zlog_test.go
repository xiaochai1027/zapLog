package zlog

import (
	"fmt"
	"go.uber.org/zap/zapcore"
	"testing"
)

//Level         string        `kconf:"level"`
//Compress      bool          `kconf:"compress"`
//MaxAge        int           `kconf:"max_age"`
//MaxSize       int           `kconf:"max_size"`
//MaxBackups    int           `kconf:"max_backups"`
//FileName      string        `kconf:"file_name"`
//FlushInterval time.Duration `kconf:"flush_interval"`
//BuffSize      int           `kconf:"buff_size"`
//AddCaller     bool          `kconf:"add_caller"`
//AddSkip       int           `kconf:"add_skip"`
func Test_zlog(t *testing.T) {
	logger := ZlogInitSplitFile(&ZlogCfg{
		Level:         "DEBUG",
		Compress:      false,
		MaxAge:        10,
		MaxSize:       1,
		MaxBackups:    30,
		FileName:      "./log/signal.log",
		FlushInterval: 300,
		BuffSize:      1,
		AddCaller:     true,
		AddSkip:       1,
	})
	defer logger.Sync()
	for i := 0; i < 1000; i++ {
		logger.Debug("is debug")
	}
	//logger.Info("is info")
	//logger.Error("is error")
	//logger.Panic("is panic")
	fmt.Println(replaceFileName("./log/signnal.log", zapcore.DebugLevel))
	t.Log("1")

}

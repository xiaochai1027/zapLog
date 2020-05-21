package zlog

import (
	"fmt"
	"go.uber.org/zap/zapcore"
	"testing"
)

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
	logger.Info("is info")
	logger.Error("is error")
	logger.Panic("is panic")
	fmt.Println(replaceFileName("./log/signnal.log", zapcore.DebugLevel))
	t.Log("1")

}

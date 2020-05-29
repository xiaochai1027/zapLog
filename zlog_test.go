package zlog

import (
	"go.uber.org/zap"
	"testing"
)

func Test_zlog(t *testing.T) {
	//logger := ZlogInitSplitFile(&ZlogCfg{
	//	Level:         "DEBUG",
	//	Compress:      false,
	//	MaxAge:        10,
	//	MaxSize:       1,
	//	MaxBackups:    30,
	//	FileName:      "./log/signal.log",
	//	FlushInterval: 300,
	//	BuffSize:      1,
	//	AddCaller:     true,
	//	AddSkip:       1,
	//})
	//defer logger.Sync()
	//for i := 0; i < 1000; i++ {
	//	logger.Debug("is debug")
	//}
	//logger.Info("is info")
	//logger.Error("is error")
	//logger.Panic("is panic")
	//fmt.Println(replaceFileName("./log/signnal.log", zapcore.DebugLevel))
	//t.Log("1")
	//viper.SetDefault("zlog.file_name","./zlog")
	//viper.SetConfigName("zlog")
	//viper.SetConfigType("yaml")
	//viper.AddConfigPath(".")
	//err := viper.ReadInConfig()
	//if err != nil {
	//	fmt.Printf("err = %v",err)
	//}
	//fmt.Println(viper.Get("zlog.level"))
	//fmt.Println(viper.Get("zlog.file_name"))
	//viper.WatchConfig()

}

func Test_zlog_config(t *testing.T){
	err,zlogConf := InitConfig("./","zloga","yaml")
	if err != nil {
		t.Error(err)
		return
	}
	zlog := ZlogInitSplitFile(zlogConf)
	defer zlog.Sync()
	zlog.Info("info",zap.String("name","info..."))

}

package zlog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type CutConf struct {
	FileName   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
	LocalTime  bool
	BufferSize int //单位为m
}

func ZapCut(c CutConf, z zapcore.EncoderConfig, l zap.AtomicLevel) zapcore.Core {
	lum := &Logger{Filename: c.FileName, MaxSize: c.MaxSize, MaxBackups: c.MaxBackups, MaxAge: c.MaxAge, Compress: c.Compress, LocalTime: c.LocalTime, bufferSize: c.BufferSize}
	zlvl := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= l.Level()
	})
	code := zapcore.NewJSONEncoder(z)
	return zapcore.NewCore(code, zapcore.AddSync(lum), zlvl)
}

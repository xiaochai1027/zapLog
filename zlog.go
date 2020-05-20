package zlog

import (
	"bytes"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var timestamp int64
var timeStr string
var mu sync.Mutex

type MultLogger map[string]*Zlog

type Zlog struct {
	Origin    *zap.Logger
	log       *zap.Logger
	AtomLv    zap.AtomicLevel
	FileExist bool
}

func (logger *Zlog) Debug(msg string, fields ...zapcore.Field) {
	logger.log.Debug(msg, fields...)
}

func (logger *Zlog) Info(msg string, fields ...zapcore.Field) {
	logger.log.Info(msg, fields...)
}

func (logger *Zlog) Warn(msg string, fields ...zapcore.Field) {
	logger.log.Warn(msg, fields...)
}
func (logger *Zlog) Error(msg string, fields ...zapcore.Field) {
	logger.log.Error(msg, fields...)
}

func (logger *Zlog) Panic(msg string, fields ...zapcore.Field) {
	logger.log.DPanic(msg, fields...)
}

func (logger *Zlog) Sync() {
	if logger.FileExist {
		logger.log.Sync()
	}
}

func (logger *Zlog) LogAppend(fields ...zapcore.Field) *Zlog {
	return &Zlog{log: logger.log.With(fields...), Origin: logger.Origin, AtomLv: logger.AtomLv}
}
func (logger *Zlog) CopyLogWithInfo(fields ...zapcore.Field) *Zlog {
	return &Zlog{log: logger.Origin.With(fields...), Origin: logger.Origin, AtomLv: logger.AtomLv}
}

func ZlogInit() *Zlog {
	//atomLv := zap.NewAtomicLevel()
	//atomLv.SetLevel(zapcore.DebugLevel)
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = CustomTimeEncoder
	dCut := CutConf{
		MaxAge:     15,
		MaxSize:    1000,
		MaxBackups: 30,
		LocalTime:  true,
		FileName:   "./log/debug.log",
	}
	debugC := ZapCut(dCut, encoderCfg, zap.NewAtomicLevelAt(zapcore.DebugLevel))
	eCut := CutConf{
		MaxAge:     15,
		MaxSize:    1000,
		MaxBackups: 30,
		LocalTime:  true,
		FileName:   "./log/error.log",
	}
	errorC := ZapCut(eCut, encoderCfg, zap.NewAtomicLevelAt(zapcore.ErrorLevel))
	//logger := zap.New(debugC, zap.AddStacktrace(zapcore.DPanicLevel), zap.ErrorOutput(nil))
	core := zapcore.NewTee(
		debugC,
		errorC,
	)
	logger := zap.New(core, zap.AddStacktrace(zapcore.DPanicLevel), zap.ErrorOutput(nil))
	log := &Zlog{Origin: logger, log: logger}
	log.Info("zlog init.......")
	log.FileExist = true
	go flushDaemon(log.Sync, 3*time.Second)
	return log
}

func ZlogInitByCfg(zcfg *ZlogCfg) *Zlog {
	atomLv := zap.NewAtomicLevel()
	level := zapcore.DebugLevel
	switch strings.ToLower(zcfg.Level) {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	default:
		level = zap.InfoLevel
	}
	atomLv.SetLevel(level)

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = CustomTimeEncoder
	cutCfg := CutConf{}
	cutCfg.Compress = zcfg.Compress     //是否压缩
	cutCfg.MaxAge = zcfg.MaxAge         //保留旧日志文件的最大天数
	cutCfg.MaxSize = zcfg.MaxSize       //文件大小 M
	cutCfg.MaxBackups = zcfg.MaxBackups //保留的旧日志文件的最大数量
	cutCfg.LocalTime = true
	cutCfg.FileName = zcfg.FileName   //文件位置
	cutCfg.BufferSize = zcfg.BuffSize //buffio大小，单位为m
	core := ZapCut(cutCfg, encoderCfg, atomLv)
	var logger *zap.Logger
	if zcfg.AddCaller {
		logger = zap.New(core, zap.AddStacktrace(zapcore.DPanicLevel), zap.ErrorOutput(nil), zap.AddCaller(), zap.AddCallerSkip(zcfg.AddSkip))
	} else {
		logger = zap.New(core, zap.AddStacktrace(zapcore.DPanicLevel), zap.ErrorOutput(nil))
	}
	log := &Zlog{Origin: logger, log: logger, AtomLv: atomLv}
	log.Info("zlog init.......")
	log.FileExist = true
	go flushDaemon(log.Sync, zcfg.FlushInterval*time.Millisecond)
	return log
}

//分级别存储文件
func ZlogInitSplitFile(zcfg *ZlogCfg) *Zlog {
	level := zapcore.InfoLevel
	switch strings.ToLower(zcfg.Level) {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zap.InfoLevel
	}
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = CustomTimeEncoder
	var cores []zapcore.Core
	cutCfg := CutConf{
		Compress:   zcfg.Compress,   //是否压缩
		MaxAge:     zcfg.MaxAge,     //保留旧日志文件的最大天数
		MaxSize:    zcfg.MaxSize,    //文件大小 M
		MaxBackups: zcfg.MaxBackups, //保留的旧日志文件的最大数量
		LocalTime:  true,
		FileName:   zcfg.FileName, //文件位置
		BufferSize: zcfg.BuffSize, //buffio大小，单位为m
	}
	switch {
	case level == zapcore.DebugLevel:
		core := cutCfg
		core.FileName = replaceFileName(core.FileName, zapcore.DebugLevel)
		cores = append(cores, ZapCut(core, encoderCfg, zap.NewAtomicLevelAt(zapcore.DebugLevel)))
		fallthrough
	case level == zapcore.InfoLevel:
		core := cutCfg
		core.FileName = replaceFileName(core.FileName, zapcore.InfoLevel)
		cores = append(cores, ZapCut(core, encoderCfg, zap.NewAtomicLevelAt(zapcore.InfoLevel)))
		fallthrough
	case level == zapcore.WarnLevel:
		core := cutCfg
		core.FileName = replaceFileName(core.FileName, zapcore.WarnLevel)
		cores = append(cores, ZapCut(core, encoderCfg, zap.NewAtomicLevelAt(zapcore.WarnLevel)))
		fallthrough
	case level == zapcore.ErrorLevel:
		core := cutCfg
		core.FileName = replaceFileName(core.FileName, zapcore.ErrorLevel)
		cores = append(cores, ZapCut(core, encoderCfg, zap.NewAtomicLevelAt(zapcore.ErrorLevel)))
		fallthrough
	case level == zapcore.DPanicLevel:
		core := cutCfg
		core.FileName = replaceFileName(core.FileName, zapcore.DPanicLevel)
		cores = append(cores, ZapCut(core, encoderCfg, zap.NewAtomicLevelAt(zapcore.DPanicLevel)))
	default:
		core := cutCfg
		core.FileName = replaceFileName(core.FileName, zapcore.InfoLevel)
		cores = append(cores, ZapCut(core, encoderCfg, zap.NewAtomicLevelAt(zapcore.InfoLevel)))
	}

	var logger *zap.Logger
	if zcfg.AddCaller {
		logger = zap.New(zapcore.NewTee(cores...), zap.AddStacktrace(zapcore.DPanicLevel), zap.AddCaller(), zap.AddCallerSkip(zcfg.AddSkip))
	} else {
		logger = zap.New(zapcore.NewTee(cores...), zap.AddStacktrace(zapcore.DPanicLevel))
	}
	log := &Zlog{Origin: logger, log: logger}
	log.Info("zlog init.......")
	log.FileExist = true
	go flushDaemon(log.Sync, zcfg.FlushInterval*time.Millisecond)
	return log
}

func CustomTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {

	if timestamp != t.Unix() {
		mu.Lock()
		if timestamp != t.Unix() {
			atomic.CompareAndSwapInt64(&timestamp, timestamp, t.Unix())
			timeStr = t.Format("2006-01-02 15:04:05")
		}
		mu.Unlock()

	}

	buffer := bytes.Buffer{}
	buffer.WriteString(timeStr)
	buffer.WriteString(".")
	millisecond := t.UnixNano() / 1e6 % 1000
	if millisecond >= 100 {
		buffer.WriteString(strconv.FormatInt(millisecond, 10))
	} else if millisecond < 100 && millisecond >= 10 {
		buffer.WriteString("0")
		buffer.WriteString(strconv.FormatInt(millisecond, 10))
	} else {
		buffer.WriteString("00")
		buffer.WriteString(strconv.FormatInt(millisecond, 10))
	}
	enc.AppendString(buffer.String())
}

func getFileName(patch string) string {
	name := filepath.Base(patch) //获取文件名带后缀
	return strings.TrimSuffix(name, path.Ext(name))
}

func replaceFileName(path string, lv zapcore.Level) string {
	name := getFileName(path)
	newName := name + "_" + lv.String()
	return strings.ReplaceAll(path, name, newName)
}

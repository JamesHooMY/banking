package log

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
	"go.elastic.co/apm/module/apmzap/v2"
	"go.elastic.co/apm/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func InitLogger(tracer *apm.Tracer) (*zap.SugaredLogger, error) {
	logMode := zapcore.InfoLevel

	// local file log
	fileCore := zapcore.NewCore(getEncoder(), getWriteSyncer(), logMode)

	// terminal log
	consoleConfig := zap.NewDevelopmentConfig()
	consoleConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleCore := zapcore.NewCore(zapcore.NewConsoleEncoder(consoleConfig.EncoderConfig), zapcore.AddSync(os.Stdout), logMode)

	// apm
	apmzapCore := &apmzap.Core{Tracer: tracer}
	apmzapWrapCore := apmzapCore.WrapCore(fileCore)

	// combine three cores
	core := zapcore.NewTee(fileCore, consoleCore, apmzapWrapCore)

	return zap.New(core).Sugar(), nil
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeTime = func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(t.Local().Format(time.DateTime))
	}

	return zapcore.NewJSONEncoder(encoderConfig)
}

func getWriteSyncer() zapcore.WriteSyncer {
	stSeparator := string(filepath.Separator)
	stRootDir, err := os.Getwd()
	if err != nil {
		panic(fmt.Errorf("Fatal error get root dir: %s \n", err))
	}

	// log file store path
	stLogFilePath := stRootDir + stSeparator + "log" + stSeparator + time.Now().Format(time.DateOnly) + ".log"

	// use lumberjack to split log file
	lumberjackSyncer := &lumberjack.Logger{
		Filename:   stLogFilePath,
		MaxSize:    viper.GetInt("log.maxSize"),
		MaxBackups: viper.GetInt("log.maxBackups"),
		MaxAge:     viper.GetInt("log.maxAge"),
		Compress:   viper.GetBool("log.compress"),
	}

	return zapcore.AddSync(lumberjackSyncer)
}

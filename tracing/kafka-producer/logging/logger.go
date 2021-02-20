package logging

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Config *zap.Config
var Log *zap.Logger
var SugaredLog *zap.SugaredLogger

func InitGlobalLogger() error {
	fmt.Println("Initialize global logger")

	cfg, cfgErr := loadConfig()
	if cfgErr != nil {
		return cfgErr
	}

	level, levelErr := getZapLevel(cfg.level)
	if levelErr != nil {
		return levelErr
	}

	Config = &zap.Config{
		Encoding:         cfg.encoding,
		Level:            zap.NewAtomicLevelAt(level),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig:    buildEncoderConfig(level),
	}
	Log, _ = Config.Build()
	SugaredLog = Log.Sugar()
	return nil
}

func getZapLevel(levelString string) (zapcore.Level, error) {
	level := zapcore.InfoLevel
	err := level.Set(levelString)
	if err != nil {
		return zapcore.InfoLevel, err
	}
	return level, nil
}

func buildEncoderConfig(level zapcore.Level) zapcore.EncoderConfig {
	if level == zapcore.DebugLevel {
		return zapcore.EncoderConfig{
			TimeKey:      "time",
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			LevelKey:     "level",
			EncodeLevel:  zapcore.CapitalLevelEncoder,
			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
			MessageKey:   "message",
		}
	} else {
		return zapcore.EncoderConfig{
			TimeKey:     "time",
			EncodeTime:  zapcore.ISO8601TimeEncoder,
			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,
			MessageKey:  "message",
		}
	}
}

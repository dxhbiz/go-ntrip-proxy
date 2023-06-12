package log

import (
	"os"

	"github.com/dxhbiz/go-ntrip-proxy/pkg/config"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// InitLog init logger instance
func Init(cfg config.LogConfig) {
	var level zapcore.Level
	switch cfg.Level {
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

	writer := getWriter(cfg)
	encoder := getEncoder(cfg)

	writers := make([]zapcore.WriteSyncer, 0, 2)
	writers = append(writers, writer)
	if cfg.Development {
		writers = append(writers, zapcore.AddSync(os.Stdout))
	}

	// create multi writeSyncer
	writerSyncer := zapcore.NewMultiWriteSyncer(writers...)

	core := zapcore.NewCore(encoder, writerSyncer, level)

	zLog := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	logger = zLog.Sugar()
}

func getWriter(cfg config.LogConfig) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   cfg.Filename,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
	}

	return zapcore.AddSync(lumberJackLogger)
}

func getEncoder(cfg config.LogConfig) zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()

	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	return zapcore.NewConsoleEncoder(encoderConfig)
}

package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

// NewLogger creates and returns a new SugaredLogger instance.
// logLevel: It uses the same text representation as the static zapcore.Levels ("debug", "info", "warn", "error", "dpanic", "panic", and "fatal").
func NewLogger(logLevel string) *zap.SugaredLogger {
	stdout := zapcore.AddSync(os.Stdout)

	file := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    100, // megabytes
		MaxBackups: 3,
		MaxAge:     28, // days
	})

	Proplevel := zap.InfoLevel
	if err := Proplevel.UnmarshalText([]byte(logLevel)); err != nil {
		return nil
	}

	level := zap.NewAtomicLevelAt(Proplevel)

	productionCfg := zap.NewProductionEncoderConfig()
	productionCfg.TimeKey = "timestamp"
	productionCfg.StacktraceKey = "stacktrace"
	productionCfg.MessageKey = "message"
	productionCfg.EncodeTime = zapcore.RFC3339TimeEncoder
	productionCfg.EncodeDuration = zapcore.StringDurationEncoder
	productionCfg.EncodeCaller = zapcore.FullCallerEncoder

	developmentCfg := zap.NewDevelopmentEncoderConfig()
	developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	productionCfg.EncodeCaller = zapcore.ShortCallerEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)
	fileEncoder := zapcore.NewJSONEncoder(productionCfg)

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, stdout, level),
		zapcore.NewCore(fileEncoder, file, level),
	)

	return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel)).Sugar()
}

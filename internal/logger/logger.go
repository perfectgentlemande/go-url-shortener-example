package logger

import (
	"fmt"
	"os"

	"github.com/perfectgentlemande/go-url-shortener-example/internal/service"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	log *zap.Logger
}

type Config struct {
	Encoder string
	Level   string
}

const (
	encoderConsole = "console"
	encoderJSON    = "json"

	levelDebug = "debug"
	levelInfo  = "info"
	levelWarn  = "warn"
	levelError = "error"
)

func selectLevel(level string) zapcore.Level {
	l, ok := map[string]zapcore.Level{
		levelDebug: zap.DebugLevel,
		levelInfo:  zap.InfoLevel,
		levelWarn:  zap.WarnLevel,
		levelError: zap.ErrorLevel,
	}[level]

	if ok {
		return l
	}

	return zap.DebugLevel
}

func selectEncoder(encoder string) func(zapcore.EncoderConfig) zapcore.Encoder {
	e, ok := map[string]func(zapcore.EncoderConfig) zapcore.Encoder{
		encoderConsole: zapcore.NewConsoleEncoder,
		encoderJSON:    zapcore.NewJSONEncoder,
	}[encoder]

	if ok {
		return e
	}

	return zapcore.NewConsoleEncoder
}

func New(conf *Config) (service.Logger, error) {
	return service.Logger(&Logger{
		log: zap.New(zapcore.NewCore(selectEncoder(conf.Encoder)(zap.NewDevelopmentEncoderConfig()), os.Stdout, selectLevel(conf.Level))),
	}), nil
}

func NewDefaultLogger() service.Logger {
	return service.Logger(&Logger{
		log: zap.New(zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()), os.Stdout, zap.DebugLevel)),
	})
}

func (l *Logger) WithField(key string, value interface{}) service.Logger {
	return service.Logger(&Logger{
		log: l.log.With(zapcore.Field{
			Type:      zapcore.ReflectType,
			Key:       key,
			Interface: value,
		}),
	})
}
func (l *Logger) WithError(err error) service.Logger {
	return service.Logger(&Logger{
		log: l.log.With(zapcore.Field{
			Type:      zapcore.ErrorType,
			Key:       "error",
			Interface: err,
		}),
	})
}
func (l *Logger) WithFields(fields map[string]interface{}) service.Logger {
	fieldsSlice := make([]string, 0, len(fields))
	for k := range fields {
		fieldsSlice = append(fieldsSlice, k)
	}

	if len(fieldsSlice) == 0 {
		return service.Logger(l)
	}

	newLog := &Logger{
		log: l.log.With(zapcore.Field{
			Type:      zapcore.ReflectType,
			Key:       fieldsSlice[0],
			Interface: fields[fieldsSlice[0]]}),
	}

	for i := 1; i < len(fieldsSlice); i++ {
		newLog.log = newLog.log.With(zapcore.Field{
			Type:      zapcore.ReflectType,
			Key:       fieldsSlice[i],
			Interface: fields[fieldsSlice[i]]})
	}

	return service.Logger(newLog)
}

func (l *Logger) Debug(message interface{}) {
	l.log.Debug(fmt.Sprintf("%v", message))
}

func (l *Logger) Error(message interface{}) {
	l.log.Info(fmt.Sprintf("%v", message))
}

func (l *Logger) Info(message interface{}) {
	l.log.Info(fmt.Sprintf("%v", message))
}

func (l *Logger) Warn(message interface{}) {
	l.log.Warn(fmt.Sprintf("%v", message))
}

func ProvideLogger(conf *Config) (service.Logger, error) {
	log, err := New(conf)
	if err != nil {
		return nil, fmt.Errorf("cannot create logger: %w", err)
	}

	return log, nil
}

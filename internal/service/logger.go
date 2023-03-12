package service

type Logger interface {
	Debug(message interface{})
	Error(message interface{})
	Info(message interface{})
	Warn(message interface{})
	WithField(key string, value interface{}) Logger
	WithError(err error) Logger
	WithFields(fields map[string]interface{}) Logger
}

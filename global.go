package log

var global = New(AddCallerSkip(1))

func SetOptions(opts ...Option) {
	global = global.WithOptions(opts...)
}

func Debug(args ...interface{}) {
	global.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	global.Debugf(format, args...)
}

func Debugln(args ...interface{}) {
	global.Debugln(args...)
}

func Debugw(msg string, keysAndValues ...interface{}) {
	global.Debugw(msg, keysAndValues...)
}

func Info(args ...interface{}) {
	global.Info(args...)
}

func Infof(format string, args ...interface{}) {
	global.Infof(format, args...)
}

func Infoln(args ...interface{}) {
	global.Infoln(args...)
}

func Infow(msg string, keysAndValues ...interface{}) {
	global.Infow(msg, keysAndValues...)
}

func Warn(args ...interface{}) {
	global.Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	global.Warnf(template, args...)
}

func Warnln(args ...interface{}) {
	global.Warnln(args...)
}

func Warnw(msg string, keysAndValues ...interface{}) {
	global.Warnw(msg, keysAndValues...)
}

func Error(args ...interface{}) {
	global.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	global.Errorf(template, args...)
}

func Errorln(args ...interface{}) {
	global.Errorln(args...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	global.Errorw(msg, keysAndValues...)
}

func DPanic(args ...interface{}) {
	global.DPanic(args...)
}

func DPanicf(template string, args ...interface{}) {
	global.DPanicf(template, args...)
}

func DPanicln(args ...interface{}) {
	global.DPanicln(args...)
}

func DPanicw(msg string, keysAndValues ...interface{}) {
	global.DPanicw(msg, keysAndValues...)
}

func Panic(args ...interface{}) {
	global.Panic(args...)
}

func Panicf(template string, args ...interface{}) {
	global.Panicf(template, args...)
}

func Panicln(args ...interface{}) {
	global.Panicln(args...)
}

func Panicw(msg string, keysAndValues ...interface{}) {
	global.Panicw(msg, keysAndValues...)
}

func Fatal(args ...interface{}) {
	global.Fatal(args...)
}

func Fatalf(template string, args ...interface{}) {
	global.Fatalf(template, args...)
}

func Fatalln(args ...interface{}) {
	global.Fatalln(args...)
}

func Fatalw(msg string, keysAndValues ...interface{}) {
	global.Fatalw(msg, keysAndValues...)
}

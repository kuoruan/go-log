package log

var global = New()

func SetOptions(opts ...Option) {
	global = global.WithOptions(opts...)
}

func Debug(args ...interface{}) {
	global.debug(global.logger, args...)
}

func Debugf(format string, args ...interface{}) {
	global.debugf(global.logger, format, args...)
}

func Debugln(args ...interface{}) {
	global.debug(global.logger, sprintln(args...))
}

func Debugw(msg string, keysAndValues ...interface{}) {
	global.debugw(global.logger, msg, keysAndValues...)
}

func Info(args ...interface{}) {
	global.info(global.logger, args...)
}

func Infof(format string, args ...interface{}) {
	global.infof(global.logger, format, args...)
}

func Infoln(args ...interface{}) {
	global.info(global.logger, sprintln(args...))
}

func Infow(msg string, keysAndValues ...interface{}) {
	global.infow(global.logger, msg, keysAndValues...)
}

func Warn(args ...interface{}) {
	global.warn(global.logger, args...)
}

func Warnf(template string, args ...interface{}) {
	global.warnf(global.logger, template, args...)
}

func Warnln(args ...interface{}) {
	global.warn(global.logger, sprintln(args...))
}

func Warnw(msg string, keysAndValues ...interface{}) {
	global.warnw(global.logger, msg, keysAndValues...)
}

func Error(args ...interface{}) {
	global.error(global.logger, args...)
}

func Errorf(template string, args ...interface{}) {
	global.errorf(global.logger, template, args...)
}

func Errorln(args ...interface{}) {
	global.error(global.logger, sprintln(args...))
}

func Errorw(msg string, keysAndValues ...interface{}) {
	global.errorw(global.logger, msg, keysAndValues...)
}

func DPanic(args ...interface{}) {
	global.dpanic(global.logger, args...)
}

func DPanicf(template string, args ...interface{}) {
	global.dpanicf(global.logger, template, args...)
}

func DPanicln(args ...interface{}) {
	global.dpanic(global.logger, sprintln(args...))
}

func DPanicw(msg string, keysAndValues ...interface{}) {
	global.dpanicw(global.logger, msg, keysAndValues...)
}

func Panic(args ...interface{}) {
	global.panic(global.logger, args...)
}

func Panicf(template string, args ...interface{}) {
	global.panicf(global.logger, template, args...)
}

func Panicln(args ...interface{}) {
	global.panic(global.logger, sprintln(args...))
}

func Panicw(msg string, keysAndValues ...interface{}) {
	global.panicw(global.logger, msg, keysAndValues...)
}

func Fatal(args ...interface{}) {
	global.fatal(global.logger, args...)
}

func Fatalf(template string, args ...interface{}) {
	global.fatalf(global.logger, template, args...)
}

func Fatalln(args ...interface{}) {
	global.fatal(global.logger, sprintln(args...))
}

func Fatalw(msg string, keysAndValues ...interface{}) {
	global.fatalw(global.logger, msg, keysAndValues...)
}

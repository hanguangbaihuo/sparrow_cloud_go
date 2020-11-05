package utils

import "github.com/kataras/iris/v12/context"

// before run app, set app logger level, iris default logger level is info
// app.Logger().SetLevel("debug")
// Available level names are:
// "disable"
// "fatal"
// "error"
// "warn"
// "info"
// "debug"

// LogDebug will print when logger's Level is debug.
func LogDebug(ctx context.Context, args ...interface{}) {
	ctx.Application().Logger().Debug(args...)
}

// LogDebugf will print when logger's Level is debug.
func LogDebugf(ctx context.Context, format string, args ...interface{}) {
	ctx.Application().Logger().Debugf(format, args...)
}

// LogInfo will print when logger's Level is info or debug.
func LogInfo(ctx context.Context, args ...interface{}) {
	ctx.Application().Logger().Info(args...)
}

// LogInfof will print when logger's Level is info or debug.
func LogInfof(ctx context.Context, format string, args ...interface{}) {
	ctx.Application().Logger().Infof(format, args...)
}

// LogWarn will print when logger's Level is warn, info or debug.
func LogWarn(ctx context.Context, args ...interface{}) {
	ctx.Application().Logger().Warn(args...)
}

// LogWarnf will print when logger's Level is warn, info or debug.
func LogWarnf(ctx context.Context, format string, args ...interface{}) {
	ctx.Application().Logger().Warnf(format, args...)
}

// LogError will print only when logger's Level is error, warn, info or debug.
func LogError(ctx context.Context, args ...interface{}) {
	ctx.Application().Logger().Error(args...)
}

// LogErrorf will print only when logger's Level is error, warn, info or debug.
func LogErrorf(ctx context.Context, format string, args ...interface{}) {
	ctx.Application().Logger().Errorf(format, args...)
}

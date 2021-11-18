/*
Copyright 2021 The tKeel Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
	http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package log

import (
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Conf struct {
	App    string
	Level  string
	Dev    bool
	Output []string
}

var resetGlobalFunc func()

func ResetGlobalLogger() {
	if resetGlobalFunc != nil {
		resetGlobalFunc()
		resetGlobalFunc = nil
	}
}

// InitLogger create new zap logger and sugared logger.
// replace global logger.
func InitLogger(app string, level string, dev bool, output ...string) error {
	c := zap.NewProductionConfig()
	c.Development = dev
	if dev {
		c.Encoding = "console"
	}
	c.Level = getLevel(level)
	if c.InitialFields == nil {
		c.InitialFields = make(map[string]interface{})
	}
	customTimeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(time.RFC3339))
	}
	c.EncoderConfig.TimeKey = "t"
	c.EncoderConfig.EncodeTime = customTimeEncoder
	c.InitialFields["app"] = app
	c.OutputPaths = append(c.OutputPaths, output...)
	logger, err := c.Build(zap.AddCallerSkip(1),
		zap.AddStacktrace(zap.DPanicLevel))
	if err != nil {
		return fmt.Errorf("error build zap log: %w", err)
	}
	resetGlobalFunc = zap.ReplaceGlobals(logger)
	return nil
}

func InitLoggerByConf(c *Conf) error {
	if err := InitLogger(c.App, c.Level, c.Dev, c.Output...); err != nil {
		return fmt.Errorf("error init logger: %w", err)
	}
	return nil
}

func getLevel(level string) zap.AtomicLevel {
	switch strings.ToLower(level) {
	case "debug":
		return zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "Info":
		return zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "warn":
		return zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "error":
		return zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	case "dpanic":
		return zap.NewAtomicLevelAt(zapcore.DPanicLevel)
	case "panic":
		return zap.NewAtomicLevelAt(zapcore.PanicLevel)
	case "fatal":
		return zap.NewAtomicLevelAt(zapcore.FatalLevel)
	}
	return zap.NewAtomicLevelAt(zapcore.InfoLevel)
}

func Debug(args ...interface{}) {
	zap.S().Debug(args...)
}

func Debugf(templateStr string, args ...interface{}) {
	zap.S().Debugf(templateStr, args...)
}

func Debugw(msg string, keysAndValues ...interface{}) {
	zap.S().Debugw(msg, keysAndValues...)
}

func Info(args ...interface{}) {
	zap.S().Info(args...)
}

func Infof(templateStr string, args ...interface{}) {
	zap.S().Infof(templateStr, args...)
}

func Infow(msg string, keysAndValues ...interface{}) {
	zap.S().Infow(msg, keysAndValues...)
}

func Warn(args ...interface{}) {
	zap.S().Warn(args...)
}

func Warnf(templateStr string, args ...interface{}) {
	zap.S().Warnf(templateStr, args...)
}

func Warnw(msg string, keysAndValues ...interface{}) {
	zap.S().Warnw(msg, keysAndValues...)
}

func Error(args ...interface{}) {
	zap.S().Error(args...)
}

func Errorf(templateStr string, args ...interface{}) {
	zap.S().Errorf(templateStr, args...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	zap.S().Errorw(msg, keysAndValues...)
}

func DPanic(args ...interface{}) {
	zap.S().DPanic(args...)
}

func DPanicf(templateStr string, args ...interface{}) {
	zap.S().DPanicf(templateStr, args...)
}

func DPanicw(msg string, keysAndValues ...interface{}) {
	zap.S().DPanicw(msg, keysAndValues...)
}

func Panic(args ...interface{}) {
	zap.S().Panic(args...)
}

func Panicf(templateStr string, args ...interface{}) {
	zap.S().Panicf(templateStr, args...)
}

func Panicw(msg string, keysAndValues ...interface{}) {
	zap.S().Panicw(msg, keysAndValues...)
}

func Fatal(args ...interface{}) {
	zap.S().Fatal(args...)
}

func Fatalf(templateStr string, args ...interface{}) {
	zap.S().Fatalf(templateStr, args...)
}

func Fatalw(msg string, keysAndValues ...interface{}) {
	zap.S().Fatalw(msg, keysAndValues...)
}

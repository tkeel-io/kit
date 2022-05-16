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

// zap.S() func
var (
	Debug   func(args ...interface{})
	Debugf  func(templateStr string, args ...interface{})
	Debugw  func(msg string, keysAndValues ...interface{})
	Info    func(args ...interface{})
	Infof   func(templateStr string, args ...interface{})
	Infow   func(msg string, keysAndValues ...interface{})
	Warn    func(args ...interface{})
	Warnf   func(templateStr string, args ...interface{})
	Warnw   func(msg string, keysAndValues ...interface{})
	Error   func(args ...interface{})
	Errorf  func(templateStr string, args ...interface{})
	Errorw  func(msg string, keysAndValues ...interface{})
	DPanic  func(args ...interface{})
	DPanicf func(templateStr string, args ...interface{})
	DPanicw func(msg string, keysAndValues ...interface{})
	Panic   func(args ...interface{})
	Panicf  func(templateStr string, args ...interface{})
	Panicw  func(msg string, keysAndValues ...interface{})
	Fatal   func(args ...interface{})
	Fatalf  func(templateStr string, args ...interface{})
	Fatalw  func(msg string, keysAndValues ...interface{})
)

func init() {
	globalLogger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(globalLogger)
	syncLogS()
}

// syncLogS global logger sync
func syncLogS() {
	globalS := zap.S()
	Debug = globalS.Debug
	Debugf = globalS.Debugf
	Debugw = globalS.Debugw
	Info = globalS.Info
	Infof = globalS.Infof
	Infow = globalS.Infow
	Warn = globalS.Warn
	Warnf = globalS.Warnf
	Warnw = globalS.Warnw
	Error = globalS.Error
	Errorf = globalS.Errorf
	Errorw = globalS.Errorw
	DPanic = globalS.DPanic
	DPanicf = globalS.DPanicf
	DPanicw = globalS.DPanicw
	Panic = globalS.Panic
	Panicf = globalS.Panicf
	Panicw = globalS.Panicw
	Fatal = globalS.Fatal
	Fatalf = globalS.Fatalf
	Fatalw = globalS.Fatalw
}

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
	logger, err := c.Build(zap.AddCallerSkip(0),
		zap.AddStacktrace(zap.DPanicLevel))
	if err != nil {
		return fmt.Errorf("error build zap log: %w", err)
	}
	resetGlobalFunc = zap.ReplaceGlobals(logger)
	syncLogS()
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

func Check(lvl zapcore.Level, msg string) *zapcore.CheckedEntry {
	return zap.L().Check(lvl, msg)
}

func L() *zap.Logger {
	return zap.L()
}

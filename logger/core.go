package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"gopkg.in/yaml.v2"
)

// Config struct
type Config struct {
	Syncer     Syncer `json:"syncer" yaml:"-"`
	Level      string `json:"level" yaml:"-"`
	FileName   string `json:"file_name" yaml:"filename"`
	MaxSize    int    `json:"max_size" yaml:"maxsize"`
	MaxAge     int    `json:"max_age" yaml:"maxage"`
	MaxBackups int    `json:"max_backups" yaml:"maxbackups"`
}

// NewCoreFile function
func NewCoreFile(cfg Config) zapcore.Core {
	b, _ := yaml.Marshal(cfg)
	var cfgRotate lumberjack.Logger
	if err := yaml.Unmarshal(b, &cfgRotate); err != nil {
		// log err
	}
	lvl := zap.NewAtomicLevel()
	lvl.UnmarshalText([]byte(cfg.Level))
	encoder := zapcore.EncoderConfig{
		TimeKey:       "timestamp",
		LevelKey:      "level",
		NameKey:       "name",
		CallerKey:     "caller",
		MessageKey:    "message",
		StacktraceKey: "stacktrace",
		EncodeTime:    TimeEncoder,
		EncodeLevel:   zapcore.CapitalLevelEncoder,
		EncodeName:    zapcore.FullNameEncoder,
	}
	if lvl.Enabled(zapcore.DebugLevel) || lvl.Enabled(zapcore.ErrorLevel) {
		encoder.EncodeCaller = zapcore.ShortCallerEncoder
	}
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoder),
		zapcore.AddSync(&cfgRotate),
		lvl,
	)

	return core
}

// NewCoreStd function
func NewCoreStd(cfg Config) zapcore.Core {
	lvl := zap.NewAtomicLevel()
	if "" != cfg.Level {
		lvl.UnmarshalText([]byte(cfg.Level))
	}
	encoder := zapcore.EncoderConfig{
		TimeKey:       "timestamp",
		LevelKey:      "level",
		NameKey:       "name",
		CallerKey:     "caller",
		MessageKey:    "message",
		StacktraceKey: "stacktrace",
		EncodeTime:    TimeEncoder,
		EncodeLevel:   zapcore.CapitalLevelEncoder,
		EncodeName:    zapcore.FullNameEncoder,
	}
	if lvl.Enabled(zapcore.DebugLevel) || lvl.Enabled(zapcore.ErrorLevel) {
		encoder.EncodeCaller = zapcore.ShortCallerEncoder
	}
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoder),
		zapcore.AddSync(os.Stdout),
		lvl,
	)

	return core
}

package logger

import (
	"reflect"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/micro/go-micro/v2/config/reader"
)

var logger *zap.SugaredLogger

// TimeEncoder serializes a time.Time to an RFC3339-formatted string without timezone.
func TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

// PanicHandler logs unexpected panic then throw it
func PanicHandler(logger *zap.SugaredLogger) {
	if r := recover(); r != nil {
		logger.Panic(r)
		panic(r)
	}
}

// NewTee function
func NewTee(value reader.Value) (zapcore.Core, error) {
	var configs []Config
	var v interface{}
	value.Scan(&v)
	t := reflect.TypeOf(v)
	switch t.Kind() {
	case reflect.String:
		configs = []Config{{Syncer: Std, Level: v.(string)}}
		break
	case reflect.Slice:
		if err := value.Scan(&configs); err != nil {
			return nil, err
		}
		break
	case reflect.Map:
		var config Config
		if err := value.Scan(&config); err != nil {
			return nil, err
		}
		configs = []Config{config}
		break
	}

	cores := make([]zapcore.Core, len(configs))
	for i := 0; i < len(configs); i++ {
		c := configs[i]
		if c.Syncer.Enabled(Std) {
			cores[i] = NewCoreStd(c)
		} else if c.Syncer.Enabled(File) {
			cores[i] = NewCoreFile(c)
		}
	}

	return zapcore.NewTee(cores...), nil
}

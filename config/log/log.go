package log

import (
	"github.com/micro/go-micro/v2/config/reader"
	log "github.com/micro/go-micro/v2/logger"

	"github.com/bimoyong/go-plugins/logger/zap/v2"
	"github.com/bimoyong/go-util/logger"
)

// Refresh defines how to deal with configuration change
func Refresh(value reader.Value) {
	var v interface{}

	if value.Scan(&v); v == nil {
		log.Error("[Config][Log] Not allow nil!")

		return
	}

	core, err := logger.NewTee(value)
	if err != nil {
		log.Errorf("[Config][Log] Read failure!: err:=[%s] val=[%+v]", err.Error(), value)

		return
	}

	l, err := zap.NewLogger(zap.WithCore(core), zap.WithCallerSkip(2))
	if err != nil {
		log.Errorf("[Config][Log] New logger failure!: err:=[%s] core=[%+v] val=[%+v]", err.Error(), core, value)

		return
	}

	log.DefaultLogger = l
}

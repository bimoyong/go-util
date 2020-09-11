package config

import (
	"reflect"
	"strings"

	"github.com/go-openapi/inflect"
	"github.com/micro/go-micro/v2/config"
	"github.com/micro/go-micro/v2/config/reader"
	"github.com/micro/go-micro/v2/config/source/file"
	log "github.com/micro/go-micro/v2/logger"

	clog "github.com/bimoyong/go-util/config/log"
	"github.com/bimoyong/go-util/str"
)

// Config struct
type Config struct {
}

// Watch struct
type Watch struct {
	Target  interface{}
	Handler reflect.Method
	Value   reader.Value
}

var watchChan = make(chan Watch)

func init() {
	err := config.Load(
		file.NewSource(
			file.WithPath("config.json"),
		),
	)
	if err != nil {
		log.Fatalf("[Config] Read failure!: err=[%s]", err.Error())
	}

	Init(NewConfig())
}

// NewConfig function
func NewConfig() *Config {
	return nil
}

// Init function
func Init(c interface{}) {
	go onRefresh()

	watch(c)
}

func onRefresh() {
	for {
		select {
		case watch := <-watchChan:
			h := watch.Handler
			log.Info("[Config] Load: ", str.Beautify(map[string]interface{}{
				"key": inflect.Underscore(h.Name),
				"val": watch.Value,
			}))
			h.Func.Call([]reflect.Value{
				reflect.ValueOf(watch.Target),
				reflect.ValueOf(watch.Value),
			})
		}
	}
}

func watch(c interface{}) {
	ct := reflect.TypeOf(c)
	for i := 0; i < ct.NumMethod(); i++ {
		m := ct.Method(i)
		keys := strings.Split(m.Name, "_")
		kPath := make([]string, len(keys))
		for i := 0; i < len(keys); i++ {
			kPath[i] = inflect.Underscore(keys[i])
		}

		watcher, err := config.Watch(kPath...)
		if err != nil {
			log.Fatalf("[Config] Watch failure!: err=[%s] key=[%s]", err.Error(), kPath)
		}

		v := config.Get(kPath...)

		go func() {
			for {
				watchChan <- Watch{Target: c, Handler: m, Value: v}
				v, _ = watcher.Next()
			}
		}()
	}
}

// Log function
func (s *Config) Log(value reader.Value) {
	clog.Refresh(value)
}

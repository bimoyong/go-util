package database

import (
	"github.com/micro/go-micro/v2/config/reader"
	log "github.com/micro/go-micro/v2/logger"

	"gitlab.com/bimoyong/go-util/database"
)

var (
	// Connection holds client object
	Connection *database.Client
	// Config stores database configuration
	Config database.Config
)

// Refresh defines how to deal with configuration change
func Refresh(value reader.Value) {
	var err error

	if err = value.Scan(&Config); err != nil {
		log.Errorf("[Config][Database] Read failure!: err:=[%s] val=[%+v]", err.Error(), value)

		return
	}

	if obsoleteConn := Connection; obsoleteConn != nil {
		defer disconnect(obsoleteConn)
	}

	if Connection, err = database.NewConnection(Config); err != nil {
		log.Errorf("[Config][Database] Connect failure!: err:=[%s] val=[%+v]", err.Error(), value)

		return
	}
	log.Info("[Config][Database] Connect success")
}

// Close disconnects as safe way to release server side resource
func Close() error {
	return disconnect(Connection)
}

func disconnect(conn *database.Client) (err error) {
	if err = conn.Disconnect(); err != nil {
		log.Errorf("[Config][Database] Disconnect failure!: err:=[%s]", err.Error())

		return
	}
	log.Info("[Config][Database] Disconnect success")

	return
}

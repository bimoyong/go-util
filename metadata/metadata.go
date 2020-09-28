package model

import (
	"github.com/micro/go-micro/v2/metadata"
)

// NewMetadata function
func NewMetadata(md metadata.Metadata) (rst metadata.Metadata) {
	rst = make(metadata.Metadata)

	var keys = []string{
		"Alias",
		"Authorization",
		"Domain",
		"ID",
		"Postback",
		"Resource",
	}

	for _, k := range keys {
		if v, ok := md.Get(k); ok {
			rst.Set(k, v)
		}
	}

	return
}

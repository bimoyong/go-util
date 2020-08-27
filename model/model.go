package model

import (
	"errors"
	"fmt"
	"strings"
)

type (
	// Header struct
	Header struct {
		ID              string `json:"id,omitempty" mapstructure:"id,omitempty"`
		Alias           string `json:"alias,omitempty" mapstructure:"alias,omitempty"`
		Resource        string `json:"resource,omitempty" mapstructure:"resource,omitempty"`
		Domain          string `json:"domain,omitempty" mapstructure:"domain,omitempty"`
		PostbackChannel string `json:"postback_channel,omitempty" mapstructure:"postback_channel,omitempty"`
	}
)

// PopPostbackChannel function
func (h *Header) PopPostbackChannel() string {
	t := h.PostbackChannel
	h.PostbackChannel = ""

	return t
}

// ParseServerName function
func ParseServerName(serverName string) (*Header, error) {
	srvName := serverName
	parts := strings.Split(srvName, ".")
	if 1 >= len(parts) {
		return nil, errors.New(fmt.Sprintf("[%s] is compliant to server name pattern: [{NAMESPACE}.{TYPE}.{ALIAS}-{DOMAIN}]", srvName))
	}

	ad := parts[2]

	sepIndex := strings.IndexByte(ad, '-')
	var h *Header
	if 0 > sepIndex {
		h = &Header{
			Alias: ad,
		}
	} else {
		h = &Header{
			Alias:  ad[:sepIndex],
			Domain: ad[sepIndex+1:],
		}
	}

	return h, nil
}

func (s *Header) String() string {
	str := s.Alias
	if 0 < len(s.Domain) {
		str += "-" + s.Domain
	}
	if 0 < len(s.Resource) {
		str += "." + s.Resource
	}

	return str
}

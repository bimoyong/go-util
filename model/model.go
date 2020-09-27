package model

import (
	"fmt"
	"strings"
)

type (
	// Header struct
	Header struct {
		ID       string `json:"id,omitempty" mapstructure:"id,omitempty"`
		Alias    string `json:"alias,omitempty" mapstructure:"alias,omitempty"`
		Resource string `json:"resource,omitempty" mapstructure:"resource,omitempty"`
		Domain   string `json:"domain,omitempty" mapstructure:"domain,omitempty"`
		Postback string `json:"postback,omitempty" mapstructure:"postback,omitempty"`
	}
)

// PopPostback function
func (h *Header) PopPostback() string {
	t := h.Postback
	h.Postback = ""

	return t
}

// ParseServerName function
func ParseServerName(serverName string) (h *Header, err error) {
	parts := strings.Split(serverName, ".")
	if 1 >= len(parts) {
		err = fmt.Errorf("[%s] is compliant to server name pattern: [{NAMESPACE}.{TYPE}.{ALIAS}-{DOMAIN}]", serverName)
		return
	}

	ad := parts[2]

	sepIndex := strings.IndexByte(ad, '-')
	if 0 > sepIndex {
		h = &Header{
			Alias: ad,
		}

		return
	}

	h = &Header{
		Alias:  ad[:sepIndex],
		Domain: ad[sepIndex+1:],
	}

	return
}

// String function
func (h *Header) String() (str string) {
	str = h.Alias
	if 0 < len(h.Domain) {
		str += "-" + h.Domain
	}
	if 0 < len(h.Resource) {
		str += "." + h.Resource
	}

	return
}

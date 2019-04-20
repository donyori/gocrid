package gocrid

import (
	"encoding/json"
	"io/ioutil"
	"time"
)

type Settings struct {
	IdPrefix          string        `json:"id_prefix,omitempty"`
	IdBodyLength      int           `json:"id_body_length,omitempty"`
	MaxAge            time.Duration `json:"max_age,omitempty"`
	EnableHostBinding bool          `json:"enable_host_binding,omitempty"`

	CookieName     string `json:"cookie_name,omitempty"`
	CookiePath     string `json:"cookie_path,omitempty"`
	CookieDomain   string `json:"cookie_domain,omitempty"`
	CookieSecure   bool   `json:"cookie_secure,omitempty"`
	CookieHttpOnly bool   `json:"cookie_http_only,omitempty"`
}

func NewSettings() *Settings {
	return &Settings{
		IdPrefix:          "",
		IdBodyLength:      16,
		MaxAge:            time.Duration(time.Minute * 30),
		EnableHostBinding: true,
		CookieName:        "gocrid",
		CookiePath:        "/",
		CookieDomain:      "",
		CookieSecure:      false,
		CookieHttpOnly:    false,
	}
}

func LoadSettings(filename string) (settings *Settings, err error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	settings = NewSettings()
	err = json.Unmarshal(data, settings)
	if err != nil {
		return nil, err
	}
	return settings, nil
}

func (s *Settings) FillEmptyFields() {
	if s == nil {
		return
	}
	if s.IdBodyLength <= 0 {
		s.IdBodyLength = 16
	}
	if s.CookieName == "" {
		s.CookieName = "gocrid"
		if s.IdPrefix != "" {
			s.CookieName += "__" + s.IdPrefix
		}
	}
	if s.CookiePath == "" {
		s.CookiePath = "/"
	}
}

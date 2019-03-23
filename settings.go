package gocrid

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"
)

type Settings struct {
	IdPrefix          string        `json:"id_prefix"`
	IdBodyLength      int           `json:"id_body_length"`
	MaxAge            time.Duration `json:"max_age"`
	EnableHostBinding bool          `json:"enable_host_binding"`

	CookieName     string `json:"cookie_name"`
	CookiePath     string `json:"cookie_path"`
	CookieDomain   string `json:"cookie_domain"`
	CookieSecure   bool   `json:"cookie_secure"`
	CookieHttpOnly bool   `json:"cookie_http_only"`
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
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close() // Ignore error.
	data, err := ioutil.ReadAll(file)
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

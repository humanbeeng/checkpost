package config

import (
	"log/slog"
	"os"
	"strings"

	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

const CheckpostConfigPrefix = "CP_"

type AppConfig struct {
	Postgres `koanf:"postgres"`
	Github   `koanf:"github"`
	Google   `koanf:"google"`
	Paseto   `koanf:"paseto"`
}

type Postgres struct {
	User     string `koanf:"user"`
	Password string `koanf:"password"`
	Host     string `koanf:"host"`
	Port     int    `koanf:"port"`
	Database string `koanf:"database"`
}

type Github struct {
	ClientId string `koanf:"client_id"`
	Secret   string `koanf:"secret"`
}

type Google struct {
	ClientId    string `koanf:"client_id"`
	Secret      string `koanf:"secret"`
	RedirectUrl string `koanf:"redirect_url"`
}

type Paseto struct {
	Key string `koanf:"key"`
}

func GetAppConfig() (*AppConfig, error) {
	k := koanf.New(".")
	if err := k.Load(file.Provider("config.toml"), toml.Parser()); os.IsNotExist(err) {
		if err := k.Load(env.Provider(CheckpostConfigPrefix, ".", func(s string) string {
			str := strings.Replace(strings.ToLower(
				strings.TrimPrefix(s, CheckpostConfigPrefix)), "_", ".", -1)
			return str
		}), nil); err != nil {
			slog.Error("unable to load config", "err", err)
			return nil, err
		}
	}

	var appConfig AppConfig

	err := k.Unmarshal("", &appConfig)
	if err != nil {
		slog.Error("unable to unmarshal config", "err", err)
		return nil, err
	}
	return &appConfig, nil
}

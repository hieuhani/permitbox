package config

import (
	"errors"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/fs"
	"github.com/knadh/koanf/v2"
	"gitlab.com/hieuhani/permitbox/pkg/strutil"
	iofs "io/fs"
	"strings"
)

const (
	EnvPrefix = "APP__"
)

type DbConfig struct {
	Host            string `koanf:"host"`
	Port            int    `koanf:"port"`
	User            string `koanf:"user"`
	Password        string `koanf:"password"`
	DbName          string `koanf:"dbName"`
	SqlDebugEnabled bool   `koanf:"sqlDebugEnabled"`
}

type AppConfig struct {
	HttpPort int      `koanf:"httpPort"`
	Db       DbConfig `koanf:"db"`
}

func InitConfig[T any](configFile iofs.FS) (T, error) {
	var config T
	k := koanf.New(".")
	configProvider := fs.Provider(configFile, "config.yaml")
	if err := k.Load(configProvider, yaml.Parser()); err != nil {
		return config, errors.New("cannot read config from file")
	}
	if err := k.Load(
		env.ProviderWithValue(
			EnvPrefix, ".", func(key string, value string) (string, any) {
				newKey := strutil.SnakeToCamel(
					strings.Replace(
						strings.ToLower(
							strings.TrimPrefix(key, EnvPrefix),
						), "__", ".", -1,
					),
				)
				if strings.Contains(value, ",") {
					return newKey, strings.Split(value, ",")
				}
				return newKey, value
			},
		), nil,
	); err != nil {
		return config, err
	}

	if err := k.Unmarshal("", &config); err != nil {
		return config, err
	}
	return config, nil
}

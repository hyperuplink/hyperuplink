package config

import (
	"net/url"

	"github.com/knadh/koanf/parsers/toml/v2"
	"github.com/knadh/koanf/providers/file"
	"github.com/mrusme/hyperuplink/errs"

	"github.com/knadh/koanf/v2"
)

type Config struct {
	cfgstr   string
	k        *koanf.Koanf
	provider koanf.Provider
}

func New(
	cfgstr string,
) (cfg *Config, err error) {
	cfg = new(Config)
	var cfgurl *url.URL

	cfg.cfgstr = cfgstr
	if cfgurl, err = cfg.parseCfgstr(); err != nil {
		return nil, err
	}

	cfg.k = koanf.New(".")

	switch cfgurl.Scheme {
	case "file":
		cfg.provider = file.Provider(cfgurl.Path)
	default:
		return nil, errs.ErrConfigTypeUnsupported
	}

	if err = cfg.k.Load(cfg.provider, toml.Parser()); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (cfg *Config) Startup() (err error) {
	cfg.provider.(*file.File).Watch(func(event interface{}, werr error) {
		if werr != nil {
			// TODO: Handle error
			return
		}

		cfg.k = koanf.New(".")
		cfg.k.Load(cfg.provider, toml.Parser())
	})

	return nil
}

func (cfg *Config) Shutdown() error {
	cfg.provider.(*file.File).Unwatch()
	return nil
}

func (cfg *Config) parseCfgstr() (u *url.URL, err error) {
	if u, err = url.Parse(cfg.cfgstr); err != nil {
		return nil, err
	}

	return u, nil
}

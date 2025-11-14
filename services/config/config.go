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
) (*Config, error) {
	var cfg *Config = new(Config)
	var err error
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

func (cfg *Config) Startup() error {
	// var err error

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

func (cfg *Config) parseCfgstr() (*url.URL, error) {
	var err error
	var u *url.URL

	if u, err = url.Parse(cfg.cfgstr); err != nil {
		return nil, err
	}

	return u, nil
}

func (cfg *Config) LoggingLevel() []byte {
	return cfg.k.Bytes("Logging.Level")
}

func (cfg *Config) DatabaseConnection() string {
	return cfg.k.String("Database.Connection")
}

func (cfg *Config) ServerBindIP() string {
	return cfg.k.String("Server.BindIP")
}

func (cfg *Config) ServerPort() int {
	return cfg.k.Int("Server.Port")
}

func (cfg *Config) ServerBodyLimit() int {
	return cfg.k.Int("Server.BodyLimit")
}

func (cfg *Config) ServerConcurrency() int {
	return cfg.k.Int("Server.Concurrency")
}

func (cfg *Config) ServerProxyHeader() string {
	return cfg.k.String("Server.ProxyHeader")
}

func (cfg *Config) ServerEnableIPValidation() bool {
	return cfg.k.Bool("Server.EnableIPValidation")
}

func (cfg *Config) ServerTrustProxy() bool {
	return cfg.k.Bool("Server.TrustProxy")
}

func (cfg *Config) ServerTrustLoopback() bool {
	return cfg.k.Bool("Server.TrustLoopback")
}

func (cfg *Config) ServerTrustProxies() []string {
	return cfg.k.Strings("Server.TrustProxies")
}

func (cfg *Config) ServerReduceMemoryUsage() bool {
	return cfg.k.Bool("Server.ReduceMemoryUsage")
}

func (cfg *Config) ServerServerHeader() string {
	return cfg.k.String("Server.ServerHeader")
}

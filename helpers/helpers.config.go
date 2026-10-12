package helpers

import (
	"xn--gckvb8fzb.com/glides/services/config"
)

const (
	APIDefaultPort           int = 3001
	APIDefaultBodyLimit      int = 4 * 1024 * 1024
	APIDefaultConcurrency    int = 256 * 1024
	APIDefaultRateLimitGuest int = 120
	APIDefaultRateLimitUser  int = 600
	APIDefaultRateLimitAuth  int = 10
	APIDefaultRequestTimeout int = 30
	WebDefaultRequestTimeout int = 30
)

type APIConfig struct {
	Enable             bool     `koanf:"Enable"`
	BindIP             string   `koanf:"BindIP"`
	Port               int      `koanf:"Port"`
	BodyLimit          int      `koanf:"BodyLimit"`
	Concurrency        int      `koanf:"Concurrency"`
	ProxyHeader        string   `koanf:"ProxyHeader"`
	EnableIPValidation bool     `koanf:"EnableIPValidation"`
	TrustProxy         bool     `koanf:"TrustProxy"`
	TrustLoopback      bool     `koanf:"TrustLoopback"`
	TrustProxies       []string `koanf:"TrustProxies"`
	ReduceMemoryUsage  bool     `koanf:"ReduceMemoryUsage"`
	ServerHeader       string   `koanf:"ServerHeader"`
	RateLimitGuest     int      `koanf:"RateLimitGuest"`
	RateLimitUser      int      `koanf:"RateLimitUser"`
	RateLimitAuth      int      `koanf:"RateLimitAuth"`
	CORSOrigins        []string `koanf:"CORSOrigins"`
	RequestTimeout     int      `koanf:"RequestTimeout"`
}

type ServerConfig struct {
	Enable             bool     `koanf:"Enable"`
	BindIP             string   `koanf:"BindIP"`
	Port               int      `koanf:"Port"`
	BodyLimit          int      `koanf:"BodyLimit"`
	Concurrency        int      `koanf:"Concurrency"`
	ProxyHeader        string   `koanf:"ProxyHeader"`
	EnableIPValidation bool     `koanf:"EnableIPValidation"`
	TrustProxy         bool     `koanf:"TrustProxy"`
	TrustLoopback      bool     `koanf:"TrustLoopback"`
	TrustProxies       []string `koanf:"TrustProxies"`
	ReduceMemoryUsage  bool     `koanf:"ReduceMemoryUsage"`
	ServerHeader       string   `koanf:"ServerHeader"`
	RequestTimeout     int      `koanf:"RequestTimeout"`
}

func NewAPIConfig() APIConfig {
	return APIConfig{
		Port:           APIDefaultPort,
		BodyLimit:      APIDefaultBodyLimit,
		Concurrency:    APIDefaultConcurrency,
		RateLimitGuest: APIDefaultRateLimitGuest,
		RateLimitUser:  APIDefaultRateLimitUser,
		RateLimitAuth:  APIDefaultRateLimitAuth,
		RequestTimeout: APIDefaultRequestTimeout,
	}
}

func APIConfigFrom(cfg *config.Config) (api APIConfig, err error) {
	api = NewAPIConfig()

	if err = cfg.Unmarshal("API", &api); err != nil {
		return NewAPIConfig(), err
	}

	return api, nil
}

func WebConfigFrom(cfg *config.Config) (web ServerConfig, err error) {
	web = ServerConfig{RequestTimeout: WebDefaultRequestTimeout}

	if err = cfg.Unmarshal("Web", &web); err != nil {
		return ServerConfig{RequestTimeout: WebDefaultRequestTimeout}, err
	}

	return web, nil
}

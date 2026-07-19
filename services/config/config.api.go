package config

const APIDefaultPort = 3001

func (cfg *Config) APIBindIP() string {
	return cfg.k.String("API.BindIP")
}

func (cfg *Config) APIPort() int {
	port := cfg.k.Int("API.Port")
	if port == 0 {
		return APIDefaultPort
	}
	return port
}

func (cfg *Config) APIBodyLimit() int {
	return cfg.k.Int("API.BodyLimit")
}

func (cfg *Config) APIConcurrency() int {
	return cfg.k.Int("API.Concurrency")
}

func (cfg *Config) APIProxyHeader() string {
	return cfg.k.String("API.ProxyHeader")
}

func (cfg *Config) APIEnableIPValidation() bool {
	return cfg.k.Bool("API.EnableIPValidation")
}

func (cfg *Config) APITrustProxy() bool {
	return cfg.k.Bool("API.TrustProxy")
}

func (cfg *Config) APITrustLoopback() bool {
	return cfg.k.Bool("API.TrustLoopback")
}

func (cfg *Config) APITrustProxies() []string {
	return cfg.k.Strings("API.TrustProxies")
}

func (cfg *Config) APIReduceMemoryUsage() bool {
	return cfg.k.Bool("API.ReduceMemoryUsage")
}

func (cfg *Config) APIServerHeader() string {
	return cfg.k.String("API.ServerHeader")
}

package config

func (cfg *Config) WebEnable() bool {
	return cfg.k.Bool("Web.Enable")
}

func (cfg *Config) WebBindIP() string {
	return cfg.k.String("Web.BindIP")
}

func (cfg *Config) WebPort() int {
	return cfg.k.Int("Web.Port")
}

func (cfg *Config) WebBodyLimit() int {
	return cfg.k.Int("Web.BodyLimit")
}

func (cfg *Config) WebConcurrency() int {
	return cfg.k.Int("Web.Concurrency")
}

func (cfg *Config) WebProxyHeader() string {
	return cfg.k.String("Web.ProxyHeader")
}

func (cfg *Config) WebEnableIPValidation() bool {
	return cfg.k.Bool("Web.EnableIPValidation")
}

func (cfg *Config) WebTrustProxy() bool {
	return cfg.k.Bool("Web.TrustProxy")
}

func (cfg *Config) WebTrustLoopback() bool {
	return cfg.k.Bool("Web.TrustLoopback")
}

func (cfg *Config) WebTrustProxies() []string {
	return cfg.k.Strings("Web.TrustProxies")
}

func (cfg *Config) WebReduceMemoryUsage() bool {
	return cfg.k.Bool("Web.ReduceMemoryUsage")
}

func (cfg *Config) WebServerHeader() string {
	return cfg.k.String("Web.ServerHeader")
}

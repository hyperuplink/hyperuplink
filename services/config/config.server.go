package config

func (cfg *Config) ServerEnable() bool {
	return cfg.k.Bool("Server.Enable")
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

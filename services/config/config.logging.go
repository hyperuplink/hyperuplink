package config

func (cfg *Config) LoggingLevel() []byte {
	return cfg.k.Bytes("Logging.Level")
}

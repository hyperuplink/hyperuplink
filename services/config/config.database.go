package config

func (cfg *Config) DatabaseConnection() string {
	return cfg.k.String("Database.Connection")
}

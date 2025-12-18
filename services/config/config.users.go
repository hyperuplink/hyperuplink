package config

func (cfg *Config) UsersPromoteAdmin() []string {
	return cfg.k.Strings("Users.PromoteAdmin")
}

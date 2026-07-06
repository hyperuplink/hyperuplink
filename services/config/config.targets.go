package config

type Target struct {
	ID    string      `koanf:"ID"`
	Type  string      `koanf:"Type"`
	Email TargetEmail `koanf:"Email,omitempty"`
	XMPP  TargetXMPP  `koanf:"XMPP,omitempty"`
}

type TargetEmail struct {
	SMTPServer    string
	SMTPAuthType  string
	SMTPTLSPolicy int
	SMTPUsername  string
	SMTPPassword  string
	From          struct {
		Email string
		Name  string
	}
}

type TargetXMPP struct {
	Server   string
	TLS      bool
	Username string
	Password string
}

type Targets []Target

func (cfg *Config) Targets() (targets Targets, err error) {
	err = cfg.k.Unmarshal("Target", &targets)
	return targets, err
}

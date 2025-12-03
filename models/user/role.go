package user

type Role string

const (
	GuestRole Role = "guest"
	UserRole  Role = "user"
	ModRole   Role = "mod"
	AdminRole Role = "admin"
)

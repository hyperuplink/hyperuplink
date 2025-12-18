package user

type Role string

const (
	GuestRole Role = "guest"
	UserRole  Role = "user"
	AdminRole Role = "admin"
)

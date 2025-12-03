package site

import "github.com/mrusme/hyperuplink/models/user"

func (s *Site) SetCurrentUser(u *user.User) {
	s.currentUser = u
}

func (s *Site) GetCurrentUserRole() string {
	if s.currentUser != nil {
		return s.currentUser.Role
	}
	return "guest"
}

func (s *Site) CurrentUserHasGuestRole() bool {
	if s.currentUser == nil {
		return true
	}
	return false
}

func (s *Site) CurrentUserHasUserRole() bool {
	if s.currentUser != nil && s.currentUser.Role == "user" {
		return true
	}
	return false
}

func (s *Site) CurrentUserHasAdminRole() bool {
	if s.currentUser != nil && s.currentUser.Role == "admin" {
		return true
	}
	return false
}

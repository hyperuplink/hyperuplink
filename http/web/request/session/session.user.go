package session

import "xn--gckvb8fzb.com/hyperuplink/models/user"

func (s *Session) SetCurrentUser(u *user.User) {
	s.currentUser = u
}

func (s *Session) GetCurrentUserUsername() string {
	if s.currentUser != nil {
		return s.currentUser.Username
	}
	return ""
}

func (s *Session) GetCurrentUserID() string {
	if s.currentUser != nil {
		return s.currentUser.ID.String()
	}
	return ""
}

func (s *Session) GetCurrentUserRole() user.Role {
	if s.currentUser != nil {
		return s.currentUser.Role
	}
	return user.GuestRole
}

func (s *Session) GetCurrentUserMemberOf() []string {
	if s.currentUser != nil {
		return s.currentUser.MemberOf
	}
	return nil
}

func (s *Session) CurrentUserHasGuestRole() bool {
	if s.currentUser == nil {
		return true
	}
	return false
}

func (s *Session) CurrentUserHasUserRole() bool {
	if s.currentUser != nil && s.currentUser.Role == user.UserRole {
		return true
	}
	return false
}

func (s *Session) CurrentUserHasAdminRole() bool {
	if s.currentUser != nil && s.currentUser.Role == user.AdminRole {
		return true
	}
	return false
}

func (s *Session) GetCurrentUserTimezone() string {
	if s.currentUser != nil {
		return s.currentUser.Timezone
	}
	return "UTC"
}

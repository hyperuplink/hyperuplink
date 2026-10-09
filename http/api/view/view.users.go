package view

import (
	"xn--gckvb8fzb.com/glides/http/route"
	"xn--gckvb8fzb.com/hyperuplink/models/apikey"
	"xn--gckvb8fzb.com/hyperuplink/models/group"
	"xn--gckvb8fzb.com/hyperuplink/models/permission"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

type Group struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type User struct {
	ID                string  `json:"id"`
	Username          string  `json:"username"`
	Role              string  `json:"role"`
	ProfilePictureURL string  `json:"profile_picture_url,omitempty"`
	SignatureText     string  `json:"signature_text"`
	SignatureHTML     string  `json:"signature_html"`
	JoinedAt          Time    `json:"joined_at"`
	Groups            []Group `json:"groups"`
	URL               string  `json:"url"`
}

type Account struct {
	User
	Email         string `json:"email"`
	EmailIsJID    bool   `json:"email_is_jid"`
	Language      string `json:"language"`
	Timezone      string `json:"timezone"`
	OTPEnabled    bool   `json:"otp_enabled"`
	NotifyOnReply bool   `json:"notify_on_reply"`
	SignInLastAt  Time   `json:"sign_in_last_at"`
	ConfirmedAt   Time   `json:"confirmed_at"`
}

type OTPEnrollment struct {
	URL       string `json:"url"`
	Secret    string `json:"secret"`
	QRDataURI string `json:"qr_data_uri"`
}

type AdminUser struct {
	Account
	MemberOf  []string `json:"member_of"`
	BannedAt  Time     `json:"banned_at"`
	DeletedAt Time     `json:"deleted_at"`
}

type Key struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	LastUsedAt Time   `json:"last_used_at"`
	CreatedAt  Time   `json:"created_at"`
}

type SessionPermissions struct {
	Admin      bool              `json:"admin"`
	Default    string            `json:"default"`
	Categories map[string]string `json:"categories"`
}

type Session struct {
	Authenticated bool               `json:"authenticated"`
	Role          string             `json:"role"`
	User          *Account           `json:"user,omitempty"`
	Key           *Key               `json:"key,omitempty"`
	Permissions   SessionPermissions `json:"permissions"`
}

func Groups(groups []group.Group) []Group {
	views := make([]Group, 0, len(groups))
	for _, grp := range groups {
		views = append(views, Group{ID: grp.ID, Name: grp.Name})
	}

	return views
}

func (m *Mapper) User(usr *user.User, groups []group.Group) User {
	return User{
		ID:                usr.ID.String(),
		Username:          usr.Username,
		Role:              string(usr.Role),
		ProfilePictureURL: m.profilePictureURL(usr.ProfilePicture),
		SignatureText:     usr.SignatureText,
		SignatureHTML:     usr.SignatureHTML,
		JoinedAt:          NewTime(usr.CreatedAt),
		Groups:            Groups(groups),
		URL: m.WebURL(route.For("User").Fill(map[string]string{
			"user": usr.Username,
		}).AsURL()),
	}
}

func (m *Mapper) Account(
	usr *user.User,
	groups []group.Group,
	profile *setting.UserProfile,
) Account {
	acc := Account{
		User:         m.User(usr, groups),
		Email:        usr.Email,
		EmailIsJID:   usr.EmailIsJID,
		Language:     usr.Language,
		Timezone:     usr.Timezone,
		OTPEnabled:   usr.OTPEnabled,
		SignInLastAt: NewTime(usr.SignInLastAt),
		ConfirmedAt:  NewTime(usr.ConfirmedAt),
	}
	if profile != nil {
		acc.NotifyOnReply = profile.NotifyOnReply
	}

	return acc
}

func KeyFrom(key *apikey.APIKey) Key {
	return Key{
		ID:         key.ID.String(),
		Name:       key.Name,
		LastUsedAt: NewTime(key.LastUsedAt),
		CreatedAt:  NewTime(key.CreatedAt),
	}
}

func Keys(keys []apikey.APIKey) []Key {
	views := make([]Key, 0, len(keys))
	for i := range keys {
		views = append(views, KeyFrom(&keys[i]))
	}

	return views
}

func SessionPermissionsFor(perms *permission.Resolution) SessionPermissions {
	sp := SessionPermissions{
		Admin:      perms.IsAdmin,
		Default:    permission.LevelToString(perms.DefaultLevel()),
		Categories: map[string]string{},
	}
	for id, lvl := range perms.Levels() {
		sp.Categories[id.String()] = permission.LevelToString(lvl)
	}

	return sp
}

func OTPEnrollmentFrom(e *user.OTPEnrollment) *OTPEnrollment {
	if e == nil {
		return nil
	}

	return &OTPEnrollment{URL: e.URL, Secret: e.Secret, QRDataURI: e.QRDataURI}
}

func (m *Mapper) AdminUser(usr *user.User) AdminUser {
	memberOf := usr.MemberOf
	if memberOf == nil {
		memberOf = []string{}
	}

	return AdminUser{
		Account:   m.Account(usr, nil, nil),
		MemberOf:  memberOf,
		BannedAt:  NewTime(usr.BannedAt),
		DeletedAt: NewTime(usr.DeletedAt),
	}
}

func (m *Mapper) AdminUsers(users []user.User) []AdminUser {
	views := make([]AdminUser, 0, len(users))
	for i := range users {
		views = append(views, m.AdminUser(&users[i]))
	}

	return views
}

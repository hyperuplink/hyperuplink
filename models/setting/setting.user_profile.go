package setting

import (
	"encoding/json"

	"github.com/google/uuid"
)

const (
	UserProfilePrefix           string = "profile-"
	UserProfileNotifyOnReplyKey string = "notify_on_reply"

	UserProfileViewBanner          string = "banner"
	UserProfileViewFooter          string = "footer"
	UserProfileViewProfilePictures string = "profile_pictures"
)

type UserProfile struct {
	NotifyOnReply       bool `json:"notify_on_reply"`
	ShowBanner          bool `json:"show_banner"`
	ShowFooter          bool `json:"show_footer"`
	ShowProfilePictures bool `json:"show_profile_pictures"`
}

func (p *UserProfile) UnmarshalJSON(data []byte) error {
	type userProfile UserProfile

	stored := userProfile(NewUserProfile())
	if err := json.Unmarshal(data, &stored); err != nil {
		return err
	}

	*p = UserProfile(stored)

	return nil
}

func UserProfileID(userID uuid.UUID) string {
	return UserProfilePrefix + userID.String()
}

func NewUserProfile() UserProfile {
	return UserProfile{
		NotifyOnReply:       true,
		ShowBanner:          true,
		ShowFooter:          true,
		ShowProfilePictures: true,
	}
}

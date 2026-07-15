package setting

import (
	"github.com/google/uuid"
)

const (
	UserProfilePrefix           string = "profile-"
	UserProfileNotifyOnReplyKey string = "notify_on_reply"
)

type UserProfile struct {
	NotifyOnReply bool `json:"notify_on_reply"`
}

func UserProfileID(userID uuid.UUID) string {
	return UserProfilePrefix + userID.String()
}

func NewUserProfile() UserProfile {
	return UserProfile{
		NotifyOnReply: true,
	}
}

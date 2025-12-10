package vforum

import "github.com/mrusme/hyperuplink/models/forum"

type VForum struct {
	forum.Forum

	Topics  int `json:"topics"`
	Replies int `json:"replies"`
}

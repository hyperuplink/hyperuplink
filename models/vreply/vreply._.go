package vreply

import (
	"github.com/jackc/pgx/v5/pgtype"
	"xn--gckvb8fzb.com/hyperuplink/models/reply"
)

type VReply struct {
	reply.Reply

	AuthorUsername       string           `json:"author_username"`
	AuthorRole           string           `json:"author_role"`
	AuthorMemberOf       []string         `json:"author_member_of"`
	AuthorEmail          string           `json:"author_email"`
	AuthorProfilePicture string           `json:"author_profile_picture"`
	AuthorSignatureText  string           `json:"author_signature_text"`
	AuthorSignatureHTML  string           `json:"author_signature_html"`
	AuthorJoinedAt       pgtype.Timestamp `json:"author_joined_at"`
}

package view

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"xn--gckvb8fzb.com/glides/http/route"
	logictopics "xn--gckvb8fzb.com/hyperuplink/logic/root/categories/forums/topics"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
	"xn--gckvb8fzb.com/hyperuplink/models/vreply"
	"xn--gckvb8fzb.com/hyperuplink/models/vsearchresult"
	"xn--gckvb8fzb.com/hyperuplink/models/vtopic"
)

type Author struct {
	ID                string `json:"id"`
	Username          string `json:"username"`
	Role              string `json:"role"`
	ProfilePictureURL string `json:"profile_picture_url,omitempty"`
	SignatureHTML     string `json:"signature_html,omitempty"`
	JoinedAt          Time   `json:"joined_at"`
}

type Attachment struct {
	ID        string `json:"id"`
	Filename  string `json:"filename"`
	MimeType  string `json:"mime_type"`
	CreatedAt Time   `json:"created_at"`
	URL       string `json:"url"`
}

type TopicViewer struct {
	Unread bool `json:"unread"`
}

type Topic struct {
	ID             string       `json:"id"`
	ShortID        string       `json:"short_id"`
	Slug           string       `json:"slug"`
	Name           string       `json:"name"`
	Kind           string       `json:"kind"`
	Pinned         bool         `json:"pinned"`
	LockedAt       Time         `json:"locked_at"`
	Text           string       `json:"text"`
	HTML           string       `json:"html"`
	CreatedAt      Time         `json:"created_at"`
	LastActivityAt Time         `json:"last_activity_at"`
	Replies        int          `json:"replies"`
	Views          int64        `json:"views"`
	Author         Author       `json:"author"`
	Forum          Ref          `json:"forum"`
	Category       Ref          `json:"category"`
	Attachments    []Attachment `json:"attachments"`
	Poll           *Poll        `json:"poll,omitempty"`
	Viewer         *TopicViewer `json:"viewer,omitempty"`
	URL            string       `json:"url"`
}

type PollOption struct {
	Index   int    `json:"index"`
	Text    string `json:"text"`
	Votes   int    `json:"votes"`
	Percent int    `json:"percent"`
}

type PollViewer struct {
	CanVote   bool `json:"can_vote"`
	HasVoted  bool `json:"has_voted"`
	Selection *int `json:"selection"`
}

type Poll struct {
	Options []PollOption `json:"options"`
	Total   int          `json:"total"`
	Ended   bool         `json:"ended"`
	EndsAt  Time         `json:"ends_at"`
	Viewer  *PollViewer  `json:"viewer,omitempty"`
}

type Reply struct {
	ID          string       `json:"id"`
	ShortID     string       `json:"short_id"`
	TopicID     string       `json:"topic_id"`
	ParentID    *string      `json:"parent_id"`
	Text        string       `json:"text"`
	HTML        string       `json:"html"`
	CreatedAt   Time         `json:"created_at"`
	Author      Author       `json:"author"`
	Attachments []Attachment `json:"attachments"`
	URL         string       `json:"url"`
}

func (m *Mapper) Author(
	id uuid.UUID,
	username string,
	role string,
	picture string,
	signatureHTML string,
	joinedAt pgtype.Timestamp,
) Author {
	return Author{
		ID:                id.String(),
		Username:          username,
		Role:              role,
		ProfilePictureURL: m.profilePictureURL(picture),
		SignatureHTML:     signatureHTML,
		JoinedAt:          NewTime(joinedAt),
	}
}

func (m *Mapper) AuthorFromUser(usr *user.User) Author {
	return m.Author(
		usr.ID,
		usr.Username,
		string(usr.Role),
		usr.ProfilePicture,
		usr.SignatureHTML,
		usr.CreatedAt,
	)
}

func (m *Mapper) Attachments(ids []uuid.UUID) []Attachment {
	atts := make([]Attachment, 0, len(ids))
	for _, id := range ids {
		att, ok := m.attachments[id]
		if !ok || att == nil {
			continue
		}
		atts = append(atts, Attachment{
			ID:        att.ID.String(),
			Filename:  att.Filename,
			MimeType:  att.MimeType,
			CreatedAt: NewTime(att.CreatedAt),
			URL: APIPath("APIV1Attachment", map[string]string{
				"attachment": att.ID.String(),
			}),
		})
	}

	return atts
}

func TopicURL(m *Mapper, categorySlug string, forumSlug string, topicSlug string) string {
	return m.WebURL(route.For("CategoriesForumsTopics").Fill(map[string]string{
		"categories": categorySlug,
		"forums":     forumSlug,
		"topics":     topicSlug,
	}).AsURL())
}

func (m *Mapper) Topic(top *vtopic.VTopic) Topic {
	return Topic{
		ID:             top.ID.String(),
		ShortID:        top.ShortID,
		Slug:           top.Slug,
		Name:           top.Name,
		Kind:           string(top.Kind),
		Pinned:         top.Pinned,
		LockedAt:       NewTime(top.LockedAt),
		Text:           top.Text,
		HTML:           top.HTML,
		CreatedAt:      NewTime(top.CreatedAt),
		LastActivityAt: NewTime(top.LastActivityAt),
		Replies:        top.Replies,
		Views:          top.Views,
		Author: m.Author(
			top.AuthorID,
			top.AuthorUsername,
			top.AuthorRole,
			top.AuthorProfilePicture,
			top.AuthorSignatureHTML,
			top.AuthorJoinedAt,
		),
		Forum:       Ref{ID: top.ForumID.String(), Name: top.ForumName, Slug: top.ForumSlug},
		Category:    Ref{Name: top.CategoryName, Slug: top.CategorySlug},
		Attachments: m.Attachments(top.AttachmentIDs),
		URL:         TopicURL(m, top.CategorySlug, top.ForumSlug, top.Slug),
	}
}

func (m *Mapper) Topics(tops []vtopic.VTopic) []Topic {
	ids := make([][]uuid.UUID, 0, len(tops))
	for i := range tops {
		ids = append(ids, tops[i].AttachmentIDs)
	}
	m.PrefetchAttachments(ids...)

	views := make([]Topic, 0, len(tops))
	for i := range tops {
		views = append(views, m.Topic(&tops[i]))
	}

	return views
}

func AttachUnread(tops []Topic, unread map[string]bool) {
	if unread == nil {
		return
	}

	for i := range tops {
		tops[i].Viewer = &TopicViewer{Unread: unread[tops[i].ID]}
	}
}

func PollFrom(poll *logictopics.Poll, ended pgtype.Timestamp, viewer bool) *Poll {
	if poll == nil {
		return nil
	}

	v := &Poll{
		Options: make([]PollOption, 0, len(poll.Options)),
		Total:   poll.Total,
		Ended:   poll.Ended,
		EndsAt:  NewTime(ended),
	}
	for _, option := range poll.Options {
		v.Options = append(v.Options, PollOption{
			Index:   option.Index,
			Text:    option.Text,
			Votes:   option.Votes,
			Percent: option.Percent,
		})
	}

	if viewer {
		pv := &PollViewer{
			CanVote:  poll.CanVote,
			HasVoted: poll.HasVoted,
		}
		if poll.HasVoted {
			selection := poll.Selection
			pv.Selection = &selection
		}
		v.Viewer = pv
	}

	return v
}

func (m *Mapper) Reply(
	rep *vreply.VReply,
	categorySlug string,
	forumSlug string,
	topicSlug string,
) Reply {
	var parentID *string
	if rep.ReplyID.Valid {
		parent := rep.ReplyID.UUID.String()
		parentID = &parent
	}

	return Reply{
		ID:        rep.ID.String(),
		ShortID:   rep.ShortID,
		TopicID:   rep.TopicID.String(),
		ParentID:  parentID,
		Text:      rep.Text,
		HTML:      rep.HTML,
		CreatedAt: NewTime(rep.CreatedAt),
		Author: m.Author(
			rep.AuthorID,
			rep.AuthorUsername,
			rep.AuthorRole,
			rep.AuthorProfilePicture,
			rep.AuthorSignatureHTML,
			rep.AuthorJoinedAt,
		),
		Attachments: m.Attachments(rep.AttachmentIDs),
		URL:         TopicURL(m, categorySlug, forumSlug, topicSlug) + "#post-" + rep.ShortID,
	}
}

func (m *Mapper) Replies(
	reps []vreply.VReply,
	categorySlug string,
	forumSlug string,
	topicSlug string,
) []Reply {
	ids := make([][]uuid.UUID, 0, len(reps))
	for i := range reps {
		ids = append(ids, reps[i].AttachmentIDs)
	}
	m.PrefetchAttachments(ids...)

	views := make([]Reply, 0, len(reps))
	for i := range reps {
		views = append(views, m.Reply(&reps[i], categorySlug, forumSlug, topicSlug))
	}

	return views
}

type Page struct {
	Name string `json:"name"`
	Text string `json:"text"`
	HTML string `json:"html"`
}

type SearchResult struct {
	Kind         string `json:"kind"`
	Topic        Ref    `json:"topic"`
	Forum        Ref    `json:"forum"`
	Category     Ref    `json:"category"`
	Author       string `json:"author"`
	CreatedAt    Time   `json:"created_at"`
	ReplyShortID string `json:"reply_short_id,omitempty"`
	Page         int    `json:"page"`
	URL          string `json:"url"`
}

func (m *Mapper) SearchResults(results []vsearchresult.VSearchResult, perPage int) []SearchResult {
	views := make([]SearchResult, 0, len(results))
	for _, res := range results {
		url := TopicURL(m, res.CategorySlug, res.ForumSlug, res.TopicSlug)
		if res.IsReply() {
			url += "#post-" + res.ReplyShortID
		}
		views = append(views, SearchResult{
			Kind:         res.Kind,
			Topic:        Ref{Name: res.TopicName, Slug: res.TopicSlug},
			Forum:        Ref{Slug: res.ForumSlug},
			Category:     Ref{Slug: res.CategorySlug},
			Author:       res.AuthorUsername,
			CreatedAt:    NewTime(res.CreatedAt),
			ReplyShortID: res.ReplyShortID,
			Page:         res.PageWithin(perPage),
			URL:          url,
		})
	}

	return views
}

package vreply

type VReplyWithTopic struct {
	VReply

	TopicName    string `json:"topic_name"`
	TopicSlug    string `json:"topic_slug"`
	ForumSlug    string `json:"forum_slug"`
	CategorySlug string `json:"category_slug"`
}

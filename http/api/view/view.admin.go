package view

import (
	"xn--gckvb8fzb.com/hyperuplink/models/postevent"
	"xn--gckvb8fzb.com/hyperuplink/models/vactivity"
	"xn--gckvb8fzb.com/hyperuplink/models/vpostevent"
)

type Log struct {
	ID        string   `json:"id"`
	Kind      string   `json:"kind"`
	Actor     string   `json:"actor"`
	Subject   string   `json:"subject"`
	SubjectID *string  `json:"subject_id"`
	Count     int      `json:"count"`
	Path      string   `json:"path,omitempty"`
	Setting   string   `json:"setting,omitempty"`
	Changed   []string `json:"changed,omitempty"`
	CreatedAt Time     `json:"created_at"`
	UpdatedAt Time     `json:"updated_at"`
}

type Report struct {
	ID           string  `json:"id"`
	Target       string  `json:"target"`
	TopicID      *string `json:"topic_id"`
	ReplyID      *string `json:"reply_id"`
	Reason       string  `json:"reason"`
	Author       string  `json:"author"`
	TargetAuthor string  `json:"target_author"`
	Snippet      string  `json:"snippet"`
	URL          string  `json:"url"`
	CreatedAt    Time    `json:"created_at"`
}

var reasons = map[postevent.ReportType]string{
	postevent.Spam:       "spam",
	postevent.Misconduct: "misconduct",
	postevent.Illegal:    "illegal",
}

func Logs(logs []vactivity.VActivity) []Log {
	views := make([]Log, 0, len(logs))
	for i := range logs {
		entry := &logs[i]
		ctx := entry.Ctx()

		var subjectID *string
		if entry.SubjectID.Valid {
			id := entry.SubjectID.UUID.String()
			subjectID = &id
		}

		views = append(views, Log{
			ID:        entry.ID.String(),
			Kind:      string(entry.Kind),
			Actor:     entry.ActorUsername,
			Subject:   string(entry.Subject),
			SubjectID: subjectID,
			Count:     entry.Count,
			Path:      ctx.Path,
			Setting:   ctx.Setting,
			Changed:   ctx.Changed,
			CreatedAt: NewTime(entry.CreatedAt),
			UpdatedAt: NewTime(entry.UpdatedAt),
		})
	}

	return views
}

func (m *Mapper) Reports(reports []vpostevent.VPostEvent) []Report {
	views := make([]Report, 0, len(reports))
	for i := range reports {
		rep := &reports[i]

		var topicID, replyID *string
		if rep.TopicID.Valid {
			id := rep.TopicID.UUID.String()
			topicID = &id
		}
		if rep.ReplyID.Valid {
			id := rep.ReplyID.UUID.String()
			replyID = &id
		}

		url := ""
		if path := rep.TargetPath(); path != "" {
			url = m.WebURL(path)
		}

		views = append(views, Report{
			ID:           rep.ID.String(),
			Target:       string(rep.Target),
			TopicID:      topicID,
			ReplyID:      replyID,
			Reason:       reasons[postevent.ReportType(rep.Selection)],
			Author:       rep.AuthorUsername,
			TargetAuthor: rep.TargetAuthorUsername,
			Snippet:      rep.TargetSnippet(),
			URL:          url,
			CreatedAt:    NewTime(rep.CreatedAt),
		})
	}

	return views
}

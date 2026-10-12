package view

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"xn--gckvb8fzb.com/glides/http/route"
	"xn--gckvb8fzb.com/glides/runtime"
	"xn--gckvb8fzb.com/glides/services/repositories/common"
	"xn--gckvb8fzb.com/glides/services/storage"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/logic/helpers/paging"
	"xn--gckvb8fzb.com/hyperuplink/models/attachment"
	"xn--gckvb8fzb.com/hyperuplink/models/permission"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
	settingRepo "xn--gckvb8fzb.com/hyperuplink/services/repositories/setting"
)

const Version = "v1"

type Time struct {
	ts pgtype.Timestamp
}

func NewTime(ts pgtype.Timestamp) Time {
	return Time{ts: ts}
}

func (t Time) Valid() bool {
	return t.ts.Valid
}

func (t Time) MarshalJSON() ([]byte, error) {
	if !t.ts.Valid {
		return []byte("null"), nil
	}

	return json.Marshal(t.ts.Time.UTC().Format(time.RFC3339Nano))
}

type Pagination struct {
	Page    int   `json:"page"`
	PerPage int   `json:"per_page"`
	Total   int64 `json:"total"`
	Pages   int   `json:"pages"`
}

func NewPagination(page int, perPage int, total int64) Pagination {
	return Pagination{
		Page:    page,
		PerPage: perPage,
		Total:   total,
		Pages:   paging.Pages(total, perPage),
	}
}

type Ref struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Slug string `json:"slug"`
}

type Permissions struct {
	Read     bool `json:"read"`
	Write    bool `json:"write"`
	Moderate bool `json:"moderate"`
}

func PermissionsFor(perms *permission.Resolution, categoryID uuid.UUID) Permissions {
	return Permissions{
		Read:     perms.CanReadID(categoryID),
		Write:    perms.CanWriteID(categoryID),
		Moderate: perms.CanModerateID(categoryID),
	}
}

type Discovery struct {
	Name    string            `json:"name"`
	API     string            `json:"api"`
	Version string            `json:"version"`
	Guests  bool              `json:"guests"`
	BaseURL string            `json:"base_url"`
	Links   map[string]string `json:"links"`
}

type Mapper struct {
	ctx         context.Context
	rt          *runtime.Runtime
	system      *setting.System
	profiles    *setting.Profiles
	attachments map[uuid.UUID]*attachment.Attachment
}

func NewMapper(ctx context.Context, rt *runtime.Runtime, system *setting.System) *Mapper {
	return &Mapper{
		ctx:         ctx,
		rt:          rt,
		system:      system,
		attachments: map[uuid.UUID]*attachment.Attachment{},
	}
}

func APIPath(id string, params map[string]string) string {
	return "/" + route.For(id).Fill(params).AsURL()
}

func (m *Mapper) WebURL(path string) string {
	return strings.TrimRight(m.system.BaseURL, "/") + "/" + path
}

func (m *Mapper) PrefetchAttachments(idLists ...[]uuid.UUID) {
	var missing []uuid.UUID
	for _, ids := range idLists {
		for _, id := range ids {
			if _, ok := m.attachments[id]; !ok {
				m.attachments[id] = nil
				missing = append(missing, id)
			}
		}
	}

	if len(missing) == 0 {
		return
	}

	atts, err := gh.Repositories(m.rt).Attachment.AllByUUIDs(m.ctx, missing, common.QueryOptions{})
	if err != nil {
		m.rt.Error("error", err)
		return
	}

	for i := range atts {
		att := &atts[i]
		m.attachments[att.ID] = att
	}
}

func (m *Mapper) profilePictureURL(id string) string {
	if id == "" {
		return ""
	}

	if m.profiles == nil {
		settingProfiles, err := settingRepo.GetByID[setting.Profiles](m.ctx,
			gh.Repositories(m.rt).Setting,
			"profiles",
		)
		if err != nil {
			return ""
		}
		m.profiles = &settingProfiles.JSONValue
	}

	if !m.profiles.EnablePicture || m.profiles.PictureStorageProviderID == "" {
		return ""
	}

	dlurl, abs, err := storage.From(m.rt).GetFileDownloadURL(
		m.profiles.PictureStorageProviderID,
		m.profiles.PictureStoragePath+"/"+id+"."+m.profiles.PictureFormat,
	)
	if err != nil {
		return ""
	}

	if abs {
		return dlurl
	}

	return m.WebURL(strings.TrimPrefix(dlurl, "/"))
}

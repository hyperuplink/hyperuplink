package view

import (
	"xn--gckvb8fzb.com/glides/http/route"
	"xn--gckvb8fzb.com/hyperuplink/models/category"
	"xn--gckvb8fzb.com/hyperuplink/models/permission"
	"xn--gckvb8fzb.com/hyperuplink/models/vforum"
)

type Category struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Slug        string      `json:"slug"`
	Position    int         `json:"position"`
	Permissions Permissions `json:"permissions"`
	Forums      []Forum     `json:"forums,omitempty"`
	URL         string      `json:"url"`
}

type Forum struct {
	ID             string      `json:"id"`
	Name           string      `json:"name"`
	Slug           string      `json:"slug"`
	Position       int         `json:"position"`
	Description    string      `json:"description"`
	Category       Ref         `json:"category"`
	Topics         int         `json:"topics"`
	Replies        int         `json:"replies"`
	LastActivityAt Time        `json:"last_activity_at"`
	Permissions    Permissions `json:"permissions"`
	URL            string      `json:"url"`
}

func (m *Mapper) Category(
	cat *category.Category,
	forums []Forum,
	perms *permission.Resolution,
) Category {
	return Category{
		ID:          cat.ID.String(),
		Name:        cat.Name,
		Slug:        cat.Slug,
		Position:    cat.Position,
		Permissions: PermissionsFor(perms, cat.ID),
		Forums:      forums,
		URL: m.WebURL(route.For("Categories").Fill(map[string]string{
			"categories": cat.Slug,
		}).AsURL()),
	}
}

func (m *Mapper) Forum(fum *vforum.VForum, perms *permission.Resolution) Forum {
	return Forum{
		ID:             fum.ID.String(),
		Name:           fum.Name,
		Slug:           fum.Slug,
		Position:       fum.Position,
		Description:    fum.Description,
		Category:       Ref{ID: fum.CategoryID.String(), Name: fum.CategoryName, Slug: fum.CategorySlug},
		Topics:         fum.Topics,
		Replies:        fum.Replies,
		LastActivityAt: NewTime(fum.LastReplyAt),
		Permissions:    PermissionsFor(perms, fum.CategoryID),
		URL: m.WebURL(route.For("CategoriesForums").Fill(map[string]string{
			"categories": fum.CategorySlug,
			"forums":     fum.Slug,
		}).AsURL()),
	}
}

func (m *Mapper) Forums(fums []vforum.VForum, perms *permission.Resolution) []Forum {
	views := make([]Forum, 0, len(fums))
	for i := range fums {
		views = append(views, m.Forum(&fums[i], perms))
	}

	return views
}

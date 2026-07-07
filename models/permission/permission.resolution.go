package permission

import "github.com/google/uuid"

type Resolution struct {
	IsAdmin     bool
	Default     byte
	levelByID   map[uuid.UUID]byte
	levelBySlug map[string]byte
}

func NewResolution() *Resolution {
	return &Resolution{
		levelByID:   map[uuid.UUID]byte{},
		levelBySlug: map[string]byte{},
	}
}

func (r *Resolution) SetCategory(id uuid.UUID, slug string, level byte) {
	if r.levelByID == nil {
		r.levelByID = map[uuid.UUID]byte{}
	}
	if r.levelBySlug == nil {
		r.levelBySlug = map[string]byte{}
	}
	r.levelByID[id] = level
	r.levelBySlug[slug] = level
}

func (r *Resolution) LevelByID(id uuid.UUID) byte {
	if r.IsAdmin {
		return Moderate
	}
	if lvl, ok := r.levelByID[id]; ok {
		return lvl
	}
	return r.Default
}

func (r *Resolution) LevelBySlug(slug string) byte {
	if r.IsAdmin {
		return Moderate
	}
	if lvl, ok := r.levelBySlug[slug]; ok {
		return lvl
	}
	return r.Default
}

func (r *Resolution) CanReadID(id uuid.UUID) bool {
	return r.LevelByID(id) >= ReadOnly
}

func (r *Resolution) CanWriteID(id uuid.UUID) bool {
	return r.LevelByID(id) >= ReadWrite
}

func (r *Resolution) CanModerateID(id uuid.UUID) bool {
	return r.LevelByID(id) >= Moderate
}

func (r *Resolution) CanReadSlug(slug string) bool {
	return r.LevelBySlug(slug) >= ReadOnly
}

func (r *Resolution) CanWriteSlug(slug string) bool {
	return r.LevelBySlug(slug) >= ReadWrite
}

func (r *Resolution) CanModerateSlug(slug string) bool {
	return r.LevelBySlug(slug) >= Moderate
}

func (r *Resolution) ReadAll() bool {
	return r.IsAdmin || r.Default >= ReadOnly
}

func (r *Resolution) CanWriteAny() bool {
	if r.IsAdmin {
		return true
	}
	if r.Default >= ReadWrite {
		return true
	}
	for _, lvl := range r.levelByID {
		if lvl >= ReadWrite {
			return true
		}
	}
	return false
}

func (r *Resolution) AllowedReadSlugs() []string {
	if r.ReadAll() {
		return nil
	}
	slugs := []string{}
	for slug, lvl := range r.levelBySlug {
		if lvl >= ReadOnly {
			slugs = append(slugs, slug)
		}
	}
	return slugs
}

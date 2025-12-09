package flash

import (
	"github.com/gofiber/fiber/v3"
)

type FlashType string

const (
	DebugFlash FlashType = "_debug"
	InfoFlash  FlashType = "_info"
	WarnFlash  FlashType = "_warn"
	ErrorFlash FlashType = "_error"
)

type Flash struct {
	c       fiber.Ctx
	flashes map[string]string
}

func New(c fiber.Ctx) (f *Flash) {
	f = new(Flash)
	f.c = c
	f.flashes = make(map[string]string)

	messages := c.Redirect().Messages()
	for _, msg := range messages {
		f.flashes[msg.Key] = msg.Value
	}

	return f
}

func (f *Flash) HasErrors() bool {
	errs := f.Errors()
	if len(errs) > 0 {
		return true
	}

	return false
}

func (f *Flash) SetError(err error) {
	f.flashes[string(ErrorFlash)] = err.Error()
}

func (f *Flash) SetInfo(info string) {
	f.flashes[string(InfoFlash)] = info
}

func (f *Flash) All() (flashes map[string]string) {
	return f.flashes
}

func (f *Flash) Clear() {
	f.flashes = make(map[string]string)
}

func (f *Flash) errors() (errs []string) {
	for key, err := range f.flashes {
		if key != string(DebugFlash) &&
			key != string(InfoFlash) &&
			key != string(WarnFlash) {
			errs = append(errs, err)
		}
	}
	return errs
}

func (f *Flash) Errors() (errs []string) {
	errs = f.errors()
	f.Clear()
	return errs
}

func (f *Flash) Get(flashType FlashType) (s string) {
	var ok bool = false

	if s, ok = f.flashes[string(flashType)]; !ok {
		return ""
	} else {
		return s
	}
}

func (f *Flash) Debug() (s string) {
	s = f.Get(DebugFlash)
	f.Clear()
	return s
}

func (f *Flash) Info() (s string) {
	s = f.Get(InfoFlash)
	f.Clear()
	return s
}

func (f *Flash) Warn() (s string) {
	s = f.Get(WarnFlash)
	f.Clear()
	return s
}

func (f *Flash) SetErrorsMap(errsmap map[string]error) {
	for field, err := range errsmap {
		f.flashes[field] = err.Error()
	}
}

func (f *Flash) ClassFor(field string) string {
	if _, exists := f.flashes[field]; exists {
		return "error"
	}

	return ""
}

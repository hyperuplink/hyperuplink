package flash

import "github.com/gofiber/fiber/v3"

type FlashType string

const (
	DebugFlash FlashType = "debug"
	InfoFlash  FlashType = "info"
	WarnFlash  FlashType = "warn"
	ErrorFlash FlashType = "error"
)

type Flash struct {
	c       fiber.Ctx
	errsmap map[string]error
	flashes map[FlashType]string
}

func New(c fiber.Ctx) (f *Flash) {
	f = new(Flash)
	f.c = c
	f.errsmap = make(map[string]error)
	f.flashes = make(map[FlashType]string)

	for _, msg := range c.Redirect().Messages() {
		f.flashes[FlashType(msg.Key)] = msg.Value
	}

	return f
}

func (f *Flash) HasErrors() bool {
	_, ok := f.flashes[ErrorFlash]
	if ok || len(f.errsmap) > 0 {
		return true
	}

	return false
}

func (f *Flash) SetError(err error) {
	f.flashes[ErrorFlash] = err.Error()
}

func (f *Flash) SetInfo(info string) {
	f.flashes[InfoFlash] = info
}

func (f *Flash) All() (flashes map[string]string) {
	flashes = make(map[string]string)
	for key, msg := range f.flashes {
		flashes[string(key)] = msg
	}
	for key, err := range f.errsmap {
		flashes[key] = err.Error()
	}

	return flashes
}

func (f *Flash) Errors() (errs []string) {
	if err, ok := f.flashes[ErrorFlash]; ok {
		errs = append(errs, err)
	}
	for _, err := range f.errsmap {
		errs = append(errs, err.Error())
	}
	return errs
}

func (f *Flash) SetErrorsMap(errsmap map[string]error) {
	f.errsmap = errsmap
}

func (f *Flash) GetErrorsMap() (errsmap map[string]error) {
	return f.errsmap
}

func (f *Flash) ClassFor(field string) string {
	if _, exists := f.errsmap[field]; exists {
		return "error"
	}

	return ""
}

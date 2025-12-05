package in

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/kaptinlin/go-i18n"
	"github.com/mrusme/hyperuplink/http/route"
)

type Internationalization struct {
	r          route.IRouteController
	c          fiber.Ctx
	acceptLang string
	I18n       *i18n.Localizer
}

func New(r route.IRouteController, c fiber.Ctx) *Internationalization {
	i := new(Internationalization)
	i.r = r
	i.c = c

	// TODO: If the user is logged in, query the lang from profile settings
	al := c.Get("Accept-Language", "en")
	als := strings.Split(al, "-")
	al = strings.ToLower(als[0])

	i.SetLang(al)

	return i
}

func (i *Internationalization) T(msg string) string {
	return i.I18n.Get(msg)
}

func (i *Internationalization) SetLang(lang string) {
	i.acceptLang = lang
	i.I18n = i.r.GetRuntime().Intnat.NewLocalizer(i.acceptLang)
}

func (i *Internationalization) Lang() string {
	return i.acceptLang
}

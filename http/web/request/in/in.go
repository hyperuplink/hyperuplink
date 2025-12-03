package in

import (
	"github.com/gofiber/fiber/v3"
	"github.com/kaptinlin/go-i18n"
	"github.com/mrusme/hyperuplink/http/route"
)

type Internationalization struct {
	r    route.IRoute
	c    fiber.Ctx
	I18n *i18n.Localizer
}

func New(r route.IRoute, c fiber.Ctx) *Internationalization {
	i := new(Internationalization)
	i.r = r
	i.c = c

	acceptLang := c.Get("Accept-Language", "en")
	i.I18n = r.GetRuntime().Intnat.NewLocalizer(acceptLang)

	return i
}

func (i *Internationalization) T(msg string) string {
	return i.I18n.Get(msg)
}

package in

import (
	"html/template"

	"github.com/gofiber/fiber/v3"
	"github.com/kaptinlin/go-i18n"
	"xn--gckvb8fzb.com/glides/http/route"
	"xn--gckvb8fzb.com/glides/services/intnat"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
)

type Internationalization struct {
	r    route.IRouteController
	c    fiber.Ctx
	lang string
	I18n *i18n.Localizer
}

func New(r route.IRouteController, c fiber.Ctx) *Internationalization {
	i := new(Internationalization)
	i.r = r
	i.c = c

	i.SetLang(c.Get(fiber.HeaderAcceptLanguage))

	return i
}

func (i *Internationalization) Ts(msg string) string {
	return i.I18n.Get(msg)
}

func (i *Internationalization) T(msg string) template.HTML {
	return template.HTML(i.Ts(msg))
}

func (i *Internationalization) SetLang(lang string) {
	i.lang = gh.Locale(i.r.GetRuntime(), lang)
	i.I18n = intnat.From(i.r.GetRuntime()).NewLocalizer(i.lang)
}

func (i *Internationalization) Lang() string {
	return i.lang
}

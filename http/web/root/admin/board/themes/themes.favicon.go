package themes

import (
	"github.com/gofiber/fiber/v3"
	logicthemes "xn--gckvb8fzb.com/hyperuplink/logic/root/admin/board/themes"
)

func (r *Route) FaviconUpload(c fiber.Ctx) (err error) {
	return r.imageUpload(c, logicthemes.Favicon, "custom_favicon")
}

func (r *Route) FaviconRemove(c fiber.Ctx) (err error) {
	return r.imageRemove(c, logicthemes.Favicon)
}

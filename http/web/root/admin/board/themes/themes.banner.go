package themes

import (
	"github.com/gofiber/fiber/v3"
	logicthemes "xn--gckvb8fzb.com/hyperuplink/logic/root/admin/board/themes"
)

func (r *Route) BannerUpload(c fiber.Ctx) (err error) {
	return r.imageUpload(c, logicthemes.Banner, "custom_banner")
}

func (r *Route) BannerRemove(c fiber.Ctx) (err error) {
	return r.imageRemove(c, logicthemes.Banner)
}

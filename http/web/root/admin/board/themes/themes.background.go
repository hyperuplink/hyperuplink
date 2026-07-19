package themes

import (
	"github.com/gofiber/fiber/v3"
	logicthemes "xn--gckvb8fzb.com/hyperuplink/logic/root/admin/board/themes"
)

func (r *Route) BackgroundUpload(c fiber.Ctx) (err error) {
	return r.imageUpload(c, logicthemes.Background, "custom_background")
}

func (r *Route) BackgroundRemove(c fiber.Ctx) (err error) {
	return r.imageRemove(c, logicthemes.Background)
}

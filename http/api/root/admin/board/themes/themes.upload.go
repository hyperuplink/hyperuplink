package themes

import (
	"github.com/gofiber/fiber/v3"
	logicthemes "xn--gckvb8fzb.com/hyperuplink/logic/root/admin/board/themes"
)

// @Summary		Upload the board banner
// @Tags			admin
// @Accept			mpfd
// @Param			image	formData	file	true	"The image"
// @Success		204
// @Failure		401	{object}	problem.Problem
// @Failure		403	{object}	problem.Problem
// @Failure		422	{object}	problem.Problem
// @Security		BearerAuth
// @Router			/admin/board/themes/banner [post]
func (r *Route) BannerUpload(c fiber.Ctx) (err error) {
	return r.imageUpload(c, logicthemes.Banner)
}

// @Summary		Remove the board banner
// @Tags			admin
// @Success		204
// @Failure		401	{object}	problem.Problem
// @Failure		403	{object}	problem.Problem
// @Security		BearerAuth
// @Router			/admin/board/themes/banner [delete]
func (r *Route) BannerRemove(c fiber.Ctx) (err error) {
	return r.imageRemove(c, logicthemes.Banner)
}

// @Summary		Upload the board favicon
// @Tags			admin
// @Accept			mpfd
// @Param			image	formData	file	true	"The image"
// @Success		204
// @Failure		401	{object}	problem.Problem
// @Failure		403	{object}	problem.Problem
// @Failure		422	{object}	problem.Problem
// @Security		BearerAuth
// @Router			/admin/board/themes/favicon [post]
func (r *Route) FaviconUpload(c fiber.Ctx) (err error) {
	return r.imageUpload(c, logicthemes.Favicon)
}

// @Summary		Remove the board favicon
// @Tags			admin
// @Success		204
// @Failure		401	{object}	problem.Problem
// @Failure		403	{object}	problem.Problem
// @Security		BearerAuth
// @Router			/admin/board/themes/favicon [delete]
func (r *Route) FaviconRemove(c fiber.Ctx) (err error) {
	return r.imageRemove(c, logicthemes.Favicon)
}

// @Summary		Upload the board background
// @Tags			admin
// @Accept			mpfd
// @Param			image	formData	file	true	"The image"
// @Success		204
// @Failure		401	{object}	problem.Problem
// @Failure		403	{object}	problem.Problem
// @Failure		422	{object}	problem.Problem
// @Security		BearerAuth
// @Router			/admin/board/themes/background [post]
func (r *Route) BackgroundUpload(c fiber.Ctx) (err error) {
	return r.imageUpload(c, logicthemes.Background)
}

// @Summary		Remove the board background
// @Tags			admin
// @Success		204
// @Failure		401	{object}	problem.Problem
// @Failure		403	{object}	problem.Problem
// @Security		BearerAuth
// @Router			/admin/board/themes/background [delete]
func (r *Route) BackgroundRemove(c fiber.Ctx) (err error) {
	return r.imageRemove(c, logicthemes.Background)
}

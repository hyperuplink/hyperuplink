package themes

import (
	"github.com/gofiber/fiber/v3"
	logicthemes "xn--gckvb8fzb.com/hyperuplink/logic/root/admin/board/themes"
)

// @Summary	Upload the board banner
// @Tags		admin
// @Accept		mpfd
// @Produce	json
// @Param		custom_banner	formData	file	true	"The image to use as the board banner"
// @Success	200				{object}	request.StatusResponse
// @Failure	401				{object}	request.ErrorResponse
// @Failure	403				{object}	request.ErrorResponse
// @Failure	422				{object}	request.ErrorResponse
// @Security	BearerAuth
// @Security	APIKeyAuth
// @Router		/admin/board/themes/banner [post]
func (r *Route) BannerUpload(c fiber.Ctx) (err error) {
	return r.imageUpload(c, logicthemes.Banner, "custom_banner")
}

// @Summary	Remove the board banner
// @Tags		admin
// @Produce	json
// @Success	200	{object}	request.StatusResponse
// @Failure	401	{object}	request.ErrorResponse
// @Failure	403	{object}	request.ErrorResponse
// @Security	BearerAuth
// @Security	APIKeyAuth
// @Router		/admin/board/themes/banner [delete]
func (r *Route) BannerRemove(c fiber.Ctx) (err error) {
	return r.imageRemove(c, logicthemes.Banner)
}

// @Summary	Upload the board favicon
// @Tags		admin
// @Accept		mpfd
// @Produce	json
// @Param		custom_favicon	formData	file	true	"The image to use as the board favicon"
// @Success	200				{object}	request.StatusResponse
// @Failure	401				{object}	request.ErrorResponse
// @Failure	403				{object}	request.ErrorResponse
// @Failure	422				{object}	request.ErrorResponse
// @Security	BearerAuth
// @Security	APIKeyAuth
// @Router		/admin/board/themes/favicon [post]
func (r *Route) FaviconUpload(c fiber.Ctx) (err error) {
	return r.imageUpload(c, logicthemes.Favicon, "custom_favicon")
}

// @Summary	Remove the board favicon
// @Tags		admin
// @Produce	json
// @Success	200	{object}	request.StatusResponse
// @Failure	401	{object}	request.ErrorResponse
// @Failure	403	{object}	request.ErrorResponse
// @Security	BearerAuth
// @Security	APIKeyAuth
// @Router		/admin/board/themes/favicon [delete]
func (r *Route) FaviconRemove(c fiber.Ctx) (err error) {
	return r.imageRemove(c, logicthemes.Favicon)
}

// @Summary	Upload the board background
// @Tags		admin
// @Accept		mpfd
// @Produce	json
// @Param		custom_background	formData	file	true	"The image to use as the board background"
// @Success	200					{object}	request.StatusResponse
// @Failure	401					{object}	request.ErrorResponse
// @Failure	403					{object}	request.ErrorResponse
// @Failure	422					{object}	request.ErrorResponse
// @Security	BearerAuth
// @Security	APIKeyAuth
// @Router		/admin/board/themes/background [post]
func (r *Route) BackgroundUpload(c fiber.Ctx) (err error) {
	return r.imageUpload(c, logicthemes.Background, "custom_background")
}

// @Summary	Remove the board background
// @Tags		admin
// @Produce	json
// @Success	200	{object}	request.StatusResponse
// @Failure	401	{object}	request.ErrorResponse
// @Failure	403	{object}	request.ErrorResponse
// @Security	BearerAuth
// @Security	APIKeyAuth
// @Router		/admin/board/themes/background [delete]
func (r *Route) BackgroundRemove(c fiber.Ctx) (err error) {
	return r.imageRemove(c, logicthemes.Background)
}

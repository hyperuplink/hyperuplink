package root

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	"xn--gckvb8fzb.com/hyperuplink/http/api/view"
)

// @Summary		Describe the API
// @Description	The board's name, the API version, the build, whether requests
// @Description	without a key are served, and the paths of the main collections.
// @Description	This is the one document a client may always read.
// @Tags			board
// @Produce		json
// @Success		200	{object}	view.Discovery
// @Router			/ [get]
func (r *Route) Index(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	version, _, _, _ := r.Runtime.GetBuild()

	return req.Respond(view.Discovery{
		Name:    req.System.Name,
		API:     view.Version,
		Version: version,
		Guests:  req.System.AllowAPIGuests,
		BaseURL: req.System.BaseURL,
		Links: map[string]string{
			"openapi": view.APIPath("APIV1OpenAPI", nil),
			"docs":    view.APIPath("APIV1Docs", nil) + "/",
			"session": view.APIPath("APIV1Session", nil),
			"board":   view.APIPath("APIV1Board", nil),
			"forums":  view.APIPath("APIV1Forums", nil),
			"topics":  view.APIPath("APIV1Topics", nil),
			"search":  view.APIPath("APIV1Search", nil),
			"users":   view.APIPath("APIV1Users", nil),
			"account": view.APIPath("APIV1Account", nil),
		},
	})
}
